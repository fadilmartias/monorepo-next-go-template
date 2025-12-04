package repositories

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[models.User]
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: NewBaseRepository[models.User](db), DB: db}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: r.BaseRepository.WithTx(tx), DB: tx}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateProviderID(user *models.User, provider string, id string) error {
	updates := map[string]any{}
	switch provider {
	case "google":
		updates["google_id"] = id
	case "discord":
		updates["discord_id"] = id
	case "facebook":
		updates["facebook_id"] = id
	case "steam":
		updates["steam_id"] = id
	case "twitch":
		updates["twitch_id"] = id
	case "apple":
		updates["apple_id"] = id
	}
	return r.DB.Model(user).Updates(updates).Error
}
