package stub

import (
	"fmt"
	"gopcp.v2/chapter6/webcrawler/errors"
	"mycha/helper/log"
	"mycha/module"
	"sync/atomic"
)

var logger = log.DLogger()

type myModule struct {
	mid             module.MID
	addr            string
	score           uint32
	scoreCalculator module.CalculateScore
	// calledCount 代表调用计数。
	calledCount uint32
	// acceptedCount 代表接受计数。
	acceptedCount uint32
	// completedCount 代表成功完成计数。
	completedCount uint32
	// handlingNumber 代表实时处理数。
	handlingNumber uint32
}

func NewModuleInternal(mid module.MID,
	scoreCalculator module.CalculateScore) (ModuleInternal, error) {

	parts, err := module.SplitMID(mid)
	if err != nil {
		return nil, errors.NewIllegalParameterError(
			fmt.Sprintf("illegal ID %q: %s", mid, err))
	}
	return &myModule{
		mid:             mid,
		addr:            parts[2],
		scoreCalculator: scoreCalculator,
	}, nil
}


func (m *myModule) ID() module.MID {
	return m.mid
}

func (m *myModule) Addr() string {
	return m.addr
}
func (m *myModule) Score() uint32 {
	return atomic.LoadUint32(&m.score)
}

func (m *myModule) SetScore(score uint32) {
	atomic.StoreUint32(&m.score, score)
}

func (m *myModule) ScoreCalculator() module.CalculateScore {
	return m.scoreCalculator
}

func (m *myModule) CalledCount() uint32 {
	return atomic.LoadUint32(&m.calledCount)
}

func (m *myModule) AcceptedCount() uint32 {
	return atomic.LoadUint32(&m.acceptedCount)
}

func (m *myModule) CompletedCount() uint32 {
	count := atomic.LoadUint32(&m.completedCount)
	return count
}

func (m *myModule) HandlingNumber() uint32 {
	return atomic.LoadUint32(&m.handlingNumber)
}

func (m *myModule) Counts() module.Counts {
	return module.Counts{
		CallNum:    atomic.LoadUint32(&m.calledCount),
		AcceptedNum:  atomic.LoadUint32(&m.acceptedCount),
		CompletedNum: atomic.LoadUint32(&m.completedCount),
		HandlingNum: atomic.LoadUint32(&m.handlingNumber),
	}
}

func (m *myModule) Summary() module.SummaryStruct {
	counts := m.Counts()
	return module.SummaryStruct{
		ID:        m.ID(),
		Called:    counts.CallNum,
		Accepted:  counts.AcceptedNum,
		Completed: counts.CompletedNum,
		Handling:  counts.HandlingNum,
		Extra:     nil,
	}
}

func (m *myModule) IncrCalledCount() {
	atomic.AddUint32(&m.calledCount, 1)
}

func (m *myModule) IncrAcceptedCount() {
	atomic.AddUint32(&m.acceptedCount, 1)
}
func (m *myModule) IncrCompletedCount() {
	atomic.AddUint32(&m.completedCount, 1)
}

func (m *myModule) IncrHandlingNumber() {
	atomic.AddUint32(&m.handlingNumber, 1)
}

func (m *myModule) DecrHandlingNumber() {
	atomic.AddUint32(&m.handlingNumber, ^uint64(0))
}

func (m *myModule) Clear() {
	atomic.StoreUint32(&m.calledCount, 0)
	atomic.StoreUint32(&m.acceptedCount, 0)
	atomic.StoreUint32(&m.completedCount, 0)
	atomic.StoreUint32(&m.handlingNumber, 0)
}






























