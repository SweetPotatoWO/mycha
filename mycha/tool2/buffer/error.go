package buffer

import "errors"

var (
	ErrorBubfferClose = errors.New("缓存已经关闭了")
	ErrorBubfferPutErr = errors.New("缓存器装入失败")
	ErrorBubfferGetErr = errors.New("缓存器获取元素失败")
	ErrorBubfferCloseNumErr  = errors.New("初始化参数错误")
)

var (
	ErrorPoolClose = errors.New("缓冲池已经关闭")
	ErrorPoolBuf = errors.New("全部缓冲器都不可以用,请重新检查")
)

