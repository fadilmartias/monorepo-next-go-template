package client

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
	"github.com/gofiber/fiber/v3"
	fiberclient "github.com/gofiber/fiber/v3/client"
	"go.uber.org/zap"
)

type HttpClient struct {
	client    *fiberclient.Client
	log       *zap.SugaredLogger
	headers   map[string]string
	query     map[string]string
	body      any
	formData  map[string]string
	timeout   time.Duration
	userAgent string
	authToken string
}

func newBaseClient() *fiberclient.Client {
	client := fiberclient.New()
	client.SetTimeout(10 * time.Second)
	client.SetHeader("Accept", "application/json")
	client.SetHeader("User-Agent", "DilzTopup-HttpClient/2.0")
	return client
}

func Http() *HttpClient {
	return &HttpClient{
		client:    newBaseClient(),
		log:       logger.Base(),
		headers:   map[string]string{},
		query:     map[string]string{},
		timeout:   10 * time.Second,
		userAgent: "DilzTopup-HttpClient/2.0",
	}
}

func HttpWithCtx(c fiber.Ctx) *HttpClient {
	return &HttpClient{
		client:    newBaseClient(),
		log:       logger.Ctx(c),
		headers:   map[string]string{},
		query:     map[string]string{},
		timeout:   10 * time.Second,
		userAgent: "DilzTopup-HttpClient/2.0",
	}
}

func (h *HttpClient) WithHeader(key, value string) *HttpClient {
	h.headers[key] = value
	return h
}

func (h *HttpClient) WithHeaders(headers map[string]string) *HttpClient {
	for key, value := range headers {
		h.headers[key] = value
	}
	return h
}

func (h *HttpClient) WithQuery(params map[string]string) *HttpClient {
	for key, value := range params {
		h.query[key] = value
	}
	return h
}

func (h *HttpClient) WithJSON(body any) *HttpClient {
	h.body = body
	h.formData = nil
	h.WithHeader("Content-Type", "application/json")
	return h
}

func (h *HttpClient) WithFormData(data map[string]string) *HttpClient {
	h.formData = make(map[string]string, len(data))
	for key, value := range data {
		h.formData[key] = value
	}
	h.body = nil
	h.WithHeader("Content-Type", "application/x-www-form-urlencoded")
	return h
}

func (h *HttpClient) WithAuthToken(token string) *HttpClient {
	h.authToken = token
	h.WithHeader("Authorization", "Bearer "+token)
	return h
}

func (h *HttpClient) WithTimeout(d time.Duration) *HttpClient {
	h.timeout = d
	h.client.SetTimeout(d)
	return h
}

func (h *HttpClient) Get(url string) (*fiberclient.Response, error) {
	return h.do("GET", url)
}

func (h *HttpClient) Post(url string) (*fiberclient.Response, error) {
	return h.do("POST", url)
}

func (h *HttpClient) Put(url string) (*fiberclient.Response, error) {
	return h.do("PUT", url)
}

func (h *HttpClient) Delete(url string) (*fiberclient.Response, error) {
	return h.do("DELETE", url)
}

func (h *HttpClient) AsJSON(resp *fiberclient.Response, out any) error {
	if resp == nil {
		return fmt.Errorf("empty response")
	}
	if err := sonic.Unmarshal(resp.Body(), out); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func (h *HttpClient) do(method, rawURL string) (*fiberclient.Response, error) {
	cfg := h.requestConfig()
	h.log.Infof("[HTTP Client Request] %s %s", method, rawURL)

	var (
		resp *fiberclient.Response
		err  error
	)

	switch method {
	case "GET":
		resp, err = h.client.Get(rawURL, cfg)
	case "POST":
		resp, err = h.client.Post(rawURL, cfg)
	case "PUT":
		resp, err = h.client.Put(rawURL, cfg)
	case "DELETE":
		resp, err = h.client.Delete(rawURL, cfg)
	default:
		resp, err = h.client.Custom(rawURL, method, cfg)
	}

	if err != nil {
		h.log.Errorf("[HTTP Client Error] %s %s | err: %v", method, rawURL, err)
		return resp, err
	}

	if resp != nil && resp.StatusCode() >= 400 {
		body := string(resp.Body())
		errMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), body)
		h.log.Errorf("[HTTP Client Failed] %d | body: %s", resp.StatusCode(), body)
		return resp, fmt.Errorf("%s", errMsg)
	}

	if resp != nil {
		h.log.Infof("[HTTP Client Success] %d | body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return resp, nil
}

func (h *HttpClient) requestConfig() fiberclient.Config {
	cfg := fiberclient.Config{Timeout: h.timeout}

	if len(h.headers) > 0 {
		cfg.Header = cloneStringMap(h.headers)
	}
	if len(h.query) > 0 {
		cfg.Param = cloneStringMap(h.query)
	}
	if len(h.formData) > 0 {
		cfg.FormData = cloneStringMap(h.formData)
	}
	if h.body != nil {
		cfg.Body = h.body
	}
	if h.authToken != "" {
		if cfg.Header == nil {
			cfg.Header = map[string]string{}
		}
		cfg.Header["Authorization"] = "Bearer " + h.authToken
	}

	return cfg
}

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]string, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}
