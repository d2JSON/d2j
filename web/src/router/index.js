import { createRouter, createWebHistory } from 'vue-router'
import ConnectionView from '../views/ConnectionView.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/connection',
      name: 'connection',
      component: ConnectionView,
    },
  ]
})

export default router
