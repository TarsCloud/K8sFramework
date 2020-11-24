<template>
  <div class="page_server_manage">

    <!-- 服务列表 -->
    <div v-if="serverList" ref="serverListLoading" style="position:relative">
      <div class="btn_group">
        <let-button theme="primary" size="small" @click="configServer">{{$t('operate.server')}}</let-button>
        <let-button theme="primary" size="small" @click="manageK8S">{{$t('operate.k8s')}}</let-button>
        <let-button theme="primary" size="small" @click="manageServant">{{$t('operate.servant')}}</let-button>
        <let-button theme="primary" size="small" @click="viewTemplate">{{$t('operate.viewTemplate')}}</let-button>
        <let-button theme="primary" size="small" @click="privateTemplateManage">{{$t('operate.privateTemplateManage')}}</let-button>
        <let-button theme="primary" size="small" @click="startServer">{{$t('operate.startServer')}}</let-button>
        <let-button theme="primary" size="small" @click="restartServer">{{$t('operate.restartServer')}}</let-button>
        <let-button theme="primary" size="small" @click="stopServer">{{$t('operate.stopServer')}}</let-button>
        <let-button theme="primary" size="small" @click="showMoreCmd">{{$t('operate.more')}}</let-button>
      </div>
      <let-table :data="serverList" :title="$t('serverList.title.serverList')" :empty-msg="$t('common.nodata')">
        <let-table-column>
          <template slot="head" slot-scope="props">
            <let-checkbox v-model="isCheckedAll"></let-checkbox>
          </template>
          <template slot-scope="scope">
            <let-checkbox v-model="scope.row.isChecked" @change="checkChange(scope.row)"></let-checkbox>
          </template>
        </let-table-column>
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
        <let-table-column :title="$t('serverList.table.th.configStatus')">
          <template slot-scope="scope">
            <span :class="getState(scope.row.SettingState)"></span>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.th.currStatus')" prop="PresentState"></let-table-column>
        <let-table-column :title="$t('serverList.table.th.currMessage')" width="120px">
          <template slot-scope="scope">
            <let-tooltip class="tooltip" placement="top" :content="scope.row.PresentMessage || '空'">
              <let-table-operation>查看</let-table-operation>
            </let-tooltip>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.th.updateTime')" prop="UpdateTime"></let-table-column>
    </let-table>

      <div style="overflow:hidden;">
        <let-pagination align="right" style="float:right;"
          :page="pagination.page" @change="gotoPage"
          :total="pagination.total">
        </let-pagination>
      </div>
    </div>

    <!-- 服务实时状态 -->
    <let-table v-if="serverNotifyList && showOthers"
      :data="serverNotifyList" :title="$t('serverList.title.serverStatus')" :empty-msg="$t('common.nodata')" ref="serverNotifyListLoading">
      <let-table-column :title="$t('common.time')" prop="NotifyTime" width="160px"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.serviceID')" prop="AppServer">
        <template slot-scope="scope">
          <span style="white-space:nowrap">{{scope.row.AppServer}}</span>
        </template>
      </let-table-column>
      <let-table-column :title="$t('serverList.table.th.podName')" prop="PodName">
        <template slot-scope="scope">
          <span style="white-space:nowrap">{{scope.row.PodName}}</span>
        </template>
      </let-table-column>
      <let-table-column :title="$t('serverList.table.th.source')" prop="NotifySource">
        <template slot-scope="scope">
          <span style="white-space:nowrap">{{scope.row.NotifySource}}</span>
        </template>
      </let-table-column>
      <let-table-column :title="$t('serverList.table.th.result')" prop="NotifyMessage">
        <template slot-scope="scope">
          <div :class="getServerNotifyLevel(scope.row.NotifyLevel)">{{ scope.row.NotifyMessage }}</div>
        </template>
      </let-table-column>
    </let-table>

    <div style="overflow:hidden;margin-bottom:20px;">
      <let-pagination align="right" style="float:right;"
        :page="notifyPagination.page" @change="notifyGotoPage"
        :total="notifyPagination.total">
      </let-pagination>
    </div>

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
          <let-input
            size="small"
            v-model="configModal.model.AsyncThread"
            :placeholder="$t('serverList.dlg.placeholder.thread')"
            required
            :pattern="configModal.model.ServerTemplate === 'tars.nodejs' ? '^[1-9][0-9]*$' : '^([3-9]|[1-9][0-9]+)$'"
            pattern-tip="$t('serverList.dlg.placeholder.thread')"
          ></let-input>
        </let-form-item>
        <!--
        <let-form-item :label="$t('serverList.dlg.startScript')">
          <let-input
            size="small"
            v-model="configModal.model.StartScript"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.dlg.stopScript')">
          <let-input
            size="small"
            v-model="configModal.model.StopScript"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.dlg.monitorScript')" itemWidth="724px">
          <let-input
            size="small"
            v-model="configModal.model.MonitorScript"
          ></let-input>
        </let-form-item>
        -->
        <let-form-item :label="$t('serverList.dlg.privateTemplate')" labelWidth="150px" itemWidth="724px">
          <let-input
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
        <let-table-column :title="$t('deployService.table.th.protocol')" width="150px">
          <template slot-scope="props">
            <let-radio v-model="props.row.IsTars" :label="true" v-if="props.row.IsTars">TARS</let-radio>
            <let-radio v-model="props.row.IsTars" :label="false" v-else>{{$t('serverList.servant.notTARS')}}</let-radio>
          </template>
        </let-table-column>
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

    <!-- Servant新增弹窗 -->
    <let-modal
      v-model="servantAddModal.show"
      :title="servantAddModal.isNew ? `${$t('operate.title.add')} Servant` : `${$t('operate.title.update')} Servant`"
      width="1200px"
      :footShow="!!servantAddModal.model"
      @on-confirm="saveServantAdd"
      @close="closeServantAddModal"
      @on-cancel="closeServantAddModal">
      <let-form ref="servantAddForm" v-if="servantAddModal.model && servantAddModal.model.ServerServant">
        <let-table :data="servantAddModal.model.ServerServant">
          <let-table-column title="OBJ" width="150px">
            <template slot="head" slot-scope="props">
              <span class="required">{{props.column.title}}</span>
            </template>
            <template slot-scope="props">
              <let-input
                size="small"
                v-model="props.row.Name"
                :placeholder="$t('deployService.form.placeholder')"
                required
                :required-tip="$t('deployService.form.objTips')"
                pattern="^[a-zA-Z0-9]+$"
                :pattern-tip="$t('deployService.form.placeholder')"
              ></let-input>
            </template>
          </let-table-column>
          <let-table-column :title="$t('deployService.table.th.port')" width="100px">
          <template slot="head" slot-scope="props">
            <span class="required">{{props.column.title}}</span>
          </template>
          <template slot-scope="props">
            <let-input
              size="small"
              type="number"
              :min="1"
              :max="30000"
              v-model="props.row.Port"
              placeholder="1-30000"
              required
              :required-tip="$t('deployService.table.tips.empty')"
            ></let-input>
          </template>
        </let-table-column>
        <let-table-column :title="$t('deployService.form.portType')" width="150px">
          <template slot="head" slot-scope="props">
            <span class="required">{{props.column.title}}</span>
          </template>
          <template slot-scope="props">
            <let-radio v-model="props.row.IsTcp" :label="true">TCP</let-radio>
            <let-radio v-model="props.row.IsTcp" :label="false">UDP</let-radio>
          </template>
        </let-table-column>
        <let-table-column :title="$t('deployService.table.th.protocol')" width="180px">
          <template slot="head" slot-scope="props">
            <span class="required">{{props.column.title}}</span>
          </template>
          <template slot-scope="props">
            <let-radio v-model="props.row.IsTars" :label="true">TARS</let-radio>
            <let-radio v-model="props.row.IsTars" :label="false">{{$t('serverList.servant.notTARS')}}</let-radio>
          </template>
        </let-table-column>
        <let-table-column :title="$t('deployService.table.th.threads')" width="80px">
          <template slot="head" slot-scope="props">
            <span class="required">{{props.column.title}}</span>
          </template>
          <template slot-scope="props">
            <let-input
              size="small"
              type="number"
              :min="0"
              v-model="props.row.Threads"
              required
              :required-tip="$t('deployService.table.tips.empty')"
            ></let-input>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.servant.connections')" width="140px">
          <template slot="head" slot-scope="props">
            <span class="required">{{props.column.title}}</span>
          </template>
          <template slot-scope="props">
            <let-input
              size="small"
              type="number"
              :min="0"
              v-model="props.row.Connections"
              required
              :required-tip="$t('deployService.table.tips.empty')"
            ></let-input>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.servant.capacity')" width="140px">
          <template slot="head" slot-scope="props">
            <span class="required">{{props.column.title}}</span>
          </template>
          <template slot-scope="props">
            <let-input
              size="small"
              type="number"
              :min="0"
              v-model="props.row.Capacity"
              required
              :required-tip="$t('deployService.table.tips.empty')"
            ></let-input>
          </template>
        </let-table-column>
        <let-table-column :title="$t('serverList.table.servant.timeout')" width="140px">
          <template slot-scope="props">
            <let-input
              size="small"
              type="number"
              :min="0"
              v-model="props.row.Timeout"
            ></let-input>
          </template>
        </let-table-column>
        <let-table-column :title="$t('operate.operates')" width="60px">
          <template slot-scope="props">
            <let-table-operation @click="addAdapter(props.row)" v-if="props.$index === servantAddModal.model.ServerServant.length - 1">{{$t('operate.add')}}</let-table-operation>
            <let-table-operation v-else
               class="danger"
              @click="servantAddModal.model.ServerServant.splice(props.$index, 1)"
            >{{$t('operate.delete')}}</let-table-operation>
          </template>
        </let-table-column>
        </let-table>
      </let-form>
    </let-modal>

    <!-- Servant编辑弹窗 -->
    <let-modal
      v-model="servantDetailModal.show"
      :title="servantDetailModal.isNew ? `${$t('operate.title.add')} Servant` : `${$t('operate.title.update')} Servant`"
      width="800px"
      :footShow="!!servantDetailModal.model"
      @on-confirm="saveServantDetail"
      @close="closeServantDetailModal"
      @on-cancel="closeServantDetailModal">
      <let-form
        v-if="servantDetailModal.model && servantDetailModal.model.ServerServant.length === 1"
        ref="servantDetailForm" itemWidth="360px" :columns="2" class="two-columns">
        <let-form-item :label="$t('serverList.servant.objName')" required>
          <let-input
            size="small"
            v-model="servantDetailModal.model.ServerServant[0].Name"
            :placeholder="$t('serverList.servant.c')"
            required
            pattern="^[A-Za-z]+$"
            :pattern-tip="$t('serverList.servant.obj')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('deployService.table.th.port')" required>
          <let-input
            size="small"
            type="number"
            :min="1"
            :max="30000"
            v-model="servantDetailModal.model.ServerServant[0].Port"
            placeholder="1-30000"
            required
            :required-tip="$t('deployService.table.tips.empty')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.numOfThread')" required>
          <let-input
            size="small"
            type="number"
            v-model="servantDetailModal.model.ServerServant[0].Threads"
            :placeholder="$t('serverList.servant.thread')"
            required
            pattern="^[1-9][0-9]*$"
            :pattern-tip="$t('serverList.servant.thread')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.maxConnecttions')" labelWidth="150px">
          <let-input
            size="small"
            type="number"
            v-model="servantDetailModal.model.ServerServant[0].Connections"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.lengthOfQueue')" labelWidth="150px">
          <let-input
            size="small"
            type="number"
            v-model="servantDetailModal.model.ServerServant[0].Capacity"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.queueTimeout')" labelWidth="150px">
          <let-input
            size="small"
            type="number"
            v-model="servantDetailModal.model.ServerServant[0].Timeout"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('serverList.servant.protocol')" required>
          <let-radio v-model="servantDetailModal.model.ServerServant[0].IsTars" :label="true">TARS</let-radio>
          <let-radio v-model="servantDetailModal.model.ServerServant[0].IsTars" :label="false">{{ $t('serverList.servant.notTARS') }}</let-radio>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.portType')" required>
          <let-radio v-model="servantDetailModal.model.ServerServant[0].IsTcp" :label="true">TCP</let-radio>
          <let-radio v-model="servantDetailModal.model.ServerServant[0].IsTcp" :label="false">UDP</let-radio>
        </let-form-item>
      </let-form>
    </let-modal>

    <!-- K8S编辑弹窗 -->
    <let-modal
      v-model="k8sDetailModel.show"
      :title="$t('operate.k8s')"
      width="800px"
      :footShow="!!k8sDetailModel.model"
      @on-confirm="savek8sDetail"
      @close="closek8sDetailModel"
      @on-cancel="closek8sDetailModel">
      <let-form 
        v-if="k8sDetailModel.model"
        ref="k8sDetailForm" itemWidth="360px" :columns="2" class="two-columns">
        <let-form-item :label="$t('deployService.table.th.nodeSelector')" itemWidth="50%">
          <let-select v-if="K8SNodeSelectorKind && K8SNodeSelectorKind.length > 0"
            size="small"
            v-model="k8sDetailModel.model.NodeSelector.Kind"
            :placeholder="$t('pub.dlg.defaultValue')"
            required
            :required-tip="$t('deployService.form.templateTips')"
          >
            <let-option v-for="d in K8SNodeSelectorKind" :key="d" :value="d">{{d}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('deployService.table.th.replicas')">
          <let-input
            size="small"
            type="number"
            :min="0"
            v-model="k8sDetailModel.model.Replicas"
            :placeholder="$t('deployService.form.placeholder')"
            required
            :required-tip="$t('deployService.table.tips.empty')"
            :pattern-tip="$t('deployService.form.placeholder')"
          ></let-input>
        </let-form-item>
        <div v-if="k8sDetailModel.model.NodeSelector.Kind === 'NodeBind'">
          <div>
            <let-form-item :label="$t('deployService.table.th.hostIpc')" itemWidth="25%">
              <let-radio v-model="k8sDetailModel.model.HostIpc" :label="true">{{ $t('common.true') }}</let-radio>
              <let-radio v-model="k8sDetailModel.model.HostIpc" :label="false">{{ $t('common.false') }}</let-radio>
            </let-form-item>
            <let-form-item :label="$t('deployService.table.th.hostNetwork')" itemWidth="25%">
              <let-radio v-model="k8sDetailModel.model.HostNetwork" :label="true">{{ $t('common.true') }}</let-radio>
              <let-radio v-model="k8sDetailModel.model.HostNetwork" :label="false">{{ $t('common.false') }}</let-radio>
            </let-form-item>
            <let-form-item :label="$t('deployService.table.th.hostPort')" itemWidth="25%">
              <let-radio v-model="k8sDetailModel.model.showHostPort" :label="true">{{ $t('common.true') }}</let-radio>
              <let-radio v-model="k8sDetailModel.model.showHostPort" :label="false">{{ $t('common.false') }}</let-radio>
            </let-form-item>
          </div>
          <div v-if="k8sDetailModel.model.showHostPort" style="padding-right:30px;">
            <let-table :data="k8sDetailModel.model.HostPortArr">
              <let-table-column title="OBJ">
                <template slot="head" slot-scope="props">
                  <span class="required">{{props.column.title}}</span>
                </template>
                <template slot-scope="props">
                  <let-input
                    size="small"
                    v-model="props.row.obj"
                    :placeholder="$t('deployService.form.placeholder')"
                    required
                    :required-tip="$t('deployService.form.objTips')"
                    pattern="^[a-zA-Z0-9]+$"
                    :pattern-tip="$t('deployService.form.placeholder')"
                  ></let-input>
                </template>
              </let-table-column>
              <let-table-column :title="$t('deployService.table.th.hostPort')">
                <template slot="head" slot-scope="props">
                  <span class="required">{{props.column.title}}</span>
                </template>
                <template slot-scope="props">
                  <let-input size="small" type="number" :min="1" :max="30000"
                             v-model="props.row.HostPort" placeholder="1-30000"
                             required :required-tip="$t('deployService.table.tips.empty')"
                  ></let-input>
                </template>
              </let-table-column>
              <let-table-column>
                <template slot-scope="props">
                  <let-button size="small" theme="primary" class="port-button" @click="generateHostPort(props.row)">
                    {{$t('deployService.table.th.generateHostPort')}}
                  </let-button>
                </template>
              </let-table-column>
            </let-table>
          </div>
        </div>
        <div>
          <let-form-item itemWidth="100%" v-if="k8sDetailModel.model.NodeSelector.Kind === 'NodeBind'">
            <let-checkbox v-model="K8SisCheckedAll" :value="K8SisCheckedAll">{{ $t('cache.config.allSelected') }}</let-checkbox>
            <div class="node_list">
              <let-checkbox class="node_item" v-for="d in K8SNodeList" :key="d" :value="K8SNodeListArr.indexOf(d) > -1" @change="K8ScheckedChange(d)">{{ d }}</let-checkbox>
            </div>
          </let-form-item>
        </div>
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
        <!-- <let-form-item itemWidth="100%">
          <let-radio v-model="moreCmdModal.model.selected" label="undeploy_tars" class="danger">{{$t('operate.undeploy')}} {{$t('common.service')}}</let-radio>
        </let-form-item> -->
      </let-form>
    </let-modal>

    <!-- 查看弹窗 -->
    <let-modal
      v-model="detailModal.show"
      :title="detailModal.title"
      width="700px"
      :footShow="false"
      @close="closeDetailModal">
      <div style="padding:20px 0 0;">
        <pre v-if="(detailModal.model && detailModal.model.detail)">{{detailModal.model.detail || $t('cfg.msg.empty')}}</pre>
        <div class="detail-loading" ref="detailModalLoading"></div>
      </div>
    </let-modal>

  </div>
</template>

<script>
import wrapper from '@/components/section-wrappper';
export default {
  name: 'ServerManage',
  components: {
    wrapper,
  },
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

      // 分页
      pagination: {
        page: 1,
        size: 10,
        total:1,
      },

      // 操作历史列表
      serverNotifyList: [],

      // 分页
      notifyPagination: {
        page: 1,
        size: 10,
        total:1,
      },

      // 默认参数
      defaultObj: {},

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

      // 新增servant
      servantAddModal: {
        show: false,
        isNew: true,
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
      
      K8SNodeSelectorKind: [],
      K8SisCheckedAll: false,
      K8SNodeList: [],
      K8SNodeListArr: [],

      // 查看弹窗
      detailModal: {
        show: false,
        title: '',
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
      failCount :0,

      isCheckedAll: false,
      checkedList: [],
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
        let url = `${domain.fromPage}?History=false&NodeIP=${data.NodeIp}&AppName=${data.ServerApp}&ServerName=${data.ServerName}&PodName=${data.PodName}`
        window.open(url)
      })
    },
    getServerId() {
      return this.treeid
    },
    getDefaultValue(){
      let { K8SNodeSelectorKind } = this
      this.$ajax.getJSON('/server/api/default', {}).then((data) => {
        this.defaultObj = data.ServerServantElem
        if(data.K8SNodeSelectorKind){
          K8SNodeSelectorKind = data.K8SNodeSelectorKind
        }
        this.K8SNodeSelectorKind = K8SNodeSelectorKind
      })
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
      if(path === 'manage'){
        that.reloadTask = setTimeout(() => {
          if(that.isreloadlist){
            that.getServerList()
          }
          that.reloadServerList()
        }, 5000)
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
      // const loading = this.$refs.serverListLoading.$loading.show();
      this.$ajax.getJSON('/server/api/pod_list', {
        ServerId: this.getServerId(),
        page: this.pagination.page,
        size: this.pagination.size,
      }).then((data) => {
        // loading.hide();
        this.serverList = []
        if (data.hasOwnProperty("Data")) {
          data.Data.forEach(item => {
            item.isChecked = false
          })
          this.serverList = data.Data
        }
        this.pagination.total = Math.ceil(data.Count.FilterCount / this.pagination.size)
      }).catch((err) => {
        // loading.hide();
        this.stopServerList();
        this.$confirm(err.err_msg || err.message || this.$t('serverList.msg.fail'), this.$t('common.alert')).then(() => {
          this.getServerList();
        });
      });
    },
    // 获取服务实时状态
    getServerNotifyList() {
      if (!this.showOthers) return;
      // const loading = this.$refs.serverNotifyListLoading.$loading.show();

      this.$ajax.getJSON('/server/api/server_notify_list', {
        ServerId: this.getServerId(),
        size: this.notifyPagination.size,
        page: this.notifyPagination.page,
      }).then((data) => {
        // loading.hide();
        this.notifyPagination.total = Math.ceil(data.Count.FilterCount/this.notifyPagination.size);
        this.serverNotifyList = data.Data;
      }).catch((err) => {
        // loading.hide();
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },
    // 获取服务通知级别
    getServerNotifyLevel(str) {
      return str && str.substring(1, str.length - 1)
    },
    // 切换服务实时状态页码
    gotoPage(num) {
      this.pagination.page = num
      this.getServerList()
    },
    notifyGotoPage(num) {
      this.notifyPagination.page = num
      this.getServerNotifyList()
    },
    // 获取服务数据
    getServerConfig(id) {
      this.stopServerList()
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
        this.startServerList()
      }).catch((err) => {
        loading.hide();
        this.closeConfigModal();
        this.startServerList()
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },
    // 编辑服务
    configServer() {
      this.stopServerList()
      this.configModal.show = true;
      
      this.$ajax.getJSON('/server/api/template_select').then((data) => {
        if (this.configModal.model) {
          this.configModal.model.templates = data.Data;
        } else {
          this.configModal.model = { templates: data.Data };
        }
        this.getServerConfig(this.getServerId());
        this.startServerList()
      }).catch((err) => {
        this.startServerList()
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
        }).catch((err) => {
          loading.hide();
          if(err.ret_code === -2){
            this.$confirm(err.err_msg, this.$t('common.alert')).then(() => {
              this.$ajax.postJSON('/server/api/server_option_update', {
                Confirmation: true,
                isBak: this.configModal.model.bak_flag,
                ...this.configModal.model,
              }).then((res) => {
                this.serverList = this.serverList.map((item) => {
                  if (item.id === res.id) {
                    return res;
                  }
                  return item;
                });
                this.closeConfigModal();
                this.$tip.success(this.$t('common.success'));
              }).catch((err) => {
                this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
              })
            });
          }else{
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          }
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
    // 启动服务
    startServer() {
      this.stopServerList();
      const checkedServerList = this.checkedList.filter(item => item.isChecked);
      if (checkedServerList.length <= 0) {
        this.$tip.warning(this.$t('pub.dlg.a'));
        return;
      }

      let podIp = checkedServerList.map(item => item.PodIp)
      let serverApp = checkedServerList[0].ServerApp
      let serverName = checkedServerList[0].ServerName
      this.$confirm(this.$t('serverList.startService.msg.startService'), this.$t('common.alert')).then(() => {
        this.$ajax.getJSON('/server/api/send_command', {
          command: 'StartServer',
          podIp,
          serverApp,
          serverName,
        }).then((data) => {
          this.$tip.success(this.$t('common.success'));
          this.startServerList()
          this.getServerNotifyList()
        }).catch((err) => {
          this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
        });
      }).catch(() => {
        this.startServerList()
      });
    },
    // 重启服务
    restartServer() {
      this.stopServerList();
      const checkedServerList = this.checkedList.filter(item => item.isChecked);
      if (checkedServerList.length <= 0) {
        this.$tip.warning(this.$t('pub.dlg.a'));
        return;
      }

      let podIp = checkedServerList.map(item => item.PodIp)
      let serverApp = checkedServerList[0].ServerApp
      let serverName = checkedServerList[0].ServerName
      this.$confirm(this.$t('serverList.restartService.msg.restartService'), this.$t('common.alert')).then(() => {
        this.$ajax.getJSON('/server/api/send_command', {
          command: 'RestartServer',
          podIp,
          serverApp,
          serverName,
        }).then((data) => {
          this.$tip.success(this.$t('common.success'));
          this.startServerList()
          this.getServerNotifyList()
        }).catch((err) => {
          this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
        });
      }).catch(() => {
        this.startServerList()
      });
    },
    // 停止服务
    stopServer() {
      this.stopServerList();
      const checkedServerList = this.checkedList.filter(item => item.isChecked);
      if (checkedServerList.length <= 0) {
        this.$tip.warning(this.$t('pub.dlg.a'));
        return;
      }

      let podIp = checkedServerList.map(item => item.PodIp)
      let serverApp = checkedServerList[0].ServerApp
      let serverName = checkedServerList[0].ServerName
      this.$confirm(this.$t('serverList.stopService.msg.stopService'), this.$t('common.alert')).then(() => {
        this.$ajax.getJSON('/server/api/send_command', {
          command: 'StopServer',
          podIp,
          serverApp,
          serverName,
        }).then((data) => {
          this.$tip.success(this.$t('common.success'));
          this.startServerList()
          this.getServerNotifyList()
        }).catch((err) => {
          this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
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

    // 查看模版
    viewTemplate() {
      const loading = this.$loading.show({
        target: this.$refs.servantModalLoading,
      });
      this.$ajax.getJSON('/server/api/server_option_template', {
        ServerId: this.getServerId(),
      }).then((data) => {
        loading.hide();
        this.detailModal.title = this.$t('cfg.title.viewTemplate');
        this.detailModal.model = {
          detail: data,
        };
        this.detailModal.show = true;
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
      });
    },

    // 管理私有模版
    privateTemplateManage() {
      const loading = this.$loading.show({
        target: this.$refs.servantModalLoading,
      });
      this.$ajax.getJSON('/server/api/server_option_select', {
        ServerId: this.getServerId(),
      }).then((data) => {
        loading.hide();
        this.detailModal.title = this.$t('cfg.title.viewProfileTemplate');
        this.detailModal.model = {
          detail: data && data[0].ServerProfile || '',
        };
        this.detailModal.show = true;
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
      });
    },

    // 管理Servant弹窗
    manageServant(server) {
      this.stopServerList()
      this.servantModal.show = true;

      this.$ajax.getJSON('/server/api/server_adapter_select', {
        ServerId: this.getServerId(),
      }).then((data) => {
        this.servantModal.model = data.Data;
        this.servantModal.currentServer = server;
      }).catch((err) => {
        this.$tip.error(`${this.$t('serverList.restart.failed')}: ${err.err_msg || err.message}`);
      });
    },

    // 管理K8S弹窗
    manageK8S(server) {
      this.stopServerList()

      const loading = this.$loading.show({
        target: this.$refs.k8sDetailModelLoading,
      });

      this.$ajax.getJSON('/server/api/server_k8s_select', {
        ServerId: this.getServerId(),
      }).then((data) => {
        loading.hide();
        data[0].showHostPort = false
        if(data && data[0]){
          let arr = []
          if(data[0].HostPort && Object.keys(data[0].HostPort).length > 0){
            data[0].showHostPort = true
            // Object -> Array
            data[0].HostPort.forEach(item => {
              arr.push({
                obj: item.nameRef,
                HostPort: Math.floor(item.port),
              })
            })
            data[0].HostPortArr = arr
          }
          if(arr.length === 0){
            this.$ajax.getJSON('/server/api/server_adapter_select', {
              ServerId: this.getServerId(),
            }).then((retdata) => {
              if(retdata && retdata.Data){
                retdata.Data.forEach(item => {
                  arr.push({
                    obj: item.Name,
                    HostPort: 0,
                  })
                })
                data[0].HostPortArr = arr
              }
            })
          }
        }

        // 兼容新接口
        if (data[0].NodeSelector.hasOwnProperty("abilityPool")) {
          data[0].NodeSelector.Kind = "AbilityPool"
          data[0].NodeSelector.Value = data[0].NodeSelector.abilityPool.values
        } else if (data[0].NodeSelector.hasOwnProperty("nodeBind")) {
          data[0].NodeSelector.Kind = "NodeBind"
          data[0].NodeSelector.Value = data[0].NodeSelector.nodeBind.values
        } else if (data[0].NodeSelector.hasOwnProperty("publishPool")) {
          data[0].NodeSelector.Kind = "PublishPool"
          data[0].NodeSelector.Value = data[0].NodeSelector.publishPool.values
        }
        this.k8sDetailModel.model = data[0];
        this.K8SNodeListArr = data[0].NodeSelector.Value
        this.k8sDetailModel.show = true;
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
    addAdapter(template) {
      this.servantAddModal.model.ServerServant.push(Object.assign({}, template, { Port: this.getPort(template.Port) }));
    },
    getPort(port) {
      const { defaultObj, servantModal } = this

      let result = 0
      if(!port){
        result = Math.floor(defaultObj.Port)

        if(servantModal.model && servantModal.model.length > 0){
          let servant = servantModal.model.sort((a, b) => {
            return b.Port - a.Port
          }) || []
          result = Math.floor(servant[0].Port)
        }
      }else{
        result = port
      }
      result++

      if(result === 19385){
        result++
      }

      return result
    },
    // 新增、编辑 servant
    configServant(id) {
      const { defaultObj } = this

      const model = {
        ServerId: this.getServerId(),
        ServerServant: [
          {
            Name: '',
            Port: this.getPort() || 0,
            Threads: defaultObj.Threads || 0,
            Connections: defaultObj.Connections || 0,
            Capacity: defaultObj.Capacity || 0,
            Timeout: defaultObj.Timeout || 0,
            IsTcp: defaultObj.IsTcp,
            IsTars: defaultObj.IsTars,
          }
        ],
      }
      if(id){
        // 编辑
        this.servantDetailModal.model = model
        this.servantDetailModal.isNew = true;

        const old = this.servantModal.model.find(item => item.AdapterId === id);
        // old.obj_name = old.servant.split('.')[2];
        // this.servantDetailModal.model = Object.assign({}, this.servantDetailModal.model, old);
        this.servantDetailModal.model = {
          AdapterId: old.AdapterId,
          ServerServant: [old],
        };
        this.servantDetailModal.isNew = false;
        this.servantDetailModal.show = true;
      }else{
        // 新增
        this.servantAddModal.model = model
        this.servantAddModal.isNew = true;
        this.servantAddModal.show = true;
      }
    },
    closeServantAddModal() {
      if (this.$refs.servantAddForm) this.$refs.servantAddForm.resetValid();
      this.servantAddModal.show = false;
      this.servantAddModal.model = null;
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
    saveServantAdd() {
      if (this.$refs.servantAddForm.validate()) {
        const loading = this.$Loading.show();
        // 新建
        if (this.servantAddModal.isNew) {
          let query = Object.assign({}, this.servantAddModal.model)

          if(query.ServerServant) {
            // Array -> Object
            let obj = {}
            query.ServerServant.forEach(item => {
              obj[item.Name] = item
            })
            query.ServerServant = obj
          }

          this.$ajax.postJSON('/server/api/server_adapter_create', query).then((res) => {
            loading.hide();
            // this.servantModal.model.unshift(res);
            this.manageServant()
            this.$tip.success(this.$t('common.success'));
            this.closeServantAddModal();
          }).catch((err) => {
            loading.hide();
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          });
        }else{
          loading.hide();
        }
      }
    },
    // 保存 servant
    saveServantDetail() {
      if (this.$refs.servantDetailForm.validate()) {
        const loading = this.$Loading.show();
        // 新建
        if (this.servantDetailModal.isNew) {
          const query = this.servantDetailModal.model;
          this.$ajax.postJSON('/server/api/server_adapter_create', query).then((res) => {
            loading.hide();
            this.servantModal.model.unshift(res);
            this.$tip.success(this.$t('common.success'));
            // this.closeServantModal();
            this.closeServantDetailModal();
          }).catch((err) => {
            loading.hide();
            this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
          });
        // 修改
        } else {
          this.servantDetailModal.model.servant=this.servantDetailModal.model.application+'.'+this.servantDetailModal.model.server_name+'.'+this.servantDetailModal.model.obj_name;

          let query = Object.assign({}, this.servantDetailModal.model)

          // if(query.ServerServant) {
          //   // Array -> Object
          //   let obj = {}
          //   query.ServerServant.forEach(item => {
          //     obj[item.Name] = item
          //   })
          //   query.ServerServant = obj
          // }

          this.$ajax.postJSON('/server/api/server_adapter_update', query).then((res) => {
            loading.hide();
          }).catch((err) => {
            loading.hide();
            if(err.ret_code === -2){
              this.$confirm(err.err_msg, this.$t('common.alert')).then(() => {
                this.$ajax.postJSON('/server/api/server_adapter_update', {
                  Confirmation: true,
                ...query,
                }).then((res) => {
                  loading.hide();
                  // this.servantModal.model = this.servantModal.model.map((item) => {
                  //   if (item.id === res.id) {
                  //     return res;
                  //   }
                  //   return item;
                  // });
                  this.$tip.success(this.$t('common.success'));
                  // this.closeServantModal();
                  this.closeServantDetailModal();
                }).catch((err) => {
                  this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
                });
              })
            }else{
              this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
            }
          });
        }
      }
    },
    adapterServerK8S(model) {
      let data = Object.assign({}, model)

      data.HostPort = []

      if(data.NodeSelector.Kind === 'NodeBind'){
        if(data.HostNetwork && data.showHostPort){
          return this.$tip.error(`${this.$t('deployService.form.portOrNetWork')}`)
        }

        if(data.showHostPort){
          data.HostPortArr && data.HostPortArr.forEach(item => {
            data.HostPort.push({
              "NameRef": item.obj,
              "Port": Math.floor(item.HostPort)
            })
          })
        }

        // 兼容新后台接口结构
        delete data.NodeSelector.abilityPool
        if(this.K8SNodeListArr.length <= 0){
          return this.$tip.error(`${this.$t('deployService.form.nodeTips')}`)
        }
        data.NodeSelector.NodeBind = {
          Value: this.K8SNodeListArr
        }
      } else if (data.NodeSelector.Kind === 'AbilityPool') {
        // 兼容新后台接口结构
        delete data.NodeSelector.nodeBind
        data.NodeSelector.AbilityPool = {
          Value: []
        }
      }

      return data
    },
    // 保存 k8s
    savek8sDetail() {
      if (this.$refs.k8sDetailForm.validate()) {
        const loading = this.$Loading.show();
        let data = this.adapterServerK8S(this.k8sDetailModel.model)
        // 新建
        if (this.k8sDetailModel.isNew) {
          this.$ajax.postJSON('/server/api/server_k8s_create', data).then((res) => {
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
          this.$ajax.postJSON('/server/api/server_k8s_update', data).then((res) => {
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
    // 自动生成端口
    generateHostPort(hostPort) {
      if (this.K8SNodeListArr.length <= 0) {
        alert("请优先选择需要绑定的节点！")
        return
      }
      this.$ajax.getJSON("/server/api/generate_host_port", {
        NodeList: this.K8SNodeListArr,
        NodePort: hostPort.HostPort
      }).then((res) => {
        hostPort.HostPort = res.port
      })
    },

    // 删除 servant
    deleteServant(id) {
      this.$confirm(this.$t('serverList.servant.a'), this.$t('common.alert')).then(() => {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/server_adapter_delete', {
          AdapterId: id,
        }).then((res) => {
          loading.hide();
          this.servantModal.model = this.servantModal.model.filter( item => item.AdapterId !== id);
          this.$tip.success(this.$t('common.success'));
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.err_msg || err.message}`);
        });
      });
    },

    closeDetailModal() {
      this.detailModal.show = false;
      this.detailModal.model = null;
    },

    // 显示更多命令
    showMoreCmd(server) {
      this.stopServerList();
      const checkedServerList = this.checkedList.filter(item => item.isChecked);
      if (checkedServerList.length <= 0) {
        this.$tip.warning(this.$t('pub.dlg.a'));
        return;
      }

      let podIp = checkedServerList.map(item => item.PodIp)
      let serverApp = checkedServerList[0].ServerApp
      let serverName = checkedServerList[0].ServerName

      this.moreCmdModal.model = {
        podIp,
        serverApp,
        serverName,
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
    sendCommand(serverApp, serverName, podIp, command, hold) {
      const loading = this.$Loading.show();
      this.$ajax.getJSON('/server/api/send_command', {
        serverApp,
        serverName,
        podIp,
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
          this.getServerNotifyList()
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
        this.undeployServer(model.serverApp, model.serverName, model.podIp);
      // 设置日志等级
      } else if (model.selected === 'setloglevel') {
        this.sendCommand(model.serverApp, model.serverName, model.podIp, `tars.setloglevel ${model.setloglevel}`);
      // push 日志文件
      } else if (model.selected === 'loadconfig' && this.$refs.moreCmdForm.validate()) {
        this.sendCommand(model.serverApp, model.serverName, model.podIp, `tars.loadconfig ${model.loadconfig}`);
      // 发送自定义命令
      } else if (model.selected === 'command' && this.$refs.moreCmdForm.validate()) {
        this.sendCommand(model.serverApp, model.serverName, model.podIp, model.command);
      // 查看服务链接
      } else if (model.selected === 'connection') {
        this.sendCommand(model.serverApp, model.serverName, model.podIp, `tars.connection`, true);
      }
    },
    closeMoreCmdModal() {
      if (this.$refs.moreCmdForm) this.$refs.moreCmdForm.resetValid();
      if (this.moreCmdModal.unwatch) this.moreCmdModal.unwatch();
      this.checkedList = []
      this.moreCmdModal.show = false;
      this.moreCmdModal.model = null;
      this.startServerList()
    },

    // 处理未发布时间显示
    handleNoPublishedTime(timeStr, noPubTip = this.$t('pub.dlg.unpublished')) {
      if (timeStr === '0000:00:00 00:00:00') {
        return noPubTip;
      }
      return timeStr;
    },

    checkChange(data) {
      debugger
      this.stopServerList()
      const checkedList = this.checkedList;

      let isChecked = data.isChecked
      if(isChecked){
        let isTrue = false
        checkedList.forEach(item => {
          if(item.PodId === data.PodId){
            isTrue = true
          }
        })
        if(!isTrue){
          checkedList.push(data)
        }
      }else{
        let isTrue = false
        let index = -1
        checkedList.forEach((iitem, iindex) => {
          if(iitem.PodId === data.PodId){
            isTrue = true
            index = iindex
          }
        })
        if(isTrue){
          checkedList.splice(index, 1)
        }
      }
    },

    K8ScheckedChange(val) {
      if(this.K8SNodeListArr.indexOf(val) > -1){
        this.K8SNodeListArr.splice(this.K8SNodeListArr.indexOf(val), 1)
      }else{
        this.K8SNodeListArr.push(val)
      }

      if(this.K8SNodeListArr.length === 0){
        this.K8SisCheckedAll = false
      }
    },

    getNodeList() {
      this.$ajax.getJSON('/server/api/node_list', { isAll: true }).then((data) => {
        this.K8SNodeList = data.Data
      })
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
    this.getNodeList();
    this.getServerList();
    this.reloadServerList()
    this.startServerList()
    this.getServerNotifyList();
    this.getDefaultValue()
  },
  watch: {
    isCheckedAll() {
      const isCheckedAll = this.isCheckedAll;
      this.serverList.forEach((item) => {
        item.isChecked = isCheckedAll;
      })

      if(isCheckedAll){
        this.checkedList = [].concat(this.serverList)
      }else{
        this.checkedList = []
      }
    },
    checkedList() {
      const checkedList = this.checkedList;
      if(checkedList.length > 0){
        this.stopServerList()
      }else{
        this.startServerList()
      }
    },
    K8SisCheckedAll() {
      let K8SisCheckedAll = this.K8SisCheckedAll;
      if(K8SisCheckedAll) {
        this.K8SNodeList.forEach(item => {
          if(this.K8SNodeListArr.indexOf(item) === -1){
            this.K8SNodeListArr.push(item)
          }
        })
      }else{
        this.K8SNodeListArr = []
      }
    },
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

  .warn { color:var(--warn-color) }
  .alarm { color:var(--alarm-color) }
  .error { color:var(--error-color) }
  .info { color:var(--info-color) }
  .normal { color:var(--normal-color) }

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

  .btn_group{position:absolute;right:0;top:0;z-index:2;}
  .btn_group .let-button + .let-button{margin-left:10px;}

  .port-button{margin-left: 10px; margin-right: 10px}

  .tooltip{display:block;font-size:0;line-height:1;}

  .let-form-item{vertical-align:top;}

  .node_list {
    border: 1px solid #e1e4eb;
    max-height: 200px;
    margin:10px 0;
    overflow: hidden;
    overflow-y: auto;
  }

  .node_item {
    box-sizing:border-box;
    padding:4px 10px;
    width:33%
  }
}

</style>
