package main

import (
	"fmt"
	"mgo"
	"strings"
	//	"os"
		"mgo/bson"
	"errors"
)

const (
	TBLAutoInc     = "A"
		schemaCounter  = "n"
	TBLclub        = "C"
	TBLevent       = "E"
		schemaSHOOTER  = "S"
		schemaAutoInc  = "U"
		schemaRANGE    = "R"
		schemaSORT     = "o"
		schemaGRADES   = "g"
	//TBLchamp       = "c" //Championship
	TBLshooter     = "S"
	TBLshooterList = "n"
)

var (
	conn                *mgo.Database
	database_status     = false
)

func DB() {
	database_status = false
	session, err := mgo.Dial("localhost:38888")
	if err != nil {
		//TODO it would be better to output the mongodb connection error
		Warning.Println("The database service is not available.")
		return
	}
	session.SetMode(mgo.Eventual, false)//false = the consistency guarantees won't be reset
	database_status = true
	conn = session.DB("local")
	//TODO defer colsing the session isn't working
//	defer session.Close()
}

func getCollection(collection_name string) []M {
	var result []M
	if database_status {
		checkErr(conn.C(collection_name).Find(nil).All(&result))
	}
	return result
}

func getClubs() []Club {
	var result []Club
	if conn != nil {
		checkErr(conn.C(TBLclub).Find(nil).All(&result))
	}
	return result
}
func getClub(id string)(Club, error){
	var result Club
	if database_status {
		err := conn.C(TBLclub).FindId(id).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get club with id: '"+id+"'")
}

func getClub_by_name(clubName string)(Club, bool){
	var result Club
	if database_status {
		//remove double spaces
		clubName = strings.Join(strings.Fields(clubName), " ")
		if clubName != "" {
			err := conn.C(TBLclub).Find(M{"n": M{"$regex": fmt.Sprintf(`^%v$`, clubName), "$options": "i"}}).One(&result)
			if err==nil {
				return result, true
			}
		}
	}
	return result, false
}

func getEvents() []Event {
	var result []Event
	if database_status {
		conn.C(TBLevent).Find(nil).All(&result)
	}
	return result
}

func getShooterLists() []NRAA_Shooter {
	var result []NRAA_Shooter
	if database_status {
		conn.C(TBLshooterList).Find(nil).All(&result)
	}
	return result
}

func getShooterList(id int) Shooter {
	var result Shooter
	if database_status {
		conn.C(TBLshooterList).FindId(id).One(&result)
	}
	return result
}

func getShooter(id int) Shooter {
	var result Shooter
	if database_status {
		conn.C(TBLshooter).FindId(id).One(&result)
	}
	return result
}

func getEvent(id string)(Event, error){
	var result Event
	//TODO is it possible to remove database_status and just use conn != nil?
	if conn != nil {
//	if database_status {
		err := conn.C(TBLevent).FindId(id).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with id: '"+id+"'")
}

func getEvent20Shooters(id string)(Event, bool){
	var result Event
	if database_status {
		err := conn.C(TBLevent).FindId(id).Select(bson.M{"S": bson.M{"$slice": -20 }}).One(&result)
		if err==nil{
			return result, false
		}
	}
	return result, true
}

func getNextId(collection_name string)(string, error){
	var result M
	if database_status {
		change := mgo.Change{
			Update:    M{"$inc": M{schemaCounter: 1}},
			Upsert:    true,
			ReturnNew: true,
		}
		_, err := conn.C(TBLAutoInc).FindId(collection_name).Apply(change, &result)
		if err != nil {
			checkErr(err)
			return "", errors.New(fmt.Sprintf("Unable to generate the next ID: %v", err))
		}
		return id_suffix(result[schemaCounter].(int))
	}
	return "", errors.New("Unable to generate the next ID")
}

func id_suffix(id int) (string, error) {
	if id < 0 {
		return "", errors.New(fmt.Sprintf("Invalid id number supplied. Id \"%v\" is out of range", id))
	}
	id = id - 1
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!*()_-."
	charset_length := 70
	temp := ""
	for id >= charset_length {
		temp = fmt.Sprintf("%c%v", charset[id%charset_length], temp)
		id = id/charset_length - 1
	}
	return fmt.Sprintf("%c%v", charset[id%charset_length], temp), nil
}

func InsertDoc(collection_name string, data interface{}) {
	checkErr(conn.C(collection_name).Insert(data))
}

func UpdateDoc_by_id(collection_name, doc_id string, data interface{}) {
	checkErr(conn.C(collection_name).UpdateId(doc_id, data))
}

//Used for database schema translation dot notation
func Dot(elem ...interface{}) string {
	var dots []string
	for _, element := range elem {
		dots = append(dots, fmt.Sprintf("%v", element))
	}
	return strings.Join(dots, ".")
}

func DB_event_add_range(event_id string, new_range Range) (int, Event) {
	change := mgo.Change{
		Update: M{
			"$push":M{schemaRANGE: new_range},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	returned := Event{}
	conn.C(TBLevent).FindId(event_id).Apply(change, &returned)
	for range_id, range_data := range returned.Ranges {
		if range_data.Name==new_range.Name&&range_data.Aggregate==new_range.Aggregate&&range_data.ScoreBoard==new_range.ScoreBoard&&range_data.Locked==new_range.Locked&&range_data.Hidden==new_range.Hidden{
			//TODO this if check seems really hacky!!!
			return range_id, returned
		}
	}
	return -1, returned
}


func event_shooter_insert(eventId string, shooter EventShooter) {
	insert := M{
		schemaSHOOTER: []EventShooter{shooter},
	}
	increment := 1
	if shooter.Grade == 8{
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
		Upsert: true,
		ReturnNew: true,
	}
	var event Event
	conn.C(TBLevent).FindId(eventId).Apply(change, &event)

	if increment == 2{
		change = mgo.Change{
			Update: M{
				"$set": M{
					Dot(schemaSHOOTER, event.AutoInc.Shooter-2, "i"): event.AutoInc.Shooter-2,
					Dot(schemaSHOOTER, event.AutoInc.Shooter-2, "l"): event.AutoInc.Shooter-1,
					Dot(schemaSHOOTER, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter-1,
					Dot(schemaSHOOTER, event.AutoInc.Shooter-1, "l"): event.AutoInc.Shooter-2,
				},
			},
		}
	}else{
		change = mgo.Change{
			Update: M{
				"$set": M{
					Dot(schemaSHOOTER, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter-1,
				},
			},
		}
	}

	conn.C(TBLevent).FindId(eventId).Apply(change, &event)
/*
	event, _ := getEvent(event_id)
	increment := 1
	insert := M{
		Dot(schemaSHOOTER, event.AutoInc.Shooter): shooter,
	}
	//Match Reserve
	if shooter.Grade == 8{
		//Shooters in grade Match Reserve also must go in grade Match Open
		shooter.LinkedId = event.AutoInc.Shooter + 1
		insert[fmt.Sprintf("%v", Dot(schemaSHOOTER, event.AutoInc.Shooter))] = shooter
		increment = 2
		duplicate := shooter
		duplicate.Grade = 7 //Match Open
		duplicate.Hidden = true
//		duplicate.LinkedId = fmt.Sprintf("%v", event.AutoInc.Shooter)	//TODO remove when not needed!
		insert[fmt.Sprintf("%v", Dot(schemaSHOOTER, event.AutoInc.Shooter + 1))] = duplicate
	}
	change := mgo.Change{
		Update: M{
			"$set": insert,
			"$inc": M{Dot(schemaAutoInc, schemaSHOOTER): increment},
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(M))
	*/
}

func eventTotalScoreUpdate(eventId string, rangeId int, shooterIds []int, score Score)Event{
	updateSetter := make(M)
	for _, shooterId := range shooterIds{
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
	checkErr(err)
	return event
}


func event_sort_aggs_with_grade(event Event, range_id string, shooter_id int){
	event_id := event.Id
	ranges_to_redo := search_for_aggs(event_id, range_id)
	//TODO this seems quite inefficent
	event = calculate_aggs(event, shooter_id, ranges_to_redo)
	//Only worry about shooters in this shooters grade
	current_grade := event.Shooters[shooter_id].Grade
	//Add the current range to the list of ranges to re-calculate
	ranges_to_redo = append(ranges_to_redo, range_id)
	for _, rangeId := range ranges_to_redo {
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
		var shooter_list []EventShooter
		for shooter_id, shooterList := range event.Shooters {
			if shooterList.Grade == current_grade {
				shooterList.Id = shooter_id
				for range_id, score := range shooterList.Scores {
					score.Position = 0
					shooterList.Scores[range_id] = score
				}
				shooter_list = append(shooter_list, shooterList)
			}
		}
		OrderedBy(total, centa, cb).Sort(shooter_list)

		rank := 0
		next_ordinal := 0
		//	score := 0
		//	center := 0
		//	countback := ""
		//	var previous_shooter Shooter
		//		shooter_length := len(shooter_list)

		//loop through the list of shooters
		for index, shooter := range shooter_list {
			this_shooter_score := shooter.Scores[rangeId]

			//			if index+1 < shooter_length {
			//			if index-1 >= 0 {

			//keep track of the next badge position number to assign when several shooters are tied-equal on the position
			next_ordinal += 1
			var next_shooter_score Score

			if index-1 >= 0 {
				next_shooter := shooter_list[index - 1]
				next_shooter_score = next_shooter.Scores[rangeId]

				//compare the shooters scores
				if this_shooter_score.Total == next_shooter_score.Total &&
					this_shooter_score.Centers == next_shooter_score.Centers &&
					this_shooter_score.CountBack1 == next_shooter_score.CountBack1 {
					//Shooters have an equal score
					if this_shooter_score.Total == 0 {
						//					shoot_equ = true
						//					if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
						rank = 0
						//					}
						//						} else {
						//							info("exact")
						//					shoot_off = true
						//					shooter_list[index].Warning = 1
						//					score_board_legend_on_off["ShootOff"] = true

					}
				} else {
					//Shooters have a different score
					if this_shooter_score.Total != 0 {
						//increase rank by 1
						rank = next_ordinal
					}else{
						rank = 0
					}
				}
			}else {
				//The very first shooter without a previous shooter assigned
				//increase rank by 1
				rank = next_ordinal
			}

			//update the database
			//TODO change this to only update once. not every loop iteration
			change := mgo.Change{
				Update: M{                                          //position
					"$set": M{Dot(schemaSHOOTER, shooter.Id, rangeId, "p"): rank},
				},
			}
			var result Event
			_, err := conn.C(TBLevent).FindId(event_id).Apply(change, &result)
			if err != nil {
				Warning.Printf("unable to update shooter rank for range: ", rangeId, ", shooter id:", shooter.Id)
			}
		}
	}
}

func event_update_range_data(event_id string, update_data M) {
	change := mgo.Change{
		Upsert: true,
		Update: update_data,
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(M))
}

func event_update_sort_scoreboard(event_id, sort_by_range string) {
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": M{schemaSORT: sort_by_range},
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(M))
}

func event_upsert_data(event_id string, data M) {
	change := mgo.Change{
		Upsert: true,
		Update: M{
			"$set": data,
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(M))
}

func nraa_upsert_shooter(shooter NRAA_Shooter) {
	_, err := conn.C("N").UpsertId(shooter.SID, &shooter)
	checkErr(err)
	Info.Printf("inserted: %v", shooter)
}
func Upsert_Doc(collection string, id interface{}, document interface{}) {
	_, err := conn.C(collection).UpsertId(id, document)
	checkErr(err)
	Info.Printf("inserted id: %v into %v", id, collection)
}

func searchShooters(query M) []Shooter {
	var result []Shooter
	checkErr(conn.C(TBLshooter).Find(query).All(&result))
	return result
}
