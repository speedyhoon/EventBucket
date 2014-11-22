package main

import (
	"fmt"
	"mgo"
	//	"os"
	"strings"
	//	"mgo/bson"
)

const (
//	DATABASE       = "eb"
	TBLAutoInc     = "A"
		schemaCounter  = "n"
	TBLclub        = "C"
	TBLevent       = "E"
		schemaSHOOTER  = "S"
		schemaAutoInc  = "U"
		schemaRANGE    = "R"
		schemaDATE     = "d"
		schemaTIME     = "t"
		schemaSORT     = "o"
		schemaGRADES   = "g"
	TBLchamp       = "c" //Championship
	TBLshooter     = "S"
	TBLshooterList = "n"
)

var (
	conn                *mgo.Database
	database_status     = false
	database_connection = 0
	//0 = not connected
	//1 = trying to connect
	//2 = connected
)

// Connect to the mongo database!
//func DB() *mgo.Database {
func DB() {
	fmt.Printf("database conn = %d\n", database_connection)
	database_connection = 1
	database_status = false
	session, err := mgo.Dial("localhost:38888")
	if err != nil {
		//TODO it would be better to output the mongodb connection error
		fmt.Printf("The database service is not reachable.")
		error_message(false, "999", "Database connection error", "The database service is not reachable. Please start the database service")
		remove_error("Initialising connection to DB")
		//		db_error_connection()
		database_connection = 0
		return
		//		os.Exit(999)
		//		return conn
	} //else{
	//		fmt.Printf("The database connected OK.")
	//	}
	//	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true) //this is supposed to be faster
	//	db_connection := session.DB(DATABASE)
	db_connection := session.DB("local")

	//	for _, table_name := range []string{TBLAutoInc, TBLclub, TBLevent, TBLchamp}{
	//		collection := db_connection.C(table_name)
	//		if collection != nil{
	//						db_error_connection()
	//			return
	//		}
	//	}
	database_status = true
	database_connection = 2
	conn = db_connection
	//	return db_connection
}
func DB_connection() {
	if database_connection == 0 {
		//		fmt.Println("Initialising connection to DB")
		//		error_message(true, "996", "Initialising connection to DB", "Initialising connection to DB")
		go DB()
	} else if database_connection == 1 {
		fmt.Println("Already connecting to DB")
		//		error_message(true, "997", "Initialising connection to DB", "Initialising connection to DB")
	} else {
		fmt.Println("connected to DB")
	}
}

func db_reconnect() {
	go DB()
}

func getCollection(collection_name string) []map[string]interface{} {
	var result []map[string]interface{}
	if database_status {
		checkErr(conn.C(collection_name).Find(nil).All(&result))
	}
	return result
}

func getClubs() []Club {
	var result []Club
	if database_status {
		checkErr(conn.C(TBLclub).Find(nil).All(&result))
	}
	return result
}
func getClub(id string) Club {
	var result Club
	if database_status {
		conn.C(TBLclub).FindId(id).One(&result)
	}
	return result
}


