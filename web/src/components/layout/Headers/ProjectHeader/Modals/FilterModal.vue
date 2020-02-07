<template>
  <div>
    <modal name="someFilter" transition="nice-modal-fade" 
        :height=600 
        :width=400 
        :draggable="true" 
        :reset="true"
        :pivotX=0.95
        :pivotY=0.4>
      <a class="pull-right exit-btn" @click="$modal.hide('someFilter')"><i class="fa fa-close"/></a>
      <div class="setting-content">
        <div>
          <h3 class="setting-heading">Filter</h3>
        </div>
        <div class="form-group">
          <label class="setting-subheading">Task name: </label>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-check"></i></span>
              <input
                class="form-control" 
                placeholder="Enter task's name" 
                type="text"/>
            </div>
        </div>

        <div class="form-group">
          <label class="setting-subheading">Asignee: </label>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-hand-paper-o"></i></span>
              <input
                class="form-control" 
                placeholder="Enter name of asignee" 
                type="text"/>
            </div>
        </div>

        <div class="form-group">
          <label class="setting-subheading">Start date: </label>
            <div class="input-group">
              <span class="input-group-addon"><i class="fa fa-fw fa-calendar-o"></i></span>
              <input
                class="form-control" 
                placeholder="Enter date" 
                type="text"/>
            </div>
        </div>

        <!-- <div class="btn-apply">
          <button class="btn btn-info btn-flat btn-lg">APPLY</button>
        </div> -->
      </div>

    </modal>
  </div>
</template>

<script>
export default {
  name: 'filterModal'
}
</script>

<style scoped>
.setting-content .btn-apply, .setting-heading{
  text-align: center
}
.setting-heading {
  margin-bottom: 60px
}
.form-group {
  margin: 20px !important;
  padding: 5px
}
.exit-btn {
    font-size: 15px;
    padding: 5px;
    color: #313233
  }
  .exit-btn:hover{
    color: #3fb0ac
  }
</style>

