package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

var (
	//Databse collection names
	tblClub    = []byte("C")
	tblEvent   = []byte("E")
	tblShooter = []byte("S")
)

/*
	event, err := getEvent3(2)
	//event, err := getEvent(2)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%T\n\n%v, %v", event, len(event.Shooters), len(event.Shooters[0].Scores))
*/

/*for _, shooter := range event.Shooters{
	log.Println(shooter.ID, shooter.FirstName, shooter.Surname)
	log.Println(len(shooter.Scores))
	for _, score := range shooter.Scores{
		log.Println("\t", score.RangeID, score.Total, score.Centres, score.Shots)
	}
}*/

func getDocument(collection []byte, ID string, result interface{}) error {
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(collection)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", collection)
		}
		byteID, err := B36toBy(ID)
		if err != nil {
			return err
		}
		trace.Println(ID, []byte(ID), byteID)

		document := bucket.Get(byteID)
		err = json.Unmarshal(document, result)
		if err != nil {
			warn.Printf("Query %p document unmarshaling failed: %#v\n", document, err)
		}
		warn.Println(string(document))
		return err
	})
	return err
}

func getEvent(ID string) (Event, error) {
	var event Event
	//	b, err := B36toUint(ID)
	//	if err != nil {
	//		return event, err
	//	}
	return event, getDocument(tblEvent, ID, &event)
}

func getClub(ID string) (Club, error) {
	var club Club
	//	b, err := B36toUint(ID)
	//	if err != nil {
	//		return club, err
	//	}
	return club, getDocument(tblClub, ID, &club)
}

func getShooter(ID string) (Shooter, error) {
	var shooter Shooter
	//	b, err := B36toUint(ID)
	//	if err != nil {
	//		return shooter, err
	//	}
	return shooter, getDocument(tblShooter, ID, &shooter)
}

/*
func insertEvent(db *bolt.DB, u Event) uint64 {
	var eventID uint64

	// store some data
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblEvent)
		if err != nil {
			return err
		}
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		eventID, _ = bucket.NextSequence()

		//log.Println("eventID", eventID)

		u.Name = string(eventID)
		u.ID = toB36(eventID) //uint64(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return bucket.Put(itob(eventID), buf)
	})

	if err != nil {
		log.Fatal(err)
	}
	return eventID
}*/
/*
func insertE(db *bolt.DB, u Event) {
	//var eventID uint64

	// store some data
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblEvent)
		if err != nil {
			return err
		}
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		eventID, _ := bucket.NextSequence()

		//log.Println("eventID", eventID)

		u.ID = toB36(eventID) //uint64(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return bucket.Put(itob(eventID), buf)
	})

	if err != nil {
		log.Fatal(err)
	}
}*/

func insertEvent(event Event) (string, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblEvent)
		if err != nil {
			return err
		}
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		event.ID, _ = bucket.NextSequence()
		//		eventID = toB36(sequence)
		//		event.ID = eventID

		// Marshal user data into bytes.
		buf, err := json.Marshal(event)
		if err != nil {
			return err
		}

		//		return bucket.Put(itob(eventID), buf)
		return bucket.Put(itob(event.ID), buf)
	})

	//	if err != nil {
	//		warn.Fatal(err)
	//	}
	return toB36(event.ID), err
}

func insertClub(club Club) (string, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblClub)
		if err != nil {
			return err
		}
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		club.ID, _ = bucket.NextSequence()
		//		clubID = toB36(sequence)
		//		club.ID = clubID

		// Marshal user data into bytes.
		buf, err := json.Marshal(club)
		if err != nil {
			return err
		}

		//		return bucket.Put(itob(clubID), buf)
		return bucket.Put(itob(club.ID), buf)
	})

	//	if err != nil {
	//		warn.Fatal(err)
	//	}
	return toB36(club.ID), err
}

/*
func getDocument_backup(collection []byte, ID uint64) ([]byte, error) {
	var document []byte
	// retrieve the data
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(collection)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", collection)
		}
		document = bucket.Get(itob(ID))
		return nil
	})
	return document, err
}

func getEvent2(eventID uint64)(Event, error){
	var event Event
	document, err := getDocument(tblEvent, eventID)
	if err != nil{
		return event, err
	}
	err = json.Unmarshal(document, &event)
	return event, err
}

func getClub(clubID uint64)(Club, error){
	var club Club
	document, err := getDocument(tblClub, clubID)
	if err != nil{
		return club, err
	}
	err = json.Unmarshal(document, &club)
	return club, err
}*/
/*
func getEvent(eventID uint64) (Event, error) {
	var event Event
	var err error
	var value []byte
	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tblEvent)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", tblEvent)
		}

		value = bucket.Get(itob(eventID))

		//log.Println("out", string(val))

		return nil

	})
	if err != nil {
		return event, err
	}
	//log.Println(string(value))
	err = json.Unmarshal(value, &event)
	return event, err
}

func getE(db *bolt.DB, collection []byte, id uint64) []byte {
	var val []byte
	// retrieve the data
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(collection)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", collection)
		}

		val = bucket.Get(itob(id))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return val
}*/

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// stob returns an 8-byte big endian representation of v.
func stob(v string) []byte {
	//	b := make([]byte, 8)
	//	binary.BigEndian.PutUint64(b, v)
	return []byte(v)
}

