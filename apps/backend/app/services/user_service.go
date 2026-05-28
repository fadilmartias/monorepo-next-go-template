package services

import (
	"errors"

	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/responses"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	UserRepository *repositories.UserRepository
	Redis          *redis.Client
}

func NewUserService(db *gorm.DB, redis *redis.Client, userRepository *repositories.UserRepository) *UserService {
	return &UserService{DB: db, UserRepository: userRepository, Redis: redis}
}

func (s *UserService) GetAll() ([]responses.UserResponse, error) {
	users, _, err := s.UserRepository.GetAll(1, 10)
	if err != nil {
		return nil, err
	}
	var userResponses []responses.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, responses.UserResponse{
			ID:              &u.ID,
			Name:            &u.Name,
			Phone:           &u.Phone,
			Email:           &u.Email,
			Role:            &u.Role,
			EmailVerifiedAt: u.EmailVerifiedAt,
			TenantID:        &u.TenantID,
			CreatedAt:       &u.CreatedAt,
			UpdatedAt:       &u.UpdatedAt,
		})
	}
	return userResponses, nil
}

// Show mengambil satu user
func (s *UserService) FindByID(id string) (*responses.UserResponse, error) {
	userDB, err := s.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &responses.UserResponse{
		ID:              &userDB.ID,
		Name:            &userDB.Name,
		Phone:           &userDB.Phone,
		Email:           &userDB.Email,
		Role:            &userDB.Role,
		EmailVerifiedAt: userDB.EmailVerifiedAt,
		TenantID:        &userDB.TenantID,
		CreatedAt:       &userDB.CreatedAt,
		UpdatedAt:       &userDB.UpdatedAt,
		Is2FAEnabled:    userDB.TOTPSecret != nil,
	}, nil
}

func (s *UserService) UpdateProfile(id string, input requests.UpdateProfileInput) (*responses.UserResponse, error) {
	user, err := s.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Map untuk fields yang berubah
	updates := map[string]any{}

	// Cek nama
	if user.Name != input.Name {
		updates["name"] = input.Name
	}

	// Cek phone
	if user.Phone != input.Phone && input.Phone != "" {
		userDB, err := s.UserRepository.FindByPhone(input.Phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if userDB != nil {
			return nil, utils.NewFormError("Nomor HP sudah digunakan", map[string]string{
				"phone": "Nomor HP sudah digunakan",
			})
		}
		updates["phone"] = input.Phone
	}

	// Cek email
	if user.Email != input.Email && input.Email != "" {
		userDB, err := s.UserRepository.FindByEmail(input.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if userDB != nil {
			return nil, utils.NewFormError("Email sudah digunakan", map[string]string{
				"email": "Email sudah digunakan",
			})
		}
		updates["email"] = input.Email
	}

	// Tidak ada perubahan
	if len(updates) == 0 {
		return nil, errors.New("tidak ada perubahan")
	}

	// Update DB
	_, err = s.UserRepository.UpdateWithModel(user, updates)
	if err != nil {
		return nil, err
	}

	return &responses.UserResponse{
		ID:              &user.ID,
		Name:            &user.Name,
		Phone:           &user.Phone,
		Email:           &user.Email,
		Role:            &user.Role,
		EmailVerifiedAt: user.EmailVerifiedAt,
		TenantID:        &user.TenantID,
		CreatedAt:       &user.CreatedAt,
		UpdatedAt:       &user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdatePassword(id string, input requests.UpdatePasswordInput) error {
	userDB, err := s.UserRepository.FindByID(id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(input.CurrentPassword)); err != nil {
		return utils.NewFormError("Password tidak sesuai", map[string]string{
			"current_password": "Password tidak sesuai",
		})
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	if err != nil {
		return err
	}
	input.NewPassword = string(bytes)
	if _, err := s.UserRepository.UpdateByID(id, map[string]any{"password": input.NewPassword}); err != nil {
		return err
	}
	return nil
}
