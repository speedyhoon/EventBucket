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
	//Database connection.
	db *bolt.DB

	//Database bucket (table) names
	tblClub    = []byte{0}
	tblEvent   = []byte{1}
	tblShooter = []byte{2}
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

func getEvent(ID string) (event Event, err error) {
	return event, getDocument(tblEvent, ID, &event)
}

func getClub(ID string) (club Club, err error) {
	return club, getDocument(tblClub, ID, &club)
}

func getShooter(ID string) (shooter Shooter, err error) {
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

func insertDocument(tblName []byte, document interface{}, assignID func(interface{}, string) interface{}) (string, error) {
	var b36 string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tblName)
		if err != nil {
			return err
		}
		var id []byte
		b36, id = nextID(bucket)
		//Generate ID for the user.
		//This returns an error only if the Tx is closed or not writable.
		//That can't happen in an Update() call so I ignore the error check.

		//Marshal user data into bytes.
		buf, err := json.Marshal(assignID(document, b36))
		if err != nil {
			return err
		}
		return bucket.Put(id, buf)
	})
	return b36, err
}

func (event Event) insert() (string, error) {
	return insertDocument(
		tblEvent,
		event,
		func(i interface{}, b36 string) interface{} {
			o := i.(Event)
			o.ID = b36
			return o
		},
	)
}

func (club Club) insert() (string, error) {
	if !club.IsDefault && !hasDefaultClub() {
		club.IsDefault = true
	}
	return insertDocument(
		tblClub,
		club,
		func(i interface{}, b36 string) interface{} {
			o := i.(Club)
			o.ID = b36
			return o
		},
	)
}

func (shooter Shooter) insert() (string, error) {
	return insertDocument(
		tblShooter,
		shooter,
		func(i interface{}, b36 string) interface{} {
			o := i.(Shooter)
			o.ID = b36
			return o
		},
	)
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
			return fmt.Errorf("'%v' Query club unmarshaling failed: \n%q\n%#v", ID, document, err)
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
	shooter.Grades = update.Grades
	shooter.AgeGroup = update.AgeGroup
	shooter.Sex = update.Sex
	return shooter
}

func updateClubDetails(decode interface{}, contents interface{}) interface{} {
	club := decode.(*Club)
	update := contents.(*Club)
	//Manually set each one otherwise it would override the existing club and its details (Mounds etc)
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

func insertClubMound(decode interface{}, mound interface{}) interface{} {
	club := decode.(*Club)
	club.Mounds = append(club.Mounds, mound.(string))
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
			r.Name = rangeDetails.Name
			if r.IsAgg {
				r.Aggs = rangeDetails.Aggs
			} else {
				r.Locked = rangeDetails.Locked
				//if r.Locked {
				//event.Shooters = addGradeSeparatorToShooterObjectAndPositions(event.Shooters, r.StrID())
				//info.Println("Recalculate range", r.ID, r.Name)
				//}
			}
			//Move range if the order has changed
			if uint(i) != rangeDetails.Order {
				//Cut range
				if i <= 0 {
					event.Ranges = append([]Range{}, event.Ranges[1:]...)
				} else if i >= len(event.Ranges)-1 {
					event.Ranges = append([]Range{}, event.Ranges[:i]...)
				} else {
					event.Ranges = append(event.Ranges[:i], event.Ranges[i+1:]...)
				}

				//Paste range
				if rangeDetails.Order <= 0 {
					event.Ranges = append([]Range{r}, event.Ranges...)
				} else if rangeDetails.Order >= uint(len(event.Ranges)-1) {
					event.Ranges = append(event.Ranges, r)
				} else {
					event.Ranges = append(event.Ranges[:rangeDetails.Order], append([]Range{r}, event.Ranges[rangeDetails.Order:]...)...)
				}

				//Now rearrange any list of Aggregates that contain this range. This saves performing a double loop in the scoreboard page because the aggregate list is now out of order.
				for w, rng := range event.Ranges {
					if rng.IsAgg {
						var rangeList []uint
						for _, t := range event.Ranges {
							for _, rngID := range rng.Aggs {
								if rngID == t.ID {
									rangeList = append(rangeList, rngID)
								}
							}
						}
						event.Ranges[w].Aggs = rangeList
					}
				}
			} else {
				event.Ranges[i] = r
			}
			break
		}
	}
	return event
}

func editMound(decode interface{}, contents interface{}) interface{} {
	club := decode.(*Club)
	mound := contents.(*Mound)
	if int(mound.ID) < len(club.Mounds) {
		club.Mounds[mound.ID] = mound.Name
	}
	return club
}

func updateEventGrades(decode interface{}, contents interface{}) interface{} {
	event := decode.(*Event)
	event.Grades = *contents.(*[]uint)
	return event
}

func getClubs() (clubs []Club, err error) {
	return clubs, db.View(func(tx *bolt.Tx) error {
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
}

func getMapClubs(clubID string) (clubs []MapClub, err error) {
	return clubs, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tblClub)
		if b == nil {
			//Club Bucket isn't created yet
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var club MapClub
			if json.Unmarshal(value, &club) == nil && clubID != "" && club.ID == clubID || clubID == "" && club.Latitude != 0 && club.Longitude != 0 {
				clubs = append(clubs, club)
			}
			return nil
		})
	})
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
	return defaultClubName() != ""
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
		AgeGroup:  newShooter.AgeGroup,
		Sex:       newShooter.Sex,
	}

SearchNextGrade:
	//Loop through the shooters selected grades & add a new shooter for each with a different grades.
	for _, gradeID := range newShooter.Grades {
		for _, s := range event.Shooters {
			if s.EID == shooter.EID && s.Grade == gradeID {
				warn.Printf("Shooter %v %v is not allowed to enter into %v event twice with the same grade %v.\n", shooter.FirstName, shooter.Surname, event.Name, globalGrades[gradeID].Name)
				continue SearchNextGrade
			}
		}

		//Assign shooter ID
		shooter.ID = event.AutoInc.Shooter
		shooter.Grade = gradeID
		event.Shooters = append(event.Shooters, shooter)

		//Increment Event Shooter ID
		event.AutoInc.Shooter++

		//Some events shooters from grade X are automatically added to grade Y, e.g. Shooters in Match Reserve are able to win prizes in the higher grade Match Open.
		for _, grade := range globalGrades[gradeID].DuplicateTo {
			//Don't add the shooter because they have already selected to enter into the duplicate grade.
			if !containsUint(newShooter.Grades, grade) {
				shooter.LinkedID = shooter.ID
				shooter.ID = event.AutoInc.Shooter
				shooter.Grade = grade
				shooter.Hidden = true
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
	event.Shooters[shooter.ID].Sex = shooter.Sex
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
	var total, centers, centers2, shootOff uint
	var countBack, countBack2 string
	for _, id := range aggRangeIDs {
		aggID := fmt.Sprintf("%d", id)
		score, ok := shooterScores[aggID]
		if ok {
			total += score.Total
			centers += score.Centers
			centers2 += score.Centers2
			countBack = score.CountBack
			countBack2 = score.CountBack2
			shootOff = score.ShootOff
		}
	}
	return Score{
		Total:      total,
		Centers:    centers,
		Centers2:   centers2,
		CountBack:  countBack,
		CountBack2: countBack2,
		ShootOff:   shootOff,
	}
}

//Converts base36 string to []byte used for bolt maps
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
	//TODO this is quite dodgy
	if err != nil && err.Error() == success {
		return club, nil
	}
	return Club{}, fmt.Errorf("Couldn't find club with name %v", clubName)
}
