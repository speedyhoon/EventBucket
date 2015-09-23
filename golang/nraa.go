package main

import (
	"code.google.com/p/go.net/html"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	nraaShooterListURL      = "http://www.nraa.com.au/nraa-shooter-list/?_p="
	nraaShooterListGradeURL = "http://www.nraa.com.au/wp-admin/admin-ajax.php?action=get-calculated-grades&shooter_id="
)

/*	Sample grade data from NRAA website
[{
	"Discipline": {"id":"1","name":"Target Rifle"},
	"Grade": {"id":"2","name":"B","threshold":"94.50"},
	"Shooter":{"id":"15","name":"Peter Evans"},
	"avg_score":"94.63879",
	"number_of_shoots":"2"
}]*/

func nraaStartUpdateShooterList(w http.ResponseWriter, r *http.Request) {
	go nraaUpdateShooterList()
}

func nraaUpdateShooterList() int {
	info.Println("Starting to download shooter list from website.")
	//TODO get the max number of pages <div class="pagination"><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=524">Last</a>
	var appendShooterIDs []int
	var i int
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
							warning.Printf(fmt.Sprintf("Unable to convert shooter id %v to int", err))
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
			for _, attr := range n.Attr {
				if attr.Key == "data-shooter-id" && attr.Val != "" {
					i = 0
					id, err = strconv.Atoi(attr.Val)
					if err == nil {
						shooter = NraaShooter{NraaID: id}
						findCells(n)
						nraaUpsertShooter(shooter)
						appendShooterIDs = append(appendShooterIDs, id)
					} else {
						warning.Printf(fmt.Sprintf("Unable to convert shooter id %v to int", err))
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
	for pageCount := 1; pageCount <= 524; pageCount++ {
		response, err = http.Get(fmt.Sprintf("%v%v", nraaShooterListURL, pageCount))
		defer response.Body.Close()
		if err != nil {
			warning.Printf("Unable to get page %v http.Get %v", pageCount, err) //TODO Improve the error framework with a helpful error message
			break
		}
		htmlBody, err = html.Parse(response.Body)
		if err != nil {
			warning.Printf("Unable to parse HTML response: %v", err)
			break
		}
		findRows(htmlBody)
	}
	info.Println("Finished downloading shooter list from website.")
	nraaShooterGrades(appendShooterIDs)
	nraaCopyShooters()
	nraaLastUpdated()
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
			warning.Printf("http.Get: %v", err)
			break
		}
		//Decode the response to JSON
		var m []NraaGrade
		err = json.NewDecoder(response.Body).Decode(&m)
		if err != nil || err == io.EOF {
			//if err != nil : The JSON returned contained an error & couldn't be decoded
			//if err == io.EOF : There was no JSON data in the returned string
			warning.Printf("json.Decode: %v", err)
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
		nraaUpdateGrading(shooterID, grades)
	}
	info.Println("Finished downloading shooter grades from website.")
}

func nraaCopyShooters() int {
	info.Println("Started inserting new shooters.")
	counter := 0
	for _, nShooter := range getShooterLists() {
		shooter := getShooterList(nShooter.SID)
		if shooter.SID != 0 && shooter.NraaID != 0 && shooter.Surname != "" && shooter.FirstName != "" && shooter.NickName != "" && shooter.Club != "" && shooter.Address != "" && shooter.Email != "" {
			upsertDoc("shooter", nShooter.SID, nShooter)
			counter++
		}
	}
	info.Println("Finished inserting new shooters.")
	return counter
}
