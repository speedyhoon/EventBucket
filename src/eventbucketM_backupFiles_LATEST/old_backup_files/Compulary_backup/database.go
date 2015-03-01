package main

import (
	"mgo"
//	"math"
	"fmt"
//	"reflect"


	"net/http"
	"mgo/bson"

)

func DB() *mgo.Database {
	session, err := mgo.Dial("localhost")
	checkErr(err)
	//	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true)//this is supposed to be faster
	return session.DB("eb")
}

func testNewCreationStructSaving(w http.ResponseWriter, r *http.Request){
	type Person struct {
		Id         bson.ObjectId   `bson:"_id,omitempty" json:"-"`
		FirstName  string          `bson:"F" json:"firstName"`
		MiddleName string          `bson:"M,omitempty" json:"middleName,omitempty"`
		LastName   string          `bson:"L" json:"lastName"`
//		Inserted   time.Time       `bson:"I" json:"-"`
	}

	data := Person{
		Id: "next",
		FirstName: "Leesa",
		MiddleName: "Loves",
		LastName: "Cam",
	}

	checkErr(conn.C("dummyTest").Insert(data))
}

func customSchema(w http.ResponseWriter, r *http.Request){

	type Score struct {
		Total			int		`bson:"t,omitempty"`
		Centers		int		`bson:"c,omitempty"`
		Shots			string	`bson:"s,omitempty"`
		Countback	uint64	`bson:"b,omitempty"`
		Xs				int		`bson:"x,omitempty"`
	}
	type Shooter struct {
		Id		bson.ObjectId		`bson:"_id,omitempty"`
		FirstName	string		`bson:"f,omitempty"`
		Surname		string		`bson:"s,omitempty"`
		MiddleName	string		`bson:"m,omitempty"`
		Age			string		`bson:"a,omitempty"`
		Club			int			`bson:"C,omitempty"`
		Class			string		`bson:"c,omitempty"`
		Grade			string		`bson:"g,omitempty"`
		Score	map[string]Score	`bson:"r,omitempty"`
	}
	type Club struct {
		Id		bson.ObjectId	`bson:"_id,omitempty"`
		Name			string	`bson:"n,omitempty"`
		Url			string	`bson:"u,omitempty"`
		Distance		string	`bson:"d,omitempty"`
		Unit			string	`bson:"u,omitempty"`
		Latitude		string	`bson:"l,omitempty"`
		Longitude	string	`bson:"o,omitempty"`
	}

	type AutoInc struct {
		Range	int	`bson:"R,omitempty"`
	}
	type Range struct{
		Name	string	`bson:"n,omitempty"`
	}
	type Event struct {
		Id			bson.ObjectId			`bson:"_id,omitempty"`
		DateTime	string					`bson:"d,omitempty"`
		AutoInc	AutoInc					`bson:"U,omitempty"`
		Range		map[string]Range		`bson:"R,omitempty"`
		Shooters	map[string]Shooter	`bson:"S,omitempty"`
		Club		int						`bson:"C,omitempty"`
	}

	data := Event{
		Id:	"hell99",
		DateTime:	"201405231300",
		Club:		44,
		Range:	map[string]Range{
			"0":	Range{
				Name:	"300 yards",
			},
			"1":	Range{
				Name:	"600 yards",
			},
			"2":	Range{
				Name:	"1200 yards",
			},
			"3":	Range{
				Name:	"700 yards",
			},
			"4":	Range{
				Name:	"Total",
			},
		},
		Shooters:	map[string]Shooter{
			"0": Shooter{
				Id:	"fd",
				FirstName: "Cam",
				Surname: "W",
				MiddleName: "F",
				Age: "",
				Club: 1,
				Class: "target",
				Grade: "B",
				Score: map[string]Score{
					"0":	Score{
						Total: 40,
						Centers: 5,
					},
					"1":	Score{
						Total: 40,
						Centers: 4,
					},
					"2":	Score{
						Total: 40,
						Centers: 3,
					},
					"3":	Score{
						Total: 40,
						Centers: 6,
					},
					"4":	Score{
						Total: 160,
						Centers: 18,
					},
				},
			},
			"1": Shooter{
				Id:	"sa",
				FirstName: "Leesa",
				Surname: "N",
				MiddleName: "M",
				Age: "U21",
				Club: 2,
				Class: "target",
				Grade: "B",
				Score: map[string]Score{
					"0":	Score{
						Total: 40,
						Centers: 5,
						Countback: 5,
					},
					"1":	Score{
						Total: 40,
						Centers: 4,
					},
					"2":	Score{
						Total: 40,
						Centers: 3,
					},
					"3":	Score{
						Total: 40,
						Centers: 7,
//						Countback: 9999999999999999999,
						Countback: 999999999999999999,
					},
					"4":	Score{
						Total: 160,
						Centers: 19,
						Countback: 4545454545,
					},
				},
			},
		},
		AutoInc:	AutoInc{
			Range:	88,
		},
	}

	checkErr(conn.C("newEvent").Insert(data))
}


