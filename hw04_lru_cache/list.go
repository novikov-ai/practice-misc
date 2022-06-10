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
	items      map[*ListItem]*ListItem
	head, tail *ListItem
}

func (l *list) Len() int {
	return len(l.items)
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
	item := appendToList(l, v)
	isSetUpped := trySetUpHeadAndTail(l, item)
	return item, isSetUpped
}

func appendToList(l *list, v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	l.items[item] = item

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

	defer delete(l.items, i)

	if l.Len() == 1 && l.head == i {
		l.head, l.tail = nil, nil
		return
	}

	current := l.items[i]

	next := current.Next
	prev := current.Prev

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

	current := l.items[i]

	prev := current.Prev
	next := current.Next

	prev.Next = next

	if l.tail == current {
		l.tail = prev
	} else {
		next.Prev = prev
	}

	prevHead := l.head
	prevHead.Prev = current

	current.Next = prevHead
	current.Prev = nil

	l.head = current
}

func NewList() List {
	return &list{
		items: make(map[*ListItem]*ListItem),
	}
}
