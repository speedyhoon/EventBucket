package main

const (
	VERSION = "0.90"

	//Server Settings
	PRODUCTION = false //False = output dev warnings, E.g. Template errors
	TEST_MODE = true	//display links to add n shooters or fillout all scores for a given range
	//Known issue - turning off minify breaks the startshooting page. moving to the next sibling in a table row return the textnode of whitespace instead of the next <td> tag
	MINIFY     = true  //turn on minify html

	//HTML Templates:
	//location "folder path/%v(filename).extension"
	PATH_HTML_MINIFIED = "htm/%v.htm"
	PATH_HTML_SOURCE   = "html/%v.html"

	//Main template html files
	TEMPLATE_HOME  = "_template_home"
	TEMPLATE_ADMIN = "_template_admin"
	TEMPLATE_EMPTY = "_template_empty"

	//folder structure
	DIR_ROOT = "./root/"
	DIR_CSS  = "/c/"
	DIR_JPEG = "/e/"
	DIR_GIF  = "/g/"
	DIR_JS   = "/j/"
	DIR_ICON = "/o/"
	DIR_PNG  = "/p/"
	DIR_SVG  = "/v/"

	FAVICON = "a"

	URL_about           = "/about"
	URL_licence         = "/licence"
	URL_licence_summary = "/licence-summary"
	URL_archive			  = "/archive"

	URL_organisers           = "/organisers"
	URL_event                = "/event/"
	URL_events               = "/events/"
	URL_eventSettings        = "/eventSettings/"
	URL_clubInsert           = "/clubInsert"
	URL_champInsert          = "/champInsert"
	URL_eventInsert          = "/eventInsert"
	URL_eventInsert2          = "/eventInsert2"
	URL_eventRangeInsert     = "/rangeInsert"
	URL_eventAggInsert       = "/aggInsert"
	URL_shooterInsert        = "/shooterInsert"
	URL_shooterListInsert    = "/shooterListInsert"
	URL_totalScores          = "/totalScores/"
	URL_totalScoresAll       = "/totalScoresAll/"
	URL_startShooting        = "/startShooting/"
	URL_startShootingAll     = "/startShootingAll/"
	URL_updateSortScoreBoard = "/updateSortScoreBoard"
	URL_updateTotalScores    = "/updateTotalScores"
	URL_updateShotScores     = "/updateShotScores"
	URL_updateEventGrades    = "/updateEventGrades"
	URL_updateEventName      = "/updateEventName/"
	URL_updateRange          = "/updateRange"
	URL_updateIsPrizeMeet    = "/updateIsPrizeMeet"
	URL_scoreboard           = "/scoreboard/"
	URL_dateUpdate           = "/dateUpdate/"
	URL_club                 = "/club/"
	URL_clubs                = "/clubs/"
	URL_club_settings        = "/clubSettings/"
	URL_club_mound_update    = "/clubMoundUpdate/"
	URL_club_mound_insert    = "/clubMoundInsert/"
	URL_updateShooterList    = "/updateShooterList"
	URL_queryShooterList     = "/queryShooterList"
	URL_eventShotsNSighters  = "/eventShotsNSighters"
	URL_rangeReport          = "/rangeReport/"
)

type ClassSettings struct {
	Desc                string
	Sighters            int
	Shots               int
	ValidShots          string
	ValidScore          string
	vcountback          string
	xcountback          string
	ValidCenta          string
	ValidSighters       string
	Valid               string
	Valid2              string
	Buttons             string
	CountbackX          bool
	CountbackValueX     int
	ShowInitialShots    int
	ShowMaxInitialShots int
	Grades              []string
}

