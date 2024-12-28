// Package roadtrip implements utility routines for reading the CSV backup
// files created by the iOS Road Trip MPG application.
package roadtrip

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"

	cvslib "github.com/tiendc/go-csvlib"
)

const (
	// Remove erroneous header fields for VEHICLE section
	// per Darren Stone 2024-12-09 via email.
	RemoveErroneousHeaders = true
)

// VehicleOptions contain the options to be used when creating a new Vehicle object.
type VehicleOptions struct {
	Logger   *slog.Logger
	LogLevel slog.Level
}

// A Vehicle holds the parsed sections contained in a Road Trip CSV backup file.
type Vehicle struct {
	Delimiters         string
	Version            int
	Language           string
	Filename           string
	Vehicles           []VehicleRecord     `roadtrip:"VEHICLE"`
	FuelRecords        []FuelRecord        `roadtrip:"FUEL RECORDS"`
	MaintenanceRecords []MaintenanceRecord `roadtrip:"MAINTENANCE RECORDS"`
	Trips              []TripRecord        `roadtrip:"ROAD TRIPS"`
	Tires              []TireRecord        `roadtrip:"TIRE LOG"`
	Valuations         []ValuationRecord   `roadtrip:"VALUATIONS"`
	Raw                []byte
	logger             *slog.Logger
	logLevel           slog.Level
}

func UnmarshalRoadtripSection(data []byte, target any) error {
	header, err := SectionHeader(target)
	if err != nil {
		return err
	}

	sectionData := GetSectionContents(data, header)

	_, err = cvslib.Unmarshal(sectionData, target)
	if err != nil {
		return err
	}

	return nil
}

// Each Road Trip "CSV" file is actually multiple, independent blocks of CSV
// data delimited by two newlines and a section header string in all capital
// letters.
//
// SectionHeaderList returns a slice of strings corresponding to each of the
// section headers found in the Road Trip data file. Currently this package
// only supports Language "en" (see known issues in the README.md file).
//
// This list is built by inspecting the `roadtrip` struct tags present in
// the [Vehicle] struct definition.
func SectionHeaderList() []string {
	var headerList []string

	vt := reflect.TypeOf(Vehicle{})
	for i := range vt.NumField() {
		field := vt.Field(i)
		sectionHeader, ok := field.Tag.Lookup("roadtrip")
		if ok {
			headerList = append(headerList, sectionHeader)
		}
	}

	return headerList
}

// SectionHeader will return the section header for any suitable target field
// in the [Vehicle] struct. It's used to identify the correct CSV block in the
// Road Trip CSV file.
func SectionHeader(target any) (string, error) {
	targetType := reflect.TypeOf(target).Elem()

	vt := reflect.TypeOf(Vehicle{})
	for i := range vt.NumField() {
		field := vt.Field(i)

		sectionHeader, ok := field.Tag.Lookup("roadtrip")

		if ok && field.Type == targetType {
			return sectionHeader, nil
		}
	}

	return "", fmt.Errorf("cannot unmarshal %s, missing roadtrip struct tag", targetType)
}

// New returns a new, empty [Vehicle] with a no-op logger.
func NewVehicle(options VehicleOptions) Vehicle {
	var v Vehicle

	if options.Logger == nil {
		options.Logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	v.logger = options.Logger
	v.logLevel = options.LogLevel

	return v
}

// NewFromFile returns a new [Vehicle] populated with data read and parsed
// from the file.
func NewVehicleFromFile(filename string, options VehicleOptions) (Vehicle, error) {
	v := NewVehicle(options)

	err := v.LoadFile(filename)
	if err != nil {
		return v, err
	}

	return v, nil
}

// SetLogger optionally binds an [slog.Logger] to a [Vehicle] for internal
// package debugging. If you do not call SetLogger, log output will be
// discarded during package operation.
func (v *Vehicle) SetLogger(l *slog.Logger) {
	v.logger = l
	v.logLevel = slog.LevelInfo
}

// SetLogLoggerLevel optionally sets the [Vehicle] logger level for internal
// package debugging.
func (v *Vehicle) SetLogLoggerLevel(levelInfo slog.Level) slog.Level {
	v.logLevel = levelInfo
	return slog.SetLogLoggerLevel(levelInfo)
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [Vehicle] object when logging.
func (v *Vehicle) LogValue() slog.Value {
	var value slog.Value

	if len(v.Vehicles) == 1 {
		if v.logLevel > -slog.LevelInfo {
			value = slog.GroupValue(
				slog.String("name", v.Vehicles[0].Name),
			)
		} else {
			value = slog.GroupValue(
				slog.String("name", v.Vehicles[0].Name),
				slog.Int("version", v.Version),
				slog.String("filename", v.Filename),
			)
		}
	}

	return value
}

// LoadFile reads and parses a file into a [Vehicle] variable.
func (v *Vehicle) LoadFile(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	v.Filename = filename

	if RemoveErroneousHeaders {
		omitHeaders := []byte(",Tank 1 Type,Tank 2 Type,Tank 2 Units")
		buf = bytes.Replace(buf, omitHeaders, []byte{}, 1)
	}

	return v.UnmarshalRoadtrip(buf)
}

func (v *Vehicle) UnmarshalRoadtrip(data []byte) error {
	v.Raw = data

	var err error

	var targets []any
	targets = append(targets, &v.Vehicles)
	targets = append(targets, &v.FuelRecords)
	targets = append(targets, &v.MaintenanceRecords)
	targets = append(targets, &v.Trips)
	targets = append(targets, &v.Tires)
	targets = append(targets, &v.Valuations)

	for _, target := range targets {
		err = UnmarshalRoadtripSection(data, target)
		if err != nil {
			return fmt.Errorf("unable to parse %s: %w", target, err)
		}
	}

	v.logger.Info("Loaded Road Trip CSV",
		"filename", v.Filename,
		"bytes", len(data),
		"vehicleRecords", len(v.Vehicles),
		"fuelRecords", len(v.FuelRecords),
		"mainteanceRecords", len(v.MaintenanceRecords),
		"trips", len(v.Trips),
		"tireLogs", len(v.Tires),
		"valuations", len(v.Valuations),
	)

	return nil
}

// Section returns a byte slice containing the raw contents of the specified section
// from the corresponding [Vehicle.Raw] object.
func GetSectionContents(data []byte, sectionHeader string) []byte {
	sectionStart := make(map[string]int)

	for _, element := range SectionHeaderList() {
		i := bytes.Index(data, []byte(element))
		sectionStart[element] = i
	}

	startPosition := sectionStart[sectionHeader]
	endPosition := len(data)

	for _, e := range sectionStart {
		if e > startPosition && e < endPosition {
			endPosition = e - 1
		}
	}

	// Don't include the section header line in the outbuf
	startPosition = startPosition + len(sectionHeader) + 1

	outbuf := data[startPosition:endPosition]

	return outbuf
}

// ParseDate parses a Road Trip styled date string and turns it into a proper
// Go [time.Time] value.
func ParseDate(dateString string) (time.Time, error) {
	t, err := time.Parse("2006-1-2 15:04", dateString)
	if err != nil {
		t, err = time.Parse("2006-1-2", dateString)
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to parse date '%s': %w", dateString, err)
		}
	}

	return t, nil
}
