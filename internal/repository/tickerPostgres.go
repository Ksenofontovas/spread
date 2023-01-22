package repository

import (
	"fmt"
	"time"

	"github.com/Ksenofontovas/spread"
	"github.com/jmoiron/sqlx"
)

type TickerPostgres struct {
	db *sqlx.DB
}

func NewTickerPostgres(db *sqlx.DB) *TickerPostgres {
	return &TickerPostgres{db: db}
}

func (r *TickerPostgres) SaveTicker(ticker spread.TickerResponse, time time.Time) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (symbol, price, exchange, time) values ($1, $2, $3, $4) RETURNING id", symbolsTable)
	row := r.db.QueryRow(query, ticker.Symbol, ticker.Price, ticker.Exchange, time)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
