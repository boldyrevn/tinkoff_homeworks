package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockAuth struct{}

func (ma MockAuth) AuthUser(name, password string) bool {
	return name == "monkey" && password == "256"
}

func TestNewJsonConfig(t *testing.T) {
	t.Run("wrong config config path", func(t *testing.T) {
		name := "test/seer_conf.json"
		_, err := NewJsonConfig(name)
		assert.Error(t, err)
	})
	t.Run("valid config", func(t *testing.T) {
		name := "../test/server_conf.json"
		c, err := NewJsonConfig(name)
		assert.NoError(t, err)
		assert.Equal(t, 2*time.Millisecond, c.IdleTimeout)
		assert.Equal(t, 2*time.Millisecond, c.WriteTimeout)
		assert.Equal(t, time.Millisecond, c.ReadTimeout)
		assert.Equal(t, time.Millisecond, c.ReadHeaderTimeout)
	})
	t.Run("wrong config format", func(t *testing.T) {
		name := "../test/server_wrong_conf.json"
		_, err := NewJsonConfig(name)
		assert.Error(t, err)
	})
}

func TestNewServer_logging(t *testing.T) {
	c := Config{
		Addr:      ":8080",
		RealmName: "test",
		LoggingOn: true,
	}

	m := http.NewServeMux()
	m.HandleFunc("/testOne", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Test answer header", "test value")
		_, _ = writer.Write([]byte("{}"))
	})
	m.HandleFunc("/testTwo", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(`{"message": "some error"}`))
	})
	s, _ := NewServer(m, c, nil)

	// log output must not be empty
	t.Run("successful request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/testOne?param=some", nil)
		req.Header.Set("Test header", "test value")
		s.Handler.ServeHTTP(recorder, req)
		assert.Equal(t, "{}", recorder.Body.String())
	})

	// log output must contain an error
	t.Run("bad request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/testTwo?param=bom", nil)
		req.Header.Set("Test header", "test value")
		s.Handler.ServeHTTP(recorder, req)
		assert.Equal(t, `{"message": "some error"}`, recorder.Body.String())
	})
}

func TestNewServer_basicAuth(t *testing.T) {
	c := Config{
		Addr:      ":8080",
		RealmName: "test",
		AuthOn:    true,
	}

	m := http.NewServeMux()
	m.HandleFunc("/auth", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("{}"))
	})

	t.Run("wrong basic auth parameters", func(t *testing.T) {
		_, err := NewServer(m, c, nil)
		assert.ErrorContains(t, err, "auth service is not provided")
		_, err = NewServer(m, Config{Addr: ":8080", AuthOn: true}, MockAuth{})
		assert.ErrorContains(t, err, "realm name must not be empty")
	})

	s, _ := NewServer(m, c, MockAuth{})

	t.Run("invalid credentials format", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/auth", nil)
		req.Header.Set("Authorization", "BasicumbW9ua2V5OjIyOA==")
		s.Handler.ServeHTTP(recorder, req)
		assert.Equal(t, 401, recorder.Code)
		assert.Equal(t, `Basic realm="test"`, recorder.Header().Get("WWW-Authenticate"))
	})

	t.Run("unauthorized user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/auth", nil)
		req.Header.Set("Authorization", "Basic bW9ua2V5OjI1NA==")
		s.Handler.ServeHTTP(recorder, req)
		assert.Equal(t, 401, recorder.Code)
		assert.Equal(t, `Basic realm="test"`, recorder.Header().Get("WWW-Authenticate"))
	})

	t.Run("successful request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:8080/auth", nil)
		req.Header.Set("Authorization", "Basic bW9ua2V5OjI1Ng==")
		s.Handler.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "", recorder.Header().Get("WWW-Authenticate"))
	})
}

func TestNewServer_userMiddleware(t *testing.T) {
	var customMwWorks bool
	fmt.Println(customMwWorks)
	myMw := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			customMwWorks = true
			handler.ServeHTTP(writer, request)
		})
	}

	c := Config{
		Addr: ":8080",
	}

	s, _ := NewServer(nil, c, nil, myMw)
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("{}"))
	})
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost:8080/test", nil)
	s.Handler.ServeHTTP(recorder, req)

	assert.Equal(t, true, customMwWorks)
}
