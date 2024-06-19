package device_test

import (
    "github.com/stretchr/testify/assert"
    app "homework/internal/device"
    "homework/internal/model"
    "testing"
)

func TestRepository_Create(t *testing.T) {
    repo := app.NewRepository()
    wantDevice := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(wantDevice)
    assert.NoError(t, err)

    gotDevice, err := repo.Get(wantDevice.SerialNum)
    assert.NoError(t, err)
    assert.Equal(t, wantDevice, gotDevice)
}

func TestRepository_CreateMultipleDevices(t *testing.T) {
    repo := app.NewRepository()
    devices := []model.Device{
        {
            SerialNum: "123",
            Model:     "model1",
            IP:        "1.1.1.1",
        },
        {
            SerialNum: "124",
            Model:     "model2",
            IP:        "1.1.1.2",
        },
        {
            SerialNum: "125",
            Model:     "model3",
            IP:        "1.1.1.3",
        },
    }

    for _, d := range devices {
        err := repo.Create(d)
        assert.NoError(t, err)
    }

    for _, wantDevice := range devices {
        gotDevice, err := repo.Get(wantDevice.SerialNum)
        assert.NoError(t, err)
        assert.Equal(t, wantDevice, gotDevice)
    }
}

func TestRepository_CreateDuplicate(t *testing.T) {
    repo := app.NewRepository()
    wantDevice := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(wantDevice)
    assert.NoError(t, err)

    err = repo.Create(wantDevice)
    assert.Error(t, err)

}

func TestRepository_GetUnexisting(t *testing.T) {
    repo := app.NewRepository()
    wantDevice := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(wantDevice)
    assert.NoError(t, err)

    _, err = repo.Get("1")
    assert.Error(t, err)
}

func TestRepository_Delete(t *testing.T) {
    repo := app.NewRepository()
    newDevice := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(newDevice)
    assert.NoError(t, err)

    err = repo.Delete(newDevice.SerialNum)
    assert.NoError(t, err)

    _, err = repo.Get(newDevice.SerialNum)
    assert.Error(t, err)
}

func TestRepository_DeleteUnexisting(t *testing.T) {
    repo := app.NewRepository()
    err := repo.Delete("123")
    assert.Error(t, err)
}

func TestRepository_Update(t *testing.T) {
    repo := app.NewRepository()
    device := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(device)
    assert.NoError(t, err)

    newDevice := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.2",
    }
    err = repo.Update(newDevice)
    assert.NoError(t, err)

    gotDevice, err := repo.Get(newDevice.SerialNum)
    assert.NoError(t, err)
    assert.Equal(t, newDevice, gotDevice)
}

func TestRepository_UpdateNonexistent(t *testing.T) {
    repo := app.NewRepository()
    device := model.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(device)
    assert.NoError(t, err)

    newDevice := model.Device{
        SerialNum: "124",
        Model:     "model1",
        IP:        "1.1.1.2",
    }
    err = repo.Update(newDevice)
    assert.Error(t, err)
}

func BenchmarkRepository_Create(b *testing.B) {
    repo := app.NewRepository()
    for i := 0; i < b.N; i++ {
        _ = repo.Create(model.Device{
            SerialNum: "123",
            Model:     "Cisco",
            IP:        "192.54.0.54",
        })
    }
}