func insertDoc(collectionName []byte, document interface{}) error {
	/*err := conn.C(collectionName).Insert(document)
	if err != nil {
		warn.Println(err)
	}
	return err*/
	return nil
}

func upsertDoc(collectionName []byte, ID string, document interface{}) error {
	/*_, err := conn.C(collectionName).UpsertId(ID, document)
	if err != nil {
		warn.Println(err)
	}
	return err*/
	return nil
}

func updateDoc(collectionName []byte, ID string, document interface{}) error {
	/*err := conn.C(collectionName).UpdateId(ID, document)
	if err != nil {
		warn.Println(err)
	}
	return err*/
	return nil
}

func getNextID(collectionName []byte) (string, error) {
	/*var result AutoID
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
	return strconv.FormatUint(result.Name, 36), nil*/
	return "", nil
}

func eventAddRange(eventID string, newRange Range) error {
	/*change := mgo.Change{
		Update: M{
			"$push": M{schemaRange: newRange},
		},
		Upsert: true,
		//		ReturnNew: true,
	}
	//	returned := Event{}
	//	conn.C(tblEvent).FindId(eventID).Apply(change, &returned)
	_, err := conn.C(tblEvent).FindId(eventID).Apply(change, &Event{})
	//	for rangeID, rangeData := range returned.Ranges {
	//		if rangeData.Name == newRange.Name && rangeData.Aggregate == newRange.Aggregate && rangeData.ScoreBoard == newRange.ScoreBoard && rangeData.Locked == newRange.Locked && rangeData.Hidden == newRange.Hidden {
	//			TODO this if check is really hacky!!!
	//			return rangeID, returned
	//		}
	//	}
	//	return -1, returned
	return err*/
	return nil
}

func getClubs() ([]Club, error) {
	var clubs []Club
	var club Club
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblClub)
		if b == nil {
			//Club Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			if json.Unmarshal(value, &club) == nil {
				clubs = append(clubs, club)
			}
			return nil
		})
	})
	return clubs, err
}

func getEvents() ([]Event, error) {
	var events []Event
	var event Event
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblEvent)
		if b == nil {
			//Event Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			if json.Unmarshal(value, &event) == nil {
				events = append(events, event)
			}
			return nil
		})
	})
	return events, err
}

func updateAll(collectionName []byte, query, update M) {
	/*_, err := conn.C(collectionName).UpdateAll(query, update)
	if err != nil {
		warn.Println(err)
	}*/
}

func collectionQty(collectionName []byte) int {
	/*qty, err := conn.C(collectionName).Count()
	if err != nil {
		warn.Println(err)
	}
	return qty*/
	return 0
}

func hasDefaultClub() bool {
	/*if conn != nil {
		qty, err := conn.C(tblClub).Find(M{schemaIsDefault: true}).Count()
		return qty > 0 && err == nil
	}*/
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

func eventShooterInsertDB(eventID string, shooter EventShooter) error {
	/*insert := M{
		schemaShooter: []EventShooter{shooter},
	}
	//If shooter is Match Reserve, duplicate them in the Match Open category
	increment := 1
	if shooter.Grade == 8 {
		increment = 2
		duplicateShooter := shooter
		duplicateShooter.Grade = 7
		duplicateShooter.Hidden = true
		insert[schemaShooter] = []EventShooter{shooter, duplicateShooter}
	}
	change := mgo.Change{
		Update: M{
			"$pushAll": insert,
			"$inc": M{
				dot(schemaAutoInc, schemaShooter): increment,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var event Event
	conn.C(tblEvent).FindId(eventID).Apply(change, &event)

	if increment == 2 {
		change = mgo.Change{
			Update: M{
				"$set": M{
					dot(schemaShooter, event.AutoInc.Shooter-2, "i"): event.AutoInc.Shooter - 2,
					dot(schemaShooter, event.AutoInc.Shooter-2, "l"): event.AutoInc.Shooter - 1,
					dot(schemaShooter, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter - 1,
					dot(schemaShooter, event.AutoInc.Shooter-1, "l"): event.AutoInc.Shooter - 2,
				},
			},
		}
	} else {
		change = mgo.Change{
			Update: M{
				"$set": M{
					dot(schemaShooter, event.AutoInc.Shooter-1, "i"): event.AutoInc.Shooter - 1,
				},
			},
		}
	}
	conn.C(tblEvent).FindId(eventID).Apply(change, &event)*/
	return nil
}

func B36toBy(a string) ([]byte, error) {
	v, err := strconv.ParseUint(a, 36, 64)
	if err != nil {
		return []byte{}, err
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b, nil
}
