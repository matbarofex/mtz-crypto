package service

import (
	"fmt"

	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
	"go.uber.org/zap"
)

type MarketDataService interface {
	GetMD(symbol string) (md model.MarketData, err error)
	ConsumeMD(mdChannel model.MdChannel)
}

type marketDataService struct {
	logger  *zap.Logger
	mdStore store.MarketDataStore
}

func NewMarketDataService(
	logger *zap.Logger,
	mdStore store.MarketDataStore,
) MarketDataService {
	return &marketDataService{
		logger:  logger,
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
			s.logger.Debug("new MD received", zap.Any("md", md))

			if err := s.mdStore.SetOrUpdateMD(md); err != nil {
				s.logger.Error("error updating MD", zap.Any("md", md), zap.Error(err))
			}
		}
	}()
}
