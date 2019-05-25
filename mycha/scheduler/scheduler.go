package scheduler

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"gopcp.v2/chapter5/cmap"
	"strings"

	"mycha/helper/log"
	"mycha/module"
	"mycha/tool/buffer"
	"net/http"
	"sync"
)

var logger = log.DLogger()



//调度器的接口
type scheduler interface {
	Init(requestArgs RequestArgs , dataArgs DataArgs,
		moduleArgs ModuleArgs) (err error) //初始化调度器
	Start(firstHTTPReq *http.Request) (err error) //启动调度器
	Stop() (err error) //停止调度器
	Status() Status //当前的状态
	ErrorChan() <-chan error  //错误通道?
	Idle() bool //用来判断所有的模块都处于空闲状态
	Summary() SchedSummary
}

func NewScheduler() scheduler{
	return &myScheduler{}
}

type myScheduler struct {
	//最大的深度
	maxDepth uint32
	//允许接受的域名哈希
	acceptedDomainMap cmap.ConcurrentMap  //用了第三方的 并发安全的map
	//registrar 代表组件注册器
	registrar module.Registrar
	//请求缓冲池
	regBufferPool buffer.Pool
	//响应缓冲池
	respBufferPool buffer.Pool
	//代表条目缓冲池
	itemBufferPool buffer.Pool
	//错误缓冲池
	errorBufferPool buffer.Pool
	//链接的字典
	urlMap cmap.ConcurrentMap
	//ctx 代表上下文 用于感知调度器停止
	ctx context.Context
	// cancelFunc 代表取消函数，用于停止调度器。
	cancelFunc context.CancelFunc
	//状态
	status Status
	//锁
	statusLock sync.RWMutex
	//摘要
	summary SchedSummary
}

func (sched *myScheduler) Stop() (err error) {
	logger.Info("Stop scheduler...")
	// 检查状态。
	logger.Info("Check status for stop...")
	var oldStatus Status
	oldStatus, err =
		sched.checkAndSetStatus(SCHED_STATUS_STOPPING)
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_STOPPED
		}
		sched.statusLock.Unlock()
	}()
	if err != nil {
		return
	}
	sched.cancelFunc()
	sched.regBufferPool.Close()
	sched.respBufferPool.Close()
	sched.itemBufferPool.Close()
	sched.errorBufferPool.Close()
	logger.Info("Scheduler has been stopped.")
	return nil
}

func (sched *myScheduler) Status() Status {
	var status Status
	sched.statusLock.RLock()
	status = sched.status
	sched.statusLock.RUnlock()
	return status
}

func (sched *myScheduler) Init(
	requestArgs RequestArgs,
	dataArgs DataArgs,
	moduleArgs ModuleArgs) (err error) {

	logger.Info("初始化调度器")
	var oldStatus Status
	oldStatus,err = sched.checkAndSetStatus(SCHED_STATUS_INITIALIZING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_INITIALIZED
		}
		sched.statusLock.Unlock()
	}()
	logger.Info("检查请求队列数组")
	if err = requestArgs.Check(); err != nil {
		return err
	}
	logger.Info("检查数据队列数组")
	if err = dataArgs.Check(); err != nil {
		return err
	}
	logger.Info("数据队列检查完毕")
	logger.Info("检查组件队列")
	if err = moduleArgs.Check(); err != nil {
		return err
	}
	logger.Info("组件队列检查完毕")
	// 初始化内部字段。
	logger.Info("检查调度器的参数...")
	if sched.registrar == nil {
		sched.registrar = module.NewRegistrar()
	} else {
		sched.registrar.Clear()
	}
	sched.maxDepth = requestArgs.MaxDepth
	logger.Infof("--最大爬取深度:%d",sched.maxDepth)
	sched.acceptedDomainMap,_ = cmap.NewConcurrentMap(1,nil)
	for _,domain := range requestArgs.AcceptedDomains {
		sched.acceptedDomainMap.Put(domain, struct{}{}) //为每个域名填充上一个空的结构体
	}
	logger.Infof("--允许的域名名单:%v",requestArgs.AcceptedDomains)
	sched.urlMap,_ = cmap.NewConcurrentMap(16,nil)
	logger.Infof("--链接的的队列长度长度:%d concurrency %d",sched.urlMap.Len(),sched.urlMap.Concurrency())
	sched.initBufferPool(dataArgs)   //一个填充数据到调度器中的方法
	sched.resetContext()  //重置上下文
	sched.summary =
		newSchedSummary(requestArgs, dataArgs, moduleArgs, sched)   //生产摘要
	// 注册组件。
	logger.Info("注册组件...")
	if err = sched.registerModules(moduleArgs); err != nil {
		return err
	}
	logger.Info("调度器已经初始化完毕.")
	return nil
}


