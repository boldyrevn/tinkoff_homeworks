package device

import (
    "errors"
    "homework/internal/model"
    "strconv"
    "strings"
)

var OctetsCountError = errors.New("IP address must have only 4 octets")
var OctetNumberError = errors.New("octet number must be an integer from 0 to 255")
var UnfilledError = errors.New("all fields must be filled")

func CheckIP(addr string) error {
    parts := strings.Split(addr, ".")
    if len(parts) != 4 {
        return OctetsCountError
    }
    for _, part := range parts {
        n, err := strconv.Atoi(part)
        if err != nil || !(0 <= n && n <= 255) {
            return OctetNumberError
        }
    }
    return nil
}

func ValidateDevice(d model.Device) error {
    if d.IP == "" || d.SerialNum == "" || d.Model == "" {
        return UnfilledError
    } else if err := CheckIP(d.IP); err != nil {
        return err
    }
    return nil
}

type UseCase interface {
    GetDevice(string) (model.Device, error)
    CreateDevice(model.Device) error
    DeleteDevice(string) error
    UpdateDevice(model.Device) error
}

type service struct {
    repo Repository
}

func NewUseCase(r Repository) UseCase {
    return &service{
        repo: r,
    }
}

func (s *service) GetDevice(num string) (model.Device, error) {
    return s.repo.Get(num)
}

func (s *service) CreateDevice(d model.Device) error {
    if err := ValidateDevice(d); err != nil {
        return err
    }
    return s.repo.Create(d)
}

func (s *service) DeleteDevice(num string) error {
    return s.repo.Delete(num)
}

func (s *service) UpdateDevice(d model.Device) error {
    err := ValidateDevice(d)
    if err != nil {
        return err
    }
    return s.repo.Update(d)
}
