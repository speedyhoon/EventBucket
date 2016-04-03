package main

import (
	"fmt"
	"regexp"
)

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
	Mound   uint `json:"schemaMound,omitempty"`
	Event   uint `json:"schemaEvent,omitempty"`
	Club    uint `json:"schemaClub,omitempty"`
	Range   uint `json:"schemaRange,omitempty"`
	Shooter uint `json:"schemaShooter,omitempty"`
}

// Mound is exported
type Mound struct {
	//	ID       string `json:"schemaID"`
	Distance uint   `json:"schemaDistance,omitempty"`
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
	Closed   bool           `json:"schemaClosed,omitempty"`
	/*Grades       []uint         `json:"schemaGrades,omitempty"`
	SortScoreboard string         `json:"o,omitempty"`
	IsPrizeMeet    bool           `json:"p,omitempty"`
	TeamCats       map[string]TeamCats     `json:"A,omitempty"`
	Teams          map[string]Team         `json:"T,omitempty"`
	*/
}

// Range is exported
type Range struct {
	ID     uint   `json:"schemaID"`
	Name   string `json:"schemaName"`
	Aggs   []uint `json:"schemaAggregate,omitempty"`
	Locked bool   `json:"schemaLocked,omitempty"`
	IsAgg  bool   `json:"schemaIsAgg,omitempty"` //Prevents aggs switching to normal ranges
	Order  int    `json:"schemaSort,omitempty"`
	Status uint8  `json:"schemaStatus,omitempty"` //ENUM change to 1 when the first shooter has recorded their first shot change to 2 when the range is finished. http://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-golang
	//Class      map[string]RangeProperty `json:"omitempty,inline"` //TODO possibly change it to optional grades per range in future
	//ScoreBoard bool                     `json:"s,omitempty"`
	//Hidden     bool                     `json:"h,omitempty"`
}

//Id method returns Range.ID as a string instead of an unsigned integer
func (r Range) strID() string {
	return fmt.Sprintf("%v", r.ID)
}

// Score is exported
type Score struct {
	//TODO the schema should change so that it can use unsigned  bit numbers instead
	Total      uint   `json:"schemaTotal"`
	Centers    uint   `json:"schemaCenters,omitempty"`
	Shots      string `json:"s,omitempty"` //Don't include this in the scoreboard struct when using a different []EventShooter
	CountBack  string `json:"v,omitempty"`
	CountBack2 string `json:"x,omitempty"`
	ShootOff   uint   `json:"f,omitempty"`
	//position  int    `json:"p,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//warning   uint8    `json:"w,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//Ordinal   string `json:"o,omitempty"`
	//Info      string `json:"i,omitempty"`
}

// EventShooter is exported
type EventShooter struct {
	ID        uint             `json:"schemaID"`
	FirstName string           `json:"schemaFirstName"` //TODO change these to point to shooters in the other shooter tables
	Surname   string           `json:"schemaSurname"`
	Club      string           `json:"schemaClub"`
	Grade     uint             `json:"schemaGrade"`
	Hidden    bool             `json:"schemaHidden,omitempty"`
	AgeGroup  uint             `json:"schemaAgeGroup,omitempty"`
	Scores    map[string]Score `json:"schemaScores,omitempty"`   //Scores   []Score `json:"schemaScores,omitempty"`   //S is not used!
	LinkedID  uint             `json:"schemaLinkedID,omitempty"` //For duplicating shooters that are in different classes with the same score
	SID       int              `json:"schemaSID,omitempty"`
	Disabled  bool             `json:"schemaDisabled,omitempty"`
	//SCOREBOARD
	position string //DON'T SAVE THIS TO DB! used for scoreboard only.
	warning  uint8  //DON'T SAVE THIS TO DB! used for scoreboard only.
	//		0 = nil
	//		1 = shoot off
	//		2 = warning, no score
	//		3 = incomplete
	//		4 = highest posible score

	//START-SHOOTING & TOTAL-SCORES
	GradeSeparator bool //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
	ClassSeparator bool //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
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

	AgeGroup uint `json:"schemaAgeGroup,omitempty"`
	Grade    uint `json:"schemaGrade"`
}

type field struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     int
	min, max, step     float32
	AutoFocus          bool
	size               uint8
	Checked            bool //only used by checkboxes
	regex              *regexp.Regexp
	internalValue      interface{}
	v8                 func(*field, ...string)
	defValue           func() []string
	valueUint          uint
	valueUintSlice     []uint
	valueFloat32       float32
}
