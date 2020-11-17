<template>
  <div class="page_base_application">
    <let-button size="small" theme="primary" style="float: right" @click="addItem">{{$t('template.btn.addApplication')}}</let-button>
    <let-form inline itemWidth="200px" @submit.native.prevent="search">
      <let-form-item :label="$t('deployService.form.app')">
        <let-input size="small" v-model="query.AppName"></let-input>
      </let-form-item>
      <let-form-item :label="$t('deployService.form.business')">
        <let-input size="small" v-model="query.BusinessName"></let-input>
      </let-form-item>
      <let-form-item>
        <let-button size="small" type="submit" theme="primary">{{$t('operate.search')}}</let-button>
      </let-form-item>
    </let-form>

    <let-table ref="table" :data="items" :empty-msg="$t('common.nodata')">
      <let-table-column :title="$t('deployService.form.app')" prop="AppName" width="25%"></let-table-column>
      <let-table-column :title="$t('deployService.form.business')" prop="BusinessName" width="25%"></let-table-column>
      <let-table-column :title="$t('deployService.form.appMark')" prop="appMark"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.createTime')" prop="CreateTime"></let-table-column>
      <let-table-column :title="$t('serverList.table.th.createPerson')" prop="CreatePerson"></let-table-column>
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
        <let-form-item :label="$t('deployService.form.app')" required>
          <let-input
            :disabled="detailModal.isNew ? false : true"
            size="small"
            v-model="detailModal.model.AppName"
            :placeholder="$t('deployService.form.placeholder')"
            required
            :required-tip="$t('deployService.form.appTips')"
            pattern="^[a-zA-Z]+$"
            :pattern-tip="$t('deployService.form.placeholder')"
          ></let-input>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.business')">
          <let-select
            size="small"
            v-model="detailModal.model.BusinessName"
            :placeholder="$t('pub.dlg.defaultValue')"
          >
            <let-option v-for="d in BusinessList" :key="d.BusinessName" :value="d.BusinessName">{{d.BusinessShow}}</let-option>
          </let-select>
        </let-form-item>
        <let-form-item :label="$t('deployService.form.appMark')">
          <let-input
            size="small"
            v-model="detailModal.model.AppMark"
          ></let-input>
        </let-form-item>
      </let-form>
    </let-modal>
  </div>
</template>

<script>
export default {
  name: 'BaseApplication',

  data() {
    return {
      query: {
        AppName: '',
        BusinessName: '',
      },
      items: [],
      BusinessList: [],
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
    this.getBusinessList();
  },

  methods: {
    gotoPage(num) {
      this.pagination.page = num
      this.fetchData()
    },
    fetchData() {
      const loading = this.$refs.table.$loading.show();
      return this.$ajax.getJSON('/server/api/application_select', Object.assign(this.query, {
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
      this.detailModal.show = true;
      this.detailModal.isNew = false;
    },

    saveItem() {
      if (this.$refs.detailForm.validate()) {
        const model = this.detailModal.model;
        const url = this.detailModal.isNew ? '/server/api/application_create' : '/server/api/application_update';

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
      this.$confirm(this.$t('template.delete.confirmTips'), this.$t('common.alert')).then(() => {
        const loading = this.$Loading.show();
        this.$ajax.getJSON('/server/api/application_delete', { AppName: d.AppName }).then(() => {
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

    getBusinessList() {
      this.$ajax.getJSON('/server/api/business_select', {}).then((data) => {
        this.BusinessList = [
          {
            BusinessMark: '',
            BusinessName: '',
            BusinessOrder: 0,
            BusinessShow: '无',
            CreatePerson: '',
            CreateTime: "2019-11-11 11:11:11"
          }
        ].concat(data.Data)
      })
    },
  },
};
</script>

<style>
.page_base_application {
  pre {
    color: #909FA3;
    margin-top: 32px;
  }

  .let_modal__body {
    overflow-y: visible;
  }
}
</style>
