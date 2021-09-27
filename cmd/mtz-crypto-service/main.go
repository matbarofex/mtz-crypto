package main

import (
	"net/http"

	"github.com/matbarofex/mtz-crypto/pkg/config"
	"github.com/matbarofex/mtz-crypto/pkg/crypto/cryptonator"
)

func main() {
	cfg := config.NewConfig(fs)

	httpClient := &http.Client{}
	cryptoClient := cryptonator.NewCryptonatorClient(cfg, httpClient)
	cryptoClient.Start()
}
