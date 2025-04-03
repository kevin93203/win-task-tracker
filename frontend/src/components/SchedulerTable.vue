<template>
  <div class="flex">
    <Sidebar ref="sidebar" />
    <div class="flex-1 p-8 transition-all duration-300" :class="{ 'ml-64': isSidebarOpen, 'ml-0': !isSidebarOpen }">
      <SchedulerHeader />
      
      <LoadingSpinner v-if="loading" />
      <ErrorMessage v-else-if="error" :message="error" />
   
      <div v-else>
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

export default {
  components: {
    Sidebar,
    SchedulerHeader,
    SchedulerFilters,
    JobCard,
    LoadingSpinner,
    ErrorMessage,
    NoResults
  },
  setup() {
    const router = useRouter();
    const sidebar = ref(null);
    const jobs = ref([]);
    const loading = ref(true);
    const error = ref(null);
    const searchQuery = ref('');
    const statusFilter = ref('Ready');
    const computerFilter = ref('');
    const isSidebarOpen = computed(() => sidebar.value?.isOpen ?? true);

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
        jobs.value = response.data;
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
      const computers = new Set(jobs.value.map(job => job.ExtraInfo.ComputerName));
      return Array.from(computers).sort();
    });

    const filteredJobs = computed(() => {
      return jobs.value.filter(job => {
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
      loading,
      error,
      searchQuery,
      statusFilter,
      computerFilter,
      uniqueComputers,
      filteredJobs,
      fetchJobs,
      refreshJobs,
      sidebar,
      isSidebarOpen
    };
  }
};
</script>