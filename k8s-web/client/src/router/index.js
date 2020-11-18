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

import Vue from 'vue';
import Router from 'vue-router';

// 重写路由的push方法
const routerPush = VueRouter.prototype.push
VueRouter.prototype.push = function push(location){
  return routerPush.call(this, location).catch(error => error)
}

// 登录管理
import Login from '@/pages/login/index';

// 服务管理
import Server from '@/pages/server/index';
import ServerManage from '@/pages/server/manage';
import HistoryManage from '@/pages/server/history';
import ServerPublish from '@/pages/server/publish';
import ServerConfig from '@/pages/server/config';
import ServerServerMonitor from '@/pages/server/monitor-server';
import ServerPropertyMonitor from '@/pages/server/monitor-property';
import userManage from '@/pages/server/user-manage';
import InterfaceDebuger from '@/pages/server/interface-debuger';

// 运维管理
import Operation from '@/pages/operation/index';
import OperationDeploy from '@/pages/operation/deploy';
import OperationApproval from '@/pages/operation/approval';
import OperationHistory from '@/pages/operation/history';
import OperationUndeploy from '@/pages/operation/undeploy';
import OperationExpand from '@/pages/operation/expand';
import OperationTemplates from '@/pages/operation/templates';
import OperationMaster from '@/pages/operation/master';
import OperationImage from '@/pages/operation/image';
import OperationStore from '@/pages/operation/store';

// 基础管理
import Base from '@/pages/base/index';
import BaseAffinity from '@/pages/base/affinity';
import BaseApplication from '@/pages/base/application';
import BaseBusiness from '@/pages/base/business';
import BaseNode from '@/pages/base/node';

// 域名管理
import Gateway from '@/pages/gateway/index';
import GatewayStation from '@/pages/gateway/station'
import GatewayUpstream from '@/pages/gateway/upstream'
import GatewayBwlist from '@/pages/gateway/bwlist'
import VueRouter from 'vue-router';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
    },
    {
      path: '/server',
      name: 'Server',
      component: Server,
      children: [
        {
          path: ':treeid/manage',
          component: ServerManage,
        },
        {
          path: ':treeid/history',
          component: HistoryManage,
        },
        {
          path: ':treeid/publish',
          component: ServerPublish,
        },
        {
          path: ':treeid/config',
          component: ServerConfig,
        },
        {
          path: ':treeid/server-monitor',
          component: ServerServerMonitor,
        },
        {
          path: ':treeid/property-monitor',
          component: ServerPropertyMonitor,
        },
        {
          path: ':treeid/interface-debuger',
          component: InterfaceDebuger,
        },
        {
          path: ':treeid/user-manage',
          component: userManage,
        },
      ],
    },
    {
      path: '/operation',
      name: 'Operation',
      component: Operation,
      redirect: '/operation/deploy',
      children: [
        {
          path: 'deploy',
          component: OperationDeploy,
        },
        {
          path: 'approval',
          component: OperationApproval,
        },
        {
          path: 'history',
          component: OperationHistory,
        },
        {
          path: 'undeploy',
          component: OperationUndeploy,
        },
        {
          path: 'expand',
          component: OperationExpand,
        },
        {
          path: 'templates',
          component: OperationTemplates,
        },
        {
          path: 'master',
          component: OperationMaster,
        },
        {
          path: 'image',
          component: OperationImage,
        },
        {
          path: 'store',
          component: OperationStore,
        },
      ],
    },
    {
      path: '/base',
      name: 'Base',
      component: Base,
      redirect: '/base/application',
      children: [
        {
          path: 'affinity',
          component: BaseAffinity,
        },
        {
          path: 'application',
          component: BaseApplication,
        },
        {
          path: 'business',
          component: BaseBusiness,
        },
        {
          path: 'node',
          component: BaseNode,
        },
      ],
    },
    {
      path: '/gateway',
      name: 'Gateway',
      component: Gateway,
      redirect: '/gateway/station',
      children: [
        {
          path: 'station',
          component: GatewayStation,
        },
        {
          path: 'upstream',
          component: GatewayUpstream,
        },
        {
          path: 'bwlist',
          component: GatewayBwlist,
        },
      ]
    },
    {
      path: '*',
      redirect: '/login',
    },
  ],
  scrollBehavior (to, from, savedPosition) {
    return {x: 0, y: 0}
  }
});
