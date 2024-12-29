package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/nugget/roadtrip-go/roadtrip"
)

var (
	logger   *slog.Logger
	logLevel *slog.LevelVar
)

func setupLogs() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func main() {
	setupLogs()

	var debugMode = flag.Bool("v", false, "Verbose logging")
	var filename = flag.String("file", "", "Road Trip vehicle CSV file")

	flag.Parse()

	if *filename == "" {
		logger.Error("no filename provided (--file)")
		os.Exit(1)
	}

	options := roadtrip.VehicleOptions{
		Logger: logger,
	}

	if *debugMode {
		// AddSource: true here
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logLevel.Set(slog.LevelDebug)
	}

	// Create a [roadtrip.Vehicle] object with contents from a Road Trip data file.
	vehicle, err := roadtrip.NewVehicleFromFile(*filename, options)
	if err != nil {
		logger.Error("unable to load Road Trip data file",
			"error", err,
		)
		os.Exit(1)
	}

	logger.Info("Loaded vehicle", "vehicle", vehicle)

	totalFuelCost := 0.00

	for i, f := range vehicle.FuelRecords {
		logger.Debug("Fuel Record",
			"index", i,
			"fuel", f,
		)
		totalFuelCost += f.TotalPrice
	}

	fmt.Printf("%s\n", vehicle.Vehicles[0].Name)
	fmt.Printf("Spent %0.02f on fuel in %d fillups\n", totalFuelCost, len(vehicle.FuelRecords))
}
