<template>
  <div class="border rounded p-4 mb-4">
    <div class="flex justify-between items-center mb-2">
      <h3 class="font-semibold">Triggers</h3>
      <button @click="openAdd" class="px-3 py-1 bg-green-600 text-white rounded">新增 Trigger</button>
    </div>
    <div v-if="triggers.length === 0" class="text-gray-500">尚無 triggers</div>
    <ul v-else class="space-y-2">
      <li v-for="(trigger, idx) in triggers" :key="idx" class="flex justify-between items-center border p-2 rounded">
        <span>{{ trigger.cron_expression }}</span>
        <div class="space-x-2">
          <button @click="openEdit(idx)" class="px-2 py-1 bg-blue-500 text-white rounded">編輯</button>
          <button @click="deleteTrigger(idx)" class="px-2 py-1 bg-red-500 text-white rounded">刪除</button>
        </div>
      </li>
    </ul>

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
import { ref, reactive, watch } from 'vue';
import { addTrigger, updateTrigger, deleteTrigger as deleteTriggerApi } from '../../services/taskService';
import TriggerForm from './TriggerForm.vue';
import { useToast } from 'vue-toastification';

const props = defineProps({
  computerId: Number,
  taskName: String,
  triggers: Array,
});

const emit = defineEmits(['refresh']);

const toast = useToast();

const showForm = ref(false);
const isEdit = ref(false);
const currentTrigger = ref({});
const currentIndex = ref(-1);

function openAdd() {
  isEdit.value = false;
  currentTrigger.value = { cron_expression: '' };
  showForm.value = true;
}

function openEdit(index) {
  isEdit.value = true;
  currentIndex.value = index;
  currentTrigger.value = { ...props.triggers[index] };
  showForm.value = true;
}

function closeForm() {
  showForm.value = false;
}

async function handleSubmit(data) {
  try {
    if (isEdit.value) {
      await updateTrigger({
        computer_id: props.computerId,
        task_name: props.taskName,
        cron_expression: data.cron_expression,
        index: currentIndex.value,
      });
      toast.success('Trigger 更新成功');
    } else {
      await addTrigger({
        computer_id: props.computerId,
        task_name: props.taskName,
        cron_expression: data.cron_expression,
      });
      toast.success('Trigger 新增成功');
    }
    emit('refresh');
    closeForm();
  } catch (err) {
    toast.error('操作失敗');
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
    toast.error('刪除失敗');
    console.error(err);
  }
}
</script>
