<template>
    <div class="main">
         <router-view :id="id"></router-view>
    </div>
</template>
<script>
export default {
  name: 'Project',
  props: ['id']
}
</script>