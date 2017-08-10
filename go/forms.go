//DO NOT EDIT! This file is generated by BuildIt.ninja form generator
//Please edit /form/forms.yml & then rebuild EventBucket
package main

import "time"

const (
	clubNew              uint8 = 0
	clubEdit             uint8 = 1
	clubMoundNew         uint8 = 2
	eventNew             uint8 = 3
	eventEdit            uint8 = 4
	eventRangeNew        uint8 = 5
	eventRangeEdit       uint8 = 6
	eventAggNew          uint8 = 7
	eventAggEdit         uint8 = 8
	eventShooterNew      uint8 = 9
	eventShooterExisting uint8 = 10
	eventTotalScores     uint8 = 13
	eventAvailableGrades uint8 = 14
	eventUpdateShotScore uint8 = 15
	shooterNew           uint8 = 19
	shootersImport       uint8 = 21
)

//TODO when forms are changed to zero index, perform a check if action is not equal to "" instead of strconv != 0
func init() {
	post("/0", 0, clubInsert)
	post("/1", 1, clubDetailsUpsert)
	post("/2", 2, clubMoundInsert)
	post("/17", 17, clubMoundUpsert)
	get("/16", 16, clubsMap)
	post("/3", 3, eventInsert)
	post("/4", 4, eventDetailsUpsert)
	post("/5", 5, eventRangeInsert)
	post("/6", 6, eventRangeUpdate)
	post("/7", 7, eventAggInsert)
	post("/8", 8, eventAggUpdate)
	post("/9", 9, eventShooterInsert)
	post("/18", 18, eventShooterUpdate)
	post("/10", 10, eventShooterExistingInsert)
	post("/14", 14, eventAvailableGradesUpsert)
	get("/11", 11, eventSearchShooters)
	post("/12", 12, shooterUpdate)
	post("/19", 19, shooterInsert)
	get(urlShooters, 20, shooters)
	//post("/21", 21, importShooters)
	post("/22", 22, settingsUpdate)
}

