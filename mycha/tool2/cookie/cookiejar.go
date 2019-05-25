package cookie

import (
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)


//构造方法
func NewCookiejar() http.CookieJar {
	options := &cookiejar.Options{PublicSuffixList:&myPublicSuffixList{}}
	cj,_ := cookiejar.New(options)
	return cj
}

type myPublicSuffixList struct {}

func (psl *myPublicSuffixList) PublicSuffix(domain string) string {
	suffix,_ := publicsuffix.PublicSuffix(domain)
	return suffix
}


func (psl *myPublicSuffixList) String() string {
	return "go语言的工具包介绍cookie 介绍"
}