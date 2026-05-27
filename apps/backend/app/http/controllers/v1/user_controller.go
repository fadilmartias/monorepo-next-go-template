package controllers_v1

import (
	"errors"

	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/golang-jwt/jwt"

	"github.com/gofiber/fiber/v3"
	// TAMBAHKAN IMPORT INI JIKA ADA YANG KURANG
)

type UserController struct {
	UserService *services.UserService
}

// NewUserController menginisialisasi UserController baru
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// Index mengambil semua user
// @Summary Ambil daftar user
// @Description Mengambil semua data pengguna dari database. (Membutuhkan hak akses admin).
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.OrderedSuccessResponse "User berhasil diambil"
// @Failure 500 {object} utils.OrderedErrorResponse "Gagal mengambil user"
// @Router /v1/users [get]
func (ctrl *UserController) Index(c fiber.Ctx) error {
	users, err := ctrl.UserService.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:      fiber.StatusInternalServerError,
			Message:   "Gagal mengambil user",
			ErrorCode: "SERVER_ERROR",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User berhasil diambil",
		Data:    users,
	})
}

// Show mengambil detail satu user
// @Summary Ambil detail user
// @Description Mengambil satu data pengguna secara spesifik menggunakan parameter ID. (Membutuhkan hak akses admin).
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID User"
// @Success 200 {object} utils.OrderedSuccessResponse "User berhasil diambil"
// @Failure 404 {object} utils.OrderedErrorResponse "User tidak ditemukan"
// @Router /v1/users/{id} [get]
func (ctrl *UserController) Show(c fiber.Ctx) error {
	id := c.Params("id")
	userDB, err := ctrl.UserService.FindByID(id)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:      fiber.StatusNotFound,
			Message:   "User tidak ditemukan",
			ErrorCode: "USER_NOT_FOUND",
		}, err)
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User berhasil diambil",
		Data:    userDB,
	})
}

// UpdateProfile memperbarui profil user yang sedang login
// @Summary Update profil user
// @Description Memperbarui data profil pengguna berdasarkan token JWT yang sedang aktif.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.UpdateProfileInput true "Data profil yang ingin diubah"
// @Success 200 {object} utils.OrderedSuccessResponse "Profile berhasil diperbarui"
// @Failure 401 {object} utils.OrderedErrorResponse "Error validasi form input"
// @Failure 500 {object} utils.OrderedErrorResponse "Terjadi kesalahan server"
// @Router /v1/users/profile [patch]
func (ctrl *UserController) UpdateProfile(c fiber.Ctx) error {
	input := c.Value("validatedBody").(requests.UpdateProfileInput)
	id := c.Value("user").(jwt.MapClaims)["id"].(string)
	user, err := ctrl.UserService.UpdateProfile(id, input)

	if err != nil {
		var formErr *utils.FormError
		if errors.As(err, &formErr) {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:      fiber.StatusUnauthorized,
				Message:   formErr.Message,
				ErrorCode: "VALIDATION_ERROR",
				Errors:    formErr.Errors,
			}, err)
		}

		// fallback: error sistem biasa
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:      fiber.StatusInternalServerError,
			Message:   "Terjadi kesalahan",
			ErrorCode: "SERVER_ERROR",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Profile berhasil diperbarui",
		Data:    user,
	})
}

// UpdatePassword memperbarui password user yang sedang login
// @Summary Update password user
// @Description Memperbarui kata sandi pengguna berdasarkan token JWT yang sedang aktif.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body requests.UpdatePasswordInput true "Data kata sandi lama dan baru"
// @Success 200 {object} utils.OrderedSuccessResponse "Password berhasil diperbarui"
// @Failure 401 {object} utils.OrderedErrorResponse "Error validasi form atau password lama salah"
// @Failure 500 {object} utils.OrderedErrorResponse "Terjadi kesalahan server"
// @Router /v1/users/password [patch]
func (ctrl *UserController) UpdatePassword(c fiber.Ctx) error {
	input := c.Value("validatedBody").(requests.UpdatePasswordInput)
	id := c.Value("user").(jwt.MapClaims)["id"].(string)

	err := ctrl.UserService.UpdatePassword(id, input)

	if err != nil {
		var formErr *utils.FormError
		if errors.As(err, &formErr) {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:      fiber.StatusUnauthorized,
				Message:   formErr.Message,
				ErrorCode: "VALIDATION_ERROR",
				Errors:    formErr.Errors,
			}, err)
		}

		// fallback: error sistem biasa
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:      fiber.StatusInternalServerError,
			Message:   "Terjadi kesalahan",
			ErrorCode: "SERVER_ERROR",
		}, err)
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password berhasil diperbarui",
	})
}