type M map[string]interface{}
func getClub_by_name(clubName string)(Club, bool){
	var result Club
	if database_status {
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
	if database_status {
		conn.C(TBLevent).Find(nil).All(&result)
	} else {
		DB_connection()
	}
	return result
}

func getShooterLists() []NRAA_Shooter {
	var result []NRAA_Shooter
	if database_status {
		conn.C(TBLshooterList).Find(nil).All(&result)
	} else {
		DB_connection()
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

func getEvent(id string)(Event, bool){
	var result Event

	if database_status {
//		checkErr(conn.C(TBLevent).FindId(id).One(&result))
		err := conn.C(TBLevent).FindId(id).One(&result)
		if err == nil{
			return result, false
		}
	}
	return result, true
}

func getNextId(collection_name string) string {
	var result map[string]interface{}
	if database_status {
		change := mgo.Change{
			Update:    map[string]interface{}{"$inc": map[string]interface{}{schemaCounter: 1}},
			Upsert:    true,
			ReturnNew: true,
		}
		_, err := conn.C(TBLAutoInc).FindId(collection_name).Apply(change, &result)
		if err != nil {
			checkErr(err)
		}
	}
	return id_suffix(result[schemaCounter].(int))
}

func id_suffix(id int) string {
	if id < 0 {
		error_message(false, "998", "Invalid id number supplied.", fmt.Sprintf("Id \"%v\" is out of range", id))
		return ""
	}
	id = id - 1
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!*()_-."
	//	fmt.Printf("charset length = %v", len(charset))
	charset_length := 70
	temp := ""
	for id >= charset_length {
		temp = fmt.Sprintf("%c%v", charset[id%charset_length], temp)
		id = id/charset_length - 1
	}
	return fmt.Sprintf("%c%v", charset[id%charset_length], temp)
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

func DB_event_add_range(event_id string, new_range Range) (string, Event) {
	event, _ := getEvent(event_id)
	change := mgo.Change{
		Update: map[string]interface{}{
			"$inc": map[string]interface{}{Dot(schemaAutoInc, schemaRANGE): 1},
			"$set": map[string]interface{}{Dot(schemaRANGE, event.AutoInc.Range): new_range},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	returned := Event{}
	conn.C(TBLevent).FindId(event_id).Apply(change, &returned)
	for range_id, range_data := range returned.Ranges {
		if range_data == new_range {
			return range_id, returned
		}
	}
	return "", returned
}

func event_shooter_insert(event_id string, shooter EventShooter) {
	event, _ := getEvent(event_id)
	change := mgo.Change{
		Update: map[string]interface{}{
			"$set": map[string]interface{}{Dot(schemaSHOOTER, event.AutoInc.Shooter): shooter},
			"$inc": map[string]interface{}{Dot(schemaAutoInc, schemaSHOOTER): 1},
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(map[string]interface{}))
}

func event_total_score_update(event_id, range_id, shooter_id string, score Score) {
	change := mgo.Change{
		Upsert: true,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{Dot(schemaSHOOTER, shooter_id, range_id): score},
		},
	}
	var event Event
	_, err := conn.C(TBLevent).FindId(event_id).Apply(change, &event)
	checkErr(err)
	//	aggs_list_to_update := search_for_aggs(event_id, range_id)

	event, _ = getEvent(event_id)
	ranges_to_redo := search_for_aggs(event_id, range_id)
	event = calculate_aggs(event, shooter_id, ranges_to_redo)
	UpdateDoc_by_id(TBLevent, event_id, event)



	//Get the up to date event
	event, _ = getEvent(event_id)

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
			//		if shooter
			//	}
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
//							fmt.Println("none")
							//					}
//						} else {
//							fmt.Println("exact")
							//					shoot_off = true
							//					shooter_list[index].Warning = 1
							//					score_board_legend_on_off["ShootOff"] = true

						}
					} else {
						//Shooters have a different score
						if this_shooter_score.Total != 0 {
							//increase rank by 1
							rank = next_ordinal
//							fmt.Println("go up")
						}else{
							rank = 0
//							fmt.Println("0=0=0")
						}
					}
				}else {
					//The very first shooter without a previous shooter assigned
					//increase rank by 1
					rank = next_ordinal
//					fmt.Println("go up")
				}
//				fmt.Println(shooter.Id, "rank:", rank, "  ", this_shooter_score.Total, " ", this_shooter_score.Centers, "  ", next_shooter_score.Total, " ", next_shooter_score.Centers, "   next:", next_ordinal)

				//update the database
				change := mgo.Change{
					Update: map[string]interface{}{                                          //position
						"$set": map[string]interface{}{Dot(schemaSHOOTER, shooter.Id, rangeId, "p"): rank},
					},
				}
				var result Event
				_, err := conn.C(TBLevent).FindId(event_id).Apply(change, &result)
				if err != nil {
					fmt.Println("unable to update shooter rank for range: ", rangeId, ", shooter id:", shooter.Id)
				}
//			}
		}
	}
}

func event_update_name(event_id, event_name string) {
	change := mgo.Change{
		Upsert: true,	//Maybe this shouldn't be upserted because name should ALWAYS be present
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"n": event_name},
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(map[string]interface{}))
}

func event_update_date(event_id, date, time string) {
	change := mgo.Change{
		Upsert: true,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{schemaDATE: date, schemaTIME: time}, //This is a separate fields because Browsers don't support a date-time field yet
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(map[string]interface{}))
}