func getCollection(collection_name string) []map[string]interface{} {
	//TODO add in support to select only the columns required
	var result []map[string]interface{}
	checkErr(conn.C(collection_name).Find(nil).All(&result))
	return result
}

func getClubs()[]Club{
	var result []Club
	checkErr(conn.C(TBLclub).Find(nil).All(&result))
	return result
}
func getClub(id string)Club{
	var result Club
	checkErr(conn.C(TBLclub).FindId(id).One(&result))
	return result
}

func getEvents()[]Event{
	var result []Event
	checkErr(conn.C(TBLevent).Find(nil).All(&result))
	return result
}

func InsertDoc(data map[string]interface{}, collection_name string)string{

	if collection_name == "club"{
		data[schemaAUTOINC] = map[string]int{
			schemaRANGE: 0,
		}
	}

	fmt.Printf(fmt.Sprintf("%v",data))


//	if data != false {
		data["_id"] = getNextId(collection_name)
//		data = [...]map[string]interface{}{"_id":getNextId(collection_name)}
		err := conn.C(collection_name).Insert(data)
		checkErr(err)
//	}
	return data["_id"].(string)
}

func InsertDoc3(data Event, collection_name string) {
//	fmt.Printf(fmt.Sprintf("%v",data))
//	data["_id"] = getNextId(collection_name)
	err := conn.C(collection_name).Insert(data)
	checkErr(err)
}

func getDocument(collection_name, id string)map[string]interface{}{
	result := make(map[string]interface{})
//	conn.C(collection_name).FindId(map[string]interface{}{"_id": id  }).Apply(change, &result)
//	checkErr(conn.C(collection_name).Find(map[string]interface{}{"_id":id}).One(&result))


//	checkErr(conn.C(collection_name).Find(map[string]interface{}{"_id":"5317198d81a27b3006872b0f"}).One(&result))
	checkErr(conn.C(collection_name).FindId(id).One(&result))
	return result
}

func getDoc_by_id(collection_name, id string)map[string]interface{}{
	result := make(map[string]interface{})
//	conn.C(collection_name).FindId(map[string]interface{}{"_id": id  }).Apply(change, &result)
//	checkErr(conn.C(collection_name).Find(map[string]interface{}{"_id":id}).One(&result))
	checkErr(conn.C(collection_name).FindId(id).One(&result))
	return result
}

func getEvent(id string)Event{
	var result Event
//	conn.C(collection_name).FindId(map[string]interface{}{"_id": id  }).Apply(change, &result)
//	checkErr(conn.C("event").Find(map[string]interface{}{"_id":"5317198d81a27b3006872b0f"}).One(&result))

//	checkErr(conn.C("event").Find(map[string]interface{}{"_id":"ahe"}).One(&result))
	checkErr(conn.C("event").FindId(id).One(&result))
	return result
}

