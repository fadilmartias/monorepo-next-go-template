package utils

import "errors"

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("invalid token")
)

var ErrorCodeMap = map[error]string{
	ErrTokenExpired: "ERR_TOKEN_EXPIRED",
	ErrTokenInvalid: "ERR_TOKEN_INVALID",
}

func TranslateErrorCode(err error) string {
	for key, code := range ErrorCodeMap {
		if errors.Is(err, key) {
			return code
		}
	}
	return "ERR_UNKNOWN"
}
