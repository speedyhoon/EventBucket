package main

import (
	"fmt"
	"math"
	"sort"
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
	t.Println(c1.Scores[rangeID].ShootOff > c2.Scores[rangeID].ShootOff, c1.Scores[rangeID].ShootOff, c2.Scores[rangeID].ShootOff)
	return c1.Scores[rangeID].ShootOff > c2.Scores[rangeID].ShootOff
}

func sortShooters(rangeID string) *multiEventShooterSorter {
	if rangeID != "" {
		return orderShooters(rangeID, sorterGrade, sorterTotal, sorterCenters, sorterCountBack, sorterCountBack2, sorterShootOff)
	}
	return &multiEventShooterSorter{}
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

//TODO change the functions to return -1, 0, 1 and reduce the number of calls for greater efficiency
func (ms *multiEventShooterSorter) Less(i, j int) bool {
	p, q := &ms.shooter[i], &ms.shooter[j]
	//Try all but the last comparison.
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

//TODO don't add grade separators for database range sorting updates
func addGradeSeparatorToShooterObjectAndPositions(eventShooters []EventShooter, rangeID string) []EventShooter {
	//Add a boolean field to each shooter in a list of ordered shooters and is true for the first shooter that has a different grade than the last
	var previousShooterGrade uint = math.MaxUint32

	var position, index uint
	var previousShooter int
	var previousScore, score Score
	var shooterTie, sameGrade bool

	//Loop through each shooter
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

		//Check if shooters grades and scores are the same
		shooterTie = sameGrade && score.Total == previousScore.Total && score.Centers == previousScore.Centers && score.Centers2 == previousScore.Centers2 && score.CountBack == previousScore.CountBack && score.CountBack2 == previousScore.CountBack2 && score.ShootOff == previousScore.ShootOff
		if shooterTie {
			previousScore.Ordinal = ordinal(position, shooterTie)
			eventShooters[previousShooter].Scores[rangeID] = previousScore
		} else {
			position++
		}

		score.Position = position
		score.Ordinal = ordinal(position, shooterTie)
		if eventShooters[shooterID].Scores == nil {
			eventShooters[shooterID].Scores = make(map[string]Score, 1)
		}
		eventShooters[shooterID].Scores[rangeID] = score
		previousShooter = shooterID
		previousScore = score
	}
	return eventShooters
}

//Ordinal gives you the input number in a rank/ordinal format. e.g. Ordinal(3, true) outputs "=3rd"
func ordinal(position uint, equal bool) string {
	suffix := "th"
	switch position % 10 {
	case 1:
		if position%100 != 11 {
			suffix = "st"
		}
	case 2:
		if position%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if position%100 != 13 {
			suffix = "rd"
		}
	}
	var sign string
	if equal {
		sign = "="
	}
	return fmt.Sprintf("%v%d%v", sign, position, suffix)
}
