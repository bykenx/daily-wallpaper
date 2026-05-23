import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import More from '@/pages/More.vue'
import Local from '@/pages/Local.vue'
import Settings from '@/pages/Settings.vue'


export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/today' },
    { name: 'Home', path: '/today', component: Home },
    { name: 'More', path: '/more', component: More },
    { name: 'Local', path: '/local', component: Local },
    { name: 'Settings', path: '/settings', component: Settings },
  ],
})