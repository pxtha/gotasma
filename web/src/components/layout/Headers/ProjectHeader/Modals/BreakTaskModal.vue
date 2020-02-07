  <template>
  <modal
    name="breakTaskModal"
    transition="pop-out"
    :height="400"
    :width="500"
    :draggable="true"
    :reset="true"
    @before-open="beforeOpen"
    @before-close="beforeClose"
  >
    <a class="pull-right exit-btn">
      <i class="fa fa-close" @click="$modal.hide('breakTaskModal')"/>
    </a>
    <div class="modal-box">
      <div class="title">
        <i class="fa fa-fw fa-edit"></i>Break Task
      </div>

      <div class="col-xs-12">
        <h4 class="myheading">Name of Task:</h4>
        <div class="input-group">
          <span class="input-group-addon">
            <i class="fa fa-fw fa-pencil"></i>
          </span>
          <input
            disabled
            class="form-control"
            placeholder="Name of task"
            type="text"
            v-model="currentTask.label"
          />
        </div>
      </div>

      <div class="col-xs-12">
        <h4 class="myheading">Assignee:</h4>
        <div class="input-group col-xs-12">
          <span class="input-group-addon">
            <i class="fa fa-fw fa-child"></i>
          </span>
          <input disabled class="form-control" placeholder="Assignee" type="text" v-model="currentTask.user" />
        </div>
      </div>

      <div class="col-xs-12">
        <h4 class="myheading">Pick date</h4>
        <div class="myPicker">
          <datepicker
            v-model="breakTaskInfo.start"
            lang="en"
            format="DD/MMM/YYYY"
            width="100%"
            :editable="false"
          ></datepicker>
        </div>
      </div>

      <div class="button-set col-xs-12">
        <button class="button-modal pull-right" @click="applyEdit(currentTask)">Apply</button>
      </div>
    </div>
  </modal>
</template>
<script>
import datepicker from 'vue2-datepicker'
import { EventBus } from '@/main.js'
export default {
  name: 'BreakTaskModal',
  components: {
    datepicker
  },
  props: ['exceptionDays'],
  data() {
    return {
      currentTask: '',
      breakTaskInfo: {
        adjacentId: '', // ip lien truoc
        parentId: '',
        id: '',
        label: '',
        start: '',
        duration: '',
        type: 'task',
        effort: '',
        collapse: true
      }
    }
  },
  shortcuts: [
    {
      onClick: () => {
        this.currentTask.startDate = [new Date(), new Date()]
      }
    }
  ],
  methods: {
    beforeOpen(event) {
        let day = new Date()
        let number = day.valueOf()
        this.currentTask = event.params.data
        this.breakTaskInfo.parentId = this.currentTask.parentId
        this.breakTaskInfo.label = this.currentTask.label
        this.breakTaskInfo.id = number
        this.breakTaskInfo.effort = this.currentTask.effort
        this.breakTaskInfo.adjacentId = this.currentTask.id // id task truoc no - de chen
    },
    beforeClose() {
      this.breakTaskInfo.start = this.breakTaskInfo.start.valueOf()
    },
    applyEdit(task) {
      if (this.breakTaskInfo.start > task.start && this.breakTaskInfo.start < task.endTime) {
        this.$modal.hide('breakTaskModal')
        let durFirstHalf = (this.breakTaskInfo.start.valueOf() - task.start)
        let durFirstHalfTemp = durFirstHalf
        let loopOfFirst = durFirstHalf / 86400000

        let timeStart = new Date(task.startTime)
        let dayofWeek = timeStart.getDay()
        let calculateTimeChart = task.startTime
        // tinh ngay lam
        for (let i = 0; i < loopOfFirst; i++) {
          let isHoliday = false
            for (let j = 0; j < this.exceptionDays.length; j++) {
              if (calculateTimeChart === this.exceptionDays[j]) {
                isHoliday = true
                break
              }
            }
            if (isHoliday) {
              durFirstHalfTemp -= 86400000
              isHoliday = false
              if (dayofWeek === 6) {
                dayofWeek = 0
              } else {
                dayofWeek += 1
              }
            } else if (dayofWeek === 6) {
              dayofWeek = 0
              durFirstHalfTemp -= 86400000
            } else if (dayofWeek === 0) {
              dayofWeek += 1
              durFirstHalfTemp -= 86400000
            } else {
              dayofWeek += 1
            }
            calculateTimeChart += 86400000
        }
        // duration of SecondHalf
        this.breakTaskInfo.duration = task.estimateDuration - durFirstHalfTemp
        EventBus.$emit('breakTask', this.breakTaskInfo)

        // tinh real duration
        timeStart = new Date(task.startTime)
        dayofWeek = timeStart.getDay()
        calculateTimeChart = task.startTime

        let durationDays = durFirstHalfTemp / 86400000
        task.duration = durFirstHalfTemp
          for (let i = 0; i < durationDays; i++) {
            let isHoliday = false
            for (let j = 0; j < this.exceptionDays.length; j++) {
              if (calculateTimeChart === this.exceptionDays[j]) {
                isHoliday = true
                break
              }
            }
            if (isHoliday) {
              task.duration += 86400000
              isHoliday = false
              durationDays++
              if (dayofWeek === 6) {
                dayofWeek = 0
              } else {
                dayofWeek += 1
              }
            } else if (dayofWeek === 6) {
              dayofWeek = 0
              task.duration += 86400000
              durationDays++
            } else if (dayofWeek === 0) {
              dayofWeek += 1
              task.duration += 86400000
              durationDays++
            } else {
              dayofWeek += 1
            }
            calculateTimeChart += 86400000
          }
          task.estimateDuration = durFirstHalfTemp
      } else {
        this.$modal.show('dialog', {
          title: 'Date invalid',
          text: 'Date not in range of task'
        })
      }
    }
  },
  beforeDestroy() {
    EventBus.$off('breakTask')
  }
}
</script>

<style scoped>
.modal-box {
  background: white;
  color: black;
  font-size: 0;
}

.modal-box .title {
  box-sizing: border-box;
  padding: 20px 20px 30px 20px;
  width: 100%;
  text-align: center;
  letter-spacing: 1px;
  font-size: 23px;
  font-weight: 300;
}

.modal-box .input-group {
  padding-bottom: 20px;
}

.modal-box .button-set :hover {
  border-color: #3fb0ac;
  color: #3fb0ac;
}

.modal-box .button-set .btn-close:hover {
  border-color: #eb4b4b;
  color: #eb4b4b;
}
.modal-box .button-set {
  margin: 15px 0 15px 0;
}
.myheading {
  margin: 5px 0 !important;
}
.exit-btn {
  font-size: 15px;
  padding: 5px;
  color: #313233;
}
.exit-btn:hover {
  color: #3fb0ac;
}

.myPicker {
  padding-bottom: 10px;
}
</style>
