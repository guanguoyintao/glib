package uset

import (
	"context"
	"github.com/emirpasic/gods/sets/treeset"
)

type TreeSet struct {
	sample interface{}
	s      *treeset.Set
}

func NewTreeSet(ctx context.Context, element interface{}) Set {
	var set *treeset.Set
	switch element.(type) {
	case int, uint, int8, uint8, int32, uint32, int64, uint64:
		set = treeset.NewWithIntComparator()
		set.Add(element)
	case string:
		set = treeset.NewWithIntComparator()
		set.Add(element.(string))
	}

	return &TreeSet{
		sample: element,
		s:      set,
	}
}

func (t *TreeSet) Empty() bool {
	return t.s.Empty()
}

func (t *TreeSet) Size() int {
	return t.s.Size()
}

func (t *TreeSet) Clear() {
	t.s.Clear()
}

func (t *TreeSet) Values() interface{} {
	values := t.s.Values()
	if values == nil || len(values) == 0 {
		return nil
	}
	switch t.sample.(type) {
	case uint64:
		res := make([]uint64, 0, len(values))
		for _, v := range values {
			res = append(res, v.(uint64))
		}
		return res
	case string:
		res := make([]string, 0, len(values))
		for _, v := range values {
			res = append(res, v.(string))
		}
		return res
	}

	return values
}

func (t *TreeSet) String() string {
	return t.s.String()
}

func (t *TreeSet) Add(elements ...interface{}) {
	t.s.Add(elements)
}

func (t *TreeSet) Remove(elements ...interface{}) {
	t.s.Remove(elements)
}

func (t *TreeSet) Contains(elements ...interface{}) bool {
	return t.s.Contains(elements)
}

func (t *TreeSet) Intersection(another Set) Set {
	newSet := NewTreeSet(context.TODO(), t.sample)
	newSet.Remove(t.sample)
	t.s.Each(func(index int, value interface{}) {
		ok := another.Contains(value)
		if ok {
			newSet.Add(value)
		}
	})

	return newSet
}

func (t *TreeSet) Union(another Set) Set {
	newSet := NewTreeSet(context.TODO(), t.sample)
	newSet.Add(another.Values())
	newSet.Add(t.Values())

	return newSet
}

func (t *TreeSet) Difference(another Set) Set {
	newSet := NewTreeSet(context.TODO(), t.sample)
	if another.Contains(t.sample) {
		newSet.Remove(t.sample)
	}
	t.s.Each(func(index int, value interface{}) {
		ok := another.Contains(value)
		if !ok {
			newSet.Add(value)
		}
	})

	return newSet
}
