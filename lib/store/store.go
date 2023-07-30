package store

import (
	"errors"
	"fmt"
)

type Store[T comparable, U any] struct {
	storeMap map[T]U
}

func NewStore[T comparable, U any]() *Store[T, U] {
	return &Store[T, U]{
		storeMap: make(map[T]U),
	}
}

func (u *Store[T, U]) Set(key T, value U) {
	u.storeMap[key] = value
}

func (u *Store[T, U]) Get(key T) (U, error) {
	value, ok := u.storeMap[key]
	if !ok {
		var zero U
		return zero, errors.New(fmt.Sprintf("Key %v does not exist in store", key))
	}
	return value, nil
}

func (u *Store[T, U]) Delete(key T) {
	delete(u.storeMap, key)
}
