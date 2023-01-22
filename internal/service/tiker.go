package service

import (
	"time"

	"github.com/Ksenofontovas/spread"
	"github.com/Ksenofontovas/spread/internal/repository"
)

type TickerService struct {
	repo repository.Ticker
}

func NewTickerService(repo repository.Ticker) *TickerService {
	return &TickerService{repo: repo}
}

func (s *TickerService) SaveTicker(tiker spread.TickerResponse, time time.Time) (int, error) {
	return s.repo.SaveTicker(tiker, time)
}
