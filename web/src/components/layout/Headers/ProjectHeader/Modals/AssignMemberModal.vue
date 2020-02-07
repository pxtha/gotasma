<template>
    <modal name="AssignMember" transition="nice-modal-fade" 
      :draggable="true" 
      :reset="true"
      height="auto"
      :clickToClose="false"
      @before-open="beforeOpen"
    >
      <div class="box box-group">
        <div class="box-header with-border dark">
          <h5 class="box-title2">Assign member</h5>
        </div>

        <div class="box-body">
          <div class="row" >
            <div class="col-xs-12">
              <h4 class="title">Assignee: </h4>
              <div class="input-group">
                <span class="input-group-addon"><i class="fa fa-fw fa-child"></i></span>
                  <multiselect 
                      name="users"
                      v-model="user"
                      :class="{ 'is-invalid': errors.has('users') }" 
                      :options="users"
                      :multiple="true"
                      :close-on-select="true"
                      :clear-on-select="false"
                      label="name"
                      track-by="name"
                      :max-height=200>
                  </multiselect>
              </div>
          </div>
          </div>
          <!--Input section-->
          <div class="row">
            <div class="col-xs-12">
                <h4 class="title" >Task ID</h4>
                <div class="input-group ">
                    <span class="input-group-addon">
                    <i class="fa fa-sun-o"></i>
                    </span>
                    <input
                    type="number"
                    v-model.number="currentTask.id"   
                    width="100%"
                    class="form-control"
                    disabled
                    />
                </div>
            </div>
            <div class="col-xs-12">
                <h4 class="title" >Task label</h4>
                <div class="input-group">
                    <span class="input-group-addon">
                    <i class="fa fa-tag"></i>
                    </span>
                    <input
                    type="text"
                    v-model="currentTask.label"
                    class="form-control"   
                    width="100%" 
                    disabled      
                    />
                </div>
            </div>
          </div>
        </div>
        <div class="box-footer">
          <button class="btn-create button-modal pull-right" @click="handleSubmit()">Assign Member</button>
          <button class="btn-close button-modal" @click="closeModal">Cancle</button>
        </div>
      </div>
    </modal>
</template>
<script>
import Multiselect from 'vue-multiselect'
import { EventBus } from '@/main.js'

export default {
    name: 'AssignMember',
    components: { Multiselect },
    props: ['users'],
    data() {
        return {
          currentTask: {},
          user: [],
          beforeEdit: ''
        }
    },
    methods: {
        beforeOpen(event) {
          this.currentTask = event.params.data
          this.beforeEdit = Object.assign({}, this.currentTask)
        },
        closeModal() {
            Object.assign(this.currentTask, this.beforeEdit)
            this.$modal.hide('AssignMember')
        },
        handleSubmit() {
          this.$validator.validate().then(valid => {
                  if (valid) {
                      let info = {
                        user: this.user,
                        taskId: this.currentTask.id
                      }
                      EventBus.$emit('assignMember', info)
                      this.$modal.hide('AssignMember')
                  } else {
                      this.$modal.show('dialog', {
                        title: 'Error',
                        text: 'Invalid input!'})
                  }
              })
           }
    },
    beforeDestroy() {
      EventBus.$off('assignMember')
    }
}
</script>
<style scoped>
.box-footer {
  margin-top: 40px
}
.title {
  padding: 6px 6px;
  font-size: 14px;
  font-weight: 600;
  margin-left: 0px;
  padding-left: 0px 
}
</style>