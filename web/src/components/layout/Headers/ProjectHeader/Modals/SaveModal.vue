<template>
    <modal name="SaveModal" transition="nice-modal-fade" 
        :draggable="true" 
        :reset="true"
        height="auto"
         @before-open="beforeOpen"
    >
      <div class="box box-group">
        <div class="box-header with-border dark">
          <h5 class="box-title2">Save project</h5>
        </div>
        <div class="box-body">
          <!--Input section-->
          <div class="input-group col-xs-12 hidden">
            <input
              type="number"
              v-model="snapshot.id"   
              width="100%"
              class="form-control"
            />
          </div>
          <h4 class="title col-xs-12" >Description</h4>
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-tag"></i>
            </span>
            <input
              type="text"
              v-validate="'required|min:3'"
              v-model="snapshot.description"
              class="form-control"
              placeholder="Enter description"
              name="description"
              :class="{ 'is-invalid':submitted &&  errors.has('description') }"
           
            />
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
              <div v-if="submitted && errors.has('description')" class="invalid-feedback">{{ errors.first('description') }}</div> 
            </transition>
          </div>
        </div>
        <div class="box-footer">
          <button class="btn-create button-modal pull-right" @click="handleSubmit(snapshot)">Save</button>
          <button class="btn-close button-modal" @click="closeModal">Cancle</button>
        </div>
      </div>
    </modal>
</template>
<script>
import { EventBus } from '@/main.js'

export default {
    name: 'SaveModal',
    props: ['projectId'],
    data() {
      return {
        submitted: false,
        snapshot: {
          description: '',
          projectId: '',
          tasks: [],
          updateDate: ''
        }
      }
    },
    methods: {
      closeModal() {
          this.$modal.hide('SaveModal')
      },
      handleSubmit(snapshot) {
        this.submitted = true
        this.$validator.validate().then(valid => {
                  if (valid) {
                    let d = new Date()
                    let n = d.valueOf()
                    this.snapshot.updateDate = n
                    EventBus.$emit('saveProject', snapshot)
                    this.$modal.hide('SaveModal')
                  } else {
                      this.$modal.show('dialog', {
                        title: 'Error',
                        text: 'Invalid input!'})
                  }
              })
      },
      beforeOpen() {
          this.snapshot.projectId = this.projectId
      }
    },
    beforeDestroy() {
      EventBus.$off('saveProject')
    }
}
</script>
<style scoped>

</style>
