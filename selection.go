package qtablam

import (
	"slices"
)

type Selection struct {
	Index int
	Elems []int
}

func newSelection() Selection {
	return Selection{
		Index: 0,
		Elems: make([]int, 0, 16),
	}
}

func (s *Selection) IsEmpty() bool {
	return len(s.Elems) == 0
}

func (s *Selection) Contains(e int) bool {
	return slices.Contains(s.Elems, e)
}

func (s *Selection) Pull(e int) bool {
	if s.Contains(e) {
		return false
	}
	s.Elems = append([]int{e}, s.Elems...)
	return true
}

func (s *Selection) Push(e int) bool {
	if s.Contains(e) {
		return false
	}
	s.Elems = append(s.Elems, e)
	return true
}

func (s *Selection) Next() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		var elem int
		if s.Index == len(s.Elems)-1 {
			elem = s.Elems[s.Index]
			s.Index = 0
		} else if s.Index < len(s.Elems)-1 {
			elem = s.Elems[s.Index]
			s.Index++
		} else {
			s.Index = 0
		}
		return elem, true
	}
}

func (s *Selection) Remove(e int) bool {
	for i, elem := range s.Elems {
		if elem == e {
			s.Elems = append(s.Elems[:i], s.Elems[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Selection) Clear() {
	s.Elems = []int{}
}
