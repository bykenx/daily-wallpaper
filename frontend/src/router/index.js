import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import More from '@/pages/More.vue'
import Local from '@/pages/Local.vue'


export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: 'Home' },
    { name: 'Home', path: '/home', component: Home },
    { name: 'More', path: '/more', component: More },
    { name: 'Local', path: '/local', component: Local },
  ],
})