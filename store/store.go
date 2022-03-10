package store

import (
	"sync"

	"github.com/darren-reddick/go-mixcloud-search/schema"
)

type Store struct {
	Data map[string]schema.Mix
	*sync.RWMutex
}

func NewStore() Store {
	m := make(map[string]schema.Mix)
	return Store{
		m,
		&sync.RWMutex{},
	}
}

func (s *Store) Put(m schema.Mix) error {
	s.Lock()
	defer s.Unlock()
	s.Data[m.Key] = m
	return nil

}
