package cryptonator

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matbarofex/mtz-crypto/pkg/config"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func (c *cryptonatorClient) TestUpdateMarketData() {
	c.updateMarketData()
}

func TestUpdateMarketData(t *testing.T) {
	// start test server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Validamos que la URL contenga el símbolo externo esperado
		assert.Regexp(t, "\\/ticker\\/externalSymbol(\\d+)$", req.URL.Path)

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")
		_, err := rw.Write([]byte(`{
			"ticker": {
			  "price": "123.456",
			  "volume": "1.0",
			  "change": "-0.01"
			},
			"timestamp": 1628610304,
			"success": true,
			"error": ""
		}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	t.Setenv("MTZ_CRYPTO_API_CRYPTONATOR_URL", server.URL)
	t.Setenv("MTZ_CRYPTO_API_CRYPTONATOR_PAIRS", "externalSymbol1;SYMBOL1")

	cfg := config.NewConfig(&flag.FlagSet{})
	mdChannel := make(chan model.MarketData)
	client := NewCryptonatorClient(cfg, server.Client(), mdChannel)
	go client.(*cryptonatorClient).updateMarketData()

	md := <-mdChannel
	assert.Equal(t, "SYMBOL1", md.Symbol)
	assert.Equal(t, decimal.RequireFromString("123.456"), md.LastPrice)
	assert.Equal(t, int64(1_628_610_304), md.LastPriceDateTime.Unix())
}
