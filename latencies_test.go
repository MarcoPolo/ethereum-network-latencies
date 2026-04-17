package simlatencies_test

import (
	"net/netip"
	"testing"

	simlatencies "github.com/marcopolo/ethereum-network-latencies"
)

var network simlatencies.Network

func init() {
	var err error
	network, err = simlatencies.NewNetworkFromFilepaths("masked-ips.txt.gz", "pairwise_predictions.csv.gz")
	if err != nil {
		panic(err)
	}
}

func TestLatencyBetweenFirstTwoIPs(t *testing.T) {
	if network.IPs[0] != netip.MustParseAddr("1.6.0.0") {
		t.Errorf("first IP is %s, want 1.6.0.0", network.IPs[0])
	}
	if network.IPs[1] != netip.MustParseAddr("1.6.0.1") {
		t.Errorf("first IP is %s, want 1.6.0.1", network.IPs[0])
	}

	latency := network.Latency(network.IPs[0], network.IPs[1])
	if latency == 0 {
		t.Errorf("latency is 0, want non-zero")
	}
	t.Logf("latency is %v\n", latency)
}

func TestAllLatencies(t *testing.T) {
	for i, a := range network.IPs {
		for j, b := range network.IPs {
			if a == b {
				continue
			}
			latency := network.Latency(a, b)
			if latency == 0 {
				t.Fatalf("latency is 0, want non-zero. Idx, i j %d %d", i, j)
			}
		}
	}
}
