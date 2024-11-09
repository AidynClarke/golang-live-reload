package utils

import "sync"

type MUVar[T any] struct {
	value T
	lock sync.Mutex
}

func NewMUFlag[T any](initialValue T) *MUVar[T] {
	return &MUVar[T]{
		value: initialValue,
		lock: sync.Mutex{},
	}
}

func (f *MUVar[T]) Set(value T) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.value = value
}

func (f *MUVar[T]) Get() T {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.value
}

func (f *MUVar[T]) GetThenSet(value T) T {
	f.lock.Lock()
	defer f.lock.Unlock()

	oldValue := f.value

	f.value = value
	return oldValue
}

func (f *MUVar[T]) UnsafeSet(value T) {
	f.value = value
}

func (f *MUVar[T]) UnsafeGet() T {
	return f.value
}

func (f *MUVar[T]) Transaction() func() {
	f.lock.Lock()
	return func() {
		f.lock.Unlock()
	}
}