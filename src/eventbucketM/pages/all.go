//to generate file run
//C:\Users\Developer\EBrepo\src\eventbucketM>TYPE .\*.go >> .\pages\all.go
//Then manually strip all imports and packages

package main
import(
"bytes"
"code.google.com/p/go.net/html"
"compress/gzip"
"encoding/json"
"fmt"
"go-randomdata-master"
"html/template"
"io"
"io/ioutil"
"io/ioutil"
"math/rand"
"mgo"
"net/http"
"net/url"
"os/exec"
"sort"
"strconv"
"strings"
"time"
"testing"
)

var urls = []string{
	"http://www.rubyconf.com/",
	"http://golang.org/",
	"http://matt.aimonetti.net/",
}
type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}
func asyncHttpGets(urls []string) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			ch <- &HttpResponse{url, resp, err}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.url)
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}
func main() {
	results := asyncHttpGets(urls)
	for _, result := range results {
		fmt.Printf("%s status: %s\n", result.url,
			result.response.Status)
	}
}

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

/*var (
	Info		*log.Logger = log.New(os.Stdout, "INFO: ",    log.Ldate|log.Ltime|log.Lshortfile)
	Warning	*log.Logger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func info(format string, a ...interface{}){
	if !PRODUCTION {
		Info.Printf(format, a...)
	}
}
func warning(format string, a ...interface{}){
	if !PRODUCTION {
		Warning.Printf(format, a...)
	}
}*/

func dev_mode_loadHTM(page_name string, existing_minified_file []byte) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name))
	checkErr(err)
	bytes = dev_mode_minifyHtml(page_name, bytes)
	existing_len := len(existing_minified_file)
	new_len := len(bytes)
	if existing_len != new_len {
		ioutil.WriteFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name), bytes, 0777)
		fmt.Printf("Page '%v' had %v bytes removed (%v percent), total: %v, from: %v", page_name, new_len-existing_len, (existing_len*100/new_len-100)*-1, existing_len, new_len)
		return bytes
	}
	return bytes
	//	return existing_minified_file
}

func dev_mode_check_form(check bool, message string){
	if !check{
		dump(message)
	}
}

func dev_mode_minifyHtml(page_name string, html []byte) []byte {
	minify := html

	if bytes.Contains(minify, []byte("ZgotmplZ")) {
		fmt.Println("Template generation error: ZgotmplZ")
		return []byte("")
	}

	remove_chars := map[string]string{
		"	": "", //Tab
		"\n": "", //new line
		"\r": "", //carriage return
	}
	//TODO remove spaces between block elements like: </div> <div> but keep between inline elements like </span> <span>
	//TODO use improved regex for better searching & replacement
	replace_chars := map[string]string{
		"  ":            " ", //double spaces
		"type=text":     "",
		"type=\"text\"": "",
		"type='text'":   "",
		" >":            ">",
		" <":            "<",
		"< ":            "<",
		">  <":          "> <",
		" />":           "/>",
		"/ >":           "/>",
		"<br/>":         "<br>",
		"</br>":         "<br>",
		"<br />":        "<br>",
	}
	for search, replace := range remove_chars {
		minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
	}

	backup := minify
	for !bytes.Equal(minify, backup) {
		for search, replace := range replace_chars {
			length := len(minify)
			minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
			if length != len(minify) {
				fmt.Printf("A dodgy character (%v) was found in the source! Please replace with (%v).", search, replace)
				//				warning("A dodgy character (%v) was found in the source! Please replace with (%v).", search, replace)
			}
		}
	}
	//TODO why is the string not being replaced here, even though it is 100% running?
	minify = bytes.Replace(minify, []byte("~~~"), []byte(" "), -1)
	return minify
}

func main() {
	DB_connection()
	file_server := http.FileServer(http.Dir(DIR_ROOT))
	//TODO make sure ALL resources don't have a . extension to save network bandwidth
	http.Handle(DIR_JS, file_headers_n_gzip(file_server, "js"))
	http.Handle(DIR_CSS, file_headers_n_gzip(file_server, "css"))
	http.Handle(DIR_PNG, file_headers_n_gzip(file_server, "png"))
	http.Handle(DIR_JPEG, file_headers_n_gzip(file_server, "jpg"))
	http.Handle(DIR_ICON, file_headers_n_gzip(file_server, "png"))
	http.Handle(DIR_SVG, file_headers_n_gzip(file_server, "svg"))

	//GET
	http.HandleFunc("/", html_headers_n_gzip(home))
	http.HandleFunc(URL_about, html_headers_n_gzip(about))
	http.HandleFunc(URL_about+"/", html_headers_n_gzip(redirectPermanent(URL_about)))
	http.HandleFunc(URL_licence, html_headers_n_gzip(licence))
	http.HandleFunc(URL_licence+"/", html_headers_n_gzip(redirectPermanent(URL_licence)))
	http.HandleFunc(URL_licence_summary, html_headers_n_gzip(licence_summary))
	http.HandleFunc(URL_licence_summary+"/", html_headers_n_gzip(redirectPermanent(URL_licence_summary)))

	http.HandleFunc(URL_organisers, html_headers_n_gzip(organisers))
	http.HandleFunc(URL_organisers+"/", html_headers_n_gzip(redirectPermanent(URL_organisers)))

	http.HandleFunc(URL_event, html_headers_n_gzip(event))
	http.HandleFunc(URL_eventSettings, html_headers_n_gzip(eventSettings))

	http.HandleFunc(URL_startShooting, html_headers_n_gzip(startShooting))
	http.HandleFunc(URL_startShootingAll, html_headers_n_gzip(startShootingAll))

	http.HandleFunc(URL_totalScores, html_headers_n_gzip(totalScores))
	http.HandleFunc(URL_totalScoresAll, html_headers_n_gzip(totalScoresAll))

	http.HandleFunc(URL_scoreboard, html_headers_n_gzip(scoreboard))
	http.HandleFunc(URL_club, html_headers_n_gzip(club))
//	http.HandleFunc(URL_clubs, html_headers_n_gzip(clubs))
	http.HandleFunc(URL_club_settings, html_headers_n_gzip(club_settings))
	//	//	http.HandleFunc("/clubs", clubs)
	//	http.HandleFunc("/events/", redirectPermanent("/events"))
	//	http.HandleFunc("/club/", html_headers_n_gzip(club))

	//Search for a shooter by first, surname & club
	http.HandleFunc(URL_queryShooterList, html_headers_n_gzip(queryShooterList))

	//POST
	http.HandleFunc(URL_clubInsert, html_headers_n_gzip(redirectVia(clubInsert, URL_organisers)))
	http.HandleFunc(URL_eventInsert, html_headers_n_gzip(redirectVia(eventInsert, URL_organisers)))

	http.HandleFunc(URL_eventInsert2, html_headers_n_gzip(eventInsert2))

	http.HandleFunc(URL_updateRange, html_headers_n_gzip(rangeUpdate2))
	http.HandleFunc(URL_dateUpdate, html_headers_n_gzip(dateUpdate))
	http.HandleFunc(URL_eventRangeInsert, html_headers_n_gzip(rangeInsert))
	http.HandleFunc(URL_eventAggInsert, html_headers_n_gzip(aggInsert))

	http.HandleFunc(URL_shooterInsert, html_headers_n_gzip(shooterInsert))
	http.HandleFunc(URL_shooterListInsert, html_headers_n_gzip(shooterListInsert))

	http.HandleFunc(URL_updateTotalScores, html_headers_n_gzip(updateTotalScores))
	http.HandleFunc(URL_updateShotScores, html_headers_n_gzip(updateShotScores))
	http.HandleFunc(URL_updateSortScoreBoard, html_headers_n_gzip(updateSortScoreBoard))

	http.HandleFunc(URL_updateEventName, html_headers_n_gzip(updateEventName))

	//Add shooters to event
	http.HandleFunc(URL_updateEventGrades, html_headers_n_gzip(updateEventGrades))
	http.HandleFunc(URL_updateShooterList, html_headers_n_gzip(redirectVia(updateShooterList, URL_organisers)))

	http.HandleFunc(URL_updateIsPrizeMeet, html_headers_n_gzip(updateIsPrizeMeet))

	http.HandleFunc(URL_eventShotsNSighters, html_headers_n_gzip(eventShotsNSighters))

	//	http.HandleFunc("/champInsert", redirectVia(champInsert, "/organisers"))
	//	http.HandleFunc("/clubMoundInsert", html_headers_n_gzip(clubMoundInsert))



	http.HandleFunc("/random-data/", html_headers_n_gzip(random_data))
	err := http.ListenAndServe(":80", nil)
	if err != nil{
//		log.Fatal("ListenAndServe: ", err)
		fmt.Printf("ListenAndServe: %v", err)
	}
	if PRODUCTION {
		url := "http://localhost/"
		err = exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", url).Start()
		if err != nil{
			fmt.Printf("Unable to open a web browser for %v", url)
		}
	}
}

func TestSqrt(t *testing.T) {
	const in, out = 4, 2
	if x := Sqrt(in); x != out {
		t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
	}
}


