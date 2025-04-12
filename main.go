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

var examples = `
Examples:
    # ping google continuously
    ping www.google.com

    # ping google 5 times
    ping -c 5 www.google.com

    # ping google 5 times at 500ms intervals
    ping -c 5 -i 500ms www.google.com

    # ping google for 10 seconds
    ping -t 10s www.google.com

    # ping google specified interface
    ping -I eth1 www.google.com

    # Send a privileged raw ICMP ping
    sudo ping --privileged www.google.com

    # Send ICMP messages with a 100-byte payload
    ping -s 100 1.1.1.1

    # Send ICMP messages with DSCP CS4 and ECN bits set to 0
    ping -Q 128 8.8.8.8

    # ping multiple hosts simultaneously
    ping www.google.com gmail.com
`

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
	timeout := flag.Duration("t", time.Second*100000, "")
	interval := flag.Duration("i", time.Second, "")
	count := flag.Int("c", -1, "")
	size := flag.Int("s", 24, "")
	ttl := flag.Int("l", 64, "TTL")
	iface := flag.String("I", "", "interface name")
	tclass := flag.Int("Q", 192, "Set Quality of Service related bits in ICMP datagrams (DSCP + ECN bits). Only decimal number supported")
	privileged := flag.Bool("privileged", false, "")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		_, err := fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		flag.PrintDefaults()
		_, err = fmt.Fprint(out, examples)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	host := flag.Arg(0)

	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Println("ERROR:", err)
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
