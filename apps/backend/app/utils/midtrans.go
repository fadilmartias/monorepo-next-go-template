package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/fadilmartias/dilz_code/apps/backend/config"
)

func MidtransBasicAuth() string {
	midtransConfig := config.LoadMidtransConfig()
	serverKey := midtransConfig.SecretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(serverKey+":"))
}

func MidtransVerifyWebhookSignature(orderId string, statusCode string, grossAmount string, signatureKey string) (bool, error) {
	midtransConfig := config.LoadMidtransConfig()
	serverKey := midtransConfig.SecretKey
	input := orderId + statusCode + grossAmount + serverKey
	hash := sha512.New()
	hash.Write([]byte(input))
	signature := hex.EncodeToString(hash.Sum(nil))
	if signature != signatureKey {
		return false, errors.New("invalid signature")
	}
	return true, nil
}
