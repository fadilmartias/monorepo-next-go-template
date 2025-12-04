package utils

func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetMap(data any) (map[string]any, bool) {
	m, ok := data.(map[string]any)
	return m, ok
}

func GetSlice(data any) ([]any, bool) {
	s, ok := data.([]any)
	return s, ok
}

// FindIndex — mirip Array.findIndex di JS
func FindIndex[T any](items []T, predicate func(T) bool) int {
	for i, item := range items {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// IndexOf — mirip Array.indexOf di JS (butuh comparable type)
func IndexOf[T comparable](items []T, value T) int {
	for i, item := range items {
		if item == value {
			return i
		}
	}
	return -1
}

// Filter — mirip Array.filter di JS
func Filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Some — mirip Array.some di JS
func Some[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Find — mirip Array.find di JS
func Find[T any](items []T, predicate func(T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Map — mirip Array.map di JS
func Map[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}
	return result
}
