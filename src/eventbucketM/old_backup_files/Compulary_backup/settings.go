package main
//truman Cell -- air purifier
//TODO: eventually replace these settings with ones that are set for each club and sometimes overridden by a clubs event settings
const (
	DEBUG		= true
	css_extension = ".css"
	icon_extension = ".png"


	rootDirectory = "/"

	//file extensions
	file_html = "htm"
	file_js = "js"
	file_css = "css"
	file_gif = "gif"
	file_jpeg = "jpg"
	file_png = "png"
	file_icon = "png"

	//folder structure
	dir_html = ""
	dir_js = "j"
	dir_css = "c"
	dir_gif = "i"
	dir_jpeg = "i"
	dir_png = "i"
	dir_icon = "i"

	nullShots = "-"	//record shots
	showMaxNumShooters = 20
	showInitialShots = 3		//the number of shots to show when a shooter is initially selected
	showMaxInitialShots = 4
	shotGroupingBorder = 3	//controlls where to place the shot separator/border between each number of shots
	borderBetweenSightersAndShots = true
	sighterGroupingBorder = 2
	indentFinish = false
	startShootingInputs = 0		//changes input text boxes to just tds for mobile input.
	allowClubNameWrap = true	//Club names have spaces between each word. false=Club names have &nbsp; between words
	startShootingDefaultSighter = "Drop All" //can select between 'Keep All' and 'Drop All'
	startShootingMaxNumShooters = 100	//can select between 'Keep All' and 'Drop All'

	//Start Shooting Page
	startShootingShowID =			-1
	startShootingShowUIN =			-2
	startShootingShowClass =		-3
	startShootingShowGrade =		4
	startShootingShowClub =			5
	startShootingShowShortName =	-6
	startShootingShowName =			7
	startShootingShowScores =		8
	startShootingShowTotal =		9
	startShootingShowReceived =	10
	//the columns to show and their order.

	scoreboardDisplayIndividuals = 1
	scoreboardCombineGrades = 0
	scoreboardShowTitle = 0 //1 = show, 0,-1 = hide titles -- show title of for syme or saturday/sunday etc
	scoreboardShowTeamsXs = 0 //1 = show, 0,-1 = hide Xs -- Agg columns if showXs == 1 display <sub>5Xs</sub>
	scoreboardShowTeamsShooters = 1 //1 = show, 0,-1 = hide Xs -- When set to 1 display Team shooters scores, When set to 0 only display teams totals.
	scoreboardShowShootOff = 0
	scoreboardShowInProgress = 1 //when enabled total score blinks while shooter is in progress

	scoreboardShowID = 1
	scoreboardShowShooterEntryID = -3//usefull to show entry id when a shooter is entered twice into the same event with different classes
	scoreboardShowUIN = -5
	scoreboardShowPosition = 100
	scoreboardShowGrade = 20
	scoreboardShowName = 30
	scoreboardShowClass = -40
	scoreboardShowClub = 70
	scoreboardShowGender = -70
	scoreboardShowAge = 80
	scoreboardShowShortName = -90
	scoreboardColorSecondRow = 0 //1=on, 0=off
	scoreboardShowRangeScores = 13

	// TODO: if one of the name options for scoreboard is not set then display the short name.
	// TODO: Add functionality to set these for javascript. output javascript code from golang. generate js file so it is cached and doesn't need to be generated on every page load.
	targetDesc = "Target Rifle 0-5 with V and X centers and able to convert Fclass scores to Target Rifle."
	targetSighters = 2
	targetShots = 10
	targetValidShots =  "012345V6X"
	targetValidScore =  "012345555"
	targetvcountback =  "012345666"
	targetxcountback =  "012345667"
	targetValidCenta =  "000000111"
	targetValidSighters = ")!@#$%v^x"
	targetValid =  ")!@#$%v^x012345V6X"
	targetValid2 = "012345V6X012345V6X"
	targetButtons = "012345VX"
	targetCountbackX = false
	targetCountbackValueX = 7
	targetShowInitialShots = 2
	targetShowMaxInitialShots = 2
	targetGrades = "A,B,C"

	matchDesc = "Match Rifle 0-5 with V and X centers and able to convert to Fclass scores to Match Rifle."
	matchSighters = 3
	matchShots = 20
	matchValidShots =		"012345V6X"
	matchValidScore =		"012345555"
	matchvcountback =		"012345666"
	matchxcountback =		"012345667"
	matchValidCenta =		"000000111"
	matchValidSighters =	")!@#$%v^x"
	matchValid =  ")!@#$%v^x012345V6X"
	matchValid2 = "012345V6X012345V6X"
	matchButtons = "012345VX"
	matchCountbackX = true
	matchCountbackValueX = 7
	matchShowInitialShots = 2
	matchShowMaxInitialShots = 2
	matchGrades = "MA,MB"

	fclassDesc = "Flcass 0-6 with X centers and able to convert Target and Match Rifle to Fclass scores."
	fclassSighters = 2
	fclassShots = 15
	fclassValidShots =		"012345V6X"
	fclassValidScore =		"012345666"
	fclassvcountback =		"012345667"
	fclassValidCenta =		"000000001"
	fclassValidSighters =	")!@#$%v^x"
	fclassValid =  ")!@#$%v^x012345V6X"
	fclassValid2 = "012345V6X012345V6X"
	fclassButtons = "0123456X"
	fclassCountbackX = false
	fclassCountbackValueX = 7
	fclassShowInitialShots = 2
	fclassShowMaxInitialShots = 2
	fclassGrades = "FA,FB,FO,FTR"

	//per Event
	shootOffSighters = 2
	shootOffShotsStart = 5
	shootOffnextShots = 3
	shootOffUseXcountback = 1	//1= true, 0=false
	shootOffUseXs = 1
	shootOffUseCountback = 1	//system settings
)
