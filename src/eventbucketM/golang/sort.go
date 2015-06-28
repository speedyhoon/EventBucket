package main

import (
	"fmt"
	"sort"
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
func sorterCentres(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].Centres > c2.Scores[rangeID].Centres
}
func sorterCountback(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].CountBack > c2.Scores[rangeID].CountBack
}
func sorterFirstName(rangeID string, c1, c2 *EventShooter) bool {
	return c1.FirstName < c2.FirstName
}
func sorterShootOff(rangeID string, c1, c2 *EventShooter) bool {
	if c1.Scores[rangeID].CountBack != "" && c1.Scores[rangeID].CountBack == c2.Scores[rangeID].CountBack {
		info.Printf("shooters scores are the same? c1= g:%v t:%v c:%v b:%v h:%v c2= g:%v t:%v c:%v b:%v h:%v", c1.Grade, c1.Scores[rangeID].Total, c1.Scores[rangeID].Centres, c1.Scores[rangeID].CountBack, c1.Scores[rangeID].ShootOff, c2.Grade, c2.Scores[rangeID].Total, c2.Scores[rangeID].Centres, c2.Scores[rangeID].CountBack, c2.Scores[rangeID].ShootOff)
		temp := c1.Scores[rangeID]
		temp.Warning = legendShootOff
		c1.Scores[rangeID] = temp

		temp = c2.Scores[rangeID]
		temp.Warning = legendShootOff
		c2.Scores[rangeID] = temp
	}
	return c1.Scores[rangeID].ShootOff < c2.Scores[rangeID].ShootOff
}

func sortShooters(rangeID string) *multiEventShooterSorter {
	if rangeID != "" {
		return orderShooters(rangeID, sorterGrade, sorterTotal, sorterCentres, sorterCountback, sorterShootOff)
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

func triggerScoreCalculation(newScore Score, rangeID int, shooter EventShooter, event Event) {
	shootFinished := hasShootFinished(newScore.Shots, shooter.Grade)

	shooterIDs := []int{shooter.ID}
	if shooter.LinkedID != nil {
		linkedShooterID := *shooter.LinkedID
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

	eventTotalScoreUpdate(event.ID, rangeID, shooterIDs, newScore)

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
	shooterID, rangeID  int
}
type calculateGrade struct {
	eventID          string
	gradeID, rangeID int
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
	for range time.NewTicker(time.Second * 4).C {
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
	var aggsFound []int
	var updateBson M
	var previousEventID string
	for _, updateData := range updates {

		if updateData.eventID != previousEventID {
			previousEventID = updateData.eventID
			event, err = getEvent(updateData.eventID)
		}
		if err == nil && event.Ranges != nil && event.Shooters != nil {
			aggsFound = searchForAggs(event.Ranges, updateData.rangeID)
			if len(aggsFound) > 0 {
				updateBson = make(M)
				for index, data := range calculateAggs(event.Shooters[updateData.shooterID].Scores, aggsFound, []int{updateData.shooterID}, event.Ranges) {
					updateBson[index] = data
				}
				if len(updateBson) > 0 {
					tableUpdateData(tblEvent, event.ID, updateBson)
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
	var updateBson M
	var shooterQty, position, shouldBePosition int
	var shootEqual, updateRequired bool
	var previousEventID, positionEqual, positionOrdinal, strRangeID string

	for _, updateData := range updates {
		//Only get the event when it is different
		if updateData.eventID != previousEventID {
			if updateRequired {
				tableUpdateData(tblEvent, updateData.eventID, updateBson)
				updateRequired = false
			}
			updateBson = make(M)
			event, err = getEvent(updateData.eventID)

			if err != nil {
				warning.Println(err)
				break
			}

			//TODO remove adding shooter ids!
			//Add shooter ids to the shooter objects
			event.Shooters = addShooterIDsToShooterObjects(event.Shooters)

			shooterQty = len(event.Shooters)
		}
		strRangeID = fmt.Sprintf("%v", updateData.rangeID)

		//sort shooters by the current rangeID
		sortShooters(strRangeID).Sort(event.Shooters)

		shouldBePosition = 0
		shootEqual = false
		positionEqual = ""
		for index, shooter := range event.Shooters {
			if shooter.Grade == updateData.gradeID {
				shouldBePosition++
				if !shootEqual {
					position = shouldBePosition
					positionEqual = ""
				} else {
					positionEqual = "="
					shootEqual = false
				}
				if shooter.Scores[strRangeID].ShootOff < 0 {
					//Shooter has the same score as the previous shooter (index-1)
					positionEqual = "="
					if index+1 < shooterQty && shooter.Grade == event.Shooters[index+1].Grade && shooter.Scores[strRangeID] == event.Shooters[index+1].Scores[strRangeID] {
						shootEqual = true
					}
				}
				positionOrdinal = positionEqual + ordinal(position)
				if shooter.Scores[strRangeID].Total != 0 && (shooter.Scores[strRangeID].Position != position || shooter.Scores[strRangeID].Ordinal != positionOrdinal) {
					updateRequired = true
					updateBson[dot("^^schemaSHOOTER^^", shooter.ID, updateData.rangeID, "o")] = positionOrdinal
					updateBson[dot("^^schemaSHOOTER^^", shooter.ID, updateData.rangeID, "p")] = position
				}
			}
		}
	}
	if updateRequired {
		tableUpdateData(tblEvent, event.ID, updateBson)
	}
	info.Println("finished grade recalculation")
}
