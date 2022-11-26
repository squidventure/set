package set

import (
	"sync"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestNew(t *testing.T) {
	s := NewSet[int]()
	if _, ok := s.(*StandardSet[int]); !ok {
		t.Fatalf("error: expected standard set but got %T\n", s)
	}
	s = NewSet[int](true)
	if _, ok := s.(*ThreadsafeSet[int]); !ok {
		t.Fatalf("error: expected threadsafe set but got %T\n", s)
	}
}

func SlicesMatch[T constraints.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStandardSet(t *testing.T) {
	set := NewSet[int]()
	for i := 0; i < 5; i++ {
		set.Add(i)
	}
	if !SlicesMatch([]int{0, 1, 2, 3, 4}, set.Slice()) {
		t.Fatalf("error: expected [0,1,2,3,4], got %v\n", set.Slice())
	}
	for i := 0; i < 5; i++ {
		set.Add(i)
	}
	if !SlicesMatch([]int{0, 1, 2, 3, 4}, set.Slice()) {
		t.Fatalf("error: expected [0,1,2,3,4] after second add, got %v\n", set.Slice())
	}
	for i := 0; i < 5; i++ {
		set.Drop(i)
	}
	if !SlicesMatch([]int{}, set.Slice()) {
		t.Fatalf("error: expected [], got %v\n", set.Slice())
	}
}

func TestThreadsafeSet(t *testing.T) {
	set := NewSet[int](true)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set.Add(i)
		}(i)
	}
	wg.Wait()
	if !SlicesMatch([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, set.Slice()) {
		t.Fatalf("error: expected [0,1,2,3,4,5,6,7,8,9], got %v\n", set.Slice())
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set.Add(i)
		}(i)
	}
	wg.Wait()
	if !SlicesMatch([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, set.Slice()) {
		t.Fatalf("error: expected [0,1,2,3,4,5,6,7,8,9] after second Add, got %v\n", set.Slice())
	}
}

var (
	standard_set   = NewStandardSet[int]()
	threadsafe_set = NewThreadsafeSet[int]()
)

func BenchmarkStandardAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		standard_set.Add(i)
	}
}

func BenchmarkStandardDrop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		standard_set.Drop(i)
	}
}

func BenchmarkThreadsafeAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		threadsafe_set.Add(i)
	}
}

func BenchmarkThreadsafeDrop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		threadsafe_set.Drop(i)
	}
}
