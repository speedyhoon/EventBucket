package main

import (
	"net/http"
	"fmt"
	"code.google.com/p/go.net/html"
	"strings"
	"strconv"
	"encoding/json"
)

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

//		Trace.Printf("page: %v\n", page_count)
		response, err := http.Get(fmt.Sprintf("%v%v", url, page_count))
		defer response.Body.Close()
		if err != nil {
			//TODO change to the error framework with a helpfull error message
			Warning.Printf("ERROR: http.Get", err)
			return 0
		}

		doc, err := html.Parse(response.Body)
		if err != nil {
			Error.Printf("%+v\n", err)
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
	Info.Println("Finished copying from website.")

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
	Info.Println("Finished inserting new shooters.")
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
	if err != nil {
		Warning.Println(err)
	}
	query := make(M, 0)
	if t.Surname != "" {
		query["s"] = M{"$regex": fmt.Sprintf(`^%v`, t.Surname), "$options": "i"}
	}
	if t.First != ""{
		query["f"] = M{"$regex": fmt.Sprintf(`^%v`, t.First), "$options": "i"}
	}
	if t.Club != ""{
		query["c"] = M{"$regex": fmt.Sprintf(`^%v`, t.Club), "$options": "i"}
	}

	//Ignore Deleted shooters. Selects not modified, updated & merged shooters
	query["$or"] = []M{{"t": nil}, {"t": M{"$lt": 3 } }}
	var option_list []Option
	for _, shooter := range searchShooters(query){
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
		Inputs: []Inputs{{
			Name: "first",
				Html:      "text",
				Label:   "First Name",
			},{
				Name: "surname",
				Html:      "text",
				Label:   "Surname",
			},{
		Name: "club",
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
			},
		},
	}
}
