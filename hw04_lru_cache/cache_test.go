package hw04lrucache

import (
	"fmt"
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

	t.Run("purge logic", func(t *testing.T) {
		cacheSize := 10
		cache := cacheInit(cacheSize)

		cache.Clear()

		for i := 0; i < cacheSize; i++ {
			key := getKeyBasedOnIndex(i)

			_, ok := cache.Get(key)
			require.False(t, ok)
		}

		newItemKey := Key("new item")
		cache.Set(newItemKey, 1)

		_, ok := cache.Get("new item")
		require.True(t, ok)

		cache.Clear()

		_, ok = cache.Get("new item")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
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

func TestCacheKnockOutOverflow(t *testing.T) {
	t.Run("knocked out if overflow", func(t *testing.T) {
		cache := cacheInit(10)

		knockoutKey := Key("bouncer")
		cache.Set(knockoutKey, 42)

		_, isFirstAtCache := cache.Get(getKeyBasedOnIndex(0))
		require.False(t, isFirstAtCache)

		_, isBouncerAtCache := cache.Get(knockoutKey)
		require.True(t, isBouncerAtCache)
	})
}

func TestCacheKnockOutUnused(t *testing.T) {
	t.Run("knocked out unused items", func(t *testing.T) {
		cacheSize := 10
		cache := cacheInit(cacheSize)

		// simulating frequent usage of even items
		for i := 0; i < cacheSize; i++ {
			if i%2 != 0 {
				continue
			}
			key := getKeyBasedOnIndex(i)
			_, ok := cache.Get(key)
			require.True(t, ok)
		}

		// unused: item #1, item #3, item #5, item #7, item #9

		for i := 10; i < 15; i++ {
			key := getKeyBasedOnIndex(i)
			cache.Set(key, i*i)
		}

		for i := 0; i < 10; i++ {
			key := getKeyBasedOnIndex(i)
			if i%2 == 0 {
				_, ok := cache.Get(key)
				require.True(t, ok)
				continue
			}

			_, ok := cache.Get(key)
			require.False(t, ok)
		}
	})
}

func getKeyBasedOnIndex(i int) Key {
	const itemKeyPrefix = "item #"
	return Key(fmt.Sprintf("%s%v", itemKeyPrefix, i))
}

func cacheInit(size int) Cache {
	c := NewCache(size)
	for i := 0; i < size; i++ {
		key := getKeyBasedOnIndex(i)
		c.Set(key, i*i)
	}
	return c
}
