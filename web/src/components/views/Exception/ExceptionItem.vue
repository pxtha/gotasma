<template>

<div>
  <transition-group enter-active-class="animated fadeInRight" leave-active-class="animated fadeOutRight">
    <div class="col-md-3" v-for="exception in exceptions" :key="exception.id">
    <div class="box box-success">
        <div class="box-header with-border">
          <h3 class="box-title"> {{ exception.title }} </h3>
      
          <div class="box-tools pull-left">
            <button type="submmit" class="btn btn-box-tool" @click="showDialog(exception.id)">
              <i class="fa fa-times"></i>
            </button>
          </div>
        </div>

        <div class="box-body">From:&nbsp;{{ exception.date[0] | momentNormalDate}}</div>
        <div class="box-body">To:&nbsp;{{ exception.date[1] | momentNormalDate}}</div>
      </div>
    </div>
  </transition-group>
  <v-dialog/>
</div>
</template>

<script>
import { mapState } from 'vuex'
export default {
    name: 'ExceptionItem',
    created() {
      this.$store.dispatch('getExceptions')
    },
    computed: {
    ...mapState([
      'exceptions'
    ])
    },
    methods: {
    showDialog(id) {
      this.$modal.show('dialog', {
        title: 'Are you sure?',
        text: 'Do you wish to delete?',
        buttons: [
          {
            title: 'OK',
            handler: () => {
              console.log(id)
              this.$store.dispatch('deleteExceptions', id)
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
.box-body {
  font-size: 16px !important;
  color: #000
}
</style>