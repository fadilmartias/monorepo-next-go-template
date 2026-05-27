// middleware/auth.go
package middleware

import (
	"errors"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

func Auth(allowedRoles []string, allowedPermissions []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userClaims := c.Locals("user")
		if userClaims == nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: User not found",
			}, errors.New("user not found"))
		}

		claims, ok := userClaims.(jwt.MapClaims)
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid user data",
			}, errors.New("invalid user data"))
		}

		// Validate role
		if len(allowedRoles) > 0 {
			role, ok := claims["role"].(string)
			if !ok || !utils.SliceContains(allowedRoles, role) {
				return utils.ErrorResponse(c, utils.ErrorResponseFormat{
					Code:    fiber.StatusForbidden,
					Message: "Forbidden: Access denied",
				}, errors.New("access denied"))
			}
		}

		// Tambahan: validate permission kalau dibutuhkan
		// if len(allowedPermissions) > 0 {
		//    ...
		// }

		return c.Next()
	}
}
