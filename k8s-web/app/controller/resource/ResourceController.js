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

const ResourceService = require('../../service/resource/ResourceService');
const path = require('path');
// const util = require('../../tools/util');
// const send = require('koa-send');
var fs = require('fs');


const ResourceController = {};

ResourceController.getTarsNode = async(ctx) => {
	// console.log('getTarsNode', ctx);

	let tgzPath = path.join(__dirname, '../../../files/tarsnode.tgz');
	let exists = fs.existsSync(tgzPath);
	if(!exists) {
		ctx.body = "#!/bin/bash \n echo 'not tarsnode.tgz exists'";
	}

	ctx.paramsObj.ip = ctx.paramsObj.ip || ctx.req.headers['x-forwarded-for'] || ctx.req.connection.remoteAddress || ctx.req.socket.remoteAddress || ctx.req.connection.socket.remoteAddress;

	//都是从web过来的请求, 用web host替换安装node脚本
	ctx.body = await ResourceService.getTarsNode(ctx.origin || ctx.request.origin, ctx.paramsObj);
}

module.exports = ResourceController;