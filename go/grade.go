package main

// ClassSettings is exported
type ClassSettings struct {
	Name                  string
	ID                    uint64
	DisplayValue          uint64
	Buttons               string
	SightersQty, ShotsQty uint64
	ValidShots            map[string]Score
	ValidSighters         []string
	//GradeQty              int
	Grades  []uint64
	Maximum Score
}

var (
	defaultClassSettings = defaultGlobalClassSettings()
)

// Grade is exported
type Grade struct {
	name, longName, className string
	classID                   int
	settings                  ClassSettings
}

func grades() []Grade {
	return []Grade{
		0: {settings: defaultClassSettings[0], classID: 0, name: "A", className: "Target", longName: "Target A"},
		1: {settings: defaultClassSettings[0], classID: 0, name: "B", className: "Target", longName: "Target B"},
		2: {settings: defaultClassSettings[0], classID: 0, name: "C", className: "Target", longName: "Target C"},
		3: {settings: defaultClassSettings[1], classID: 1, name: "FA", className: "F Class", longName: "F Class A"},
		4: {settings: defaultClassSettings[1], classID: 1, name: "FB", className: "F Class", longName: "F Class B"},
		5: {settings: defaultClassSettings[1], classID: 1, name: "F Open", className: "F Class", longName: "F Class Open"},
		6: {settings: defaultClassSettings[1], classID: 1, name: "F/TR", className: "F Class", longName: "F/TR"},
		7: {settings: defaultClassSettings[2], classID: 2, name: "Open", className: "Match", longName: "Match Open"},
		8: {settings: defaultClassSettings[2], classID: 2, name: "Reserve", className: "Match", longName: "Match Reserve"},
		9: {settings: defaultClassSettings[1], classID: 1, name: "Rifle", className: "303", longName: "303 Rifle"},
	}
}

func defaultGlobalClassSettings() []ClassSettings {
	target := ClassSettings{
		ID:          0,
		Name:        "Target",
		Buttons:     "012345VX",
		SightersQty: 2,
		ShotsQty:    10,
		Maximum:     Score{Total: 5, Centers: 1},
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
		},
		ValidSighters: []string{")", "!", "@", "#", "$", "%", "v", "^", "x"},
		Grades:        []uint64{0, 1, 2, 9},
	}
	fclass := target
	fclass.ID = 1
	fclass.Name = "F Class"
	fclass.Buttons = "0123456X"
	fclass.Maximum = Score{Total: 6, Centers: 1}
	fclass.ValidShots["V"] = Score{Total: 5, Centers: 0, CountBack: "6"}
	fclass.ValidShots["6"] = Score{Total: 6, Centers: 0, CountBack: "6"}
	fclass.ValidShots["X"] = Score{Total: 6, Centers: 1, CountBack: "7"}
	fclass.Grades = []uint64{3, 4, 5, 6}

	match := target
	match.ID = 2
	match.Name = "Match"
	match.SightersQty = 3
	match.ShotsQty = 15
	match.Grades = []uint64{7, 8}
	return []ClassSettings{target, fclass, match}
}
