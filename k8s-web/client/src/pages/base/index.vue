<template>
  <div class="page_base">
    <el-menu class="tree" :default-active="$route.path" @select="selectTree" v-if="treeData && treeData.length > 0">
        <el-menu-item v-for="(item, index) in treeData" :key="index" :index="item.path" @click="clickEvent(item.path)">{{ item.title }}</el-menu-item>
    </el-menu>
    <div class="right-view">
      <router-view ref="childView" class="page_base_children" :key="$route.params.treeid"></router-view>
    </div>
  </div>
</template>

<script>
let oldPath = '/base/business';

import { baseTree } from '@/store/index.js'

export default {
  name: 'Base',

  data() {
    return {
      treeData: [],
    }
  },

  beforeRouteEnter(to, from, next) {
    if (to.path === '/base') {
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
      this.treeData = baseTree || []
    },
    selectTree(nodeKey) {
      if (this.$route.path === '/base') {
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
.page_base {
  display: flex;

  &_children {
    padding: 30px 0;    
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
    flex: 1;
    margin-left: 40px;
    margin-top: -10px;
  }
}
</style>
