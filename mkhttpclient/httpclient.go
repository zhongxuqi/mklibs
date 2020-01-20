package mkhttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
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
	PostFileEx(ctx context.Context, path string, files map[string][]byte, res interface{}, header map[string]string, options ...HTTPClientOption) error
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
	ml.Infof("%s%s", s.host, path)
	return s.do(ctx, http.MethodGet, fmt.Sprintf("%s%s", s.host, path), nil, res, header, options...)
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
	ml.Infof("url %+v%+v", s.host, path)
	return s.do(ctx, http.MethodGet, fmt.Sprintf("%s%s", s.host, path), nil, res, nil)
}

// DeleteEx ...
func (s *httpClient) DeleteEx(ctx context.Context, path string, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	ml.Infof("url [%s]%s%s header %+v", http.MethodDelete, s.host, path, header)
	return s.do(ctx, http.MethodDelete, fmt.Sprintf("%s%s", s.host, path), nil, res, header, options...)
}

// Delete ...
func (s *httpClient) Delete(ctx context.Context, path string, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	ml.Infof("url [%s]%s%s", http.MethodDelete, s.host, path)
	return s.do(ctx, http.MethodDelete, fmt.Sprintf("%s%s", s.host, path), nil, res, nil)
}

// PutJSONEx ...
func (s *httpClient) PutJSONEx(ctx context.Context, path string, params interface{}, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	ml.Infof("url [%s]%s%s params %+v header %+v", http.MethodPut, s.host, path, params, header)
	paramsByte := make([]byte, 0)
	if params != nil {
		var err error
		paramsByte, err = json.Marshal(params)
		if err != nil {
			ml.Errorf("json.Marshal error %+v", err.Error())
			return err
		}
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = ContentTypeJSON
	ml.Infof("params %+v", string(paramsByte))
	return s.do(ctx, http.MethodPut, fmt.Sprintf("%s%s", s.host, path), paramsByte, res, header, options...)
}

// PutJSON ...
func (s *httpClient) PutJSON(ctx context.Context, path string, params interface{}, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	ml.Infof("url [%s]%s%s params %+v", http.MethodPut, s.host, path, params)
	paramsByte := make([]byte, 0)
	if params != nil {
		var err error
		paramsByte, err = json.Marshal(params)
		if err != nil {
			ml.Errorf("json.Marshal error %+v", err.Error())
			return err
		}
	}
	ml.Infof("params %+v", string(paramsByte))
	return s.do(ctx, http.MethodPut, fmt.Sprintf("%s%s", s.host, path), paramsByte, res, map[string]string{"Content-Type": ContentTypeJSON})
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
	ml.Infof("params %+v", string(paramsByte))
	header["Content-Type"] = ContentTypeJSON
	return s.do(ctx, http.MethodPost, fmt.Sprintf("%s%s", s.host, path), paramsByte, res, header, options...)
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
	ml.Infof("[POST]%+v body: %+v", fmt.Sprintf("%s%s", s.host, path), string(paramsByte))
	return s.do(ctx, http.MethodPost, fmt.Sprintf("%s%s", s.host, path), paramsByte, res, map[string]string{"Content-Type": ContentTypeJSON})
}

// PostEx ...
func (s *httpClient) PostEx(ctx context.Context, path string, params map[string]string, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	paramsByte := make([]byte, 0)
	if params != nil {
		values := make(url.Values)
		for k, v := range params {
			values.Add(k, v)
		}
		paramsByte = []byte(values.Encode())
	}
	ml.Infof("[POST]%+v body: %+v", fmt.Sprintf("%s%s", s.host, path), string(paramsByte))
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = ContentTypeForm
	return s.do(ctx, http.MethodPost, fmt.Sprintf("%s%s", s.host, path), paramsByte, res, header, options...)
}

// Post ...
func (s *httpClient) Post(ctx context.Context, path string, params map[string]string, res interface{}) error {
	ml := mklog.NewWithContext(ctx)
	paramsByte := make([]byte, 0)
	if params != nil {
		values := make(url.Values)
		for k, v := range params {
			values.Add(k, v)
		}
		paramsByte = []byte(values.Encode())
	}
	ml.Infof("[POST]%+v body: %+v", fmt.Sprintf("%s%s", s.host, path), string(paramsByte))
	return s.do(ctx, http.MethodPost, fmt.Sprintf("%s%s", s.host, path), paramsByte, res, map[string]string{"Content-Type": ContentTypeForm})
}

func (s *httpClient) PostFileEx(ctx context.Context, path string, files map[string][]byte, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	ml.Infof("url [%s]%s%s", http.MethodPost, s.host, path)
	body := bytes.NewBuffer([]byte(""))
	writer := multipart.NewWriter(body)
	var contentLength int64
	for fileName, fileContent := range files {
		part, err := writer.CreateFormFile(fileName, fileName)
		if err != nil {
			ml.Errorf("writer.CreateFormFile error %+v", err.Error())
			return err
		}
		_, err = io.Copy(part, bytes.NewReader(fileContent))
		if err != nil {
			ml.Errorf("bytes.NewReader error %+v", err.Error())
			return err
		}
		contentLength += int64(len(fileContent))
		err = writer.Close()
		if err != nil {
			ml.Errorf("writer.Close error %+v", err.Error())
			return err
		}
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = writer.FormDataContentType()
	//header["Content-Length"] = fmt.Sprintf("%d", contentLength)
	return s.do(ctx, http.MethodPost, fmt.Sprintf("%s%s", s.host, path), body.Bytes(), res, header, options...)
}

func (s *httpClient) do(ctx context.Context, method, url string, bodyByte []byte, res interface{}, header map[string]string, options ...HTTPClientOption) error {
	ml := mklog.NewWithContext(ctx)
	allOptions := make([]HTTPClientOption, 0, len(s.options)+len(options))
	allOptions = append(allOptions, s.options...)
	allOptions = append(allOptions, options...)
	currConfig := parseBaseConfig(defaultBaseConfig, allOptions)
	var httpRes *http.Response

	resChan := make(chan *http.Response)
	for i := 0; i < currConfig.RetryTimes+1; i++ {
		var buf *bytes.Buffer = nil
		if bodyByte != nil {
			buf = bytes.NewBuffer(bodyByte)
		} else {
			buf = bytes.NewBuffer([]byte(""))
		}
		req, err := http.NewRequest(method, url, buf)
		if err != nil {
			ml.Errorf("http.NewRequest error %+v", err.Error())
			return err
		}

		// add http headers
		req.Header.Set(common.HttpLogID, ml.GetLogID())
		for k, v := range header {
			req.Header.Add(k, v)
		}

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
				if httpStatusCode, err := parseRes(ctx, httpRes, res); httpStatusCode/100 != 2 && currConfig.RetryHttpError {
					ml.Errorf("parseRes error %+v", err)
				} else if err != nil {
					ml.Errorf("parseRes error %+v", err)
					return err
				} else {
					return nil
				}
			case err := <-errChan:
				ml.Errorf("client.Do error %+v", err)
			case <-time.After(currConfig.RetryTimeout):
				ml.Infof("retry %+v", i)
			}
		} else {
			select {
			case httpRes = <-resChan:
				if _, err := parseRes(ctx, httpRes, res); err != nil {
					ml.Errorf("parseRes error %+v", err)
					return err
				}
				return nil
			case err := <-errChan:
				ml.Errorf("http error %+v", err)
				return err
			case <-time.After(currConfig.TotalTimeout):
				ml.Errorf("req %+v error timeout", req)
				return fmt.Errorf("req %+v error timeout", req)
			}
		}
	}
	bodyStr := ""
	if bodyByte != nil {
		bodyStr = string(bodyByte)
	}
	ml.Errorf("req %+v error timeout", bodyStr)
	return fmt.Errorf("req %+v error timeout", bodyStr)
}

func parseRes(ctx context.Context, httpRes *http.Response, res interface{}) (int, error) {
	ml := mklog.NewWithContext(ctx)
	bodyByte, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		ml.Errorf("ioutil.ReadAll error %+v", err)
		return httpRes.StatusCode, err
	}
	ml.Infof("response body: %+v", string(bodyByte))
	if httpRes.StatusCode/100 != 2 {
		return httpRes.StatusCode, fmt.Errorf("http error %d %s", httpRes.StatusCode, httpRes.Status)
	}
	err = json.Unmarshal(bodyByte, res)
	if err != nil {
		ml.Errorf("json.Unmarshal error %+v", err)
	}
	if httpRes.StatusCode/100 != 2 {
		return httpRes.StatusCode, fmt.Errorf("http error %d %s", httpRes.StatusCode, httpRes.Status)
	}
	return httpRes.StatusCode, nil
}
