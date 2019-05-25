package buffer

import (
	"sync"
	"sync/atomic"
)

type Buffer interface {

	Put(datum interface{}) (flag bool,err error)  //压入函数
	Get() (datum interface{},err error)   //取出
	Len() (num uint32,err error) //当前容量
	Cap() (num uint32,err error)  //最大容量
	Close() (flag bool,err error)  //关闭
	Closed() (flag bool)   //关闭是否 true 表示关闭 false 表示未关闭
}


type buffer struct {
	ch chan interface{}   //通道
	closed uint32  //是否关闭
	bufferLock sync.RWMutex  //读写锁
}

func NewBuffer(cap uint32) Buffer {

	if cap == 0 {
		return nil
	}
	return &buffer{
		ch:make(chan interface{},cap),
	}
}


func (b *buffer) Put(datum interface{}) (flag bool,err error) {
	b.bufferLock.Lock()
	if b.Closed() {
		return false,ErrorBubfferClose
	}
	defer func() {
		b.bufferLock.Unlock()
	}()
	select {
		case b.ch<-datum:
			return true,nil
	default:
		return false,ErrorBubfferPutErr
	}
}


func (b *buffer) Get() (datum interface{},err error) {
	b.bufferLock.RLock()
	if b.Closed() {
		return nil,ErrorBubfferClose
	}
	defer func() {
		b.bufferLock.RUnlock()
	}()
	select {
	case res := <-b.ch:
		return res,nil
	default:
		return nil,ErrorBubfferGetErr
	}
}


func (b *buffer) Len() (num uint32,err error) {
	if b.Closed() {
		return 0,ErrorBubfferClose
	}
	return uint32(len(b.ch)),nil
}

func (b *buffer) Cap() (num uint32,err error) {
	if b.Closed() {
		return 0,ErrorBubfferClose
	}
	return uint32(cap(b.ch)),nil
}


func (b *buffer) Close() (flag bool,err error) {
	if atomic.CompareAndSwapUint32(&b.closed,0,1) {
		b.bufferLock.Lock()
		close(b.ch)
		b.bufferLock.Unlock()
		return true,nil
	}
	return false,ErrorBubfferCloseNumErr
}

func (b *buffer) Closed() (flag bool) {
	if atomic.LoadUint32(&b.closed) == 0 {
		return false
	}
	return true
}





