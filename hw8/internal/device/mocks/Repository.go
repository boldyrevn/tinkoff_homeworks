// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
    mock "github.com/stretchr/testify/mock"
    "homework/internal/device"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: d
func (_m *Repository) Create(d device.Device) error {
	ret := _m.Called(d)

	var r0 error
	if rf, ok := ret.Get(0).(func(device.Device) error); ok {
		r0 = rf(d)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: num
func (_m *Repository) Delete(num string) error {
	ret := _m.Called(num)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(num)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: num
func (_m *Repository) Get(num string) (device.Device, error) {
	ret := _m.Called(num)

	var r0 device.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (device.Device, error)); ok {
		return rf(num)
	}
	if rf, ok := ret.Get(0).(func(string) device.Device); ok {
		r0 = rf(num)
	} else {
		r0 = ret.Get(0).(device.Device)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(num)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: d
func (_m *Repository) Update(d device.Device) error {
	ret := _m.Called(d)

	var r0 error
	if rf, ok := ret.Get(0).(func(device.Device) error); ok {
		r0 = rf(d)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
