package main

import (
	"log"
	"os"
)

const (
	//GET
	urlHome     = "/"
	urlAbout    = "/about"
	urlClubs    = "/clubs"
	urlLicence  = "/licence"
	urlArchive  = "/archive"
	urlShooters = "/shooters"
	//GET with PARAMETERS
	urlEvent            = "/event/" //eventID
	urlClub             = "/club/"
	urlEventSettings    = "/eventSettings/"
	urlScoreboard       = "/scoreboard/" //eventID/rangeID
	urlTotalScores      = "/totalScores/"
	urlTotalScoresAll   = "/totalScoresAll/"
	urlStartShooting    = "/startShooting/"
	urlStartShootingAll = "/startShootingAll/"
	//POST
	urlQueryShooterList     = "/queryShooterList"
	urlQueryShooterGrade    = "/queryShooterGrade"
	urlEventUpdateShooter   = "/eventUpdateShooter"
	urlClubInsert           = "/clubInsert"
	urlEventInsert          = "/eventInsert"
	urlEventRangeInsert     = "/rangeInsert"
	urlEventAggInsert       = "/aggInsert"
	urlShooterInsert        = "/shooterInsert"
	urlShooterListInsert    = "/shooterListInsert"
	urlUpdateSortScoreBoard = "/updateSortScoreBoard"
	urlUpdateTotalScores    = "/updateTotalScores"
	urlUpdateShotScores     = "/updateShotScores"
	urlUpdateEventGrades    = "/updateEventGrades"
	urlUpdateRange          = "/updateRange"
	urlUpdateIsPrizeMeet    = "/updateIsPrizeMeet"
	urlClubMoundInsert      = "/clubMoundInsert/"
	urlClubDetailsUpsert    = "/clubDetailsUpsert/"
	urlUpdateShooterList    = "/updateShooterList"
	urlEventShotsNSighters  = "/eventShotsNSighters/"
	//urlClubMoundUpdate      = "/clubMoundUpdate/"
	//urlChampInsert          = "/champInsert"
	//urlRangeReport          = "/rangeReport/"

	//Global program settings
	versionNumber    = "^^VersionNumber^^"
	buildDate        = "^^BuildDate^^"
	pathHTMLMinified = "h/%v.htm"
	//Main template html files
	templateHome   = "_template_home"
	templateAdmin  = "_template_admin"
	templateEmpty  = "_template_empty"
	idCharset      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!*()_-."
	idCharsetRegex = `\w~!\*\(\)\-\.`
	//Folder structure
	dirCSS  = "^^dirCSS^^"
	dirJPEG = "^^dirJPEG^^"
	dirJS   = "^^dirJS^^"
	dirPNG  = "^^dirPNG^^"
	dirSVG  = "^^dirSVG^^"
	dirWOF  = "^^dirWOF^^"
	dirWOF2 = "^^dirWOF2^^"
	//Barcodes
	//QRCODE     = "qr"
	//DATAMATRIX = "dm"

	//Scoreboard Settings
	//SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES = true //true = Don't award shooters a place if they haven't submitted a score, false = shooter without a score is awarded last place (gets 5th when beaten by 4 other shooters)

	//Total Scores Settings
	errorEnterScoresInAgg = "<p>This range is an aggregate. Can't enter scores!</p>"

	v8MaxEventID     = 100
	v8MinEventID     = 1
	v8MaxStringInput = 100
	v8MinStringInput = 1
	v8MinShots       = 90
	v8Minhots        = 1
	v8MaxIntegerID   = 999
	v8MinIntegerID   = 0
)

// ClassSettings is exported
type ClassSettings struct {
	Name                  string
	Display               string
	DisplayValue          int
	Buttons               string
	SightersQty, ShotsQty int
	ValidShots            map[string]Score
	ValidSighters         []string
	GradeQty              int
	Grades                []int
	Maximum               Score
}

