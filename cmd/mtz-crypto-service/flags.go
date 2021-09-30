package main

import (
	"flag"
	"time"

	"github.com/spf13/pflag"
)

var fs = flag.NewFlagSet("crypto", flag.ExitOnError)

// Configs generales
var (
	_ = fs.Bool("crypto.debug.mode", false, "Activar modo debug")
	_ = fs.String("crypto.http.addr", ":8000", "Puerto HTTP del servicio")
	_ = fs.Duration("crypto.http.shutdown.timeout", 15*time.Second, "HTTP server graceful shutdown timeout")
)

// Cryptonator (API externa)
var (
	_ = fs.String("crypto.api.cryptonator.url", "https://api.cryptonator.com/api", "URL API de servicio cryptonator")
	_ = fs.Duration("crypto.api.cryptonator.poll.interval", 15*time.Second, "Intervalo de consulta")
	_ = fs.Duration("crypto.api.cryptonator.timeout", 15*time.Second, "Timeout para solicitudes a Cryptonator API")
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
