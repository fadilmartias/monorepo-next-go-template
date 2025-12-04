package utils

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/davecgh/go-spew/spew"
)

// Dump menggunakan spew untuk nge-print detail object
func Dump(v ...any) {
	for _, val := range v {
		spew.Dump(val)
	}
}

// Sdump mengembalikan string hasil spew.Sdump
func Sdump(v any) string {
	return spew.Sdump(v)
}

// PrettyJSON print data dalam format JSON terformat (jika bisa di-marshal)
func PrettyJSON(v any) {
	b, err := sonic.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("[PrettyJSON] error:", err)
		return
	}
	fmt.Println(string(b))
}

// JSONString mengembalikan string JSON terformat
func JSONString(v any) string {
	b, err := sonic.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("[JSONString] error: %v", err)
	}
	return string(b)
}

// SmartDump mencoba print sebagai JSON, kalau gagal fallback ke spew
func SmartDump(v any) {
	b, err := sonic.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
		return
	}
	spew.Dump(v)
}
