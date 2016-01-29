package main

const (
	//Databse collection names
	tblClub    = "C"
	tblEvent   = "E"
	tblAutoInc = "U"

	//Collection property names
	schemaClub      = "schemaClub"
	schemaEvent     = "schemaEvent"
	schemaMound     = "schemaMound"
	schemaRange     = "schemaRange"
	schemaShooter   = "schemaShooter"
	schemaAutoInc   = "schemaAutoInc"
	schemaID        = "schemaID"
	schemaAddress   = "schemaAddress"
	schemaIsDefault = "schemaIsDefault"
	schemaClose     = "schemaClose"
	schemaDate      = "schemaDate"

	schemaGrade          = "schemaGrade"
	schemaIsPrizeMeet    = "schemaIsPrizeMeet"
	schemaLongName       = "schemaLongName"
	schemaName           = "schemaName"
	schemaSortScoreboard = "schemaSortScoreboard"
	schemaPostcode       = "schemaPostcode"
	schemaSort           = "schemaSort"
	schemaTime           = "schemaTime"
	schemaURL            = "schemaURL"

	schemaTown      = "schemaTown"
	schemaLatitude  = "schemaLatitude"
	schemaLongitude = "schemaLongitude"
)

// Club is exported
type Club struct {
	ID        string `bson:"schemaID"`
	Name      string `bson:"schemaName"`
	IsDefault bool   `bson:"schemaIsDefault,omitempty"`
	//	LongName  string  `bson:"l,omitempty"`
	//	Mounds    []Mound `bson:"M,omitempty"`
	//	Latitude  string  `bson:"t,omitempty"`
	//	Longitude string  `bson:"g,omitempty"`
	//	URL       string  `bson:"u,omitempty"`
	//	Address   string  `bson:"a,omitempty"`
	//	Town      string  `bson:"w,omitempty"`
	//	PostCode  string  `bson:"p,omitempty"`
	//	AutoInc   AutoInc `bson:"^^schemaAutoInc^^,omitempty"`
}

// Event is exported
type Event struct {
	ID   string `bson:"schemaID"`
	Name string `bson:"schemaName"`
	Club string `bson:"schemaClub"`
}
