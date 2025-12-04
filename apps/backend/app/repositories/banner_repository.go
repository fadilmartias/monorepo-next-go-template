package repositories

import (
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

type BannerRepository struct {
	BaseRepository[models.Banner]
	DB *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{BaseRepository: NewBaseRepository[models.Banner](db), DB: db}
}

func (r *BannerRepository) WithTx(tx *gorm.DB) *BannerRepository {
	return &BannerRepository{BaseRepository: r.BaseRepository.WithTx(tx), DB: tx}
}

func (r *BannerRepository) GetActiveAndValidBanners(limit int) ([]models.Banner, error) {
	var banners []models.Banner
	query := r.DB.Where("is_active = ?", true).
		Where("valid_from <= ?", time.Now()).
		Where("valid_until >= ?", time.Now()).
		Order("`order` ASC").
		Order("created_at DESC")

	// kalau limit > 0 baru dipasang
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&banners).Error
	return banners, err
}

func (r *BannerRepository) FindActiveAndValidBannerById(id string) (*models.Banner, error) {
	var banner models.Banner
	err := r.DB.Where("id = ?", id).Where("is_active = ?", true).Where("valid_from <= ?", time.Now()).Where("valid_until >= ?", time.Now()).First(&banner).Error
	return &banner, err
}
