<template>
  <div>
    <Sidebar ref="sidebar" />
    <div 
      class="p-6 transition-all duration-300" 
      :class="contentClass"
    >
      <h1 class="text-2xl font-bold mb-6">遠端電腦管理</h1>
      
      <!-- Error Message -->
      <div v-if="error" class="mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">錯誤：</strong>
        <span class="block sm:inline">{{ error }}</span>
        <button @click="error = null" class="absolute top-0 bottom-0 right-0 px-4 py-3">
          <svg class="fill-current h-6 w-6 text-red-500" role="button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
            <title>關閉</title>
            <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"/>
          </svg>
        </button>
      </div>
      
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
                  {{ cred.username }} #{{ cred.id }}
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
                    認證資訊
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    建立時間
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    更新時間
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
                    <div v-if="computer.editing" class="flex items-center space-x-2">
                      <select 
                        v-model="computer.selectedCredentialId"
                        class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm"
                      >
                        <option value="">無認證資訊</option>
                        <option v-for="cred in credentials || []" :key="cred.id" :value="cred.id">
                          {{ cred.username }} #{{ cred.id }}
                        </option>
                      </select>
                      <div class="flex space-x-1">
                        <button 
                          @click="saveCredentialMapping(computer)"
                          class="bg-green-600 text-white p-1 rounded hover:bg-green-700 transition-colors"
                          title="儲存"
                          :disabled="loading"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                          </svg>
                        </button>
                        <button 
                          @click="cancelCredentialEdit(computer)"
                          class="bg-gray-500 text-white p-1 rounded hover:bg-gray-600 transition-colors"
                          title="取消"
                          :disabled="loading"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                      </div>
                    </div>
                    <div v-else class="flex items-center justify-between">
                      <span class="text-sm" :class="computer.credential_id ? 'text-gray-900' : 'text-gray-400 italic'">
                        {{ computer.credential_username ? 
                          `${computer.credential_username} `  : '無認證資訊' }}
                        <span v-if="computer.credential_id" class="inline-block rounded-md bg-gray-200 px-2 py-0.5 text-xs text-gray-700 font-mono ml-1">
                          ID: {{ computer.credential_id }}
                        </span>
                      </span>
                      <button 
                        @click="editCredentialMapping(computer)"
                        class="ml-2 text-blue-600 hover:text-blue-800 transition-colors"
                        title="編輯認證資訊"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                      </button>
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">{{ formatDate(computer.created_at) }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">{{ formatDate(computer.updated_at) }}</div>
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
                    <div class="text-sm font-medium text-gray-900">
                      {{ credential.username }}
                      <span class="inline-block rounded-md bg-gray-200 px-2 py-0.5 text-xs text-gray-700 font-mono ml-1">
                        ID: {{ credential.id }}
                      </span>
                    </div>
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
    
    // Fetch data - now using just two API calls
    this.fetchAllComputerMappings()
      .then(() => this.fetchCredentials())
      .catch(error => {
        console.error('初始化數據失敗:', error);
        this.error = '載入數據時發生錯誤，請重新整理頁面或聯繫系統管理員。';
      })
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
    
    async fetchAllComputerMappings() {
      try {
        const response = await fetch('http://localhost:8080/api/computers/credential-mappings', {
          method: 'GET',
          credentials: 'include'
        });
        
        if (!response.ok) {
          throw new Error('無法獲取遠端電腦與認證資訊');
        }
        
        const mappingsData = await response.json();
        
        // Check if mappingsData is null or undefined
        if (!mappingsData) {
          console.warn('API返回空數據');
          this.computers = [];
          return [];
        }
        
        // Process the data to create our computers array with credential information
        this.computers = Array.isArray(mappingsData) ? mappingsData.map(item => ({
          id: item.computer_id,
          name: item.computer_name,
          created_at: item.computer_created_at,
          updated_at: item.mapping_updated_at,
          credential_id: item.credential_id,
          credential_username: item.credential_username,
          mappingId: item.mapping_id,
          editing: false,
          selectedCredentialId: item.credential_id || ''
        })) : [];
        
        return this.computers;
      } catch (error) {
        console.error('獲取遠端電腦與認證資訊失敗:', error);
        this.error = error.message;
        this.computers = [];
        return [];
      }
    },
    
    // Keep fetchCredentials as we still need the full list for the dropdown
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
        
        // Clear the form
        this.newComputer.name = '';
        this.newComputer.credential_id = '';
        
        // Refresh the list with the combined API that gets both computers and mappings
        await this.fetchAllComputerMappings();
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
        
        await this.fetchAllComputerMappings();
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
    },
    
    editCredentialMapping(computer) {
      computer.editing = true;
      computer.selectedCredentialId = computer.credential_id || '';
    },
    
    async saveCredentialMapping(computer) {
      this.loading = true;
      try {
        // Check if we have an existing mapping by looking at the mappingId property
        if (computer.mappingId && computer.selectedCredentialId) {
          // Update an existing mapping
          const response = await fetch('http://localhost:8080/api/computers/map-credential/update', {
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
              computer_credential_mapping_id: computer.mappingId,
              credential_id: computer.selectedCredentialId
            })
          });
          
          if (!response.ok) {
            throw new Error('更新認證資訊失敗');
          }
        } else if (computer.selectedCredentialId) {
          // Create a new mapping
          const response = await fetch('http://localhost:8080/api/computers/map-credential', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
              computer_id: computer.id,
              credential_id: computer.selectedCredentialId
            })
          });
          
          if (!response.ok) {
            throw new Error('新增認證資訊映射失敗');
          }
        }
        
        // Exit edit mode
        computer.editing = false;
        
        // Save the selected credential ID before refreshing
        const selectedCredId = computer.selectedCredentialId;
        
        // Refresh the computers to get updated mappings
        await this.fetchAllComputerMappings();
        
        // Update the credential username display after refresh if needed
        // This might not be necessary if the fetchAllComputerMappings already updates everything correctly
      } catch (error) {
        console.error('更新認證資訊失敗:', error);
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
    
    cancelCredentialEdit(computer) {
      computer.editing = false;
    }
  }
};
</script> 