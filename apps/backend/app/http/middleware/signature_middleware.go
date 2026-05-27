package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
)

var excludedSignatureRoutes = []string{
	"/v1/midtrans/webhook",
	"/v1/digiflazz/webhook",
	"/v1/product-orders/digiflazz-webhook",
	"/v1/product-orders/ipaymu-webhook",
	"/v1/midtrans/qr-gopay",
	"/v1/auth/google/redirect",
	"/v1/auth/google/callback",
	"/v1/auth/google/one-tap",
	"/v1/auth/discord/redirect",
	"/v1/auth/discord/callback",
	"/v1/auth/facebook/redirect",
	"/v1/auth/facebook/callback",
	"/v1/auth/steam/redirect",
	"/v1/auth/steam/callback",
	"/v1/auth/twitch/redirect",
	"/v1/auth/twitch/callback",
}

func Signature() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Bypass excluded route
		path := strings.Split(c.OriginalURL(), "?")[0]
		for _, route := range excludedSignatureRoutes {
			if strings.Contains(path, route) {
				return c.Next()
			}
		}

		// 2. Bypass for Postman test
		if c.Get("X-Dev-Key") == os.Getenv("DEV_KEY") {
			return c.Next()
		}

		// 3. Required headers
		signature := c.Get("X-Signature")
		timestamp := c.Get("X-Timestamp")

		if signature == "" || timestamp == "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusBadRequest,
				Message: "Missing required headers",
			})
		}

		// 4. Timestamp validation (±5 minutes)
		reqTime, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid timestamp format",
			}, err)
		}
		if math.Abs(time.Since(reqTime).Minutes()) > 5 {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusBadRequest,
				Message: "Request timestamp expired",
			})
		}

		// 5. Compute body hash (SHA256 lowercase hex)
		var bodyHash string

		contentType := c.Get("Content-Type")

		if c.Method() == "GET" {
			rawQuery := string(c.Request().URI().QueryString())
			if rawQuery != "" {
				// 1. pisah per "&"
				pairs := strings.Split(rawQuery, "&")

				flat := map[string]string{}
				keys := []string{}

				for _, pair := range pairs {
					if pair == "" {
						continue
					}

					kv := strings.SplitN(pair, "=", 2)

					keyEnc := kv[0]
					valEnc := ""
					if len(kv) > 1 {
						valEnc = kv[1]
					}

					// decode key+value sama seperti FE
					key, _ := url.QueryUnescape(keyEnc)
					val, _ := url.QueryUnescape(valEnc)

					flat[key] = val
					keys = append(keys, key)
				}

				// sort
				sort.Strings(keys)

				parts := []string{}
				for _, key := range keys {
					parts = append(parts, fmt.Sprintf("%s=%s", key, flat[key]))
				}

				raw := strings.Join(parts, "&")

				logger.Infof("Raw query: %s", raw)

				if raw != "" {
					hash := sha256.Sum256([]byte(raw))
					bodyHash = strings.ToLower(hex.EncodeToString(hash[:]))
				}
			}
		} else if c.Method() != "GET" && len(c.Body()) > 0 {
			// Untuk POST/PUT/PATCH: hash body JSON
			if strings.HasPrefix(contentType, "multipart/form-data") {
				form, err := c.MultipartForm()
				if err != nil {
					return c.Status(400).JSON(fiber.Map{"error": "invalid multipart form"})
				}

				flat := map[string]string{}
				keys := []string{}

				// Ambil field selain file
				for key, vals := range form.Value {
					if len(vals) > 0 {
						flat[key] = vals[0]
						keys = append(keys, key)
					}
				}

				// NOTE: file TIDAK di-hash

				sort.Strings(keys)

				parts := []string{}
				for _, k := range keys {
					parts = append(parts, fmt.Sprintf("%s=%s", k, flat[k]))
				}

				raw := strings.Join(parts, "&")

				if raw != "" {
					sum := sha256.Sum256([]byte(raw))
					bodyHash = strings.ToLower(hex.EncodeToString(sum[:]))
				}

			} else {
				raw := strings.TrimSpace(string(c.Body()))
				if raw != "" && raw != "{}" { // jaga biar {} dianggap kosong
					hash := sha256.Sum256([]byte(raw))
					bodyHash = strings.ToLower(hex.EncodeToString(hash[:]))
				}
			}
		}

		// 6. Prepare stringToSign
		method := strings.ToUpper(c.Method())
		endpoint := c.Path()
		stringToSign := fmt.Sprintf("%s:%s:%s:%s", method, endpoint, bodyHash, timestamp)

		// 7. HMAC SHA512 with clientSecret → base64
		mac := hmac.New(sha512.New, []byte(os.Getenv("CLIENT_KEY")))
		mac.Write([]byte(stringToSign))
		expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		// 8. Compare (gunakan hmac.Equal untuk security)
		if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusForbidden,
				Message: "Invalid signature",
			})
		}

		return c.Next()
	}
}
