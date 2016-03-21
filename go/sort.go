package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type sortEventShooter func(r1 string, p1, p2 *EventShooter) bool

type multiEventShooterSorter struct {
	shooter []EventShooter
	sort    []sortEventShooter
	rangeID string
}

func sorterGrade(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Grade < c2.Grade
}
func sorterTotal(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].Total > c2.Scores[rangeID].Total
}
func sorterCenters(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].Centers > c2.Scores[rangeID].Centers
}
func sorterCountback(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].CountBack > c2.Scores[rangeID].CountBack
}
func sorterFirstName(rangeID string, c1, c2 *EventShooter) bool {
	return c1.FirstName < c2.FirstName
}
func sorterShootOff(rangeID string, c1, c2 *EventShooter) bool {
	if c1.Scores[rangeID].CountBack != "" && c1.Scores[rangeID].CountBack == c2.Scores[rangeID].CountBack {
		info.Printf("shooters scores are the same? c1= g:%v t:%v c:%v b:%v h:%v c2= g:%v t:%v c:%v b:%v h:%v", c1.Grade, c1.Scores[rangeID].Total, c1.Scores[rangeID].Centers, c1.Scores[rangeID].CountBack, c1.Scores[rangeID].ShootOff, c2.Grade, c2.Scores[rangeID].Total, c2.Scores[rangeID].Centers, c2.Scores[rangeID].CountBack, c2.Scores[rangeID].ShootOff)
		temp := c1.Scores[rangeID]
		//temp.Warning = legendShootOff
		c1.Scores[rangeID] = temp

		temp = c2.Scores[rangeID]
		//temp.Warning = legendShootOff
		c2.Scores[rangeID] = temp
	}
	return c1.Scores[rangeID].ShootOff < c2.Scores[rangeID].ShootOff
}

func sortShooters(rangeID string) *multiEventShooterSorter {
	if rangeID != "" {
		return orderShooters(rangeID, sorterGrade, sorterTotal, sorterCenters, sorterCountback, sorterShootOff)
	}
	var temp multiEventShooterSorter
	return &temp
}

func orderShooters(rangeID string, sort ...sortEventShooter) *multiEventShooterSorter {
	return &multiEventShooterSorter{
		sort:    sort,
		rangeID: rangeID,
	}
}

func (ms *multiEventShooterSorter) Sort(shooter []EventShooter) {
	ms.shooter = shooter
	sort.Sort(ms)
}

func (ms *multiEventShooterSorter) Len() int {
	return len(ms.shooter)
}

func (ms *multiEventShooterSorter) Swap(i, j int) {
	ms.shooter[i], ms.shooter[j] = ms.shooter[j], ms.shooter[i]
}

func (ms *multiEventShooterSorter) Less(i, j int) bool {
	p, q := &ms.shooter[i], &ms.shooter[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.sort)-1; k++ {
		sort := ms.sort[k]
		switch {
		case sort(ms.rangeID, p, q):
			return true
		case sort(ms.rangeID, q, p):
			return false
		}
	}
	return ms.sort[k](ms.rangeID, p, q)
}

func triggerScoreCalculation(newScore Score, rangeID uint64, shooter EventShooter, event Event) {
	shootFinished := hasShootFinished(newScore.Shots, shooter.Grade)

	shooterIDs := []uint64{shooter.ID}
	if shooter.LinkedID != 0 {
		linkedShooterID := shooter.LinkedID
		//Add the linked shooter to the update
		shooterIDs = append(shooterIDs, linkedShooterID)

		//Add the linked shooter to the aggregate & grade update ticker
		recalculateShooters.mu.Lock()
		recalculateShooters.values[fmt.Sprintf("%v.%v.%v", event.ID, linkedShooterID, rangeID)] = calculateShooter{
			eventID:    event.ID,
			shooterID:  linkedShooterID,
			strRangeID: fmt.Sprintf("%v", rangeID),
			rangeID:    rangeID,
		}
		recalculateShooters.mu.Unlock()

		if shootFinished {
			recalculateGrades.mu.Lock()
			recalculateGrades.values[fmt.Sprintf("%v.%v.%v", event.ID, shooter.Grade, rangeID)] = calculateGrade{
				eventID: event.ID,
				gradeID: shooter.Grade,
				rangeID: rangeID,
			}
			recalculateGrades.mu.Unlock()
		}
	}

	//eventTotalScoreUpdate(event.ID, rangeID, shooterIDs, newScore)

	aggsFound := searchForAggs(event.Ranges, rangeID)
	aggsFound = append(aggsFound, rangeID)

	recalculateShooters.mu.Lock()
	recalculateShooters.values[fmt.Sprintf("%v.%v.%v", event.ID, shooter.ID, rangeID)] = calculateShooter{
		eventID:    event.ID,
		shooterID:  shooter.ID,
		strRangeID: fmt.Sprintf("%v", rangeID),
		rangeID:    rangeID,
	}
	recalculateShooters.mu.Unlock()

	for _, thisRangeID := range aggsFound {
		if shootFinished {
			recalculateGrades.mu.Lock()
			recalculateGrades.values[fmt.Sprintf("%v.%v.%v", event.ID, shooter.Grade, thisRangeID)] = calculateGrade{
				eventID: event.ID,
				gradeID: shooter.Grade,
				rangeID: thisRangeID,
			}
			recalculateGrades.mu.Unlock()
		}
	}
}

