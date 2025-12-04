package utils

import (
	"fmt"
	"net/url"

	"github.com/bytedance/sonic"
)

func BuildQueryURLSonic(check_account_sku string, jsonDef string, input map[string]string) (string, error) {
	var arr []map[string]any
	if err := sonic.Unmarshal([]byte(jsonDef), &arr); err != nil {
		return "", err
	}

	values := url.Values{}

	for _, field := range arr {
		var paramKey string
		// Ambil key param
		if v, ok := field["check_name"].(string); ok && v != "" {
			paramKey = v
		}

		// Masukkan value ke query param
		if name, ok := field["name"].(string); ok {
			if val, exists := input[name]; exists && paramKey != "" {
				values.Set(paramKey, val)
			}
		}
	}

	return fmt.Sprintf("https://api.isan.eu.org/nickname/%s?%s&decode=false", check_account_sku, values.Encode()), nil
}
