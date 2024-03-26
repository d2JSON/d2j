import './assets/style.css'
import 'highlight.js/styles/github.css';

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(router)
app.mount('#app')