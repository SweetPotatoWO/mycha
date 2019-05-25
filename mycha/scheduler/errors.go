package scheduler

import (
	"mycha/errors"
	"mycha/module"
	"mycha/tool/buffer"
)

//生产爬虫错误  传递字符串
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER,errMsg)
}

//产生爬虫错误 传递error类型
func genErrorByError(err error) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER,err.Error())
}

//产生参数错误
func genParameterError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,errMsg)
}



//发送错误到缓存池子中
func sendError(err error, mid module.MID,errorBufferPool buffer.Pool) bool {
	if err == nil || errorBufferPool == nil || errorBufferPool.Closed() {
		return false
	}
	var crawlerError errors.CrawlerError
	var ok bool
	crawlerError,ok = err.(errors.CrawlerError)
	if !ok {
		var moduleType module.Type
		var errorType errors.ErrorType
		ok,moduleType = module.GetType(mid)
		if !ok {
			errorType = errors.ERROR_TYPE_SCHEDULER
		} else {
			switch moduleType {
			case module.TYPE_DOWNLOADER:
				errorType = errors.ERROR_TYPE_DOWNLOADER
			case module.TYPE_PIPELINE:
				errorType = errors.ERROR_TYPE_PIPELINE
			case module.TYPE_ANALYZER:
				errorType = errors.ERROR_TYPE_ANALYZER
			}
		}
		crawlerError = errors.NewCrawlerError(errorType,err.Error())
	}
	if errorBufferPool.Closed() {
		return false
	}
	go func(crawlerError errors.CrawlerError) {
		if err := errorBufferPool.Put(crawlerError); err != nil {
			logger.Warnln("错误缓存池子关闭了 请检查程序是否有问题")
		}
	}(crawlerError)
	return true
}

























