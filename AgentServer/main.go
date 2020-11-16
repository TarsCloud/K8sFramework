package main

import (
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"tafagent/common"
	"tafagent/crond"
	"tafagent/monitor"
)

func main() {
	// http 服务端，机器资源监控
	server := gin.Default()

	server.Use(common.JSONAppErrorReporter())
	server.Use(common.CORSMiddleware())

	server.GET("/cpu",  	monitor.CPU)
	server.GET("/memory", monitor.Memory)
	server.GET("/disk", 	monitor.Disk)
	server.GET("/host", 	monitor.Host)
	server.GET("/net", 	monitor.Net)
	server.GET("/port", 	monitor.Port)

	// crond日志清理
	dirs := []string{
		"/usr/local/app/taf/app_log/",
		"/usr/local/app/taf/app_log/"}
	patterns := []string{
		"*.log",
		"core.*"}
	crond.WalkDirRemove(dirs, patterns)

	port := os.Getenv("TAF_AGENT_PORT")
	if port == "" {
		port = "8000"
	}
	_ = server.Run(net.JoinHostPort("0.0.0.0", port))
}
