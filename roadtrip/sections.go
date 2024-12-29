package roadtrip

import (
	"log/slog"
)

// A FuelRecord contains a single fuel CSV row from the underlying Road Trip
// data file and represents a single vehicle fuel fillup and all of its
// associated attributes.
//
// A file will contain zero or more Fuel records in the FUEL RECORDS section of
// the file.
type FuelRecord struct {
	Odometer     float64 `csv:"Odometer (mi)"`
	TripDistance float64 `csv:"Trip Distance,omitempty"`
	Date         string  `csv:"Date"`
	FillAmount   float64 `csv:"Fill Amount,omitempty"`
	FillUnits    string  `csv:"Fill Units"`
	PricePerUnit float64 `csv:"Price per Unit,omitempty"`
	TotalPrice   float64 `csv:"Total Price,omitempty"`
	PartialFill  string  `csv:"Partial Fill,omitempty"`
	MPG          float64 `csv:"MPG,omitempty"`
	Note         string  `csv:"Note"`
	Octane       string  `csv:"Octane"`
	Location     string  `csv:"Location"`
	Payment      string  `csv:"Payment"`
	Conditions   string  `csv:"Conditions"`
	Reset        string  `csv:"Reset"`
	Categories   string  `csv:"Categories"`
	Flags        string  `csv:"Flags"`
	CurrencyCode int     `csv:"Currency Code,omitempty"`
	CurrencyRate int     `csv:"Currency Rate,omitempty"`
	Latitude     float64 `csv:"Latitude,omitempty"`
	Longitude    float64 `csv:"Longitude,omitempty"`
	ID           int     `csv:"ID,omitempty"`
	FuelEconomy  string  `csv:"Trip Comp Fuel Economy"`
	AvgSpeed     string  `csv:"Trip Comp Avg. Speed"`
	Temperature  float64 `csv:"Trip Comp Temperature,omitempty"`
	DriveTime    string  `csv:"Trip Comp Drive Time"`
	TankNumber   int     `csv:"Tank Number,omitempty"`
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [FuelRecord] object when logging.
func (v FuelRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Float64("odometer", v.Odometer),
		slog.String("date", v.Date),
		slog.String("location", v.Location),
		slog.Float64("totalPrice", v.TotalPrice),
	)
}

// A MaintenceRecord is a single CSV row from the Road Trip data file and
// represents a distinct vehicle maintenance activity with all of its
// associated attributes.
//
// A file will contain zero or more MaintenanceRecord rows in the MAINTENANCE
// RECORDS section of the file.
type MaintenanceRecord struct {
	Description          string  `csv:"Description"`
	Date                 string  `csv:"Date"`
	Odometer             float64 `csv:"Odometer (mi.),omitempty"`
	Cost                 float64 `csv:"Cost,omitempty"`
	Note                 string  `csv:"Note"`
	Location             string  `csv:"Location"`
	Type                 string  `csv:"Type"`
	Subtype              string  `csv:"Subtype"`
	Payment              string  `csv:"Payment"`
	Categories           string  `csv:"Categories"`
	ReminderInterval     string  `csv:"Reminder Interval"`
	ReminderDistance     float64 `csv:"Reminder Distance,omitempty"`
	Flags                string  `csv:"Flags"`
	CurrencyCode         int     `csv:"Currency Code,omitempty"`
	CurrencyRate         int     `csv:"Currency Rate,omitempty"`
	Latitude             float64 `csv:"Latitude,omitempty"`
	Longitude            float64 `csv:"Longitude,omitempty"`
	ID                   int     `csv:"ID,omitempty"`
	NotificationInterval string  `csv:"Notification Interval"`
	NotificationDistance float64 `csv:"Notification Distance,omitempty"`
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [MaintenanceRecord] object when logging.
func (v MaintenanceRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Float64("odometer", v.Odometer),
		slog.String("date", v.Date),
		slog.String("description", v.Description),
		slog.String("location", v.Location),
		slog.Float64("Cost", v.Cost),
	)
}

