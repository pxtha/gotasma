<template>
  <li class="user-menu">
    <a @click="$modal.show('about')" class="dropdown-toggle specialalala">
      <!-- The user image in the navbar-->
      <!-- hidden-xs hides the username on small devices so only the image appears. -->
      <img :src="admin.avatar" class="pull-left avatar"/>
      <div class="pull-left info">
          <div>
            <span class="hidden-xs">{{admin.fullName}}</span>
          </div>
            <i class="fa fa-circle text-success"></i> <span>{{admin.role}}</span>
      </div>
    </a>
  </li>
</template>

<script>

export default {
  name: 'AdminMenu',
  props: ['admin']
}
</script>
<style scoped>
  .dropdown-toggle .user-image{
    border-radius: 5px !important;
    width: 30px;
    height: 40px;
 
  }
  .hidden-xs{
    font-size: 20px;
  }
  .avatar {
    margin: 5px;
    margin-top: 2px;
    border-radius: 8px;
    width: 40px;
  }
  .myColor {
    color: #587477 !important
  }
  .specialalala {
    cursor: pointer;
    padding: 0px !important;
    margin: 5px;
  }
</style>
