package main

/* Schema Rules:
	lowercase letters MUST be used for struct properties
	Uppercase letters MUST be used for a sub struct
*/
type Event struct{
	Id string									`bson:"_id"`
	Club string									`bson:"c"`
	Name string									`bson:"n"`
	Datetime string							`bson:"d,omitempty"`
	Ranges map[string]Range					`bson:"R,omitempty"`
	Shooters map[string]EventShooter		`bson:"S,omitempty"`
	TeamCat map[string]TeamCat				`bson:"A,omitempty"`
	Teams map[string]Team					`bson:"T,omitempty"`
	AutoInc Autoinc							`bson:"U,omitempty"`
}

type Autoinc struct{
	Mound int									`bson:"M,omitempty"`
//	Event int									`bson:"E,omitempty"`
//	Club int										`bson:"C,omitempty"`
}

type Club struct{
	Id string									`bson:"_id"`
	Name string									`bson:"n"`
	Mounds map[string]Mound					`bson:"M,omitempty"`
	Latitude string							`bson:"t,omitempty"`
	Longitude string							`bson:"g,omitempty"`
	Url string									`bson:"u,omitempty"`
	AutoInc Autoinc				`bson:"U,omitempty"`
}

type Range struct{
	Name string									`bson:"n"`
	Aggregate []int							`bson:"a,omitempty"`
	ScoreBoard bool							`bson:"s,omitempty"`
	Enabled bool								`bson:"e,omitempty"`
}

type Score struct{
	total int64									`bson:"t"`
	shots string								`bson:"s"`
	centers int64								`bson:"c,omitempty"`
	VCountBack string							`bson:"v,omitempty"`
	XCountBack string							`bson:"x,omitempty"`
}

type EventShooter struct{
	FirstName string							`bson:"f"`
	Surname string								`bson:"s"`
	Class string								`bson:"c"`
	Grade string								`bson:"g"`
	AgeGroup string							`bson:"a,omitempty"`
	Scores map[string]Score					`bson:"S,omitempty"`
}

type Shooter struct{
	Id string									`bson:"n"`
	FirstName string							`bson:"f"`
	Surname string								`bson:"s"`
	Address string								`bson:"a"`
	Email string								`bson:"e"`
	Skill map[string]Skill	//Grading set by the VRA for each class
}

type TeamCat struct{
	Name string									`bson:"n"`
}

type Team struct{
	name string									`bson:"n"`
	teamCat []int								`bson:"t"`
	shooters []int								`bson:"s,omitempty"`
}

type Mound struct{
	Distance int								`bson:"d"`
	Unit string									`bson:"u"`
	Name string									`bson:"n,omitempty"`
	Notes string								`bson:"o,omitempty"`
}

type Skill struct{
	Grade string
	Class string
	Percentage float64 //TODO would prefer an unsigned float here
}

const(
	TBLautoinc = "autoinc"
	TBLclub = "club"
	TBLevent = "event"
	TBLchamp = "champ"		//Championship
)

//TODO remove all code below this line!!!!
func schema(key string)string{
	temp := map[string]string{
		"id": "_id",
		"autoinc": "U",

		//EVENT
		"name": "N",
		"datetime": "D",
		"club": "C",
		"shooters": "S",
		"shooter": "S",
			"class": "c",
			"grade": "g",
			"age": "a",

		//SCRORES
			"total": "t",
			"totals": "t",
			"centers": "c",
			"center": "c",
			"vCountBack": "v",
			"xCountBack": "x",
			"shots": "s",
			"shot": "s",
		"ranges": "R",
		"range": "R",
			"aggs": "a",
			"agg": "a",
		"scoreboard": "B",
		"disabled": "D",

		//Teams
		"teamCat": "E",
		"team": "T",
		"teams": "T",


		//CLUB
		//	"name": "N",
		"url": "L",
		"distance": "d",
		"unit": "J",
		"latitude": "A",
		"longitude": "G",


		//SHOOTER
		"firstname": "F",
		"surname": "M",
	}

	if out, ok := temp[key]; ok{
		return out
	}else{
		return ""
	}
}

const (
	schemaID					= "_id"
	schemaAUTOINC			= "U"
	schemaMOUND				= "M"
	schemaNAME				= "N"
	schemaCLUB				= "C"		//Club_id
	schemaSHOOTER			= "S"
	schemaRANGE				= "R"
	schemaAGG			= "a"
	schemaCLASS			= "c"
	schemaGRADE			= "g"
	schemaTOTAL			= "t"
	schemaCENTER		= "c"
	schemaLATITUDE			= "t"
	schemaLONGITUDE		= "g"
)
