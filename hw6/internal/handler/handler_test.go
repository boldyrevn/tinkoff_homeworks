package handler_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    "homework/internal/device/mocks"
    "homework/internal/handler"
    "homework/internal/model"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
)

var devs = [][]byte{
    []byte(`{ "serialNum": "1", "model": "model1", "IP": "8.8.8.8" }`), // successful
    []byte(`{ "serialNum": "2", "model": "model2", "IP": "8.8.8.8" }`), // service error
    []byte(`{ "serialNum": "3", ???, +"IP": "8.8.8.8" }`),              // wrong json format
}

var dev1 = model.Device{
    SerialNum: "1",
    Model:     "model1",
    IP:        "8.8.8.8",
}

var dev2 = model.Device{
    SerialNum: "2",
    Model:     "model2",
    IP:        "8.8.8.8",
}

func makeMessage(s string) []byte {
    res, _ := json.MarshalIndent(model.Message{Text: s}, "", "  ")
    return res
}

type HandlerSuite struct {
    suite.Suite
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
}

func (suite *HandlerSuite) TearDownSuite() {
}

func (suite *HandlerSuite) TestHandler_GetDevice() {
    getTests := []Test{
        {
            req:    httptest.NewRequest("GET", "/test?num=1", nil),
            expect: []byte("{}"),
        },
        {
            req:    httptest.NewRequest("GET", "/test?num=2", nil),
            expect: makeMessage("some error"),
        },
        {
            req:    httptest.NewRequest("GET", "/test", nil),
            expect: makeMessage("query must contain `num` parameter"),
        },
    }

    uc := new(mocks.UseCase)
    uc.EXPECT().GetDevice("1").Return(model.Device{}, nil)
    uc.EXPECT().GetDevice("2").Return(model.Device{}, errors.New("some error"))
    h := handler.Handler{UseCase: uc}
    runTests(suite, getTests, h.GetDevice)
}

func (suite *HandlerSuite) TestHandler_DeleteDevice() {
    deleteTests := []Test{
        {
            req:    httptest.NewRequest("DELETE", "/test?num=1", nil),
            expect: makeMessage("device was deleted successfully"),
        },
        {
            req:    httptest.NewRequest("DELETE", "/test?num=2", nil),
            expect: makeMessage("some error"),
        },
        {
            req:    httptest.NewRequest("DELETE", "/test", nil),
            expect: makeMessage("query must contain `num` parameter"),
        },
    }

    uc := new(mocks.UseCase)
    uc.EXPECT().DeleteDevice("1").Return(nil)
    uc.EXPECT().DeleteDevice("2").Return(errors.New("some error"))
    h := handler.Handler{UseCase: uc}
    runTests(suite, deleteTests, h.DeleteDevice)
}

func (suite *HandlerSuite) TestHandler_CreateDevice() {
    createTests := []Test{
        {
            req:    httptest.NewRequest("POST", "/test", bytes.NewReader(devs[0])),
            expect: makeMessage("device was created successfully"),
        },
        {
            req:    httptest.NewRequest("POST", "/test", bytes.NewReader(devs[1])),
            expect: makeMessage("some error"),
        },
        {
            req:    httptest.NewRequest("POST", "/test", bytes.NewReader(devs[2])),
            expect: makeMessage("wrong json format"),
        },
    }
    uc := new(mocks.UseCase)
    uc.EXPECT().CreateDevice(dev1).Return(nil)
    uc.EXPECT().CreateDevice(dev2).Return(errors.New("some error"))
    h := handler.Handler{UseCase: uc}
    runTests(suite, createTests, h.CreateDevice)
}

func (suite *HandlerSuite) TestHandler_UpdateDevice() {
    updateTests := []Test{
        {
            req:    httptest.NewRequest("PUT", "/test", bytes.NewReader(devs[0])),
            expect: makeMessage("device was updated successfully"),
        },
        {
            req:    httptest.NewRequest("PUT", "/test", bytes.NewReader(devs[1])),
            expect: makeMessage("some error"),
        },
        {
            req:    httptest.NewRequest("PUT", "/test", bytes.NewReader(devs[2])),
            expect: makeMessage("wrong json format"),
        },
    }
    uc := new(mocks.UseCase)
    uc.EXPECT().UpdateDevice(dev1).Return(nil)
    uc.EXPECT().UpdateDevice(dev2).Return(errors.New("some error"))
    h := handler.Handler{UseCase: uc}
    runTests(suite, updateTests, h.UpdateDevice)
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
    uc.On("GetDevice", mock.AnythingOfType("string")).Return(model.Device{}, nil)
    uc.On("DeleteDevice", mock.AnythingOfType("string")).Return(nil)
    uc.On("CreateDevice", mock.AnythingOfType("model.Device")).Return(nil)
    uc.On("UpdateDevice", mock.AnythingOfType("model.Device")).Return(nil)
    h := handler.Handler{UseCase: uc}

    f.Add([]byte("{some json}"))
    f.Fuzz(func(t *testing.T, b []byte) {
        h.CreateDevice(makeReqRec(b))
        h.DeleteDevice(makeReqRec(b))
        h.GetDevice(makeReqRec(b))
        h.UpdateDevice(makeReqRec(b))
    })
}
