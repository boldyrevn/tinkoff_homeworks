package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// WrapperRoundTripper adds functionality to other http.RoundTripper
type WrapperRoundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)

	// SetNext sets next RoundTripper that will be called.
	SetNext(http.RoundTripper)
}

// Config represents client configuration parameters
type Config struct {
	BreakerSettings  gobreaker.Settings `json:"breakerSettings"`
	Timeout          time.Duration      `json:"timeout,omitempty"`
	LoggingOn        bool               `json:"loggingOn,omitempty"`
	CircuitBreakerOn bool               `json:"circuitBreakerOn,omitempty"`
}

// NewJsonConfig makes client config from .json file. Timeouts
// must be specified in milliseconds
func NewJsonConfig(fileName string) (Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	var conf Config
	if err = json.NewDecoder(file).Decode(&conf); err != nil {
		return Config{}, err
	}
	conf.Timeout *= 1e6
	conf.BreakerSettings.Timeout *= 1e6
	conf.BreakerSettings.Interval *= 1e6
	return conf, nil
}

// NewClient returns new *http.Client. Timeout is specified in milliseconds
func NewClient(conf Config, baseTransport http.RoundTripper, userRoundTrippers ...WrapperRoundTripper) *http.Client {
	curTripper := baseTransport
	for i := len(userRoundTrippers) - 1; i >= 0; i-- {
		userRoundTrippers[i].SetNext(curTripper)
		curTripper = userRoundTrippers[i]
	}
	if conf.CircuitBreakerOn {
		b := &breakerRoundTripper{
			next:    curTripper,
			breaker: gobreaker.NewCircuitBreaker(conf.BreakerSettings),
		}
		curTripper = b
	}
	if conf.LoggingOn {
		l := &loggingRoundTripper{next: curTripper}
		curTripper = l
	}
	return &http.Client{
		Transport: curTripper,
		Timeout:   conf.Timeout,
	}
}

type loggingRoundTripper struct {
	next http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var logString string

	logString += fmt.Sprintf("%s %s \n", request.Method, request.URL.Path)
	logString += "-- Request headers --\n"
	for header, value := range request.Header {
		logString += fmt.Sprintf("%s: %s\n", header, strings.Join(value, ", "))
	}

	logString += "-- Query parameters --\n"
	for param, value := range request.URL.Query() {
		logString += fmt.Sprintf("%s = %s\n", param, strings.Join(value, " "))
	}

	resp, err := l.next.RoundTrip(request)

	logString += "-- Status --\n"
	if err != nil {
		logString += fmt.Sprintf("ERROR %s \n", err.Error())
		log.Print(logString)
		return resp, err
	}

	logString += fmt.Sprintf("OK %v \n", resp.StatusCode)
	logString += "-- Answer headers --\n"
	for header, value := range resp.Header {
		logString += fmt.Sprintf("%s: %s\n", header, strings.Join(value, ", "))
	}

	log.Print(logString)
	return resp, err
}

type breakerRoundTripper struct {
	next    http.RoundTripper
	breaker *gobreaker.CircuitBreaker
}

func (l *breakerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := l.breaker.Execute(func() (interface{}, error) {
		return l.next.RoundTrip(req)
	})
	if err != nil {
		return nil, err
	}
	ans, ok := resp.(*http.Response)
	if !ok {
		return nil, errors.New("one of the round trippers returned wrong response format")
	}
	return ans, nil
}
