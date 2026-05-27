package controllers_v0

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type UploadController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewUploadController(db *gorm.DB, redis *redis.Client) *UploadController {
	return &UploadController{DB: db, Redis: redis}
}

// Upload menangani unggahan file
// @Summary Upload file (Gambar, Dokumen, atau TinyMCE)
// @Description Endpoint untuk mengunggah file. Mendukung penerimaan file melalui field form-data: 'image', 'file', atau 'tinymce'. Hanya mendukung ekstensi jpg, png, webp, gif, dan svg.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file false "File gambar yang akan diupload"
// @Param file formData file false "File dokumen/umum yang akan diupload"
// @Param tinymce formData file false "File gambar khusus dari editor TinyMCE"
// @Success 200 {object} utils.OrderedSuccessResponse "Berhasil mengunggah file"
// @Failure 400 {object} utils.OrderedErrorResponse "Gagal memproses form, field tidak ditemukan, atau tipe file tidak didukung"
// @Failure 500 {object} utils.OrderedErrorResponse "Kesalahan konfigurasi server atau gagal menyimpan file"
// @Router /v0/uploads [post]
func (ctrl *UploadController) Upload(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Failed to parse form",
		})
	}

	files := form.File
	acceptedFields := []string{"image", "file", "tinymce"}

	var uploadedFile *multipart.FileHeader
	var fieldName string

	for _, field := range acceptedFields {
		fileHeaders, ok := files[field]
		if ok && len(fileHeaders) > 0 {
			uploadedFile = fileHeaders[0]
			fieldName = field
			break
		}
	}

	if uploadedFile == nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "No valid file field found (image, file)",
		})
	}

	contentType := uploadedFile.Header.Get("Content-Type")
	allowedTypes := map[string]string{
		"image/jpeg":    ".jpg",
		"image/png":     ".png",
		"image/webp":    ".webp",
		"image/gif":     ".gif",
		"image/svg+xml": ".svg",
	}

	ext := filepath.Ext(uploadedFile.Filename)
	if ext == "" || ext == "blob" {
		var ok bool
		ext, ok = allowedTypes[contentType]
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: "Unsupported file type"})
		}
	}

	// Tentukan target direktori
	var uploadDir string
	if fieldName == "tinymce" {
		uploadDir = filepath.Join("storage", "uploads", "tinymce")
	} else {
		// e.g. "storage/uploads/tmp/images"
		uploadDir = filepath.Join("storage", "uploads", "tmp", fmt.Sprintf("%ss", strings.ToLower(fieldName)))
	}

	// GARANSI PEMBUATAN DIREKTORI
	// 0755 memastikan direktori bisa dibaca/dieksekusi publik, dan ditulis oleh owner (user golang)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		// Log error di server console untuk kemudahan tracing ops
		fmt.Printf("[UPLOAD ERROR] Gagal membuat direktori %s: %v\n", uploadDir, err)
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Server configuration error: failed to initialize storage directory",
		})
	}

	timestamp := time.Now().Unix()
	randomStr := utils.GenerateShortID(6)
	newFileName := fmt.Sprintf("%s-%d-%s%s", fieldName, timestamp, randomStr, ext)
	savePath := filepath.Join(uploadDir, newFileName)

	// Save file ke local disk
	if err := c.SaveFile(uploadedFile, savePath); err != nil {
		fmt.Printf("[UPLOAD ERROR] Gagal menyimpan file ke %s: %v\n", savePath, err)
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to write file to disk",
		})
	}

	var fileURL string
	if fieldName == "tinymce" {
		fileURL = fmt.Sprintf("%s/uploads/tinymce/%s", os.Getenv("APP_URL"), newFileName)
	} else {
		fileURL = fmt.Sprintf("%s/uploads/tmp/%ss/%s", os.Getenv("APP_URL"), strings.ToLower(fieldName), newFileName)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Upload successfully",
		Data: fiber.Map{
			"file_url":  fileURL,
			"file_name": newFileName,
		},
	})
}
