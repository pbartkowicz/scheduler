package university

import (
	"errors"
	"strconv"
	"time"
)

var (
	// ErrWrongClassType is returned when a passed class type is incorrect.
	ErrWrongClassType = errors.New("incorrect class type, available types: Class, Lecture, Laboratory")
	// ErrWrongWeekday is returned when a passed weekday is incorrect.
	ErrWrongWeekday = errors.New("incorrect weekday, available weekdays: Monday, Tuesday, Wednesday, Thursday, Friday")
)

// GroupError represents an error struct returned when creating new Group.
type GroupError struct {
	Err error
}

func (e *GroupError) Error() string {
	return "failed to create group: " + e.Err.Error()
}

// ClassType defines type of a class.
type ClassType string

const (
	// Class - class.
	Class ClassType = "Class"
	// Lecture - lecture.
	Lecture ClassType = "Lecture"
	// Laboratory - laboratory.
	Laboratory ClassType = "Laboratory"
)

var types = map[string]ClassType{
	"Class":      Class,
	"Lecture":    Lecture,
	"Laboratory": Laboratory,
}

var weekdays = map[string]time.Weekday{
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
}

// Group represents a single students group for one subject.
// It implements sort.Interface based on students' happiness in a slice containing students.
type Group struct {
	Type             ClassType
	Teacher          string
	Weekday          time.Weekday
	StartTime        time.Time
	EndTime          time.Time
	Place            string
	StartDate        time.Time
	Frequency        int
	Name             string
	Capacity         int
	Students         []*Student
	PriorityStudents []*Student
	SubGroups        []*Group
}

func (g *Group) Len() int {
	return len(g.Students)
}

func (g *Group) Less(i, j int) bool {
	return g.Students[i].GetHappiness() > g.Students[j].GetHappiness()
}

func (g *Group) Swap(i, j int) {
	g.Students[i], g.Students[j] = g.Students[j], g.Students[i]
}

// NewGroup creates a new instance of Group.
// It returns GroupError when passed parameters are invalid.
// subjects:
// 0 - subject name
// 1 - class type [Class, Lecture, Laboratory]
// 2 - teacher
// 3 - weekday [Monday, Tuesday, Wednesday, Thursday, Friday]
// 4 - start time, format: 14:00
// 5 - end time, format: 15:30
// 6 - place
// 7 - start date, format: 03-05-20 (5th of March 2020)
// 8 - frequency, format: number
// 9 - group name
// 10 - capacity, format: number
func NewGroup(subjects []string) (*Group, error) {
	t := types[subjects[1]]
	if t == "" {
		return nil, &GroupError{Err: ErrWrongClassType}
	}

	w := weekdays[subjects[3]]
	if w == 0 {
		return nil, &GroupError{Err: ErrWrongWeekday}
	}

	st, err := time.Parse(timeLayout, subjects[4])
	if err != nil {
		return nil, &GroupError{Err: err}
	}
	et, err := time.Parse(timeLayout, subjects[5])
	if err != nil {
		return nil, &GroupError{Err: err}
	}

	d, err := time.Parse(dateLayout, subjects[7])
	if err != nil {
		return nil, &GroupError{Err: err}
	}

	f, err := strconv.Atoi(subjects[8])
	if err != nil {
		return nil, &GroupError{Err: err}
	}
	c, err := strconv.Atoi(subjects[10])
	if err != nil {
		return nil, &GroupError{Err: err}
	}

	return &Group{
		Type:      t,
		Teacher:   subjects[2],
		Weekday:   time.Weekday(w),
		StartTime: st,
		EndTime:   et,
		Place:     subjects[6],
		StartDate: d,
		Frequency: f,
		Name:      subjects[9],
		Capacity:  c,
	}, nil
}

// Conflicts calculate the number of conflicts within a group.
// Conflict exists when the number of students in a group is too high.
func (g *Group) Conflicts() int {
	return (len(g.PriorityStudents) + len(g.Students)) - g.Capacity
}

// Collide checks if groups are not in the same time.
func (g *Group) Collide(a *Group) bool {
	// TODO: groups can be twice a week, so they can be on the same weekday and at the same hour
	// Frequency & start date

	if g.Weekday != a.Weekday {
		return false
	}
	// There is a fixed time schedule at the university
	if g.StartTime != a.StartTime {
		return false
	}
	return true
}

// RemoveStudent removes student from group.
func (g *Group) RemoveStudent(st *Student) {
	newStudents := []*Student{}
	for _, s := range g.Students {
		if st.Name != s.Name {
			newStudents = append(newStudents, s)
		}
	}
	g.Students = newStudents
}

// Save creates a slice with students who will attend this group.
func (g *Group) Save() [][]string {
	res := make([][]string, len(g.PriorityStudents)+len(g.Students))
	var i int
	saveStudents(g.PriorityStudents, res, &i)
	saveStudents(g.Students, res, &i)
	return res
}

func saveStudents(sts []*Student, res [][]string, i *int) {
	for _, st := range sts {
		r := make([]string, 1)
		r[0] = st.Name
		res[*i] = r
		*i++
	}
}
