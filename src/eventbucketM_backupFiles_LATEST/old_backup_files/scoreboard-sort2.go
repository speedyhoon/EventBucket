package main

import (
	"fmt"
	"sort"
)

type Score struct{
	Total int									`bson:"t"`
	Shots string								`bson:"s,omitempty"`
	Centers int									`bson:"c"`
	CountBack1 string							`bson:"v,omitempty"`
	CountBack2 string							`bson:"x,omitempty"`
}
type EventShooter struct{
	FirstName string							`bson:"f"`
	Surname string								`bson:"s"`
	Club string									`bson:"b"`//TODO should possibly change to "C"??
	Class string								`bson:"c"`
	Grade string								`bson:"g"`
	AgeGroup string							`bson:"a,omitempty"`
	Scores map[string]Score					`bson:"S,omitempty,inline"` //S is not used!
}
type Event struct{
	Id string									`bson:"_id"`
	Club string									`bson:"c"`
	Name string									`bson:"n"`
	Datetime string							`bson:"d,omitempty"`
	Shooters map[string]EventShooter		`bson:"S,omitempty"`
}

type lessFunc func(p1, p2 *EventShooter) bool

type multiSorter struct {
	changes []EventShooter
	less    []lessFunc
}

func (ms *multiSorter) Sort(changes []EventShooter) {
	ms.changes = changes
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}

func main() {
	sortByRange := "4"

	// Closures that order the Change structure.
	clas := func(c1, c2 *EventShooter) bool {
		return c1.Class > c2.Class
	}
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	total := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].Total > c2.Scores[sortByRange].Total
	}
	centa := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].Centers > c2.Scores[sortByRange].Centers
	}
	cb := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].CountBack1 > c2.Scores[sortByRange].CountBack1
	}

	var changes = []EventShooter{}
	for _, shooterList := range event.Shooters {
		changes = append(changes, shooterList)
	}

	OrderedBy(clas, grade, total, centa, cb).Sort(changes)
	for _, item := range changes{
		fmt.Println(fmt.Sprintf("%v",item))
	}

	fmt.Println(len(changes))
}







var event = Event{
	Shooters:	map[string]EventShooter{
		"0": EventShooter{
			FirstName: "Cam",
			Surname: "W",
			//			MiddleName: "F",
			//			Age: "",
			Club: "1",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 160,
					Centers: 18,
					CountBack1: "1234567890",
				},
			},
		},
		"1": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 160,
					Centers: 19,
					CountBack1: "1234567890",
				},
			},
		},
		"2": EventShooter{
			FirstName: "Sam",
			Surname: "N",
			//			MiddleName: "Smith",
			//			Age: "",
			Club: "2",
			Class: "fclass",
			Grade: "FB",
			Scores: map[string]Score{
				"4":	Score{
					Total: 200,
					Centers: 25,
					CountBack1: "1234567890",
				},
			},
		},
		"3": EventShooter{
			FirstName: "Larry",
			Surname: "N",
			//			MiddleName: "Perkins",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "A",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 13,
					CountBack1: "1234567890",
				},
			},
		},
		"4": EventShooter{
			FirstName: "Peter",
			Surname: "N",
			//			MiddleName: "Brock",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "A",
			Scores: map[string]Score{
				"4":	Score{
					Total: 160,
					Centers: 14,
					CountBack1: "1234567890",
				},
			},
		},
		"5": EventShooter{
			FirstName: "Darrel",
			Surname: "N",
			//			MiddleName: "Nepia",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 158,
					Centers: 6,
					CountBack1: "1234567890",
				},
			},
		},
		"6": EventShooter{
			FirstName: "Raymond",
			Surname: "N",
			//			MiddleName: "Tapps",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "C",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 5,
					CountBack1: "1234567890",
				},
			},
		},
		"7": EventShooter{
			FirstName: "Leo",
			Surname: "N",
			//			MiddleName: "Stigger",
			//			Age: "U21",
			Club: "2",
			Class: "fclass",
			Grade: "FA",
			Scores: map[string]Score{
				"4":	Score{
					Total: 201,
					Centers: 3,
					CountBack1: "1234567890",
				},
			},
		},
		"8": EventShooter{
			FirstName: "Tim",
			Surname: "N",
			//			MiddleName: "Murray",
			//			Age: "U21",
			Club: "2",
			Class: "fclass",
			Grade: "FA",
			Scores: map[string]Score{
				"4":	Score{
					Total: 199,
					Centers: 16,
					CountBack1: "1234567890",
				},
			},
		},
		"9": EventShooter{
			FirstName: "Garry",
			Surname: "N",
			//			MiddleName: "Green",
			//			Age: "U21",
			Club: "2",
			Class: "fclass",
			Grade: "FB",
			Scores: map[string]Score{
				"4":	Score{
					Total: 95,
					Centers: 40,
					CountBack1: "1234567890",
				},
			},
		},
		"10": EventShooter{
			FirstName: "Fred",
			Surname: "N",
			//			MiddleName: "Webb",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 187,
					Centers: 60,
					CountBack1: "1234567890",
				},
			},
		},
		"11": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1034567894",
				},
			},
		},
		"12": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567855",
				},
			},
		},
		"13": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
				},
			},
		},
		"14": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567990",
				},
			},
		},
		"15": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567897",
				},
			},
		},
		"16": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "2134567890",
				},
			},
		},
		"17": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567809",
				},
			},
		},
		"18": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234557890",
				},
			},
		},
		"19": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567893",
				},
			},
		},
		"20": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "0234567890",
				},
			},
		},
		"21": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567892",
				},
			},
		},
		"22": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567891",
				},
			},
		},
		"23": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567890",
				},
			},
		},
		"24": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567890B",
				},
			},
		},
		"25": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567890A",
				},
			},
		},
		"26": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567890b",
				},
			},
		},
		"27": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1234567890a",
				},
			},
		},
		"28": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "1",
				},
			},
		},
		"29": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "2",
				},
			},
		},
		"30": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "3",
				},
			},
		},
		"31": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "4",
				},
			},
		},
		"32": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "5",
				},
			},
		},
		"33": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "6",
				},
			},
		},
		"34": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "7",
				},
			},
		},
		"35": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "8",
				},
			},
		},
		"36": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "9",
				},
			},
		},
		"37": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "10",
				},
			},
		},
		"38": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "11",
				},
			},
		},
		"39": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "12",
				},
			},
		},
		"40": EventShooter{
			FirstName: "Leesa",
			Surname: "N",
			//			MiddleName: "M",
			//			Age: "U21",
			Club: "2",
			Class: "target",
			Grade: "B",
			Scores: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					CountBack1: "13",
				},
			},
		},
	},
}
