package monitor

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tafagent/common"
)

func Port(c *gin.Context) {
	pd := PortData{}

	req := pd.NewPortReq()

	req.Host = c.DefaultQuery("host", "")
	if req.Host == "" {
		req.Host = "127.0.0.1"
	}
	req.Port, _ = strconv.Atoi(c.DefaultQuery("port", ""))

	var rsp *PortRsp

	if req.Port == 0 {
		rsp = pd.GetAvailPort(req.Host, 0)
	} else {
		rsp = pd.CheckPortInUse(req.Host, req.Port)
	}

	common.WriteJsonRsp(c.Writer, rsp)
}

func CPU(c *gin.Context) {
	td := TopData{}
	common.WriteJsonRsp(c.Writer, td.getCpuInfo())
}
func Memory(c *gin.Context) {
	td := TopData{}
	common.WriteJsonRsp(c.Writer, td.getMemInfo())
}
func Disk(c *gin.Context) {
	td := TopData{}
	common.WriteJsonRsp(c.Writer, td.getDiskInfo())
}

func Host(c *gin.Context) {
	td := TopData{}
	common.WriteJsonRsp(c.Writer, td.getHostInfo())
}

func Net(c *gin.Context) {
	td := TopData{}
	common.WriteJsonRsp(c.Writer, td.getNetInfo())
}
