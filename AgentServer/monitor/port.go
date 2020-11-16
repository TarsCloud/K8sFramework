package monitor

import (
	"net"
	"strconv"
	"strings"
	"tafagent/common"
	"time"
)

type PortReq struct {
	Host string 	`json:"host"`
	Port int    	`json:"port"`
}

type PortRsp struct {
	ErrInfo common.ErrorInfo 	`json:"err_info"`
	Port 	int					`json:"port"`
	InUse   bool      			`json:"in_use"`
}

type PortData struct {
	req *PortReq
	rsp *PortRsp
}

func (pd * PortData) NewPortReq() *PortReq {
	pd.req = &PortReq{}
	return pd.req
}

func (pd * PortData) NewPortRsp() *PortRsp {
	pd.rsp = &PortRsp{}
	return pd.rsp
}


func (pd * PortData) GetAvailPort(ip string, port int) *PortRsp {
	rsp := pd.NewPortRsp()

	host := []string {ip, strconv.Itoa(port)}
	addr, err := net.ResolveTCPAddr("tcp", strings.Join(host, ":"))
	if err != nil {
		rsp.ErrInfo.ErrCode = -1
		rsp.ErrInfo.ErrMsg 	= err.Error()
		return rsp
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		rsp.ErrInfo.ErrCode = -1
		rsp.ErrInfo.ErrMsg 	= err.Error()
		return rsp
	}
	_ = l.Close()

	rsp.Port 	= l.Addr().(*net.TCPAddr).Port

	return rsp
}

func (pd * PortData) CheckPortInUse(host string, port int) *PortRsp {
	rsp := pd.NewPortRsp()
	rsp.Port = port

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), time.Second)
	if err == nil && conn != nil {
		conn.Close()
		rsp.InUse = true
	}

	return rsp
}
