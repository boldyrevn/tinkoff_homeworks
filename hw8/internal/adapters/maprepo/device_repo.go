package maprepo

import (
    "errors"
    "homework/internal/device"
    "sync"
)

type repository struct {
    data map[string]device.Device
    mux  sync.RWMutex
}

func NewRepository() device.Repository {
    return &repository{data: make(map[string]device.Device)}
}

func (r *repository) Get(num string) (device.Device, error) {
    r.mux.RLock()
    defer r.mux.RUnlock()
    d, ok := r.data[num]
    if !ok {
        return d, errors.New("device is not found")
    }
    return d, nil
}

func (r *repository) Create(d device.Device) error {
    r.mux.Lock()
    defer r.mux.Unlock()
    if _, ok := r.data[d.SerialNum]; ok {
        return errors.New("device already exists")
    }
    r.data[d.SerialNum] = d
    return nil
}

func (r *repository) Delete(num string) error {
    r.mux.Lock()
    defer r.mux.Unlock()
    if _, ok := r.data[num]; !ok {
        return errors.New("device is not found")
    }
    delete(r.data, num)
    return nil
}

func (r *repository) Update(d device.Device) error {
    r.mux.Lock()
    defer r.mux.Unlock()
    if _, ok := r.data[d.SerialNum]; !ok {
        return errors.New("device is not found")
    }
    r.data[d.SerialNum] = d
    return nil
}
