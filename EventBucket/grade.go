package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

//Discipline separates different types of shooting so the number of shots & sighters can be easily changed while still using the same targets and Mark as another Discipline, e.g. Target rifles and Match rifles are vastly different disciplines but use the same scoring standard.
type Discipline struct {
	Name        string  `json:"name,omitempty"`
	ID          uint    `json:"id,omitempty"`
	QtySighters uint    `json:"sightersQty,omitempty"`
	QtyShots    uint    `json:"shotsQty,omitempty"`
	Colspan     uint    `json:"-"`
	Grades      []Grade `json:"grades,omitempty"`
	Marking     Mark    `json:"marking,omitempty"`
	ShootOff    bool    `json:"shootOff,omitempty"`
	TopShot     uint    `json:"-"`
	TopTotal    uint    `json:"-"`
}

//Mark is a group of settings associated with possible shooter scores on a target also known as "marking". Each type of target scoring standard can be specified by a Mark and be reused within several Disciplines.
type Mark struct {
	Buttons      string          `json:"buttons,omitempty"`
	Shots        map[string]Shot `json:"shots,omitempty"`
	DoCountBack2 bool            `json:"doCountBack2,omitempty"`
}

//Grade are subcategories of each discipline that shooters can be grouped together by similar skill levels.
type Grade struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"` //Target A, F Class B, Match Reserve
	Abbr        string `json:"abbr,omitempty"` //Abbreviation of Name: A,B,C,FA,FB,FO,MO,MR
	ClassID     uint   `json:"classID,omitempty"`
	DuplicateTo []uint `json:"duplicateTo,omitempty"`
}

//Shot is exported
type Shot struct {
	Value      uint   `json:"value"`
	Center     uint   `json:"center,omitempty"`
	Center2    uint   `json:"center2,omitempty"`
	Shot       string `json:"shot,omitempty"`
	Sighter    string `json:"sighter,omitempty"`
	CountBack  string `json:"countBack,omitempty"`
	CountBack2 string `json:"countBack2,omitempty"`
}

var (
	globalDisciplines     []Discipline
	globalGrades          []Grade
	globalGradesDataList  []option
	globalAvailableGrades []option
)

func redoGlobals(disciplines []Discipline) {
	if len(disciplines) > 0 {
		globalDisciplines = disciplines
	} else {
		globalDisciplines = defaultGlobalDisciplines()
	}
	globalGrades = defaultGrades(globalDisciplines)
	globalGradesDataList = dataListGrades(globalGrades)
	globalAvailableGrades = availableGrades([]uint{})
}

func defaultGrades(classes []Discipline) []Grade {
	var grades []Grade
	for classID, class := range classes {
		for _, grade := range class.Grades {
			grade.ClassID = uint(classID)
			grades = append(grades, grade)
		}
	}
	return grades
}

func dataListGrades(grades []Grade) []option {
	var options []option
	for id, grade := range grades {
		options = append(options, option{Value: fmt.Sprintf("%d", id), Label: grade.Name})
	}
	return options
}

func availableGrades(grades []uint) []option {
	var options []option
	for id, grade := range globalGrades {
		selected := len(grades) == 0
		if !selected {
			for _, gradeID := range grades {
				if uint(id) == gradeID {
					selected = true
					break
				}
			}
		}
		options = append(options, option{Value: fmt.Sprintf("%d", id), Label: grade.Name, Selected: selected})
	}
	return options
}

func eventGrades(grades []uint) []option {
	if len(grades) == 0 {
		return globalGradesDataList
	}
	var options []option
	for id, grade := range globalGrades {
		for _, gradeID := range grades {
			if uint(id) == gradeID {
				options = append(options, option{Value: fmt.Sprintf("%d", id), Label: grade.Name})
				break
			}
		}
	}
	return options
}

func (d *Discipline) QtyTotal() uint {
	return d.QtySighters + d.QtyShots
}

