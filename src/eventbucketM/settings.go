package main

import "regexp"
const (
//Global program settings
	PRODUCTION = false //False = output dev warnings, E.g. Template errors
	//TEST_MODE = false //display links to add n shooters or fillout all scores for a given range
	//Known issue - turning off minify breaks the startshooting page. moving to the next sibling in a table row return the textnode of whitespace instead of the next <td> tag
	//	MINIFY     = true  //turn on minify html
	//HTML Templates:
	//location "folder path/%v(filename).extension"
	PATH_HTML_MINIFIED = "htm/%v.htm"
	PATH_HTML_SOURCE   = "html/%v.html"
	//Main template html files
	TEMPLATE_HOME  = "_template_home"
	TEMPLATE_ADMIN = "_template_admin"
	TEMPLATE_EMPTY = "_template_empty"
	VERSION = 0.58
	BUILDDATE = "Compiled on February 17, 2015  by Cam Webb"
	ID_CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!*()_-."
	ID_CHARSET_REGEX = `\w~!\*\(\)\-\.`
//Folder structure
	DIR_ROOT = "./root/"
	DIR_CSS  = "/c/"
	DIR_JPEG = "/e/"
	//	DIR_GIF  = "/g/"
	DIR_JS   = "/j/"
	DIR_PNG  = "/p/"
	DIR_SVG  = "/v/"
	//	DIR_WEBP = "/w/"

	//Scoreboard
	SCOREBOARD_SHOW_WARNING_FOR_ZERO_SCORES    = true
	SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES = false


	//truman Cell -- air purifier
	//TODO: eventually replace these settings with ones that are set for each club and sometimes overridden by a clubs event settings
/*	nullShots                     = "-" //record shots
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

	TARGET_Desc                = "Target Rifle 0-5 with V and X centers and able to convert Fclass scores to Target Rifle."
	MATCH_Desc                = "Match Rifle 0-5 with V and X centers and able to convert to Fclass scores to Match Rifle."
	FCLASS_Desc                = "Flcass 0-6 with X centers and able to convert Target and Match Rifle to Fclass scores."

	//per Event
	SHOOTOFF_Sighters      = 2
	SHOOTOFF_ShotsStart    = 5
	SHOOTOFF_nextShots     = 3
	SHOOTOFF_UseXcountback = 1 //1= true, 0=false
	SHOOTOFF_UseXs         = 1
	SHOOTOFF_UseCountback  = 1 //system settings
*/
)

type Page struct{
	TemplateFile, Title string
	Theme string	//TODO change type to struct enum to select between "TEMPLATE_HOME", "TEMPLATE_ADMIN" & "TEMPLATE_EMPTY"
	Data M
	v8Url *regexp.Regexp
}

var (
	VURL_home                = regexp.MustCompile("^/$")
	VURL_event               = regexp.MustCompile("^"+URL_event+"(["+ID_CHARSET_REGEX+"]+)$")
	VURL_eventShotsNSighters = regexp.MustCompile("^"+URL_eventShotsNSighters+"(["+ID_CHARSET_REGEX+"]+)$")
	VURL_club                = regexp.MustCompile("^"+URL_club+"(["+ID_CHARSET_REGEX+"]+)$")
	VURL_eventSettings       = regexp.MustCompile("^"+URL_eventSettings+"(["+ID_CHARSET_REGEX+"]+)$")
	VURL_scoreboard          = regexp.MustCompile("^"+URL_scoreboard+"(["+ID_CHARSET_REGEX+"]+)$")
)

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
}