func generateForm2(form Form) string {
	var output, attributes, element, options string
	for name, input := range form.Inputs {
		element = ""
		attributes = ""
		options = ""
		dev_mode_check_form(input.Html!="submit"||input.Html!="number"||input.Html!="text"||input.Html!="range"||input.Html!="datalist"||input.Html!="select"||input.Html!="date"||input.Html!="hidden", "don't use element "+input.Html)

		if input.Html != "submit" {
			if name != "" {
				attributes += " name="+name
				dev_mode_check_form(name == addQuotes(name), "names can't have spaces")
			}
			if input.Value != "" {
				attributes += " value="+addQuotes(input.Value)
				dev_mode_check_form(input.Html != "select", "select boxes shouldn't have a value attribute")
			}
		}else{
			dev_mode_check_form(input.Value != "", "submits should have a value")
		}
		if input.Required {
			attributes += " required"
			dev_mode_check_form(input.Html=="number"||input.Html=="text"||input.Html=="range"||input.Html=="datalist"||input.Html=="select"||input.Html=="date", "this element shouldn't have required")
		}
		if input.Placeholder != "" {
			attributes += " placeholder="+addQuotes(input.Placeholder)
			dev_mode_check_form(input.Html == "text"||input.Html == "number"||input.Html == "range"||input.Html == "datalist", "placeholders are only allowed on text, datalist, number and ranges")
		}
		if input.Min != "" {
			attributes += fmt.Sprintf(" min=%v", input.Min)
			dev_mode_check_form(input.Html == "number" || input.Html == "range", "min is only allowed on type  number and range")
		}
		if input.Max != ""{
			attributes += fmt.Sprintf(" max=%v", input.Max)
			dev_mode_check_form(input.Html == "number" || input.Html == "range", "max is only allowed on type  number and range")
		}
		if input.Step != 0 {
			attributes += fmt.Sprintf(" step=%v", input.Step)
			dev_mode_check_form(input.Html == "number" || input.Html == "range", "step is only allowed on type  number and range")
		}
		if input.Checked {
			attributes += " checked"
			dev_mode_check_form(input.Html == "radio" || input.Html == "checkbox", "checked is only valid on radio buttons and checkboxes")
		}
		if input.Size > 0 {
			attributes += fmt.Sprintf(" size=%d", input.Size)
			dev_mode_check_form(input.Html == "select", "size is only allowed on select tags")
			dev_mode_check_form(input.Size >= 4, "size should be >= 4")
		}
		if input.AutoComplete != "" {
			attributes += " autocomplete="+input.AutoComplete
			dev_mode_check_form(input.Html == "datalist", "autocomplete is only allowed on datalist tags")
		}

		if input.MultiSelect {
			attributes += " multiple"
			if len(input.Options) > 4 {
				attributes += fmt.Sprintf(" size=%d", len(input.Options))
			}
			dev_mode_check_form(input.Html == "select", "multiple is only available on select boxes")
			dev_mode_check_form(input.Html != "submit", "buttons and submits shouldn't have multiple")
		}
		if len(input.Options) > 0 {
			options = draw_options(input, name)
		}
		if input.Help != "" {
			attributes += "title=" + addQuotes(input.Help)
		}


		if input.Html == "select" {
			element += "<select"+attributes+">"+options+"</select>"
		}else if input.Html == "submit" {
			output += "<button"+attributes+">"+input.Value+"</button>"
		}else {
			if input.Html == "datalist" && options != ""{
				attributes += " type=datalist id=" + name
			}
			if input.Html != "text" {
				attributes += " type="+input.Html
			}
			element += "<input"+attributes+">"+options
		}
		if input.Label != "" {
			output += "<label>"+input.Label+": "+element+"</label>"
			dev_mode_check_form(input.Html != "submit"||input.Html != "button", "submits and buttons shouldn't have lables")
		}
	}
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}else {
		dev_mode_check_form(false, "all forms should have a title")
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

type Form struct {
	Action string
	Title  string
	Inputs map[string]Inputs
}
type Inputs struct {
	//AutoComplete values can be: "off" or "on"
	Html, Label, Help, Value, Pattern, Placeholder, AutoComplete, Min, Max  string
	Checked, MultiSelect, Required bool
	Size                           int
	Options                        []Option
	Step                           float64
}
type Option struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

func draw_options(input Inputs, name string)string{
	export(input)
	dev_mode_check_form(len(input.Options) <= 0, "select should have at least one option to select from for element='"+name+"' type='"+input.Html+"'")
	dev_mode_check_form(len(input.Options) <= 0 && input.Required, "select shouldn't be required with no available options to select")
	output := ""
	if input.Placeholder != "" && input.Html != "datalist"{
		output += "<option selected value disabled>"+input.Placeholder+"</option>"
	}
	for _, option := range input.Options {
		output += "<option"
		if option.Selected {
			output += " selected"
			dev_mode_check_form(input.Html != "datalist", "datalist shouldn't have any selected values! change it to a value attribute")
			dev_mode_check_form(!(input.Placeholder != "" && input.Html != "datalist"), "shouldn't set a placeholder when options are already selected")
		}
		if option.Value != "" {
			output += " value" + addQuotesEquals(option.Value)
		}else {
			dev_mode_check_form(false, "option values shouldn't be empty")
		}
		output += ">" + option.Display + "</option>"
		dev_mode_check_form(option.Display != "", "option must have display text")
	}
	if input.Html == "datalist"{
		output = "<datalist id=" + name + ">"+output+"</datalist>"
		dev_mode_check_form(false,"make sure datalist id='"+name+"' is unique!")
	}
	return output
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}

func generateForm2(form Form) string {
	var output, attributes, element, options string
	for name, input := range form.Inputs {
		element = ""
		attributes = ""
		options = ""

		if name != "" {
			attributes += " name=" + name
			dev_mode_check_form(name == addQuotes(name), "names can't have spaces")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have names")
		}
		if input.Value != ""{
			if input.Html != "submit" {
				attributes += " value=" + addQuotes(input.Value))
			}
			dev_mode_check_form(input.Html != "select", "select boxes shouldn't have a value attribute")
		}else {
			dev_mode_check_form(input.Html != "submit", "submits should have a value")
		}
		if input.Required {
			attributes += " required"
			dev_mode_check_form(input.Html != "hidden" && input.Html != "button", "hidden inputs are not allowed to have required attributes")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have required")
		}
		if input.Placeholder != "" {
			attributes += " placeholder="+addQuotes(input.Placeholder)
			dev_mode_check_form(input.Html == "text"||input.Html == "number"||input.Html == "range", "placeholders are only allowed on text, number and ranges")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have placeholder")
		}
		if input.Min > -1 {
			attributes += fmt.Sprintf(" min=%v", input.Min)
			dev_mode_check_form(input.Html == "number" || input.Html == "range", "min is only allowed on type  number and range")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have min")
		}
		if input.Max > -1 {
			attributes += fmt.Sprintf(" max=%v", input.Max)
			dev_mode_check_form(input.Html == "number" || input.Html == "range", "max is only allowed on type  number and range")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have max")
		}
		if input.Step > 0 {
			attributes += fmt.Sprintf(" step=%v", input.Step)
			dev_mode_check_form(input.Html == "number" || input.Html == "range", "step is only allowed on type  number and range")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have step")
		}
		if input.Checked {
			attributes += " checked"
			dev_mode_check_form(input.Html == "radio" || input.Html == "checkbox", "checked is only valid on radio buttons and checkboxes")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have checked")
		}
		if input.Size > 0 {
			attributes += " size=%d" + input.Size
			dev_mode_check_form(input.Html == "select", "size is only allowed on select tags")
			dev_mode_check_form(input.Size >= 4, "size should be >= 4")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have size")
		}
		if input.AutoComplete != "" {
			attributes += " autocomplete="+input.AutoComplete
			dev_mode_check_form(input.Html == "datalist", "autocomplete is only allowed on datalist tags")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have autocomplete")
		}

		if input.MultiSelect {
			attributes += " multiple"
			if len(input.Options) > 4 {
				attributes += fmt.Sprintf(" size=%d", len(input.Options))
			}
			dev_mode_check_form(input.Html == "select", "multiple is only available on select boxes")
			dev_mode_check_form(input.Html != "submit" && input.Html != "button", "buttons and submits shouldn't have multiple")
		}
		if input.Html == "datalist"{
			attributes += " id=" + name
		}
		options = draw_options(input, name)
		if input.Help != "" {
			attributes += "title=" + addQuotes(input.Help)
		}


		if input.Html == "select" {
			element += "<select"+attributes+">"+options+"</select>"
		}else if input.Html == "submit" {
			element += "<button"+attributes+">"+input.Value+options+"</button>"
		}else {
			if input.Html != "text" {
				attributes += " type="+input.Html
			}
			element += "<input"+attributes+">"
		}
		if input.Label != "" {
			output += "<label>"+input.Label+": "+element+"</label>"
			dev_mode_check_form(input.Html != "submit"||input.Html != "button", "submits and buttons shouldn't have lables")
		}
	}
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}else {
		dev_mode_check_form(false, "all forms should have a title")
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

type Form struct {
	Action string
	Title  string
	Inputs map[string]Inputs
}
type Inputs struct {
	//AutoComplete values can be: "off" or "on"
	Html, Label, Help, Value, Pattern, Placeholder, AutoComplete  string
	Checked, MultiSelect, Required bool
	Min, Max, Size                 int
	Options                        []Options
	Step                           float64
}
type Options struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

func draw_options(input Inputs, name string)string{
	if len(input.Options) <= 0 {
		dev_mode_check_form(false, "select should have at least one option to select from")
		dev_mode_check_form(input.Required, "select shouldn't be required with no available options to select")
		return ""
	}
	output := ""
	if input.Placeholder != "" && input.Html != "datalist"{
		output += "<option selected value disabled>"+input.Placeholder+"</option>"
	}
	for _, option := range input.Options {
		output += "<option"
		if option.Selected {
			output += " selected"
			dev_mode_check_form(input.Html != "datalist", "datalist shouldn't have any selected values! change it to a value attribute")
			dev_mode_check_form(!(input.Placeholder != "" && input.Html != "datalist"), "shouldn't set a placeholder when options are already selected")
		}
		if option.Value != "" {
			output += " value" + addQuotesEquals(option.Value)
		}else {
			dev_mode_check_form(false, "option values shouldn't be empty")
		}
		output += ">" + option.Display + "</option>"
		dev_mode_check_form(option.Display != "", "option must have display text")
	}
	if input.Html == "datalist"{
		output = "<datalist id=" + name + ">"+output+"</datalist>"
		dev_mode_check_form(false,"make sure datalist id='"+name+"' is unique!")
	}
	return output
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
//6,493 bytes

func _required(attr bool)string{
	if attr {
		return " required"
	}
	return ""
}
func _name(attr string)string{
	if attr != ""{
		return " name=" + addQuotes(inputName)
	}
	return ""
}
func _multiSelect(inputData Inputs)string{
	output := ""
	if inputData.MultiSelect {
		output += " multiple"
		if len(inputData.SelectedValues) > 1 {
			output += fmt.Sprintf(" size=%d", len(inputData.SelectedValues))
		} else if len(inputData.SelectValues) > 1 {
			output += fmt.Sprintf(" size=%d", len(inputData.SelectValues))
		} else if len(inputData.Select) > 1 {
			output += fmt.Sprintf(" size=%d", len(inputData.Select))
		}
	}
	return output
}
func _help(attr string)string{
	if attr != "" {
		output += fmt.Sprintf("title=%v", addQuotes(attr))
	}
	return ""
}
func _value(attr string)string{
	if attr != "" {
		return fmt.Sprintf(" value=%v", addQuotes(attr))
	}
	return ""
}
func _autoComplete(attr string)string{
	if attr != ""{
		return " autocomplete=" + addQuotes(attr)
	}
	return ""
}
func _placeHolder(attr string)string{
	if attr != ""{
		return " placeholder=" + addQuotes(attr)
	}
	return ""
}
func _type(attr string)string{
	if attr != "text" {
		return " type=" + inputData.Html
	}
	return ""
}

func generateForm2(form Form) string {
	output := ""
	for inputName, inputData := range form.Inputs {
		if inputData.Html != "submit" && inputData.Label != "" {
			output += " <label>" + inputData.Label + ": "
		}
		if inputData.Html == "select" {
			output += "<select"
			output += _name(inputName)
			output += _multiSelect(inputData)
			output += _required(inputData.Required)
			output += _help(inputData.Help)
			output += ">"
			//TODO <option value="" disabled selected>Select your option</option>
			options, selected_options := draw_list_box(inputData.SelectedValues)
			if inputData.Placeholder != "" && !selected_options {
				output += fmt.Sprintf("<option disabled selected value>%v</option>", inputData.Placeholder)
			}
			for _, option := range inputData.Select {
				output += fmt.Sprintf("<option>%v</option>", option)
			}
			output += build_options_deprecated(inputData.SelectValues)

			output += options
			output += "</select>"
		} else if inputData.Html == "submit" {
			//TODO change all the submit labels to values
			if inputData.Value != "" && inputData.Label != "" {
				output += inputData.Label + " <button>" + inputData.Value + "</button>"
			} else {
				output += "<button>" + inputData.Label + "</button>"
			}
		} else if inputData.Html == "datalist" {
			output += "<input"
			output += _name(inputName)
			output += _required(inputData.Required)
			output += _value(inputData.Value)
			output += _autoComplete(inputData.AutoComplete)
			output += _placeHolder(inputData.Placeholder)
			if len(inputData.Options) > 0 {
				output += " list=" + inputName + "><datalist id=" + inputName + ">"
				for _, option := range inputData.Options {
					output += fmt.Sprintf("<option value=%v>", addQuotes(option))
				}
				output += "</datalist>"
			}else if len(inputData.SelectedValues) > 0 {
				//TODO remove this old datalist generator
				output += " list=" + inputName + ">"

				output += "<datalist id=" + inputName + ">"
				for value, option := range inputData.SelectValues {
					if value != "" {
						output += fmt.Sprintf("<option value=%v>%v</option>", addQuotes(value), option)
					} else {
						output += fmt.Sprintf("<option>%v</option>", option)
					}
				}
				output += "</datalist>"
			} else {
				output += ">"
			}
		} else {
			output += "<input"
			output += _type(inputData.Html)

			if inputData.Html != "submit" {
				output += _name(inputName)
			}
			output += _autoComplete(inputData.AutoComplete)
			if inputData.Html == "number" || inputData.Html == "range" {
				if inputData.Min > -1 {
					output += " min=" + echo(inputData.Min)
				}
				if inputData.Max > -1 {
					output += " max=" + echo(inputData.Max)
				}
				if inputData.Step > 0 {
					output += " step=" + echo(inputData.Step)
				}
			}
			if inputData.Checked {
				output += " checked"
			}
			if inputData.Size > 0 {
				output += fmt.Sprintf(" size=%d", inputData.Size)
			}
			output += _placeHolder(inputData.Placeholder)
			if inputData.Required {
				if inputData.Html != "hidden" {
					output += " required"
				} else {
					fmt.Println("\nhidden inputs are not allowed to have required attributes\n")
				}
			}
			output += _help(inputData.Help)
			output += _value(inputData.Value)
			output += ">"
		}
		if inputData.Html != "submit" && inputData.Label != "" {
			output += "</label>"
		}
	}
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

//TODO change map to slice of SelectedValues
//func build_options(options []SelectedValues)string{
func build_options_deprecated(options map[string]string) string {
	output := ""
	for value, option := range options {
		output += fmt.Sprintf("<option value%v>%v</option>", addQuotesEquals(value), option)
	}
	return output
}
func build_options(options []SelectedValues) string {
	output := ""
	selected := ""
	for _, option := range options {
		if option.Selected {
			selected = " selected"
		}
		output += fmt.Sprintf("<option%v value%v>%v</option>", selected, addQuotesEquals(option.Value), option.Display)
		selected = ""
	}
	return output
}

type SelectedValues struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

type Inputs struct {
	Html, Label, Help, Value       string
	Placeholder                    string
	Select                         []string
	SelectValues                   map[string]string
	SelectedValues                 []SelectedValues
	Options								 []string
	Checked, MultiSelect, Required bool
	Min, Max                       int
	Size                           int
	Step                           float64
	Pattern                        string

	//TODO maybe use bool? but how to detect if it is not set?
	AutoComplete						string	//0="off", 1="on"
}
type Form struct {
	Action string
	Title  string
	Inputs map[string]Inputs
}

func draw_list_box(options []SelectedValues) (string, bool) {
	output := ""
	selected_option := false
	for _, option := range options {
		output += "<option"
		if option.Selected {
			output += " selected"
			selected_option = true
		}
		if option.Value != "" {
			output += " value=" + addQuotes(option.Value)
		}
		output += ">" + option.Display + "</option>"
	}
	return output, selected_option
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}

func file_headers_n_gzip(h http.Handler, content_type string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RequestURI)
		http_headers(w, []string{"expire", "cache", content_type})
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(gzr, r)
		} else {
			h.ServeHTTP(w, r)
			fmt.Println("This request does not support gzip")
			//			Info.Println("This request does not support gzip")
		}
	}
}

func html_headers_n_gzip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RequestURI)
		http_headers(w, []string{"html", "nocache0", "nocache1", "nocache2"})
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(gzr, r)
		} else {
			h.ServeHTTP(w, r)
			fmt.Println("This request does not support gzip")
			//			Info.Println("This request does not support gzip")
		}
	}
}

func http_headers(w http.ResponseWriter, set_headers []string) {
	//	mins_in_year := 525949	//	hours_in_year := 8765.81
	//	one_year := time.Duration(mins_in_year)*time.Minute
	headers := map[string][2]string{
		"expire": [2]string{"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO it should return GMT time I think
		//		"expire":[2]string{"Expires", time.Now().UTC().Add(one_year).Format(time.RFC1123)},//TODO it should return GMT time I think
		//		"0cache":[2]string{"Expires", time.Now().UTC().Format(time.RFC1123)},//TODO it should return GMT time I think
		"nocache0": [2]string{"Cache-Control", "no-cache, no-store, must-revalidate"},
		"nocache1": [2]string{"Expires", "0"},
		"nocache2": [2]string{"Pragma", "no-cache"},
		"cache":    [2]string{"Vary", "Accept-Encoding"},
		"csp":      [2]string{"Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';"}, //content-security-policy.com
		"gzip":     [2]string{"Content-Encoding", "gzip"},
		"html":     [2]string{"Content-Type", "text/html; charset=utf-8"},
		"css":      [2]string{"Content-Type", "text/css; charset=utf-8"},
		//TODO which mime type is best for javascript?
		//"js":		[2]string{"Content-Type", "application/javascript"},
		"js":  [2]string{"Content-Type", "text/javascript"},
		"png": [2]string{"Content-Type", "image/png"},
		"jpg": [2]string{"Content-Type", "image/jpeg"},
		"gif": [2]string{"Content-Type", "image/gif"},
		//TODO find valid mime type for webp & svg
		"webp": [2]string{"Content-Type", "text/webp"},
		"svg":  [2]string{"Content-Type", "image/svg+xml"},
	}
	for _, lookup := range set_headers {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func redirectPermanent(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently) //Search engine Optimisation
	}
}

func redirectVia(runThisFirst func(http.ResponseWriter, *http.Request), path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		runThisFirst(w, r)
		http.Redirect(w, r, path, http.StatusSeeOther) //303 mandating the change of request type to GET
	}
}

