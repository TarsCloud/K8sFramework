<template>
  <div class="page_operation_approval">
    <let-form inline itemWidth="200px" @submit.native.prevent="search">
      <let-form-item :label="$t('deployService.form.app')">
        <let-input size="small" v-model="query.ServerApp"></let-input>
      </let-form-item>
      <let-form-item :label="$t('deployService.form.serviceName')">
        <let-input size="small" v-model="query.ServerName"></let-input>
      </let-form-item>
      <let-form-item>
        <let-button size="small" type="submit" theme="primary">{{$t('operate.search')}}</let-button>
      </let-form-item>
    </let-form>

    <let-table ref="table" :data="items" :empty-msg="$t('common.nodata')">
      <let-table-column :title="$t('deployService.form.app')" prop="ServerApp" width="25%"></let-table-column>
      <let-table-column :title="$t('deployService.form.serviceName')" prop="ServerName" width="25%"></let-table-column>
      <let-table-column :title="$t('serviceApproval.RequestPerson')" prop="RequestPerson"></let-table-column>
      <!-- <let-table-column :title="$t('serviceApproval.ApprovalResult')" prop="ApprovalResult">
        <template slot-scope="scope">
          <span v-if="scope.row.ApprovalResult">{{ scope.row.ApprovalResult ? '通过' : '驳回' }}</span>
        </template>
      </let-table-column>
      <let-table-column :title="$t('serviceApproval.ApprovalMark')" prop="ApprovalMark"></let-table-column> -->
      <let-table-column :title="$t('operate.operates')" width="180px">
        <template slot-scope="scope">
          <let-table-operation @click="editItem(scope.row)">{{$t('operate.update')}}</let-table-operation>
          <let-table-operation @click="approvalItem(scope.row)">{{$t('operate.approval')}}</let-table-operation>
          <let-table-operation @click="removeItem(scope.row)">{{$t('operate.delete')}}</let-table-operation>
        </template>
      </let-table-column>
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

    <let-modal v-model="viewModal.show" :title="$t('template.view.title')" width="800px">
      <pre v-if="viewModal.model">{{viewModal.model}}</pre>
      <div slot="foot"></div>
    </let-modal>

    <let-modal
      v-model="detailModal.show"
      :title="detailModal.isNew ? this.$t('dialog.title.add') : this.$t('dialog.title.edit')"
      width="80%"
      @on-confirm="saveItem"
      @on-cancel="closeDetailModal"
    >
      <let-form ref="detailForm" inline v-if="detailModal.model">
        <let-form-item :label="$t('deployService.form.app')" itemWidth="45%" required>
          <let-input
            size="small"
            v-model="detailModal.model.ServerApp"
            :placeholder="$t('deployService.form.placeholder')"
            required
            :required-tip="$t('deployService.form.appTips')"
            pattern="^[a-zA-Z]+$"
            :pattern-tip="$t('deployService.form.placeholder')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.serviceName')" itemWidth="45%" required>
          <let-input
            size="small"
            v-model="detailModal.model.ServerName"
            :placeholder="$t('deployService.form.serviceFormatTips')"
            required
            :required-tip="$t('deployService.form.serviceTips')"
            pattern="^[a-zA-Z]([a-zA-Z0-9]+)?$"
            :pattern-tip="$t('deployService.form.serviceFormatTips')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.template')" itemWidth="45%" required>
          <let-select
            size="small"
            v-model="detailModal.model.ServerOption.ServerTemplate"
            required
            :required-tip="$t('deployService.form.templateTips')"
          >
            <let-option v-for="d in templates" :key="d.TemplateName" :value="d.TemplateName">{{d.TemplateName}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.serviceMark')" itemWidth="45%">
          <let-input
            size="small"
            v-model="detailModal.model.ServerMark"
            :placeholder="$t('deployService.form.serviceMark')"
          >
          </let-input>
        </let-form-item>
        <let-table :data="detailModal.model.ServerServant">
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
                :min="0"
                :max="65535"
                v-model="props.row.Port"
                placeholder="0-65535"
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
              <let-radio v-model="props.row.IsTaf" :label="true">TAF</let-radio>
              <let-radio v-model="props.row.IsTaf" :label="false">{{$t('serverList.servant.notTAF')}}</let-radio>
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
          <let-table-column :title="$t('serverList.table.servant.connections')" width="100px">
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
          <let-table-column :title="$t('serverList.table.servant.capacity')" width="120px">
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
          <let-table-column :title="$t('serverList.table.servant.timeout')" width="150px">
            <template slot-scope="props">
              <let-input
                size="small"
                type="number"
                :min="0"
                v-model="props.row.Timeout"
              ></let-input>
            </template>
          </let-table-column>
          <let-table-column width="10px"></let-table-column>
        </let-table>
        <let-form-group class="let-box" title="K8S" inline label-position="top">
          <let-form-item :label="$t('deployService.table.th.replicas')" itemWidth="30%">
            <let-input
              size="small"
              type="number"
              :min="0"
              v-model="detailModal.model.ServerK8S.Replicas"
              :placeholder="$t('deployService.form.placeholder')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
              :pattern-tip="$t('deployService.form.placeholder')"
            ></let-input>
          </let-form-item>
          <let-form-item :label="$t('deployService.table.th.netMode')" itemWidth="30%">
            <let-select v-if="K8SNetModeOptional && K8SNetModeOptional.length > 0"
              size="small"
              v-model="detailModal.model.ServerK8S.NetMode"
              :placeholder="$t('pub.dlg.defaultValue')"
              required
              :required-tip="$t('deployService.form.templateTips')"
            >
              <let-option v-for="d in K8SNetModeOptional" :key="d" :value="d">{{d}}</let-option>
            </let-select>
            <let-input v-else
              size="small"
              :min="0"
              v-model="detailModal.model.ServerK8S.NetMode"
              :placeholder="$t('deployService.form.placeholder')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
              :pattern-tip="$t('deployService.form.placeholder')"
            ></let-input>
          </let-form-item>
          <let-form-item :label="$t('deployService.table.th.nodeSelector')" itemWidth="30%">
            <let-select v-if="K8SNodeSelectorOptional && K8SNodeSelectorOptional.length > 0"
              size="small"
              v-model="detailModal.model.ServerK8S.NodeSelector"
              :placeholder="$t('pub.dlg.defaultValue')"
              required
              :required-tip="$t('deployService.form.templateTips')"
            >
              <let-option v-for="d in K8SNodeSelectorOptional" :key="d" :value="d">{{d}}</let-option>
            </let-select>
            <let-input v-else
              size="small" 
              :min="0"
              v-model="detailModal.model.ServerK8S.NodeSelector"
              :placeholder="$t('deployService.form.placeholder')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
              :pattern-tip="$t('deployService.form.placeholder')"
            ></let-input>
          </let-form-item>
        </let-form-group>
      </let-form>
    </let-modal>

    <let-modal
      v-model="approvalModal.show"
      :title="approvalModal.isNew ? this.$t('deployService.title.approval') : this.$t('deployService.title.approval')"
      width="80%"
      @on-confirm="saveApprovalItem"
      @on-cancel="closeApprovalModal"
    >
      <let-form ref="approvalForm" inline v-if="approvalModal.model">
        <let-form-item :label="$t('deployService.form.app')" itemWidth="45%" required>
          <let-input disabled
            size="small"
            v-model="approvalModal.model.ServerApp"
            :placeholder="$t('deployService.form.placeholder')"
            required
            :required-tip="$t('deployService.form.appTips')"
            pattern="^[a-zA-Z]+$"
            :pattern-tip="$t('deployService.form.placeholder')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.serviceName')" itemWidth="45%" required>
          <let-input disabled
            size="small"
            v-model="approvalModal.model.ServerName"
            :placeholder="$t('deployService.form.serviceFormatTips')"
            required
            :required-tip="$t('deployService.form.serviceTips')"
            pattern="^[a-zA-Z]([a-zA-Z0-9]+)?$"
            :pattern-tip="$t('deployService.form.serviceFormatTips')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.template')" itemWidth="45%" required>
          <let-select disabled
            size="small"
            v-model="approvalModal.model.ServerOption.ServerTemplate"
            required
            :required-tip="$t('deployService.form.templateTips')"
          >
            <let-option v-for="d in templates" :key="d.TemplateName" :value="d.TemplateName">{{d.TemplateName}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.serviceMark')" itemWidth="45%">
          <let-input disabled
            size="small"
            v-model="approvalModal.model.ServerMark"
            :placeholder="$t('deployService.form.serviceMark')"
          >
          </let-input>
        </let-form-item>
        <let-table :data="approvalModal.model.ServerServant">
          <let-table-column title="OBJ" width="150px">
            <template slot="head" slot-scope="props">
              <span class="required">{{props.column.title}}</span>
            </template>
            <template slot-scope="props">
              <let-input disabled
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
              <let-input disabled
                size="small"
                type="number"
                :min="0"
                :max="65535"
                v-model="props.row.Port"
                placeholder="0-65535"
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
              <let-radio disabled v-model="props.row.IsTcp" :label="true">TCP</let-radio>
              <let-radio disabled v-model="props.row.IsTcp" :label="false">UDP</let-radio>
            </template>
          </let-table-column>
          <let-table-column :title="$t('deployService.table.th.protocol')" width="180px">
            <template slot="head" slot-scope="props">
              <span class="required">{{props.column.title}}</span>
            </template>
            <template slot-scope="props">
              <let-radio disabled v-model="props.row.IsTaf" :label="true">TAF</let-radio>
              <let-radio disabled v-model="props.row.IsTaf" :label="false">{{$t('serverList.servant.notTAF')}}</let-radio>
            </template>
          </let-table-column>
          <let-table-column :title="$t('deployService.table.th.threads')" width="80px">
            <template slot="head" slot-scope="props">
              <span class="required">{{props.column.title}}</span>
            </template>
            <template slot-scope="props">
              <let-input disabled
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
              <let-input disabled
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
              <let-input disabled
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
              <let-input disabled
                size="small"
                type="number"
                :min="0"
                v-model="props.row.Timeout"
              ></let-input>
            </template>
          </let-table-column>
          <let-table-column width="10px"></let-table-column>
        </let-table>
        <let-form-group class="let-box" title="K8S" inline label-position="top">
          <let-form-item :label="$t('deployService.table.th.replicas')" itemWidth="30%">
            <let-input disabled
              size="small"
              type="number"
              :min="0"
              v-model="approvalModal.model.ServerK8S.Replicas"
              :placeholder="$t('deployService.form.placeholder')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
              :pattern-tip="$t('deployService.form.placeholder')"
            ></let-input>
          </let-form-item>
          <let-form-item :label="$t('deployService.table.th.netMode')" itemWidth="30%">
            <let-select disabled v-if="K8SNetModeOptional && K8SNetModeOptional.length > 0"
              size="small"
              v-model="approvalModal.model.ServerK8S.NetMode"
              :placeholder="$t('pub.dlg.defaultValue')"
              required
              :required-tip="$t('deployService.form.templateTips')"
            >
              <let-option v-for="d in K8SNetModeOptional" :key="d" :value="d">{{d}}</let-option>
            </let-select>
            <let-input disabled v-else
              size="small"
              :min="0"
              v-model="approvalModal.model.ServerK8S.NetMode"
              :placeholder="$t('deployService.form.placeholder')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
              :pattern-tip="$t('deployService.form.placeholder')"
            ></let-input>
          </let-form-item>
          <let-form-item :label="$t('deployService.table.th.nodeSelector')" itemWidth="30%">
            <let-select v-if="K8SNodeSelectorOptional && K8SNodeSelectorOptional.length > 0"
              size="small" disabled
              v-model="approvalModal.model.ServerK8S.NodeSelector"
              :placeholder="$t('pub.dlg.defaultValue')"
              required
              :required-tip="$t('deployService.form.templateTips')"
            >
              <let-option v-for="d in K8SNodeSelectorOptional" :key="d" :value="d">{{d}}</let-option>
            </let-select>
            <let-input v-else
              size="small" disabled
              :min="0"
              v-model="approvalModal.model.ServerK8S.NodeSelector"
              :placeholder="$t('deployService.form.placeholder')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
              :pattern-tip="$t('deployService.form.placeholder')"
            ></let-input>
          </let-form-item>
        </let-form-group>
        <div>
          <let-form-item :label="$t('serviceApproval.ApprovalResult')" required>
            <let-select
              size="small"
              v-model="approvalModal.model.ApprovalResult"
              required
              :required-tip="$t('deployService.table.tips.empty')"
            >
              <let-option v-for="d in approvalResult" :key="d.label" :value="d.value">{{ d.label }}</let-option>
            </let-select>
          </let-form-item>
        </div>
        <div>
          <let-form-item :label="$t('serviceApproval.ApprovalMark')" required>
            <let-input
              type="textarea"
              :rows="3"
              v-model="approvalModal.model.ApprovalMark"
              :placeholder="$t('serviceApproval.ApprovalMark')"
              required
              :required-tip="$t('deployService.table.tips.empty')"
            >
            </let-input>
          </let-form-item>
        </div>
      </let-form>
    </let-modal>
  </div>
</template>

<script>
const approvalResult = [
  { label: '通过', value: 'true' },
  { label: '驳回', value: 'false' },
];

export default {
  name: 'OperationApproval',

  data() {
    return {
      approvalResult,
      templates: [],
      query: {
        ServerApp: '',
        ServerName: '',
      },
      items: [],
      // 分页
      pagination: {
        page: 1,
        size: 10,
        total:1,
      },
      K8SNetModeOptional: [],
      K8SNodeSelectorOptional: [],
      viewModal: {
        show: false,
        model: null,
      },
      detailModal: {
        show: false,
        model: null,
        isNew: false
      },
      approvalModal: {
        show: false,
        model: null,
        isNew: false
      },
    };
  },

  mounted() {
    this.getServerTemplate()
    this.fetchData();
    this.getDefaultValue();
  },

  methods: {
    // 切换服务实时状态页码
    gotoPage(num) {
      this.pagination.page = num
      this.fetchData()
    },
    getDefaultValue(){
      let { K8SNetModeOptional, K8SNodeSelectorOptional } = this
      this.$ajax.getJSON('/server/api/default', {}).then((data) => {
        if(data.K8SNetModeOptional){
          K8SNetModeOptional = data.K8SNetModeOptional
        }
        if(data.K8SNodeSelectorOptional){
          K8SNodeSelectorOptional = data.K8SNodeSelectorOptional
        }
        this.K8SNetModeOptional = K8SNetModeOptional
        this.K8SNodeSelectorOptional = K8SNodeSelectorOptional
      }).catch((err) => {
        this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
      });
    },
    fetchData() {
      const loading = this.$refs.table.$loading.show();

      return this.$ajax.getJSON('/server/api/deploy_select', {
        ServerApp: this.query.ServerApp,
        ServerName: this.query.ServerName,
        page: this.pagination.page,
        size: this.pagination.size,
      }).then((data) => {
        loading.hide();
        this.items = data.Data;
      }).catch((err) => {
        loading.hide();
        this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
      });
    },

    search() {
      this.fetchData();
    },

    closeDetailModal() {
      this.$refs.detailForm.resetValid();
      this.detailModal.show = false;
      this.detailModal.model = null;
    },
    closeApprovalModal() {
      this.$refs.approvalForm.resetValid();
      this.approvalModal.show = false;
      this.approvalModal.model = null;
    },

    addItem() {
      this.detailModal.model = {};
      this.detailModal.show = true;
      this.detailModal.isNew = true;
    },

    viewItem(d) {
      this.viewModal.model = d;
      this.viewModal.show = true;
    },

    editItem(d) {
      this.detailModal.model = d;
      this.detailModal.show = true;
      this.detailModal.isNew = false;
    },

    approvalItem(d) {
      this.approvalModal.model = d;
      this.approvalModal.show = true;
      this.approvalModal.isNew = false;
    },

    saveItem() {
      if (this.$refs.detailForm.validate()) {
        const model = this.detailModal.model;
        // if(`${model.ApprovalResult}` === 'true'){
        //   model.ApprovalResult = true
        // }else{
        //   model.ApprovalResult = false
        // }
        const loading = this.$Loading.show();
        this.$ajax.postJSON('/server/api/deploy_update', model).then(() => {
          loading.hide();
          this.fetchData().then(() => {
            this.detailModal.show = false;
            this.$tip.success(this.$t('common.success'));
          });
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });
      }
    },

    saveApprovalItem() {
      if (this.$refs.approvalForm.validate()) {
        const model = this.approvalModal.model;
        // if(`${model.ApprovalResult}` === 'true'){
        //   model.ApprovalResult = true
        // }else{
        //   model.ApprovalResult = false
        // }
        const loading = this.$Loading.show();
        this.$ajax.postJSON('/server/api/approval_create', model).then(() => {
          loading.hide();
          this.fetchData().then(() => {
            this.approvalModal.show = false;
            this.$tip.success(this.$t('common.success'));
          });
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });
      }
    },

    removeItem(d) {
      this.$confirm(this.$t('template.delete.confirmTips'), this.$t('common.alert')).then(() => {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/deploy_delete', { DeployId: d.DeployId }).then(() => {
          loading.hide();
          this.fetchData().then(() => {
            this.$tip.success(this.$t('common.success'));
          });
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });
      }).catch(() => {});
    },

    getServerTemplate() {
      this.$ajax.getJSON('/server/api/template_select', {
        isAll: true,
      }).then((data) => {
        this.templates = data.Data;
      }).catch((err) => {
        this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
      });
    },
  },
};
</script>

<style>
.page_operation_approval {
  pre {
    color: #909FA3;
    margin-top: 32px;
  }

  .let_modal__body {
    overflow-y: visible;
  }

  .let_modal__body .let-form .let-box .let-form-item:last-child {
    margin-bottom:20px;
  }
}
</style>
