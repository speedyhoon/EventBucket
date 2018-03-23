//DO NOT EDIT! This file is generated by BuildIt.ninja form generator
//Please edit /form/forms.yml & then rebuild EventBucket
package main

import (
	"time"

	"github.com/speedyhoon/v8"
	"github.com/speedyhoon/forms"
)

const (
	clubNew              uint8 = 0
	clubEdit             uint8 = 1
	clubMoundNew         uint8 = 2
	eventNew             uint8 = 5
	eventEdit            uint8 = 6
	eventRangeNew        uint8 = 7
	eventRangeEdit       uint8 = 8
	eventAggNew          uint8 = 9
	eventAggEdit         uint8 = 10
	eventShooterNew      uint8 = 11
	eventShooterExisting uint8 = 13
	eventTotalScores     uint8 = 14
	eventAvailableGrades uint8 = 15
	eventUpdateShotScore uint8 = 16
	shooterNew           uint8 = 17
	shooterSearch        uint8 = 20
	shootersImport       uint8 = 21
	settings             uint8 = 22
)

func init() {
	post("/0", 0, clubInsert)
	post("/1", 1, clubDetailsUpsert)
	post("/2", 2, clubMoundInsert)
	post("/3", 3, clubMoundUpsert)
	get("/4", 4, clubsMap)
	post("/5", 5, eventInsert)
	post("/6", 6, eventDetailsUpsert)
	post("/7", 7, eventRangeInsert)
	post("/8", 8, eventRangeUpdate)
	post("/9", 9, eventAggInsert)
	post("/10", 10, eventAggUpdate)
	post("/11", 11, eventShooterInsert)
	post("/12", 12, eventShooterUpdate)
	post("/13", 13, eventShooterExistingInsert)
	post("/15", 15, eventAvailableGradesUpsert)
	post("/17", 17, shooterInsert)
	post("/18", 18, shooterInsert)
	post("/19", 19, shooterUpdate)
	get(urlShooters, 20, shooters)
	post("/22", 22, settingsUpdate)
}

