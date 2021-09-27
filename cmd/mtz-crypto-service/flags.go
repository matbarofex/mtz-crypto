package main

import (
	"flag"

	"github.com/spf13/pflag"
)

var fs = flag.NewFlagSet("crypto", flag.ExitOnError)

// Cryptonator (API externa)
var (
	_ = fs.String("crypto.api.cryptonator.url", "https://api.cryptonator.com/api", "URL API de servicio cryptonator")
	_ = pflag.StringSlice("crypto.api.cryptonator.pairs", []string{
		"btc-usd;BTCUSD",
		"eth-usd;ETHUSD",
		"ada-usd;ADAUSD",
		"dot-usd;DOTUSD",
	}, "Pares 'simbolo externo;simbolo interno'")
)
