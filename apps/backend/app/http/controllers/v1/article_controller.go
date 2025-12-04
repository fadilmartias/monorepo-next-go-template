package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

type ArticleController struct {
	ArticleService *services.ArticleService
}

func NewArticleController(articleService *services.ArticleService) *ArticleController {
	return &ArticleController{ArticleService: articleService}
}

func (ctrl *ArticleController) Index(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	articles, err := ctrl.ArticleService.GetPublishedArticles(limit)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal mendapatkan data article",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data article",
		Data:    articles,
	})
}

func (ctrl *ArticleController) Show(c *fiber.Ctx) error {
	slug := c.Params("slug")
	related := c.QueryBool("related", false)
	article, err := ctrl.ArticleService.GetPublishedArticleBySlug(slug, related)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Artikel tidak ditemukan",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data article",
		Data:    article,
	})
}

func (ctrl *ArticleController) Popular(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	except := c.Query("except", "")
	articles, err := ctrl.ArticleService.GetPopularArticles(limit, except)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Artikel tidak ditemukan",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data article",
		Data:    articles,
	})
}

func (ctrl *ArticleController) Process(c *fiber.Ctx) error {
	article := models.Article{}
	if err := c.BodyParser(&article); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Gagal memproses data artikel",
		}, err)
	}

	article.TenantID = "1"
	result, err := ctrl.ArticleService.Process(&article)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Gagal memproses artikel",
		}, err)
	}

	msg := "Berhasil menambahkan artikel"
	if article.ID != "" {
		msg = "Berhasil mengupdate artikel"
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: msg,
		Data:    result,
	})
}
