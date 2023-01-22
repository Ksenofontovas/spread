package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Ksenofontovas/spread"
)

const (
	garantexTickerURL = "api/v2/trades"
	binanceTickerURL  = "api/v3/ticker/price"
)

type GarantexTickerResponse struct {
	Price  string `json:"price"`
	Market string `json:"market"`
}

type BinanceTickerResponse struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}

type Exchange struct {
	host       string
	httpClient *http.Client
	name       string
}

func NewExchange(name, host string, timeout time.Duration) *Exchange {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Exchange{
		name:       name,
		host:       host,
		httpClient: client,
	}
}

func (e *Exchange) do(method, endpoint string, params map[string]string) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", e.host, endpoint)
	req, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()
	return e.httpClient.Do(req)
}

type Exchanges struct {
	exchanges []*Exchange
}

func NewExchanges(exchanges []*Exchange) *Exchanges {
	return &Exchanges{
		exchanges: exchanges,
	}
}

func (r *Exchanges) GetTickers(pairs []string) (resp []spread.TickerResponse, err error) {

	for _, e := range r.exchanges {
		switch e.name {
		case "Garantex":
			res, err := GarantexGetTickers(e, pairs)
			if err != nil {
				return nil, err
			}
			resp = append(resp, res...)
		case "Binance":
			res, err := BinanceGetTickers(e, pairs)
			if err != nil {
				return nil, err
			}
			resp = append(resp, res...)
		default:
			return nil, errors.New("undefine exchange")
		}
	}
	return
}

func GarantexGetTickers(exchange *Exchange, pairs []string) (resp []spread.TickerResponse, err error) {

	market := make(map[string]string)
	for _, pair := range pairs {
		var garantexResponse []GarantexTickerResponse
		market["market"] = pair
		res, err := exchange.do(http.MethodGet, garantexTickerURL, market)
		if err != nil {
			return resp, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return resp, err
		}
		if err = json.Unmarshal(body, &garantexResponse); err != nil {
			return resp, err
		}
		tickerResponse := spread.TickerResponse{Symbol: garantexResponse[0].Market, Price: garantexResponse[0].Price, Exchange: exchangeName}
		resp = append(resp, tickerResponse)
	}
	return
}

func BinanceGetTickers(exchange *Exchange, pairs []string) (resp []spread.TickerResponse, err error) {
	var symbols string
	for i := 0; i < len(pairs); i++ {
		switch {
		case i == 0:
			symbols += fmt.Sprintf("[\"%v\",", pairs[i])
		case i == len(pairs)-1:
			symbols += fmt.Sprintf("\"%v\"]", pairs[i])
		default:
			symbols += fmt.Sprintf("\"%v\",", pairs[i])
		}
	}
	//	fmt.Println(strings.ToUpper(symbols))
	pairsMap := make(map[string]string)
	pairsMap["symbols"] = strings.ToUpper(symbols)
	// pairsMap := make(map[string]string)
	// pairsMap["symbols"] = `["BTCUSDT","BTCRUB","USDTRUB"]`
	res, err := exchange.do(http.MethodGet, binanceTickerURL, pairsMap)

	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	var tr []BinanceTickerResponse
	if err = json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}
	for i := 0; i < len(tr); i++ {
		tickerResponse := spread.TickerResponse{Symbol: tr[i].Symbol, Price: tr[i].Price, Exchange: "Binance"}
		resp = append(resp, tickerResponse)
	}
	return
}
