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

// Postgres
var (
	_ = fs.String("crypto.postgres.host", "localhost", "Host de la base de Postgres")
	_ = fs.Int("crypto.postgres.port", 5432, "Puerto de la base de Postgres")
	_ = fs.String("crypto.postgres.dbname", "crypto_db_dev", "Nombre la base de Postgres")
	_ = fs.String("crypto.postgres.username", "postgres", "Nombre de usuario para la conexión a Postgres")
	_ = fs.String("crypto.postgres.password", "postgres", "Contraseña para la conexión a Postgres")
	_ = fs.Int("crypto.postgres.maxidleconns", 10, "Máximo de conexiones inactivas de Postgres")
	_ = fs.Bool("crypto.postgres.sqldebug", false, "Debug del SQL de Postgres")
	_ = fs.String("crypto.postgres.extraopts", "connect_timeout=10 application_name=crypto sslmode=disable", "Opciones extra para la conexión a Postgres")
)
