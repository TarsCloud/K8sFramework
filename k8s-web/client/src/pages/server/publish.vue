<template>
  <div class="page_server_publish">
    <!-- 服务列表 -->
    <div v-if="!showHistory">
      <let-table ref="table" :data="totalServerList" :title="$t('serverList.title.patchList')" :empty-msg="$t('common.nodata')">
        <!-- <let-table-column>
          <template slot="head" slot-scope="props">
            <let-checkbox v-model="isCheckedAll"></let-checkbox>
          </template>
          <template slot-scope="scope">
            <let-checkbox v-model="scope.row.isChecked" :value="scope.row.StateId"></let-checkbox>
          </template>
        </let-table-column> -->
        <let-table-column :title="$t('deployService.form.app')" prop="ServerApp"></let-table-column>
        <let-table-column :title="$t('deployService.form.serviceName')" prop="ServerName"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.version')" prop="ServiceVersion"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.currStatus')">
          <template slot-scope="scope">
            <span :class="getState(scope.row.Enabled)"></span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('deployService.form.serviceMark')" prop="ServiceMark"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.createPerson')" prop="CreatePerson"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.createTime')" prop="CreateTime"></let-table-column>
        <!-- <let-table-column :title="$t('serverList.table.th.enablePerson')" prop="EnablePerson"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.time')" prop="EnableTime"></let-table-column> -->
        <!-- <let-table-column :title="$t('serverList.table.th.enableSet')">
          <template slot-scope="scope">
            <span>{{ scope.row.enable_set ? $t('common.enable') : $t('common.disable') }}</span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('common.set.setName')" prop="set_name"></let-table-column>
        <let-table-column :title="$t('common.set.setArea')" prop="set_area"></let-table-column>
        <let-table-column :title="$t('common.set.setGroup')" prop="set_group"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.configStatus')">
          <template slot-scope="scope">
            <span :class="scope.row.setting_state == 'active' ? 'status-active' : 'status-off'"></span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.th.currStatus')">
           <template slot-scope="scope">
            <span :class="scope.row.present_state == 'active' ? 'status-active' : 'status-off'"></span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.th.version')" prop="patch_version"></let-table-column> -->
        <!-- <let-table-column :title="$t('serverList.table.th.time')">
          <template slot-scope="scope">
            <span style="word-break: break-word">{{handleNoPublishedTime(scope.row.patch_time)}}</span>
          </template>
        </let-table-column> -->
      </let-table>
      <div slot="operations" style="overflow:hidden;">
        <let-pagination align="right" style="float:right;"
          :page="pagination.page" @change="gotoPage"
          :total="pagination.total">
        </let-pagination>
        <div class="btn_group" style="float:left;">
          <let-button theme="primary" size="small" @click="showUploadModal">{{$t('pub.dlg.upload')}}</let-button>
          <let-button theme="primary" size="small" @click="openPublishVersionModal">{{$t('pub.btn.pub')}}</let-button>
          <!-- <let-button size="small" v-if="serverList && serverList.length > 0" @click="gotoHistory">{{$t('pub.btn.history')}}</let-button> -->
        </div>
      </div>
      <!-- 发布服务弹出框 -->
      <let-modal
        v-model="publishModal.show"
        :title="$t('index.rightView.tab.patch')"
        width="880px"
        :footShow="false"
        @close="closePublishModal"
        @on-confirm="savePublishServer">
          <let-form
            v-if="publishModal.model"
            ref="publishForm"
            itemWidth="100%">
              <!-- <let-form-item :label="$t('pub.dlg.ip')">
                <div v-for="server in publishModal.model.serverList" :key="server.id">{{server.node_name}}</div>
              </let-form-item> -->
              <!-- <let-form-item :label="$t('pub.dlg.patchType')" v-if="patchRadioData.length>1">
                <let-radio-group type="button" size="small" @change="patchChange" v-model="patchType" :data="patchRadioData">
                </let-radio-group>
              </let-form-item> -->
              <template v-if="publishModal.model.show">
                <let-form-item :label="$t('pub.dlg.releaseVersion')" required>
                  <let-select
                    size="small"
                    v-model="publishModal.model.ServiceId"
                    required
                    :required-tip="$t('pub.dlg.ab')"
                  >
                    <let-option v-for="d in publishModal.model.patchList" :key="d.ServiceId" :value="d.ServiceId">
                      <i class="let-icon let-icon-gou let-tag_success" v-if="`${d.Enabled}` === 'true'"></i>
                      <span>{{d.ServiceVersion}} | {{d.CreateTime}} | {{d.ServiceImage}}</span>
                    </let-option>
                  </let-select>
                </let-form-item>
                <let-form-item :label="$t('deployService.table.th.replicas')">
                  <let-input
                    size="small"
                    type="number"
                    :min="1"
                    v-model="publishModal.model.Replicas"
                    :placeholder="$t('deployService.form.placeholder')"
                    required
                    :required-tip="$t('deployService.table.tips.empty')"
                    :pattern-tip="$t('deployService.form.placeholder')"
                  ></let-input>
                </let-form-item>
                <let-form-item :label="$t('serverList.servant.comment')">
                  <let-input v-model="publishModal.model.EnableMark"></let-input>
                  <!-- <let-button theme="primary" size="small" class="mt10" @click="showUploadModal">{{$t('pub.dlg.upload')}}</let-button> -->
                  <!-- <br> -->
                  <let-button theme="primary" size="small" class="mt20" @click="savePublishServer">{{$t('common.patch')}}</let-button>
                </let-form-item>
              </template>
              <template v-else>
                <let-form-item :label="$t('serverList.table.th.version')">
                  <let-select size="small" required :required-tip="$t('deployService.table.tips.empty')"
                    v-model="tagVersion"
                    requred>
                    <let-option v-for="it in tagList" :key="`${it.version}`" :value="it.path +'--'+ it.version">{{it.version}}</let-option>
                  </let-select>
                  <let-button theme="primary" size="small" class="mt10" @click="addCompileTask">{{$t('pub.dlg.compileAndPublish')}}</let-button>
                  <let-button size="small" class="mt10" @click="openPubConfModal" v-if="false">{{$t('pub.dlg.conf')}}</let-button>
                </let-form-item>
              </template>
            </let-form>
      </let-modal>

     <!-- 上传发布包弹出框 -->
      <let-modal
        v-model="uploadModal.show"
        :title="$t('pub.dlg.upload')"
        width="880px"
        :footShow="false"
        @on-cancel="closeUploadModal">
        <let-form
          v-if="uploadModal.model"
          ref="uploadForm"
          itemWidth="100%"
          @submit.native.prevent="uploadPatchPackage">
            <let-form-item :label="$t('pub.dlg.releasePkg')" itemWidth="400px">
               <let-uploader
                :placeholder="$t('pub.dlg.defaultValue')"
                @upload="uploadFile">{{$t('common.submit')}}
                </let-uploader>
                <span v-if="uploadModal.model.file">{{uploadModal.model.file.name}}</span>
            </let-form-item>
            <let-form-item :label="$t('deployService.form.serviceType')" required>
              <let-select
                size="small"
                v-model="uploadModal.model.ServerType"
                required
                :required-tip="$t('deployService.form.serviceTypeTips')"
              >
                <let-option v-for="d in serverType" :key="d" :value="d">{{d}}</let-option>
              </let-select>
            </let-form-item>
            <let-form-item :label="$t('serverList.servant.comment')">
              <let-input
                type="textarea"
                :rows="3"
                v-model="uploadModal.model.comment"
              >
              </let-input>
            </let-form-item>
            <let-button type="submit" theme="primary">{{$t('serverList.servant.upload')}}</let-button>
        </let-form>
      </let-modal>

     <!-- 发布结果弹出框 -->
     <let-modal
        v-model="finishModal.show"
        :title="$t('serverList.table.th.result')"
        width="880px"
        :footShow="false"
        @on-cancel="closeFinishModal">
        <let-table
          v-if="finishModal.model"
          :title="$t('serverList.servant.taskID') + ': ' + finishModal.model.task_no"
          :data="finishModal.model.items">
          <let-table-column :title="$t('serverList.table.th.buildID')" prop="BuildId"></let-table-column>
          <let-table-column :title="$t('serverList.table.th.buildTime')" width="180">
            {{ formatDate(new Date(), 'YYYY-MM-DD HH:mm:ss') }}
          </let-table-column>
          <let-table-column :title="$t('common.status')">
            <template slot-scope="scope">
              <let-tag style="margin-right:15px;white-space:nowrap;"
                :theme="scope.row.BuildStatus === 'done' ? 'success' : (scope.row.BuildStatus === 'error' ? 'danger' : '')" checked>
                {{ scope.row.BuildMessage || scope.row.BuildStatus }}
              </let-tag>
            </template>
          </let-table-column>
        </let-table>
      </let-modal>
    </div>

    <!-- 发布历史 -->
    <div v-if="showHistory">
      <let-form inline itemWidth="300px" @submit.native.prevent="getHistoryList">
        <let-form-item :label="$t('pub.date')">
          <let-date-range-picker :start.sync="startTime" :end.sync="endTime"></let-date-range-picker>
        </let-form-item>
        <let-form-item>
          <let-button type="submit" theme="primary">{{$t('operate.search')}}</let-button>
        </let-form-item>
      </let-form>
      <let-table ref="historyTable" v-if="totalHistoryList && totalHistoryList.length > 0" :data="totalHistoryList" :title="$t('historyList.title')" :empty-msg="$t('common.nodata')">
        <let-table-column :title="$t('serverList.servant.taskID')" prop="task_no"></let-table-column>
        <let-table-column :title="$t('historyList.table.th.c2')">
          <template slot-scope="scope">
            <span>{{scope.row.serial ? $t('common.yes'): $t('common.no')}}</span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.dlg.title.taskStatus')">
          <template slot-scope="scope">
            <span>{{ statusMap[scope.row.status] || '-'}}</span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('historyList.table.th.c4')">
          <template slot-scope="scope">
            <let-table-operation @click="viewTask(scope.row.task_no)">{{$t('operate.view')}}</let-table-operation>
          </template>
        </let-table-column>
        <let-pagination slot="pagination" align="right"
          :total="historyTotalPage" :page="historyPage" @change="changeHistoryPage">
      </let-pagination>
        <div slot="operations" style="margin-left: -15px;">
          <let-button theme="primary" size="small" @click="showHistory=false">{{$t('operate.goback')}}</let-button>
        </div>
      </let-table>


      <!-- 子任务详情弹出框 -->
      <let-modal
        v-model="taskModal.show"
        :title="$t('historyList.table.th.c4')"
        width="880px"
        :footShow="false"
        @on-cancel="taskModal.show = false">
        <let-table
          v-if="taskModal.model"
          :data="taskModal.model.items">
          <let-table-column :title="$t('historyList.dlg.th.c1')" prop="item_no"></let-table-column>
          <let-table-column :title="$t('historyList.dlg.th.c2')" prop="application"></let-table-column>
          <let-table-column :title="$t('historyList.dlg.th.c3')" prop="server_name"></let-table-column>
          <let-table-column :title="$t('historyList.dlg.th.c4')" prop="node_name"></let-table-column>
          <let-table-column :title="$t('historyList.dlg.th.c5')" prop="command"></let-table-column>
          <let-table-column :title="$t('monitor.search.start')" prop="start_time"></let-table-column>
          <let-table-column :title="$t('monitor.search.end')" prop="end_time"></let-table-column>
          <let-table-column :title="$t('common.status')" prop="status_info"></let-table-column>
          <let-table-column :title="$t('historyList.dlg.th.c7')" prop="execute_info"></let-table-column>
        </let-table>
      </let-modal>
   </div>

  <!-- 配置编译接口 -->
   <let-modal
        v-model="publishUrlConfModal.show"
        :title="$t('pub.dlg.conf')"
        width="800px"
        :footShow="true"
        @on-confirm="saveCompilerUrl"
        @on-cancel="publishUrlConfModal.show = false">
        <let-form ref="compilerForm"
          itemWidth="100%"
          v-if="publishUrlConfModal.model"
          required>
          <let-form-item :label="$t('pub.dlg.tag')">
            <let-input size="small"
              v-model="publishUrlConfModal.model.tag"
              :placeholder="$t('pub.tips.tag')"
              :required-tip="$t('deployService.table.tips.empty')"
              required ></let-input>
          </let-form-item>
        </let-form>
    </let-modal>

    <!-- 编译进度 -->
    <let-modal
        v-model="compilerModal.show"
        :title="$t('pub.dlg.compileProgress')"
        width="880px"
        :footShow="false">
        <let-table
          v-if="compilerModal.model"
          :data="compilerModal.model.progress">
            <let-table-column :title="$t('historyList.dlg.th.c2')" prop="application"></let-table-column>
            <let-table-column :title="$t('historyList.dlg.th.c3')" prop="server_name"></let-table-column>
            <let-table-column :title="$t('historyList.dlg.th.c4')" prop="node"></let-table-column>
            <let-table-column :title="$t('historyList.dlg.th.c8')" prop="status">
              <template slot-scope="scope">
                <span v-if="scope.row.state=='1'" class="running">{{scope.row.status}}</span>
                <span v-else-if="scope.row.state=='2'" class="success">{{scope.row.status}}</span>
                <span v-else class="stop">{{scope.row.status}}</span>
              </template>
            </let-table-column>
            <let-table-column :title="$t('monitor.search.start')" prop="start_time"></let-table-column>
            <let-table-column :title="$t('monitor.search.end')" prop="end_time"></let-table-column>
        </let-table>
    </let-modal>

  </div>
