import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/Login.vue';
import SchedulerTable from '../components/SchedulerTable.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: SchedulerTable,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard
router.beforeEach((to, from, next) => {
  // Check if route requires authentication
  if (to.matched.some(record => record.meta.requiresAuth !== false)) {
    // Check if JWT cookie exists
    const cookies = document.cookie.split(';');
    const hasJwtCookie = cookies.some(cookie => cookie.trim().startsWith('jwt='));
    if (!hasJwtCookie) {
      // No JWT cookie found, redirect to login
      next({ name: 'Login', query: { redirect: to.fullPath } });
      return;
    }

    // JWT cookie exists, verify it
    next();
  } else {
    next();
  }
});

export default router;
