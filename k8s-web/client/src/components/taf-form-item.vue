<template>
  <let-form-item class="taf-form-item">
    <label
      class="let-form-item__label let-label__position_top clickable"
      @click="$emit('onLabelClick')"
    >{{label}}</label>
    <slot></slot>
  </let-form-item>
</template>

<script>
  export default {
    name: 'TafFormItem',
    props: {
      label: String,
    },
  };
</script>

<style>
  .taf-form-item {

  .clickable {
    color: #3f5ae0;
  }

  }
</style>
