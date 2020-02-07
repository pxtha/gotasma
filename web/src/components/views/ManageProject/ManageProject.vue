<template>
  <section class="content">
    <div class="row center-block">
      <div class="header col-md-12">
        <span>Projects</span> <br>
        <button
          type="submit"
          class="btn btn-info pull-right special"
          @click="$modal.show('createNewProj')">
          Create New Project</button>
        <i> Here you can manage all your projects. You can also create new project by clicking a single button.</i>
      </div>
    </div>
    <tableProject :projects="projects" ></tableProject>
    <NewProject :projects="projects"></NewProject>
  </section>  
</template>


<script>
import NewProject from './NewProject'
import TableProject from './TableProject'
import { mapState } from 'vuex'

export default {
  name: 'manageProject',
  components: {
    NewProject,
    TableProject
  },
  created() {
    this.$store.dispatch('getProjects')
  },
  computed: {
    ...mapState([
      'projects'
    ])
  }
}
</script>
<style scoped>
span {
  font-size: 24px;
  color: #242e35;
  font-weight: bold;
}
.btn-spacing {
  margin: 20px 0;
  vertical-align: middle;
}
.content {
  margin: 20px !important;
  min-height: 1102px !important
}
.vdp-datepicker .vdp-datepicker__calendar{
  color: #000 !important;
}

</style>
