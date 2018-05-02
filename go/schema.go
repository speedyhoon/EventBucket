package main

import (
	"fmt"
	"strings"
	"time"
)

//AutoInc is an auto increment counter
type AutoInc struct {
	Event   uint `json:"E,omitempty"`
	Club    uint `json:"C,omitempty"`
	Range   uint `json:"R,omitempty"`
	Shooter uint `json:"S,omitempty"`
}

//Club is exported
type Club struct {
	ID        string   `json:"I"`
	Name      string   `json:"n"`
	IsDefault bool     `json:"d,omitempty"`
	Mounds    []string `json:"M,omitempty"`
	Latitude  float32  `json:"x,omitempty"`
	Longitude float32  `json:"y,omitempty"`
	URL       string   `json:"u,omitempty"`
	Address   string   `json:"a,omitempty"`
	Town      string   `json:"w,omitempty"`
	Postcode  string   `json:"p,omitempty"`
	AutoInc   AutoInc  `json:"U,omitempty"`
}

//Mound could in future contain additional details like distance
type Mound struct {
	Name string
	ID   uint
}

//Event is exported
type Event struct {
	ID       string         `json:"I"`
	Name     string         `json:"n"`
	ClubID   string         `json:"C,omitempty"`
	Club     Club           `json:"-"`
	Date     string         `json:"d,omitempty"`
	Time     string         `json:"t,omitempty"`
	ISO      time.Time      `json:"-"` //TODO Change Date & Time to time.Time?
	Ranges   []Range        `json:"R,omitempty"`
	AutoInc  AutoInc        `json:"U,omitempty"`
	Shooters []EventShooter `json:"S,omitempty"`
	Grades   []uint         `json:"G,omitempty"`
	Closed   bool           `json:"z,omitempty"`
	/*TeamCats map[string]TeamCats `json:"M,omitempty"` //TeamCategories
	Teams    map[string]Team     `json:"T,omitempty"`*/
}

//RangeNames returns a string of comma separated event range names excluding aggregate ranges
func (e Event) RangeNames() string {
	var names []string
	for _, r := range e.Ranges {
		if !r.IsAgg {
			names = append(names, r.Name)
		}
	}
	return strings.Join(names, ", ")
}

/*Championship is exported
type Championship struct {
	ID       string         `json:"I"`
	Name     string         `json:"n"`
	Events   []uint         `json:"E,omitempty"`
	Club     string         `json:"C,omitempty"`
	AutoInc  AutoInc        `json:"U,omitempty"`
	Grades   []uint         `json:"G,omitempty"`
	Closed   bool           `json:"z,omitempty"`
	TeamCats map[string]TeamCats `json:"M,omitempty"` //TeamCategories
	Teams    map[string]Team     `json:"T,omitempty"`
}*/

//Range is exported
type Range struct {
	ID     uint   `json:"I"`
	Name   string `json:"n"`
	Aggs   []uint `json:"A,omitempty"`
	Order  uint   `json:"-"`
	Locked bool   `json:"k,omitempty"`
	IsAgg  bool   `json:"g,omitempty"` //Prevents aggs switching to normal ranges
	//Status uint8  `json:"u,omitempty"` //ENUM change to 1 when the first shooter has recorded their first shot. Change to 2 when the range is finished.
}

//StrID returns Range.ID as a string instead of an unsigned integer
func (r Range) StrID() string {
	return fmt.Sprintf("%v", r.ID)
}

type rID uint //rangeID

//StrID returns Range.ID as a string instead of an unsigned integer
func (r rID) StrID() string {
	return fmt.Sprintf("%d", r)
}

//Score is exported
type Score struct {
	Total      uint   `json:"t,omitempty"`
	Centers    uint   `json:"c,omitempty"`
	Centers2   uint   `json:"2,omitempty"`
	Shots      string `json:"s,omitempty"` //Don't include this in the scoreboard struct when using a different []EventShooter
	Sighters   string `json:"i,omitempty"`
	CountBack  string `json:"v,omitempty"`
	CountBack2 string `json:"x,omitempty"`
	ShootOff   uint   `json:"h,omitempty"`
	Ordinal    string `json:"o,omitempty"` //Used for scoreboard only.
	Position   uint   `json:"p,omitempty"`
	//Warning    uint8  `json:"w,omitempty"`
}

//ScoreMap represents a list of Scores indexed by rangeIDs
type ScoreMap map[string]Score

//get returns a Score given a rangeID index
func (s ScoreMap) get(id uint) (sc Score, ok bool) {
	sc, ok = s[fmt.Sprintf("%d", id)]
	return
}

type sID uint //shooterID

//EventShooter is exported
type EventShooter struct {
	ID             uint     `json:"I"`
	FirstName      string   `json:"f"`
	Surname        string   `json:"s"`
	Club           string   `json:"C,omitempty"`
	Grade          uint     `json:"G,omitempty"`
	AgeGroup       uint     `json:"r,omitempty"`
	Scores         ScoreMap `json:"O,omitempty"` //Using string instead of uint as an index because JSON doesn't support map[uint]Score
	LinkedID       uint     `json:"l,omitempty"` //For duplicating shooters that are in different classes with the same score
	EID            string   `json:"V,omitempty"` //Points to EventBucket Shooter ID
	Hidden         bool     `json:"h,omitempty"`
	Disabled       bool     `json:"d,omitempty"`
	Sex            bool     `json:"x,omitempty"`
	GradeSeparator bool     `json:"-"` //Used for enterShots and enterRangeTotals only.
	//position  string           `json:"-"` //Used for scoreboard only.
	//warning   uint8            `json:"-"` //Used for scoreboard only.
	//0 = nil
	//1 = highest possible score
	//2 = no score
	//3 = incomplete
	//4 = shoot off
}

//Shooter is exported
type Shooter struct {
	ID        string           `json:"I"` //EventBucket shooters ID
	SID       uint             `json:"J,omitempty"`
	NID       uint             `json:"N,omitempty"` //NRAA sequential integer id.
	FirstName string           `json:"f,omitempty"`
	Surname   string           `json:"s,omitempty"`
	NickName  string           `json:"n,omitempty"`
	Club      string           `json:"C,omitempty"`
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
	Grades    []uint    `json:"G,omitempty"`
	Sex       bool      `json:"x,omitempty"`
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

type serverSettings struct {
	theme bool
}