//注册组件 注册的registrar中 生产序号 保存到列表中
func (sched *myScheduler) registerModules(moduleArgs ModuleArgs) error {
	for _,d:= range moduleArgs.Downloaders {
		if d == nil {
			continue
		}
		ok,err := sched.registrar.Register(d)
		if err != nil {
			return genErrorByError(err)
		}
		if !ok {
			errMsg := fmt.Sprintf("存在不能够注册成功的下载器实体,出错的MID为%d",d.ID())
			return genError(errMsg)
		}
	}
	logger.Infof("全部的下载器已经注册了(num%d)",len(moduleArgs.Downloaders))
	for _,a:=range moduleArgs.Analyzers {
		if a== nil {
			continue
		}
		ok,err := sched.registrar.Register(a)
		if err != nil {
			return genErrorByError(err)
		}
		if !ok {
			errMsg := fmt.Sprintf("存在不能够注册成功的分析器主题，出错的MID为%d",a.ID())
			return genError(errMsg)
		}
	}
	logger.Infof("全部的分析器已经注册了(num%d)",len(moduleArgs.Analyzers))
	for _,p := range moduleArgs.Pipelines {
		if p == nil {
			continue
		}
		ok,err := sched.registrar.Register(p)
		if err != nil  {
			return genErrorByError(err)
		}
		if !ok  {
			errMsg := fmt.Sprintf("存在条目管道实例注册不成功，出错的MID为%d",p.ID())
			return genError(errMsg)
		}
	}
	logger.Infof("全部的条目通道都注册成功(num%d)",len(moduleArgs.Pipelines))
	return nil
}



//重置上下文  直接将ctx 重置为祖节点
func (sched *myScheduler) resetContext() {
	sched.ctx,sched.cancelFunc = context.WithCancel(context.Background())
}


//DataArgs 包含各个容器的配置参数 好像是 加个？
func (sched *myScheduler) initBufferPool(dataArgs DataArgs) {
	if sched.regBufferPool != nil && !sched.regBufferPool.Closed() {   //请求缓存池不为空 且不为关闭状态 关闭 并重置
		sched.regBufferPool.Close() //关闭
	}
	sched.regBufferPool,_ = buffer.NewPool(dataArgs.ReqBufferCap,dataArgs.ReqMaxBufferNumber)
	logger.Infof("-- 请求缓存池子: bufferCap(容量): %d, maxBufferNumber(当前数): %d",
		sched.regBufferPool.BufferCap(), sched.regBufferPool.MaxBufferNumber())
	// 初始化响应缓冲池。
	if sched.respBufferPool != nil && !sched.respBufferPool.Closed() {
		sched.respBufferPool.Close()
	}
	sched.respBufferPool, _ = buffer.NewPool(
		dataArgs.RespBufferCap, dataArgs.RespMaxBufferNumber)
	logger.Infof("-- =响应缓存池 : bufferCap: %d, maxBufferNumber: %d",
		sched.respBufferPool.BufferCap(), sched.respBufferPool.MaxBufferNumber())
	// 初始化条目缓冲池。
	if sched.itemBufferPool != nil && !sched.itemBufferPool.Closed() {
		sched.itemBufferPool.Close()
	}
	sched.itemBufferPool, _ = buffer.NewPool(
		dataArgs.ItemBufferCap, dataArgs.ItemMaxBufferNumber)
	logger.Infof("-- 条目缓存池: bufferCap: %d, maxBufferNumber: %d",
		sched.itemBufferPool.BufferCap(), sched.itemBufferPool.MaxBufferNumber())
	// 初始化错误缓冲池。
	if sched.errorBufferPool != nil && !sched.errorBufferPool.Closed() {
		sched.errorBufferPool.Close()
	}
	sched.errorBufferPool, _ = buffer.NewPool(
		dataArgs.ErrorBufferCap, dataArgs.ErrorMaxBufferNumber)
	logger.Infof("-- 错误缓存池: bufferCap: %d, maxBufferNumber: %d",
		sched.errorBufferPool.BufferCap(), sched.errorBufferPool.MaxBufferNumber())

}


