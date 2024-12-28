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

// A Vehicle holds the parsed sections contained in a Road Trip CSV backup file.
type Vehicle struct {
	Delimiters         string
	Version            int
	Language           string
	Filename           string
	Vehicles           []VehicleRecord     `section:"VEHICLE"`
	FuelRecords        []FuelRecord        `section:"FUEL RECORDS"`
	MaintenanceRecords []MaintenanceRecord `section:"MAINTENANCE RECORDS"`
	Trips              []TripRecord        `section:"ROAD TRIPS"`
	Tires              []TireRecord        `section:"TIRE LOG"`
	Valuations         []ValuationRecord   `section:"VALUATIONS"`
	Raw                []byte
}

// NewFromFile returns a new [Vehicle] populated with data read and parsed
// from the file.
func NewFromFile(filename string) (Vehicle, error) {
	var v Vehicle

	err := v.LoadFile(filename)
	if err != nil {
		return v, err
	}

	return v, nil
}

// LoadFile reads and parses a file into a [Vehicle] variable.
func (v *Vehicle) LoadFile(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	v.Filename = filename

	if true {
		// Remove erroneous header fields for VEHICLE section
		// per Darren Stone 9-Dec-2024 via email
		omitHeaders := []byte(",Tank 1 Type,Tank 2 Type,Tank 2 Units")
		v.Raw = bytes.Replace(buf, omitHeaders, []byte{}, 1)
	} else {
		v.Raw = buf
	}

	err = v.Parse("FUEL RECORDS", &v.FuelRecords)
	if err != nil {
		return fmt.Errorf("unable to parse FuelRecords: %w", err)
	}

	err = v.Parse("MAINTENANCE RECORDS", &v.MaintenanceRecords)
	if err != nil {
		return fmt.Errorf("MaintenanceRecords: %w", err)
	}

	err = v.Parse("ROAD TRIPS", &v.Trips)
	if err != nil {
		return fmt.Errorf("unable to parse Trips: %w", err)
	}

	err = v.Parse("VEHICLE", &v.Vehicles)
	if err != nil {
		return fmt.Errorf("unable to parse Vehicle: %w", err)
	}

	err = v.Parse("TIRE LOG", &v.Tires)
	if err != nil {
		return fmt.Errorf("unable to parse TireLogs: %w", err)
	}

	err = v.Parse("VALUATIONS", &v.Valuations)
	if err != nil {
		return fmt.Errorf("unable to parse Valuations: %w", err)
	}

	slog.Debug("Loaded Road Trip CSV",
		"filename", v.Filename,
		"bytes", len(buf),
		"vehicleRecords", len(v.Vehicles),
		"fuelRecords", len(v.FuelRecords),
		"mainteanceRecords", len(v.MaintenanceRecords),
		"Trips", len(v.Trips),
		"tireLogs", len(v.Tires),
		"valuations", len(v.Valuations),
	)

	return nil
}

// Section returns a byte slice containing the raw contents of the specified section
// from the corresponding [CSV] object.
func (v *Vehicle) Section(sectionHeader string) []byte {
	slog.Debug("Fetching Section from Raw",
		"sectionHeader", sectionHeader,
	)

	sectionStart := make(map[string]int)

	for index, element := range SectionHeaders {
		i := bytes.Index(v.Raw, []byte(SectionHeaders[index]))
		sectionStart[element] = i

		slog.Debug("Section Start detected",
			"element", element,
			"sectionStart", i,
		)
	}

	startPosition := sectionStart[sectionHeader]
	endPosition := len(v.Raw)

	for _, e := range sectionStart {
		if e > startPosition && e < endPosition {
			endPosition = e - 1
		}
	}

	// Don't include the section header line in the outbuf
	startPosition = startPosition + len(sectionHeader) + 1

	outbuf := v.Raw[startPosition:endPosition]

	slog.Debug("Section Range calculated",
		"sectionHeader", sectionHeader,
		"startPosition", startPosition,
		"endPosition", endPosition,
		"sectionBytes", len(outbuf),
	)

	return outbuf
}

// Parse unmarshalls the raw byte slice of the specified section from the underlying [CSV]
// object and transforms it into the struct used by this package.
func (v *Vehicle) Parse(sectionHeader string, target interface{}) error {
	if _, err := csvlib.Unmarshal(v.Section(sectionHeader), target); err != nil {
		return err
	}

	return nil
}

// ParseDate parses a Road Trip styled date string and turns it into a proper
// Go [time.Time] value.
func ParseDate(dateString string) time.Time {
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
