package module

import "net/http"


type Data interface {
	Valid() bool
}


//自己封装的一个请求的对象
type Request struct {
	httpReq *http.Request
	depth uint32
}


func NewRequest(httpReq *http.Request,depth uint32) *Request {
	return &Request{
		httpReq:httpReq,
		depth:depth,
	}
}


func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}


func (req *Request) Depth() uint32 {
	return req.depth
}

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}


//自己封装的一个响对象
type Response struct {
	httpResp *http.Response
	depth uint32
}


// NewResponse 用于创建一个新的响应实例。
func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp: httpResp, depth: depth}
}

// HTTPResp 用于获取HTTP响应。
func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

// Depth 用于获取响应深度。
func (resp *Response) Depth() uint32 {
	return resp.depth
}

// Valid 用于判断响应是否有效。
func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}


//条目
type Item map[string]interface{}


// Valid 用于判断条目是否有效。
func (item Item) Valid() bool {
	return item != nil
}
