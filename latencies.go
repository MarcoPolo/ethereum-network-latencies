package simlatencies

import (
	"bufio"
	"embed"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"net/netip"
	"strconv"
	"time"
)

//go:embed masked-ips.txt.gz pairwise_predictions.csv.gz
var data embed.FS

var IPs []netip.Addr

var ipToID map[netip.Addr]int
var latencies [][]time.Duration

var initialized bool

func Init(maskedIPs, pairwisePredictions io.Reader) error {
	IPs = make([]netip.Addr, 0, 7000)
	ipToID = make(map[netip.Addr]int, 7000)

	scanner := bufio.NewScanner(maskedIPs)
	var i int
	for scanner.Scan() {
		ipString := scanner.Text()
		ip, err := netip.ParseAddr(ipString)
		if err != nil {
			return fmt.Errorf("parsing IP %q: %w", ipString, err)
		}
		IPs = append(IPs, ip)
		ipToID[ip] = i
		i++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning masked-ips.txt.gz: %w", err)
	}

	if latencies == nil {
		latencies = make([][]time.Duration, len(IPs))
		for i := range latencies {
			latencies[i] = make([]time.Duration, len(IPs))
		}
	}

	r := csv.NewReader(pairwisePredictions)
	if _, err := r.Read(); err != nil {
		return fmt.Errorf("reading pairwise_predictions.csv.gz header: %w", err)
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("reading pairwise_predictions.csv.gz record: %w", err)
		}
		srcID, err := strconv.Atoi(rec[0])
		if err != nil {
			return fmt.Errorf("parsing src ID %q: %w", rec[0], err)
		}
		destID, err := strconv.Atoi(rec[1])
		if err != nil {
			return fmt.Errorf("parsing dest ID %q: %w", rec[1], err)
		}
		latencyMS, err := strconv.ParseFloat(rec[2], 64)
		if err != nil {
			return fmt.Errorf("parsing latency %q: %w", rec[2], err)
		}
		latencies[srcID][destID] = time.Duration(latencyMS * float64(time.Millisecond))
	}

	initialized = true
	return nil
}

func MustInit(maskedIPs, pairwisePredictions io.Reader) {
	if err := Init(maskedIPs, pairwisePredictions); err != nil {
		panic(err)
	}
}

func Latency(src, dest netip.Addr) (time.Duration, error) {
	if !initialized {
		return 0, fmt.Errorf("latencies not initialized. Must call simlatencies.Init()")
	}

	srcID, ok := ipToID[src]
	if !ok {
		return 0, fmt.Errorf("source IP not found")
	}
	destID, ok := ipToID[dest]
	if !ok {
		return 0, fmt.Errorf("destination IP not found")
	}
	latency := latencies[srcID][destID]
	if latency == 0 {
		return 0, fmt.Errorf("latency not found")
	}
	return latency, nil
}

func MustLatency(src, dest netip.Addr) time.Duration {
	latency, err := Latency(src, dest)
	if err != nil {
		panic(err)
	}
	return latency
}
