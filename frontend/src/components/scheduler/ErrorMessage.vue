<template>
  <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded mb-6">
    <p class="font-bold">載入失敗</p>
    <p>{{ message }}</p>
  </div>
</template>

<script>
export default {
  name: 'ErrorMessage',
  props: {
    message: {
      type: String,
      required: true
    }
  }
}
</script> 