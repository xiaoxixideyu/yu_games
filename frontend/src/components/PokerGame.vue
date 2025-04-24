<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useGameStore } from '../stores/game'
import { ElMessage } from 'element-plus'

const props = defineProps({
  socket: Object,
  roomId: String
})

const gameStore = useGameStore()
const currentUser = JSON.parse(localStorage.getItem('userInfo') || '{}')

// 游戏状态
const gameState = ref({
  status: 'waiting', // waiting, playing, ended
  players: [],
  currentPlayer: null,
  communityCards: [],
  pot: 0,
  currentBet: 0,
  dealerPosition: 0,
  round: 'preflop' // preflop, flop, turn, river, showdown
})

// 玩家手牌
const playerHand = ref([])

// 玩家操作
const playerActions = ref({
  canCheck: false,
  canCall: false,
  canRaise: false,
  canFold: true
})

// 下注金额
const betAmount = ref(0)
const minBet = computed(() => gameState.value.currentBet * 2 || 10)

// 是否轮到当前玩家操作
const isPlayerTurn = computed(() => {
  return gameState.value.currentPlayer && 
         gameState.value.currentPlayer.id === currentUser.id
})

// 监听游戏状态变化
onMounted(() => {
  if (props.socket) {
    // 监听游戏状态更新
    props.socket.on('gameState', (state) => {
      gameState.value = state
      
      // 更新玩家可执行的操作
      updatePlayerActions()
    })
    
    // 监听发牌
    props.socket.on('dealCards', (cards) => {
      playerHand.value = cards
    })
    
    // 监听游戏结果
    props.socket.on('gameResult', (result) => {
      gameState.value.status = 'ended'
      gameState.value.result = result
      
      // 显示游戏结果
      ElMessage({
        message: `游戏结束: ${result.message}`,
        type: 'success',
        duration: 5000,
        showClose: true
      })
      
      // 3秒后显示详细结果
      setTimeout(() => {
        ElMessage({
          message: `获胜者: ${result.winner?.username || '无'}, 牌型: ${getHandRankName(result.winningHand?.rank)}`,
          type: 'success',
          duration: 5000,
          showClose: true
        })
      }, 3000)
    })
    
    // 监听错误消息
    props.socket.on('error', (error) => {
      ElMessage.error(error.message)
    })
    
    // 监听玩家准备状态变化
    props.socket.on('playerReady', (data) => {
      // 更新玩家准备状态
      const playerIndex = gameState.value.players.findIndex(p => p.id === data.playerId)
      if (playerIndex !== -1) {
        gameState.value.players[playerIndex].ready = data.ready
      }
    })
  }
})

onUnmounted(() => {
  if (props.socket) {
    props.socket.off('gameState')
    props.socket.off('dealCards')
    props.socket.off('gameResult')
    props.socket.off('error')
  }
})

// 更新玩家可执行的操作
const updatePlayerActions = () => {
  if (!isPlayerTurn.value || gameState.value.status !== 'playing') {
    playerActions.value = {
      canCheck: false,
      canCall: false,
      canRaise: false,
      canFold: false
    }
    return
  }
  
  const player = gameState.value.players.find(p => p.id === currentUser.id)
  if (!player) return
  
  const playerBet = player.bet || 0
  const currentBet = gameState.value.currentBet || 0
  
  playerActions.value = {
    canCheck: playerBet === currentBet,
    canCall: playerBet < currentBet,
    canRaise: true,
    canFold: true
  }
}

// 玩家操作函数
const check = () => {
  if (!isPlayerTurn.value || !playerActions.value.canCheck) return
  props.socket.emit('playerAction', { action: 'check', roomId: props.roomId })
  // 显示操作提示
  ElMessage.success('看牌')
}

const call = () => {
  if (!isPlayerTurn.value || !playerActions.value.canCall) return
  props.socket.emit('playerAction', { action: 'call', roomId: props.roomId })
  // 显示操作提示
  ElMessage.success('跟注')
}

