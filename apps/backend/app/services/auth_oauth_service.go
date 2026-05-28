package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"gorm.io/gorm"
)

type AuthOAuthService struct {
	DB                 *gorm.DB
	UserRepository     *repositories.UserRepository
	ActivityLogService *ActivityLogService
	TokenIssuer        *AuthTokenIssuer
}

func NewAuthOAuthService(db *gorm.DB, userRepository *repositories.UserRepository, activityLogService *ActivityLogService, tokenIssuer *AuthTokenIssuer) *AuthOAuthService {
	return &AuthOAuthService{
		DB:                 db,
		UserRepository:     userRepository,
		ActivityLogService: activityLogService,
		TokenIssuer:        tokenIssuer,
	}
}

func (s *AuthOAuthService) CreateNewUser(trx *gorm.DB, input requests.RegisterInput, tenantID string) (*models.User, error) {
	newUser := models.User{
		Name:     input.Name,
		TenantID: tenantID,
		Email:    input.Email,
		Phone:    input.Phone,
		Role:     "user",
	}

	if err := newUser.HashPassword(input.Password); err != nil {
		return nil, err
	}

	if err := s.UserRepository.WithTx(trx).Create(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *AuthOAuthService) HandleOAuth(ctx context.Context, provider, id, email, name string) (*models.User, string, string, string, bool, error) {
	var isRegister bool
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isRegister = true
			trx := s.DB.Begin()
			if trx.Error != nil {
				return nil, "", "", "", isRegister, trx.Error
			}
			defer func() {
				if r := recover(); r != nil {
					trx.Rollback()
				}
			}()

			input := requests.RegisterInput{
				Name:     name,
				Email:    email,
				Phone:    "",
				Password: strconv.Itoa(utils.GenerateRandomNumber(16)),
			}

			newUser, err := s.CreateNewUser(trx, input, "T1")
			if err != nil {
				trx.Rollback()
				return nil, "", "", "", isRegister, err
			}

			s.applyOAuthProviderID(newUser, provider, id)
			emailVerifiedAt := time.Now()
			newUser.EmailVerifiedAt = &emailVerifiedAt

			if err := trx.Save(&newUser).Error; err != nil {
				trx.Rollback()
				return nil, "", "", "", isRegister, err
			}
			if err := trx.Commit().Error; err != nil {
				return nil, "", "", "", isRegister, err
			}
			user = newUser
		} else {
			return nil, "", "", "", isRegister, err
		}
	} else {
		s.applyOAuthProviderID(user, provider, id)
		if _, err := s.UserRepository.UpdateByModel(user); err != nil {
			return nil, "", "", "", isRegister, err
		}
	}

	if user.TOTPSecret != nil {
		tempToken, err := s.TokenIssuer.IssueTempToken(user.Email)
		if err != nil {
			return nil, "", "", "", isRegister, err
		}
		return user, "", "", tempToken, isRegister, Err2FARequired
	}

	accessToken, refreshToken, err := s.TokenIssuer.IssueTokenPair(user, user.TOTPSecret != nil)
	if err != nil {
		return nil, "", "", "", isRegister, err
	}

	user.RefreshToken = &refreshToken
	if _, err := s.UserRepository.UpdateByModel(user); err != nil {
		return nil, "", "", "", isRegister, err
	}

	entry := models.ActivityLog{
		LogName:     "auth",
		Desc:        "Login with " + provider,
		Event:       "login",
		SubjectType: "user",
		SubjectID:   &user.ID,
		CauserType:  "user",
		CauserID:    &user.ID,
		Properties:  nil,
	}
	_ = s.ActivityLogService.Log(ctx, &entry)

	return user, accessToken, refreshToken, "", isRegister, nil
}

func (s *AuthOAuthService) applyOAuthProviderID(user *models.User, provider, id string) {
	switch provider {
	case "google":
		user.GoogleID = models.NewNullString(id)
	case "facebook":
		user.FacebookID = models.NewNullString(id)
	}
}
