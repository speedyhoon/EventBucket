package main

import("mgo/bson")

const (
	Ename = "N"
	Eid = "0"
	Edatetime = "D"
	Eshooters = "s"
		EshootersId = "0"
		EshootersClass = "c"
		EshootersGrade = "g"
	Escores = "S"
		EscoresTotal = "T"
		EscoresCenters = "C"
		EscoresVCountBack = "V"
		EscoresXCountBack = "X"
		Eshots = "S"
)

type Club struct{
//	Id bson.ObjectId `bson:"_id"`
	name, rego, url string
}

type Event struct{
	Id bson.ObjectId `bson:"_id"`
	id, name string
	datetime string
	shooters []eventShooter
}
type ranges struct{
	id int
	name string
	aggs []int
	scoreBoard, enabled bool
}
type eventShooter struct{
	class, grade string
	scores []rangeScores
}

type rangeScores struct{
	total, centers, vcountBack, xcountBack int
	shots string
}
type grades struct{
	class, grade string
}
type teamCat struct{
	name string
}
type team struct{
	teamCat []int
	name string
	shooters []int
}

var _ = map[string]interface{}{
	"_id":	"123456789012345678901234",
	Ename:	"Event Name",
	Eid:		"aB",
	Edatetime:	"389044567",
	Eshooters:	map[string]interface{}{
		"shooterId": map[string]interface{}{
			EshootersClass:	"Target",
			EshootersGrade:	"C",
			Escores:		map[string]interface{}{
				"rangeID":	map[string]interface{}{
					EscoresTotal:	49,
					EscoresCenters:4,
					Eshots:	"345XV5555455",
					EscoresVCountBack: 554555566543,
					EscoresXCountBack: 554555567543,
				},
			},
		},
	},
}
