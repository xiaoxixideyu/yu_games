# 游戏平台API接口文档

## 概述

本文档描述了游戏平台前端与后端交互的所有API接口。所有API请求都使用统一的响应格式，便于前端处理。

## 统一响应格式

所有API响应都遵循以下统一格式：

```json
{
  "code": 0,       // 状态码，0表示成功，非0表示失败
  "data": {},     // 响应数据，具体格式根据接口不同而变化
  "msg": ""       // 响应消息，成功时为空或成功提示，失败时为错误信息
}
```

前端通过axios拦截器统一处理这种响应格式，提取`data`字段中的数据：

```javascript
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
```

## 用户认证接口

### 用户登录

- **URL**: `/api/login`
- **方法**: POST
- **请求参数**:

```json
{
  "username": "用户名",
  "password": "密码"
}
```

- **响应数据**:

```json
{
  "code": 0,
  "data": {
    "token": "认证令牌",
    "user": {
      "id": "用户ID",
      "username": "用户名",
      // 其他用户信息字段
    }
  },
  "msg": "登录成功"
}
```

- **前端调用示例**:

```javascript
async login(username, password) {
  try {
    const response = await axios.post('/api/login', {
      username,
      password
    })
    // 响应拦截器已处理统一格式，直接获取data中的数据
    const { token, user } = response
    // 存储token和用户信息
    localStorage.setItem('token', token)
    localStorage.setItem('userInfo', JSON.stringify(user))
    return true
  } catch (error) {
    console.error('Login failed:', error)
    return false
  }
}
```

### 用户注册

- **URL**: `/api/register`
- **方法**: POST
- **请求参数**:

```json
{
  "username": "用户名",
  "password": "密码"
}
```

- **响应数据**:

```json
{
  "code": 0,
  "data": null,
  "msg": "注册成功"
}
```

- **前端调用示例**:

```javascript
async register(username, password) {
  try {
    await axios.post('/api/register', {
      username,
      password
    })
    return true
  } catch (error) {
    console.error('Registration failed:', error)
    return false
  }
}
```

## 通用游戏接口

### 获取游戏列表

- **URL**: `/api/games`
- **方法**: GET
- **请求参数**: 无
- **响应数据**:

```json
{
  "code": 0,
  "data": [
    {
      "id": "游戏ID",
      "name": "游戏名称",
      "description": "游戏描述",
      "image": "游戏图片URL",
      // 其他游戏信息字段
    },
    // 更多游戏...
  ],
  "msg": ""
}
```

- **前端调用示例**:

```javascript
async fetchGameList() {
  try {
    const data = await axios.get('/api/games')
    // 响应拦截器已处理统一格式，直接获取data中的数据
    this.gameList = data
    return true
  } catch (error) {
    console.error('Failed to fetch game list:', error)
    return false
  }
}
```

### 获取游戏房间列表

- **URL**: `/api/games/{gameId}/rooms`
- **方法**: GET
- **请求参数**: 路径参数 `gameId`
- **响应数据**:

```json
{
  "code": 0,
  "data": [
    {
      "id": "房间ID",
      "name": "房间名称",
      "gameId": "游戏ID",
      "playerCount": "当前玩家数",
      "maxPlayers": "最大玩家数",
      "status": "房间状态",
      // 其他房间信息字段
    },
    // 更多房间...
  ],
  "msg": ""
}
```

- **前端调用示例**:

```javascript
async fetchGameLobby(gameId) {
  try {
    const data = await axios.get(`/api/games/${gameId}/rooms`)
    // 响应拦截器已处理统一格式，直接获取data中的数据
    this.lobbyRooms = data
    this.currentGame = this.gameList.find(game => game.id === gameId)
    return true
  } catch (error) {
    console.error('Failed to fetch game lobby:', error)
    return false
  }
}
```

### 创建游戏房间

- **URL**: `/api/games/{gameId}/rooms`
- **方法**: POST
- **请求参数**:
  - 路径参数 `gameId`
  - 请求体:

```json
{
  "name": "房间名称"
}
```

- **响应数据**:

```json
{
  "code": 0,
  "data": {
    "id": "房间ID",
    "name": "房间名称",
    "gameId": "游戏ID",
    "creator": {
      "id": "创建者ID",
      "username": "创建者用户名"
    },
    "players": [
      // 当前房间内的玩家列表
    ],
    "status": "房间状态",
    // 其他房间信息字段
  },
  "msg": "创建成功"
}
```

- **前端调用示例**:

```javascript
async createRoom(gameId, roomName) {
  try {
    const data = await axios.post(`/api/games/${gameId}/rooms`, {
      name: roomName
    })
    // 响应拦截器已处理统一格式，直接获取data中的数据
    this.currentRoom = data
    return true
  } catch (error) {
    console.error('Failed to create room:', error)
    return false
  }
}
```

### 加入游戏房间

