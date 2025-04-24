import './assets/main.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

const app = createApp(App)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 配置axios响应拦截器处理统一响应格式
axios.interceptors.response.use(
  response => {
    // 如果响应已经是我们期望的格式，直接返回
    if (response.data && response.data.hasOwnProperty('code')) {
      // 如果code不为0，说明有错误
      if (response.data.code !== 0) {
        // 可以在这里统一处理错误，比如显示错误消息
        console.error(response.data.msg || '请求失败')
        return Promise.reject(response.data)
      }
      // 返回data字段中的数据
      return response.data.data
    }
    // 如果不是统一格式，保持原样返回
    return response
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 全局挂载axios
app.config.globalProperties.$axios = axios

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

app.mount('#app')