func event_update_range_data(event_id string, update_data map[string]interface{}) {
	change := mgo.Change{
		Upsert: true,
		Update: update_data,
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(map[string]interface{}))
}

func event_update_sort_scoreboard(event_id, sort_by_range string) {
	change := mgo.Change{
		Upsert: true,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{schemaSORT: sort_by_range},
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(map[string]interface{}))
}

func event_upsert_data(event_id string, data map[string]interface{}) {
	change := mgo.Change{
		Upsert: true,
		Update: map[string]interface{}{
			"$set": data,
		},
	}
	conn.C(TBLevent).FindId(event_id).Apply(change, make(map[string]interface{}))
}

func nraa_upsert_shooter(shooter NRAA_Shooter) {
	_, err := conn.C("N").UpsertId(shooter.SID, &shooter)
	checkErr(err)
	fmt.Printf("inserted: %v\n", shooter)
}
func Upsert_Doc(collection string, id interface{}, document interface{}) {
	_, err := conn.C(collection).UpsertId(id, document)
	checkErr(err)
	fmt.Printf("inserted id: %v into %v\n", id, collection)
}

//func searchShooters(criteria Shooter)[]Shooter{
func searchShooters(query map[string]interface{}) []Shooter {
	//	var query map[string]interface{}
	/*	query := make(map[string]interface{}, 0)
		if criteria.Surname != "" {
	//		query["s"] = map[string]interface{}{"$regex": bson.RegEx{fmt.Sprintf(`^%v`, criteria.Surname), "i"}}
			query["s"] = criteria.Surname
		}
	//		query["s"] = map[string]interface{}{"$regex": bson.RegEx{fmt.Sprintf(`/^%v/i`, criteria.Surname), ""}}
		if criteria.FirstName != ""{
	//		query["f"] = map[string]interface{}{"$regex": bson.RegEx{fmt.Sprintf(`/^%v/i`, criteria.FirstName), ""}}
			query["f"] = criteria.FirstName
		}
		if criteria.Club != ""{
	//		query["c"] = map[string]interface{}{"$regex": bson.RegEx{fmt.Sprintf(`/^%v/i`, criteria.Club), ""}}
			query["c"] = criteria.Club
		}
	*/
	var result []Shooter

	//	integer, err := conn.C(TBLshooter).Find(bson.M{"s": bson.M{"$regex": bson.RegEx{`//Webb//`, ""}}}).Count()
	//	         er2 := conn.C(TBLshooter).Find(bson.M{"s": bson.M{"$regex": bson.RegEx{`Webb`, ""}}}).One(&result)
	//												 .Find(bson.M{"nm":bson.M{"$regex": bson.RegEx{`Andy.*`, ""}}}).One(&person)

	//	integer, err := conn.C(TBLshooter).Find(bson.M{"s": `\Webb\`}).Count()
	//	integer, err := conn.C(TBLshooter).Find(bson.M{"s": bson.M{"$regex": bson.RegEx{`Webb`, ""}}}).Count()
	//	err := conn.C(TBLshooter).Find(query).All(&result)
	//	                               map[string]interface{}{"s": map[string]interface{}{"$regex": "^Webb", "$options": "i"}, "f":map[string]interface{}{"$regex": "^C",       "$options": "i"}}
	//	err := conn.C(TBLshooter).Find(map[string]interface{}{"s": map[string]interface{}{"$regex": `^Webb`, "$options": "i"}, "f":map[string]interface{}{"$regex": `^cAmErOn`, "$options": "i"}}).All(&result)
	//	err := conn.C(TBLshooter).Find(query).All(&result)
	//	dump("\n\n\n\n fffffffffffffffffff:")
	//	export(query)
	//	dump("\n fffffffffffffffffff <<<\n\n\n")
	err := conn.C(TBLshooter).Find(query).All(&result)
	checkErr(err)

	//	dump("length:")
	//	dump(len(result))

	//	fmt.Printf("\nloggit \n%v\n...", integer)
	//	fmt.Printf("\nloggit \n%v\n...", er2)
	//	dump("search for\n")
	//	export(result)
	//	dump("done")
	return result

	//	err = c.Find(bson.M{"path": bson.M{"$regex": bson.RegEx{`^\\[^\\]*\\$`, ""}}}).All(&nodeList)
}
