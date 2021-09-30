package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matbarofex/mtz-crypto/mocks"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWalletControllerOK(t *testing.T) {
	svcReq := model.GetWalletValueRequest{ID: "wallet1"}
	ts, _ := time.Parse(time.RFC3339, "2021-08-03T12:34:56Z")
	svcResp := model.GetWalletValueResponse{
		ID:       "wallet1",
		Value:    decimal.NullDecimal{Decimal: decimal.RequireFromString("123.456"), Valid: true},
		DateTime: &ts,
	}

	walletServiceMock := new(mocks.WalletService)
	walletServiceMock.On("GetWalletValue", svcReq).Return(svcResp, nil)

	logger := zap.NewNop()
	walletController := NewWalletController(logger, walletServiceMock)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/wallet/value", walletController.GetWalletValue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wallet/value?wallet=wallet1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"walletId":"wallet1","value":"123.456","dateTime":"2021-08-03T12:34:56Z"}`, w.Body.String())
	walletServiceMock.AssertExpectations(t)
}

func TestWalletControllerWithoutWallet(t *testing.T) {
	logger := zap.NewNop()
	walletServiceMock := new(mocks.WalletService)
	walletController := NewWalletController(logger, walletServiceMock)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/wallet/value", walletController.GetWalletValue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wallet/value", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"wallet is required"}`, w.Body.String())
}

func TestWalletControllerWithoutValue(t *testing.T) {
	walletServiceMock := new(mocks.WalletService)
	svcReq := model.GetWalletValueRequest{ID: "wallet1"}
	svcResp := model.GetWalletValueResponse{ID: "wallet1", Value: decimal.NullDecimal{Valid: false}}
	walletServiceMock.On("GetWalletValue", svcReq).Return(svcResp, nil)

	logger := zap.NewNop()
	walletController := NewWalletController(logger, walletServiceMock)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/wallet/value", walletController.GetWalletValue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wallet/value?wallet=wallet1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"walletId":"wallet1","value":null}`, w.Body.String())
	walletServiceMock.AssertExpectations(t)
}