type calculateShooter struct {
	eventID, strRangeID string
	shooterID, rangeID  uint64
}
type calculateGrade struct {
	eventID          string
	gradeID, rangeID uint64
}

type myShooterMutex struct {
	mu     sync.Mutex
	values map[string]calculateShooter
}

type myGradeMutex struct {
	mu     sync.Mutex
	values map[string]calculateGrade
}

var (
	recalculateShooters = myShooterMutex{
		values: make(map[string]calculateShooter),
	}
	recalculateGrades = myGradeMutex{
		values: make(map[string]calculateGrade),
	}
)

func startTicker() {
	//TODO off load the grade position and shoot off calculations into another process
	//TODO only calculate the grade positions once a shooter has finished the range
	var shooters map[string]calculateShooter
	var grades map[string]calculateGrade
	for range time.NewTicker(time.Second * 10).C {
		recalculateShooters.mu.Lock()
		if len(recalculateShooters.values) > 0 {
			shooters = recalculateShooters.values
			recalculateShooters.values = make(map[string]calculateShooter)
			recalculateShootersAggs(shooters)
		}
		recalculateShooters.mu.Unlock()
		recalculateGrades.mu.Lock()
		if len(recalculateGrades.values) > 0 {
			grades = recalculateGrades.values
			recalculateGrades.values = make(map[string]calculateGrade)
			recalculateGradePositions(grades)
		}
		recalculateGrades.mu.Unlock()
	}
}

func recalculateShootersAggs(updates map[string]calculateShooter) {
	trace.Println("executing recalculateShooterAggs")
	var event Event
	var err error
	var aggsFound []uint64
	var updateBson map[string]interface{}
	var previousEventID string
	for _, updateData := range updates {

		if updateData.eventID != previousEventID {
			previousEventID = updateData.eventID
			event, err = getEvent(updateData.eventID)
		}
		if err == nil && event.Ranges != nil && event.Shooters != nil {
			aggsFound = searchForAggs(event.Ranges, updateData.rangeID)
			if len(aggsFound) > 0 {
				updateBson = make(map[string]interface{})
				for index, data := range calculateAggs(event.Shooters[updateData.shooterID].Scores, aggsFound, []uint64{updateData.shooterID}, event.Ranges) {
					updateBson[index] = data
				}
				if len(updateBson) > 0 {
					//tableUpdateData(tblEvent, event.ID, updateBson)
				}
			}
		}
	}
	trace.Println("finished recalculateShooterAggs")
}

