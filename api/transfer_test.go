package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

func TestServer_MakeTransfer(t *testing.T) {
	account1 := db.Account{
		ID:       util.RandomInt(1, 1024),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	account2 := db.Account{
		ID:       util.RandomInt(1025, 2048),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	transfer := db.Transfer{
		ID:            util.RandomInt(1, 2048),
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}

	tests := []struct {
		name          string
		param         api.MakeTransferRequest
		buildStub     func(store *mocks.Store)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: api.MakeTransferRequest{
				FromAccountID: transfer.FromAccountID,
				ToAccountID:   transfer.ToAccountID,
				Amount:        transfer.Amount,
			},
			buildStub: func(store *mocks.Store) {
				store.On("TransferTx", mock.Anything, db.TransferTxParams{
					FromAccountID: transfer.FromAccountID,
					ToAccountID:   transfer.ToAccountID,
					Amount:        transfer.Amount,
				}).Return(db.TransferTxResult{
					Transfer: db.Transfer{
						ID:            transfer.ID,
						FromAccountID: transfer.FromAccountID,
						ToAccountID:   transfer.ToAccountID,
						Amount:        transfer.Amount,
					},
				}, nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, transfer, bytesToTransfer(t, resp.Body.Bytes()))
			},
		},
		{
			name: "InvalidID",
			param: api.MakeTransferRequest{
				FromAccountID: 0,
				ToAccountID:   transfer.ToAccountID,
				Amount:        transfer.Amount,
			},
			buildStub: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			name: "AccountNotFound",
			param: api.MakeTransferRequest{
				FromAccountID: transfer.FromAccountID,
				ToAccountID:   transfer.ToAccountID,
				Amount:        transfer.Amount,
			},
			buildStub: func(store *mocks.Store) {
				store.On("TransferTx", mock.Anything, db.TransferTxParams{
					FromAccountID: transfer.FromAccountID,
					ToAccountID:   transfer.ToAccountID,
					Amount:        transfer.Amount,
				}).Return(db.TransferTxResult{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "InternalError",
			param: api.MakeTransferRequest{
				FromAccountID: transfer.FromAccountID,
				ToAccountID:   transfer.ToAccountID,
				Amount:        transfer.Amount,
			},
			buildStub: func(store *mocks.Store) {
				store.On("TransferTx", mock.Anything, db.TransferTxParams{
					FromAccountID: transfer.FromAccountID,
					ToAccountID:   transfer.ToAccountID,
					Amount:        transfer.Amount,
				}).Return(db.TransferTxResult{}, sql.ErrConnDone)
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
			url := "/transfers"
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

func bytesToTransfer(t *testing.T, b []byte) db.Transfer {
	t.Helper()

	var transfer db.TransferTxResult
	err := json.Unmarshal(b, &transfer)
	require.NoError(t, err)

	return transfer.Transfer
}
