package main
import (
	"net/http"
//	"html/template"
	"fmt"
	"mgo"
	"strconv"
	//	"mgo/bson"
	//	"reflect"
//	"net/url"
)
var conn *mgo.Database//global variables that can't be a constant
func main(){
	conn = DB()
	http.HandleFunc("/startScoring", startScoring)
	http.HandleFunc("/save", save)
	http.ListenAndServe(":80", nil)
}
func getCollection(collection_name string)[]map[string]interface{}{
	var result []map[string]interface{}
	c := conn.C(collection_name)
	err := c.Find(nil).All(&result)
	checkErr(err)
	return result
}
func checkErr(err error){
	if err != nil {
		panic(err)
	}
}

type Scores struct{
	bowler int
//	bowler_wickets int
//	bowler_total int
	nb, wide, runs, stroke, wickets int
	batsman int
//	batsman_runs int
	facing int
	batsman2 int
//	batsman2_runs int
	how_out string
//	total_wickets int
//	total_runs int
}
type Tally struct{
//	bowler int
	bowler_wickets int
	bowler_total int
//	nb, wide, runs, stroke, wickets int
//	batsman int
	batsman_runs int
	batsman_wickets int
//	facing int
//	batsman2 int
//	batsman2_runs int
//	how_out string
	total_wickets int
	total_runs int
}
func header()string{
	output := "<tr>"
	output += "<td>Over</td>"
	output += "<td>Bowler</td>"
	output += "<td>NB</td>"
	output += "<td>Wide</td>"
	output += "<td>Bye</td>"
	output += "<td>LBye</td>"
	output += "<td>Total</td>"
	output += "<td>Batsman</td>"
	output += "<td>Stroke</td>"
	output += "<td>Runs</td>"
	output += "<td>Total</td>"
	output += "<td>Batsman</td>"
	output += "<td>Stroke</td>"
	output += "<td>Runs</td>"
	output += "<td>Total</td>"
	output += "<td>How Out</td>"
	output += "<td>Score</td>"
	output += "</tr>"
	return output
}

func InsertDoc(data interface{}, collection string){
	if data != false{
		err := conn.C(collection).Insert(data)
		checkErr(err)
	}
}
func DB()*mgo.Database{
	session, err := mgo.Dial("localhost")
	checkErr(err)
	//	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true)//this is supposed to be faster
	return session.DB("cricket")
}

func save(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		fmt.Print("\ndefinatly a POST")
		r.ParseForm()

//		t_bowler, _ := strconv.Atoi(r.Form["bowler"][0])
//		t_nb, _ := strconv.Atoi(r.Form["nb"][0])
//		t_wide, _:= strconv.Atoi(r.Form["wide"][0])
//		t_runs, _:= strconv.Atoi(r.Form["runs"][0])
//		t_stroke, _:= strconv.Atoi(r.Form["stroke"][0])
//		t_batsman, _:= strconv.Atoi(r.Form["batsman"][0])
//		t_batsman2, _:= strconv.Atoi(r.Form["batsman2"][0])
//		t_facing, _:= strconv.Atoi(r.Form["facing"][0])
		tryThis :=  make(map[string]string)
		for index, value := range r.Form{
			tryThis[index] = value[0]
		}

//		tryThis := map[string]string{
//			"bowler"   : r.Form["bowler"][0]    ,
//			"nb"       : r.Form["nb"][0]        ,
//			"wide"     : r.Form["wide"][0]      ,
//			"runs"     : r.Form["runs"][0]      ,
//			"stroke"   : r.Form["stroke"][0]    ,
//			"batsman"  : r.Form["batsman"][0]   ,
//			"batsman2" : r.Form["batsman2"][0]  ,
//			"facing"   : r.Form["facing"][0]    ,
//			"how_out"  : r.Form["howout"][0]    ,
//		}
		fmt.Print("\nInsert Data\n")
		fmt.Print(r.Form)
		InsertDoc(tryThis, "event")
	}else{
		fmt.Print("did not receive a http POST")
	}
	startScoring(w, r)
}

