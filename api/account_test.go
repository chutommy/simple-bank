package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chutified/simple-bank/api"
	"github.com/chutified/simple-bank/db/mocks"
	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/chutified/simple-bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_GetAccountByID(t *testing.T) {
	account := db.Account{
		ID:       util.RandomInt(1, 2048),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
		// CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		apiRequest    api.GetAccountByIDRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			apiRequest: api.GetAccountByIDRequest{ID: account.ID},
			buildStub: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, account, bytesToAccount(t, recorder.Body))
			},
		},
		{
			name:       "InvalidId",
			apiRequest: api.GetAccountByIDRequest{ID: 0},
			buildStub:  func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			apiRequest: api.GetAccountByIDRequest{ID: account.ID},
			buildStub: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			apiRequest: api.GetAccountByIDRequest{ID: account.ID},
			buildStub: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			url := fmt.Sprintf("/accounts/%d", test.apiRequest.ID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			// serve
			server.Srv.Handler.ServeHTTP(recorder, req)

			// check result
			test.checkResponse(t, recorder)
			mockStore.AssertExpectations(t)
		})
	}
}

func bytesToAccount(t *testing.T, data *bytes.Buffer) db.Account {
	t.Helper()

	var a db.Account
	err := json.Unmarshal(data.Bytes(), &a)
	require.NoError(t, err)

	return a
}
