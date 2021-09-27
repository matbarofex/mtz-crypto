package model

import (
	"errors"
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
	ID       string              `json:"walletId"`
	Value    decimal.NullDecimal `json:"value"`
	DateTime *time.Time          `json:"dateTime,omitempty"`
}

type MdChannel chan MarketData

var (
	// TODO agregar el resto de los errores
	ErrWalletIsRequired = errors.New("wallet is required")
	ErrUnexpected       = errors.New("unexpected error")
)