- **URL**: `/api/rooms/{roomId}/join`
- **方法**: POST
- **请求参数**: 路径参数 `roomId`
- **响应数据**:

```json
{
  "code": 0,
  "data": {
    "id": "房间ID",
    "name": "房间名称",
    "gameId": "游戏ID",
    "players": [
      // 当前房间内的玩家列表，包括新加入的玩家
    ],
    "status": "房间状态",
    // 其他房间信息字段
  },
  "msg": "加入成功"
}
```

- **前端调用示例**:

```javascript
async joinRoom(roomId) {
  try {
    const data = await axios.post(`/api/rooms/${roomId}/join`)
    // 响应拦截器已处理统一格式，直接获取data中的数据
    this.currentRoom = data
    return true
  } catch (error) {
    console.error('Failed to join room:', error)
    return false
  }
}
```

## 通用WebSocket接口

所有游戏都使用WebSocket进行实时通信的基本功能，下面是所有游戏通用的WebSocket事件说明。

### 连接建立

- **连接URL**: `{VITE_API_BASE_URL}`
- **连接参数**: 
  - `roomId`: 房间ID

- **前端连接示例**:

```javascript
// 连接到WebSocket服务器
socket.value = io(import.meta.env.VITE_API_BASE_URL, {
  query: {
    roomId: route.params.roomId
  }
})
```

### 通用服务器发送的事件

#### 房间更新事件

- **事件名**: `roomUpdate`
- **数据格式**:

```json
{
  "id": "房间ID",
  "name": "房间名称",
  "gameId": "游戏ID",
  "players": [
    {
      "id": "玩家ID",
      "username": "玩家用户名",
      "ready": true/false,  // 玩家准备状态
      "avatar": "头像URL"   // 可选
    },
    // 更多玩家...
  ],
  "status": "房间状态",  // waiting, playing, ended
  "maxPlayers": 6       // 最大玩家数
}
```

- **前端处理示例**:

```javascript
socket.value.on('roomUpdate', (room) => {
  gameStore.currentRoom = room
  updateReadyStatus()
})
```

#### 玩家准备状态变化事件

- **事件名**: `playerReady`
- **数据格式**:

```json
{
  "playerId": "玩家ID",
  "ready": true/false     // 准备状态
}
```

- **前端处理示例**:

```javascript
socket.value.on('playerReady', (data) => {
  // 更新玩家准备状态
  const playerIndex = gameState.value.players.findIndex(p => p.id === data.playerId)
  if (playerIndex !== -1) {
    gameState.value.players[playerIndex].ready = data.ready
  }
})
```

#### 错误消息事件

- **事件名**: `error`
- **数据格式**:

```json
{
  "message": "错误信息"
}
```

- **前端处理示例**:

```javascript
socket.value.on('error', (error) => {
  ElMessage.error(error.message)
})
```

### 通用客户端发送的事件

#### 玩家准备状态切换事件

- **事件名**: `toggleReady`
- **数据格式**:

```json
{
  "roomId": "房间ID"      // 房间标识
}
```

- **前端发送示例**:

```javascript
const toggleReady = () => {
  if (!gameStore.currentRoom) return
  
  socket.value.emit('toggleReady', { 
    roomId: gameStore.currentRoom.id 
  })
}
```

## 德州扑克游戏接口

德州扑克游戏除了使用通用的游戏接口和WebSocket接口外，还有以下特定的WebSocket事件。

### 德州扑克服务器发送的事件

#### 游戏状态更新事件

- **事件名**: `gameState`
- **数据格式**:

```json
{
  "status": "游戏状态",       // waiting, playing, ended
  "players": [
    {
      "id": "玩家ID",
      "username": "玩家用户名",
      "chips": 1000,          // 玩家筹码
      "bet": 10,             // 当前下注金额
      "folded": false,        // 是否弃牌
      "handSize": 2           // 手牌数量（对其他玩家隐藏具体牌面）
    },
    // 更多玩家...
  ],
  "currentPlayer": {         // 当前行动的玩家
    "id": "玩家ID",
    "username": "玩家用户名"
  },
  "communityCards": [        // 公共牌
    {
      "rank": "ace",         // 牌面点数
      "suit": "hearts"       // 牌面花色
    },
    // 更多公共牌...
  ],
  "pot": 100,               // 底池金额
  "currentBet": 20,         // 当前回合最高下注额
  "dealerPosition": 0,       // 庄家位置
  "round": "preflop"         // 当前回合: preflop, flop, turn, river, showdown
}
```

- **前端处理示例**:

```javascript
socket.value.on('gameState', (state) => {
  gameState.value = state
  
  // 更新玩家可执行的操作
  updatePlayerActions()
})
```

#### 发牌事件

- **事件名**: `dealCards`
- **数据格式**:

```json
[
  {
    "rank": "ace",         // 牌面点数
    "suit": "hearts"       // 牌面花色
  },
  {
    "rank": "king",
    "suit": "spades"
  }
]
```

