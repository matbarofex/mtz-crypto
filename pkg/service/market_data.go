package service

import (
	"fmt"

	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
)

type MarketDataService interface {
	GetMD(symbol string) (md model.MarketData, err error)
	ConsumeMD(mdChannel model.MdChannel)
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

func (s *marketDataService) ConsumeMD(mdChannel model.MdChannel) {
	go func() {
		for md := range mdChannel {
			fmt.Println("Nueva MD", md)

			if err := s.mdStore.SetOrUpdateMD(md); err != nil {
				// TODO log de error
				fmt.Println("Error", err)
			}
		}
	}()
}