func defaultGlobalDisciplines() []Discipline {
	XV5 := Mark{Buttons: "012345VX",
		DoCountBack2: true,
		Shots: map[string]Shot{
			"-": {Value: 0, Center: 0, CountBack: "0", Center2: 0, CountBack2: "0", Shot: "-", Sighter: "-"},
			"0": {Value: 0, Center: 0, CountBack: "0", Center2: 0, CountBack2: "0", Shot: "0", Sighter: "a"},
			"1": {Value: 1, Center: 0, CountBack: "1", Center2: 0, CountBack2: "1", Shot: "1", Sighter: "b"},
			"2": {Value: 2, Center: 0, CountBack: "2", Center2: 0, CountBack2: "2", Shot: "2", Sighter: "c"},
			"3": {Value: 3, Center: 0, CountBack: "3", Center2: 0, CountBack2: "3", Shot: "3", Sighter: "d"},
			"4": {Value: 4, Center: 0, CountBack: "4", Center2: 0, CountBack2: "4", Shot: "4", Sighter: "e"},
			"5": {Value: 5, Center: 0, CountBack: "5", Center2: 0, CountBack2: "5", Shot: "5", Sighter: "f"},
			"V": {Value: 5, Center: 1, CountBack: "6", Center2: 0, CountBack2: "6", Shot: "V", Sighter: "v"},
			"6": {Value: 5, Center: 1, CountBack: "6", Center2: 0, CountBack2: "6", Shot: "V", Sighter: "g"},
			"X": {Value: 5, Center: 1, CountBack: "6", Center2: 1, CountBack2: "7", Shot: "X", Sighter: "x"},
			"a": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
			"b": {Value: 0, Center: 0, CountBack: "0", Shot: "1", Sighter: "b"},
			"c": {Value: 0, Center: 0, CountBack: "0", Shot: "2", Sighter: "c"},
			"d": {Value: 0, Center: 0, CountBack: "0", Shot: "3", Sighter: "d"},
			"e": {Value: 0, Center: 0, CountBack: "0", Shot: "4", Sighter: "e"},
			"f": {Value: 0, Center: 0, CountBack: "0", Shot: "5", Sighter: "f"},
			"g": {Value: 0, Center: 0, CountBack: "0", Shot: "6", Sighter: "g"},
			"v": {Value: 0, Center: 0, CountBack: "0", Shot: "V", Sighter: "v"},
			"x": {Value: 0, Center: 0, CountBack: "0", Shot: "X", Sighter: "x"},
			//TODO sort isn't sorting by countback 2 descending.
			//TODO precedence is taken over the last X shot rather than the most X's shot
		}}
	/*Disciplines F Standard, F Open and F/TR have been merged together because they all use the same scoring method (0123456X). Although they are technically separate disciplines, most events assign the same number of shots and sighters for all three.  If the disciplines need to be independent, these settings can be overwritten using the command line flag -grades and specifying a JSON settings file. e.g. EventBucket.exe -grades my_new_grades.json
	EventBucket will not import or remember the settings file next time you start the application. Adding command line flags to an EventBucket shortcut is the easiest way to specify settings every time EventBucket is started.*/
	return []Discipline{{
		ID:          0,
		Name:        "Target Rifle",
		QtySighters: 2,
		QtyShots:    10,
		TopShot:     5,
		TopTotal:    50,
		Marking:     XV5,
		//Target rifle is traditionally scored up to 5 (bullseye) which is has a larger area than 6 on an F class target.
		//This causes significantly more shoot-offs for winning a range than F Class.
		Grades: []Grade{{ID: 0, Abbr: "A", Name: "Target A"},
			{ID: 1, Abbr: "B", Name: "Target B"},
			{ID: 2, Abbr: "C", Name: "Target C"}},
	}, {
		ID:          1,
		Name:        "F Class",
		QtyShots:    10,
		QtySighters: 2,
		TopShot:     6,
		TopTotal:    72,
		Marking: Mark{
			Buttons: "0123456X",
			Shots: map[string]Shot{
				"-": {Value: 0, Center: 0, CountBack: "0", Shot: "-", Sighter: "-"},
				"0": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
				"1": {Value: 1, Center: 0, CountBack: "1", Shot: "1", Sighter: "b"},
				"2": {Value: 2, Center: 0, CountBack: "2", Shot: "2", Sighter: "c"},
				"3": {Value: 3, Center: 0, CountBack: "3", Shot: "3", Sighter: "d"},
				"4": {Value: 4, Center: 0, CountBack: "4", Shot: "4", Sighter: "e"},
				"5": {Value: 5, Center: 0, CountBack: "5", Shot: "5", Sighter: "f"},
				"V": {Value: 6, Center: 0, CountBack: "6", Shot: "6", Sighter: "g"},
				"6": {Value: 6, Center: 0, CountBack: "6", Shot: "6", Sighter: "g"},
				"X": {Value: 6, Center: 1, CountBack: "7", Shot: "X", Sighter: "x"},
				"a": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
				"b": {Value: 0, Center: 0, CountBack: "0", Shot: "1", Sighter: "b"},
				"c": {Value: 0, Center: 0, CountBack: "0", Shot: "2", Sighter: "c"},
				"d": {Value: 0, Center: 0, CountBack: "0", Shot: "3", Sighter: "d"},
				"e": {Value: 0, Center: 0, CountBack: "0", Shot: "4", Sighter: "e"},
				"f": {Value: 0, Center: 0, CountBack: "0", Shot: "5", Sighter: "f"},
				"g": {Value: 0, Center: 0, CountBack: "0", Shot: "6", Sighter: "g"},
				"v": {Value: 0, Center: 0, CountBack: "0", Shot: "V", Sighter: "v"},
				"x": {Value: 0, Center: 0, CountBack: "0", Shot: "X", Sighter: "x"},
			}},
		Grades: []Grade{{ID: 3, Abbr: "FA", Name: "F Standard A"},
			{ID: 4, Abbr: "FB", Name: "F Standard B"},
			{ID: 5, Abbr: "FO", Name: "F Open"},
			{ID: 6, Abbr: "FTR", Name: "F/TR"}},
	}, {
		ID:          2,
		Name:        "Match Rifle",
		QtySighters: 3,
		QtyShots:    15,
		TopShot:     5,
		TopTotal:    75,
		Marking:     XV5,
		Grades: []Grade{{ID: 7, Abbr: "MO", Name: "Match Open"},
			{ID: 8, Abbr: "MR", Name: "Match Reserve", DuplicateTo: []uint{7}}}, //If shooter is Match Reserve, duplicate them in the Match Open category. Used for Victorian Match Rifle Championships.
	}, {
		ID:          3,
		Name:        "Service Rifle",
		QtySighters: 1,
		QtyShots:    8,
		TopShot:     5,
		TopTotal:    45,
		Marking: Mark{
			Buttons: "012345V",
			Shots: map[string]Shot{
				"-": {Value: 0, Center: 0, CountBack: "0", Shot: "-", Sighter: "-"},
				"0": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
				"1": {Value: 1, Center: 0, CountBack: "1", Shot: "1", Sighter: "b"},
				"2": {Value: 2, Center: 0, CountBack: "2", Shot: "2", Sighter: "c"},
				"3": {Value: 3, Center: 0, CountBack: "3", Shot: "3", Sighter: "d"},
				"4": {Value: 4, Center: 0, CountBack: "4", Shot: "4", Sighter: "e"},
				"5": {Value: 5, Center: 0, CountBack: "5", Shot: "5", Sighter: "f"},
				"V": {Value: 5, Center: 1, CountBack: "6", Shot: "V", Sighter: "v"},
				"6": {Value: 5, Center: 1, CountBack: "6", Shot: "V", Sighter: "v"},
				"X": {Value: 5, Center: 1, CountBack: "6", Shot: "V", Sighter: "v"},
				"a": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
				"b": {Value: 0, Center: 0, CountBack: "0", Shot: "1", Sighter: "b"},
				"c": {Value: 0, Center: 0, CountBack: "0", Shot: "2", Sighter: "c"},
				"d": {Value: 0, Center: 0, CountBack: "0", Shot: "3", Sighter: "d"},
				"e": {Value: 0, Center: 0, CountBack: "0", Shot: "4", Sighter: "e"},
				"f": {Value: 0, Center: 0, CountBack: "0", Shot: "5", Sighter: "f"},
				"g": {Value: 0, Center: 0, CountBack: "0", Shot: "6", Sighter: "g"},
				"v": {Value: 0, Center: 0, CountBack: "0", Shot: "V", Sighter: "v"},
				"x": {Value: 0, Center: 0, CountBack: "0", Shot: "X", Sighter: "x"},
			}},
		Grades: []Grade{{ID: 9, Abbr: "303", Name: "303 Rifle"}},
	}}
}

