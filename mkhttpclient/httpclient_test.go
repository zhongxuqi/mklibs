package mkhttpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestHttpOption(t *testing.T) {
	if c, ok := NewHTTPClient("http://domain").(*httpClient); !ok {
		t.Fatalf("NewHttpClient error")
	} else {
		currConfig := parseBaseConfig(defaultBaseConfig, c.options)
		if c.host != "http://domain" || currConfig.MaxIdleConns != 10 || currConfig.IdleConnTimeout != time.Duration(30*time.Second) ||
			currConfig.DisableCompression != true || currConfig.RetryTimeout != 5*time.Second || currConfig.RetryTimes != 0 ||
			currConfig.TotalTimeout != 5*time.Second {
			t.Fatalf("NewHttpClient data error %+v", c)
		}
	}

	if c, ok := NewHTTPClient("http://domain", WithMaxIdleConns(100), WithIdleConnTimeout(time.Duration(300*time.Second)),
		WithDisableCompression(false), WithRetryTimes(3), WithRetryTimeout(time.Duration(10*time.Second)),
		WithTotalTimeout(time.Duration(15*time.Second))).(*httpClient); !ok {
		t.Fatalf("NewHttpClient error")
	} else {
		currConfig := parseBaseConfig(defaultBaseConfig, c.options)
		if c.host != "http://domain" || currConfig.MaxIdleConns != 100 || currConfig.IdleConnTimeout != time.Duration(300*time.Second) ||
			currConfig.DisableCompression != false || currConfig.RetryTimeout != 10*time.Second || currConfig.RetryTimes != 3 ||
			currConfig.TotalTimeout != 15*time.Second {
			t.Fatalf("NewHttpClient data error %+v", c)
		}
	}
}

type testRes struct {
	ErrNo  int64  `json:"errno"`
	ErrMsg string `json:"errmsg"`
}

func TestHttpRpc(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(testRes{
			ErrNo:  1,
			ErrMsg: r.Method,
		})
		if (r.Method == http.MethodPut || r.Method == http.MethodPost) && r.Header.Get("Content-Type") != ContentTypeJSON {
			b, _ = json.Marshal(testRes{
				ErrNo:  2,
				ErrMsg: "error",
			})
		}
		w.Write(b)
	})
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		server.ListenAndServe()
	}()
	defer server.Shutdown(context.TODO())

	client := NewHTTPClient("http://127.0.0.1:8080")

	// 测试Get
	var res testRes
	if err := client.Get(context.TODO(), "/rpc", map[string]interface{}{}, &res); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodGet {
		t.Fatalf("client.Get data error %+v", res)
	}
	if err := client.Delete(context.TODO(), "/rpc", &res); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodDelete {
		t.Fatalf("client.Get data error %+v", res)
	}
	if err := client.PostJSON(context.TODO(), "/rpc", map[string]interface{}{}, &res); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodPost {
		t.Fatalf("client.Get data error %+v", res)
	}
	if err := client.PutJSON(context.TODO(), "/rpc", map[string]interface{}{}, &res); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodPut {
		t.Fatalf("client.Get data error %+v", res)
	}
}

func TestHttpRpcEx(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(testRes{
			ErrNo:  1,
			ErrMsg: r.Method,
		})
		if (r.Method == http.MethodPut || r.Method == http.MethodPost) && r.Header.Get("Content-Type") != ContentTypeJSON {
			b, _ = json.Marshal(testRes{
				ErrNo:  2,
				ErrMsg: "error",
			})
		}
		if r.Header.Get("header-test") != "httprpc" {
			b, _ = json.Marshal(testRes{
				ErrNo:  2,
				ErrMsg: "error",
			})
		}
		w.Write(b)
	})
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		server.ListenAndServe()
	}()
	defer server.Shutdown(context.TODO())

	client := NewHTTPClient("http://127.0.0.1:8080")

	// 测试Get
	var res testRes
	if err := client.GetEx(context.TODO(), "/rpc", map[string]interface{}{}, &res, map[string]string{"header-test": "httprpc"}); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodGet {
		t.Fatalf("client.Get data error %+v", res)
	}
	if err := client.DeleteEx(context.TODO(), "/rpc", &res, map[string]string{"header-test": "httprpc"}); err != nil {
		t.Fatalf("client.Delete error %+v", err)
	} else if res.ErrMsg != http.MethodDelete {
		t.Fatalf("client.Delete data error %+v", res)
	}
	if err := client.PostJSONEx(context.TODO(), "/rpc", map[string]interface{}{}, &res, map[string]string{"header-test": "httprpc"}); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodPost {
		t.Fatalf("client.Get data error %+v", res)
	}
	if err := client.PutJSONEx(context.TODO(), "/rpc", map[string]interface{}{}, &res, map[string]string{"header-test": "httprpc"}); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodPut {
		t.Fatalf("client.Get data error %+v", res)
	}
}

func TestHttpRpcRetry(t *testing.T) {
	hasRPC := false
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		if !hasRPC {
			hasRPC = true
			time.Sleep(2 * time.Second)
		}
		b, _ := json.Marshal(testRes{
			ErrNo:  1,
			ErrMsg: r.Method,
		})
		w.Write(b)
	})
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		server.ListenAndServe()
	}()
	defer server.Shutdown(context.TODO())

	client := NewHTTPClient("http://127.0.0.1:8080", WithRetryTimeout(time.Second), WithTotalTimeout(time.Second))

	// 测试Get
	var res testRes
	if err := client.GetEx(context.TODO(), "/rpc", map[string]interface{}{}, &res, map[string]string{}, WithRetryTimes(1)); err != nil {
		t.Fatalf("client.Get error %+v", err)
	} else if res.ErrMsg != http.MethodGet {
		t.Fatalf("client.Get data error %+v", res)
	}
}
