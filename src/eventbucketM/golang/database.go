package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// TBLAutoInc shouldn't be exported
	TBLAutoInc    = "A"
	schemaCounter = "n"
	// TBLclub shouldn't be exported
	TBLclub = "C"
	// TBLevent shouldn't be exported
	TBLevent      = "E"
	schemaSHOOTER = "^^schemaSHOOTER^^"
	schemaAutoInc = "^^schemaAutoInc^^"
	schemaRANGE   = "^^schemaRANGE^^"
	schemaSORT    = "^^schemaSORT^^"
	schemaGRADES  = "^^schemaGRADES^^"
	//TBLshooter     = "S"
	// TBLnraaShooter shouldn't be exported
	TBLnraaShooter = "N"
	//TBLchamp       = "H"
)

var conn *mgo.Database

func startDatabase() {
	databasePath := os.Getenv("ProgramData") + "/EventBucket"
	if !dirExists(databasePath) {
		Error.Printf("Can't find folder %v", databasePath)
		os.Mkdir(databasePath, os.ModeDir)
	}
	cmd := exec.Command("^^DbArgs^^")
	//TODO output mongodb errors/logs to stdOut
	/*stdout, err := cmd.StdoutPipe()
	if err != nil {
		Error.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Error.Println(err)
	}
	err = cmd.Start()*/
	cmd.Start()
	/*if err != nil {
		dump("exporting!")
		dump(stdout)
		dump(stderr)
		return
	}
	fmt.Println("Result: Err")
	dump(stderr)
	fmt.Println("Result: stdn out")
	dump(stdout)*/
	//TODO is a new goroutine really needed for DB call?
	go db()
}

func db() {
	session, err := mgo.Dial("localhost:38888")
	if err != nil {
		//TODO it would be better to output the mongodb connection error
		Error.Println("The database service is not available.", err)
		return
	}
	//session.SetMode(mgo.Eventual, true) //false = the consistency guarantees won't be reset
	session.SetMode(mgo.Strong, true)
	conn = session.DB("local")
	//TODO defer closing the session isn't working
	//defer session.Close()
}

/*func getCollection(collectionName string) []M {
	var result []M
	if conn != nil {
		err := conn.C(collectionName).Find(nil).All(&result)
		if err != nil {
			Warning.Println(err)
		}
	}
	return result
}*/

func getClubs() []Club {
	var result []Club
	if conn != nil {
		err := conn.C(TBLclub).Find(nil).All(&result)
		if err != nil {
			Warning.Println(err)
		}
	}
	return result
}

