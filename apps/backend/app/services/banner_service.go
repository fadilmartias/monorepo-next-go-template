package services

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"gorm.io/gorm"
)

type BannerService struct {
	DB               *gorm.DB
	BannerRepository *repositories.BannerRepository
	Redis            *config.RedisClient
}

func NewBannerService(db *gorm.DB, redis *config.RedisClient, bannerRepository *repositories.BannerRepository) *BannerService {
	return &BannerService{DB: db, BannerRepository: bannerRepository, Redis: redis}
}

func (s *BannerService) GetActiveAndValidBanners(limit int) ([]models.Banner, error) {
	return s.BannerRepository.GetActiveAndValidBanners(limit)
}

func (s *BannerService) FindActiveAndValidBannerById(id string) (*models.Banner, error) {
	return s.BannerRepository.FindActiveAndValidBannerById(id)
}

func (s *BannerService) Process(banner *models.Banner) (*models.Banner, error) {
	finalDir := "uploads/images/banners"

	// Pindahkan thumbnail utama
	if err := utils.UploadImageToR2IfExists(*banner.Img, finalDir, "Gagal memindahkan image"); err != nil {
		return nil, err
	}

	// Update
	if banner.ID != "" {
		existing, err := s.BannerRepository.FindByID(banner.ID)
		if err != nil {
			return nil, err
		}
		if _, err := s.BannerRepository.UpdateWithExistingModel(existing, *banner); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create
	if err := s.BannerRepository.Create(banner); err != nil {
		return nil, err
	}
	return banner, nil
}
