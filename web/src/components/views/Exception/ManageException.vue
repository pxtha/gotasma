<template>
  <section class="content">
    <div class="row center-block">
      <div class="header col-md-12">
        <div class="mySpacing">
          <span>Holidays, Exceptions Management</span>
          <button
          type="submit"
          class="btn btn-lg btn-info pull-right special">
          Recalculated</button>
          <br />
          <i>Here you can manage exception date like Holiday, day-off, etc. The form below help you to create an exception.</i>
        </div>
        <div class="box box-primary">
          <div class="box-header with-border">
            <h3 class="box-title">Create Exception Date</h3>
          </div>
          <!-- /.box-header -->
          <!-- form start -->
            <div class="box-body">
              <div class="form-group">
                <label for="excepttitle">Title</label>
                <input
                  type="text"
                  class="form-control"
                  id="exceptTitle"
                  placeholder="Enter Title"
                  name="Title"
                  v-validate="'required|min:5'" 
                  v-model="exceptDate.title"
                  :class="{ 'is-invalid':submitted &&  errors.has('Title') }"
                />
                <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
                  <div class="invalid-feedback special" v-if="submitted && errors.has('Title')">{{ errors.first('Title')}}</div>
                </transition>
              </div>
              

              <div class="form-group">
                <label>Choose Date</label>
                <datepicker
                  v-model="exceptDate.date"
                  :editable="false"
                  range
                  lang="en"
                  format="DD/MMM/YYYY"
                  width="100%"
                  data-vv-name="Date"
                  v-validate="'required'"
                  :class="{ 'is-invalid':submitted &&  errors.has('Date') }"></datepicker>
                <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
                  <div class="invalid-feedback special" v-if="submitted && errors.has('Date')">{{ errors.first('Date')}}</div>
                </transition>
              </div>
            </div>
            <!-- /.box-body -->

            <div class="box-footer">
              <button type="submit" class="btn btn-info" @click="addException">Submit</button>
            </div>
        </div>
        <h2>Exceptions</h2>
        <exception-item></exception-item>
      </div>
    </div>
  </section>
</template>


<script>
import moment from 'moment'
import datepicker from 'vue2-datepicker'
import ExceptionItem from './ExceptionItem'

export default {
  name: 'manageExceptions',
  components: {
    datepicker,
    ExceptionItem
  },
  data() {
    return {
         submitted: false,
      exceptDate: {
        title: '',
        date: ''
      }
    }
  },
   mounted() {
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case 'ADD_EXCEPTION':
          this.$store.dispatch('getExceptions')
          break
        case 'DELETE_EXCEPTION':
          this.$store.dispatch('getExceptions')
         break
      }
    })
  },
  methods: {
    addException() {
      this.submitted = true
      this.$validator.validateAll().then(result => {
          if (result) {
            this.exceptDate.date[0] = moment(this.exceptDate.date[0]).valueOf()
            this.exceptDate.date[1] = moment(this.exceptDate.date[1]).valueOf()
            this.$store.dispatch('addExceptions', this.exceptDate)
            this.exceptDate.title = ''
            this.exceptDate.date = ''
          }
      })
  },
  shortcuts: [
    {
      onClick: () => {
        this.exceptDate.date = [new Date(), new Date()]
      }
    }
  ]
}
}
</script>

<style scoped>
span {
  font-size: 24px;
  color: #242e35;
  font-weight: bold;
}
.mySpacing {
  padding-bottom: 20px
}
.form-group {
  margin-bottom: 25px !important
}
.special {
  margin-top: 0px
}
</style>
