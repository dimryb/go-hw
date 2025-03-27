package main

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	total       int64
	current     int64
	lastPercent int
}

func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total: total,
	}
}

func formatProgressBar(percent int) string {
	totalWidth := 50
	filled := percent * totalWidth / 100
	arrow := ""
	if filled < totalWidth {
		arrow = ">"
	}
	empty := strings.Repeat("_", 50-filled-len(arrow))
	return fmt.Sprintf("\rProgress: [%s%s%s] %d%%", strings.Repeat("=", filled), arrow, empty, percent)
}

func (pb *ProgressBar) Update(current int64) {
	pb.current = current
	percent := int(float64(pb.current) / float64(pb.total) * 100)

	if percent != pb.lastPercent {
		pb.lastPercent = percent
		bar := formatProgressBar(percent)
		fmt.Print(bar)
	}
}

func (pb *ProgressBar) Finish() {
	fmt.Println("\nProgress completed!")
}
