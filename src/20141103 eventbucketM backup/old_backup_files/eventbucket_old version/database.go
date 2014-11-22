package main

import (
	"mgo"
	"mgo/bson"
	"fmt"
//	"reflect"
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

func getCollection(collection_name string) []map[string]interface{} {
	//TODO add in support to select only the columns required
	var result []map[string]interface{}
	checkErr(conn.C(collection_name).Find(nil).All(&result))
	return result
}

type Person struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"-"`
	name string `bson:"firstName" json:"firstName"`
}
type Hello struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"-"`
	name string `bson:"firstName" json:"firstName"`
}

//func getShit(collection_name, id string) map[string]interface{} {
//
//	searchResults := []Person{}
//	results := []Hello{}
//	conn.C(collection_name).Find(nil).All(&searchResults)
//
//	for index, row := range searchResults{
//		fmt.Printf("%v",index)
//		fmt.Printf("\t\t")
//		fmt.Printf("%v",row)
//		fmt.Printf("\n")
//
//
//		fmt.Printf("Loop: == \n")
//		conn.C(collection_name).FindId(bson.M{"_id":row}).All(&results)
//		for index, row := range results{
//			fmt.Printf("%v",index)
//			fmt.Printf("\t\t")
//			fmt.Printf("%v",row)
//			fmt.Printf("\n")
//		}
//		fmt.Printf("END Loop: == \n")
//	}





//	var result map[string]interface{}
//	var mapkey interface{}
//	temp := getCollection("event")
//	for key, row := range temp{
//		fmt.Printf("\n%v",key)
//		fmt.Printf("\n%v",row)
//		for column, value := range row{
//			fmt.Printf("\ncol:%v",column)
//			fmt.Printf("\nval:%v",value)
//		}
//		mapkey = row["_id"]
//		break
//	}
//
//	fmt.Printf("\nmapkey\t%v\n\n", mapkey)
//	fmt.Printf("\nmapkey2\t%v\n\n", bson.ObjectIdHex("4fe90ea6eaf9553e4d114a77"))
//	fmt.Printf("\nmapkey3\t%v\n\n",  reflect.TypeOf(mapkey))
//
//
//	var temp2 bson.ObjectId
//
//
//	temp2 = mapkey
//
//
//
////	err := conn.C(collection_name).FindId(mapkey).All(&result)
//	err := conn.C(collection_name).Find(bson.M{"_id":mapkey}).All(&result)
//	fmt.Printf("\n%v",err)
//	fmt.Printf("\n%v",result)


//	return map[string]interface{}{
//		"Id": "temps",
//	}
//}


func getDocument(collection_name, id string) []map[string]interface{} {
	//TODO add in support to select only the columns required
	var result []map[string]interface{}
	var resulted bson.M
//	conn.C(collection_name).FindId(bson.ObjectIdHex(id)).One(&result)

//	conn.C(collection_name).FindId(bson.ObjectIdHex(id)).Iter().All(&result)

	c := conn.C("event")

	err := conn.C(collection_name).Find(bson.M{"_id": "52a9a1ffff7f0c7aacacbe09"}).Iter().All(&result)

	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")

	err = conn.C(collection_name).Find(bson.M{"_id": `ObjectIdHex("52a9a1ffff7f0c7aacacbe09")`}).Iter().All(&result)

	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")

	err = conn.C("event").FindId(bson.ObjectIdHex("5309e14f8242ae6b6cbe0adc")).One(&result)

	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")

	err = c.FindId(bson.ObjectIdHex("5309e14f8242ae6b6cbe0adc")).Iter().All(&result)

	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")





//	query := collection.Find(bson.M{"_id": id})

	err = c.Find(bson.M{"name": "newEventName!!!"}).One(&resulted)
	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")

//collection.Find(bson.M{"_id": id})
	err = c.Find(bson.M{"_id": "5309e14f8242ae6b6cbe0adc"}).One(&result)
	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")


	err = c.FindId("5309e14f8242ae6b6cbe0adc").One(&result)
	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")


	err = c.FindId(bson.ObjectId("5309e14f8242ae6b6cbe0adc")).One(&result)
	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("%v",result)
	fmt.Printf("\n")


	var temper *mgo.Query
//	temper = c.Find(bson.M{"_id": "5309e14f8242ae6b6cbe0adc"})
	temper = c.Find(bson.M{"_id": bson.ObjectId("5309e14f8242ae6b6cbe0adc")})
//	temper = c.FindId(bson.ObjectId("5309e14f8242ae6b6cbe0adc"))

//	err = c.Find(bson.M{"name": "newEventName!!!"}).All(&result)



//	err = c.Find(bson.M{"_id": bson.ObjectId("5309e14f8242ae6b6cbe0adc")}).All(&result)
	err = c.FindId("23").All(&result)
	fmt.Printf("%v",temper)
	fmt.Printf("\n")
	fmt.Printf("%v",bson.ObjectId("5309e14f8242ae6b6cbe0adc"))
	fmt.Printf("\t\t")
	fmt.Printf("\n")
	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("\nFUCK this shit!!!!\n")
	fmt.Printf("%v",result)
	fmt.Printf("\n")


	fmt.Printf("\n")
	fmt.Printf("\n")


	fmt.Printf("%v\n", bson.M{"_id": bson.ObjectIdHex("5309e14f8242ae6b6cbe0adc")})
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("%v\n", bson.ObjectId("5309e14f8242ae6b6cbe0adc"))
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("Not again!!\n")

	//	err = c.FindId(bson.ObjectIdHex("5309e14f8242ae6b6cbe0adc")).All(&result)
	err = c.Find(bson.M{"_id": bson.ObjectId("5309e14f8242ae6b6cbe0adc")}).All(&result)
	fmt.Printf("%v",err)
	fmt.Printf("\t\t")
	fmt.Printf("\n")
	fmt.Printf("%v",result)
	fmt.Printf("\n")









	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("\n")
	return result
}

func InsertDoc(data interface{}, collection string) {
	if data != false {
		err := conn.C(collection).Insert(data)
		checkErr(err)
	}
}




//
//
//
//session, err := mgo.Dial("localhost")
//checkErr(err)
////	defer session.Close()
//// Optional. Switch the session to a monotonic behavior.
////	session.SetMode(mgo.Monotonic, true)
//session.SetMode(mgo.Eventual, true)//this is supposed to be faster
//return session.DB("eb")
//
//
//func getSession () *mgo.Session {
//    if mgoSession == nil {
//        var err error
//        mgoSession, err = mgo.Dial("localhost")
//        if err != nil {
//             panic(err) // no, not really
//        }
//    }
//    return mgoSession.Clone()
//}
//
//
//
//
//
// func withCollection(collection string, s func(*mgo.Collection) error) error {
//    session := getSession()
//    defer session.Close()
//    c := session.DB(databaseName).C(collection)
//    return s(c)
//}
////The withCollection() function takes the name of the collection, along with a function that expects the connection object to that collection, and can execute access functions on it.
////
////Here’s how the “Person” collection can be searched, using the withCollection() function:
//
//func SearchPerson (q interface{}, skip int, limit int) (searchResults []Person, searchErr string) {
//    searchErr     = ""
//    searchResults = []Person{}
//    query := func(c *mgo.Collection) error {
//        fn := c.Find(q).Skip(skip).Limit(limit).All(&searchResults)
//        if limit < 0 {
//            fn = c.Find(q).Skip(skip).All(&searchResults)
//        }
//        return fn
//    }
//    search := func() error {
//        return withCollection("person", query)
//    }
//    err := search()
//    if err != nil {
//        searchErr = "Database Error"
//    }
//    return
//}
