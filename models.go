package spread

type TickerResponse struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	Price    string `json:"price"`
}

type Symbols struct {
	Id       int     `json:"-"`
	Exchange string  `json:"exchange" binding:"required"`
	Symbol   string  `json:"symbol" binding:"required"`
	Price    float32 `json:"price" binding:"required"`
}
