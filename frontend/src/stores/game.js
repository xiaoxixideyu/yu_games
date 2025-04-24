import { defineStore } from 'pinia'
import axios from 'axios'

export const useGameStore = defineStore('game', {
  state: () => ({
    gameList: [],
    currentGame: null,
    lobbyRooms: [],
    currentRoom: null
  }),
  
  actions: {
    async fetchGameList() {
      try {
        // 这里替换为实际的API地址
        // 响应拦截器已经处理了统一格式，直接获取data中的数据
        const data = await axios.get('/api/games')
        this.gameList = data
        return true
      } catch (error) {
        console.error('Failed to fetch game list:', error)
        return false
      }
    },

    async fetchGameLobby(gameId) {
      try {
        // 这里替换为实际的API地址
        // 响应拦截器已经处理了统一格式，直接获取data中的数据
        const data = await axios.get(`/api/games/${gameId}/rooms`)
        this.lobbyRooms = data
        this.currentGame = this.gameList.find(game => game.id === gameId)
        return true
      } catch (error) {
        console.error('Failed to fetch game lobby:', error)
        return false
      }
    },

    async joinRoom(roomId) {
      try {
        // 这里替换为实际的API地址
        // 响应拦截器已经处理了统一格式，直接获取data中的数据
        const data = await axios.post(`/api/rooms/${roomId}/join`)
        this.currentRoom = data
        return true
      } catch (error) {
        console.error('Failed to join room:', error)
        return false
      }
    },

    async createRoom(gameId, roomName) {
      try {
        // 这里替换为实际的API地址
        // 响应拦截器已经处理了统一格式，直接获取data中的数据
        const data = await axios.post(`/api/games/${gameId}/rooms`, {
          name: roomName
        })
        this.currentRoom = data
        return true
      } catch (error) {
        console.error('Failed to create room:', error)
        return false
      }
    },

    leaveRoom() {
      this.currentRoom = null
    },
    
    // 玩家准备状态切换 - 不再使用REST API，改为使用WebSocket
    // 此方法保留为空，实际的toggleReady操作将在GameRoom.vue中通过socket.emit实现
    toggleReady() {
      // 方法保留但不执行任何操作，因为准备状态通过WebSocket事件处理
      // 房间状态将通过roomUpdate事件自动更新
      return true
    },
    
    // 检查房间状态
    canJoinRoom(room) {
      // 只能加入未开始且未满员的房间
      return room.status === 'waiting' && room.players.length < room.maxPlayers
    },
    
    // 检查游戏是否可以开始
    canStartGame(room) {
      if (!room) return false
      
      // 游戏开始条件：参与人数大于1且所有玩家都已准备
      return room.players.length > 1 && 
             room.players.every(player => player.ready === true)
    }
  }
})