const raise = () => {
  if (!isPlayerTurn.value || !playerActions.value.canRaise || betAmount.value < minBet.value) return
  props.socket.emit('playerAction', { 
    action: 'raise', 
    amount: betAmount.value,
    roomId: props.roomId 
  })
  // 显示操作提示
  ElMessage.success(`加注 ${betAmount.value}`)
}

const fold = () => {
  if (!isPlayerTurn.value || !playerActions.value.canFold) return
  props.socket.emit('playerAction', { action: 'fold', roomId: props.roomId })
  // 显示操作提示
  ElMessage.success('弃牌')
}

// 获取回合名称的中文显示
const getRoundName = (round) => {
  const roundNames = {
    'preflop': '前翻牌',
    'flop': '翻牌',
    'turn': '转牌',
    'river': '河牌',
    'showdown': '摊牌'
  }
  return roundNames[round] || round
}

// 获取牌型名称
const getHandRankName = (rank) => {
  if (!rank) return '无效牌型'
  
  const handRanks = {
    'high_card': '高牌',
    'pair': '一对',
    'two_pair': '两对',
    'three_of_a_kind': '三条',
    'straight': '顺子',
    'flush': '同花',
    'full_house': '葫芦',
    'four_of_a_kind': '四条',
    'straight_flush': '同花顺',
    'royal_flush': '皇家同花顺'
  }
  
  return handRanks[rank] || rank
}

// 动态生成卡牌SVG
const getCardImage = (card) => {
  if (!card) return '/cards/back.png'
  
  // 生成卡牌SVG的函数
  const generateCardSVG = (rank, suit) => {
    // 卡牌颜色
    const isRed = suit === 'hearts' || suit === 'diamonds'
    const color = isRed ? '#e53935' : '#212121'
    
    // 卡牌符号
    const suitSymbols = {
      'hearts': '♥',
      'diamonds': '♦',
      'clubs': '♣',
      'spades': '♠'
    }
    
    // 卡牌数值
    const rankDisplay = {
      'ace': 'A',
      'king': 'K',
      'queen': 'Q',
      'jack': 'J',
      '10': '10',
      '9': '9',
      '8': '8',
      '7': '7',
      '6': '6',
      '5': '5',
      '4': '4',
      '3': '3',
      '2': '2'
    }
    
    // 生成SVG
    return `data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="169" height="245" viewBox="0 0 169 245">
      <rect width="169" height="245" rx="12" fill="white" stroke="#cccccc" stroke-width="1" />
      <text x="20" y="40" font-family="Arial" font-size="30" font-weight="bold" fill="${color}">${rankDisplay[rank]}</text>
      <text x="20" y="70" font-family="Arial" font-size="30" fill="${color}">${suitSymbols[suit]}</text>
      <text x="84.5" y="130" font-family="Arial" font-size="60" text-anchor="middle" fill="${color}">${suitSymbols[suit]}</text>
      <text x="149" y="205" font-family="Arial" font-size="30" font-weight="bold" fill="${color}" text-anchor="end">${rankDisplay[rank]}</text>
      <text x="149" y="175" font-family="Arial" font-size="30" fill="${color}" text-anchor="end">${suitSymbols[suit]}</text>
    </svg>`
  }
  
  return generateCardSVG(card.rank, card.suit)
}

// 获取玩家位置样式
const getPlayerPosition = (index, totalPlayers) => {
  const positions = [
    { top: '75%', left: '50%' }, // 底部中央（当前玩家）
    { top: '60%', left: '20%' }, // 左下
    { top: '25%', left: '20%' }, // 左上
    { top: '10%', left: '50%' }, // 顶部中央
    { top: '25%', left: '80%' }, // 右上
    { top: '60%', left: '80%' }  // 右下
  ]
  
  // 确保当前玩家始终在底部中央
  const currentPlayerIndex = gameState.value.players.findIndex(p => p.id === currentUser.id)
  if (currentPlayerIndex === -1) return positions[index % positions.length]
  
  // 计算相对位置
  const relativeIndex = (index - currentPlayerIndex + totalPlayers) % totalPlayers
  return positions[relativeIndex % positions.length]
}
</script>

