import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/Login.vue';
import SchedulerTable from '../components/SchedulerTable.vue';
import RemoteComputers from '../components/RemoteComputers.vue';
import authService from '../services/authService';

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
  {
    path: '/remote-computers',
    name: 'RemoteComputers',
    component: RemoteComputers,
    meta: { requiresAuth: true }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard
router.beforeEach(async (to, from, next) => {
  // Check if route requires authentication
  if (to.matched.some(record => record.meta.requiresAuth !== false)) {
    // Check authentication status with the backend
    const isAuthenticated = await authService.isAuthenticated();
    
    if (!isAuthenticated) {
      // Not authenticated, redirect to login
      next({ name: 'Login', query: { redirect: to.fullPath } });
      return;
    }

    // Authenticated, proceed
    next();
  } else {
    next();
  }
});

export default router;
