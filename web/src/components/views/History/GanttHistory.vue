<template>
      <gantt-elastic
        v-if="snapshot.tasks && exceptionDays"
        :options="defaultOptions"
        :tasks="snapshot.tasks"
        :exceptionDays="exceptionDays">
      </gantt-elastic>
    <div v-else class="notfound" >
      <img class="notice" :src="this.adminInfo.avatar"/>
      <h3 class="col-sm-12">Choose history from list</h3>
    </div>
</template>
<script>
import { mapState, mapGetters } from 'vuex'
import dayjs from 'dayjs'
import GanttElastic from 'gantt-elastic'

export default {
  name: 'GanttHistory',
  props: ['snapshot'],
  data() {
    return {
      defaultOptions: {
        scope: {
            before: 1,
            after: 80
        },
        maxRows: 1000,
        maxHeight: 480,
        times: {
            timeZoom: 21
        },
        row: {
            height: 25
        },
        calendar: {
            hour: {
                display: false
            },
            workingDays: [1, 2, 3, 4, 5]
        },
        chart: {
            text: {
                display: false
            },
            expander: {
                display: true
            }
        },
        taskList: {
            expander: {
                straight: true
            },
            columns: [{
                    id: 2,
                    label: 'Description',
                    value: 'label',
                    width: 200,
                    expander: true,
                    style: {
                        'task-list-item-value-container': { 'font-weight': 'bold' },
                        'cursor': 'pointer'
                    }
                },
                {
                    id: 3,
                    label: 'Assignee',
                    value: task => `${this.getMember(task)}`,
                    width: 90,
                    html: true
                },
                {
                    id: 3,
                    label: 'Start',
                    value: task => dayjs(task.start).format('DD-MM-YYYY'),
                    width: 78
                },
                {
                    id: 4,
                    label: 'Duration (estimated)',
                    value: task => (task.estimateDuration / 86400000) + 'd',
                    // value: task => dayjs(task.endTime).format('DD-MM-YYYY'),
                    width: 45,
                    style: {
                        'task-list-header-label': { 'text-align': 'center' }
                    }
                },
                {
                    id: 5,
                    label: 'Duration (real)',
                    value: task => task.duration / 86400000 + 'd',
                    // value: task => dayjs(task.endTime).format('DD-MM-YYYY'),
                    width: 45,
                    style: {
                        'task-list-header-label': { 'text-align': 'center' }
                    }
                }
            ]
        }
      }
    }
  },
  components: { GanttElastic },
  computed: {
    ...mapState([
      'resources',
      'adminInfo'
      ],
    ),
    ...mapGetters([
      'exceptionDays'
    ])
  },
  created() {
    this.$store.dispatch('getResources')
    this.$store.dispatch('getExceptions')
  },
  methods: {
    getMember(task) {
      let arrName = []
      for (let i = 0; i < this.resources.length; i++) {
        for (let j = 0; j < task.user.length; j++) {
          if (this.resources[i].id === task.user[j]) {
            arrName.push(this.resources[i].name)
          }
        }
      }
      return arrName
     }
  }
}
</script>
<style scoped>
.notice {
  width: 240px;
}
.notfound {
  padding-left: 20em;
}
</style>