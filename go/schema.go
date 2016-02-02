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
	Latitude  string  `bson:"schemaLatitude,omitempty"`
	Longitude string  `bson:"schemaLongitude,omitempty"`
	URL       string  `bson:"schemaURL,omitempty"`
	Address   string  `bson:"schemaAddress,omitempty"`
	Town      string  `bson:"schemaTown,omitempty"`
	Postcode  string  `bson:"schemaPostcode,omitempty"`
	AutoInc   AutoInc `bson:"schemaAutoInc,omitempty"`
}

// AutoInc is a auto increment counter
type AutoInc struct {
	Mound   int `bson:"schemaMound,omitempty"`
	Event   int `bson:"schemaEvent,omitempty"`
	Club    int `bson:"schemaClub,omitempty"`
	Range   int `bson:"schemaRange,omitempty"`
	Shooter int `bson:"schemaShooter,omitempty"`
}

// Mound is exported
type Mound struct {
	ID       int    `bson:"schemaID"`
	Distance int    `bson:"schemaDistance,omitempty"`
	Unit     string `bson:"schemaUnit,omitempty"`
	Name     string `bson:"schemaName,omitempty"`
}

// Event is exported
type Event struct {
	ID   string `bson:"schemaID"`
	Name string `bson:"schemaName"`
	Club string `bson:"schemaClub"`
	Date string `bson:"schemaDate"`
	Time string `bson:"schemaTime"`
}
