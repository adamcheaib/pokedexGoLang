package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	interval := 5 * time.Second
	cases := []struct {
		key      string
		expected []byte
	}{
		{
			key:      "https://jsonplaceholder.typicode.com/todos/1",
			expected: []byte("1"),
		},
		{
			key:      "https://jsonplaceholder.typicode.com/todos/3",
			expected: []byte("3"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Running test number %v", i+1), func(t *testing.T) {
			fmt.Println("Running here")
			cache := NewCache(interval)
			cache.Add(c.key, c.expected)
			data, status := cache.Get(c.key)
			if !status {
				t.Fatalf("Expected to find key %v", c.key)
			}
			if string(data) != string(c.expected) {
				t.Fatalf("Expected %v, got %v", string(c.expected), string(data))
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	baseTime := 5 * time.Second
	waitTime := baseTime + 1*time.Second

	cache := NewCache(waitTime)
	cache.Add("tester", []byte("test"))

	if _, status := cache.Get("tester"); !status {
		t.Fatalf("Expected to find key tester")
	}

	time.Sleep(waitTime)

	if _, status := cache.Get("tester"); status {
		t.Fatalf("Expected to not find key tester")
	}

}