// A TripRecord  is a single CSV row from the Road Trip data file and
// represents a road trip activity with all of its associated attributes. It is
// date and odometer range bound with a start and end value for each of those
// fields corresponding to the vehicle's service dates and odometer readings.
//
// A file will contain zero or more TripRecord rows in the ROAD TRIPS section
// of the file.
type TripRecord struct {
	Name          string  `csv:"Name"`
	StartDate     string  `csv:"Start Date"`
	StartOdometer float64 `csv:"Start Odometer (mi.),omitempty"`
	EndDate       string  `csv:"End Date"`
	EndOdometer   float64 `csv:"End Odometer,omitempty"`
	Note          string  `csv:"Note"`
	Distance      float64 `csv:"Distance,omitempty"`
	ID            int     `csv:"ID,omitempty"`
	Type          string  `csv:"Type"`
	Categories    string  `csv:"Categories"`
	Flags         string  `csv:"Flags"`
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [TripRecord] object when logging.
func (v TripRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Float64("startOdometer", v.StartOdometer),
		slog.String("startDate", v.StartDate),
		slog.String("note", v.Note),
		slog.Float64("distance", v.Distance),
	)
}

// A VehicleRecord is a single CSV row from the Road Trip data file and
// represents a the vehicle for this file with all of its associated
// attributes.
//
// A file is expected to only contain a single row in the VEHICLE section.
type VehicleRecord struct {
	Name                string  `csv:"Name"`
	Odometer            string  `csv:"Odometer"`
	Units               string  `csv:"Units"`
	Notes               string  `csv:"Notes"`
	TankCapacity        float64 `csv:"Tank Capacity,omitempty"`
	Tank1Units          string  `csv:"Tank Units"`
	HomeCurrency        string  `csv:"Home Currency"`
	Flags               string  `csv:"Flags"`
	IconID              string  `csv:"IconID"`
	FuelUnits           string  `csv:"FuelUnits"`
	TripCompUnits       string  `csv:"TripComp Units"`
	TripCompSpeed       string  `csv:"TripComp Speed"`
	TripCompTemperature string  `csv:"TripComp Temperature"`
	TripCompTimeEnabled string  `csv:"TripComp Time Enabled"`
	OdometerShift       string  `csv:"Odometer Shift"`
	Tank1Type           string  `csv:"Tank 1 Type,optional"`
	Tank2Type           string  `csv:"Tank 2 Type,optional"`
	Tank2Units          string  `csv:"Tank 2 Units,optional"`
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [VehicleRecord] object when logging.
func (v VehicleRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", v.Name),
		slog.String("odometer", v.Odometer),
	)
}

// A VehicleRecord is a single CSV row from the Road Trip data file and

// A TireRecord is a single CSV row from the Road Trip data file and represents
// a set of tires installed on the vehicle with all of its associated
// attributes. It is date and odometer range bound with a start value for each
// of those fields corresponding to the vehicle's service dates and odometer
// readings.
//
// A file will contain zero or more Tire records in the TIRE LOG section of the
// file.
type TireRecord struct {
	Name           string  `csv:"Name"`
	StartDate      string  `csv:"Start Date"`
	StartOdometer  float64 `csv:"Start Odometer (mi.),omitempty"`
	Size           string  `csv:"Size"`
	SizeCorrection string  `csv:"Size Correction"`
	Distance       float64 `csv:"Distance,omitempty"`
	Age            string  `csv:"Age"`
	Note           string  `csv:"Note"`
	Flags          string  `csv:"Flags"`
	ID             int     `csv:"ID,omitempty"`
	ParentID       int     `csv:"ParentID,omitempty"`
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [TireRecord] object when logging.
func (v TireRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", v.Name),
		slog.Float64("startOdometer", v.StartOdometer),
		slog.String("startDate", v.StartDate),
		slog.Float64("distance", v.Distance),
	)
}

// A ValuationRecord is a single CSV row from the Road Trip data file and
// represents the market value of the vehicle at a specified Odometer reading
// and date.
//
// A file will contain zero or more Valuation records in the VALUATIONS section
// of the file.
type ValuationRecord struct {
	Type     string  `csv:"Type"`
	Date     string  `csv:"Date"`
	Odometer float64 `csv:"Odometer,omitempty"`
	Price    string  `csv:"Price"`
	Notes    string  `csv:"Notes"`
	Flags    string  `csv:"Flags"`
}

// LogValue is the handler for [log.slog] to emit structured output for a
// [ValuationRecord] object when logging.
func (v ValuationRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("type", v.Type),
		slog.String("date", v.Date),
		slog.Float64("odometer", v.Odometer),
		slog.String("price", v.Price),
		slog.String("flags", v.Flags),
	)
}
