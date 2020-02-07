<template>
  <div class="col-md-12 box box-body">
            <div class="col-sm-12 table-responsive">
              <table
                aria-describedby="resourcesTable_info"
                role="grid"
                id="resourcesTable"
                class="table   table-hover dataTable"
              >
                <thead>
                  <tr role="row">
                    <th
                      aria-label="Rendering engine: activate to sort column descending"
                      aria-sort="ascending"
                      style="width: 167px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="resourcesTable"
                      tabindex="0"
                      class="sorting_asc"
                    ></th>
                    <th
                      aria-label="Rendering engine: activate to sort column descending"
                      aria-sort="ascending"
                      style="width: 167px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="resourcesTable"
                      tabindex="0"
                      class="sorting_asc"
                    >Bagde ID</th>
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
                      aria-label="Platform(s): activate to sort column ascending"
                      style="width: 182px;"
                      colspan="1"
                      rowspan="1"
                      aria-controls="resourcesTable"
                      tabindex="0"
                      class="sorting"
                    >Email</th>
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
                  <tr class="even" role="row" v-for="resource of resources" :key="resource.badgeID">
                    <td>  <avatar :username="resource.name" :size="40"></avatar></td>
                    <td class="sorting_1">
                      {{resource.badgeID}}
                    </td>
                    <td>{{resource.name}}</td>
                    <td>{{resource.email}}</td>
                    <td>
                      <div class="external-event bg-red" v-for="project in getProjectsOfResource(resource.id)" :key="project">{{project}}</div>
                    </td>
                    <td ><a class="btn btn-app" @click="$modal.show('newresource', {resource})"><i class="fa fa-edit"></i></a>
                        <a class="btn btn-app" style="color:#c70707c2" @click="showDialogMember(resource.id)"><i class="fa fa-remove"></i></a>           
                    </td> 
                  </tr>                        
                </tbody>
              </table>
            </div>
  <v-dialog/>
  </div> 
  </template>
<script>
import NewResource from './NewResource'
import $ from 'jquery'
import { mapState } from 'vuex'
import Avatar from 'vue-avatar'

require('datatables.net-bs')

export default {
  name: 'resource-table',
  components: {
    NewResource,
    Avatar
  },
  methods: {
    showDialogMember(id) {
      this.$modal.show('dialog', {
        title: 'Are you sure?',
        text: 'Do you wish to delete resource',
        buttons: [
          {
            title: 'OK',
            default: true,
            handler: () => {
              this.$store.dispatch('deleteResource', id)
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
      }
  },
  updated() {
      this.$nextTick(() => {
        $('#resourcesTable').DataTable()
      })
  },
  created() {
    this.$store.dispatch('getResources')
    this.$store.dispatch('getProjects')
  },
  computed: {
   ...mapState([
     'resources',
     'projects'
   ])
  },
  mounted() {
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case 'ADD_RESOURCE':
          this.$store.dispatch('getResources')
          break
        case 'DELETE_RESOURCE':
          this.$store.dispatch('getResources')
          break
        case 'EDIT_RESOURCE':
        this.$store.dispatch('getResources')
        break
      }
    })
  }
}
</script>
<style>
.btn-edit-remove{
  font-size: 30px;
  margin-left: 10px;
}
table{
  color: #242E35;
  border-radius: 10px;
  font-size: 16px !important
}
td a {
  min-width: 30px !important;
  width: 40px;
  height: 40px
}
td a i{
  padding: 0px;
  margin-bottom: 10px;
}
</style>
