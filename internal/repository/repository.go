package repository

import (
	"time"

	"github.com/Ksenofontovas/spread"
	"github.com/jmoiron/sqlx"
)

const (
	symbolsTable = "symbols"
)

type Exchanger interface {
	GetTickers(pairs []string) (resp []spread.TickerResponse, err error)
}

type Ticker interface {
	SaveTicker(tiker spread.TickerResponse, time time.Time) (int, error)
}

type Repository struct {
	Exchanger
	//	Binance
	Ticker
	//Exchange []Exchange
}

func NewRepository(db *sqlx.DB, exchanges []*Exchange) *Repository {

	return &Repository{
		Exchanger: NewExchanges(exchanges),
		//Binance: NewBinanceClient(hos),
		Ticker: NewTickerPostgres(db),
	}
}
