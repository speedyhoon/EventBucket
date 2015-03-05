// +build dev

package main

import (
	"bytes"
	"fmt"
	"go-randomdata-master"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	//"github.com/yosssi/ace"
	"github.com/yvasiyarov/gorelic"
)

const (
	NEWRELIC         = false //Send logging data to New Relic
	PATH_HTML_SOURCE = "html/%v.html"
	URL_randomData   = "/randomData/"
)

var (
	agent = gorelic.NewAgent()
	//Use io.Writer >>> ioutil.Discard to disable logging any output
	Trace   = log.New(os.Stdout, "TRACE:   ", log.Ltime|log.Lshortfile)
	Info    = log.New(os.Stdout, "INFO:    ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func serveDir(contentType string) {
	http.Handle(contentType,
		http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			//If url is a directory return a 404 to prevent displaying a directory listing
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			httpHeaders(w, []string{"expire", "cache", contentType, "public"})
			Gzip(http.FileServer(http.Dir("^^DIR_ROOT^^")), w, r)
		})))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		defer dev_mode_timeTrack(time.Now(), r.RequestURI)
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

func dev_mode_check_form(check bool, message string) {
	if !check {
		Warning.Println(message)
	}
}

func loadHTM(page_name string) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name))
	if err == nil {
		existingLength := len(bytes)
		bytes = dev_mode_minifyHtml(page_name, bytes)
		newLength := len(bytes)
		if existingLength != newLength {
			//Trace.Printf("Page '%v' had %v bytes removed (%v percent), from: %v, to: %v", page_name, existingLength-newLength, 100-newLength*100/existingLength, existingLength, newLength)
		}
		if newLength > existingLength {
			Error.Println("How did this page get bigger?")
		}
	} else {
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

func randomData(w http.ResponseWriter, r *http.Request) {
	eventId := ""
	rangeId := ""
	shooterQty := 0
	totalScores := false
	startShooting := false
	var properties []string
	//	random_grades := []string{/*"a","b","c",*/"d", "e", "f", "g", "h", "i", "j"}
	for _, request := range strings.Split(strings.Replace(r.RequestURI, URL_randomData, "", -1), "&") {
		properties = strings.Split(request, "=")
		switch properties[0] {
		case "eventId":
			eventId = properties[1]
			break
		case "rangeId":
			rangeId = properties[1]
			break
		case "totalScores":
			totalScores = true
			break
		case "startShooting":
			totalScores = true
		case "shooterQty":
			shooterQty, _ = strconv.Atoi(properties[1])
		}
	}
	if eventId == "" {
		Error.Printf("Need a valid event Id to proceed.")
	}
	if shooterQty > 0 {
		go randomDataShooterQty(shooterQty, eventId)
	}
	if totalScores {
		go randomDataTotalScores(eventId, rangeId)
	}
	if startShooting {
		go randomDataStartShooting(eventId, rangeId)
	}
}

func randomDataStartShooting(eventId, rangeId string) {
	event, _ := getEvent(eventId)
	for shooterId, shooter := range event.Shooters {
		http.PostForm("http://localhost/updateShotScores",
			url.Values{
				"shots":     {randomShooterScores(shooter.Grade)},
				"shooterid": {fmt.Sprintf("%v", shooterId)},
				"rangeid":   {rangeId},
				"eventid":   {eventId},
			},
		)
	}
}

func randomDataTotalScores(eventId, rangeId string) {
	event, _ := getEvent(eventId)
	for shooter_id := range event.Shooters {
		rand.Seed(time.Now().UnixNano()) //Use rand.Seed(90) with a constant number to make the same number
		/*resp, _ := http.PostForm("http://localhost/updateTotalScores",
		go http.PostForm("http://localhost/updateTotalScores",
			url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
			"score":			{fmt.Sprintf("%v.%v",rand.Intn(51),rand.Intn(11))},
			"shooterid":	{shooterid},
			"rangeid":		{rangeid},
			"eventid":		{eventid},
		})
		resp.Body.Close()*/
		range_Id, err := strToInt(rangeId)
		if err == nil {
			eventTotalScoreUpdate(eventId, range_Id, []int{shooter_id}, Score{
				Total:   rand.Intn(51),
				Centers: rand.Intn(11),
			})
		}
	}
}

func randomDataShooterQty(shooterQty int, eventId string) {
	counter := 0
	for counter < shooterQty {
		//make some requests for x number of shooters
		counter += 1
		Trace.Printf("inserting shooter :%v", counter)
		event_shooter_insert(eventId, EventShooter{
			FirstName: randomdata.FirstName(randomdata.RandomGender),
			Surname:   randomdata.LastName(),
			Club:      randomdata.State(randomdata.Large),
			AgeGroup:  "N",
			Grade:     rand.Intn(8),
		})
	}
}

func randomShooterScores(shooterGrade int) string {
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

func main() {
	if NEWRELIC {
		agent.Verbose = true
		agent.CollectHTTPStat = true
		agent.NewrelicLicense = "abf730f5454a9a1e78af7a75bfe04565e9e0d3f1"
		agent.Run()
	}
	start()
	Post(URL_randomData, randomData)
	Info.Println("ready to go")
	Warning.Println("ListenAndServe: %v", http.ListenAndServe(":81", nil))
}
