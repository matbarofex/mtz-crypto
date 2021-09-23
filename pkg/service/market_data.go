package service

import (
	"fmt"

	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
)

type MarketDataService interface {
	GetMD(symbol string) (md model.MarketData, err error)
	// TODO update MD
}

type marketDataService struct {
	mdStore store.MarketDataStore
}

func NewMarketDataService(mdStore store.MarketDataStore) MarketDataService {
	return &marketDataService{
		mdStore: mdStore,
	}
}

func (s *marketDataService) GetMD(symbol string) (md model.MarketData, err error) {
	if symbol == "" {
		return md, fmt.Errorf("symbol is required")
	}

	return s.mdStore.GetMD(symbol)
}
