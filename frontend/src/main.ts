import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'

import '@exampledev/new.css/new.css'

import Axios from 'axios'
Axios.defaults.baseURL = '/guardmech'

createApp(App).use(router).mount('#app')
