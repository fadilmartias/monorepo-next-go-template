package middleware

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func GetUser() fiber.Handler {
	return GetUserWithConfig(GetUserConfig{
		SupportCookie: true,
		CookieName:    defaultCookieName(),
		TestTokens: map[string]jwt.MapClaims{
			"TestAdmin123": {
				"id":    "A1",
				"name":  "Admin",
				"email": "admin@gmail.com",
				"phone": "08123456789",
				"role":  "admin",
			},
			"TestUser123": {
				"id":    "A2",
				"name":  "User",
				"email": "user@gmail.com",
				"phone": "08123456780",
				"role":  "user",
			},
		},
	})
}

type GetUserConfig struct {
	SupportCookie bool
	SupportBearer bool
	SupportAPIKey bool

	CookieName   string
	BearerName   string
	BearerPrefix string
	APIKeyName   string

	UserContextKey string
	AuthSourceKey  string

	ValidateToken func(string) (jwt.MapClaims, error)
	TestTokens    map[string]jwt.MapClaims
	APIKey        APIKeyConfig
}

type APIKeyLookupFunc func(ctx context.Context, key string) (any, bool, error)

type APIKeyConfig struct {
	DB             *gorm.DB
	Lookup         APIKeyLookupFunc
	Table          string
	Model          any
	Column         string
	Selected       string
	Joins          []string
	Preloads       []string
	MergeRelations []string
}

func (cfg APIKeyConfig) WithUserRelation(relation string) APIKeyConfig {
	if relation == "" {
		relation = "User"
	}
	cfg.Preloads = appendUnique(cfg.Preloads, relation)
	cfg.MergeRelations = appendUnique(cfg.MergeRelations, relation)
	return cfg
}

func GetUserWithConfig(cfg GetUserConfig) fiber.Handler {
	resolveDefaults(&cfg)

	return func(c fiber.Ctx) error {
		if cfg.SupportAPIKey {
			if claims, ok, err := authenticateAPIKey(c.Context(), c.Get(cfg.APIKeyName), cfg); err != nil {
				return err
			} else if ok {
				setUserLocals(c, cfg, claims, "api_key")
				return c.Next()
			}
		}

		if cfg.SupportBearer {
			if claims, ok := authenticateBearer(c.Get(cfg.BearerName), cfg); ok {
				setUserLocals(c, cfg, claims, "bearer")
				return c.Next()
			}
		}

		if cfg.SupportCookie {
			if claims, ok := authenticateCookie(c.Cookies(cfg.CookieName), cfg); ok {
				setUserLocals(c, cfg, claims, "cookie")
				return c.Next()
			}
		}

		return c.Next()
	}
}

func resolveDefaults(cfg *GetUserConfig) {
	if cfg.UserContextKey == "" {
		cfg.UserContextKey = "user"
	}
	if cfg.AuthSourceKey == "" {
		cfg.AuthSourceKey = "auth_source"
	}
	if cfg.CookieName == "" {
		cfg.CookieName = defaultCookieName()
	}
	if cfg.BearerName == "" {
		cfg.BearerName = "Authorization"
	}
	if cfg.BearerPrefix == "" {
		cfg.BearerPrefix = "Bearer"
	}
	if cfg.APIKeyName == "" {
		cfg.APIKeyName = "X-API-Key"
	}
	if cfg.ValidateToken == nil {
		cfg.ValidateToken = utils.ValidateToken
	}
}

func defaultCookieName() string {
	return "access_token_" + os.Getenv("APP_ENV")
}

func authenticateCookie(token string, cfg GetUserConfig) (jwt.MapClaims, bool) {
	return authenticateJWT(token, cfg)
}

func authenticateBearer(header string, cfg GetUserConfig) (jwt.MapClaims, bool) {
	if header == "" {
		return nil, false
	}

	prefix := cfg.BearerPrefix + " "
	if !strings.HasPrefix(strings.ToLower(header), strings.ToLower(prefix)) {
		return nil, false
	}

	token := strings.TrimSpace(header[len(prefix):])
	return authenticateJWT(token, cfg)
}