func ShotsToValue(shot string) string {
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

var (
	DEFAULT_CLASS_SETTINGS = []ClassSettings{
		ClassSettings{
			Name:         "target",
			Display:      "Target",
			DisplayValue: 0,
			Buttons:      "012345VX",
			SightersQty:  2,
			ShotsQty:     10,
			ValidShots: map[string]Score{
				"-": Score{Total: 0, Centers: 0, CountBack1: "0"},
				"0": Score{Total: 0, Centers: 0, CountBack1: "0"},
				"1": Score{Total: 1, Centers: 0, CountBack1: "1"},
				"2": Score{Total: 2, Centers: 0, CountBack1: "2"},
				"3": Score{Total: 3, Centers: 0, CountBack1: "3"},
				"4": Score{Total: 4, Centers: 0, CountBack1: "4"},
				"5": Score{Total: 5, Centers: 0, CountBack1: "5"},
				"V": Score{Total: 5, Centers: 1, CountBack1: "6"},
				"6": Score{Total: 5, Centers: 1, CountBack1: "6"},
				"X": Score{Total: 5, Centers: 1, CountBack1: "6"},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      3,
			Grades:        []int{0, 1, 2},
		},
		ClassSettings{
			Name:         "fclass",
			Display:      "F Class",
			DisplayValue: 1,
			Buttons:      "0123456X",
			SightersQty:  2,
			ShotsQty:     10,
			ValidShots: map[string]Score{
				"-": Score{Total: 0, Centers: 0, CountBack1: "0"},
				"0": Score{Total: 0, Centers: 0, CountBack1: "0"},
				"1": Score{Total: 1, Centers: 0, CountBack1: "1"},
				"2": Score{Total: 2, Centers: 0, CountBack1: "2"},
				"3": Score{Total: 3, Centers: 0, CountBack1: "3"},
				"4": Score{Total: 4, Centers: 0, CountBack1: "4"},
				"5": Score{Total: 5, Centers: 0, CountBack1: "5"},
				"V": Score{Total: 5, Centers: 0, CountBack1: "6"},
				"6": Score{Total: 6, Centers: 0, CountBack1: "6"},
				"X": Score{Total: 6, Centers: 1, CountBack1: "7"},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      4,
			Grades:        []int{3, 4, 5, 6, 9},
		},
		ClassSettings{
			Name:         "match",
			Display:      "Match",
			DisplayValue: 2,
			Buttons:      "012345VX",
			SightersQty:  2,
			ShotsQty:     15,
			ValidShots: map[string]Score{
				"-": Score{Total: 0, Centers: 0, CountBack1: "0"},
				"0": Score{Total: 0, Centers: 0, CountBack1: "0"},
				"1": Score{Total: 1, Centers: 0, CountBack1: "1"},
				"2": Score{Total: 2, Centers: 0, CountBack1: "2"},
				"3": Score{Total: 3, Centers: 0, CountBack1: "3"},
				"4": Score{Total: 4, Centers: 0, CountBack1: "4"},
				"5": Score{Total: 5, Centers: 0, CountBack1: "5"},
				"V": Score{Total: 5, Centers: 1, CountBack1: "6"},
				"6": Score{Total: 5, Centers: 1, CountBack1: "6"},
				"X": Score{Total: 5, Centers: 1, CountBack1: "6"},
			},
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			GradeQty:      2,
			Grades:        []int{7, 8},
		},
	}

	CLASSES = []string{
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

type Grade struct{
	Name, LongName, ClassName string
	ClassId int
	Settings ClassSettings
}
func gradeList()[]int{
	return []int{0,1,2,3,4,5,6,7,8,9}
}
func grades()[]Grade{
	return []Grade{
		0: Grade{Settings: DEFAULT_CLASS_SETTINGS[0], ClassId: 0,Name:"A", ClassName:"Target", LongName: "Target A"},
		1: Grade{Settings: DEFAULT_CLASS_SETTINGS[0], ClassId: 0,Name:"B", ClassName:"Target", LongName: "Target B"},
		2: Grade{Settings: DEFAULT_CLASS_SETTINGS[0], ClassId: 0,Name:"C", ClassName:"Target", LongName: "Target C"},
		3: Grade{Settings: DEFAULT_CLASS_SETTINGS[1], ClassId: 1,Name:"FA", ClassName:"F Class", LongName: "F Class A"},
		4: Grade{Settings: DEFAULT_CLASS_SETTINGS[1], ClassId: 1,Name:"FB", ClassName:"F Class", LongName: "F Class B"},
		5: Grade{Settings: DEFAULT_CLASS_SETTINGS[1], ClassId: 1,Name:"F Open", ClassName:"F Class", LongName: "F Class Open"},
		6: Grade{Settings: DEFAULT_CLASS_SETTINGS[1], ClassId: 1,Name:"F/TR", ClassName:"F Class", LongName: "F/TR"},
		7: Grade{Settings: DEFAULT_CLASS_SETTINGS[2], ClassId: 2,Name:"Open", ClassName:"Match", LongName: "Match Open"},
		8: Grade{Settings: DEFAULT_CLASS_SETTINGS[2], ClassId: 2,Name:"Reserve", ClassName:"Match", LongName: "Match Reserve"},
		9: Grade{Settings: DEFAULT_CLASS_SETTINGS[1], ClassId: 1,Name:"Rifle", ClassName:"303", LongName: "303 Rifle"},
	}
}

func AGE_GROUPS2() []Option{
	return []Option{
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
}

type Legend struct{
	cssClass, name string
}
func scoreBoardLegend()[7]Legend{
	return [7]Legend{
		Legend{cssClass:"ST", name:"First"},
		Legend{cssClass:"ND", name:"Second"},
		Legend{cssClass:"TH", name:"Third"},
		Legend{cssClass:"w4", name:"Highest Possible Score"},
		Legend{cssClass:"w1", name:"Shoot Off"},
		Legend{cssClass:"w3", name:"Incomplete Score"},
		Legend{cssClass:"w2", name:"No Score"},
	}
}
