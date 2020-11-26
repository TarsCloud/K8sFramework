/**
 * Tencent is pleased to support the open source community by making Tars available.
 *
 * Copyright (C) 2016THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

const NodeProxy = require("../util/rpcClient/rpcProxy/NodeProxy");
// const NodeProxy = require("../util/rpcClient/tarsProxy/NodeProxy");
const { nodePrx, RPCClientPrx } = require('../util/rpcClient');

const NodeService = {};

NodeService.startServer = async (application, serverName, podIp) => {
    podIp = podIp.split(',') || []
    let rets = [];
    for (var i = 0, len = podIp.length; i < len; i++) {
        let target = podIp[i];
        let ret = {};
        try {
            let nodePrx = await RPCClientPrx(NodeProxy, 'tars', 'Node', `tars.tarsnode.NodeObj@tcp -h ${target} -p 19385 -t 3000`)
            ret = await nodePrx.startServer(application, serverName);
        } catch (e) {
            ret = {
                __return: -1,
                result: e
            }
        }

        rets.push({
            application: application,
            server_name: serverName,
            target,
            ret_code: ret.__return,
            err_msg: ret.result
        });
    }
    return rets;
}

NodeService.stopServer = async (application, serverName, podIp) => {
    podIp = podIp.split(',') || []
    let rets = [];
    for (var i = 0, len = podIp.length; i < len; i++) {
        let target = podIp[i];
        let ret = {};
        try {
            let nodePrx = await RPCClientPrx(NodeProxy, 'tars', 'Node', `tars.tarsnode.NodeObj@tcp -h ${target} -p 19385 -t 3000`)
            ret = await nodePrx.stopServer(application, serverName);
        } catch (e) {
            ret = {
                __return: -1,
                result: e
            }
        }

        rets.push({
            application: application,
            server_name: serverName,
            target,
            ret_code: ret.__return,
            err_msg: ret.result
        });
    }
    return rets;
}

NodeService.restartServer = async (application, serverName, podIp) => {
    podIp = podIp.split(',') || []
    let rets = [];
    for (var i = 0, len = podIp.length; i < len; i++) {
        let target = podIp[i];
        let ret = {};
        try {
            let nodePrx = await RPCClientPrx(NodeProxy, 'tars', 'Node', `tars.tarsnode.NodeObj@tcp -h ${target} -p 19385 -t 3000`)
            ret = await nodePrx.restartServer(application, serverName);
        } catch (e) {
            ret = {
                __return: -1,
                result: e
            }
        }

        rets.push({
            application: application,
            server_name: serverName,
            target,
            ret_code: ret.__return,
            err_msg: ret.result
        });
    }
    return rets;
}

NodeService.doCommand = async (application, serverName, podIp, command) => {
    podIp = podIp.split(',') || []
    let rets = [];
    for (var i = 0, len = podIp.length; i < len; i++) {
        let target = podIp[i];
        let ret = {};
        try {
            let nodePrx = await RPCClientPrx(NodeProxy, 'tars', 'Node', `tars.tarsnode.NodeObj@tcp -h ${target} -p 19385 -t 3000`)
            ret = await nodePrx.notifyServer(application, serverName, command);
        } catch (e) {
            ret = {
                __return: -1,
                result: e
            }
        }

        rets.push({
            application: application,
            server_name: serverName,
            target,
            ret_code: ret.__return,
            err_msg: ret.result,
            ret,
        });
    }
    return rets;
}

module.exports = NodeService;