package main

import (
	"errors"
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

// setupLogs populates the global logger and logLevel variables. It
// also chooses an appropriate log level based on runtime attributes.
func setupLogs(verboseLogs *bool) {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelInfo)

	if *verboseLogs {
		// AddSource: true here
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logLevel.Set(slog.LevelDebug)
	}
}

// run is the real main, but one where we can exit with an error.
func run(args []string) error {
	myFlags := flag.NewFlagSet("myFlags", flag.ExitOnError)

	var verboseLogs = myFlags.Bool("v", false, "Verbose logging")
	var filename = myFlags.String("file", "", "Road Trip vehicle CSV file")

	err := myFlags.Parse(args[1:])
	if err != nil {
		return err
	}

	setupLogs(verboseLogs)

	if *filename == "" {
		return errors.New("no filename provided (--file)")
	}

	//
	// Here's where the interesting stuff starts as far as this library is concerned.
	//

	// The roadtrip package will happily use your log/slog Logger if you have one.
	options := roadtrip.VehicleOptions{Logger: logger}

	// Create a [roadtrip.Vehicle] object with contents from a Road Trip data file.
	vehicle, err := roadtrip.NewVehicleFromFile(*filename, options)
	if err != nil {
		return err
	}

	logger.Info("Loaded vehicle", "vehicle", vehicle)
	fmt.Printf("-- \n\n")
	fmt.Printf("%s\n\n", vehicle.Vehicles[0].Name)

	var (
		totalFuelCost float64
		totalUnits    float64
	)

	for i, f := range vehicle.FuelRecords {
		logger.Debug("Fuel Record",
			"index", i,
			"fuel", f,
		)

		totalFuelCost += f.TotalPrice
		totalUnits += f.FillAmount
	}

	startOdometer := vehicle.FuelRecords[0].Odometer
	endOdometer := vehicle.FuelRecords[len(vehicle.FuelRecords)-1].Odometer
	totalMiles := endOdometer - startOdometer

	fmt.Printf(" * Drove %.0f miles averaging %0.02f mpg\n",
		totalMiles,
		(totalMiles / totalUnits),
	)

	fmt.Printf(" * Spent $%0.02f on %0.0f gallons of fuel in %d fillups\n",
		totalFuelCost,
		totalUnits,
		len(vehicle.FuelRecords),
	)

	fmt.Printf("\n")

	return nil
}

// main does as little as we can get away with.
func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
