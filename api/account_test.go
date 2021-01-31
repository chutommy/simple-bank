package api_test

import (
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
	account := randomAccount()

	tests := []struct {
		name          string
		accountID     int64
		buildStub     func(*mocks.Store)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStub: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, account, bytesToAccount(t, recorder.Body.Bytes()))
			},
		},
		{
			name:      "InvalidId",
			accountID: -1,
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStub: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
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
			url := fmt.Sprintf("/accounts/%d", test.accountID)
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

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 2048),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
		// CreatedAt: time.Now(),
	}
}

func bytesToAccount(t *testing.T, data []byte) db.Account {
	t.Helper()

	var a db.Account
	err := json.Unmarshal(data, &a)
	require.NoError(t, err)

	return a
}
