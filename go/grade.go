package main

// Discipline separates different types of shooting so the number of shots & sighters can be easily changed while still using the same targets and Mark as another Discipline, e.g. Target rifles and Match rifles are vastly different disciplines but use the same scoring standard.
type Discipline struct {
	Name                  string
	ID                    uint
	QtySighters, QtyShots uint
	Grades                []Grade
	Marking               Mark
	ShootOff              bool
}

// Mark is a group of settings associated with possible shooter scores on a target also known as "marking". Each type of target scoring standard can be specified by a Mark and be reused within several Disciplines.
type Mark struct {
	Buttons      string
	Shots        map[string]Score
	Sighters     []string
	DoCountBack2 bool
}

var (
	globalDiscipline = defaultGlobalDiscipline()
	globalGrades     = defaultGrades(globalDiscipline)
)

// Grade are subcategories of each discipline that shooters can be grouped together by similar skill levels.
type Grade struct {
	name    string //Target A, F Class B, Match Reserve
	abbr    string //Abbreviation of name: A,B,C,FA,FB,FO,MO,MR
	classID uint
}

func defaultGrades(classes []Discipline) []Grade {
	var grades []Grade
	for classID, class := range classes {
		for _, grade := range class.Grades {
			grade.classID = uint(classID)
			grades = append(grades, grade)
		}
	}
	return grades
}

func defaultGlobalDiscipline() []Discipline {
	XV5 := Mark{Buttons: "012345VX",
		DoCountBack2: true,
		Sighters:     []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
		Shots: map[string]Score{
			"-": {Total: 0, Centers: 0, CountBack: "0", CountBack2: "0"},
			"0": {Total: 0, Centers: 0, CountBack: "0", CountBack2: "0"},
			"1": {Total: 1, Centers: 0, CountBack: "1", CountBack2: "1"},
			"2": {Total: 2, Centers: 0, CountBack: "2", CountBack2: "2"},
			"3": {Total: 3, Centers: 0, CountBack: "3", CountBack2: "3"},
			"4": {Total: 4, Centers: 0, CountBack: "4", CountBack2: "4"},
			"5": {Total: 5, Centers: 0, CountBack: "5", CountBack2: "5"},
			"V": {Total: 5, Centers: 1, CountBack: "6", CountBack2: "6"},
			"6": {Total: 5, Centers: 1, CountBack: "6", CountBack2: "6"},
			"X": {Total: 5, Centers: 1, CountBack: "6", CountBack2: "7"},
		}}
	return []Discipline{{
		ID:          0,
		Name:        "Target Rifle",
		QtySighters: 2,
		QtyShots:    10,
		Marking:     XV5,
		Grades: []Grade{{abbr: "A", name: "Target A"},
			{abbr: "B", name: "Target B"},
			{abbr: "C", name: "Target C"}},
	}, {
		ID:          1,
		Name:        "F Class",
		QtyShots:    12,
		QtySighters: 3,
		Marking: Mark{
			Buttons:  "0123456X",
			Sighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			Shots: map[string]Score{
				"-": {Total: 0, Centers: 0, CountBack: "0"},
				"0": {Total: 0, Centers: 0, CountBack: "0"},
				"1": {Total: 1, Centers: 0, CountBack: "1"},
				"2": {Total: 2, Centers: 0, CountBack: "2"},
				"3": {Total: 3, Centers: 0, CountBack: "3"},
				"4": {Total: 4, Centers: 0, CountBack: "4"},
				"5": {Total: 5, Centers: 0, CountBack: "5"},
				"V": {Total: 6, Centers: 0, CountBack: "6"},
				"6": {Total: 6, Centers: 0, CountBack: "6"},
				"X": {Total: 6, Centers: 1, CountBack: "7"},
			}},
		Grades: []Grade{{abbr: "FA", name: "F Class A"},
			{abbr: "FB", name: "F Class B"},
			{abbr: "FO", name: "F Class Open"},
			{abbr: "F/TR", name: "F/TR"}},
	}, {
		ID:          2,
		Name:        "Match Rifle",
		QtySighters: 3,
		QtyShots:    15,
		Marking:     XV5,
		Grades: []Grade{{abbr: "Open", name: "Match Open"},
			{abbr: "Reserve", name: "Match Reserve"}},
	}, {
		ID:          3,
		Name:        "Service Rifle",
		QtySighters: 1,
		QtyShots:    8,
		Marking: Mark{
			Buttons:  "012345V",
			Sighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			Shots: map[string]Score{
				"-": {Total: 0, Centers: 0, CountBack: "0"},
				"0": {Total: 0, Centers: 0, CountBack: "0"},
				"1": {Total: 1, Centers: 0, CountBack: "1"},
				"2": {Total: 2, Centers: 0, CountBack: "2"},
				"3": {Total: 3, Centers: 0, CountBack: "3"},
				"4": {Total: 4, Centers: 0, CountBack: "4"},
				"5": {Total: 5, Centers: 0, CountBack: "5"},
				"V": {Total: 5, Centers: 1, CountBack: "6"},
				"6": {Total: 5, Centers: 1, CountBack: "6"},
				"X": {Total: 5, Centers: 1, CountBack: "6"},
			}},
		Grades: []Grade{{abbr: "303", name: "303 Rifle"}},
	}}
}
