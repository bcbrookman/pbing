package main

import (
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	probing "github.com/prometheus-community/pro-bing"
)

func TestColorizeRTT(t *testing.T) {
	// Bypass check for non-tty output streams
	color.NoColor = false

	stddev := time.Millisecond * 5
	avg := time.Millisecond * 50

	tests := []struct {
		name        string
		stats       *probing.Statistics
		pktrtt      time.Duration
		expectColor string // "", "green", "red"
	}{
		{
			name: "RTT significantly lower than avg",
			stats: &probing.Statistics{
				AvgRtt:    avg,
				StdDevRtt: stddev,
			},
			pktrtt:      time.Millisecond * 40,
			expectColor: "green",
		},
		{
			name: "RTT significantly higher than avg",
			stats: &probing.Statistics{
				AvgRtt:    avg,
				StdDevRtt: stddev,
			},
			pktrtt:      time.Millisecond * 60,
			expectColor: "red",
		},
		{
			name: "RTT within stddev range (no color)",
			stats: &probing.Statistics{
				AvgRtt:    avg,
				StdDevRtt: stddev,
			},
			pktrtt:      time.Millisecond * 53,
			expectColor: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colored := ColorizeRTT(tt.stats, tt.pktrtt)

			hasGreen := strings.Contains(colored, "\x1b[32m")
			hasRed := strings.Contains(colored, "\x1b[31m")

			switch tt.expectColor {
			case "green":
				if !hasGreen {
					t.Errorf("expected green color in output: %q", colored)
				}
			case "red":
				if !hasRed {
					t.Errorf("expected red color in output: %q", colored)
				}
			case "":
				if hasGreen || hasRed {
					t.Errorf("expected no color in output: %q", colored)
				}
			}
		})
	}
}

func TestColorizePacketDelta(t *testing.T) {
	// Bypass check for non-tty output streams
	color.NoColor = false

	interval := time.Second

	tests := []struct {
		name        string
		delta       time.Duration
		expectColor string // "", "yellow", "red"
	}{
		{"No color (below 2x)", time.Second + 1, ""},
		{"Yellow (>= 2x)", time.Second * 2, "yellow"},
		{"Red (>= 3x)", time.Second * 3, "red"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizePacketDelta(interval, tt.delta)

			hasYellow := strings.Contains(result, "\x1b[33m")
			hasRed := strings.Contains(result, "\x1b[31m")

			switch tt.expectColor {
			case "yellow":
				if !hasYellow {
					t.Errorf("expected yellow, got %q", result)
				}
			case "red":
				if !hasRed {
					t.Errorf("expected red, got %q", result)
				}
			case "":
				if hasYellow || hasRed {
					t.Errorf("expected no color, got %q", result)
				}
			default:
				t.Fatalf("unknown expected color: %q", tt.expectColor)
			}
		})
	}
}