func recalculateGradePositions(updates map[string]calculateGrade) {
	trace.Println("executing grade recalculation")
	var event Event
	var err error
	//var updateBson map[string]interface{}
	var shooterQty /*position, */, shouldBePosition uint64
	var shootEqual, updateRequired bool
	var previousEventID /* positionEqual,*/ /*positionOrdinal,*/, strRangeID string

	for _, updateData := range updates {
		//Only get the event when it is different
		if updateData.eventID != previousEventID {
			if updateRequired {
				//tableUpdateData(tblEvent, updateData.eventID, updateBson)
				updateRequired = false
			}
			//updateBson = make(map[string]interface{})
			event, err = getEvent(updateData.eventID)

			if err != nil {
				warn.Println(err)
				break
			}

			//TODO remove adding shooter ids!
			//Add shooter ids to the shooter objects
			//event.Shooters = addShooterIDsToShooterObjects(event.Shooters)

			shooterQty = uint64(len(event.Shooters))
		}
		strRangeID = fmt.Sprintf("%v", updateData.rangeID)

		//sort shooters by the current rangeID
		sortShooters(strRangeID).Sort(event.Shooters)

		shouldBePosition = 0
		shootEqual = false
		//positionEqual = ""
		for index, shooter := range event.Shooters {
			if shooter.Grade == updateData.gradeID {
				shouldBePosition++
				if !shootEqual {
					//position = shouldBePosition
					//positionEqual = ""
				} else {
					//positionEqual = "="
					shootEqual = false
				}
				if shooter.Scores[strRangeID].ShootOff < 0 {
					//Shooter has the same score as the previous shooter (index-1)
					//positionEqual = "="
					if uint64(index+1) < shooterQty && shooter.Grade == event.Shooters[index+1].Grade && shooter.Scores[strRangeID] == event.Shooters[index+1].Scores[strRangeID] {
						shootEqual = true
					}
				}
				/*positionOrdinal = positionEqual + ordinal(position)
				if shooter.Scores[strRangeID].Total != 0 && (shooter.Scores[strRangeID].Position != position || shooter.Scores[strRangeID].Ordinal != positionOrdinal) {
					updateRequired = true
					updateBson[dot("S", shooter.ID, updateData.rangeID, "o")] = positionOrdinal
					updateBson[dot("S", shooter.ID, updateData.rangeID, "p")] = position
				}*/
			}
		}
	}
	if updateRequired {
		//tableUpdateData(tblEvent, event.ID, updateBson)
	}
	info.Println("finished grade recalculation")
}

func hasShootFinished(shots string, grade uint64) bool {
	classSettings := defaultClassSettings[grades()[grade].classID]
	return uint64(len(strings.Replace(shots[classSettings.SightersQty:], "-", "", -1))) == classSettings.ShotsQty
}

func searchForAggs(ranges []Range, rangeID uint64) []uint64 {
	var aggsFound []uint64
	//var foundRangeID uint64
	//var err error
	for indexRangeID, rangeObj := range ranges {
		if len(rangeObj.Aggs) > 0 {
			for range rangeObj.Aggs {
				//for _, rangeID := range rangeObj.Aggs {
				//foundRangeID, err = strconv.Atoi(strRangeID)
				//if err == nil /*&& rangeID == rangeID*/ {
				aggsFound = append(aggsFound, uint64(indexRangeID))
				//}
			}
		}
	}
	return aggsFound
}

/*func eventTotalScoreUpdate(eventID string, rangeID uint64, shooterIDs []uint64, score Score) Event {
	var event Event
	if conn != nil {
		updateSetter := make(map[string]interface{})
		for _, shooterID := range shooterIDs {
			updateSetter[dot(schemaSHOOTER, shooterID, rangeID)] = score
		}
		change := mgo.Change{
			Upsert: true,
			Update: map[string]interface{}{
				"$set": updateSetter,
			},
		}
		_, err := conn.C(tblEvent).FindId(eventID).Apply(change, &event)
		//TODO better error handling would be nice
		if err != nil {
			warn.Println(err)
		}
	}
	return event
}*/

// Ordinal gives you the input number in a rank/ordinal format.
// Ordinal(3) -> 3rd
//author github.com/dustin/go-humanize/blob/master/ordinals.go
func ordinal(x uint64) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%v", x, suffix)
}

func calculateAggs(shooterScores map[string]Score, ranges []uint64, shooterIDs []uint64, eventRanges []Range) map[string]interface{} {
	updateBson := make(map[string]interface{})
	if shooterScores == nil {
		return updateBson
	}
	var total, centres uint64
	var countBack string

	for _, aggID := range ranges {
		total = 0
		centres = 0
		for _, rID := range eventRanges[aggID].Aggs {
			rangeID := fmt.Sprintf("%d", rID)
			total += shooterScores[rangeID].Total
			centres += shooterScores[rangeID].Centers
			countBack = shooterScores[rangeID].CountBack
		}
		for _, shooterID := range shooterIDs {
			updateBson[dot("S", shooterID, aggID)] = Score{Total: total, Centers: centres, CountBack: countBack}
		}
	}
	return updateBson
}

/*func tableUpdateData(collectionName, documentID string, data map[string]interface{}) {
	if conn != nil {
		_, err := conn.C(collectionName).FindId(documentID).Apply(mgo.Change{
			Upsert: true,
			Update: M{"$set": data},
		}, make(map[string]interface{}))
		if err != nil {
			warn.Println(err)
		}
	}
}*/

//Used for database schema translation dot notation
func dot(elem ...interface{}) string {
	var dots []string
	for _, element := range elem {
		dots = append(dots, fmt.Sprintf("%v", element))
	}
	return strings.Join(dots, ".")
}