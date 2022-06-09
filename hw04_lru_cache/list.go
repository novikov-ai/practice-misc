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
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	l.items[node] = node

	if l.Len() == 1 {
		l.head = node
		l.tail = node
		return node
	}

	node.Next = l.head

	l.head.Prev = node

	l.head = node

	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	l.items[node] = node

	if l.Len() == 1 {
		l.head = node
		l.tail = node
		return node
	}

	node.Prev = l.tail

	l.tail.Next = node

	l.tail = node

	return node
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 0 {
		return
	}

	current := l.items[i]

	delete(l.items, i)

	if l.Len() == 0 && l.head == i {
		l.head, l.tail = nil, nil
		return
	}

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

	if l.tail == i && l.Len() == 2 {
		prevHead := l.head

		prevHead.Next = nil
		prevHead.Prev = i

		l.tail = prevHead
		l.head = i

		i.Next = prevHead
		i.Prev = nil

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
