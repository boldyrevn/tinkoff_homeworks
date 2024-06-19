package httpvanilla

import (
    "bytes"
    "encoding/json"
    "errors"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    "homework/internal/app/mocks"
    "homework/internal/device"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
)

func makeMessage(s string) []byte {
    res, _ := json.Marshal(Message{Text: s})
    return res
}

var jsonError = "wrong json format: invalid character '?' looking for beginning of object key string"

type HandlerSuite struct {
    suite.Suite
    paramRequests []*http.Request
    bodyRequests  []*http.Request
    uc            *mocks.UseCase
    devs          [][]byte
}

type Test struct {
    req    *http.Request
    expect []byte
}

func runTests(suite *HandlerSuite, tests []Test, hf func(http.ResponseWriter, *http.Request)) {
    for _, test := range tests {
        rec := httptest.NewRecorder()
        hf(rec, test.req)
        ans, _ := io.ReadAll(rec.Result().Body)
        suite.Equal(string(test.expect), string(ans))
    }
}

func (suite *HandlerSuite) SetupSuite() {
    suite.paramRequests = []*http.Request{
        httptest.NewRequest("GET", "/test?num=1", nil), // successful
        httptest.NewRequest("GET", "/test?num=2", nil), // service error
        httptest.NewRequest("GET", "/test", nil),       // without params
    }
    suite.devs = [][]byte{
        []byte(`{ "serialNum": "1", "model": "model1", "IP": "8.8.8.8" }`), // successful
        []byte(`{ "serialNum": "2", "model": "model2", "IP": "8.8.8.8" }`), // service error
        []byte(`{ "serialNum": "3", ???, +"IP": "8.8.8.8" }`),              // wrong json format
    }
    suite.uc = &mocks.UseCase{}
    suite.uc.On("GetDevice", "1").Return(device.Device{}, nil)
    suite.uc.On("GetDevice", "2").Return(device.Device{}, errors.New("some error"))
    suite.uc.On("DeleteDevice", "1").Return(nil)
    suite.uc.On("DeleteDevice", "2").Return(errors.New("some error"))

    var dev1 device.Device
    var dev2 device.Device
    _ = json.Unmarshal(suite.devs[0], &dev1)
    _ = json.Unmarshal(suite.devs[1], &dev2)
    for _, act := range []string{"Create", "Update"} {
        suite.uc.On(act+"Device", dev1).Return(nil)
        suite.uc.On(act+"Device", dev2).Return(errors.New("some error"))
    }
}

func (suite *HandlerSuite) TearDownSuite() {
    suite.paramRequests = nil
    suite.devs = nil
    suite.uc = nil
}

func (suite *HandlerSuite) SetupTest() {
    suite.bodyRequests = []*http.Request{
        httptest.NewRequest("POST", "/test", bytes.NewReader(suite.devs[0])),
        httptest.NewRequest("POST", "/test", bytes.NewReader(suite.devs[1])),
        httptest.NewRequest("POST", "/test", bytes.NewReader(suite.devs[2])),
    }
}

func (suite *HandlerSuite) TestHandler_GetDevice() {
    getTests := []Test{
        {req: suite.paramRequests[0], expect: []byte("{}")},
        {req: suite.paramRequests[1], expect: makeMessage("some error")},
        {req: suite.paramRequests[2], expect: makeMessage("query must contain `num` parameter")},
    }
    runTests(suite, getTests, GetDevice(suite.uc))
}

func (suite *HandlerSuite) TestHandler_DeleteDevice() {
    deleteTests := []Test{
        {req: suite.paramRequests[0], expect: makeMessage("")},
        {req: suite.paramRequests[1], expect: makeMessage("some error")},
        {req: suite.paramRequests[2], expect: makeMessage("query must contain `num` parameter")},
    }
    runTests(suite, deleteTests, DeleteDevice(suite.uc))
}

func (suite *HandlerSuite) TestHandler_CreateDevice() {
    createTests := []Test{
        {req: suite.bodyRequests[0], expect: makeMessage("")},
        {req: suite.bodyRequests[1], expect: makeMessage("some error")},
        {req: suite.bodyRequests[2], expect: makeMessage(jsonError)},
    }
    runTests(suite, createTests, CreateDevice(suite.uc))
}

func (suite *HandlerSuite) TestHandler_UpdateDevice() {
    updateTests := []Test{
        {req: suite.bodyRequests[0], expect: makeMessage("")},
        {req: suite.bodyRequests[1], expect: makeMessage("some error")},
        {req: suite.bodyRequests[2], expect: makeMessage(jsonError)},
    }
    runTests(suite, updateTests, UpdateDevice(suite.uc))
}

func TestHandler(t *testing.T) {
    suite.Run(t, new(HandlerSuite))
}

func makeReqRec(b []byte) (http.ResponseWriter, *http.Request) {
    rec := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/test", bytes.NewReader(b))
    return rec, req
}

func FuzzCreateDevice(f *testing.F) {
    uc := &mocks.UseCase{}
    uc.On("GetDevice", mock.AnythingOfType("string")).Return(device.Device{}, nil)
    uc.On("DeleteDevice", mock.AnythingOfType("string")).Return(nil)
    uc.On("CreateDevice", mock.AnythingOfType("model.Device")).Return(nil)
    uc.On("UpdateDevice", mock.AnythingOfType("model.Device")).Return(nil)

    f.Add([]byte("{some json}"))
    f.Fuzz(func(t *testing.T, b []byte) {
        CreateDevice(uc)(makeReqRec(b))
        CreateDevice(uc)(makeReqRec(b))
        DeleteDevice(uc)(makeReqRec(b))
        GetDevice(uc)(makeReqRec(b))
        UpdateDevice(uc)(makeReqRec(b))
    })
}
