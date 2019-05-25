package scheduler

import "mycha/module"

//Args 代表参数容器的统一接口类型
type Args interface {
	Check() error
}

//请求参数的容器
type RequestArgs struct {
	//代表可以接受的域名列表
	AcceptedDomains []string `json:"accepted_primary_domains"`
	//代表可以爬取的最大深度
	MaxDepth uint32 `json:"max_depth"`
}

func (req *RequestArgs) Check() error {
	if req.AcceptedDomains == nil {
		return genError("限定域名列表为空")
	}
	return nil
}

//检查当前的请求参数和另外一个是否一致
func (args *RequestArgs) Same(another *RequestArgs) bool {
	if another == nil {
		return false
	}
	if another.MaxDepth != args.MaxDepth {
		return false
	}
	anotherDomains := another.AcceptedDomains
	anotherDomainsLen := len(another.AcceptedDomains)
	if anotherDomainsLen != len(args.AcceptedDomains) {
		return false
	}
	if anotherDomainsLen > 0 {
		for i,domain := range anotherDomains {
			if domain != args.AcceptedDomains[i] {
				return false
			}
		}
	}
	return true
}

//DataArgs 代表数据相关的参数容器
type DataArgs struct {
	// ReqBufferCap 代表请求缓冲器的容量。
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	// ReqMaxBufferNumber 代表请求缓冲器的最大数量。
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	// RespBufferCap 代表响应缓冲器的容量。
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	// RespMaxBufferNumber 代表响应缓冲器的最大数量。
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_number"`
	// ItemBufferCap 代表条目缓冲器的容量。
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	// ItemMaxBufferNumber 代表条目缓冲器的最大数量。
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	// ErrorBufferCap 代表错误缓冲器的容量。
	ErrorBufferCap uint32 `json:"error_buffer_cap"`
	// ErrorMaxBufferNumber 代表错误缓冲器的最大数量。
	ErrorMaxBufferNumber uint32 `json:"error_max_buffer_number"`
}



func (args *DataArgs) Check() error {
	if args.ReqBufferCap == 0 {
		return genError("请求缓冲器容量为空")
	}
	if args.ReqMaxBufferNumber == 0 {
		return genError("请求缓冲器最大限制为空")
	}
	if args.RespBufferCap == 0 {
		return genError("响应缓冲器容量为空")
	}
	if args.RespMaxBufferNumber == 0 {
		return genError("响应缓冲器最大限制为空")
	}
	if args.ItemBufferCap == 0 {
		return genError("条目缓冲器容量为空")
	}
	if args.ItemMaxBufferNumber == 0 {
		return genError("条目缓冲器最大限制为空")
	}
	if args.ErrorBufferCap == 0 {
		return genError("错误缓冲器容量为空")
	}
	if args.ErrorMaxBufferNumber == 0 {
		return genError("错误缓冲器最大限制为空")
	}
	return nil
}


//信息结构体
type ModuleArgsSummary struct {
	DownloaderListSize int `json:"downloader_list_size"`
	AnalyzerListSize int `json:"analyzer_list_size"`
	PipelineListSize int `json:"pipeline_list_size"`
}


// ModuleArgs 代表组件相关的参数容器的类型。
type ModuleArgs struct {
	// Downloaders 代表下载器列表。
	Downloaders []module.Downloader
	// Analyzers 代表分析器列表。
	Analyzers []module.Analyzer
	// Pipelines 代表条目处理管道管道列表。
	Pipelines []module.Pipeline
}

// Check 用于当前参数容器的有效性。
func (args *ModuleArgs) Check() error {
	if len(args.Downloaders) == 0 {
		return genError("empty downloader list")
	}
	if len(args.Analyzers) == 0 {
		return genError("empty analyzer list")
	}
	if len(args.Pipelines) == 0 {
		return genError("empty pipeline list")
	}
	return nil
}

func (args *ModuleArgs) Summary() ModuleArgsSummary {
	return ModuleArgsSummary{
		DownloaderListSize: len(args.Downloaders),
		AnalyzerListSize:   len(args.Analyzers),
		PipelineListSize:   len(args.Pipelines),
	}
}



















