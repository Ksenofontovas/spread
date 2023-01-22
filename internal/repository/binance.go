package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Ksenofontovas/spread"
)

const (
	tickerURL    = "api/v2/trades"
	exchangeName = "Garantex"
)

// type BinanceTickerSResponse struct {
// 	// ID        int       `json:"id"`
// 	Price string `json:"price"`
// 	// Volume    string    `json:"volume"`
// 	// Funds     string    `json:"funds"`
// 	Market string `json:"market"`
// 	// CreatedAt time.Time `json:"created_at"`

// }

type BinanceClient struct {
	host       string
	httpClient *http.Client
}

func NewBinanceClient(host string, timeout time.Duration) *BinanceClient {
	client := &http.Client{
		Timeout: timeout,
	}
	return &BinanceClient{
		host:       host,
		httpClient: client,
	}
}

func (c *BinanceClient) do(method, endpoint string, params map[string]string) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", c.host, endpoint)
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
	return c.httpClient.Do(req)
}

func (c *BinanceClient) GetTickers(pairs []string) (resp []spread.TickerResponse, err error) {
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

	pairsMap := make(map[string]string)
	pairsMap["symbols"] = symbols
	res, err := c.do(http.MethodGet, tickerURL, pairsMap)
	fmt.Println(res)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}
	return
}
