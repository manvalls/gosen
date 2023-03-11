package selectorcache

import (
	"sync"

	"github.com/andybalholm/cascadia"
)

type SelectorCache struct {
	cache map[string]cascadia.Sel
	mux   sync.RWMutex
}

func New() *SelectorCache {
	return &SelectorCache{
		cache: make(map[string]cascadia.Sel),
	}
}

func (s *SelectorCache) Get(selector string, dynamic bool) (cascadia.Sel, error) {
	if dynamic {
		return cascadia.Parse(selector)
	}

	s.mux.RLock()

	if cached, ok := s.cache[selector]; ok {
		s.mux.RUnlock()
		return cached, nil
	}

	s.mux.RUnlock()

	sel, err := cascadia.Parse(selector)
	if err != nil {
		return nil, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	s.cache[selector] = sel
	return sel, nil
}
