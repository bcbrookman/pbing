package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	probing "github.com/prometheus-community/pro-bing"
)

const helpUsage string = `
Usage:
  pbing [options] <destination>
`

const helpOptions string = `
Options:
  <destination>
        dns name or ip address to ping
  -h, -help
        print this help message
`

const helpExamples string = `
Examples:
  pbing example.com                    # ping continuously
  pbing -c 5 example.com               # ping 5 times
  pbing -c 5 -i 500ms example.com      # ping 5 times at 500ms intervals
  pbing -T 10s example.com             # ping for 10 seconds
  pbing -I eth0 example.com            # ping from a specific interface
  sudo pbing -privileged example.com   # ping using raw ICMP pings
  pbing -s 100 example.com             # ping with 100-byte payloads
  pbing -Q 128 example.com             # ping with DSCP CS4 and ECN 0
`

var Version = "development"

func ColorizeRTT(stats *probing.Statistics, pktrtt time.Duration) string {
	// calculate difference of average RTT and current packet RTT
	avgRttDiff := time.Duration(stats.AvgRtt - pktrtt)

	// determine whether the difference is within current standard deviation
	avgRttDiffWithinStdDev := avgRttDiff.Abs() >= stats.StdDevRtt

	result := pktrtt.String()
	if avgRttDiff >= 0 { // if positive, the pktrtt is lower than the current avgrtt
		if avgRttDiffWithinStdDev {
			result = color.GreenString(pktrtt.String())
		}
	} else { // if negative, the pktrtt is higher than the current avgrtt
		if avgRttDiffWithinStdDev {
			result = color.RedString(pktrtt.String())
		}
	}
	return result
}

func ColorizePacketDelta(interval time.Duration, delta time.Duration) string {
	result := fmt.Sprintf("%0.1fs", delta.Seconds())

	if delta >= interval*3 {
		result = color.RedString(result)
	} else if delta >= interval*2 {
		result = color.YellowString(result)
	}
	return result
}

func main() {
	iface := flag.String("I", "", "`interface` name to source pings from")
	tclass := flag.Int("Q", 192, "QoS `tclass` (DSCP + ECN bits) as a decimal number")
	timeout := flag.Duration("T", time.Second*100000, "maximum `time` to ping before exiting")
	count := flag.Int("c", -1, "maximum `count` of pings before exiting")
	ttl := flag.Int("t", 64, "time to live (`TTL`) value")
	interval := flag.Duration("i", time.Second, "time `interval` between pings")
	privileged := flag.Bool("privileged", false, "enable privileged mode to send raw ICMP rather than UDP")
	size := flag.Int("s", 24, "payload `size` in bytes")
	version := flag.Bool("V", false, "print `version` and exit")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		_, err := fmt.Fprint(out, helpUsage)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		_, err = fmt.Fprint(out, helpOptions)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		flag.PrintDefaults()
		_, err = fmt.Fprint(out, helpExamples, "\n")
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
	flag.Parse()

	if *version {
		fmt.Println("pbing", Version)
		os.Exit(0)
	}

	host := flag.Arg(0)

	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Println("ERROR:", err)
		fmt.Println("See 'pbing -h' for usage")
		return
	}

	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	var lastSuccessfulPacketTime time.Time

	pinger.OnRecv = func(pkt *probing.Packet) {

		currPacketTimeStamp := time.Now()
		var timeSinceLastSuccessfulPacket time.Duration
		if !lastSuccessfulPacketTime.IsZero() {
			timeSinceLastSuccessfulPacket = currPacketTimeStamp.Sub(lastSuccessfulPacketTime)
		}

		stats := pinger.Statistics()

		fmt.Printf("%s (\u0394%v): %d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
			currPacketTimeStamp.Format(time.DateTime), ColorizePacketDelta(pinger.Interval, timeSinceLastSuccessfulPacket),
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, ColorizeRTT(stats, pkt.Rtt), pkt.TTL)

		lastSuccessfulPacketTime = currPacketTimeStamp
	}
	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {

		currPacketTimeStamp := time.Now()
		var timeSinceLastSuccessfulPacket time.Duration
		if !lastSuccessfulPacketTime.IsZero() {
			timeSinceLastSuccessfulPacket = currPacketTimeStamp.Sub(lastSuccessfulPacketTime)
		}

		stats := pinger.Statistics()

		fmt.Printf("%v (\u0394%v): %d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			currPacketTimeStamp.Format(time.DateTime), ColorizePacketDelta(pinger.Interval, timeSinceLastSuccessfulPacket),
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, ColorizeRTT(stats, pkt.Rtt), pkt.TTL)

		lastSuccessfulPacketTime = currPacketTimeStamp
	}
	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = *count
	pinger.Size = *size
	pinger.Interval = *interval
	pinger.Timeout = *timeout
	pinger.TTL = *ttl
	pinger.InterfaceName = *iface
	pinger.SetPrivileged(*privileged)
	pinger.SetTrafficClass(uint8(*tclass))

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		fmt.Println("Failed to ping target host:", err)
	}
}
