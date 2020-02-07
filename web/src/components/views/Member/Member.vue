<template>
  <section class="content">
    <div class="row center-block" v-if="project">
      <div class="header col-md-12">
         <ol class="breadcrumb">
          <li>
            <router-link :to="'/project/' + id">
            <a>
              <i class="fa fa-home"></i> &nbsp;{{project.name}}
            </a>
            </router-link>
          </li>
          <li class="active"> <i class="fa fa-pagelines"></i> &nbsp;{{$route.name.toUpperCase()}}</li>
        </ol>
        <br>
        <button type="button" @click="showTableResources = true" class="btn btn-info pull-right special">Choose from member</button>
        <button type="button" @click="showTableResources = false" class="btn btn-info pull-right special">Hide table</button>
        <i> Here you manage all your project members. You can choose from the already invited team members.</i>
      </div> 
      <member-table :resources="getResourceOfProject" :currentProject="project"></member-table>  
      <transition enter-active-class="animated slideInRight" leave-active-class="animated slideOutRight">     
          <resources-table :showTableResources="showTableResources" :availableResources="availableResources" :projects="projects" :currentProject="project"></resources-table>
      </transition> 
    </div>
  </section>
</template>
<script>
import ResourcesTable from './ResourcesTable'
import MemberTable from './MemberTable'
import { mapState, mapGetters } from 'vuex'
// Require needed datatables modules

export default {
  name: 'member',
  props: ['id'],
  components: { ResourcesTable, MemberTable },
  data() {
    return {
      showTableResources: false
    }
  },
  created() {
    this.$store.dispatch('getProjectById', this.id)
    this.$store.dispatch('getResources')
    this.$store.dispatch('getProjects')
  },
  computed: {
    ...mapState([
      'project',
      'resources',
      'projects'
    ]),
    ...mapGetters([
      'availableResources',
      'getResourceOfProject'
    ])
  },
  mounted() {
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case 'ADD_USER_TO_PROJECT':
          this.$store.dispatch('getResources')
          break
        case 'DELETE_USER_TO_PROJECT':
          this.$store.dispatch('getResources')
          break
      }
    })
  }
}
</script>

<style scoped>
.delete{
  font-size: 14px !important;
}
.center-block{
  height: 100%;
}
.content {
  padding: 60px;
}
.header span {
  font-size: 24px;
  color: #242e35;
  font-weight: bold;
}
.header p {
  font-size: 14px;
  font-style: italic;
  color: #242e35;
  margin-bottom: 20px;
}
.email{
  height: 46px;
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px
}
.special {
  margin: 10px;
  font-size: 16px !important 
}

</style>
