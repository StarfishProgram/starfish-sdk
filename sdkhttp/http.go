package sdkhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"
)

type Config struct {
	URL                 string
	Method              string
	Header              http.Header
	Params              url.Values
	Timeout             time.Duration
	SuccessResponseCode []int
	Data                any
}

type Options func(config *Config)

// WithUrl URL
func WithUrl(url string) Options {
	return func(config *Config) {
		config.URL = url
	}
}

// WithGetMethod GET请求
func WithGetMethod() Options {
	return func(config *Config) {
		config.Method = http.MethodGet
	}
}

// WithPostMethod POST请求
func WithPostMethod() Options {
	return func(config *Config) {
		config.Method = http.MethodPost
	}
}

// WithParam 请求参数
func WithParam(key, val string) Options {
	return func(config *Config) {
		config.Params.Add(key, val)
	}
}

// WithHeader 请求头
func WithHeader(key, val string) Options {
	return func(config *Config) {
		config.Header.Add(key, val)
	}
}

// WithJsonContentType Content-Type : application/json
func WithJsonContentType() Options {
	return func(config *Config) {
		config.Header.Add("Content-Type", "application/json")
	}
}

// WithTimeout 超时时间
func WithTimeout(t time.Duration) Options {
	return func(config *Config) {
		config.Timeout = t
	}
}

// WithSuccessResponseCode 超时时间
func WithSuccessResponseCode(code ...int) Options {
	return func(config *Config) {
		config.SuccessResponseCode = append(config.SuccessResponseCode, code...)
	}
}

// WithData body data
func WithData(data any) Options {
	return func(config *Config) {
		config.Data = data
	}
}

func Request[R any](opts ...Options) (*R, error) {
	config := Config{
		URL:                 "",
		Method:              http.MethodGet,
		Header:              map[string][]string{},
		Params:              map[string][]string{},
		Timeout:             time.Second * 60,
		SuccessResponseCode: []int{http.StatusOK},
		Data:                nil,
	}
	for _, opt := range opts {
		opt(&config)
	}
	var requestUrl string
	if args := config.Params.Encode(); len(args) != 0 {
		requestUrl = config.URL + "?" + args
	} else {
		requestUrl = config.URL
	}
	var requestBody io.Reader
	if config.Data != nil {
		jsonData, err := json.Marshal(config.Data)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(jsonData)
	}
	client := http.Client{Timeout: config.Timeout}
	request, err := http.NewRequest(config.Method, requestUrl, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header = config.Header
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(config.SuccessResponseCode, response.StatusCode) {
		return nil, fmt.Errorf("响应状态码错误 : %d, 正常响应码 : %v", response.StatusCode, config.SuccessResponseCode)
	}

	var r R
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