func appendRange(id, range_name, agg string){
	var result map[string]interface{}
//	`db.event.update({"_id":"m"}, {$set:{"S.re": {"F":"new","M":"shooter","C":123} }})`

	new_range := map[string]interface{}{schemaNAME: range_name}
	if agg != ""{
		new_range[schemaAGG] = agg
	}

	result = getDoc_by_id("event", id)
//	conn.C("event").Find(map[string]interface{}{"_id": id  }).One(&result)

//	dump(result)


//	dump(result[fmt.Sprintf("%v.%v",schemaAUTOINC,schemaRANGE)])
//	result[schemaAUTOINC].(map[string]interface{})[schemaRANGE]
//	dump(result[schemaAUTOINC].(map[string]interface{})[schemaRANGE])

	//TODO change this to get the range id from the club settings rather than make a new range all the time
//	new_range[schemaID] = result[schemaAUTOINC].(map[string]interface{})[schemaRANGE]
	new_index := fmt.Sprintf("%v.%v", schemaRANGE, result[schemaAUTOINC].(map[string]interface{})[schemaRANGE])

	change := mgo.Change{
		Update: map[string]interface{}{
//			"$push": map[string]interface{}{schemaRANGE: new_range },
			"$set": map[string]interface{}{new_index: new_range},
			"$inc": map[string]interface{}{fmt.Sprintf("%v.%v",schemaAUTOINC,schemaRANGE): 1},
		},
		ReturnNew: true,
	}
	conn.C("event").Find(map[string]interface{}{"_id": id  }).Apply(change, &result)
}

func appendShooter(event_id string, shooter_obj map[string]interface{}){
	var result map[string]interface{}
	change := mgo.Change{
		Update: map[string]interface{}{
			"$set": map[string]interface{}{fmt.Sprintf("%v.%v", schemaSHOOTER, shooter_obj[schemaID]): shooter_obj },
//			"$push": map[string]interface{}{schemaSHOOTER: shooter_obj },
//			"$inc": map[string]interface{}{fmt.Sprintf("%v.%v",schemaAUTOINC,schemaSHOOTER): 1},
		},
		ReturnNew: true,
	}
	conn.C("event").FindId(event_id).Apply(change, &result)
//	dump("rewq###\n")
//	dump(result)
//	dump("^^^rewq\n")
//	conn.C("event").UpdateId(event_id, map[string]interface{}{schemaSHOOTER})
}

//func getNexId(w http.ResponseWriter, r *http.Request){
func getNextId(collection_name string)string{
	var result map[string]interface{}
//	checkErr(conn.C(collection_name).Find(map[string]interface{}{"_id":collection_name}).All(&result))
//	current_number := result["inc"]

//	collection_name := "


	change := mgo.Change{
		Update: map[string]interface{}{"$inc": map[string]interface{}{"inc": 1}},
		ReturnNew: true,
//		Select: map[string]interface{}{"inc": 1},
	}
//	info, err := conn.C(collection_name).Find(map[string]interface{}{"_id":"event"}).Apply(change, &result)
	conn.C("autoinc").Find(map[string]interface{}{"_id":collection_name}).Apply(change, &result)
//	err := conn.C(collection_name).Find(map[string]interface{}{"_id":collection_name}) //.Apply(change, &result)


//	log10(345) / log10(65)

//	rego := 345
//	base := 65

//	for index, value := range info{
//		fmt.Printf(index)
//		fmt.Printf(value)
//	}


//	chars := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","~","!","@","#","$","%","^","&","*","(",")","_","-","+","{","}","[","]","|","\\",":",";",",",".","/","?","`"}
	chars := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","~","!","*","(",")","_","-","."}

	newId := chars[ result["inc"].(int) ]


//	fmt.Println(info)


//	fmt.Printf(fmt.Sprintf("%v",info.N))
	fmt.Println("<<newId\n\n")

	fmt.Println(newId)
//	fmt.Println("<<info\n\n")
//
//
//
//	checkErr(err)
//	fmt.Println(result["inc"])
//	fmt.Println("<<result\n\n")
//	fmt.Println("\n\n\n")
//
//
//	fmt.Println("geewiz::",genId(result["inc"].(int)))




	return newId
}


func club_insert_mound(club_id string, new_mound Mound){
	var club Club
	change := mgo.Change{
		Update: map[string]interface{}{
			"$inc": map[string]interface{}{fmt.Sprintf("%v.%v", schemaAUTOINC, schemaMOUND): 1},
		},
		ReturnNew: true,
	}
	conn.C(TBLclub).FindId(club_id).Apply(change, &club)
	change = mgo.Change{
		Update: map[string]interface{}{
			"$set": map[string]interface{}{fmt.Sprintf("%v.%v", schemaMOUND, club.AutoInc.Mound): new_mound},
		},
		ReturnNew: true,
	}
	conn.C(TBLclub).FindId(club_id).Apply(change, &club)
}
