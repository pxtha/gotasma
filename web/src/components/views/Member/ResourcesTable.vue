<template>
  <div class="col-sm-12 table-responsive" v-if="this.showTableResources && availableResources && currentProject">
    <table
      aria-describedby="resourcesTable_info"
      role="grid"
      id="resourcesTable"
      class="table table-striped dataTable box box-body"
    >
      <thead>
        <tr role="row">
          <th
            aria-label="Browser: activate to sort column ascending"
            style="width: 207px;"
            colspan="1"
            rowspan="1"
            aria-controls="resourcesTable"
            tabindex="0"
            class="sorting"
          >Full Name</th>
          <th
            aria-label="Browser: activate to sort column ascending"
            style="width: 207px;"
            colspan="1"
            rowspan="1"
            aria-controls="resourcesTable"
            tabindex="0"
            class="sorting"
          >Email Name</th>
          <th
            aria-label="Engine version: activate to sort column ascending"
            style="width: 142px;"
            colspan="1"
            rowspan="1"
            aria-controls="resourcesTable"
            tabindex="0"
            class="sorting"
          >Current Projects</th>
          <th 
            aria-label="CSS grade: activate to sort column ascending"
            style="width: 101px;"
            colspan="1"
            rowspan="1"
            aria-controls="resourcesTable"
            tabindex="0"
            class="sorting"
          ></th>
        </tr>
      </thead>
      <tbody>
        <tr class="even" role="row" v-for="resource of availableResources" :key="resource.id">
          <td>{{resource.name}}</td>
          <td>{{resource.email}}</td>
          <td>
          <div class="external-event bg-yellow" v-for="project in getProjectsOfResource(resource.id)" :key="project">{{project}}</div>
          </td>
          <td >
              <a class="btn" @click="addResource(resource, currentProject)" ><i class="fa fa-user-plus special"></i></a>           
          </td> 
        </tr>                        
      </tbody>
    </table>
    <i class="notice" v-if="availableResources.length === 0">No availabe resources</i>
      <v-dialog/>
  </div>
</template>
<script>
import $ from 'jquery'

require('datatables.net-bs')

export default {
  name: 'resource-table',
  props: ['showTableResources', 'availableResources', 'projects', 'currentProject'],
  methods: {
    getProjectsOfResource(idResource) {
      let projectsName = []
      for (let i = 0; i < this.projects.length; i++) {
        let current = this.projects[i]
          for (let j = 0; j < current.users.length; j++) {
            if (idResource === current.users[j]) {
              projectsName.push(current.name)
              break
            }
          }
      }
      return projectsName
    },
    addResource(resource, project) {
      this.$modal.show('dialog', {
        title: 'Are you sure?',
        text: 'Do you wish to add this member project?',
        buttons: [
          {
            title: 'OK',
            default: true,
            handler: () => {
                project.users.push(resource.id)
                let addInfo = {
                    id: project.id,
                    newInfo: project.users
                }
                if (typeof resource.projects === 'undefined') {
                  resource.projects = []
                }
                resource.projects.push(project.id)
                this.$store.dispatch('addResourceToProject', addInfo)
                    addInfo.id = resource.id
                    addInfo.newInfo = resource.projects
                this.$store.dispatch('addProjectToResource', addInfo)
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
    }
  },
  updated() {
      this.$nextTick(() => {
        $('#resourcesTable').DataTable()
      })
  }
}
</script>
<style scoped>
.special{
    font-size: 20px !important;
    color:green;
}
table{
  font-size: 16px !important
}
.notice {
  color: red;
}
</style>
