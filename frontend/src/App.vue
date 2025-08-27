<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class='container'>
    <router-view />
  </div>
</template>

<style scoped>
.container {
  padding: 15px;
}
</style>
