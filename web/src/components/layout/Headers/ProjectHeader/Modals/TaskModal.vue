  <template>
  <modal name="taskModal" transition="pop-out" 
      height="auto" 
      :scrollable="true"
      :width=500 
      :draggable="true" 
      :reset="true"
      :clickToClose="false"
      @before-open="beforeOpen" 
       >
      <a class="pull-right exit-btn" @click="cancelEdit"><i class="fa fa-close"/></a>
      <div class="modal-box">
        <div class="title"><i class="fa fa-fw fa-edit"></i> Task info</div>
        
          <div class="col-xs-12">
            <h4 class="myheading">Name of Task: </h4>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-pencil"></i></span>
                <input
                  class="form-control" 
                  placeholder="Name of task" 
                  type="text"
                  v-model="newTaskInfo.label"
                  >
            </div>
          </div>
          <div class="col-xs-12" v-if="currentTask.type === 'task' ||currentTask.type === 'milestone' ">
            <h4 class="myheading">Start date</h4>
            <div class="myPicker">
              <datepicker
              v-model="newTaskInfo.start"
              lang="en" 
              format="DD/MMM/YYYY"
              width="100%"
              :editable="false">
              </datepicker>
            </div>
          </div>
          <div class="col-xs-6" v-if="currentTask.type === 'task'">
            <h4 class="myheading">Duration (estimated)</h4>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-hourglass-3"></i></span>
                <input
                  class="form-control" 
                  placeholder="Duration"
                  type="number"
                  min="0"
                  v-model.number="newTaskInfo.estimateDuration">
            </div>
          </div>

          <div class="col-xs-6" v-if="currentTask.type === 'task'">
            <h4 class="myheading">Effort</h4>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-check-square-o"></i></span>
                <input
                  class="form-control" 
                  placeholder="Effort"
                  type="number"
                  min="0"
                  max="100"
                  v-model.number="newTaskInfo.effort">
              </div>
          </div>

          <div class="col-xs-12">
            <h4 class="myheading">Type</h4>
            <div class="input-group col-xs-12">
              <span class="input-group-addon">
                <i class="fa fa-fw fa-hand-pointer-o"></i>
              </span>
              <select
              disabled
              class="form-control"
              v-model="currentTask.type">
                <option v-for="type in taskType" :value="type.value" :key="type">
                    {{ type.text }}
                </option>
              </select>
            </div>          
          </div>

          <div class="col-xs-6">
              <h4 class="myheading">Duration (real-time)</h4>
              <div class="input-group">
                <span class="input-group-addon"><i class="fa fa-fw fa-hourglass"></i></span>
                  <input
                    class="form-control" 
                    disabled
                    v-model="newTaskInfo.duration">
              </div>
            </div>

            <div class="col-xs-6">
              <h4 class="myheading">End date</h4>
              <div>
                <datepicker
                v-model="currentTask.endTime"
                lang="en" 
                format="DD/MMM/YYYY" 
                width="100%"
                :editable="false"
                disabled>
                </datepicker>
              </div>
            </div>

          <div class="button-set col-xs-12">
            <button class="btn-create button-modal pull-right" @click="applyEdit(newTaskInfo)">Apply</button>
            <button class="btn-close button-modal" @click="deleteTask(currentTask.id)">Delete</button>
          </div>

      </div>
  </modal>
</template>
<script>
import datepicker from 'vue2-datepicker'
import { EventBus } from '@/main.js'
import Multiselect from 'vue-multiselect'

export default {
  name: 'TaskModal',
  components: { datepicker, Multiselect },
  props: ['exceptionDays'],
  data() {
    return {
      currentTask: {},
      newTaskInfo: {
          id: '',
          label: '',
          start: '',
          duration: '',
          estimateDuration: '',
          effort: '',
          type: ''
      },
      beforeEdit: '',
      taskType: [
        { text: 'Parent Task', value: 'project' },
        { text: 'Task', value: 'task' },
        { text: 'Milestone', value: 'milestone' }
      ]
    }
  },
  shortcuts: [{
      onClick: () => {
        this.currentTask.startDate = [ new Date(), new Date() ]
      }
    }
  ],
  methods: {
    beforeOpen(event) {
      this.currentTask = event.params.data
      this.newTaskInfo.type = this.currentTask.type
      this.newTaskInfo.label = this.currentTask.label
      this.newTaskInfo.start = this.currentTask.start
      this.newTaskInfo.duration = this.currentTask.duration / 86400000
      this.newTaskInfo.estimateDuration = this.currentTask.estimateDuration / 86400000
      this.newTaskInfo.effort = this.currentTask.effort
      this.newTaskInfo.id = this.currentTask.id
      this.beforeEdit = Object.assign({}, this.currentTask)
    },
    applyEdit(newTaskInfo) {
        EventBus.$emit('editTask', newTaskInfo)
        this.$modal.hide('taskModal')
    },
    deleteTask(idTask) {
      this.$modal.show('dialog', {
        title: 'Are you sure?',
        text: 'This task will be deleted permanantly',
        buttons: [
          {
            title: 'OK',
            default: true,
            handler: () => {
              EventBus.$emit('deleteThisTask', idTask)
              this.$modal.hide('dialog')
              this.$modal.hide('taskModal')
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
    cancelEdit() {
      Object.assign(this.currentTask, this.beforeEdit)
      this.$modal.hide('taskModal')
    }
  },
  beforeDestroy() {
    EventBus.$off('deleteThisTask')
    EventBus.$off('editTask')
  }
}
</script>
<style>
  .multiselect__single {
    color: #000;
  }
</style>
<style scoped>
  .modal-box {
    background: white;
    color: black;
    font-size: 0;
  }

  .modal-box .title {
    box-sizing: border-box;
    padding: 20px;
    width: 100%;
    text-align: center;
    letter-spacing: 1px;
    font-size: 23px;
    font-weight: 300;
  }

  .modal-box .input-group{
    padding-bottom: 10px
  }

  .modal-box .button-set {
    margin: 15px 0 15px 0;
  }
  .myheading {
    margin: 5px 0 !important
  }
  .exit-btn {
    font-size: 15px;
    padding: 5px;
    color: #313233
  }
  .exit-btn:hover{
    color: #3fb0ac
  }

  .myPicker {
    padding-bottom: 10px
  }
</style>
