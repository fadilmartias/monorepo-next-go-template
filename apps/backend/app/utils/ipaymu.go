package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
)

func GenerateIPaymuSign(postBody any) (string, error) {
	ipaymuVa := config.LoadIPaymuConfig().VA
	ipaymuKey := config.LoadIPaymuConfig().APIKey
	//generate signature
	bodyBytes, err := sonic.Marshal(postBody)
	if err != nil {
		return "", err
	}
	bodyHash := sha256.Sum256(bodyBytes)
	bodyHashToString := hex.EncodeToString(bodyHash[:])
	stringToSign := "POST:" + ipaymuVa + ":" + strings.ToLower(string(bodyHashToString)) + ":" + ipaymuKey

	h := hmac.New(sha256.New, []byte(ipaymuKey))
	if _, err := h.Write([]byte(stringToSign)); err != nil {
		return "", err
	}
	signature := hex.EncodeToString(h.Sum(nil))
	return signature, nil
}

// end generate signatrure
