<template>
<div >
 
  <aside class="control-sidebar control-sidebar-dark control-sidebar-open">
    <!-- Tab panes -->
    <div class="tab-content">
      <p class="current-time">{{day}} {{date}}</p>
      <form v-on:submit.prevent class="sidebar-form" id="searchForm">
        <div class="input-group" id="searchContainer">
          <span class="input-group-btn">
            <input type="text"
            name="search"
            id="search"
            class="search form-control"
            data-toggle="hideseek" p
            laceholder="Search Menus"
            data-list=".control-sidebar-menu">
            <button type="submit" name="search" id="search-btn" class="btn btn-flat">
              <i class="fa fa-search"></i>
            </button>
          </span>
        </div>
      </form>
      <div class="tab-home" id="control-sidebar-home-tab">
        <h3 class="control-sidebar-heading">Recent Activity</h3>
        <ul class="control-sidebar-menu" v-if="history.length > 0">
          <li v-for="snapshot in history" :key="snapshot">
            <a class="cursor" @click="sendSnapshotID(snapshot.id)">
              <i class="menu-icon fa fa-book bg-green"></i>
              <div class="menu-info">
                <h4 class="control-sidebar-subheading">{{snapshot.description}}</h4>
                <p>{{snapshot.updateDate | momentDetailDate}}</p>
              </div>
            </a>
          </li>
        </ul>
        <ul class="control-sidebar-menu" v-else>
          <li>
            <a>
              <i class="menu-icon fa fa-bug"></i>
              <div class="menu-info">
                <h4 class="control-sidebar-subheading">Empty history</h4>
              </div>
            </a>
          </li>
        </ul>
      </div>
    </div>
  </aside>
<div class="control-sidebar-bg"></div>
</div>
</template>
<script>
export default {
    name: 'SnapshotList',
    props: ['history'],
    components: {
        day: new Date(),
        date: new Date()
    },
    data() {
    return {
        date: this.$dn.date(new Date(), 'dd-mm-yyyy', '-'),
        day: this.$dn.dayText(new Date(), 'en'),
        time: this.$dn.time()
    }
  },
  methods: {
    sendSnapshotID: function(id) {
      this.$emit('clicked', id)
    }
  }
}
</script>
<style scoped>
.cursor{
  cursor: pointer
}
ul li {
  border-bottom: 1px solid #242e35 ;
}
.control-sidebar-heading {
  text-align: center;
}
.current-time
{
    text-align: center;
    font-size: 16px;
    padding-top: 20px;
    font-weight: bold;
}
</style>