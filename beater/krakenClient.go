package beater

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type KrakenTransactions struct {
	transactions []krakenTransaction
	since        time.Time
}

type krakenTransaction struct {
	price     float64
	volume    float64
	timestamp time.Time
}

type Krakenclient interface {
	poll(pairs []string) []KrakenTransactions
}

type KrakenHTTPClient struct{}

type krakenJson struct {
	Error  []string `json:"error"`
	Result map[string]interface{} `json:"result"`
}

func (k *KrakenHTTPClient) Poll(pairs []string) KrakenTransactions {
	transactions := KrakenTransactions{transactions: []krakenTransaction{}}
	for _, pair := range pairs {
		url := fmt.Sprintf("https://api.kraken.com/0/public/Trades?pair=%s", pair)
		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		parsedBody := krakenJson{}
		if err := json.Unmarshal(body, &parsedBody); err != nil {
			panic(err)
		}
		for _, transaction := range parsedBody.Result[pair].([]interface{}) {
			price, _ := strconv.ParseFloat(transaction.([]interface{})[0].(string),64)
			volume,_ := strconv.ParseFloat(transaction.([]interface{})[1].(string),64)
			transactions.transactions = append(transactions.transactions, krakenTransaction{
				price:     price,
				volume:    volume,
				timestamp: time.Unix(int64(transaction.([]interface{})[2].(float64)), 0),
			})
		}
		since, err := strconv.ParseInt(parsedBody.Result["last"].(string),10,64)
		if err != nil {
			panic(err)
		}
		transactions.since = time.Unix(0, since)
	}
	return transactions
}
