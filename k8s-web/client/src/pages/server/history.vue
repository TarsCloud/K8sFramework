<template>
  <div class="page_server_manage">

    <!-- 服务列表 -->
    <let-table v-if="serverList" :data="serverList" :title="$t('serverList.title.serverList')" :empty-msg="$t('common.nodata')">
      <!-- <let-table-column :title="$t('serverList.table.th.podID')" prop="PodId"></let-table-column> -->
      <let-table-column :title="$t('deployService.form.app')" prop="ServerApp"></let-table-column>
      <let-table-column :title="$t('deployService.form.serviceName')" prop="ServerName"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.podName')" prop="PodName">
        <template slot-scope="scope">
          <let-table-operation @click="gotoLog(scope.row)">{{ scope.row.PodName }}</let-table-operation>
        </template>
      </let-table-column>
      <let-table-column :title="$t('serverList.table.th.podIP')" prop="PodIp"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.ip')" prop="NodeIp"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.version')" prop="ServiceVersion"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.createTime')" prop="CreateTime"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.deleteTime')" prop="DeleteTime"></let-table-column>
    </let-table>

    <div style="overflow:hidden;">
      <let-pagination align="right" style="float:right;"
        :page="pagination.page" @change="gotoPage"
        :total="pagination.total">
      </let-pagination>
      <!-- <div style="float:left;">
        <let-button theme="primary" size="small" @click="configServer">{{$t('operate.update')}}</let-button>
        <let-button theme="primary" size="small" @click="manageServant">{{$t('operate.servant')}}</let-button>
        <let-button theme="primary" size="small" @click="manageK8S">{{$t('operate.k8s')}}</let-button>
      </div> -->
    </div>

    <!-- 服务实时状态 -->
    <let-table v-if="serverNotifyList && showOthers"
      :data="serverNotifyList" :title="$t('serverList.title.serverStatus')" :empty-msg="$t('common.nodata')" ref="serverNotifyListLoading">
      <let-table-column :title="$t('common.time')" prop="notifytime"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.serviceID')" prop="server_id"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.threadID')" prop="thread_id"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.result')" prop="result"></let-table-column>
    </let-table>

    <!-- 编辑服务弹窗 -->
    <let-modal
      v-model="configModal.show"
      :title="$t('serverList.dlg.title.editService')"
      width="800px"
      :footShow="!!(configModal.model && configModal.model.ServerId)"
      @on-confirm="saveConfig"
      @close="closeConfigModal"
      @on-cancel="closeConfigModal">
      <let-form
        v-if="!!(configModal.model && configModal.model.ServerId)"
        ref="configForm" itemWidth="360px" :columns="2" class="two-columns">
        <!-- <let-form-item :label="$t('common.service')">{{configModal.model.server_name}}</let-form-item>
        <let-form-item :label="$t('common.ip')">{{configModal.model.node_name}}</let-form-item>
        <let-form-item :label="$t('serverList.dlg.isBackup')" required>
          <let-radio-group
            size="small"
            v-model="configModal.model.bak_flag"
            :data="[{ value: true, text: $t('common.yes') }, { value: false, text: $t('common.no') }]">
          </let-radio-group>
        </let-form-item> -->
        <let-form-item :label="$t('common.template')" required>
          <let-select
            size="small"
            v-model="configModal.model.ServerTemplate"
            v-if="configModal.model.templates && configModal.model.templates.length"
            required>
            <let-option v-for="t in configModal.model.templates" :key="t.TemplateName" :value="t.TemplateName">{{t.TemplateName}}</let-option>
          </let-select>
          <span v-else>{{configModal.model.ServerTemplate}}</span>
        </let-form-item>
        <let-form-item :label="$t('serverList.dlg.asyncThread')" required>
          <let-input disabled
            size="small"
            v-model="configModal.model.AsyncThread"
            :placeholder="$t('serverList.dlg.placeholder.thread')"
            required
            :pattern="configModal.model.server_type === 'taf.nodejs' ? '^[1-9][0-9]*$' : '^([3-9]|[1-9][0-9]+)$'"
            pattern-tip="$t('serverList.dlg.placeholder.thread')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.dlg.privateTemplate')" labelWidth="150px" itemWidth="724px">
          <let-input disabled
            size="large"
            type="textarea"
            :rows="4"
            v-model="configModal.model.ServerProfile"
          ></let-input>
        </let-form-item>
      </let-form>
      <div v-else class="loading-placeholder" ref="configFormLoading"></div>
    </let-modal>

    <!-- Servant管理弹窗 -->
    <let-modal
      v-model="servantModal.show"
      :title="$t('serverList.table.servant.title')"
      width="1200px"
      :footShow="false"
      @close="closeServantModal">
      <let-button size="small" theme="primary" class="tbm16" @click="configServant()">{{$t('operate.add')}} Servant</let-button>
      <let-table v-if="servantModal.model" :data="servantModal.model" :empty-msg="$t('common.nodata')">
        <let-table-column title="OBJ" prop="Name"></let-table-column>
        <let-table-column :title="$t('deployService.table.th.port')" prop="Port"></let-table-column>
        <let-table-column :title="$t('deployService.form.portType')" width="150px">
          <template slot-scope="props">
            <let-radio v-model="props.row.IsTcp" :label="true" v-if="props.row.IsTcp">TCP</let-radio>
            <let-radio v-model="props.row.IsTcp" :label="false" v-else>UDP</let-radio>
          </template>
        </let-table-column>
        <let-table-column :title="$t('deployService.table.th.threads')" width="80px">
          <template slot-scope="props">{{ props.row.Threads }}</template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.servant.connections')" width="140px">
          <template slot-scope="props">{{ props.row.Connections }}</template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.servant.capacity')" width="140px">
          <template slot-scope="props">{{ props.row.Capacity }}</template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.servant.timeout')" width="140px">
          <template slot-scope="props">{{ props.row.Timeout }}</template>
        </let-table-column>
        <let-table-column :title="$t('operate.operates')" width="90px">
          <template slot-scope="scope">
            <let-table-operation @click="configServant(scope.row.AdapterId)">{{$t('operate.update')}}</let-table-operation>
            <let-table-operation class="danger" @click="deleteServant(scope.row.AdapterId)">{{$t('operate.delete')}}</let-table-operation>
          </template>
        </let-table-column>
      </let-table>
      <div v-else class="loading-placeholder" ref="servantModalLoading"></div>
    </let-modal>

    <!-- Servant新增、编辑弹窗 -->
    <let-modal
      v-model="servantDetailModal.show"
      :title="servantDetailModal.isNew ? `${$t('operate.title.add')} Servant` : `${$t('operate.title.update')} Servant`"
      width="800px"
      :footShow="!!servantDetailModal.model"
      @on-confirm="saveServantDetail"
      @close="closeServantDetailModal"
      @on-cancel="closeServantDetailModal">
      <let-form
        v-if="servantDetailModal.model"
        ref="servantDetailForm" itemWidth="360px" :columns="2" class="two-columns">
        <!-- <let-form-item :label="$t('serverList.servant.appService')" itemWidth="724px">
          <span>{{servantDetailModal.model.application}}·{{servantDetailModal.model.server_name}}</span>
        </let-form-item> -->
        <let-form-item :label="$t('serverList.servant.objName')" required>
          <let-input
            size="small"
            v-model="servantDetailModal.model.Name"
            :placeholder="$t('serverList.servant.c')"
            required
            pattern="^[A-Za-z]+$"
            :pattern-tip="$t('serverList.servant.obj')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.numOfThread')" required>
          <let-input
            size="small"
            v-model="servantDetailModal.model.Threads"
            :placeholder="$t('serverList.servant.thread')"
            required
            pattern="^[1-9][0-9]*$"
            :pattern-tip="$t('serverList.servant.thread')"
          ></let-input>
        </let-form-item>
        <!-- <let-form-item :label="$t('serverList.table.servant.adress')" required itemWidth="724px">
          <let-input
            ref="endpoint"
            size="small"
            v-model="servantDetailModal.model.endpoint"
            placeholder="tcp -h 127.0.0.1 -t 60000 -p 12000"
            required
            :extraTip="isEndpointValid ? '' :
              $t('serverList.servant.error')"
          ></let-input>
        </let-form-item> -->
        <let-form-item :label="$t('serverList.servant.maxConnecttions')" labelWidth="150px">
          <let-input
            size="small"
            v-model="servantDetailModal.model.Connections"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.lengthOfQueue')" labelWidth="150px">
          <let-input
            size="small"
            v-model="servantDetailModal.model.Capacity"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.queueTimeout')" labelWidth="150px">
          <let-input
            size="small"
            v-model="servantDetailModal.model.Timeout"
          ></let-input>
        </let-form-item>
        <!-- <let-form-item :label="$t('serverList.servant.allowIP')">
          <let-input
            size="small"
            v-model="servantDetailModal.model.allow_ip"
          ></let-input>
        </let-form-item> -->
        <let-form-item :label="$t('serverList.servant.protocol')" required>
          <let-radio v-model="servantDetailModal.model.IsTaf" :label="true">TAF</let-radio>
          <let-radio v-model="servantDetailModal.model.IsTaf" :label="false">{{ $t('serverList.servant.notTAF') }}</let-radio>
        </let-form-item>
        <!-- <let-form-item :label="$t('serverList.servant.treatmentGroup')" labelWidth="150px">
          <let-input
            size="small"
            v-model="servantDetailModal.model.handlegroup"
          ></let-input>
        </let-form-item> -->
      </let-form>
    </let-modal>

    <!-- K8S编辑弹窗 -->
    <let-modal
      v-model="k8sDetailModel.show"
      :title="`${$t('operate.title.update')} Servant`"
      width="800px"
      :footShow="!!k8sDetailModel.model"
      @on-confirm="savek8sDetail"
      @close="closek8sDetailModel"
      @on-cancel="closek8sDetailModel">
      <let-form
        v-if="k8sDetailModel.model"
        ref="k8sDetailForm" itemWidth="360px" :columns="2" class="two-columns">
        <let-form-item :label="$t('deployService.table.th.replicas')">
          <let-input
            size="small"
            type="number"
            :min="0"
            v-model="k8sDetailModel.model.ServerK8S.Replicas"
            :placeholder="$t('deployService.form.placeholder')"
            required
            :required-tip="$t('deployService.table.tips.empty')"
            :pattern-tip="$t('deployService.form.placeholder')"
          ></let-input>
        </let-form-item>
      </let-form>
      <div v-else class="loading-placeholder" ref="k8sDetailModelLoading"></div>
    </let-modal>

    <!-- 更多命令弹窗 -->
    <let-modal
      v-model="moreCmdModal.show"
      :title="$t('operate.title.more')"
      width="700px"
      class="more-cmd"
      @on-confirm="invokeMoreCmd"
      @close="closeMoreCmdModal"
      @on-cancel="closeMoreCmdModal">
      <let-form v-if="moreCmdModal.model" ref="moreCmdForm">
        <let-form-item itemWidth="100%">
          <let-radio v-model="moreCmdModal.model.selected" label="setloglevel">{{$t('serverList.servant.logLevel')}}</let-radio>
          <let-select
            size="small"
            :disabled="moreCmdModal.model.selected !== 'setloglevel'"
            v-model="moreCmdModal.model.setloglevel">
            <let-option v-for="l in logLevels" :key="l" :value="l">{{l}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item itemWidth="100%">
          <let-radio v-model="moreCmdModal.model.selected" label="loadconfig">{{$t('serverList.servant.pushFile')}}</let-radio>
          <let-select
            size="small"
            :placeholder="moreCmdModal.model.configs && moreCmdModal.model.configs.length ? $t('pub.dlg.defaultValue') : $t('pub.dlg.noConfFile')"
            :disabled="!(moreCmdModal.model.configs && moreCmdModal.model.configs.length)
              || moreCmdModal.model.selected !== 'loadconfig'"
            v-model="moreCmdModal.model.loadconfig"
            :required="moreCmdModal.model.selected === 'loadconfig'">
            <let-option v-for="l in moreCmdModal.model.configs" :key="l.filename" :value="l.filename">{{l.filename}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item itemWidth="100%">
          <let-radio v-model="moreCmdModal.model.selected" label="command">{{$t('serverList.servant.sendCommand')}}</let-radio>
          <let-input
            size="small"
            :disabled="moreCmdModal.model.selected !== 'command'"
            v-model="moreCmdModal.model.command"
            :required="moreCmdModal.model.selected === 'command'"
          ></let-input>
        </let-form-item>
        <let-form-item itemWidth="100%">
          <let-radio v-model="moreCmdModal.model.selected" label="connection">{{$t('serverList.servant.serviceLink')}}</let-radio>
        </let-form-item>
        <let-form-item itemWidth="100%">
          <let-radio v-model="moreCmdModal.model.selected" label="undeploy_tars" class="danger">{{$t('operate.undeploy')}} {{$t('common.service')}}</let-radio>
        </let-form-item>
      </let-form>
    </let-modal>
  </div>
</template>

<script>
export default {
  name: 'ServerManage',
  data() {
    return {
      // 当前页面信息
      serverData: {
        level: 5,
        application: '',
        server_name: '',
        set_name: '',
        set_area: '',
        set_group: '',
      },

      // 服务列表
      isreloadlist: false,
      reloadTask: null,
      reloadlist: null,
      serverList: [],

      // 操作历史列表
      serverNotifyList: [],

      // 分页
      pagination: {
        page: 1,
        size: 10,
        total:1,
      },

      // 编辑服务
      serverTypes: [
        'tars_cpp',
        'tars_java',
        'tars_php',
        'tars_nodejs',
        'not_tars',
        'tars_go'
      ],
      configModal: {
        show: false,
        model: null,
      },

      // 编辑servant
      servantModal: {
        show: false,
        model: null,
        currentServer: null,
      },
      servantDetailModal: {
        show: false,
        isNew: true,
        model: null,
      },

      // 编辑K8S
      k8sDetailModel: {
        show: false,
        isNew: false,
        model: null,
      },

      // 更多命令
      logLevels: [
        'NONE',
        'DEBUG',
        'INFO',
        'WARN',
        'ERROR',
      ],
      moreCmdModal: {
        show: false,
        model: null,
        currentServer: null,
      },

      // 失败重试次数
      failCount :0
    };
  },
  props: ['treeid'],
  computed: {
    showOthers() {
      return this.serverData.level === 5;
    },
    isEndpointValid() {
      if (!this.servantDetailModal.model || !this.servantDetailModal.model.endpoint) {
        return false;
      }
      return this.checkServantEndpoint(this.servantDetailModal.model.endpoint);
    },
  },
  methods: {
    gotoLog(data) {
      this.$ajax.getJSON("/server/api/shell_domain").then(domain => {
        let url = `${domain.fromPage}?History=true&NodeIP=${data.NodeIp}&AppName=${data.ServerApp}&ServerName=${data.ServerName}&PodName=${data.PodName}`
        window.open(url)
      })
    },
    getServerId() {
      return this.treeid
    },
    // 状态对应class
    getState(state) {
      let result = ''
      switch(`${state}`){
        case 'Active':
          result = 'status-active'
          break;
        case 'Inactive':
          result = 'status-inactive'
          break;
        case 'Activating':
          result = 'status-activeing'
          break;
        case 'Deactivating':
          result = 'status-deactivating'
          break;
        case 'Unknown':
          result = 'unknown'
          break;
      }
      return result
    },
    reloadServerList() {
      let that = this
      
      let route = that.$route.path
      let path = route.substring(route.lastIndexOf('/') + 1, route.length)
      if(path === 'history'){
        that.reloadTask = setTimeout(() => {
          if(that.isreloadlist){
            that.getServerList()
          }
          that.reloadServerList()
        }, 1000 * 60)
      }else{
        that.stopServerList()
      }
    },
    startServerList(){
      this.isreloadlist = true
    },
    stopServerList(){
      this.isreloadlist = false
    },
    // 获取服务列表
    getServerList() {
      this.$ajax.getJSON('/server/api/pod_history_list', {
        ServerId: this.getServerId(),
        page: this.pagination.page,
        size: this.pagination.size,
      }).then((data) => {
        this.serverList = data.Data
        this.pagination.total = Math.ceil(data.Count.FilterCount / this.pagination.size)
      }).catch((err) => {
        this.$confirm(err.err_msg || err.message || this.$t('serverList.msg.fail'), this.$t('common.alert')).then(() => {
          this.getServerList();
        });
      });
    },
    // 获取服务实时状态
    getServerNotifyList(curr_page) {
      if (!this.showOthers) return;
      const loading = this.$refs.serverNotifyListLoading.$loading.show();

      this.$ajax.getJSON('/server/api/server_notify_list', {
        ServerId: this.getServerId(),
        page_size: this.pageSize,
        curr_page: curr_page,
      }).then((data) => {
        loading.hide();
        this.pageNum = curr_page;
        this.total = Math.ceil(data.Count/this.pageSize);
        this.serverNotifyList = data.rows;
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },
    // 切换服务实时状态页码
    gotoPage(num) {
      this.pagination.page = num
      this.getServerList()
    },
    // 获取服务数据
    getServerConfig(id) {
      const loading = this.$loading.show({
        target: this.$refs.configFormLoading,
      });

      this.$ajax.getJSON('/server/api/server_option_select', {
        ServerId: id,
      }).then((data) => {
        loading.hide();
        if (this.configModal.model) {
          this.configModal.model = Object.assign({}, this.configModal.model, data[0]);
        } else {
          data.templates = [];
          this.configModal.model = data[0];
        }
      }).catch((err) => {
        loading.hide();
        this.closeConfigModal();
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },
    // 编辑服务
    configServer() {
      this.stopServerList()
      this.configModal.show = true;
      
      this.$ajax.getJSON('/server/api/template_select').then((data) => {
        if (this.configModal.model) {
          this.configModal.model.templates = data;
        } else {
          this.configModal.model = { templates: data };
        }
        this.getServerConfig(this.getServerId());
      }).catch((err) => {
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },
    saveConfig() {
      if (this.$refs.configForm.validate()) {
        const loading = this.$Loading.show();
        this.$ajax.postJSON('/server/api/server_option_update', {
          isBak: this.configModal.model.bak_flag,
          ...this.configModal.model,
        }).then((res) => {
          loading.hide();
          this.serverList = this.serverList.map((item) => {
            if (item.id === res.id) {
              return res;
            }
            return item;
          });
          this.closeConfigModal();
          this.$tip.success(this.$t('serverList.restart.success'));
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.message || err.err_msg}`);
        });
      }
    },
    closeConfigModal() {
      if (this.$refs.configForm) this.$refs.configForm.resetValid();
      this.configModal.show = false;
      this.configModal.model = null;
      this.startServerList()
    },

    // 检查任务状态
    checkTaskStatus(taskid, isRetry) {
      return new Promise((resolve, reject) => {
        this.$ajax.getJSON('/server/api/task', {
          task_no: taskid,
        }).then((data) => {
          // 进行中，1秒后重试
          if (data.status === 1 || data.status === 0) {
            setTimeout(() => {
              resolve(this.checkTaskStatus(taskid));
            }, 3000);
          // 成功
          } else if (data.status === 2) {
            resolve(`taskid: ${data.task_no}`);
          // 失败
          }  else {
            reject(new Error(`taskid: ${data.task_no}`));
          }
        }).catch((err) => {
          // 网络问题重试1次
          if (isRetry) {
            reject(new Error(err.err_msg || err.message || this.$t('common.networkErr')));
          } else {
            setTimeout(() => {
              resolve(this.checkTaskStatus(taskid, true));
            }, 3000);
          }
        });
      });
    },
    // 添加任务
    addTask(id, command, tipObj) {
      const loading = this.$Loading.show();
      this.$ajax.postJSON('/server/api/server_command', {
        // serial: true, // 是否串行
        // items: [{
        //   server_id: id,
        //   command,
        // }],
        ServiceId: id,
        ServiceCommand: command,
      }).then((res) => {  // eslint-disable-line
        return this.checkTaskStatus(res).then((info) => {
          loading.hide();
          // 任务成功重新拉取列表
          this.getServerList();
          this.startServerList()
          this.$tip.success({
            title: tipObj.success,
            message: info,
          });
        }).catch((err) => {
          throw err;
        });
      }).catch((err) => {
        loading.hide();
        // 任务失败也重新拉取列表
        this.getServerList();
        this.startServerList()
        this.$tip.error({
          title: tipObj.error,
          message: err.err_msg || err.message || this.$t('common.networkErr'),
        });
      });
    },
    // 重启服务
    restartServer(id) {
      this.stopServerList()
      this.addTask(id, 'restart', {
        success: this.$t('serverList.restart.success'),
        error: this.$t('serverList.restart.failed'),
      });
    },
    // 停止服务
    stopServer(id) {
      this.stopServerList()
      this.$confirm(this.$t('serverList.stopService.msg.stopService'), this.$t('common.alert')).then(() => {
        this.addTask(id, 'stop', {
          success: this.$t('serverList.restart.success'),
          error: this.$t('serverList.restart.failed'),
        });
      }).catch(() => {
        this.startServerList()
      });
    },
    // 下线服务
    undeployServer(id) {
      this.$confirm(this.$t('serverList.dlg.msg.undeploy'), this.$t('common.alert')).then(() => {
        this.addTask(id, 'undeploy_tars', {
          success: this.$t('serverList.restart.success'),
          error: this.$t('serverList.restart.failed'),
        });
        this.closeMoreCmdModal();
      });
    },

    // 管理Servant弹窗
    manageServant(server) {
      this.stopServerList()
      this.servantModal.show = true;

      const loading = this.$loading.show({
        target: this.$refs.servantModalLoading,
      });

      this.$ajax.getJSON('/server/api/server_adapter_select', {
        ServerId: this.getServerId(),
      }).then((data) => {
        loading.hide();
        this.servantModal.model = data;
        this.servantModal.currentServer = server;
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },

    // 管理K8S弹窗
    manageK8S(server) {
      this.stopServerList()
      this.k8sDetailModel.show = true;

      const loading = this.$loading.show({
        target: this.$refs.k8sDetailModelLoading,
      });

      this.$ajax.getJSON('/server/api/server_k8s_select', {
        ServerId: this.getServerId(),
      }).then((data) => {
        loading.hide();
        this.k8sDetailModel.model = data[0];
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },
    closek8sDetailModel() {
      this.k8sDetailModel.show = false;
      this.k8sDetailModel.model = null;
      this.startServerList()
    },

    closeServantModal() {
      this.servantModal.show = false;
      this.servantModal.model = null;
      this.servantModal.currentServer = null;
      this.startServerList()
    },
    // 新增、编辑 servant
    configServant(id) {
      // 新增
      this.servantDetailModal.model = {
        ServerId: this.getServerId(),
        ServerServant: [
          {
            Name: '',
            Port: '',
            Threads: 0,
            Connections: 0,
            Capacity: 0,
            Timeout: 0,
            IsTcp: true,
            IsTaf: false,
          }
        ],
      };
      this.servantDetailModal.isNew = true;
      // 编辑
      if (id) {
        const old = this.servantModal.model.find(item => item.AdapterId === id);
        // old.obj_name = old.servant.split('.')[2];
        // this.servantDetailModal.model = Object.assign({}, this.servantDetailModal.model, old);
        this.servantDetailModal.model = old;
        this.servantDetailModal.isNew = false;
      }
      this.servantDetailModal.show = true;
    },
    closeServantDetailModal() {
      if (this.$refs.servantDetailForm) this.$refs.servantDetailForm.resetValid();
      this.servantDetailModal.show = false;
      this.servantDetailModal.model = null;
    },
    // 检查绑定地址
    checkServantEndpoint(endpoint) {
      const tmp = endpoint.split(/\s-/);
      const regProtocol = /^tcp|udp$/i;
      let regHost = /^h\s(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/i;
      let regT = /^t\s([1-9]|[1-9]\d+)$/i;
      let regPort = /^p\s\d{4,5}$/i;

      let check = true;
      if (regProtocol.test(tmp[0])) {
        let flag = 0;
        for (let i = 1; i < tmp.length; i++) {  // eslint-disable-line
          // 验证 -h
          if (regHost && regHost.test(tmp[i])) {
            flag++;  // eslint-disable-line
            // 提取参数
            this.servantDetailModal.model.node_name = tmp[i].split(/\s/)[1];
            regHost = null;
          }
          // 验证 -t
          if (regT && regT.test(tmp[i])) {
            flag++;  // eslint-disable-line
            regT = null;
          }
          // 验证 -p
          if (regPort && regPort.test(tmp[i])) {
            const port = tmp[i].substring(2);
            if (!(port < 0 || port > 65535)) {
              flag++;  // eslint-disable-line
            }
            regPort = null;
          }
        }
        check = flag === 3;
      } else {
        check = false;
      }
      return check;
    },
    // 保存 servant
    saveServantDetail() {
      if (this.$refs.servantDetailForm.validate()) {
        const loading = this.$Loading.show();
        // 新建
        if (this.servantDetailModal.isNew) {
          const query = this.servantDetailModal.model;
          // query.servant = [query.application, query.server_name, query.obj_name].join('.');
          this.$ajax.postJSON('/server/api/server_adapter_create', query).then((res) => {
            loading.hide();
            this.servantModal.model.unshift(res);
            this.$tip.success(this.$t('common.success'));
            this.closeServantDetailModal();
          }).catch((err) => {
            loading.hide();
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          });
        // 修改
        } else {
          this.servantDetailModal.model.servant=this.servantDetailModal.model.application+'.'+this.servantDetailModal.model.server_name+'.'+this.servantDetailModal.model.obj_name;
          this.$ajax.postJSON('/server/api/server_adapter_update', this.servantDetailModal.model).then((res) => {
            loading.hide();
            this.servantModal.model = this.servantModal.model.map((item) => {
              if (item.id === res.id) {
                return res;
              }
              return item;
            });
            this.$tip.success(this.$t('common.success'));
            this.closeServantDetailModal();
          }).catch((err) => {
            loading.hide();
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          });
        }
      }
    },
    // 保存 k8s
    savek8sDetail() {
      if (this.$refs.k8sDetailForm.validate()) {
        const loading = this.$Loading.show();
        // 新建
        if (this.k8sDetailModel.isNew) {
          const query = this.k8sDetailModel.model;
          // query.servant = [query.application, query.server_name, query.obj_name].join('.');
          this.$ajax.postJSON('/server/api/server_k8s_create', query).then((res) => {
            loading.hide();
            this.servantModal.model.unshift(res);
            this.$tip.success(this.$t('common.success'));
            this.closek8sDetailModel();
          }).catch((err) => {
            loading.hide();
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          });
        // 修改
        } else {
          this.$ajax.postJSON('/server/api/server_k8s_update', this.k8sDetailModel.model).then((res) => {
            loading.hide();
            this.$tip.success(this.$t('common.success'));
            this.getServerList()
            this.closek8sDetailModel();
          }).catch((err) => {
            loading.hide();
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          });
        }
      }
    },
    // 删除 servant
    deleteServant(id) {
      this.$confirm(this.$t('serverList.servant.a'), this.$t('common.alert')).then(() => {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/server_adapter_delete', {
          AdapterId: id,
        }).then((res) => {
          loading.hide();
          this.servantModal.model = this.servantModal.model.filter( item => item.id !== id);
          this.$tip.success(this.$t('common.success'));
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
        });
      });
    },

    // 显示更多命令
    showMoreCmd(server) {
      this.moreCmdModal.model = {
        selected: 'setloglevel',
        setloglevel: 'NONE',
        loadconfig: '',
        command: '',
        configs: null,
      };
      this.moreCmdModal.unwatch = this.$watch('moreCmdModal.model.selected', () => {
        if (this.$refs.moreCmdForm) this.$refs.moreCmdForm.resetValid();
      });
      this.moreCmdModal.show = true;
      this.moreCmdModal.currentServer = server;

      // this.$ajax.getJSON('/server/api/config_file_list', {
      //   level: 5,
      //   application: server.application,
      //   server_name: server.server_name,
      // }).then((data) => {
      //   if (this.moreCmdModal.model) this.moreCmdModal.model.configs = data;
      // }).catch((err) => {
      //   this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
      // });
    },
    sendCommand(id, command, hold) {
      const loading = this.$Loading.show();
      this.$ajax.getJSON('/server/api/send_command', {
        server_ids: id,
        command,
      }).then((res) => {
        loading.hide();
        const msg = res[0].err_msg.replace(/\n/g, '<br>');
        if (res[0].ret_code === 0) {
          const opt = {
            title: this.$t('common.success'),
            message: msg,
          };
          if (hold) opt.duration = 0;
          this.$tip.success(opt);
        } else {
          throw new Error(msg);
        }
      }).catch((err) => {
        loading.hide();
        this.$tip.error({
          title: this.$t('common.error'),
          message: err.err_msg || err.message,
        });
      });
    },
    invokeMoreCmd() {
      const model = this.moreCmdModal.model;
      const server = this.moreCmdModal.currentServer;
      // 下线服务
      if (model.selected === 'undeploy_tars') {
        this.undeployServer(server.id);
      // 设置日志等级
      } else if (model.selected === 'setloglevel') {
        this.sendCommand(server.id, `taf.setloglevel ${model.setloglevel}`);
      // push 日志文件
      } else if (model.selected === 'loadconfig' && this.$refs.moreCmdForm.validate()) {
        this.sendCommand(server.id, `taf.loadconfig ${model.loadconfig}`);
      // 发送自定义命令
      } else if (model.selected === 'command' && this.$refs.moreCmdForm.validate()) {
        this.sendCommand(server.id, model.command);
      // 查看服务链接
      } else if (model.selected === 'connection') {
        this.sendCommand(server.id, `taf.connection`, true);
      }
    },
    closeMoreCmdModal() {
      if (this.$refs.moreCmdForm) this.$refs.moreCmdForm.resetValid();
      if (this.moreCmdModal.unwatch) this.moreCmdModal.unwatch();
      this.moreCmdModal.show = false;
      this.moreCmdModal.model = null;
    },

    // 处理未发布时间显示
    handleNoPublishedTime(timeStr, noPubTip = this.$t('pub.dlg.unpublished')) {
      if (timeStr === '0000:00:00 00:00:00') {
        return noPubTip;
      }
      return timeStr;
    },
  },
  created() {
    this.serverData = this.$parent.getServerData();

    this.$once('hook:beforeDestroy', () => {            
      clearInterval(this.reloadTask)
      this.reloadTask = null
    })
  },
  mounted() {
    this.getServerList();
    this.reloadServerList()
    this.startServerList()
    this.getServerNotifyList(1);
  },
};
</script>

<style>
@import '../../assets/css/variable.css';

.page_server_manage {
  .tbm16 {
    margin: 16px 0;
  }
  .danger {
    color: var(--off-color);
  }

  .more-cmd {
    .let-form-item__content {
      display: flex;
      align-items: center;
    }
    span.let-radio {
      margin-right: 5px;
    }
    label.let-radio {
      width: 200px;
    }
  }
}

</style>
