package store

import "github.com/matbarofex/mtz-crypto/pkg/model"

type WalletStore interface {
	GetWallet(id string) (rs model.Wallet, err error)
}

type MarketDataStore interface {
	GetMD(symbol string) (rs model.MarketData, err error)
	SetOrUpdateMD(md model.MarketData) (err error)
}
