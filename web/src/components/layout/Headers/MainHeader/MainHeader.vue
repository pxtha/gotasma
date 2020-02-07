<template>
  <header class="main-header">
    <span class="logo-mini">
       <router-link to="/">
      <a href="/">
        <img
          src="/static/img/kitty.svg"
          alt="Logo"
          class="img-responsive center-block logo"
        />
      </a>
      </router-link>
    </span>
    <nav class="navbar navbar-static-top" role="navigation">
            <a href="javascript:;" class="sidebar-toggle" data-toggle="offcanvas" role="button">
        <span class="sr-only">Toggle navigation</span>
      </a>
    <div class="navbar-custom-menu">
        <ul class="nav navbar-nav">
          <admin-menu :admin="admin"></admin-menu>
        </ul>
    </div>
    </nav>

  </header>
</template>
<script>
import AdminMenu from './AdminMenu'

export default {
  name: 'dash-header',
  components: { AdminMenu },
  props: ['admin'],
  data() {
    return {
    }
  }
}
</script>
