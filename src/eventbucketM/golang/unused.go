package main

/**/

//DATABASE
func getCollection(collectionName string) []M {
	var result []M
	if conn != nil {
		err := conn.C(collectionName).Find(nil).All(&result)
		if err != nil {
			Warning.Println(err)
		}
	}
	return result
}

func getShooter(ID int) Shooter {
	var result Shooter
	if conn != nil {
		conn.C(TBLshooter).FindId(ID).One(&result)
	}
	return result
}

func getEvent20Shooters(ID string) (Event, error) {
	var result Event
	if conn != nil {
		err := conn.C(TBLevent).FindId(ID).Select(M{"S": M{"$slice": -20}}).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with ID: '" + ID + "'")
}

func eventSortAggsWithGrade(event Event, rangeID string, shooterID int) {
	eventID := event.ID
	rangesToRedo := eventSearchForAggs(eventID, rangeID)
	//TODO this seems quite inefficient
	event = eventCalculateAggs(event, shooterID, rangesToRedo)
	//Only worry about shooters in this shooters grade
	currentGrade := event.Shooters[shooterID].Grade
	//Add the current range to the list of ranges to re-calculate
	rangesToRedo = append(rangesToRedo, rangeID)
	for _, rangeID := range rangesToRedo {
		// Closures that order the Change structure.
		//	grade := func(c1, c2 *EventShooter) bool {
		//		return c1.Grade < c2.Grade
		//	}
		total := func(c1, c2 *EventShooter) bool {
			return c1.Scores[rangeID].Total > c2.Scores[rangeID].Total
		}
		centa := func(c1, c2 *EventShooter) bool {
			return c1.Scores[rangeID].Centres > c2.Scores[rangeID].Centres
		}
		cb := func(c1, c2 *EventShooter) bool {
			return c1.Scores[rangeID].CountBack > c2.Scores[rangeID].CountBack
		}

		//convert the map[string] to a slice of EventShooters
		var eventShooterList []EventShooter
		for thisShooterID, shooterList := range event.Shooters {
			if shooterList.Grade == currentGrade {
				shooterList.ID = thisShooterID
				for thisRangeID, score := range shooterList.Scores {
					score.Position = 0
					shooterList.Scores[thisRangeID] = score
				}
				eventShooterList = append(eventShooterList, shooterList)
			}
		}
		orderedBy(total, centa, cb).Sort(eventShooterList)

		rank := 0
		nextOrdinal := 0
		//	score := 0
		//	centre := 0
		//	countback := ""
		//	var previousShooter Shooter
		//		shooterLength := len(shooterList)

		//loop through the list of shooters
		for index, shooter := range eventShooterList {
			thisShooterScore := shooter.Scores[rangeID]

			//			if index+1 < shooterLength {
			//			if index-1 >= 0 {

			//keep track of the next badge position number to assign when several shooters are tied-equal on the position
			nextOrdinal++
			var nextShooterScore Score

			if index-1 >= 0 {
				nextShooter := eventShooterList[index-1]
				nextShooterScore = nextShooter.Scores[rangeID]

				//compare the shooters scores
				if thisShooterScore.Total == nextShooterScore.Total &&
					thisShooterScore.Centres == nextShooterScore.Centres &&
					thisShooterScore.CountBack == nextShooterScore.CountBack {
					//Shooters have an equal score
					if thisShooterScore.Total == 0 {
						//					shootEqu = true
						//					if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
						rank = 0
						//					}
						//						} else {
						//							info("exact")
						//					shootOff = true
						//					shooterList[index].Warning = 1
						//					scoreBoardLegendOnOff["ShootOff"] = true
					}
				} else {
					//Shooters have a different score
					if thisShooterScore.Total != 0 {
						//increase rank by 1
						rank = nextOrdinal
					} else {
						rank = 0
					}
				}
			} else {
				//The very first shooter without a previous shooter assigned
				//increase rank by 1
				rank = nextOrdinal
			}

			//update the database
			//TODO change this to only update once. not every loop iteration
			change := mgo.Change{
				Update: M{ //position
					"$set": M{dot(schemaSHOOTER, shooter.ID, rangeID, "p"): rank},
				},
			}
			var result Event
			_, err := conn.C(TBLevent).FindID(eventID).Apply(change, &result)
			if err != nil {
				Warning.Printf("unable to update shooter rank for range: ", rangeID, ", shooter ID:", shooter.ID)
			}
		}
	}
}

//DEV

func slice_to_map_bool(input []string) map[string]bool {
	output := make(map[string]bool)
	for _, value := range input {
		output[value] = true
	}
	return output
}

//HTTPGZIP

//TODO maybe add page redirects for pages. replace uppercase urls with lowercase
//Is this needed or not??????????????????????????
func redirectToUppercase() {
	//listOfPages := map[string]*func(){
	//	urlEventSettings: eventSettings,
	//}
	listOfPages := map[string]string{
		urlEventSettings: strings.ToLower(urlEventSettings),
	}
	for from, to := range listOfPages {
		GetRedirectPermanentTo(from, to)
	}
}
func GetRedirectPermanentTo(from, to string) {
	http.Handle(from, http.RedirectHandler(to, http.StatusMovedPermanently))
}

//CLUB
func getClubSelectionBox(clubList []Club) []Option {
	var dropDown []Option
	for _, club := range clubList {
		dropDown = append(dropDown, Option{Display: club.Name, Value: club.ID})
	}
	return dropDown
}

//PROD
func devModeCheckForm(check bool, message string) {
	if !check {
		Warning.Println(message)
	}
}

//TOTALSCORES
func eventSearchForAggs(eventID, rangeID string) []string {
	var aggsToCalculate []string
	event, _ := getEvent(eventID)
	for aggID, rangeData := range event.Ranges {
		if len(rangeData.Aggregate) > 0 {
			for _, thisRangeID := range rangeData.Aggregate {
				if string(thisRangeID) == rangeID {
					aggsToCalculate = append(aggsToCalculate, fmt.Sprintf("%v", aggID))
				}
			}
		}
	}
	return aggsToCalculate
}

//SCHEMA
type NraaGrading struct {
	Class string  `bson:"n,omitempty"`
	Grade      string  `bson:"r,omitempty"`
	Threshold string  `bson:"t,omitempty"`
	AvgScore       float64 `bson:"a,omitempty"`
	ShootQty       int     `bson:"s,omitempty"`
}

type TeamCat struct {
	Name string `bson:"n"`
}

type Team struct {
	name     string `bson:"n"`
	teamCat  []int  `bson:"t"`
	shooters []int  `bson:"s,omitempty"`
}

type Skill struct {
	Grade      string
	Percentage float64 //TODO would prefer an unsigned float here
}

//SETTINGS
var (
	//truman Cell -- air purifier
	//TODO: eventually replace these settings with ones that are set for each club and sometimes overridden by a clubs event settings
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

	TARGET_Desc = "Target Rifle 0-5 with V and X centres and able to convert Fclass scores to Target Rifle."
	MATCH_Desc  = "Match Rifle 0-5 with V and X centres and able to convert to Fclass scores to Match Rifle."
	FCLASS_Desc = "Flcass 0-6 with X centres and able to convert Target and Match Rifle to Fclass scores."

	//per Event
	SHOOTOFF_Sighters      = 2
	SHOOTOFF_ShotsStart    = 5
	SHOOTOFF_nextShots     = 3
	SHOOTOFF_UseXcountback = 1 //1= true, 0=false
	SHOOTOFF_UseXs         = 1
	SHOOTOFF_UseCountback  = 1 //system settings
)

func calculateHPS4Class(classID, numberOdShots int) Score {
	return Score{
		Total:   defaultClassSettings[classID].Maximum.Total * numberOdShots,
		Centres: defaultClassSettings[classID].Maximum.Centres * numberOdShots,
	}
}

func calculateHighestPossibleScores(numberOdShots int) []Score {
	var classHPS []Score
	for _, class := range defaultClassSettings {
		classHPS = append(classHPS, Score{
			Total:   class.Maximum.Total * numberOdShots,
			Centres: class.Maximum.Centres * numberOdShots,
		})
	}
	return classHPS
}

//UTILS
	//"bytes"
	//"encoding/base64"
	//"errors"
	//"github.com/boombuler/barcode"
	//"github.com/boombuler/barcode/datamatrix"
	//"github.com/boombuler/barcode/qr"
	//"image/png"

func exists(dict M, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func ternary(condition bool, True, False string) string {
	if condition {
		return True
	}
	return False
}

func imgBarcode(width, height int, barcodeType, value string) string {
	var data barcode.Barcode
	var err error
	switch barcodeType {
	case QRCODE:
		data, err = qr.Encode(value, qr.L, qr.Auto)
		break
	case DATAMATRIX:
		data, err = datamatrix.Encode(value)
		break
	default:
		err = errors.New("barcode type " + barcodeType + " is not implemented!")
		break
	}
	if err == nil {
		data, err = barcode.Scale(data, width, height)
		if err == nil {
			var buf bytes.Buffer
			err = png.Encode(&buf, data)
			if err == nil {
				return fmt.Sprintf("<img src=\"data:image/png;base64,%v\" width=%v height=%v alt=%v/>", base64.StdEncoding.EncodeToString(buf.Bytes()), width, height, value)
			}
		}
	}
	Error.Println(err)
	return ""
}





