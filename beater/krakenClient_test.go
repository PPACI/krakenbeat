package beater

import (
	"testing"
	"time"
)

func TestKrakenClientPoll(t *testing.T) {
	client := KrakenHTTPClient{}
	since := map[string]time.Time{"BCHEUR":time.Now().Add(-time.Duration(2 * time.Minute))}
	t.Logf("transaction since : %v\n", since)
	t.Logf("%v\n", since["BCHEUR"].UnixNano())
	transactions, err := client.Poll([]string{"BCHEUR"}, since)
	if err != nil {
		t.Error(err)
	}
	for _, transaction := range transactions.transactions {
		t.Logf("%+v", transaction)
	}
	if len(transactions.transactions) == 0 {
		t.Error("No transaction recorded !")
	}
	zero_time := time.Time{}
	if transactions.since["BCHEUR"] == zero_time {
		t.Error("No last timestamp recorded")
	}
	if len(transactions.transactions) > 0 && transactions.since["BCHEUR"] != zero_time {
		t.Logf("%d transactions polled. Last transaction at %v\n", len(transactions.transactions), transactions.since)
	}
}
