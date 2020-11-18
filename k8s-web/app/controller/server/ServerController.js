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

const logger = require('../../logger');
const NodeService = require('../../service/node/NodeService');
const util = require('../../tools/util');

const serverConfStruct = {
	id: '',
	application: '',
	server_name: '',
	node_name: '',
	server_type: '',
	enable_set: {
		formatter: (value) => {
			return value == 'Y' ? true : false;
		}
	},
	set_name: '',
	set_area: '',
	set_group: '',
	setting_state: '',
	present_state: '',
	bak_flag: {
		formatter: (value) => {
			return value == 0 ? false : true;
		}
	},
	template_name: '',
	profile: '',
	async_thread_num: '',
	base_path: '',
	exe_path: '',
	start_script_path: '',
	stop_script_path: '',
	monitor_script_path: '',
	patch_time: {formatter: util.formatTimeStamp},
	patch_version: "",
	process_id: '',
	posttime: {formatter: util.formatTimeStamp}
};

const ServerController = {};

ServerController.sendCommand = async (ctx) => {
	let params = ctx.paramsObj;
	let application = params.serverApp
	let serverName = params.serverName
	let podIp = params.podIp
	let command = params.command

	try{
		let ret
		switch (`${command}`) {
			case 'StartServer':
				ret = await NodeService.startServer(application, serverName, podIp);
				break;
			case 'StopServer':
				ret = await NodeService.stopServer(application, serverName, podIp);
				break;
			case 'RestartServer':
				ret = await NodeService.restartServer(application, serverName, podIp);
				break;
			default:
				ret = await NodeService.doCommand(application, serverName, podIp, command);
				break;
		}
		ctx.makeResObj(200, '', ret);
	}catch(e) {
		logger.error('[sendCommand]', e, ctx);
		ctx.makeErrResObj();
	}
}

module.exports = ServerController;