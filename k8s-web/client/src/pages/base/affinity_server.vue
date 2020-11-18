<template>
  <div class="page_base_affinity">
    <div style="float: right">
      <let-button size="small" theme="primary" @click="addItem">{{$t('filter.btn.addServer')}}</let-button>
    </div>
    <let-form inline itemWidth="200px" @submit.native.prevent="search">
      <let-form-item :label="$t('filter.title.node')">
        <let-select
          size="small"
          v-model="query.NodeName"
          :placeholder="$t('pub.dlg.defaultValue')"
        >
          <let-option v-for="d in NodeList" :key="d.NodeName" :value="d.NodeName">{{d.NodeName}}</let-option>
        </let-select>
      </let-form-item>
      <let-form-item :label="$t('filter.title.app')">
        <let-select
          size="small"
          v-model="query.ServerApp"
          :placeholder="$t('pub.dlg.defaultValue')"
          required
          :required-tip="$t('deployService.form.templateTips')"
        >
          <let-option v-for="d in AppList" :key="d.AppName" :value="d.AppName">{{d.AppName}}</let-option>
        </let-select>
      </let-form-item>
      <let-form-item>
        <let-button size="small" type="submit" theme="primary">{{$t('operate.search')}}</let-button>
      </let-form-item>
    </let-form>

    <let-table ref="table" :data="items" :empty-msg="$t('common.nodata')">
      <let-table-column :title="$t('filter.title.app')">
        <template slot-scope="scope">
          <div class="table_con">
            <let-tag theme="success" checked>{{ scope.row.ServerApp }}</let-tag>
          </div>
        </template>
      </let-table-column>
      <let-table-column :title="$t('filter.title.node')">
        <template slot-scope="scope">
          <div class="table_con">
            <let-tag v-for="item in scope.row.NodeName" :key="item" :value="item" theme="success" checked>{{ item }}</let-tag>
          </div>
        </template>
      </let-table-column>
      <let-table-column :title="$t('operate.operates')" width="150px">
        <template slot-scope="scope">
          <let-table-operation @click="editItem(scope.row)">{{$t('operate.update')}}</let-table-operation>
        </template>
      </let-table-column>
    </let-table>

    <let-modal
      v-model="detailModal.show"
      :title="this.$t('filter.title.node')"
      width="800px"
      @on-confirm="addServerItem"
      @on-cancel="closeDetailModal"
    >
      <let-form ref="detailForm" v-if="detailModal.model" itemWidth="700px">
        <let-form-item :label="$t('filter.title.app')" required>
          <let-select
            size="small"
            v-model="detailModal.model.ServerApp"
            :placeholder="$t('pub.dlg.defaultValue')"
            required
            :required-tip="$t('deployService.form.templateTips')"
          >
            <let-option v-for="d in AppList" :key="d.AppName" :value="d.AppName">{{d.AppName}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('filter.title.node')">
          <let-checkbox v-model="isCheckedAll" :value="isCheckedAll">全选</let-checkbox>
          <div class="checkbox_list">
            <let-checkbox class="checkbox_item" v-for="d in NodeList" :key="d.NodeName" :value="checkBoxAddList.indexOf(d.NodeName) > -1" @change="checkedChange(d.NodeName)">{{ d.NodeName }}</let-checkbox>
          </div>
          <!-- <let-select
            size="small"
            v-model="ServerAppStr"
            :placeholder="$t('pub.dlg.defaultValue')"
            @change="appChange"
          >
            <let-option v-for="d in NodeList" :key="d.NodeName" :value="d.NodeName">{{d.NodeName}}</let-option>
          </let-select>
          <div style="border:1px solid #c0c4cc;border-radius:4px;box-sizing:border-box;margin-top:16px;padding:10px 10px 0;" v-if="ServerAppArr && ServerAppArr.length > 0">
            <let-tag style="margin:0 10px 10px 0;"
              v-for="item in ServerAppArr"
              :key="item"
              checked
              closable
              theme="success"
              @close="handleClose(item)"
            >
              {{ item }}
            </let-tag>
          </div> -->
        </let-form-item>
      </let-form>
    </let-modal>
  </div>
</template>

<script>
export default {
  name: 'BaseAffinityServer',

  data() {
    return {
      // 全选
      isCheckedAll: false,
      checkBoxAddList: [],
      checkBoxDelList: [],
      query: {
        NodeName: '',
        ServerApp: '',
      },
      totalPage: 0,
      pageSize: 20,
      page: 1,
      items: [],
      AppList: [],
      NodeList: [],
      ServerAppArr: [],
      ServerAppStr: '',
      ServerAppDeleteArr: [],
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
    this.getAppList()
    this.getNodeList()
  },

  watch: {
    isCheckedAll() {
      let isCheckedAll = this.isCheckedAll;
      if(isCheckedAll) {
        this.NodeList.forEach(item => {
          if(this.checkBoxAddList.indexOf(item.NodeName) === -1){
            this.checkBoxAddList.push(item.NodeName)
          }
          this.checkBoxDelList = []
        })
      }else{
        this.NodeList.forEach(item => {
          if(this.checkBoxDelList.indexOf(item.NodeName) === -1){
            this.checkBoxDelList.push(item.NodeName)
          }
          this.checkBoxAddList = []
        })
      }
    },
  },

  methods: {
    checkedChange(val) {
      if(this.checkBoxAddList.indexOf(val) === -1){
        this.checkBoxAddList.push(val)
        this.checkBoxDelList.splice(this.checkBoxDelList.indexOf(val), 1)
      }else{
        this.checkBoxDelList.push(val)
        this.checkBoxAddList.splice(this.checkBoxAddList.indexOf(val), 1)
      }

      if(this.checkBoxAddList.length === 0){
        this.isCheckedAll = false
      }
    },
    changePage(page) {
      this.page = page;
    },
    fetchData() {
      const loading = this.$refs.table.$loading.show();
      return this.$ajax.getJSON('/server/api/affinity_list_server', this.query).then((data) => {
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

    addItem() {
      this.checkBoxAddList = []
      this.detailModal.model = {};
      this.detailModal.show = true;
      this.detailModal.isNew = true;
    },

    viewItem(d) {
      this.viewModal.model = d;
      this.viewModal.show = true;
    },

    editItem(d) {
      this.checkBoxAddList = d.NodeName || []
      this.detailModal.model = d;
      this.detailModal.show = true;
      this.detailModal.isNew = false;
    },

    addServerItem() {
      if (this.$refs.detailForm.validate()) {
        const model = this.detailModal.model;
        const url = '/server/api/affinity_add_node';

        const loading = this.$Loading.show();

        if(this.checkBoxDelList && this.checkBoxDelList.length > 0){
          this.$ajax.postJSON('/server/api/affinity_del_node', {
            ServerApp: model.ServerApp,
            NodeName: this.checkBoxDelList,
          })
        }

        this.$ajax.postJSON(url, {
          ServerApp: model.ServerApp,
          NodeName: this.checkBoxAddList,
        }).then(() => {
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
      this.$confirm(this.$t('template.delete.confirmTips'), this.$t('common.alert')).then(() => {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/affinity_del_node', { AppName: d.AppName }).then(() => {
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

    getAppList() {
      this.$ajax.getJSON('/server/api/application_select', {}).then((data) => {
        this.AppList = data.Data
      })
    },

    getNodeList() {
      this.$ajax.getJSON('/server/api/node_select', {}).then((data) => {
        this.NodeList = data.Data
      })
    },
  },
};
</script>

<style>
.page_base_affinity {
  pre {
    color: #909FA3;
    margin-top: 32px;
  }

  .let_modal__body {
    overflow-y: visible;
  }
}

.table_con{margin-right:20px;max-height:200px;overflow:auto;padding:10px 0 0;}
.table_con .let-tag{margin:0 10px 10px 0;}
</style>