func (sched *myScheduler) Start(firstHTTPReq *http.Request) (err error) {
	defer func() {   //处理恐慌
		if p:= recover(); p!= nil {
			errMsg := fmt.Sprintf("出现调度器错误: %sched", p)
			logger.Fatal(errMsg)
			err = genError(errMsg)
		}
	}()
	logger.Info("开始调度")
	logger.Info("检查调度器状态")
	var oldStatus Status
	oldStatus,err = sched.checkAndSetStatus(SCHED_STATUS_STARTING)  //检查并且设置
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()

	if err != nil {
		return
	}
	//检查参数
	logger.Info("检查第一个HTTP请求")
	if firstHTTPReq == nil {
		err = genParameterError("第一个请求的参数错误")
		return
	}
	logger.Info("第一个HTTP请求检查完毕")
	logger.Info("获取请求主域名")
	logger.Infof("--host:%s",firstHTTPReq.Host)
	var primaryDomain string
	primaryDomain,err = getPrimaryDomain(firstHTTPReq.Host)  //检查域名是否正确
	if err != nil {
		return
	}
	logger.Infof("主域名为:%s",primaryDomain)
	sched.acceptedDomainMap.Put(primaryDomain, struct {}{})
	if err  = sched.checkBufferPoolForStart(); err != nil {
		return
	}
	sched.download()   //循环的读取缓存池子的参数
	sched.analyze()
	sched.pick()
	logger.Info("调度器启动成功")
	firstReq := module.NewRequest(firstHTTPReq,0)
	sched.sendReq(firstReq)  //读取第一个请求放入池子中
	return nil
}


func (sched *myScheduler) download() {
	go func() {
		for {
			if sched.canceled() {  //检查上下文是否关闭 如果关闭代表取消全部的goroutine
				break
			}
			datum,err := sched.respBufferPool.Get()  //从池子中获取到一个节点
			if err != nil {
				logger.Warnf("请求的缓存池子被关闭了")
				break
			}
			req,ok := datum.(*module.Request)
			if !ok {
				errMsg := fmt.Sprintf("缓存池子中的节点无法转换成正常的请求类型 该类型为 %T",datum)
				sendError(errors.New(errMsg),"",sched.errorBufferPool)
			}
			sched.downloadOne(req)
		}
	}()
}


func (sched *myScheduler) downloadOne(req *module.Request) {
	if req == nil {
		return
	}
	if sched.canceled() {
		return
	}
	m,err := sched.registrar.Get(module.TYPE_DOWNLOADER)
	if err != nil || m == nil {
		errMsg := fmt.Sprintf("无法获取到下载器: %s", err)
		sendError(errors.New(errMsg), "", sched.errorBufferPool)
		sched.sendReq(req)
		return
	}
	downloader,ok := m.(module.Downloader)  //类型断言 断言为Downloader 方法 同个结构体继承多个接口时需要
	if !ok {
		errMsg := fmt.Sprintf("断言下载器类型是 类型和编号为: %T (MID: %s)",
			m, m.ID())
		sendError(errors.New(errMsg), m.ID(), sched.errorBufferPool)
		sched.sendReq(req)
		return
	}
	resp,err := downloader.Download(req)
	if resp != nil {
		sendResq(resp,sched.respBufferPool)
	}
	if err != nil {
		sendError(err, m.ID(), sched.errorBufferPool)
	}

}





//analyze 会从响应缓存池子中取出响应并解析
//然后把得到最终的结果 比如一些对的格式的程序 或者 页面存在的新的请求
func (sched *myScheduler) analyze()  {
	go func() {
		for {
			if sched.canceled() {
				break
			}
			datum,err := sched.respBufferPool.Get()
			if err != nil {
				logger.Warnln("响应缓存池已经关闭，丢弃这个响应请求")
				break
			}
			resp,ok := datum.(*module.Response)
			if !ok {
				errMsg := fmt.Sprintf("无法转换 response type: %T", datum)
				sendError(errors.New(errMsg), "", sched.errorBufferPool)
			}
			sched.analyzeOne(resp)
		}
	}()
}

