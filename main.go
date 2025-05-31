package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var hardcodedCloudflareIpRanges = []string{
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
}

var cloudflareIpRangesChanged = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "cloudflare_ip_ranges_changed",
	Help: "1 if there has been a change in https://www.cloudflare.com/ips-v4, 0 otherwise",
})

func fetchLiveIpRanges() (map[string]struct{}, error) {
	rangesUrl := "https://www.cloudflare.com/ips-v4"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(rangesUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	liveIpRanges := make(map[string]struct{})
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" {
			liveIpRanges[ip] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return liveIpRanges, nil
}

func detectIpRangesChange() {
	liveIpRanges, err := fetchLiveIpRanges()
	if err != nil {
		fmt.Println("error fetching live IP ranges:", err)
		return
	}
	changed := false

	if len(liveIpRanges) != len(hardcodedCloudflareIpRanges) {
		changed = true
	} else {
		for _, ip := range hardcodedCloudflareIpRanges {
			_, exists := liveIpRanges[ip]
			if !exists {
				changed = true
				break
			}
		}
	}

	if changed {
		cloudflareIpRangesChanged.Set(1)
		fmt.Println("cloudflare ip ranges have changed, go check https://www.cloudflare.com/ips-v4")
	} else {
		cloudflareIpRangesChanged.Set(0)
		fmt.Println("cloudflare ip ranges have NOT changed.")
	}

}

func main() {
	go func() {
		for {
			detectIpRangesChange()
			time.Sleep(6 * time.Hour)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