func loadGrades(filePath string) error {
	if filePath == "" {
		return errors.New("File specified is empty")
	}
	contents, err := ioutil.ReadFile(filePath)
	//If file is empty, try to write a new JSON file.
	if err != nil {
		warn.Println(err)
		return err
	}
	var disciplines []Discipline
	err = json.Unmarshal(contents, &disciplines)
	if err != nil {
		//Unable to unmarshal settings from JSON file.
		warn.Println("error:", err)
		return fmt.Errorf("Error: %v, File: %v", err, filePath)
	}
	info.Println("Loaded grade settings from:", filePath)
	redoGlobals(disciplines)
	return nil
}

func buildGradesFile(filePath string) error {
	//Generate JSON from globalDisciplines
	src, err := json.MarshalIndent(globalDisciplines, "", "\t")
	if err != nil {
		//Output marshal errors
		warn.Println("error:", err)
	}
	if !strings.HasSuffix(filePath, ".json") {
		filePath += ".json"
	}

	err = ioutil.WriteFile(filePath, src, 0777)
	if err != nil {
		warn.Println(err, "Unable to write to file", filePath)
		return err
	}
	info.Println("Created grades settings file:", filePath)
	return nil
}

func findGrade(index uint) Grade {
	if index < uint(len(globalGrades)) {
		return globalGrades[index]
	}
	return Grade{}
}
