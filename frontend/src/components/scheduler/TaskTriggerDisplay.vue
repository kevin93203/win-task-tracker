<template>
  <div class="flex items-start gap-x-3">
    <div class="flex-shrink-0">
      <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </div>
    <div class="flex-1 space-y-1.5">
      <div class="flex justify-between items-center">
        <p class="text-sm font-medium text-gray-900">觸發程序</p>
        <button @click="openAdd" class="px-2 py-0.5 text-xs bg-green-600 text-white rounded hover:bg-green-700 transition-colors">新增</button>
      </div>
      <div class="space-y-2">
        <div v-if="(formattedTriggers.calendarTriggers.length === 0) && (formattedTriggers.timeTriggers.length === 0)" class="text-sm text-gray-500 italic px-3 py-2">
          無觸發程序
        </div>
        <div v-else class="space-y-2">
                <div v-for="(trigger, index) in formattedTriggers.timeTriggers" :key="index"
                  class="flex items-center gap-2 px-3 py-2 rounded-lg border border-blue-200 bg-blue-50 text-blue-700">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span class="text-sm break-words">{{ trigger }}</span>
                  <div class="space-x-1 flex-shrink-0">
                    <button @click="openEdit(idx)" class="px-1.5 py-0.5 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">編輯</button>
                    <button @click="deleteTrigger(idx)" class="px-1.5 py-0.5 text-xs bg-red-500 text-white rounded hover:bg-red-600 transition-colors">刪除</button>
                  </div>
                </div>
              <div v-for="(trigger, index) in formattedTriggers.calendarTriggers" :key="index"
                  class="flex items-center gap-2 px-3 py-2 rounded-lg border border-purple-200 bg-purple-50 text-purple-700">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  <span class="text-sm break-words">{{ trigger }}</span>
                  <div class="space-x-1 flex-shrink-0">
                    <button @click="openEdit(idx)" class="px-1.5 py-0.5 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">編輯</button>
                    <button @click="deleteTrigger(idx)" class="px-1.5 py-0.5 text-xs bg-red-500 text-white rounded hover:bg-red-600 transition-colors">刪除</button>
                  </div>
              </div>
        </div>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded shadow w-96">
        <TriggerForm
          :model-value="currentTrigger"
          :is-edit="isEdit"
          @submit="handleSubmit"
          @cancel="closeForm"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { addTrigger, updateTrigger, deleteTrigger as deleteTriggerApi } from '../../services/taskService';
import TriggerForm from './TriggerForm.vue';
import { useToast } from 'vue-toastification';

const props = defineProps({
  computerId: Number,
  taskName: String,
  triggers: Array, // 原始的 job.Triggers 資料
});

const emit = defineEmits(['refresh']);
const toast = useToast();

const showForm = ref(false);
const isEdit = ref(false);
const currentTrigger = ref({});
const currentIndex = ref(-1);

// // --- 來自 JobCard 的 computed 邏輯，稍作修改 ---
// const formattedTriggers = computed(() => {
//   if (!props.triggers || props.triggers.length === 0) {
//     return [];
//   }
//   // 假設後端回傳的 triggers 陣列結構是 { cron_expression: '...' }
//   // 如果結構不同，需要調整這裡
//   return props.triggers.map(trigger => ({
//     display: trigger.cron_expression, // 直接顯示 cron 字串
//     original: trigger // 保留原始物件供編輯使用
//   }));
// });
function formatDateTime(dateTime) {
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
    }


// --- 來自 JobCard 的 computed 邏輯，稍作修改 ---
const formattedTriggers = computed(() => ({
  calendarTriggers: (props.triggers.CalendarTriggers || []).map(trigger => {
            let text = `${formatDateTime(trigger.StartBoundary)}`;
            if (trigger.ScheduleByDay) {
              text += ` (每 ${trigger.ScheduleByDay.DaysInterval} 天)`;
            } else if (trigger.ScheduleByWeek) {
              let weekdays = []
              if (trigger.ScheduleByWeek.DaysOfWeek.Sunday) {
                weekdays.push('日')
              }
              if (trigger.ScheduleByWeek.DaysOfWeek.Monday) {
                weekdays.push('一')
              }
              if (trigger.ScheduleByWeek.DaysOfWeek.Tuesday) {
                weekdays.push('二')
              }
              if (trigger.ScheduleByWeek.DaysOfWeek.Wednesday) {
                weekdays.push('三')
              }
              if (trigger.ScheduleByWeek.DaysOfWeek.Thursday) {
                weekdays.push('四')
              }
              if (trigger.ScheduleByWeek.DaysOfWeek.Friday) {
                weekdays.push('五')
              }
              if (trigger.ScheduleByWeek.DaysOfWeek.Saturday) {
                weekdays.push('六')
              }
              text += ` (每 ${trigger.ScheduleByWeek.WeeksInterval} 個星期的 星期${weekdays.join('、')})`;
            }
            return text;
          }),
  timeTriggers: (props.triggers.TimeTriggers || []).map(trigger => {
            let text = `${formatDateTime(trigger.StartBoundary)}`;
            if (trigger.Repetition && trigger.Repetition.Interval) {
              text += ` (每 ${trigger.Repetition.Interval.replace('PT', '').replace('H', ' 小時').replace('M', ' 分鐘')})`;
            }
            return text;
          })
}));

// --- 來自 TaskTriggers 的 CRUD 邏輯 ---
function openAdd() {
  isEdit.value = false;
  currentTrigger.value = { cron_expression: '' }; // 表單只需要 cron_expression
  showForm.value = true;
}

function openEdit(index) {
  isEdit.value = true;
  currentIndex.value = index;
  // 從 formattedTriggers 取得原始物件來填充表單
  currentTrigger.value = { cron_expression: formattedTriggers.value[index].original.cron_expression };
  showForm.value = true;
}

function closeForm() {
  showForm.value = false;
}

async function handleSubmit(formData) {
  try {
    const payload = {
      computer_id: props.computerId,
      task_name: props.taskName,
      cron_expression: formData.cron_expression,
    };
    if (isEdit.value) {
      payload.index = currentIndex.value;
      await updateTrigger(payload);
      toast.success('Trigger 更新成功');
    } else {
      await addTrigger(payload);
      toast.success('Trigger 新增成功');
    }
    emit('refresh');
    closeForm();
  } catch (err) {
    toast.error(`操作失敗: ${err.response?.data?.message || err.message}`);
    console.error(err);
  }
}

async function deleteTrigger(index) {
  if (!confirm('確定刪除這個 Trigger?')) return;
  try {
    await deleteTriggerApi({
      computer_id: props.computerId,
      task_name: props.taskName,
      index,
    });
    toast.success('Trigger 刪除成功');
    emit('refresh');
  } catch (err) {
    toast.error(`刪除失敗: ${err.response?.data?.message || err.message}`);
    console.error(err);
  }
}

// --- 來自 JobCard 的 formatDateTime (如果需要顯示時間) ---
// function formatDateTime(dateTime) { ... }
// 目前直接顯示 cron 字串，暫不需要
</script>
