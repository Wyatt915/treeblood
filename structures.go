package treeblood

import "fmt"

type stack[T any] struct {
	data []T
	top  int
}

func newStack[T any]() *stack[T] {
	return &stack[T]{
		data: make([]T, 0),
		top:  -1,
	}
}

func (s *stack[T]) Push(val T) {
	s.top++
	if len(s.data) <= s.top { // Check if we need to grow the slice
		newSize := len(s.data) * 2
		if newSize == 0 {
			newSize = 1 // Start with a minimum capacity if the stack is empty
		}
		newData := make([]T, newSize)
		copy(newData, s.data) // Copy old elements to new slice
		s.data = newData
	}
	s.data[s.top] = val
}

func (s *stack[T]) Peek() (val T) {
	val = s.data[s.top]
	return
}

func (s *stack[T]) Pop() (val T) {
	val = s.data[s.top]
	s.top--
	return
}

func (s *stack[T]) empty() bool {
	return s.top < 0
}

////////////////////////////////////////////////////////////////////////////////

var (
	ErrEmptyQueue = fmt.Errorf("popping empty queue")
)

// nice implementation from https://stackoverflow.com/a/50418813
type queue[T any] struct {
	data []T
	head int
	tail int
	sz   int
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{
		data: make([]T, 256),
		head: 0,
		tail: 0,
		sz:   0,
	}
}

func (q *queue[T]) Empty() bool {
	return q.sz == 0
}

func (q *queue[T]) next(i int) int {
	return (i + 1) & (len(q.data) - 1)
}

func (q *queue[T]) prev(i int) int {
	return (i - 1) & (len(q.data) - 1)
}

func (q *queue[T]) growIfFull() {
	if q.sz < len(q.data) {
		return
	}
	newBuf := make([]T, q.sz<<1)
	if q.tail > q.head {
		copy(newBuf, q.data[q.head:q.tail])
	} else {
		n := copy(newBuf, q.data[q.head:])
		copy(newBuf[n:], q.data[:q.tail])
	}
	q.head = 0
	q.tail = q.sz
	q.data = newBuf
}

func (q *queue[T]) PushFront(item T) {
	q.growIfFull()
	q.head = q.prev(q.head)
	q.data[q.head] = item
	q.sz++
}

func (q *queue[T]) PopFront() (T, error) {
	var result T
	if q.sz < 1 {
		return result, ErrEmptyQueue
	}
	result = q.data[q.head]
	q.head = q.next(q.head)
	q.sz--
	return result, nil
}

func (q *queue[T]) PeekFront() (T, error) {
	if q.sz < 1 {
		return *new(T), ErrEmptyQueue
	}
	return q.data[q.head], nil
}

// Return the first element that does not satisfy the condition
func (q *queue[T]) PopFrontWhile(condition func(T) bool) (T, error) {
	result, err := q.PopFront()
	for condition(result) && err == nil {
		result, err = q.PopFront()
	}
	return result, err
}

func (q *queue[T]) PushBack(item T) {
	q.growIfFull()
	q.data[q.tail] = item
	q.tail = q.next(q.tail)
	q.sz++
}

func (q *queue[T]) PopBack() (T, error) {
	var result T
	if q.sz < 1 {
		return result, ErrEmptyQueue
	}
	result = q.data[q.tail]
	q.tail = q.prev(q.tail)
	q.sz--
	return result, nil
}