//truman Cell -- air purifier
//TODO: eventually replace these settings with ones that are set for each club and sometimes overridden by a clubs event settings
const (
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

	//Scoreboard
	SCOREBOARD_SHOW_WARNING_FOR_ZERO_SCORES    = true
	SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES = false

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
	TARGET_Desc                = "Target Rifle 0-5 with V and X centers and able to convert Fclass scores to Target Rifle."
	TARGET_Sighters            = 2
	TARGET_Shots               = 10
	TARGET_ValidShots          = "012345V6X"
	TARGET_ValidScore          = "012345555"
	TARGET_vcountback          = "012345666"
	TARGET_xcountback          = "012345667"
	TARGET_ValidCenta          = "000000111"
	TARGET_ValidSighters       = ")!@#$%v^x"
	TARGET_Valid               = ")!@#$%v^x012345V6X"
	TARGET_Valid2              = "012345V6X012345V6X"
	TARGET_Buttons             = "012345VX"
	TARGET_CountbackX          = false
	TARGET_CountbackValueX     = 7
	TARGET_ShowInitialShots    = 2
	TARGET_ShowMaxInitialShots = 2
	TARGET_Grades              = "A,B,C"

	MATCH_Desc                = "Match Rifle 0-5 with V and X centers and able to convert to Fclass scores to Match Rifle."
	MATCH_Sighters            = 3
	MATCH_Shots               = 20
	MATCH_ValidShots          = "012345V6X"
	MATCH_ValidScore          = "012345555"
	MATCH_vcountback          = "012345666"
	MATCH_xcountback          = "012345667"
	MATCH_ValidCenta          = "000000111"
	MATCH_ValidSighters       = ")!@#$%v^x"
	MATCH_Valid               = ")!@#$%v^x012345V6X"
	MATCH_Valid2              = "012345V6X012345V6X"
	MATCH_Buttons             = "012345VX"
	MATCH_CountbackX          = true
	MATCH_CountbackValueX     = 7
	MATCH_ShowInitialShots    = 2
	MATCH_ShowMaxInitialShots = 2
	MATCH_Grades              = "MA,MB"

	FCLASS_Desc                = "Flcass 0-6 with X centers and able to convert Target and Match Rifle to Fclass scores."
	FCLASS_Sighters            = 2
	FCLASS_Shots               = 15
	FCLASS_ValidShots          = "012345V6X"
	FCLASS_ValidScore          = "012345666"
	FCLASS_vcountback          = "012345667"
	FCLASS_ValidCenta          = "000000001"
	FCLASS_ValidSighters       = ")!@#$%v^x"
	FCLASS_Valid               = ")!@#$%v^x012345V6X"
	FCLASS_Valid2              = "012345V6X012345V6X"
	FCLASS_Buttons             = "0123456X"
	FCLASS_CountbackX          = false
	FCLASS_CountbackValueX     = 7
	FCLASS_ShowInitialShots    = 2
	FCLASS_ShowMaxInitialShots = 2
	FCLASS_Grades              = "FA,FB,FO,FTR"

	//per Event
	SHOOTOFF_Sighters      = 2
	SHOOTOFF_ShotsStart    = 5
	SHOOTOFF_nextShots     = 3
	SHOOTOFF_UseXcountback = 1 //1= true, 0=false
	SHOOTOFF_UseXs         = 1
	SHOOTOFF_UseCountback  = 1 //system settings

	ERROR_ENTER_SCORES_IN_AGG = "<p>This range is an aggregate. Can't enter scores!</p>"
	ERROR_NO_SHOOTERS         = "<p>No Shooters entered in this event.</p>"
	ERROR_NO_EVENTS           = "<p>No upcoming events listed.</p>"
)

type ClassSettings2 struct {
	Name                  string
	Display               string
	DisplayValue          int
	Buttons               string
	SightersQty, ShotsQty int
	ValidShots            map[string]Score
	ValidSighters         []string
	GradeQty              int
	Grades                []string
}

