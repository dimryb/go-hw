package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	// Place your code here.
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int { return 0 }

func (l *list) Front() *ListItem {
	return &ListItem{
		Value: nil,
		Next:  nil,
		Prev:  nil,
	}
}

func (l *list) Back() *ListItem {
	return &ListItem{
		Value: nil,
		Next:  nil,
		Prev:  nil,
	}
}

func (l *list) PushFront(v interface{}) *ListItem {
	return &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
}

func (l *list) PushBack(v interface{}) *ListItem {
	return &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}
}
