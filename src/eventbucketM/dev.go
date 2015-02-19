// +build dev

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
var (
	agent = gorelic.NewAgent()
	//Use io.Writer >>> ioutil.Discard to disable logging any output
	Trace   = log.New(os.Stdout, "TRACE:   ", log.Ltime|log.Lshortfile)
	Info    = log.New(os.Stdout, "INFO:    ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func serveDir(contentType string){
	http.Handle(contentType,
		http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			//If url is a directory return a 404 to prevent displaying a directory listing
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			httpHeaders(w, []string{"expire", "cache", contentType, "public"})
			Gzip(http.FileServer(http.Dir(DIR_ROOT)), w, r)
		})))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc{
	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer dev_mode_timeTrack(time.Now(), r.RequestURI)
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
		Gzip(h, w, r)
	}))
}

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
	Trace.Printf("%s took %s", requestURI, time.Since(start))
}

func dev_mode_check_form(check bool, message string){
	if !check{
		Warning.Println(message)
	}
}

func loadHTM(page_name string) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name))
	if err == nil{
		existingLength := len(bytes)
		bytes = dev_mode_minifyHtml(page_name, bytes)
		newLength := len(bytes)
		if existingLength != newLength {
			Trace.Printf("Page '%v' had %v bytes removed (%v percent), from: %v, to: %v", page_name, existingLength-newLength, 100-newLength*100/existingLength, existingLength, newLength)
		}
		if newLength > existingLength {
			Error.Println("How did this page get bigger?")
		}
	}else{
		ioutil.WriteFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name), bytes, 0777)
	}
	ioutil.WriteFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name), bytes, 0777)
	return bytes
}

func dev_mode_minifyHtml(pageName string, minify []byte) []byte {
	if bytes.Contains(minify, []byte("ZgotmplZ")) {
		Warning.Println("Template generation error: ZgotmplZ")
		return []byte("")
	}
	removeChars := []string{
		"	", //Tab
		"\n", //new line
		"\r", //carriage return
	}
	//TODO remove spaces between block elements like: </div> <div> but keep between inline elements like </span> <span>
	//TODO use improved regex for better searching & replacement
	replaceChars := map[string]string{
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
	for _, search := range removeChars {
		minify = bytes.Replace(minify, []byte(search), []byte(""), -1)
	}
	for search, replace := range replaceChars {
		length := len(minify)
		minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
		if length != len(minify) {
			Warning.Printf("A dodgy character (%v) was found in '%v'! Please replace with (%v).", search, pageName, replace)
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

func exists(dict M, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func main(){
	agent.Verbose = true
	agent.CollectHTTPStat = true
	agent.NewrelicLicense = "abf730f5454a9a1e78af7a75bfe04565e9e0d3f1"
	agent.Run()
	start()
	Post(URL_randomData, dev_mode_random_data)
	Info.Println("ready to go")
	Warning.Println("ListenAndServe: %v", http.ListenAndServe(":81", nil))
}
