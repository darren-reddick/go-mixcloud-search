package mixcloud

import (
	"fmt"
	"sync"
)

type storeError struct {
	Term string
	msg  string
}

func (i *storeError) Error() string {
	return fmt.Sprintf("%s: %s", i.Term, i.msg)
}

type Store struct {
	Data  map[string]Mix
	limit int
	*sync.RWMutex
}

func NewStore(limit int) Store {
	m := make(map[string]Mix)
	return Store{
		m,
		limit,
		&sync.RWMutex{},
	}
}

func (s *Store) Put(m Mix) error {
	s.Lock()
	defer s.Unlock()
	if len(s.Data) >= s.limit && s.limit != 0 {
		return &storeError{"StoreFull", "Store is full"}
	}
	s.Data[m.Key] = m

	return nil

}
