<template>
  <ul class="sidebar-menu">
    <li class="header">Projects managment</li>
    <router-link tag="li" class="pageLink" to="/manageproject">
      <a>
        <i class="fa fa-suitcase"></i>
        <span class="page">View all project</span>
      </a>
    </router-link>
    <!-- nho fix nhe -->
    <li class="treeview">
      <a href="#">
        <i class="fa fa-star-o"></i>
        <span class="treeview-title">Highlighted Projects</span>
        <span class="pull-right-container pull-right">
          <i class="fa fa-angle-left fa-fw"></i>
        </span>
      </a>
      <ul class="treeview-menu" v-if="highlightedProjects.length">
        <li v-for="hProj in highlightedProjects" :key="hProj">
          <a :href="'/project/' + hProj.id + '/'" >
            <i class="fa fa-file"></i> {{ hProj.name }}
          </a>
        </li>
      </ul>
    </li>
    <!--  -->
    <li class="header">Resources</li>
    <router-link tag="li" class="pageLink" to="/resources">
      <a>
        <i class="fa fa-users"></i>
        <span class="page">View all resources</span>
      </a>
    </router-link>

    <li class="header">Exceptions</li>
    <router-link tag="li" class="pageLink" to="/exception">
      <a>
        <i class="fa fa-calendar-times-o  "></i>
        <span class="page">View Exceptions</span>
      </a>
    </router-link>
  </ul>
</template>

<script>
import { mapState } from 'vuex'
export default {
  name: 'SidebarMenu',
  computed: {
    ...mapState([ 'highlightedProjects' ])
  },
  created() {
    if (this.highlightedProjects.length === 0) {
      this.$store.dispatch('getHighlightedProjects')
    }
  },
  mounted() {
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case 'HIGHLIGHT_PROJECT':
          this.$store.dispatch('getHighlightedProjects')
          break
      }
    })
  }
}
</script>

<style>
/* override default */
.sidebar-menu > li > a {
  padding: 12px 15px 12px 15px;
}

.sidebar-menu li.active > a > .fa-angle-left,
.sidebar-menu li.active > a > .pull-right-container > .fa-angle-left {
  animation-name: rotate;
  animation-duration: 0.2s;
  animation-fill-mode: forwards;
}

.treeview-title {
  z-index: 1;
}

@keyframes rotate {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(-90deg);
  }
}

span {
  color: #eeeee7;
}
.page{
  font-weight: bold;
}
</style>
