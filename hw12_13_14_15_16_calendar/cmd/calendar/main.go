package main

import (
	"flag"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/app"
)

var (
	configFile string
	migrate    bool
)

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
	flag.BoolVar(&migrate, "migrate", false, "Migrate DB")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	app.Run(configFile, migrate)
}
