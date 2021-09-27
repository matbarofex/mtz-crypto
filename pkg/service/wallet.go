package service

import (
	"time"

	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
	"github.com/shopspring/decimal"
)

type WalletService interface {
	GetWalletValue(req model.GetWalletValueRequest) (rs model.GetWalletValueResponse, err error)
}

type walletService struct {
	mdService   MarketDataService
	walletStore store.WalletStore
}

func NewWalletService(walletStore store.WalletStore, mdService MarketDataService) WalletService {
	return &walletService{
		mdService:   mdService,
		walletStore: walletStore,
	}
}

func (s *walletService) GetWalletValue(req model.GetWalletValueRequest) (rs model.GetWalletValueResponse, err error) {
	wallet, err := s.walletStore.GetWallet(req.ID)
	if err != nil {
		return rs, err
	}

	rs.ID = req.ID

	var datetime time.Time

	valueIsNull := true
	value := decimal.Zero

	for _, item := range wallet.Items {
		md, err := s.mdService.GetMD(item.Symbol)
		if err != nil {
			return rs, err
		}

		valueIsNull = false
		value = value.Add(md.LastPrice.Mul(item.Quantity))

		if md.LastPriceDateTime.After(datetime) {
			datetime = md.LastPriceDateTime
		}
	}

	if !valueIsNull {
		rs.Value = decimal.NullDecimal{Valid: true, Decimal: value}
		rs.DateTime = &datetime
	}

	return rs, err
}
