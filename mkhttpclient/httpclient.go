package mkhttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zhongxuqi/mklibs/common"
	"github.com/zhongxuqi/mklibs/mklog"
)

// const ...
const (
	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

var (
	defaultBaseConfig = baseConfig{
		MaxIdleConns:       10,
		IdleConnTimeout:    time.Duration(30 * time.Second),
		DisableCompression: true,
		RetryTimes:         0,
		RetryTimeout:       time.Duration(5 * time.Second),
		TotalTimeout:       time.Duration(5 * time.Second),
	}
)

// HTTPClient ...
type HTTPClient interface {
	GetEx(ctx context.Context, path string, params map[string]interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error
	Get(ctx context.Context, path string, params map[string]interface{}, res interface{}) error
	DeleteEx(ctx context.Context, path string, res interface{}, header map[string]string, options ...HTTPClientOption) error
	Delete(ctx context.Context, path string, res interface{}) error
	PutJSONEx(ctx context.Context, path string, params interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error
	PutJSON(ctx context.Context, path string, params interface{}, res interface{}) error
	PostJSONEx(ctx context.Context, path string, params interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error
	PostJSON(ctx context.Context, path string, params interface{}, res interface{}) error
	PostEx(ctx context.Context, path string, params map[string]string, res interface{}, header map[string]string, options ...HTTPClientOption) error
	Post(ctx context.Context, path string, params map[string]string, res interface{}) error
}

type httpClient struct {
	host    string
	options []HTTPClientOption
	client  *http.Client
}

// NewHTTPClient ...
func NewHTTPClient(host string, options ...HTTPClientOption) HTTPClient {
	currConfig := parseBaseConfig(defaultBaseConfig, options)
	return &httpClient{
		host:    host,
		options: options,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       currConfig.MaxIdleConns,
				IdleConnTimeout:    currConfig.IdleConnTimeout,
				DisableCompression: currConfig.DisableCompression,
			},
		},
	}
}

// GetEx ...
func (s *httpClient) GetEx(ctx context.Context, path string, params map[string]interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	if len(params) > 0 {
		path += "?"
		for k, v := range params {
			path += fmt.Sprintf("%s=%+v&", k, v)
		}
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", s.host, path), nil)
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, header, options...)
}

// Get ...
func (s *httpClient) Get(ctx context.Context, path string, params map[string]interface{}, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	if len(params) > 0 {
		path += "?"
		for k, v := range params {
			path += fmt.Sprintf("%s=%+v&", k, v)
		}
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", s.host, path), nil)
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, nil)
}

// DeleteEx ...
func (s *httpClient) DeleteEx(ctx context.Context, path string, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", s.host, path), nil)
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, header, options...)
}

// Delete ...
func (s *httpClient) Delete(ctx context.Context, path string, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", s.host, path), nil)
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, nil)
}

// PutJSONEx ...
func (s *httpClient) PutJSONEx(ctx context.Context, path string, params interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	paramsByte := make([]byte, 0)
	if params != nil {
		var err error
		paramsByte, err = json.Marshal(params)
		if err != nil {
			ml.Errorf("json.Marshal error %+v", err)
			return err
		}
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = ContentTypeJSON
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", s.host, path), bytes.NewBuffer(paramsByte))
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, header, options...)
}

// PutJSON ...
func (s *httpClient) PutJSON(ctx context.Context, path string, params interface{}, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	paramsByte := make([]byte, 0)
	if params != nil {
		var err error
		paramsByte, err = json.Marshal(params)
		if err != nil {
			ml.Errorf("json.Marshal error %+v", err)
			return err
		}
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", s.host, path), bytes.NewBuffer(paramsByte))
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, map[string]string{"Content-Type": ContentTypeJSON})
}

// PostJSONEx ...
func (s *httpClient) PostJSONEx(ctx context.Context, path string, params interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	paramsByte := make([]byte, 0)
	if params != nil {
		var err error
		paramsByte, err = json.Marshal(params)
		if err != nil {
			ml.Errorf("json.Marshal error %+v", err)
			return err
		}
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = ContentTypeJSON
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", s.host, path), bytes.NewBuffer(paramsByte))
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, header, options...)
}

