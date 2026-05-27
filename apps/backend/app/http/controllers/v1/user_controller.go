package controllers_v1

import (
	"errors"

	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/golang-jwt/jwt"

	"github.com/gofiber/fiber/v3"
	// TAMBAHKAN IMPORT INI
)

type UserController struct {
	UserService *services.UserService
}

// Ubah fungsi NewUserController untuk menerima koneksi DB
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// Index mengambil semua user
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

// Show mengambil satu user
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
