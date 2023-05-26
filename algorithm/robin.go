package algorithm

import "sync"

type RobinSeed struct {
	Index int
	Lock  *sync.Mutex
}

var (
	Robin *RobinSeed
)

func init() {
	Robin = NewRobinSeed()
}

func NewRobinSeed() *RobinSeed {
	return &RobinSeed{0, &sync.Mutex{}}
}

func (r *RobinSeed) NextN(max int) int {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.Index += 1
	if r.Index >= max {
		r.Index = 0
	}
	return r.Index
}
