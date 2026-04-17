package simlatencies_test

import (
	"compress/gzip"
	"fmt"
	"net/netip"
	"os"
	"testing"

	simlatencies "github.com/marcopolo/ethereum-network-latencies"
)

func init() {
	f, err := os.Open("masked-ips.txt.gz")
	if err != nil {
		panic(fmt.Errorf("opening masked-ips.txt.gz: %w", err))
	}
	ipReader, err := gzip.NewReader(f)
	if err != nil {
		panic(fmt.Errorf("creating gzip reader for masked-ips.txt.gz: %w", err))
	}

	f, err = os.Open("pairwise_predictions.csv.gz")
	if err != nil {
		panic(fmt.Errorf("opening pairwise_predictions.csv.gz: %w", err))
	}
	pairwisePredictions, err := gzip.NewReader(f)
	if err != nil {
		panic(fmt.Errorf("creating gzip reader for pairwise_predictions.csv.gz: %w", err))
	}
	simlatencies.MustInit(ipReader, pairwisePredictions)
}

func TestLatencyBetweenFirstTwoIPs(t *testing.T) {
	if simlatencies.IPs[0] != netip.MustParseAddr("1.6.0.0") {
		t.Errorf("first IP is %s, want 1.6.0.0", simlatencies.IPs[0])
	}
	if simlatencies.IPs[1] != netip.MustParseAddr("1.6.0.1") {
		t.Errorf("first IP is %s, want 1.6.0.1", simlatencies.IPs[0])
	}

	latency, err := simlatencies.Latency(simlatencies.IPs[0], simlatencies.IPs[1])
	if err != nil {
		t.Errorf("Latency error: %v", err)
	}
	if latency == 0 {
		t.Errorf("latency is 0, want non-zero")
	}
	t.Logf("latency is %v\n", latency)
}

func TestAllLatencies(t *testing.T) {
	for _, a := range simlatencies.IPs {
		for _, b := range simlatencies.IPs {
			if a == b {
				continue
			}
			latency, err := simlatencies.Latency(a, b)
			if err != nil {
				t.Errorf("Latency error: %v", err)
			}
			if latency == 0 {
				t.Errorf("latency is 0, want non-zero")
			}
		}
	}
}
