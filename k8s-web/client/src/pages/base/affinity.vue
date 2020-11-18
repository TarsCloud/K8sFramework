<template>
  <div class="page_base_affinity">
    <!-- <let-tabs> -->
      <!-- <let-tab-pane tab="节点聚合">
        <AffinityNode></AffinityNode>
      </let-tab-pane> -->
      <!-- <let-tab-pane tab="应用聚合"> -->
        <AffinityServer></AffinityServer>
      <!-- </let-tab-pane> -->
    <!-- </let-tabs> -->
  </div>
</template>

<script>
import AffinityNode from './affinity_node'
import AffinityServer from './affinity_server'
export default {
  name: 'BaseAffinity',

  components: {
    AffinityNode,
    AffinityServer
  },

  data() {
    return {}
  },

  mounted() {
    
  },

  methods: {
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
</style>
