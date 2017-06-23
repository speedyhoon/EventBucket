package main

import (
	"math"
	"time"
)

//TODO when forms are changed to zero index, perform a check if action is not equal to "" instead of strconv != 0
const (
	clubNew              uint8 = 0
	clubDetails          uint8 = 1
	clubMoundNew         uint8 = 2
	eventNew             uint8 = 3
	eventDetails         uint8 = 4
	eventRangeNew        uint8 = 5
	eventAggNew          uint8 = 6
	eventShooterNew      uint8 = 7
	eventShooterExisting uint8 = 8
	eventShooterSearch   uint8 = 9
	shooterNew           uint8 = 10
	shooterDetails       uint8 = 11
	shooterSearch        uint8 = 12
	eventTotalScores     uint8 = 13
	eventAvailableGrades uint8 = 14
	eventUpdateShotScore uint8 = 15
	importShooter        uint8 = 16
	mapResults           uint8 = 17
	clubMoundEdit        uint8 = 18
	eventUpdateRange     uint8 = 19
	eventUpdateAgg       uint8 = 20
	eventEditShooter     uint8 = 21
)

func getForm(id uint8) form {
	const formQty = 22
	if id >= formQty {
		return form{}
	}
	return form{
		action: id,
		Fields: [formQty][]field{{ //0 New Club
			{name: "n", v8: v8StrReq},
		}, { //1 Club Details
			{name: "n", v8: v8StrReq},
			{name: "a", v8: v8Str},
			{name: "w", v8: v8Str},
			{name: "p", v8: v8Str},
			//TODO add a struct isValid{ max: 64}
			//TODO add methods to isValid struct so they can be passed as a value like: isValid.Float32
			{name: "x", v8: v8Float32, min: -90, max: 90, step: .000001},
			{name: "y", v8: v8Float32, min: -180, max: 180, step: .000001},
			{name: "b", v8: v8Bool},
			{name: "u", v8: v8Str},
			{name: "C", v8: v8Regex, regex: regexID},
		}, { //2 New Shooting Mound
			{name: "n", v8: v8StrReq},
			{name: "C", v8: v8Regex, regex: regexID},
		}, { //3 New Event
			{name: "C", v8: v8Str, Value: defaultClubName(), Required: hasDefaultClub(), minLen: 1, Options: clubsDataList()},
			{name: "n", v8: v8StrReq},
			{name: "d", v8: v8Str, Value: time.Now().Format("2006-01-02"), maxLen: 10},
			{name: "t", v8: v8Str, Value: time.Now().Format("15:04"), maxLen: 5},
		}, { //4 Event Details
			{name: "C", v8: v8StrReq, Options: clubsDataList()},
			{name: "n", v8: v8StrReq},
			{name: "d", v8: v8Str, maxLen: 10, minLen: 1},
			{name: "t", v8: v8Str, maxLen: 5, minLen: 5},
			{name: "c", v8: v8Bool},
			{name: "E", v8: v8Regex, regex: regexID},
		}, { //5 Add Range
			{name: "n", v8: v8StrReq},
			{name: "E", v8: v8Regex, regex: regexID},
		}, { //6 Add Aggregate Range
			{name: "n", v8: v8StrReq},
			{name: "R", v8: v8UintList, Required: true, minLen: 2, min: 1, max: math.MaxInt8},
			{name: "E", v8: v8Regex, regex: regexID},
		}, { //7 Shooter Entry
			{name: "f", v8: v8StrReq},
			{name: "s", v8: v8StrReq},
			{name: "C", v8: v8StrReq, Options: clubsDataList(), Placeholder: defaultClubName()},
			{name: "S", v8: v8Str, Options: searchShootersOptions("", "", defaultClubName())},
			{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
			{name: "x", v8: v8Bool},
			{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
			{name: "E", v8: v8Regex, regex: regexID},
		}, { //7 Existing Shooter Entry
			{name: "S", v8: v8RegexReq, regex: regexID},
			{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
			{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
			{name: "E", v8: v8RegexReq, regex: regexID},
		}, { //8 Shooter Search
			{name: "f", v8: v8Str},
			{name: "s", v8: v8Str},
			{name: "C", v8: v8Str},
		}, { //9 New Shooter
			{name: "f", v8: v8StrReq},
			{name: "s", v8: v8StrReq},
			//TODO change datalist functions to variables do that the number of database calls can be reduced? But this would use more RAM which is not good
			{name: "C", v8: v8StrReq, Value: defaultClubName(), Options: clubsDataList(), Placeholder: defaultClubName()},
			{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
			{name: "x", v8: v8Bool},
			{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
		}, { //10 Shooter Details
			{name: "f", v8: v8StrReq},
			{name: "s", v8: v8StrReq},
			{name: "C", v8: v8StrReq},
			{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
			{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
			{name: "x", v8: v8Bool},
			{name: "I", v8: v8RegexReq, regex: regexID},
		}, { //11 Shooter Search
			{name: "f", v8: v8Str},
			{name: "s", v8: v8Str},
			{name: "C", v8: v8Str},
		}, { //12 Enter Range Totals
			{name: "t", v8: v8UintReq, max: 120},
			{name: "c", v8: v8Uint, max: 20},
			{name: "E", v8: v8RegexReq, regex: regexID},
			{name: "R", v8: v8UintReq, min: 1, max: math.MaxInt8},
			{name: "S", v8: v8UintReq, max: math.MaxInt8},
			{name: "h", v8: v8Uint, max: 100},
		}, { //13 Grades Available
			{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: availableGrades([]uint{})},
			{name: "I", v8: v8Regex, regex: regexID},
		}, { //14 Update Shooter Shots (Scorecards)
			{name: "s", v8: v8StrReq, maxLen: 12, minLen: 1},
			{name: "E", v8: v8RegexReq, regex: regexID},
			{name: "R", v8: v8UintReq, min: 1, max: math.MaxInt8},
			{name: "S", v8: v8UintReq, max: math.MaxInt8},
		}, { //15 Import Shooters
			//TODO add file validation
			{name: "f", Required: true},
		}, { //16 Map Clubs
			{name: "C", v8: v8Regex, regex: regexID},
		}, { //17 Edit Shooting Mound
			{name: "n", v8: v8StrReq},
			{name: "I", v8: v8Uint, max: math.MaxInt8},
			{name: "C", v8: v8Regex, regex: regexID},
		}, { //18 Update Range
			{name: "E", v8: v8Regex, regex: regexID},
			{name: "I", v8: v8UintReq, min: 1, max: math.MaxInt8},
			{name: "n", v8: v8StrReq},
			{name: "k", v8: v8Bool},
			{name: "o", v8: v8UintReq, max: math.MaxInt8},
		}, { //20 Update Agg
			{name: "E", v8: v8Regex, regex: regexID},
			{name: "I", v8: v8UintReq, min: 1, max: math.MaxInt8},
			{name: "n", v8: v8StrReq},
			{name: "R", v8: v8UintList, Required: true, minLen: 2, min: 1, max: math.MaxInt8},
			{name: "o", v8: v8UintReq, max: math.MaxInt8},
		}, { //21 Entries Edit Shooter Details
			{name: "S", v8: v8UintReq, max: math.MaxInt8},
			{name: "E", v8: v8RegexReq, regex: regexID},
			{name: "f", v8: v8StrReq},
			{name: "s", v8: v8StrReq},
			{name: "C", v8: v8StrReq},
			//TODO i don't think using a max & min is incorrect here. instead maybe pass in a list of available/valid ids instead?
			{name: "g", v8: v8UintReq, max: len(globalGrades) - 1, Options: globalGradesDataList},
			{name: "r", v8: v8Uint, max: len(dataListAgeGroup()) - 1, Options: dataListAgeGroup()},
			{name: "x", v8: v8Bool},
			{name: "k", v8: v8Bool},
		},
		}[id],
	}
}
