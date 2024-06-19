// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	model "homework/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

type UseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *UseCase) EXPECT() *UseCase_Expecter {
	return &UseCase_Expecter{mock: &_m.Mock}
}

// CreateDevice provides a mock function with given fields: _a0
func (_m *UseCase) CreateDevice(_a0 model.Device) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Device) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UseCase_CreateDevice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateDevice'
type UseCase_CreateDevice_Call struct {
	*mock.Call
}

// CreateDevice is a helper method to define mock.On call
//   - _a0 model.Device
func (_e *UseCase_Expecter) CreateDevice(_a0 interface{}) *UseCase_CreateDevice_Call {
	return &UseCase_CreateDevice_Call{Call: _e.mock.On("CreateDevice", _a0)}
}

func (_c *UseCase_CreateDevice_Call) Run(run func(_a0 model.Device)) *UseCase_CreateDevice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.Device))
	})
	return _c
}

func (_c *UseCase_CreateDevice_Call) Return(_a0 error) *UseCase_CreateDevice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UseCase_CreateDevice_Call) RunAndReturn(run func(model.Device) error) *UseCase_CreateDevice_Call {
	_c.Call.Return(run)
	return _c
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

// UseCase_DeleteDevice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteDevice'
type UseCase_DeleteDevice_Call struct {
	*mock.Call
}

// DeleteDevice is a helper method to define mock.On call
//   - _a0 string
func (_e *UseCase_Expecter) DeleteDevice(_a0 interface{}) *UseCase_DeleteDevice_Call {
	return &UseCase_DeleteDevice_Call{Call: _e.mock.On("DeleteDevice", _a0)}
}

func (_c *UseCase_DeleteDevice_Call) Run(run func(_a0 string)) *UseCase_DeleteDevice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UseCase_DeleteDevice_Call) Return(_a0 error) *UseCase_DeleteDevice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UseCase_DeleteDevice_Call) RunAndReturn(run func(string) error) *UseCase_DeleteDevice_Call {
	_c.Call.Return(run)
	return _c
}

// GetDevice provides a mock function with given fields: _a0
func (_m *UseCase) GetDevice(_a0 string) (model.Device, error) {
	ret := _m.Called(_a0)

	var r0 model.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (model.Device, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) model.Device); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(model.Device)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UseCase_GetDevice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDevice'
type UseCase_GetDevice_Call struct {
	*mock.Call
}

// GetDevice is a helper method to define mock.On call
//   - _a0 string
func (_e *UseCase_Expecter) GetDevice(_a0 interface{}) *UseCase_GetDevice_Call {
	return &UseCase_GetDevice_Call{Call: _e.mock.On("GetDevice", _a0)}
}

func (_c *UseCase_GetDevice_Call) Run(run func(_a0 string)) *UseCase_GetDevice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UseCase_GetDevice_Call) Return(_a0 model.Device, _a1 error) *UseCase_GetDevice_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UseCase_GetDevice_Call) RunAndReturn(run func(string) (model.Device, error)) *UseCase_GetDevice_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDevice provides a mock function with given fields: _a0
func (_m *UseCase) UpdateDevice(_a0 model.Device) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Device) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UseCase_UpdateDevice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDevice'
type UseCase_UpdateDevice_Call struct {
	*mock.Call
}

// UpdateDevice is a helper method to define mock.On call
//   - _a0 model.Device
func (_e *UseCase_Expecter) UpdateDevice(_a0 interface{}) *UseCase_UpdateDevice_Call {
	return &UseCase_UpdateDevice_Call{Call: _e.mock.On("UpdateDevice", _a0)}
}

func (_c *UseCase_UpdateDevice_Call) Run(run func(_a0 model.Device)) *UseCase_UpdateDevice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.Device))
	})
	return _c
}

func (_c *UseCase_UpdateDevice_Call) Return(_a0 error) *UseCase_UpdateDevice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UseCase_UpdateDevice_Call) RunAndReturn(run func(model.Device) error) *UseCase_UpdateDevice_Call {
	_c.Call.Return(run)
	return _c
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