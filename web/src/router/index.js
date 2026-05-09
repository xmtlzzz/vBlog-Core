import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../blog/Home.vue') },
  { path: '/post/:id', component: () => import('../blog/Post.vue') },
  { path: '/archives', component: () => import('../blog/Archives.vue') },
  { path: '/tags', component: () => import('../blog/Tags.vue') },
  { path: '/about', component: () => import('../blog/About.vue') },
  { path: '/admin/login', component: () => import('../admin/Login.vue') },
  {
    path: '/admin',
    component: () => import('../admin/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', component: () => import('../admin/Dashboard.vue') },
      { path: 'posts', component: () => import('../admin/Posts.vue') },
      { path: 'tags', component: () => import('../admin/Tags.vue') },
      { path: 'comments', component: () => import('../admin/Comments.vue') },
      { path: 'custom', component: () => import('../admin/Custom.vue') },
      { path: 'settings', component: () => import('../admin/Settings.vue') },
    ]
  }
]

const router = createRouter({ history: createWebHistory(), routes })
router.beforeEach((to, from, next) => {
  if (to.matched.some(r => r.meta.requiresAuth)) {
    if (!localStorage.getItem('vblog-token')) next('/admin/login')
    else next()
  } else next()
})
export default router
