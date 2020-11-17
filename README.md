#  关于

- K8S 是为了将 tars 部署在 K8S平台上而做的适应性改造项目,项目还处于早期开发阶段,当前能够预览效果，但很多细节需要完善.

- K8SFramework 改造的原则:
    1. 改造仅限于 framework 程序（tarsregistry ,tarsnode , tarsnotify , tarsconfig）
    2. framework 程序对业务服务暴露接口保持兼容

- 基于以上原则,现有的 tars 服务如果没有特别的依赖情况,基本可以无缝迁移到 K8SFramework平台

- K8SFramework 暂时不能部署 DCache 服务,但最终会支持


# K8SFramework 的基础组成

+ tarsregistry
    > k8s tarsRegistry 与 tars 框架中的　tarsRegistry　服务提供相同的功能. 在实现上是通过 k8s list-watch  机制来感知集群内的各个服务Pod状态.以及通过 k8s patch 机制来变更服务Pod状态

+ tarsweb
    > k8s tarsweb 是对外提供的管理界面, tars-k8s 框架使用者或管理员通过 tarssweb 提供的界面来管理或运维集群

+ tarsadmin
    > tarsadmin 与 tars 中的　tarsAdminRegistry 服务提供类似功能. 是 tarsweb 操作 k8s,读写 tars_db 的桥梁.
             
+ tarsagent
  > tarsagent 提供tars 节点宿主机的网络栈查询及周边辅助功能，包括日志清理等后续扩展.

+ tarsnotify ,tarsproprety ,tarsstat, tarsqueryproperty ,tarsquerystat,tarslog
    > 以上服务与 tars 框架中的服务同名服务提供相同的功能
    
+ tarsnode
    > tarsnode 程序经过轻量化改造，集成在每一个业务服务镜像中. 作为业务服务的守护进程，以及命令转发桥梁. 每个业务服务镜像中，只能有一个业务服务程序

+ tarsimage
    > 此服务程序是 tars-k8s  框架新增的服务，用于提供镜像生成功能.

+ 删除了 tarspatch , tarsAdminRegistry


# K8SFramework 安装

## 编译代码

请首先编译tars-cpp.

然后编译服务代码.


```bash
    mkdir build && cd build
    cmake ..
    make -j4
    make install
```

因为 基础镜像采用 debian:stretch-slim(debian9 精简版),为避免出现兼容性问题,建议使用 debian 9 环境编译

## 执行安装 install

- 本目录存放部署的模板文件

- 安装前提:
   1. 需要一个可用的 k8s 集群 ,且部署机能通过 kubectl 操作该集群
   2. 需要三个 mysql database ,分别作为 框架基础数据库(_TARS_DB)，服务监控数据库(_TARS_STAT_)，特性监控数据库(_TARS_PROPERTY_)
   3. 需要可用的 docker 镜像仓库，以及合法的用户名，密码
   
- 安装步骤:
```
cd /usr/local/tars/cpp/deploy

#准备环境
sh k8s-pre-install.sh

#执行k8s-install.sh
./k8s-install.sh DOCKER_REGISTRY_URL DOCKER_REGISTRY_USER DOCKER_REGISTRY_PASSWORD DB_tars_HOST DB_tars_PORT DB_tars_USER DB_tars_PASSWORD NODEIP";

#等待执行完毕, 查看 tars-tarsweb service 的 nodeport 端口
kubectl get service -n tars

# 通过 nodeport 端口，可以访问 tars-tarsweb
```

# K8SFramework 使用

## K8SFramework 的使用方式变化
    
    虽然在改造中尽力避免，但仍然有些使用方式发生改变:

+ 寻址结果变化
    >  在 tars 中， 服务名寻址得到的结果都是 ip:port 列表, 在 tars-k8s 中，服务名寻址得到的结果可能是 host:port 列表，也可能是 ip:port 列表.

+ 节点配置的变化
    > 在 tars 中 , 可以针对每个部署节点设定差异化的配置内容 ，在 tars-k8s 中，改为依赖 pod 域名来定制差异化配置

+ 持久存储
    > 因为 k8s 数据的持久存储需要外部资源配合 ,改造过程中暂时忽略此问题

+ 暂停 set 功能
    > set 功能后续完成

+ 新增节点亲和性调度
    > tars-k8s 利用 k8s 本身的特性,可以限定业务服务 pod 只能调度到指定的 k8s node 上, tars 要完成类似工作只能依靠人为约定

+ 下线服务
    > tars 的下线服务是取消某个节点上的某个服务部署.其他节点上的服务部署不受影响.
    > tars-k8s 的下线服务是指在整个 tars-k8s 内删除此服务所有数据.

+ 新增节点
    > tars 的新增节点是新增一台独立主机,在其上部署 tars相关程序,然后注册节点信息
      
    > tars-k8s 新增节点步骤如下:
      
    1. 按 k8s 的使用方式添加 node 
    2. 在新增的 node  上 增加 tars.io/node 标签
    3. 在 tarsadmin程序会自动将 k8s node 那纳为可用节点

# K8SFramework Todo

## 总体
+ 增加 tars-k8s 的统一对外网关, 包括 tcp网关和 http 网关
+ 打通 tars-k8s 部署与 tars 的调用链 (依赖代理服务)
+ 研究 DCache 服务的部署方案
+ 对齐 tars set 功能
+ 增加 持久存储功能

## tarsweb
+ 增加用户登陆,用户管理,用户权限功能
+ 根据 tarsadmin 功能提供情况开发对应界面功能

## tarsadmin
+ 增加更多的 k8s 选项参数 例如节点绑定,调度策略,存储挂载,资源限制等
+ 跟踪 k8s 版本 ,利用新特性

## tarsnode
+ 继续轻量化,在保持功能前提下降低资源占用
+ 根据参数 (tars-web平台配置,并注入到容器运行时) 执行不同的策略

## tarsimage
+ tarsimage 暂时只支持单副本工作,后续考虑增加多副本工作
+ 增加并发任务数量限制功能
+ 使用 docker api 代替 docker 客户端进行镜像生成，上传等工作

## tarregistry
+ 跟踪 k8s 版本,与k8s api 保持兼容

## tarnotify
+ 跟踪 tars 版本同名程序特性,保持接口兼容

## tarsoperator
+ 开发和使用 tarsoperator,用于自动监控和调度 tars 服务