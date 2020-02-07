<template>
  <modal name="createNewProj" transition="pop-out" height="auto"
  :draggable="true" 
  @before-open="beforeOpen">
      <div class="box box-group">
        <div class="box-header with-border dark">
          <h3 class="box-title">CREATE NEW PROJECT</h3>
        </div>
        <div class="box-body">
          <!--Input section-->
          <h4 class="title col-xs-12" >Project name</h4>
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-fw fa-check"></i>
            </span>
                <input v-model="project.name" 
                  class="form-control" 
                  placeholder="Name of project" 
                  type="text" 
                  v-validate="'required|min:5'" 
                  name="Project Name" 
                  :class="{'is-invalid':submitted && errors.has('Project Name')}">
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
              <div class="invalid-feedback" v-if="submitted && errors.has('Project Name')">{{ errors.first('Project Name')}}</div>
            </transition>
          </div>
          <h4 class="title col-xs-12" >Woring days</h4>
          <div class="input-group col-xs-12" disable>
              <span class="input-group-addon"><i class="fa fa-fw fa-calendar"></i></span>
                  <multiselect 
                      disabled
                      name="workingDays"
                      v-model="dayOfWeek"
                      :options="dayOfWeek"
                      :multiple="true"
                      label="name"
                    >
                  </multiselect>
          </div>
          <h4 class="title col-xs-12" >Start date</h4>
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-envelope"></i>
            </span>
            <datepicker v-model="project.startDate" 
              lang="en" 
              format="MMM/DD/YYYY" 
              width="100%"
              data-vv-name="Start Date"
              v-validate="'required'"
              :class="{'is-invalid':submitted && errors.has('Start Date')}"
              :editable="false">
            </datepicker>
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
                <div class="invalid-feedback specialalala" v-if="submitted && errors.has('Start Date')">{{ errors.first('Start Date')}}</div>
            </transition>
          </div>
        </div>
        <div class="box-footer">
          <button class="btn-create button-modal pull-right" @click="createProject">Create</button>
          <button class="btn-close button-modal" @click="cancelCreate">Cancel</button>
        </div>
      </div>
  </modal>
</template>
<script>
import datepicker from 'vue2-datepicker'
import Multiselect from 'vue-multiselect'

export default {
  name: 'NewProject',
  components: {
    datepicker, Multiselect
  },
  data() {
    return {
      submitted: false,
      project: {
        highlighted: false,
        name: '',
        startDate: '',
        updateDate: '',
        tasks: [],
        users: []
      },
      dayOfWeek: [{
          name: 'Monday',
          value: '1'
        },
        {
          name: 'Tuesday',
          value: '2'
        },
        {
          name: 'Wednesday',
          value: '3'
        },
        {
          name: 'Thursday',
          value: '4'
        },
        {
          name: 'Friday',
          value: '5'
        }]
    }
  },
  shortcuts: [{
      onClick: () => {
        this.project.startDate = [ new Date(), new Date() ]
        }
    }
  ],
  methods: {
    createProject() {
      this.submitted = true
      this.$validator.validateAll()
      .then(result => {
        if (result) {
          this.project.startDate = (this.project.startDate).valueOf()
          this.$store.dispatch('addProject', this.project)
          this.$modal.hide('createNewProj')
        }
      })
      .catch(error => {
        console.log(error)
      })
    },
    beforeOpen() {
        let d = new Date()
        d.setHours(0, 0, 0, 0)
        let n = d.valueOf()
        this.project.name = ''
        this.project.effort = ''
        this.project.startDate = ''
        this.project.updateDate = n
        this.project.tasks = []
        this.project.users = []
    },
    cancelCreate() {
      this.$modal.hide('createNewProj')
    }
  }
}
</script>

<style scoped>
.specialalala {
  margin-top: 0
}
.box-title {
  padding: 5px;
  letter-spacing: 1px;
  font-family: "Open Sans", sans-serif;
  font-weight: 400;
  color: #313233;
  text-transform: uppercase;
  transition: 0.1s all;
  font-size: 16px;
}
</style>
