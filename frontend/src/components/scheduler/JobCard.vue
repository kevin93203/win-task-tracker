<template>
  <div class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
    <div class="p-4 border-b border-gray-100">
      <div class="flex justify-between items-center mb-2">
        <h2 class="font-bold text-lg text-gray-800 truncate" :title="job.ExtraInfo.TaskName">
          {{ job.ExtraInfo.TaskName }}
        </h2>
        <span :class="getStatusClass(job.ExtraInfo.State)" class="px-2 py-1 rounded-full text-xs font-bold">
          {{ getStatusChinese(job.ExtraInfo.State) }}
        </span>
      </div>
      <p v-if="job.RegistrationInfo.Description" class="text-sm text-gray-600">
        {{ job.RegistrationInfo.Description }}
      </p>
    </div>
    
    <div class="p-4 space-y-3">
      <div class="flex items-start">
        <span class="text-gray-500 w-24 flex-shrink-0">執行指令:</span>
        <div class="flex-1">
          <div v-if="commands.length > 0" class="space-y-2">
            <div v-for="(command, index) in commands" :key="index"
              class="group">
              <div class="text-xs text-gray-500 mb-1">指令 #{{ index + 1 }}</div>
              <code class="block w-full bg-gray-50 text-sm p-3 rounded-lg border border-gray-200 font-mono text-gray-800 break-all whitespace-pre-wrap group-hover:bg-gray-100 transition-colors">{{ command.Command || '無' }}</code>
              <div v-if="command.Arguments" class="mt-2">
                <span class="text-xs text-gray-500">參數:</span>
                <code class="ml-2 px-2 py-1 bg-gray-50 text-xs rounded border border-gray-200 font-mono">{{ command.Arguments }}</code>
              </div>
            </div>
          </div>
          <div v-else class="text-gray-500 italic">無執行指令</div>
        </div>
      </div>
      
      <div class="flex items-start">
        <span class="text-gray-500 w-24 flex-shrink-0">觸發程序:</span>
        <div class="flex-1 space-y-2">
          <div v-if="triggers.timeTriggers.length" class="space-y-1">
            <div v-for="(trigger, index) in triggers.timeTriggers" :key="index"
              class="flex items-center gap-2 bg-blue-50 px-3 py-2 rounded-md">
              <svg class="h-4 w-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="text-sm text-blue-700">{{ trigger }}</span>
            </div>
          </div>
          <div v-if="triggers.calendarTriggers.length" class="space-y-1">
            <div v-for="(trigger, index) in triggers.calendarTriggers" :key="index"
              class="flex items-center gap-2 bg-purple-50 px-3 py-2 rounded-md">
              <svg class="h-4 w-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="text-sm text-purple-700">{{ trigger }}</span>
            </div>
          </div>
          <div v-if="!triggers.timeTriggers.length && !triggers.calendarTriggers.length"
            class="text-gray-500 italic text-sm">
            無觸發程序
          </div>
        </div>
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
</template>

<script>
export default {
  name: 'JobCard',
  props: {
    job: {
      type: Object,
      required: true
    }
  },
  computed: {
    commands() {
      if (this.job.Actions.Execs) {
        return this.job.Actions.Execs;
      } else if (this.job.Actions.Exec) {
        return [this.job.Actions.Exec];
      }
      return [];
    },
    triggers() {
      const triggers = this.job.Triggers;
      const result = {
        timeTriggers: [],
        calendarTriggers: []
      };

      if (triggers && triggers.length > 0) {
        if (triggers[0].TimeTriggers) {
          result.timeTriggers = triggers[0].TimeTriggers.map(trigger => {
            let text = `${this.formatDateTime(trigger.StartBoundary)}`;
            if (trigger.Repetition.Interval) {
              text += ` (每 ${trigger.Repetition.Interval.replace('PT', '').replace('H', ' 小時').replace('M', ' 分鐘')})`;
            }
            return text;
          });
        }

        if (triggers[0].CalendarTriggers) {
          result.calendarTriggers = triggers[0].CalendarTriggers.map(trigger => {
            let text = `${this.formatDateTime(trigger.StartBoundary)}`;
            if (trigger.ScheduleByDay.DaysInterval) {
              text += ` (每 ${trigger.ScheduleByDay.DaysInterval} 天)`;
            }
            return text;
          });
        }
      }

      return result;
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
          return 'bg-green-100 text-green-800';
        case 'Running':
          return 'bg-blue-100 text-blue-800';
        case 'Disabled':
          return 'bg-gray-100 text-gray-800';
        default:
          return 'bg-gray-100 text-gray-800';
      }
    },
    getResultClass(result) {
      if (!result || result === '0') {
        return 'text-green-600 font-medium';
      }
      return 'text-red-600 font-medium';
    }
  }
}
</script> 