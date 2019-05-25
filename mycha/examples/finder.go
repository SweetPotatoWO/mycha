package main

import (
	"flag"
	"fmt"


	"mycha/helper/log"
    sched "mycha/scheduler"
	"os"
	"strings"
)

var (
	firstURL string
	domains string
	depth uint
	dirPath string
)

var logger = log.DLogger()

func init() {
	flag.StringVar(&firstURL, "first", "http://zhihu.sogou.com/zhihu?query=golang+logo",
		"你想抓取的第一个链接.")
	flag.StringVar(&domains, "domains", "zhihu.com",
		"The primary domains which you accepted. "+
			"Please using comma-separated multiple domains.")
	flag.UintVar(&depth, "depth", 3,
		"The depth for crawling.")
	flag.StringVar(&dirPath, "dir", "./pictures",
		"The path which you want to save the image files.")
}


func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tfinder [flags] \n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()  //命令赋值
	//新建一个调度器
	scheduler := sched.NewScheduler()
	//能接受的域名列表 主域名
	domainParts := strings.Split(domains,",")
	acceptedDomains := []string{}
	for _,domain := range domainParts {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			acceptedDomains = append(acceptedDomains,domain)
		}
	}

	//请求的实例的配置？？
	requestArgs := sched.RequestArgs{
		AcceptedDomains:acceptedDomains,
		MaxDepth:uint32(depth),
	}
	//数据的实例的配置??
	dataArgs := sched.DataArgs{
		ReqBufferCap:         50,
		ReqMaxBufferNumber:   1000,
		RespBufferCap:        50,
		RespMaxBufferNumber:  10,
		ItemBufferCap:        50,
		ItemMaxBufferNumber:  100,
		ErrorBufferCap:       50,
		ErrorMaxBufferNumber: 1,
	}

	//downloaders,err :=



}



























