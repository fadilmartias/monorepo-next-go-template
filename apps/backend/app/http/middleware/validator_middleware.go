package middleware

import (
	"reflect"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("email_if_target_type_email", emailIfTargetTypeEmail)
}

func emailIfTargetTypeEmail(fl validator.FieldLevel) bool {
	parent := fl.Parent()
	targetType := parent.FieldByName("TargetType")
	target := fl.Field().String()

	if targetType.IsValid() && targetType.String() == "email" {
		// validate as email
		return validator.New().Var(target, "email") == nil
	}
	return true
}

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	if f, ok := structType.FieldByName(fieldName); ok {
		tag := f.Tag.Get("json")
		if tag != "" && tag != "-" {
			return splitJSONTag(tag)
		}
	}
	return fieldName // fallback to struct field name
}

func splitJSONTag(tag string) string {
	// only take the first part, ignore omitempty
	if commaIdx := indexComma(tag); commaIdx != -1 {
		return tag[:commaIdx]
	}
	return tag
}

func indexComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

func Validator[T any]() fiber.Handler {
	return func(c fiber.Ctx) error {
		var body T

		if err := c.Bind().Body(&body); err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:      fiber.StatusBadRequest,
				Message:   "Data yang diberikan tidak valid",
				ErrorCode: "VALIDATION_ERROR",
			}, err)
		}

		if err := validate.Struct(body); err != nil {
			errs := make(map[string]string)
			valErrs := err.(validator.ValidationErrors)

			t := reflect.TypeOf(body)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			for _, e := range valErrs {
				jsonKey := getJSONFieldName(t, e.StructField())
				errs[jsonKey] = e.Tag()
			}

			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:      fiber.StatusUnprocessableEntity,
				Message:   "Data yang diberikan tidak valid",
				ErrorCode: "VALIDATION_ERROR",
				Errors:    errs,
			}, err)
		}

		c.Locals("validatedBody", body)
		return c.Next()
	}
}
