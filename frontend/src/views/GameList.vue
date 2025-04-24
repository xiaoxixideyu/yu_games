<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useGameStore } from '../stores/game'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

const loading = ref(false)

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const enterGameLobby = (gameId) => {
  router.push(`/lobby/${gameId}`)
}

const fetchGames = async () => {
  loading.value = true
  try {
    const success = await gameStore.fetchGameList()
    if (!success) {
      ElMessage.error('获取游戏列表失败')
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchGames()
})
</script>

<template>
  <div class="game-list-container">
    <el-header class="header">
      <div class="header-content">
        <h2>游戏大厅</h2>
        <div class="user-info">
          <span>{{ userStore.userInfo.username }}</span>
          <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
      </div>
    </el-header>

    <el-main>
      <el-row :gutter="20">
        <el-col
          v-for="game in gameStore.gameList"
          :key="game.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card
            class="game-card"
            :body-style="{ padding: '0px' }"
            @click="enterGameLobby(game.id)"
          >
            <img :src="game.image" class="game-image" />
            <div class="game-info">
              <h3>{{ game.name }}</h3>
              <p>{{ game.description }}</p>
              <div class="game-stats">
                <span>在线: {{ game.onlinePlayers }}</span>
                <span>房间: {{ game.roomCount }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </div>
</template>

<style scoped>
.game-list-container {
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

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.el-main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.game-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.3s;
}

.game-card:hover {
  transform: translateY(-5px);
}

.game-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.game-info {
  padding: 15px;
}

.game-info h3 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.game-info p {
  margin: 10px 0;
  color: #606266;
  font-size: 14px;
}

.game-stats {
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 12px;
}
</style>