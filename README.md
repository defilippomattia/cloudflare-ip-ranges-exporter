[![ci](https://github.com/defilippomattia/cloudflare-ip-ranges-exporter/actions/workflows/cicd.yml/badge.svg?branch=main)](https://github.com/defilippomattia/cloudflare-ip-ranges-exporter/actions/workflows/cicd.yml)

# Introduction

Prometheus exporter that monitors Cloudflare's published IPv4 address ranges (https://www.cloudflare.com/ips-v4/) and detects changes compared to a hardcoded baseline. It periodically fetches the current IP list from Cloudflare's public URL and exposes a Prometheus metric (cloudflare_ip_ranges_changed) indicating whether the list has changed.

This can be useful if you have a server which allows only connections coming from Cloudflare (e.g in your firewall rich rules) and you want to be informed if there were some changes in the IP ranges.

# How to run
(wip, check issues)

```
go run main.go
If you want to specify different config params:
go run main.go --port=2541 --url=https://www.cloudflare.com/ips-v4 --interval=6
http://localhost:2541/metrics
```

## Flags
| Flag         | Type   | Default Value                      | Description                                                                 |
|--------------|--------|------------------------------------|-----------------------------------------------------------------------------|
| `--port`     | string | `2541`                             | The port number on which the HTTP metrics server will listen. Must be between `1` and `65535`. |
| `--url`      | string | `https://www.cloudflare.com/ips-v4`| The URL to fetch the list of Cloudflare IPv4 ranges from. Must be a valid HTTP(S) URL (this is not meant to be changed, this flag was added so I can pass some local server URL which simulates cloudflares response for testing. Check Mocking Cloudflare response chapter.) |
| `--interval` | string | `6`                                | Interval in hours to check for changes in the IP ranges. Must be a positive integer. |


## Mocking Cloudflare response

For testing purposes, you can mock the repsonse by adding a text file in the `cf-mock` folder and running 

`python -m http.server 8081`

Visiting http://localhost:8081/mock-1.txt will show you the content of the `mock-1.txt` file.