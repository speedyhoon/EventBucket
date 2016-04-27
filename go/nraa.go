package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	"code.google.com/p/go.net/html"
)

const (
	nraaDomain              = "http://www.nraa.com.au"
	nraaShooterListURL      = nraaDomain + "/nraa-shooter-list/?_p="
	nraaShooterListGradeURL = nraaDomain + "/wp-admin/admin-ajax.php?action=get-calculated-grades&shooter_id="
)

/*	Sample grade data from NRAA website
[
	{
		"Discipline": {
			"id":"2",
			"name":"F Standard"
		},
		"Grade": {
			"id":"4",
			"name":"A",
			"threshold":"94.50"
		},
		"Shooter":{
			"id":"2336",
			"name":"Barry Roennfeldt"
		},
		"avg_score":"96.99557",
		"number_of_shoots":"8"
	},
	{
		"Discipline":{
			"id":"3",
			"name":"F Open"
		},
		"Grade":{
			"id":"6",
			"name":"FO",
			"threshold":"0.00"
		},
		"Shooter":{
			"id":"2336",
			"name":"Barry Roennfeldt"
		},
		"avg_score":"96.50634",
		"number_of_shoots":"2"
	}
]*/

// NraaShooter is exported
type NraaShooter struct {
	ID        int                    `json:"I"`
	NID       int                    `json:"N"` //NRAA sequential integer id.
	SID       int                    `json:"M,omitempty"`
	Surname   string                 `json:"s,omitempty"`
	FirstName string                 `json:"f,omitempty"`
	NickName  string                 `json:"n,omitempty"`
	Club      string                 `json:"c,omitempty"`
	Grades    map[string]NraaGrading `json:"g,omitempty,inline"`
}

// NraaGrade is exported
type NraaGrade struct {
	Discipline NraaDetails `json:"Discipline,omitempty"`
	Grade      NraaDetails `json:"Grade,omitempty"`
	AvgScore   string      `json:"avg_score,omitempty"`
	ShootQty   string      `json:"number_of_shoots,omitempty"`
	//ignore Shooter
}

// NraaDetails is exported
type NraaDetails struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Threshold string `json:"threshold,omitempty"`
}

// NraaGrading is exported
type NraaGrading struct {
	Class     string  `json:"c,omitempty"`
	Grade     string  `json:"g,omitempty"`
	Threshold string  `json:"t,omitempty"`
	AvgScore  float64 `json:"a,omitempty"`
	ShootQty  int     `json:"s,omitempty"`
}

//Converts NRAA discipline id to EventBucket discipline id. Designed to be separated so when NRAA changes their JSON interface again, EventBucket doesn't have to change any core structure.
func nraaDiscipline(disciplineID string) uint {
	id, ok := map[string]uint{
		"1": 0,
		"2": 1,
		"3": 1,
	}[disciplineID]
	if !ok {
		warn.Printf("Unable to find discipline id '%v'\n", disciplineID)
	}
	return id
}

func nraaStartUpdateShooterList(w http.ResponseWriter, r *http.Request) {
	go nraaUpdateShooterList()
}

