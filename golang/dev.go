// +build dev

package main

import (
	"fmt"
	"github.com/yvasiyarov/gorelic"
	"go-randomdata-master"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/handler/gzip"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	newRelic      = false //Send logging data to New Relic
	urlRandomData = "/randomData/"
)

var (
	agent = gorelic.NewAgent()
	//Use io.Writer >>> ioutil.Discard to disable logging any output
	trace   = log.New(os.Stdout, "TRACE:   ", log.Lshortfile)
	info    = log.New(os.Stdout, "INFO:    ", log.Lshortfile)
	warning = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//shootersMakeList = generateForm(makeShooterList())
)

func main() {
	if newRelic {
		agent.Verbose = true
		agent.CollectHTTPStat = true
		agent.NewrelicLicense = "abf730f5454a9a1e78af7a75bfe04565e9e0d3f1"
		agent.Run()
	}
	go startDatabase(false)
	//	http.Handle(dirJS, http.StripPrefix("/j", http.FileServer(http.Dir(dirRoot))))
	//	http.Handle(dirJS, gzip.FileServer(http.Dir(dirRoot)))
	//	http.Handle("/j/", http.FileServer(http.Dir(dirRoot)))
	http.Handle("/j/", gzip.FileServer(http.Dir(dirRoot)))

	start()
	post(urlRandomData, randomData)
	//Nraa
	post(urlMakeShooterList, postVia(nraaStartUpdateShooterList, urlShooters))
	warning.Println("ListenAndServe: %v", http.ListenAndServe(":81", nil))
}

func serveDir(contentType string) {
	//	http.Handle(contentType, gzip.FileServer(http.Dir(dirRoot)))
	//	http.Handle("/gzip/", http.StripPrefix("/gzip", gzip.FileServer(http.Dir(*base))))
	http.Handle(contentType,
		http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//defer devModeTimeTrack(time.Now(), r.RequestURI)
			if strings.HasSuffix(r.URL.Path, "/") {
				//If url is a directory return a 404 to prevent displaying a directory listing
				http.NotFound(w, r)
				return
			}
			httpHeaders(w, []string{"expire", "cache", contentType, "public"})
			gzipperFiles(http.FileServer(http.Dir(dirRoot)), w, r)
		})))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer devModeTimeTrack(time.Now(), r.RequestURI)
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
		gzipper(h, w, r)
	}))
}

func dump(input ...interface{}) {
	for _, print := range input {
		trace.Printf("%v", print)
	}
}

/*func vardump(input ...interface{}) {
	for _, print := range input {
		trace.Printf("\n%+v", print) //map field names included
	}
}*/
func export(input ...interface{}) {
	for _, print := range input {
		trace.Printf("%#v", print) //can copy and declare new variable with it. Most ouput available
	}
}

/*func devModeTimeTrack(start time.Time, requestURI string) {
	trace.Printf("%s took %s", requestURI, time.Since(start))
}

func devModeCheckForm(check bool, message string) {
	if !check {
		warning.Println(message)
	}
}*/

func loadHTM(pageName string) string {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(pathHTMLMinified, pageName))
	if err != nil {
		warning.Println(err)
	}
	return string(bytes)
}

func randomData(w http.ResponseWriter, r *http.Request) {
	var eventID, rangeID string
	var totalScores, startShooting bool
	var properties []string
	var shooterQty int
	for _, request := range strings.Split(strings.Replace(r.RequestURI, urlRandomData, "", -1), "&") {
		properties = strings.Split(request, "=")
		switch strings.ToLower(properties[0]) {
		case "eventid":
			eventID = properties[1]
			break
		case "rangeid":
			rangeID = properties[1]
			break
		case "totalscores":
			totalScores = true
			break
		case "startshooting":
			startShooting = true
		case "shooterqty":
			shooterQty = str2Int(properties[1])
		}
	}
	if eventID == "" {
		warning.Printf("Need a valid event ID to proceed.")
	}
	if shooterQty > 0 {
		randomDataShooterQty(shooterQty, eventID)
	}
	if startShooting {
		randomDataStartShooting(eventID, rangeID, w)
	} else if totalScores {
		randomDataTotalScores(eventID, rangeID, w)
	}
}

func randomDataStartShooting(eventID, rangeID string, w http.ResponseWriter) {
	event, _ := getEvent(eventID)
	for shooterID, shooter := range event.Shooters {
		http.PostForm("http://localhost:81/updateShotScores",
			url.Values{
				"shots":     {randomShooterScores(shooter.Grade)},
				"shooterid": {fmt.Sprintf("%v", shooterID)},
				"rangeid":   {rangeID},
				"eventid":   {eventID},
			},
		)
	}
	fmt.Fprintf(w, "<p>Finished StartShooting for all shooters in event %v", eventID)
}

func randomDataTotalScores(eventID, rangeID string, w http.ResponseWriter) {
	event, _ := getEvent(eventID)
	for shooterID := range event.Shooters {
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
		rangeID, err := strconv.Atoi(rangeID)
		if err == nil {
			eventTotalScoreUpdate(eventID, rangeID, []int{shooterID}, Score{
				Total:   rand.Intn(51),
				Centres: rand.Intn(11),
			})
		}
	}
	fmt.Fprintf(w, "<p>Finished TotalScores for all shooters in event %v", eventID)
}

func randomDataShooterQty(shooterQty int, eventID string) {
	counter := 0
	for counter < shooterQty {
		//make some requests for x number of shooters
		counter++
		trace.Printf("inserting shooter #%v", counter)
		eventShooterInsert(eventID, EventShooter{
			FirstName: randomdata.FirstName(randomdata.RandomGender),
			Surname:   randomdata.LastName(),
			Club:      randomdata.State(randomdata.Large),
			AgeGroup:  "N",
			Grade:     rand.Intn(8),
		})
	}
	info.Println("Finished inserting shooters")
}

func randomShooterScores(shooterGrade int) string {
	shooterClass := grades()[shooterGrade].settings
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
