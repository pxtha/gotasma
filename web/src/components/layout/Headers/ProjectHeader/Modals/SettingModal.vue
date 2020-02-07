<template>
  <div v-if="project">
    <modal name="settingModal" 
        transition="nice-modal-fade" 
        :height=600 
        :width=400 
        :draggable="true" 
        :reset="true"
        @before-open="beforeOpen"
        @before-close="beforeClose"
        :pivotX=0.95
        :pivotY=0.4>
      <a class="pull-right exit-btn" @click="$modal.hide('settingModal')"><i class="fa fa-close"/></a>
      <div class="setting-content">
        <div>
          <h3 class="setting-heading">Project Setting</h3>
        </div>
        <div class="form-group">
          <label class="setting-subheading">Project's name: </label>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-file-o"></i></span>
              <input
                class="form-control" 
                v-model="project.name"
                type="text"/>
            </div>
        </div>

        <div class="form-group">
          <label class="setting-subheading">Last update:</label>
          <div>
            <datepicker v-model="project.updateDate" 
            disabled
            lang="en" 
            format="MMM/DD/YYYY" 
            width="100%"
            :editable="false">
            </datepicker>
          </div>
        </div>
        
        <div class="form-group">
          <button class="button-modal btn-create pull-right" 
          @click="applyEdit(project)">APPLY</button>
        </div>
      </div>

    </modal>
  </div>
</template>

<script>
import datepicker from 'vue2-datepicker'
export default {
  name: 'settingModal',
  props: ['id', 'project'],
  components: {
    datepicker
  },
  data() {
    return {
      beforeEdit: '',
      isChanged: false
    }
  },
  methods: {
    applyEdit(project) {
      this.isChanged = true
      this.$store.dispatch('editProject', project)
      this.$modal.hide('settingModal')
    },
    beforeOpen() {
      this.beforeEdit = Object.assign({}, this.project)
    },
    beforeClose() {
      if (this.isChanged === false) {
        Object.assign(this.project, this.beforeEdit)
      }
    }
  }
}
</script>

<style scoped>
.setting-heading {
  margin-bottom: 50px;
  text-align: center
}
.form-group {
  margin: 20px !important;
  padding: 5px
}
.exit-btn {
  padding: 5px;
}
</style>

