<template>
  <footer class="main-footer">
    <strong>	&copy; {{ year }}
    <span>TMA solution  </span>DC-14</strong>
    <span class="right-align">Version 1.0.0 || <a @click="$modal.show('about')" class="about"><i class="fa fa-fw fa-gg"></i>About Us</a></span>
    <about :admin="admin"></about>
  </footer>
</template>

<script>
import About from './About'

export default {
  name: 'DashFooter',
  props: ['admin'],
  components: {
    About
  },
  data: function () {
    return {
      year: new Date().getFullYear()
    }
  }
}
</script>

<style scoped>
span {
color: #173e43;
}
.right-align{
  position: absolute;
  right: 0;
}
.about:hover{
  cursor: pointer
}
</style>
