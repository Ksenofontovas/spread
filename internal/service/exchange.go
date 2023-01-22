package service

import (
	"github.com/Ksenofontovas/spread"
	"github.com/Ksenofontovas/spread/internal/repository"
)

type ExchangeService struct {
	repo repository.Exchanger
}

func NewExchangeService(repo repository.Exchanger) *ExchangeService {
	return &ExchangeService{repo: repo}
}

func (s *ExchangeService) GetTickers(pairs []string) (resp []spread.TickerResponse, err error) {
	return s.repo.GetTickers(pairs)
}
