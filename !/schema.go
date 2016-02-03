package main

const (
	//Databse collection names
	tblClub    = "C"
	tblEvent   = "E"
	tblAutoInc = "U"

	//Collection property names
	schemaClub           = "C"
	schemaEvent          = "E"
	schemaMound          = "M"
	schemaRange          = "R"
	schemaShooter        = "S"
	schemaAutoInc        = "A"
	schemaID             = "_id"
	schemaAddress        = "a"
	schemaIsDefault      = "b"
	schemaClose          = "c"
	schemaDate           = "d"
	schemaDistance       = "schemaDistance"
	schemaGrade          = "g"
	schemaIsPrizeMeet    = "i"
	schemaLongName       = "l"
	schemaName           = "n"
	schemaSortScoreboard = "o"
	schemaPostcode       = "schemaPostcode"
	schemaSort           = "s"
	schemaTime           = "t"
	schemaURL            = "u"
	schemaUnit           = "schemaUnit"
	schemaTown           = "w"
	schemaLatitude       = "x"
	schemaLongitude      = "y"
)

// Club is exported
type Club struct {
	ID        string  `bson:"_id"`
	Name      string  `bson:"n"`
	IsDefault bool    `bson:"b,omitempty"`
	Mounds    []Mound `bson:"M,omitempty"`
	Latitude  string  `bson:"x,omitempty"`
	Longitude string  `bson:"y,omitempty"`
	URL       string  `bson:"u,omitempty"`
	Address   string  `bson:"a,omitempty"`
	Town      string  `bson:"w,omitempty"`
	Postcode  string  `bson:"schemaPostcode,omitempty"`
	AutoInc   AutoInc `bson:"A,omitempty"`
}

// AutoInc is a auto increment counter
type AutoInc struct {
	Mound   int `bson:"M,omitempty"`
	Event   int `bson:"E,omitempty"`
	Club    int `bson:"C,omitempty"`
	Range   int `bson:"R,omitempty"`
	Shooter int `bson:"S,omitempty"`
}

// Mound is exported
type Mound struct {
	ID       int    `bson:"_id"`
	Distance int    `bson:"schemaDistance,omitempty"`
	Unit     string `bson:"schemaUnit,omitempty"`
	Name     string `bson:"n,omitempty"`
}

// Event is exported
type Event struct {
	ID   string `bson:"_id"`
	Name string `bson:"n"`
	Club string `bson:"C"`
	Date string `bson:"d"`
	Time string `bson:"t"`
}
