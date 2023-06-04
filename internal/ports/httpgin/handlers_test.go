package httpgin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"goads/internal/app"
	"goads/internal/app/mocks"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type test struct {
	name   string
	body   map[string]string
	method string
	url    string
	code   int
}

type adTest struct {
	test
	want adResponse
}

type userTest struct {
	test
	want userResponse
}

type GinTestSuite struct {
	suite.Suite
	adTests   []adTest
	userTests []userTest
	router    *gin.Engine
	ads       *mocks.Ads
	users     *mocks.Users
}

func (s *GinTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.New()
	SetRoutes(s.router.Group("/api"), s.ads, s.users)
}

func (s *GinTestSuite) runTest(t test) *httptest.ResponseRecorder {
	data, _ := json.Marshal(t.body)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(t.method, fmt.Sprintf("/api/%s", t.url), bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	s.router.ServeHTTP(rec, req)
	return rec
}

func (s *GinTestSuite) TestAds() {
	for _, tt := range s.adTests {
		s.T().Run(tt.name, func(t *testing.T) {
			rec := s.runTest(tt.test)
			assert.Equal(t, tt.code, rec.Code)
			if tt.code != 200 || tt.method == "DELETE" {
				return
			}
			respData, err := io.ReadAll(rec.Body)
			assert.NoError(t, err)

			var response struct {
				Data adResponse
			}
			err = json.Unmarshal(respData, &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, response.Data)
		})
	}
}

func (s *GinTestSuite) TestUsers() {
	for _, tt := range s.userTests {
		s.T().Run(tt.name, func(t *testing.T) {
			rec := s.runTest(tt.test)
			assert.Equal(t, tt.code, rec.Code)
			if tt.code != 200 || tt.method == "DELETE" {
				return
			}
			respData, err := io.ReadAll(rec.Body)
			assert.NoError(t, err)

			var response struct {
				Data userResponse
			}
			err = json.Unmarshal(respData, &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, response.Data)
		})
	}
}

func newSuite(am *mocks.Ads, um *mocks.Users, adTests []adTest, userTests []userTest) *GinTestSuite {
	s := new(GinTestSuite)
	s.userTests = userTests
	s.adTests = adTests
	s.users = um
	s.ads = am
	return s
}

func TestCreate(t *testing.T) {
	amCreate := mocks.NewAds(t)
	amCreate.
		On("Create", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(func(_ context.Context, title string, text string, authorID int64) (ads.Ad, error) {
			if title == "error" {
				return ads.Ad{}, fmt.Errorf("test error")
			}
			return ads.Ad{
				AuthorID: authorID,
				Title:    title,
				Text:     text,
			}, nil
		})

	umCreate := mocks.NewUsers(t)
	umCreate.
		On("Create", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(func(_ context.Context, email string, name string) (users.User, error) {
			if name == "error" {
				return users.User{}, fmt.Errorf("test error")
			}
			return users.User{
				Name:  name,
				Email: email,
			}, nil
		})

	suite.Run(t, newSuite(
		amCreate, umCreate,
		[]adTest{
			{
				test: test{
					name: "Valid request",
					body: map[string]string{
						"title": "123",
						"text":  " 123",
					},
					method: "POST",
					url:    "ads",
					code:   200,
				},
				want: adResponse{
					Title: "123",
					Text:  " 123",
				},
			},
			{
				test: test{
					name: "Creating error",
					body: map[string]string{
						"title": "error",
						"text":  " error",
					},
					method: "POST",
					url:    "ads",
					code:   500,
				},
			},
		},
		[]userTest{
			{
				test: test{
					method: "POST",
					url:    "users",
					name:   "Valid request",
					body: map[string]string{
						"name":  "123",
						"email": " 123",
					},
					code: 200,
				},
				want: userResponse{
					Name:  "123",
					Email: " 123",
				},
			},
			{
				test: test{
					method: "POST",
					url:    "users",
					name:   "Creating error",
					body: map[string]string{
						"name":  "error",
						"email": " error",
					},
					code: 500,
				},
			},
		},
	))
	suite.Run(t, newSuite(
		nil, nil,
		[]adTest{
			{
				test: test{
					name:   "Invalid binding",
					code:   400,
					method: "POST",
					url:    "ads",
				},
			},
		},
		[]userTest{
			{
				test: test{
					method: "POST",
					url:    "users",
					name:   "Invalid binding",
					code:   400,
				},
			},
		},
	))
}

func TestUpdate(t *testing.T) {
	amChangeGet := mocks.NewAds(t)
	amChange := mocks.NewAds(t)

	amChange.
		On("ChangeStatus", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(app.ErrNotFound)

	amChangeGet.
		On("ChangeStatus", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("bool")).
		Return(nil)
	amChangeGet.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(ads.Ad{}, func(_ context.Context, id int64) error {
			if id == 123 {
				return app.ErrPermissionDenied
			}
			return nil
		})

	amUpdateGet := mocks.NewAds(t)
	amUpdate := mocks.NewAds(t)
	amUpdate.
		On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(app.ErrNotFound)

	amUpdateGet.
		On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	amUpdateGet.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(ads.Ad{}, func(_ context.Context, id int64) error {
			if id == 123 {
				return app.ErrPermissionDenied
			}
			return nil
		})

	suite.Run(t, newSuite(amUpdateGet, nil, []adTest{
		{
			test: test{
				method: "PUT",
				url:    "ads/0",
				name:   "Valid request",
				body: map[string]string{
					"title": "1",
					"text":  "1",
				},
				code: 200,
			},
		},
		{
			test: test{
				method: "PUT",
				url:    "ads/123",
				name:   "Invalid get",
				body: map[string]string{
					"title": "1",
					"text":  "1",
				},
				code: 403,
			},
		},
	}, nil))

	suite.Run(t, newSuite(amUpdate, nil, []adTest{
		{
			test: test{
				method: "PUT",
				url:    "ads/123",
				name:   "Invalid change status",
				body: map[string]string{
					"title": "1",
					"text":  "1",
				},
				code: 404,
			},
		},
	}, nil))

	suite.Run(t, newSuite(amChangeGet, nil, []adTest{
		{
			test: test{
				name:   "Valid request",
				method: "PUT",
				url:    "ads/0/status",
				code:   200,
			},
		},
		{
			test: test{
				name:   "Invalid get",
				method: "PUT",
				url:    "ads/123/status",
				code:   403,
			},
		},
	}, nil))

	suite.Run(t, newSuite(amChange, nil, []adTest{
		{
			test: test{
				name:   "Invalid change status",
				method: "PUT",
				url:    "ads/123/status",
				code:   404,
			},
		},
	}, nil))

	suite.Run(t, newSuite(nil, nil, []adTest{
		{
			test: test{
				method: "PUT",
				url:    "ads/123",
				name:   "Invalid request",
				code:   400,
			},
		},
		{
			test: test{
				method: "PUT",
				url:    "ads/_",
				name:   "Invalid ad id",
				code:   400,
			},
		},
		{
			test: test{
				name:   "Invalid ad id",
				method: "PUT",
				url:    "ads/_/status",
				code:   400,
			},
		},
	}, nil))
}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	am := mocks.NewAds(t)

	am.
		On("Delete", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64, _ int64) error {
			if id == 123 {
				return app.ErrPermissionDenied
			}
			return nil
		})
	um := mocks.NewUsers(t)

	um.
		On("Delete", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) error {
			if id == 123 {
				return app.ErrPermissionDenied
			}
			return nil
		})

	suite.Run(t, newSuite(
		am, um, []adTest{
			{
				test: test{
					name:   "Valid request",
					method: "DELETE",
					url:    "ads/0",
					code:   200,
				},
			},
			{
				test: test{
					name:   "Invalid delete",
					code:   403,
					method: "DELETE",
					url:    "ads/123",
				},
			},
		},
		[]userTest{
			{
				test: test{
					name:   "Valid user delete",
					method: "DELETE",
					url:    "users/0",
					code:   200,
				},
			},
			{
				test: test{
					name:   "Invalid delete",
					method: "DELETE",
					url:    "users/123",
					code:   403,
				},
			},
		},
	))

	suite.Run(t, newSuite(
		nil, nil,
		[]adTest{
			{
				test: test{
					name:   "Invalid ad id",
					code:   400,
					method: "DELETE",
					url:    "ads/_",
				},
			},
		},
		[]userTest{
			{
				test: test{
					name:   "Invalid user id",
					code:   400,
					method: "DELETE",
					url:    "users/_",
				},
			},
		},
	))
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	umGet := mocks.NewUsers(t)

	umGet.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) (users.User, error) {
			if id == 123 {
				return users.User{}, fmt.Errorf("test error")
			}
			return users.User{
				ID: id,
			}, nil
		})

	suite.Run(t, newSuite(nil, umGet, nil, []userTest{
		{
			test: test{
				name:   "Valid request",
				code:   200,
				method: "GET",
				url:    "users/0",
			},
		},
		{
			test: test{
				code:   500,
				method: "GET",
				name:   "Getting error",
				url:    "users/123",
			},
		},
	}))

	suite.Run(t, newSuite(nil, nil, nil, []userTest{
		{
			test: test{
				name:   "Invalid user id",
				code:   400,
				method: "GET",
				url:    "users/_",
			},
		},
	}))
}
