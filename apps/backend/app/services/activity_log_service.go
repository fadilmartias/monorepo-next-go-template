package services

import (
	"context"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"gorm.io/gorm"
)

type ActivityLogService struct {
	db                    *gorm.DB
	activityLogRepository *repositories.ActivityLogRepository
}

func NewActivityLogService(db *gorm.DB, activityLogRepository *repositories.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{
		db:                    db,
		activityLogRepository: activityLogRepository,
	}
}

func (s *ActivityLogService) Log(ctx context.Context, entry *models.ActivityLog) error {
	traceID := utils.GetTraceID(ctx)
	ip := utils.GetIP(ctx)
	userAgent := utils.GetUserAgent(ctx)
	entry.TraceID = &traceID
	entry.IP = &ip
	entry.UserAgent = &userAgent
	return s.activityLogRepository.Create(entry)
}

func (s *ActivityLogService) FindByID(id string) (*models.ActivityLog, error) {
	return s.activityLogRepository.FindByIDHasPreload(id)
}
