package model

import "time"

type MarketData struct {
	Symbol            string
	LastPrice         float64
	LastPriceDateTime time.Time
}

type Wallet struct {
	ID    string
	Items []WalletItem
}

type WalletItem struct {
	Symbol   string
	Quantity float64
}

type GetWalletValueRequest struct {
	ID string
}

type GetWalletValueResponse struct {
	ID       string
	Value    *float64
	DateTime *time.Time
}
