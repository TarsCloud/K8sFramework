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

const ServerDao = require('../../dao/ServerDao');
const BusinessDao = require('../../dao/BusinessDao');
const BusinessRelationDao = require('../../dao/BusinessRelationDao');
const AuthService = require('../auth/AuthService');

const TreeService = {};

TreeService.getTreeNodes = async (searchKey, uid, applicationList, serverNameList) => {
	let serverList;
	if (await AuthService.hasAdminAuth(uid)) {
		serverList = await ServerDao.getServerConf4Tree(applicationList, serverNameList, false);
	} else {
		let appCond = [];
		let serverCond = [];
		let authList = await AuthService.getAuthListByUid(uid);
		authList.forEach((auth) => {
			let application = auth.application;
			let serverName = auth.serverName;
			if (serverName) {
				serverCond.push(application + '.' + serverName);
			} else {
				appCond.push(application);
			}
		});
		serverList = await ServerDao.getServerConf4Tree(appCond, serverCond, false);
	}

	let treeNodeList = [];
	let treeNodeMap = {};
	let rootNode = [];
	serverList.forEach(function (server) {
		server = server.dataValues;
		let id;
		if (server.enable_set == 'Y') {
			id = '1' + server.application + '.' + '2' + server.set_name + '.' + '3' + server.set_area + '.' + '4' + server.set_group + '.' + '5' + server.server_name;
		} else {
			id = '1' + server.application + '.' + '5' + server.server_name;
		}
		let treeNode = {};
		treeNode.id = id;
		treeNode.name = server.server_name;
		treeNode.pid = id.substring(0, id.lastIndexOf('.'));
		treeNode.is_parent = false;
		treeNode.open = false;
		treeNode.children = [];
		treeNodeList.push(treeNode);
		treeNodeMap[id] = treeNode;

		TreeService.parents(treeNodeMap, treeNode, rootNode);
	});
	// return rootNode;

	let result = await TreeService.parentsBusiness(rootNode, searchKey)
	return result	
};

/**
 * 将应用服务转换成层级数据
 * @param treeList
 */
TreeService.parents = (treeNodeMap, treeNode, rootNodes) => {
	let id = treeNode.pid;

	if (id == 'root') {
		rootNodes.push(treeNode);
		return;
	}

	if (treeNodeMap[id]) {
		treeNodeMap[id].children.push(treeNode);
		return;
	}

	let pid, name;
	let idx = id.lastIndexOf('.');
	if (idx === -1) {
		pid = 'root';
		name = id;
	} else {
		pid = id.substring(0, idx);
		name = id.substring(idx + 1);
	}
	let newTreeNode = {};
	newTreeNode.id = id;
	newTreeNode.name = name.substring(1);
	newTreeNode.pid = pid;
	newTreeNode.is_parent = true;
	newTreeNode.open = false;
	newTreeNode.children = [];
	newTreeNode.children.push(treeNode);

	treeNodeMap[id] = newTreeNode;
	TreeService.parents(treeNodeMap, newTreeNode, rootNodes);

};

/**
 * 添加业务层级数据
 * @param treeList
 */
TreeService.parentsBusiness = async (data, searchKey) => {
	let businessList = await BusinessDao.getList('') || []
	let businessRelationList = await BusinessRelationDao.getList('', '') || []

	// 查询过滤
	if(searchKey){
		businessList = businessList.filter(item => 
			item.f_name.toLowerCase().indexOf(searchKey.toLowerCase()) > -1
			|| item.f_show_name.toLowerCase().indexOf(searchKey.toLowerCase()) > -1
		)
		if(businessList && businessList.length > 0){
			let newData = []
			businessRelationList && businessRelationList.forEach(jitem => {
				businessList && businessList.forEach(kitem => {
					data && data.forEach(item => {
						if(item.name === jitem.f_application_name && jitem.f_business_name === kitem.f_name){
							newData.push(item)
						}
					})
					
				})
			})
			data = newData
		}else{
			data = data.filter(item => item.name.toLowerCase().indexOf(searchKey.toLowerCase()) > -1)
		}
	}

	let result = []
	data && data.forEach((item, index) => {
		let relationIsTrue = false
		businessRelationList && businessRelationList.forEach(jitem => {
			businessList && businessList.forEach(kitem => {
				if(item.name === jitem.f_application_name && jitem.f_business_name === kitem.f_name){
					relationIsTrue = true
					let obj = {
						name: kitem.f_show_name,
						pid: kitem.f_name,
						is_parent: true,
						open: false,
						children: [item],
						order: kitem.f_order,
					}

					let isTrue = false
					let resultIndex = 0
					result.forEach((litem, lindex) => {
						if(litem.pid === obj.pid){
							isTrue = true
							resultIndex = lindex
						}
					})

					if(isTrue){
						result[resultIndex].children.push(obj.children[0])
					}else{
						result.push(obj)
					}
				}
			})
		})
		if(!relationIsTrue && item.name){
			item.order = 0
			result.push(item)
		}
	})
	
	result = result.sort((a, b) => {
		if(b.order === a.order){
			let aName = a.name.toLowerCase()
			let bName = b.name.toLowerCase()
			if(aName < bName){ return -1 }
			if(aName > bName){ return 1 }
			return 0
		}
		return b.order - a.order
	})

	return result
}

module.exports = TreeService;