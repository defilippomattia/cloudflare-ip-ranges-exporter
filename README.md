# Introduction

Prometheus exporter that monitors Cloudflare's published IPv4 address ranges (https://www.cloudflare.com/ips-v4/) and detects changes compared to a hardcoded baseline. It periodically fetches the current IP list from Cloudflare's public URL and exposes a Prometheus metric (cloudflare_ip_ranges_changed) indicating whether the list has changed.

This can be useful if you have a server which allows only connections coming from Cloudflare (e.g in your firewall rich rules) and you want to be informed if there were some changes in the IP ranges.

# How to run

(wip, check issues)

```
go run main.go
http://localhost:2112/metrics
```