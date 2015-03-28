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

	//truman Cell -- air purifier
	//TODO: eventually replace these settings with ones that are set for each club and sometimes overridden by a clubs event settings
	/*
		nullShots                     = "-" //record shots
		showMaxNumShooters            = 20
		showInitialShots              = 3 //the number of shots to show when a shooter is initially selected
		showMaxInitialShots           = 4
		shotGroupingBorder            = 3 //controlls where to place the shot separator/border between each number of shots
		borderBetweenSightersAndShots = true
		sighterGroupingBorder         = 2
		indentFinish                  = false
		startShootingInputs           = 0          //changes input text boxes to just tds for mobile input.
		allowClubNameWrap             = true       //Club names have spaces between each word. false=Club names have &nbsp; between words
		startShootingDefaultSighter   = "Drop All" //can select between 'Keep All' and 'Drop All'
		startShootingMaxNumShooters   = 100        //can select between 'Keep All' and 'Drop All'

		//Start Shooting Page
		STARTSHOOTING_COL_ID        = -1
		STARTSHOOTING_COL_UIN       = -2
		STARTSHOOTING_COL_CLASS     = -3
		STARTSHOOTING_COL_GRADE     = 4
		STARTSHOOTING_COL_CLUB      = 5
		STARTSHOOTING_COL_SHORTNAME = -6
		STARTSHOOTING_COL_NAME      = 7
		STARTSHOOTING_COL_SCORES    = 8
		STARTSHOOTING_COL_TOTAL     = 9
		STARTSHOOTING_COL_RECEIVED  = 10
		//the columns to show and their order.

		SCOREBOARD_COL_ID               = 1
		SCOREBOARD_COL_SHOOTERENTRYID   = -3 //usefull to show entry id when a shooter is entered twice into the same event with different classes
		SCOREBOARD_COL_UIN              = -5
		SCOREBOARD_COL_POSITION         = 100
		SCOREBOARD_COL_GRADE            = 20
		SCOREBOARD_COL_NAME             = 30
		SCOREBOARD_COL_CLASS            = -40
		SCOREBOARD_COL_CLUB             = 70
		SCOREBOARD_COL_GENDER           = -70
		SCOREBOARD_COL_AGE              = 80
		SCOREBOARD_COL_SHORTNAME        = -90
		SCOREBOARD_COL_RANGESCORES      = 13
		SCOREBOARD_ALTERNATE_ROW_COLOUR = 0 //colour every nth row, 0 = off
		SCOREBOARD_DISPLAY_INDIVIDUALs  = 1
		SCOREBOARD_COMBINE_GRADES       = 0
		SCOREBOARD_SHOW_TITLE           = 0 //1 = show, 0,-1 = hide titles -- show title of for syme or saturday/sunday etc
		SCOREBOARD_SHOW_TEAMS_XS        = 0 //1 = show, 0,-1 = hide Xs -- Agg columns if showXs == 1 display <sub>5Xs</sub>
		SCOREBOARD_SHOWTEAMS_SHOOTERS   = 1 //1 = show, 0,-1 = hide Xs -- When set to 1 display Team shooters scores, When set to 0 only display teams totals.
		SCOREBOARD_SHOW_SHOOTOFF        = 0
		SCOREBOARD_SHOW_IN_PROGRESS     = 1 //when enabled total score blinks while shooter is in progress

		// TODO: if one of the name options for scoreboard is not set then display the short name.
		// TODO: Add functionality to set these for javascript. output javascript code from golang. generate js file so it is cached and doesn't need to be generated on every page load.

		TARGET_Desc                = "Target Rifle 0-5 with V and X centres and able to convert Fclass scores to Target Rifle."
		MATCH_Desc                = "Match Rifle 0-5 with V and X centres and able to convert to Fclass scores to Match Rifle."
		FCLASS_Desc                = "Flcass 0-6 with X centres and able to convert Target and Match Rifle to Fclass scores."

		//per Event
		SHOOTOFF_Sighters      = 2
		SHOOTOFF_ShotsStart    = 5
		SHOOTOFF_nextShots     = 3
		SHOOTOFF_UseXcountback = 1 //1= true, 0=false
		SHOOTOFF_UseXs         = 1
		SHOOTOFF_UseCountback  = 1 //system settings
	*/

	//Strings
	v8MaxEventID     = 100
	v8MinEventID     = 1
	v8MaxStringInput = 100
	v8MinStringInput = 1
	v8MinShots       = 90
	v8Minhots        = 1
	//Integers
	v8MaxIntegerID = 999
	v8MinIntegerID = 0
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
			SightersQty:  2,
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

/*func calculateHPS4Class(classID, numberOdShots int) Score {
	return Score{
		Total:   defaultClassSettings[classID].Maximum.Total * numberOdShots,
		Centres: defaultClassSettings[classID].Maximum.Centres * numberOdShots,
	}
}*/

func isScoreHighestPossibleScore(classID, numberOdShots, total, centres int) bool {
	return defaultClassSettings[classID].Maximum.Total*numberOdShots == total && defaultClassSettings[classID].Maximum.Centres*numberOdShots == centres
}

/*func calculateHighestPossibleScores(numberOdShots int) []Score {
	var classHPS []Score
	for _, class := range defaultClassSettings {
		classHPS = append(classHPS, Score{
			Total:   class.Maximum.Total * numberOdShots,
			Centres: class.Maximum.Centres * numberOdShots,
		})
	}
	return classHPS
}*/

// Grade is exported
type Grade struct {
	Name, LongName, ClassName string
	ClassID                   int
	Settings                  ClassSettings
}

func gradeList() []int {
	return []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
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

//	LEGEND_FIRST                  = 5
//	LEGEND_SECOND                 = 6
//	LEGEND_THIRD                  = 7
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
