package controllers_v0

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/processors"
	"github.com/fadilmartias/dilz_code/apps/backend/app/registry"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils" // Ganti dengan path utils Anda
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type GenericController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewGenericController(db *gorm.DB, redis *redis.Client) *GenericController {
	return &GenericController{DB: db, Redis: redis}
}

// Index menangani GET /:model
// @Summary Ambil daftar data model
// @Description Mengambil daftar data berdasarkan nama model secara dinamis dengan dukungan paginasi dan filter query params.
// @Tags Generic API
// @Accept json
// @Produce json
// @Param model path string true "Nama Model (contoh: users, posts)"
// @Param cache query bool false "Gunakan cache Redis (true/false)"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengambil data"
// @Failure 404 {object} utils.OrderedErrorResponse "Model tidak ditemukan di registry"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v0/{model} [get]
func (ctrl *GenericController) Index(c fiber.Ctx) error {
	modelName := c.Params("model")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	queryParams := c.Request().URI().QueryArgs()
	urlValues := make(url.Values)
	queryParams.VisitAll(func(key, value []byte) {
		urlValues.Add(string(key), string(value))
	})
	params := utils.NewQueryParams(urlValues)
	tx := ctrl.DB.Model(modelInfo.Instance)
	tx = utils.BuildGormQuery(tx, urlValues, false)

	isCache := false

	// Cek query (optional override juga kalau mau)
	if c.Query("cache") == "true" {
		isCache = true
	}

	var cacheKey string
	if isCache {
		cacheKey = fmt.Sprintf("api:%s:%s", modelName, c.Request().URI().QueryString())
	}

	apiResponse, err := utils.FetchAndCacheDynamic(
		c.Context(), ctrl.Redis, tx, params, cacheKey, 1*time.Minute, false,
		modelInfo.Instance, modelInfo.NewSlice,
	)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	switch resp := apiResponse.(type) {
	case utils.PaginatedResponse[any]:
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code: http.StatusOK, Message: "Data retrieved successfully", Data: resp.Data, Pagination: &resp.Pagination,
		})
	default:
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: "Invalid response type"})
	}
}

// Show menangani GET /:model/:id
// @Summary Ambil detail data model
// @Description Mengambil satu baris data berdasarkan ID atau slug secara dinamis.
// @Tags Generic API
// @Accept json
// @Produce json
// @Param model path string true "Nama Model"
// @Param id path string true "ID atau Slug data"
// @Param slug query bool false "Set true jika mencari berdasarkan slug"
// @Param cache query bool false "Gunakan cache Redis (true/false)"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengambil detail data"
// @Failure 404 {object} utils.OrderedErrorResponse "Model tidak ditemukan"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v0/{model}/{id} [get]
func (ctrl *GenericController) Show(c fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	urlValues := make(url.Values)
	queryArgs := c.Request().URI().QueryArgs()
	queryArgs.VisitAll(func(key, value []byte) {
		urlValues.Add(string(key), string(value))
	})
	params := utils.NewQueryParams(urlValues)

	// Tambahkan filter by ID atau slug
	tx := ctrl.DB.Model(modelInfo.Instance)
	if c.Query("slug") == "true" {
		tx = tx.Where("slug = ?", id)
	} else {
		tx = tx.Where("id = ?", id)
	}
	tx = utils.BuildGormQuery(tx, urlValues, true)

	isCache := false

	// Cek query (optional override juga kalau mau)
	if c.Query("cache") == "true" {
		isCache = true
	}

	var cacheKey string
	if isCache {
		cacheKey = fmt.Sprintf("api:%s:%s:%s", modelName, id, c.Request().URI().QueryString())
	}

	apiResponse, err := utils.FetchAndCacheDynamic(
		c.Context(), ctrl.Redis, tx, params, cacheKey, 5*time.Minute, true,
		modelInfo.Instance, modelInfo.NewSlice,
	)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	switch resp := apiResponse.(type) {
	case utils.SingleResponse[any]:
		baseURL := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
		metaData, _ := processors.GenericPostProcessor(resp.Data, baseURL)

		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    http.StatusOK,
			Message: "Data retrieved successfully",
			Data:    resp.Data,
			Meta:    metaData, // Kirim meta yang sudah diproses
		})
	default:
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: "Invalid response type"})
	}
}

