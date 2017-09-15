package beater

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"
)

type KrakenTransactions struct {
	transactions []krakenTransaction
	since        map[string]time.Time
}

type krakenTransaction struct {
	price     float64
	volume    float64
	timestamp time.Time
	pair      string
}

type Krakenclient interface {
	Poll(pairs []string, since map[string]time.Time) (KrakenTransactions, error)
}

type KrakenHTTPClient struct{}

type krakenJson struct {
	Error  []string               `json:"error"`
	Result map[string]interface{} `json:"result"`
}

func (k *KrakenHTTPClient) Poll(pairs []string, lastPolls map[string]time.Time) (KrakenTransactions, error) {
	transactions := KrakenTransactions{transactions: []krakenTransaction{}, since: lastPolls}
	for _, pair := range pairs {
		req, err := http.NewRequest("GET", "https://api.kraken.com/0/public/Trades", nil)
		if err != nil {
			return KrakenTransactions{}, err
		}
		req.URL.RawQuery = url2.Values{
			"pair":  []string{pair},
			"since": []string{strconv.FormatInt(lastPolls[pair].UnixNano(), 10)},
		}.Encode()
		resp, err := (&http.Client{}).Do(req)
		defer resp.Body.Close()
		if err != nil {
			return KrakenTransactions{}, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return KrakenTransactions{}, err
		}
		parsedBody := krakenJson{}
		if err := json.Unmarshal(body, &parsedBody); err != nil {
			return KrakenTransactions{}, err
		}
		for _, transaction := range parsedBody.Result[pair].([]interface{}) {
			price, _ := strconv.ParseFloat(transaction.([]interface{})[0].(string), 64)
			volume, _ := strconv.ParseFloat(transaction.([]interface{})[1].(string), 64)
			timestamp := transaction.([]interface{})[2].(float64)
			transactions.transactions = append(transactions.transactions, krakenTransaction{
				price:     price,
				volume:    volume,
				timestamp: krakenTimestampToUnixTime(timestamp),
				pair:      pair,
			})
		}
		since, err := strconv.ParseInt(parsedBody.Result["last"].(string), 10, 64)
		if err != nil {
			return KrakenTransactions{}, err
		}
		transactions.since[pair] = time.Unix(0, since)
	}
	return transactions, nil
}

func krakenTimestampToUnixTime(timestamp float64) time.Time {
	second := int64(timestamp)
	nano_part := int64(1000 * (timestamp - float64(int64(timestamp))))
	nanoduration := time.Duration(nano_part) * time.Millisecond
	nanosecond := nanoduration.Nanoseconds()
	return time.Unix(second, nanosecond)
}
