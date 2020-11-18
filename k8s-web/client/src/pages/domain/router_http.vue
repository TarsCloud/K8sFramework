<template>
  <div class="page_router_http">
    <let-button size="small" theme="primary" style="float: right" @click="addItem">{{$t('filter.btn.add')}}</let-button>
    <let-form inline itemWidth="200px" @submit.native.prevent="search">
      <let-form-item :label="$t('filter.title.domain')">
        <let-input size="small" v-model="query.domain"></let-input>
      </let-form-item>
      <let-form-item>
        <let-button size="small" type="submit" theme="primary">{{$t('filter.btn.search')}}</let-button>
      </let-form-item>
    </let-form>

    <let-table ref="table" :data="items" :empty-msg="$t('common.nodata')">
      <let-table-column :title="$t('filter.title.domain')" prop="DomainValue"></let-table-column>
      <let-table-column :title="$t('filter.title.router')" prop="MatchValue"></let-table-column>
      <let-table-column :title="$t('filter.title.routerType')" prop="MatchType">
        <template slot-scope="scope">
          <span v-if="scope.row.MatchType === 'Path'">全量匹配</span>
          <span v-else-if="scope.row.MatchType === 'Regex'">全量正则匹配</span>
          <span v-else-if="scope.row.MatchType === 'Prefix'">前缀匹配</span>
        </template>
      </let-table-column>
      <!-- <let-table-column :title="$t('filter.title.serverId')" prop="BackendServerId"></let-table-column> -->
      <let-table-column :title="$t('filter.title.app')" prop="BackendServerApp"></let-table-column>
      <let-table-column :title="$t('filter.title.business')" prop="BackendServerName"></let-table-column>
      <!-- <let-table-column :title="$t('filter.title.adapterId')" prop="BackendAdapterId"></let-table-column> -->
      <let-table-column :title="$t('filter.title.adapterName')" prop="BackendAdapterName"></let-table-column>
      <let-table-column :title="$t('operate.operates')" width="180px">
        <template slot-scope="scope">
          <let-table-operation @click="editItem(scope.row)">{{$t('operate.update')}}</let-table-operation>
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
      <pre v-if="viewModal.model">{{viewModal.model.TemplateContent}}</pre>
      <div slot="foot"></div>
    </let-modal>

    <let-modal
      v-model="detailModal.show"
      :title="detailModal.isNew ? this.$t('dialog.title.add') : this.$t('dialog.title.edit')"
      width="800px"
      @on-confirm="saveItem"
      @on-cancel="closeDetailModal"
    >
      <let-form ref="detailForm" v-if="detailModal.model" itemWidth="700px">
        <let-form-item :label="$t('filter.title.domain')" required>
          <let-select
            size="small"
            v-model="detailModal.model.DomainValue"
            :placeholder="$t('common.notEmpty')"
            required
            :required-tip="$t('common.notEmpty')"
          >
            <let-option v-for="d in DomainList" :key="d.DomainValue" :value="d.DomainValue">{{d.DomainValue}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('filter.title.router')" required>
          <let-input
            size="small"
            v-model="detailModal.model.MatchValue"
            :placeholder="$t('common.notEmpty')"
            required
            :required-tip="$t('common.notEmpty')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('filter.title.routerType')" required>
          <let-select
            size="small"
            v-model="detailModal.model.MatchType"
            :placeholder="$t('common.notEmpty')"
            required
            :required-tip="$t('common.notEmpty')"
          >
            <let-option v-for="d in MatchValueList" :key="d" :value="d">
              <template v-if="d === 'Path'">全量匹配</template>
              <template v-else-if="d === 'Regex'">全量正则匹配</template>
              <template v-else-if="d === 'Prefix'">前缀匹配</template>
            </let-option>
          </let-select> 
        </let-form-item>
        <let-form-item :label="$t('filter.title.serverId')" required>
          <let-select @change="ServerAdapterEvent"
            size="small"
            v-model="detailModal.model.BackendServerId"
            :placeholder="$t('common.notEmpty')"
            required
            :required-tip="$t('common.notEmpty')"
          >
            <let-option v-for="d in ServerList" :key="d.ServerId" :value="d.ServerId">{{d.ServerApp}}.{{d.ServerName}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('filter.title.adapterId')" required>
          <let-select
            size="small"
            v-model="detailModal.model.BackendAdapterId"
            :placeholder="$t('common.notEmpty')"
            required
            :required-tip="$t('common.notEmpty')"
          >
            <let-option v-for="d in ServerAdapter" :key="d.AdapterId" :value="d.AdapterId">{{d.Port}}|{{d.Name}}</let-option>
          </let-select>
        </let-form-item>
      </let-form>
    </let-modal>
  </div>
</template>

<script>
export default {
  name: 'RouterHttp',

  data() {
    return {
      query: {
        domain: '',
      },
      MatchValueList: [],
      items: [],
      DomainList: [],
      ServerList: [],
      ServerAdapter: [],
      // 分页
      pagination: {
        page: 1,
        size: 10,
        total:1,
      },
      viewModal: {
        show: false,
        model: null,
      },
      detailModal: {
        show: false,
        model: null,
        isNew: false
      },
    };
  },

  mounted() {
    this.fetchData();
    this.getDomainList()
    this.getServerList()
    this.getDefaultValue()
  },

  methods: {
    gotoPage(num) {
      this.pagination.page = num
      this.fetchData()
    },
    fetchData() {
      const loading = this.$refs.table.$loading.show();
      return this.$ajax.getJSON('/server/api/router_http_select', Object.assign(this.query, {
        page: this.pagination.page,
        size: this.pagination.size,
      })).then((data) => {
        loading.hide();
        this.items = data.Data
        this.pagination.total = Math.ceil(data.Count.FilterCount / this.pagination.size)
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

      if(d.BackendServerId){
        this.ServerAdapterEvent(d.BackendServerId)
      }

      this.detailModal.show = true;
      this.detailModal.isNew = false;
    },

    saveItem() {
      if (this.$refs.detailForm.validate()) {
        const model = this.detailModal.model;
        const url = this.detailModal.isNew ? '/server/api/router_http_create' : '/server/api/router_http_update';

        const loading = this.$Loading.show();
        this.$ajax.postJSON(url, model).then(() => {
          loading.hide();
          this.$tip.success(this.$t('common.success'));
          this.closeDetailModal();
          this.fetchData();
        }).catch((err) => {
          loading.hide();
          this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
        });
      }
    },

    removeItem(d) {
      this.$confirm(this.$t('dialog.tips.delete'), this.$t('common.alert')).then(() => {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/router_http_delete', { RouterId: d.RouterId }).then(() => {
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

    getDomainList() {
      this.$ajax.getJSON('/server/api/domain_select', {
        isAll: true,
      }).then((data) => {
        this.DomainList = data.Data
      })
    },

    getServerList() {
      this.$ajax.getJSON('/server/api/server_list', {
        isAll: true,
      }).then((data) => {
        this.ServerList = data
      })
    },

    getServerAdapter(ServerId) {
      this.$ajax.getJSON('/server/api/server_adapter_select', {
        ServerId,
        isAll: true,
        IsTaf: false,
      }).then((data) => {
        this.ServerAdapter = data.Data
      })
    },

    ServerAdapterEvent(data) {
      this.ServerList.forEach(item => {
        if(`${item.ServerId}` === `${data}`){
          this.getServerAdapter(`${item.ServerApp}.${item.ServerName}`)
        }
      })
    },

    getDefaultValue(){
      let { MatchValueList } = this
      this.$ajax.getJSON('/server/api/default', {}).then((data) => {
        if(data.HttpRouterMatchTypeOptional){
          MatchValueList = data.HttpRouterMatchTypeOptional
        }
        this.MatchValueList = MatchValueList
      }).catch((err) => {
        this.$tip.error(`${this.$t('common.error')}: ${err.message || err.err_msg}`);
      });
    },
  },
};
</script>

<style>
.page_router_http {
  pre {
    color: #909FA3;
    margin-top: 32px;
  }

  .let_modal__body {
    overflow-y: visible;
  }
}
</style>
