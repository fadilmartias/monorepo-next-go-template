package services

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/responses"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ArticleService struct {
	DB                *gorm.DB
	ArticleRepository *repositories.ArticleRepository
	Redis             *redis.Client
}

func NewArticleService(db *gorm.DB, redis *redis.Client, articleRepository *repositories.ArticleRepository) *ArticleService {
	return &ArticleService{DB: db, ArticleRepository: articleRepository, Redis: redis}
}

func (s *ArticleService) GetPublishedArticles(limit int) ([]models.Article, error) {
	return s.ArticleRepository.GetPublishedArticles(limit)
}

func (s *ArticleService) GetPublishedArticleById(id string) (*models.Article, error) {
	return s.ArticleRepository.GetPublishedArticleById(id)
}

func (s *ArticleService) GetPublishedArticleBySlug(slug string, related bool) (*responses.ArticleResponse, error) {
	article, err := s.ArticleRepository.GetPublishedArticleBySlug(slug)
	if err != nil {
		return nil, err
	}
	var relatedArticles []models.Article
	if related {
		relatedArticles, err = s.ArticleRepository.GetRelatedArticles(article)
		if err != nil {
			return nil, err
		}
	} else {
		relatedArticles = nil
	}
	return &responses.ArticleResponse{
		ID:              article.ID,
		TenantID:        article.TenantID,
		CategoryID:      article.CategoryID,
		Category:        article.Category,
		Title:           article.Title,
		Author:          article.Author,
		Slug:            article.Slug,
		Content:         article.Content,
		Img:             article.Img,
		Views:           article.Views,
		Likes:           article.Likes,
		MetaDescription: article.MetaDescription,
		MetaKeywords:    article.MetaKeywords,
		MetaImg:         article.MetaImg,
		PublishedAt:     article.PublishedAt,
		IsActive:        article.IsActive,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
		RelatedArticles: relatedArticles,
	}, nil
}

func (s *ArticleService) GetPopularArticles(limit int, except string) ([]models.Article, error) {
	return s.ArticleRepository.GetPopularArticles(limit, except)
}

func (s *ArticleService) Process(article *models.Article) (*models.Article, error) {
	// finalDir := "public/uploads/images/articles"
	finalDir := "uploads/images/articles"

	// Pindahkan thumbnail utama
	if err := utils.UploadImageToR2IfExists(article.Img, finalDir, "Gagal memindahkan image"); err != nil {
		return nil, err
	}

	// Pindahkan meta image kalau ada
	if article.MetaImg != nil {
		if err := utils.UploadImageToR2IfExists(*article.MetaImg, finalDir, "Gagal memindahkan meta image"); err != nil {
			return nil, err
		}
	}

	// Update
	if article.ID != "" {
		existing, err := s.ArticleRepository.FindByID(article.ID)
		if err != nil {
			return nil, err
		}
		if _, err := s.ArticleRepository.UpdateWithExistingModel(existing, *article); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create
	if err := s.ArticleRepository.Create(article); err != nil {
		return nil, err
	}
	return article, nil
}
