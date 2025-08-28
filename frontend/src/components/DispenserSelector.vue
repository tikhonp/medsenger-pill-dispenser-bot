<script setup>

import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class='chooser-page'>
    <div class='dispenser-chooser shadow p-4 rounded'>
      <h5 class='mb-3 text-primary fw-bold text-center'>Выберите размер таблетницы...</h5>
      <router-link class='btn btn-info btn-block' to='/confirm-layout/fill-2x2-dispenser'>2 x 2</router-link>
      <router-link class='btn btn-info btn-block ' to='/confirm-layout/fill-4x7-dispenser'>4 x 7</router-link>
    </div>
  </div>

</template>

<style scoped>
.chooser-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.dispenser-chooser {
  display: flex;
  flex-direction: column;
  max-width: 500px;
  background-color: #ffffff;
  border-radius: 6px;
  animation: fade-in 0.3s ease;
}

.btn.btn-primary {
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.btn.btn-primary:hover {
  box-shadow: 0 4px 12px rgba(0, 108, 136, 0.2);
  transform: translateY(-1px);
}

h5 {
  margin-top: 0;
}
</style>
