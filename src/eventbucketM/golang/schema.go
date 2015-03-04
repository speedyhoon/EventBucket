package main

//TODO research would changing map[string]interface{} to M{} have any affect on the size of the compiled application?
type M map[string]interface{}

/* Schema Rules:
lowercase letters MUST be used for struct properties
Uppercase letters MUST be used for a sub struct*/
type Event struct {
	Id             string         `bson:"_id"`
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
	AutoInc        AutoInc        `bson:"U"`
	//	TeamCat        map[string]TeamCat      `bson:"A,omitempty"`
	//	Teams          map[string]Team         `bson:"T,omitempty"`
	//Datetime string				`bson:"d,omitempty"`		No browser currently supports date time, so settling for separate fields that google chrome allows
}

type AutoInc struct {
	Mound   int `bson:"M,omitempty"`
	Event   int `bson:"E,omitempty"`
	Club    int `bson:"C,omitempty"`
	Range   int `bson:"^^schemaRANGE^^,omitempty"`
	Shooter int `bson:"^^schemaSHOOTER^^,omitempty"`
}

type Club struct {
	Id        string  `bson:"_id"`
	Name      string  `bson:"n"`
	LongName  string  `bson:"l,omitempty"`
	Mounds    []Mound `bson:"M,omitempty"`
	Latitude  string  `bson:"t,omitempty"`
	Longitude string  `bson:"g,omitempty"`
	Url       string  `bson:"u,omitempty"`
	Address   string  `bson:"a,omitempty"`
	Town      string  `bson:"w,omitempty"`
	PostCode  string  `bson:"p,omitempty"`
	AutoInc   AutoInc `bson:"^^schemaAutoInc^^,omitempty"`
}

type Range struct {
	Name       string                   `bson:"n"`
	Aggregate  string                   `bson:"a,omitempty"`
	ScoreBoard bool                     `bson:"s,omitempty"`
	Locked     bool                     `bson:"l,omitempty"`
	Hidden     bool                     `bson:"h,omitempty"`
	Order      int                      `bson:"^^schemaSORT^^,omitempty"`
	Status     int                      `bson:"t,omitempty"`      //ENUM change to 1 when the first shooter has recorded their first shot change to 2 when the range is finished. http://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-golang
	Class      map[string]RangeProperty `bson:"omitempty,inline"` //TODO possibly change it to optional grades per range in future
	Id         *int                     `bson:"i,omitempty"`
}
type RangeProperty struct {
	ShotsQty    int `bson:"s,omitempty"`
	SightersQty int `bson:"i,omitempty"`
}

type Score struct {
	//TODO the schema should change so that it can use unsigned 64 bit numbers instead
	Total int `bson:"t"`
	//	Total uint64								`bson:"t"`
	Shots   string `bson:"s,omitempty"`
	Centers int    `bson:"c,omitempty"`
	//	Centers uint64								`bson:"c"`
	CountBack1 string `bson:"v,omitempty"`
	//	CountBack2 string							`bson:"x,omitempty"`
	//	Xs string									`bson:"u,omitempty"` //This might be handy for the future?
	Position int `bson:"p,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
}

type NRAA_Shooter struct {
	SID       int    `bson:"_id,omitempty"`
	NRAA_Id   int    `bson:"i,omitempty"`
	Surname   string `bson:"s,omitempty"`
	FirstName string `bson:"f,omitempty"`
	NickName  string `bson:"n,omitempty"`
	Club      string `bson:"c,omitempty"`
}

type EventShooter struct {
	FirstName string           `bson:"f"` //TODO change these to point to shooters in the other shooter tables
	Surname   string           `bson:"s"`
	Club      string           `bson:"b"` //TODO should possibly change to "C"??
	Grade     int              `bson:"g"`
	Hidden    bool             `bson:"h,omitempty"`
	AgeGroup  string           `bson:"a,omitempty"`
	Scores    map[string]Score `bson:"omitempty,inline"` //S is not used!
	LinkedId  *int             `bson:"l,omitempty"`      //For duplicating shooters that are in different classes with the same score
	SID       int              `bson:"d,omitempty"`
	//SCOREBOARD
	Id       int    `bson:"i,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	Position string `bson:"x,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	Warning  int8   `bson:"y,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//		0 = nil
	//		1 = shoot off
	//		2 = warning, no score
	//		3 = incomplete
	//		4 = highest posible score

	//START-SHOOTING & TOTAL-SCORES
	GradeSeparator bool `bson:"z,omitempty"` //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
	ClassSeparator bool `bson:"o,omitempty"` //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
	//	Id string									`bson:"w,omitempty"`//DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
}

type Shooter struct {
	SID int `bson:"_id,omitempty"`
	//	SID int					`bson:"s,omitempty"`
	NRAA_Id   int    `bson:"i,omitempty"`
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
}

type TeamCat struct {
	Name string `bson:"n"`
}

type Team struct {
	name     string `bson:"n"`
	teamCat  []int  `bson:"t"`
	shooters []int  `bson:"s,omitempty"`
}

type Mound struct {
	Id       int    `bson:"i"`
	Distance int    `bson:"d"`
	Unit     string `bson:"u"`
	Name     string `bson:"n,omitempty"`
	Notes    string `bson:"o,omitempty"`
}

type Skill struct {
	Grade      string
	Percentage float64 //TODO would prefer an unsigned float here
}
