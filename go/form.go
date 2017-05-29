package main

import (
	"math"
	"time"
)

type form struct {
	action uint8
	Fields []field
	Error  error
	expiry time.Time
}

type option struct {
	Label, Value string
	Selected     bool
}

const (
	clubNew              uint8 = 1
	clubDetails          uint8 = 2
	clubMoundNew         uint8 = 3
	eventNew             uint8 = 4
	eventDetails         uint8 = 5
	eventRangeNew        uint8 = 6
	eventAggNew          uint8 = 7
	eventShooterNew      uint8 = 8
	eventShooterExisting uint8 = 9
	eventShooterSearch   uint8 = 10
	shooterNew           uint8 = 11
	shooterDetails       uint8 = 12
	shooterSearch        uint8 = 13
	eventTotalScores     uint8 = 14
	eventAvailableGrades uint8 = 15
	eventUpdateShotScore uint8 = 16
	importShooter        uint8 = 17
	mapResults           uint8 = 18
	clubMoundEdit        uint8 = 19
	eventUpdateRange     uint8 = 20
	eventUpdateAgg       uint8 = 21
	eventEditShooter     uint8 = 22
	maxLen               int   = 64
	pageError            uint8 = 255
)

func dataListAgeGroup() []option {
	return []option{
		{Value: "0", Label: "None"},
		{Value: "1", Label: "U21"},
		{Value: "2", Label: "U25"},
		{Value: "3", Label: "Veteran"},
		{Value: "4", Label: "Super Veteran"},
	}
}

func defaultDate() string {
	return time.Now().Format("2006-01-02")
}

func defaultTime() string {
	return time.Now().Format("15:04")
}

