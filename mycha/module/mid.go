package module

import (
	"fmt"
	"mycha/errors"
	"net"
	"strconv"
	"strings"
)

var DefaultSNGen = NewSNGenertor(1,0)  //生产默认的一个编号生产器

type MID string  //设置默认的类型

var midTemplate = "%s%d|%s"

//参数新的随机序列号
func GenMID(mtype Type, sn uint32,maddr net.Addr)(MID,error) {
	if !LegalType(mtype) {
		msg := fmt.Sprint("不存在这样定义的组件类型: %s", mtype)
		return "",errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,msg)
	}
	letter := legalTypeLetterMap[mtype]  //获取到简写
	var midstr string
	if maddr == nil {
		midstr = fmt.Sprint(midTemplate,letter,sn,"")
		midstr = midstr[:len(midstr)-1]
	} else {
		midstr = fmt.Sprintf(midTemplate,letter,sn,maddr.String())
	}

	return MID(midstr),nil  //强制转换类型
}

//判断mid是否合法
func LegalMID(mid MID) bool {
	if _,err := SplitMID(mid); err == nil {
		return true
	}
	return false
}


//拆解MID 获取到其中的信息
func SplitMID(mid MID)([]string ,error) {
	var letter string
	var ok bool
	var snStr string
	var addr string
	midStr := string(mid)
	if len(midStr) <= 1 {
		return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"拆解的MID长度不正确")
	}
	letter = midStr[:1]
	if _,ok = legalLetterTypeMap[letter];ok {
		return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"MID的组件类型前缀不正确")
	}
	snAndAddr := midStr[1:]
	index := strings.LastIndex(snAndAddr,"|")
	if index < 0  {
		snStr = snAndAddr
		if !legalSN(snStr) {
			return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"码值不正确")
		}
	} else {
		snStr = snAndAddr[:index]
		if !legalSN(snStr) {
			return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"组件的序列号非法")
		}
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr,":")
		if index<=0 {
			return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"组件的序列号的组件地址非法")
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"组件的序列号的地址不正确")
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64); err != nil {
			return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"组件的序列号的地址不正确")
		}

	}
	return []string{letter,snStr,addr},nil
}

//其实就是简单转换为数字 以10进制  转换为64位的数字 如果不能转换 代表其中存在字符串等非法字符
func legalSN(snStr string) bool {
	_,err := strconv.ParseUint(snStr,10,64)
	if err != nil {
		return false
	}
	return  true
}







