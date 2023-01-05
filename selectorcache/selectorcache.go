package selectorcache

import (
	"fmt"
	"sync"

	"github.com/andybalholm/cascadia"
)

type SelectorCache struct {
	cache map[string]cascadia.Sel
	mux   sync.Mutex
}

func New() *SelectorCache {
	return &SelectorCache{
		cache: make(map[string]cascadia.Sel),
	}
}

func (s *SelectorCache) Get(selector string, args []interface{}) (cascadia.Sel, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if len(args) > 0 {
		return cascadia.Parse(fmt.Sprintf(selector, args...))
	}

	if cached, ok := s.cache[selector]; ok {
		return cached, nil
	}

	sel, err := cascadia.Parse(selector)
	if err != nil {
		return nil, err
	}

	s.cache[selector] = sel
	return sel, nil
}
