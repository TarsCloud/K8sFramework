module tafagent

go 1.14

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/shirou/gopsutil v2.20.8+incompatible
	golang.org/x/sys v0.0.0-20200908134130-d2e65c121b96 // indirect
)

replace (
	tafagent/common v0.0.0 => ./common
	tafagent/crond v0.0.0 => ./crond
	tafagent/monitor v0.0.0 => ./monitor
)
