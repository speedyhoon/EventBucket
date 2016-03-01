package main

import "regexp"

const (
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
	schemaFirstName      = "f"
	schemaSurname        = "s"
	schemaAgeGroup       = "r"
)

// Club is exported
type Club struct {
	ID        string  `bson:"_id"`
	Name      string  `bson:"n"`
	IsDefault bool    `bson:"b,omitempty"`
	Mounds    []Mound `bson:"M,omitempty"`
	Latitude  float32 `bson:"x,omitempty"`
	Longitude float32 `bson:"y,omitempty"`
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
	ID       string         `bson:"_id"`
	Name     string         `bson:"n"`
	Club     string         `bson:"C"`
	DateTime string         `bson:"d"`
	Ranges   []Range        `bson:"R,omitempty"`
	AutoInc  AutoInc        `bson:"A"`
	Shooters []EventShooter `bson:"schemaSHOOTER,omitempty"`
}

/*
// Event is exported
type Event struct {
	ID             string         `bson:"_id"`
	Club           string         `bson:"c"`
	Name           string         `bson:"n"`
	Date           string         `bson:"d,omitempty"`
	Time           string         `bson:"t,omitempty"`
	Grades         []int          `bson:"^^schemaGRADES^^,omitempty"`
	SortScoreboard string         `bson:"o,omitempty"`
	IsPrizeMeet    bool           `bson:"p,omitempty"`
	Closed         bool           `bson:"l,omitempty"`
	Ranges         []Range        `bson:"^^schemaRANGE^^,omitempty"`
	Shooters       []EventShooter `bson:"^^schemaSHOOTER^^,omitempty"`
	//TeamCat        map[string]TeamCat      `bson:"A,omitempty"`
	//Teams          map[string]Team         `bson:"T,omitempty"`
	//Datetime string				`bson:"d,omitempty"`		No browser currently supports date time, so settling for separate fields that google chrome allows
}*/

// Range is exported
type Range struct {
	Name      string `bson:"n"`
	Aggregate string `bson:"a,omitempty"` //TODO Maybe this one could be a pointer to prevent it from being removed?
	//	ScoreBoard bool                     `bson:"s,omitempty"`
	//	Locked     bool                     `bson:"l,omitempty"`
	//	Hidden     bool                     `bson:"h,omitempty"`
	//	Order      int                      `bson:"^^schemaSORT^^,omitempty"`
	//	Status     int                      `bson:"t,omitempty"`      //ENUM change to 1 when the first shooter has recorded their first shot change to 2 when the range is finished. http://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-golang
	//	Class      map[string]RangeProperty `bson:"omitempty,inline"` //TODO possibly change it to optional grades per range in future
	//	ID         *int                     `bson:"i,omitempty"`
	IsAgg bool `bson:"g,omitempty"` //Prevents aggs switching to normal ranges //TODO is there a better way to determine an empty agg rather than having this separate column?
}

// EventShooter is exported
type EventShooter struct {
	FirstName string `bson:"f"` //TODO change these to point to shooters in the other shooter tables
	Surname   string `bson:"s"`
	Club      string `bson:"c"`
	Grade     uint64 `bson:"g"`
	Hidden    bool   `bson:"h,omitempty"`
	AgeGroup  uint64 `bson:"a,omitempty"`
	//	Scores    map[string]Score `bson:"omitempty,inline"` //S is not used!
	LinkedID *int `bson:"l,omitempty"` //For duplicating shooters that are in different classes with the same score
	SID      int  `bson:"d,omitempty"`
	Disabled bool `bson:"b,omitempty"`
	//SCOREBOARD
	ID       uint64 `bson:"i,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	Position string `bson:"x,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	Warning  uint8  `bson:"y,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//		0 = nil
	//		1 = shoot off
	//		2 = warning, no score
	//		3 = incomplete
	//		4 = highest posible score

	//START-SHOOTING & TOTAL-SCORES
	GradeSeparator bool `bson:"z,omitempty"` //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
	//	ID string									`bson:"w,omitempty"`//DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
}

// Shooter is exported
type Shooter struct {
	SID       int    `bson:"_id"`
	NraaID    int    `bson:"i,omitempty"`
	Surname   string `bson:"s,omitempty"`
	FirstName string `bson:"f,omitempty"`
	NickName  string `bson:"n,omitempty"`
	Club      string `bson:"c,omitempty"`
	//Skill map[string]Skill	//Grading set by the VRA for each class
	Address string `bson:"a,omitempty"`
	Email   string `bson:"e,omitempty"`
	//Shooter details 0=not modified, 1=updated, 2=merged, 3=deleted
	Status int `bson:"t,omitempty"`
	//If shooter details are merged with another existing shooter then this is the other NRAA_SID it was merged with
	//When merging set one record to merged, the other to deleted.
	//Both records must set MergedSID to the other corresponding shooter SID
	MergedSID int `bson:"m,omitempty"`

	AgeGroup uint64 `bson:"a,omitempty"`
	ID       string `bson:"i,omitempty"`
	Grade    uint64 `bson:"g"`
}

type field struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     uint
	min, max, step     float32
	AutoFocus          bool
	size               uint8
	Checked            bool //only used by checkboxes
	regex              *regexp.Regexp
	internalValue      interface{}
	v8                 func([]string, field) (interface{}, string)
	defValue           func() []string
}
