package model

import "time"

type MarketData struct {
	Symbol            string
	LastPrice         float64
	LastPriceDateTime time.Time
}
