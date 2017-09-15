package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/PPACI/krakenbeat/config"
)

type Krakenbeat struct {
	done         chan struct{}
	config       config.Config
	client       publisher.Client
	krakenClient Krakenclient
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	krakenclient := KrakenHTTPClient{}
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Krakenbeat{
		done:         make(chan struct{}),
		config:       config,
		krakenClient: &krakenclient,
	}
	return bt, nil
}

func (bt *Krakenbeat) Run(b *beat.Beat) error {
	logp.Info("krakenbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	lastPoll := time.Now()

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		krakenTransactions, err := bt.krakenClient.Poll(bt.config.Pairs, lastPoll)
		if err != nil {
			continue
		}
		for _, transaction := range krakenTransactions.transactions {
			logp.Info("%+v", transaction)
			event := common.MapStr{
				"timestamp_transaction": common.Time(transaction.timestamp),
				"pair":      transaction.pair,
				"price":     transaction.price,
				"volume":    transaction.volume,
				"type":      b.Name,
				"@timestamp": common.Time(time.Now()),
			}
			bt.client.PublishEvent(event)
		}
		logp.Info("%d Event sent", len(krakenTransactions.transactions))
		lastPoll = krakenTransactions.since
	}
}

func (bt *Krakenbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
