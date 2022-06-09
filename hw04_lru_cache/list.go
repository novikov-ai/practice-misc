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
	Items      map[*ListItem]*ListItem
	Head, Tail *ListItem
}

func (l *list) Len() int {
	return len(l.Items)
}

func (l *list) Front() *ListItem {
	return l.Head
}

func (l *list) Back() *ListItem {
	return l.Tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.Items == nil {
		l.Items = make(map[*ListItem]*ListItem)
	}

	l.Items[node] = node

	if l.Len() == 1 {
		l.Head = node
		l.Tail = node
		return node
	}

	node.Next = l.Head

	l.Head.Prev = node

	l.Head = node

	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.Items == nil {
		l.Items = make(map[*ListItem]*ListItem)
	}

	l.Items[node] = node

	if l.Len() == 1 {
		l.Head = node
		l.Tail = node
		return node
	}

	node.Prev = l.Tail

	l.Tail.Next = node

	l.Tail = node

	return node
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 0 {
		return
	}

	current := l.Items[i]

	delete(l.Items, i)

	if l.Len() == 0 && l.Head == i {
		l.Head, l.Tail = nil, nil
		return
	}

	next := current.Next
	prev := current.Prev

	if next == nil {
		prev.Next = nil
		l.Tail = prev
		return
	}

	if prev == nil {
		next.Prev = nil
		l.Head = next
		return
	}

	prev.Next = next
	next.Prev = prev
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Len() < 2 || i == l.Head {
		return
	}

	if l.Tail == i && l.Len() == 2 {
		prevHead := l.Head

		prevHead.Next = nil
		prevHead.Prev = i

		l.Tail = prevHead
		l.Head = i

		i.Next = prevHead
		i.Prev = nil

		return
	}

	current := l.Items[i]

	prev := current.Prev
	next := current.Next

	prev.Next = next

	if l.Tail == current {
		l.Tail = prev
	} else {
		next.Prev = prev
	}

	prevHead := l.Head
	prevHead.Prev = current

	current.Next = prevHead
	current.Prev = nil

	l.Head = current
}

func NewList() List {

	return new(list)
}
