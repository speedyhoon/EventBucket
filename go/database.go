package main

import (
	"errors"
	"fmt"
	"os/exec"

	"gopkg.in/mgo.v2"
)

var conn *mgo.Database

func startDB() {
	const (
		mgoPort = "38888"
		mgoDial = "localhost:" + mgoPort
	)
	mkDir(databasePath)
	cmd := exec.Command("mongod", "--dbpath", databasePath, "--port", mgoPort, "--nssize", "1", "--smallfiles", "--noscripting", "--nohttpinterface")
	//TODO output mongodb errors/logs to stdOut
	/*stdout, err := cmd.StdoutPipe()
	if err != nil {
		warning.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		warning.Println(err)
	}
	err = cmd.Start()*/
	cmd.Start()
	/*if err != nil {
		dump("exporting!")
		dump(stdout)
		dump(stderr)
		return
	}
	fmt.Println("Result: Err")
	dump(stderr)
	fmt.Println("Result: stdn out")
	dump(stdout)*/

	session, err := mgo.Dial(mgoDial)
	if err != nil {
		//TODO it would be better to output the mongodb connection error
		warn.Println("The database service is not available.", err)
		return
	}
	session.SetMode(mgo.Strong, false) //false = the consistency guarantees won't be reset

	//Can't directly change global variables in a go routine, so call an external function.
	dbConnection(session)
}

func dbConnection(session *mgo.Session) {
	conn = session.DB("local")
}

func upsertDoc(collectionName string, ID string, document interface{}) error {
	_, err := conn.C(collectionName).UpsertId(ID, document)
	if err != nil {
		warn.Println(err)
	}
	return err
}

func updateDoc(collectionName string, ID string, document interface{}) error {
	err := conn.C(collectionName).UpdateId(ID, document)
	if err != nil {
		warn.Println(err)
	}
	return err
}

func getNextID(collectionName string) (string, error) {
	var result M
	if conn == nil {
		return "", errors.New("Unable to generate the next ID. No database connection.")
	}

	change := mgo.Change{
		//		Update:    M{"$inc": M{schemaCounter: 1}},
		Update:    M{"$inc": M{schemaName: 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	_, err := conn.C(tblAutoInc).FindId(collectionName).Apply(change, &result)
	if err != nil {
		warn.Println(err)
		return "", fmt.Errorf("Unable to generate the next ID: '%v'", err)
	}
	//	return idSuffix(result[schemaCounter].(int))
	return idSuffix(result[schemaName].(int))
}

func getClubs() ([]Club, error) {
	var result []Club
	if conn != nil {
		err := conn.C(tblClub).Find(nil).All(&result)
		if err != nil {
			warn.Println(err)
		}
		return result, err
	}
	return result, errors.New("Unable to get clubs")
}

func getClub(ID string) (Club, error) {
	var result Club
	if conn != nil {
		err := conn.C(tblClub).FindId(ID).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get club with ID: '" + ID + "'")
}

func updateAll(collectionName string, query, update M) {
	_, err := conn.C(collectionName).UpdateAll(query, update)
	if err != nil {
		warn.Println(err)
	}
}

func collectionQty(collectionName string) int {
	qty, err := conn.C(collectionName).Count()
	if err != nil {
		warn.Println(err)
	}
	return qty
}

func hasDefaultClub() bool {
	if conn != nil {
		qty, err := conn.C(tblClub).Find(M{schemaIsDefault: true}).Count()
		return qty > 0 && err == nil
	}
	return false
}

/*func getDefaultClub() (Club, error) {
	var result Club
	if conn != nil {
		err := conn.C(tblClub).Find(M{schemaIsDefault: true}).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with ID: '" + ID + "'")
}*/