- **前端处理示例**:

```javascript
socket.value.on('dealCards', (cards) => {
  playerHand.value = cards
})
```

#### 游戏结果事件

- **事件名**: `gameResult`
- **数据格式**:

```json
{
  "message": "游戏结果描述",
  "winner": {
    "id": "获胜者ID",
    "username": "获胜者用户名"
  },
  "winningHand": {
    "rank": "straight_flush",  // 获胜牌型
    "cards": [                 // 组成获胜牌型的牌
      {
        "rank": "ace",
        "suit": "hearts"
      },
      // 更多牌...
    ]
  },
  "players": [                 // 所有玩家的最终结果
    {
      "id": "玩家ID",
      "username": "玩家用户名",
      "hand": [                // 玩家手牌
        {
          "rank": "ace",
          "suit": "hearts"
        },
        // 更多牌...
      ],
      "handRank": "straight_flush", // 玩家牌型
      "winnings": 200             // 玩家赢得的筹码
    },
    // 更多玩家...
  ]
}
```

- **前端处理示例**:

```javascript
socket.value.on('gameResult', (result) => {
  gameState.value.status = 'ended'
  gameState.value.result = result
  
  // 显示游戏结果
  ElMessage({
    message: `游戏结束: ${result.message}`,
    type: 'success',
    duration: 5000,
    showClose: true
  })
})
```

### 德州扑克客户端发送的事件

#### 玩家操作事件

- **事件名**: `playerAction`
- **数据格式**:

```json
{
  "action": "操作类型",    // check, call, raise, fold
  "amount": 50,          // 仅在raise操作时需要
  "roomId": "房间ID"      // 房间标识
}
```

- **前端发送示例**:

```javascript
// 看牌操作
const check = () => {
  if (!isPlayerTurn.value || !playerActions.value.canCheck) return
  props.socket.emit('playerAction', { action: 'check', roomId: props.roomId })
}

// 跟注操作
const call = () => {
  if (!isPlayerTurn.value || !playerActions.value.canCall) return
  props.socket.emit('playerAction', { action: 'call', roomId: props.roomId })
}

// 加注操作
const raise = () => {
  if (!isPlayerTurn.value || !playerActions.value.canRaise || betAmount.value < minBet.value) return
  props.socket.emit('playerAction', { 
    action: 'raise', 
    amount: betAmount.value,
    roomId: props.roomId 
  })
}

// 弃牌操作
const fold = () => {
  if (!isPlayerTurn.value || !playerActions.value.canFold) return
  props.socket.emit('playerAction', { action: 'fold', roomId: props.roomId })
}
```

## 错误处理

当API请求失败时，响应格式如下：

```json
{
  "code": 错误码,  // 非0的错误码
  "data": null,   // 通常为null
  "msg": "错误信息" // 具体的错误描述
}
```

前端通过axios拦截器统一处理错误响应，当`code`不为0时，会在控制台输出错误信息，并通过Promise.reject返回错误，便于调用方进行错误处理。

## 认证机制

系统使用基于token的认证机制：

1. 用户登录成功后，服务器返回token
2. 前端将token保存在localStorage中
3. 前端在后续请求中通过请求头携带token
4. 路由守卫检查localStorage中是否存在token，不存在则重定向到登录页面

### 建议实现的请求拦截器

当前项目中未实现请求拦截器来自动添加认证头，建议添加以下代码到main.js中：

```javascript
// 配置axios请求拦截器添加认证头
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)
```

## 接口状态码说明

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 用户未登录或token无效 |
| 1002 | 用户名或密码错误 |
| 1003 | 用户名已存在 |
| 2001 | 游戏不存在 |
| 2002 | 房间不存在 |
| 2003 | 房间已满 |
| 2004 | 已在房间中 |
| 9999 | 服务器内部错误 |

## 环境配置

### 开发环境

开发环境API基础路径为相对路径，通过Vite的代理功能转发到后端服务：

```javascript
// vite.config.js
export default defineConfig({
  // ...
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  }
})
```

### 生产环境

生产环境可以通过环境变量配置API基础路径：

```javascript
// 在.env.production文件中
VITE_API_BASE_URL=https://api.yourgameserver.com

// 在main.js中配置axios基础路径
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || ''
```

## 注意事项

1. 所有API请求都需要在登录状态下进行，除了登录和注册接口
2. 前端需要处理token过期的情况，可能需要在响应拦截器中检测特定的错误码（如1001），并引导用户重新登录
3. WebSocket连接需要在用户加入房间后建立，并在离开房间时断开连接
4. 游戏状态更新和玩家操作主要通过WebSocket事件进行，确保实时性
5. 玩家操作需要在服务器端进行验证，防止客户端作弊
6. 游戏结果需要在服务器端计算，客户端只负责展示
7. 在处理WebSocket事件时，需要注意事件的顺序和依赖关系，确保游戏状态的一致性