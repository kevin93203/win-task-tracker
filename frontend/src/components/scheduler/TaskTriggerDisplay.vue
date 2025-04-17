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
        <button @click="openAdd" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-600 hover:bg-gray-50 rounded transition-colors group relative" title="新增">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
      <div class="space-y-2">
        <div v-if="(formattedTriggers.calendarTriggers.length === 0) && (formattedTriggers.timeTriggers.length === 0)" class="text-sm text-gray-500 italic px-3 py-2">
          無觸發程序
        </div>
        <div v-else class="space-y-2">
          <div v-for="(trigger, index) in formattedTriggers.timeTriggers" :key="index"
            class="group rounded-lg border border-gray-200 overflow-hidden bg-white hover:border-gray-300 transition-colors">
            <div class="flex items-center justify-between px-4 py-2 bg-gray-50 border-b border-gray-200">
              <div class="flex items-center gap-2">
                <svg class="h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="text-xs font-medium text-gray-600">時間觸發 #{{ index + 1 }}</span>
              </div>
              <div class="space-x-2 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click="openEdit(index)" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded transition-colors group relative" title="編輯">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button @click="deleteTrigger(index)" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded transition-colors group relative" title="刪除">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
            <div class="p-3">
              <p class="text-sm text-gray-800 break-words">{{ trigger }}</p>
            </div>
          </div>
          <div v-for="(trigger, index) in formattedTriggers.calendarTriggers" :key="index"
            class="group rounded-lg border border-gray-200 overflow-hidden bg-white hover:border-gray-300 transition-colors">
            <div class="flex items-center justify-between px-4 py-2 bg-gray-50 border-b border-gray-200">
              <div class="flex items-center gap-2">
                <svg class="h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <span class="text-xs font-medium text-gray-600">日曆觸發 #{{ index + 1 }}</span>
              </div>
              <div class="space-x-2 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click="openEdit(index)" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded transition-colors group relative" title="編輯">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button @click="deleteTrigger(index)" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded transition-colors group relative" title="刪除">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
            <div class="p-3">
              <p class="text-sm text-gray-800 break-words">{{ trigger }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-lg w-96">
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
  triggers: Object, // 原始的 job.Triggers 資料
});

const emit = defineEmits(['refresh']);
const toast = useToast();

const showForm = ref(false);
const isEdit = ref(false);
const currentTrigger = ref({});
const currentIndex = ref(-1);

// --- 來自 JobCard 的 computed 邏輯，稍作修改 ---
const formattedTriggers = computed(() => ({
  calendarTriggers: (props.triggers.CalendarTriggers || []).map(trigger => convertToChineseSummary(trigger)),
  timeTriggers: (props.triggers.TimeTriggers || []).map(trigger => convertToChineseSummary(trigger))
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

function convertToChineseSummary(scheduleObj) {
  let summary = "";
  
  // 處理開始時間
  const startTime = new Date(scheduleObj.StartBoundary);
  const formattedStartTime = `${startTime.getFullYear()}/${startTime.getMonth() + 1}/${startTime.getDate()} ${startTime.getHours()}:${String(startTime.getMinutes()).padStart(2, '0')}`;
  
  // 處理重複設定
  let intervalDesc = "";
  if (scheduleObj.Repetition && scheduleObj.Repetition.Interval) {
    const interval = scheduleObj.Repetition.Interval;
    if (interval.includes("PT")) {
      const timeValue = interval.replace("PT", "");
      if (timeValue.includes("H")) {
        intervalDesc = `每${timeValue.replace("H", "")}小時`;
      } else if (timeValue.includes("M")) {
        intervalDesc = `每${timeValue.replace("M", "")}分鐘`;
      }
    }
  }
  
  // 處理按日排程
  if (scheduleObj.ScheduleByDay) {
    const daysInterval = scheduleObj.ScheduleByDay.DaysInterval
    summary = `從${formattedStartTime}開始，${intervalDesc}執行一次，每${daysInterval}天重複`;
  }
  // 處理按週排程
  else if (scheduleObj.ScheduleByWeek) {
    const daysOfWeek = scheduleObj.ScheduleByWeek.DaysOfWeek;
    const selectedDays = [];
    const weeksInterval = scheduleObj.ScheduleByWeek.WeeksInterval
    
    if (daysOfWeek.Sunday) selectedDays.push("日");
    if (daysOfWeek.Monday) selectedDays.push("一");
    if (daysOfWeek.Tuesday) selectedDays.push("二");
    if (daysOfWeek.Wednesday) selectedDays.push("三");
    if (daysOfWeek.Thursday) selectedDays.push("四");
    if (daysOfWeek.Friday) selectedDays.push("五");
    if (daysOfWeek.Saturday) selectedDays.push("六");
    
    summary = `從${formattedStartTime}開始，${intervalDesc}執行一次，每${weeksInterval}週${selectedDays.join("、")}重複`;
  }
  // 處理按月排程
  else if (scheduleObj.ScheduleByMonth) {
    const months = scheduleObj.ScheduleByMonth.Months;
    const selectedMonths = [];
    
    if (months.January) selectedMonths.push("1");
    if (months.February) selectedMonths.push("2");
    if (months.March) selectedMonths.push("3");
    if (months.April) selectedMonths.push("4");
    if (months.May) selectedMonths.push("5");
    if (months.June) selectedMonths.push("6");
    if (months.July) selectedMonths.push("7");
    if (months.August) selectedMonths.push("8");
    if (months.September) selectedMonths.push("9");
    if (months.October) selectedMonths.push("10");
    if (months.November) selectedMonths.push("11");
    if (months.December) selectedMonths.push("12");
    
    const daysOfMonth = scheduleObj.ScheduleByMonth.DaysOfMonth.Days.join("、");
    
    summary = `從${formattedStartTime}開始，${intervalDesc}執行一次，每年${selectedMonths.join("、")}月的${daysOfMonth}號重複`;
  }
  // 如果沒有特定排程，只有重複設定
  else {
    const duration = scheduleObj.Repetition && scheduleObj.Repetition.Duration ? "，持續一天" : "";
    summary = `從${formattedStartTime}開始，${intervalDesc}執行一次${duration}`;
  }
  
  return summary;
}

</script>
