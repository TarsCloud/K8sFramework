
## Select 接口定义

> ServerTree 的 Select请求比较特殊. columns 仅作为占位符使用

> 请求格式

```json
    {
        "params": {
            "columns": ["ServerTree"] 
        }
    } 
```


> 响应结果

```json
{
    "jsonrpc": "2.0",
    "id": "1573698508292",
    "result": [
        {
            "BusinessName": "ZSZQ",
            "BusinessShow": "招商证券",
            "App": [
                {
                    "AppName": "Test",
                    "Server": []
                }
            ]
        },
        {
            "BusinessName": "YHZQ",
            "BusinessShow": "银河证券2",
            "App": []
        },
        {
            "BusinessName": "ZXZQ",
            "BusinessShow": "中信证券",
            "App": [
                {
                    "AppName": "Agent",
                    "Server": []
                }
            ]
        },
        {
            "BusinessName": "",
            "BusinessShow": "",
            "App": [
                {
                    "AppName": "Login",
                    "Server": []
                },
                {
                    "AppName": "News",
                    "Server": []
                }
            ]
        }
    ]
}
```