func redirecter(path string, w http.ResponseWriter, r *http.Request) {
	//	return func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, path, http.StatusSeeOther)
	//	}
}

func get_id_from_url(r *http.Request, page_url string) string {
	//TODO add validation checking for id using regex pattens
	/*TODO add a http layer function between p_page functions and main.go so that the event_id or club_id can
	be validated and the p_page functions don't have to interact with http at all*/
	//	var validID = regexp.MustCompile(`\A` + page_url + `[0-9a-f]{24}\z`)
	url := fmt.Sprintf("%v", r.URL)
	//	if validID.MatchString(url) {
	//		templator("admin", eventSettings_HTML(), eventSettings_Data(url[len(page_url):]), w)
	//	}else {
	//		redirectPermanent("/events")
	//		fmt.Println("redirected user " + url)
	//	}
	return url[len(page_url):]
}

func timeTrack(start time.Time, requestURI string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", requestURI, elapsed)
}

//A better solution to gzip over http!
//package main
//
//import (
//"compress/gzip"
//"flag"
//"io"
//"http"
//"log"
//"os"
//"path"
//"strings"
//)
//
//var rootdir *string = flag.String("rootdir", "/home/pkf/intraday", "chroot to this directory.")
//var prefix *string = flag.String("prefix", "/intraday/", "prefix path in URLs")
//
//func checkencoding(req *http.Request) bool {
//	encoding := req.Header["Accept-Encoding"]
//	if encoding != "" {
//		for _, v := range strings.Split(encoding, ",", 10) {
//			if strings.TrimSpace(v) == "gzip" {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//// send a file, with optional compression
//func handler(c *http.Conn, req *http.Request) {
//
//	// we only support GET
//	if req.Method != "GET" {
//		log.Stderrf("req.Method is %s", req.Method)
//		return
//	}
//
//	// should we compress?
//	compress := checkencoding(req)
//	log.Stderrf("%s; compression requested: %v", req.URL.Path, compress)
//
//	// clean the path and make sure the file exists
//	// handles the case where people try to fetch /../../../../etc/passwd or something
//	cpath := path.Clean(req.URL.Path)
//	if !strings.HasPrefix(cpath, *prefix) {
//		http.NotFound(c, req)
//		return
//	}
//
//	cpath = cpath[len(*prefix):]
//
//	file, _ := os.Open(cpath, os.O_RDONLY, 0)
//	if file == nil {
//		http.NotFound(c, req)
//		return
//	}
//	defer file.Close()
//
//	// write it out
//	c.SetHeader("Content-Type", "text-plain; charset=us-ascii")
//
//	if compress {
//		comp, _ := gzip.NewWriter(c)
//		defer comp.Close()
//
//		c.SetHeader("Content-Encoding", "gzip")
//
//		buf := make([]byte, 1048576)
//		for {
//			nbytes, err := file.Read(buf)
//			if err != nil {
//				return
//			}
//			comp.Write(buf[0:nbytes])
//		}
//	} else {
//		io.Copy(c, file)
//	}
//}
//
//func main() {
//	flag.Parse()
//	err := os.Chdir(*rootdir)
//	if err != nil {
//		log.Stderr(err.String())
//		os.Exit(1)
//	}
//	log.Stderrf("starting in %s, expecting urls to begin with %s", *rootdir, *prefix)
//	http.ListenAndServe(":12345", http.HandlerFunc(handler))
//}

func about(w http.ResponseWriter, r *http.Request) {
	templator(TEMPLATE_HOME, "about", aboutData(), w)
}

func aboutData() map[string]interface{} {
	return map[string]interface{}{
		"Version":  VERSION,
		"PageName": "About",
		"Menu":     home_menu("/about",HOME_MENU_ITEMS),
	}
}

func club_settings(w http.ResponseWriter, r *http.Request) {
	club_id := get_id_from_url(r, URL_club)
	templator(TEMPLATE_ADMIN, "club-settings", club_settings_Data(club_id), w)
}

func club_settings_Data(club_id string) map[string]interface{} {
	this_club := getClub(club_id)
	var temporary []string
	for _, mound := range this_club.Mounds {
		temporary = append(temporary, generateForm2(clubMoundUpdateForm(this_club.Id, mound)))
	}
	return map[string]interface{}{
		"Name":            this_club.Name,
		"Id":              this_club.Id,
		"ListMounds":      temporary,
		"Latitude":        this_club.Latitude,
		"Longitude":       this_club.Longitude,
		"InsertRangeForm": generateForm2(clubMoundInsertForm(this_club.Id)),
		"MapForm":         generateForm2(clubMapUpsertForm(this_club.Id)),
	}
}

//func clubs(w http.ResponseWriter, r *http.Request) {
//	templator(TEMPLATE_HOME, "clubs", clubs_Data(), w)
//}

func clubMoundInsertForm(club_id string) Form {
	return Form{
		Action: "clubMoundInsert",
		Title:  "Insert Mound",
		Inputs: map[string]Inputs{
			"clubid": {
				Html:  "hidden",
				Value: club_id,
			},
			"distance": {
				Html:     "number",
				Label:    "Distance",
				Required: true,
				Min:      "1",
			},
			"unit": {
				Html:     "select",
				Required: true,
				Label:    "Unit",
				Options:   []Option{Option{Value:"Yards",Display:"Yards"},Option{Value:"Metres",Display:"Metres"}},
			},
			"submit": {
				Html:  "submit",
				Value: "Insert New Mound",
			},
		},
	}
}

//if inputData.Min != 0 {
//output += fmt.Sprintf(" min=%f",inputData.Min)
//}
//if inputData.Max != 0{
//output += fmt.Sprintf(" max=%f",inputData.Max)
//}
//if inputData.Step != 0{
//output += fmt.Sprintf(" step=%f",inputData.Step)
//}

func clubMapUpsertForm(club_id string) Form {
	return Form{
		Action: "clubMapUpsert",
		Title:  "Update Map Location",
		Inputs: map[string]Inputs{
			"clubid": {
				Html:  "hidden",
				Value: club_id,
			},
			"latitude": {
				Html:     "number",
				Label:    "Latitude",
				Required: true,
				Min:      "-90",
				Step:     0.000001,
				Max:      "90",
			},
			"longitude": {
				Html:     "number",
				Required: true,
				Label:    "Longitude",
				Min:      "-180",
				Step:     0.000001,
				Max:      "180",
			},
			"submit": {
				Html:  "submit",
				Value: "Update Co-ordinates",
			},
		},
	}
}

func clubMoundUpdateForm(club_id string, mound Mound) Form {
	return Form{
		Action: "clubMoundUpdate",
		Inputs: map[string]Inputs{
			"clubid": {
				Html:  "hidden",
				Value: club_id,
			},
			"distance": {
				Html:     "number",
				Label:    "Distance",
				Required: true,
				Value:    echo(mound.Distance),
			},
			"unit": {
				Html:     "select",
				Required: true,
				Label:    "Unit",
				Options:   []Option{Option{Value:"Yards", Display:"Yards"},Option{Value:"Metres", Display:"Metres"}},
			},
			"submit": {
				Html:  "submit",
				Value: "Update Mound",
			},
		},
	}
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(clubMoundInsertForm("").Inputs, r)
	club_id := validated_values["clubid"]
	redirecter(fmt.Sprintf("/club/%v", club_id), w, r)
	var new_mound Mound
	new_mound.Distance = str_to_int(validated_values["distance"])
	new_mound.Unit = validated_values["unit"]
	//	club_insert_mound(club_id, new_mound)
}

func club(w http.ResponseWriter, r *http.Request) {
	club_id := get_id_from_url(r, URL_club)
	//TODO change club url to be the club.Url instead of club.Id
	templator(TEMPLATE_HOME, "club", club_Data(club_id), w)
}

func club_Data(club_id string) map[string]interface{} {
	this_club := getClub(club_id)
	//	menu_items := []Menu{
	//		Menu{
	//			Name:		"Home",
	//			Link:		"/",
	//		},
	//		Menu{
	//			Name:		"Clubs",
	//			Link:		"/clubs",
	//		},
	//		Menu{
	//			Name:		"Events",
	//			Link:		"/events",
	//		},
	//		Menu{
	//			Name:		"Event Archive",
	//			Link:		"/archive",
	//		},
	//		Menu{
	//			Name:		"Organisers",
	//			Link:		"/organisers",
	//		},
	//	}
	var temporary []string
	for _, mound := range this_club.Mounds {
		temporary = append(temporary, generateForm2(clubMoundUpdateForm(this_club.Id, mound)))
	}
	return map[string]interface{}{
		"Name":       this_club.Name,
		"Menu":       home_menu("/clubs",HOME_MENU_ITEMS),
		"ClubId":     this_club.Id,
		"ListMounds": this_club.Mounds,
		"Latitude":   this_club.Latitude,
		"Longitude":  this_club.Longitude,
	}
}

//func clubMoundInsertForm(club_id string) Form{
//	return Form{
//		Action: URL_club_mound_insert,
//		Title: "foobar",
//		Inputs: map[string]Inputs{
//			"clubid": {
//				Html:  "hidden",
//				Value: club_id,
//			},
//			"distance": {
//				Html:   "number",
//				Label: "Distance",
//				Required: true,
//				Min: 1,
//			},
//			"unit": {
//				Html: "select",
//				Required: true,
//				Label: "Unit",
//				Select: []string{"Yards", "Metres"},
//			},
//			"submit": {
//				Html:  "submit",
//				Value: "Insert New Mound",
//			},
//		},
//	}
//}
//
//func clubMoundUpdateForm(club_id string, mound Mound) Form {
//	return Form{
//		Action: URL_club_mound_update,
//		Title: "ellh!!",
//		Inputs: map[string]Inputs{
//			"clubid": {
//				Html:  "hidden",
//				Value: club_id,
//			},
//			"distance": {
//				Html:   "number",
//				Label: "Distance",
//				Required: true,
//				Value: echo(mound.Distance),
//			},
//			"unit": {
//				Html: "select",
//				Required: true,
//				Label: "Unit",
//				Select: []string{"Yards", "Metres"},
//			},
//			"submit": {
//				Html:  "submit",
//				Value: "Update Mound",
//			},
//		},
//	}
//}

//func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
//	validated_values := check_form(clubMoundInsertForm(""), r)
//	club_id := validated_values["clubid"]
//	redirecter(fmt.Sprintf("/club/%v",club_id), w, r)
//	var new_mound Mound
//	new_mound.Distance = str_to_int(validated_values["distance"])
//	new_mound.Unit = validated_values["unit"]
//	club_insert_mound(club_id, new_mound)
//}

//
//import (
//	"net/http"
//	"fmt"
//)
//
//func clubs(w http.ResponseWriter, r *http.Request) {
//	templator(TEMPLATE_HOME, "club", clubs_Data(), w)
//}
//
//func clubs_Data(club_id string) map[string]interface{} {
//	this_club := getClub(club_id)
//	var temporary []string
//	for _, mound := range this_club.Mounds{
//		temporary = append(temporary, generateForm("clubMoundUpdate", clubMoundUpdateForm(this_club.Id, mound)))
//	}
//	return map[string]interface{}{
//		"Title": this_club.Name,
//		"Id": this_club.Id,
//		"ListMounds": temporary,
//		"Latitude": this_club.Latitude,
//		"Longitude": this_club.Longitude,
//		"InsertRangeForm": generateForm("clubMoundInsert", clubMoundInsertForm(this_club.Id)),
//	}
//}
//
//func clubMoundInsertForm(club_id string) map[string]Inputs {
//	return map[string]Inputs{
//		"clubid": {
//			Html:  "hidden",
//			Value: club_id,
//		},
//		"distance": {
//			Html:	"number",
//			Label: "Distance",
//			Required: true,
//			Min: 1,
//		},
//		"unit": {
//			Html: "select",
//			Required: true,
//			Label: "Unit",
//			Select: []string{"Yards","Metres"},
//		},
//		"submit": {
//			Html:  "submit",
//			Label: "Insert New Mound",
//		},
//	}
//}
//
//func clubMoundUpdateForm(club_id string, mound Mound) map[string]Inputs {
//	return map[string]Inputs{
//		"clubid": {
//			Html:  "hidden",
//			Value: club_id,
//		},
//		"distance": {
//			Html:	"number",
//			Label: "Distance",
//			Required: true,
//			Value: echo(mound.Distance),
//		},
//		"unit": {
//			Html: "select",
//			Required: true,
//			Label: "Unit",
//			Select: []string{"Yards","Metres"},
//		},
//		"submit": {
//			Html:  "submit",
//			Label: "Update Mound",
//		},
//	}
//}
//
//func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
//	validated_values := check_form(clubMoundInsertForm(""), r)
//	club_id := validated_values["clubid"]
//	redirecter(fmt.Sprintf("/club/%v",club_id), w, r)
//	var new_mound Mound
//	new_mound.Distance = str_to_int(validated_values["distance"])
//	new_mound.Unit = validated_values["unit"]
//	club_insert_mound(club_id, new_mound)
//}

func updateEventName(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_event_name("","").Inputs, r)
	event_id := validated_values["event_id"]
	event_update_name(event_id, validated_values["name"])
	redirecter(URL_eventSettings+event_id, w, r)
}

func rangeInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_rangeForm("").Inputs, r)
	range_agg_insert(validated_values)
	event_id := validated_values["event_id"]
	referer := URL_event
	if strings.Contains(r.Header["Referer"][:1][0], URL_eventSettings) {
		referer = URL_eventSettings
	}
	redirecter(referer+event_id, w, r)
}

func aggInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_aggForm("", []Option{}).Inputs, r)
	range_agg_insert(validated_values)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
}
func range_agg_insert(validated_values map[string]string) {
	var new_range Range
	new_range.Name = validated_values["name"]
	//	dump("range_agg_insert")
	//	export(validated_values["agg"])
	if validated_values["agg"] != "" {
		//		new_range.Aggregate = []int64{}
		//		for _, range_name := range strings.Split(validated_values["agg"], ","){
		//			new_range.Aggregate = append(new_range.Aggregate, str_to_int64(range_name))
		//		}
		//		new_range.Aggregate = strings.Split(validated_values["agg"], ",")
		new_range.Aggregate = validated_values["agg"]
	}
	event_id := validated_values["event_id"]
	range_id, event_data := DB_event_add_range(event_id, new_range)
	if range_id != "" {
		go calc_new_agg_range_scores(event_id, range_id, event_data)
	}
}
func calc_new_agg_range_scores(event_id, range_id string, event Event) {
	ranges := []string{range_id}
	for shooter_id := range event.Shooters {
		event = calculate_aggs(event, shooter_id, ranges)
	}
	UpdateDoc_by_id(TBLevent, event_id, event)
}

func rangeUpdate2(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_update_range("", "").Inputs, r)
	dump(validated_values)
	//	event_id := get_id_from_url(r, URL_eventSettings)
	//	templator(TEMPLATE_ADMIN, "event-settings", eventSettings_Data(event_id), w)

	event_id := validated_values["event_id"]
	range_id := validated_values["range_id"]
