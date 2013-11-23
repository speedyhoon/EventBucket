package main
import(
	"fmt"
	"mgo"
	"mgo/bson"
)

type Person struct{
	Name string
	Phone string
}

func main(){
	session, err := mgo.Dial("localhost")
//	session, err := mgo.Dial("server1.example.com,server2.example.com")
	checkErr(err)
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
//	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true)//this is supposed to be faster


	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
						&Person{"Cla", "+55 53 8402 8510"})

	checkErr(err)

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	checkErr(err)
	fmt.Println("Phone:", result.Phone)
}
func checkErr(err error){
	if err != nil {
		panic(err)
	}
}
