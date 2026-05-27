package controllers_v1

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"

	"github.com/gofiber/fiber/v3"
)

type ArticleController struct {
	ArticleService *services.ArticleService
}

func NewArticleController(articleService *services.ArticleService) *ArticleController {
	return &ArticleController{ArticleService: articleService}
}

// Index mengambil daftar artikel yang dipublikasikan
// @Summary Ambil daftar artikel publik
// @Description Mengambil daftar artikel yang sudah dipublikasikan (published) untuk blog CodeXtory. Mendukung limitasi jumlah data menggunakan query parameter.
// @Tags Articles
// @Accept json
// @Produce json
// @Param limit query int false "Batasan jumlah artikel yang dikembalikan" default(10)
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data article"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mendapatkan data article"
// @Router /v1/articles [get]
func (ctrl *ArticleController) Index(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
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

// Show mengambil detail artikel berdasarkan slug
// @Summary Ambil detail artikel
// @Description Mengambil satu data artikel yang dipublikasikan berdasarkan slug-nya. Terdapat opsi untuk memuat artikel terkait (related articles).
// @Tags Articles
// @Accept json
// @Produce json
// @Param slug path string true "Slug dari artikel"
// @Param related query bool false "Set true untuk menyertakan artikel terkait" default(false)
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data article"
// @Failure 404 {object} utils.OrderedErrorResponse "Artikel tidak ditemukan"
// @Router /v1/articles/{slug} [get]
func (ctrl *ArticleController) Show(c fiber.Ctx) error {
	slug := c.Params("slug")
	related := fiber.Query[bool](c, "related", false)
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

// Popular mengambil daftar artikel terpopuler
// @Summary Ambil artikel populer
// @Description Mengambil daftar artikel yang paling banyak dibaca atau berstatus populer. Bisa mengecualikan artikel tertentu menggunakan parameter 'except'.
// @Tags Articles
// @Accept json
// @Produce json
// @Param limit query int false "Batasan jumlah artikel yang dikembalikan" default(10)
// @Param except query string false "ID atau slug artikel yang ingin dikecualikan dari hasil"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mendapatkan data article"
// @Failure 404 {object} utils.OrderedErrorResponse "Artikel tidak ditemukan"
// @Router /v1/articles/popular [get]
func (ctrl *ArticleController) Popular(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
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

// Process melakukan operasi Insert atau Update artikel
// @Summary Tambah atau Update Artikel
// @Description Memproses data payload artikel. Jika field ID pada payload kosong, maka akan melakukan Insert (Tambah Artikel). Jika field ID terisi, akan melakukan Update pada artikel tersebut.
// @Tags Articles
// @Accept json
// @Produce json
// @Param payload body models.Article true "Data JSON Artikel"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil menambahkan atau mengupdate artikel"
// @Failure 400 {object} utils.OrderedErrorResponse "Gagal memproses data artikel dari body request"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal memproses/menyimpan artikel ke database"
// @Router /v1/articles/process [post]
func (ctrl *ArticleController) Process(c fiber.Ctx) error {
	article := models.Article{}
	if err := c.Bind().Body(&article); err != nil {
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
