<template>
  <div class="p-4">
    <h2 class="text-lg font-semibold mb-2">{{ isEdit ? '編輯指令' : '新增指令' }}</h2>
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block mb-1 font-medium">執行檔路徑</label>
        <input
          v-model="form.execute"
          type="text"
          class="w-full border rounded px-3 py-2"
          placeholder="例如：C:\\Windows\\System32\\notepad.exe"
          required
        />
      </div>
      <div>
        <label class="block mb-1 font-medium">參數</label>
        <input
          v-model="form.args"
          type="text"
          class="w-full border rounded px-3 py-2"
          placeholder="例如：/A test.txt"
        />
      </div>
      <div>
        <label class="block mb-1 font-medium">工作目錄</label>
        <input
          v-model="form.working_directory"
          type="text"
          class="w-full border rounded px-3 py-2"
          placeholder="例如：C:\\MyFolder"
        />
      </div>
      <div class="flex justify-end space-x-2">
        <button type="button" @click="$emit('cancel')" class="px-4 py-2 border rounded">取消</button>
        <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded">{{ isEdit ? '更新' : '新增' }}</button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue';

const props = defineProps({
  modelValue: Object,
  isEdit: Boolean,
});

const emit = defineEmits(['submit', 'cancel']);

const form = reactive({
  execute: '',
  args: '',
  working_directory: '',
});

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      form.execute = val.Command || '';
      form.args = val.Arguments || '';
      form.working_directory = val.WorkingDirectory || '';
    }
  },
  { immediate: true }
);

function handleSubmit() {
  emit('submit', { ...form });
}
</script>
