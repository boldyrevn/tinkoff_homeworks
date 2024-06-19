package device

import (
    "errors"
    "homework/internal/model"
    "sync"
)

type Repository interface {
    Get(num string) (model.Device, error)
    Create(d model.Device) error
    Delete(num string) error
    Update(d model.Device) error
}

type repository struct {
    mux  sync.RWMutex
    data map[string]model.Device
}

func NewRepository() Repository {
    return &repository{data: make(map[string]model.Device)}
}

func (m *repository) Get(num string) (model.Device, error) {
    m.mux.RLock()
    defer m.mux.RUnlock()
    d, ok := m.data[num]
    if !ok {
        return d, errors.New("device is not found")
    }
    return d, nil
}

func (m *repository) Create(d model.Device) error {
    m.mux.Lock()
    defer m.mux.Unlock()
    if _, ok := m.data[d.SerialNum]; ok {
        return errors.New("device already exists")
    }
    m.data[d.SerialNum] = d
    return nil
}

func (m *repository) Delete(num string) error {
    m.mux.Lock()
    defer m.mux.Unlock()
    if _, ok := m.data[num]; !ok {
        return errors.New("device is not found")
    }
    delete(m.data, num)
    return nil
}

func (m *repository) Update(d model.Device) error {
    m.mux.Lock()
    defer m.mux.Unlock()
    if _, ok := m.data[d.SerialNum]; !ok {
        return errors.New("device is not found")
    }
    m.data[d.SerialNum] = d
    return nil
}
