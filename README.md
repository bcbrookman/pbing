# pbing

An enhanced ping utility written in Go using the [prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing) library.

## Features

- **Timestamps** so you don't have to watch the clock
- **RTT color-coding**  to highlight latency and jitter
- **Time deltas** between responses so you don't have to do math

## Quickstart

1. Install using Go

   ```shell
   go install github.com/bcbrookman/pbing@latest
   ```

2. Run `pbing` like you would `ping`.

   ```shell
   $ pbing github.com
   PING github.com (X.X.X.X):
   2025-04-13 20:50:39 (Δ0.0s): 32 bytes from X.X.X.X: icmp_seq=0 time=42.479811ms ttl=52
   2025-04-13 20:50:40 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=1 time=43.59449ms ttl=52
   2025-04-13 20:50:41 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=2 time=41.741389ms ttl=52
   2025-04-13 20:50:42 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=3 time=42.7836ms ttl=52
   2025-04-13 20:50:43 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=4 time=41.19848ms ttl=52
   2025-04-13 20:50:44 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=5 time=40.197365ms ttl=52
   2025-04-13 20:50:45 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=6 time=41.741094ms ttl=52
   2025-04-13 20:50:46 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=7 time=39.466592ms ttl=52
   2025-04-13 20:50:47 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=8 time=42.702654ms ttl=52
   2025-04-13 20:50:48 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=9 time=40.856067ms ttl=52
   2025-04-13 20:50:49 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=10 time=42.968379ms ttl=52
   2025-04-13 20:50:50 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=11 time=43.368166ms ttl=52
   2025-04-13 20:50:51 (Δ1.0s): 32 bytes from X.X.X.X: icmp_seq=12 time=43.300946ms ttl=52
   ^C
   --- github.com ping statistics ---
   13 packets transmitted, 13 packets received, 0 duplicates, 0% packet loss
   round-trip min/avg/max/stddev = 39.466592ms/42.030695ms/43.59449ms/1.244149ms
   ```

## Manual Install

1. Download a binary from [releases](https://github.com/bcbrookman/pbing/releases/) that matches your OS and arch.
2. Move the downloaded binary to a location in your PATH.

   For example, on most Linux systems:

   ```shell
   sudo mv ~/Downloads/pbing_* /usr/local/bin/pbing
   sudo chmod +x /usr/local/bin/pbing
   ```
