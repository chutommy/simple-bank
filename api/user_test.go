package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chutified/simple-bank/api"
	"github.com/chutified/simple-bank/db/mocks"
	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/chutified/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestServer_CreateUser(t *testing.T) {
	user := db.User{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomOwner(),
		FirstName:      util.RandomOwner(),
		LastName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	tests := []struct {
		name          string
		param         api.CreateUserRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: api.CreateUserRequest{
				Username:  user.Username,
				Password:  user.HashedPassword,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateUser", mock.Anything, db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					FirstName:      user.FirstName,
					LastName:       user.LastName,
					Email:          user.Email,
				}).Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				u := bytesToUser(t, resp.Body.Bytes())
				assert.Equal(t, user, u)
			},
		},
		{
			name:      "InvalidRequest",
			param:     api.CreateUserRequest{},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "UniqueKeyViolation",
			param: api.CreateUserRequest{
				Username:  user.Username,
				Password:  user.HashedPassword,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateUser", mock.Anything, db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					FirstName:      user.FirstName,
					LastName:       user.LastName,
					Email:          user.Email,
				}).Return(db.User{}, pq.Error{
					Code:    "23505",
					Message: "unique_violation",
				})
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, resp.Code)
			},
		},
		{
			name: "InternalError",
			param: api.CreateUserRequest{
				Username:  user.Username,
				Password:  user.HashedPassword,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateUser", mock.Anything, db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					FirstName:      user.FirstName,
					LastName:       user.LastName,
					Email:          user.Email,
				}).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct a server with a mock db.Store
			store := new(mocks.Store)
			server := api.NewServer(store)
			test.buildStub(store)

			// construct request and response recorder
			url := "/users"
			b, err := json.Marshal(test.param)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			recorder := httptest.NewRecorder()

			// server
			server.Srv.Handler.ServeHTTP(recorder, req)

			// check response
			test.checkResponse(t, recorder)
			store.AssertExpectations(t)
		})
	}
}

func TestServer_GetUser(t *testing.T) {
	user := db.User{
		Username:           util.RandomOwner(),
		HashedPassword:     util.RandomOwner(),
		FirstName:          util.RandomOwner(),
		LastName:           util.RandomOwner(),
		Email:              util.RandomEmail(),
		PasswordModifiedAt: time.Now(),
		CreatedAt:          time.Now(),
	}

	tests := []struct {
		name          string
		param         api.GetUserRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: api.GetUserRequest{
				Username: user.Username,
			},
			buildStub: func(store *mocks.Store) {
				store.On("GetUser", mock.Anything, user.Username).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				u := bytesToUser(t, resp.Body.Bytes())
				assert.Equal(t, user, u)
			},
		},
		{
			name:      "InvalidRequest",
			param:     api.GetUserRequest{},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "NotFound",
			param: api.GetUserRequest{
				Username: user.Username,
			},
			buildStub: func(store *mocks.Store) {
				store.On("GetUser", mock.Anything, user.Username).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "InternalError",
			param: api.GetUserRequest{
				Username: user.Username,
			},
			buildStub: func(store *mocks.Store) {
				store.On("GetUser", mock.Anything, user.Username).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct server with mock db.Store
			mockStore := new(mocks.Store)
			test.buildStub(mockStore)
			server := api.NewServer(mockStore)

			// construct a request and response recorder
			url := "/users"
			b, err := json.Marshal(test.param)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodGet, url, bytes.NewReader(b))
			recorder := httptest.NewRecorder()

			// serve
			server.Srv.Handler.ServeHTTP(recorder, req)

			// check result
			test.checkResponse(t, recorder)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_UpdateUserPassword(t *testing.T) {
	user1 := db.User{
		Username:           util.RandomOwner(),
		HashedPassword:     util.RandomOwner(),
		FirstName:          util.RandomOwner(),
		LastName:           util.RandomOwner(),
		Email:              util.RandomEmail(),
		PasswordModifiedAt: time.Now().Add(-10 * time.Hour),
		CreatedAt:          time.Now().Add(-10 * time.Hour),
	}
	user2 := db.User{
		Username:           util.RandomOwner(),
		HashedPassword:     util.RandomOwner(),
		FirstName:          util.RandomOwner(),
		LastName:           util.RandomOwner(),
		Email:              util.RandomEmail(),
		PasswordModifiedAt: time.Now(),
		CreatedAt:          time.Now(),
	}

	tests := []struct {
		name          string
		params        api.UpdateUserPasswordRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: api.UpdateUserPasswordRequest{
				Username:     user1.Username,
				OldPassoword: user1.HashedPassword,
				NewPassword:  user2.HashedPassword,
			},
			buildStub: func(store *mocks.Store) {
				store.On("UpdateUserPassword", mock.Anything, db.UpdateUserPasswordParams{
					Username:         user1.Username,
					HashedPassword:   user1.HashedPassword,
					HashedPassword_2: user2.HashedPassword,
				}).
					Return(user2, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				u := bytesToUser(t, resp.Body.Bytes())
				assert.Equal(t, user2, u)
			},
		},
		{
			name:      "InvalidRequest",
			params:    api.UpdateUserPasswordRequest{},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "NotFound",
			params: api.UpdateUserPasswordRequest{
				Username:     user1.Username,
				OldPassoword: user1.HashedPassword,
				NewPassword:  user2.HashedPassword,
			},
			buildStub: func(store *mocks.Store) {
				store.On("UpdateUserPassword", mock.Anything, db.UpdateUserPasswordParams{
					Username:         user1.Username,
					HashedPassword:   user1.HashedPassword,
					HashedPassword_2: user2.HashedPassword,
				}).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "InternalError",
			params: api.UpdateUserPasswordRequest{
				Username:     user1.Username,
				OldPassoword: user1.HashedPassword,
				NewPassword:  user2.HashedPassword,
			},
			buildStub: func(store *mocks.Store) {
				store.On("UpdateUserPassword", mock.Anything, db.UpdateUserPasswordParams{
					Username:         user1.Username,
					HashedPassword:   user1.HashedPassword,
					HashedPassword_2: user2.HashedPassword,
				}).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct server with mock db.Store
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := "/users"
			b, err := json.Marshal(test.params)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(b))
			recorder := httptest.NewRecorder()

			// server
			server.Srv.Handler.ServeHTTP(recorder, req)

			// check response
			test.checkResponse(t, recorder)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_DeleteUser(t *testing.T) {
	user := db.User{
		Username:           util.RandomOwner(),
		HashedPassword:     util.RandomOwner(),
		FirstName:          util.RandomOwner(),
		LastName:           util.RandomOwner(),
		Email:              util.RandomEmail(),
		PasswordModifiedAt: time.Now(),
		CreatedAt:          time.Now(),
	}

	tests := []struct {
		name          string
		params        api.DeleteUserRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: api.DeleteUserRequest{
				Username: user.Username,
				Password: user.HashedPassword,
			},
			buildStub: func(store *mocks.Store) {
				store.On("DeleteUser", mock.Anything, db.DeleteUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
				}).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct a server with mock Store
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := "/users"
			b, err := json.Marshal(test.params)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodDelete, url, bytes.NewReader(b))
			recorder := httptest.NewRecorder()

			// server
			server.Srv.Handler.ServeHTTP(recorder, req)

			// check result
			test.checkResponse(t, recorder)
			mockStore.AssertExpectations(t)
		})
	}
}

func bytesToUser(t *testing.T, b []byte) db.User {
	var u db.User
	err := json.Unmarshal(b, &u)
	require.NoError(t, err)

	return u
}
