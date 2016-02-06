package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"gopkg.in/mgo.v2"
)

var conn *mgo.Database

func startDB() {
	const (
		mgoPort = "38888"
		mgoDial = "localhost:" + mgoPort
	)
	databasePath := os.Getenv("ProgramData") + subDir
	mkDir(databasePath)
	//TODO remove db args after comment tag once profiling is finished.
	cmd := exec.Command("mongod", "--dbpath", databasePath, "--port", mgoPort, "--nssize", "1", "--smallfiles", "--noscripting", "--nohttpinterface", "--quiet", "--notablescan" /*, "--slowms", "25", "--profile", "1"*/)
	//Combine standard output of command into EventBucket standard output source://stackoverflow.com/questions/8875038/redirect-stdout-pipe-of-child-process-in-golang
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	session, err := mgo.Dial(mgoDial)
	if err != nil {
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
		Update:    M{"$inc": M{schemaName: 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	_, err := conn.C(tblAutoInc).FindId(collectionName).Apply(change, &result)
	if err != nil {
		warn.Println(err)
		return "", fmt.Errorf("Unable to generate the next ID: '%v'", err)
	}

	//Convert integer to a alpha-numeric (0-9a-z / 36 base) string
	return strconv.FormatUint(uint64(result[schemaName].(int)), 36), nil
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

func getEvents() ([]Event, error) {
	var result []Event
	if conn != nil {
		err := conn.C(tblEvent).Find(nil).All(&result)
		if err != nil {
			warn.Println(err)
		}
		return result, err
	}
	return result, errors.New("Unable to get events")
}

func getClub(ID string) (Club, error) {
	var result Club
	if conn != nil {
		err := conn.C(tblClub).FindId(ID).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get club with ID: '" + ID + "'")
}

func getEvent(ID string) (Event, error) {
	var result Event
	if conn != nil {
		err := conn.C(tblEvent).FindId(ID).One(&result)
		return result, err
	}
	return result, errors.New("Unable to get event with ID: '" + ID + "'")
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
