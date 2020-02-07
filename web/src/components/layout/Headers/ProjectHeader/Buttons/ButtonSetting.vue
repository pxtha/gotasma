<template>
  <li>
    <a @click="show" title="Setting">
      <span class="fa fa-gears myIcon" ></span>
    </a>
  </li>
</template>

<script>
export default {
  name: 'buttonSetting',
  props: ['id'],
  methods: {
  show () {
    this.$modal.show('settingModal')
  },
  hide () {
    this.$modal.hide('settingModal')
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

