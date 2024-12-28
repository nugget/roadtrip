package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/nugget/roadtrip-go/roadtrip"
)

func main() {
	_ = slog.SetLogLoggerLevel(slog.LevelInfo)

	var debugMode = flag.Bool("v", false, "Verbose logging")
	var filename = flag.String("file", "", "Road Trip vehicle CSV file")

	flag.Parse()

	if *debugMode {
		_ = slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if *filename == "" {
		slog.Error("No filename provided (--file)")
		os.Exit(1)
	}

	vehicle, err := roadtrip.NewFromFile(*filename)
	if err != nil {
		slog.Error("Unable to load CSV file", "error", err)
		os.Exit(1)
	}

	slog.Info("Loaded vehicle", "vehicle", vehicle)
}
