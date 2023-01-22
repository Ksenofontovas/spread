package service

import (
	"time"

	"github.com/Ksenofontovas/spread"
	"github.com/Ksenofontovas/spread/internal/repository"
)

type Exchanger interface {
	GetTickers(pairs []string) (resp []spread.TickerResponse, err error)
}

type Binance interface {
	GetTickers(pairs []string) (resp []spread.TickerResponse, err error)
}

type Ticker interface {
	SaveTicker(tiker spread.TickerResponse, time time.Time) (int, error)
}

type Service struct {
	Exchanger
	Ticker
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Exchanger: NewExchangeService(repos.Exchanger),
		Ticker:    NewTickerService(repos.Ticker),
	}
}
