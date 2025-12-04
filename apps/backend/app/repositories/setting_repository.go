package repositories

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

type SettingRepository struct {
	BaseRepository[models.Setting]
	DB *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{BaseRepository: NewBaseRepository[models.Setting](db), DB: db}
}

func (r *SettingRepository) WithTx(tx *gorm.DB) *SettingRepository {
	return &SettingRepository{BaseRepository: r.BaseRepository.WithTx(tx), DB: tx}
}

func (r *SettingRepository) GetByKey(key string) (string, error) {
	var value string
	err := r.DB.Model(&models.Setting{}).
		Where("`key` = ?", key).
		Select("value").
		First(&value).Error
	return value, err
}

func (r *SettingRepository) UpdateValueByKey(key string, value string) error {
	return r.DB.Model(&struct{}{}).
		Table("settings").
		Where("`key` = ?", key).
		Update("value", value).Error
}
