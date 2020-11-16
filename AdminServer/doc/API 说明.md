# RPC 说明

```
    RPC 风格以JsonRPC为基础，参考　Restful 思想，将所有操作抽象成对资源的增删改查(Create ,Delete ,Select,Update,Patch)方法, 针对某些不太容易转换会资源操作的情况,额外提供了 do 方法.
    每种方法又定义了如下的模板请求格式.
```

## 操作

### create
> create 表示新建某种资源的对象

+ 请求
```json
{
    "jsonrpc":2.0,
    "id":1234,
    "method":RPCRequestMethodCreate,
    "params":{
        "token":"xxxx",
        "kind":"deploy",
        "metadata":{"Name":"李四","Age":11}
    }
}
```
+ 请求格式说明
```text
    "method":RPCRequestMethodCreate ,是 指请求的的目标是创建一种资源
    "token":"xxxx", token　是在登录操作之后返回的　登录标识符．所有的 api操作都要携带此字段
    "kind"， 是指资源的类型
    "metadata" 是创建资源对象需要提供的参数 ,具体字段在资源描述文档内定义
    示例 请求等价 SQL: insert into deploy (name,age) values('李四‘,11)
```

+ 响应

```json
{
    "jsonrpc" : 2.0,
    "id" : 1234,
    "result" : 0
}
```

+ 响应格式说明
```text
    一般情况下，create 资源成功后，会返回 "result": 0 ,  如果有特殊情况,会在资源描述文档中说明
    如果出错，则参考　jsonrpc 错误返回情况
```

### select
> select 用于查询某种资源的对象
+ 请求
```json
{
    "jsonrpc":2.0,
    "id":1234,
    "method":"select",
    "params": {
        "token":"xxxx",
        "kind":"deploy",
        "columns":["key1","key2","key3"],
        "filter": {
            "eq":{
                "key1":"abc"
            },
            "ne":{
            },
            "gt":{
                "key2":3
            },
            "ge":{
            },
            "lt":{
            },
            "le":{
            },
            "in":{
            },
            "like":{
            }
        },
        "order":[
            {"Key":"xxx","order":"asc"},
            {"Key":"xxx","order":"desc"}
        ],
        "limiter" : {
            "offset":0,
            "rows":10
        },
        "attach":[
            "attachKey"
        ]
    }
}
```

+ 请求格式说明
```text
    "method":"seletc" ,是 指请求的的目标是查询一种资源
    "kind" :"deploy"  ,是 指资源的类型是　deploy　,
    "columns": 是指请求的列, 必须不为空.
    "filter" 是 指查询时设定的条件, 类似 SQl 的 where　子句. select 预设了 8种关系. 分别为 eq(等于),ne(不等于),gt(大于),ge(大于等于),lt(小于),le(小于等于),in(sqL in),like (sql like)
    每种关系之间 是 and 关系,不支持 or 关系. 只支持在单个字段上设定条件，不支持在组合字段上设定条件.
    比如需要查询　　
    name!="张山" and age >=18 的信息
    "filter":{
        "ne":{
            "name":"张三"
        }
        "ge":{
            "age":18
        }
    }

    查询注册时间在　19-08-10 和　19-08-20 之间的.

    "filter": {
        "ge": {
            "RegistTime":"19-08-10"
        },
        "le":{
            "RegistTime":"19-08-20"
        }
    }
    并不是所有资源字段都支持 filter ,具体会在资源描述文档中指出.
    order 是指 select 操作时的 排序顺序,与 sql语意相同
    limter 是指查询时的分页条件
    attach 是指查询资源时附加选项,附加选项的具体含义在资源文档中定义
```

+ 响应
```json
{
    "jsonrpc" : 2.0,
    "id" : 1234,
    "result" : {
        "attach":{
            "attachKey":123
        },
        "data":[
            {
                "name":"张山",
                "age":18,
                "RegistyTIme":"19-08-11"
            },
            {
                "name":"李四",
                "age":30,
                "RegistyTIme":"19-08-12"
            }
        ]
    }
}
```

