package simlatencies

import (
	"bufio"
	"compress/gzip"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strconv"
	"time"
)

func NewNetworkFromFilepaths(maskedIPsGZFile, pairwisePredictionsGZFile string) (Network, error) {
	ipsFile, err := os.Open(maskedIPsGZFile)
	if err != nil {
		return Network{}, err
	}
	ipReader, err := gzip.NewReader(ipsFile)
	if err != nil {
		return Network{}, fmt.Errorf("creating gzip reader for masked-ips.txt.gz: %w", err)
	}

	latenciesFile, err := os.Open(pairwisePredictionsGZFile)
	if err != nil {
		return Network{}, err
	}
	latencies, err := gzip.NewReader(latenciesFile)
	if err != nil {
		panic(fmt.Errorf("creating gzip reader for pairwise_predictions.csv.gz: %w", err))
	}
	return NewNetworkFromReaders(ipReader, latencies)
}

func NewNetworkFromReaders(maskedIPsReader, latenciesReader io.Reader) (Network, error) {
	ips := make([]netip.Addr, 0, 7000)

	scanner := bufio.NewScanner(maskedIPsReader)
	var i int
	for scanner.Scan() {
		ipString := scanner.Text()
		ip, err := netip.ParseAddr(ipString)
		if err != nil {
			return Network{}, fmt.Errorf("parsing IP %q: %w", ipString, err)
		}
		ips = append(ips, ip)
		i++
	}
	if err := scanner.Err(); err != nil {
		return Network{}, fmt.Errorf("scanning masked-ips.txt.gz: %w", err)
	}

	latencies := make([][]time.Duration, len(ips))
	for i := range latencies {
		latencies[i] = make([]time.Duration, len(ips))
	}

	r := csv.NewReader(latenciesReader)
	if _, err := r.Read(); err != nil {
		return Network{}, fmt.Errorf("reading pairwise_predictions.csv.gz header: %w", err)
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return Network{}, fmt.Errorf("reading pairwise_predictions.csv.gz record: %w", err)
		}
		srcID, err := strconv.Atoi(rec[0])
		if err != nil {
			return Network{}, fmt.Errorf("parsing src ID %q: %w", rec[0], err)
		}
		destID, err := strconv.Atoi(rec[1])
		if err != nil {
			return Network{}, fmt.Errorf("parsing dest ID %q: %w", rec[1], err)
		}
		latencyMS, err := strconv.ParseFloat(rec[2], 64)
		if err != nil {
			return Network{}, fmt.Errorf("parsing latency %q: %w", rec[2], err)
		}
		latencies[srcID][destID] = time.Duration(latencyMS * float64(time.Millisecond))
	}

	return Network{
		IPs:       ips,
		Latencies: latencies,
	}, nil
}

type Network struct {
	IPs       []netip.Addr
	Latencies [][]time.Duration
}

func (n *Network) Latency(src, dest netip.Addr) time.Duration {
	src4 := src.As4()
	srcID := int(src4[2])*256 + int(src4[3])
	dest4 := dest.As4()
	destID := int(dest4[2])*256 + int(dest4[3])
	if srcID == destID {
		return 0
	}
	if srcID > destID {
		// Latencies are only stored for the smaller index to the larger index.
		srcID, destID = destID, srcID
	}
	return n.Latencies[srcID][destID]
}
