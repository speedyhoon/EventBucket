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

	clubNew              = 0
	clubDetails          = 1
	clubMoundNew         = 2
	eventNew             = 3
	eventDetails         = 4
	eventRangeNew        = 5
	eventAggNew          = 6
	eventShooterExisting = 7
	eventShooterNew      = 8
)

func dataListGrades() []option {
	return []option{
		{},
		{Value: "1", Label: "Target A"},
		{Value: "2", Label: "Target B"},
		{Value: "3", Label: "Target C"},
		{Value: "4", Label: "F Class A"},
		{Value: "5", Label: "F Class B"},
		{Value: "6", Label: "F Class Open"},
		{Value: "7", Label: "F/TR"},
		{Value: "8", Label: "Match Open"},
		{Value: "9", Label: "Match Reserve"},
		{Value: "10", Label: "303 Rifle"},
	}
}

func dataListAgeGroup() []option {
	return []option{
		{},
		{Value: "1", Label: "Junior U21"},
		{Value: "2", Label: "Junior U25"},
		{Value: "3", Label: "Veteran"},
		{Value: "4", Label: "Super Veteran"},
	}
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
			name: "C", v8: isValidID,
		}}
	case 2:
		return []field{{
			name: "n", maxLen: 64, v8: isValidStr,
		}, {
			name: "e", Required: true, min: 10, max: 65535, step: .01, v8: isValidFloat32,
		}, {
			name: "z", Required: true, v8: isValidStr, Options: []option{{Label: "Metres", Value: "m", Selected: false}, {Label: "Kilometres", Value: "k", Selected: false}, {Label: "Yards", Value: "y", Selected: false}, {Label: "Feet", Value: "f", Selected: false}},
		}, {
			name: "C", v8: isValidID,
		}}
	case 3:
		return []field{{
			name: "C", Value: defaultClubName(), Required: hasDefaultClub(), maxLen: 64, v8: isValidStr, Options: getDataListClubs(),
		}, {
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "d", Value: defaultDate(), v8: isValidStr,
		}, {
			name: "t", Value: defaultTime(), step: 300, v8: isValidStr,
		}}
	case 4:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "C", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "d", v8: isValidStr,
		}, {
			name: "t", v8: isValidStr,
		}, {
			name: "E", v8: isValidID,
		}}
	case 5:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "E", v8: isValidID,
		}}
	case 6:
		return []field{{
			name: "n", Required: true, maxLen: 64, v8: isValidStr,
		}, {
			name: "R", Required: true, v8: isValidStr,
		}, {
			name: "E", v8: isValidID,
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
			name: "g", Required: true, v8: isValidUint64, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Target A", Value: "1", Selected: false}, {Label: "Target B", Value: "2", Selected: false}, {Label: "Target C", Value: "3", Selected: false}, {Label: "F Class A", Value: "4", Selected: false}, {Label: "F Class B", Value: "5", Selected: false}, {Label: "F Class Open", Value: "6", Selected: false}, {Label: "F/TR", Value: "7", Selected: false}, {Label: "Match Open", Value: "8", Selected: false}, {Label: "Match Reserve", Value: "9", Selected: false}, {Label: "303 Rifle", Value: "10", Selected: false}},
		}, {
			name: "r", Required: true, v8: isValidUint64, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Junior U21", Value: "1", Selected: false}, {Label: "Junior U25", Value: "2", Selected: false}, {Label: "Veteran", Value: "3", Selected: false}, {Label: "Super Veteran", Value: "4", Selected: false}},
		}, {
			name: "E",
		}, {
			name: "E",
		}}
	case 8:
		return []field{{
			name: "f", Required: true, maxLen: 64,
		}, {
			name: "s", Required: true, maxLen: 64,
		}, {
			name: "C", Required: true, maxLen: 64,
		}, {
			name: "S",
		}, {
			name: "g", Required: true, v8: isValidUint64, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Target A", Value: "1", Selected: false}, {Label: "Target B", Value: "2", Selected: false}, {Label: "Target C", Value: "3", Selected: false}, {Label: "F Class A", Value: "4", Selected: false}, {Label: "F Class B", Value: "5", Selected: false}, {Label: "F Class Open", Value: "6", Selected: false}, {Label: "F/TR", Value: "7", Selected: false}, {Label: "Match Open", Value: "8", Selected: false}, {Label: "Match Reserve", Value: "9", Selected: false}, {Label: "303 Rifle", Value: "10", Selected: false}},
		}, {
			name: "r", Required: true, v8: isValidUint64, Options: []option{{Label: "", Value: "", Selected: false}, {Label: "Junior U21", Value: "1", Selected: false}, {Label: "Junior U25", Value: "2", Selected: false}, {Label: "Veteran", Value: "3", Selected: false}, {Label: "Super Veteran", Value: "4", Selected: false}},
		}, {
			name: "E",
		}, {
			name: "E",
		}}
	}
	return []field{}
}
