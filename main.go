package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

func fetchLiveIpRanges(rangesUrl string) (map[string]struct{}, error) {
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

func detectIpRangesChange(rangesUrl string) {
	liveIpRanges, err := fetchLiveIpRanges(rangesUrl)
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
		log.Printf("cloudflare ip ranges have changed, go check https://www.cloudflare.com/ips-v4")
	} else {
		cloudflareIpRangesChanged.Set(0)
		log.Printf("cloudflare ip ranges have NOT changed.")
	}

}

func main() {
	portStr := flag.String("port", "2541", "Port to expose metrics on")
	url := flag.String("url", "https://www.cloudflare.com/ips-v4", "URL to fetch Cloudflare IP ranges from")
	intervalStr := flag.String("interval", "6", "Interval (in hours) to check for IP range changes")
	flag.Parse()
	exitApp := false

	port, err := strconv.Atoi(*portStr)
	if err != nil || port < 1 || port > 65535 {
		log.Printf("[ERROR] invalid port: %s (must be between 1 and 65535)", *portStr)
		exitApp = true
	}
	intervalHours, err := strconv.Atoi(*intervalStr)
	if err != nil || intervalHours <= 0 {
		log.Printf("[ERROR] invalid interval: %s (must be a positive integer)", *intervalStr)
		exitApp = true
	}

	if exitApp {
		os.Exit(1)
	}
	log.Printf("####### start configuration values #######")
	log.Printf("Port: %d", port)
	log.Printf("URL: %s", *url)
	log.Printf("Interval: %d hour(s)", intervalHours)
	log.Printf("####### end configuration values #######")

	interval := time.Duration(intervalHours) * time.Hour

	go func() {
		for {
			log.Printf("detecting changes..")
			detectIpRangesChange(*url)
			log.Printf("sleeping for %d hour(s)", intervalHours)
			time.Sleep(interval)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf("0.0.0.0:%s", *portStr)
	log.Printf("starting metrics server on %s\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("[ERROR] failed to start server: %v", err)
		os.Exit(1)
	}
}
