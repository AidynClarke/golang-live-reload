package utils

import (
	"sync"
)

type MUArray[T any] struct {
	array []T
	lock  sync.RWMutex
}

func NewMUArray[T any]() *MUArray[T] {
	return &MUArray[T]{
		array: []T{},
		lock:  sync.RWMutex{},
	}
}

func (a *MUArray[T]) Get() []T {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.array
}

func (a *MUArray[T]) Set(value []T) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.array = value
}

func (a *MUArray[T]) Append(value T) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.array = append(a.array, value)
}

func (a *MUArray[T]) Transaction() ([]T, func()) {
	a.lock.Lock()
	return a.array, func() {
		a.lock.Unlock()
	}
}