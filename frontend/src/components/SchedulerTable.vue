<template>
  <div class="container mx-auto px-4 py-8">
    <header class="mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Windows 排程器工作資訊</h1>
      <p class="text-gray-600 mt-2">顯示所有排程工作的狀態和詳細資訊</p>
    </header>
    
    <!-- 加載狀態 -->
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
    </div>
    
    <!-- 錯誤提示 -->
    <div v-else-if="error" class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded mb-6">
      <p class="font-bold">載入失敗</p>
      <p>{{ error }}</p>
    </div>
 
    <!-- 主要內容 -->
    <div v-else>
      <!-- 搜尋和過濾功能 -->
      <div class="flex flex-col md:flex-row justify-between mb-6 gap-4">
        <div class="relative w-full md:w-1/3">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜尋工作名稱..."
            class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <svg class="absolute left-3 top-3 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
          </svg>
        </div>
        
        <div class="flex gap-4">
          <select v-model="statusFilter" class="rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="">全部狀態</option>
            <option value="Ready">就緒</option>
            <option value="Running">執行中</option>
            <option value="Disabled">已停用</option>
          </select>
          
          <button @click="refreshJobs" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg transition-colors flex items-center gap-2">
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
            </svg>
            重新整理
          </button>
        </div>
      </div>
 
      <!-- 工作卡片視圖 -->
      <div v-if="filteredJobs.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="job in filteredJobs" :key="job.RegistrationInfo.URI" class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
          <div class="p-4 border-b border-gray-100 flex justify-between items-center">
            <h2 class="font-bold text-lg text-gray-800 truncate" :title="job.RegistrationInfo.URI">
              {{ job.RegistrationInfo.URI }}
            </h2>
            <span :class="getStatusClass(job.ExtraInfo.State)" class="px-2 py-1 rounded-full text-xs font-bold">
              {{ job.ExtraInfo.State }}
            </span>
          </div>
          
          <div class="p-4 space-y-3">
            <div class="flex items-start">
              <span class="text-gray-500 w-24 flex-shrink-0">執行指令:</span>
              <span class="text-gray-800 break-all">{{ job.Actions.Exec.Command || '無' }}</span>
            </div>
            
            <div class="flex items-start">
              <span class="text-gray-500 w-24 flex-shrink-0">觸發程序:</span>
              <span class="text-gray-800">{{ getTriggers(job) }}</span>
            </div>
            
            <div class="flex items-start">
              <span class="text-gray-500 w-24 flex-shrink-0">電腦名稱:</span>
              <span class="text-gray-800">{{ job.ExtraInfo.ComputerName }}</span>
            </div>
            
            <div class="flex items-start">
              <span class="text-gray-500 w-24 flex-shrink-0">作者:</span>
              <span class="text-gray-800">{{ job.RegistrationInfo.Author }}</span>
            </div>
          </div>
          
          <div class="bg-gray-50 p-4 space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">下次執行:</span>
              <span class="text-gray-800 font-medium">{{ formatDateTime(job.ExtraInfo.NextRunTime) }}</span>
            </div>
            
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">上次執行:</span>
              <span class="text-gray-800 font-medium">{{ formatDateTime(job.ExtraInfo.LastRunTime) }}</span>
            </div>
            
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">執行結果:</span>
              <span :class="getResultClass(job.ExtraInfo.LastTaskResult)">{{ job.ExtraInfo.LastTaskResult || '無' }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 無結果顯示 -->
      <div v-else class="bg-gray-50 rounded-lg p-8 text-center text-gray-500">
        <svg class="h-16 w-16 mx-auto text-gray-400 mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-xl font-medium mb-2">沒有找到符合條件的工作</h3>
        <p>嘗試調整搜尋條件或重新整理資料</p>
      </div>
    </div>
  </div>
 </template>
 
 <script>
 import axios from 'axios';
 import { ref, computed, onMounted } from 'vue';
 
 export default {
  setup() {
    const jobs = ref([]);
    const loading = ref(true);
    const error = ref(null);
    const searchQuery = ref('');
    const statusFilter = ref('Ready');
 
    const fetchJobs = async () => {
      loading.value = true;
      error.value = null;
      
      try {
        const response = await axios.get('http://localhost:8080/api/tasks');
        jobs.value = response.data;
      } catch (err) {
        console.error(err);
        error.value = '無法載入排程工作資訊，請稍後再試。';
      } finally {
        loading.value = false;
      }
    };
 
    const refreshJobs = () => {
      fetchJobs();
    };
 
    const filteredJobs = computed(() => {
      return jobs.value.filter(job => {
        const matchesSearch = searchQuery.value === '' ||
          job.RegistrationInfo.URI.toLowerCase().includes(searchQuery.value.toLowerCase());
        
        const matchesStatus = statusFilter.value === '' ||
          job.ExtraInfo.State === statusFilter.value;
        
        return matchesSearch && matchesStatus;
      });
    });
 
    const getTriggers = (job) => {
      const triggers = job.Triggers;
      let triggerText = '';
 
      if (triggers && triggers.length > 0) {
        if (triggers[0].TimeTriggers) {
          triggerText += triggers[0].TimeTriggers.map(trigger => `定時: ${formatDateTime(trigger.StartBoundary)}`).join(', ');
        }
        if (triggers[0].CalendarTriggers) {
          triggerText += triggers[0].CalendarTriggers.map(trigger => `日曆: ${formatDateTime(trigger.StartBoundary)}`).join(', ');
        }
      }
 
      return triggerText || '無';
    };
 
    const formatDateTime = (dateTime) => {
      if (!dateTime) return '無';
      
      try {
        const date = new Date(dateTime);
        return new Intl.DateTimeFormat('zh-TW', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          hour12: false
        }).format(date);
      } catch (e) {
        return dateTime;
      }
    };
 
    const getStatusClass = (status) => {
      switch (status) {
        case '已啟用':
          return 'bg-green-100 text-green-800';
        case '已停用':
          return 'bg-gray-100 text-gray-800';
        default:
          return 'bg-blue-100 text-blue-800';
      }
    };
 
    const getResultClass = (result) => {
      if (!result || result === '0') {
        return 'text-green-600 font-medium';
      }
      return 'text-red-600 font-medium';
    };
 
    onMounted(() => {
      fetchJobs();
    });
 
    return {
      jobs,
      loading,
      error,
      searchQuery,
      statusFilter,
      filteredJobs,
      getTriggers,
      formatDateTime,
      getStatusClass,
      getResultClass,
      refreshJobs
    };
  }
 };
 </script>