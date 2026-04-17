package simlatencies_test

import (
	"net/netip"
	"testing"

	simlatencies "github.com/marcopolo/ethereum-network-latencies"
)

func init() {
	simlatencies.MustInit()
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
