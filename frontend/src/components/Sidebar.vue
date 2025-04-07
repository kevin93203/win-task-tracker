<template>
  <div>
    <!-- Sidebar -->
    <div class="sidebar bg-gray-200 shadow-lg h-screen fixed left-0 top-0 transition-all duration-300 z-20" :class="{ 'w-64': isOpen, 'w-16': !isOpen }">
      <!-- Logo/Brand -->
      <div class="p-4 border-b border-gray-200 flex justify-between items-center">
        <h2 class="text-xl font-bold text-gray-800" :class="{ 'hidden': !isOpen }">任務管理系統</h2>
        <!-- Mobile close button (remains) -->
        <button @click="toggleSidebar" class="md:hidden text-gray-500 hover:text-gray-700">
          <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <!-- Desktop toggle button (moved here) -->
        <button
          @click="toggleSidebar"
          class="hidden md:block text-gray-500 hover:text-gray-700 p-1 rounded hover:bg-gray-100"
          title="Toggle Sidebar"
        >
          <svg
            class="h-6 w-6 transition-transform duration-300 flex-shrink-0"
            :class="{ 'rotate-180': !isOpen }"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>

      <!-- Navigation Links -->
      <nav class="p-4">
        <router-link
          to="/"
          class="flex items-center text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-lg mb-2"
          :class="[
            { 'bg-blue-50 text-blue-600': $route.path === '/' },
            isOpen ? 'px-4 py-2' : 'p-3 justify-center'
          ]"
          :title="isOpen ? '' : '排程工作'"
        >
          <svg class="h-5 w-5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM11 13a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
          <span :class="{ 'hidden': !isOpen, 'ml-3': isOpen }">排程工作</span>
        </router-link>
        
        <router-link
          to="/remote-computers"
          class="flex items-center text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-lg mb-2"
          :class="[
            { 'bg-blue-50 text-blue-600': $route.path === '/remote-computers' },
            isOpen ? 'px-4 py-2' : 'p-3 justify-center'
          ]"
           :title="isOpen ? '' : '遠端電腦管理'"
       >
          <svg class="h-5 w-5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M3 5a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2h-2.22l.123.489.804.804A1 1 0 0113 18H7a1 1 0 01-.707-1.707l.804-.804L7.22 15H5a2 2 0 01-2-2V5zm5.771 7H5V5h10v7H8.771z" clip-rule="evenodd" />
          </svg>
          <span :class="{ 'hidden': !isOpen, 'ml-3': isOpen }">遠端電腦管理</span>
        </router-link>
      </nav>

      <!-- Bottom Section with Logout -->
      <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200">
        <button 
          @click="handleLogout" 
          class="flex items-center w-full text-red-600 hover:bg-red-50 rounded-lg"
           :class="[
            isOpen ? 'px-4 py-2' : 'p-3 justify-center'
          ]"
          :title="isOpen ? '' : '登出'"
        >
          <svg class="h-5 w-5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M3 3a1 1 0 00-1 1v12a1 1 0 001 1h12a1 1 0 001-1V4a1 1 0 00-1-1H3zm11 4.414l-4.293 4.293a1 1 0 01-1.414 0L4 7.414V15h10V7.414z" clip-rule="evenodd" />
          </svg>
          <span :class="{ 'hidden': !isOpen, 'ml-3': isOpen }">登出</span>
        </button>
      </div>
    </div>

    <!-- Toggle Button for Mobile -->
    <button 
      @click="toggleSidebar" 
      class="fixed top-4 left-4 z-30 md:hidden bg-white p-2 rounded-lg shadow-lg hover:bg-gray-100"
      v-show="!isOpen"
    >
      <svg class="h-6 w-6 text-gray-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>

    <!-- Removed Desktop Toggle Button -->
  </div>
</template>

<style scoped>
.sidebar {
  transition: transform 0.3s ease;
}
</style>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { globalState } from '../main';
import authService from '../services/authService';

export default {
  name: 'Sidebar',
  setup() {
    const router = useRouter();

    const toggleSidebar = () => {
      // Update the global state directly
      globalState.sidebarOpen = !globalState.sidebarOpen;
      console.log('Toggling sidebar, new state:', globalState.sidebarOpen);
      
      // Still dispatch the event for backward compatibility
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('sidebar-toggle', { 
          detail: { isOpen: globalState.sidebarOpen },
          bubbles: true,
          composed: true
        }));
      }
    };

    const handleLogout = async () => {
      try {
        await authService.logout();
        router.push('/login');
      } catch (error) {
        console.error('Logout failed:', error);
      }
    };

    return {
      // Use the global state directly
      get isOpen() {
        return globalState.sidebarOpen;
      },
      toggleSidebar,
      handleLogout
    };
  }
};
</script>

<style scoped>
.sidebar {
  z-index: 1000;
}
</style>
