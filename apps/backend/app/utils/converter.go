package utils

import (
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
)

func StringPtr(v any) *string {
	if v == nil {
		return nil
	}

	// Coba konversi langsung ke string
	if str, ok := v.(string); ok {
		if str == "" || str == "<nil>" {
			return nil
		}
		return &str
	}

	// Kalau bukan string tapi bisa diformat jadi string
	str := fmt.Sprintf("%v", v)
	if str == "" || str == "<nil>" || str == "map[]" {
		return nil
	}

	return &str
}

func IntPtr(v any) *int {
	if v == nil {
		return nil
	}

	// Coba konversi langsung ke int
	if num, ok := v.(int); ok {
		return &num
	}

	// Kalau bukan int tapi bisa diformat jadi int
	str := fmt.Sprintf("%v", v)
	if str == "" || str == "<nil>" || str == "map[]" {
		return nil
	}

	// Coba konversi string ke int
	if num, err := strconv.Atoi(str); err == nil {
		return &num
	}

	return nil
}

func StructToMap(data any) map[string]any {
	var result map[string]any
	jsonBytes, _ := sonic.Marshal(data)
	sonic.Unmarshal(jsonBytes, &result)
	return result
}

// helper biar gak nulis ParseFloat terus
func ToFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0 // atau bisa return error kalau mau strict
	}
	return f
}

func ToInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // atau bisa panic / return error kalau mau strict
	}
	return i
}

func ToBool(s string) bool {
	if s == "" {
		return false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false // atau bisa panic / return error kalau mau strict
	}
	return b
}

func BoolPtr(s string) *bool {
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil // atau bisa panic / return error kalau mau strict
	}
	return &b
}

func MarshalToString(data any) (string, error) {
	jsonBytes, err := sonic.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func DerefOrEmpty[T any](val *T) T {
	if val == nil {
		var empty T
		return empty
	}
	return *val
}

func IsNotNil[T any](val *T) bool {
	return val != nil
}
