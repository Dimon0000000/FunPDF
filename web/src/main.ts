import { createApp } from 'vue'
import { createPinia } from 'pinia'
import '@fortawesome/fontawesome-free/css/all.min.css'
import './styles/global.css'
import App from './App.vue'

createApp(App)
  .use(createPinia())
  .mount('#app')
