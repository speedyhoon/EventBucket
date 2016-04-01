package main

// ClassSettings is exported
type ClassSettings struct {
	Name                  string
	ID                    uint
	QtySighters, QtyShots uint
	Grades                []Grade
	Shotput               Shotput
	ShootOff              bool
}

//TODO rename to something useful
type Shotput struct {
	Buttons       string
	ValidShots    map[string]Score
	ValidSighters []string
}

var (
	globalClassSettings = defaultGlobalClassSettings()
	globalGrades        = defaultGrades(globalClassSettings)
)

// Grade is exported
type Grade struct {
	name    string //Target A, F Class B, Match Reserve
	short   string //A,B,C,FA,FB,FO,MO,MR
	classID uint
}

func defaultGrades(classes []ClassSettings) []Grade {
	var grades []Grade
	for classID, class := range classes {
		for _, grade := range class.Grades {
			grade.classID = uint(classID)
			grades = append(grades, grade)
		}
	}
	return grades
}

func defaultGlobalClassSettings() []ClassSettings {
	V5 := Shotput{
		Buttons:       "012345V",
		ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
		ValidShots: map[string]Score{
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
		}}
	return []ClassSettings{{
		ID:          0,
		Name:        "Target Rifle",
		QtySighters: 2,
		QtyShots:    10,
		Shotput: Shotput{Buttons: "012345VX",
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			ValidShots: map[string]Score{
				"-": {Total: 0, Centers: 0, CountBack: "0"},
				"0": {Total: 0, Centers: 0, CountBack: "0"},
				"1": {Total: 1, Centers: 0, CountBack: "1"},
				"2": {Total: 2, Centers: 0, CountBack: "2"},
				"3": {Total: 3, Centers: 0, CountBack: "3"},
				"4": {Total: 4, Centers: 0, CountBack: "4"},
				"5": {Total: 5, Centers: 0, CountBack: "5"},
				"V": {Total: 5, Centers: 1, CountBack: "6"},
				"6": {Total: 5, Centers: 1, CountBack: "6"},
				"X": {Total: 5, Centers: 1, CountBack: "7"},
			}},
		Grades: []Grade{{short: "A", name: "Target A"},
			{short: "B", name: "Target B"},
			{short: "C", name: "Target C"}},
	}, {
		ID:          1,
		Name:        "F Class",
		QtyShots:    12,
		QtySighters: 3,
		Shotput: Shotput{
			Buttons:       "0123456X",
			ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
			ValidShots: map[string]Score{
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
		Grades: []Grade{{short: "FA", name: "F Class A"},
			{short: "FB", name: "F Class B"},
			{short: "FO", name: "F Class Open"},
			{short: "F/TR", name: "F/TR"}},
	}, {
		ID:          2,
		Name:        "Match Rifle",
		QtySighters: 3,
		QtyShots:    15,
		Shotput:     V5,
		Grades: []Grade{{short: "Open", name: "Match Open"},
			{short: "Reserve", name: "Match Reserve"}},
	}, {
		ID:          3,
		Name:        "Service Rifle",
		QtySighters: 1,
		QtyShots:    8,
		Shotput:     V5,
		Grades:      []Grade{{short: "303", name: "303 Rifle"}},
	}}
}