//	hide := validated_values["hide"]
//	lock := validated_values["lock"]

	unset := make(map[string]interface{})
	set := map[string]interface{}{
		Dot("R",range_id,"n"): validated_values["name"],
	}

	did_unset := false

	if validated_values["hide"] == "on"{
//		set update["$set"][Dot("R",range_id,"h")] = true
		set[Dot("R",range_id,"h")] = true
	}else{
//		update["$unset"][Dot("R",range_id,"h")] = ""
		unset[Dot("R",range_id,"h")] = ""
		did_unset = true
	}
	if validated_values["lock"] == "on"{
//		update["$set"][Dot("R",range_id,"l")] = true
		set[Dot("R",range_id,"l")] = true
	}else{
//		update["$unset"][Dot("R",range_id,"l")] = ""
		unset[Dot("R",range_id,"l")] = ""
		did_unset = true
	}

	update := map[string]interface{}{
		"$set": set,
	}
	if did_unset{
		update["$unset"] = unset
	}
	event_update_range_data(event_id, update)
	redirecter(URL_eventSettings+event_id, w, r)
}
func eventSettings_update_range(event_id, range_id string) Form {
	return Form{
		Action: URL_updateRange,
		Inputs: map[string]Inputs{
			"name": Inputs{
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"range_id": Inputs{
				Html:  "hidden",
				Value: range_id,
			},
			"hide": Inputs{
				Html:    "checkbox",
				Checked: false,
			},
			"lock": Inputs{
				Html:    "checkbox",
				Checked: false,
			},
//			"aggs": Inputs{
//				Html:        "select",
//				MultiSelect: true,
//			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func eventSettings(w http.ResponseWriter, r *http.Request) {
	event_id := get_id_from_url(r, URL_eventSettings)
	templator(TEMPLATE_ADMIN, "event-settings", eventSettings_Data(event_id), w)
}
func eventSettings_Data(event_id string) map[string]interface{} {
	event, _ := getEvent(event_id)
	var event_ranges []Option
	for range_id, item := range event.Ranges {
		if item.Aggregate == "" {
			event_ranges = append(event_ranges, Option{Value:range_id, Display:item.Name})
		} else {
			var list_of_ranges = []Option{}
			agg_list := strings.Split(item.Aggregate, ",")
			for agg_id, agg := range event.Ranges {
				if agg.Aggregate == "" {
					ok := stringInSlice(agg_id, agg_list)
					list_of_ranges = append(list_of_ranges, Option{
						Value:    agg_id,
						Display:  agg.Name,
						Selected: ok,
					})
				}
			}
			var tmp = event.Ranges[range_id]
			select_options := draw_options(Inputs{Options:list_of_ranges},"")
			tmp.Aggregate = fmt.Sprintf("<select name=aggs form=range%v multiple size=%v>%v</select>", range_id, len(list_of_ranges), select_options)
			event.Ranges[range_id] = tmp
		}
	}
	var add_agg string
	if len(event.Ranges) >= 2 {
		add_agg = generateForm2(eventSettings_add_aggForm(event_id, event_ranges))
	}
	return map[string]interface{}{
		"Title":          "Event Settings",
		"EventName":      event.Name,
		"Id":             event_id,
		"AddRange":       generateForm2(eventSettings_add_rangeForm(event_id)),
		"AddAgg":         add_agg,
		"ListRanges":     event.Ranges,
		"ListGrades":     CLASSES,
		"isPrizemeeting":	generateForm2(eventSettings_isPrizeMeet(event_id, event.IsPrizeMeet)),
//		"AddDate":        generateForm2(eventSettings_add_dateForm(event_id, event.Date, event.Time)),
		"menu":           event_menu(event_id, event.Ranges, URL_eventSettings),
		"EventGrades":    generateForm2(eventSettings_class_grades(event)),
//		"ChangeName":     generateForm2(eventSettings_event_name(event.Name, event_id)),
		"AllEventGrades": DEFAULT_CLASS_SETTINGS,
		"SortScoreboard": generateForm2(eventSettings_sort_scoreboard(event_id, event.SortScoreboard, event.Ranges)),
		"FormNewEvent": 	generateForm2(home_form_new_event(getClubs(), event.Name,event.Club,event.Date,event.Time)),
	}
}
func eventSettings_add_rangeForm(event_id string) Form {
	return Form{
		Action: URL_eventRangeInsert,
		Title:  "Add Range",
		Inputs: map[string]Inputs{
			"name": Inputs{
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func updateSortScoreBoard(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_sort_scoreboard("", "", make(map[string]Range)).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
	event_update_sort_scoreboard(event_id, validated_values["sort"])
}

func eventSettings_sort_scoreboard(event_id, existing_sort string, ranges map[string]Range) Form {
	var sort_by_ranges []Option
	var sort_by bool
	for index, Range := range ranges {
		sort_by = false
		if index == existing_sort {
			sort_by = true
		}
		sort_by_ranges = append(sort_by_ranges, Option{Display: Range.Name, Value: index, Selected: sort_by})
	}
	//	export(sort_by_ranges)
	return Form{
		Action: URL_updateSortScoreBoard,
		Title:  "Sort Scoreboard",
		Inputs: map[string]Inputs{
			"sort": Inputs{
				Html:           "select",
				Label:          "Sort Scoreboard by Range",
				Required:       true,
				Options: sort_by_ranges,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

func dateUpdate(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_dateForm("", "", "").Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
	event_update_date(event_id, validated_values["date"], validated_values["time"])
}
func eventSettings_add_dateForm(event_id, date, hour_minute string) Form {
	if date == "" {
		date = time.Now().Format("2006-02-01")
	}
	if hour_minute == "" {
		hour_minute = time.Now().Format("15:04")
	}
	return Form{
		Action: URL_dateUpdate,
		Title:  "Date &amp; Time",
		Inputs: map[string]Inputs{
			"date": Inputs{
				Html:     "date",
				Label:    "Date",
				Required: true,
				Value:    date,
			},
			"time": Inputs{
				Html:  "time",
				Label: "Time",
				Value: hour_minute,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save Date",
			},
		},
	}
}
func eventSettings_add_aggForm(event_id string, event_ranges []Option) Form {
	return Form{
		Action: URL_eventAggInsert,
		Title:  "Add Aggregate Range",
		Inputs: map[string]Inputs{
			"name": Inputs{
				Html:     "text",
				Label:    "Aggregate Name",
				Required: true,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"agg": Inputs{
				Html:         "select",
				MultiSelect:  true,
				Options: event_ranges,
				Label:        "Sum up ranges",
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Create Aggregate",
			},
		},
	}
}

func updateEventGrades(w http.ResponseWriter, r *http.Request) {
	var event Event
	validated_values := check_form(eventSettings_class_grades(event).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
	event_upsert_data(event_id, map[string]interface{}{schemaGRADES: validated_values["grades"]})
}

func slice_to_map_bool(input []string) map[string]bool {
	output := make(map[string]bool)
	for _, value := range input {
		output[value] = true
	}
	return output
}

func eventSettings_class_grades(event Event) Form {
	var grades []Option
	selected := false
	var grade_list map[string]bool
	selected_grades := strings.Split(event.Grades, ",")
	no_grades_selected := event.Grades == ""
	if !no_grades_selected {
		grade_list = slice_to_map_bool(selected_grades)
	}

	for _, class_settings := range DEFAULT_CLASS_SETTINGS {
		for _, grade_id := range class_settings.Grades {
			selected = false
			if grade_list[grade_id] || no_grades_selected {
				selected = true
			}
			grades = append(grades, Option{
				Value:    grade_id,
				Display:  CLASSES[grade_id],
				Selected: selected,
			})
		}
	}
	var event_id string
	if event.Id != "" {
		event_id = event.Id
	}
	return Form{
		Action: URL_updateEventGrades,
		Title:  "Classes &amp; Grades",
		Inputs: map[string]Inputs{
			"grades": Inputs{
				Html:           "select",
				Label:          "select Classes &amp; Grades in this event",
				MultiSelect:    true,
				Options: grades,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

func eventSettings_event_name(event_name, event_id string) Form {
	return Form{
		Action: URL_updateEventName,
		Title:  "Event name",
		Inputs: map[string]Inputs{
			"name": Inputs{
				Html:        "text",
				Label:       "Change event name",
				Value:       event_name,
				Placeholder: event_name,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

//func totalScores_update(event_id, shooter_id, range_id string)  Form {
//	return Form{
//		Action: URL_updateTotalScores,
//		Inputs: map[string]Inputs{
//			schemaTOTAL:Inputs{
//				Html:      "number",
//				Label:   "Total",
//				Required: true,
//				Min: 0,
//				Max: 60,
//			},
//			schemaCENTER:Inputs{
//				Html:      "number",
//				Label:   "Centers",
//				Required: true,
//				Min: 0,
//				Max: 60,
//			},
//			"event_id":Inputs{
//				Html: "hidden",
//				Value: event_id,
//				Required: true,
//			},
//			"shooter_id":Inputs{
//				Html: "hidden",
//				Value: shooter_id,
//				Required: true,
//			},
//			"range_id":Inputs{
//				Html: "hidden",
//				Value: range_id,
//				Required: true,
//			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Save",
//			},
//		},
//	}
//}




func eventShotsNSighters(w http.ResponseWriter, r *http.Request) {
//	var event Event



	r.ParseForm()
	form := r.Form


	fmt.Println("event_id:::", r.Form["event_id"])
	fmt.Println("shots:::", r.Form["shots"])
	fmt.Println("sight:::", r.Form["sight"])








	fmt.Println("form:")
	export(form)

	if event_id, ok := form["event_id"]; ok && len(event_id) > 0{
		fmt.Println("event_id=",event_id)
		if shots, ok := form["shots"]; ok {
			fmt.Println("shots...")
			for range_id, range_data := range shots {
				for class_id, shot_value := range range_data {
					fmt.Println("range=",range_id," class=",class_id," value=",shot_value)
				}
			}
		}else{
			fmt.Println("shots not found")
		}
	}

//	validated_values := check_form(eventSettings_class_grades(event).Inputs, r)

	dump(form)

//	event_id := validated_values["event_id"]
//	redirecter(URL_eventSettings+event_id, w, r)
//	event_upsert_data(event_id, map[string]interface{}{schemaGRADES: validated_values["grades"]})
}





func updateIsPrizeMeet(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_isPrizeMeet("", false).Inputs, r)
	event_id := validated_values["event_id"]
	prizemeet := false
	if "on" ==validated_values["prizemeet"]{
		prizemeet = true
	}
	redirecter(URL_eventSettings+event_id, w, r)
	event_upsert_data(event_id, map[string]interface{}{"p": prizemeet})
}

func eventSettings_isPrizeMeet(event_id string, checked bool) Form {
	return Form{
		Action: URL_updateIsPrizeMeet,
		Title:  "Prize Meeting Event",
		Inputs: map[string]Inputs{
			"prizemeet": Inputs{
				Html:     "checkbox",
				Label:    "Is this Event a Prize Meeting?",
				Checked:   checked,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

func shooterInsert(w http.ResponseWriter, r *http.Request) {
	var event Event
	validated_values := check_form(event_add_shooterForm(event).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_event + event_id, w, r)
	var new_shooter EventShooter
	new_shooter.FirstName = validated_values["first"]
	new_shooter.Surname = validated_values["surname"]
	new_shooter.Club = validated_values["club"]
	new_shooter.Grade = validated_values["grade"]
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	event_shooter_insert(event_id, new_shooter)
}

func event(w http.ResponseWriter, r *http.Request) {
	event_id := get_id_from_url(r, "/event/")
	templator(TEMPLATE_ADMIN, "event", event_Data(event_id), w)
}

func event_Data(event_id string) map[string]interface{} {
	event, err := getEvent(event_id)
	if err{
		return map[string]interface{}{
			"Title": "Event not found",
			"Menu":  standard_menu(ORGANISERS_MENU_ITEMS),
			"Valid": false,
		}
	}
	return map[string]interface{}{
		"Title": event.Name,
		"EventId": event_id,
		"ListRanges": event.Ranges,
//		"AddShooter": generateForm2(event_add_shooterForm(event)),
		"ListShooters": event.Shooters,
		"Menu": event_menu(event_id, event.Ranges, URL_event),
		"AddRange": generateForm2(eventSettings_add_rangeForm(event_id)),
		"ExistingShooterEntry": URL_shooterListInsert,
		"NewShooterEntry": URL_shooterInsert,
//		"GradeOptions": build_options(available_classes_grades(event)),
		"GradeOptions": draw_options(Inputs{Options:available_classes_grades(event)}, ""),
		//TODO add ClubOptions when club textbox is changed to a datalist
//		"AgeOptions": build_options(AGE_GROUPS2),
		"AgeOptions": draw_options(Inputs{Options:AGE_GROUPS2}, ""),
		"Valid": true,
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	var event Event
	validated_values := check_form(event_add_shooterListForm(event).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_event + event_id, w, r)
	var new_shooter EventShooter
	new_shooter.SID = str_to_int(validated_values["sid"])

	//	new_shooter.FirstName = validated_values["first"]
	//	new_shooter.Surname = validated_values["surname"]
	//	new_shooter.Club = validated_values["club"]
	new_shooter.Grade = validated_values["grade"]
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	temp_shooter := getShooterList(new_shooter.SID)
	new_shooter.FirstName = temp_shooter.NickName
	new_shooter.Surname = temp_shooter.Surname
	new_shooter.Club = temp_shooter.Club
	event_shooter_insert(event_id, new_shooter)
}

func event_add_shooterForm(event Event) Form {

	var event_id string
	if event.Id != ""{
		event_id = event.Id
	}

	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
			"first":Inputs{
				Html:      "text",
				Label:   "First Name",
				Required: true,
			},
			"surname":Inputs{
				Html:      "text",
				Label:   "Surname",
				Required: true,
			},
			"club":Inputs{
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
				Required: true,
			},
			"age":Inputs{
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2,
				Required: true,
			},
			"grade":Inputs{
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				Options: available_classes_grades(event),
			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Add Shooter",
//			},
			"event_id":Inputs{
				Html: "hidden",
				Value: event_id,
			},
		},
	}
}

func event_add_shooterListForm(event Event) Form {

	var event_id string
	if event.Id != ""{
		event_id = event.Id
	}

	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
//			"first":Inputs{
//				Html:      "text",
//				Label:   "First Name",
//				Required: true,
//			},
//			"surname":Inputs{
//				Html:      "text",
//				Label:   "Surname",
//				Required: true,
//			},
//			"club":Inputs{
//				Html:      "text",
//				TODO change club to a data-list
//				//SelectValues:   getClubSelectBox(eventsCollection),
//				Label:   "Club",
//				Required: true,
//			},
			"age":Inputs{
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2,
				Required: true,
			},
			"grade":Inputs{
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				Options: available_classes_grades(event),
			},
			"sid":Inputs{
				Html:      "select",
//				Label: "Class & Grade",
//				Placeholder: "Class & Grade",
				Required: true,
//				SelectedValues: available_classes_grades(event),
			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Add Shooter",
//			},
			"event_id":Inputs{
				Html: "hidden",
				Value: event_id,
			},
		},
	}
}

func available_classes_grades(event Event)[]Option{
	var grades []Option
	var grade_list map[string]bool
	selected_grades := strings.Split(event.Grades, ",")
	no_grades_selected := event.Grades == ""
	if !no_grades_selected {
		grade_list = slice_to_map_bool(selected_grades)
	}

	for _, class_settings := range DEFAULT_CLASS_SETTINGS {
		for _, grade_id := range class_settings.Grades {
			if grade_list[grade_id] || no_grades_selected {
				grades = append(grades, Option{
						Value: grade_id,
						Display: CLASSES[grade_id],
					})
			}
		}
	}
	return grades
}

type HomeCalendar struct {
	Id, Name, Club, ClubId, Time string
	Day                          string
	Date                         string
	Month                        time.Month
	Year                         int
}

func home(w http.ResponseWriter, r *http.Request) {
	//TODO change sending a string filename to the URL_page and devmode handles it automatically
	templator(TEMPLATE_HOME, "home", homeData(getEvents()), w)
}

func homeData(events []Event) map[string]interface{} {
	calendar_events := []HomeCalendar{}
	for _, event := range events {
		calendar_event := HomeCalendar{
			Id:     event.Id,
			Name:   event.Name,
			ClubId: event.Club,
			Club:   getClub(event.Club).Name,
			Time:   event.Time,
		}
		if event.Date != "" {
			//			export( event.Date)
			date_obj, err := time.Parse("2006-01-02", event.Date)
			checkErr(err)
			calendar_event.Day = date_obj.Weekday().String()
			calendar_event.Date = ordinal(date_obj.Day())
			calendar_event.Month = date_obj.Month()
			calendar_event.Year = date_obj.Year()
		}
		calendar_events = append(calendar_events, calendar_event)
	}

	//TODO change getClubs to simpler DB lookup getClubNames
	clubs := getClubs()
	return map[string]interface{}{
		"Events":   calendar_events,
		"PageName": "Calendar",
		"Menu":     home_menu("/", HOME_MENU_ITEMS),
		"FormNewEvent": generateForm2(home_form_new_event(clubs, "","","","")),
	}
}

func home_form_new_event(clubs []Club, name, club, date, eventTime string) Form {
	var action, title, save string
	if name != "" || club != "" || date != "" || eventTime != ""{
		title = "Event Details"
		save = "Update Event"
		//TODO change update to a new save function
		action = URL_eventInsert2
	}else {
		action = URL_eventInsert2
		title = "Create Event"
		save = "Save Event"
		date = time.Now().Format("2006-01-02")
		eventTime = time.Now().Format("15:04")
	}

	var clubName string

	var club_list []Option
	for _, club_data := range clubs {
		if club_data.Id == club {
			clubName = club_data.Name
		}
		club_list = append(club_list, Option{
				Value: club_data.Id,
				Display: club_data.Name,
			})
	}

	return Form{
		Action: action,
		Title:  title,
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Event Name",
				Required: true,
//				AutoComplete: "off",
				Value: name,
			},
			"club": {
				Html: "datalist",
				Label: "Host Club",
				Placeholder: "Club Name",
				Options: club_list,
				Required: true,
				AutoComplete: "off",
				Value: clubName,
			},
			"date": {
				Html:     "date",
				Label:    "Date",
				Required: true,
				Value:    date,
			},
			"time": {
				Html:  "time",
				Label: "Time",
				Value: eventTime,
			},
			"submit": {
				Html:  "submit",
				Value: save,
			},
		},
	}
}


func eventInsert2(w http.ResponseWriter, r *http.Request) {
	var clubs []Club
	validated_values := check_form(home_form_new_event(clubs,"","","","").Inputs, r)

	var newEvent Event
	newEvent.Name = validated_values["name"]

	club_name := validated_values["club"]
	club, ok := getClub_by_name(club_name)
	if ok {
		newEvent.Club = club.Id
	}else{
		newEvent.Club = insert_new_club(club_name)
	}

	if validated_values["date"] != ""{
		newEvent.Date = validated_values["date"]
	}

	if validated_values["time"] != ""{
		newEvent.Time = validated_values["time"]
	}

	//Add default ranges and aggregate ranges
	newEvent = default_event_settings(newEvent)

	newEvent.Id = getNextId(TBLevent)
	InsertDoc(TBLevent, newEvent)

	//redirect user to event settings
	redirecter(URL_eventSettings+newEvent.Id, w, r)
}

func licence(w http.ResponseWriter, r *http.Request) {
	test := map[string]interface{}{
		"Menu":     home_menu("/licence",HOME_MENU_ITEMS),
	}
	templator(TEMPLATE_HOME, "licence", test, w)
}

func licence_summary(w http.ResponseWriter, r *http.Request) {
	test := map[string]interface{}{
		"Menu":     home_menu("/licence-summary",HOME_MENU_ITEMS),
	}
	templator(TEMPLATE_HOME, "licence-summary", test, w)
}

func organisers(w http.ResponseWriter, r *http.Request) {
	templator(TEMPLATE_ADMIN, "organisers", organisers_Data(), w)
}

func organisers_Data() map[string]interface{} {
	clubs := getClubs()
	return map[string]interface{}{
		"Title":        "Organisers",
		"Events":       generateForm2(organisers_eventForm(clubs)),
		"EventList":    eventList(),
		"Clubs":        generateForm2(organisers_clubForm()),
		"ClubList":     clubs,
		"Championship": generateForm2(organisers_champForm()),
		"Menu":         standard_menu(ORGANISERS_MENU_ITEMS),
		"ShooterList":  generateForm2(organisers_update_shooter_list("")),
	}
}

func organisers_clubForm() Form {
	//TODO add validation to
	return Form{
		Action: URL_clubInsert,
		Title:  "Create Club",
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Club Name",
				Required: true,
			},
			"submit": {
				Html:  "submit",
				Value: "Add Club",
			},
		},
	}
}

func organisers_eventForm(clubs []Club) Form {
	club_name := "club"
	club := Inputs{
		Label:    "Host Club",
		Required: true,
	}
	if len(clubs) > 0 {
		club.Html = "select"
		club.Options = getClubSelectionBox(clubs)
//		club.Placeholder = "Select Club
//		club.Html = "select"
//		club.SelectValues = getClubSelectionBox(clubs)
		if len(clubs) > 1 {
			club.Placeholder = "Select Club"
		}
	} else {
		club_name = "club_insert"
		club.Html = "text"
		club.Label = "Host Club"
		club.Placeholder = "Club Name"
	}
	return Form{
		Action: URL_eventInsert,
		Title:  "Create Event",
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Event Name",
				Required: true,
			},

			club_name: club,

			"submit": {
				Html:  "submit",
				Value: "Add Event",
			},
		},
	}
}

func organisers_update_shooter_list(last_updated string) Form {
	if last_updated == "" {
		last_updated = "Never"
	}
	return Form{
		Action: URL_updateShooterList,
		Title:  "Update Shooter List",
		Inputs: map[string]Inputs{
			"submit": {
				Html:  "submit",
				Label: "Last updated: " + last_updated,
				Value: "Update",
			},
		},
	}
}

func clubInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(organisers_clubForm().Inputs, r)
	insert_new_club(validated_values["name"])
}

func insert_new_club(club_name string) string {
	var newClub Club
	newClub.Name = club_name
	newClub.Id = getNextId(TBLclub)
	newClub.AutoInc.Mound = 1
	InsertDoc(TBLclub, newClub)
	return newClub.Id
}

func eventInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(organisers_eventForm(getClubs()).Inputs, r)
	//	export(validated_values)
	var newEvent Event
	newEvent.Name = validated_values["name"]
	if club_name, ok := validated_values["club_insert"]; ok {
		insert_new_club(club_name)
	} else if club_name, ok := validated_values["club"]; ok {
		newEvent.Club = club_name
	}
	newEvent = default_event_settings(newEvent)
	newEvent.Id = getNextId(TBLevent)
	InsertDoc(TBLevent, newEvent)
}

func default_event_settings(event Event) Event {
	//TODO add club settings for default ranges and aggs to create
	event.Ranges = map[string]Range{
		"1": Range{Name: "New Range 1"},
		"2": Range{Name: "New Range 2"},
		"3": Range{Name: "New Aggregate 1", Aggregate: "1,2"},
	}
	event.SortScoreboard = "3"
	event.AutoInc.Range = 4
	return event
}

func getClubSelectionBox(club_list []Club) []Option {
	var drop_down []Option
	for _, club := range club_list {
		drop_down = append(drop_down, Option{Display: club.Name, Value:club.Id})
	}
	return drop_down
}

func eventList() []Club {
	events := getEvents()
	event_list := []Club{}
	for _, row := range events {
		event_list = append(event_list, Club{
			Name: row.Name,
			Url:  "/event/" + row.Id,
		})
	}
	return event_list
}

func organisers_champForm() Form {
	return Form{
		Action: URL_champInsert,
		Title:  "Create Championship",
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Championship Name",
				Required: true,
			},
			"submit": {
				Html:  "submit",
				Value: "Add Championship",
			},
		},
	}
}

func scoreboard(w http.ResponseWriter, r *http.Request) {
	event_id := get_id_from_url(r, URL_scoreboard)
	templator(TEMPLATE_EMPTY, "scoreboard", scoreboard_Data(event_id), w)
}
func scoreboard_Data(url string) map[string]interface{} {
	//	export(url)
	arr := strings.Split(url, "/")
	event_id := arr[0]

	event, _ := getEvent(event_id)
	var sortByRange string
	if event.SortScoreboard != "" {
		sortByRange = event.SortScoreboard
	} else if len(event.Ranges) >= 1 {
		for event_range := range event.Ranges {
			sortByRange = event_range
			break
		}
	}

	score_board_legend_on_off := make(map[string]bool)
	for _, legend_name := range SCOREBOARD_LEGEND {
		score_board_legend_on_off[legend_name] = false
	}

	// Closures that order the Change structure.
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	total := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].Total > c2.Scores[sortByRange].Total
	}
	centa := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].Centers > c2.Scores[sortByRange].Centers
	}
	cb := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].CountBack1 > c2.Scores[sortByRange].CountBack1
	}

	var shooter_list []EventShooter
	for shooter_id, shooterList := range event.Shooters {
		shooterList.Id = shooter_id
		for range_id, score := range shooterList.Scores {
			//			vardump(score)
			//			export(score)
//			score.Position = 0
			shooterList.Scores[range_id] = score
			//			dump("\n\n\n")
		}
		shooter_list = append(shooter_list, shooterList)
		//		vardump(shooterList)
	}
	if sortByRange != "" {
		OrderedBy(grade, total, centa, cb).Sort(shooter_list)
	}

	previous_grade := ""
	previous_class := ""
	position := 0
	should_be_position := 0
	shoot_off := false
	shoot_equ := false
	shooter_length := len(shooter_list)
	for index, shooter := range shooter_list {
		should_be_position += 1
		if shooter.Grade != previous_grade {
			//reset position back to 1st
			position = 1
			should_be_position = 1
			shooter_list[index].GradeSeparator = true
			previous_grade = shooter.Grade
			if class_translation(shooter.Grade) != previous_class {
				previous_class = class_translation(shooter.Grade)
				shooter_list[index].ClassSeparator = true
			}
		} else if !shoot_off && !shoot_equ {
			position = should_be_position
		}
		var display string
		if shoot_off {
			score_board_legend_on_off["ShootOff"] = true
			display = "="
			shoot_off = false
			shoot_equ = false
			shooter_list[index].Warning = 1
		}
		if shoot_equ {
			display = "="
			shoot_equ = false
		}

		this_shooter_score := shooter.Scores[sortByRange]
		if SCOREBOARD_SHOW_WARNING_FOR_ZERO_SCORES && this_shooter_score.Total == 0 && this_shooter_score.Centers == 0 {
			score_board_legend_on_off["NoScore"] = true
			shooter_list[index].Warning = 2
			if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
				position = 0
			}
		}
		if this_shooter_score.Centers == 10 && ((this_shooter_score.Total == 60 && class_translation(shooter.Grade) == "F Class") || (this_shooter_score.Total == 50 && class_translation(shooter.Grade) == "Target")) {
			shooter_list[index].Warning = 4
			score_board_legend_on_off["HighestPossibleScore"] = true
		}
		if index+1 < shooter_length {
			next_shooter := shooter_list[index+1]
			next_shooter_score := next_shooter.Scores[sortByRange]
			if shooter.Grade == next_shooter.Grade &&
				this_shooter_score.Total == next_shooter_score.Total &&
				this_shooter_score.Centers == next_shooter_score.Centers &&
				this_shooter_score.CountBack1 == next_shooter_score.CountBack1 {
				display = "="
				if this_shooter_score.Total == 0 {
					shoot_equ = true
					if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
						position = 0
					}
				} else {
					shoot_off = true
					shooter_list[index].Warning = 1
					score_board_legend_on_off["ShootOff"] = true
				}
			}
		}
		if position > 0 {
			shooter_list[index].Position = fmt.Sprintf("%v%v", display, ordinal(position))
		}
	}

	return map[string]interface{}{
		"Title":        "Scoreboard",
		"EventId":      arr[0],
		"EventName":    event.Name,
		"ListShooters": shooter_list,
		"ListRanges":   event.Ranges,
		"Css":          "scoreboard.css",
		"Legend":       render_legend(score_board_legend_on_off),
		"SortByRange":  event.SortScoreboard,
		"menu":         event_menu(event_id, event.Ranges, URL_scoreboard),
	}
}

func render_legend(items_status map[string]bool) string {
	labels := []string{}
	for _, legend_name := range SCOREBOARD_LEGEND {
		if items_status[legend_name] {
			labels = append(labels, fmt.Sprintf("<label class=%v>%v</label>", SCOREBOARD_LEGEND_CSS_CLASSES[legend_name][0], SCOREBOARD_LEGEND_CSS_CLASSES[legend_name][1]))
		}
	}
	return strings.Join(labels, " ")
}

type lessFunc func(p1, p2 *EventShooter) bool

type multiSorter struct {
	changes []EventShooter
	less    []lessFunc
}

func (ms *multiSorter) Sort(changes []EventShooter) {
	ms.changes = changes
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}

func updateShooterList(w http.ResponseWriter, r *http.Request) {
	go updateShooterList2()
}
func updateShooterList2()int{
	/* TODO:
		ckeck if there is another page

		get a shooters grades
		translate a shooters grades
		save a shooters grades
	*/

	url := "http://www.nraa.com.au/nraa-shooter-list/?_p="
//	max_pages := 514
	max_pages := 0
	for page_count := 1; page_count <= max_pages; page_count += 1 {

//		fmt.Printf("page: %v\n", page_count)
		response, err := http.Get(fmt.Sprintf("%v%v", url, page_count))
		defer response.Body.Close()
		if err != nil {
			//TODO change to the error framework with a helpfull error message
			fmt.Println("ERROR: http.Get", err)
			return 0
		}

		doc, err := html.Parse(response.Body)
		if err != nil {
			vardump(err)
		}

		var i int = 0
		var trim_space string
		var shooter NRAA_Shooter

		var find_cells func(*html.Node)
		find_cells = func(n *html.Node) {
			if n.Type == html.TextNode {
				trim_space = strings.TrimSpace(n.Data)
				if trim_space != "" {
					if i >= 1 {
						switch {
						case i == 1:
							shooter.SID, _ = strconv.Atoi(trim_space)
							break
						case i == 2:
							shooter.Surname = trim_space
							break
						case i == 3:
							shooter.FirstName = trim_space
							break
						case i == 4:
							shooter.NickName = trim_space
							break
						case i == 5:
							shooter.Club = trim_space
							break
						}
					}
					i += 1
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				find_cells(c)
			}
		}

		var find_rows func(*html.Node)
		find_rows = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "tr" {
				for _, a := range n.Attr {
					if a.Key == "data-shooter-id" && a.Val != "" {
						i = 0
						id, _ := strconv.Atoi(a.Val)
						shooter = NRAA_Shooter{NRAA_Id: id}
						find_cells(n)
						nraa_upsert_shooter(shooter)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				find_rows(c)
			}
		}
		find_rows(doc)
	}
	fmt.Println("Finished copying from website.")

	return copyNewEntries()
}

func copyNewEntries()int{
	counter := 0
	shooter_list := getShooterLists()
	for _, n_shooter := range shooter_list{
		shooter := getShooterList(n_shooter.SID)
		if shooter.SID!=0&&shooter.NRAA_Id!=0&&shooter.Surname!=""&&shooter.FirstName!=""&&shooter.NickName!=""&&shooter.Club!=""&&shooter.Address!=""&&shooter.Email!=""{
			Upsert_Doc("shooter", n_shooter.SID, n_shooter)
			counter += 1
		}
	}
	fmt.Println("Finished inserting new shooters.")
	return counter
}

type Fdsa struct{
	Surname string	`json:"surname"`
	First string	`json:"first"`
	Club string		`json:"club"`
}

func queryShooterList(w http.ResponseWriter, r *http.Request) {
	var t Fdsa
	err := json.NewDecoder(r.Body).Decode(&t)
	checkErr(err)

	query := make(map[string]interface{}, 0)
	if t.Surname != "" {
		query["s"] = map[string]interface{}{"$regex": fmt.Sprintf(`^%v`, t.Surname), "$options": "i"}
	}
	if t.First != ""{
		query["f"] = map[string]interface{}{"$regex": fmt.Sprintf(`^%v`, t.First), "$options": "i"}
	}
	if t.Club != ""{
		query["c"] = map[string]interface{}{"$regex": fmt.Sprintf(`^%v`, t.Club), "$options": "i"}
	}

	//TODO research would changing map[string]interface{} to M{} have any affect on the size of the compiled application?
//	type M map[string]interface{}
//	type S []M

	//Ignore Deleted shooters. Selects not modified, updated & merged shooters
//	query["$or"] = M{"$or": S{{"t": nil}, {"t": 0 }}}
	query["$or"] = []map[string]interface{}{{"t": nil}, {"t": map[string]interface{}{"$lt": 3 } }}
//	query["t"] = map[string]interface{}{"$or": []map[string]interface{}{{"t": nil}, {"t", map[string]interface{}{"$lt": 3 }}}}
//	query["t"] = map[string]interface{}{"$or": []map[string]interface{}{{"_id": "id1"}, {"_id": "id2"}}}



//	db.S.find({"s": {"$regex": /^Webb$/}      ,"$or": [ {"t": {"$lt": 3}}, {"t": null}]        })


//	v := M{"$or": S{{"_id": "id1"}, {"_id": "id2"}}}




//	tree := []map[string]interface{}{}

	//	type M map[string]interface{}
	//	type S []M

//	v := map[string]interface{}{"$or": []map[string]interface{}{{"_id": "id1"}, {"_id": "id2"}}}
//	v := map[string]interface{}{"$or": []map[string]interface{}{{"t": nil}, {"t", {"$lt": 3 }}}}
//
//{ age: { $lt: 25 } }


	found := searchShooters(query)

	var option_list []Option
	for _, shooter := range found{
		option_list = append(option_list, Option{
				Value: fmt.Sprintf("%v", shooter.SID),
				Display: fmt.Sprintf("%v %v, ~~ %v", shooter.FirstName, shooter.Surname, shooter.Club),
			})
	}
	fmt.Fprint(w, draw_options(Inputs{Options:option_list}, ""))
}
func event_query_shooterForm() Form {
	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
			"first":Inputs{
				Html:      "text",
				Label:   "First Name",
//				Required: true,
			},
			"surname":Inputs{
				Html:      "text",
				Label:   "Surname",
//				Required: true,
			},
			"club":Inputs{
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
//				Required: true,
			},
		},
	}
}

func event_add_existing_shooterForm() Form {
	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
			"sid":Inputs{
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
				Required: true,
			},
			"age":Inputs{
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2,
				Required: true,
			},
			"grade":Inputs{
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				//				SelectedValues: available_classes_grades(event),
			},
			"event_id":Inputs{
				Html: "hidden",
//				Value: event_id,
				Value: "",
			},
		},
	}
}

func startShooting(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_startShooting)
	templator(TEMPLATE_ADMIN, "start-shooting", startShooting_Data(data, false), w)
}

func startShootingAll(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_startShootingAll)
	templator(TEMPLATE_ADMIN, "start-shooting", startShooting_Data(data, true), w)
}
func startShooting_Data(data string, showAll bool) map[string]interface{} {
	arr := strings.Split(data, "/")
	event_id := arr[0]
	range_id := arr[1]
	event, _ := getEvent(event_id)

	available_class_shots := map[string][]string{
		"F Class": []string{
			"S1", "S2", "S3", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
		},
		"Target": []string{
			"S1", "S2", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		},
		"Match": []string{
			"S1", "S2", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		},
	}

	class_shots := map[string][]string{}
	class_shots_length := map[string]int{}
	var long_shots []string
	var temp_grade string
	var shooter_list []EventShooter
	for shooter_id, shooter_data := range event.Shooters {

		if showAll || (!showAll && ((event.IsPrizeMeet && len(shooter_data.Scores[range_id].Shots) <= 0)||(!event.IsPrizeMeet && shooter_data.Scores[range_id].Total <= 0))) {

			temp_grade = CLASS[shooter_data.Grade]
			class_shots[temp_grade] = available_class_shots[temp_grade]


			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club Inc.", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club Inc", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club.", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club", "", -1)
			shooter_data.Id = shooter_id
			shooter_list = append(shooter_list, shooter_data)

		}


	}
	for temp_grade, shots_array := range class_shots {
		class_shots_length[temp_grade] = len(shots_array)
		if len(long_shots) < len(shots_array) {
			long_shots = shots_array
		}
	}
	first_class := ""
	for _, shooter := range event.Shooters {
		first_class = CLASS[shooter.Grade]
		break
	}

	//Sort the list of shooters by grade only
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	name := func(c1, c2 *EventShooter) bool {
		return c1.FirstName < c2.FirstName
	}
	OrderedBy(grade, name).Sort(shooter_list, )


	var totalScores_link string
	if showAll{
		totalScores_link = fmt.Sprintf("<a href=%v/%v/%v>View Incompleted Shooters</a>", URL_startShooting, event_id, range_id)
	}else{
		totalScores_link = fmt.Sprintf("<a href=%v/%v/%v>View All Shooters</a>", URL_startShootingAll, event_id, range_id)
	}

	return map[string]interface{}{
		"Title":              "Start Shooting",
		"EventId":            event_id,
		"LinkToPage": 				totalScores_link,
		"RangeName":          event.Ranges[range_id].Name,
		"class_shots":        class_shots,
		"menu":               event_menu(event_id, event.Ranges, URL_startShooting),
		"RangeId":            range_id,
		"first_class":        first_class,
		"longest_shots":      long_shots,
		"class_shots_length": class_shots_length,
		//		"ListShooters": event.Shooters,
		"ListShooters": shooter_list,
		"Css":          "admin.css",
		"Js":           "start-shooting.js",
	}
}

