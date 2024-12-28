package main

import (
	"flag"
	"fmt"
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
		slog.Error("no filename provided (--file)")
		os.Exit(1)
	}

	vehicle, err := roadtrip.NewFromFile(*filename)
	if err != nil {
		slog.Error("unable to load CSV file", "error", err)
		os.Exit(1)
	}

	// slog.Info("Loaded vehicle", "vehicle", vehicle)

	totalFuelCost := 0.00

	for i, f := range vehicle.FuelRecords {
		slog.Debug("processing FuelRecord", "i", i, "f", f)
		totalFuelCost += f.TotalPrice
	}

	fmt.Printf("Spent %0.02f on fuel in %d fillups\n", totalFuelCost, len(vehicle.FuelRecords))
}
