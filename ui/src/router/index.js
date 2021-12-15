import {createRouter, createWebHistory, createWebHashHistory} from 'vue-router';
import config from '@/config/config'

const router = createRouter({
  history: createWebHistory(process.env.NODE_ENV === "production" ? config.publicPath : ""), // 路由的history模式
  routes: [
    {
      path: "/",
      name: "index",
      component: () => import('@/view/index/index.vue'),
      meta: {
        title: "index"
      }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/view/404/404.vue')
    },

  ],
})

export default router;
