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

		document := bucket.Get(byteID)
		err = json.Unmarshal(document, result)
		if err != nil {
			warn.Printf("Query %p document unmarshaling failed: %#v\n", document, err)
		}
		return err
	})
	return err
}

func getEvent(ID string) (Event, error) {
	var event Event
	return event, getDocument(tblEvent, ID, &event)
}

func getClub(ID string) (Club, error) {
	var club Club
	return club, getDocument(tblClub, ID, &club)
}

func getShooter(ID string) (Shooter, error) {
	var shooter Shooter
	return shooter, getDocument(tblShooter, ID, &shooter)
}

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
		// Marshal user data into bytes.
		buf, err := json.Marshal(event)
		if err != nil {
			return err
		}
		return bucket.Put(itob(event.ID), buf)
	})
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
		// Marshal user data into bytes.
		buf, err := json.Marshal(club)
		if err != nil {
			return err
		}
		return bucket.Put(itob(club.ID), buf)
	})
	return toB36(club.ID), err
}

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