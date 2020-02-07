<template>
  <li>
    <a @click="showText" title="Show/hide text">
      <span class="fa fa-list-ul myIcon"></span>
    </a>
  </li>
</template>
<script>
import { EventBus } from '@/main.js'
import { mapState } from 'vuex'
export default {
    name: 'textChart',
  methods: {
    showText() {
      EventBus.$emit('chartText', this.chartText)
    }
  },
  computed: {
      ...mapState([ 'chartText' ])
  },
  beforeDestroy() {
    EventBus.$off('chartText')
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