// Store menangani POST /:model
// @Summary Buat data baru
// @Description Menyimpan entri data baru ke database untuk model yang ditentukan.
// @Tags Generic API
// @Accept json
// @Produce json
// @Param model path string true "Nama Model"
// @Param payload body map[string]interface{} true "Data JSON yang akan disimpan"
// @Success 201 {object} utils.OrderedSuccessResponse "Berhasil membuat data"
// @Failure 400 {object} utils.OrderedErrorResponse "Invalid request body"
// @Failure 404 {object} utils.OrderedErrorResponse "Model tidak ditemukan"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v0/{model} [post]
func (ctrl *GenericController) Store(c fiber.Ctx) error {
	// 1. Dapatkan info model dari registry
	modelName := c.Params("model")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	// 2. Buat instance baru dari model menggunakan reflection
	newInstance := reflect.New(reflect.TypeOf(modelInfo.Instance).Elem()).Interface()

	// 3. Parse request body JSON ke dalam instance baru tersebut
	if err := c.Bind().Body(newInstance); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusBadRequest, Message: "Invalid request body"})
	}

	// 4. Set ID jika belum ada
	val := reflect.ValueOf(newInstance).Elem()
	idField := val.FieldByName("Id")
	if idField.IsValid() && idField.Kind() == reflect.String && idField.String() == "" {
		idField.SetString(utils.GenerateShortID(7))
	}

	// 5. Simpan instance ke database
	if result := ctrl.DB.Create(newInstance); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	// 6. Kembalikan data yang baru dibuat
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusCreated,
		Message: "Data created successfully",
		Data:    newInstance,
	})
}

// Update menangani PUT /:model/:id
// @Summary Update keseluruhan data
// @Description Melakukan update data yang ada berdasarkan ID. Menggunakan PUT biasanya merepresentasikan penggantian keseluruhan baris.
// @Tags Generic API
// @Accept json
// @Produce json
// @Param model path string true "Nama Model"
// @Param id path string true "ID data"
// @Param payload body map[string]interface{} true "Data JSON untuk update"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengupdate data"
// @Failure 400 {object} utils.OrderedErrorResponse "Invalid request body"
// @Failure 404 {object} utils.OrderedErrorResponse "Data tidak ditemukan"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v0/{model}/{id} [put]
func (ctrl *GenericController) Update(c fiber.Ctx) error {
	// 1. Dapatkan info model dan parameter
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	// 2. Buat instance untuk menampung data update dari body
	updateData := reflect.New(reflect.TypeOf(modelInfo.Instance).Elem()).Interface()
	if err := c.Bind().Body(updateData); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusBadRequest, Message: "Invalid request body"})
	}

	// 3. Lakukan update di database
	result := ctrl.DB.Model(modelInfo.Instance).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	// 4. Periksa apakah ada baris yang terpengaruh
	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: "Data not found or no changes made"})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusOK,
		Message: "Data updated successfully",
		Data:    fiber.Map{"rows_affected": result.RowsAffected},
	})
}

// Patch menangani PATCH /:model/:id
// @Summary Update parsial data
// @Description Melakukan update data secara parsial berdasarkan field yang dikirimkan.
// @Tags Generic API
// @Accept json
// @Produce json
// @Param model path string true "Nama Model"
// @Param id path string true "ID data"
// @Param payload body map[string]interface{} true "Data JSON parsial untuk update"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengupdate data parsial"
// @Failure 400 {object} utils.OrderedErrorResponse "Invalid request body"
// @Failure 404 {object} utils.OrderedErrorResponse "Data tidak ditemukan"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v0/{model}/{id} [patch]
func (ctrl *GenericController) Patch(c fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	updateMap := map[string]any{}
	if err := c.Bind().Body(&updateMap); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	result := ctrl.DB.Model(modelInfo.Instance).Where("id = ?", id).Updates(updateMap)
	if result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: "Data not found or no changes made"})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusOK,
		Message: "Data updated successfully",
		Data:    fiber.Map{"rows_affected": result.RowsAffected},
	})
}

// Destroy menangani DELETE /:model/:id
// @Summary Hapus data
// @Description Menghapus (atau soft delete) data berdasarkan ID.
// @Tags Generic API
// @Accept json
// @Produce json
// @Param model path string true "Nama Model"
// @Param id path string true "ID data"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil menghapus data"
// @Failure 404 {object} utils.OrderedErrorResponse "Data atau model tidak ditemukan"
// @Failure 500 {object} utils.OrderedErrorResponse "Internal Server Error"
// @Router /v0/{model}/{id} [delete]
func (ctrl *GenericController) Destroy(c fiber.Ctx) error {
	// 1. Dapatkan info model dan parameter
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	// 2. Hapus data dari database
	result := ctrl.DB.Where("id = ?", id).Delete(modelInfo.Instance)
	if result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	// 3. Periksa apakah ada baris yang terpengaruh
	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: "Data not found"})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusOK,
		Message: "Data deleted successfully",
	})
}
