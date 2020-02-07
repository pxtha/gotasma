<template>
  <li>
    <a @click="showConfirm">
      <span class="fa fa-download" id="myIcon" />
      <span class="hidden-xs myIcon">SAVE</span>
    </a>
  </li>
</template>

<script>
export default {
  name: 'buttonSave',
  methods: {
    showConfirm() {
      this.$modal.show('SaveModal')
    }
  }
}
</script>

<style scoped>
.myIcon {
  color: rgb(43, 42, 42);
  margin: 0 5px 0 5px
}
a{
  cursor: pointer
}
</style>

