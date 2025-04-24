import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}')
  }),
  
  actions: {
    async login(username, password) {
      try {
        // 这里替换为实际的API地址
        const response = await axios.post('/api/login', {
          username,
          password
        })
        // 响应拦截器已经处理了统一格式，直接获取data中的数据
        const { token, user } = response
        this.token = token
        this.userInfo = user
        localStorage.setItem('token', token)
        localStorage.setItem('userInfo', JSON.stringify(user))
        return true
      } catch (error) {
        console.error('Login failed:', error)
        return false
      }
    },

    async register(username, password) {
      try {
        // 这里替换为实际的API地址
        // 响应拦截器已经处理了统一格式，不需要额外处理response.data
        await axios.post('/api/register', {
          username,
          password
        })
        return true
      } catch (error) {
        console.error('Registration failed:', error)
        return false
      }
    },

    logout() {
      this.token = ''
      this.userInfo = {}
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  }
})