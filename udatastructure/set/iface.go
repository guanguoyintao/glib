package uset

type Set interface {
	Empty() bool
	Size() int
	Clear()
	Values() interface{}
	String() string
	Add(elements ...interface{})
	Remove(elements ...interface{})
	Contains(elements ...interface{}) bool
	Intersection(another Set) Set
	Union(another Set) Set
	Difference(another Set) Set
}
