package errlog

import (
	"encoding/json"
	"sync"
)

type skippedEntity = map[string]SkipReason

type SkipLog struct {
	policies map[string]skippedEntity
	lock     sync.Mutex
}

func NewSkipLog() *SkipLog {
	return &SkipLog{
		policies: make(map[string]skippedEntity),
	}
}

func (s *SkipLog) Add(policyName string, entityName string, skipReason SkipReason) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.policies[policyName]; !ok {
		s.policies[policyName] = make(skippedEntity)
	}
	s.policies[policyName][entityName] = skipReason
}

func (p *SkipLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.policies)
}

func (p *SkipLog) Empty() bool {
	return len(p.policies) == 0
}