func updateShotScores(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(startShooting_Form("", "", "", "").Inputs, r) //total_scores_Form(event_id, range_id, shooter_id, shots
	event_id := validated_values["event_id"]
	event, _ := getEvent(event_id)
	range_id := validated_values["range_id"]
	if !event.Ranges[range_id].Locked {
		shooter_id := validated_values["shooter_id"]
		shots := validated_values["shots"]
		//TODO change this to use grade to Int
		new_score := calc_total_centers(shots, 0)
		dump(new_score)
		if new_score.Centers > 0 {
			generator(w, fmt.Sprintf("%v.%v", new_score.Total, new_score.Centers), make(map[string]interface{}))
		} else {
			generator(w, fmt.Sprintf("%v", new_score.Total), make(map[string]interface{}))
		}
		event_total_score_update(event_id, range_id, shooter_id, new_score)
	} else {
		fmt.Println("BAD updateShotScores. Current Range is locked!")
	}
}

func calc_total_centers(shots string, class int) Score {
	//This function assumes all validation on input "shots" has at least been done!
	//AND input "shots" is verified to contain all characters in settings[class].validShots!
	//TODO need validation to check that the shots given match the required vailation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot

	total := 0
	centers := 0
	//	xs := 0
	countback1 := ""
	//	countback2 := ""

	relevant_settings := DEFAULT_CLASS_SETTINGS[class].ValidShots
	for _, shot := range strings.Split(shots[DEFAULT_CLASS_SETTINGS[class].SightersQty:], "") {
		total += relevant_settings[shot].Total
		centers += relevant_settings[shot].Centers
		countback1 = relevant_settings[shot].CountBack1 + countback1
		//		countback2 = relevant_settings[shot].CountBack2 + countback2
	}
	return Score{Total: total, Centers: centers, Shots: shots /*Xs: xs,*/, CountBack1: countback1 /*CountBack2: countback2*/}
}

