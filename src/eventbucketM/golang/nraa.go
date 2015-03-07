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
	NRAA_SHOOTER_LIST_URL       = "http://www.nraa.com.au/nraa-shooter-list/?_p="
	NRAA_SHOOTER_LIST_GRADE_URL = "http://www.nraa.com.au/wp-admin/admin-ajax.php?action=get-calculated-grades&shooter_id="
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
	Info.Println("Starting to download shooter list from website.")
	//TODO get the max number of pages <div class="pagination"><a href="http://www.nraa.com.au/nraa-shooter-list/?_p=524">Last</a>
	maxPages := 3 //maxPages := 524
	var appendShooterIds []int
	for pageCount := 2; pageCount <= maxPages; pageCount += 1 {
		response, err := http.Get(fmt.Sprintf("%v%v", NRAA_SHOOTER_LIST_URL, pageCount))
		defer response.Body.Close()
		if err != nil {
			Warning.Printf("Unable to get page %v http.Get %v", pageCount, err) //TODO Improve the error framework with a helpful error message
			break
		}
		htmlBody, err := html.Parse(response.Body)
		if err != nil {
			Warning.Printf("Unable to parse HTML response: ", err)
			break
		}
		var i int
		var trimSpace string
		var shooter NraaShooter

		var findCells func(*html.Node)
		//TODO this would run faster declaring findCells & findRows outside the loop
		findCells = func(n *html.Node) {
			if n.Type == html.TextNode {
				trimSpace = strings.TrimSpace(n.Data)
				if trimSpace != "" {
					if i >= 1 {
						switch i {
						case 1:
							var err error
							shooter.SID, err = strconv.Atoi(trimSpace)
							if err != nil {
								Error.Printf(fmt.Sprintf("Unable to convert shooter id %v to int", err))
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
					i += 1
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				findCells(c)
			}
		}

		var findRows func(*html.Node)
		findRows = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "tr" {
				for _, attr := range n.Attr {
					if attr.Key == "data-shooter-id" && attr.Val != "" {
						i = 0
						id, err := strconv.Atoi(attr.Val)
						if err == nil {
							shooter = NraaShooter{NraaId: id}
							findCells(n)
							nraaUpsertShooter(shooter)
							appendShooterIds = append(appendShooterIds, id)
						} else {
							Warning.Printf(fmt.Sprintf("Unable to convert shooter id %v to int", err))
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				findRows(c)
			}
		}
		findRows(htmlBody)
	}
	Info.Println("Finished downloading shooter list from website.")
	nraaShooterGrades(appendShooterIds)
	nraaCopyShooters()
	nraaLastUpdated()
	return 1
}

func nraaShooterGrades(shooterIdList []int) {
	Info.Println("Starting to download shooter grades from website.")
	for _, shooterId := range shooterIdList {
		//Query the server for shooterId's grades
		response, err := http.Get(fmt.Sprintf("%v%v", NRAA_SHOOTER_LIST_GRADE_URL, shooterId))
		defer response.Body.Close()
		if err != nil {
			//Unable to contact the server
			Warning.Printf("http.Get", err)
			break
		}
		//Decode the response to JSON
		var m []NraaGrade
		err = json.NewDecoder(response.Body).Decode(&m)
		if err != nil || err == io.EOF {
			//if err != nil : The JSON returned contained an error & couldn't be decoded
			//if err == io.EOF : There was no JSON data in the returned string
			Warning.Printf("json.Decode", err)
			break
		}
		var grades []NraaGrading
		for _, data := range m {
			grades = append(grades, NraaGrading{
				DisciplineId:   data.Discipline.Id,
				DisciplineName: data.Discipline.Name,
				GradeId:        data.Grade.Id,
				GradeName:      data.Grade.Name,
				GradeThreshold: data.Grade.Threshold,
				AvgScore:       str2float(data.AvgScore),
				ShootQty:       str2Int(data.ShootQty),
			})
		}
		nraaUpdateGrading(shooterId, grades)
	}
	Info.Println("Finished downloading shooter grades from website.")
}

func nraaCopyShooters() int {
	Info.Println("Started inserting new shooters.")
	counter := 0
	shooter_list := getShooterLists()
	for _, n_shooter := range shooter_list {
		shooter := getShooterList(n_shooter.SID)
		if shooter.SID != 0 && shooter.NraaId != 0 && shooter.Surname != "" && shooter.FirstName != "" && shooter.NickName != "" && shooter.Club != "" && shooter.Address != "" && shooter.Email != "" {
			UpsertDoc("shooter", n_shooter.SID, n_shooter)
			counter += 1
		}
	}
	Info.Println("Finished inserting new shooters.")
	return counter
}
