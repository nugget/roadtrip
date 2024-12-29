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
	// Supported Road Trip vehicle data file version "1500,en" only.
	SupportedVersion int = 1500

	// Remove erroneous header fields for VEHICLE section
	// per Darren Stone 2024-12-09 via email.
	RemoveErroneousHeaders = true
)

// RawFileData contains the raw contents read from a single Road Trip data
// file.
type RawFileData []byte

// RawSectionData contains the raw contents read from a single section of a
// single Road Trip data file.
type RawSectionData []byte

// VehicleOptions contain the options to be used when creating a new Vehicle object.
type VehicleOptions struct {
	Logger *slog.Logger
}

// A Vehicle holds the parsed sections contained in a Road Trip vehicle data file.
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
	Raw                RawFileData
	logger             *slog.Logger
}

// NewVehicle returns a new, empty [Vehicle] object.
func NewVehicle(options VehicleOptions) Vehicle {
	var v Vehicle

	if options.Logger == nil {
		options.Logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	v.logger = options.Logger

	return v
}

// NewVehicleFromFile returns a new [Vehicle] object populated with data read
// and parsed from the file.
func NewVehicleFromFile(filename string, options VehicleOptions) (Vehicle, error) {
	v := NewVehicle(options)

	err := v.LoadFile(filename)
	if err != nil {
		return v, err
	}

	return v, nil
}

// Each Road Trip "CSV" file is actually multiple, independent blocks of CSV
// data delimited by two newlines and a section header string in all capital
// letters.
//
// SectionHeaderList returns a slice of strings corresponding to each of the
// section headers expected in the Road Trip vehicle data file. Currently this
// package only supports Language "en" (see known issues in the README.md
// file).
//
// This list is built by inspecting the `roadtrip` struct tags present in the
// [Vehicle] struct definition.
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

// SectionHeaderForTarget will return the section header for any suitable
// target field in the [Vehicle] struct. It's used to identify the correct CSV
// block in the Road Trip vehicle data file.
func SectionHeaderForTarget(target any) (string, error) {
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

// GetSectionContents evaluates the raw content from a Road Trip data file and extracts only
// the single section block identified by the supplied section header string value.
func (fileData *RawFileData) GetSectionContents(sectionHeader string) RawSectionData {
	sectionStart := make(map[string]int)

	dataBytes := reflect.ValueOf(*fileData).Bytes()

	for _, element := range SectionHeaderList() {
		i := bytes.Index(dataBytes, []byte(element))
		sectionStart[element] = i
	}

	startPosition := sectionStart[sectionHeader]
	endPosition := len(dataBytes)

	for _, e := range sectionStart {
		if e > startPosition && e < endPosition {
			endPosition = e - 1
		}
	}

	// Don't include the section header line in the outbuf
	startPosition = startPosition + len(sectionHeader) + 1

	outbuf := dataBytes[startPosition:endPosition]

	return outbuf
}

// UnmarshalRoadtripSection takes the raw contents of a Road Trip vehicle data
// file, extracts only the relevant section block, and then parses it into
// appropriate struct field based on the type of the target variable.
//
// This relies on an accurate struct tag on the [Vehicle] field in question
// which instructs the function on which section header line to look for.
func (fileData *RawFileData) UnmarshalRoadtripSection(target any) error {
	header, err := SectionHeaderForTarget(target)
	if err != nil {
		return err
	}

	sectionData := fileData.GetSectionContents(header)

	_, err = cvslib.Unmarshal(sectionData, target)
	if err != nil {
		return err
	}

	return nil
}

// SetLogger optionally sets the [Vehicle] logger for internal package
// debugging.
func (v *Vehicle) SetLogger(l *slog.Logger) {
	v.logger = l
}

// LogValue is the handler for [log.slog] to emit structured output for the
// [Vehicle] object when logging.
func (v *Vehicle) LogValue() slog.Value {
	var value slog.Value

	if len(v.Vehicles) == 1 {
		value = slog.GroupValue(
			slog.String("name", v.Vehicles[0].Name),
			slog.Int("version", v.Version),
			slog.String("filename", v.Filename),
		)
	}

	return value
}

// LoadFile reads and parses a file into the [Vehicle] object.
func (v *Vehicle) LoadFile(filename string) error {
	var buf RawFileData

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

// UnmarshalRoadtrip takes the raw contents of a Road Trip data file and
// and populates the [Vehicle] object with what it finds inside.
func (v *Vehicle) UnmarshalRoadtrip(data RawFileData) error {
	v.Raw = data

	var err error

	// This seems ripe for future improvement, it should be possible
	// to generate the targets array by reflecting through v and finding
	// the correct pointers to append.
	var targets []any
	targets = append(targets, &v.Vehicles)
	targets = append(targets, &v.FuelRecords)
	targets = append(targets, &v.MaintenanceRecords)
	targets = append(targets, &v.Trips)
	targets = append(targets, &v.Tires)
	targets = append(targets, &v.Valuations)

	for _, target := range targets {
		err = data.UnmarshalRoadtripSection(target)
		if err != nil {
			return fmt.Errorf("unable to parse %s: %w", target, err)
		}
	}

	v.logger.Debug("Loaded Road Trip vehicle data file",
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
