package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/speedyhoon/utl"
)

type sortEventShooter func(r1 string, p1, p2 *EventShooter) bool

type multiEventShooterSorter struct {
	shooter []EventShooter
	sort    []sortEventShooter
	rangeID string
}

func sorterGrade(_ string, c1, c2 *EventShooter) bool {
	return c1.Grade < c2.Grade
}
func sorterTotal(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].Total > c2.Scores[rangeID].Total
}
func sorterCenters(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].Centers > c2.Scores[rangeID].Centers
}
func sorterCountBack(rangeID string, c1, c2 *EventShooter) bool {
	return c1.Scores[rangeID].CountBack > c2.Scores[rangeID].CountBack
}
func sorterCountBack2(rangeID string, c1, c2 *EventShooter) bool {
	//t.Println("sorterCountBack2", rangeID, c1.ID, c1.FirstName, c1.Surname, c2.ID, c2.FirstName, c2.Surname)
	return c1.Scores[rangeID].CountBack2 > c2.Scores[rangeID].CountBack2
}

/*func sorterFirstName(rangeID string, c1, c2 *EventShooter) bool {
	return c1.FirstName < c2.FirstName
}*/

func sorterShootOff(rangeID string, c1, c2 *EventShooter) bool {
	//t.Println("sorterShootOff", rangeID, c1.ID, c1.FirstName, c1.Surname, c2.ID, c2.FirstName, c2.Surname)
	/*if c1.Scores[rangeID].CountBack != "" && c1.Scores[rangeID].CountBack == c2.Scores[rangeID].CountBack {
		t.Printf("shooters scores are the same? c1= g:%v t:%v c:%v b:%v h:%v c2= g:%v t:%v c:%v b:%v h:%v", c1.Grade, c1.Scores[rangeID].Total, c1.Scores[rangeID].Centers, c1.Scores[rangeID].CountBack, c1.Scores[rangeID].ShootOff, c2.Grade, c2.Scores[rangeID].Total, c2.Scores[rangeID].Centers, c2.Scores[rangeID].CountBack, c2.Scores[rangeID].ShootOff)
		temp := c1.Scores[rangeID]
		//temp.Warning = legendShootOff
		c1.Scores[rangeID] = temp

		temp = c2.Scores[rangeID]
		//temp.Warning = legendShootOff
		c2.Scores[rangeID] = temp
	}*/
	trc.Println(c1.Scores[rangeID].ShootOff > c2.Scores[rangeID].ShootOff, c1.Scores[rangeID].ShootOff, c2.Scores[rangeID].ShootOff)
	return c1.Scores[rangeID].ShootOff > c2.Scores[rangeID].ShootOff
}

func sortShooters(rangeID rID) *multiEventShooterSorter {
	return orderShooters(fmt.Sprintf("%v", rangeID), sorterGrade, sorterTotal, sorterCenters, sorterCountBack, sorterCountBack2, sorterShootOff)
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
	ms.shooter = addGradeSeparatorToShooterObjectAndPositions(ms.shooter, ms.rangeID)
}

func (ms *multiEventShooterSorter) Len() int {
	return len(ms.shooter)
}

func (ms *multiEventShooterSorter) Swap(i, j int) {
	ms.shooter[i], ms.shooter[j] = ms.shooter[j], ms.shooter[i]
}

// TODO change the functions to return -1, 0, 1 and reduce the number of calls for greater efficiency.
func (ms *multiEventShooterSorter) Less(i, j int) bool {
	p, q := &ms.shooter[i], &ms.shooter[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.sort)-1; k++ {
		sortFunc := ms.sort[k]
		switch {
		case sortFunc(ms.rangeID, p, q):
			return true
		case sortFunc(ms.rangeID, q, p):
			return false
		}
	}
	return ms.sort[k](ms.rangeID, p, q)
}

// TODO don't add grade separators for database range sorting updates.
func addGradeSeparatorToShooterObjectAndPositions(eventShooters []EventShooter, rangeID string) []EventShooter {
	// Add a boolean field to each shooter in a list of ordered shooters and is true for the first shooter that has a different grade than the last.
	var previousShooterGrade uint = math.MaxUint32

	var position, index uint
	var previousShooter int
	var previousScore, score Score
	var shooterTie, sameGrade bool

	// Loop through each shooter.
	for shooterID := range eventShooters {
		sameGrade = eventShooters[shooterID].Grade == previousShooterGrade
		if !sameGrade {
			eventShooters[shooterID].GradeSeparator = true
			position = 0
			index = 0
			previousShooterGrade = eventShooters[shooterID].Grade
		}

		score = eventShooters[shooterID].Scores[rangeID]

		index++

		// Check if shooters grades and scores are the same.
		shooterTie = sameGrade && score.Total == previousScore.Total && score.Centers == previousScore.Centers && score.Centers2 == previousScore.Centers2 && score.CountBack == previousScore.CountBack && score.CountBack2 == previousScore.CountBack2 && score.ShootOff == previousScore.ShootOff
		if shooterTie {
			previousScore.Ordinal = utl.Ordinal(position, shooterTie)
			eventShooters[previousShooter].Scores[rangeID] = previousScore
		} else {
			position++
		}

		score.Position = position
		score.Ordinal = utl.Ordinal(position, shooterTie)
		if eventShooters[shooterID].Scores == nil {
			eventShooters[shooterID].Scores = make(ScoreMap, 1)
		}
		eventShooters[shooterID].Scores[rangeID] = score
		previousShooter = shooterID
		previousScore = score
	}
	return eventShooters
}
