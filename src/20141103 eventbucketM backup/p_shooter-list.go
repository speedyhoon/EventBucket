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

//		fmt.Printf("page: %v\n", page_count)
		response, err := http.Get(fmt.Sprintf("%v%v", url, page_count))
		defer response.Body.Close()
		if err != nil {
			//TODO change to the error framework with a helpfull error message
			fmt.Println("ERROR: http.Get", err)
			return 0
		}

		doc, err := html.Parse(response.Body)
		if err != nil {
			vardump(err)
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
	fmt.Println("Finished copying from website.")

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
	fmt.Println("Finished inserting new shooters.")
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
	checkErr(err)

	query := make(map[string]interface{}, 0)
	if t.Surname != "" {
		query["s"] = map[string]interface{}{"$regex": fmt.Sprintf(`^%v`, t.Surname), "$options": "i"}
	}
	if t.First != ""{
		query["f"] = map[string]interface{}{"$regex": fmt.Sprintf(`^%v`, t.First), "$options": "i"}
	}
	if t.Club != ""{
		query["c"] = map[string]interface{}{"$regex": fmt.Sprintf(`^%v`, t.Club), "$options": "i"}
	}

	//TODO research would changing map[string]interface{} to M{} have any affect on the size of the compiled application?
//	type M map[string]interface{}
//	type S []M

	//Ignore Deleted shooters. Selects not modified, updated & merged shooters
//	query["$or"] = M{"$or": S{{"t": nil}, {"t": 0 }}}
	query["$or"] = []map[string]interface{}{{"t": nil}, {"t": map[string]interface{}{"$lt": 3 } }}
//	query["t"] = map[string]interface{}{"$or": []map[string]interface{}{{"t": nil}, {"t", map[string]interface{}{"$lt": 3 }}}}
//	query["t"] = map[string]interface{}{"$or": []map[string]interface{}{{"_id": "id1"}, {"_id": "id2"}}}



//	db.S.find({"s": {"$regex": /^Webb$/}      ,"$or": [ {"t": {"$lt": 3}}, {"t": null}]        })


//	v := M{"$or": S{{"_id": "id1"}, {"_id": "id2"}}}




//	tree := []map[string]interface{}{}

	//	type M map[string]interface{}
	//	type S []M

//	v := map[string]interface{}{"$or": []map[string]interface{}{{"_id": "id1"}, {"_id": "id2"}}}
//	v := map[string]interface{}{"$or": []map[string]interface{}{{"t": nil}, {"t", {"$lt": 3 }}}}
//
//{ age: { $lt: 25 } }


	found := searchShooters(query)

	var option_list []Option
	for _, shooter := range found{
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
		Inputs: map[string]Inputs{
			"first":Inputs{
				Html:      "text",
				Label:   "First Name",
//				Required: true,
			},
			"surname":Inputs{
				Html:      "text",
				Label:   "Surname",
//				Required: true,
			},
			"club":Inputs{
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
//				Required: true,
			},
		},
	}
}

func event_add_existing_shooterForm() Form {
	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
			"sid":Inputs{
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
				Required: true,
			},
			"age":Inputs{
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2,
				Required: true,
			},
			"grade":Inputs{
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				//				SelectedValues: available_classes_grades(event),
			},
			"event_id":Inputs{
				Html: "hidden",
//				Value: event_id,
				Value: "",
			},
		},
	}
}
