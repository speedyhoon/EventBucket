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
	pageError            uint8 = 255
)

func dataListAgeGroup() []option {
	return []option{
		{Label: "None"},
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
	case 1:
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "b", v8: isValidBool,
		}}
	case 2:
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
	case 3:
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", v8: isValidRegex, regex: regexID,
		}}
	case 4:
		return []field{{
			name: "C", Value: defaultClubName(), Required: hasDefaultClub(), maxLen: 64, v8: isValidStr, Options: clubsDataList(),
		}, {
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "d", Value: defaultDate(), maxLen: 10, v8: isValidStr,
		}, {
			name: "t", Value: defaultTime(), maxLen: 5, step: 300, v8: isValidStr,
		}}
	case 5:
		return []field{{
			name: "n", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "d", maxLen: 10, minLen: 1, v8: isValidStr,
		}, {
			name: "t", maxLen: 5, minLen: 5, v8: isValidStr,
		}, {
			name: "c", v8: isValidBool,
		}, {
			name: "a", v8: isValidBool,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 6:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 7:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "R", Required: true, minLen: 2, min: 1, max: 65535, step: 1, v8: listUint,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 8:
		return []field{{
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr, Options: clubsDataList(),
		}, {
			name: "S", v8: isValidStr,
		}, {
			name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "x", v8: isValidBool,
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "E", v8: isValidRegex, regex: regexID,
		}}
	case 9:
		return []field{{
			name: "S", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}}
	case 10:
		return []field{{
			name: "f", maxLen: 64, v8: isValidStr,
		}, {
			name: "s", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", maxLen: 64, v8: isValidStr, manyRequired: []int{0, 1, 2}, manyRequiredQty: 1, Options: clubsDataList(),
		}}
	case 11:
		return []field{{
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}}
	case 12:
		return []field{{
			name: "f", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "s", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, minLen: 1, v8: isValidStr,
		}, {
			name: "g", Required: true, max: float32(len(globalGrades) - 1), step: 1, v8: isValidUint, Options: globalGradesDataList,
		}, {
			name: "r", max: 4, step: 1, v8: isValidUint, Options: dataListAgeGroup(),
		}, {
			name: "I", Required: true, maxLen: 64, v8: isValidRegex, regex: regexID,
		}}
	case 13:
		return []field{{
			name: "f", maxLen: 64, v8: isValidStr,
		}, {
			name: "s", maxLen: 64, v8: isValidStr,
		}, {
			name: "C", maxLen: 64, v8: isValidStr, manyRequired: []int{0, 1, 2}, manyRequiredQty: 1, Options: clubsDataList(),
		}}
	case 14:
		return []field{{
			name: "t", Required: true, max: 60, step: 1, v8: isValidUint,
		}, {
			name: "c", Required: true, max: 10, step: 1, v8: isValidUint,
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "R", Required: true, min: 1, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "S", Required: true, max: 65535, step: 1, v8: isValidUint,
		}}
	case 15:
		return []field{{
			name: "g", Required: true, minLen: 1, max: 65535, step: 1, v8: listUint, Options: availableGrades([]uint{}),
		}, {
			name: "I", v8: isValidRegex, regex: regexID,
		}}
	case 16:
		return []field{{
			name: "s", Required: true, maxLen: 12, minLen: 1, v8: isValidStr,
		}, {
			name: "E", Required: true, v8: isValidRegex, regex: regexID,
		}, {
			name: "R", Required: true, min: 1, max: 65535, step: 1, v8: isValidUint,
		}, {
			name: "S", Required: true, max: 65535, step: 1, v8: isValidUint,
		}}
	}
	return []field{}
}