//每次实际的处理
func (sched *myScheduler) analyzeOne(resp *module.Response)  {
	if resp == nil {
		return
	}
	if sched.canceled() {
		return
	}

	m, err := sched.registrar.Get(module.TYPE_ANALYZER)
	if err != nil || m == nil {
		errMsg := fmt.Sprintf("无法获取到分析器: %s", err)
		sendError(errors.New(errMsg), "", sched.errorBufferPool)
		sendResq(resp, sched.respBufferPool)
		return
	}
	analyzer, ok := m.(module.Analyzer)
	if !ok {
		errMsg := fmt.Sprintf("incorrect analyzer type: %T (MID: %s)",
			m, m.ID())
		sendError(errors.New(errMsg), m.ID(), sched.errorBufferPool)
		sendResq(resp, sched.respBufferPool)
		return
	}
	dataList, errs := analyzer.Analyze(resp)
	if dataList != nil {
		for _, data := range dataList {
			if data == nil {
				continue
			}
			switch d := data.(type) {
			case *module.Request:
				sched.sendReq(d)
			case module.Item:
				sendItem(d, sched.itemBufferPool)
			default:
				errMsg := fmt.Sprintf("Unsupported data type %T! (data: %#v)", d, d)
				sendError(errors.New(errMsg), m.ID(), sched.errorBufferPool)
			}
		}
	}

	if errs != nil {
		for _, err := range errs {
			sendError(err, m.ID(), sched.errorBufferPool)
		}
	}
}



// sendItem 会向条目缓冲池发送条目。
func sendItem(item module.Item, itemBufferPool buffer.Pool) bool {
	if item == nil || itemBufferPool == nil || itemBufferPool.Closed() {
		return false
	}
	go func(item module.Item) {
		if err := itemBufferPool.Put(item); err != nil {
			logger.Warnln("The item buffer pool was closed. Ignore item sending.")
		}
	}(item)
	return true
}

//判断各个组件是否处于空闲阶段
func (sched *myScheduler) Idle() bool {
	moduleMap := sched.registrar.GetAll()
	for _, module := range moduleMap {
		if module.Handling() > 0 {
			return false
		}
	}
	if sched.regBufferPool.Total() > 0 ||
		sched.respBufferPool.Total() > 0 ||
		sched.itemBufferPool.Total() > 0 {
		return false
	}
	return true
}


//信息摘要
func (sched *myScheduler) Summary() SchedSummary {
	return sched.summary
}






// pick 会从条目缓冲池取出条目并处理。
func (sched *myScheduler) pick() {
	go func() {
		for {
			if sched.canceled() {
				break
			}
			datum, err := sched.itemBufferPool.Get()
			if err != nil {
				logger.Warnln("The item buffer pool was closed. Break item reception.")
				break
			}
			item, ok := datum.(module.Item)
			if !ok {
				errMsg := fmt.Sprintf("incorrect item type: %T", datum)
				sendError(errors.New(errMsg), "", sched.errorBufferPool)
			}
			sched.pickOne(item)
		}
	}()
}

// pickOne 会处理给定的条目。
func (sched *myScheduler) pickOne(item module.Item) {
	if sched.canceled() {
		return
	}
	m, err := sched.registrar.Get(module.TYPE_PIPELINE)
	if err != nil || m == nil {
		errMsg := fmt.Sprintf("couldn't get a pipeline pipline: %s", err)
		sendError(errors.New(errMsg), "", sched.errorBufferPool)
		sendItem(item, sched.itemBufferPool)
		return
	}
	pipeline, ok := m.(module.Pipeline)
	if !ok {
		errMsg := fmt.Sprintf("incorrect pipeline type: %T (MID: %s)",
			m, m.ID())
		sendError(errors.New(errMsg), m.ID(), sched.errorBufferPool)
		sendItem(item, sched.itemBufferPool)
		return
	}
	errs := pipeline.Send(item)
	if errs != nil {
		for _, err := range errs {
			sendError(err, m.ID(), sched.errorBufferPool)
		}
	}
}







//发送响应的内容到响应缓存池
func sendResq(resp *module.Response,respBufferPool buffer.Pool) bool {
	if resp == nil || respBufferPool == nil || respBufferPool.Closed() {
		return false
	}
	go func(resp *module.Response) {
		if err := respBufferPool.Put(resp); err != nil {
			logger.Warnln("响应缓存池子意外关闭")
		}
	} (resp)
	return true
}


