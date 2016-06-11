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
	/*TeamCats       map[string]TeamCats     `json:"A,omitempty"`
	Teams          map[string]Team         `json:"T,omitempty"`*/
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
	//position   int    `json:"-"` //Used for scoreboard only.
	//warning    uint8  `json:"-"` //Used for scoreboard only.
	//Ordinal    string `json:"o,omitempty"`
	//Info       string `json:"i,omitempty"`
}

//EventShooter is exported
type EventShooter struct {
	ID        uint             `json:"I"`
	FirstName string           `json:"f"`
	Surname   string           `json:"s"`
	Club      string           `json:"C,omitempty"`
	ClubID    string           `json:"c,omitempty"`
	Grade     uint             `json:"g,omitempty"`
	Hidden    bool             `json:"h,omitempty"`
	AgeGroup  uint             `json:"r,omitempty"`
	Scores    map[string]Score `json:"S,omitempty"` //Using string instead of uint as an index because JSON doesn't support map[uint]Score
	LinkedID  uint             `json:"l,omitempty"` //For duplicating shooters that are in different classes with the same score
	EID       string           `json:"M,omitempty"` //Points to EventBucket Shooter ID
	Disabled  bool             `json:"d,omitempty"`
	Ladies    bool             `json:"x,omitempty"`
	position  string           `json:"-"` //Used for scoreboard only.
	warning   uint8            `json:"-"` //Used for scoreboard only.
	//0 = nil
	//1 = shoot off
	//2 = no score
	//3 = incomplete
	//4 = highest possible score
	GradeSeparator bool `json:"-"` //Used for enterShots and enterRangeTotals only.
	ClassSeparator bool `json:"-"` //Used for enterShots and enterRangeTotals only.
}

//Shooter is exported
type Shooter struct {
	ID        string           `json:"I"`
	SID       uint             `json:"M,omitempty"`
	NID       uint             `json:"N,omitempty"` //NRAA sequential integer id.
	FirstName string           `json:"f,omitempty"`
	Surname   string           `json:"s,omitempty"`
	NickName  string           `json:"n,omitempty"`
	Club      string           `json:"C,omitempty"`
	ClubID    string           `json:"c,omitempty"`
	Skill     map[string]Skill `json:"K,omitempty"` //Grading set by the NRAA for each class
	Address   string           `json:"a,omitempty"`
	Email     string           `json:"e,omitempty"`
	Status    int              `json:"v,omitempty"` //Shooter details 0=not modified, 1=updated, 2=merged, 3=deleted
	//If shooter details are merged with another existing shooter then this is the other NRAA_SID it was merged with
	//When merging set one record to merged, the other to deleted.
	//Both records must set MergedSID to the other corresponding shooter SID
	MergedSID int       `json:"m,omitempty"`
	Modified  time.Time `json:"o,omitempty"`
	AgeGroup  uint      `json:"r,omitempty"`
	Ladies    bool      `json:"l,omitempty"`
	Grade     []uint    `json:"g,omitempty"`
}

//Skill is used for importing shooters from JSON files
type Skill struct {
	AvgScore  float64 `json:"a,omitempty"`
	ShootQty  int     `json:"q,omitempty"`
	Threshold float32 `json:"t,omitempty"`
}

type shooterScore struct {
	rangeID string
	id      uint
	score   Score
}

type field struct {
	name, Error, Value, Placeholder       string
	Required, Disable, AutoFocus, Checked bool
	Options                               []option
	maxLen, minLen                        int
	min, max, step, valueFloat32          float32
	size                                  uint8
	regex                                 *regexp.Regexp
	v8                                    func(*field, ...string)
	valueUint                             uint
	valueUintSlice                        []uint
}
