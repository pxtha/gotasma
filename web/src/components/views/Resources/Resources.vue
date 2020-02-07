<template>
  <section class="content">
    <div class="row center-block">
      <div class="header col-md-12">
        <span>Team and resources</span> <br>
          <button type="button" @click="$modal.show('newresource')" class="btn btn-info pull-right special">Invite new member</button>
        <i>Here you manage people from all your projects. You can also invite new members by email.</i>
              
      </div>
            <new-resource></new-resource>
            <resource-table></resource-table>
      </div>
  </section>
</template>

<script>
import ResourceTable from './ResourceTable'
import NewResource from './NewResource'

export default {
  name: 'resource',
  components: {
    ResourceTable,
    NewResource
  }
}
</script>
<style scoped>
.special{
  margin-bottom: 10px;
}
.content {
  margin: 20px !important;
  min-height: 1102px !important
}
.header span {
  font-size: 24px;
  color: #242e35;
  font-weight: bold;
}
.header p {
  font-size: 14px;
  font-style: italic;
  color: #242e35;
  margin-bottom: 20px;
}

.email{
  height: 46px;
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px
}
</style>
