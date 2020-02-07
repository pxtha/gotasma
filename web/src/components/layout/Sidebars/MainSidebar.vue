<template>
  <aside class="main-sidebar" >
    <!-- sidebar: style can be found in sidebar.less -->
    <section class="sidebar">
      <!-- Sidebar Menu -->
      <sidebar-menu />
      <!-- /.sidebar-menu -->
      <div class="user-panel">
        <div class="pull-left image">
          <img :src="admin.avatar" />
        </div>
        <div class="pull-left info">
          <div>
            <p class="white">{{admin.fullName}}</p>
          </div>
          <a href="javascript:;">
            <i class="fa fa-circle text-success"></i> DC-14
          </a>
        </div>
      </div>
    </section>
    <!-- /.sidebar -->
  </aside>
</template>
<script>

import SidebarMenu from './SidebarElements/SidebarMenu'

export default {
  name: 'Sidebar',
  props: ['admin'],
  components: { SidebarMenu },
  mounted: function() {
    window
      .jQuery('[data-toggle="hideseek"]')
      .off()
      .hideseek()
  }
}
</script>
<style scope="local">
.user-panel {
  position: absolute;
  bottom: 0;
  background-color: #d8d3c0;
  height: 40px;
  padding-bottom: 50px;
}
.user-panel .image img{
  border-radius: 10px;
}
#searchForm {
  padding-left: 0em;
  padding-right: 0em;
}
#searchContainer {
  height: 100%;
  padding-bottom: 0em;
}
#search {
  width: 80%;
  float: right
}

#search-btn {
  width: 20%;
}
.main-sidebar{
 margin-top: 5px;
  height: 100%; 
   -moz-box-shadow:    0px 0 6px 0 #555;
  -webkit-box-shadow:  0px 0 6px 0 #555;
  box-shadow:          0px 0 6px 0 #555;
}
</style>