<template>
  <div class="poker-game">
    <!-- 游戏状态 -->
    <div class="game-status">
      <h3>{{ gameState.status === 'waiting' ? '等待游戏开始' : 
             gameState.status === 'playing' ? `当前回合: ${getRoundName(gameState.round)}` : '游戏结束' }}</h3>
      <div v-if="gameState.status === 'playing'" class="pot">
        底池: {{ gameState.pot }}
      </div>
      <div v-if="gameState.status === 'playing' && gameState.currentBet > 0" class="current-bet">
        当前下注: {{ gameState.currentBet }}
      </div>
      <div v-if="gameState.status === 'playing' && isPlayerTurn" class="turn-indicator">
        轮到你行动了!
      </div>
      <!-- 游戏结果显示 -->
      <div v-if="gameState.status === 'ended' && gameState.result" class="game-result">
        <h4>游戏结果</h4>
        <div class="winner-info">
          <span class="winner-label">获胜者:</span>
          <span class="winner-name">{{ gameState.result.winner?.username || '平局' }}</span>
        </div>
        <div class="hand-info" v-if="gameState.result.winningHand">
          <span class="hand-label">获胜牌型:</span>
          <span class="hand-name">{{ getHandRankName(gameState.result.winningHand.rank) }}</span>
        </div>
      </div>
    </div>
    
    <!-- 公共牌 -->
    <div class="community-cards" v-if="gameState.status === 'playing'">
      <div v-for="(card, index) in gameState.communityCards" :key="index" class="card">
        <img :src="getCardImage(card)" :alt="card ? `${card.rank} of ${card.suit}` : 'Card back'">
      </div>
    </div>
    
    <!-- 玩家区域 -->
    <div class="players-container">
      <div v-for="(player, index) in gameState.players" :key="player.id"
           class="player"
           :class="{
             'current-player': player.id === currentUser.id,
             'active-player': gameState.currentPlayer && gameState.currentPlayer.id === player.id
           }"
           :style="getPlayerPosition(index, gameState.players.length)">
        <div class="player-avatar">
          <el-avatar :size="50">
            {{ player.username.charAt(0).toUpperCase() }}
          </el-avatar>
        </div>
        <div class="player-info">
          <div class="player-name">{{ player.username }}</div>
          <div class="player-chips">筹码: {{ player.chips || 0 }}</div>
          <div v-if="player.bet" class="player-bet">下注: {{ player.bet }}</div>
          <div v-if="player.folded" class="player-folded">已弃牌</div>
          <div v-if="gameState.dealerPosition === index" class="dealer-button">D</div>
          <div v-if="gameState.status === 'waiting' && player.ready" class="player-ready">已准备</div>
        </div>
        
        <!-- 当前玩家的手牌 -->
        <div v-if="player.id === currentUser.id" class="player-hand">
          <div v-for="(card, cardIndex) in playerHand" :key="cardIndex" class="card">
            <img :src="getCardImage(card)" :alt="`${card.rank} of ${card.suit}`">
          </div>
        </div>
        
        <!-- 其他玩家的手牌（背面） -->
        <div v-else class="player-hand">
          <div v-for="i in (player.handSize || 0)" :key="i" class="card">
            <img src="/cards/back.png" alt="Card back">
          </div>
        </div>
      </div>
    </div>
    
    <!-- 玩家操作区域 -->
    <div class="player-actions" v-if="gameState.status === 'playing' && isPlayerTurn">
      <el-button :disabled="!playerActions.canCheck" @click="check">看牌</el-button>
      <el-button :disabled="!playerActions.canCall" @click="call">跟注</el-button>
      <el-button :disabled="!playerActions.canRaise" @click="raise">加注</el-button>
      <el-button :disabled="!playerActions.canFold" @click="fold">弃牌</el-button>
      
      <div class="bet-slider" v-if="playerActions.canRaise">
        <el-slider v-model="betAmount" :min="minBet" :max="100" :step="5" show-input></el-slider>
      </div>
    </div>
  </div>
</template>

