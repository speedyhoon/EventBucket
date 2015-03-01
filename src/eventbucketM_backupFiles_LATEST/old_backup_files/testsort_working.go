package main

import (
	"fmt"
	"sort"
)

type Ranges struct {
	total int
	centers int
}

// A Change is a record of source code changes, recording user, language, and delta size.
type Change struct {
//	user     string
//	language string
//	lines    int


	clas string
	grade string
//	_id string
//	club string
	ranges []Ranges
	first string
	surname string
}

type lessFunc func(p1, p2 *Change) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	changes []Change
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(changes []Change) {
//func (ms *multiSorter) Sort(changes map[string]Change) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

//var changes = map[string]Change{
var changes = []Change{
//	"0":Change{"target", "A", []Ranges{{30,9}}},
//	"1":Change{"target", "A", []Ranges{{30,10}}},
//	"2":Change{"target", "A", []Ranges{{30,11}}},
//	"3":Change{"target", "A", []Ranges{{30,14}}},
//	"4":Change{"target", "A", []Ranges{{30,8}}},
//	"5":Change{"target", "A", []Ranges{{60,12}}},
//
//	"10":Change{"target", "B", []Ranges{{30,9}}},
//	"11":Change{"target", "B", []Ranges{{30,10}}},
//	"12":Change{"target", "B", []Ranges{{30,11}}},
//	"13":Change{"target", "B", []Ranges{{30,14}}},
//	"14":Change{"target", "B", []Ranges{{30,8}}},
//	"15":Change{"target", "B", []Ranges{{60,12}}},
//
//	"20":Change{"target", "C", []Ranges{{30,9}}},
//	"21":Change{"target", "C", []Ranges{{30,10}}},
//	"22":Change{"target", "C", []Ranges{{30,11}}},
//	"23":Change{"target", "C", []Ranges{{30,14}}},
//	"24":Change{"target", "C", []Ranges{{30,8}}},
//	"25":Change{"target", "C", []Ranges{{60,12}}},
//
//	"30":Change{"fclass", "A", []Ranges{{30,9}}},
//	"31":Change{"fclass", "A", []Ranges{{30,10}}},
//	"32":Change{"fclass", "A", []Ranges{{30,11}}},
//	"33":Change{"fclass", "A", []Ranges{{30,14}}},
//	"34":Change{"fclass", "A", []Ranges{{30,8}}},
//	"35":Change{"fclass", "A", []Ranges{{60,12}}},
//
//	"40":Change{"fclass", "B", []Ranges{{30,9}}},
//	"41":Change{"fclass", "B", []Ranges{{30,10}}},
//	"42":Change{"fclass", "B", []Ranges{{30,11}}},
//	"43":Change{"fclass", "B", []Ranges{{30,14}}},
//	"44":Change{"fclass", "B", []Ranges{{30,8}}},
//	"45":Change{"fclass", "B", []Ranges{{60,12}}},
//
//	"50":Change{"fclass", "C", []Ranges{{30,9}}},
//	"51":Change{"fclass", "C", []Ranges{{30,10}}},
//	"52":Change{"fclass", "C", []Ranges{{30,11}}},
//	"53":Change{"fclass", "C", []Ranges{{30,14}}},
//	"54":Change{"fclass", "C", []Ranges{{30,8}}},
//	"55":Change{"fclass", "C", []Ranges{{60,12}}},
	{"target", "A", []Ranges{{30,9}}, "fist", "last"},
	{"target", "A", []Ranges{{30,10}}, "fist", "last"},
	{"target", "A", []Ranges{{30,11}}, "fist", "last"},
	{"target", "A", []Ranges{{30,14}}, "fist", "last"},
	{"target", "A", []Ranges{{30,8}}, "fist", "last"},
	{"target", "A", []Ranges{{60,12}}, "fist", "last"},

	{"target", "B", []Ranges{{30,9}}, "fist", "last"},
	{"target", "B", []Ranges{{30,10}}, "fist", "last"},
	{"target", "B", []Ranges{{30,11}}, "fist", "last"},
	{"target", "B", []Ranges{{30,14}}, "fist", "last"},
	{"target", "B", []Ranges{{30,8}}, "fist", "last"},
	{"target", "B", []Ranges{{60,12}}, "fist", "last"},

	{"target", "C", []Ranges{{30,9}}, "fist", "last"},
	{"target", "C", []Ranges{{30,10}}, "fist", "last"},
	{"target", "C", []Ranges{{30,11}}, "fist", "last"},
	{"target", "C", []Ranges{{30,14}}, "fist", "last"},
	{"target", "C", []Ranges{{30,8}}, "fist", "last"},
	{"target", "C", []Ranges{{60,12}}, "fist", "last"},

	{"fclass", "A", []Ranges{{30,9}}, "fist", "last"},
	{"fclass", "A", []Ranges{{30,10}}, "fist", "last"},
	{"fclass", "A", []Ranges{{30,11}}, "fist", "last"},
	{"fclass", "A", []Ranges{{30,14}}, "fist", "last"},
	{"fclass", "A", []Ranges{{30,8}}, "fist", "last"},
	{"fclass", "A", []Ranges{{60,12}}, "fist", "last"},

	{"fclass", "B", []Ranges{{30,9}}, "fist", "last"},
	{"fclass", "B", []Ranges{{30,10}}, "fist", "last"},
	{"fclass", "B", []Ranges{{30,11}}, "fist", "last"},
	{"fclass", "B", []Ranges{{30,14}}, "fist", "last"},
	{"fclass", "B", []Ranges{{30,8}}, "fist", "last"},
	{"fclass", "B", []Ranges{{60,12}}, "fist", "last"},

	{"fclass", "C", []Ranges{{30,9}}, "fist", "last"},
	{"fclass", "C", []Ranges{{30,10}}, "fist", "last"},
	{"fclass", "C", []Ranges{{30,11}}, "fist", "last"},
	{"fclass", "C", []Ranges{{30,14}}, "fist", "last"},
	{"fclass", "C", []Ranges{{30,8}}, "fist", "last"},
	{"fclass", "C", []Ranges{{60,12}}, "fist", "last"},
//
//
//var changes = []Change{
//	{"target", "A", []Ranges{{30,9}}},
//	{"target", "A", []Ranges{{30,10}}},
//	{"target", "B", []Ranges{{30,10}}},
//	{"fclass", "FB", []Ranges{{30,10}}},
//	{"fclass", "FO", []Ranges{{30,10}}},
//	{"fclass", "FO", []Ranges{{60,10}}},



//	{"gri", "Go", 100},
//	{"ken", "C", 150},
//	{"glenda", "Go", 200},
//	{"rsc", "Go", 200},
//	{"r", "Go", 100},
//	{"ken", "Go", 200},
//	{"dmr", "C", 100},
//	{"r", "C", 150},
//	{"gri", "Smalltalk", 80},
}

// ExampleMultiKeys demonstrates a technique for sorting a struct type using different
// sets of multiple fields in the comparison. We chain together "Less" functions, each of
// which compares a single field.
func main() {
	// Closures that order the Change structure.
	total := func(c1, c2 *Change) bool {
		return c1.ranges[0].total > c2.ranges[0].total
	}
	centa := func(c1, c2 *Change) bool {
		return c1.ranges[0].centers > c2.ranges[0].centers
	}
	clas := func(c1, c2 *Change) bool {
		return c1.clas < c2.clas
	}
	grade := func(c1, c2 *Change) bool {
		return c1.grade < c2.grade
	}


	//decreasingLines := func(c1, c2 *Change) bool {
	//	return c1.lines > c2.lines // Note: > orders downwards.
	//}

	// Simple use: Sort by user.
	//OrderedBy(user).Sort(changes)
	//fmt.Println("By user:", changes)
    //
	//// More examples.
	//OrderedBy(user, increasingLines).Sort(changes)
	//fmt.Println("By user,<lines:", changes)
    //
	//OrderedBy(user, decreasingLines).Sort(changes)
	//fmt.Println("By user,>lines:", changes)
    //
	//OrderedBy(language, increasingLines).Sort(changes)
	//fmt.Println("By language,<lines:", changes)

//	OrderedBy(language, increasingLines, user).Sort(changes)


	OrderedBy(clas, grade, total, centa).Sort(changes)
	fmt.Println("By language,<lines,user:")
	for _, item := range changes{
		fmt.Println(fmt.Sprintf("%v",item))
	}

}
