package main

import (
	"time"
)

type form struct {
	action uint8
	Fields []field
	Error  string
	expiry time.Time
}

type option struct {
	Label, Value string
	Selected     bool
}

const (
	fieldMaxLen = 64

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
	eventTotal           uint8 = 13
)

func dataListAgeGroup() []option {
	return []option{
		{},
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
	case 0:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "b", v8: isValidBool,
		}}
	case 1:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
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
			name: "C", v8: isValidID, regex: regexID,
		}}
	case 2:
		return []field{{
			name: "e", Required: true, min: .01, max: 65535, step: .01, v8: isValidFloat32,
		}, {
			name: "z", Required: true, v8: isValidStr, Options: []option{{Label: "Metres", Value: "m", Selected: false}, {Label: "Kilometres", Value: "k", Selected: false}, {Label: "Yards", Value: "y", Selected: false}, {Label: "Feet", Value: "f", Selected: false}},
		}, {
			name: "C", v8: isValidID, regex: regexID,
		}}
	case 3:
		return []field{{
			name: "C", Value: defaultClubName(), Required: hasDefaultClub(), maxLen: 64, v8: isValidStr, Options: getDataListClubs(),
		}, {
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "d", Value: defaultDate(), maxLen: 10, v8: isValidStr,
		}, {
			name: "t", Value: defaultTime(), maxLen: 5, step: 300, v8: isValidStr,
		}}
	case 4:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "d", maxLen: 10, v8: isValidStr,
		}, {
			name: "t", maxLen: 5, v8: isValidStr,
		}, {
			name: "c", v8: isValidBool,
		}, {
			name: "E", v8: isValidID, regex: regexID,
		}}
	case 5:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "E", v8: isValidID, regex: regexID,
		}}
	case 6:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "R", Required: true, min: 1, max: 999, step: 1, v8: listUint,
		}, {
			name: "E", v8: isValidID, regex: regexID,
		}}
	case 7:
		return []field{{
			name: "f", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "S", v8: isValidStr,
		}, {
			name: "g", Required: true, min: 1, max: float32(len(globalGradesDataList)), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", min: 1, max: 5, step: 1, v8: isValidUint, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Junior U21", Value: "1", Selected: false}, {Label: "Junior U25", Value: "2", Selected: false}, {Label: "Veteran", Value: "3", Selected: false}, {Label: "Super Veteran", Value: "4", Selected: false}},
		}, {
			name: "E", Required: true, v8: isValidID, regex: regexID,
		}, {
			name: "E", v8: isValidID, regex: regexID,
		}}
	case 8:
		return []field{{
			name: "S", Required: true, v8: isValidID, regex: regexID,
		}, {
			name: "g", Required: true, min: 1, max: float32(len(globalGradesDataList)), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", min: 1, max: 5, step: 1, v8: isValidUint, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Junior U21", Value: "1", Selected: false}, {Label: "Junior U25", Value: "2", Selected: false}, {Label: "Veteran", Value: "3", Selected: false}, {Label: "Super Veteran", Value: "4", Selected: false}},
		}, {
			name: "E", Required: true, v8: isValidID, regex: regexID,
		}}
	case 9:
		return []field{{
			name: "f", maxLen: 64, v8: isValidStr,
		}, {
			name: "s", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", maxLen: 64, v8: isValidStr,
		}}
	case 10:
		return []field{{
			name: "f", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "g", Required: true, min: 1, max: float32(len(globalGradesDataList)), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", min: 1, max: 5, step: 1, v8: isValidUint, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Junior U21", Value: "1", Selected: false}, {Label: "Junior U25", Value: "2", Selected: false}, {Label: "Veteran", Value: "3", Selected: false}, {Label: "Super Veteran", Value: "4", Selected: false}},
		}}
	case 11:
		return []field{{
			name: "f", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "g", Required: true, min: 1, max: float32(len(globalGradesDataList)), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", min: 1, max: 5, step: 1, v8: isValidUint, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Junior U21", Value: "1", Selected: false}, {Label: "Junior U25", Value: "2", Selected: false}, {Label: "Veteran", Value: "3", Selected: false}, {Label: "Super Veteran", Value: "4", Selected: false}},
		}, {
			name: "I", Required: true, maxLen: 64, v8: isValidID, regex: regexID,
		}}
	case 12:
		return []field{{
			name: "f", maxLen: 64, v8: isValidStr,
		}, {
			name: "s", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", maxLen: 64, v8: isValidStr,
		}}
	case 13:
		return []field{{
			name: "t", Required: true, max: 60, step: 1, v8: isValidUint,
		}, {
			name: "c", Required: true, max: 10, step: 1, v8: isValidUint,
		}, {
			name: "E", step: 1, v8: isValidID, regex: regexID,
		}, {
			name: "R", min: 1, step: 1, v8: isValidUint,
		}, {
			name: "S", step: 1, v8: isValidUint,
		}}
	}
	return []field{}
}
