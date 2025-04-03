<template>
  <div>
    <Sidebar ref="sidebar" />
    <div 
      class="p-6 transition-all duration-300" 
      :class="contentClass"
    >
      <h1 class="text-2xl font-bold mb-6">遠端電腦管理</h1>
      
      <!-- Loading indicator -->
      <div v-if="loading" class="flex justify-center items-center py-10">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
      
      <div v-else>
        <!-- Add New Computer Form -->
        <div class="mb-8 bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-xl font-semibold mb-4">新增遠端電腦</h2>
          <form @submit.prevent="addComputer">
            <div class="mb-4">
              <label for="computerName" class="block text-sm font-medium text-gray-700">電腦名稱</label>
              <input
                type="text"
                id="computerName"
                v-model="newComputer.name"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="輸入電腦名稱或IP"
                required
              />
            </div>
            
            <div class="mb-4">
              <label for="credential" class="block text-sm font-medium text-gray-700">選擇認證資訊（選填）</label>
              <select
                id="credential"
                v-model="newComputer.credential_id"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              >
                <option value="">選擇認證資訊</option>
                <option v-for="cred in credentials || []" :key="cred.id" :value="cred.id">
                  {{ cred.username }}
                </option>
              </select>
            </div>
            
            <div class="flex justify-end">
              <button
                type="submit"
                class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                新增電腦
              </button>
            </div>
          </form>
        </div>
        
        <!-- Add New Credential Form -->
        <div class="mb-8 bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-xl font-semibold mb-4">新增認證資訊</h2>
          <form @submit.prevent="addCredential">
            <div class="mb-4">
              <label for="username" class="block text-sm font-medium text-gray-700">使用者名稱</label>
              <input
                type="text"
                id="username"
                v-model="newCredential.username"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="輸入使用者名稱"
                required
              />
            </div>
            
            <div class="mb-4">
              <label for="password" class="block text-sm font-medium text-gray-700">密碼</label>
              <input
                type="password"
                id="password"
                v-model="newCredential.password"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="輸入密碼"
                required
              />
            </div>
            
            <div class="flex justify-end">
              <button
                type="submit"
                class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                新增認證資訊
              </button>
            </div>
          </form>
        </div>
        
        <!-- Remote Computers List -->
        <div class="bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-xl font-semibold mb-4">遠端電腦列表</h2>
          <div v-if="!computers || computers.length === 0" class="text-center py-4 text-gray-500">
            尚未新增任何遠端電腦
          </div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    名稱
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    建立時間
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="computer in computers || []" :key="computer.id">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{{ computer.name }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">{{ formatDate(computer.created_at) }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button
                      @click="deleteComputer(computer.id)"
                      class="text-red-600 hover:text-red-900"
                    >
                      刪除
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        
        <!-- Credentials List -->
        <div class="mt-8 bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-xl font-semibold mb-4">認證資訊列表</h2>
          <div v-if="!credentials || credentials.length === 0" class="text-center py-4 text-gray-500">
            尚未新增任何認證資訊
          </div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    使用者名稱
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    建立時間
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="credential in credentials || []" :key="credential.id">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{{ credential.username }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">{{ formatDate(credential.created_at) }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button
                      @click="deleteCredential(credential.id)"
                      class="text-red-600 hover:text-red-900"
                    >
                      刪除
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Sidebar from './Sidebar.vue';
import { globalState } from '../main';

export default {
  name: 'RemoteComputers',
  components: {
    Sidebar
  },
  data() {
    return {
      computers: [],
      credentials: [],
      newComputer: {
        name: '',
        credential_id: ''
      },
      newCredential: {
        username: '',
        password: ''
      },
      loading: false,
      error: null,
      isMobile: false
    };
  },
  computed: {
    sidebarOpen() {
      // Access sidebar state from global state
      return globalState.sidebarOpen;
    },
    contentClass() {
      // Determine content class based on sidebar state and screen size
      return {
        'md:ml-64': this.sidebarOpen && !this.isMobile,
        'ml-0': !this.sidebarOpen || this.isMobile
      };
    }
  },
  mounted() {
    // Set loading state before fetching data
    this.loading = true;
    
    // Fetch data sequentially
    this.fetchComputers()
      .then(() => this.fetchCredentials())
      .finally(() => {
        this.loading = false;
        this.checkMobileView();
      });
    
    // Listen for resize events to update isMobile
    window.addEventListener('resize', this.checkMobileView);
  },
  beforeUnmount() {
    // Only need to remove resize listener now
    window.removeEventListener('resize', this.checkMobileView);
  },
  methods: {
    checkMobileView() {
      this.isMobile = window.innerWidth < 768;
    },
    
    async fetchComputers() {
      try {
        const response = await fetch('http://localhost:8080/api/computers/list', {
          method: 'GET',
          credentials: 'include'
        });
        
        if (!response.ok) {
          throw new Error('無法獲取遠端電腦列表');
        }
        
        this.computers = await response.json();
        return this.computers;
      } catch (error) {
        console.error('獲取遠端電腦失敗:', error);
        this.error = error.message;
        this.computers = [];
        return [];
      }
    },
    
    async fetchCredentials() {
      try {
        const response = await fetch('http://localhost:8080/api/credentials/list', {
          method: 'GET',
          credentials: 'include'
        });
        
        if (!response.ok) {
          throw new Error('無法獲取認證資訊列表');
        }
        
        this.credentials = await response.json();
        return this.credentials;
      } catch (error) {
        console.error('獲取認證資訊失敗:', error);
        this.error = error.message;
        this.credentials = [];
        return [];
      }
    },
    
    async addComputer() {
      this.loading = true;
      try {
        const payload = {
          name: this.newComputer.name
        };
        
        if (this.newComputer.credential_id) {
          payload.credential_id = parseInt(this.newComputer.credential_id);
        }
        
        const response = await fetch('http://localhost:8080/api/computers', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify(payload)
        });
        
        if (!response.ok) {
          throw new Error('新增遠端電腦失敗');
        }
        
        // Clear the form and refresh the list
        this.newComputer.name = '';
        this.newComputer.credential_id = '';
        await this.fetchComputers();
      } catch (error) {
        console.error('新增遠端電腦失敗:', error);
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
    
    async addCredential() {
      this.loading = true;
      try {
        const response = await fetch('http://localhost:8080/api/credentials', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify(this.newCredential)
        });
        
        if (!response.ok) {
          throw new Error('新增認證資訊失敗');
        }
        
        // Clear the form and refresh the list
        this.newCredential.username = '';
        this.newCredential.password = '';
        await this.fetchCredentials();
      } catch (error) {
        console.error('新增認證資訊失敗:', error);
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteComputer(computerId) {
      if (!confirm('確定要刪除這台遠端電腦嗎?')) {
        return;
      }
      
      this.loading = true;
      try {
        const response = await fetch('http://localhost:8080/api/computers/delete', {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify({ computer_id: computerId })
        });
        
        if (!response.ok) {
          throw new Error('刪除遠端電腦失敗');
        }
        
        await this.fetchComputers();
      } catch (error) {
        console.error('刪除遠端電腦失敗:', error);
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteCredential(credentialId) {
      if (!confirm('確定要刪除這個認證資訊嗎?')) {
        return;
      }
      
      this.loading = true;
      try {
        const response = await fetch('http://localhost:8080/api/credentials/delete', {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify({ credential_id: credentialId })
        });
        
        if (!response.ok) {
          throw new Error('刪除認證資訊失敗');
        }
        
        await this.fetchCredentials();
      } catch (error) {
        console.error('刪除認證資訊失敗:', error);
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
    
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleString();
    }
  }
};
</script> 