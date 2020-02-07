<template>
  <li>
    <router-link :to="'/project/' + id + '/history'">
        <a title="History (snapshots list)">
          <span class="fa fa-history myIcon"></span>
        </a>
    </router-link>
  </li>
</template>

<script>
export default {
  name: 'buttonHistory',
  props: ['id']
}
</script>

<style scoped>
.myIcon {
  color:#3fb0ac;
  margin: 0 5px 0 5px
}
</style>

