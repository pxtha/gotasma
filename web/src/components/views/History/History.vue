<template>
<section class="content"> 
   <div class="row center-block">
      <div class="header col-md-12">
         <ol class="breadcrumb">
          <li>
            <router-link :to="'/project/' + id">
            <a>
              <i class="fa fa-home"></i> &nbsp;Current project
            </a>
            </router-link>
          </li>
          <li class="active"> <i class="fa fa-pagelines"></i> &nbsp;{{$route.name.toUpperCase()}}</li>
        </ol>
        <p> Here you see can all the states of your project at a specific point in time, <br> you navigate around the history list and choose see what you have saved</p>
      </div>
      <snapshot-list :history="history" @clicked="eventChild" ></snapshot-list>
</div>
 <gantt-history :snapshot="snapshot"></gantt-history>
</section>
</template>

<script>
import GanttHistory from './GanttHistory'
import SnapshotList from './SnapshotList'
import { mapState } from 'vuex'
export default {
  name: 'history',
  data() {
    return {
      snapshot: {}
    }
  },
  props: ['id'],
  components: {
    GanttHistory,
    SnapshotList
  },
  methods: {
    eventChild(snapshotId) {
      this.snapshot = this.$store.getters.getHistoryById(snapshotId)
    }
  },
  computed: {
    ...mapState(['history'])
  },
  created() {
    this.$store.dispatch('getHistoryById', this.id)
  }
}
</script>
<style scoped>
span {
  font-size: 24px;
  color: #242e35;
  font-weight: bold;
}
</style>
