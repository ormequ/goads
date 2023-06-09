// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	ads "goads/internal/entities/ads"

	context "context"

	filters "goads/internal/filters"

	mock "github.com/stretchr/testify/mock"
)

// Ads is an autogenerated mock type for the Ads type
type Ads struct {
	mock.Mock
}

// ChangeStatus provides a mock function with given fields: ctx, id, userID, published
func (_m *Ads) ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error {
	ret := _m.Called(ctx, id, userID, published)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, bool) error); ok {
		r0 = rf(ctx, id, userID, published)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, title, text, authorID
func (_m *Ads) Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error) {
	ret := _m.Called(ctx, title, text, authorID)

	var r0 ads.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64) (ads.Ad, error)); ok {
		return rf(ctx, title, text, authorID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64) ads.Ad); ok {
		r0 = rf(ctx, title, text, authorID)
	} else {
		r0 = ret.Get(0).(ads.Ad)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int64) error); ok {
		r1 = rf(ctx, title, text, authorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id, userID
func (_m *Ads) Delete(ctx context.Context, id int64, userID int64) error {
	ret := _m.Called(ctx, id, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, id, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Ads) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
	ret := _m.Called(ctx, id)

	var r0 ads.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (ads.Ad, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) ads.Ad); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(ads.Ad)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFiltered provides a mock function with given fields: ctx, filter
func (_m *Ads) GetFiltered(ctx context.Context, filter filters.AdsOptions) ([]ads.Ad, error) {
	ret := _m.Called(ctx, filter)

	var r0 []ads.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, filters.AdsOptions) ([]ads.Ad, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, filters.AdsOptions) []ads.Ad); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ads.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, filters.AdsOptions) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: ctx, title
func (_m *Ads) Search(ctx context.Context, title string) ([]ads.Ad, error) {
	ret := _m.Called(ctx, title)

	var r0 []ads.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]ads.Ad, error)); ok {
		return rf(ctx, title)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []ads.Ad); ok {
		r0 = rf(ctx, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ads.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, userID, title, text
func (_m *Ads) Update(ctx context.Context, id int64, userID int64, title string, text string) error {
	ret := _m.Called(ctx, id, userID, title, text)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string, string) error); ok {
		r0 = rf(ctx, id, userID, title, text)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAds interface {
	mock.TestingT
	Cleanup(func())
}

// NewAds creates a new instance of Ads. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAds(t mockConstructorTestingTNewAds) *Ads {
	mock := &Ads{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
