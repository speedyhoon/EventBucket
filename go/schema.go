package main

import (
	"fmt"
	"regexp"
	"time"
)

//Club is exported
type Club struct {
	ID        string  `json:"I"`
	Name      string  `json:"n"`
	IsDefault bool    `json:"d,omitempty"`
	Mounds    []Mound `json:"M,omitempty"`
	Latitude  float32 `json:"x,omitempty"`
	Longitude float32 `json:"y,omitempty"`
	URL       string  `json:"u,omitempty"`
	Address   string  `json:"a,omitempty"`
	Town      string  `json:"t,omitempty"`
	Postcode  string  `json:"p,omitempty"`
	AutoInc   AutoInc `json:"A,omitempty"`
}

//AutoInc is a auto increment counter
type AutoInc struct {
	Mound   uint `json:"M,omitempty"`
	Event   uint `json:"E,omitempty"`
	Club    uint `json:"C,omitempty"`
	Range   uint `json:"R,omitempty"`
	Shooter uint `json:"S,omitempty"`
}

//Mound could in future contain additional details like distance
type Mound struct {
	Name string `json:"n,omitempty"`
	ID   uint   `json:"-"`
}

//Event is exported
type Event struct {
	ID       string         `json:"I"`
	Name     string         `json:"n"`
	ClubID   string         `json:"C,omitempty"`
	Club     string         `json:"c,omitempty"`
	Date     string         `json:"d,omitempty"`
	Time     string         `json:"t,omitempty"`
	Ranges   []Range        `json:"R,omitempty"`
	AutoInc  AutoInc        `json:"A,omitempty"`
	Shooters []EventShooter `json:"S,omitempty"`
	Closed   bool           `json:"z,omitempty"`
	Grades   []uint         `json:"g,omitempty"`
	/*SortScoreboard string         `json:"o,omitempty"`
	IsPrizeMeet    bool           `json:"p,omitempty"`
	TeamCats       map[string]TeamCats     `json:"A,omitempty"`
	Teams          map[string]Team         `json:"T,omitempty"`
	AverTwin bool           `json:"a,omitempty"` //TODO remove support for allowing shooter with the same details to enter an event twice or more.*/
}

//Range is exported
type Range struct {
	ID     uint   `json:"I"`
	Name   string `json:"n"`
	Aggs   []uint `json:"a,omitempty"`
	Locked bool   `json:"l,omitempty"`
	IsAgg  bool   `json:"i,omitempty"` //Prevents aggs switching to normal ranges
	Order  int    `json:"o,omitempty"`
	Status uint8  `json:"u,omitempty"` //ENUM change to 1 when the first shooter has recorded their first shot change to 2 when the range is finished. http://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-golang
	//Class      map[string]RangeProperty `json:"omitempty,inline"` //TODO possibly change it to optional grades per range in future
}

//StrID returns Range.ID as a string instead of an unsigned integer
func (r Range) StrID() string {
	return fmt.Sprintf("%v", r.ID)
}

//Score is exported
type Score struct {
	Total      uint   `json:"t,omitempty"`
	Centers    uint   `json:"c,omitempty"`
	Centers2   uint   `json:"n,omitempty"`
	Shots      string `json:"s,omitempty"` //Don't include this in the scoreboard struct when using a different []EventShooter
	CountBack  string `json:"v,omitempty"`
	CountBack2 string `json:"x,omitempty"`
	ShootOff   uint   `json:"h,omitempty"`
	//position  int    `json:"p,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//warning   uint8    `json:"w,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//Ordinal   string `json:"o,omitempty"`
	//Info      string `json:"i,omitempty"`
}

//EventShooter is exported
type EventShooter struct {
	ID        uint             `json:"I"`
	FirstName string           `json:"f"` //TODO change these to point to shooters in the other shooter tables? Would require extra look ups though :(
	Surname   string           `json:"s"`
	Club      string           `json:"C"` //TODO change Club to struct Club{ID uint, Name string} ??
	ClubID    string           `json:"c"`
	Grade     uint             `json:"g"`
	Hidden    bool             `json:"h,omitempty"`
	AgeGroup  uint             `json:"r,omitempty"`
	Scores    map[string]Score `json:"S,omitempty"` //TODO look into using uint as index instead of string. //TODO look into inline the json field and if it would work with uint indexes //Scores   []Score `json:"schemaScores,omitempty"`   //S is not used!
	LinkedID  uint             `json:"l,omitempty"` //For duplicating shooters that are in different classes with the same score
	EID       string           `json:"M,omitempty"` //Points to EventBucket Shooter ID
	Disabled  bool             `json:"d,omitempty"`
	Ladies    bool             `json:"x,omitempty"`
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

//Shooter is exported
type Shooter struct {
	ID        string         `json:"I"`
	SID       uint           `json:"M,omitempty"`
	NID       uint           `json:"N,omitempty"` //NRAA sequential integer id.
	FirstName string         `json:"f,omitempty"`
	Surname   string         `json:"s,omitempty"`
	NickName  string         `json:"n,omitempty"`
	Club      string         `json:"C,omitempty"`
	ClubID    string         `json:"c,omitempty"`
	Skill     map[uint]Skill `json:"K,omitempty"` //Grading set by the VRA for each class
	Address   string         `json:"a,omitempty"`
	Email     string         `json:"e,omitempty"`
	Status    int            `json:"v,omitempty"` //Shooter details 0=not modified, 1=updated, 2=merged, 3=deleted
	//If shooter details are merged with another existing shooter then this is the other NRAA_SID it was merged with
	//When merging set one record to merged, the other to deleted.
	//Both records must set MergedSID to the other corresponding shooter SID
	MergedSID int       `json:"m,omitempty"`
	Modified  time.Time `json:"o,omitempty"`
	AgeGroup  uint      `json:"r,omitempty"`
	Ladies    bool      `json:"l,omitempty"`
	Grade     []uint    `json:"g,omitempty"` //TODO change to a slice in future
}

//Skill is exported
type Skill struct {
	Threshold string  `json:"t,omitempty"`
	AvgScore  float64 `json:"a,omitempty"`
	ShootQty  int     `json:"s,omitempty"`
}

type shooterScore struct {
	rangeID string
	id      uint
	score   Score
}

type field struct {
	name, Error, Value string
	Required, Disable  bool
	Options            []option
	maxLen, minLen     int
	min, max, step     float32
	AutoFocus          bool
	size               uint8
	Checked            bool //only used by checkboxes
	regex              *regexp.Regexp
	v8                 func(*field, ...string)
	defValue           func() []string
	valueUint          uint
	valueUintSlice     []uint
	valueFloat32       float32
	manyRequired       []int
	manyRequiredQty    uint
	Placeholder        string
}
