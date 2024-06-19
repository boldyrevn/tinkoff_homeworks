package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ClientSuit struct {
	suite.Suite
	Server *httptest.Server
}

func (c *ClientSuit) TestNewJsonConfig() {
	c.Run("wrong config path", func() {
		name := `asdfdfs/sadfsdf`
		_, err := NewJsonConfig(name)
		assert.Error(c.T(), err)
	})
	c.Run("valid config", func() {
		name := `../test/client_conf.json`
		conf, err := NewJsonConfig(name)
		assert.NoError(c.T(), err)
		assert.Equal(c.T(), 3*time.Second, conf.Timeout)
		assert.Equal(c.T(), 30*time.Second, conf.BreakerSettings.Timeout)
		assert.Equal(c.T(), 60*time.Second, conf.BreakerSettings.Interval)
		assert.Equal(c.T(), "inner", conf.BreakerSettings.Name)
	})
	c.Run("wrong config format", func() {
		name := `../test/client_wrong_conf.json`
		_, err := NewJsonConfig(name)
		assert.Error(c.T(), err)
	})
}

func (c *ClientSuit) TestNewClient_logging() {
	client := NewClient(Config{LoggingOn: true}, c.Server.Client().Transport)
	c.Run("successful request", func() {
		req, _ := http.NewRequest("GET", c.Server.URL+"/testOne?test=value", nil)
		req.Header.Add("Test-Header", "test-value")
		resp, err := client.Do(req)
		// logging output must not be empty
		if err != nil {
			c.T().Fatal(err)
		}
		b, _ := io.ReadAll(resp.Body)
		assert.Equal(c.T(), "{}", string(b))
	})
	c.Run("error request", func() {
		req, _ := http.NewRequest("GET", "http://localhost:8080/testTwo?test=value", nil)
		req.Header.Add("Test-Header", "test value")
		_, err := client.Do(req)
		// logging output must not be empty
		assert.Error(c.T(), err)
	})
}

func (c *ClientSuit) TestNewClient_breaker() {
	c.Run("wrong request", func() {
		client := NewClient(Config{CircuitBreakerOn: true}, c.Server.Client().Transport)
		for i := 0; i < 10; i++ {
			_, err := client.Get("http://localhost:8080/wrongPath")
			assert.Error(c.T(), err)
		}
		_, err := client.Get("http://localhost:8080/wrongPath")
		assert.ErrorContains(c.T(), err, "circuit breaker is open")
	})

	c.Run("successful request", func() {
		client := NewClient(Config{CircuitBreakerOn: true}, c.Server.Client().Transport)
		for i := 0; i < 10; i++ {
			_, err := client.Get(c.Server.URL + "/testOne")
			assert.NoError(c.T(), err)
		}
	})
}

type myTripper struct {
	next   http.RoundTripper
	number int
}

func (m *myTripper) SetNext(t http.RoundTripper) {
	m.next = t
}

func (m *myTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.number = 777
	return m.next.RoundTrip(req)
}

func (c *ClientSuit) TestNewClient_userTripper() {
	t := &myTripper{}
	client := NewClient(Config{}, http.DefaultTransport, t)
	resp, err := client.Get(c.Server.URL + "/testOne")
	b, _ := io.ReadAll(resp.Body)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), "{}", string(b))
	assert.Equal(c.T(), 777, t.number)
}

func (c *ClientSuit) SetupSuite() {
	http.HandleFunc("/testOne", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Test answer header", "test value")
		_, _ = writer.Write([]byte("{}"))
	})
	http.HandleFunc("/testTwo", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(`{"message": "some error"}`))
	})
	http.HandleFunc("/testThree", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(`{"message": "some error"}`))
	})
	s := httptest.NewServer(http.DefaultServeMux)
	c.Server = s
}

func (c *ClientSuit) TearDownSuite() {
	c.Server.Close()
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuit))
}