func startShooting_Form(event_id, range_id, shooter_id, shots string) Form {
	return Form{
		Action: URL_updateTotalScores,
		Inputs: map[string]Inputs{
			"shots": Inputs{
				Html:     "number",
				Label:    "Total",
				Required: true,
				Value:    shots,
				//TODO add min and max for validation on fclass and target
			},
			"shooter_id": Inputs{
				Html:  "hidden",
				Value: shooter_id,
			},
			"range_id": Inputs{
				Html:  "hidden",
				Value: range_id,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}
func totalScores(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_totalScores)
	templator(TEMPLATE_ADMIN, "total-scores", totalScores_Data(data, false), w)
}
func totalScoresAll(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_totalScoresAll)
	templator(TEMPLATE_ADMIN, "total-scores", totalScores_Data(data, true), w)
}
func totalScores_Data(data string, show_all bool) map[string]interface{} {
	arr := strings.Split(data, "/")
	event_id := arr[0]
	range_id := arr[1]
	event, _ := getEvent(event_id)
	selected_range :=  event.Ranges[range_id]

	var totalScores_link string
	if show_all{
		totalScores_link = fmt.Sprintf("<a href=%v%v>View Incompleted Shooters</a>", URL_totalScores, data)
	}else{
		totalScores_link = fmt.Sprintf("<a href=%v%v>View All Shooters</a>", URL_totalScoresAll, data)
	}

	if len(selected_range.Aggregate) > 0{
		return map[string]interface{}{
			"Title": "Total Scores",
			"LinkToPage": totalScores_link,
			"EventId": event_id,
			"RangeName": selected_range.Name,
			"Message": ERROR_ENTER_SCORES_IN_AGG,
			"menu": event_menu(event_id, event.Ranges, URL_totalScores),
		}
	}

	//Sort the list of shooters by grade only
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	name := func(c1, c2 *EventShooter) bool {
		return c1.FirstName < c2.FirstName
	}

	var shooter_list []EventShooter
	shooters_forms := make(map[string]string)
	for shooter_id,shooter_data := range event.Shooters{
		var score string
		if shooter_data.Scores[range_id].Total > 0 {
			score = fmt.Sprintf("%v", shooter_data.Scores[range_id].Total)
		}
		if shooter_data.Scores[range_id].Centers > 0 {
			score += fmt.Sprintf(".%v", shooter_data.Scores[range_id].Centers)
		}
		if show_all || (!show_all && score == "") {
			shooters_forms[shooter_id] = generateForm2(total_scores_Form(event_id, range_id, shooter_id, score))
			shooter_data.Id = shooter_id
			shooter_list = append(shooter_list, shooter_data)
		}
	}

	OrderedBy(grade, name).Sort(shooter_list)

	return map[string]interface{}{
		"Title": "Total Scores",
		"LinkToPage": totalScores_link,
		"EventId": event_id,
		"RangeName": selected_range.Name,
		"RangeId": range_id,
		"ListRanges": event.Ranges,

//		"ListShooters": event.Shooters,
		"ListShooters": shooter_list,
		"menu": event_menu(event_id, event.Ranges,URL_totalScores),
		"FormTotalScores": shooters_forms,
		"Js": "total-scores.js",
	}
}

//func updateTotalScores(w http.ResponseWriter, r *http.Request){
//	validated_values := check_form(total_scores_Form("", "", "", "").Inputs, r) //total_scores_Form(event_id, range_id, shooter_id, total, centers
//	event_id := validated_values["event_id"]
//	range_id := validated_values["range_id"]
//	shooter_id := validated_values["shooter_id"]
////	redirecter(URL_totalScores + event_id + "/" + range_id, w, r)
//	var new_score Score
//	score := strings.Split(validated_values["score"], ".")
//	if total, ok := score[0]
//	new_score.Total = str_to_int(validated_values["total"])
//	new_score.Centers = str_to_int(validated_values["centers"])
//	event_total_score_update(event_id, range_id, shooter_id, new_score)
//
//	if new_score.Centers > 0 {
//		generator(w, fmt.Sprintf("%v.%v", new_score.Total, new_score.Centers), make(map[string]interface{}))
//	}else {
//		generator(w, fmt.Sprintf("%v", new_score.Total), make(map[string]interface{}))
//	}
//}

func updateTotalScores(w http.ResponseWriter, r *http.Request){
	validated_values := check_form(total_scores_Form("", "", "", "").Inputs, r) //total_scores_Form(event_id, range_id, shooter_id, total, centers
	event_id := validated_values["event_id"]
	range_id := validated_values["range_id"]
	shooter_id := validated_values["shooter_id"]
	//	redirecter(URL_totalScores + event_id + "/" + range_id, w, r)
	score := strings.Split(validated_values["score"], ".")
	//return_string is used for ajax returning the send value
//	return_string := "0"
	total := str_to_int(score[0])
	if total > 0{
		new_score := Score{Total: total}
		if len(score) > 1 && score[1] != "" && str_to_int(score[1]) > 0{
			centers := str_to_int(score[1])
			new_score.Centers = centers
//			return_string = fmt.Sprintf("%v.%v", new_score.Total, new_score.Centers)
//		}else{
//			return_string = fmt.Sprintf("%v", new_score.Total)
		}
		go event_total_score_update(event_id, range_id, shooter_id, new_score)
	}
//	generator(w, return_string, make(map[string]interface{}))
	redirecter(URL_totalScores+event_id+"/"+range_id, w, r)
}



func search_for_aggs(event_id, range_id string)[]string{
	var aggs_to_calculate []string
	event, _ := getEvent(event_id)
	for agg_id, range_data := range event.Ranges{
		if len(range_data.Aggregate) > 0{
			for _, this_range_id := range range_data.Aggregate{
				if string(this_range_id) == range_id{
					aggs_to_calculate = append(aggs_to_calculate, agg_id)
				}
			}
		}
	}
	return aggs_to_calculate
}
func calculate_aggs(event Event, shooter_id string, ranges []string)Event{

//	if xx, ok := event.Shooters[shooter_id]; ok {
//		xx.count = 2
//		m["x"] = xx
//	} else {
//		panic("X isn't in the map")
//	}



//	if event.Shooters[shooter_id].Scores != nil{
//		dump("new val is not none")
//	}else{
	if event.Shooters[shooter_id].Scores == nil{
//		dump("shooter's scores is empty")
		temp_kkk := event.Shooters[shooter_id]
		temp_kkk.Scores = map[string]Score{}
		event.Shooters[shooter_id] = temp_kkk
	}
	for _, agg_id := range ranges {
		total := 0
		centers := 0
		count_back1 := ""
		range_id := ""
		for _, rangeId := range event.Ranges[agg_id].Aggregate {
			range_id = string(rangeId)
			total += event.Shooters[shooter_id].Scores[range_id].Total
			centers += event.Shooters[shooter_id].Scores[range_id].Centers
			count_back1 += event.Shooters[shooter_id].Scores[range_id].CountBack1
		}
		event.Shooters[shooter_id].Scores[agg_id] = Score{Total: total, Centers: centers, CountBack1: count_back1}
	}

//		event.Shooters[shooter_id].Scores = make([]Score{}, 1)
//		agg_total := event.Shooters[shooter_id].Scores[agg_id]
//		agg_total.Total = 0
//		agg_total.Centers = 0
//		agg_total.CountBack1 = ""
//		for _, rangeId := range event.Ranges[agg_id].Aggregate{
//			range_id := string(rangeId)
//			agg_total.Total += event.Shooters[shooter_id].Scores[range_id].Total
//			agg_total.Centers += event.Shooters[shooter_id].Scores[range_id].Centers
//			agg_total.CountBack1 += event.Shooters[shooter_id].Scores[range_id].CountBack1
//		}
//		event.Shooters[shooter_id].Scores[agg_id] = agg_total
//	}
	return event
}


func total_scores_Form(event_id, range_id, shooter_id, score string) Form {
	return Form{
		Action: URL_updateTotalScores,
		Inputs: map[string]Inputs{
			"score":Inputs{
				Html:      "tel",
//				Label:   "Total",
				Required: true,
				Value: score,
				//TODO add min and max for validation on fclass and taget
				Size: 4,
//				Min: 0,
//				Step: 0.01,
//				Max: 50,
				Pattern: "[0-9]{1,2}(.[0-9]{1,2}){0,1}",
			},
//			"centers":Inputs{
//				Html:      "number",
//				Label:   "Centers",
//				Required: true,
//				Value: centers,
//				Size: 4,
//				Min: 0,
//				Max: 10,
//				//TODO add html5 validation for centers based on total.
//				//TODO add min = 0, max = parseInt(  total / max(class_valid_shots) )
//			},
			"shooter_id":Inputs{
				Html: "hidden",
				Value: shooter_id,
			},
			"range_id":Inputs{
				Html: "hidden",
				Value: range_id,
			},
			"event_id":Inputs{
				Html: "hidden",
				Value: event_id,
			},
			"submit":Inputs{
				Html:    "submit",
				Value:   "Save",
			},
		},
	}
}

func random_data(w http.ResponseWriter, r *http.Request) {
	event_id := "L"
	range_id := "2"
	attributes := strings.Split(strings.Replace(r.RequestURI, "/random-data/", "", -1), "&")

	for _, request := range attributes {
		properties := strings.Split(request, "=")

		switch properties[0] {
		case "shooterQty":
			shooterQty := str_to_int(properties[1])
			counter := 0
			for counter <= shooterQty {
				counter += 1
				//make some requests for x number of shooters
//				resp, _ := http.PostForm("http://localhost/shooterInsert",
//				go http.PostForm("http://localhost/shooterInsert",
				http.PostForm("http://localhost/shooterInsert",
					url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
					"surname":   {randomdata.LastName()},
					"club":      {randomdata.State(randomdata.Large)},
					"age":       {"N"},
					"grade":     {"a"},
					"event_id":  {event_id},
				})
//				defer resp.Body.Close()
			}
		case "totalScores":
			event, _ := getEvent(event_id)
			for shooter_id, _ := range event.Shooters{
				// rand.Seed(time.Now().UnixNano())
				fmt.Println(shooter_id)
				rand.Seed(90)
//				resp, _ := http.PostForm("http://localhost/updateTotalScores",
//				go http.PostForm("http://localhost/updateTotalScores",
//					url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
//					"score":			{fmt.Sprintf("%v.%v",rand.Intn(51),rand.Intn(11))},
//					"shooter_id":	{shooter_id},
//					"range_id":		{range_id},
//					"event_id":		{event_id},
//				})
//				defer resp.Body.Close()
				event_total_score_update(event_id, range_id, shooter_id, Score{
						Total: rand.Intn(51),
						Centers: rand.Intn(11),
				})
			}
		}
	}
/*
	_, ok = form["start_shooting"]

	if ok {
		//fill out every shooters range
	}*/

}

/* Schema Rules:
lowercase letters MUST be used for struct properties
Uppercase letters MUST be used for a sub struct
*/
type Event struct {
	Id   string `bson:"_id"`
	Club string `bson:"c"`
	Name string `bson:"n"`
	//	Datetime string							`bson:"d,omitempty"`		No browser currently supports date time, so settling for separate fields that google chrome allows
	Date           string                  `bson:"d,omitempty"`
	Time           string                  `bson:"t,omitempty"`
	Grades         string                  `bson:"g,omitempty"`
	SortScoreboard string                  `bson:"o,omitempty"`
	IsPrizeMeet    bool                    `bson:"p,omitempty"`
	Ranges         map[string]Range        `bson:"R,omitempty"`
	Shooters       map[string]EventShooter `bson:"S,omitempty"`
	TeamCat        map[string]TeamCat      `bson:"A,omitempty"`
	Teams          map[string]Team         `bson:"T,omitempty"`
	AutoInc        AutoInc                 `bson:"U"`
}

type AutoInc struct {
	Mound   int `bson:"M,omitempty"`
	Event   int `bson:"E,omitempty"`
	Club    int `bson:"C,omitempty"`
	Range   int `bson:"R,omitempty"`
	Shooter int `bson:"S,omitempty"`
}

type Club struct {
	Id        string           `bson:"_id"`
	Name      string           `bson:"n"`
	LongName  string           `bson:"l,omitempty"`
	Mounds    map[string]Mound `bson:"M,omitempty"`
	Latitude  string           `bson:"t,omitempty"`
	Longitude string           `bson:"g,omitempty"`
	Url       string           `bson:"u,omitempty"`
	AutoInc   AutoInc          `bson:"U"`
}

type Range struct {
	Name       string `bson:"n"`
	Aggregate  string `bson:"a,omitempty"`
	ScoreBoard bool   `bson:"s,omitempty"`
	Locked     bool   `bson:"l,omitempty"`
	Hidden     bool   `bson:"h,omitempty"`
}

type Score struct {
	//TODO the schema should change so that it can use unsigned 64 bit numbers instead
	Total int `bson:"t"`
	//	Total uint64								`bson:"t"`
	Shots   string `bson:"s,omitempty"`
	Centers int    `bson:"c,omitempty"`
	//	Centers uint64								`bson:"c"`
	CountBack1 string `bson:"v,omitempty"`
	//	CountBack2 string							`bson:"x,omitempty"`
	//	Xs string									`bson:"u,omitempty"` //This might be handy for the future?
	Position int `bson:"p,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
}

type NRAA_Shooter struct {
	SID       int    `bson:"_id,omitempty"`
	NRAA_Id   int    `bson:"i,omitempty"`
	Surname   string `bson:"s,omitempty"`
	FirstName string `bson:"f,omitempty"`
	NickName  string `bson:"n,omitempty"`
	Club      string `bson:"c,omitempty"`
}

type EventShooter struct {
	FirstName string           `bson:"f"`
	Surname   string           `bson:"s"`
	Club      string           `bson:"b"` //TODO should possibly change to "C"??
	Grade     string           `bson:"g"`
	AgeGroup  string           `bson:"a,omitempty"`
	Scores    map[string]Score `bson:"omitempty,inline"` //S is not used!

	SID int `bson:"d,omitempty"`
	//SCOREBOARD
	Id       string `bson:"i,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	Position string `bson:"x,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	Warning  int8   `bson:"y,omitempty"` //DON'T SAVE THIS TO DB! used for scoreboard only.
	//		0 = nil
	//		1 = shoot off
	//		2 = warning, no score
	//		3 = incomplete
	//		4 = highest posible score

	//START-SHOOTING & TOTAL-SCORES
	GradeSeparator bool `bson:"z,omitempty"` //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
	ClassSeparator bool `bson:"o,omitempty"` //DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
	//	Id string									`bson:"w,omitempty"`//DON'T SAVE THIS TO DB! used for start-shooting and total-scores only.
}

type Shooter struct {
	SID int `bson:"_id,omitempty"`
	//	SID int					`bson:"s,omitempty"`
	NRAA_Id   int    `bson:"i,omitempty"`
	Surname   string `bson:"s,omitempty"`
	FirstName string `bson:"f,omitempty"`
	NickName  string `bson:"n,omitempty"`
	Club      string `bson:"c,omitempty"`
	Address   string `bson:"a,omitempty"`
	Email     string `bson:"e,omitempty"`
	//Shooter details 0=not modified, 1=updated, 2=merged, 3=deleted
	Status    int    `bson:"t,omitempty"`
	//If shooter details are merged with another existing shooter then this is the other NRAA_SID it was merged with
	//When merging set one record to merged, the other to deleted.
	//Both records must set MergedSID to the other corresponding shooter SID
	MergedSID int    `bson:"m,omitempty"`
}

