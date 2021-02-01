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
	"github.com/go-playground/assert/v2"
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
			url := fmt.Sprintf("/entries/%d", test.params.ID)
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

func bufferToEntry(t *testing.T, b *bytes.Buffer) db.Entry {
	t.Helper()

	var entry db.Entry

	require.NoError(t, json.Unmarshal(b.Bytes(), &entry))

	return entry
}
