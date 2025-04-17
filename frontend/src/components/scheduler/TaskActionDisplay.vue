<template>
  <div class="flex items-start gap-x-3">
    <div class="flex-shrink-0">
      <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
    </div>
    <div class="flex-1 space-y-1.5">
      <div class="flex justify-between items-center">
        <p class="text-sm font-medium text-gray-900">執行指令</p>
        <button @click="openAdd" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-600 hover:bg-gray-50 rounded transition-colors group relative" title="新增">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
      <div v-if="formattedActions.length === 0" class="text-sm text-gray-500 italic">無執行指令</div>
      <div v-else class="space-y-3">
        <div v-for="(action, idx) in formattedActions" :key="idx"
          class="group rounded-lg border border-gray-200 overflow-hidden bg-white hover:border-gray-300 transition-colors">
          <div class="flex items-center justify-between px-4 py-2 bg-gray-50 border-b border-gray-200">
            <span class="text-xs font-medium text-gray-600">指令 #{{ idx + 1 }}</span>
            <div class="space-x-2 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
              <button @click="openEdit(idx)" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded transition-colors group relative" title="編輯">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button @click="deleteAction(idx)" class="inline-flex items-center justify-center w-6 h-6 text-xs font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded transition-colors group relative" title="刪除">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
          <div class="p-3">
            <code class="block text-sm font-mono text-gray-800 break-all whitespace-pre-wrap">{{ action.original.Command || '無' }}</code>
            <div v-if="action.original.Arguments" class="mt-2 flex items-center gap-2">
              <span class="text-xs font-medium text-gray-500">參數:</span>
              <code class="text-xs font-mono text-gray-600 bg-gray-50 px-2 py-1 rounded break-all">{{ action.original.Arguments }}</code>
            </div>
            <div v-if="action.original.WorkingDirectory" class="mt-1 flex items-center gap-2">
              <span class="text-xs font-medium text-gray-500">目錄:</span>
              <code class="text-xs font-mono text-gray-600 bg-gray-50 px-2 py-1 rounded break-all">{{ action.original.WorkingDirectory }}</code>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-lg w-96">
        <ActionForm
          :model-value="currentAction"
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
import { addAction, updateAction, deleteAction as deleteActionApi } from '../../services/taskService';
import ActionForm from './ActionForm.vue';
import { useToast } from 'vue-toastification';

const props = defineProps({
  computerId: Number,
  taskName: String,
  actions: Object, // 原始的 job.Actions 資料
});

const emit = defineEmits(['refresh']);
const toast = useToast();

const showForm = ref(false);
const isEdit = ref(false);
const currentAction = ref({});
const currentIndex = ref(-1);

// --- 來自 JobCard 的 computed 邏輯，稍作修改 ---
const formattedActions = computed(() => {
  if (!props.actions || props.actions.length === 0) {
    return [];
  }
  // 假設後端回傳的 actions 陣列結構是 { execute: '...', args: '...', working_directory: '...' }
  // 如果結構不同，需要調整這裡
  return props.actions.Execs.map(action => ({
        original: action 
    }));
});

// --- 來自 TaskActions 的 CRUD 邏輯 ---
function openAdd() {
  isEdit.value = false;
  currentAction.value = { execute: '', args: '', working_directory: '' };
  showForm.value = true;
}

function openEdit(index) {
  isEdit.value = true;
  currentIndex.value = index;
  // 從 formattedActions 取得原始物件來填充表單
  currentAction.value = { ...formattedActions.value[index].original };
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
      execute: formData.execute,
      args: formData.args,
      working_directory: formData.working_directory,
    };
    if (isEdit.value) {
      payload.index = currentIndex.value;
      await updateAction(payload);
      toast.success('Action 更新成功');
    } else {
      await addAction(payload);
      toast.success('Action 新增成功');
    }
    emit('refresh');
    closeForm();
  } catch (err) {
     toast.error(`操作失敗: ${err.response?.data?.message || err.message}`);
    console.error(err);
  }
}

async function deleteAction(index) {
  if (!confirm('確定刪除這個 Action?')) return;
  try {
    await deleteActionApi({
      computer_id: props.computerId,
      task_name: props.taskName,
      index,
    });
    toast.success('Action 刪除成功');
    emit('refresh');
  } catch (err) {
    toast.error(`刪除失敗: ${err.response?.data?.message || err.message}`);
    console.error(err);
  }
}
</script>
