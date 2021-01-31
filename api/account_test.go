package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chutified/simple-bank/api"
	"github.com/chutified/simple-bank/db/mocks"
	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/chutified/simple-bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_CreateAccount(t *testing.T) {
	account := db.Account{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}

	tests := []struct {
		name          string
		apiRequest    api.CreateAccountRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			apiRequest: api.CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateAccount", mock.Anything, db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}).Return(db.Account{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				resultAccount := bytesToAccount(t, recorder.Body)
				assert.Equal(t, account.Owner, resultAccount.Owner)
				assert.Equal(t, account.Currency, resultAccount.Currency)
				assert.Equal(t, int64(0), resultAccount.Balance)
			},
		},
		{
			name: "InvalidCurrency",
			apiRequest: api.CreateAccountRequest{
				Owner:    account.Owner,
				Currency: strings.ToLower(account.Currency),
			},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			apiRequest: api.CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateAccount", mock.Anything, db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			url := "/accounts"
			b, err := json.Marshal(test.apiRequest)
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

func TestServer_ListAccounts(t *testing.T) {
	accounts := []db.Account{
		{
			ID:       util.RandomInt(1, 1024),
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		},
		{
			ID:       util.RandomInt(1025, 2048),
			Owner:    util.RandomOwner(),
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		},
	}

	tests := []struct {
		name           string
		accountRequest api.ListAccountsRequest
		buildStub      func(store *mocks.Store)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountRequest: api.ListAccountsRequest{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
				store.On("ListAccounts", mock.Anything, db.ListAccountsParams{
					Limit:  10,
					Offset: 0,
				}).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				accountsResult := bytesToAccounts(t, recorder.Body)
				assert.Equal(t, accounts, accountsResult)
			},
		},
		{
			name: "InvalidPageNumber",
			accountRequest: api.ListAccountsRequest{
				PageNum:  0,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			accountRequest: api.ListAccountsRequest{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
				store.On("ListAccounts", mock.Anything, db.ListAccountsParams{
					Limit:  10,
					Offset: 0,
				}).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			accountRequest: api.ListAccountsRequest{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
				store.On("ListAccounts", mock.Anything, db.ListAccountsParams{
					Limit:  10,
					Offset: 0,
				}).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct a server with mock db.Store
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := fmt.Sprintf(
				"/accounts?page_num=%d&page_size=%d",
				test.accountRequest.PageNum,
				test.accountRequest.PageSize,
			)
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

func bytesToAccounts(t *testing.T, data *bytes.Buffer) []db.Account {
	t.Helper()

	var aa []db.Account
	err := json.Unmarshal(data.Bytes(), &aa)
	require.NoError(t, err)

	return aa
}
