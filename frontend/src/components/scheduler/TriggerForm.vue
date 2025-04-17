<template>
  <div class="p-4">
    <h2 class="text-lg font-semibold mb-2">{{ isEdit ? '編輯 Trigger' : '新增 Trigger' }}</h2>
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block mb-1 font-medium">Cron 表達式</label>
        <input
          v-model="form.cron_expression"
          type="text"
          class="w-full border rounded px-3 py-2"
          placeholder="例如：0 8 * * *"
          required
        />
        <p class="text-sm text-gray-500 mt-1">格式：分 時 日 月 星期</p>
      </div>
      <div class="flex justify-end space-x-2">
        <button type="button" @click="$emit('cancel')" class="px-4 py-2 border rounded">取消</button>
        <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded">{{ isEdit ? '更新' : '新增' }}</button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, watch, toRefs } from 'vue';

const props = defineProps({
  modelValue: Object,
  isEdit: Boolean,
});

const emit = defineEmits(['submit', 'cancel']);

const form = reactive({
  cron_expression: '',
});

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      form.cron_expression = val.cron_expression || '';
    }
  },
  { immediate: true }
);

function handleSubmit() {
  emit('submit', { ...form });
}
</script>
