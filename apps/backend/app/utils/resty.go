package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HttpClient wrapper untuk Resty dengan fitur tambahan
type HttpClient struct {
	client  *resty.Client
	request *resty.Request
	logger  *zap.Logger
}

// newBaseClient konfigurasi dasar Resty
func newBaseClient() *resty.Client {
	client := resty.New().
		SetTimeout(10*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "DilzTopup-HttpClient/2.0")

	// Exponential backoff
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		return r.StatusCode() >= 500
	})

	return client
}

// Http mengembalikan instance baru siap pakai
func Http() *HttpClient {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // warna juga kalau mau
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	))

	client := newBaseClient()

	h := &HttpClient{
		client:  client,
		request: client.R(),
		logger:  logger,
	}

	// Logging Interceptor
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		h.logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("url", r.URL),
			zap.Any("body", r.Body),
		)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		h.logger.Info("HTTP Response",
			zap.String("method", r.Request.Method),
			zap.String("url", r.Request.URL),
			zap.String("body", PrettySonicJSON(r.Body())),
			zap.Int("status", r.StatusCode()),
			zap.Duration("duration", r.Time()),
		)
		return nil
	})

	return h
}

// -------------------- CHAINABLE METHODS -------------------- //

func (h *HttpClient) WithHeader(key, value string) *HttpClient {
	h.request.SetHeader(key, value)
	return h
}

func (h *HttpClient) WithHeaders(headers map[string]string) *HttpClient {
	for k, v := range headers {
		h.request.SetHeader(k, v)
	}
	return h
}

func (h *HttpClient) WithQuery(params map[string]string) *HttpClient {
	h.request.SetQueryParams(params)
	return h
}

func (h *HttpClient) WithJSON(body any) *HttpClient {
	h.request.SetBody(body)
	h.request.SetHeader("Content-Type", "application/json")
	return h
}

func (h *HttpClient) WithFormData(data map[string]string) *HttpClient {
	h.request.SetFormData(data)
	h.request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return h
}

func (h *HttpClient) WithAuthToken(token string) *HttpClient {
	h.request.SetAuthToken(token)
	return h
}

func (h *HttpClient) WithTimeout(d time.Duration) *HttpClient {
	h.client.SetTimeout(d)
	return h
}

// -------------------- REQUEST METHODS -------------------- //

func (h *HttpClient) Get(url string) (*resty.Response, error) {
	resp, err := h.request.Get(url)
	return h.wrapError(resp, err)
}

func (h *HttpClient) Post(url string) (*resty.Response, error) {
	resp, err := h.request.Post(url)
	return h.wrapError(resp, err)
}

func (h *HttpClient) Put(url string) (*resty.Response, error) {
	resp, err := h.request.Put(url)
	return h.wrapError(resp, err)
}

func (h *HttpClient) Delete(url string) (*resty.Response, error) {
	resp, err := h.request.Delete(url)
	return h.wrapError(resp, err)
}

// -------------------- UTILITIES -------------------- //

func (h *HttpClient) AsJSON(resp *resty.Response, out any) error {
	if resp == nil {
		return fmt.Errorf("empty response")
	}
	if err := json.Unmarshal(resp.Body(), out); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

// wrapError membuat error lebih informatif
func (h *HttpClient) wrapError(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		h.logger.Error("HTTP Error", zap.Error(err))
		return resp, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), resp.String())
		h.logger.Error("HTTP Response Error", zap.String("body", resp.String()))
		return resp, fmt.Errorf(errMsg)
	}
	return resp, nil
}
