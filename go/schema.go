package main

const (
	//Databse collection names
	tblClub    = "C"
	tblEvent   = "E"
	tblAutoInc = "U"

	//Collection property names
	schemaClub           = "schemaClub"
	schemaEvent          = "schemaEvent"
	schemaMound          = "schemaMound"
	schemaRange          = "schemaRange"
	schemaShooter        = "schemaShooter"
	schemaAutoInc        = "schemaAutoInc"
	schemaID             = "schemaID"
	schemaAddress        = "schemaAddress"
	schemaIsDefault      = "schemaIsDefault"
	schemaClose          = "schemaClose"
	schemaDate           = "schemaDate"
	schemaDistance       = "schemaDistance"
	schemaGrade          = "schemaGrade"
	schemaIsPrizeMeet    = "schemaIsPrizeMeet"
	schemaLongName       = "schemaLongName"
	schemaName           = "schemaName"
	schemaSortScoreboard = "schemaSortScoreboard"
	schemaPostcode       = "schemaPostcode"
	schemaSort           = "schemaSort"
	schemaTime           = "schemaTime"
	schemaURL            = "schemaURL"
	schemaUnit           = "schemaUnit"
	schemaTown           = "schemaTown"
	schemaLatitude       = "schemaLatitude"
	schemaLongitude      = "schemaLongitude"
)

// Club is exported
type Club struct {
	ID        string  `bson:"schemaID"`
	Name      string  `bson:"schemaName"`
	IsDefault bool    `bson:"schemaIsDefault,omitempty"`
	Mounds    []Mound `bson:"schemaMound,omitempty"`
	Latitude  float64 `bson:"schemaLatitude,omitempty"`
	Longitude float64 `bson:"schemaLongitude,omitempty"`
	URL       string  `bson:"schemaURL,omitempty"`
	Address   string  `bson:"schemaAddress,omitempty"`
	Town      string  `bson:"schemaTown,omitempty"`
	Postcode  string  `bson:"schemaPostcode,omitempty"`
	AutoInc   AutoInc `bson:"schemaAutoInc,omitempty"`
}

// AutoInc is a auto increment counter
type AutoInc struct {
	Mound   uint64 `bson:"schemaMound,omitempty"`
	Event   uint64 `bson:"schemaEvent,omitempty"`
	Club    uint64 `bson:"schemaClub,omitempty"`
	Range   uint64 `bson:"schemaRange,omitempty"`
	Shooter uint64 `bson:"schemaShooter,omitempty"`
}

// Mound is exported
type Mound struct {
	//	ID       string `bson:"schemaID"`
	Distance uint64 `bson:"schemaDistance,omitempty"`
	Unit     string `bson:"schemaUnit,omitempty"`
}

// Event is exported
type Event struct {
	ID   string `bson:"schemaID"`
	Name string `bson:"schemaName"`
	Club string `bson:"schemaClub"`
	Date string `bson:"schemaDate"`
	Time string `bson:"schemaTime"`
}
