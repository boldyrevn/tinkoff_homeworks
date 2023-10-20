package app

import (
	"errors"
	"homework/model"
)

type Service interface {
	GetDevice(string) (model.Device, error)
	CreateDevice(model.Device) error
	DeleteDevice(string) error
	UpdateDevice(model.Device) error
}

type service struct {
	devices map[string]model.Device
}

func NewService() Service {
	return &service{
		devices: make(map[string]model.Device),
	}
}

func (s *service) GetDevice(num string) (model.Device, error) {
	d, ok := s.devices[num]
	if !ok {
		return d, errors.New("device is not found")
	}
	return d, nil
}

func (s *service) CreateDevice(d model.Device) error {
	_, ok := s.devices[d.SerialNum]
	if ok {
		return errors.New("device already exists")
	}
	s.devices[d.SerialNum] = model.Device{
		SerialNum: d.SerialNum,
		Model:     d.Model,
		IP:        d.IP,
	}
	return nil
}

func (s *service) DeleteDevice(num string) error {
	d, ok := s.devices[num]
	if !ok {
		return errors.New("device is not found")
	}
	delete(s.devices, d.SerialNum)
	return nil
}

func (s *service) UpdateDevice(d model.Device) error {
	_, ok := s.devices[d.SerialNum]
	if !ok {
		return errors.New("device is not found")
	}
	s.devices[d.SerialNum] = d
	return nil
}
