<template>
  <li class="notification-item">
    <a href="javascript:;">
      <div class="pull-left">
         <avatar :username="user.name" :size="30" ></avatar>
      </div>
      <h4>
        <span class="name">{{ user.name }}</span>
        <span class="time pull-right">
          <i class="fa fa-tags"></i>
          <span class="name">{{ user.badgeID }}</span>
        </span>
      </h4>
      <p>{{ user.email }}</p>
    </a>
  </li>
</template>

<script>
import Avatar from 'vue-avatar'

export default {
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  components: { Avatar }
}
</script>

<style scoped>
.special-fa{
  font-size: 25px;
}
.name { 
    color: black !important
}
.fa-tags{
    color: green !important
}
.navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > h4, .navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > p {
  margin: 0;
  white-space: normal;
  word-break: break-word;
}
.navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > h4 {
  font-size: 15px;
}
.navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > h4 small {
  font-size: 10px;
}
.navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > p {
  font-size: 12px;
}
.navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > h4 > span {
  display: inline-block;
}
.navbar-custom-menu > .navbar-nav > li.notifications-menu > .dropdown-menu > li > ul.menu li > a > h4 > small.time {
  margin-top: 3px;
}
</style>
