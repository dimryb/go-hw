package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	Run(configFile)
}

func Run(configPath string) {
	cfg, err := config.NewSchedulerConfig(configPath)
	if err != nil {
		log.Fatalf("SchedulerConfig error: %s", err)
	}

	fmt.Println("Scheduler Config:", *cfg)
}
