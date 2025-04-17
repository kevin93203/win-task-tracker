<template>
  <div class="border rounded p-4 mb-4">
    <div class="flex justify-between items-center mb-2">
      <h3 class="font-semibold">Actions</h3>
      <button @click="openAdd" class="px-3 py-1 bg-green-600 text-white rounded">新增 Action</button>
    </div>
    <div v-if="actions.length === 0" class="text-gray-500">尚無 actions</div>
    <ul v-else class="space-y-2">
      <li v-for="(action, idx) in actions" :key="idx" class="flex justify-between items-center border p-2 rounded">
        <span>{{ action.execute }}</span>
        <div class="space-x-2">
          <button @click="openEdit(idx)" class="px-2 py-1 bg-blue-500 text-white rounded">編輯</button>
          <button @click="deleteAction(idx)" class="px-2 py-1 bg-red-500 text-white rounded">刪除</button>
        </div>
      </li>
    </ul>

    <div v-if="showForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded shadow w-96">
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
import { ref } from 'vue';
import { addAction, updateAction, deleteAction as deleteActionApi } from '../../services/taskService';
import ActionForm from './ActionForm.vue';
import { useToast } from 'vue-toastification';

const props = defineProps({
  computerId: Number,
  taskName: String,
  actions: Array,
});

const emit = defineEmits(['refresh']);

const toast = useToast();

const showForm = ref(false);
const isEdit = ref(false);
const currentAction = ref({});
const currentIndex = ref(-1);

function openAdd() {
  isEdit.value = false;
  currentAction.value = { execute: '', args: '', working_directory: '' };
  showForm.value = true;
}

function openEdit(index) {
  isEdit.value = true;
  currentIndex.value = index;
  currentAction.value = { ...props.actions[index] };
  showForm.value = true;
}

function closeForm() {
  showForm.value = false;
}

async function handleSubmit(data) {
  try {
    if (isEdit.value) {
      await updateAction({
        computer_id: props.computerId,
        task_name: props.taskName,
        execute: data.execute,
        args: data.args,
        working_directory: data.working_directory,
        index: currentIndex.value,
      });
      toast.success('Action 更新成功');
    } else {
      await addAction({
        computer_id: props.computerId,
        task_name: props.taskName,
        execute: data.execute,
        args: data.args,
        working_directory: data.working_directory,
      });
      toast.success('Action 新增成功');
    }
    emit('refresh');
    closeForm();
  } catch (err) {
    toast.error('操作失敗');
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
    toast.error('刪除失敗');
    console.error(err);
  }
}
</script>