var (
	GRADE_TO_INT = map[string]int{
		"a": 0,
		"b": 0,
		"c": 0,
		"d": 1,
		"e": 1,
		"f": 1,
		"g": 1,
		"h": 2,
		"i": 2,
		"j": 1,
	}

	ShotsToValue = map[string]string{
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
	}

	//	GRADE_ORDER = string{"a","b","c","d","e","f","g","h","i"}
	DEFAULT_CLASS_SETTINGS = []ClassSettings2{
		ClassSettings2{
			Name:         "target",
			Display:      "Target",
			DisplayValue: 0,
			Buttons:      "012345VX",
			SightersQty:  2,
			ShotsQty:     10,
			ValidShots: map[string]Score{
				"-": Score{Total: 0, Centers: 0, CountBack1: "0" /*, CountBack2:"0"*/},
				"0": Score{Total: 0, Centers: 0, CountBack1: "0" /*, CountBack2:"0"*/},
				"1": Score{Total: 1, Centers: 0, CountBack1: "1" /*, CountBack2:"1"*/},
				"2": Score{Total: 2, Centers: 0, CountBack1: "2" /*, CountBack2:"2"*/},
				"3": Score{Total: 3, Centers: 0, CountBack1: "3" /*, CountBack2:"3"*/},
				"4": Score{Total: 4, Centers: 0, CountBack1: "4" /*, CountBack2:"4"*/},
				"5": Score{Total: 5, Centers: 0, CountBack1: "5" /*, CountBack2:"5"*/},
				"V": Score{Total: 5, Centers: 1, CountBack1: "6" /*, CountBack2:"6"*/},
				"6": Score{Total: 5, Centers: 1, CountBack1: "6" /*, CountBack2:"6"*/},
				"X": Score{Total: 5, Centers: 1, CountBack1: "6" /*, CountBack2:"7"*/},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      3,
			Grades:        []string{"a", "b", "c"},
		},
		ClassSettings2{
			Name:         "fclass",
			Display:      "F Class",
			DisplayValue: 1,
			Buttons:      "0123456X",
			SightersQty:  2,
			ShotsQty:     10,
			ValidShots: map[string]Score{
				"-": Score{Total: 0, Centers: 0, CountBack1: "0" /*, CountBack2:"0"*/},
				"0": Score{Total: 0, Centers: 0, CountBack1: "0" /*, CountBack2:"0"*/},
				"1": Score{Total: 1, Centers: 0, CountBack1: "1" /*, CountBack2:"1"*/},
				"2": Score{Total: 2, Centers: 0, CountBack1: "2" /*, CountBack2:"2"*/},
				"3": Score{Total: 3, Centers: 0, CountBack1: "3" /*, CountBack2:"3"*/},
				"4": Score{Total: 4, Centers: 0, CountBack1: "4" /*, CountBack2:"4"*/},
				"5": Score{Total: 5, Centers: 0, CountBack1: "5" /*, CountBack2:"5"*/},
				"V": Score{Total: 5, Centers: 0, CountBack1: "6" /*, CountBack2:"6"*/},
				"6": Score{Total: 6, Centers: 0, CountBack1: "6" /*, CountBack2:"6"*/},
				"X": Score{Total: 6, Centers: 1, CountBack1: "7" /*, CountBack2:"7"*/},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      4,
			Grades:        []string{"d", "e", "f", "g", "j"},
		},
		ClassSettings2{
			Name:         "match",
			Display:      "Match",
			DisplayValue: 2,
			Buttons:      "012345VX",
			SightersQty:  2,
			ShotsQty:     15,
			ValidShots: map[string]Score{
				"-": Score{Total: 0, Centers: 0, CountBack1: "0" /*, CountBack2:"0"*/},
				"0": Score{Total: 0, Centers: 0, CountBack1: "0" /*, CountBack2:"0"*/},
				"1": Score{Total: 1, Centers: 0, CountBack1: "1" /*, CountBack2:"1"*/},
				"2": Score{Total: 2, Centers: 0, CountBack1: "2" /*, CountBack2:"2"*/},
				"3": Score{Total: 3, Centers: 0, CountBack1: "3" /*, CountBack2:"3"*/},
				"4": Score{Total: 4, Centers: 0, CountBack1: "4" /*, CountBack2:"4"*/},
				"5": Score{Total: 5, Centers: 0, CountBack1: "5" /*, CountBack2:"5"*/},
				"V": Score{Total: 5, Centers: 1, CountBack1: "6" /*, CountBack2:"6"*/},
				"6": Score{Total: 5, Centers: 1, CountBack1: "6" /*, CountBack2:"6"*/},
				"X": Score{Total: 5, Centers: 1, CountBack1: "6" /*, CountBack2:"7"*/},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      2,
			Grades:        []string{"h", "i"},
		},
	}

	ALL_CLASS_SETTINGS = map[string]ClassSettings{
		"target": ClassSettings{
			Desc:                "Target Rifle 0-5 with V and X centers and able to convert Fclass scores to Target Rifle.",
			Sighters:            2,
			Shots:               10,
			ValidShots:          "012345V6X",
			ValidScore:          "012345555",
			vcountback:          "012345666",
			xcountback:          "012345667",
			ValidCenta:          "000000111",
			ValidSighters:       ")!@#$%v^x",
			Valid:               ")!@#$%v^x012345V6X",
			Valid2:              "012345V6X012345V6X",
			Buttons:             "012345VX",
			CountbackX:          false,
			CountbackValueX:     7,
			ShowInitialShots:    2,
			ShowMaxInitialShots: 2,
			Grades:              []string{"a", "b", "c"},
		},
	}
	SCOREBOARD_LEGEND = []string{
		//Also sets the order for the legend
		"First",
		"Second",
		"Third",
		"HighestPossibleScore", //4 css class=w4 etc.
		"ShootOff",             //1
		"Incomplete",           //3
		"NoScore",              //2
	}
	SCOREBOARD_LEGEND_CSS_CLASSES = map[string][2]string{
		"First":                [2]string{0: "ST", 1: "First"},
		"Second":               [2]string{0: "ND", 1: "Second"},
		"Third":                [2]string{0: "TH", 1: "Third"},
		"HighestPossibleScore": [2]string{0: "w4", 1: "Highest Possible Score"},
		"ShootOff":             [2]string{0: "w1", 1: "Shoot Off"},
		"Incomplete":           [2]string{0: "w3", 1: "Incomplete Score"},
		"NoScore":              [2]string{0: "w2", 1: "No Score"},
	}
	//TODO make these dynamic from club settings
	CLASSES = map[string]string{
		"a": "Target A",
		"b": "Target B",
		"c": "Target C",
		"d": "F Class A",
		"e": "F Class B",
		"f": "F Class Open",
		"g": "F/TR",
		"h": "Match A",
		"i": "Match B",
		"j": "303 Rifle",
	}
	//TODO make these dynamic from club settings
	GRADE = map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
		"d": "FA",
		"e": "FB",
		"f": "F Open",
		"g": "F/TR",
		"h": "Open",
		"i": "Reserve",
		"j": "Rifle",
	}
	//TODO make these dynamic from club settings
	CLASS = map[string]string{
		"a": "Target",
		"b": "Target",
		"c": "Target",
		"d": "F Class",
		"e": "F Class",
		"f": "F Class",
		"g": "F Class",
		"h": "Match",
		"i": "Match",
		"j": "303",
	}
	CLASS_LONG = map[string]string{
		"a": "Target Rifle",
		"b": "Target Rifle",
		"c": "Target Rifle",
		"d": "F Class",
		"e": "F Class",
		"f": "F Class",
		"g": "F Class",
		"h": "Match Rifle",
		"i": "Match Rifle",
		"j": "303",
	}
	//TODO make these dynamic from club settings
	//TODO delete and use grade to int in start shooting instead
	JSCLASS = map[string]string{
		"a": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["a"]].Name,
		"b": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["b"]].Name,
		"c": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["c"]].Name,
		"d": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["d"]].Name,
		"e": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["e"]].Name,
		"f": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["f"]].Name,
		"g": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["g"]].Name,
		"h": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["h"]].Name,
		"i": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["i"]].Name,
		"j": DEFAULT_CLASS_SETTINGS[GRADE_TO_INT["j"]].Name,
	}
	//TODO make these dynamic from club settings
	AGE_GROUPS = map[string]string{
		"N":    "None",
		"U21": "Junior (U21)",
		"V":   "Veteran",
		"SV":  "Super Vet",
		"U25": "Junior (U25)",
	}
	AGE_GROUPS2 = []Option{
		0: Option{
			Display:  "None",
			Value: "N",
			Selected: true,
		},
		1: Option{
			Display: "Junior (U21)",
			Value:   "U21",
		},
		2: Option{
			Display: "Junior (U25)",
			Value:   "U25",
		},
		3: Option{
			Display: "Veteran",
			Value:   "V",
		},
		4: Option{
			Display: "Super Veteran",
			Value:   "SV",
		},
	}
)

func class_translation(class string) string {
	return CLASS[class]
}
func class_long_translation(class string) string {
	return CLASS_LONG[class]
}
func js_class_translation(class string) string {
	return JSCLASS[class]
}
func grade_translation(grade string) string {
	return GRADE[grade]
}