func getFields(id uint8) []forms.Field {
	switch id {
	case 0: //New Club
		return []forms.Field{
			{Name: "n", V8: v8.StrReq},
		}
	case 1: //Club Details
		return []forms.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "a", V8: v8.Str},
			{Name: "w", V8: v8.Str},
			{Name: "p", V8: v8.Str},
			{Name: "x", V8: v8.Float32, Min: -90, Max: 90, Step: 1e-06},
			{Name: "y", V8: v8.Float32, Min: -180, Max: 180, Step: 1e-06},
			{Name: "b", V8: v8.Bool},
			{Name: "u", V8: v8.Str},
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case 2: //New Shooting Mound
		return []forms.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case 3: //Edit Shooting Mound
		return []forms.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "I", V8: v8.Uint, Max: 65535},
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case 4: //Map Clubs
		return []forms.Field{
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case 5: //New Event
		clubName := defaultClubName()
		return []forms.Field{
			{Name: "C", V8: v8.Str, Value: clubName, Required: clubName == "", MinLen: 1, Options: clubsDataList()},
			{Name: "n", V8: v8.StrReq},
			{Name: "d", V8: v8.Str, Value: time.Now().Format("2006-01-02"), MaxLen: 10},
			{Name: "t", V8: v8.Str, Value: time.Now().Format("15:04"), MaxLen: 5},
		}
	case 6: //Event Details
		return []forms.Field{
			{Name: "C", V8: v8.StrReq, Options: clubsDataList()},
			{Name: "n", V8: v8.StrReq},
			{Name: "d", V8: v8.Str, MaxLen: 10, MinLen: 1},
			{Name: "t", V8: v8.Str, MaxLen: 5, MinLen: 5},
			{Name: "c", V8: v8.Bool},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case 7: //Add Range
		return []forms.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case 8: //Update Range
		return []forms.Field{
			{Name: "I", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "n", V8: v8.StrReq},
			{Name: "k", V8: v8.Bool},
			{Name: "o", V8: v8.UintReq, Max: 65535},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case 9: //Add Aggregate Range
		return []forms.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "R", V8: v8.UintList, Required: true, MaxLen: 5, MinLen: 2, Min: 1, Max: 65535},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case 10: //Update Agg
		return []forms.Field{
			{Name: "E", V8: v8.Regex, Regex: regexID},
			{Name: "I", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "n", V8: v8.StrReq},
			{Name: "R", V8: v8.UintList, Required: true, MinLen: 2, Min: 1, Max: 65535},
			{Name: "o", V8: v8.UintReq, Max: 65535},
		}
	case 11: //Shooter Entry
		clubName := defaultClubName()
		return []forms.Field{
			{Name: "f", V8: v8.StrReq},
			{Name: "s", V8: v8.StrReq},
			{Name: "C", V8: v8.Str, Placeholder: clubName, Options: clubsDataList()},
			{Name: "S", V8: v8.Str, Options: searchShootersOptions("", "", clubName)},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "x", V8: v8.Bool},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case 12: //Entries Edit Shooter Details
		return []forms.Field{
			{Name: "S", V8: v8.UintReq, Max: 65535},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
			{Name: "f", V8: v8.StrReq},
			{Name: "s", V8: v8.StrReq},
			{Name: "C", V8: v8.Regex, Regex: regexID},
			{Name: "g", V8: v8.UintReq, Max: len(globalGrades) - 1, Options: globalGradesDataList},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "x", V8: v8.Bool},
			{Name: "k", V8: v8.Bool},
		}
	case 13: //Existing Shooter Entry
		return []forms.Field{
			{Name: "S", V8: v8.RegexReq, Regex: regexID},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
		}
	case 14: //Enter Range Totals
		return []forms.Field{
			{Name: "t", V8: v8.UintReq, Max: 120},
			{Name: "c", V8: v8.Uint, Max: 20},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
			{Name: "R", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "S", V8: v8.UintReq, Max: 65535},
			{Name: "h", V8: v8.Uint, Max: 100},
		}
	case 15: //Grades Available
		return []forms.Field{
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: availableGrades([]uint{})},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case 16: //Update Shooter Shots (Scorecards)
		return []forms.Field{
			{Name: "s", V8: v8.StrReq, MaxLen: 12},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
			{Name: "R", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "S", V8: v8.UintReq, Max: 65535},
		}
	case 17: //New Shooter
		clubName := defaultClubName()
		return []forms.Field{
			{Name: "f", V8: v8.StrReq},
			{Name: "s", V8: v8.StrReq},
			{Name: "C", V8: v8.Str, Placeholder: clubName, Required: clubName == "", MinLen: 1, Options: clubsDataList()},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "x", V8: v8.Bool},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
		}
	case 18: //Shooter Details
		return []forms.Field{
			{Name: "f", V8: v8.StrReq},
			{Name: "s", V8: v8.StrReq},
			{Name: "C", V8: v8.StrReq},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "x", V8: v8.Bool},
			{Name: "I", V8: v8.RegexReq, Regex: regexID},
		}
	case 19: //Shooter Update
		return []forms.Field{
			{Name: "f", V8: v8.Str},
			{Name: "s", V8: v8.Str},
			{Name: "C", V8: v8.Str},
		}
	case 20: //Shooter Search
		clubName := defaultClubName()
		return []forms.Field{
			{Name: "f", V8: v8.Str},
			{Name: "s", V8: v8.Str},
			{Name: "C", V8: v8.Str, Placeholder: clubName, Required: clubName == "", Options: clubsDataList()},
		}
	case 21: //Import Shooters
		return []forms.Field{
			{Name: "f", V8: v8.FileReq},
		}
	case 22: //Settings
		return []forms.Field{
			{Name: "t", V8: v8.Bool},
		}
	}
	return []forms.Field{}
}
