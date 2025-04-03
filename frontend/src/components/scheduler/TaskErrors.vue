<template>
  <div v-if="errors && errors.length > 0" class="mb-6">
    <div class="bg-red-50 border border-red-200 rounded-lg p-4">
      <h3 class="text-red-800 font-medium mb-2">執行錯誤</h3>
      <ul class="space-y-2">
        <li v-for="error in errors" :key="error.computer_name + error.username" class="text-red-700">
          <span class="font-medium">{{ error.computer_name }}</span>
          <span class="text-red-600"> ({{ error.username }}):</span>
          <span class="text-red-600 ml-1">{{ error.error }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    errors: {
      type: Array,
      required: true
    }
  }
}
</script>
