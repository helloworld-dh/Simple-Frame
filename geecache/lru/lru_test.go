package lru

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	lru := New(3, nil)
	lru.Add("1", 1)
	if v, ok := lru.Get("1"); !ok {
		if v != 1 {
			t.Fatalf("Get error number %d \n", v)
		}
		t.Fatal("Get is error")
	}
	if _, ok := lru.Get("2"); ok {
		t.Fatal("Get number from nothing")
	}
}

func TestRemove(t *testing.T) {
	lru := New(3, nil)
	lru.Remove()
	if lru.ll.Len() != 0 {
		t.Fatalf("remove ")
	}
	lru.Remove()
}

func TestOnEvicted(t *testing.T) {
	keys := make([]Key, 0)
	callback := func(key Key, val interface{}) {
		keys = append(keys, key)
	}
	lru := New(10, callback)
	lru.OnEvicted = callback
	for i := 0; i < 11; i++ {
		lru.Add(Key(fmt.Sprintf("myKey%d", i)), i)
	}
	if len(keys) != 1 {
		t.Fatalf("wrong callback")
	}
}
