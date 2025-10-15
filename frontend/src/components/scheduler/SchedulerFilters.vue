<template>
  <div class="flex flex-col md:flex-row justify-between mb-6 gap-4">
    <div class="relative w-full md:w-1/3">
      <input
        :value="searchQuery"
        @input="$emit('update:searchQuery', $event.target.value)"
        type="text"
        placeholder="搜尋工作名稱..."
        class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <svg class="absolute left-3 top-3 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
      </svg>
    </div>
    
    <div class="flex gap-4">
      <select 
        :value="statusFilter"
        @change="$emit('update:statusFilter', $event.target.value)"
        class="rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">全部狀態</option>
        <option value="Ready">就緒</option>
        <option value="Running">執行中</option>
        <option value="ReadyAndRunning">就緒+執行中</option>
        <option value="Disabled">已停用</option>
      </select>

      <select 
        :value="computerFilter"
        @change="$emit('update:computerFilter', $event.target.value)"
        class="rounded-lg border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">全部電腦</option>
        <option v-for="computer in computers" :key="computer" :value="computer">
          {{ computer }}
        </option>
      </select>
      
      <button @click="$emit('refresh')" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg transition-colors flex items-center gap-2">
        <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
        </svg>
        重新整理
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SchedulerFilters',
  props: {
    searchQuery: {
      type: String,
      required: true
    },
    statusFilter: {
      type: String,
      required: true
    },
    computerFilter: {
      type: String,
      required: true
    },
    computers: {
      type: Array,
      required: true
    }
  },
  emits: ['update:searchQuery', 'update:statusFilter', 'update:computerFilter', 'refresh']
}
</script>