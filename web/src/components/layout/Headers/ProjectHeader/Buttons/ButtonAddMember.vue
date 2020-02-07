<template>
  <li class="dropdown messages-menu" v-if="users">
    <a href="javascript:;" title="Add members to project" class="dropdown-toggle" data-toggle="dropdown">
      <span class="fa  fa-street-view myIcon"></span>
      <span class="label special" >{{ users | count}}</span>
    </a>
     <ul class="dropdown-menu">
      <li class="header">
         There are {{ users | count }} member(s) in your project
      </li>
      <li>
        <ul class="menu">
          <member-modal v-for="user in users" :key="user.id" :user="user"></member-modal>
          <li v-if=" users.length == 0">
            <span class="center-block text-center alert-error">Please assign member to your project</span>
          </li>
        </ul>
      </li>
      <li class="footer">
        <router-link :to="'/project/' + id + '/member'">
        <a >
          <i class="fa fa-users"></i>
          See all member of project
        </a>
        </router-link>
      </li>
    </ul>
</li>
</template>

<script>
import MemberModal from '../Modals/MemberModal.vue'

export default {
  name: 'buttonAddMember',
  components: { MemberModal },
  props: ['id', 'users']
}
</script>

<style scoped>
.special {
  padding: 0px;
  margin: 0px;
  text-decoration-line: underline
}
.myIcon {
  color: rgb(43, 42, 42);
  margin: 0 5px 0 5px
}
.dropdown-toggle:hover {
  background-color: #f3f5f7 !important;
  border-radius: none !important
}
</style>

