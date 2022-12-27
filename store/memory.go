package store

import (
	"errors"
	"sync"
)

type InMemoryUrlStore struct {
	data  []string
	len   int
	mutex sync.Mutex
}

func (s *InMemoryUrlStore) Store(url string) int {
	s.mutex.Lock()
	s.data = append(s.data, url)
	idx := len(s.data)
	s.len = s.len + 1
	s.mutex.Unlock()
	return idx
}

func (s *InMemoryUrlStore) GetById(id int) (string, error) {
	if id > s.len {
		return "", errors.New("record does not exist")
	}
	return s.data[id-1], nil
}

func (s *InMemoryUrlStore) GetSize() int {
	return s.len
}

func NewInMemoryUrlStore() *InMemoryUrlStore {
	return &InMemoryUrlStore{data: make([]string, 0), len: 0}
}
