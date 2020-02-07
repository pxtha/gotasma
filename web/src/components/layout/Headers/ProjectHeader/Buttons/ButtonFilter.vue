<template>
  <li>
    <a @click="show" title="Filter">
      <span class="fa fa-sort-amount-asc myIcon"></span>
    </a>
  </li>
</template>

<script>
export default {
  name: 'buttonFilter',
  methods: {
  show () {
    this.$modal.show('someFilter')
  },
  hide () {
    this.$modal.hide('someFilter')
  }
}
}
</script>

<style scoped>
.myIcon {
  color: rgb(43, 42, 42);
  margin: 0 5px 0 5px
}
a {
  cursor: pointer
}
</style>

