package main

import (
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	"log"
	"time"
	"net/http"
	"net/url"
	"strings"
	"go-randomdata-master"
	"math/rand"
	"strconv"

	"github.com/yvasiyarov/gorelic"
)
const dev_mode_DEBUG = false	//Send system metric data to NewRelic.com
var (
	agent = gorelic.NewAgent()
	//Use io.Writer >>> ioutil.Discard to disable logging any output
	Timimgs = log.New(ioutil.Discard, "TIMINGS: ", log.Ltime|log.Lshortfile)
	Trace   = log.New(os.Stdout,      "TRACE:   ", log.Ltime|log.Lshortfile)
	Info    = log.New(os.Stdout,      "INFO:    ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr,      "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error   = log.New(os.Stderr,      "ERROR:   ", log.Ldate|log.Ltime|log.Lshortfile)
)

func dump(input interface{}) {
	Trace.Printf("%v", input)
}
func vardump(input interface{}) {
	Trace.Printf("%+v", input) //map field names included
}
func export(input interface{}) {
	Trace.Printf("%#v", input) //can copy and declare new variable with it. Most ouput available
}

func dev_mode_timeTrack(start time.Time, requestURI string) {
	elapsed := time.Since(start)
	Timimgs.Printf("%s took %s", requestURI, elapsed)
}

func dev_mode_check_form(check bool, message string){
	if !check{
		Warning.Println(message)
	}
}

func dev_mode_NewRelicDebugging(){
	if dev_mode_DEBUG{
		agent.Verbose = true
		agent.CollectHTTPStat = true
		agent.NewrelicLicense = "abf730f5454a9a1e78af7a75bfe04565e9e0d3f1"
		agent.Run()
	}
}

func dev_mode_loadHTM(page_name string, existing_minified_file []byte) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name))
	if err == nil{
		bytes = dev_mode_minifyHtml(page_name, bytes)
		existing_len := len(existing_minified_file)
		new_len := len(bytes)
		if existing_len != new_len {
//			Warning.Printf("Page '%v' had %v bytes removed (%v percent), total: %v, from: %v", page_name, existing_len-new_len, (existing_len*100/new_len-100)*-1, new_len, existing_len)
			Info.Printf("Page '%v' had %v bytes removed (%v percent), total: %v, from: %v", page_name, existing_len-new_len, existing_len*100/new_len-100, new_len, existing_len)
		}
		Info.Printf("sizes==", new_len, existing_len)
		if new_len > existing_len {
			Error.Println("How did this page get bigger?")
		}
	}else{
		ioutil.WriteFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name), bytes, 0777)
	}
	ioutil.WriteFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name), bytes, 0777)
	return bytes
}

func dev_mode_minifyHtml(page_name string, minify []byte) []byte {
	if bytes.Contains(minify, []byte("ZgotmplZ")) {
		Warning.Println("Template generation error: ZgotmplZ")
		return []byte("")
	}

	remove_chars := []string{
		"	", //Tab
		"\n", //new line
		"\r", //carriage return
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
	for _, search := range remove_chars {
		minify = bytes.Replace(minify, []byte(search), []byte(""), -1)
	}
	for search, replace := range replace_chars {
		length := len(minify)
		minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
		if length != len(minify) {
			Warning.Printf("A dodgy character (%v) was found in the source! Please replace with (%v).", search, replace)
		}
	}
	//TODO why is the string not being replaced here, even though it is 100% running?
	return bytes.Replace(minify, []byte("~~~"), []byte(" "), -1)
}

func dev_mode_random_data(w http.ResponseWriter, r *http.Request) {
	eventId := ""
	rangeId := ""
	shooterQty := 0
	totalScores := false
	startShooting := false
	var properties []string
//	random_grades := []string{/*"a","b","c",*/"d", "e", "f", "g", "h", "i", "j"}
	attributes := strings.Split(strings.Replace(r.RequestURI, "/random-data/", "", -1), "&")
	for _, request := range attributes {
		properties = strings.Split(request, "=")
		if properties[0] == "event_id" || properties[0] == "eventId" {
			eventId = properties[1]
		} else if properties[0] == "range_id" || properties[0] == "rangeId" {
			rangeId = properties[1]
		} else if properties[0] == "totalScores"{
			totalScores = true
		} else if properties[0] == "startShooting"{
			totalScores = true
		} else if properties[0] == "shooterQty" {
			shooterQty, _ = strconv.Atoi(properties[1])
		}
	}
	if shooterQty > 0{
		go dev_mode_random_data_shooterQty(shooterQty, eventId)
	}
	if totalScores{
		go dev_mode_random_data_totalScores(eventId, rangeId)
	}
	if startShooting{
		go dev_mode_random_data_startShooting(eventId, rangeId)
	}
}

func dev_mode_random_data_startShooting(eventId, rangeId string){
	event, _ := getEvent(eventId)
	for shooterId, shooter := range event.Shooters {
		http.PostForm("http://localhost/updateShotScores",
			url.Values{
			"shots":			{randomShooterScores(shooter.Grade)},
			"shooter_id":	{fmt.Sprintf("%v",shooterId)},
			"range_id":		{rangeId},
			"event_id":		{eventId},
		},
		)
	}
}

func dev_mode_random_data_totalScores(eventId, rangeId string){
	event, _ := getEvent(eventId)
	for shooter_id, _ := range event.Shooters {
		rand.Seed(time.Now().UnixNano())	//Use rand.Seed(90) with a constant number to make the same number
		/*resp, _ := http.PostForm("http://localhost/updateTotalScores",
		go http.PostForm("http://localhost/updateTotalScores",
			url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
			"score":			{fmt.Sprintf("%v.%v",rand.Intn(51),rand.Intn(11))},
			"shooter_id":	{shooter_id},
			"range_id":		{range_id},
			"event_id":		{event_id},
		})
		resp.Body.Close()*/
		range_Id, err := strToInt(rangeId)
		if err == nil {
			eventTotalScoreUpdate(eventId, range_Id, []int{shooter_id}, Score{
				Total: rand.Intn(51),
				Centers: rand.Intn(11),
			})
		}
	}
}

func dev_mode_random_data_shooterQty(shooterQty int, eventId string){
//	var resp *http.Response
	counter := 0
	for counter <= shooterQty {
		//make some requests for x number of shooters
		counter += 1

		event_shooter_insert(eventId, EventShooter{
			FirstName: randomdata.FirstName(randomdata.RandomGender),
			Surname: randomdata.LastName(),
			Club: randomdata.State(randomdata.Large),
			AgeGroup: "N",
			Grade: rand.Intn(8),
		})

		/*resp, _ = http.PostForm("http://localhost/shooterInsert",
			url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
			"surname":   {randomdata.LastName()},
			"club":      {randomdata.State(randomdata.Large)},
			"age":       {"N"},
				//					"grade":     {random_grades[rand.Intn(len(random_grades)-1)]},
			"grade":     {fmt.Sprintf("%v", rand.Intn(8))},
			"event_id":  {eventId},
		})
		resp.Body.Close()*/
	}
}

func randomShooterScores(shooterGrade int)string{
	shooterClass := grades()[shooterGrade].Settings
	score := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < shooterClass.SightersQty; i++ {
		score += "-"
	}
	availableShots := len(shooterClass.Buttons)
	for i := 0; i < shooterClass.ShotsQty; i++ {
		score += string(shooterClass.Buttons[rand.Intn(availableShots)])
	}
	return score
}
