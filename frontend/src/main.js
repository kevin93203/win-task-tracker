import { createApp } from 'vue'
import { reactive } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

// Create a reactive global state object
export const globalState = reactive({
  sidebarOpen: true
})

const app = createApp(App)
app.use(router)
app.mount('#app')
