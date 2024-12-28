// Package roadtrip implements utility routines for reading the CSV backup
// files created by the iOS Road Trip MPG application.
package roadtrip

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/tiendc/go-csvlib"
)

// A CSV holds the parsed sections contained in a Road Trip CSV backup file.
type CSV struct {
	Delimiters         string
	Version            int
	Language           string
	Filename           string
	Vehicle            []Vehicle
	FuelRecords        []Fuel
	MaintenanceRecords []Maintenance
	RoadTrips          []RoadTrip
	TireLogs           []Tire
	Valuations         []Valuation
	Raw                []byte
}

// NewFromFile returns a new [CSV] populated with data read and parsed
// from the file.
func NewFromFile(filename string) (CSV, error) {
	var rt CSV

	err := rt.LoadFile(filename)
	if err != nil {
		return rt, err
	}

	return rt, nil
}

// LoadFile reads and parses a file into a [CSV] variable.
func (rt *CSV) LoadFile(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	rt.Filename = filename

	if true {
		// Remove erroneous header fields for VEHICLE section
		// per Darren Stone 9-Dec-2024 via email
		omitHeaders := []byte(",Tank 1 Type,Tank 2 Type,Tank 2 Units")
		rt.Raw = bytes.Replace(buf, omitHeaders, []byte{}, 1)
	} else {
		rt.Raw = buf
	}

	if err := rt.Parse("FUEL RECORDS", &rt.FuelRecords); err != nil {
		return fmt.Errorf("FuelRecords: %w", err)
	}

	if err := rt.Parse("MAINTENANCE RECORDS", &rt.MaintenanceRecords); err != nil {
		return fmt.Errorf("MaintenanceRecords: %w", err)
	}

	if err := rt.Parse("ROAD TRIPS", &rt.RoadTrips); err != nil {
		return fmt.Errorf("RoadTrips: %w", err)
	}

	if err := rt.Parse("VEHICLE", &rt.Vehicle); err != nil {
		return fmt.Errorf("Vehicle: %w", err)
	}

	if err := rt.Parse("TIRE LOG", &rt.TireLogs); err != nil {
		return fmt.Errorf("TireLogs: %w", err)
	}

	if err := rt.Parse("VALUATIONS", &rt.Valuations); err != nil {
		return fmt.Errorf("Valuations: %w", err)
	}

	slog.Info("Loaded Road Trip CSV",
		"filename", rt.Filename,
		"bytes", len(buf),
		"vehicleRecords", len(rt.Vehicle),
		"fuelRecords", len(rt.FuelRecords),
		"mainteanceRecords", len(rt.MaintenanceRecords),
		"roadTrips", len(rt.RoadTrips),
		"tireLogs", len(rt.TireLogs),
		"valuations", len(rt.Valuations),
	)

	return nil
}

// Section returns a byte slice containing the raw contents of the specified section
// from the corresponding [CSV] object.
func (rt *CSV) Section(sectionHeader string) (outbuf []byte) {
	slog.Debug("Fetching Section from Raw",
		"sectionHeader", sectionHeader,
	)

	sectionStart := make(map[string]int)

	for index, element := range Sections {
		i := bytes.Index(rt.Raw, []byte(Sections[index]))
		sectionStart[element] = i

		slog.Debug("Section Start detected",
			"element", element,
			"sectionStart", i,
		)
	}

	startPosition := sectionStart[sectionHeader]
	endPosition := len(rt.Raw)

	for _, e := range sectionStart {
		if e > startPosition && e < endPosition {
			endPosition = e - 1
		}
	}

	// Don't include the section header line in the outbuf
	startPosition = startPosition + len(sectionHeader) + 1

	outbuf = rt.Raw[startPosition:endPosition]

	slog.Debug("Section Range calculated",
		"sectionHeader", sectionHeader,
		"startPosition", startPosition,
		"endPosition", endPosition,
		"sectionBytes", len(outbuf),
	)

	return
}

// Parse unmarshalls the raw byte slice of the specified section from the underlying [CSV]
// object and transforms it into the struct used by this package.
func (rt *CSV) Parse(sectionHeader string, target interface{}) error {
	if _, err := csvlib.Unmarshal(rt.Section(sectionHeader), target); err != nil {
		return err
	}

	return nil
}

// ParseDate parses a Road Trip styled date string and turns it into a proper
// Go [time.Time] value
func ParseDate(dateString string) (t time.Time) {
	t, err := time.Parse("2006-1-2 15:04", dateString)
	if err != nil {
		t, err = time.Parse("2006-1-2", dateString)
		if err != nil {
			slog.Debug("Can't parse Road Trip date string",
				"error", err,
				"dateString", dateString,
			)
		}
	}

	return t
}
