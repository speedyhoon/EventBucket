package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var (
	//Databse bucket (table) names
	tblClub    = []byte("C")
	tblEvent   = []byte("E")
	tblShooter = []byte("S")
)

const (
	eNoBucket   = "Bucket %q not found!"
	eNoDocument = "'%v' document is empty / doesn't exist %q"
)

func makeBuckets() {
	db.Update(func(tx *bolt.Tx) error {
		for index, bucketName := range [][]byte{tblClub, tblEvent, tblShooter} {
			_, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				warn.Printf("Unable to create table %v in database", []string{"club", "event", "shooter"}[index])
			}
		}
		return nil
	})
}

func getDocument(bucketName []byte, ID string, result interface{}) error {
	byteID, err := b36toBy(ID)
	if err != nil {
		return err
	}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf(eNoBucket, bucketName)
		}

		document := bucket.Get(byteID)
		if len(document) == 0 {
			return fmt.Errorf(eNoDocument, ID, document)
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
	num, err := bucket.NextSequence()
	if err != nil {
		warn.Println("Failed to get the next sequence number.", err)
	}
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
		//Generate ID for the user.
		//This returns an error only if the Tx is closed or not writeable.
		//That can't happen in an Update() call so I ignore the error check.
		event.ID = b36

		//Marshal user data into bytes.
		buf, err := json.Marshal(event)
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return b36, err
}

func insertClub(club Club) (string, error) {
	if !club.IsDefault && !hasDefaultClub() {
		club.IsDefault = true
	}

	var b36 string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblClub)
		if err != nil {
			return err
		}
		var id []byte
		b36, id = nextID(bucket)
		//Generate ID for the user.
		//This returns an error only if the Tx is closed or not writeable.
		//That can't happen in an Update() call so I ignore the error check.
		club.ID = b36
		//Marshal user data into bytes.
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
		//Generate ID for the user.
		//This returns an error only if the Tx is closed or not writeable.
		//That can't happen in an Update() call so I ignore the error check.
		shooter.ID = b36
		//Marshal user data into bytes.
		buf, err := json.Marshal(shooter)
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return b36, err
}

func updateDocument(bucketName []byte, b36ID string, update interface{}, decode interface{}, function func(interface{}, interface{}) interface{}) error {
	ID, err := b36toBy(b36ID)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		bucket, err = tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		document := bucket.Get(ID)
		if len(document) == 0 {
			return fmt.Errorf(eNoDocument, ID, document)
		}

		err = json.Unmarshal(document, &decode)
		if err != nil {
			return fmt.Errorf("'%v' Query club unmarshaling failed: \n%q\n%#v\n", ID, document, err)
		}

		document, err = json.Marshal(function(decode, update))
		if err != nil {
			return err
		}

		return bucket.Put(ID, document)
	})
	return err
}

func updateShooterDetails(decode interface{}, contents interface{}) interface{} {
	shooter := decode.(*Shooter)
	update := *contents.(*Shooter)
	shooter.FirstName = update.FirstName
	shooter.Surname = update.Surname
	shooter.Club = update.Club
	shooter.Grade = update.Grade
	shooter.AgeGroup = update.AgeGroup
	shooter.Ladies = update.Ladies
	return shooter
}

func updateClubDetails(decode interface{}, contents interface{}) interface{} {
	club := decode.(*Club)
	update := contents.(*Club)
	//Manually set each one otherwise it would override the existing club and its details (Ranges, Shooters & their scores) since the form doesn't already have that info.
	club.Name = update.Name
	club.Address = update.Address
	club.Town = update.Town
	club.Postcode = update.Postcode
	club.Latitude = update.Latitude
	club.Longitude = update.Longitude
	club.IsDefault = update.IsDefault
	club.URL = update.URL
	return club
}

func updateClubDefault(decode interface{}, contents interface{}) interface{} {
	club := decode.(*Club)
	club.IsDefault = contents.(*Club).IsDefault
	return club
}

func insertClubMound(decode interface{}, contents interface{}) interface{} {
	club := decode.(*Club)
	club.Mounds = append(club.Mounds, *contents.(*Mound))
	return club
}

func updateEventDetails(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	update := contents.(*Event)
	//Manually set each one otherwise it would override the existing event and its details (Ranges, Shooters & their scores) since the form doesn't already have that info.
	event.Name = update.Name
	event.Club = update.Club
	event.Date = update.Date
	event.Time = update.Time
	event.Closed = update.Closed
	return event
}

