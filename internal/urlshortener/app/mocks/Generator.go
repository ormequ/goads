// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// Generate provides a mock function with given fields: ctx
func (_m *Generator) Generate(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGenerator interface {
	mock.TestingT
	Cleanup(func())
}

// NewGenerator creates a new instance of Generator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGenerator(t mockConstructorTestingTNewGenerator) *Generator {
	mock := &Generator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}