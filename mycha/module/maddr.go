package module

import (
	"mycha/errors"
	"net"
	"strconv"
)

type mAddr struct {
	network string
	address string
}

func (m *mAddr) Network() string {

	return m.network
}

func (m *mAddr) String() string {
	return m.address
}

func NewAddr(network string,ip string,port uint32) (net.Addr,error) {
	if network!= "http" && network != "https" {
		errMsg := errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"错误的地址")
		return nil,errMsg
	}
	if parseIP := net.ParseIP(ip);parseIP == nil {
		errMsg := errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"错误的地址")
		return nil,errMsg
	}

	return &mAddr{
		network:network,
		address:ip+":"+strconv.Itoa(int(port)),
	},nil
}