package main

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

//		ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890abcdefghijklmnopqrstuvwxyz
//		AB   FGHIJKLM OPQR   VWXYZ01234567890abcdefghijklmnopqrstuvwxyz
const (
	schemaID						= "_id"
	schemaAUTOINC				= "Y"
//Event
	schemaNAME				= "N"
	schemaDATETIME			= "D"
	schemaCLUB				= "C"		//Club_id
	schemaSHOOTER			= "S"
//Ranges
	schemaRANGE				= "R"
		schemaAGG			= "a"
		schemaDISABLE		= "d"
		schemaSCOREBOARD	= "s"
//Shooter
		schemaFIRSTNAME	= "f"
		schemaMIDDLENAME	= "m"
		schemaSURNAME		= "u"
		//schemaCLUB		= "C"		//Club_id
		schemaCLASS			= "c"
		schemaGRADE			= "g"
		schemaAGE			= "a"
	//Scores
		schemaTOTAL			= "t"
		schemaCENTER		= "c"
		schemaVCOUNTBACK	= "v"
		schemaXCOUNTBACK	= "x"
		schemaSHOT			= "s"
//Teams
	schemaTEAMCAT			= "E"
	schemaTEAM				= "T"

//Club
	schemaURL				= "U"
	schemaRANGEDISTANCE	= "d"
	schemaUNIT				= "j"
	schemaLATITUDE			= "t"
	schemaLONGITUDE		= "g"
	//Ranges
		// ...
)

//const (
//	schemaID						= "_id"
//	schemaAUTOINC				= "Y"
////Event
//	schemaEV_NAME				= "N"
//	schemaEV_DATETIME			= "D"
//	schemaEV_CLUB				= "C"		//Club_id
//	schemaEV_SHOOTER			= "S"
//		schemaEV_CLASS			= "c"
//		schemaEV_GRADE			= "g"
//		schemaEV_AGE			= "a"
//		schemaEV_FIRSTNAME	= "f"
//		schemaEV_SURNAME		= "u"
//
//		schemaEV_TOTAL			= "t"
//		schemaEV_CENTER		= "c"
//		schemaEV_VCOUNTBACK	= "v"
//		schemaEV_XCOUNTBACK	= "x"
//		schemaEV_SHOT			= "s"
//	schemaEV_RANGE				= "R"
//		schemaEV_AGG			= "a"
//		schemaEV_DISABLE		= "d"
//		schemaEV_SCOREBOARD	= "s"
////Teams
//	schemaEV_TEAMCAT			= "E"
//	schemaEV_TEAM				= "T"
//
////Club
//	schemaCL_URL				= "U"
//	schemaCL_RANGEDISTANCE	= "d"
//	schemaCL_UNIT				= "j"
//	schemaCL_LATITUDE			= "t"
//	schemaCL_LONGITUDE		= "g"
//
////Shooter
//	schemaSH_FIRSTNAME		= "f"
//	schemaSH_MIDDLETNAME		= "m"
//	schemaSH_SURNAME			= "u"
//	schemaSH_CLUB				= "C"		//Club_id
//)



//import("mgo/bson")

//const (
////	EVENT
//	Id = "_id"
//	Ename = "N"
//	Eid = "I"
//	Edatetime = "D"
//	Eshooters = "S"
//	Eclub = "C"
////		EshootersId = Eid
//		EshootersClass = "C"
//		EshootersGrade = "G"
////		Eranges = "R"
//		Etotal = "T"
//		Ecenters = "C"
//		EvCountBack = "V"
//		ExCountBack = "X"
//		Eshots = "S"
//	Eranges = "R"
////		Ename = "N"
//		Eaggs = "A"
//		Etype = "T"
//		Escoreboard = "S"
//		Edisabled = "D"
//	EteamCat = "E"
////		Ename = "N"
//	Eteam = "T"
////		Ename = "N"
////		Eshooters = "S"
//
////CLUB
//	Cname = "N"
//	Cid = "I"
//	Curl = "U"
//	Clatitude = "A"
//	Clongitude = "G"
//)



type TeamCat struct{
	name string
}
type Team struct{
	teamCat []int
	name string
	shooters []int
}
type Club struct{
//	Id bson.ObjectId `bson:"_id"`
	name string`bson:"n"`
	rego string`bson:"r"`
	url string`bson:"u"`
}

type Event struct{
	Id string `bson:"_id"`
	Club string `bson:"c"`
	Name string `bson:"n"`
	Datetime string `bson:"d"`
//	Ranges Ranges `bson:"R"`
//	Ranges []Ranges `bson:"R,inline"`
	Ranges map[string]Ranges `bson:"R,inline"`
//	Shooters []EventShooter `bson:"S"`
//	Shooters []EventShooter `bson:"S,inline"`
	Shooters map[string]EventShooter `bson:"S,inline"`
//	TeamCat []TeamCat `bson:"t"`
//	TeamCat []TeamCat `bson:"A,inline"`
	TeamCat map[string]TeamCat `bson:"A,inline"`
//	Teams []Team `bson:"T"`
//	Teams []Team `bson:"T,inline"`
	Teams map[string]Team `bson:"T,inline"`
}
type Ranges struct{
	id int
	Name string `bson:"N"`
	Type string `bson:"T"`
	Aggs []int
	ScoreBoard bool
	Enabled bool
}
type EventShooter struct{
	class, grade string
	scores []RangeScores
}

type RangeScores struct{
	total, centers, vcountBack, xcountBack int
	shots string
}
type Grades struct{
	class, grade string
}

//func try(w http.ResponseWriter, r *http.Request){
//	var newEvent = mpa[string]interface{}{
//		"_id": "abc",
//		Ranges: map[string]Ranges{
//			"0":	Ranges{
//				Name: "hello",
//				Type: "agg",
//			},
//			"1":	Ranges{
//				Name: "Neo",
//				Type: "range",
//			},
//		},
//		//	Ranges: []Ranges{
//		//		0:	Ranges{
//		//			Name: "hello",
//		//			Type: "agg",
//		//		},
//		//		1:	Ranges{
//		//			Name: "Neo",
//		//			Type: "range",
//		//		},
//		//	},
//		Shooters: map[string]EventShooter{
//			"0":	EventShooter{
//				class: "Target",
//				grade: "A",
//			},
//			"1":	EventShooter{
//				class: "Fclass",
//				grade: "FA",
//			},
//		},
//		//	Shooters: []EventShooter{
//		//		0:	EventShooter{
//		//			class: "Target",
//		//			grade: "A",
//		//		},
//		//		1:	EventShooter{
//		//			class: "Fclass",
//		//			grade: "FA",
//		//		},
//		//	},
//	}
//	checkErr(conn.C("event").Insert(newEvent))
//}




//var _ = map[string]interface{}{
//	"_id":	"123456789012345678901234",
//	Ename:	"Event Name",
//	Eid:		"aB",
//	Edatetime:	"389044567",
//	Eshooters:	map[string]interface{}{
//		"shooterId": map[string]interface{}{
//			EshootersClass:	"Target",
//			EshootersGrade:	"C",
//			"Escores":		map[string]interface{}{
//				"rangeID":	map[string]interface{}{
//					Etotal:	49,
//					Ecenters:4,
//					Eshots:	"345XV5555455",
//					EvCountBack: 554555566543,
//					ExCountBack: 554555567543,
//				},
//			},
//		},
//	},
//}
