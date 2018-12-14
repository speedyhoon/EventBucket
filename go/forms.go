//DO NOT EDIT! This file is generated by BuildIt.ninja form generator
//Please edit /form/frm.yml & then rebuild EventBucket
package main

import (
	"time"

	"github.com/speedyhoon/frm"
	"github.com/speedyhoon/v8"
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
	post("/0", clubNew, clubInsert)
	post("/1", clubEdit, clubDetailsUpsert)
	post("/2", clubMoundNew, clubMoundInsert)
	post("/3", 3, clubMoundUpsert)
	get("/4", 4, clubsMap)
	post("/5", eventNew, eventInsert)
	post("/6", eventEdit, eventDetailsUpsert)
	post("/7", eventRangeNew, eventRangeInsert)
	post("/8", eventRangeEdit, eventRangeUpdate)
	post("/9", eventAggNew, eventAggInsert)
	post("/10", eventAggEdit, eventAggUpdate)
	post("/11", eventShooterNew, eventShooterInsert)
	post("/12", 12, eventShooterUpdate)
	post("/13", eventShooterExisting, eventShooterExistingInsert)
	post("/15", eventAvailableGrades, eventAvailableGradesUpsert)
	post("/17", shooterNew, shooterInsert)
	post("/18", 18, shooterInsert)
	post("/19", 19, shooterUpdate)
	get(urlShooters, shooterSearch, shooters)
	post("/22", settings, settingsUpdate)
}

func getFields(id uint8) []frm.Field {
	const dateTime = "2006-01-02 15:04"
	switch id {
	case clubNew: //New Club
		return []frm.Field{
			{Name: "n", V8: v8.StrReq},
		}
	case clubEdit: //Club Details
		return []frm.Field{
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
	case clubMoundNew: //New Shooting Mound
		return []frm.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case 3: //Edit Shooting Mound
		return []frm.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "I", V8: v8.Uint, Max: 65535},
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case 4: //Map Clubs
		return []frm.Field{
			{Name: "C", V8: v8.Regex, Regex: regexID},
		}
	case eventNew: //New Event
		club := defaultClub()
		return []frm.Field{
			{Name: "C", V8: v8.Str, Value: club.Name, Required: club.IsDefault, MinLen: 1, Options: clubsDataList()},
			{Name: "n", V8: v8.StrReq},
			{Name: "d", V8: v8.DateTime, Value: time.Now().Format(dateTime), Placeholder: dateTime, MaxLen: 16, MinLen: 14},
		}
	case eventEdit: //Event Details
		return []frm.Field{
			{Name: "C", V8: v8.StrReq, Options: clubsDataList()},
			{Name: "n", V8: v8.StrReq},
			{Name: "d", V8: v8.DateTime, Placeholder: dateTime, MaxLen: 16, MinLen: 14},
			{Name: "c", V8: v8.Bool},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case eventRangeNew: //Add Range
		return []frm.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case eventRangeEdit: //Update Range
		return []frm.Field{
			{Name: "I", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "n", V8: v8.StrReq},
			{Name: "k", V8: v8.Bool},
			{Name: "o", V8: v8.UintReq, Max: 65535},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case eventAggNew: //Add Aggregate Range
		return []frm.Field{
			{Name: "n", V8: v8.StrReq},
			{Name: "R", V8: v8.UintList, Required: true, MaxLen: 5, MinLen: 2, Min: 1, Max: 65535},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case eventAggEdit: //Update Agg
		return []frm.Field{
			{Name: "E", V8: v8.Regex, Regex: regexID},
			{Name: "I", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "n", V8: v8.StrReq},
			{Name: "R", V8: v8.UintList, Required: true, MinLen: 2, Min: 1, Max: 65535},
			{Name: "o", V8: v8.UintReq, Max: 65535},
		}
	case eventShooterNew: //Shooter Entry
		clubName := defaultClub().Name
		return []frm.Field{
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
		return []frm.Field{
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
	case eventShooterExisting: //Existing Shooter Entry
		return []frm.Field{
			{Name: "S", V8: v8.RegexReq, Regex: regexID},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
		}
	case eventTotalScores: //Enter Range Totals
		return []frm.Field{
			{Name: "t", V8: v8.UintReq, Max: 120},
			{Name: "c", V8: v8.Uint, Max: 20},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
			{Name: "R", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "S", V8: v8.UintReq, Max: 65535},
			{Name: "h", V8: v8.Uint, Max: 100},
		}
	case eventAvailableGrades: //Grades Available
		return []frm.Field{
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: availableGrades([]uint{})},
			{Name: "E", V8: v8.Regex, Regex: regexID},
		}
	case eventUpdateShotScore: //Update Shooter Shots (Scorecards)
		return []frm.Field{
			{Name: "s", V8: v8.StrReq, MaxLen: 12},
			{Name: "E", V8: v8.RegexReq, Regex: regexID},
			{Name: "R", V8: v8.UintReq, Min: 1, Max: 65535},
			{Name: "S", V8: v8.UintReq, Max: 65535},
		}
	case shooterNew: //New Shooter
		club := defaultClub()
		return []frm.Field{
			{Name: "f", V8: v8.StrReq},
			{Name: "s", V8: v8.StrReq},
			{Name: "C", V8: v8.Str, Placeholder: club.Name, Required: club.IsDefault, MinLen: 1, Options: clubsDataList()},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "x", V8: v8.Bool},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
		}
	case 18: //Shooter Details
		return []frm.Field{
			{Name: "f", V8: v8.StrReq},
			{Name: "s", V8: v8.StrReq},
			{Name: "C", V8: v8.StrReq},
			{Name: "g", V8: v8.UintList, Required: true, Max: len(globalGrades) - 1, Options: globalGradesDataList},
			{Name: "r", V8: v8.UintOpt, Options: dataListAgeGroup()},
			{Name: "x", V8: v8.Bool},
			{Name: "I", V8: v8.RegexReq, Regex: regexID},
		}
	case 19: //Shooter Update
		return []frm.Field{
			{Name: "f", V8: v8.Str},
			{Name: "s", V8: v8.Str},
			{Name: "C", V8: v8.Str},
		}
	case shooterSearch: //Shooter Search
		club := defaultClub()
		return []frm.Field{
			{Name: "f", V8: v8.Str},
			{Name: "s", V8: v8.Str},
			{Name: "C", V8: v8.Str, Placeholder: club.Name, Required: club.IsDefault, Options: clubsDataList()},
		}
	case shootersImport: //Import Shooters
		return []frm.Field{
			{Name: "f", V8: v8.FileReq},
		}
	case settings: //Settings
		return []frm.Field{
			{Name: "t", V8: v8.Bool},
		}
	}
	return []frm.Field{}
}
