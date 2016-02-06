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
	schemaDistance       = "e"
	schemaGrade          = "g"
	schemaIsPrizeMeet    = "i"
	schemaLongName       = "l"
	schemaName           = "n"
	schemaSortScoreboard = "o"
	schemaPostcode       = "p"
	schemaSort           = "s"
	schemaTime           = "t"
	schemaURL            = "u"
	schemaUnit           = "z"
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
	Latitude  float64 `bson:"x,omitempty"`
	Longitude float64 `bson:"y,omitempty"`
	URL       string  `bson:"u,omitempty"`
	Address   string  `bson:"a,omitempty"`
	Town      string  `bson:"w,omitempty"`
	Postcode  string  `bson:"p,omitempty"`
	AutoInc   AutoInc `bson:"A,omitempty"`
}

// AutoInc is a auto increment counter
type AutoInc struct {
	Mound   uint64 `bson:"M,omitempty"`
	Event   uint64 `bson:"E,omitempty"`
	Club    uint64 `bson:"C,omitempty"`
	Range   uint64 `bson:"R,omitempty"`
	Shooter uint64 `bson:"S,omitempty"`
}

// Mound is exported
type Mound struct {
	//	ID       string `bson:"_id"`
	Distance uint64 `bson:"e,omitempty"`
	Unit     string `bson:"z,omitempty"`
}

// Event is exported
type Event struct {
	ID   string `bson:"_id"`
	Name string `bson:"n"`
	Club string `bson:"C"`
	Date string `bson:"d"`
	Time string `bson:"t"`
}