func (sched *myScheduler) sendReq(req *module.Request) bool {
	if req == nil {
		return false
	}
	if sched.canceled() {  //上下文是否关闭
		return false
	}
	httpReq:=req.HTTPReq()
	if httpReq == nil {
		logger.Warnln("忽略这个请求，这个请求是空的")
		return false
	}
	reqURL := httpReq.URL
	if reqURL == nil {
		logger.Warnln("忽悠这个请求，这个请求的链接为空")
		return false
	}
	scheme := strings.ToLower(reqURL.Scheme)
	if scheme != "http" && scheme != "https" {
		logger.Warnf("忽悠这个请求! 链接的前缀为 %q, 但必须是 %q or %q. (URL: %s)\n",
			scheme, "http", "https", reqURL)
		return false
	}
	if v := sched.urlMap.Get(reqURL.String()); v != nil { //一个链接对应一个请求的内容 如果获取到的话 就说明成功
		logger.Warnf("忽略这个请求! 请求的链接已经请求过 . (URL: %s)\n", reqURL)
		return false
	}
	pd, _ := getPrimaryDomain(httpReq.Host)   //获取到全部的完整域名
	if sched.acceptedDomainMap.Get(pd) == nil {   //是否在允许请求的域名列表
		if pd == "bing.net" {
			panic(httpReq.URL)
		}
		logger.Warnf("Ignore the request! Its host %q is not in accepted primary domain map. (URL: %s)\n",
			httpReq.Host, reqURL)
		return false
	}
	if req.Depth() > sched.maxDepth {   //请求的深度
		logger.Warnf("Ignore the request! Its depth %d is greater than %d. (URL: %s)\n",
			req.Depth(), sched.maxDepth, reqURL)
		return false
	}
	go func(req *module.Request) {
		if err := sched.regBufferPool.Put(req); err != nil {
			logger.Warnln("The request buffer pool was closed. Ignore request sending.")
		}
	}(req)
	sched.urlMap.Put(reqURL.String(), struct{}{})
	return true


}





















//检查上下文是否取消
func (sched *myScheduler) canceled() bool {
	select {
	case <-sched.ctx.Done():
		return true
	default:
		return false
	}
}



// checkBufferPoolForStart 会检查缓冲池是否已为调度器的启动准备就绪。
// 如果某个缓冲池不可用，就直接返回错误值报告此情况。
// 如果某个缓冲池已关闭，就按照原先的参数重新初始化它。
func (sched *myScheduler) checkBufferPoolForStart() error {
	if sched.regBufferPool == nil {
		return genError("空的请求缓存池")
	}
	if sched.regBufferPool != nil && sched.regBufferPool.Closed() {
		sched.regBufferPool,_ = buffer.NewPool(
			sched.regBufferPool.BufferCap(),
			sched.regBufferPool.MaxBufferNumber())
	}
	if sched.respBufferPool == nil {
		return genError("空的响应缓存池")
	}
	if sched.respBufferPool != nil && sched.respBufferPool.Closed() {
		sched.respBufferPool, _ = buffer.NewPool(
			sched.respBufferPool.BufferCap(), sched.respBufferPool.MaxBufferNumber())
	}
	// 检查条目缓冲池。
	if sched.itemBufferPool == nil {
		return genError("空的条目缓存池")
	}
	if sched.itemBufferPool != nil && sched.itemBufferPool.Closed() {
		sched.itemBufferPool, _ = buffer.NewPool(
			sched.itemBufferPool.BufferCap(), sched.itemBufferPool.MaxBufferNumber())
	}
	// 检查错误缓冲池。
	if sched.errorBufferPool == nil {
		return genError("空的错误缓存池")
	}
	if sched.errorBufferPool != nil && sched.errorBufferPool.Closed() {
		sched.errorBufferPool, _ = buffer.NewPool(
			sched.errorBufferPool.BufferCap(), sched.errorBufferPool.MaxBufferNumber())
	}
	return nil

}





// checkAndSetStatus 用于状态的检查，并在条件满足时设置状态。
func (sched *myScheduler) checkAndSetStatus(wantedStatus Status) (oldStatus Status, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	oldStatus = sched.status
	err = checkStatus(oldStatus, wantedStatus, nil)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}



func (sched *myScheduler) ErrorChan() <-chan error {
	errBuffer := sched.errorBufferPool
	errCh := make(chan error,errBuffer.BufferCap())
	go func(errBuffer buffer.Pool,errCh chan error) {
		for {
			if sched.canceled() {
				close(errCh)
				break
			}

			datum, err := errBuffer.Get()
			if err != nil {
				logger.Warnln("The error buffer pool was closed. Break error reception.")
				close(errCh)
				break
			}
			err, ok := datum.(error)
			if !ok {
				errMsg := fmt.Sprintf("incorrect error type: %T", datum)
				sendError(errors.New(errMsg), "", sched.errorBufferPool)
				continue
			}
			if sched.canceled() {
				close(errCh)
				break
			}
			errCh <- err
		}
	}(errBuffer,errCh)
	return errCh
}






































































