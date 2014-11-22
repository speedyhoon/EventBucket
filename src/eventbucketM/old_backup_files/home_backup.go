

func home(w http.ResponseWriter, r *http.Request) {
	templator("home", home_temp, homeData(getCollection("event")), w)

}
func homeData(event []map[string]interface{}) map[string]interface{} {

	//	datetime := "201302021252"

	//	test, err := time.Parse("2006/02/01 15:04", datetime)
	//	test, err := time.Parse("20060201", fmt.Sprintf("%v", time.Now()))


	this := fmt.Sprintf("%v", time.Now())

	all := []HomeCalendar{}
	for _, row := range event {
		all = append(all, HomeCalendar{
				Name:      row["name"].(string),
				Club:      row["clubId"].(string),
				Day:      row["datetime"].(string),
				Date:     tryTime(row["datetime"]),
				//				Date:      fmt.Sprintf("%v", test),
				//				Date:      fmt.Sprintf("%v", test),
				//				Date:      row["datetime"].(string),
				//				Month:  row["datetime"].(string),
				Month:  this,
				Time:      row["datetime"].(string),
			})
	}

	return map[string]interface{}{
		"Source": all,
		"PageName": "Calendar",
		"Menu": "Menu is not built yet",
	}

}
func tryTime(datetime interface{}) string {
	test, err := time.Parse("200601021504", fmt.Sprintf("%v", datetime))
	//	test, err := time.Parse("2006-01-02 15:04", "2011-01-19 22:15")
	checkError(err)
	//	return test.Month
	return fmt.Sprintf("%v", test.Month)
}

const (
	stdLongMonth      = "January"
	stdMonth          = "Jan"
	stdNumMonth       = "1"
	stdZeroMonth      = "01"
	stdLongWeekDay    = "Monday"
	stdWeekDay        = "Mon"
	stdDay            = "2"
	stdUnderDay       = "_2"
	stdZeroDay        = "02"
	stdHour           = "15"
	stdHour12         = "3"
	stdZeroHour12     = "03"
	stdMinute         = "4"
	stdZeroMinute     = "04"
	stdSecond         = "5"
	stdZeroSecond     = "05"
	stdLongYear       = "2006"
	stdYear           = "06"
	stdPM             = "PM"
	stdpm             = "pm"
	stdTZ             = "MST"
	stdISO8601TZ      = "Z0700"  // prints Z for UTC
	stdISO8601ColonTZ = "Z07:00" // prints Z for UTC
	stdNumTZ          = "-0700"  // always numeric
	stdNumShortTZ     = "-07"    // always numeric
	stdNumColonTZ     = "-07:00" // always numeric
)

const home_temp = `<h1>{{ .PageName}}</h1>
{{if .Menu}} {{XTC .Menu}}{{end}}
{{with .Source}}
    {{range .}}
        Name is {{.Name}} <br>
        Club is {{.Club}} <br>
        Day: {{.Day}}<br>
        Date {{.Date}}<br>
        Time {{.Time}}<br>
        Month {{.Month}}<br>
    {{end}}
{{end}}
================================
`

//<a href=organisers>Organisers</a>
/*
Thursday
3rd
May
EditTest Event
Test Club
3.00am
*/

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

type HomeCalendar struct {
	Name, Club, Day, Date, Month, Time string
}







































//
//
//func (f FooId) GetBSON() (interface{}, error) {
//	return bson.ObjectId(f), nil
//}

