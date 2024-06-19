// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
    mock "github.com/stretchr/testify/mock"
    "homework/internal/device"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// CreateDevice provides a mock function with given fields: _a0
func (_m *UseCase) CreateDevice(_a0 device.Device) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(device.Device) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDevice provides a mock function with given fields: _a0
func (_m *UseCase) DeleteDevice(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDevice provides a mock function with given fields: _a0
func (_m *UseCase) GetDevice(_a0 string) (device.Device, error) {
	ret := _m.Called(_a0)

	var r0 device.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (device.Device, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) device.Device); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(device.Device)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDevice provides a mock function with given fields: _a0
func (_m *UseCase) UpdateDevice(_a0 device.Device) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(device.Device) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
