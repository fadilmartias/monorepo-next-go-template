package utils

import (
	"reflect"

	"github.com/bytedance/sonic"
)

func JSONParse(input any) map[string]any {
	var result map[string]any
	var raw []byte

	switch v := input.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return map[string]any{}
	}

	err := sonic.Unmarshal(raw, &result)
	if err != nil {
		return map[string]any{}
	}
	return result
}

func JSONStringify(v any) string {
	str, err := sonic.MarshalString(v)
	if err != nil {
		return ""
	}
	return str
}

func PrettySonicJSON(b []byte) string {
	var v any

	// Coba parse level pertama
	if err := sonic.Unmarshal(b, &v); err != nil {
		return string(b)
	}

	// Jika hasilnya string, berarti JSON di dalam string -> unmarshal ulang
	if str, ok := v.(string); ok {
		var inner any
		if err := sonic.Unmarshal([]byte(str), &inner); err == nil {
			v = inner
		} else {
			return str // fallback ke string mentah
		}
	}

	pretty, _ := sonic.MarshalIndent(v, "", "  ")
	return string(pretty)
}

// ToMap: converts struct or map into map[string]any
func ToMap(v any) map[string]any {
	out := map[string]any{}
	if v == nil {
		return out
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return out
		}
		rv = rv.Elem()
	}
	// Use JSON round-trip to convert reliably
	b, _ := sonic.Marshal(v)
	_ = sonic.Unmarshal(b, &out)
	return out
}

// Diff returns { "old": {k:v}, "new": {k:v} } for changed keys only
func Diff(oldObj, newObj any) map[string]map[string]any {
	oldMap := ToMap(oldObj)
	newMap := ToMap(newObj)

	res := map[string]map[string]any{
		"old": {},
		"new": {},
	}

	// collect keys
	keys := map[string]struct{}{}
	for k := range oldMap {
		keys[k] = struct{}{}
	}
	for k := range newMap {
		keys[k] = struct{}{}
	}

	for k := range keys {
		ov, ook := oldMap[k]
		nv, nok := newMap[k]
		if !rookEqual(ov, nv) || ook != nok {
			// only include when different (or missing)
			if ook {
				res["old"][k] = ov
			} else {
				res["old"][k] = nil
			}
			if nok {
				res["new"][k] = nv
			} else {
				res["new"][k] = nil
			}
		}
	}

	// if no changes, return empty
	if len(res["old"]) == 0 && len(res["new"]) == 0 {
		return nil
	}
	return res
}

// rookEqual: robust equality for JSON-serializable values
func rookEqual(a, b any) bool {
	ba, _ := sonic.Marshal(a)
	bb, _ := sonic.Marshal(b)
	return string(ba) == string(bb)
}