func eventAddRange(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	newRange := contents.(*Range)
	newRange.ID = event.AutoInc.Range
	event.AutoInc.Range++
	event.Ranges = append(event.Ranges, *newRange)
	return event
}

func eventAddAgg(decode interface{}, contents interface{}) interface{} {
	event := eventAddRange(decode, contents).(*Event)
	aggRange := contents.(*Range)
	rangeID := aggRange.StrID()
	for sID, shooter := range event.Shooters {
		if shooter.Scores != nil {
			event.Shooters[sID].Scores[rangeID] = calcShooterAgg(aggRange.Aggs, shooter.Scores)
		}
	}
	return event
}

func editRange(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	rangeDetails := contents.(*Range)
	for i, r := range event.Ranges {
		if r.ID == rangeDetails.ID {
			event.Ranges[i].Name = rangeDetails.Name
			if event.Ranges[i].IsAgg {
				event.Ranges[i].Aggs = rangeDetails.Aggs
			} else {
				event.Ranges[i].Locked = rangeDetails.Locked
			}
			break
		}
	}
	return event
}

func editMound(decode interface{}, contents interface{}) interface{} {
	club := decode.(*Club)
	details := contents.(*Mound)
	if int(details.ID) < len(club.Mounds) {
		club.Mounds[details.ID].Name = details.Name
	}
	return club
}

func updateEventGrades(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	event.Grades = *contents.(*[]uint)
	return event
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

func clubsDataList() []option {
	var clubs []option
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblClub)
		if b == nil {
			//Club Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var club Club
			if json.Unmarshal(value, &club) == nil {
				clubs = append(clubs, option{Value: club.ID, Label: club.Name, Selected: club.IsDefault})
			}
			return nil
		})
	})
	if err != nil {
		warn.Println(err)
	}
	return clubs
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

func getCalendarEvents() ([]CalendarEvent, error) {
	var events []CalendarEvent
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblEvent)
		if b == nil {
			//Event Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var event CalendarEvent
			if json.Unmarshal(value, &event) == nil && !event.Closed {
				if event.Date != "" {
					event.ISO, _ = time.Parse("2006-01-02", event.Date)
				}
				events = append(events, event)
			}
			return nil
		})
	})
	return events, err
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

func eventShooterInsertDB(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	newShooter := *contents.(*Shooter)

	shooter := EventShooter{
		EID:       newShooter.ID,
		FirstName: newShooter.FirstName,
		Surname:   newShooter.Surname,
		Club:      newShooter.Club,
		ClubID:    newShooter.ClubID,
		AgeGroup:  newShooter.AgeGroup,
		Ladies:    newShooter.Ladies,
	}

SearchNextGrade:
	//Loop through the shooters selected grades & add a new shooter for each with a different grades.
	for _, gradeID := range newShooter.Grade {
		for _, s := range event.Shooters {
			if s.EID == shooter.EID && s.Grade == gradeID {
				warn.Printf("Shooter %v %v is not allowed to enter into %v event twice with the same grade %v.\n", shooter.FirstName, shooter.Surname, event.Name, globalGrades[gradeID].Name)
				continue SearchNextGrade
			}
		}

		//Assign shooter ID
		linkedID := event.AutoInc.Shooter
		shooter.ID = linkedID
		shooter.Grade = gradeID
		event.Shooters = append(event.Shooters, shooter)

		//Increment Event Shooter ID
		event.AutoInc.Shooter++

		//Some events shooters from grade X are automatically added to grade Y, e.g. Shooters in Match Reserve are able to win prizes in the higher grade Match Open.
		for _, grade := range globalGrades[gradeID].DuplicateTo {
			//Don't add the shooter because they have already selected to enter into the duplicate grade.
			if !containsUint(newShooter.Grade, grade) {
				shooter.ID = event.AutoInc.Shooter
				shooter.Grade = grade
				shooter.Hidden = true
				shooter.LinkedID = linkedID
				event.Shooters = append(event.Shooters, shooter)
				event.AutoInc.Shooter++
			}
		}
	}
	return event
}

func containsUint(list []uint, searchFor uint) bool {
	for _, x := range list {
		if x == searchFor {
			return true
		}
	}
	return false
}

func eventShooterUpdater(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	shooter := *contents.(*EventShooter)
	event.Shooters[shooter.ID].FirstName = shooter.FirstName
	event.Shooters[shooter.ID].Surname = shooter.Surname
	event.Shooters[shooter.ID].Club = shooter.Club
	event.Shooters[shooter.ID].Grade = shooter.Grade
	event.Shooters[shooter.ID].AgeGroup = shooter.AgeGroup
	event.Shooters[shooter.ID].Ladies = shooter.Ladies
	event.Shooters[shooter.ID].Disabled = shooter.Disabled
	return event
}

