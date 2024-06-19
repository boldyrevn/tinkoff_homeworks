package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AuthService interface {
	AuthUser(name, password string) bool
}

type Middleware func(http.Handler) http.Handler

// Config represents server configuration parameters.
type Config struct {
	Addr              string        `json:"addr,omitempty"`
	RealmName         string        `json:"realmName,omitempty"`
	ReadTimeout       time.Duration `json:"readTimeout,omitempty"`
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout,omitempty"`
	WriteTimeout      time.Duration `json:"writeTimeout,omitempty"`
	IdleTimeout       time.Duration `json:"idleTimeout,omitempty"`
	AuthOn            bool          `json:"authOn,omitempty"`
	LoggingOn         bool          `json:"loggingOn,omitempty"`
}

// NewJsonConfig makes server config from .json file. Timeouts
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
	conf.ReadTimeout *= 1e6
	conf.ReadHeaderTimeout *= 1e6
	conf.WriteTimeout *= 1e6
	conf.IdleTimeout *= 1e6
	return conf, nil
}

type loggingWriter struct {
	http.ResponseWriter
	status     int
	errMessage string
	isWritten  bool
}

func (lw *loggingWriter) WriteHeader(status int) {
	if lw.isWritten {
		return
	}
	lw.isWritten = true
	lw.ResponseWriter.WriteHeader(status)
	lw.status = status
}

func (lw *loggingWriter) Write(b []byte) (int, error) {
	if !lw.isWritten {
		lw.WriteHeader(http.StatusOK)
	}
	res, err := lw.ResponseWriter.Write(b)
	if lw.status >= 400 {
		lw.errMessage = string(b)
	}
	return res, err
}

// NewServer creates new instance of Server with logging/auth/user middleware.
// Returns an error if Basic auth is on and AuthService is not provided. Uses
// http.DefaultServeMux if mux is nil
func NewServer(mux http.Handler, conf Config, as AuthService, userMw ...Middleware) (*http.Server, error) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	for _, mw := range userMw {
		mux = mw(mux)
	}
	if conf.AuthOn {
		if as == nil {
			return &http.Server{}, errors.New("auth service is not provided")
		} else if conf.RealmName == "" {
			return &http.Server{}, errors.New("realm name must not be empty")
		}
		mux = basicAuthMiddleware(mux, conf.RealmName, as)
	}
	if conf.LoggingOn {
		mux = loggingMiddleware(mux)
	}
	return &http.Server{
		Addr:              conf.Addr,
		Handler:           mux,
		ReadTimeout:       conf.ReadTimeout,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
	}, nil
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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

		logWriter := &loggingWriter{ResponseWriter: writer}
		start := time.Now()
		h.ServeHTTP(logWriter, request)
		end := time.Since(start).String()

		logString += "-- Status --\n"
		if logWriter.status >= 400 {
			logString += fmt.Sprintf("ERROR %v %v %s \n", logWriter.status, logWriter.errMessage, end)
		} else {
			logString += fmt.Sprintf("OK %v %s \n", logWriter.status, end)
		}

		logString += "-- Answer headers --\n"
		for header, value := range writer.Header() {
			logString += fmt.Sprintf("%s: %s\n", header, strings.Join(value, ", "))
		}
		log.Print(logString)
	})
}

func basicAuthMiddleware(h http.Handler, realmName string, as AuthService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, password, ok := request.BasicAuth()
		if !ok {
			writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realmName))
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ok = as.AuthUser(user, password)
		if !ok {
			writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realmName))
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(writer, request)
	})
}
