package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("pop last because of size", func(t *testing.T) {
		c := NewCache(3)
		items := []cacheItem{
			{key: "first", value: 100},
			{key: "second", value: 200},
			{key: "third", value: 300},
			{key: "fourth", value: 400},
		}

		for _, item := range items {
			wasInCache := c.Set(item.key, item.value)
			require.False(t, wasInCache)
		}

		val, ok := c.Get("second")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("third")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("fourth")
		require.True(t, ok)
		require.Equal(t, 400, val)

		val, ok = c.Get("first")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("pop least recently used", func(t *testing.T) {
		c := NewCache(5)
		items := []cacheItem{
			{key: "first", value: 100},
			{key: "second", value: 200},
			{key: "third", value: 300},
			{key: "fourth", value: 400},
			{key: "fifth", value: 500},
		}

		for _, item := range items {
			wasInCache := c.Set(item.key, item.value)
			require.False(t, wasInCache)
		}
		// queque: "fifth", "fourth", "third", "second", "first"

		val, ok := c.Get("second")
		require.True(t, ok)
		require.Equal(t, 200, val)
		// queque:  "second", "fifth", "fourth", "third", "first"

		wasInCache := c.Set("first", 10_000)
		require.True(t, wasInCache)
		// queque: "first", "second", "fifth", "fourth", "third"

		val, ok = c.Get("fifth")
		require.True(t, ok)
		require.Equal(t, 500, val)
		// queque: "fifth", "first", "second", "fourth", "third"

		val, ok = c.Get("first")
		require.True(t, ok)
		require.Equal(t, 10_000, val)
		// queque: "first", "fifth", "second", "fourth", "third"

		wasInCache = c.Set("sixth", 600)
		require.False(t, wasInCache)
		// queque: "sixth", "first", "fifth", "second", "fourth"

		// "third" has to be popped
		val, ok = c.Get("third")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(5)
		items := []cacheItem{
			{key: "first", value: 100},
			{key: "second", value: 200},
			{key: "third", value: 300},
			{key: "fourth", value: 400},
			{key: "fifth", value: 500},
		}

		for _, item := range items {
			wasInCache := c.Set(item.key, item.value)
			require.False(t, wasInCache)
		}
		// queque: "fifth", "fourth", "third", "second", "first"

		c.Clear()

		// cache empty
		for _, item := range items {
			val, ok := c.Get(item.key)
			require.False(t, ok)
			require.Nil(t, val)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
