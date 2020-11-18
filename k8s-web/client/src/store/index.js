export const operationTree = [
    { title: '服务部署', path: '/operation/deploy' },
    { title: '待审批服务', path: '/operation/approval' },
    { title: '已审批服务', path: '/operation/history' },
    { title: '下线服务', path: '/operation/undeploy' },
    { title: '模版管理', path: '/operation/templates' },
    // { title: '主控配置', path: '/operation/master' },
    // { title: '镜像管理', path: '/operation/image' },
    // { title: '仓库配置', path: '/operation/store' },
]

export const baseTree = [
    { title: '业务管理', path: '/base/business' },
    { title: '应用管理', path: '/base/application' },
    { title: '节点管理', path: '/base/node' },
    { title: '能力管理', path: '/base/affinity' },
]

export const gatewayTree = [
    { title: '站点配置', path: '/gateway/station' },
    { title: 'upstream配置', path: '/gateway/upstream' },
    { title: '全局黑名单配置', path: '/gateway/bwlist' },
]