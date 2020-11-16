
## do 接口定义

### Action StartServer
> 用于通知指定的 指定Pod 执行 StartServer 操作
```gotemplate
    type CommandMetadata struct {
    	ServerApp  string   `json:"ServerApp"`
   		ServerName string `json:"ServerName"`
        PodIps     []string  `json:"PodIps"`
    }
```

> 响应结果
```json
{
    "result" : [
        {
            "PodIp": "xxx",
            "Result":"xxx"
        },
        {
            "PodIp": "xxx",
            "Result":"xxx"
        }
    ]
}
```

### Action StopServer
> 用于通知指定的 指定Pod 执行 StopServer 操作
```gotemplate
    type CommandMetadata struct {
    	ServerApp  string   `json:"ServerApp"`
   		ServerName string `json:"ServerName"`
        PodIps     []string  `json:"PodIps"`
    }
```

> 响应结果
```json
{
    "result" : [
        {
            "PodIp": "xxx",
            "Result":"xxx"
        },
        {
            "PodIp": "xxx",
            "Result":"xxx"
        }
    ]
}
```

### Action NotifyServer
> 用于通知指定的 指定Pod 执行 NotifyServer 操作
```gotemplate
    type CommandMetadata struct {
    	ServerApp  string   `json:"ServerApp"`
   		ServerName string   `json:"ServerName"`
        Command    string   `json:"Command"`
        PodIps     []string  `json:"PodIps"`
    }
```

> 响应结果
```json
{
    "result" : [
        {
            "PodIp": "xxx",
            "Result":"xxx"
        },
        {
            "PodIp": "xxx",
            "Result":"xxx"
        }
    ]
}
```

