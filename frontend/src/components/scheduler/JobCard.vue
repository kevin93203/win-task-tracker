<template>
  <div class="bg-white rounded-xl shadow-sm hover:shadow-md transition-all duration-300 border border-gray-200">
    <div class="p-5">
      <div class="flex justify-between items-start gap-4 mb-3">
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <span 
              class="w-2 h-2 rounded-full" 
              :class="{
                'bg-green-500': job.ExtraInfo.State === 'Ready',
                'bg-yellow-500': job.ExtraInfo.State === 'Running',
                'bg-red-500': job.ExtraInfo.State === 'Disabled'
              }"
            ></span>
            <h2 class="font-semibold text-lg text-gray-800 truncate group-hover:text-blue-600 transition-colors" :title="job.ExtraInfo.TaskName">
              {{ job.ExtraInfo.TaskName }}
            </h2>
          </div>
          <p class="text-sm text-gray-500 mt-1">{{ job.ExtraInfo.ComputerName }}</p>
        </div>
        
        <div class="flex items-center gap-2">
          <button 
            v-if="job.ExtraInfo.State === 'Disabled'"
            @click="enableTask" 
            class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-green-600 hover:text-green-700 hover:bg-green-50 rounded transition-colors group relative cursor-pointer"
            :disabled="isEnabling"
            :title="isEnabling ? '啟用中...' : '啟用'"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
          <button 
            v-if="job.ExtraInfo.State === 'Ready'"
            @click="startTask" 
            class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-green-600 hover:text-green-700 hover:bg-green-50 rounded transition-colors group relative cursor-pointer"
            :disabled="isStarting"
            :title="isStarting ? '啟動中...' : '啟動'"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M5 3v18l14-9z" />
            </svg>
          </button>
          <button 
            v-if="job.ExtraInfo.State === 'Running'"
            @click="stopTask" 
            class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-red-500 hover:text-red-600 hover:bg-red-50 rounded transition-colors group relative cursor-pointer"
            :disabled="isStopping"
            :title="isStopping ? '停止中...' : '停止'"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <rect x="5" y="5" width="14" height="14" rx="1" />
            </svg>
          </button>
          <button 
            v-if="job.ExtraInfo.State !== 'Disabled'"
            @click="disableTask" 
            class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-red-500 hover:text-red-600 hover:bg-red-50 rounded transition-colors group relative cursor-pointer"
            :disabled="isDisabling"
            :title="isDisabling ? '停用中...' : '停用'"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
          </button>
        </div>
      </div>
      <div v-if="job.RegistrationInfo.Description" class="mt-2 text-sm text-gray-500 line-clamp-2">
        {{ job.RegistrationInfo.Description }}
      </div>

      <div class="mt-6 space-y-4">
        <TaskActionDisplay
          :computer-id="job.ExtraInfo.ComputerID"
          :task-name="job.ExtraInfo.TaskName"
          :actions="job.Actions || []"
          @refresh="$emit('refresh', {
            computerID: job.ExtraInfo.ComputerID,
            taskName: job.ExtraInfo.TaskName
          })"
        />
        <TaskTriggerDisplay
          :computer-id="job.ExtraInfo.ComputerID"
          :task-name="job.ExtraInfo.TaskName"
          :triggers="job.Triggers || []"
          @refresh="$emit('refresh', {
            computerID: job.ExtraInfo.ComputerID,
            taskName: job.ExtraInfo.TaskName
          })"
        />
        <div class="flex items-start gap-x-3">
          <div class="flex-shrink-0">
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-900">作者</p>
            <p class="mt-1 text-sm text-gray-500">{{ job.RegistrationInfo.Author || '未知' }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="border-t border-gray-200 bg-gray-50 p-5">
      <div class="grid grid-cols-3 gap-4">
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500">下次執行</p>
          <p class="text-sm font-medium text-gray-900 break-words">{{ job.ExtraInfo.NextRunTime || '無' }}</p>
        </div>
        
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500">上次執行</p>
          <p class="text-sm font-medium text-gray-900 break-words">{{ job.ExtraInfo.LastRunTime || '無' }}</p>
        </div>
        
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500">執行結果</p>
          <p :class="['text-sm font-medium break-words', getResultClass(job.ExtraInfo.LastTaskResult)]">{{ job.ExtraInfo.LastTaskResult || '無' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 導入必要組件和庫
import TaskTriggerDisplay from './TaskTriggerDisplay.vue';
import TaskActionDisplay from './TaskActionDisplay.vue';
import { ref } from 'vue';
import { useToast } from 'vue-toastification';
import axios from 'axios';

// 定義 props
const props = defineProps({
  job: {
    type: Object,
    required: true
  }
});

// 定義 emits
const emit = defineEmits(['task-disabled', 'task-enabled', 'task-started', 'task-stopped', 'refresh']);

// 創建refs
const isDisabling = ref(false);
const isEnabling = ref(false);
const isStarting = ref(false);
const isStopping = ref(false);
const toast = useToast();

// 啟用任務
async function enableTask() {
  if (isEnabling.value) return;
  console.log('Enabling task:', props.job.ExtraInfo.TaskName);

  try {
    isEnabling.value = true;

    // 取得 JWT token
    const cookies = document.cookie.split(';');
    const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='));
    const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null;

    const response = await axios.post(
      'http://localhost:8080/api/tasks/enable',
      {
        computer_id: props.job.ExtraInfo.ComputerID,
        task_name: props.job.ExtraInfo.TaskName
      },
      {
        withCredentials: true,
        headers: {
          'Authorization': `Bearer ${jwt}`
        }
      }
    );

    const result = response.data;

    if (result.success) {
      toast.success('任務已啟用');
      emit('task-enabled', {
        computerID: props.job.ExtraInfo.ComputerID,
        taskName: props.job.ExtraInfo.TaskName
      });
    } else {
      throw new Error(result.error || '啟用任務失敗');
    }
  } catch (error) {
    if (error.response) {
      // 後端回傳的錯誤
      toast.error(error.response.data.error || '啟用任務失敗');
    } else if (error.request) {
      // 網路連線問題
      toast.error('網路連線失敗，請稍後再試');
    } else {
      // 其他錯誤
      toast.error(error.message || '啟用任務失敗');
    }
  } finally {
    isEnabling.value = false;
  }
}

// 啟動任務
async function startTask() {
  if (isStarting.value) return;
  console.log('Starting task:', props.job.ExtraInfo.TaskName);

  try {
    isStarting.value = true;

    // 取得 JWT token
    const cookies = document.cookie.split(';');
    const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='));
    const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null;

    const response = await axios.post(
      'http://localhost:8080/api/tasks/start',
      {
        computer_id: props.job.ExtraInfo.ComputerID,
        task_name: props.job.ExtraInfo.TaskName
      },
      {
        withCredentials: true,
        headers: {
          'Authorization': `Bearer ${jwt}`
        }
      }
    );

    const result = response.data;

    if (result.success) {
      toast.success('任務已啟動');
      emit('task-started', {
        computerID: props.job.ExtraInfo.ComputerID,
        taskName: props.job.ExtraInfo.TaskName
      });
    } else {
      throw new Error(result.error || '啟動任務失敗');
    }
  } catch (error) {
    if (error.response) {
      // 後端回傳的錯誤
      toast.error(error.response.data.error || '啟動任務失敗');
    } else if (error.request) {
      // 網路連線問題
      toast.error('網路連線失敗，請稍後再試');
    } else {
      // 其他錯誤
      toast.error(error.message || '啟動任務失敗');
    }
  } finally {
    isStarting.value = false;
  }
}

// 停止任務
async function stopTask() {
  if (isStopping.value) return;
  console.log('Stopping task:', props.job.ExtraInfo.TaskName);

  try {
    isStopping.value = true;

    // 取得 JWT token
    const cookies = document.cookie.split(';');
    const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='));
    const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null;

    const response = await axios.post(
      'http://localhost:8080/api/tasks/stop',
      {
        computer_id: props.job.ExtraInfo.ComputerID,
        task_name: props.job.ExtraInfo.TaskName
      },
      {
        withCredentials: true,
        headers: {
          'Authorization': `Bearer ${jwt}`
        }
      }
    );

    const result = response.data;

    if (result.success) {
      toast.success('任務已停止');
      emit('task-stopped', {
        computerID: props.job.ExtraInfo.ComputerID,
        taskName: props.job.ExtraInfo.TaskName
      });
    } else {
      throw new Error(result.error || '停止任務失敗');
    }
  } catch (error) {
    if (error.response) {
      // 後端回傳的錯誤
      toast.error(error.response.data.error || '停止任務失敗');
    } else if (error.request) {
      // 網路連線問題
      toast.error('網路連線失敗，請稍後再試');
    } else {
      // 其他錯誤
      toast.error(error.message || '停止任務失敗');
    }
  } finally {
    isStopping.value = false;
  }
}

// 停用任務
async function disableTask() {
  if (isDisabling.value) return;
  console.log('Disabling task:', props.job.ExtraInfo.TaskName);

  try {
    isDisabling.value = true;

    // 取得 JWT token
    const cookies = document.cookie.split(';');
    const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='));
    const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null;

    const response = await axios.post(
      'http://localhost:8080/api/tasks/disable',
      {
        computer_id: props.job.ExtraInfo.ComputerID,
        task_name: props.job.ExtraInfo.TaskName
      },
      {
        withCredentials: true,
        headers: {
          'Authorization': `Bearer ${jwt}`
        }
      }
    );

    const result = response.data;

    if (result.success) {
      toast.success('任務已停用');
      emit('task-disabled', {
        computerID: props.job.ExtraInfo.ComputerID,
        taskName: props.job.ExtraInfo.TaskName
      });
    } else {
      throw new Error(result.error || '停用任務失敗');
    }
  } catch (error) {
    if (error.response) {
      // 後端回傳的錯誤
      toast.error(error.response.data.error || '停用任務失敗');
    } else if (error.request) {
      // 網路連線問題
      toast.error('網路連線失敗，請稍後再試');
    } else {
      // 其他錯誤
      toast.error(error.message || '停用任務失敗');
    }
  } finally {
    isDisabling.value = false;
  }
}

// 獲取結果樣式
function getResultClass(result) {
  return result === 0 ? 'text-green-700' : 'text-red-700';
}

// 定義狀態相關的方法
function getStatusChinese(status) {
  switch (status) {
    case 'Ready':
      return '就緒';
    case 'Running':
      return '執行中';
    case 'Disabled':
      return '已停用';
    default:
      return status;
  }
}

function getStatusClass(status) {
  switch (status) {
    case 'Ready':
      return 'bg-green-50 text-green-700';
    case 'Running':
      return 'bg-yellow-50 text-yellow-700';
    case 'Disabled':
      return 'bg-red-50 text-red-700';
    default:
      return 'bg-gray-50 text-gray-700';
  }
}
</script>
