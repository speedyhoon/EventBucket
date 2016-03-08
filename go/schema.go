package main

import "regexp"

// Club is exported
type Club struct {
	ID        string  `json:"schemaID"`
	Name      string  `json:"schemaName"`
	IsDefault bool    `json:"schemaIsDefault,omitempty"`
	Mounds    []Mound `json:"schemaMound,omitempty"`
	Latitude  float32 `json:"schemaLatitude,omitempty"`
	Longitude float32 `json:"schemaLongitude,omitempty"`
	URL       string  `json:"schemaURL,omitempty"`
	Address   string  `json:"schemaAddress,omitempty"`
	Town      string  `json:"schemaTown,omitempty"`
	Postcode  string  `json:"schemaPostcode,omitempty"`
	AutoInc   AutoInc `json:"schemaAutoInc,omitempty"`
}

// AutoInc is a auto increment counter
type AutoInc struct {
	Mound   uint64 `json:"schemaMound,omitempty"`
	Event   uint64 `json:"schemaEvent,omitempty"`
	Club    uint64 `json:"schemaClub,omitempty"`
	Range   uint64 `json:"schemaRange,omitempty"`
	Shooter uint64 `json:"schemaShooter,omitempty"`
}

// Mound is exported
type Mound struct {
	//	ID       string `json:"schemaID"`
	Distance uint64 `json:"schemaDistance,omitempty"`
	Unit     string `json:"schemaUnit,omitempty"`
}

// Event is exported
type Event struct {
	ID       string         `json:"schemaID"`
	Name     string         `json:"schemaName"`
	Club     string         `json:"schemaClub"`
	Date     string         `json:"schemaDate"`
	Time     string         `json:"schemaTime"`
	Ranges   []Range        `json:"schemaRange,omitempty"`
	AutoInc  AutoInc        `json:"schemaAutoInc"`
	Shooters []EventShooter `json:"schemaShooter,omitempty"`
	/*Grades         []uint          `json:"schemaGrades,omitempty"`
	SortScoreboard string         `json:"o,omitempty"`
	IsPrizeMeet    bool           `json:"p,omitempty"`
	Closed         bool           `json:"l,omitempty"`
	TeamCat        map[string]TeamCat      `json:"A,omitempty"`
	Teams          map[string]Team         `json:"T,omitempty"`
	*/
}

// Range is exported
type Range struct {
	Name      string `json:"schemaName"`
	Aggregate string `json:"schemaAggregate,omitempty"` //TODO Maybe this one could be a pointer to prevent it from being removed?
	//	ScoreBoard bool                     `json:"s,omitempty"`
	//	Locked     bool                     `json:"l,omitempty"`
	//	Hidden     bool                     `json:"h,omitempty"`
	//	Order      int                      `json:"schemaSort,omitempty"`
	//	Status     int                      `json:"t,omitempty"`      //ENUM change to 1 when the first shooter has recorded their first shot change to 2 when the range is finished. http://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-golang
	//	Class      map[string]RangeProperty `json:"omitempty,inline"` //TODO possibly change it to optional grades per range in future
	//	ID         *int                     `json:"i,omitempty"`
	IsAgg bool `json:"schemaIsAgg,omitempty"` //Prevents aggs switching to normal ranges //TODO is there a better way to determine an empty agg rather than having this separate column?
}

// EventShooter is exported
type EventShooter struct {
	ID        uint64 `json:"schemaID"`
	FirstName string `json:"schemaFirstName"` //TODO change these to point to shooters in the other shooter tables
	Surname   string `json:"schemaSurname"`
	Club      string `json:"schemaClub"`
	Grade     uint64 `json:"schemaGrade"`
	Hidden    bool   `json:"schemaHidden,omitempty"`
	AgeGroup  uint64 `json:"schemaAgeGroup,omitempty"`
	//	Scores    map[string]Score `json:"omitempty,inline"` //S is not used!
	LinkedID *int `json:"schemaLinkedID,omitempty"` //For duplicating shooters that are in different classes with the same score
	SID      int  `json:"schemaSID,omitempty"`
	Disabled bool `json:"schemaDisabled,omitempty"`
	//SCOREBOARD
	position string //DON'T SAVE THIS TO DB! used for scoreboard only.
	warning  uint8  //DON'T SAVE THIS TO DB! used for scoreboard only.
	//		0 = nil
	//		1 = shoot off
	//		2 = warning, no score
	//		3 = incomplete
	//		4 = highest posible score

	//START-SHOOTING & TOTAL-SCORES
	gradeSeparator bool //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
}

// Shooter is exported
type Shooter struct {
	ID        string `json:"schemaID"`
	SID       int    `json:"schemaSID,omitempty"`
	NraaID    int    `json:"schemaNraaID,omitempty"`
	Surname   string `json:"schemaSurname,omitempty"`
	FirstName string `json:"schemaFirstName,omitempty"`
	NickName  string `json:"schemaNickName,omitempty"`
	Club      string `json:"schemaClub,omitempty"`
	//Skill map[string]Skill	//Grading set by the VRA for each class
	Address string `json:"schemaAddress,omitempty"`
	Email   string `json:"schemaEmail,omitempty"`
	//Shooter details 0=not modified, 1=updated, 2=merged, 3=deleted
	Status int `json:"schemaStatus,omitempty"`
	//If shooter details are merged with another existing shooter then this is the other NRAA_SID it was merged with
	//When merging set one record to merged, the other to deleted.
	//Both records must set MergedSID to the other corresponding shooter SID
	MergedSID int `json:"schemaMergedSID,omitempty"`

	AgeGroup uint64 `json:"schemaAgeGroup,omitempty"`
	Grade    uint64 `json:"schemaGrade"`
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