func nraaUpdateShooterList() int {
	info.Println("Starting to download shooter list from website.")
	//TODO get the max number of pages <div class="pagination"><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=524">Last</a>
	var appendShooterIDs []int
	var i int
	var shooterQty uint
	var trimSpace string
	var shooter NraaShooter
	var err error
	var findCells func(*html.Node)
	findCells = func(n *html.Node) {
		if n.Type == html.TextNode {
			trimSpace = strings.TrimSpace(n.Data)
			if trimSpace != "" {
				if i >= 1 {
					switch i {
					case 1:
						shooter.SID, err = strconv.Atoi(trimSpace)
						if err != nil {
							warn.Printf(fmt.Sprintf("Unable to convert shooter id %v to int", err))
							shooter = NraaShooter{} //Clear the NraaShooter so bad data doesn't get save to the DB
							return
						}
					case 2:
						shooter.Surname = trimSpace
					case 3:
						shooter.FirstName = trimSpace
					case 4:
						shooter.NickName = trimSpace
					case 5:
						shooter.Club = trimSpace
					}
				}
				i++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findCells(c)
		}
	}

	var findRows func(*html.Node)
	findRows = func(n *html.Node) {
		var id int
		if n.Type == html.ElementNode && n.Data == "tr" {
			shooterQty++
			for _, attr := range n.Attr {
				if attr.Key == "data-shooter-id" && attr.Val != "" {
					i = 0
					id, err = strconv.Atoi(attr.Val)
					if err == nil {
						shooter = NraaShooter{NID: id}
						findCells(n)
						//nraaUpsertShooter(shooter)
						appendShooterIDs = append(appendShooterIDs, id)
					} else {
						warn.Printf(fmt.Sprintf("Unable to convert shooter id %v to int", err))
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findRows(c)
		}
	}

	var htmlBody *html.Node
	var response *http.Response
	for pageCount := 1; pageCount <= 1; /*524*/ pageCount++ {
		response, err = http.Get(fmt.Sprintf("%v%v", nraaShooterListURL, pageCount))
		defer response.Body.Close()
		if err != nil {
			warn.Printf("Unable to get page %v http.Get %v", pageCount, err) //TODO Improve the error framework with a helpful error message
			break
		}
		htmlBody, err = html.Parse(response.Body)
		if err != nil {
			warn.Printf("Unable to parse HTML response: %v", err)
			break
		}
		findRows(htmlBody)
	}
	info.Println("Finished downloading shooter list from website.")
	nraaShooterGrades(appendShooterIDs)
	//nraaCopyShooters()
	nraaLastUpdated(NRAAUpdated{date: time.Now(), qty: shooterQty})
	return 1
}

func nraaShooterGrades(shooterIDList []int) {
	info.Println("Starting to download shooter grades from website.")
	for _, shooterID := range shooterIDList {
		//Query the server for shooterID's grades
		response, err := http.Get(fmt.Sprintf("%v%v", nraaShooterListGradeURL, shooterID))
		defer response.Body.Close()
		if err != nil {
			//Unable to contact the server
			warn.Printf("http.Get: %v", err)
			break
		}
		//Decode the response to JSON
		var m []NraaGrade
		err = json.NewDecoder(response.Body).Decode(&m)
		if err != nil || err == io.EOF {
			//if err != nil : The JSON returned contained an error & couldn't be decoded
			//if err == io.EOF : There was no JSON data in the returned string
			warn.Printf("json.Decode: %v", err)
			break
		}
		grades := map[string]NraaGrading{}
		for _, data := range m {
			if data.Discipline.ID != "" && data.Grade.Name != "" {
				grades[data.Discipline.ID] = NraaGrading{
					Class:     data.Discipline.Name,
					Grade:     data.Grade.Name,
					Threshold: data.Grade.Threshold,
					AvgScore:  str2float(data.AvgScore),
					ShootQty:  str2Int(data.ShootQty),
				}
			}
		}
		//nraaUpdateGrading(shooterID, grades)
	}
	info.Println("Finished downloading shooter grades from website.")
}

/*func nraaCopyShooters() int {
	info.Println("Started inserting new shooters.")
	counter := 0
	for _, nShooter := range getShooterLists() {
		shooter := getShooterList(nShooter.SID)
		if shooter.SID != 0 && shooter.NID != 0 && shooter.Surname != "" && shooter.FirstName != "" && shooter.NickName != "" && shooter.Club != "" && shooter.Address != "" && shooter.Email != "" {
			upsertDoc("shooter", nShooter.SID, nShooter)
			counter++
		}
	}
	info.Println("Finished inserting new shooters.")
	return counter
}*/

func str2Int(input interface{}) int {
	number, err := strconv.Atoi(fmt.Sprintf("%v", input))
	if err != nil {
		warn.Println(err)
	}
	return number
}

func str2float(input interface{}) float64 {
	float, err := strconv.ParseFloat(fmt.Sprintf("%v", input), 64)
	if err != nil {
		warn.Println(err)
	}
	return float
}

type NRAAUpdated struct {
	ID   string
	date time.Time
	qty  uint
}

func nraaLastUpdated(updated NRAAUpdated) error {
	//	var b36 string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblNRAAUpdated)
		if err != nil {
			return err
		}
		//		var id []byte
		b36, id := nextID(bucket)
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		updated.ID = b36

		// Marshal user data into bytes.
		buf, err := json.Marshal(updated)
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return err
}
