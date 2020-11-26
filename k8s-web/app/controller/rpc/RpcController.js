const logger = require('../../logger')
const request = require('request')
const axios = require('axios');
const fs = require('fs')
// const LdapServer = require('../../../proxy/LdapServer')

const { pkgUploadPath, rpcDomain, uploadDomain } = global.WebConf

const ajax = ({ options = {} }) => {
    return new Promise((resolve) => {
        let result = {
            ret: -1,
            msg: 'error'
        }

        if(Object.keys(options).length === 0){
            result.msg = '[rpc_ajax]: argument invalid'
            return resolve(result)
        }
        
        try{
            console.log('[rpc_ajax]', options)

            request(options, (error, response, body) => {
                if (!error && response.statusCode === 200) {
                    // swagger服务端框架默认客户端使用默认值初始化。即，变量的默认值(result=0，""，false等)将不会组包发送
                    if (Object.keys(body).length === 0) {
                        result.ret = 0
                        result.data = 0
                        result.msg = 'success'
                    } else {
                        let rsp = body && JSON.parse(body) || body
                        if(rsp && (rsp.result === 0 || rsp.result)){
                            result.ret = 0
                            result.data = rsp.result
                            result.msg = 'success'
                        }else{
                            result.ret = 0
                            result.data = rsp
                            result.msg = 'success'
                        }
                    }
                    resolve(result)
                }else{
                    if (body.hasOwnProperty("message")) {
                        result.msg = body.message
                    } else {
                        result.msg = body
                    }
                    resolve(result)
                }
            })
        } catch (e) {
            logger.error(JSON.stringify(e))
            result.msg = e.message
            resolve(result)
        }
    })
}

const upload = ({ formData = {} }) => {
    return new Promise((resolve) => {
        let result = {
            ret: -1,
            msg: 'error'
        }

        if(Object.keys(formData).length === 0){
            result.msg = '[rpc_upload]: argument invalid'
            return resolve(result)
        }
        
        try{
            const options = {
                url: `${uploadDomain}/upload`,
                method: 'POST',
                formData: formData,
                timeout: 3000,
            }
            request(options, (error, response, body) => {
                if (!error && response.statusCode == 200) {
                    let rsp = body && JSON.parse(body) || body
                    if(rsp && rsp.UploadStatus === 'done'){
                        result.ret = 0
                        result.data = rsp.FileName
                        result.msg = 'success'
                    }else{
                        result.msg = rsp.UploadMessage
                    }
                    resolve(result)
                } else {
                    if (error && error.hasOwnProperty("message")) {
                        result.msg = error.message
                    } else {
                        result.msg = response.body.toString()
                    }
                    resolve(result)
                    logger.warn('[rpc_upload]', error, data)
                }
            })
        } catch (e) {
            result.msg = e.message
            resolve(result)
            logger.error('[rpc_upload]', e, data)
        }
    })
}

const build = ({ data = {} }) => {
    return new Promise((resolve) => {
        let result = {
            ret: -1,
            msg: 'error'
        }

        if(Object.keys(data).length === 0){
            result.msg = '[rpc_build]: argument invalid'
            return resolve(result)
        }
        
        try{
            const options = {
                url: `${uploadDomain}/build`,
                method: 'POST',
                header: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
                timeout: 3000,
            }
            request(options, (error, response, body) => {
                if (!error && response.statusCode == 200) {
                    let rsp = body && JSON.parse(body) || body
                    if(rsp && (rsp.BuildStatus === 'done' || rsp.BuildStatus === 'working')){
                        result.ret = 0
                        result.data = {
                            BuildId: rsp.BuildId || '',
                            BuildStatus: rsp.BuildStatus,
                            BuildImage: rsp.BuildImage || '',
                            BuildMessage: rsp.BuildMessage || '',
                        }
                    }
                    result.msg = rsp.BuildMessage
                    resolve(result)
                }else{
                    if (error && error.hasOwnProperty("message")) {
                        result.msg = error.message
                    } else {
                        result.msg = response.body.toString()
                    }
                    resolve(result)
                    logger.warn('[rpc_build]', error, data)
                }
            })
        } catch (e) {
            result.msg = e.message
            resolve(result)
            logger.error('[rpc_build]', e, data)
        }
    })
}

const buildStatus = ({ BuildId = '' }) => {
    return new Promise((resolve) => {
        let result = {
            ret: -1,
            msg: 'error'
        }

        if(!BuildId){
            result.msg = '[rpc_buildStatus]: argument invalid'
            return resolve(result)
        }
        
        try{
            const options = {
                url: `${uploadDomain}/buildStatus`,
                method: 'POST',
                header: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    BuildId,
                }),
                timeout: 3000,
            }
            request(options, (error, response, body) => {
                if (!error && response.statusCode == 200) {
                    let rsp = body && JSON.parse(body) || body
                    if(rsp && rsp.BuildStatus){
                        result.ret = 0
                        result.data = {
                            BuildId: rsp.BuildId || '',
                            BuildStatus: rsp.BuildStatus,
                            BuildMessage: rsp.BuildMessage || '',
                        }
                        result.msg = rsp.BuildMessage
                        resolve(result)
                    }else{
                        result.msg = rsp.BuildMessage
                        resolve(result)
                    }
                }else{
                    if (error && error.hasOwnProperty("message")) {
                        result.msg = error.message
                    } else {
                        result.msg = response.body.toString()
                    }
                    resolve(result)
                    logger.warn('[rpc_buildStatus]', error, BuildId)
                }
            })
        } catch (e) {
            result.msg = e.message
            resolve(result)
            logger.error('[rpc_buildStatus]', e, data)
        }
    })
}

const treeNode = ({ data = [] }) => {
    let result = []
    if(data && data.length > 0){
        data.forEach((item, index) => {
            let obj = {}
            if(item.BusinessName){
                obj = {
                    id: item.BusinessName,
                    name: item.BusinessShow,
                    type: '0',
                    is_parent: true,
                    open: true,
                    pid: 'root',
                    children: [],
                }

                if(item.App && item.App.length > 0){
                    let tempArr = []
                    item.App.forEach(AppItem => {
                        let obj = {
                            children: [],
                            id: `${AppItem.AppName}`,
                            name: AppItem.AppName,
                            type: '1',
                            is_parent: false,
                            open: 'false',
                            pid: item.BusinessName,
                        }
                        if(AppItem.Server && AppItem.Server.length > 0){
                            AppItem.Server.forEach(serverItem => {
                                obj.children.push({
                                    id: `${AppItem.AppName}.${serverItem.ServerName}`,
                                    name: serverItem.ServerName,
                                    type: '2',
                                    is_parent: false,
                                    open: 'false',
                                    pid: AppItem.AppName,
                                    children: [],
                                })
                            })
                        }
                        tempArr.push(obj)
                    })
                    obj.children = tempArr
                }

                result.push(obj)
            }else{
                if(item.App && item.App.length > 0){
                    item.App.forEach(AppItem => {
                        if(AppItem.AppName){
                            obj = {
                                id: AppItem.AppName,
                                name: AppItem.AppName,
                                type: '1',
                                is_parent: true,
                                open: true,
                                pid: 'root',
                                children: [],
                            }
                            if(AppItem.Server && AppItem.Server.length > 0){
                                AppItem.Server.forEach(serverItem => {
                                    obj.children.push({
                                        id: `${AppItem.AppName}.${serverItem.ServerName}`,
                                        name: serverItem.ServerName,
                                        type: '2',
                                        is_parent: false,
                                        open: 'false',
                                        pid: AppItem.AppName,
                                        children: [],
                                    })
                                })
                            }
                            result.push(obj)
                        }
                    })
                }
            }
        })
    }
    return result
}

