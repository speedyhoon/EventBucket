package main

import (
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
	getDisciplines       uint8 = 254
	pageError            uint8 = 255
)

func dataListAgeGroup() []option {
	return []option{
		{Value: "0", Label: "None"},
		{Value: "1", Label: "Junior U21"},
		{Value: "2", Label: "Junior U25"},
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

func getForm(id uint8) []field {
	switch id {
	case 1://New Club
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}}
	case 2://Club Details
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "a", maxLen: 64, v8: isValidStr,
		}, {
			name: "w", maxLen: 64, v8: isValidStr,
		}, {
			name: "p", maxLen: 64, v8: isValidStr,
		}, {
			name: "x", min: -90, max: 90, step: .000001, v8: isValidFloat32,
		}, {
			name: "y", min: -180, max: 180, step: .000001, v8: isValidFloat32,
		}, {
			name: "b", v8: isValidBool,
		}, {
			name: "u", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", v8: isValidRegex, regex: regexID,
		}}
	case 3://New Shooting Mound
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", v8: isValidRegex, regex: regexID,
		}}
	case 4://New Event
		return []field{{
			name: "C", Value: defaultClubName(), Required: hasDefaultClub(), maxLen: 64, minLen: 1, v8: isValidStr, Options: clubsDataList(),
		}, {
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "d", Value: defaultDate(), maxLen: 10, v8: isValidStr,
		}, {
			name: "t", Value: defaultTime(), maxLen: 5, step: 300, v8: isValidStr,
		}}
	case 5://Event Details
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr, Options: clubsDataList(),
		}, {
			name: "d", maxLen: 10, minLen: 1, v8: isValidStr,
		}, {
			name: "t", maxLen: 5, minLen: 5, v8: isValidStr,
		}, {
			name: "c", v8: isValidBool,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 6://Add Range
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 7://Add Aggregate Range
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "R", Required: true, minLen: 2, min: 1, max: 65535, step: 1, v8: listUint,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 8://Shooter Entry
		return []field{{
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Value: defaultClubName(), Required: true, maxLen: 64, minLen: 1, v8: isValidStr, Options: clubsDataList(), Placeholder: defaultClubName(),
		}, {
			name: "S", v8: isValidStr, Options: searchShootersOptions("", "", defaultClubName()),
		}, {
			name: "g", Required: true, minLen: 1, max: float32(len(globalGrades)-1), step: 1, v8: listUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "x", v8: isValidBool,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 9://Existing Shooter Entry
		return []field{{
			name: "S", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "g", Required: true, minLen: 1, max: float32(len(globalGrades)-1), step: 1, v8: listUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}}
	case 10://Shooter Search
		return []field{{
			name: "f", maxLen: 64, v8: isValidStr,
		}, {
			name: "s", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", maxLen: 64, v8: isValidStr,
		}}
	case 11://New Shooter
		return []field{{
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Value: defaultClubName(), Required: true, maxLen: 64, minLen: 1, v8: isValidStr, Options: clubsDataList(),
		}, {
			name: "g", Required: true, max: float32(len(globalGrades)-1), step: 1, v8: listUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "x", v8: isValidBool,
		}}
	case 12://Shooter Details
		return []field{{
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "g", Required: true, max: float32(len(globalGrades)-1), step: 1, v8: listUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "x", v8: isValidBool,
		}, {
			name: "I", Required: true, v8: isValidRegex, regex: regexID,
		}}
	case 13://Shooter Search
		return []field{{
			name: "f", maxLen: 64, v8: isValidStr,
		}, {
			name: "s", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", maxLen: 64, v8: isValidStr,
		}}
	case 14://Enter Range Totals
		return []field{{
			name: "t", Required: true, max: 60, step: 1, v8: isValidUint,
		}, {
			name: "c", max: 10, step: 1, v8: isValidUint,
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "R", Required: true, min: 1, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "S", Required: true, max: 65535, step: 1, v8: isValidUint,
		}}
	case 15://Grades Available
		return []field{{
			name: "g", Required: true, minLen: 1, max: float32(len(globalGrades)-1), step: 1, v8: listUint, Options: availableGrades([]uint{}),
		}, {
			name: "I", v8: isValidRegex, regex: regexID,
		}}
	case 16://Update Shooter Shots (Scorecards)
		return []field{{
			name: "s", Required: true, maxLen: 12, minLen: 1, v8: isValidStr,
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "R", Required: true, min: 1, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "S", Required: true, max: 65535, step: 1, v8: isValidUint,
		}}
	case 17://Import Shooters
		return []field{{
			name: "f", Required: true, maxLen: 64,
		}}
	case 18://Map Clubs
		return []field{{
			name: "C", v8: isValidRegex, regex: regexID,
		}}
	case 19://Edit Shooting Mound
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "I", max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "C", v8: isValidRegex, regex: regexID,
		}}
	case 20://Update Range
		return []field{{
			name: "E", v8: isValidRegex, regex: regexID,
		}, {
			name: "I", Required: true, min: 1, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "k", v8: isValidBool,
		}}
	case 21://Update Agg
		return []field{{
			name: "E", v8: isValidRegex, regex: regexID,
		}, {
			name: "I", Required: true, min: 1, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "R", Required: true, minLen: 2, min: 1, max: 65535, step: 1, v8: listUint,
		}}
	case 22://Entries Edit Shooter Details
		return []field{{
			name: "S", Required: true, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "g", Required: true, max: float32(len(globalGrades)-1), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "x", v8: isValidBool,
		}, {
			name: "k", v8: isValidBool,
		}}
	}
	return []field{}
}
