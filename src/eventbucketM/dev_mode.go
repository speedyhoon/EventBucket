package main

import (
	"fmt"
	//	"os"
	"bytes"
	"io/ioutil"
	//	"log"
	"time"
	"net/http"
	"net/url"
	"strings"
	"go-randomdata-master"
	"math/rand"
	"strconv"



	"github.com/yvasiyarov/gorelic"
)

const (
	dev_mode_DEBUG = false	//Send system metric data to NewRelic.com
)
var agent = gorelic.NewAgent()

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

func dev_mode_timeTrack(start time.Time, requestURI string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", requestURI, elapsed)
}

func dev_mode_check_form(check bool, message string){
	if !check{
		fmt.Printf("%v\n", message)
	}
}

func dev_mode_loadHTM(page_name string, existing_minified_file []byte) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name))
	checkErr(err)
	bytes = dev_mode_minifyHtml(page_name, bytes)
	existing_len := len(existing_minified_file)
	new_len := len(bytes)
	if existing_len != new_len {
		fmt.Printf("Page '%v' had %v bytes removed (%v percent), total: %v, from: %v", page_name, new_len-existing_len, (existing_len*100/new_len-100)*-1, existing_len, new_len)

//		return bytes
	}
	ioutil.WriteFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name), bytes, 0777)
	return bytes
	//	return existing_minified_file
}

func dev_mode_minifyHtml(page_name string, html []byte) []byte {
//	if MINIFY {
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
//	}else{
//		return html
//	}
}

func dev_mode_random_data(w http.ResponseWriter, r *http.Request) {
	eventId := ""
	rangeId := ""
	random_grades := []string{/*"a","b","c",*/"d", "e", "f", "g", "h", "i", "j"}
	attributes := strings.Split(strings.Replace(r.RequestURI, "/random-data/", "", -1), "&")
	for _, request := range attributes {
		properties := strings.Split(request, "=")
		if properties[0] == "event_id" {
			eventId = properties[1]
		} else if properties[0] == "range_id" {
			rangeId = properties[1]
		}
	}

	for _, request := range attributes {
		properties := strings.Split(request, "=")
		switch properties[0] {
		case "shooterQty":
			shooterQty, _ := strconv.Atoi(properties[1])
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
					"grade":     {random_grades[rand.Intn(len(random_grades)-1)]},
					"event_id":  {eventId},
				})
				//				defer resp.Body.Close()
			}
		case "totalScores":
			event, _ := getEvent(eventId)
			for shooter_id, _ := range event.Shooters {
				rand.Seed(time.Now().UnixNano())
				//				rand.Seed(90)
				fmt.Println(shooter_id)
				//				resp, _ := http.PostForm("http://localhost/updateTotalScores",
				//				go http.PostForm("http://localhost/updateTotalScores",
				//					url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
				//					"score":			{fmt.Sprintf("%v.%v",rand.Intn(51),rand.Intn(11))},
				//					"shooter_id":	{shooter_id},
				//					"range_id":		{range_id},
				//					"event_id":		{event_id},
				//				})
				//				defer resp.Body.Close()
				range_Id, _ := strToInt(rangeId)
				eventTotalScoreUpdate(eventId, range_Id, []string{shooter_id}, Score{
					Total: rand.Intn(51),
					Centers: rand.Intn(11),
				})
			}
		case "startShooting":
			event, _ := getEvent(eventId)
			for shooterId, shooter := range event.Shooters {
				go http.PostForm("http://localhost/updateShotScores",
						url.Values{
							"shots":			{randomShooterScores(shooter.Grade)},
							"shooter_id":	{shooterId},
							"range_id":		{rangeId},
							"event_id":		{eventId},
						},
				)
			}
		}
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
