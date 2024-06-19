package app

import (
    "homework/internal/device"
    "homework/pkg/utils"
)

type UseCase interface {
    GetDevice(string) (device.Device, error)
    CreateDevice(device.Device) error
    DeleteDevice(string) error
    UpdateDevice(device.Device) error
}

type service struct {
    repo device.Repository
}

func NewUseCase(r device.Repository) UseCase {
    return &service{
        repo: r,
    }
}

func (s *service) GetDevice(num string) (device.Device, error) {
    return s.repo.Get(num)
}

func (s *service) CreateDevice(d device.Device) error {
    if err := ValidateDevice(d); err != nil {
        return err
    }
    return s.repo.Create(d)
}

func (s *service) DeleteDevice(num string) error {
    return s.repo.Delete(num)
}

func (s *service) UpdateDevice(d device.Device) error {
    err := ValidateDevice(d)
    if err != nil {
        return err
    }
    return s.repo.Update(d)
}

func ValidateDevice(d device.Device) error {
    if d.IP == "" || d.SerialNum == "" || d.Model == "" {
        return utils.UnfilledError
    } else if err := utils.ValidateIP(d.IP); err != nil {
        return err
    }
    return nil
}
