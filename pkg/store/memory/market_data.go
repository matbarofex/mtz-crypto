package memory

import (
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
)

type marketDataStore struct {
	// FIXME - Atención! data race
	data map[string]model.MarketData
}

func NewMarketDataStore() store.MarketDataStore {
	return &marketDataStore{
		data: make(map[string]model.MarketData),
	}
}

func (s *marketDataStore) GetMD(symbol string) (rs model.MarketData, err error) {
	return s.data[symbol], nil
}

func (s *marketDataStore) SetOrUpdateMD(md model.MarketData) (err error) {
	s.data[md.Symbol] = md

	return nil
}
