package cryptonator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/matbarofex/mtz-crypto/pkg/config"
	"github.com/matbarofex/mtz-crypto/pkg/crypto"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type cryptonatorClient struct {
	config      *config.Config
	logger      *zap.Logger
	httpClient  *http.Client
	symbolPairs []cryptonatorSymbolPair
	baseURL     string
	mdChannel   model.MdChannel
}

type cryptonatorSymbolPair struct {
	Symbol         string
	ExternalSymbol string
}

type cryptonatorTicker struct {
	Base   string              `json:"base"`
	Target string              `json:"target"`
	Price  decimal.NullDecimal `json:"price"`
}

type cryptonatorResponse struct {
	Ticker    cryptonatorTicker `json:"ticker"`
	Timestamp int               `json:"timestamp"`
	Success   bool              `json:"success"`
	Error     string            `json:"error"`
}

func NewCryptonatorClient(
	config *config.Config,
	logger *zap.Logger,
	httpClient *http.Client,
	mdChannel model.MdChannel,
) crypto.Client {
	baseURL := config.GetString("crypto.api.cryptonator.url")
	symbolPairsStringSlice := config.GetStringSlice("crypto.api.cryptonator.pairs")

	symbolPairs := []cryptonatorSymbolPair{}
	for _, pairStr := range symbolPairsStringSlice {
		pairs := strings.Split(pairStr, ";")
		symbolPairs = append(symbolPairs, cryptonatorSymbolPair{
			ExternalSymbol: pairs[0],
			Symbol:         pairs[1],
		})
	}

	return &cryptonatorClient{
		config:      config,
		logger:      logger,
		httpClient:  httpClient,
		baseURL:     baseURL,
		symbolPairs: symbolPairs,
		mdChannel:   mdChannel,
	}
}

func (c *cryptonatorClient) Start() {
	// Actualización inicial
	c.updateMarketData()
	c.logger.Info("cryptonatorClient started")

	// Actualización periódica de la market data
	interval := c.config.GetDuration("crypto.api.cryptonator.poll.interval")
	ticker := time.NewTicker(interval)
	workers := c.config.GetInt("crypto.api.cryptonator.workers")

	go func() {
		for range ticker.C {
			c.updateMarketDataWithWorkers(workers)
		}
	}()
}

// updateMarketData Obtiene la MD de todos los activos y la publica en el channel
func (c *cryptonatorClient) updateMarketData() {
	wg := sync.WaitGroup{}
	wg.Add(len(c.symbolPairs))

	for _, p := range c.symbolPairs {
		go func(pair cryptonatorSymbolPair) {
			c.logger.Debug("requesting MD", zap.String("externalSymbol", pair.ExternalSymbol))

			md, err := c.retrieveMD(pair.ExternalSymbol, pair.Symbol)
			if err != nil {
				c.logger.Error("error requesting MD",
					zap.String("externalSymbol", pair.ExternalSymbol),
					zap.Error(err))
			}

			c.mdChannel <- md
			wg.Done()
		}(p)
	}

	wg.Wait()
}

// updateMarketDataWithWorkers Obtiene la MD limitando la concurrencia a una cantidad de workers
func (c *cryptonatorClient) updateMarketDataWithWorkers(workers int) {
	ch := make(chan cryptonatorSymbolPair)
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for pair := range ch {
				c.logger.Debug("requesting MD",
					zap.Int("worker", workerID),
					zap.String("externalSymbol", pair.ExternalSymbol))

				md, err := c.retrieveMD(pair.ExternalSymbol, pair.Symbol)
				if err != nil {
					c.logger.Error("error requesting MD",
						zap.String("externalSymbol", pair.ExternalSymbol),
						zap.Error(err))
				}

				c.logger.Debug("sending MD to channel",
					zap.Int("worker", workerID),
					zap.Any("md", md))

				c.mdChannel <- md
			}

			wg.Done()
		}(i)
	}

	for _, p := range c.symbolPairs {
		ch <- p
	}

	close(ch)
	wg.Wait()
}

// retrieveMD Obtiene la Market data de un activo
func (c *cryptonatorClient) retrieveMD(externalSymbol, symbol string) (md model.MarketData, err error) {
	httpResp, err := c.httpClient.Get(fmt.Sprintf("%s/ticker/%s", c.baseURL, externalSymbol))
	if err != nil {
		return md, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return md, fmt.Errorf("invalid HTTP status code: %d", httpResp.StatusCode)
	}

	var resp cryptonatorResponse

	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return md, err
	}

	if !resp.Success {
		return md, errors.New(resp.Error)
	}

	if !resp.Ticker.Price.Valid {
		return md, fmt.Errorf("last price not found")
	}

	timestamp := time.Unix(int64(resp.Timestamp), 0)
	md = model.MarketData{
		Symbol:            symbol,
		LastPrice:         resp.Ticker.Price.Decimal,
		LastPriceDateTime: timestamp,
	}

	return md, nil
}
