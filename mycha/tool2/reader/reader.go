package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)


//多重读取器的接口
type MultipleReader interface {
	Reader() io.ReadCloser  //IO接口的阅读和关闭
}

//实现类型
type multipleReader struct {
	data []byte
}

//构造方法
func NewMultipleReader(reader io.Reader) (MultipleReader,error) {
	var data []byte
	var err error

	if reader != nil {
		data,err = ioutil.ReadAll(reader)
		if err != nil {
			return nil,fmt.Errorf("不能创建一个多重读取器")
		}

	} else {
		data = []byte{}
	}

	return &multipleReader{
		data:data,
	},nil

}


func (r *multipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(r.data))
}











