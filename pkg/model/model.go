package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type MarketData struct {
	Symbol            string
	LastPrice         decimal.Decimal
	LastPriceDateTime time.Time
}

type Wallet struct {
	ID    string
	Items []WalletItem
}

type WalletItem struct {
	Symbol   string
	Quantity decimal.Decimal
}

type GetWalletValueRequest struct {
	ID string
}

type GetWalletValueResponse struct {
	ID       string
	Value    decimal.NullDecimal
	DateTime *time.Time
}

type MdChannel chan MarketData
