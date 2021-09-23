package service

import (
	"testing"
	"time"

	"github.com/matbarofex/mtz-crypto/mocks"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store/memory"
	"github.com/stretchr/testify/assert"
)

func TestGetWalletValue(t *testing.T) {
	items := []model.WalletItem{
		{Symbol: "SYM1", Quantity: 0.1},
		{Symbol: "SYM2", Quantity: 0.1},
		{Symbol: "SYM3", Quantity: 0.1},
	}
	wallet := model.Wallet{
		ID:    "wallet1",
		Items: items,
	}

	mdStore := memory.NewMarketDataStore()
	mdService := NewMarketDataService(mdStore)

	ts1, _ := time.Parse(time.RFC3339, "2021-09-23T12:34:56Z")
	mdStore.SetOrUpdateMD(model.MarketData{
		Symbol:            "SYM1",
		LastPrice:         1.0,
		LastPriceDateTime: ts1,
	})
	ts2, _ := time.Parse(time.RFC3339, "2021-09-23T13:34:56Z")
	mdStore.SetOrUpdateMD(model.MarketData{
		Symbol:            "SYM2",
		LastPrice:         1.0,
		LastPriceDateTime: ts2,
	})
	ts3, _ := time.Parse(time.RFC3339, "2021-09-23T14:34:56Z")
	mdStore.SetOrUpdateMD(model.MarketData{
		Symbol:            "SYM3",
		LastPrice:         1.0,
		LastPriceDateTime: ts3,
	})

	walletStoreMock := new(mocks.WalletStore)
	walletStoreMock.On("GetWallet", "wallet1").Return(wallet, nil)

	walletService := NewWalletService(walletStoreMock, mdService)

	req := model.GetWalletValueRequest{ID: "wallet1"}
	resp, err := walletService.GetWalletValue(req)

	assert.NoError(t, err)
	assert.Equal(t, 0.3, *resp.Value)
	assert.Equal(t, ts3, *resp.DateTime)
}