<style scoped>
.poker-game {
  position: relative;
  width: 100%;
  height: 100%;
  background-color: #1b5e20;
  background-image: radial-gradient(#2e7d32 15%, #1b5e20 60%);
  border-radius: 8px;
  overflow: hidden;
  color: white;
  box-shadow: inset 0 0 60px rgba(0, 0, 0, 0.5);
}

.game-status {
  text-align: center;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.5);
  border-bottom: 2px solid rgba(255, 255, 255, 0.1);
}

.pot {
  font-size: 18px;
  font-weight: bold;
  margin-top: 5px;
  text-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
}

.current-bet {
  font-size: 16px;
  margin-top: 5px;
  color: #ffeb3b;
}

.turn-indicator {
  font-size: 18px;
  font-weight: bold;
  margin-top: 10px;
  color: #ff9800;
  animation: pulse 1.5s infinite;
}

.game-result {
  margin-top: 15px;
  padding: 10px 15px;
  background-color: rgba(0, 0, 0, 0.6);
  border-radius: 8px;
  border: 1px solid rgba(255, 215, 0, 0.5);
  box-shadow: 0 0 10px rgba(255, 215, 0, 0.3);
  animation: fadeIn 0.5s ease-in-out;
}

.game-result h4 {
  color: #ffd700;
  margin-top: 0;
  margin-bottom: 10px;
  text-align: center;
  font-size: 18px;
}

.winner-info, .hand-info {
  margin: 5px 0;
  display: flex;
  justify-content: space-between;
}

.winner-label, .hand-label {
  color: #e0e0e0;
}

.winner-name, .hand-name {
  color: #4caf50;
  font-weight: bold;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes pulse {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}

.community-cards {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin: 20px 0;
  perspective: 1000px;
}

.card {
  width: 80px;
  height: 120px;
  border-radius: 5px;
  background-color: white;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.5);
  transition: all 0.3s ease;
  transform-style: preserve-3d;
}

.card:hover {
  transform: translateY(-5px) rotateY(5deg);
  box-shadow: 5px 8px 12px rgba(0, 0, 0, 0.5);
}

.card img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 5px;
}

.players-container {
  position: relative;
  width: 100%;
  height: 400px;
}

.player {
  position: absolute;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: all 0.3s ease;
  filter: drop-shadow(0 0 5px rgba(0, 0, 0, 0.5));
}

.player-avatar {
  margin-bottom: 5px;
  border: 2px solid rgba(255, 255, 255, 0.5);
  border-radius: 50%;
  transition: all 0.3s ease;
}

.player-info {
  text-align: center;
  background-color: rgba(0, 0, 0, 0.7);
  padding: 8px 12px;
  border-radius: 8px;
  margin-bottom: 10px;
  min-width: 120px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.player-hand {
  display: flex;
  gap: 5px;
  transition: transform 0.3s ease;
}

.player-hand .card {
  width: 60px;
  height: 90px;
  transform: rotate(0);
  transition: transform 0.3s ease;
}

.player-hand:hover .card:first-child {
  transform: rotate(-5deg) translateX(-5px);
}

.player-hand:hover .card:last-child {
  transform: rotate(5deg) translateX(5px);
}

.current-player {
  z-index: 10;
}

.active-player .player-avatar {
  box-shadow: 0 0 15px 5px gold;
  border-radius: 50%;
  animation: glow 1.5s infinite alternate;
}

@keyframes glow {
  from { box-shadow: 0 0 10px 2px gold; }
  to { box-shadow: 0 0 20px 8px gold; }
}

.player-folded {
  color: #ff5252;
  font-weight: bold;
}

.player-ready {
  color: #4caf50;
  font-weight: bold;
  margin-top: 5px;
}

.dealer-button {
  position: absolute;
  top: -15px;
  right: -15px;
  width: 25px;
  height: 25px;
  background-color: white;
  color: black;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.5);
  border: 1px solid #ffd700;
}

.player-actions {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.7);
  padding: 15px;
  border-radius: 12px;
  width: 80%;
  max-width: 600px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.player-actions .el-button {
  transition: all 0.2s ease;
  min-width: 80px;
}

.player-actions .el-button:not(:disabled):hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

.bet-slider {
  width: 100%;
  margin-top: 15px;
  padding: 5px 10px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}
</style>