func startScoring(w http.ResponseWriter, r *http.Request){
	output := "<html><head><style>input[type=\"number\"],.short{width:45px;}</style></head><body>"
	output += "<table border=1>"
	output += header()
	batsmen := []string{
		"T.Hicky",
		"D.Donnel",
		"V.Taylor",
		"E.Sainz",
		"R.Harley",
	}
	bowlers := []string{
		"A.Maveric",
		"D.Powell",
		"C.Urn",
		"D.Junkie",
//		"A.Maveric",
	}
	balls_per_over := 6
	scores := map[int]Scores{
//		1:	[]Scores{
			1:Scores{
				bowler:	0,
				nb:		1,
				wide:		1,
				runs:		1,
				stroke:	3,
				batsman:	0,
				wickets:	0,
				batsman2:	1,
				facing:	0,
				how_out: "LBW",
			},
			2:Scores{
				nb:		1,
				wide:		1,
				runs:		1,
				stroke:	3,
				batsman:	0,
				wickets:	1,
				batsman2:	1,
				facing:	1,
				how_out: "LBW",
			},
			3:Scores{
				nb:		1,
				wide:		1,
				runs:		1,
				stroke:	3,
				batsman:	0,
				wickets:	0,
				batsman2:	1,
				facing:	1,
				how_out: "LBW",
			},
			4:Scores{
				nb:		0,
				wide:		0,
				runs:		1,
				stroke:	3,
				batsman:	0,
				wickets:	0,
				batsman2:	1,
				facing:	1,
				how_out: "",
			},
			5:Scores{
				nb:		1,
				wide:		1,
				runs:		1,
				stroke:	3,
				batsman:	0,
				wickets:	0,
				batsman2:	1,
				facing:	1,
				how_out: "",
			},
			6:Scores{
				nb:		1,
				wide:		1,
				runs:		1,
				stroke:	3,
				batsman:	0,
				wickets:	1,
				batsman2:	1,
				facing:	1,
				how_out: "",
			},
//		},
	}

	facing := -1
	batsman := -1
	batsman2 := -1
	bowler := -1
	over_number := 1
	over_ball_number := 1
	new_bowler := ""
	bowler_tally := map[int]Tally{
		0:Tally{
		bowler_wickets: 0,
		bowler_total: 0,
		batsman_runs : 0,
		batsman_wickets : 0,
		total_wickets : 0,
		total_runs : 0,
		},
		1:Tally{
			bowler_wickets: 0,
			bowler_total: 0,
			batsman_runs : 0,
			batsman_wickets : 0,
			total_wickets : 0,
			total_runs : 0,
		},
		2:Tally{
			bowler_wickets: 0,
			bowler_total: 0,
			batsman_runs : 0,
			batsman_wickets : 0,
			total_wickets : 0,
			total_runs : 0,
		},
		3:Tally{
			bowler_wickets: 0,
			bowler_total: 0,
			batsman_runs : 0,
			batsman_wickets : 0,
			total_wickets : 0,
			total_runs : 0,
		},
	}
//	var batsmen_tally map[int]Tally
	total_runs := 0
	total_wickets := 0


//	for over_index, over := range scores{
	for _, balls := range scores{
		facing = balls.facing
		batsman = balls.batsman
		batsman2 = balls.batsman2
		bowler = balls.bowler
		output += "<tr>"
		output += fmt.Sprintf("<td>%d.%d</td>", over_number, over_ball_number)
		if new_bowler != bowlers[balls.bowler]{
			new_bowler = bowlers[balls.bowler]
			output += fmt.Sprintf("\n\t<td rowspan=%d>%v</td>", overLength(over_ball_number,balls.bowler, scores), new_bowler)
		}
		output += "<td>"
		if balls.nb > 0{
			total_runs += balls.nb
			output += fmt.Sprintf("%d", balls.nb)
		}
		output += "</td><td>"
		if balls.wide > 0{
			output += fmt.Sprintf("%d", balls.wide)
			total_runs += balls.wide
		}
		output += "</td><td>"
		if balls.wide > 0{
			output += fmt.Sprintf("%d", balls.wide)
			total_runs += balls.wide
		}
		output += "</td><td>"
		if balls.wide > 0{
			output += fmt.Sprintf("%d", balls.wide)
			total_runs += balls.wide
		}

		if balls.wickets > 0{
//			bowler_tally[bowler].bowler_wickets = balls.wickets
			total_wickets += balls.wickets
//			bowler_tally[bowler].bowler_wickets += balls.wickets
		}
		if balls.runs > 0{
//			bowler_tally[bowler] = Tally{bowler_total: balls.runs}
			total_runs += balls.runs
		}
		if balls.stroke > 0{
			total_runs += balls.stroke
		}

		output += fmt.Sprintf("<td>%d/%d</td>", bowler_tally[bowler].bowler_wickets,  bowler_tally[bowler].bowler_total)
		FS := balls.facing
		BS := balls.batsman
		output += fmt.Sprintf("<td%v>%v</td><td>%v</td><td>%v</td><td>%d</td>",facingBatsman(balls.facing,balls.batsman),  batsmen[balls.batsman], batsman_scores(FS,BS,balls.stroke), batsman_scores(FS,BS,balls.runs), 0 )   // balls.batsman_runs
		BS = balls.batsman2
		output += fmt.Sprintf("<td%v>%v</td><td>%v</td><td>%v</td><td>%d</td>",facingBatsman(balls.facing,balls.batsman2),  batsmen[balls.batsman2],batsman_scores(FS,BS,balls.stroke), batsman_scores(FS,BS,balls.runs), 0 ) // balls.batsman2_runs
		output += "<td>"+balls.how_out+"</td>"
		output += fmt.Sprintf("<td>%d/%d</td>", total_wickets, total_runs)
		output += "</tr>"
		if balls.nb == 0 && balls.wide == 0{
			over_ball_number += 1
			if over_ball_number > balls_per_over{
			   over_ball_number = 1
				over_number += 1
				new_bowler = ""
			}
		}
	}

//	output += "<tr><td colspan=17>_</td></tr>"



	tempScores := getCollection("cricket")
	fmt.Print("\n\n")
	fmt.Print(tempScores)
	for _, loop_balls := range tempScores{
		balls := Scores{
			bowler:     strInt( loop_balls["bowler"]   ),
			nb:         strInt( loop_balls["nb"]       ),
			wide:       strInt( loop_balls["wide"]     ),
			runs:       strInt( loop_balls["runs"]     ),
			stroke:     strInt( loop_balls["stroke"]   ),
			batsman:    strInt( loop_balls["batsman"]  ),
			wickets:    strInt( loop_balls["wickets"]  ),
			batsman2:   strInt( loop_balls["batsman2"] ),
			facing:     strInt( loop_balls["facing"]   ),
			how_out:    fmt.Sprintf("%v",loop_balls["how_out"]),
		}
		fmt.Print("\n\n")
		fmt.Print(balls)


		facing = balls.facing
		batsman = balls.batsman
		batsman2 = balls.batsman2
		bowler = balls.bowler
		output += "<tr>"
		output += fmt.Sprintf("<td>%d.%d</td>", over_number, over_ball_number)
		if new_bowler != bowlers[balls.bowler]{
			new_bowler = bowlers[balls.bowler]
			output += fmt.Sprintf("\n\t<td rowspan=%d>%v</td>", overLength(over_ball_number,balls.bowler, scores), new_bowler)
		}
		output += "<td>"
		if balls.nb > 0{
			total_runs += balls.nb
			output += fmt.Sprintf("%d", balls.nb)
		}
		output += "</td><td>"
		if balls.wide > 0{
			output += fmt.Sprintf("%d", balls.wide)
			total_runs += balls.wide
		}
		output += "</td><td>"
		if balls.wide > 0{
			output += fmt.Sprintf("%d", balls.wide)
			total_runs += balls.wide
		}
		output += "</td><td>"
		if balls.wide > 0{
			output += fmt.Sprintf("%d", balls.wide)
			total_runs += balls.wide
		}

		if balls.wickets > 0{
			//			bowler_tally[bowler].bowler_wickets = balls.wickets
			total_wickets += balls.wickets
			//			bowler_tally[bowler].bowler_wickets += balls.wickets
		}
		if balls.runs > 0{
			//			bowler_tally[bowler] = Tally{bowler_total: balls.runs}
			total_runs += balls.runs
		}
		if balls.stroke > 0{
			total_runs += balls.stroke
		}

//		output += fmt.Sprintf("<td>%d/%d</td>", bowler_tally[bowler].bowler_wickets,  bowler_tally[bowler].bowler_total)
		FS := balls.facing
		BS := balls.batsman
//		output += fmt.Sprintf("<td%v>%v</td><td>%v</td><td>%v</td><td>%d</td>",facingBatsman(balls.facing,balls.batsman),  batsmen[balls.batsman], batsman_scores(FS,BS,balls.stroke), batsman_scores(FS,BS,balls.runs), 0 )   // balls.batsman_runs
		output += fmt.Sprintf("<td%v>%v</td><td>%v</td><td>%v fdsafdsafsdfsa</td>",facingBatsman(balls.facing,balls.batsman),  batsmen[balls.batsman], batsman_scores(FS,BS,balls.stroke), batsman_scores(FS,BS,balls.runs) )   // balls.batsman_runs
		BS = balls.batsman2
//		output += fmt.Sprintf("<td%v>%v</td><td>%v</td><td>%v</td><td>%d</td>",facingBatsman(balls.facing,balls.batsman2),  batsmen[balls.batsman2],batsman_scores(FS,BS,balls.stroke), batsman_scores(FS,BS,balls.runs), 0 ) // balls.batsman2_runs
		output += fmt.Sprintf("<td%v>%v</td><td>%v</td><td>%v</td>",facingBatsman(balls.facing,balls.batsman2),  batsmen[balls.batsman2],batsman_scores(FS,BS,balls.stroke), batsman_scores(FS,BS,balls.runs) ) // balls.batsman2_runs
		output += "<td>"+balls.how_out+"</td>"
		output += fmt.Sprintf("<td>%d/%d</td>", total_wickets, total_runs)
		output += "</tr>"
		if balls.nb == 0 && balls.wide == 0{
			over_ball_number += 1
			if over_ball_number > balls_per_over{
				over_ball_number = 1
				over_number += 1
				new_bowler = ""
			}
		}
	}

	output += "<tr><form method=POST action=save>"
//	output += fmt.Sprintf("<td><input class=short disabled value=%d.%d></td>", over_number, over_ball_number)
	output += fmt.Sprintf("<td>%d.%d</td>", over_number, over_ball_number)
	output += fmt.Sprintf("<td>%v</td>", select_bowler(bowlers, bowler, over_ball_number))
//	output += "<td><input name=nb type=checkbox></td>"
//	output += "<td><input name=wide type=checkbox></td>"
	output += "<td colspan=2><select><option></option><option>No Ball</option><option>Wide</option></select></td>"
	output += "<td colspan=2><select><option></option><option>Bye</option><option>Leg Bye</option></select></td>"
	output += "<td></td>"
	output += fmt.Sprintf("<td>%v</td>", select_batsman(batsmen, facing,  batsman, "batsman"))
	output += fmt.Sprintf("<td>%v</td>", select_strokes(facing, batsman))
	output += fmt.Sprintf("<td>%v</td>", input_runs(facing==batsman, "runs"))
	output += "<td></td>"
	output += fmt.Sprintf("<td>%v</td>", select_batsman(batsmen, facing, batsman2, "batsman2"))
	output += fmt.Sprintf("<td>%v</td>", select_strokes(facing, batsman2))
	output += fmt.Sprintf("<td>%v</td>", input_runs(facing==batsman2, "runs"))
	output += "<td></td>"
	output += "<td><input name=howout placeholder=\"How Out\"></td>"
	output += "<td><input type=submit value=Save></td>"
	output += "</form></tr>"
	output += "</table><br>"




	output += "<table border=1 style=text-align:right>"
	output += "<tr><td>Bowlers</td><td>Overs</td><td>Maiden</td><td>Wides</td><td>No Balls</td><td>Wickets</td><td>Runs</td></tr>"
	output += "<tr><td>A.Maveric</td><td>50</td><td>1</td><td>4</td><td>6</td><td>11</td><td>56</td></tr>"
	output += "</table><br>"

	output += "<table border=1 style=text-align:right>"
	output += "<tr><td>Batsmen</td><td>Total Runs</td><td>How Out</td><td>Bowler</td></tr>"
	output += "<tr><td>T.Hicky</td><td>56</td><td>LBW</td><td>A.Maveric</td></tr>"
	output += "</table></body></html>"
	fmt.Fprint(w, output)
}
func strInt(input interface{})int{
	i, _ := strconv.Atoi(fmt.Sprintf("%v",input))
	return i
}
func batsman_scores(facing, batsman, score int)string{
	if facing == batsman{
		return fmt.Sprintf("%v",score)
	}
	return ""
}
func overLength(ball_number, bowler_id int, scores map[int]Scores)int{
//	bowler_id := scores[ball_number].bowler
	output := 0
	for ball_index, balls := range scores{
		if ball_index >= ball_number && balls.bowler != bowler_id {
			return ball_index - ball_number +1
		}else{
			output = ball_index
		}
	}
	return output - ball_number +1
}