module.exports = {
    /**
     * 登录、退出
     * @param  {String}  Account             帐号
     * @param  {String}  Password            密码
     * @param  {String}  Token               登录签名
     */
    async AccessToken(ctx) {
        let { Account = '', Password = '', Token = '' } = ctx.paramsObj

        try {
            const method = Token ? 'delete' : 'create'
            const metadata = Token ? { Token } : { Account, Password }

            const params = {
                kind: 'AccessToken',
                metadata,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AccessToken]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 目录树
     * @param  {String}  Token                登录签名
     */
    async ServerTree(ctx) {
        let { Token = '' } = ctx.paramsObj

        try {
            const options = {
                url: `${rpcDomain}/trees`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                let treeData = treeNode({
                    data: result.data,
                })
                ctx.makeResObj(200, result.msg, treeData)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerTree]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 默认值
     * @param  {String}  Token                登录签名
     */
    async DefaultValue(ctx) {
        let { Token = '' } = ctx.paramsObj

        try {
            const options = {
                url: `${rpcDomain}/defaults`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_DefaultValue]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 发布包上传
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务
     * @param  {String}  ServerType           服务类型(tars_java,tars_cpp,tars_node_pkg,tars_node,tars_node8,tars_node10)
     */
    async uploadPatchPackage(ctx) {
        const that = module.exports

        let { Token = '', ServerId = '', ServerType = '' } = ctx.paramsObj
        let file = ctx.req.files[0]

        try {
            if(!file){
                return ctx.makeResObj(500, 'no files')
            }
            let baseUploadPath = pkgUploadPath.path
            let uploadTgzName = `${baseUploadPath}/${file.filename}.gz`
            fs.renameSync(`${baseUploadPath}/${file.filename}`, uploadTgzName)
            let fileBuffer = fs.createReadStream(uploadTgzName)
            let uploadData = await upload({
                formData: {
                    uploadfile: fileBuffer
                }
            })
            if(uploadData && uploadData.ret === 0 && uploadData.data){
                // 查询服务信息
                let ServerApp = ServerId.substring(0, ServerId.indexOf('.')) || '',
                    ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length) || ''

                if(!ServerApp || !ServerName){
                    return ctx.makeResObj(500, 'ServerId有误')
                }

                // 生成镜像
                let buildData = await build({
                    data: {
                        ServerApp,
                        ServerName,
                        ServerType,
                        ServerTGZ: uploadData.data,
                    }
                })

                // 查询镜像状态
                if(buildData && buildData.ret === 0 && buildData.data){
                    ctx.makeResObj(200, buildData.msg, buildData.data)
                }else{
                    ctx.makeResObj(500, buildData.msg)
                }
            }else{
                ctx.makeResObj(500, uploadData.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerDeployCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 发布包上传状态
     * @param  {String}  Token                登录签名
     * @param  {Number}  BuildId              编译ID
     */
    async uploadPatchStatus(ctx) {
        const that = module.exports

        let { Token = '', ServerId = '', BuildId = 0 } = ctx.paramsObj

        try {
            if(ServerId){
                let serverData = await that.ServerInfoSelect(ctx)
                if(serverData && serverData.ret === 0 && serverData.data) {
                    ctx.paramsObj.ServerId = serverData.data.Data[0].ServerId
                }
            }

            let buildStatusData = await buildStatus({ BuildId })

            if(buildStatusData && buildStatusData.ret === 0 && buildStatusData.data){
                if(buildStatusData.data.BuildStatus === 'done'){
                    that.ServicePoolCreate(ctx)
                }
                ctx.makeResObj(200, buildStatusData.msg, buildStatusData.data)
            }else{
                ctx.makeResObj(500, buildStatusData.msg)
            }
        } catch (e) {
            logger.error('[rpc_uploadPatchStatus]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务名
     * @param  {String}  ServerApp            应用名
     * @param  {String}  ServerName           服务名
     */
    async ServerInfoSelect(ctx) {
        return new Promise(async (resolve) => {
            let { Token = '',
                ServerId = '',
                page = 1,
                isAll = false,
            } = ctx.paramsObj

            let filter = {
                eq: {},
            }

            let pageIndex = Math.floor(page) || 1
            let pageSize = 10

            let limiter = {}
            if(!isAll){
                limiter = {
                    offset: (pageIndex - 1) * pageSize,
                    rows: pageSize,
                }
            }

            if(ServerId){
                if(ServerId.indexOf('.') === -1){
                    filter.eq.ServerApp = ServerId
                }else{
                    filter.eq.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                    filter.eq.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
                }
            }

            try {
                const options = {
                    url: `${rpcDomain}/servers?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
                }

                let result = await ajax({ options })
                resolve(result)
            } catch (e) {
                logger.error('[rpc_ServerInfoSelect]', e, ctx)
                resolve({
                    ret: -500,
                    msg: e.message,
                })
            }
        })
    },
    /**
     * 服务列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务名
     * @param  {String}  ServerApp            应用名
     * @param  {String}  ServerName           服务名
     * @param  {String}  ServerType           类型
     * @param  {String}  ServerMark           备注
     */
    async ServerSelect(ctx) {
        let { Token = '', ServerId = '', ServerApp = '', ServerName = '',
            page = 1, isAll = false, } = ctx.paramsObj
        let filter = {
            eq: {},
            like: {},
        }

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        if(ServerId){
            if(ServerId.indexOf('.') === -1){
                filter.eq.ServerApp = ServerId
            }else{
                filter.eq.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                filter.eq.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }
        }

        if(ServerApp){
            filter.like.ServerApp = `.*${ServerApp}.*`
        }
        if(ServerName){
            filter.like.ServerName = `.*${ServerName}.*`
        }

        try {
            const options = {
                url: `${rpcDomain}/servers?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServerId             服务ID
     * @param  {String}  ServerType           类型
     * @param  {String}  ServerMark           备注
     */
    async ServerUpdate(ctx) {
        let { Token = '', ServerId = '',
            ServerType = '', ServerMark = '',
        } = ctx.paramsObj
        const that = module.exports

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                ServerId = serverData.data.Data[0].ServerId
            }
        }

        try {
            const metadata = {
                ServerId,
            }
            const target = {
                ServerType,
                ServerMark,
            }

            const params = {
                metadata,
                target,
            }
            const options = {
                url: `${rpcDomain}/servers`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务下线
     * @param  {String}  Token                登录签名
     * @param  {Array}   ServerId             服务ID
     */
    async ServerUndeploy(ctx) {
        let { Token = '', ServerId = [] } = ctx.paramsObj

        try {
            const metadata = {
                ServerId,
            }

            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/servers`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerUndeploy]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务部署创建
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerApp            应用名
     * @param  {String}  ServerName           服务名
     * @param  {String}  ServerMark           服务部署备注
     * @param  {Object}  ServerK8S            服务K8S
     * @param  {Object}  ServerServant        服务详情
     * @param  {Object}  ServerOption         服务配置
     */
    async ServerDeployCreate(ctx) {
        let { Token = '',
            ServerApp = '', ServerName = '', ServerMark = '', ServerK8S = {}, ServerServant = {}, ServerOption = {}
        } = ctx.paramsObj

        if(ServerServant){
            for(let item in ServerServant){
                delete ServerServant[item].HostPort
                ServerServant[item].Port = Math.floor(ServerServant[item].Port) || 0
                ServerServant[item].Threads = Math.floor(ServerServant[item].Threads) || 0
                ServerServant[item].Connections = Math.floor(ServerServant[item].Connections) || 0
                ServerServant[item].Capacity = Math.floor(ServerServant[item].Capacity) || 0
                ServerServant[item].Timeout = Math.floor(ServerServant[item].Timeout) || 0
            }
        }

        if(ServerK8S){
            ServerK8S.Replicas = ServerK8S.Replicas ? Math.floor(ServerK8S.Replicas) : 1
        }

        try {
            const metadata = {
                ServerApp,
                ServerName,
                ServerMark,
                ServerK8S,
                ServerServant,
                ServerOption,
            }

            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/deploys`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerDeployCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务部署列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerApp            应用名
     * @param  {String}  ServerName           服务名
     */
    async ServerDeploySelect(ctx) {
        let { Token = '',
            ServerApp = '', ServerName = '',
            page = 1, isAll = false,
        } = ctx.paramsObj

        let filter = {
            like: {},
        }

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        if(ServerApp){
            filter.like.ServerApp = `.*${ServerApp}.*`
        }
        if(ServerName){
            filter.like.ServerName = `.*${ServerName}.*`
        }

        try {
            const order = [
                { column: 'RequestTime', order: 'desc' },
            ]

            const options = {
                url: `${rpcDomain}/deploys?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerDeploySelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务部署更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  DeployId             服务ID
     * @param  {String}  ServerMark           服务备注
     * @param  {Object}  ServerK8S            服务K8S
     * @param  {Array}   ServerServant        服务详情
     * @param  {Object}  ServerOption         服务配置
     */
    async ServerDeployUpdate(ctx) {
        let { Token = '', DeployId = '',
            ServerK8S = {}, ServerServant = {}, ServerOption = {},
        } = ctx.paramsObj

        try {
            const metadata = {
                DeployId,
            }
            const target = {
                ServerK8S,
                ServerServant,
                ServerOption,
            }
            const params = {
                metadata,
                target,
            }

            const options = {
                url: `${rpcDomain}/deploys`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerDeployUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务部署删除
     * @param  {String}  Token                登录签名
     * @param  {Number}  DeployId             服务ID
     */
    async ServerDeployDelete(ctx) {
        let { Token = '', DeployId = '' } = ctx.paramsObj

        try {
            const metadata = {
                DeployId: Math.floor(DeployId),
            }
            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/deploys`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerDeployDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务审批
     * @param  {String}  Token                登录签名
     * @param  {Number}  DeployId             部署ID
     * @param  {Boolean} ApprovalResult       提交时间
     * @param  {String}  ApprovalMark         审批备注
     */
    async ServerApprovalCreate(ctx) {
        let { Token = '', DeployId = 0,
            ApprovalResult = false, ApprovalMark = '',
        } = ctx.paramsObj

        if(ApprovalResult && `${ApprovalResult}` === 'true') {
            ApprovalResult = true
        }else if(ApprovalResult && `${ApprovalResult}` === 'false'){
            ApprovalResult = false
        }
        
        try {
            const metadata = {
                DeployId,
                ApprovalResult,
                ApprovalMark,
            }

            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/approvals`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerApprovalCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务审批列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务ID
     * @param  {String}  ServerApp            应用名
     * @param  {String}  ServerName           服务名
     */
    async ServerApprovalSelect(ctx) {
        const that = module.exports

        let { Token = '',
            ServerId = '', ServerApp = '', ServerName = '',
            page = 1, isAll = false,
        } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            eq: {},
            like: {},
        }

        if(ServerId){
            filter.eq.ServerId = ServerId
        }

        if(ServerApp){
            filter.like.ServerApp = `.*${ServerApp}.*`
        }
        if(ServerName){
            filter.like.ServerName = `.*${ServerName}.*`
        }
        
        try {
            const order = [
                { column: 'RequestTime', order: 'desc' },
            ]

            const options = {
                url: `${rpcDomain}/approvals?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerApprovalSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性创建
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             节点名
     * @param  {String}  AppServer            应用
     */
    async AffinityCreate(ctx) {
        let { Token = '',
            NodeName = '', AppServer = '',
        } = ctx.paramsObj
        
        try {
            const method = 'create'
            const metadata = {
                NodeName,
                AppServer,
            }
            const params = {
                kind: 'Affinity',
                token: Token,
                metadata,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性列表
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             节点名
     * @param  {String}  AppServer            应用
     */
    async AffinitySelect(ctx) {
        let { Token = '', page = 1, isAll = false,
            NodeName = '', AppServer = '',
        } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            like: {},
        }

        if(NodeName){
            filter.like.NodeName = `.*${NodeName}.*`
        }

        if(AppServer){
            filter.like.AppServer = `.*${AppServer}.*`
        }
        
        try {
            const options = {
                url: `${rpcDomain}/affinities?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data.Data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinitySelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性聚合Node
     * @param  {String}  Token                登录签名
     */
    async AffinityListByNode(ctx) {
        let { Token = '',
            NodeName = []
        } = ctx.paramsObj
        
        try {
            const options = {
                url: `${rpcDomain}/affinities/nodes?NodeName=${JSON.stringify(NodeName)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityListByNode]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性聚合Server
     * @param  {String}  Token                登录签名
     */
    async AffinityListByServer(ctx) {
        let { Token = '',
            ServerApp = []
        } = ctx.paramsObj
        
        try {
            const options = {
                url: `${rpcDomain}/affinities/servers?ServerApp=${JSON.stringify(ServerApp)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityListByServer]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性批量增加server
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             节点名
     * @param  {Array}   ServerApp            应用
     */
    async AffinityAddServer(ctx) {
        let { Token = '', NodeName = '', ServerApp = [] } = ctx.paramsObj
        
        try {
            const metadata = {
                NodeName,
                ServerApp,
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/affinities/nodes`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityAddServer]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性批量增加Node
     * @param  {String}  Token                登录签名
     * @param  {Array}   NodeName             节点名
     * @param  {String}  ServerApp            应用
     */
    async AffinityAddNode(ctx) {
        let { Token = '', NodeName = [], ServerApp = '' } = ctx.paramsObj
        
        try {
            const metadata = {
                NodeName,
                ServerApp,
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/affinities/servers`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityAddNode]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性批量删除server
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             节点名
     * @param  {Array}   ServerApp            应用
     */
    async AffinityDeleteServer(ctx) {
        let { Token = '', NodeName = '', ServerApp = [] } = ctx.paramsObj
        
        try {
            const metadata = {
                NodeName,
                ServerApp,
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/affinities/nodes`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityDeleteServer]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性批量删除Node
     * @param  {String}  Token                登录签名
     * @param  {Array}   NodeName             节点名
     * @param  {String}  ServerApp            应用
     */
    async AffinityDeleteNode(ctx) {
        let { Token = '', NodeName = [], ServerApp = '' } = ctx.paramsObj
        
        try {
            const metadata = {
                NodeName,
                ServerApp,
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/affinities/servers`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityDeleteNode]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 亲和性删除
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             节点名
     * @param  {String}  ServerApp            应用
     */
    async AffinityDelete(ctx) {
        let { Token = '', NodeName = '', ServerApp = '' } = ctx.paramsObj
        
        try {
            const method = 'delete'
            const metadata = {
                NodeName,
                ServerApp,
            }
            const params = {
                kind: 'Affinity',
                token: Token,
                metadata,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_AffinityDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 应用创建
     * @param  {String}  Token                登录签名
     * @param  {String}  AppName              应用名
     * @param  {String}  AppMark              应用备注
     * @param  {String}  BusinessName         业务名
     */
    async ApplicationCreate(ctx) {
        let { Token = '',
            AppName = '', AppMark = '', BusinessName = '',
        } = ctx.paramsObj
        
        try {
            const metadata = {
                AppName,
                AppMark,
                BusinessName,
            }
            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/applications`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ApplicationCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 应用列表
     * @param  {String}  Token                登录签名
     * @param  {String}  AppName              应用名
     * @param  {String}  AppMark              应用备注
     * @param  {String}  BusinessName         业务名
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  CreatePerson         创建人
     */
    async ApplicationSelect(ctx) {
        let { Token = '', page = 1, isAll = false,
            AppName = '', BusinessName = '',
        } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            like: {},
        }

        if(AppName){
            filter.like.AppName = `.*${AppName}.*`
        }

        if(BusinessName){
            filter.like.BusinessName = `.*${BusinessName}.*`
        }
        
        try {
            url = `${rpcDomain}/applications?Filter=${JSON.stringify(filter)}`
            const options = {
                url: url,
            }

            if(!isAll){
                options.url = `${url}&Limiter=${JSON.stringify(limiter)}`
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ApplicationSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 应用更新
     * @param  {String}  Token                登录签名
     * @param  {String}  AppName              应用名
     * @param  {String}  AppMark              应用备注
     * @param  {String}  BusinessName         业务名
     */
    async ApplicationUpdate(ctx) {
        let { Token = '', AppName = '', AppMark = '', BusinessName = '' } = ctx.paramsObj
        
        try {
            const metadata = {
                AppName,
            }
            let target = {
                AppMark,
                BusinessName,
            }

            const params = {
                metadata,
                target,
            }

            const options = {
                url: `${rpcDomain}/applications`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ApplicationUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 应用删除
     * @param  {String}  Token                登录签名
     * @param  {String}  AppName              应用名
     */
    async ApplicationDelete(ctx) {
        let { Token = '', AppName = '' } = ctx.paramsObj
        
        try {
            const metadata = {
                AppName,
            }
            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/applications`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ApplicationDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务创建
     * @param  {String}  Token                登录签名
     * @param  {String}  BusinessName         名称
     * @param  {String}  BusinessShow         显示内容
     * @param  {String}  BusinessMark         备注
     * @param  {Number}  BusinessOrder        排序
     */
    async BusinessCreate(ctx) {
        let { Token = '',
            BusinessName = '', BusinessShow = '', BusinessMark = '', BusinessOrder = 0,
        } = ctx.paramsObj

        BusinessOrder = Math.floor(BusinessOrder) || 0
        
        try {
            const metadata = {
                BusinessName,
                BusinessShow,
                BusinessMark,
                BusinessOrder,
            }
            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/businesses`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务列表
     * @param  {String}  Token                登录签名
     * @param  {String}  BusinessName         名称
     * @param  {String}  BusinessShow         显示内容
     * @param  {String}  BusinessMark         备注
     * @param  {Number}  BusinessOrder        排序
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  CreatePerson         创建人
     */
    async BusinessSelect(ctx) {
        let { Token = '', page = 1, isAll = false,
            BusinessName = '', BusinessShow = '', BusinessMark = '',
        } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            like: {},
        }

        if(BusinessName){
            filter.like.BusinessName = `.*${BusinessName}.*`
        }

        if(BusinessShow){
            filter.like.BusinessShow = `.*${BusinessShow}.*`
        }

        if(BusinessMark){
            filter.like.BusinessMark = `.*${BusinessMark}.*`
        }
        
        try {
            const options = {
                url: `${rpcDomain}/businesses?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务更新
     * @param  {String}  Token                登录签名
     * @param  {String}  BusinessName         名称
     * @param  {String}  BusinessShow         显示内容
     * @param  {String}  BusinessMark         备注
     * @param  {Number}  BusinessOrder        排序
     */
    async BusinessUpdate(ctx) {
        let { Token = '',
            BusinessName = '', BusinessShow = '', BusinessMark = '', BusinessOrder = 0,
        } = ctx.paramsObj

        BusinessOrder = Math.floor(BusinessOrder) || 0
        
        try {
            const metadata = {
                BusinessName,
            }
            let target = {
                BusinessShow,
                BusinessMark,
                BusinessOrder,
            }

            const params = {
                metadata,
                target,
            }
            const options = {
                url: `${rpcDomain}/businesses`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务删除
     * @param  {String}  Token                登录签名
     * @param  {String}  BusinessName         名称
     */
    async BusinessDelete(ctx) {
        let { Token = '', BusinessName = '' } = ctx.paramsObj
        
        try {
            const metadata = {
                BusinessName,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/businesses`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务添加APP
     * @param  {String}  Token                登录签名
     * @param  {String}  BusinessName         业务名称
     * @param  {Array}   AppName              应用名称
     */
    async BusinessAddApp(ctx) {
        let { Token = '',
            BusinessName = '', AppName = [],
        } = ctx.paramsObj
        
        try {
            const metadata = {
                BusinessName,
                AppName,
            }
            const params = {
                metadata
            }
            const options = {
                url: `${rpcDomain}/businesses/apps`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessAddApp]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务聚合App
     * @param  {String}  Token                登录签名
     * @param  {Array}  BusinessName         业务名称
     */
    async BusinessListByApp(ctx) {
        let { Token = '',
            BusinessName = [],
        } = ctx.paramsObj
        
        try {
            const options = {
                url: `${rpcDomain}/businesses/apps?BusinessName=${JSON.stringify(BusinessName)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data[0])
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessListByApp]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 业务删除APP
     * @param  {String}  Token                登录签名
     * @param  {String}  BusinessName         业务名称
     * @param  {Array}   AppName              应用名称
     */
    async BusinessDeleteApp(ctx) {
        let { Token = '',
            BusinessName = '', AppName = [],
        } = ctx.paramsObj
        
        try {
            const metadata = {
                BusinessName,
                AppName,
            }
            const params = {
                metadata
            }
            const options = {
                url: `${rpcDomain}/businesses/apps`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_BusinessDeleteApp]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 节点创建
     * @param  {String}  Token                登录签名
     * @param  {String}   NodeName             名称
     * @param  {String}  NodeMark             备注
     */
    async NodeCreate(ctx) {
        let { Token = '',
            NodeName = [], NodeMark = '',
        } = ctx.paramsObj
        
        try {
            const method = 'create'
            const metadata = {
                NodeName,
                NodeMark,
            }
            const params = {
                kind: 'Node',
                token: Token,
                metadata,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NodeCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 节点列表
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             名称
     * @param  {String}  NodeMark             备注
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  CreatePerson         创建人
     */
    async NodeSelect(ctx) {
        let { Token = '', page = 1, isAll = false,
            NodeName = '', NodeMark = '',
        } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            like: {},
        }

        if(NodeName){
            filter.like.NodeName = `.*${NodeName}.*`
        }

        if(NodeMark){
            filter.like.NodeMark = `.*${NodeMark}.*`
        }
        
        try {
            const options = {
                url: `${rpcDomain}/nodes?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NodeSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 节点更新
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             名称
     */
    async NodeUpdate(ctx) {
        let { Token = '',
            NodeName = '', NodePublic = '',
        } = ctx.paramsObj
        
        try {
            const method = 'update'
            const metadata = {
                NodeName,
            }
            let target = {
                NodePublic,
            }

            const params = {
                kind: 'Node',
                token: Token,
                metadata,
                target,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NodeUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 节点删除
     * @param  {String}  Token                登录签名
     * @param  {String}  NodeName             名称
     */
    async NodeDelete(ctx) {
        let { Token = '', NodeName = '' } = ctx.paramsObj
        
        try {
            const method = 'delete'
            const metadata = {
                NodeName,
            }
            const params = {
                kind: 'Node',
                token: Token,
                metadata,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NodeDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 节点列表
     * @param  {String}  Token                登录签名
     */
    async NodeList(ctx) {
        let { Token = '',

        } = ctx.paramsObj

        try {
            const options = {
                url: `${rpcDomain}/nodes/cluster`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NodeList]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 节点更新
     * @param  {String}  Token                登录签名
     * @param  {Array}   NodeName             名称
     */
    async NodeAdd(ctx) {
        let { Token = '',
            NodeName = []
        } = ctx.paramsObj
        
        try {
            const method = 'do'
            const action = {
                key: 'AddNode',
                value: {
                    NodeName
                },
            }

            const params = {
                kind: 'Node',
                token: Token,
                action,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NodeAdd]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 启用公用节点
     * @param  {String}  Token                登录签名
     * @param  {Array}   NodeName             名称
     */
    async NodeStartPublic(ctx) {
        let { Token = '',
            NodeName = []
        } = ctx.paramsObj
        
        try {
            const metadata = {
                NodeName
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/nodes/public`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_SetPublicNode]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 停用公用节点
     * @param  {String}  Token                登录签名
     * @param  {Array}   NodeName             名称
     */
    async NodeStopPublic(ctx) {
        let { Token = '',
            NodeName = []
        } = ctx.paramsObj
        
        try {
            const metadata = {
                NodeName
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/nodes/public`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_DeletePublicNode]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 模版创建
     * @param  {String}  Token                登录签名
     * @param  {String}  TemplateName         模板名
     * @param  {String}  TemplateParent       父模板名
     * @param  {String}  TemplateContent      模板内容
     * @param  {String}  CreateMark           模板备注
     */
    async TemplateCreate(ctx) {
        let { Token = '',
            TemplateName = '', TemplateParent = '', TemplateContent = '', CreateMark = '',
        } = ctx.paramsObj
        
        try {
            const metadata = {
                TemplateName,
                TemplateParent,
                TemplateContent,
                CreateMark,
            }
            const params = {
                metadata,
            }

            const options = {
                url: `${rpcDomain}/templates`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_TemplateCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 模版列表
     * @param  {String}  Token                登录签名
     * @param  {String}  TemplateId           模版ID
     * @param  {String}  TemplateName         模版名
     * @param  {String}  TemplateParent       父模版名
     * @param  {String}  TemplateContent      模版内容
     * @param  {String}  CreatePerson         创建人
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  CreateMark           创建备注
     * @param  {String}  UpdatePerson         更新人
     * @param  {String}  UpdateTime           更新时间
     * @param  {String}  UpdateMark           更新备注
     */
    async TemplateSelect(ctx) {
        let { Token = '', page = 1, isAll = false, TemplateName = '' } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            like: {},
        }

        if(TemplateName){
            filter.like.TemplateName = `.*${Math.floor(TemplateName)}.*`
        }
        
        try {
            url = `${rpcDomain}/templates?Filter=${JSON.stringify(filter)}`
            const options = {
                url: url,
            }

            if(!isAll){
                options.url = `${url}&Limiter=${JSON.stringify(limiter)}`
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_TemplateSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务模版更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  TemplateId           模版ID
     * @param  {String}  TemplateParent       父模板名
     * @param  {String}  TemplateContent      模板内容
     * @param  {String}  CreateMark           备注
     */
    async TemplateUpdate(ctx) {
        let { Token = '', TemplateId = '', TemplateParent = '', TemplateContent = '', CreateMark = '' } = ctx.paramsObj

        try {
            const metadata = {
                TemplateId,
            }
            let target = {
                TemplateParent,
                TemplateContent,
                CreateMark,
            }

            const params = {
                metadata,
                target,
            }
            const options = {
                url: `${rpcDomain}/templates`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_TemplateUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务模版删除
     * @param  {String}  Token                登录签名
     * @param  {Number}  TemplateId           模版ID
     */
    async TemplateDelete(ctx) {
        let { Token = '', TemplateId = '' } = ctx.paramsObj

        try {
            const metadata = {
                TemplateId,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/templates`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_TemplateDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务版本创建
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServerId             服务ID
     * @param  {Number}  ServiceVersion       版本
     * @param  {String}  ServiceImage         镜像地址
     * @param  {Object}  ImageDetail          镜像详情
     * @param  {String}  CreateMark           版本备注
     */
    async ServicePoolCreate(ctx) {
        const that = module.exports

        let { Token = '',
            ServerId = '', ServerType = '', ServiceImage = '', ImageDetail = {}, CreateMark = '',
        } = ctx.paramsObj


        try {
            const metadata = {
                ServerId,
                ServerType,
                ServiceImage,
                ImageDetail,
                ServiceMark: CreateMark,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/releases`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServicePoolCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务版本列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务名
     */
    async ServicePoolSelect(ctx) {
        const that = module.exports

        let { Token = '', ServerId = '', page = 1, isAll = false } = ctx.paramsObj

        let filter = {
            eq: {},
        }

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        if(ServerId){
            if(ServerId.indexOf('.') === -1){
                filter.eq.ServerApp = ServerId
            }else{
                filter.eq.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                filter.eq.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }
        }
        
        try {
            const order = [
                { column: 'ServiceVersion', order: 'desc' },
            ]
            const options = {
                url: `${rpcDomain}/releases?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServicePoolSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务版本更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServiceId            版本ID
     * @param  {Number}  EnableMark           备注
     */
    async ServicePoolUpdate(ctx) {
        const that = module.exports
        let { Token = '', ServerId, ServiceId = '', Replicas = 1, EnableMark = '' } = ctx.paramsObj

        Replicas = Math.floor(Replicas) || 1
        
        try {
            const metadata = {
                ServerId,
                ServiceId,
                Replicas,
                EnableMark,
            }

            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/releases`,
                method: "PUT",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServicePoolUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务版本删除
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServiceId            服务版本ID
     */
    async ServicePoolDelete(ctx) {
        let { Token = '', ServiceId = 0 } = ctx.paramsObj
        
        try {
            const method = 'delete'
            const metadata = {
                ServiceId,
            }
            const params = {
                kind: 'ServicePool',
                token: Token,
                metadata,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServicePoolDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务版本号启用
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServiceId            版本ID
     * @param  {String}  EnableMark           备注
     */
    async ServicePoolState(ctx) {
        let { Token = '',
            ServiceId = 0, EnableMark = '',
        } = ctx.paramsObj

        ServiceId = Math.floor(ServiceId) || 0
        
        try {
            const method = 'do'
            const action = {
                key: 'EnableService',
                value: {
                    ServiceId,
                    EnableMark,
                },
            }

            const params = {
                kind: 'ServicePool',
                token: Token,
                action,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServicePoolState]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务版本号生成
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServerId             服务ID
     */
    async ServicePoolVersionID(ctx) {
        const that = module.exports
        let { Token = '', ServerId = 0 } = ctx.paramsObj

        ServerId = Math.floor(ServerId) || 0

        try {
            const method = 'make'
            const metadata = {
                ServerId,
            }
            const params = {
                kind: 'ServicePool',
                token: Token,
                metadata,
                makeTarget: ['ServiceVersion']
            }
            return await ajax({ method, params })
        } catch (e) {
            logger.error('[rpc_ServicePoolVersionID]', e, ctx)
            return {
                ret: -500,
                msg: e.message,
            }
        }
    },
    /**
     * 服务servant创建
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务ID
     * @param  {Object}   ServerServant        servant列表
     */
    async ServerAdapterCreate(ctx) {
        const that = module.exports
        let { Token = '', ServerId = '', ServerServant = {} } = ctx.paramsObj

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                ServerId = serverData.data.Data[0].ServerId
            }
        }
        
        if(ServerServant){
            for(let item in ServerServant){
                delete ServerServant[item].HostPort
                ServerServant[item].Port = Math.floor(ServerServant[item].Port) || 0
                ServerServant[item].Threads = Math.floor(ServerServant[item].Threads) || 0
                ServerServant[item].Connections = Math.floor(ServerServant[item].Connections) || 0
                ServerServant[item].Capacity = Math.floor(ServerServant[item].Capacity) || 0
                ServerServant[item].Timeout = Math.floor(ServerServant[item].Timeout) || 0
            }
        }

        try {
            const metadata = {
                ServerId,
                Servant: ServerServant,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/servers/servants`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerAdapterCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务servant列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务ID
     */
    async ServerAdapterSelect(ctx) {
        const that = module.exports
        let { Token = '', ServerId = '', AdapterId = 0, page = 1, isAll = false, isTars = '', isTcp = false } = ctx.paramsObj

        let filter = {
            eq: {},
        }

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                filter.eq.ServerId = serverData.data.Data[0].ServerId
            }
        }
        if(AdapterId){
            filter.eq.AdapterId = Math.floor(AdapterId)
        }

        if(isTars){
            if(isTars === 'true'){
                filter.eq.IsTars = true
            }else if(isTars === 'false'){
                filter.eq.IsTars = false
            }
        }

        if(isTcp){
            filter.eq.IsTcp = true
        }
        
        try {
            const options = {
                url: `${rpcDomain}/servers/servants?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerAdapterSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务servant更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  AdapterId            servantID
     * @param  {String}  Name                 名称
     * @param  {Number}  Port                 端口(1-30000)
     * @param  {Number}  Threads              线程数(1-30)
     * @param  {Number}  Connections          连接数(1024-100000)
     * @param  {Number}  Capacity             队列长度(100-100000)
     * @param  {Number}  Timeout              超时时间(0-100000)
     * @param  {Boolean} IsTars               是否TARS(true, false)
     * @param  {Boolean} IsTcp                是否TCP(true, false)
     */
    async ServerAdapterUpdate(ctx) {
        let { Token = '', AdapterId = '', ServerServant = [], Confirmation = false } = ctx.paramsObj

        ServerServant.forEach(item => {
            item.Port = Math.floor(item.Port)
            item.Threads = Math.floor(item.Threads)
            item.Connections = Math.floor(item.Connections)
            item.Capacity = Math.floor(item.Capacity)
            item.Timeout = Math.floor(item.Timeout)
        })
        
        try {
            const metadata = {
                AdapterId,
            }
            const target = {
                Name: ServerServant[0].Name,
                Threads: ServerServant[0].Threads,
                Connections: ServerServant[0].Connections,
                Port: ServerServant[0].Port,
                Capacity: ServerServant[0].Capacity,
                Timeout: ServerServant[0].Timeout,
                IsTars: ServerServant[0].IsTars,
                IsTcp: ServerServant[0].IsTcp,
            }
            const params = {
                metadata,
                target,
                Confirmation,
            }
            const options = {
                url: `${rpcDomain}/servers/servants`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerAdapterUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务servant删除
     * @param  {String}  Token                登录签名
     * @param  {Number}  AdapterId            servantID
     */
    async ServerAdapterDelete(ctx) {
        let { Token = '', AdapterId = '' } = ctx.paramsObj

        try {
            const metadata = {
                AdapterId,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/servers/servants`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerAdapterDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务编辑列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务
     */
    async ServerOptionSelect(ctx) {
        const that = module.exports
        let { Token = '', ServerId = '', page = 1, isAll = false } = ctx.paramsObj

        let filter = {
            eq: {},
        }

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                filter.eq.ServerId = serverData.data.Data[0].ServerId
            }
        }
        
        try {
            const options = {
                url: `${rpcDomain}/servers/options?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data.Data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerOptionSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务编辑列表
     * @param  {String}  Token                 登录签名
     * @param  {Number}  ServerId              服务ID
	 * @param  {Number}  ServerImportant       服务等级
     * @param  {String}  ServerTemplate        服务模版
	 * @param  {String}  ServerProfile         私有模版
     * @param  {String}  StartScript           开始脚本路径
	 * @param  {String}  StopScript            停止脚本路径
	 * @param  {String}  MonitorScript         监控脚本路径
	 * @param  {Number}  AsyncThread           异步线程数
	 * @param  {Boolean} RemoteLogEnable       远程日志启用( true, false )
	 * @param  {Number}  RemoteLogReserveTime  远程日志保留时间
	 * @param  {Number}  RemoteLogCompressTIme 远程日志压缩时间
     */
    async ServerOptionUpdate(ctx) {
        const that = module.exports
        let { Token = '', ServerId = '',
            ServerImportant = 0, ServerTemplate = '', ServerProfile = '',
            StartScript = '', StopScript = '',
            MonitorScript = '', AsyncThread = '',
            Confirmation = false
        } = ctx.paramsObj

        if (AsyncThread === '') {
            AsyncThread = 3
        } else {
            AsyncThread = parseInt(AsyncThread)
        }

        try {
            const metadata = {
                ServerId,
            }
            let target = {
                ServerTemplate,
                ServerProfile,
                StartScript,
                StopScript,
                MonitorScript,
                AsyncThread,
                ServerImportant,
            }

            const params = {
                metadata,
                target,
                Confirmation,
            }
            const options = {
                url: `${rpcDomain}/servers/options`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerOptionUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 查看服务模版内容
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务ID
     */
    async ServerOptionTemplate(ctx) {
        const that = module.exports
        let { Token = '', ServerId = '' } = ctx.paramsObj

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                ServerId = serverData.data.Data[0].ServerId
            }
        }

        try {
            const options = {
                url: `${rpcDomain}/servers/options/templates?ServerId=${ServerId}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerOptionTemplate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置创建
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId            应用名或服务名(ServerApp、ServerApp.ServerName)
     * @param  {String}  ConfigName           配置名
     * @param  {String}  ConfigContent        配置内容
     * @param  {String}  ConfigMark           创建备注
     */
    async ServerConfigCreate(ctx) {
        const that = module.exports

        let { Token = '', PodSeq = '', ServerId = '', ConfigName = '', ConfigContent = '', ConfigMark = '' } = ctx.paramsObj

        try {
            const metadata = {
                ConfigName,
                ConfigContent,
                ConfigMark,
            }

            if(ServerId){
                if (ServerId.indexOf('.') === -1) {
                    metadata.ServerApp = ServerId
                } else {
                    metadata.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                    metadata.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
                }
            }

            if (PodSeq) {
                metadata.PodSeq = PodSeq
            }

            const params = {
                kind: 'ServerConfig',
                token: Token,
                metadata,
            }
            const options = {
                url: `${rpcDomain}/configs`,
                method: "POST",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigCreate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置列表
     * @param  {String}  Token                登录签名
     * @param  {String}  AppServer            应用名或服务名(ServerApp、ServerApp.ServerName)
     * @param  {Number}  ConfigId             配置ID
     * @param  {String}  ConfigName           配置名
     * @param  {Number}  ConfigVersion        配置版本
     * @param  {String}  ConfigContent        配置内容
     * @param  {String}  CreatePerson         创建人
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  ConfigMark           创建备注
     */
    async ServerConfigSelect(ctx) {
        const that = module.exports

        let { Token = '', page = 1, ServerId = '', ConfigName = '', isAll = false } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            eq: {},
        }

        if(ServerId){
            if (ServerId.indexOf('.') === -1) {
                filter.eq.ServerApp = ServerId
            } else {
                filter.eq.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                filter.eq.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }

            // select configName的节点配置
            if(ConfigName !== ''){
                filter.eq.ConfigName = ConfigName
            }
        }
        
        try {
            const options = {
                url: `${rpcDomain}/configs?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  ConfigId             配置ID
     * @param  {String}  ConfigContent        配置内容
     */
    async ServerConfigUpdate(ctx) {
        let { Token = '', ConfigId = '',  ConfigMark = '', ConfigContent = '' } = ctx.paramsObj

        try {
            const metadata = {
                ConfigId,
            }
            let target = {
                ConfigMark,
                ConfigContent,
            }

            const params = {
                metadata,
                target,
            }
            const options = {
                url: `${rpcDomain}/configs`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置删除
     * @param  {String}  Token                登录签名
     * @param  {Number}  ConfigId             配置ID
     */
    async ServerConfigDelete(ctx) {
        let { Token = '', ConfigId = '' } = ctx.paramsObj

        try {
            const metadata = {
                ConfigId,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/configs`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 预览节点配置与主配置合并后的内容
     * @param  {String}  Token                登录签名
     * @param  {String}  AppServer            应用名或服务名(ServerApp、ServerApp.ServerName)
     * @param  {String}  ConfigName           配置名
     * @param  {Number}  PodSeq               节点序号
     */
    async ServerConfigContent(ctx) {
        let { Token = '', ServerId = '', ConfigName = '', PodSeq = '', ServerApp = '', ServerName = '' } = ctx.paramsObj


        try {
            if (ServerId.indexOf('.') === -1) {
                ServerApp = ServerId
            } else {
                ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }

            const options = {
                url: `${rpcDomain}/configs/join?ServerApp=${ServerApp}&ServerName=${ServerName}&ConfigName=${ConfigName}&PodSeq=${PodSeq}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigContent]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置历史列表
     * @param  {String}  Token                登录签名
     * @param  {String}  AppServer            应用名或服务名(ServerApp、ServerApp.ServerName)
     * @param  {Number}  HistoryId            配置历史ID
     * @param  {Number}  ConfigId             配置ID
     * @param  {String}  ConfigName           配置名
     * @param  {Number}  ConfigVersion        配置版本
     * @param  {String}  ConfigContent        配置内容
     * @param  {String}  CreatePerson         创建人
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  ConfigMark           创建备注
     */
    async ServerConfigHistroySelect(ctx) {
        let { Token = '', page = 1, ConfigId = '', ConfigName = '' } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {
            offset: (pageIndex - 1) * pageSize,
            rows: pageSize,
        }

        let filter = {
            eq: {},
        }

        if(ConfigId){
            filter.eq.ConfigId = ConfigId
        }
        if(ConfigName){
            filter.eq.ConfigName = ConfigName
        }
        
        try {
            const order = [
                { column: 'ConfigVersion', order: 'desc' },
            ]
            const options = {
                url: `${rpcDomain}/configs/versions?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigHistorySelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置历史删除
     * @param  {String}  Token                登录签名
     * @param  {Number}  HistoryId            配置历史ID
     */
    async ServerConfigHistroyDelete(ctx) {
        let { Token = '', HistoryId = 0 } = ctx.paramsObj

        HistoryId = Math.floor(HistoryId) || 0
        
        try {
            const metadata = {
                HistoryId,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/configs/versions`,
                method: "DELETE",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigHistoryDelete]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务配置回滚
     * @param  {String}  Token                登录签名
     * @param  {Number}  HistoryId            历史ID
     */
    async ServerConfigHistoryBack(ctx) {
        let { Token = '', HistoryId = 0 } = ctx.paramsObj

        HistoryId = Math.floor(HistoryId) || 0
        
        try {
            const metadata = {
                HistoryId,
            }
            const params = {
                metadata,
            }
            const options = {
                url: `${rpcDomain}/configs`,
                method: "PUT",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerConfigHistoryBack]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * K8S列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             配置ID
     * @param  {Object}  ServerK8S            配置
     * @param  {Number}  ServerK8S.Replicas   副本数量
     */
    async ServerK8SSelect(ctx) {
        const that = module.exports

        let { Token = '', page = 1, ServerId = '', isAll = false } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            eq: {},
        }

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                filter.eq.ServerId = serverData.data.Data[0].ServerId
            }
        }
        
        try {
            const options = {
                url: `${rpcDomain}/servers/k8s?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data.Data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerK8SSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * ServerK8S更新
     * @param  {String}  Token                登录签名
     * @param  {Number}  ServerId             配置ID
     * @param  {Number}  Replicas             副本
     * @param  {String}  NodeSelector         节点
     */
    async ServerK8SUpdate(ctx) {
        let { Token = '', ServerId = '', Replicas = 0, HostIpc = false, HostNetwork = false, HostPort = {}, NodeSelector = '' } = ctx.paramsObj

        Replicas = Math.floor(Replicas) || 0
        
        try {
            const metadata = {
                ServerId,
            }
            let target = {
                Replicas,
                HostIpc,
                HostNetwork,
                HostPort,
                NodeSelector,
            }

            const params = {
                metadata,
                target,
            }
            const options = {
                url: `${rpcDomain}/servers/k8s`,
                method: "PATCH",
                json: true,
                headers: {
                    "content-type": "application/json",
                },
                body: params
            }
            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServerK8SUpdate]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    async ServerK8SGenerateHostPort(ctx) {
        let NodeList = ctx.paramsObj.NodeList.split(",")
        let NodePort = ctx.paramsObj.NodePort

        if (NodeList.length <= 0) {
            ctx.makeErrResObj()
            return
        }

        // 最多重试3次
        let count = 0
        while (count < 3) {
            try {
                // 首先发出第一个请求，查找候选的空闲端口
                let response = await axios.get(`${rpcDomain}/hostPorts?NodeName=${NodeList[0]}&Port=${NodePort}`)

                let result = response.data
                if (response.status === 500) {
                    ctx.makeResObj(500, result)
                    return

                }
                if (!result.available) {
                    ctx.makeResObj(500, `${result.port} is not in use now.`, result)
                    return
                }

                // 如果存在多个节点，并发请求候选节点是否符合要求
                if (NodeList.length > 1) {
                    let NodeReuqest = []
                    for (let i = 1; i < NodeList.length-1; i++) {
                        NodeList.forEach((node, index) => {
                            NodeReuqest.push(
                                axios.get(`${rpcDomain}/hostPorts?NodeName=${NodeList[i]}&Port=${result.port}`)
                            )
                        })
                    }

                    let available  = true
                    for (let i = 0; i < NodeReuqest.length; i++) {
                        response    = await NodeReuqest[i]
                        if (response.status === 500) {
                            available = false
                        } else {
                            available = response.data.available && available
                        }
                    }
                    if (!available) {
                        logger.info(`[rpc_ServerK8SUpdate] port:${result.port} is in use, try again.`)
                        count++;
                        continue
                    }
                }

                ctx.makeResObj(200, '', result)
                return
            } catch (e) {
                logger.error('[rpc_GenerateHostPort]', e, ctx)
                ctx.makeErrResObj()
                return
            }
        }
    },
    /**
     * Pod存活列表
     * @param  {String}  Token                登录签名
     * @param  {String}  PodId                Pod ID
     * @param  {String}  PodName              Pod名称
     * @param  {String}  PodIp                Pod IP
     * @param  {String}  ServerId             服务ID
     * @param  {String}  ServerApp            服务应用
     * @param  {String}  ServerName           服务名称
     * @param  {Number}  ServiceVersion       服务版本
     * @param  {String}  SettingState         设置状态
     * @param  {String}  PresentState         当前状态
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  UpdateTime           更新时间
     */
    async PodAliveSelect(ctx) {
        const that = module.exports

        let { Token = '', page = 1, size = 10, PodId = '', ServerId = '', isAll = false } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = Math.floor(size)

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            eq: {},
        }
        
        if(PodId){
            filter.eq.PodId = PodId
        }
        if(ServerId){
            if(ServerId.indexOf('.') === -1){
                filter.eq.ServerApp = ServerId
            }else{
                filter.eq.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                filter.eq.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }
        }
        
        try {
            const order = [
                { column: 'PodName', order: 'asc' },
            ]
            const options = {
                url: `${rpcDomain}/servers/alivePods?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_PodAliveSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * Pod历史列表
     * @param  {String}  Token                登录签名
     * @param  {String}  PodId                Pod ID
     * @param  {String}  PodName              Pod名称
     * @param  {String}  PodIp                Pod IP
     * @param  {Number}  ServerId             服务ID
     * @param  {String}  ServerApp            服务应用
     * @param  {String}  ServerName           服务名称
     * @param  {String}  CreateTime           创建时间
     * @param  {String}  DeleteTime           删除时间
     */
    async PodPerishedSelect(ctx) {
        const that = module.exports

        let { Token = '', page = 1, size = 10, PodId = '', ServerId = 0, isAll = false } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = Math.floor(size)

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            eq: {},
        }
        
        if(PodId){
            filter.eq.PodId = PodId
        }
        if(ServerId){
            if(ServerId.indexOf('.') === -1){
                filter.eq.ServerApp = ServerId
            }else{
                filter.eq.ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                filter.eq.ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }
        }
        
        try {
            const order = [
                { column: 'DeleteTime', order: 'desc' },
            ]
            const options = {
                url: `${rpcDomain}/servers/perishedPods?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_PodPerishedSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务启用列表
     * @param  {String}  Token                登录签名
     * @param  {String}  ServerId             服务器ID
     * @param  {String}  ServerApp            服务器应用
     * @param  {String}  ServerName           服务器名称
     * @param  {Number}  ServiceId            服务ID
     * @param  {Number}  ServiceVersion       服务版本 
     * @param  {String}  ServiceImage         服务资源
     * @param  {String}  ImageDetail          资源详情
     * @param  {String}  EnablePerson         启用人
     * @param  {String}  EnableTime           启用时间
     * @param  {String}  EnableMark           启用备注
     * @param  {String}  CreatePerson         创建人
     * @param  {String}  CreateTime           创建时间
     */
    async ServiceEnabledSelect(ctx) {
        const that = module.exports

        let { Token = '', page = 1, ServerId = '', isAll = false } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {}
        if(!isAll){
            limiter = {
                offset: (pageIndex - 1) * pageSize,
                rows: pageSize,
            }
        }

        let filter = {
            eq: {},
        }

        if(ServerId){
            let serverData = await that.ServerInfoSelect(ctx)
            if(serverData && serverData.ret === 0 && serverData.data) {
                filter.eq.ServerId = serverData.data.Data[0].ServerId
            }
        }
        
        try {
            const options = {
                url: `${rpcDomain}/servers/releases?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_ServiceEnabledSelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 服务状态列表
     * @param  {String}  Token                登录签名
     * @param  {Number}  NotifyId             路由ID
     * @param  {Number}  AppServer            应用
     * @param  {Number}  ServerId             服务ID
     * @param  {Number}  ThreadId             线程ID
     * @param  {Number}  Command              
     * @param  {Number}  Result               结果
     * @param  {String}  NotifyTime           时间
     */
    async NotifySelect(ctx) {
        const that = module.exports

        let { Token = '', page = 1, ServerId = '' } = ctx.paramsObj

        let pageIndex = Math.floor(page) || 1
        let pageSize = 10

        let limiter = {
            offset: (pageIndex - 1) * pageSize,
            rows: pageSize,
        }

        let filter = {
            eq: {},
        }
        
        if(ServerId){
            filter.eq.AppServer = ServerId
        }

        let order = [
            { column: 'NotifyId', order: 'desc' },
        ]

        try {
            const options = {
                url: `${rpcDomain}/notifies?Filter=${JSON.stringify(filter)}&Limiter=${JSON.stringify(limiter)}&Order=${JSON.stringify(order)}`,
            }

            let result = await ajax({ options })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_NotifySelect]', e, ctx)
            ctx.makeErrResObj()
        }
    },
    /**
     * 命令
     * @param  {String}  Token                登录签名
     * @param  {String}  command              命令名称 (StartServer, StopServer)
     * @param  {String}  server_ids           服务ID
     */
    async Command(ctx) {
        let { Token = '', command = '', server_ids = [] } = ctx.paramsObj

        try {
            let ServerApp = '', ServerName = ''
            if (ServerId.indexOf('.') === -1) {
                ServerApp = ServerId
            } else {
                ServerApp = ServerId.substring(0, ServerId.indexOf('.'))
                ServerName = ServerId.substring(ServerId.indexOf('.') + 1, ServerId.length)
            }

            if(!ServerApp || !ServerName){
                return ctx.makeResObj(500, 'ServerId有误')
            }

            const method = 'do'
            const action = {
                key: command,
                value: {
                    ServerApp,
                    ServerName,
                    PodIps,
                },
            }

            const params = {
                kind: 'Command',
                token: Token,
                action,
            }
            let result = await ajax({ method, params })
            if(result && result.ret === 0){
                ctx.makeResObj(200, result.msg, result.data)
            }else{
                ctx.makeResObj(500, result.msg)
            }
        } catch (e) {
            logger.error('[rpc_Command]', e, ctx)
            ctx.makeErrResObj()
        }
    },
}
