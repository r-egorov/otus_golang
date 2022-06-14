package hw04lrucache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

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

		expectedList := []int{70, 80, 60, 40, 10, 30, 50}
		assertList(t, l, expectedList)
	})

	t.Run("PushFront empty list", func(t *testing.T) {
		l := NewList()

		l.PushFront("first") // ["first"]
		require.Equal(t, 1, l.Len())

		front := l.Front()
		back := l.Back()

		require.Same(t, front, back)
		require.Equal(t, "first", front.Value)

		expectedList := []string{"first"}
		assertList(t, l, expectedList)
	})

	t.Run("PushBack empty list", func(t *testing.T) {
		l := NewList()

		l.PushBack("first") // ["first"]
		require.Equal(t, 1, l.Len())

		front := l.Front()
		back := l.Back()

		require.Same(t, front, back)
		require.Equal(t, "first", front.Value)

		expectedList := []string{"first"}
		assertList(t, l, expectedList)
	})

	t.Run("Remove, list of one element", func(t *testing.T) {
		l := NewList()
		l.PushBack("first")

		l.Remove(l.Front())
		require.Equal(t, 0, l.Len())
		assertList(t, l, []string{})
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("Remove", func(t *testing.T) {
		l := NewList()

		l.PushBack("first")
		l.PushBack("second")
		l.PushBack("third")
		l.PushBack("fourth")
		l.PushBack("fifth") // ["first", "second", "third", "fourth", "fifth"]

		l.Remove(l.Front())
		expectedList := []string{"second", "third", "fourth", "fifth"}
		assertList(t, l, expectedList)

		l.Remove(l.Front().Next)
		expectedList = []string{"second", "fourth", "fifth"}
		assertList(t, l, expectedList)

		l.Remove(l.Back().Prev)
		expectedList = []string{"second", "fifth"}
		assertList(t, l, expectedList)

		l.Remove(l.Back())
		expectedList = []string{"second"}
		assertList(t, l, expectedList)
	})

	t.Run("MoveToFront", func(t *testing.T) {
		l := NewList()
		l.PushBack("first")
		l.MoveToFront(l.Front())
		assertList(t, l, []string{"first"})

		l.PushBack("second")
		l.PushBack("third")
		l.PushBack("fourth")
		l.PushBack("fifth")

		expectedList := []string{"first", "second", "third", "fourth", "fifth"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Back())
		expectedList = []string{"fifth", "first", "second", "third", "fourth"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Back())
		expectedList = []string{"fourth", "fifth", "first", "second", "third"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Back())
		expectedList = []string{"third", "fourth", "fifth", "first", "second"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Back())
		expectedList = []string{"second", "third", "fourth", "fifth", "first"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Back())
		expectedList = []string{"first", "second", "third", "fourth", "fifth"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Front())
		expectedList = []string{"first", "second", "third", "fourth", "fifth"}
		assertList(t, l, expectedList)

		l.MoveToFront(l.Front().Next.Next) // "third"
		expectedList = []string{"third", "first", "second", "fourth", "fifth"}
		assertList(t, l, expectedList)
	})
}

// Takes a slice of `expected` values (strings or ints) in the specific order
// and bidirectionally checks if the list has these values in the same order.
func assertList(t *testing.T, li List, expected interface{}) {
	t.Helper()
	// expected is a slice of values (int or string) in the list

	expectedLen := reflect.ValueOf(expected).Len()
	require.Equal(t, li.Len(), expectedLen)

	switch expected.(type) {
	case []string:
		elems := make([]string, 0, li.Len())
		for i := li.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(string))
		}
		require.Equal(t, expected, elems)
		// reverse expected and check from the back using Prev
		reverseSlice(expected)
		revElems := make([]string, 0, li.Len())
		for i := li.Back(); i != nil; i = i.Prev {
			revElems = append(revElems, i.Value.(string))
		}
		require.Equal(t, expected, revElems)
		reverseSlice(expected) // reverse expected back to initial state
	case []int:
		elems := make([]int, 0, li.Len())
		for i := li.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, expected, elems)

		// reverse expected and check from the back using Prev
		reverseSlice(expected)
		revElems := make([]int, 0, li.Len())
		for i := li.Back(); i != nil; i = i.Prev {
			revElems = append(revElems, i.Value.(int))
		}
		require.Equal(t, expected, revElems)
		reverseSlice(expected) // reverse expected back to initial state
	}
}

func reverseSlice(s interface{}) {
	sLen := reflect.ValueOf(s).Len()
	sMid := sLen / 2
	sSwap := reflect.Swapper(s)

	for i := 0; i < sMid; i++ {
		j := sLen - i - 1

		sSwap(i, j)
	}
}
