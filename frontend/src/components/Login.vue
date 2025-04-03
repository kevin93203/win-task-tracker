<template>
  <div class="login-container">
    <div class="login-box">
      <h2>登入</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">帳號</label>
          <input
            type="text"
            id="username"
            v-model="username"
            required
            placeholder="請輸入帳號"
          />
        </div>
        <div class="form-group">
          <label for="password">密碼</label>
          <input
            type="password"
            id="password"
            v-model="password"
            required
            placeholder="請輸入密碼"
          />
        </div>
        <button type="submit" :disabled="isLoading">
          {{ isLoading ? '登入中...' : '登入' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import authService from '../services/authService';

export default {
  name: 'Login',
  setup() {
    const username = ref('');
    const password = ref('');
    const isLoading = ref(false);
    const router = useRouter();

    const handleLogin = async () => {
      try {
        isLoading.value = true;
        await authService.login({
          username: username.value,
          password: password.value,
        });
        
        // Redirect to home page after successful login
        router.push('/');
      } catch (error) {
        console.error('Login failed:', error);
        alert('登入失敗，請檢查帳號密碼是否正確');
      } finally {
        isLoading.value = false;
      }
    };

    return {
      username,
      password,
      isLoading,
      handleLogin,
    };
  },
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-box {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: #666;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

button {
  width: 100%;
  padding: 0.75rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #45a049;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>