//
//type Shooter struct{
//	Id string									`bson:"n"`
//	FirstName string							`bson:"f"`
//	Surname string								`bson:"s"`
//	Skill map[string]Skill	//Grading set by the VRA for each class
//}
//
type TeamCat struct {
	Name string `bson:"n"`
}

type Team struct {
	name     string `bson:"n"`
	teamCat  []int  `bson:"t"`
	shooters []int  `bson:"s,omitempty"`
}

type Mound struct {
	Distance int    `bson:"d"`
	Unit     string `bson:"u"`
	Name     string `bson:"n,omitempty"`
	Notes    string `bson:"o,omitempty"`
}

type Skill struct {
	Grade      string
	Percentage float64 //TODO would prefer an unsigned float here
}

const (
	VERSION = "0.90"

	//Server Settings
	PRODUCTION = false //False = output dev warnings, E.g. Template errors
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
			Grades:        []string{"d", "e", "f", "g"},
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
	}
	//TODO make these dynamic from club settings
	JSCLASS = map[string]string{
		"a": "target",
		"b": "target",
		"c": "target",
		"d": "fclass",
		"e": "fclass",
		"f": "fclass",
		"g": "fclass",
		"h": "match",
		"i": "match",
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

type ErrorMsg struct {
	Title, Message string
	Info           bool
}

var error_queue []ErrorMsg

func error_message(info bool, index, title, message string) {
	for _, queue_item := range error_queue {
		if queue_item.Title == title {
			return
		}
	}
	error_queue = append(error_queue, ErrorMsg{
		//	error_queue[index] = ErrorMsg{
		Title:   title,
		Message: message,
		Info:    info,
	})
	//	}
}
func remove_error(title string) {
	for index, queue_item := range error_queue {
		if queue_item.Title == title {
			error_queue[index] = ErrorMsg{}
		}
	}
}

func templator(main_template string, content_template string, data map[string]interface{}, w http.ResponseWriter) {
	//Ajax responses should not use this function! Instead use "generator(w, body, data)"
	source := loadHTM(main_template)
	error_list := []byte(render_errors())
	remove_chars := map[string][]byte{
		"^^DIR_JS^^":       []byte(DIR_JS),
		"^^DIR_CSS^^":      []byte(DIR_CSS),
		"^^DIR_ICON^^":     []byte(DIR_ICON),
		"^^FAVICON^^":      []byte(FAVICON),
		"^^CURRENT_YEAR^^": []byte(fmt.Sprintf("%v", time.Now().Year())),
		"^^ERROR^^":        error_list,
		"^^BODY^^":         loadHTM(content_template),
	}
	for search, replace := range remove_chars {
		source = bytes.Replace(source, []byte(search), replace, -1)
	}
	generator(w, string(source), data)
}

func render_errors() string {
	var output string
	if len(error_queue) >= 1 {
		for _, error := range error_queue {
			if error.Title != "" && error.Message != "" {
				var class = "error"
				if error.Info {
					class = "info"
				}
				output += fmt.Sprintf("<div class=%v><h2>%v:</h2>%v</div>", class, error.Title, error.Message)
			}
		}
		error_queue = []ErrorMsg{}
	}
	return output
}

func loadHTM(page_name string) []byte {
	//TODO add all html sources to []byte constant in a new go file
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name))
	checkErr(err)
	return dev_mode_loadHTM(page_name, bytes)
}

func generator(w http.ResponseWriter, fillin string, data map[string]interface{}) {
	my_html := template.New("my_template").Funcs(template.FuncMap{
		"HTM": func(x string) template.HTML {
			return template.HTML(x)
		},
		"CLASS": func(class string) string {
			return class_translation(class)
		},
		"CLASSLONG": func(class string) string {
			return class_long_translation(class)
		},
		"JSCLASS": func(class string) string {
			return js_class_translation(class)
		},
		"GRADE": func(grade string) string {
			return grade_translation(grade)
		},
		"Fieldset": func(title string) template.HTML {
			return template.HTML(field_set(title))
		},
		"EndFieldset": func() template.HTML {
			return template.HTML("</fieldset>")
		},
		"COLSPAN": func(longest_shots []string, short_shots int) template.HTMLAttr {
			if len(longest_shots) > short_shots {
				return template.HTMLAttr(fmt.Sprintf(" colspan=%v", len(longest_shots)-short_shots+1))
			}
			return template.HTMLAttr("")
		},
		"CSSclass": func(class_name1, class_name2 interface{}) template.HTMLAttr {
			if class_name1 != "" && class_name2 != "" {
				return template.HTMLAttr(fmt.Sprintf(" class=%v%v", class_name1, class_name2))
			}
			return template.HTMLAttr("")
		},
		"DisplayShot": func(shot_index int, score Score) string {
			if len(score.Shots) > shot_index {
				return fmt.Sprintf("%s", ShotsToValue[string(score.Shots[shot_index])])
			}
			return ""
		},
		"NOSHOOTERS": func() template.HTML {
			return template.HTML(ERROR_NO_SHOOTERS)
		},
		"ERROR_NO_EVENTS": func() template.HTML {
			return template.HTML(ERROR_NO_EVENTS)
		},
		"POSITION": func(position int) template.HTMLAttr {
//			return template.HTMLAttr(fmt.Sprintf(" class=p%v", position))
						if position > 0{
			//				if position <= 3 {
								return template.HTMLAttr(fmt.Sprintf(" class=p%v", position))
			//				}
						}
						return template.HTMLAttr("")
		},
//		"START_SHOOTING_SHOTS": func(score Score) template.HTML {
//			var output string
//			for _, shot := range strings.Split(score.Shots, "") {
//				output += fmt.Sprintf("<td>%v</td>", ShotsToValue[shot])
//			}
//			return template.HTML(output)
//		},
	})

	t := template.Must(my_html.Parse(fillin))
	checkErr(t.Execute(w, data))
}

type Menu struct {
	Name, Link string
	Ranges     bool
}

var EventMenuItems = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
//	Menu{
//		Name: "Organisers",
//		Link: URL_organisers,
//	},
	Menu{
		Name: "Event",
		Link: URL_event,
	},
	Menu{
		Name: "Event Settings",
		Link: URL_eventSettings,
	},
	Menu{
		Name: "Scoreboard",
		Link: URL_scoreboard,
	},
	Menu{
		Name:   "Total Scores",
		Link:   URL_totalScores,
		Ranges: true,
	},
	Menu{
		Name:   "Start Shooting",
		Link:   URL_startShooting,
		Ranges: true,
	},
//	Menu{
//		Name: "Close Menu",
//		Link: "#",
//	},
}

func event_menu(event_id string, event_ranges map[string]Range, page_url string) string {
	menu := "<ul>"
	selected := ""
	for _, menu_item := range EventMenuItems {
		if menu_item.Link == page_url{
			selected = " class=v"
		}

		if menu_item.Ranges {
			class := "m"
			if menu_item.Link == page_url{
				class = `"v m"`
			}


			//The a tag is needed for my ipad
			menu += fmt.Sprintf("<li class=%v><a href=#>%v</a><ul>", class, menu_item.Name)
			for range_id, range_item := range event_ranges {
				if len(range_item.Aggregate) == 0 && !range_item.Hidden {
					menu += fmt.Sprintf("<li><a href=%v%v/%v>%v</a></li>", menu_item.Link, event_id, range_id, range_item.Name)
				}
			}
			menu += "</ul></li>"
		} else {
			if menu_item.Link == "/"{
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, menu_item.Link, menu_item.Name)
			}else if menu_item.Link[len(menu_item.Link)-1:] == "/" {
				menu += fmt.Sprintf("<li%v><a href=%v%v>%v</a></li>", selected, menu_item.Link, event_id, menu_item.Name)
			} else {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, menu_item.Link, menu_item.Name)
			}
		}
		selected = ""
	}
	menu += "</ul>"
	return menu
}

var HOME_MENU_ITEMS = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
//	Menu{
//		Name: "Features",
//		Link: "/features",
//	},
//	Menu{
//		Name: "Events",
//		Link: "/events",
//	},
//	Menu{
//		Name: "Clubs",
//		Link: URL_clubs,
//	},
	Menu{
		Name: "Organisers",
		Link: URL_organisers,
	},
//	Menu{
//		Name: "Event Archive",
//		Link: URL_archive,
//	},
	Menu{
		Name: "About",
		Link: URL_about,
	},
}

var ORGANISERS_MENU_ITEMS = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
	Menu{
		Name: "Clubs",
		Link: URL_clubs,
	},
	Menu{
		Name: "Events",
		Link: URL_events,
	},
	Menu{
		Name: "Event Archive",
		Link: URL_archive,
	},
	Menu{
		Name: "Organisers",
		Link: URL_organisers,
	},
}

func standard_menu(menu_items []Menu) string {
	menu := "<ul>"
	for _, menu_item := range menu_items {
		menu += fmt.Sprintf("<li><a href=%v>%v</a></li>", menu_item.Link, menu_item.Name)
	}
	menu += "</ul>"
	return menu
}

func home_menu(page string, menu_items []Menu) string {
	menu := "<ul id=menu>"
	for _, menu_item := range menu_items {
		if page != menu_item.Link {
			menu += fmt.Sprintf("<li><a href=%v>%v</a></li>", menu_item.Link, menu_item.Name)
		}else{
			menu += fmt.Sprintf("<li class=v><a href=%v>%v</a></li>", menu_item.Link, menu_item.Name)
		}
	}
	menu += "</ul>"
	return menu
}

//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func dump(input interface{}) {
	fmt.Printf("%v\n", input)
}
func vardump(input interface{}) {
	fmt.Printf("%+v\n", input) //map field names included
}
func export(input interface{}) {
	fmt.Printf("\n%#v\n\n", input) //can copy and declare new variable with it. Most ouput available
}

func exists(dict map[string]interface{}, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func echo(input interface{}) string {
	return fmt.Sprintf("%v", input)
}
func str_to_int(input string) int {
	output, err := strconv.Atoi(input)
	checkErr(err)
	return output
}
func str_to_int64(input string) int64 {
	output, err := strconv.ParseInt(input, 10, 64)
	checkErr(err)
	return output
}

func addQuotes(input string) string {
	if strings.Contains(input, " ") {
		return "\"" + input + "\""
	}
	return input
}
func addQuotesEquals(input string) string {
	if input != "" {
		if strings.Contains(input, " ") {
			return "=\"" + input + "\""
		}
		return "=" + input
	}
	dump("addQuotesEquals had an empty parameter!")
	return ""
}

// Ordinal gives you the input number in a rank/ordinal format.
// Ordinal(3) -> 3rd
//author github.com/dustin/go-humanize/blob/master/ordinals.go
func ordinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(x) + suffix
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func check_form(options map[string]Inputs, r *http.Request)(map[string]string){
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return new_values
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	new_values := make(map[string]string)
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok && ((options[option].Required && array[0] != "") || !options[option].Required) {
				if len(array) > 1{
					new_values[option] = strings.Join(array, ",")
				}else{
					new_values[option] = array[0]
				}
			}else {
				fmt.Printf("options[%v] is REQUIRED OR is not in array", option)
//				warning("options[%v] is REQUIRED OR is not in array", option)
			}
		}
	}
	return new_values
}

func generateForm2(form Form) string {
	var output, attributes, element, options string
	for name, input := range form.Inputs {
		element = ""
		attributes = ""
		options = ""

		if name != "" {
			attributes += " name=" + name
		}
		if input.Value != ""{
			if input.Html != "submit" {
				attributes += " value=" + addQuotes(input.Value))
			}
		}
		if input.Required {
			attributes += " required"
		}
		if input.Placeholder != "" {
			attributes += " placeholder="+addQuotes(input.Placeholder)
		}
		if input.Min > -1 {
			attributes += fmt.Sprintf(" min=%v", input.Min)
		}
		if input.Max > -1 {
			attributes += fmt.Sprintf(" max=%v", input.Max)
		}
		if input.Step > 0 {
			attributes += fmt.Sprintf(" step=%v", input.Step)
		}
		if input.Checked {
			attributes += " checked"
		}
		if input.Size > 0 {
			attributes += " size=%d" + input.Size
		}
		if input.AutoComplete != "" {
			attributes += " autocomplete="+input.AutoComplete
		}

		if input.MultiSelect {
			attributes += " multiple"
			if len(input.Options) > 4 {
				attributes += fmt.Sprintf(" size=%d", len(input.Options))
			}
		}
		if input.Html == "datalist"{
			attributes += " id=" + name
		}
		options = draw_options(input, name)
		if input.Help != "" {
			attributes += "title=" + addQuotes(input.Help)
		}


		if input.Html == "select" {
			element += "<select"+attributes+">"+options+"</select>"
		}else if input.Html == "submit" {
			element += "<button"+attributes+">"+input.Value+options+"</button>"
		}else {
			if input.Html != "text" {
				attributes += " type="+input.Html
			}
			element += "<input"+attributes+">"
		}
		if input.Label != "" {
			output += "<label>"+input.Label+": "+element+"</label>"
		}
	}
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

type Form struct {
	Action string
	Title  string
	Inputs map[string]Inputs
}
type Inputs struct {
	//AutoComplete values can be: "off" or "on"
	Html, Label, Help, Value, Pattern, Placeholder, AutoComplete  string
	Checked, MultiSelect, Required bool
	Min, Max, Size                 int
	Options                        []Options
	Step                           float64
}
type Options struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

func draw_options(input Inputs, name string)string{
	if len(input.Options) <= 0 {
		return ""
	}
	output := ""
	if input.Placeholder != "" && input.Html != "datalist"{
		output += "<option selected value disabled>"+input.Placeholder+"</option>"
	}
	for _, option := range input.Options {
		output += "<option"
		if option.Selected {
			output += " selected"
		}
		if option.Value != "" {
			output += " value" + addQuotesEquals(option.Value)
		}
		output += ">" + option.Display + "</option>"
	}
	if input.Html == "datalist"{
		output = "<datalist id=" + name + ">"+output+"</datalist>"
	}
	return output
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
//6,493 bytes
