package repositories

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	BaseRepository[models.ActivityLog]
	DB *gorm.DB
}

func (r *ActivityLogRepository) WithTx(tx *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{BaseRepository: r.BaseRepository.WithTx(tx), DB: tx}
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{BaseRepository: NewBaseRepository[models.ActivityLog](db), DB: db}
}

func (r *ActivityLogRepository) FindByIDHasPreload(id string) (*models.ActivityLog, error) {
	var log models.ActivityLog

	if err := r.DB.First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Load subject
	switch log.SubjectType {
	case "user":
		var user models.User
		r.DB.First(&user, log.SubjectID)
		log.Subject = user
	}
	// Load causer
	switch log.CauserType {
	case "user":
		var user models.User
		r.DB.First(&user, log.CauserID)
		log.Causer = user
	}
	return &log, nil
}