var (
	// Error is exported
	Error = log.New(os.Stderr, "ERROR:  ", log.Ldate|log.Ltime|log.Lshortfile)
	//TODO move the below to a constant if possible
	latitudeMin  = -90
	latitudeMax  = 90
	longitudeMin = -180
	longitudeMax = 180

	//BARCODE_TYPE = QRCODE

	homeMenuItems = []Menu{
		{
			Name: "Home",
			Link: urlHome,
		}, {
			Name: "Archive",
			Link: urlArchive,
		}, {
			Name: "Clubs",
			Link: urlClubs,
		}, {
			Name: "About",
			Link: urlAbout,
		}, {
			Name: "Shooters",
			Link: urlShooters,
		}, {
			Name:   "Licence",
			Link:   urlLicence,
			Hidden: true,
		},
	}

	defaultClassSettings = []ClassSettings{
		{
			Name:         "target",
			Display:      "Target",
			DisplayValue: 0,
			Buttons:      "012345VX",
			SightersQty:  2,
			ShotsQty:     10,
			Maximum:      Score{Total: 5, Centres: 1},
			ValidShots: map[string]Score{
				"-": {Total: 0, Centres: 0, CountBack: "0"},
				"0": {Total: 0, Centres: 0, CountBack: "0"},
				"1": {Total: 1, Centres: 0, CountBack: "1"},
				"2": {Total: 2, Centres: 0, CountBack: "2"},
				"3": {Total: 3, Centres: 0, CountBack: "3"},
				"4": {Total: 4, Centres: 0, CountBack: "4"},
				"5": {Total: 5, Centres: 0, CountBack: "5"},
				"V": {Total: 5, Centres: 1, CountBack: "6"},
				"6": {Total: 5, Centres: 1, CountBack: "6"},
				"X": {Total: 5, Centres: 1, CountBack: "6"},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      3,
			Grades:        []int{0, 1, 2},
		},
		{
			Name:         "fclass",
			Display:      "F Class",
			DisplayValue: 1,
			Buttons:      "0123456X",
			SightersQty:  2,
			ShotsQty:     10,
			Maximum:      Score{Total: 6, Centres: 1},
			ValidShots: map[string]Score{
				"-": {Total: 0, Centres: 0, CountBack: "0"},
				"0": {Total: 0, Centres: 0, CountBack: "0"},
				"1": {Total: 1, Centres: 0, CountBack: "1"},
				"2": {Total: 2, Centres: 0, CountBack: "2"},
				"3": {Total: 3, Centres: 0, CountBack: "3"},
				"4": {Total: 4, Centres: 0, CountBack: "4"},
				"5": {Total: 5, Centres: 0, CountBack: "5"},
				"V": {Total: 5, Centres: 0, CountBack: "6"},
				"6": {Total: 6, Centres: 0, CountBack: "6"},
				"X": {Total: 6, Centres: 1, CountBack: "7"},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      4,
			Grades:        []int{3, 4, 5, 6, 9},
		},
		{
			Name:         "match",
			Display:      "Match",
			DisplayValue: 2,
			Buttons:      "012345VX",
			SightersQty:  3,
			ShotsQty:     15,
			Maximum:      Score{Total: 5, Centres: 1},
			ValidShots: map[string]Score{
				"-": {Total: 0, Centres: 0, CountBack: "0"},
				"0": {Total: 0, Centres: 0, CountBack: "0"},
				"1": {Total: 1, Centres: 0, CountBack: "1"},
				"2": {Total: 2, Centres: 0, CountBack: "2"},
				"3": {Total: 3, Centres: 0, CountBack: "3"},
				"4": {Total: 4, Centres: 0, CountBack: "4"},
				"5": {Total: 5, Centres: 0, CountBack: "5"},
				"V": {Total: 5, Centres: 1, CountBack: "6"},
				"6": {Total: 5, Centres: 1, CountBack: "6"},
				"X": {Total: 5, Centres: 1, CountBack: "6"},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      2,
			Grades:        []int{7, 8},
		},
	}

	// ClassNamesList is exported
	ClassNamesList = []string{
		0: "Target A",
		1: "Target B",
		2: "Target C",
		3: "F Class A",
		4: "F Class B",
		5: "F Class Open",
		6: "F/TR",
		7: "Match Open",
		8: "Match Reserve",
		9: "303 Rifle",
	}
)

func isScoreHighestPossibleScore(classID, numberOdShots, total, centres int) bool {
	return defaultClassSettings[classID].Maximum.Total*numberOdShots == total && defaultClassSettings[classID].Maximum.Centres*numberOdShots == centres
}

// Grade is exported
type Grade struct {
	Name, LongName, ClassName string
	ClassID                   int
	Settings                  ClassSettings
}

func gradeList() []int {
	return []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
}
func grade2Class(grade int) int {
	return []Grade{
		0: {ClassID: 0},
		1: {ClassID: 0},
		2: {ClassID: 0},
		3: {ClassID: 1},
		4: {ClassID: 1},
		5: {ClassID: 1},
		6: {ClassID: 1},
		7: {ClassID: 2},
		8: {ClassID: 2},
		9: {ClassID: 1},
	}[grade].ClassID
}
func grades() []Grade {
	return []Grade{
		0: {Settings: defaultClassSettings[0], ClassID: 0, Name: "A", ClassName: "Target", LongName: "Target A"},
		1: {Settings: defaultClassSettings[0], ClassID: 0, Name: "B", ClassName: "Target", LongName: "Target B"},
		2: {Settings: defaultClassSettings[0], ClassID: 0, Name: "C", ClassName: "Target", LongName: "Target C"},
		3: {Settings: defaultClassSettings[1], ClassID: 1, Name: "FA", ClassName: "F Class", LongName: "F Class A"},
		4: {Settings: defaultClassSettings[1], ClassID: 1, Name: "FB", ClassName: "F Class", LongName: "F Class B"},
		5: {Settings: defaultClassSettings[1], ClassID: 1, Name: "F Open", ClassName: "F Class", LongName: "F Class Open"},
		6: {Settings: defaultClassSettings[1], ClassID: 1, Name: "F/TR", ClassName: "F Class", LongName: "F/TR"},
		7: {Settings: defaultClassSettings[2], ClassID: 2, Name: "Open", ClassName: "Match", LongName: "Match Open"},
		8: {Settings: defaultClassSettings[2], ClassID: 2, Name: "Reserve", ClassName: "Match", LongName: "Match Reserve"},
		9: {Settings: defaultClassSettings[1], ClassID: 1, Name: "Rifle", ClassName: "303", LongName: "303 Rifle"},
	}
}

// AgeGroupDisplay is exported
func AgeGroupDisplay(value string) string {
	if value != "N" {
		for _, ageGroup := range ageGroups() {
			if value == ageGroup.Value {
				return ageGroup.Display
			}
		}
	}
	return ""
}

func ageGroups() []Option {
	return []Option{
		0: {
			Display:  "None",
			Value:    "N",
			Selected: true,
		},
		1: {
			Display: "Junior (U21)",
			Value:   "U21",
		},
		2: {
			Display: "Junior (U25)",
			Value:   "U25",
		},
		3: {
			Display: "Veteran",
			Value:   "V",
		},
		4: {
			Display: "Super Veteran",
			Value:   "SV",
		},
	}
}

//Return age group select box options with the shooters value selected
func shooterAgeGroupSelectbox(shooter EventShooter) []Option {
	options := ageGroups()
	for _, ageGroup := range options {
		ageGroup.Value = shooter.AgeGroup == ageGroup.Value
	}
	return options
}

// Legend is exported
type Legend struct {
	//To access a field in HTML a struct, it must start with an uppercase letter. Other wise it will output error: xxx is an unexported field of struct type main.Legend
	CSSClass, Name string
}

const (
	legendShootOff             = 1
	legendNoScore              = 2
	legendIncompleteScore      = 3
	legendHighestPossibleScore = 4
	//LEGEND_FIRST               = 5
	//LEGEND_SECOND              = 6
	//LEGEND_THIRD               = 7
)

func scoreBoardLegend() []Legend {
	return []Legend{
		{CSSClass: "w4", Name: "Highest Possible Score"},
		{CSSClass: "w1", Name: "Shoot Off"},
		{CSSClass: "w3", Name: "Incomplete Score"},
		{CSSClass: "w2", Name: "No Score"},
		{CSSClass: "p1", Name: "1st"},
		{CSSClass: "p2", Name: "2nd"},
		{CSSClass: "p3", Name: "3rd"},
	}
}

func shotsToValue(shot string) string {
	return map[string]string{
		"-": "",
		"0": "0",
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
		"V": "V",
		"6": "6",
		"X": "X",
		")": "0",
		"!": "1",
		"@": "2",
		"#": "3",
		"$": "4",
		"%": "5",
		"v": "V",
		"^": "6",
		"x": "X",
	}[shot]
}
