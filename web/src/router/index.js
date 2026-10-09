import { createRouter, createWebHistory } from 'vue-router'
import { isTokenExpired } from '../utils/auth'

const routes = [
  { path: '/', component: () => import('../blog/Home.vue') },
  { path: '/post/:id', component: () => import('../blog/Post.vue') },
  { path: '/archives', component: () => import('../blog/Archives.vue') },
  { path: '/modules', component: () => import('../blog/Modules.vue') },
  { path: '/tags', component: () => import('../blog/Tags.vue') },
  { path: '/about', component: () => import('../blog/About.vue') },
  { path: '/admin/login', component: () => import('../admin/Login.vue') },
  { path: '/admin/register', redirect: '/admin/login' },
  {
    path: '/admin',
    component: () => import('../admin/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', component: () => import('../admin/Dashboard.vue') },
      { path: 'posts', component: () => import('../admin/Posts.vue') },
      { path: 'posts/new', component: () => import('../admin/EditPost.vue'), meta: { breadcrumb: '新建文章' } },
      { path: 'posts/:id/edit', component: () => import('../admin/EditPost.vue'), meta: { breadcrumb: '编辑文章' } },
      { path: 'tags', component: () => import('../admin/Tags.vue') },
      { path: 'comments', component: () => import('../admin/Comments.vue') },
      { path: 'custom', component: () => import('../admin/Custom.vue') },
      { path: 'trash', component: () => import('../admin/Trash.vue') },
      { path: 'settings', component: () => import('../admin/Settings.vue') },
    ]
  }
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  if (to.matched.some(r => r.meta.requiresAuth)) {
    const token = localStorage.getItem('vblog-token')
    if (!token || isTokenExpired(token)) {
      if (token) {
        localStorage.removeItem('vblog-token')
        sessionStorage.setItem('vblog_auth_expired_msg', '登录已过期，请重新登录')
      }
      next({
        path: '/admin/login',
        query: to.fullPath && to.fullPath !== '/admin' ? { redirect: to.fullPath } : {}
      })
      return
    }
    next()
  } else {
    next()
  }
})

export default router
