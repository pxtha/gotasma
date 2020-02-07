<template>
    <modal name="AddTask" transition="nice-modal-fade" 
        :draggable="true" 
        :reset="true"
        height="auto"
        :resizable="true"
        @before-open="beforeOpen"
    >
      <div class="box box-group">
        <div class="box-header with-border dark">
          <h5 class="box-title2">Add new task</h5>
        </div>
        <div class="box-body">
          <div class="row">
             <div class="col-xs-12"> 
                <h4 class="title" >Parent task ID</h4>
                <div class="input-group">
                  <span class="input-group-addon">
                    <i class="fa fa-shield"></i>
                  </span>
                  <input
                    class="form-control"
                    v-model="currentTask.id"
                    disabled
                  />
                </div>
            </div>
            <div class="col-xs-6" hidden>
              <h4 class="title" >Task ID</h4>
              <div class="input-group ">
                <span class="input-group-addon">
                  <i class="fa fa-sun-o"></i>
                </span>
                <input
                  type="number"
                  v-model.number="newTaskInfo.id"
                  class="form-control"
                  disabled
                />
              </div>
            </div>
          </div>
          <h4 class="title col-xs-12" >Task label</h4>
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-tag"></i>
            </span>
            <input
              v-validate="'required|min:3'"
              v-model="newTaskInfo.label"
              class="form-control"
              placeholder="Enter task label"
              name="Task label"
              :class="{ 'is-invalid':submitted &&  errors.has('Task label') }"
            />
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
              <div v-if="submitted && errors.has('Task label')" class="invalid-feedback">{{ errors.first('Task label') }}</div> 
            </transition>
          </div>
            <h4 class="title col-xs-12">Type</h4>
            <div class="input-group col-xs-12">
              <span class="input-group-addon">
                <i class="fa fa-fw fa-hand-pointer-o"></i>
              </span>
              <select
              class="form-control"
              v-model="newTaskInfo.type">
                <option v-for="type in taskType" :value="type.value" :key="type">
                    {{ type.text }}
                </option>
              </select>
            </div>

          <div class="row">
            <div class="col-xs-6">
              <h4 class="title" >Start date</h4>
              <div class="input-group">
                <span class="input-group-addon">
                  <i class="fa fa-calendar"></i>
                </span>
                <div>
                <datepicker   
                  lang="en" 
                  v-model="newTaskInfo.start"
                  format="DD/MMM/YYYY"
                  width="100%"
                  data-vv-name="start"
                  v-validate="'required'"
                  :class="{ 'is-invalid': submitted &&  errors.has('start') }"
                  >
                </datepicker>
                <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
                  <div v-if="submitted && errors.has('start')" class="invalid-feedback sitLenMotTi">{{ errors.first('start') }}</div> 
                </transition>
                </div>
              </div>
            </div>
            <div class="col-xs-6" v-if="newTaskInfo.type === 'task'"> 
              <h4 class="title" >Estimate duration</h4>
              <div class="input-group">
                <span class="input-group-addon">
                  <i class="fa fa-hourglass-half"></i>
                </span>
                <div>
                <input
                  type="number"
                  min = 0          
                  v-validate="'required'"
                  v-model.number="newTaskInfo.duration"
                  class="form-control"
                  placeholder="Enter estimated duration"
                  name="Estimate duraion"
                  :class="{ 'is-invalid': submitted &&  errors.has('Estimate duraion') }"
                />
                <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
                  <div v-if="submitted && errors.has('Estimate duraion')" class="invalid-feedback">{{ errors.first('Estimate duraion') }}</div> 
                </transition>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="box-footer">
          <button class="btn-create button-modal pull-right" @click="handleSubmit(newTaskInfo, currentTask)"> Add task</button>
          <button class="btn-close button-modal" @click="closeModal"> Cancel</button>
        </div>
      </div>
    </modal>
</template>
<script>
import datepicker from 'vue2-datepicker'
import { EventBus } from '@/main.js'
export default {
    data() {
      return {
        currentTask: {},
        submitted: false,
        newTaskInfo: {
          parentId: '',
          id: '',
          label: '',
          parents: [],
          start: '',
          duration: '',
          type: 'task',
          collapse: true,
          endTime: ''
        },
        taskType: [
          { text: 'Task', value: 'task' },
          { text: 'Milestone', value: 'milestone' }
        ]
      }
    },
    components: { datepicker },
    shortcuts: [{
      onClick: () => {
        this.values = [ new Date(), new Date() ]
        }
      }
      ],
    methods: {
      closeModal() {
          this.$modal.hide('AddTask')
      },
      handleSubmit(newTaskInfo, currentTask) {
      this.submitted = true
      this.$validator.validate().then(valid => {
                if (valid) {
                    if (newTaskInfo.type === 'milestone') {
                      EventBus.$emit('addMilestone', newTaskInfo)
                    } else {
                      EventBus.$emit('addTask', newTaskInfo)
                      currentTask.user = ''
                      currentTask.type = 'project'
                    }
                      this.$modal.hide('AddTask')
                } else {
                  this.$modal.show('dialog', {
                  title: 'Error',
                  text: 'Invalid input' })
                }
            })
      },
      beforeOpen(event) {
          var d = new Date()
          var n = d.valueOf()
          this.currentTask = event.params.data
          this.newTaskInfo.id = n
          this.newTaskInfo.parentId = this.currentTask.id
          this.newTaskInfo.start = this.currentTask.start
          this.newTaskInfo.endTime = this.currentTask.endTime
      }
    },
    beforeDestroy() {
      EventBus.$off('addTask')
      EventBus.$off('addMilestone')
    }
}
</script>
<style src="vue-multiselect/dist/vue-multiselect.min.css">
</style>
