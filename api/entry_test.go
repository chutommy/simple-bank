package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chutommy/simple-bank/api"
	"github.com/chutommy/simple-bank/db/mocks"
	db "github.com/chutommy/simple-bank/db/sqlc"
	"github.com/chutommy/simple-bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_GetEntryByID(t *testing.T) {
	entry := db.Entry{
		ID:        util.RandomInt(1, 2048),
		AccountID: util.RandomInt(1, 2048),
		Amount:    util.RandomAmount(),
	}

	tests := []struct {
		name          string
		params        api.GetEntryByIDRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: api.GetEntryByIDRequest{ID: entry.ID},
			buildStub: func(store *mocks.Store) {
				store.On("GetEntry", mock.Anything, entry.ID).Return(entry, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, entry, bufferToEntry(t, resp.Body))
			},
		},
		{
			name:      "InvalidID",
			params:    api.GetEntryByIDRequest{ID: 0},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:   "NotFound",
			params: api.GetEntryByIDRequest{ID: entry.ID},
			buildStub: func(store *mocks.Store) {
				store.On("GetEntry", mock.Anything, entry.ID).Return(db.Entry{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:   "InternalError",
			params: api.GetEntryByIDRequest{ID: entry.ID},
			buildStub: func(store *mocks.Store) {
				store.On("GetEntry", mock.Anything, entry.ID).Return(db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct a server with a mock db.Store
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := fmt.Sprintf("/entries/id/%d", test.params.ID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			resp := httptest.NewRecorder()

			// server request
			server.Srv.Handler.ServeHTTP(resp, req)

			// check result
			test.checkResponse(t, resp)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_ListEntries(t *testing.T) {
	account := db.Account{
		ID: util.RandomInt(0, 2048),
	}

	entries := []db.Entry{
		{
			ID:        util.RandomInt(1, 1024),
			AccountID: account.ID,
			Amount:    util.RandomAmount(),
		},
		{
			ID:        util.RandomInt(1025, 2048),
			AccountID: account.ID,
			Amount:    util.RandomAmount(),
		},
	}

	tests := []struct {
		name          string
		paramURI      api.ListEntriesRequestURI
		paramQuery    api.ListEntriesRequestQuery
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			paramURI: api.ListEntriesRequestURI{
				AccountID: account.ID,
			},
			paramQuery: api.ListEntriesRequestQuery{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
				store.On("ListEntries", mock.Anything, db.ListEntriesParams{
					AccountID: account.ID,
					Limit:     10,
					Offset:    0,
				}).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, entries, bufferToEntries(t, resp.Body))
			},
		},
		{
			name: "InvalidID",
			paramURI: api.ListEntriesRequestURI{
				AccountID: 0,
			},
			paramQuery: api.ListEntriesRequestQuery{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "InvalidPageNum",
			paramURI: api.ListEntriesRequestURI{
				AccountID: account.ID,
			},
			paramQuery: api.ListEntriesRequestQuery{
				PageNum:  0,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "NotFound",
			paramURI: api.ListEntriesRequestURI{
				AccountID: account.ID,
			},
			paramQuery: api.ListEntriesRequestQuery{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
				store.On("ListEntries", mock.Anything, db.ListEntriesParams{
					AccountID: account.ID,
					Limit:     10,
					Offset:    0,
				}).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "InternalError",
			paramURI: api.ListEntriesRequestURI{
				AccountID: account.ID,
			},
			paramQuery: api.ListEntriesRequestQuery{
				PageNum:  1,
				PageSize: 10,
			},
			buildStub: func(store *mocks.Store) {
				store.On("ListEntries", mock.Anything, db.ListEntriesParams{
					AccountID: account.ID,
					Limit:     10,
					Offset:    0,
				}).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct server with mock db
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := fmt.Sprintf("/entries/accountid/%d?page_num=%d&page_size=%d",
				test.paramURI.AccountID, test.paramQuery.PageNum, test.paramQuery.PageSize)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			resp := httptest.NewRecorder()

			// serve
			server.Srv.Handler.ServeHTTP(resp, req)

			// check result
			test.checkResponse(t, resp)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_CreateEntry(t *testing.T) {
	entry := db.Entry{
		ID:        util.RandomInt(1, 1024),
		AccountID: util.RandomInt(0, 2048),
		Amount:    util.RandomAmount(),
	}

	tests := []struct {
		name          string
		param         api.CreateEntryRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: api.CreateEntryRequest{
				AccountID: entry.AccountID,
				Amount:    entry.Amount,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateEntry", mock.Anything, db.CreateEntryParams{
					AccountID: entry.AccountID,
					Amount:    entry.Amount,
				}).Return(entry, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, entry, bufferToEntry(t, resp.Body))
			},
		},
		{
			name: "InvalidRequest",
			param: api.CreateEntryRequest{
				Amount: entry.Amount,
			},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "AccountNotFound",
			param: api.CreateEntryRequest{
				AccountID: entry.AccountID,
				Amount:    entry.Amount,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateEntry", mock.Anything, db.CreateEntryParams{
					AccountID: entry.AccountID,
					Amount:    entry.Amount,
				}).Return(db.Entry{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "InternalError",
			param: api.CreateEntryRequest{
				AccountID: entry.AccountID,
				Amount:    entry.Amount,
			},
			buildStub: func(store *mocks.Store) {
				store.On("CreateEntry", mock.Anything, db.CreateEntryParams{
					AccountID: entry.AccountID,
					Amount:    entry.Amount,
				}).Return(db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct server with mock db
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := "/entries"
			b, err := json.Marshal(test.param)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))
			resp := httptest.NewRecorder()

			// serve
			server.Srv.Handler.ServeHTTP(resp, req)

			// check result
			test.checkResponse(t, resp)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_UpdateEntry(t *testing.T) {
	account := db.Account{
		ID: util.RandomInt(1, 2048),
	}

	entry1 := db.Entry{
		ID:        util.RandomInt(1, 1024),
		AccountID: account.ID,
		Amount:    util.RandomAmount(),
	}

	entry2 := db.Entry{
		ID:        entry1.ID,
		AccountID: entry1.AccountID,
		Amount:    util.RandomAmount(),
	}

	tests := []struct {
		name          string
		paramURI      api.UpdateEntryRequestURI
		paramJSON     api.UpdateEntryRequestJSON
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			paramURI:  api.UpdateEntryRequestURI{ID: entry1.ID},
			paramJSON: api.UpdateEntryRequestJSON{Amount: entry2.Amount},
			buildStub: func(store *mocks.Store) {
				store.On("UpdateEntryAmount", mock.Anything, db.UpdateEntryAmountParams{
					ID:     entry1.ID,
					Amount: entry2.Amount,
				}).Return(entry2, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, entry2, bufferToEntry(t, resp.Body))
			},
		},
		{
			name:      "InvalidURI",
			paramURI:  api.UpdateEntryRequestURI{ID: 0},
			paramJSON: api.UpdateEntryRequestJSON{Amount: entry2.Amount},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:      "InvalidJSON",
			paramURI:  api.UpdateEntryRequestURI{ID: entry1.ID},
			paramJSON: api.UpdateEntryRequestJSON{},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:      "NotFound",
			paramURI:  api.UpdateEntryRequestURI{ID: entry1.ID},
			paramJSON: api.UpdateEntryRequestJSON{Amount: entry2.Amount},
			buildStub: func(store *mocks.Store) {
				store.On("UpdateEntryAmount", mock.Anything, db.UpdateEntryAmountParams{
					ID:     entry1.ID,
					Amount: entry2.Amount,
				}).Return(db.Entry{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:      "InternalError",
			paramURI:  api.UpdateEntryRequestURI{ID: entry1.ID},
			paramJSON: api.UpdateEntryRequestJSON{Amount: entry2.Amount},
			buildStub: func(store *mocks.Store) {
				store.On("UpdateEntryAmount", mock.Anything, db.UpdateEntryAmountParams{
					ID:     entry1.ID,
					Amount: entry2.Amount,
				}).Return(db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct server with mock db
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := fmt.Sprintf("/entries/%d", test.paramURI.ID)
			b, err := json.Marshal(test.paramJSON)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(b))
			resp := httptest.NewRecorder()

			// serve
			server.Srv.Handler.ServeHTTP(resp, req)

			// check result
			test.checkResponse(t, resp)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestServer_DeleteEntry(t *testing.T) {
	entry := db.Entry{
		ID:        util.RandomInt(1, 1024),
		AccountID: util.RandomInt(1, 2048),
		Amount:    util.RandomAmount(),
	}

	tests := []struct {
		name          string
		param         api.DeleteEntryRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: api.DeleteEntryRequest{ID: entry.ID},
			buildStub: func(store *mocks.Store) {
				store.On("DeleteEntry", mock.Anything, entry.ID).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:      "InvalidID",
			param:     api.DeleteEntryRequest{ID: 0},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name:  "NotFound",
			param: api.DeleteEntryRequest{ID: entry.ID},
			buildStub: func(store *mocks.Store) {
				store.On("DeleteEntry", mock.Anything, entry.ID).Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:  "InternalError",
			param: api.DeleteEntryRequest{ID: entry.ID},
			buildStub: func(store *mocks.Store) {
				store.On("DeleteEntry", mock.Anything, entry.ID).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// construct server with mock db
			mockStore := new(mocks.Store)
			server := api.NewServer(mockStore)
			test.buildStub(mockStore)

			// prepare request and response recorder
			url := fmt.Sprintf("/entries/%d", test.param.ID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			resp := httptest.NewRecorder()

			// serve
			server.Srv.Handler.ServeHTTP(resp, req)

			// check result
			test.checkResponse(t, resp)
			mockStore.AssertExpectations(t)
		})
	}
}

func bufferToEntry(t *testing.T, b *bytes.Buffer) db.Entry {
	t.Helper()

	var entry db.Entry

	require.NoError(t, json.Unmarshal(b.Bytes(), &entry))

	return entry
}

func bufferToEntries(t *testing.T, b *bytes.Buffer) []db.Entry {
	t.Helper()

	var entries []db.Entry
	err := json.Unmarshal(b.Bytes(), &entries)
	require.NoError(t, err)

	return entries
}
