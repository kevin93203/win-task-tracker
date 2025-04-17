<template>
  <div class="bg-white rounded-xl shadow-sm hover:shadow-md transition-all duration-300 border border-gray-100">
    <div class="p-5">
      <div class="flex justify-between items-start gap-4 mb-3">
        <div class="flex-1 min-w-0">
          <h2 class="font-semibold text-lg text-gray-800 truncate group-hover:text-indigo-600 transition-colors" :title="job.ExtraInfo.TaskName">
            {{ job.ExtraInfo.TaskName }}
          </h2>
          <p class="text-sm text-gray-500 mt-1">{{ job.ExtraInfo.ComputerName }}</p>
        </div>
        
        <div class="flex items-center flex-wrap gap-3">
          <button 
            v-if="job.ExtraInfo.State === 'Disabled'"
            @click="enableTask" 
            class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-green-700 hover:text-white bg-green-50 hover:bg-green-600 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed group"
            :disabled="isEnabling"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span>{{ isEnabling ? '啟用中...' : '啟用' }}</span>
          </button>
          <button 
            v-if="job.ExtraInfo.State === 'Ready'"
            @click="startTask" 
            class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-blue-700 hover:text-white bg-blue-50 hover:bg-blue-600 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed group"
            :disabled="isStarting"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ isStarting ? '啟動中...' : '啟動' }}</span>
          </button>
          <button 
            v-if="job.ExtraInfo.State === 'Running'"
            @click="stopTask" 
            class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-orange-700 hover:text-white bg-orange-50 hover:bg-orange-600 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed group"
            :disabled="isStopping"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
            </svg>
            <span>{{ isStopping ? '停止中...' : '停止' }}</span>
          </button>
          <button 
            v-if="job.ExtraInfo.State !== 'Disabled'"
            @click="disableTask" 
            class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-red-700 hover:text-white bg-red-50 hover:bg-red-600 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed group"
            :disabled="isDisabling"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
            <span>{{ isDisabling ? '停用中...' : '停用' }}</span>
          </button>
          <span 
            :class="[getStatusClass(job.ExtraInfo.State), 
                   'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-sm font-medium']">
            <span class="w-2 h-2 rounded-full" :class="{
              'bg-green-400': job.ExtraInfo.State === 'Ready',
              'bg-yellow-400': job.ExtraInfo.State === 'Running',
              'bg-red-400': job.ExtraInfo.State === 'Disabled'
            }"></span>
            {{ getStatusChinese(job.ExtraInfo.State) }}
          </span>
        </div>
      </div>
      <div v-if="job.RegistrationInfo.Description" class="mt-2 text-sm text-gray-500 line-clamp-2">
        {{ job.RegistrationInfo.Description }}
      </div>

      <div class="mt-6 space-y-4">
        <!-- 新的 Action Display 元件 -->
        <TaskActionDisplay
          :computer-id="job.ExtraInfo.ComputerID"
          :task-name="job.ExtraInfo.TaskName"
          :actions="job.Actions || []"
          @refresh="$emit('refresh')"
        />
        <!-- 新的 Trigger Display 元件 -->
        <TaskTriggerDisplay
          :computer-id="job.ExtraInfo.ComputerID"
          :task-name="job.ExtraInfo.TaskName"
          :triggers="job.Triggers || []"
          @refresh="$emit('refresh')"
        />
        <!-- 作者區塊 -->
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
    
    <div class="border-t border-gray-100 bg-gray-50 p-5">
      <div class="grid grid-cols-3 gap-4">
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500">下次執行</p>
          <p class="text-sm font-medium text-gray-900 break-words">{{ formatDateTime(job.ExtraInfo.NextRunTime) || '無' }}</p>
        </div>
        
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500">上次執行</p>
          <p class="text-sm font-medium text-gray-900 break-words">{{ formatDateTime(job.ExtraInfo.LastRunTime) || '無' }}</p>
        </div>
        
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500">執行結果</p>
          <p :class="['text-sm font-medium break-words', getResultClass(job.ExtraInfo.LastTaskResult)]">{{ job.ExtraInfo.LastTaskResult || '無' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useToast } from 'vue-toastification'
import axios from 'axios'

export default {
  name: 'JobCard',
  emits: ['task-disabled', 'task-enabled', 'task-started', 'task-stopped'],
  props: {
    job: {
      type: Object,
      required: true
    }
  },
  setup(props, { emit }) {
    const toast = useToast()
    const isDisabling = ref(false)
    const isEnabling = ref(false)
    const isStarting = ref(false)
    const isStopping = ref(false)

    const disableTask = async () => {
      if (isDisabling.value) return

      try {
        isDisabling.value = true

        // 取得 JWT token
        const cookies = document.cookie.split(';')
        const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='))
        const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null

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
        )

        const result = response.data

        if (result.success) {
          toast.success('任務已停用')
          emit('task-disabled', {
            computerID: props.job.ExtraInfo.ComputerID,
            taskName: props.job.ExtraInfo.TaskName
          })
        } else {
          throw new Error(result.error || '停用任務失敗')
        }
      } catch (error) {
        if (error.response) {
          // 後端回傳的錯誤
          toast.error(error.response.data.error || '停用任務失敗')
        } else if (error.request) {
          // 網路連線問題
          toast.error('網路連線失敗，請稍後再試')
        } else {
          // 其他錯誤
          toast.error(error.message || '停用任務失敗')
        }
      } finally {
        isDisabling.value = false
      }
    }

    const enableTask = async () => {
      if (isEnabling.value) return

      try {
        isEnabling.value = true

        // 取得 JWT token
        const cookies = document.cookie.split(';')
        const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='))
        const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null

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
        )

        const result = response.data

        if (result.success) {
          toast.success('任務已啟用')
          emit('task-enabled', {
            computerID: props.job.ExtraInfo.ComputerID,
            taskName: props.job.ExtraInfo.TaskName
          })
        } else {
          throw new Error(result.error || '啟用任務失敗')
        }
      } catch (error) {
        if (error.response) {
          // 後端回傳的錯誤
          toast.error(error.response.data.error || '啟用任務失敗')
        } else if (error.request) {
          // 網路連線問題
          toast.error('網路連線失敗，請稍後再試')
        } else {
          // 其他錯誤
          toast.error(error.message || '啟用任務失敗')
        }
      } finally {
        isEnabling.value = false
      }
    }

    const startTask = async () => {
      if (isStarting.value) return

      try {
        isStarting.value = true

        // 取得 JWT token
        const cookies = document.cookie.split(';')
        const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='))
        const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null

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
        )

        const result = response.data

        if (result.success) {
          toast.success('任務已啟動')
          emit('task-started', {
            computerID: props.job.ExtraInfo.ComputerID,
            taskName: props.job.ExtraInfo.TaskName
          })
        } else {
          throw new Error(result.error || '啟動任務失敗')
        }
      } catch (error) {
        if (error.response) {
          // 後端回傳的錯誤
          toast.error(error.response.data.error || '啟動任務失敗')
        } else if (error.request) {
          // 網路連線問題
          toast.error('網路連線失敗，請稍後再試')
        } else {
          // 其他錯誤
          toast.error(error.message || '啟動任務失敗')
        }
      } finally {
        isStarting.value = false
      }
    }

    const stopTask = async () => {
      if (isStopping.value) return

      try {
        isStopping.value = true

        // 取得 JWT token
        const cookies = document.cookie.split(';')
        const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='))
        const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null

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
        )

        const result = response.data

        if (result.success) {
          toast.success('任務已停止')
          emit('task-stopped', {
            computerID: props.job.ExtraInfo.ComputerID,
            taskName: props.job.ExtraInfo.TaskName
          })
        } else {
          throw new Error(result.error || '停止任務失敗')
        }
      } catch (error) {
        if (error.response) {
          // 後端回傳的錯誤
          toast.error(error.response.data.error || '停止任務失敗')
        } else if (error.request) {
          // 網路連線問題
          toast.error('網路連線失敗，請稍後再試')
        } else {
          // 其他錯誤
          toast.error(error.message || '停止任務失敗')
        }
      } finally {
        isStopping.value = false
      }
    }

    return {
      isDisabling,
      isEnabling,
      isStarting,
      isStopping,
      disableTask,
      enableTask,
      startTask,
      stopTask
    }
  },

  methods: {
    formatDateTime(dateTime) {
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
    },
    getStatusChinese(status) {
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
    },
    getStatusClass(status) {
      switch (status) {
        case 'Ready':
          return 'bg-green-50 text-green-700'
        case 'Running':
          return 'bg-yellow-50 text-yellow-700'
        case 'Disabled':
          return 'bg-red-50 text-red-700'
        default:
          return 'bg-gray-50 text-gray-700'
      }
    },
    getResultClass(result) {
      return result === 0 ? 'text-green-700' : 'text-red-700'
    }
  }
}
</script>

<script setup>
// 移除 TaskTriggers 和 TaskActions 的導入
// import TaskTriggers from './TaskTriggers.vue';
// import TaskActions from './TaskActions.vue';
// 導入新的 Display 元件
import TaskTriggerDisplay from './TaskTriggerDisplay.vue';
import TaskActionDisplay from './TaskActionDisplay.vue';
</script>
