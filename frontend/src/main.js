import { createApp } from 'vue'
import { reactive } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

// Import Toast
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

// Toast options
const toastOptions = {
  position: 'top-right',
  timeout: 3000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: 'button',
  icon: true,
  rtl: false
}

// Create a reactive global state object
export const globalState = reactive({
  sidebarOpen: true
})

const app = createApp(App)
app.use(router)
app.use(Toast, toastOptions)
app.mount('#app')
