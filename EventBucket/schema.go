package main

import "regexp"

// Club is exported
type Club struct {
	ID        string  `json:"I"`
	Name      string  `json:"n"`
	IsDefault bool    `json:"b,omitempty"`
	Mounds    []Mound `json:"M,omitempty"`
	Latitude  float32 `json:"x,omitempty"`
	Longitude float32 `json:"y,omitempty"`
	URL       string  `json:"u,omitempty"`
	Address   string  `json:"a,omitempty"`
	Town      string  `json:"w,omitempty"`
	Postcode  string  `json:"p,omitempty"`
	AutoInc   AutoInc `json:"A,omitempty"`
}

// AutoInc is a auto increment counter
type AutoInc struct {
	Mound   uint64 `json:"M,omitempty"`
	Event   uint64 `json:"E,omitempty"`
	Club    uint64 `json:"C,omitempty"`
	Range   uint64 `json:"R,omitempty"`
	Shooter uint64 `json:"S,omitempty"`
}

// Mound is exported
type Mound struct {
	//	ID       string `json:"I"`
	Distance uint64 `json:"e,omitempty"`
	Unit     string `json:"z,omitempty"`
}

// Event is exported
type Event struct {
	ID       string         `json:"I"`
	Name     string         `json:"n"`
	Club     string         `json:"C"`
	Date     string         `json:"d"`
	Time     string         `json:"t"`
	Ranges   []Range        `json:"R,omitempty"`
	AutoInc  AutoInc        `json:"A"`
	Shooters []EventShooter `json:"S,omitempty"`
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
	ID        uint64 `json:"I"`
	Name      string `json:"n"`
	Aggregate string `json:"schemaAggregate,omitempty"`
	Locked    bool   `json:"k,omitempty"`
	IsAgg     bool   `json:"schemaIsAgg,omitempty"` //Prevents aggs switching to normal ranges
	//	ScoreBoard bool                     `json:"s,omitempty"`
	//	Hidden     bool                     `json:"h,omitempty"`
	//	Order      int                      `json:"s,omitempty"`
	//	Status     int                      `json:"t,omitempty"`      //ENUM change to 1 when the first shooter has recorded their first shot change to 2 when the range is finished. http://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-golang
	//	Class      map[string]RangeProperty `json:"omitempty,inline"` //TODO possibly change it to optional grades per range in future
}

// EventShooter is exported
type EventShooter struct {
	ID        uint64 `json:"I"`
	FirstName string `json:"f"` //TODO change these to point to shooters in the other shooter tables
	Surname   string `json:"s"`
	Club      string `json:"C"`
	Grade     uint64 `json:"g"`
	Hidden    bool   `json:"h,omitempty"`
	AgeGroup  uint64 `json:"r,omitempty"`
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
	ID        string `json:"I"`
	SID       int    `json:"schemaSID,omitempty"`
	NraaID    int    `json:"schemaNraaID,omitempty"`
	Surname   string `json:"s,omitempty"`
	FirstName string `json:"f,omitempty"`
	NickName  string `json:"schemaNickName,omitempty"`
	Club      string `json:"C,omitempty"`
	//Skill map[string]Skill	//Grading set by the VRA for each class
	Address string `json:"a,omitempty"`
	Email   string `json:"schemaEmail,omitempty"`
	//Shooter details 0=not modified, 1=updated, 2=merged, 3=deleted
	Status int `json:"schemaStatus,omitempty"`
	//If shooter details are merged with another existing shooter then this is the other NRAA_SID it was merged with
	//When merging set one record to merged, the other to deleted.
	//Both records must set MergedSID to the other corresponding shooter SID
	MergedSID int `json:"m,omitempty"`

	AgeGroup uint64 `json:"r,omitempty"`
	Grade    uint64 `json:"g"`
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
