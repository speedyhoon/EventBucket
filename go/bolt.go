package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	//Databse collection names
	tblClub    = []byte("C")
	tblEvent   = []byte("E")
	tblShooter = []byte("S")
)

func getDocument(collection []byte, ID string, result interface{}) error {
	byteID, err := b36toBy(ID)
	if err != nil {
		return err
	}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(collection)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", collection)
		}

		document := bucket.Get(byteID)
		if len(document) == 0 {
			return fmt.Errorf("'%v' document is empty / doesn't exist %q", ID, document)
		}
		err = json.Unmarshal(document, &result)
		if err != nil {
			warn.Printf("'%v' Query document unmarshaling failed: \n%q\n%#v\n", ID, document, err)
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

func nextID(bucket *bolt.Bucket) (string, []byte) {
	num, _ := bucket.NextSequence()
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return strconv.FormatUint(num, 36), b
}

func insertEvent(event Event) (string, error) {
	var b36 string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblEvent)
		if err != nil {
			return err
		}
		var id []byte
		b36, id = nextID(bucket)
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		event.ID = b36

		// Marshal user data into bytes.
		buf, err := json.Marshal(event)
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return b36, err
}

func insertClub(club Club) (string, error) {
	var b36 string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblClub)
		if err != nil {
			return err
		}
		var id []byte
		b36, id = nextID(bucket)
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		club.ID = b36
		// Marshal user data into bytes.
		buf, err := json.Marshal(club)
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return b36, err
}

func insertShooter(shooter Shooter) (string, error) {
	var b36 string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblShooter)
		if err != nil {
			return err
		}
		var id []byte
		b36, id = nextID(bucket)
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		shooter.ID = b36
		// Marshal user data into bytes.
		buf, err := json.Marshal(shooter)
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return b36, err
}

func updateShooter(shooter Shooter, eventID string) error {
	var sID /*, eID*/ []byte
	var err error
	var buf []byte
	sID, err = b36toBy(shooter.ID)
	if err != nil {
		return err
	}
	// Marshal user data into bytes.
	buf, err = json.Marshal(shooter)
	if err != nil {
		return err
	}
	/*if eventID != "" {
		eID, err = b36toBy(eventID)
		if err != nil {
			return err
		}
	}*/

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tblShooter)
		if bucket == nil {
			//Shooter Bucket isn't created yet
			return nil
		}
		//TODO This will destroy all the shooters scores. needs a fix!
		return bucket.Put(sID, buf)
	})
	return err
}

func updateDoc(collectionName []byte, ID string, document interface{}) error {
	/*err := conn.C(collectionName).UpdateId(ID, document)
	if err != nil {
		warn.Println(err)
	}
	return err*/
	return nil
}

func updateEventDetails(update Event) error {
	eID, err := b36toBy(update.ID)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tblEvent)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", tblEvent)
		}

		document := bucket.Get(eID)
		if len(document) == 0 {
			return fmt.Errorf("'%v' document is empty / doesn't exist %q", update.ID, document)
		}
		var event Event
		err = json.Unmarshal(document, &event)
		if err != nil {
			return fmt.Errorf("'%v' Query event unmarshaling failed: \n%q\n%#v\n", update.ID, document, err)
		}
		//Manually set each one otherwise it would override the existing event and its details (Ranges, Shooters & their scores) since the form doesn't already have that info.
		event.Name = update.Name
		event.Club = update.Club
		event.Date = update.Date
		event.Time = update.Time
		event.Closed = update.Closed

		document, err = json.Marshal(event)
		if err != nil {
			return err
		}

		return bucket.Put(eID, document)
	})
	return err
}

func eventAddRange(eventID string, newRange Range) error {
	eID, err := b36toBy(eventID)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tblEvent)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", tblEvent)
		}

		document := bucket.Get(eID)
		if len(document) == 0 {
			return fmt.Errorf("'%v' document is empty / doesn't exist %q", eventID, document)
		}
		var event Event
		err = json.Unmarshal(document, &event)
		if err != nil {
			return fmt.Errorf("'%v' Query event unmarshaling failed: \n%q\n%#v\n", eventID, document, err)
		}
		//Manually set each one otherwise it would override the existing event and its details (Ranges, Shooters & their scores) since the form doesn't already have that info.
		newRange.ID = event.AutoInc.Range
		event.Ranges = append(event.Ranges, newRange)
		event.AutoInc.Range++

		document, err = json.Marshal(event)
		if err != nil {
			warn.Println("error", err)
			return err
		}

		return bucket.Put(eID, document)
	})
	return err
}

