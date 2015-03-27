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
	TBLAutoInc     = "A"
	schemaCounter  = "n"
	TBLclub        = "C"
	TBLevent       = "E"
	schemaSHOOTER  = "^^schemaSHOOTER^^"
	schemaAutoInc  = "^^schemaAutoInc^^"
	schemaRANGE    = "^^schemaRANGE^^"
	schemaSORT     = "^^schemaSORT^^"
	schemaGRADES   = "^^schemaGRADES^^"
	TBLchamp       = "H"
	TBLshooter     = "S"
	TBLnraaShooter = "N"
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
	go DB()
}

func DB() {
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

func getClub(id string) (Club, error) {
	var result Club
	if conn != nil {
		err := conn.C(TBLclub).FindId(id).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get club with id: '" + id + "'")
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

func getShooterList(id int) Shooter {
	var result Shooter
	if conn != nil {
		conn.C(TBLnraaShooter).FindId(id).One(&result)
	}
	return result
}

func getNraaShooter(id int) NraaShooter {
	var result NraaShooter
	if conn != nil {
		conn.C(TBLnraaShooter).FindId(id).One(&result)
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

func nraaUpdateGrading(shooterId int, grades []NraaGrading) {
	change := mgo.Change{
		Upsert: false,
		Update: M{"$set": M{"g": grades}},
	}
	conn.C(TBLnraaShooter).FindId(shooterId).Apply(change, make(M))
}

func nraaLastUpdated() {
	conn.C(TBLAutoInc).FindId(TBLnraaShooter).Apply(mgo.Change{
		Upsert: true,
		Update: M{"$set": M{"n": time.Now().Format("January 2, 2006")}},
	}, make(M))
}

func getShooter(id int) Shooter {
	var result Shooter
	if conn != nil {
		conn.C(TBLshooter).FindId(id).One(&result)
	}
	return result
}

func getEvent(id string) (Event, error) {
	var result Event
	if conn != nil {
		err := conn.C(TBLevent).FindId(id).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with id: '" + id + "'")
}

func getEvent20Shooters(id string) (Event, error) {
	var result Event
	if conn != nil {
		err := conn.C(TBLevent).FindId(id).Select(M{"S": M{"$slice": -20}}).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with id: '" + id + "'")
}

func getNextId(collectionName string) (string, error) {
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
			return "", errors.New(fmt.Sprintf("Unable to generate the next ID: %v", err))
		}
		return idSuffix(result[schemaCounter].(int))
	}
	return "", errors.New("Unable to generate the next ID")
}

func idSuffix(id int) (string, error) {
	if id < 0 {
		return "", errors.New(fmt.Sprintf("Invalid id number supplied. Id \"%v\" is out of range", id))
	}
	id = id - 1
	charset := ID_CHARSET
	charsetLength := 70
	temp := ""
	for id >= charsetLength {
		temp = fmt.Sprintf("%c%v", charset[id%charsetLength], temp)
		id = id/charsetLength - 1
	}
	return fmt.Sprintf("%c%v", charset[id%charsetLength], temp), nil
}

func InsertDoc(collectionName string, data interface{}) {
	err := conn.C(collectionName).Insert(data)
	if err != nil {
		Error.Println(err)
	}
}

func UpdateDocById(collectionName, docId string, data interface{}) {
	err := conn.C(collectionName).UpdateId(docId, data)
	if err != nil {
		Error.Println(err)
	}
}

//Used for database schema translation dot notation
func Dot(elem ...interface{}) string {
	var dots []string
	for _, element := range elem {
		dots = append(dots, fmt.Sprintf("%v", element))
	}
	return strings.Join(dots, ".")
}

func EventAddRange(eventId string, newRange Range) (int, Event) {
	change := mgo.Change{
		Update: M{
			"$push": M{schemaRANGE: newRange},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	returned := Event{}
	conn.C(TBLevent).FindId(eventId).Apply(change, &returned)
	for rangeId, rangeData := range returned.Ranges {
		if rangeData.Name == newRange.Name && rangeData.Aggregate == newRange.Aggregate && rangeData.ScoreBoard == newRange.ScoreBoard && rangeData.Locked == newRange.Locked && rangeData.Hidden == newRange.Hidden {
			//TODO this if check is really hacky!!!
			return rangeId, returned
		}
	}
	return -1, returned
}

func eventShooterInsert(eventId string, shooter EventShooter) {
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
				Dot(schemaAutoInc, schemaSHOOTER): increment,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var event Event
	conn.C(TBLevent).FindId(eventId).Apply(change, &event)

	if increment == 2 {
		change = mgo.Change{
			Update: M{
				"$set": M{
					Dot(schemaSHOOTER, event.AutoInc.Shooter-2, "i"): event.AutoInc.Shooter - 2,
					Dot(schemaSHOOTER, event.AutoInc.Shooter-2, "l"): event.AutoInc.Shooter - 1,
					Dot(schemaSHOOTER, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter - 1,
					Dot(schemaSHOOTER, event.AutoInc.Shooter-1, "l"): event.AutoInc.Shooter - 2,
				},
			},
		}
	} else {
		change = mgo.Change{
			Update: M{
				"$set": M{
					Dot(schemaSHOOTER, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter - 1,
				},
			},
		}
	}
	conn.C(TBLevent).FindId(eventId).Apply(change, &event)
}

func eventTotalScoreUpdate(eventId string, rangeId int, shooterIds []int, score Score) Event {
	updateSetter := make(M)
	for _, shooterId := range shooterIds {
		updateSetter[Dot(schemaSHOOTER, shooterId, rangeId)] = score
	}
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": updateSetter,
		},
	}
	var event Event
	_, err := conn.C(TBLevent).FindId(eventId).Apply(change, &event)
	//TODO better error handling would be nice
	if err != nil {
		Warning.Println(err)
	}
	return event
}

func eventSortAggsWithGrade(event Event, rangeId string, shooterId int) {
	eventId := event.Id
	rangesToRedo := eventSearchForAggs(eventId, rangeId)
	//TODO this seems quite inefficient
	event = eventCalculateAggs(event, shooterId, rangesToRedo)
	//Only worry about shooters in this shooters grade
	currentGrade := event.Shooters[shooterId].Grade
	//Add the current range to the list of ranges to re-calculate
	rangesToRedo = append(rangesToRedo, rangeId)
	for _, rangeId := range rangesToRedo {
		// Closures that order the Change structure.
		//	grade := func(c1, c2 *EventShooter) bool {
		//		return c1.Grade < c2.Grade
		//	}
		total := func(c1, c2 *EventShooter) bool {
			return c1.Scores[rangeId].Total > c2.Scores[rangeId].Total
		}
		centa := func(c1, c2 *EventShooter) bool {
			return c1.Scores[rangeId].Centers > c2.Scores[rangeId].Centers
		}
		cb := func(c1, c2 *EventShooter) bool {
			return c1.Scores[rangeId].CountBack1 > c2.Scores[rangeId].CountBack1
		}

		//convert the map[string] to a slice of EventShooters
		var eventShooterList []EventShooter
		for thisShooterId, shooterList := range event.Shooters {
			if shooterList.Grade == currentGrade {
				shooterList.Id = thisShooterId
				for thisRangeId, score := range shooterList.Scores {
					score.Position = 0
					shooterList.Scores[thisRangeId] = score
				}
				eventShooterList = append(eventShooterList, shooterList)
			}
		}
		OrderedBy(total, centa, cb).Sort(eventShooterList)

		rank := 0
		nextOrdinal := 0
		//	score := 0
		//	center := 0
		//	countback := ""
		//	var previousShooter Shooter
		//		shooterLength := len(shooterList)

		//loop through the list of shooters
		for index, shooter := range eventShooterList {
			thisShooterScore := shooter.Scores[rangeId]

			//			if index+1 < shooterLength {
			//			if index-1 >= 0 {

			//keep track of the next badge position number to assign when several shooters are tied-equal on the position
			nextOrdinal += 1
			var nextShooterScore Score

			if index-1 >= 0 {
				nextShooter := eventShooterList[index-1]
				nextShooterScore = nextShooter.Scores[rangeId]

				//compare the shooters scores
				if thisShooterScore.Total == nextShooterScore.Total &&
					thisShooterScore.Centers == nextShooterScore.Centers &&
					thisShooterScore.CountBack1 == nextShooterScore.CountBack1 {
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
					"$set": M{Dot(schemaSHOOTER, shooter.Id, rangeId, "p"): rank},
				},
			}
			var result Event
			_, err := conn.C(TBLevent).FindId(eventId).Apply(change, &result)
			if err != nil {
				Warning.Printf("unable to update shooter rank for range: ", rangeId, ", shooter id:", shooter.Id)
			}
		}
	}
}

func eventUpdateRangeData(eventId string, updateData M) {
	conn.C(TBLevent).FindId(eventId).Apply(mgo.Change{
		Upsert: false,
		Update: updateData,
	}, make(M))
}

func eventUpdateSortScoreboard(eventId, sortByRange string) {
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": M{schemaSORT: sortByRange},
		},
	}
	conn.C(TBLevent).FindId(eventId).Apply(change, make(M))
}

func eventUpsertData(eventId string, data M) {
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": data,
		},
	}
	conn.C(TBLevent).FindId(eventId).Apply(change, make(M))
}

func tableUpdateData(collectionName, documentId string, data M) {
	change := mgo.Change{
		Upsert: false,
		Update: M{"$set": data},
	}
	conn.C(collectionName).FindId(documentId).Apply(change, make(M))
}
func UpsertDoc(collection string, id interface{}, document interface{}) {
	_, err := conn.C(collection).UpsertId(id, document)
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
