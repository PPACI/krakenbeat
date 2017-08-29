package beater

import (
	"testing"
	"time"
)

func TestKrakenClientPoll(t *testing.T) {
	client := KrakenHTTPClient{}
	since := time.Now().Add(-time.Duration(2*time.Minute))
	t.Logf("transaction since : %v\n", since)
	transactions := client.Poll([]string{"BCHEUR"}, since)
	if len(transactions.transactions) == 0{
		t.Error("No transaction recorded !")
	}
	zero_time := time.Time{}
	if transactions.since == zero_time{
		t.Error("No last timestamp recorded")
	}
	if len(transactions.transactions) > 0 && transactions.since != zero_time {
	t.Logf("%d transactions polled. Last transaction at %v\n", len(transactions.transactions), transactions.since)
	}
}
