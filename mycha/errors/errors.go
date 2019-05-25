package errors

import (
	"bytes"
	"fmt"
)

type ErrorType string


//错误的类型
const (
	ERROR_TYPE_DOWNLOADER ErrorType  = "下载器错误"
	ERROR_TYPE_ANALYZER  ErrorType = "分析器错误"
	ERROR_TYPE_PIPELINE ErrorType = "条目处理器错误"
	ERROR_TYPE_SCHEDULER ErrorType = "调度器错误"
	ERROR_TYPE_PARAMETER ErrorType = "参数错误"
)


type CrawlerError interface {
	Type() ErrorType
	Error() string
}


//实现全部的错误类型
type myCrawlerError struct {
	errType ErrorType
	errMsg string
	fullMsg string
}

//根据字符串创建
func NewCrawlerError(typeMsg ErrorType,msg string) CrawlerError {
	return &myCrawlerError{
		errType:typeMsg,
		errMsg: msg,
	}
}

//根据错误类型来判断
func NewCrawlerErrorByErr(typeMsg ErrorType,err error) CrawlerError {
	return NewCrawlerError(typeMsg,err.Error())
}


//返回错误的类型
func (merr *myCrawlerError) Type() ErrorType {
	return merr.errType
}


//返回详细的错误
func (merr *myCrawlerError) Error() string {
	if merr.fullMsg == "" {
		merr.creatFullMsg()
	}
	return merr.fullMsg
}


//创建详细的错误信息
func (merr *myCrawlerError) creatFullMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("出现一个错误 类型：")
	if merr.errType != "" {
		buffer.WriteString(string(merr.errType))
	}
	buffer.WriteString(merr.errMsg)
	merr.fullMsg = fmt.Sprintf("%s",buffer.String())
	return
}

