func getForm(id uint8) form {
	return form{
		action: id,
		Fields: func(id uint8) []field {
			switch id {
			case 0: // New Club
				return []field{
					{name: "n", v8: v8StrReq},
				}
			case 1: //Club Details
				return []field{
					{name: "n", v8: v8StrReq},
					{name: "a", v8: v8Str},
					{name: "w", v8: v8Str},
					{name: "p", v8: v8Str},
					{name: "x", v8: v8Float32, min: -90, max: 90, step: .000001},
					{name: "y", v8: v8Float32, min: -180, max: 180, step: .000001},
					{name: "b", v8: v8Bool},
					{name: "u", v8: v8Str},
					{name: "C", v8: v8Regex, regex: regexID},
				}
			case 2: //New Shooting Mound
				return []field{
					{name: "n", v8: v8StrReq},
					{name: "C", v8: v8Regex, regex: regexID},
				}
			case 3: //New Event
				clubName := defaultClubName()
				return []field{
					{name: "C", v8: v8Str, Value: clubName, Required: clubName == "", minLen: 1, Options: clubsDataList()},
					{name: "n", v8: v8StrReq},
					{name: "d", v8: v8Str, Value: time.Now().Format("2006-01-02"), maxLen: 10},
					{name: "t", v8: v8Str, Value: time.Now().Format("15:04"), maxLen: 5},
				}
			case 4: //Event Details
				return []field{
					{name: "C", v8: v8StrReq, Options: clubsDataList()},
					{name: "n", v8: v8StrReq},
					{name: "d", v8: v8Str, maxLen: 10, minLen: 1},
					{name: "t", v8: v8Str, maxLen: 5, minLen: 5},
					{name: "c", v8: v8Bool},
					{name: "E", v8: v8Regex, regex: regexID},
				}
			case 5: //Add Range
				return []field{
					{name: "n", v8: v8StrReq},
					{name: "E", v8: v8Regex, regex: regexID},
				}
			case 6: //Update Range
				return []field{
					{name: "E", v8: v8Regex, regex: regexID},
					{name: "I", v8: v8UintReq, min: 1, max: 65535},
					{name: "n", v8: v8StrReq},
					{name: "k", v8: v8Bool},
					{name: "o", v8: v8UintReq, max: 65535},
				}
			case 7: //Add Aggregate Range
				return []field{
					{name: "n", v8: v8StrReq},
					{name: "R", v8: v8UintList, Required: true, maxLen: 5, minLen: 2, min: 1, max: 65535},
					{name: "E", v8: v8Regex, regex: regexID},
				}
			case 8: //Update Agg
				return []field{
					{name: "E", v8: v8Regex, regex: regexID},
					{name: "I", v8: v8UintReq, min: 1, max: 65535},
					{name: "n", v8: v8StrReq},
					{name: "R", v8: v8UintList, Required: true, minLen: 2, min: 1, max: 65535},
					{name: "o", v8: v8UintReq, max: 65535},
				}
			case 9: //Shooter Entry
				return []field{
					{name: "f", v8: v8StrReq},
					{name: "s", v8: v8StrReq},
					{name: "C", v8: v8StrReq, Options: clubsDataList()},
					{name: "S", v8: v8Str, Options: searchShootersOptions("", "", "")},
					{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
					{name: "x", v8: v8Bool},
					{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
					{name: "E", v8: v8Regex, regex: regexID},
				}
			case 10: //Existing Shooter Entry
				return []field{
					{name: "S", v8: v8RegexReq, regex: regexID},
					{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
					{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
					{name: "E", v8: v8RegexReq, regex: regexID},
				}
			case 11: //Shooter Details
				return []field{
					{name: "f", v8: v8StrReq},
					{name: "s", v8: v8StrReq},
					{name: "C", v8: v8StrReq},
					{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
					{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
					{name: "x", v8: v8Bool},
					{name: "I", v8: v8RegexReq, regex: regexID},
				}
			case 12: //Shooter Search
				return []field{
					{name: "f", v8: v8Str},
					{name: "s", v8: v8Str},
					{name: "C", v8: v8Str},
				}
			case 13: //Enter Range Totals
				return []field{
					{name: "t", v8: v8UintReq, max: 120},
					{name: "c", v8: v8Uint, max: 20},
					{name: "E", v8: v8RegexReq, regex: regexID},
					{name: "R", v8: v8UintReq, min: 1, max: 65535},
					{name: "S", v8: v8UintReq, max: 65535},
					{name: "h", v8: v8Uint, max: 100},
				}
			case 14: //Grades Available
				return []field{
					{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: availableGrades([]uint{})},
					{name: "E", v8: v8Regex, regex: regexID},
				}
			case 15: //Update Shooter Shots (Scorecards)
				return []field{
					{name: "s", v8: v8StrReq, maxLen: 12},
					{name: "E", v8: v8RegexReq, regex: regexID},
					{name: "R", v8: v8UintReq, min: 1, max: 65535},
					{name: "S", v8: v8UintReq, max: 65535},
				}
			case 16: //Map Clubs
				return []field{
					{name: "C", v8: v8Regex, regex: regexID},
				}
			case 17: //Edit Shooting Mound
				return []field{
					{name: "n", v8: v8StrReq},
					{name: "I", v8: v8Uint, max: 65535},
					{name: "C", v8: v8Regex, regex: regexID},
				}
			case 18: //Entries Edit Shooter Details
				ageGroups := dataListAgeGroup()
				return []field{
					{name: "S", v8: v8UintReq, max: 65535},
					{name: "E", v8: v8RegexReq, regex: regexID},
					{name: "f", v8: v8StrReq},
					{name: "s", v8: v8StrReq},
					{name: "C", v8: v8Regex, regex: regexID},
					{name: "g", v8: v8UintReq, max: len(globalGrades) - 1, Options: globalGradesDataList},
					{name: "r", v8: v8Uint, max: len(ageGroups) - 1, Options: ageGroups},
					{name: "x", v8: v8Bool},
					{name: "k", v8: v8Bool},
				}
			case 19: //New Shooter
				clubName := defaultClubName()
				return []field{
					{name: "f", v8: v8StrReq},
					{name: "s", v8: v8StrReq},
					{name: "C", v8: v8Str, Value: clubName, Placeholder: clubName, Required: clubName == "", minLen: 1, Options: clubsDataList()},
					{name: "r", v8: v8Uint, max: 4, Options: dataListAgeGroup()},
					{name: "x", v8: v8Bool},
					{name: "g", v8: v8UintList, Required: true, max: len(globalGrades) - 1, Options: globalGradesDataList},
				}
			case 20: //Shooter Search
				return []field{
					{name: "f", v8: v8Str},
					{name: "s", v8: v8Str},
					{name: "C", v8: v8Str},
				}
			case 21: //Import Shooters
				return []field{
					{name: "f", v8: v8File, Required: true},
				}
			case 22: //Settings
				return []field{
					{name: "t", v8: v8StrReq, Options: []option{
						{Label: "Lite", Selected: !masterTemplate.IsDarkTheme},
						{Label: "Dark", Selected: masterTemplate.IsDarkTheme},
					}},
				}
			}
			return []field{}
		}(id),
	}
}