import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'

import Axios from 'axios'

const curSrc = document.currentScript as HTMLScriptElement
const baseUrl = new URL(curSrc.src)
const base = baseUrl.pathname.substring(0, baseUrl.pathname.indexOf('/admin/js/'))

Axios.defaults.baseURL = base

createApp(App).use(router).mount('#app')
