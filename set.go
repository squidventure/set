package set

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

type Set[T constraints.Ordered] interface {
	Add(t ...T)
	Drop(t ...T)
	Slice() []T
	Contains(t T) bool
	Reset() []T
}

func NewSet[T constraints.Ordered](threadsafe ...bool) Set[T] {
	if len(threadsafe) == 1 && threadsafe[0] {
		return NewThreadsafeSet[T]()
	} else {
		return NewStandardSet[T]()
	}
}

type StandardSet[T constraints.Ordered] struct {
	contents map[T]any
}

func NewStandardSet[T constraints.Ordered]() *StandardSet[T] {
	return &StandardSet[T]{contents: make(map[T]any)}
}

func (s *StandardSet[T]) Add(t ...T) {
	for _, tmp := range t {
		s.contents[tmp] = struct{}{}
	}
}

func (s *StandardSet[T]) Drop(t ...T) {
	for _, tmp := range t {
		delete(s.contents, tmp)
	}
}

func (s *StandardSet[T]) Slice() []T {
	results := make([]T, 0, len(s.contents))
	for t := range s.contents {
		results = append(results, t)
	}
	sort.Slice(results, func(i, j int) bool { return results[i] < results[j] })
	return results
}

func (s *StandardSet[T]) Reset() []T {
	contents := s.Slice()
	s.contents = make(map[T]any)
	return contents
}

func (s *StandardSet[T]) Contains(t T) bool {
	_, found := s.contents[t]
	return found
}

type ThreadsafeSet[T constraints.Ordered] struct {
	mutex sync.RWMutex
	set   *StandardSet[T]
}

func NewThreadsafeSet[T constraints.Ordered]() *ThreadsafeSet[T] {
	return &ThreadsafeSet[T]{
		set: NewStandardSet[T](),
	}
}

func (s *ThreadsafeSet[T]) Add(t ...T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set.Add(t...)
}

func (s *ThreadsafeSet[T]) Drop(t ...T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set.Drop(t...)
}

func (s *ThreadsafeSet[T]) Slice() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.set.Slice()
}

func (s *ThreadsafeSet[T]) Contains(t T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.set.Contains(t)
}

func (s *ThreadsafeSet[T]) Reset() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set.Reset()
}
