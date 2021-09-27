package cryptonator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/matbarofex/mtz-crypto/pkg/crypto"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/shopspring/decimal"
)

type cryptonatorClient struct {
	httpClient  *http.Client
	symbolPairs []cryptonatorSymbolPair
	baseURL     string
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
	baseURL string,
	httpClient *http.Client,
) crypto.Client {
	// TODO externalizar configs
	symbolPairs := []cryptonatorSymbolPair{
		{Symbol: "BTCUSD", ExternalSymbol: "btc-usd"},
		// TODO resto de los activos
	}

	return &cryptonatorClient{
		httpClient:  httpClient,
		baseURL:     baseURL,
		symbolPairs: symbolPairs,
	}
}

func (c *cryptonatorClient) Start() {
	// Actualización inicial
	c.updateMarketData()

	// TODO Actualización periódica de la market data
}

// updateMarketData Actualiza la Market Data de todos los activos, obteniéndolos
// del servicio externo
func (c *cryptonatorClient) updateMarketData() {
	for _, pair := range c.symbolPairs {
		md, err := c.retrieveMD(pair.ExternalSymbol, pair.Symbol)
		if err != nil {
			// TODO log
		}

		// TODO actualizar MD (enviar a channel)
		fmt.Println("MD", md)
	}
}

// retrieveAndUpdateMD Obtiene la Market data de 1 activo y la actualiza
// localmente, invocando al marketDataService
func (c *cryptonatorClient) retrieveMD(externalSymbol, symbol string) (md model.MarketData, err error) {
	httpResp, err := c.httpClient.Get(fmt.Sprintf("%s/ticker/%s", c.baseURL, externalSymbol))
	if err != nil {
		return md, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return md, fmt.Errorf("invalid HTTP status code: %d", httpResp.StatusCode)
	}

	var resp cryptonatorResponse
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
