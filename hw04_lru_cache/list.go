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
	Count      int
	Head, Tail *ListItem
}

func (l list) Len() int {
	return l.Count
}

func (l list) Front() *ListItem {
	return l.Head
}

func (l list) Back() *ListItem {
	return l.Tail
}

func (l list) PushFront(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.Items == nil {
		l.Items = make(map[*ListItem]*ListItem)
	}

	l.Items[node] = node

	l.Count++

	if l.Count == 1 {
		l.Head = node
		l.Tail = node
		return node
	}

	node.Next = l.Head

	l.Head.Prev = node

	l.Head = node

	return node
}

func (l list) PushBack(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.Items == nil {
		l.Items = make(map[*ListItem]*ListItem)
	}

	l.Items[node] = node

	l.Count++

	if l.Count == 1 {
		l.Head = node
		l.Tail = node
		return node
	}

	node.Prev = l.Tail

	l.Tail.Next = node

	l.Tail = node

	return node
}

// todo: O(1). It's O(n) now
func (l list) Remove(i *ListItem) {
	if l.Count == 0 {
		return
	}

	if l.Count == 1 && l.Head == i {
		l.Head, l.Tail = nil, nil
		l.Count = 0
		return
	}

	current := l.Head
	for {
		if current == i {
			next := current.Next
			next.Prev = nil

			current.Next = nil

			l.Head = next
			l.Count--
			return
		}

		current = current.Next
		if current == nil {
			return
		}
	}
}

// todo: O(1). It's O(n) now
func (l list) MoveToFront(i *ListItem) {
	if l.Count < 2 || i == l.Head {
		return
	}

	if l.Tail == i && l.Count == 2 {
		prevHead := l.Head

		prevHead.Next = nil
		prevHead.Prev = i

		l.Tail = prevHead
		l.Head = i

		i.Next = prevHead
		i.Prev = nil

		return
	}

	current := l.Head
	for {
		if current == i {
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

			return
		}

		current = current.Next
		if current == nil {
			return
		}
	}

}

func NewList() List {
	return new(list)
}
