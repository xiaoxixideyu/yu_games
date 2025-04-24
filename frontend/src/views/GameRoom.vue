<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGameStore } from '../stores/game'
import { ElMessage } from 'element-plus'
import PokerGame from '../components/PokerGame.vue'
import io from 'socket.io-client'

const route = useRoute()
const router = useRouter()
const gameStore = useGameStore()
const socket = ref(null)
const currentUser = JSON.parse(localStorage.getItem('userInfo') || '{}')
const isReady = ref(false)

// 检查当前用户是否已准备
const updateReadyStatus = () => {
  if (!gameStore.currentRoom) return
  const currentPlayer = gameStore.currentRoom.players.find(p => p.id === currentUser.id)
  isReady.value = currentPlayer?.ready || false
}

// 切换准备状态 - 使用WebSocket事件
const toggleReady = () => {
  if (!gameStore.currentRoom) return
  
  // 通过WebSocket发送toggleReady事件
  socket.value.emit('toggleReady', { 
    roomId: gameStore.currentRoom.id 
  })
  // 房间状态将通过roomUpdate事件自动更新
}

// 游戏状态
const gameStatus = computed(() => {
  return gameStore.currentRoom?.status || 'waiting'
})

// 检查是否所有玩家都已准备
const allPlayersReady = computed(() => {
  return gameStore.canStartGame(gameStore.currentRoom)
})

// 建立WebSocket连接
const connectSocket = () => {
  // 连接到WebSocket服务器
  socket.value = io(import.meta.env.VITE_API_BASE_URL, {
    query: {
      roomId: route.params.roomId
    }
  })
  
  // 监听房间更新事件
  socket.value.on('roomUpdate', (room) => {
    gameStore.currentRoom = room
    updateReadyStatus()
  })
  
  // 监听游戏状态更新事件
  socket.value.on('gameState', (state) => {
    // 游戏状态更新处理
    ElMessage.info(`游戏状态更新: ${state.round || '等待中'}`)
  })
  
  // 监听游戏开始事件
  socket.value.on('gameStart', () => {
    ElMessage.success('游戏开始！')
  })
  
  // 监听玩家准备状态变化事件
  socket.value.on('playerReady', (data) => {
    // 此事件可能不需要特别处理，因为roomUpdate事件会更新整个房间状态
    // 但如果需要特定的UI反馈，可以在这里处理
    const playerName = gameStore.currentRoom?.players.find(p => p.id === data.playerId)?.username || '玩家'
    ElMessage.info(`${playerName} ${data.ready ? '已准备' : '取消准备'}`)
  })
  
  // 监听游戏结果事件
  socket.value.on('gameResult', (result) => {
    ElMessage({
      message: `游戏结束: ${result.message}`,
      type: 'success',
      duration: 5000,
      showClose: true
    })
  })
  
  // 监听错误消息
  socket.value.on('error', (error) => {
    ElMessage.error(error.message)
  })
}

// 离开房间
const handleLeaveRoom = () => {
  if (socket.value) {
    socket.value.disconnect()
    socket.value = null
  }
  gameStore.leaveRoom()
  router.push(`/lobby/${gameStore.currentGame.id}`)
}

const goBack = () => {
  handleLeaveRoom()
}

onMounted(() => {
  connectSocket()
  updateReadyStatus()
})

onBeforeUnmount(() => {
  handleLeaveRoom()
})
</script>

<template>
  <div class="room-container">
    <el-header class="header">
      <div class="header-content">
        <div class="header-left">
          <el-button @click="goBack" icon="Back">退出房间</el-button>
          <h2>{{ gameStore.currentRoom?.name }}</h2>
        </div>
        <div class="room-info">
          <span>房间号：{{ route.params.roomId }}</span>
        </div>
      </div>
    </el-header>

    <el-main>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="player-list">
            <template #header>
              <div class="card-header">
                <span>玩家列表</span>
              </div>
            </template>
            <div
              v-for="player in gameStore.currentRoom?.players"
              :key="player.id"
              class="player-item"
            >
              <el-avatar :size="32" :src="player.avatar">
                {{ player.username.charAt(0).toUpperCase() }}
              </el-avatar>
              <span class="player-name">{{ player.username }}</span>
              <el-tag
                size="small"
                :type="player.ready ? 'success' : 'info'"
              >
                {{ player.ready ? '已准备' : '未准备' }}
              </el-tag>
            </div>
          </el-card>
        </el-col>
        <el-col :span="18">
          <el-card class="game-area">
            <!-- 游戏区域 -->
            <div v-if="gameStatus === 'playing'">
              <PokerGame :socket="socket" :roomId="route.params.roomId" />
            </div>
            <!-- 等待区域 -->
            <div v-else class="game-placeholder">
              <div class="waiting-area">
                <h3>{{ gameStore.currentRoom?.name }} - 等待玩家准备</h3>
                <p>当所有玩家准备就绪且人数大于1时，游戏将自动开始</p>
                <el-button 
                  type="primary" 
                  size="large" 
                  @click="toggleReady"
                  :type="isReady ? 'danger' : 'primary'"
                >
                  {{ isReady ? '取消准备' : '准备' }}
                </el-button>
                <div class="ready-status" v-if="gameStore.currentRoom?.players.length > 0">
                  <p>准备状态: {{ gameStore.currentRoom.players.filter(p => p.ready).length }}/{{ gameStore.currentRoom.players.length }}</p>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </div>
</template>

<style scoped>
.room-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.room-info {
  color: #606266;
}

.el-main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.player-list {
  margin-bottom: 20px;
}

.player-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}

.player-item:last-child {
  border-bottom: none;
}

.player-name {
  margin: 0 10px;
  flex: 1;
}

.game-area {
  height: 600px;
}

.game-placeholder {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.waiting-area {
  text-align: center;
  padding: 30px;
}

.waiting-area h3 {
  margin-bottom: 15px;
  font-size: 24px;
}

.waiting-area p {
  margin-bottom: 20px;
  color: #606266;
}

.ready-status {
  margin-top: 20px;
  font-size: 16px;
  font-weight: bold;
}

.ready-status p {
  color: #409EFF;
}
</style>