func getClub(ID string) (Club, error) {
	var result Club
	if conn != nil {
		err := conn.C(TBLclub).FindId(ID).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get club with ID: '" + ID + "'")
}

func getClubByName(clubName string) (Club, bool) {
	var result Club
	if conn != nil {
		//remove double spaces
		clubName = strings.Join(strings.Fields(clubName), " ")
		if clubName != "" {
			err := conn.C(TBLclub).Find(M{"n": M{"$regex": fmt.Sprintf(`^%v$`, clubName), "$options": "i"}}).One(&result)
			if err == nil {
				return result, true
			}
		}
	}
	return result, false
}

func getEvents() []Event {
	var result []Event
	if conn != nil {
		conn.C(TBLevent).Find(nil).All(&result)
	}
	return result
}

func getShooterLists() []NraaShooter {
	var result []NraaShooter
	if conn != nil {
		conn.C(TBLnraaShooter).Find(nil).All(&result)
	}
	return result
}

func getShooterList(ID int) Shooter {
	var result Shooter
	if conn != nil {
		conn.C(TBLnraaShooter).FindId(ID).One(&result)
	}
	return result
}

func getNraaShooter(ID int) NraaShooter {
	var result NraaShooter
	if conn != nil {
		conn.C(TBLnraaShooter).FindId(ID).One(&result)
	}
	return result
}

func nraaGetLastUpdated() string {
	var result map[string]string
	conn.C(TBLAutoInc).FindId(TBLnraaShooter).One(&result)
	return result["n"]
}

func nraaUpsertShooter(shooter NraaShooter) {
	_, err := conn.C(TBLnraaShooter).UpsertId(shooter.SID, &shooter)
	if err != nil {
		Warning.Println(err)
	}
}

func nraaUpdateGrading(shooterID int, grades []NraaGrading) {
	change := mgo.Change{
		Upsert: false,
		Update: M{"$set": M{"g": grades}},
	}
	conn.C(TBLnraaShooter).FindId(shooterID).Apply(change, make(M))
}

func nraaLastUpdated() {
	conn.C(TBLAutoInc).FindId(TBLnraaShooter).Apply(mgo.Change{
		Upsert: true,
		Update: M{"$set": M{"n": time.Now().Format("January 2, 2006")}},
	}, make(M))
}

/*func getShooter(ID int) Shooter {
	var result Shooter
	if conn != nil {
		conn.C(TBLshooter).FindId(ID).One(&result)
	}
	return result
}*/

func getEvent(ID string) (Event, error) {
	var result Event
	if conn != nil {
		err := conn.C(TBLevent).FindId(ID).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with ID: '" + ID + "'")
}

/*func getEvent20Shooters(ID string) (Event, error) {
	var result Event
	if conn != nil {
		err := conn.C(TBLevent).FindId(ID).Select(M{"S": M{"$slice": -20}}).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with ID: '" + ID + "'")
}*/

func getNextID(collectionName string) (string, error) {
	var result M
	if conn != nil {
		change := mgo.Change{
			Update:    M{"$inc": M{schemaCounter: 1}},
			Upsert:    true,
			ReturnNew: true,
		}
		_, err := conn.C(TBLAutoInc).FindId(collectionName).Apply(change, &result)
		if err != nil {
			Error.Println(err)
			return "", fmt.Errorf("Unable to generate the next ID: %v", err)
		}
		return idSuffix(result[schemaCounter].(int))
	}
	return "", errors.New("Unable to generate the next ID")
}

func idSuffix(ID int) (string, error) {
	if ID < 0 {
		return "", fmt.Errorf("Invalid ID number supplied. ID \"%v\" is out of range", ID)
	}
	ID--
	charset := idCharset
	charsetLength := 70
	temp := ""
	for ID >= charsetLength {
		temp = fmt.Sprintf("%c%v", charset[ID%charsetLength], temp)
		ID = ID/charsetLength - 1
	}
	return fmt.Sprintf("%c%v", charset[ID%charsetLength], temp), nil
}

func insertDoc(collectionName string, data interface{}) {
	err := conn.C(collectionName).Insert(data)
	if err != nil {
		Error.Println(err)
	}
}

func updateDocByID(collectionName, docID string, data interface{}) {
	err := conn.C(collectionName).UpdateId(docID, data)
	if err != nil {
		Error.Println(err)
	}
}

//Used for database schema translation dot notation
func dot(elem ...interface{}) string {
	var dots []string
	for _, element := range elem {
		dots = append(dots, fmt.Sprintf("%v", element))
	}
	return strings.Join(dots, ".")
}

func eventAddRange(eventID string, newRange Range) (int, Event) {
	change := mgo.Change{
		Update: M{
			"$push": M{schemaRANGE: newRange},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	returned := Event{}
	conn.C(TBLevent).FindId(eventID).Apply(change, &returned)
	for rangeID, rangeData := range returned.Ranges {
		if rangeData.Name == newRange.Name && rangeData.Aggregate == newRange.Aggregate && rangeData.ScoreBoard == newRange.ScoreBoard && rangeData.Locked == newRange.Locked && rangeData.Hidden == newRange.Hidden {
			//TODO this if check is really hacky!!!
			return rangeID, returned
		}
	}
	return -1, returned
}

func eventShooterInsert(eventID string, shooter EventShooter) {
	insert := M{
		schemaSHOOTER: []EventShooter{shooter},
	}
	//If shooter is Match Reserve, duplicate them in the Match Open category
	increment := 1
	if shooter.Grade == 8 {
		increment = 2
		duplicateShooter := shooter
		duplicateShooter.Grade = 7
		duplicateShooter.Hidden = true
		insert[schemaSHOOTER] = []EventShooter{shooter, duplicateShooter}
	}
	change := mgo.Change{
		Update: M{
			"$pushAll": insert,
			"$inc": M{
				dot(schemaAutoInc, schemaSHOOTER): increment,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var event Event
	conn.C(TBLevent).FindId(eventID).Apply(change, &event)

	if increment == 2 {
		change = mgo.Change{
			Update: M{
				"$set": M{
					dot(schemaSHOOTER, event.AutoInc.Shooter-2, "i"): event.AutoInc.Shooter - 2,
					dot(schemaSHOOTER, event.AutoInc.Shooter-2, "l"): event.AutoInc.Shooter - 1,
					dot(schemaSHOOTER, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter - 1,
					dot(schemaSHOOTER, event.AutoInc.Shooter-1, "l"): event.AutoInc.Shooter - 2,
				},
			},
		}
	} else {
		change = mgo.Change{
			Update: M{
				"$set": M{
					dot(schemaSHOOTER, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter - 1,
				},
			},
		}
	}
	conn.C(TBLevent).FindId(eventID).Apply(change, &event)
}

func eventTotalScoreUpdate(eventID string, rangeID int, shooterIDs []int, score Score) Event {
	updateSetter := make(M)
	for _, shooterID := range shooterIDs {
		updateSetter[dot(schemaSHOOTER, shooterID, rangeID)] = score
	}
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": updateSetter,
		},
	}
	var event Event
	_, err := conn.C(TBLevent).FindId(eventID).Apply(change, &event)
	//TODO better error handling would be nice
	if err != nil {
		Warning.Println(err)
	}
	return event
}

/*func eventSortAggsWithGrade(event Event, rangeID string, shooterID int) {
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
}*/

func eventUpdateRangeData(eventID string, updateData M) {
	conn.C(TBLevent).FindId(eventID).Apply(mgo.Change{
		Upsert: false,
		Update: updateData,
	}, make(M))
}

func eventUpdateSortScoreboard(eventID, sortByRange string) {
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": M{schemaSORT: sortByRange},
		},
	}
	conn.C(TBLevent).FindId(eventID).Apply(change, make(M))
}

func eventUpsertData(eventID string, data M) {
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": data,
		},
	}
	conn.C(TBLevent).FindId(eventID).Apply(change, make(M))
}

func tableUpdateData(collectionName, documentID string, data M) {
	change := mgo.Change{
		Upsert: false,
		Update: M{"$set": data},
	}
	conn.C(collectionName).FindId(documentID).Apply(change, make(M))
}
func upsertDoc(collection string, ID interface{}, document interface{}) {
	_, err := conn.C(collection).UpsertId(ID, document)
	if err != nil {
		Warning.Println(err)
	}
}

func searchShooters(query M) []Shooter {
	var result []Shooter
	err := conn.C(TBLnraaShooter).Find(query).All(&result) //TODO switch back to shooter
	//	err := conn.C(TBLshooter).Find(query).All(&result)
	if err != nil {
		Warning.Println(err)
	}
	return result
}
