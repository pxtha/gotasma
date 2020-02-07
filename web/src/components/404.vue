<template>
  <div id="notFound">
    <img src="/static/img/stock/giphy.gif" class="center-block logo2">

    <div class="text-center col-sm-12">
      <h1>You are lost.</h1>
      <h4>This page doesn't exist.</h4>
      <router-link to="/" class="vertical-5p lead">Take me home.</router-link>
    </div>
  </div>
</template>
<script>
export default {
  name: 'NotFound'
}
</script>

<style scoped>
.logo2 {
  width: 250px
}
#notFound {
  padding: 10em;
 
}
</style>
