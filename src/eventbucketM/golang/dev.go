// +build dev

package main

import (
	//	"bytes"
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
	Trace   = log.New(os.Stdout, "TRACE:   ", log.Lshortfile)
	Info    = log.New(os.Stdout, "INFO:    ", log.Lshortfile)
	Warning = log.New(os.Stderr, "WARNING: ", log.Lshortfile)
)

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

func serveDir(contentType string) {
	http.Handle(contentType,
		http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//defer devModeTimeTrack(time.Now(), r.RequestURI)
			if strings.HasSuffix(r.URL.Path, "/") {
				//If url is a directory return a 404 to prevent displaying a directory listing
				http.NotFound(w, r)
				return
			}
			httpHeaders(w, []string{"expire", "cache", contentType, "public"})
			Gzip(http.FileServer(http.Dir("^^DIR_ROOT^^")), w, r)
		})))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer devModeTimeTrack(time.Now(), r.RequestURI)
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
		Gzip(h, w, r)
	}))
}

func dump(input ...interface{}) {
	for _, print := range input {
		Trace.Printf("\n%v", print)
	}
}
func vardump(input ...interface{}) {
	for _, print := range input {
		Trace.Printf("\n%+v", print) //map field names included
	}
}
func export(input ...interface{}) {
	for _, print := range input {
		Trace.Printf("\n%#v", print) //can copy and declare new variable with it. Most ouput available
	}
}

func devModeTimeTrack(start time.Time, requestURI string) {
	Trace.Printf("%s took %s", requestURI, time.Since(start))
}

func devModeCheckForm(check bool, message string) {
	if !check {
		Warning.Println(message)
	}
}

func loadHTM(pageName string) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_MINIFIED, pageName))
	if err != nil {
		Error.Println(err)
	}
	return bytes
}

func randomData(w http.ResponseWriter, r *http.Request) {
	var eventId, rangeId string
	var totalScores, startShooting bool
	var properties []string
	var shooterQty int
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
			startShooting = true
		case "shooterQty":
			shooterQty, _ = strconv.Atoi(properties[1])
		}
	}
	if eventId == "" {
		Error.Printf("Need a valid event Id to proceed.")
	}
	if shooterQty > 0 {
		randomDataShooterQty(shooterQty, eventId)
	}
	if startShooting {
		randomDataStartShooting(eventId, rangeId, w)
	} else if totalScores {
		randomDataTotalScores(eventId, rangeId, w)
	}
}

func randomDataStartShooting(eventId, rangeId string, w http.ResponseWriter) {
	event, _ := getEvent(eventId)
	for shooterId, shooter := range event.Shooters {
		http.PostForm("http://localhost:81/updateShotScores",
			url.Values{
				"shots":     {randomShooterScores(shooter.Grade)},
				"shooterid": {fmt.Sprintf("%v", shooterId)},
				"rangeid":   {rangeId},
				"eventid":   {eventId},
			},
		)
	}
	fmt.Fprintf(w, "<p>Finished StartShooting for all shooters in event %v", eventId)
}

func randomDataTotalScores(eventId, rangeId string, w http.ResponseWriter) {
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
		range_Id, err := strconv.Atoi(rangeId)
		if err == nil {
			eventTotalScoreUpdate(eventId, range_Id, []int{shooter_id}, Score{
				Total:   rand.Intn(51),
				Centers: rand.Intn(11),
			})
		}
	}
	fmt.Fprintf(w, "<p>Finished TotalScores for all shooters in event %v", eventId)
}

func randomDataShooterQty(shooterQty int, eventId string) {
	counter := 0
	for counter < shooterQty {
		//make some requests for x number of shooters
		counter += 1
		Trace.Printf("inserting shooter :%v", counter)
		eventShooterInsert(eventId, EventShooter{
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
	var shots string
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < shooterClass.SightersQty; i++ {
		shots += string(shooterClass.ValidSighters[rand.Intn(len(shooterClass.ValidSighters))])
	}
	availableShots := shooterClass.Buttons + "-"
	for i := 0; i < shooterClass.ShotsQty; i++ {
		shots += string(availableShots[rand.Intn(len(availableShots))])
	}
	return shots
}

func slice_to_map_bool(input []string) map[string]bool {
	output := make(map[string]bool)
	for _, value := range input {
		output[value] = true
	}
	return output
}
