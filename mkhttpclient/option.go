package mkhttpclient

import "time"

// OptionKey ...
type OptionKey int

// const ...
const (
	OptionKeyMaxIdleConns OptionKey = iota
	OptionKeyIdleConnTimeout
	OptionKeyDisableCompression
	OptionKeyRetryTimes
	OptionKeyRetryTimeout
	OptionKeyTotalTimeout
	OptionKeyOdinMetrics
)

// HTTPClientOption ...
type HTTPClientOption interface {
	OptionKey() OptionKey
	OptionValue() interface{}
}

// httpClientOption ...
type httpClientOption struct {
	optionKey   OptionKey
	optionValue interface{}
}

// OptionKey ...
func (s *httpClientOption) OptionKey() OptionKey {
	return s.optionKey
}

// OptionValue ...
func (s *httpClientOption) OptionValue() interface{} {
	return s.optionValue
}

// baseConfig ...
type baseConfig struct {
	MaxIdleConns       int
	IdleConnTimeout    time.Duration
	DisableCompression bool
	RetryTimes         int
	RetryTimeout       time.Duration
	TotalTimeout       time.Duration
}

// parseBaseConfig ...
func parseBaseConfig(defaultOption baseConfig, options []HTTPClientOption) baseConfig {
	res := defaultOption
	for _, option := range options {
		switch option.OptionKey() {
		case OptionKeyMaxIdleConns:
			if v, ok := option.OptionValue().(int); ok {
				res.MaxIdleConns = v
			}
		case OptionKeyIdleConnTimeout:
			if v, ok := option.OptionValue().(time.Duration); ok {
				res.IdleConnTimeout = v
			}
		case OptionKeyDisableCompression:
			if v, ok := option.OptionValue().(bool); ok {
				res.DisableCompression = v
			}
		case OptionKeyRetryTimes:
			if v, ok := option.OptionValue().(int); ok {
				res.RetryTimes = v
			}
		case OptionKeyRetryTimeout:
			if v, ok := option.OptionValue().(time.Duration); ok {
				res.RetryTimeout = v
			}
		case OptionKeyTotalTimeout:
			if v, ok := option.OptionValue().(time.Duration); ok {
				res.TotalTimeout = v
			}
		}
	}
	return res
}

// WithMaxIdleConns ...
func WithMaxIdleConns(maxIdleConns int) HTTPClientOption {
	return &httpClientOption{
		optionKey:   OptionKeyMaxIdleConns,
		optionValue: maxIdleConns,
	}
}

// WithIdleConnTimeout ...
func WithIdleConnTimeout(idleConnTimeout time.Duration) HTTPClientOption {
	return &httpClientOption{
		optionKey:   OptionKeyIdleConnTimeout,
		optionValue: idleConnTimeout,
	}
}

// WithDisableCompression ...
func WithDisableCompression(disableCompression bool) HTTPClientOption {
	return &httpClientOption{
		optionKey:   OptionKeyDisableCompression,
		optionValue: disableCompression,
	}
}

// WithRetryTimes ...
func WithRetryTimes(retryTimes int) HTTPClientOption {
	return &httpClientOption{
		optionKey:   OptionKeyRetryTimes,
		optionValue: retryTimes,
	}
}

// WithRetryTimeout ...
func WithRetryTimeout(timeout time.Duration) HTTPClientOption {
	return &httpClientOption{
		optionKey:   OptionKeyRetryTimeout,
		optionValue: timeout,
	}
}

// WithTotalTimeout ...
func WithTotalTimeout(timeout time.Duration) HTTPClientOption {
	return &httpClientOption{
		optionKey:   OptionKeyTotalTimeout,
		optionValue: timeout,
	}
}
