<template>
  <div class="flex">
    <Sidebar ref="sidebar" />
    <div class="flex-1 p-8 transition-all duration-300" :class="{ 'ml-64': isSidebarOpen, 'ml-0': !isSidebarOpen }">
      <SchedulerHeader />
      
      <LoadingSpinner v-if="loading" />
      <ErrorMessage v-else-if="error" :message="error" />
   
      <div v-else>
        <TaskErrors :errors="taskErrors" />
        <SchedulerFilters
          v-model:searchQuery="searchQuery"
          v-model:statusFilter="statusFilter"
          v-model:computerFilter="computerFilter"
          :computers="uniqueComputers"
          @refresh="refreshJobs"
        />
 
        <div v-if="filteredJobs.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <JobCard
            v-for="job in filteredJobs"
            :key="job.ExtraInfo.ComputerName + job.RegistrationInfo.URI"
            :job="job"
            @task-disabled="handleTaskDisabled"
            @task-enabled="handleTaskEnabled"
            @task-started="handleTaskStarted"
            @task-stopped="handleTaskStopped"
          />
        </div>
        
        <NoResults v-else />
      </div>
    </div>
  </div>
</template>
 
<script>
import axios from 'axios';
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Sidebar from './Sidebar.vue';
import SchedulerHeader from './scheduler/SchedulerHeader.vue';
import SchedulerFilters from './scheduler/SchedulerFilters.vue';
import JobCard from './scheduler/JobCard.vue';
import LoadingSpinner from './scheduler/LoadingSpinner.vue';
import ErrorMessage from './scheduler/ErrorMessage.vue';
import NoResults from './scheduler/NoResults.vue';
import TaskErrors from './scheduler/TaskErrors.vue';
import { globalState } from '../main';

export default {
  components: {
    Sidebar,
    SchedulerHeader,
    SchedulerFilters,
    JobCard,
    LoadingSpinner,
    ErrorMessage,
    NoResults,
    TaskErrors
  },
  setup() {
    const router = useRouter();
    const jobs = ref([]);
    const taskErrors = ref([]);
    const loading = ref(true);
    const error = ref(null);
    const searchQuery = ref('');
    const statusFilter = ref('Ready');
    const computerFilter = ref('');
    
    // 確保jobs.value永遠是陣列
    const ensureJobsArray = computed(() => jobs.value || []);
    const isSidebarOpen = computed(() => globalState.sidebarOpen);

    const handleTaskDisabled = ({ computerID, taskName }) => {
      // 更新本地狀態
      jobs.value = jobs.value.map(job => {
        if (job.ExtraInfo.ComputerID === computerID && job.ExtraInfo.TaskName === taskName) {
          return {
            ...job,
            ExtraInfo: {
              ...job.ExtraInfo,
              State: 'Disabled'
            }
          }
        }
        return job
      })
    }

    const handleTaskEnabled = ({ computerID, taskName }) => {
      // 更新本地狀態
      jobs.value = jobs.value.map(job => {
        if (job.ExtraInfo.ComputerID === computerID && job.ExtraInfo.TaskName === taskName) {
          return {
            ...job,
            ExtraInfo: {
              ...job.ExtraInfo,
              State: 'Ready'
            }
          }
        }
        return job
      })
    }

    const handleTaskStarted = ({ computerID, taskName }) => {
      // 更新本地狀態
      jobs.value = jobs.value.map(job => {
        if (job.ExtraInfo.ComputerID === computerID && job.ExtraInfo.TaskName === taskName) {
          return {
            ...job,
            ExtraInfo: {
              ...job.ExtraInfo,
              State: 'Running'
            }
          }
        }
        return job
      })
    }

    const handleTaskStopped = ({ computerID, taskName }) => {
      // 更新本地狀態
      jobs.value = jobs.value.map(job => {
        if (job.ExtraInfo.ComputerID === computerID && job.ExtraInfo.TaskName === taskName) {
          return {
            ...job,
            ExtraInfo: {
              ...job.ExtraInfo,
              State: 'Ready'
            }
          }
        }
        return job
      })
    }

    const fetchJobs = async () => {
      loading.value = true;
      error.value = null;
      
      try {
        const cookies = document.cookie.split(';');
        const jwtCookie = cookies.find(cookie => cookie.trim().startsWith('jwt='));
        const jwt = jwtCookie ? jwtCookie.split('=')[1].trim() : null;

        const response = await axios.get('http://localhost:8080/api/tasks', {
          withCredentials: true,
          headers: {
            'Authorization': `Bearer ${jwt}`
          }
        });
        
        // 處理新的回應格式
        const { tasks = [], errors = [] } = response.data || {};
        jobs.value = Array.isArray(tasks) ? tasks : [];
        taskErrors.value = Array.isArray(errors) ? errors : [];
      } catch (err) {
        console.error(err);
        if (err.response && err.response.status === 401) {
          router.push({ name: 'SchedulerTable', query: { redirect: router.currentRoute.value.fullPath } });
        } else {
          error.value = '無法載入排程工作資訊，請稍後再試。';
        }
      } finally {
        loading.value = false;
      }
    };
 
    const refreshJobs = () => {
      fetchJobs();
    };

    const uniqueComputers = computed(() => {
      const computers = new Set(ensureJobsArray.value
        .filter(job => job && job.ExtraInfo && job.ExtraInfo.ComputerName)
        .map(job => job.ExtraInfo.ComputerName));
      return Array.from(computers).sort();
    });

    const filteredJobs = computed(() => {
      return ensureJobsArray.value.filter(job => {
        const matchesSearch = searchQuery.value === '' ||
          job.RegistrationInfo.URI.toLowerCase().includes(searchQuery.value.toLowerCase());
        
        const matchesStatus = statusFilter.value === '' ||
          job.ExtraInfo.State === statusFilter.value;

        const matchesComputer = computerFilter.value === '' ||
          job.ExtraInfo.ComputerName === computerFilter.value;
        
        return matchesSearch && matchesStatus && matchesComputer;
      });
    });
 
    onMounted(() => {
      fetchJobs();
    });
 
    return {
      jobs,
      taskErrors,
      loading,
      error,
      searchQuery,
      statusFilter,
      computerFilter,
      filteredJobs,
      uniqueComputers,
      isSidebarOpen,
      refreshJobs,
      handleTaskDisabled,
      handleTaskEnabled,
      handleTaskStarted,
      handleTaskStopped
    };
  }
};
</script>