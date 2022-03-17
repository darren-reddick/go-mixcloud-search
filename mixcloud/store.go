package mixcloud

import (
	"sync"
)

type Store struct {
	Data map[string]Mix
	*sync.RWMutex
}

func NewStore() Store {
	m := make(map[string]Mix)
	return Store{
		m,
		&sync.RWMutex{},
	}
}

func (s *Store) Put(m Mix) {
	s.Lock()
	defer s.Unlock()
	s.Data[m.Key] = m

}