func upsertScore(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	shooter := *contents.(*shooterScore)

	if event.Shooters[shooter.id].Scores == nil {
		event.Shooters[shooter.id].Scores = make(map[string]Score)
	}
	event.Shooters[shooter.id].Scores[shooter.rangeID] = shooter.score

	event.Shooters[shooter.id].Scores = calcShooterAggs(event.Ranges, event.Shooters[shooter.id].Scores)
	return event
}

func calcShooterAggs(ranges []Range, shooterScores map[string]Score) map[string]Score {
	for _, r := range ranges {
		if r.IsAgg {
			shooterScores[r.StrID()] = calcShooterAgg(r.Aggs, shooterScores)
		}
	}
	return shooterScores
}

func calcShooterAgg(aggRangeIDs []uint, shooterScores map[string]Score) Score {
	var total, centers uint
	var countBack, countBack2 string
	for _, id := range aggRangeIDs {
		aggID := fmt.Sprintf("%d", id)
		score, ok := shooterScores[aggID]
		if ok {
			total += score.Total
			centers += score.Centers
			countBack = score.CountBack
			countBack2 = score.CountBack2
		}
	}
	return Score{
		Total:      total,
		Centers:    centers,
		CountBack:  countBack,
		CountBack2: countBack2,
	}
}

//Converts base36 string to binary used for bolt maps
func b36toBy(id string) ([]byte, error) {
	num, err := strconv.ParseUint(id, 36, 64)
	if err != nil {
		return []byte{}, err
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return b, nil
}

func getSearchShooters(firstName, surname, club string) ([]Shooter, uint, error) {
	var shooters []Shooter
	var totalQty uint

	//Search for shooters in the default club if all search values are empty.
	if firstName == "" && surname == "" && club == "" {
		club = defaultClubName()
	}

	firstName = strings.ToLower(firstName)
	surname = strings.ToLower(surname)
	club = strings.ToLower(club)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblShooter)
		if b == nil {
			return fmt.Errorf(eNoBucket, tblShooter)
		}
		totalQty = uint(tx.Bucket(tblShooter).Stats().KeyN)
		return b.ForEach(func(_, value []byte) error {
			var shooter Shooter
			//strings.Contains returns true when sub-string is "" (empty string)
			if json.Unmarshal(value, &shooter) == nil && strings.Contains(strings.ToLower(shooter.FirstName), firstName) && strings.Contains(strings.ToLower(shooter.Surname), surname) && strings.Contains(strings.ToLower(shooter.Club), club) {
				shooters = append(shooters, shooter)
			}
			return nil
		})
	})
	return shooters, totalQty, err
}

func searchShootersOptions(firstName, surname, club string) []option {
	if firstName == "" && surname == "" && club == "" {
		club = defaultClubName()
	}

	firstName = strings.ToLower(firstName)
	surname = strings.ToLower(surname)
	club = strings.ToLower(club)

	shooters := []option{{}}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblShooter)
		if b == nil {
			return fmt.Errorf(eNoBucket, tblShooter)
		}
		return b.ForEach(func(_, value []byte) error {
			var shooter Shooter
			//strings.Contains returns true when sub-string is "" (empty string)
			if json.Unmarshal(value, &shooter) == nil && strings.Contains(strings.ToLower(shooter.FirstName), firstName) && strings.Contains(strings.ToLower(shooter.Surname), surname) && strings.Contains(strings.ToLower(shooter.Club), club) {
				shooters = append(shooters, option{Value: shooter.ID, Label: shooter.FirstName + " " + shooter.Surname + ", " + shooter.Club})
			}
			return nil
		})
	})
	if err != nil {
		warn.Println(err)
	}
	return shooters
}

func getClubByName(clubName string) (Club, error) {
	var club Club
	const success = "Found the club you were looking for"
	clubName = strings.ToLower(clubName)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblClub)
		if b == nil {
			return fmt.Errorf(eNoBucket, tblClub)
		}
		return b.ForEach(func(_, value []byte) error {
			//Case insensitive search
			if json.Unmarshal(value, &club) == nil && strings.ToLower(club.Name) == clubName {
				return fmt.Errorf(success)
			}
			return nil
		})
	})
	if err != nil && err.Error() == success {
		return club, nil
	}
	return Club{}, fmt.Errorf("Couldn't find club with name %v", clubName)
}