func getForm(id uint8) form {
	const formQty = 23
	if id >= formQty {
		return form{}
	}
	return form{
		action: id,
		Fields: [formQty][]field{
			{}, { //New Club
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
			}, { //Club Details
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "a", maxLen: maxLen, v8: isValidStr},
				{name: "w", maxLen: maxLen, v8: isValidStr},
				{name: "p", maxLen: maxLen, v8: isValidStr},
				{name: "x", min: -90, max: 90, step: .000001, v8: isValidFloat32},
				{name: "y", min: -180, max: 180, step: .000001, v8: isValidFloat32},
				{name: "b", v8: isValidBool},
				{name: "u", maxLen: maxLen, v8: isValidStr},
				{name: "C", v8: isValidRegex, regex: regexID},
			}, { //New Shooting Mound
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "C", v8: isValidRegex, regex: regexID},
			}, { //New Event
				{name: "C", Value: defaultClubName(), Required: hasDefaultClub(), maxLen: maxLen, minLen: 1, v8: isValidStr, Options: clubsDataList()},
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "d", Value: defaultDate(), maxLen: 10, v8: isValidStr},
				{name: "t", Value: defaultTime(), maxLen: 5, v8: isValidStr},
			}, { //Event Details
				{name: "C", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr, Options: clubsDataList()},
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "d", maxLen: 10, minLen: 1, v8: isValidStr},
				{name: "t", maxLen: 5, minLen: 5, v8: isValidStr},
				{name: "c", v8: isValidBool},
				{name: "E", v8: isValidRegex, regex: regexID},
			}, { //Add Range
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "E", v8: isValidRegex, regex: regexID},
			}, { //Add Aggregate Range
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "R", Required: true, minLen: 2, min: 1, max: math.MaxInt8, step: 1, v8: listUint},
				{name: "E", v8: isValidRegex, regex: regexID},
			}, { //Shooter Entry
				{name: "f", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "s", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "C", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr, Options: clubsDataList(), Placeholder: defaultClubName()},
				{name: "S", v8: isValidStr, Options: searchShootersOptions("", "", defaultClubName())},
				{name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup()},
				{name: "x", v8: isValidBool},
				{name: "g", Required: true, minLen: 1, max: float32(len(globalGrades) - 1), step: 1, v8: listUint, Options: globalGradesDataList},
				{name: "E", v8: isValidRegex, regex: regexID},
			}, { //Existing Shooter Entry
				{name: "S", Required: true, v8: isValidRegex, regex: regexID},
				{name: "g", Required: true, minLen: 1, max: float32(len(globalGrades) - 1), step: 1, v8: listUint, Options: globalGradesDataList},
				{name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup()},
				{name: "E", Required: true, v8: isValidRegex, regex: regexID},
			}, { //Shooter Search
				{name: "f", maxLen: maxLen, v8: isValidStr},
				{name: "s", maxLen: maxLen, v8: isValidStr},
				{name: "C", maxLen: maxLen, v8: isValidStr},
			}, { //New Shooter
				{name: "f", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "s", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "C", Value: defaultClubName(), Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr, Options: clubsDataList(), Placeholder: defaultClubName()},
				{name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup()},
				{name: "x", v8: isValidBool},
				{name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: listUint, Options: globalGradesDataList},
			}, { //Shooter Details
				{name: "f", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "s", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "C", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: listUint, Options: globalGradesDataList},
				{name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup()},
				{name: "x", v8: isValidBool},
				{name: "I", Required: true, v8: isValidRegex, regex: regexID},
			}, { //Shooter Search
				{name: "f", maxLen: maxLen, v8: isValidStr},
				{name: "s", maxLen: maxLen, v8: isValidStr},
				{name: "C", maxLen: maxLen, v8: isValidStr},
			}, { //Enter Range Totals
				{name: "t", Required: true, max: 120, step: 1, v8: isValidUint},
				{name: "c", max: 20, step: 1, v8: isValidUint},
				{name: "E", Required: true, v8: isValidRegex, regex: regexID},
				{name: "R", Required: true, min: 1, max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "S", Required: true, max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "h", max: 100, step: 1, v8: isValidUint},
			}, { //Grades Available
				{name: "g", Required: true, minLen: 1, max: float32(len(globalGrades) - 1), step: 1, v8: listUint, Options: availableGrades([]uint{})},
				{name: "I", v8: isValidRegex, regex: regexID},
			}, { //Update Shooter Shots (Scorecards)
				{name: "s", Required: true, maxLen: 12, minLen: 1, v8: isValidStr},
				{name: "E", Required: true, v8: isValidRegex, regex: regexID},
				{name: "R", Required: true, min: 1, max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "S", Required: true, max: math.MaxInt8, step: 1, v8: isValidUint},
			}, { //Import Shooters
				{name: "f", Required: true, maxLen: maxLen},
			}, { //Map Clubs
				{name: "C", v8: isValidRegex, regex: regexID},
			}, { //Edit Shooting Mound
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "I", max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "C", v8: isValidRegex, regex: regexID},
			}, { //Update Range
				{name: "E", v8: isValidRegex, regex: regexID},
				{name: "I", Required: true, min: 1, max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "k", v8: isValidBool},
				{name: "o", Required: true, max: math.MaxInt8, step: 1, v8: isValidUint},
			}, { //Update Agg
				{name: "E", v8: isValidRegex, regex: regexID},
				{name: "I", Required: true, min: 1, max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "n", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "R", Required: true, minLen: 2, min: 1, max: math.MaxInt8, step: 1, v8: listUint},
				{name: "o", Required: true, max: math.MaxInt8, step: 1, v8: isValidUint},
			}, { //Entries Edit Shooter Details
				{name: "S", Required: true, max: math.MaxInt8, step: 1, v8: isValidUint},
				{name: "E", Required: true, v8: isValidRegex, regex: regexID},
				{name: "f", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "s", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "C", Required: true, maxLen: maxLen, minLen: 1, v8: isValidStr},
				{name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: isValidUint, Options: globalGradesDataList},
				{name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup()},
				{name: "x", v8: isValidBool},
				{name: "k", v8: isValidBool},
			},
		}[id],
	}
}