func authenticateJWT(token string, cfg GetUserConfig) (jwt.MapClaims, bool) {
	if token == "" {
		return nil, false
	}

	if claims, ok := cfg.TestTokens[token]; ok {
		return cloneClaims(claims), true
	}

	claims, err := cfg.ValidateToken(token)
	if err != nil {
		return nil, false
	}

	return claims, true
}

func authenticateAPIKey(ctx context.Context, key string, cfg GetUserConfig) (jwt.MapClaims, bool, error) {
	if key == "" {
		return nil, false, nil
	}

	if cfg.APIKey.Lookup != nil {
		value, found, err := cfg.APIKey.Lookup(ctx, key)
		if err != nil || !found {
			return nil, found, err
		}
		return claimsFromValue(value), true, nil
	}

	if cfg.APIKey.DB == nil || cfg.APIKey.Model == nil || cfg.APIKey.Column == "" {
		return nil, false, fmt.Errorf("api key auth is enabled but APIKey.DB, APIKey.Model, and APIKey.Column must be set")
	}

	dest := newModelDestination(cfg.APIKey.Model)
	query := cfg.APIKey.DB.Model(dest)
	for _, join := range cfg.APIKey.Joins {
		query = query.Joins(join)
	}
	for _, preload := range cfg.APIKey.Preloads {
		query = query.Preload(preload)
	}
	if cfg.APIKey.Table != "" {
		query = cfg.APIKey.DB.Table(cfg.APIKey.Table)
	}

	err := query.Where(cfg.APIKey.Column+" = ?", key).First(dest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	claims := claimsFromValue(reflect.Indirect(reflect.ValueOf(dest)).Interface())
	for _, relation := range cfg.APIKey.MergeRelations {
		claims = mergeClaims(claims, claimsFromField(dest, relation))
	}

	return claims, true, nil
}

func newModelDestination(model any) any {
	if model == nil {
		return nil
	}

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		return reflect.New(v.Type().Elem()).Interface()
	}

	return reflect.New(v.Type()).Interface()
}

func claimsFromValue(value any) jwt.MapClaims {
	if value == nil {
		return jwt.MapClaims{}
	}

	if claims, ok := value.(jwt.MapClaims); ok {
		return cloneClaims(claims)
	}
	if claims, ok := value.(map[string]any); ok {
		return jwt.MapClaims(claims)
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return jwt.MapClaims{}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return jwt.MapClaims{"value": value}
	}

	rt := rv.Type()
	claims := jwt.MapClaims{}
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = strings.ToLower(field.Name[:1]) + field.Name[1:]
		} else if comma := strings.Index(name, ","); comma >= 0 {
			name = name[:comma]
		}
		if name == "" || name == "-" {
			continue
		}
		claims[name] = rv.Field(i).Interface()
	}

	return claims
}

func claimsFromField(value any, fieldName string) jwt.MapClaims {
	if value == nil || fieldName == "" {
		return jwt.MapClaims{}
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return jwt.MapClaims{}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return jwt.MapClaims{}
	}

	field := rv.FieldByName(fieldName)
	if !field.IsValid() {
		for i := 0; i < rv.NumField(); i++ {
			if strings.EqualFold(rv.Type().Field(i).Name, fieldName) {
				field = rv.Field(i)
				break
			}
		}
	}
	if !field.IsValid() {
		return jwt.MapClaims{}
	}

	return claimsFromValue(field.Interface())
}

func mergeClaims(dst, src jwt.MapClaims) jwt.MapClaims {
	if dst == nil {
		dst = jwt.MapClaims{}
	}
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func appendUnique(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func cloneClaims(src jwt.MapClaims) jwt.MapClaims {
	if src == nil {
		return jwt.MapClaims{}
	}
	out := make(jwt.MapClaims, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}

func setUserLocals(c fiber.Ctx, cfg GetUserConfig, claims jwt.MapClaims, authSource string) {
	if claims == nil {
		claims = jwt.MapClaims{}
	}
	c.Locals(cfg.UserContextKey, claims)
	c.Locals(cfg.AuthSourceKey, authSource)
}
