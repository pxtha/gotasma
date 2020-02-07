<template>
  <li>
    <a @click="showModal" title="Add Parent Task">
      <span class="fa fa-plus myIcon"></span>
    </a>
  </li>
</template>
<script>
export default {
  methods: {
    showModal() {
        this.$modal.show('AddParentTaskModal')
    }
  }
}
</script>

<style scoped>
.myIcon {
  color: rgb(43, 42, 42);
  margin: 0 5px 0 5px;
  font-size: 20px
}
a {
  cursor: pointer
}
</style>