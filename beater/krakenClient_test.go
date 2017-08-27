package beater

import (
	"testing"
	"time"
)

func TestKrakenClientPoll(t *testing.T) {
	client := KrakenHTTPClient{}
	transactions := client.Poll([]string{"BCHEUR"})
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
