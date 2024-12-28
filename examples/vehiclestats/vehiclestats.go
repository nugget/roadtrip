package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/nugget/roadtrip-go/roadtrip"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var debugMode = flag.Bool("v", false, "Verbose logging")
	var filename = flag.String("file", "", "Road Trip vehicle CSV file")

	flag.Parse()

	if *filename == "" {
		slog.Error("no filename provided (--file)")
		os.Exit(1)
	}

	options := roadtrip.VehicleOptions{
		Logger:   logger,
		LogLevel: slog.LevelInfo,
	}

	if *debugMode {
		// AddSource: true here
		options.LogLevel = slog.LevelDebug
	}

	vehicle, err := roadtrip.NewVehicleFromFile(*filename, options)
	if err != nil {
		slog.Error("unable to load CSV file", "error", err)
		os.Exit(1)
	}

	// logger.Info("Loaded vehicle", "vehicle", vehicle)

	totalFuelCost := 0.00

	for i, f := range vehicle.FuelRecords {
		totalFuelCost += f.TotalPrice
	}

	fmt.Printf("Spent %0.02f on fuel in %d fillups\n", totalFuelCost, len(vehicle.FuelRecords))
}
