<template>
  <div class="col-md-12 box box-body special">
            <div class="col-sm-12 table-responsive">
              <table
                aria-describedby="tableProjects_info"
                role="grid"
                id="tableProjects"
                class="table table-bordered table-hover dataTable"
              >
                <thead>
                  <tr role="row">
                    <th
                      aria-label="Rendering engine: activate to sort column descending"
                      aria-sort="ascending"
                      style="width: 130px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="tableProjects"
                      tabindex="0"
                      class="sorting_asc"
                    >Project name</th>
                    <th
                      aria-label="Browser: activate to sort column ascending"
                      style="width: 20px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="tableProjects"
                      tabindex="0"
                      class="sorting"
                    >Members</th>
                    <th
                      aria-label="Browser: activate to sort column ascending"
                      style="width: 10px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="tableProjects"
                      tabindex="0"
                      class="sorting"
                    >Tasks</th>
                    <th
                      aria-label="Browser: activate to sort column ascending"
                      style="width: 150px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="tableProjects"
                      tabindex="0"
                      class="sorting"
                    >Start date</th>
                    <th
                      aria-label="Engine version: activate to sort column ascending"
                      style="width: 150px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="tableProjects"
                      tabindex="0"
                      class="sorting"
                    >Last change (Last Update)</th>
                    <th
                      aria-label="Engine version: activate to sort column ascending"
                      style="width: 80px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="tableProjects"
                      tabindex="0"
                      class="sorting"
                    ></th>
                  </tr>
                </thead>
                <tbody>
                  <tr class="even" role="row" v-for="project of projects" :key="project.id" >
                    <td class="sorting_1">
                      <router-link :to="'../project/' + project.id ">
                          <a>{{ project.name }}</a>
                      </router-link></td>
                    <td v-if="project.users">{{ project.users | count }}</td>
                    <td v-else> 0 </td>
                    <td v-if="project.tasks">{{ project.tasks | count }}</td>
                    <td>{{ project.startDate | momentNormalDate }}</td>
                    <td>{{ project.updateDate | momentDetailDate }}</td>
                    <td>
                      <a class="btn btn-app del-btn" title="Delete project" @click="showDialog(project.id)"><i class="fa fa-remove"></i></a>
                      <a 
                          v-if="project.highlighted === false"
                          class="btn btn-app star-btn" 
                          title="Highlight project" 
                          @click="highlightProject(project)">
                          <i class="fa fa-star-o"></i>
                      </a>
                      <a  
                          v-else
                          class="btn btn-app star-btn" 
                          title="Highlight project"
                          @click="highlightProject(project)">
                          <i class="fa fa-star"></i>
                      </a>
                    </td>
                  </tr>
                </tbody>  
                <tfoot>
                  <tr>
                    <th colspan="1" rowspan="1">Project name</th>
                    <th colspan="1" rowspan="1">Number of members</th>
                    <th colspan="1" rowspan="1">Max effort</th>
                    <th colspan="1" rowspan="1">Start date</th>
                    <th colspan="1" rowspan="1">Last change (Last Update)</th>
                    <th colspan="1" rowspan="1"></th>
                  </tr>
                </tfoot>
              </table>
      </div>
      <v-dialog/>
  </div>
</template>

<script>
import $ from 'jquery'
require('datatables.net-bs')

export default {
  name: 'tableProject',
  props: ['projects'],
  mounted() {
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case 'ADD_PROJECT':
          this.$store.dispatch('getProjects')
          break
        case 'DELETE_PROJECT':
          this.$store.dispatch('getProjects')
         break
      }
    })
  },
  methods: {
    showDialog(id) {
      this.$modal.show('dialog', {
        title: 'Are you sure?',
        text: 'Do you wish to delete?',
        buttons: [
          {
            title: 'OK',
            default: true,
            handler: () => {
              this.$store.dispatch('deleteProject', id)
              this.$modal.hide('dialog')
            }
          },
          {
            title: 'CANCEL',
            handler: () => {
              this.$modal.hide('dialog')
            }
          }
        ]
      })
    },
    highlightProject(project) {
      if (project.highlighted === false) {
        project.highlighted = true
        this.$store.dispatch('highlightProjects', project)
      } else {
        project.highlighted = false
        this.$store.dispatch('highlightProjects', project)
      }
    }
  },
  beforeUpdate() {
      this.$nextTick(() => {
        $('#tableProjects').DataTable()
      })
  }
}
</script>

<style scoped>
.special{
  margin: 20px;
  margin-right: 20px !important
}
table {
  color: #242e35;
  font-size: 16px !important;
}
td a {
  min-width: 30px !important;
  width: 40px;
  height: 40px
}
.del-btn {
  color: #c70707c2
}
.star-btn {
  color: #ffca00
}
</style>