// PostJSON ...
func (s *httpClient) PostJSON(ctx context.Context, path string, params interface{}, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	paramsByte := make([]byte, 0)
	if params != nil {
		var err error
		paramsByte, err = json.Marshal(params)
		if err != nil {
			ml.Errorf("json.Marshal error %+v", err)
			return err
		}
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", s.host, path), bytes.NewBuffer(paramsByte))
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, map[string]string{"Content-Type": ContentTypeJSON})
}

// PostEx ...
func (s *httpClient) PostEx(ctx context.Context, path string, params map[string]string, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	paramsStr := ""
	for k, v := range params {
		if paramsStr != "" {
			paramsStr += "&"
		}
		paramsStr += k + "=" + v
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = ContentTypeForm
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", s.host, path), bytes.NewBufferString(paramsStr))
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, header, options...)
}

// Post ...
func (s *httpClient) Post(ctx context.Context, path string, params map[string]string, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	paramsStr := ""
	for k, v := range params {
		if paramsStr != "" {
			paramsStr += "&"
		}
		paramsStr += k + "=" + v
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", s.host, path), bytes.NewBufferString(paramsStr))
	if err != nil {
		ml.Errorf("http.NewRequest error %+v", err)
		return err
	}
	return s.do(ctx, req, res, map[string]string{"Content-Type": ContentTypeForm})
}

func (s *httpClient) do(ctx context.Context, req *http.Request, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	allOptions := make([]HTTPClientOption, 0, len(s.options)+len(options))
	allOptions = append(allOptions, s.options...)
	allOptions = append(allOptions, options...)
	currConfig := parseBaseConfig(defaultBaseConfig, allOptions)
	var httpRes *http.Response

	// add http headers
	req.Header.Set(common.HttpLogID, ml.GetLogID())
	for k, v := range header {
		req.Header.Add(k, v)
	}

	resChan := make(chan *http.Response)
	for i := 0; i < currConfig.RetryTimes+1; i++ {
		errChan := make(chan error)
		go func() {
			res, err := s.client.Do(req)
			if err != nil {
				ml.Errorf("client.Do error %+v", err)
				errChan <- err
				return
			}
			resChan <- res
		}()
		if i < currConfig.RetryTimes {
			select {
			case httpRes = <-resChan:
				if err := parseRes(ctx, httpRes, res); err != nil {
					ml.Errorf("parseRes error %+v", err)
					return err
				}
				return nil
			case err := <-errChan:
				ml.Errorf("client.Do error %+v", err)
			case <-time.Tick(currConfig.RetryTimeout):
				ml.Infof("retry %+v", i)
			}
		} else {
			select {
			case httpRes = <-resChan:
				if err := parseRes(ctx, httpRes, res); err != nil {
					ml.Errorf("parseRes error %+v", err)
					return err
				}
				return nil
			case err := <-errChan:
				ml.Errorf("http error %+v", err)
				return err
			case <-time.Tick(currConfig.TotalTimeout):
				ml.Errorf("req %+v error timeout", req)
				return fmt.Errorf("req %+v error timeout", req)
			}
		}
	}
	ml.Errorf("req %+v error timeout", req)
	return fmt.Errorf("req %+v error timeout", req)
}

func parseRes(ctx context.Context, httpRes *http.Response, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	bodyByte, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		ml.Errorf("ioutil.ReadAll error %+v", err)
	}
	err = json.Unmarshal(bodyByte, res)
	if err != nil {
		ml.Errorf("json.Unmarshal error %+v", err)
	}
	if httpRes.StatusCode/100 != 2 {
		return fmt.Errorf("http error %d %s", httpRes.StatusCode, httpRes.Status)
	}
	if httpRes.StatusCode/100 != 2 {
		return fmt.Errorf("http error %d %s %s", httpRes.StatusCode, httpRes.Status, string(bodyByte))
	}
	return nil
}