+ 响应格式说明
```text
    如果请求中携带了 attach.attachKey 字段, 且该字段要求有响应值.则响应结果中有该字段,否则就没有. 
```



### update
> update 表示对资源的某个对象做一次变更
>
+ 请求
```json
{
    "jsonrpc":2.0,
    "id":1234,
    "method":"update",
    "params":{
        "kind":"deploy",
        "token":"",
        "metadata":{"name":"李四"},
        "target": {
            "age":4,
            "name":"张三"
        }
    }
}
```
+ 请求格式说明
```text
    变更操作由后端保证幂等
    "metadata"是指更新操作的必填项
    "target" 指更新字段以及其目标值
    示例请求等价于 SQL : update deploy set age=4 where name='李四'
    update 支持更新多个字段值
```

+ 响应
```json
{
    "jsonrpc" : 2.0,
    "id" : 1234,
    "result" : 0
}
```


### patch
> patch 表示对资源的某个对象做一次变更

+ 请求
```json
{
    "jsonrpc":2.0,
    "id":1234,
    "method":"patch",
    "params":{
        "kind":"deploy",
        "token":"",
        "metadata":{"name":"李四"},
        "target": {
            "age":4
        }
    }
}
```
+ 请求格式说明
```text
    "metadata"是指更新操作的必填项
    "target" 指更新字段以及其目标值
    示例请求等价于 SQL : update deploy set age=4 where name='李四'
    一次 patch 请求只支持更改单个对象的单个值
```

+ 响应
```json
{
    "jsonrpc" : 2.0,
    "id" : 1234,
    "result" : 0
}
```


+ 响应格式说明
```text
    一般情况下，patch 资源　成功后，会返回　"result" : 0 ,没有其他信息.
    如果出错，则参考　jsonrpc 错误返回情况
```

### delete

> delete 用来表示删除某种资源

+ 请求
```json
{
    "jsonrpc":2.0,
    "id":1234,
    "method":"delete",
    "params": {
        "kind":"deploy",
        "token":"",
        "metadata":{"name":"李四"}
    }
}
```

+ 请求格式说明
```text
    "metadata"  表示　delete　操作需要填写的必填项
    示例请求 等价于 SQl delete from deploy where name='李四'
```

+ 响应

```json
{
    "jsonrpc" : 2.0,
    "id" : 1234,
    "result" : 0
}
```

+ 响应格式说明
```text
    一般情况下，delete 资源　成功后，会返回　"result" : 0 , 没有其他信息.
    如果出错，则参考　jsonrpc 错误返回情况
```

### do

> do 某些操作不太容易抽象成资源操作时启用,
+ 请求
```json
{
    "jsonrpc":2.0,
    "id":1234,
    "method":"do",
    "params": {
        "kind":"deploy",
        "token":"",
        "action":{
          "key": "DeleteBatch",
          "value":[1,2,3,4]
        }
    }
}
```

+ 请求格式说明

```text
    "action" 是 do 操作的必选 参数
    "action.key " 是 action 的名字,具体在资源描述文件中定义
    "action.value" 是 action操作需要的参数,随 action.key 的变化而变化,具体在资源描述文件中定义
```

> 响应

```text
   do 操作没有固定的响应格式,具体在资源描述文件中定义.如果错误,错误响应说明
```


## 错误响应

+ 当后端处理API请求发现错误或其他情况时,会停止处理并返回如下格式

```json
{

    "jsonrpc" : "2.0",
    "id" : "1234",
    "error":{
        "code": -1,
        "message":"描述",
        "data": "暂未使用"
    }
}
```

```text
    code 表示错误类型, 取值 有 -1, -2
    -1 表示请求出现错误,具体错误在 message中描述.
    -2 表示服务端认为请求参数存在很大风险,需要用户对用户警告并,并希望用户二次确认请求.警告信息在 message中描述
       二次请求需要在请求中添加 params.confirmation:true 字段
```