package maprepo

import (
    app "homework/internal/device"
    "testing"
)

func TestRepository_Create(t *testing.T) {
    repo := NewRepository()
    wantDevice := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(wantDevice)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    gotDevice, err := repo.Get(wantDevice.SerialNum)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if wantDevice != gotDevice {
        t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
    }
}

func TestRepository_CreateMultipleDevices(t *testing.T) {
    repo := NewRepository()
    devices := []app.Device{
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
        if err != nil {
            t.Errorf("unexpected error: %v", err)
        }
    }

    for _, wantDevice := range devices {
        gotDevice, err := repo.Get(wantDevice.SerialNum)
        if err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if wantDevice != gotDevice {
            t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
        }
    }
}

func TestRepository_CreateDuplicate(t *testing.T) {
    repo := NewRepository()
    wantDevice := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(wantDevice)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    err = repo.Create(wantDevice)
    if err == nil {
        t.Errorf("want error, but got nil")
    }

}

func TestRepository_GetUnexisting(t *testing.T) {
    repo := NewRepository()
    wantDevice := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(wantDevice)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    _, err = repo.Get("1")
    if err == nil {
        t.Error("want error, but got nil")
    }
}

func TestRepository_Delete(t *testing.T) {
    repo := NewRepository()
    newDevice := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(newDevice)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    err = repo.Delete(newDevice.SerialNum)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    _, err = repo.Get(newDevice.SerialNum)
    if err == nil {
        t.Error("want error, but got nil")
    }
}

func TestRepository_DeleteUnexisting(t *testing.T) {
    repo := NewRepository()
    err := repo.Delete("123")
    if err == nil {
        t.Errorf("want error, but got nil")
    }
}

func TestRepository_Update(t *testing.T) {
    repo := NewRepository()
    device := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(device)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    newDevice := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.2",
    }
    err = repo.Update(newDevice)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    gotDevice, err := repo.Get(newDevice.SerialNum)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if gotDevice != newDevice {
        t.Errorf("new device %+#v not equal got device %+#v", newDevice, gotDevice)
    }
}

func TestRepository_UpdateNonexistent(t *testing.T) {
    repo := NewRepository()
    device := app.Device{
        SerialNum: "123",
        Model:     "model1",
        IP:        "1.1.1.1",
    }

    err := repo.Create(device)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    newDevice := app.Device{
        SerialNum: "124",
        Model:     "model1",
        IP:        "1.1.1.2",
    }
    err = repo.Update(newDevice)
    if err == nil {
        t.Errorf("want err, but got nil")
    }
}

func BenchmarkRepository_Create(b *testing.B) {
    repo := NewRepository()
    for i := 0; i < b.N; i++ {
        _ = repo.Create(app.Device{
            SerialNum: "123",
            Model:     "Cisco",
            IP:        "192.54.0.54",
        })
    }
}
