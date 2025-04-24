<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGameStore } from '../stores/game'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const gameStore = useGameStore()

const loading = ref(false)
const createRoomDialogVisible = ref(false)
const newRoomName = ref('')

const fetchLobbyData = async () => {
  loading.value = true
  try {
    const success = await gameStore.fetchGameLobby(route.params.gameId)
    if (!success) {
      ElMessage.error('获取房间列表失败')
    }
  } finally {
    loading.value = false
  }
}

const handleCreateRoom = async () => {
  if (!newRoomName.value) {
    ElMessage.warning('请输入房间名称')
    return
  }

  try {
    const success = await gameStore.createRoom(route.params.gameId, newRoomName.value)
    if (success) {
      ElMessage.success('创建房间成功')
      createRoomDialogVisible.value = false
      newRoomName.value = ''
      await fetchLobbyData()
    } else {
      ElMessage.error('创建房间失败')
    }
  } catch (error) {
    ElMessage.error('创建房间失败')
  }
}

const handleJoinRoom = async (roomId) => {
  try {
    // 获取要加入的房间
    const room = gameStore.lobbyRooms.find(r => r.id === roomId)
    
    // 检查房间是否可以加入
    if (!room || !gameStore.canJoinRoom(room)) {
      ElMessage.warning('该房间已开始游戏或已满员，无法加入')
      return
    }
    
    const success = await gameStore.joinRoom(roomId)
    if (success) {
      router.push(`/room/${roomId}`)
    } else {
      ElMessage.error('加入房间失败')
    }
  } catch (error) {
    ElMessage.error('加入房间失败')
  }
}

const goBack = () => {
  router.push('/games')
}

onMounted(() => {
  fetchLobbyData()
})
</script>

<template>
  <div class="lobby-container">
    <el-header class="header">
      <div class="header-content">
        <div class="header-left">
          <el-button @click="goBack" icon="Back">返回</el-button>
          <h2>{{ gameStore.currentGame?.name }} - 游戏大厅</h2>
        </div>
        <el-button type="primary" @click="createRoomDialogVisible = true">
          创建房间
        </el-button>
      </div>
    </el-header>

    <el-main>
      <el-table
        v-loading="loading"
        :data="gameStore.lobbyRooms"
        style="width: 100%"
      >
        <el-table-column prop="name" label="房间名称" />
        <el-table-column prop="playerCount" label="当前人数" width="120">
          <template #default="{ row }">
            {{ row.playerCount }}/{{ row.maxPlayers }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'waiting' ? 'success' : 'warning'">
              {{ row.status === 'waiting' ? '等待中' : '游戏中' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="120">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :disabled="!gameStore.canJoinRoom(row)"
              @click="handleJoinRoom(row.id)"
            >
              加入
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-main>

    <el-dialog
      v-model="createRoomDialogVisible"
      title="创建房间"
      width="30%"
    >
      <el-form>
        <el-form-item label="房间名称">
          <el-input v-model="newRoomName" placeholder="请输入房间名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createRoomDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateRoom">创建</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.lobby-container {
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

.el-main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
</style>