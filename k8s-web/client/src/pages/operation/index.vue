<template>
  <div class="page_operation">
    <el-menu class="tree" :default-active="$route.path" @select="selectTree" v-if="treeData && treeData.length > 0">
        <el-menu-item v-for="(item, index) in treeData" :key="index" :index="item.path" @click="clickEvent(item.path)">{{ item.title }}</el-menu-item>
    </el-menu>
    <div class="right-view">
      <router-view ref="childView" class="page_operation_children" :key="$route.params.treeid"></router-view>
    </div>
  </div>
</template>


<!--
<template>
  <div class="page_operation">
    <let-tabs @click="onTabClick" :activekey="$route.path">
      <let-tab-pane :tab="$t('deployService.title.deploy')" tabkey="/operation/deploy"></let-tab-pane>
      <let-tab-pane :tab="$t('deployService.title.approval')" tabkey="/operation/approval"></let-tab-pane>
      <let-tab-pane :tab="$t('deployService.title.history')" tabkey="/operation/history"></let-tab-pane>
      !-- <let-tab-pane :tab="$t('deployService.title.expand')" tabkey="/operation/expand"></let-tab-pane> --
      <let-tab-pane :tab="$t('deployService.title.template')" tabkey="/operation/templates"></let-tab-pane>
    </let-tabs>

    <router-view class="page_operation_children"></router-view>
  </div>
</template>
-->

<script>
let oldPath = '/operation/deploy';

import { operationTree } from '@/store/index.js'

export default {
  name: 'Oparetion',

  data() {
    return {
      treeData: [],
    }
  },

  beforeRouteEnter(to, from, next) {
    if (to.path === '/operation') {
      next(oldPath);
    } else {
      next();
    }
  },

  beforeRouteLeave(to, from, next) {
    oldPath = from.path;
    next();
  },
  mounted() {
    this.loadData()
  },
  methods: {
    loadData() {
      this.treeData = operationTree || []
    },
    selectTree(nodeKey) {
      if (this.$route.path === '/operation') {
        this.$router.push(oldPath);
      } else {
        this.$router.push({
          params: {
            treeid: nodeKey,
          }
        });
      }
    },
    clickEvent(tabkey) {
      this.$router.replace(tabkey);
    },
  },
};
</script>

<style>
.page_operation {
  display: flex;

  &_children {
    padding: 20px 0;    
  }

  .tree{
    padding-top:20px;
    width:160px;

    .el-menu-item {
      height:32px;
      line-height:32px;
    }

    .el-menu-item.is-active {
      color:#3f5ae0
    }
  }

  /*右侧窗口*/
  .right-view {
    display:flex;
    flex-flow: column;
    flex: 1;
    margin-left: 40px;
  }

  .page_operation_children{box-sizing:border-box;display:block;flex:1;overflow:auto;}
}
</style>
