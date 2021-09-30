package memory

import (
	"errors"
	"sync"

	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
)

type marketDataStore struct {
	data sync.Map
}

func NewMarketDataStore() store.MarketDataStore {
	return &marketDataStore{
		data: sync.Map{},
	}
}

func (s *marketDataStore) GetMD(symbol string) (rs model.MarketData, err error) {
	value, ok := s.data.Load(symbol)
	if !ok {
		return rs, errors.New("symbol not found")
	}

	return value.(model.MarketData), nil
}

func (s *marketDataStore) SetOrUpdateMD(md model.MarketData) (err error) {
	s.data.Store(md.Symbol, md)

	return nil
}
