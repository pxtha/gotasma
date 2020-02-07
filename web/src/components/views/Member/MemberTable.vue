<template>
    <div class="col-md-12" v-if="resources">
          <!-- /.box-header -->
          <div class="box box-body no-padding table-responsive">
            <table class="table table-hover">
            <thead>
                  <th>Member</th>
                  <th>Badge ID</th>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Remove</th>
            </thead>
            <tbody v-if="resources.length > 0">
                <tr v-for="resource in resources" :key="resource.id">
                    <td><avatar :username="resource.name" :size="40"></avatar></td>
                    <td>{{resource.badgeID}}</td>
                    <td>{{resource.name}}</td>
                    <td>{{resource.email}}</td>
                    <td> <a class="btn-remove" @click="showDialogModal(resource)">remove</a></td>
                </tr>
            </tbody>
            <tbody v-else>
                <p class="notice"> No member available in this team! Please choose member from member list</p>
            </tbody>
            </table>
          </div>
          <!-- /.box-body -->
      <v-dialog></v-dialog>
    </div>
</template>
<script>
import Avatar from 'vue-avatar'

export default {
    name: 'team-table',
    props: ['resources', 'currentProject'],
    components: { Avatar },
    methods: {
    showDialogModal(resource) {
      this.$modal.show('dialog', {
        title: 'Are you sure?',
        text: 'Do you want to remove this user?',
        buttons: [
          {
            title: 'OK',
            default: true,
            handler: () => {
              this.currentProject.users.splice(this.currentProject.users.findIndex(deleteUser => deleteUser === resource.id), 1)
              let info = {
                  id: this.currentProject.id,
                  info: this.currentProject.users
              }
              this.$store.dispatch('deleteUserToProject', info)
              resource.project.splice(resource.project.findIndex(deleteProject => deleteProject === this.currentProject.id), 1)
              info.id = resource.id
              info.info = resource.projects
              this.$store.dispatch('deleteProjectofResource', info)
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
    }
}
</script>

<style scoped>
.notice {
  padding: 10px;
  font-size: 14px;
  font-style: italic;
  color: rgb(238, 61, 61)
}
.btn-remove:hover{
  cursor: pointer
}
table>thead>th{
  padding: 10px
}
table {
  border-radius: 10px;
  font-size: 16px !important;
}
</style>