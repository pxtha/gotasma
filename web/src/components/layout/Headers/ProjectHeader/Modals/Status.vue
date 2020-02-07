<template>
   <li class="u-pull-left">
    <a title="Status" class="notice">
      <span v-if="getStatus" > Last update time: {{ getStatus }} </span>
    </a>
  </li>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'status',
  computed: {
    ...mapGetters([ 'getStatus' ])
  }
}
</script>

<style scoped>
  .notice {
    font-size: 12px;
  }
</style>

