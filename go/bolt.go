package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/speedyhoon/frm"
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
	eNoDocument = "'%v' document is empty / doesn't exist %q"
)

func startDB(dbPath string) {
	//Create database directory if needed.
	err := mkDir(filepath.Dir(dbPath))
	if err != nil {
		log.Fatal(err)
	}

	//Database save location
	db, err = bolt.Open(dbPath, 0644, &bolt.Options{Timeout: time.Second * 8})
	if err != nil {
		log.Fatalln("Connection timeout. Unable to open", dbPath)
	}

	//Prepare database by creating all buckets (tables) needed. Otherwise view (read only) transactions will fail.
	err = db.Update(func(tx *bolt.Tx) error {
		for index, bucketName := range [][]byte{tblClub, tblEvent, tblShooter} {
			_, err = tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				log.Printf("Unable to create table %v in database", []string{"club", "event", "shooter"}[index])
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

//view checks if a bucket exists before executing the provided function myCall
func view(bucketName []byte, myCall func(*bolt.Bucket) error) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", bucketName)
		}
		return myCall(bucket)
	})
}

//search executes view() first and then checks if each item can be unmarshalled before executing the provided function myCall
func search(table []byte, object interface{}, myCall func(interface{}) error) error {
	return view(table, func(b *bolt.Bucket) error {
		return b.ForEach(func(_, value []byte) error {
			if err := json.Unmarshal(value, object); err != nil {
				return err
			}
			return myCall(object)
		})
	})
}

//tblQty returns the total number of records contained in the bucket (table)
func tblQty(bucketName []byte) (qty uint) {
	err := view(bucketName, func(bucket *bolt.Bucket) error {
		qty = uint(bucket.Stats().KeyN)
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return
}

func getDocument(bucketName []byte, ID string, result interface{}) error {
	byteID, err := b36toBy(ID)
	if err != nil {
		return err
	}

	return view(bucketName, func(bucket *bolt.Bucket) error {
		document := bucket.Get(byteID)
		if len(document) == 0 {
			return fmt.Errorf(eNoDocument, ID, document)
		}
		err = json.Unmarshal(document, &result)
		if err != nil {
			log.Printf("'%v' Query document unmarshaling failed: \n%q\n%#v\n", ID, document, err)
		}
		return err
	})
}

func getEvent(ID string) (event Event, err error) {
	err = getDocument(tblEvent, ID, &event)
	if err != nil {
		return
	}
	event.Club, err = getClub(event.ClubID)
	return
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
		log.Println("Failed to get the next sequence number.", err)
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return strconv.FormatUint(num, 36), b
}

func insertDocument(tblName []byte, document interface{}, assignID func(interface{}, string) interface{}) (b36 string, err error) {
	return b36, db.Update(func(tx *bolt.Tx) error {
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
	if !club.IsDefault && !defaultClub().IsDefault {
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
	event.DateTime = update.DateTime
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
			event.Shooters[sID].Scores[rangeID] = shooter.Scores.calcAgg(aggRange.Aggs)
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
	return clubs, search(tblClub, &Club{}, func(c interface{}) error {
		clubs = append(clubs, *c.(*Club))
		return nil
	})
}

func getMapClubs(clubID string) (clubs []Club, err error) {
	return clubs, search(tblClub, &Club{}, func(c interface{}) error {
		club := *c.(*Club)
		if clubID != "" && club.ID == clubID || clubID == "" && club.Latitude != 0 && club.Longitude != 0 {
			clubs = append(clubs, club)
		}
		return nil
	})
}

func clubsDataList() (clubs []frm.Option) {
	err := search(tblClub, &Club{}, func(c interface{}) error {
		club := *c.(*Club)
		clubs = append(clubs, frm.Option{Value: club.ID, Label: club.Name, Selected: club.IsDefault})
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return
}

func getEvents(query func(Event) bool) ([]Event, error) {
	var events []Event
	return events, search(tblEvent, &Event{}, func(e interface{}) error {
		event := *e.(*Event)
		if query(event) {
			event.Club, _ = getClub(event.ClubID)
			events = append(events, event)
		}
		return nil
	})
}

func defaultClub() Club {
	const success = "1"
	var club Club
	err := search(tblClub, &club, func(interface{}) error {
		if club.IsDefault {
			return fmt.Errorf(success)
		}
		return nil
	})
	if err != nil {
		if err.Error() == success {
			return club
		}
		log.Println(err)
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
				log.Printf("Shooter %v %v is not allowed to enter into %v event twice with the same grade %v.\n", shooter.FirstName, shooter.Surname, event.Name, globalGrades[gradeID].Name)
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
		event.Shooters[shooter.id].Scores = make(ScoreMap)
	}
	event.Shooters[shooter.id].Scores[shooter.rangeID] = shooter.score

	event.Shooters[shooter.id].Scores.calcShooterAggs(event.Ranges)
	return event
}

func (shooterScores ScoreMap) calcShooterAggs(ranges []Range) ScoreMap {
	for _, r := range ranges {
		if r.IsAgg {
			shooterScores[r.StrID()] = shooterScores.calcAgg(r.Aggs)
		}
	}
	return shooterScores
}

func (shooterScores ScoreMap) calcAgg(aggRangeIDs []uint) (total Score) {
	for _, id := range aggRangeIDs {
		if score, ok := shooterScores.get(id); ok {
			total.Total += score.Total
			total.Centers += score.Centers
			total.Centers2 += score.Centers2
			total.CountBack = score.CountBack
			total.CountBack2 = score.CountBack2
			total.ShootOff = score.ShootOff
		}
	}
	return total
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

func searchShooters(firstName, surname, club string) (shooters []Shooter) {
	//Search for shooters in the default club if all search values are empty.
	if firstName == "" && surname == "" && club == "" {
		club = defaultClub().Name
	}

	firstName = strings.ToLower(firstName)
	surname = strings.ToLower(surname)
	club = strings.ToLower(club)

	err := search(tblShooter, &Shooter{}, func(s interface{}) error {
		shooter := *s.(*Shooter)
		//strings.Contains returns true when sub-string is "" (empty string)
		if strings.Contains(strings.ToLower(shooter.FirstName), firstName) && strings.Contains(strings.ToLower(shooter.Surname), surname) {
			clubs, err := getClub(shooter.Club)
			if err == nil && strings.Contains(strings.ToLower(clubs.Name), club) {
				shooters = append(shooters, shooter)
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return
}

func searchShootersOptions(firstName, surname, club string) (options []frm.Option) {
	for _, s := range searchShooters(firstName, surname, club) {
		options = append(options, frm.Option{Value: s.ID, Label: s.FirstName + " " + s.Surname + ", " + s.Club})
	}
	return
}

func getClubByName(clubName string) (club Club, ok bool) {
	clubName = strings.ToLower(clubName)
	const success = "1"

	err := search(tblClub, &club, func(club interface{}) error {
		//Case insensitive search
		if strings.ToLower(club.(*Club).Name) == clubName {
			//Return a successful error to stop searching any further
			return fmt.Errorf(success)
		}
		return nil
	})
	return club, err != nil && err.Error() == success
}