func tempTry(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Try Executed\n")
	//	tryThis := valid8(r.Form, organisers_clubForm())


	//	temp3 := []interface{}{bson.M{"_id": "club"}, bson.M{"next":1}	}

	//	531293d28242ae6b6cbe0ade

	var _ = bson.M{}



//	newDoc := map[int]interface{}{
//		0:"name",
//		1: "url",
//		2: map[int]interface{}{ //shooters
//			0:"i like this!",
//		},
//	}


var result []map[string]interface{}
var result2 []map[string]interface{}

//	result = getCollection("club")





// 	conn.C("club").Find(nil).All(&result)     //select * from club
//	conn.C("club").Find(bson.M{"name":"hello"}).All(&result)     //select * from club where name = hello
conn.C("club").Find(bson.M{ "name":"hello"}).All(&result)     //select * from club where name = hello

//	conn.C("club").FindId("531293d28242ae6b6cbe0ade").One(&result)



//	conn.C("club").Insert(bson.M{"_id":autoInc("club"), "name":"two"})
//	conn.C("club").Insert(bson.M{"_id":autoInc("club"), "name":"three"})



//	checkErr(conn.C("autoinc").Find(bson.M{"_id": "club", "next":1}).One(&result))

for index, row := range result {
fmt.Printf("%v", index)
fmt.Printf("\t\t")
fmt.Printf("%v", row)
fmt.Printf("\n")
fmt.Printf("\n")
//		fmt.Printf("%v", bson.ObjectId( fmt.Sprintf("%v", row["_id"])   ).String() )
conn.C("club").FindId(row["_id"]).All(&result2)
fmt.Printf("\n")
fmt.Printf("%v", row["_id"])
fmt.Printf("\n")
fmt.Printf("%v", row["name"])
fmt.Printf("\n")
fmt.Printf("%v", row["url"])
fmt.Printf("\n")


}


fmt.Printf("\nRESULT TWO::::::::")

for index, row := range result2 {
fmt.Printf("%v", index)
fmt.Printf("\t\t")
fmt.Printf("%v", row)
fmt.Printf("\n")
fmt.Printf("\n")
//		fmt.Printf("%v", bson.ObjectId( fmt.Sprintf("%v", row["_id"])   ).String() )
//		conn.C("club").FindId(row["_id"]).All(&result2)
fmt.Printf("\n")
fmt.Printf("%v", row["_id"])
fmt.Printf("\n")
fmt.Printf("%v", row["name"])
fmt.Printf("\n")
fmt.Printf("%v", row["url"])
fmt.Printf("\n")


}

//	fmt.Printf("\nRESULT THEE::::::::")



//	newClub := Club{
//		name: "hello",
//		url: "my_url",
//		rego: "cF",
//	}
//
//	face := map[int]interface{}{
//		0: newClub,
//	}
//
//	InsertStruct(face, "club")
//
//
//	getItBackOut := Club{}
//
//
////	var result3 []map[interface{}]interface{}
//
//	//	result = getCollection("club")
////	conn.C("club").Find(bson.M{"0":"name"}).All(&getItBackOut)
//	conn.C("club").Find(bson.M{"url": "my_url"}).One(&getItBackOut)
//
//
////	for index, row := range getItBackOut {
////		fmt.Printf("%v", index)
////		fmt.Printf("\t\t")
//		fmt.Printf("%v", getItBackOut)
//		fmt.Printf("\n")
//		fmt.Printf("\n")
//		//		fmt.Printf("%v", bson.ObjectId( fmt.Sprintf("%v", row["_id"])   ).String() )
//		//		conn.C("club").FindId(row["_id"]).All(&result2)
//		fmt.Printf("\n")
////		fmt.Printf("%v", row["_id"])
//		fmt.Printf("\n")
//		fmt.Printf("%v", getItBackOut.name)
//		fmt.Printf("\n")
//		fmt.Printf("%v", getItBackOut.url)
//		fmt.Printf("\n")


//	}



//	result := Person{}
//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("Phone:", result.Phone)




//
//
//	var result3 []map[interface{}]interface{}
//
//	//	result = getCollection("club")
//	conn.C("club").Find(bson.M{"0":"name"}).All(&result3)
//
//
//	for index, row := range result3 {
//		fmt.Printf("%v", index)
//		fmt.Printf("\t\t")
//		fmt.Printf("%v", row)
//		fmt.Printf("\n")
//		fmt.Printf("\n")
//		//		fmt.Printf("%v", bson.ObjectId( fmt.Sprintf("%v", row["_id"])   ).String() )
//		//		conn.C("club").FindId(row["_id"]).All(&result2)
//		fmt.Printf("\n")
//		fmt.Printf("%v", row["_id"])
//		fmt.Printf("\n")
//		fmt.Printf("%v", row["name"])
//		fmt.Printf("\n")
//		fmt.Printf("%v", row["url"])
//		fmt.Printf("\n")
//	}









//
//
//
//	function counter(name) {
//	var ret = db.counters.findAndModify({query:{_id:name}, update:{$inc : {next:1}}, "new":true, upsert:true});
//// ret == { "_id" : "users", "next" : 1 }
//return ret.next;
//}
//
//db.users.insert({_id:counter("users"), name:"Sarah C."}) // _id : 1
//db.users.insert({_id:counter("users"), name:"Bob D."}) // _id : 2


//	tryThis["_id"] = result["next"] + 1

}
