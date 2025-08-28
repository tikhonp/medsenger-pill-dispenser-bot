import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Fill2x2DispenserView from '../views/Fill2x2DispenserView.vue'
import Fill4x7DispenserView from '../views/Fill4x7DispenserView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/confirm-layout',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/confirm-layout/fill-2x2-dispenser',
      name: '2x2',
      component: Fill2x2DispenserView,
    },
    {
      path: '/confirm-layout/fill-4x7-dispenser',
      name: '4x7',
      component: Fill4x7DispenserView,
    },
  ],
})

export default router
