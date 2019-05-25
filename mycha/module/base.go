package module

import "net/http"

type Counts struct {
	//调用次数
	CallNum uint32
	//当前框架接受了多少个任务
	AcceptedNum uint32
	//当前完成多少个任务
	CompletedNum uint32
	//当前正在处理多少个任务
	HandlingNum uint32
}

//组件摘要的结构体
type SummaryStruct struct {
	ID MID		`json:"id"`
	Called uint32 `json:"called"`
	Accepted uint32 `json:"accepted"`
	Completed uint32 `json:"completed"`
	Handling uint32 `json:"handling"`
	Extra interface{} `json:"extra,omitempty"`
}
//统一的组件接口
type Module interface {
	ID() MID   //某个组件的ID
	Addr() string //某个组件的地址
	Score() uint32  //获取到当前组件的评分
	SetScore( score uint32)  //设置组件的评分
	ScoreCalculator() CalculateScore //设置组件的评分规制
	CallCount() uint32 //获取调用的计数
	AcceptedCount() uint32 //获取到接受的计数
	Completed() uint32 //完成的计数
	Handling() uint32 //实现的计数
	Counts() Counts  //获取到的全部的计数对象
	Summary() SummaryStruct //获取信息摘要
}

//下载器的继承接口
type Downloader interface {
	Module
	Download(req *Request) (*Response,error) //下载函数
}


//分析器的继承接口
type Analyzer interface {
	Module
	//每个请求可能对应不同的处理函数切片
	RespParsers() []ParseResponse
	//分析函数
	Analyze(resp *Response) ([]Data,[]error)
}
//函数列表类型
type ParseResponse func(httpResq *http.Response,respDepth uint32) ([]Data,[]error)

//处理条目通道
type Pipeline interface {
	Module
	ItemProcessors() []ProcessItem  //处理条目的函数
	Send(item Item) [] error  //发送条目？
	FailFast() bool // 是否快速错误?
	SetFailFast(failFast bool) //设置快速错误?
}
//条目处理函数
type ProcessItem func(item Item) (result Item,err error)














































