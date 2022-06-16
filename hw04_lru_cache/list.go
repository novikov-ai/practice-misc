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
	head, tail *ListItem
	length     int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item, isSetUpped := pushToList(l, v)
	if isSetUpped {
		return item
	}

	item.Next = l.head
	l.head.Prev = item

	l.head = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item, isSetUpped := pushToList(l, v)
	if isSetUpped {
		return item
	}

	item.Prev = l.tail
	l.tail.Next = item

	l.tail = item

	return item
}

func pushToList(l *list, v interface{}) (*ListItem, bool) {
	l.length++
	item := newListItem(v)
	isSetUpped := trySetUpHeadAndTail(l, item)
	return item, isSetUpped
}

func newListItem(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	return item
}

func trySetUpHeadAndTail(l *list, item *ListItem) bool {
	if l.Len() != 1 {
		return false
	}

	l.tail = item
	l.head = item
	return true
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 0 {
		return
	}

	l.length--

	if l.Len() == 0 {
		l.head, l.tail = nil, nil
		return
	}

	next := i.Next
	prev := i.Prev

	if next == nil {
		prev.Next = nil
		l.tail = prev
		return
	}

	if prev == nil {
		next.Prev = nil
		l.head = next
		return
	}

	prev.Next = next
	next.Prev = prev
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Len() < 2 || i == l.head {
		return
	}

	prev := i.Prev
	next := i.Next

	prev.Next = next

	if l.tail == i {
		l.tail = prev
	} else {
		next.Prev = prev
	}

	prevHead := l.head
	prevHead.Prev = i

	i.Next = prevHead
	i.Prev = nil

	l.head = i
}

func NewList() List {
	return new(list)
}
