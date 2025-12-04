package repositories

import (
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	BaseRepository[models.Article]
	DB *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{BaseRepository: NewBaseRepository[models.Article](db), DB: db}
}

func (r *ArticleRepository) WithTx(tx *gorm.DB) *ArticleRepository {
	return &ArticleRepository{BaseRepository: r.BaseRepository.WithTx(tx), DB: tx}
}

func (r *ArticleRepository) GetPublishedArticles(limit int) ([]models.Article, error) {
	var articles []models.Article
	query := r.DB.Preload("Category").Where("is_active = ?", true).Where("published_at <= ?", time.Now()).Order("published_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *ArticleRepository) GetPublishedArticleById(id string) (*models.Article, error) {
	var article models.Article
	if err := r.DB.Preload("Category").Where("id = ?", id).Where("is_active = ?", true).Where("published_at <= ?", time.Now()).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) GetPublishedArticleBySlug(slug string) (*models.Article, error) {
	var article models.Article
	if err := r.DB.Preload("Category").Where("slug = ?", slug).Where("is_active = ?", true).Where("published_at <= ?", time.Now()).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) GetRelatedArticles(article *models.Article) ([]models.Article, error) {
	var relatedArticles []models.Article
	if err := r.DB.Where("category_id = ?", article.CategoryID).Where("id != ?", article.ID).Where("is_active = ?", true).Where("published_at <= ?", time.Now()).Order("published_at DESC").Limit(3).Find(&relatedArticles).Error; err != nil {
		return nil, err
	}
	return relatedArticles, nil
}

func (r *ArticleRepository) GetPopularArticles(limit int, except string) ([]models.Article, error) {
	var articles []models.Article
	query := r.DB.Preload("Category").Where("is_active = ?", true).Where("published_at <= ?", time.Now()).Order("published_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if except != "" {
		query = query.Where("slug != ?", except)
	}
	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}