</template>

<script>
import { formatDate } from '@/lib/date';

export default {
  name: 'ServerPublish',
  data() {
    return {
      activeKey: '',
      treeData: [],
      totalServerList: [],
      serverList: [],
      isCheckedAll: false,
      // 分页
      pagination: {
        page: 1,
        size: 10,
        total:1,
      },
      publishModal: {
        show: false,
        model: {
          patchList: [],
        },
      },
      finishModal: {
        show: false,
        model: {
          task_no : '',
          items : []
        },
      },
      serverK8S: {},
      serverType: [],
      statusConfig: {
        0: this.$t('serverList.restart.notStart'),
        1: this.$t('serverList.restart.running'),
        2: this.$t('serverList.restart.success'),
        3: this.$t('serverList.restart.failed'),
        4: this.$t('serverList.restart.cancel'),
        5: this.$t('serverList.restart.parial'),
      },
      statusMap: {
        0: 'EM_T_NOT_START',
        1: 'EM_T_RUNNING',
        2: 'EM_T_SUCCESS',
        3: 'EM_T_FAILED',
        4: 'EM_T_CANCEL',
        5: 'EM_T_PARIAL',
      },
      showHistory: false,
      startTime: '',
      endTime: '',
      totalHistoryList: [],
      historyList: [],
      historyTotalPage: 0,
      historyPage: 1,
      historyPageSize: 20,
      taskModal: {
        show: false,
        modal: true,
      },
      uploadModal: {
        show: false,
        model: null,
      },
      patchType: 'patch',
      patchRadioData: [
        {value:'patch', text:this.$t('pub.dlg.upload')}

      ],
      tagList: [],
      tagVersion: '',
      publishUrlConfModal: {
        show: false,
        model: {tag:'',compiler:'',task:''},
      },
      compilerModal: {
        show: false,
        model: null
      },
      pkgUpload: {
        show : false,
        model: null
      }
    };
  },
  props: ['treeid'],
  methods: {
    getServerId() {
      return this.treeid
    },
    getServerType() {
      return this.$ajax.getJSON('/server/api/server_list', {
        ServerId: this.getServerId(),
      })
    },
    getDefault() {
      this.$ajax.getJSON('/server/api/default').then(data => {
        this.serverK8S = data.ServerK8S || {}
        this.serverType = data.ServerTypeOptional || []        
      })
    },
    // 状态对应class
    getState(state) {
      let result = ''
      switch(`${state}`){
        case 'true':
          result = 'status-active'
          break;
      }
      return result
    },
    formatDate(date, formatter){
      return formatDate(date, 'YYYY-MM-DD HH:mm:ss')
    },
    getCompileConf(){
      this.$ajax.getJSON('/server/api/get_compile_conf').then((data) =>{
        if(data.enable) {
          this.patchRadioData.push({value:'compile', text:this.$t('pub.dlg.compileAndPublish')});
        }
      }).catch((err) =>{
        this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
      })
    },
    getServerList() {
      // 获取服务列表
      this.$ajax.getJSON('/server/api/patch_list', {
        ServerId: this.getServerId(),
        page: this.pagination.page,
        size: this.pagination.size,
      }).then((data) => {
        if (!data.hasOwnProperty("Data")) {
          return
        }
        this.getPatchEnabled().then(dataEnabled => {
          const items = data.Data || [];
          this.pagination.total = Math.ceil(data.Count.FilterCount / this.pagination.size)
          if (dataEnabled.hasOwnProperty("Data")) {
            items.forEach((item) => {
              if(`${item.ServiceId}` === `${dataEnabled.Data[0].ServiceId}`){
                item.Enabled = true
              }
            });
          }
          this.totalServerList = items;
          this.totalPage = Math.ceil(this.totalServerList.length / this.pageSize);
          this.page = 1;
          this.updateServerList();
        }).catch((err) => {
          this.$confirm(err.err_msg || err.message || this.$t('serverList.table.msg.fail')).then(() => {
            this.getServerList();
          });
        })
      }).catch((err) => {
        this.$confirm(err.err_msg || err.message || this.$t('serverList.table.msg.fail')).then(() => {
          this.getServerList();
        });
      });
    },
    // 切换服务实时状态页码
    gotoPage(num) {
      this.pagination.page = num
      this.getServerList()
    },
    openPublishVersionModal() {
      let ServerId = this.getServerId()
      this.publishModal.model = {
        ServerId,
        Replicas: this.serverK8S.Replicas || 1,
        show: true
      };
      this.getPatchList(ServerId, 1, 50).then((data) => {
        this.getK8SData().then(k8sData => {
          this.publishModal.model.Replicas = k8sData[0].Replicas || this.serverK8S.Replicas || 1
        })
        if (!data.hasOwnProperty("Data")) {
          return
        }
        this.getPatchEnabled().then(dataEnabled => {
          if (dataEnabled.hasOwnProperty("Data")) {
            data.Data.forEach(item => {
              if(`${item.ServiceId}` === `${dataEnabled.Data[0].ServiceId}`){
                item.Enabled = true
              }
            })
          }
          this.publishModal.model.patchList = data.Data
          window.setTimeout(() => this.publishModal.show = true, 300)
        })
      });
    },
    getPatchList(ServerId, currPage, pageSize) {
      return this.$ajax.getJSON('/server/api/patch_list', {
        ServerId,
        curr_page : currPage,
        page_size : pageSize
      });
    },
    getPatchEnabled() {
      return this.$ajax.getJSON('/server/api/patch_enabled', {
        ServerId: this.getServerId(),
      });
    },
    getK8SData() {
      return this.$ajax.getJSON('/server/api/server_k8s_select', {
        ServerId: this.getServerId(),
      });
    },
    closePublishModal() {
      // 关闭发布弹出框
      this.publishModal.show = false;
      this.publishModal.modal = null;
      this.patchType = 'patch';
      this.$refs.publishForm.resetValid();
    },
    savePublishServer() {
      // 发布
      const loading = this.$Loading.show();
      if (this.$refs.publishForm.validate()) {
        loading.hide()
        this.$ajax.postJSON('/server/api/patch_publish', {
          ServerId: this.publishModal.model.ServerId,
          ServiceId: this.publishModal.model.ServiceId,
          Replicas: this.publishModal.model.Replicas,
          EnableMark: this.publishModal.model.EnableMark,
        }).then((data) => {
          loading.hide();
          this.getServerList()
          this.closePublishModal();
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });



        // this.publishModal.model.patchList.forEach((item) => {
        //   items.push({
        //     ServerId: item.ServerId.toString(),
        //     command: 'patch_tars',
        //     parameters: {
        //       ServerVersion: this.publishModal.model.ServiceVersion.toString(),
        //       bak_flag: item.bak_flag,
        //       update_text: this.publishModal.model.update_text,
        //     },
        //   });
        // });
        // const loading = this.$Loading.show();
        // this.$ajax.postJSON('/server/api/add_task', {
        //   serial: true,
        //   items,
        // }).then((data) => {
        //   loading.hide();
        //   this.closePublishModal();
        //   this.finishModal.model.task_no = data;
        //   this.finishModal.show = true;
        //   // 实时更新状态
        //   this.getTaskRepeat(data);
        // }).catch((err) => {
        //   loading.hide();
        //   this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        // });
      }
    },
    closeFinishModal() {
      // 关闭发布结果弹出框
      this.finishModal.show = false;
      this.finishModal.modal = null;
      this.$refs.finishForm.resetValid();
    },
    getTaskRepeat({ ServerId, ServerType, BuildId, ServiceImage, CreateMark }) {
      let timerId;
      timerId && clearTimeout(timerId);
      const getTask = () => {
        this.$ajax.getJSON('/server/api/patch_upload_status', {
          ServerId,
          ServerType,
          BuildId,
          ServiceImage,
          CreateMark,
        }).then((data) => {
          let done = true;
          if(data.BuildStatus === 'working'){
            done = false
          }
          if(done){
            clearTimeout(timerId)
            this.getServerList()
          }else{
            timerId = setTimeout(getTask, 2000);
          }
          this.finishModal.model.items = [data]
        }).catch((err) => {
          clearTimeout(timerId);
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });
      };
      getTask();
    },
    updateServerList() {
      // 更新服务列表
      const start = (this.page - 1) * this.pageSize;
      const end = this.page * this.pageSize;
      this.serverList = this.totalServerList.slice(start, end);
    },
    updateServerInfo(data) {
      this.$ajax.postJSON('/server/api/server_update', {
        ServerId: data.ServerId,
        ServerType: data.ServerType,
        ServerMark: data.ServerMark,
      })
    },
    gotoHistory() {
      // 切换到发布历史
      this.showHistory = true;
      this.getHistoryList(1);
    },
    getHistoryList(curr_page) {
      // 更新历史记录
      if(typeof curr_page != 'number'){
        curr_page = 1;
      }
      const loading = this.$Loading.show();
      const params = {
        application: this.serverList[0].application || '',
        server_name: this.serverList[0].server_name || '',
        from: this.startTime,
        to: this.endTime,
        page_size: this.historyPageSize,
        curr_page: curr_page
      };
      this.historyPage = curr_page;
      this.$ajax.getJSON(`/server/api/task_list`, params).then((data) => {
        loading.hide();
        this.totalHistoryList = data.rows || [];
        this.historyTotalPage = Math.ceil(data.Count / this.historyPageSize);
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
      });
    },
    viewTask(taskId) {
      this.$ajax.getJSON('/server/api/task', {
        task_no: taskId,
      }).then((data) => {
        this.taskModal.model = data;
        this.taskModal.show = true;
      });
    },
    changeHistoryPage(page) {
      this.getHistoryList(page);
    },
    showUploadModal() {
      this.getServerType().then(data => {
        let ServerId = this.getServerId()
        this.uploadModal.model = {
          ServerId,
          file: null,
          ServerMark: '',
          ServerType: data.Data[0].ServerType,
        };
        this.uploadModal.show = true;
      })
    },
    closeUploadModal() {
      // 关闭上传文件弹出框
      this.uploadModal.show = false;
      this.uploadModal.model = null;
      this.$refs.uploadForm.resetValid();
    },
    uploadFile(file) {
      this.uploadModal.model.file = file;
    },
    uploadPatchPackage() {
      let ServerId = this.getServerId()

      // 上传发布包
      if (this.$refs.uploadForm.validate()) {
        const loading = this.$Loading.show();
        const formdata = new FormData();
        formdata.append('ServerId', ServerId);
        formdata.append('ServerType', this.uploadModal.model.ServerType);
        formdata.append('suse', this.uploadModal.model.file);
        // if(this.uploadModal.model.comment){
        //   formdata.append('CreateMark', this.uploadModal.model.comment);
        // }
        // formdata.append('task_id', new Date().getTime());
        this.$ajax.postForm('/server/api/patch_upload', formdata).then((data) => {
          this.getPatchList(ServerId,1,10).then((data) => {
            loading.hide();
            this.publishModal.model.patchList = data.Data;
            this.closeUploadModal();
          });

          this.getServerList()

          // 发布状态
          this.finishModal.model.task_no = data.data.BuildId;
          this.finishModal.show = true;
          // 实时更新状态
          this.getTaskRepeat({
            ServerId,
            ServerType: this.uploadModal.model.ServerType,
            BuildId: data.data.BuildId,
            ServiceImage: data.data.BuildImage,
            CreateMark: this.uploadModal.model.comment
          });
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });
      }
    },

    // 处理未发布时间显示
    handleNoPublishedTime(timeStr, noPubTip = this.$t('pub.dlg.unpublished')) {
      if (timeStr === '0000:00:00 00:00:00') {
        return noPubTip;
      }
      return timeStr;
    },
    // 选择上传方式或编译发布方式
    patchChange() {
      if(this.patchType=='patch'){
        this.publishModal.model.show = true;
      }else {
        this.publishModal.model.show = false;
        this.getCodeVersion();
      }
    },
    getCodeVersion() {
      this.$ajax.get('/server/api/get_tag_list',{
        application: this.publishModal.model.application,
        server_name: this.publishModal.model.server_name
      }).then(data => {
        if(data.data=='') {
          this.openPubConfModal();
        }else {
          this.tagList = data.data;
        }
      }).catch(e=> {
        this.tagList = [];
        this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
      })
    },
    openPubConfModal() {
      this.publishUrlConfModal.show = true;
      this.$ajax.getJSON('/server/api/get_tag_conf', {
        application : this.publishModal.model.application,
        server_name : this.publishModal.model.server_name,
      }).then(data => {
          this.publishUrlConfModal.model.tag = data.path;
      }).catch(err => {
        this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
      })
    },
    saveCompilerUrl() {
      if(this.$refs.compilerForm.validate()) {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/set_tag_conf', {
          path: this.publishUrlConfModal.model.tag,
          application : this.publishModal.model.application,
          server_name : this.publishModal.model.server_name,
          }).then(data => {
          loading.hide();
          this.$tip.success(this.$t('common.success'));
          this.publishUrlConfModal.show = false;
          this.getCodeVersion();
        }).catch(err => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
        })
      }
    },
    addCompileTask() {
      this.$ajax.getJSON('/server/api/get_compile_conf').then(data => {
          const compileUrl = data.getVersionList;
          if(!compileUrl) {
            this.openPubConfModal();
            return;
          }else {
            let nodes = this.publishModal.model.serverList.map(item => item.node_name);
            let opts = {
              application : this.publishModal.model.application,
              server_name : this.publishModal.model.server_name,
              node : nodes.join(';'),
              path : this.tagVersion.split('--')[0],
              version : this.tagVersion.split('--')[1],
              comment : this.publishModal.model.update_text || '',
              compileUrl : compileUrl
            };
            const loading = this.$Loading.show();
            this.$ajax.postJSON('/server/api/do_compile', opts).then(data => {
                loading.hide();
                this.compilerModal.show = true;
                const taskNo = typeof data === 'string' ? data : data.data;
                this.getStatus(taskNo);
                //this.taskStatus(taskNo);
            }).catch(err => {
                loading.hide();
                this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
            })
          }
      }).catch(err => {
        this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
      })
    },
    taskStatus(taskNo) {
        this.getStatus(taskNo);
    },
    getStatus(taskNo) {
        const f = () => {
          let t = null;
          t && clearTimeout(t);
          this.$ajax.getJSON('/server/api/compiler_task', {taskNo}).then(data =>{
            const ret = typeof data === 'array' ? data : data.data;
            ret[0].status = this.statusConfig[ret[0].state];
            if(ret[0].state==1){
                t = setTimeout(f, 2000);
            }
            if(this.compilerModal.model) {
              Object.assign(this.compilerModal.model, {progress : ret})
            }else {
              this.compilerModal.model = {progress : ret};
            }
            // 编译成功后轮询发布包回传情况
            if(ret[0].state==2){
                const loading = this.$Loading({text:'回传发布包'});
                loading.show();
                this.compilerModal.show = false;
                let timer = ()=>{
                  this.$ajax.getJSON('/server/api/get_server_patch', {task_id : taskNo}).then(data => {
                    if(Object.keys(data).length !== 0) {
                      loading.hide();
                      this.publishModal.model.patch_id = data.id;
                      this.publishModal.show = false;
                      this.savePublishServer();
                    }else {
                      setTimeout(timer, 2000);
                    }
                  }).catch(err => {
                    loading.hide();
                    this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
                  });
                }
                setTimeout(timer, 2000);
            }
          }).catch(err =>{
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          })
        }
        f();
    }
  },
  mounted() {
    this.getServerList();
    this.getDefault()
    // this.getCompileConf();
  },
  watch: {
    isCheckedAll() {
      const isCheckedAll = this.isCheckedAll;
      this.serverList.forEach((item) => {
        item.isChecked = isCheckedAll;
      });
    },
    page() {
      this.updateServerList();
    }
  },
};
</script>

<style>
@import '../../assets/css/variable.css';

.page_server_publish {
  padding-bottom: 32px;
  .mt20 {
    margin-top: 20px;
  }
  .running {
    color:#3f5ae0
  }
  .success {color:#6accab}
  .stop {color: #f56c77}

  .btn_group .let-button + .let-button{margin-left:10px;}
}
</style>
