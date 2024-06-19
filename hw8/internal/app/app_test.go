package app

import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "homework/internal/device"
    "homework/internal/device/mocks"
    "homework/pkg/utils"
    "testing"
)

type Test struct {
    Name   string
    Dev    device.Device
    Expect error
}

func TestService_GetDevice(t *testing.T) {
    d := device.Device{
        SerialNum: "123X",
        Model:     "Abakus 3000",
        IP:        "123.45.228.23",
    }
    repo := &mocks.Repository{}
    repo.On("Get", "123X").Return(d, nil)

    uc := NewUseCase(repo)
    res, err := uc.GetDevice(d.SerialNum)
    assert.Equal(t, d, res)
    assert.Equal(t, err, nil)
}

func TestService_DeleteDevice(t *testing.T) {
    repo := &mocks.Repository{}
    repo.On("Delete", mock.AnythingOfType("string")).Return(nil)

    uc := NewUseCase(repo)
    err := uc.DeleteDevice("something")
    assert.Equal(t, nil, err)
}

func TestService_CreateDevice(t *testing.T) {
    tests := []Test{
        {
            "Right device format",
            device.Device{SerialNum: "58.12", Model: "Cisco", IP: "123.05.73.12"},
            nil,
        },
        {
            "IP with too big number",
            device.Device{SerialNum: "234F", Model: "Huawei", IP: "666.23.1.3"},
            utils.IncorrectOctetNumber,
        },
        {
            "IP with negative number",
            device.Device{SerialNum: "124", Model: "Bobus", IP: "172.-12.55.30"},
            utils.IncorrectOctetNumber,
        },
        {
            "Just wrong IP",
            device.Device{SerialNum: "ZXC", Model: "Breakdown", IP: "hehe.23bombom"},
            utils.OctetsCountError,
        },
        {
            "Empty fields",
            device.Device{SerialNum: "", Model: "Bamborghini", IP: "192.168.0.1"},
            utils.UnfilledError,
        },
        {
            "Too much IP octets",
            device.Device{SerialNum: "23-XF", Model: "HP Thinkcentre", IP: "23.234.15.1.56"},
            utils.OctetsCountError,
        },
    }

    repo := &mocks.Repository{}
    repo.On("Create", mock.AnythingOfType("device.Device")).Return(nil)

    uc := NewUseCase(repo)

    for _, test := range tests {
        test := test
        t.Run(test.Name, func(t *testing.T) {
            t.Parallel()
            err := uc.CreateDevice(test.Dev)
            assert.ErrorIs(t, err, test.Expect)
        })
    }
}

func TestService_UpdateDevice(t *testing.T) {
    tests := []Test{
        {
            "Successful update",
            device.Device{SerialNum: "92X", Model: "Asus 26", IP: "192.168.0.2"},
            nil,
        },
        {
            "Wrong device data",
            device.Device{SerialNum: "", Model: "", IP: "65.0.0.99"},
            utils.UnfilledError,
        },
    }

    repo := &mocks.Repository{}
    repo.On("Update", mock.AnythingOfType("device.Device")).Return(nil)

    uc := NewUseCase(repo)
    for _, test := range tests {
        test := test
        t.Run(test.Name, func(t *testing.T) {
            t.Parallel()
            err := uc.UpdateDevice(test.Dev)
            assert.ErrorIs(t, err, test.Expect)
        })
    }
}

func BenchmarkValidateDevice(b *testing.B) {
    d := device.Device{
        SerialNum: "57X",
        Model:     "Huawei",
        IP:        "78.235.9.23",
    }
    for i := 0; i < b.N; i++ {
        _ = ValidateDevice(d)
    }
}
