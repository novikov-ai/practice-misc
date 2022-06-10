package hw04lrucache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func listInit() List {
	l := NewList()
	for i := 0; i < 15; i++ {
		if i%2 == 0 {
			l.PushBack(i)
		} else {
			l.PushFront(i)
		}
	}
	return l // 13 11 9 7 5 3 1 0 2 4 6 8 10 12 14
}

func getItemByValue(l List, v int) *ListItem {
	item := l.Front()
	for {
		if item.Value.(int) == v {
			return item
		}
		item = item.Next
		if item == nil {
			break
		}
	}
	return nil
}

func TestListLen(t *testing.T) {
	t.Run("count items", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())

		l = listInit()

		require.Equal(t, 15, l.Len())
	})
}

func TestListGetHead(t *testing.T) {
	t.Run("get head", func(t *testing.T) {
		l := NewList()

		require.Nil(t, l.Front())

		l = listInit() // 13 11 9 7 5 3 1 0 2 4 6 8 10 12 14

		require.Equal(t, 13, l.Front().Value.(int))
	})
}

func TestListGetTail(t *testing.T) {
	t.Run("get tail", func(t *testing.T) {
		l := NewList() // 0 items

		require.Nil(t, l.Back())

		l = listInit() // 15 items

		require.Equal(t, 14, l.Back().Value.(int))
	})
}

func TestListPushFront(t *testing.T) {
	t.Run("push front", func(t *testing.T) {
		l := NewList()
		require.Equal(t, 0, l.Len())

		headValue := 50
		l.PushFront(headValue)

		require.Equal(t, 1, l.Len())
		require.Equal(t, headValue, l.Front().Value.(int))

		headValueText := "something"
		l.PushFront(headValueText)

		require.Equal(t, 2, l.Len())
		require.Equal(t, headValueText, l.Front().Value.(string))
	})
}

func TestListPushBack(t *testing.T) {
	t.Run("push back", func(t *testing.T) {
		l := NewList()
		require.Equal(t, 0, l.Len())

		tailValue := 50
		l.PushFront(tailValue)

		require.Equal(t, 1, l.Len())
		require.Equal(t, tailValue, l.Back().Value.(int))

		tailValueText := "something"
		l.PushBack(tailValueText)

		require.Equal(t, 2, l.Len())
		require.Equal(t, tailValueText, l.Back().Value.(string))
	})
}

func TestListRemove(t *testing.T) {
	t.Run("remove item", func(t *testing.T) {
		l := listInit() // 13 11 9 7 5 3 1 0 2 4 6 8 10 12 14
		require.Equal(t, 15, l.Len())

		l.Remove(l.Front())
		require.Equal(t, 11, l.Front().Value.(int))
		require.Equal(t, 14, l.Len())

		l.Remove(l.Back())
		require.Equal(t, 12, l.Back().Value.(int))
		require.Equal(t, 13, l.Len())

		middleItem := getItemByValue(l, 1)
		if middleItem == nil {
			fmt.Println("Item wasn't found. Tests interrupted.")
			return
		}

		l.Remove(middleItem)
		require.Equal(t, 11, l.Front().Value.(int))
		require.Equal(t, 12, l.Back().Value.(int))
		require.Equal(t, 12, l.Len())

		for i := 0; i < 12; i++ {
			l.Remove(l.Front())
		}

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}

func TestListMoveToFront(t *testing.T) {
	t.Run("move to front item", func(t *testing.T) {
		l := listInit() // 13 11 9 7 5 3 1 0 2 4 6 8 10 12 14
		require.Equal(t, 15, l.Len())

		require.Equal(t, 13, l.Front().Value.(int))
		require.Equal(t, 14, l.Back().Value.(int))

		l.MoveToFront(l.Back())
		require.Equal(t, 14, l.Front().Value.(int))
		require.Equal(t, 12, l.Back().Value.(int))

		middleItemValue := 1
		middleItem := getItemByValue(l, middleItemValue)
		if middleItem == nil {
			fmt.Println("Item wasn't found. Tests interrupted.")
			return
		}

		l.MoveToFront(middleItem)
		require.Equal(t, middleItemValue, l.Front().Value.(int))
		require.Equal(t, 12, l.Back().Value.(int))
		require.Equal(t, 15, l.Len())
	})
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
