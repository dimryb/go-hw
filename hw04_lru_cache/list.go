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

type Empty struct{}

type list struct {
	data  map[*ListItem]Empty
	front *ListItem
	back  *ListItem
}

type position int

const (
	front position = iota
	back
)

func NewList() List {
	l := &list{
		data: make(map[*ListItem]Empty),
	}
	return l
}

func (l *list) Len() int {
	return len(l.data)
}

func (l *list) Front() *ListItem {
	if l.Len() == 0 {
		return nil
	}
	return l.front
}

func (l *list) Back() *ListItem {
	if l.Len() == 0 {
		return nil
	}
	return l.back
}

func (l *list) push(v interface{}, position position) *ListItem {
	listItem := &ListItem{Value: v}

	switch {
	case l.Len() == 0:
		l.front = listItem
		l.back = listItem
	case position == front:
		listItem.Next = l.front
		l.front.Prev = listItem
		l.front = listItem
	case position == back:
		listItem.Prev = l.back
		l.back.Next = listItem
		l.back = listItem
	}
	l.data[listItem] = Empty{}

	return listItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	return l.push(v, front)
}

func (l *list) PushBack(v interface{}) *ListItem {
	return l.push(v, back)
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	delete(l.data, i)
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.front == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.back == i {
		l.back = i.Prev
	}

	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}
