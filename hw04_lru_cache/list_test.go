package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func toList(l List) []int {
	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	return elems
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("simple push front", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushFront(20) // [20, 10]
		l.PushFront(30) // [30, 20, 10]

		require.Equal(t, 3, l.Len())

		elems := toList(l)
		require.Equal(t, []int{30, 20, 10}, elems)
	})

	t.Run("simple push back", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]
		l.PushBack(20) // [10, 20]
		l.PushBack(30) // [10, 20, 30]

		require.Equal(t, 3, l.Len())

		elems := toList(l)
		require.Equal(t, []int{10, 20, 30}, elems)
	})

	t.Run("remove", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)
		l.PushBack(50)

		require.Equal(t, 5, l.Len())

		middle := l.Front().Next.Next
		l.Remove(middle)

		require.Equal(t, 4, l.Len())

		elems := toList(l)
		require.Equal(t, []int{10, 20, 40, 50}, elems)
	})

	t.Run("move to front", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)

		item := l.Front().Next.Next
		l.MoveToFront(item) // [30, 10, 20, 40]

		elems := toList(l)
		require.Equal(t, []int{30, 10, 20, 40}, elems)
	})

	t.Run("move to front end item", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)

		l.MoveToFront(l.Back()) // [40, 10, 20, 30]

		elems := toList(l)
		assert.Equal(t, []int{40, 10, 20, 30}, elems)
		require.Equal(t, 30, l.Back().Value.(int))
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

		elems := toList(l)
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
