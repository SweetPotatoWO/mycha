package buffer

import (
	"sync"
	"sync/atomic"
)

type Pool interface {
	BufferCap() (num uint32)   //缓存池中的缓存器统一容量
	BufferNumber() (num uint32)  //缓冲器数量
	BufferMaxNumber() (num uint32) //缓冲器最大的数量
	PoolPut(datum interface{}) (err error)  //数据压入数量
	PoolGet() (datum interface{},err error) //数据获取
	Close() (flag bool)  //关闭缓存池子
	Closed() (flag bool)  //缓存池子是否关闭
	DatumNumber() (num uint32) //数据的总类型
}

type pool struct {
	cap uint32
	bufferNumber uint32
	bufferMaxNumber uint32
	datumNumber uint32
	bufferCh chan Buffer
	poolLock sync.RWMutex
	closed uint32
}

func NewBufferPool(cap uint32,max uint32) Pool {
	bu := NewBuffer(cap)
	poolBuf := make(chan Buffer,max)
	poolBuf<-bu
	return &pool{
		cap:cap,
		bufferNumber:0,
		bufferCh:poolBuf,
		bufferMaxNumber:max,
		datumNumber:0,
	}
}

func (p *pool) BufferCap() (num uint32) {
	return p.cap
}

func (p *pool) BufferNumber() (num uint32) {
	return atomic.LoadUint32(&p.bufferNumber)
}

func (p *pool) BufferMaxNumber() (num uint32) {
	return p.bufferMaxNumber
}

func (p *pool) PoolPut(datum interface{}) (err error) {
	if p.Closed() {
		return ErrorPoolClose
	}
	var i uint32
	for val:= range p.bufferCh {
		ok,err := p.putData(val,datum,p.bufferMaxNumber)
		if ok && err != nil {
			break
		}
		i++
		if i == p.bufferMaxNumber {
			err = ErrorPoolBuf
			break
		}
	}
	return
}


func (p *pool) putData(buf Buffer,datum interface{},max uint32) (ok bool,err error) {
	if p.Closed() {
		return false,ErrorPoolClose
	}

	defer func() {   //无论成功还是失败 都将缓存器回收 但新建的缓存器的时候不允许进行这样的操作
		p.poolLock.RLock()
		p.bufferCh<-buf
		p.poolLock.RUnlock()
	}()

	p.poolLock.Lock()
	ok,err = buf.Put(datum)
	if ok && err!=nil {
		atomic.AddUint32(&p.datumNumber,1)
		return ok,err
	}
	if max != p.bufferNumber {
		for  {
			newBuf := NewBuffer(p.cap)
			ok ,err = newBuf.Put(datum)
			if ok && err!= nil{
				atomic.AddUint32(&p.datumNumber,1)
				atomic.AddUint32(&p.bufferNumber,1)
				p.bufferCh<-newBuf
				p.poolLock.Unlock()
				return
			}
		}
	}
	p.poolLock.Unlock()
	return
}

func (p *pool) PoolGet() (datum interface{},err error) {
	if p.Closed() {
		return nil,ErrorPoolClose
	}
	for val := range p.bufferCh {
		res,err := p.getData(val)
		if res != nil && err == nil {
			return res,err
		}
	}
}


func (p *pool) getData(buf Buffer) (datum interface{},err error) {
	if p.Closed() {
		return nil,ErrorPoolClose
	}

	defer func() {

		num,_ := buf.Len()
		if num == 0 && p.bufferNumber>1 {  //当缓存器中的内容为0
			buf.Close()
			atomic.AddUint32(&p.bufferNumber,^uint32(0))
			return
		}

		p.poolLock.RLock()
		p.bufferCh<-buf
		p.poolLock.RUnlock()

	}()

	datum ,err = buf.Get()
	if datum != nil {
		atomic.AddUint32(&p.datumNumber,^uint32(0))
		return
	}
	if err != nil {
		return
	}
	return
}



func (p *pool) Close() (flag bool) {
	if !atomic.CompareAndSwapUint32(&p.closed,0,1) {
		return false
	}
	p.poolLock.Lock()
	defer p.poolLock.Unlock()
	close(p.bufferCh)

	for buf := range p.bufferCh {
		buf.Close()
	}
	return true
}

func (p *pool) Closed() (flag bool) {
	if atomic.LoadUint32(&p.closed) == 1 {
		return true
	}
	return false
}

func (p *pool) DatumNumber() (num uint32) {

	return atomic.LoadUint32(&p.datumNumber)
}


