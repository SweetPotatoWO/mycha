package module

import (
	"math"
	"sync"
)

//序列号生产接口
type SNGenertor interface {
	Start() uint32  //起步的参数
	Max() uint32  //最大的参数
	Next() uint32 //获取到下一个
	CycleCount() uint32
	Get() uint32
}

type mySNGenertor struct {
	start uint32
	max uint32
	next uint32
	cycleCount uint32
	lock sync.RWMutex
}


func NewSNGenertor(start uint32,max uint32) SNGenertor {
	if max == 0 {
		max = math.MaxUint64
	}
	return &mySNGenertor{
		start :start,
		max:max,
		next:start,
	}
}


func (gen *mySNGenertor) Start() uint32 {
	return gen.start
}

func (gen *mySNGenertor) Max() uint32 {
	return gen.max

}

func (gen *mySNGenertor) Next() uint32 {
	gen.lock.RLock()
	defer gen.lock.RUnlock()
	return gen.next
}

func (gen *mySNGenertor) CycleCount() uint32 {
	gen.lock.RLock()
	defer gen.lock.RUnlock()
	return gen.cycleCount
}
func (gen *mySNGenertor) Get() uint32 {
	gen.lock.Lock()
	defer gen.lock.Unlock()
	id := gen.next
	if id == gen.max {
		gen.next = gen.start
		gen.cycleCount ++
	} else {
			gen.next++
	}
	return id
}




