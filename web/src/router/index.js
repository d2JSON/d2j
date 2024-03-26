import { createRouter, createWebHistory } from 'vue-router'
import ConnectionView from '../views/ConnectionView.vue'
import ConvertView from '../views/ConvertView.vue'
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
    {
      path: '/convert',
      name: 'convert',
      component: ConvertView,
    },
  ]
})

export default router