func getClubs() ([]Club, error) {
	var clubs []Club
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblClub)
		if b == nil {
			//Club Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var club Club
			if json.Unmarshal(value, &club) == nil {
				clubs = append(clubs, club)
			}
			return nil
		})
	})
	return clubs, err
}

func getEvents(query func(Event) bool) ([]Event, error) {
	var events []Event
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblEvent)
		if b == nil {
			//Event Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var event Event
			if json.Unmarshal(value, &event) == nil && query(event) {
				events = append(events, event)
			}
			return nil
		})
	})
	return events, err
}

func getShooters() ([]Shooter, error) {
	var shooters []Shooter
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblShooter)
		if b == nil {
			//Shooter Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var shooter Shooter
			if json.Unmarshal(value, &shooter) == nil {
				shooters = append(shooters, shooter)
			}
			return nil
		})
	})
	return shooters, err
}

func hasDefaultClub() bool {
	return getDefaultClub().Name != ""
}

func defaultClubName() string {
	return getDefaultClub().Name
}

func getDefaultClub() Club {
	var club Club
	var found bool
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblClub)
		if b == nil {
			//Club Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			if json.Unmarshal(value, &club) == nil && club.IsDefault {
				found = true
				return fmt.Errorf("no error")
			}
			return nil
		})
	})
	if found {
		return club
	}
	return Club{}
}

func eventShooterInsertDB(ID string, shooter EventShooter) error {
	byteID, err := b36toBy(ID)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tblEvent)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", tblEvent)
		}

		document := bucket.Get(byteID)
		if len(document) == 0 {
			return fmt.Errorf("'%v' document is empty / doesn't exist %q", ID, document)
		}
		var event Event
		err = json.Unmarshal(document, &event)
		if err != nil {
			return err
		}

		//Assign shooter ID
		shooter.ID = event.AutoInc.Shooter
		event.Shooters = append(event.Shooters, shooter)

		//Increment Event Shooter ID
		event.AutoInc.Shooter++

		//If shooter is Match Reserve, duplicate them in the Match Open category
		if shooter.Grade == 8 {
			shooter.ID = event.AutoInc.Shooter
			shooter.Grade = 7
			shooter.Hidden = true
			event.Shooters = append(event.Shooters, shooter)
			event.AutoInc.Shooter++
		}

		document, err = json.Marshal(event)
		if err != nil {
			return err
		}

		return bucket.Put(byteID, document)
	})
	return err
}

func upsertScore(eventID, rID string, sID uint64, score Score) error {
	byteID, err := b36toBy(eventID)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tblEvent)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", tblEvent)
		}

		document := bucket.Get(byteID)
		if len(document) == 0 {
			return fmt.Errorf("'%v' document is empty / doesn't exist %q", eventID, document)
		}
		var event Event
		err = json.Unmarshal(document, &event)
		if err != nil {
			return err
		}

		if event.Shooters[sID].Scores == nil {
			event.Shooters[sID].Scores = make(map[string]Score)
		}
		event.Shooters[sID].Scores[rID] = score

		document, err = json.Marshal(event)
		if err != nil {
			return err
		}

		return bucket.Put(byteID, document)
	})
	return err
}

//Converts base36 string to uint64
func b36tou(id string) (uint64, error) {
	return strconv.ParseUint(id, 36, 64)
}

//Converts base36 string to binary used for bolt maps
func b36toBy(id string) ([]byte, error) {
	num, err := b36tou(id)
	if err != nil {
		return []byte{}, err
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return b, nil
}

func getSearchShooters(firstName, surname, club string) ([]Shooter, error) {
	var shooters []Shooter
	var shooter Shooter
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblShooter)
		if b == nil {
			return fmt.Errorf("Bucket %q not found!", tblShooter)
		}
		return b.ForEach(func(_, value []byte) error {
			//strings.Contains returns true when substr is "" (empty string)
			if json.Unmarshal(value, &shooter) == nil && strings.Contains(strings.ToLower(shooter.FirstName), strings.ToLower(firstName)) && strings.Contains(strings.ToLower(shooter.Surname), strings.ToLower(surname)) && strings.Contains(strings.ToLower(shooter.Club), strings.ToLower(club)) {
				shooters = append(shooters, shooter)
			}
			return nil
		})
	})
	return shooters, err
}
