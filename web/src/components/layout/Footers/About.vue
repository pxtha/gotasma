<template>
    <modal name="about" :draggable="true" height="auto"> 
    <div class="box box-widget widget-user">
            <!-- Add the bg color to the header using any of the bg-* classes -->
            <div class="widget-user-header bg-black" style="background: url('../../../../static/img/about.jpg') center center;">
              <h3 class="widget-user-username">{{admin.fullName}}</h3>
              <h5 class="widget-user-desc">{{admin.country}}</h5>
            </div>
            <div class="widget-user-image">
              <img class="img-circle" :src="admin.aboutAvatar" alt="User Avatar">
            </div>
            <div class="box-footer">
              <div class="row">
                <div class="col-sm-4 border-right">
                  <div class="description-block">
                    <h5 class="description-header">200</h5>
                    <span class="description-text">PROjECT</span>
                  </div>
                  <!-- /.description-block -->
                </div>
                <!-- /.col -->
                <div class="col-sm-4 border-right">
                  <div class="description-block">
                    <h5 class="description-header">2000$</h5>
                    <span class="description-text">Income</span>
                  </div>
                  <!-- /.description-block -->
                </div>
                <!-- /.col -->
                <div class="col-sm-4">
                  <div class="description-block">
                    <h5 class="description-header">DC-14</h5>
                    <span class="description-text">TEAM</span>
                  </div>
                  <!-- /.description-block -->
                </div>
                <!-- /.col -->
              </div>
              <!-- /.row -->
            </div>
       </div>
    </modal>
</template>
<script>
export default {
  name: 'about',
  props: ['admin']
}
</script>
<style scoped>
.box-footer h5,span{
  color: #000;
}
.box-footer{
  margin-top: 15px;
}
.widget-user-username{
  padding-top: 40px;
  font-weight: 500;
  font-size: 2em;
}
.widget-user-desc{
  font-weight: 500;
  font-size: 18px;
}
</style>