func facingBatsman(facing, batsman int)string{
	if facing == batsman{
		return " style=background:0f7"
	}
	return ""
}
func select_batsman(batsmen []string, facing, selected int, name string)string{
	output := "<select name="+name+" "+facingBatsman(facing, selected)+"><option selected></option>"
	for index, _ := range batsmen{
		selected_index := ""
		if index == selected+1{
//			selected_index = " selected"
		}
		output += fmt.Sprintf("<option value=%v%v>%v</option>",index, selected_index, batsmen[index])
	}
	output += "</select>"
	return output
}
func select_bowler(bowlers []string, selected, balls int)string{
	output := "<select name=bowler>"

	for index, _ := range bowlers{
		selected_index := ""
		if index == selected && balls == 6{
			selected_index = " selected"
		}
		output += fmt.Sprintf("<option value=%v%v>%v</option>",index, selected_index, bowlers[index])
	}
	output += "</select>"
	return output
}
func select_strokes(facing int, batsman int)string{
	output := ""
	if facing==batsman{
		output = fmt.Sprintf("<input type=hidden name=facing value=%d><select class=short name=stroke><option>0</option><option>4</option><option>6</option></select>", batsman)
	}
	return output
}
func input_runs(facing bool, name string)string{
	output := ""
	if facing{
		output = fmt.Sprintf("<input type=number name=%v value=0 min=0 max=10>", name)
	}
	return output
}
