# pbing

An enhanced ping utility written in Go using the [prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing) library.

## Features

- **Timestamps** so you don't have to watch the clock
- **RTT color-coding**  to highlight latency and jitter
- **Time deltas** between responses so you don't have to do math

## Installation

Installation requires the following general steps:

1. Downloading a binary from [releases](https://github.com/bcbrookman/pbing/releases/) matching your OS and arch.
2. Marking the binary as executable.
3. Optionally, moving the binary to a location in your system's PATH.

You can install the latest version using the **Quickstart** commands below.

<details open>
<summary><strong>Linux Quickstart</strong></summary>

Run the following in a terminal:

```shell
curl -Lo /tmp/pbing https://github.com/bcbrookman/pbing/releases/latest/download/pbing_$(uname -s)_$(uname -m)
sudo chmod +x /tmp/pbing
sudo mv /tmp/pbing /usr/local/bin/pbing
pbing -h
```

</details>

<details open>
<summary><strong>Windows Quickstart</strong></summary>

Run the following in PowerShell as an Administrator:

```powershell
New-Item -Path "$env:ProgramFiles\pbing" -ItemType Directory -Force
Invoke-WebRequest `
   -Uri "https://github.com/bcbrookman/pbing/releases/latest/download/pbing_Windows_x86_64.exe" `
   -OutFile "$env:ProgramFiles\pbing\pbing.exe"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:ProgramFiles\pbing", [EnvironmentVariableTarget]::Machine)
pbing -h
```

</details>

> [!IMPORTANT]
> - On **Linux**, it may be necessary to set `sysctl -w net.ipv4.ping_group_range="0 2147483647"` to allow "unprivileged" pings via UDP.
> - On **Windows**, you must **ALWAYS** use `-privileged`, and accessing TTL values is not supported due to low-level package limitations.
> - MacOS is not officially supported.
> - For more information, see [pro-bing's supported operating systems](https://github.com/prometheus-community/pro-bing?tab=readme-ov-file#supported-operating-systems).

## Usage

Use `pbing` as you would `ping`.

```text
$ pbing -h

Usage:
  pbing [options] <destination>

Options:
  <destination>
        dns name or ip address to ping
  -h, -help
        print this help message
  -I interface
        interface name to source pings from
  -Q tclass
        QoS tclass (DSCP + ECN bits) as a decimal number (default 192)
  -T time
        maximum time to ping before exiting (default 27h46m40s)
  -V version
        print version and exit
  -c count
        maximum count of pings before exiting (default -1)
  -i interval
        time interval between pings (default 1s)
  -privileged
        enable privileged mode to send raw ICMP rather than UDP
  -s size
        payload size in bytes (default 24)
  -t TTL
        time to live (TTL) value (default 64)

Examples:
  pbing example.com                    # ping continuously
  pbing -c 5 example.com               # ping 5 times
  pbing -c 5 -i 500ms example.com      # ping 5 times at 500ms intervals
  pbing -T 10s example.com             # ping for 10 seconds
  pbing -I eth0 example.com            # ping from a specific interface
  sudo pbing -privileged example.com   # ping using raw ICMP pings
  pbing -s 100 example.com             # ping with 100-byte payloads
  pbing -Q 128 example.com             # ping with DSCP CS4 and ECN 0

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
