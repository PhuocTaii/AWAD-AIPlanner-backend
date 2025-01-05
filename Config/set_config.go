package config

type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates and returns a new Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{elements: make(map[T]struct{})}
}

// Add adds an element to the Set
func (s *Set[T]) Add(element T) {
	s.elements[element] = struct{}{}
}

// Remove removes an element from the Set
func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

// Contains checks if the Set contains an element
func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Size returns the number of elements in the Set
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Clear removes all elements from the Set
func (s *Set[T]) Clear() {
	s.elements = make(map[T]struct{})
}

// ToSlice converts the Set to a slice
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.elements))
	for key := range s.elements {
		result = append(result, key)
	}
	return result
}
