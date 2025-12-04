package responses

import (
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
)

type ArticleResponse struct {
	ID              string           `json:"id,omitempty"`
	TenantID        string           `json:"tenant_id,omitempty"`
	CategoryID      string           `json:"category_id,omitempty"`
	Category        models.Category  `json:"category,omitempty"`
	Title           string           `json:"title,omitempty"`
	Author          string           `json:"author,omitempty"`
	Slug            string           `json:"slug,omitempty"`
	Content         string           `json:"content,omitempty"`
	Img             string           `json:"img,omitempty"`
	Views           int              `json:"views,omitempty"`
	Likes           int              `json:"likes,omitempty"`
	MetaDescription *string          `json:"meta_description,omitempty"`
	MetaKeywords    *string          `json:"meta_keywords,omitempty"`
	MetaImg         *string          `json:"meta_img,omitempty"`
	PublishedAt     time.Time        `json:"published_at,omitempty"`
	IsActive        bool             `json:"is_active,omitempty"`
	CreatedAt       time.Time        `json:"created_at,omitempty"`
	UpdatedAt       time.Time        `json:"updated_at,omitempty"`
	RelatedArticles []models.Article `json:"related_articles,omitempty"`
}
