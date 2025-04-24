package services

import (
	"encoding/json"
	"fmt"
	"games/dao"
	"games/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WSService WebSocket服务
type WSService struct {
	// 连接管理，外层map是房间ID，内层map是用户ID对应的连接
	connections  map[uint]map[uint]*websocket.Conn
	mu           sync.RWMutex // 保护上述map的互斥锁
	roomService  *RoomService
	pokerService *PokerService
	roomDAO      *dao.RoomDAO
	userDAO      *dao.UserDAO
	wsUpgrader   websocket.Upgrader
}

// ClientMessage 客户端消息
type ClientMessage struct {
	Type   string          `json:"type"`
	RoomID uint            `json:"room_id"`
	Data   json.RawMessage `json:"data"`
}

// ToggleReadyMessage 玩家准备消息
type ToggleReadyMessage struct {
	RoomID uint `json:"room_id"`
}

// PlayerActionMessage 玩家操作消息
type PlayerActionMessage struct {
	Action string `json:"action"`
	Amount int    `json:"amount,omitempty"`
	RoomID uint   `json:"room_id"`
}

// NewWSService 创建WebSocket服务实例
func NewWSService(roomService *RoomService, pokerService *PokerService) *WSService {
	return &WSService{
		connections:  make(map[uint]map[uint]*websocket.Conn),
		roomService:  roomService,
		pokerService: pokerService,
		roomDAO:      dao.NewRoomDAO(),
		userDAO:      dao.NewUserDAO(),
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源的WebSocket连接
			},
		},
	}
}

// HandleConnection 处理WebSocket连接
func (s *WSService) HandleConnection(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		return
	}

	// 获取房间ID
	roomIDStr := c.Query("roomId")
	var roomID uint
	if roomIDStr != "" {
		if _, err := fmt.Sscanf(roomIDStr, "%d", &roomID); err != nil {
			return
		}
	} else {
		return
	}

	// 检查用户是否在房间中
	if !s.roomDAO.IsPlayerInRoom(roomID, userID.(uint)) {
		return
	}

	// 升级HTTP连接为WebSocket
	ws, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 记录连接
	s.mu.Lock()
	if _, ok := s.connections[roomID]; !ok {
		s.connections[roomID] = make(map[uint]*websocket.Conn)
	}
	s.connections[roomID][userID.(uint)] = ws
	s.mu.Unlock()

	// 向所有玩家广播房间更新信息
	s.BroadcastRoomUpdate(roomID)

	// 处理连接
	s.handleMessages(ws, userID.(uint), roomID)
}

// handleMessages 处理WebSocket消息
func (s *WSService) handleMessages(ws *websocket.Conn, userID, roomID uint) {
	defer func() {
		ws.Close()
		s.mu.Lock()
		delete(s.connections[roomID], userID)
		if len(s.connections[roomID]) == 0 {
			delete(s.connections, roomID)
		}
		s.mu.Unlock()
	}()

	for {
		// 读取消息
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// 解析消息
		var clientMsg ClientMessage
		err = json.Unmarshal(message, &clientMsg)
		if err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		// 处理不同类型的消息
		switch clientMsg.Type {
		case "toggleReady":
			var readyMsg ToggleReadyMessage
			err = json.Unmarshal(clientMsg.Data, &readyMsg)
			if err != nil {
				log.Println("JSON parse error:", err)
				continue
			}
			s.handleToggleReady(userID, readyMsg.RoomID)

		case "playerAction":
			var actionMsg PlayerActionMessage
			err = json.Unmarshal(clientMsg.Data, &actionMsg)
			if err != nil {
				log.Println("JSON parse error:", err)
				continue
			}
			s.handlePlayerAction(userID, actionMsg)

		default:
			log.Println("Unknown message type:", clientMsg.Type)
		}
	}
}

// handleToggleReady 处理玩家准备消息
func (s *WSService) handleToggleReady(userID, roomID uint) {
	// 切换玩家准备状态
	newStatus, err := s.roomService.ToggleReady(roomID, userID)
	if err != nil {
		s.SendErrorToUser(roomID, userID, err.Error())
		return
	}

	// 广播玩家准备状态
	s.BroadcastPlayerReady(roomID, userID, newStatus)

	// 检查是否所有玩家都已准备
	room, err := s.roomService.GetRoomByID(roomID)
	if err != nil {
		return
	}

	// 检查玩家数量和准备状态
	readyStatus, err := s.roomDAO.GetPlayerReadyStatus(roomID)
	if err != nil {
		return
	}

	allReady := true
	for _, player := range room.Players {
		if status, ok := readyStatus[player.ID]; !ok || !status {
			allReady = false
			break
		}
	}

	// 如果所有玩家都已准备，且玩家数量足够，开始游戏
	if allReady && len(room.Players) >= 2 {
		s.StartGame(roomID)
	}
}

// handlePlayerAction 处理玩家游戏操作
func (s *WSService) handlePlayerAction(userID uint, msg PlayerActionMessage) {
	// 根据操作类型处理
	err := s.pokerService.PlayerAction(msg.RoomID, userID, msg.Action, msg.Amount)
	if err != nil {
		s.SendErrorToUser(msg.RoomID, userID, err.Error())
		return
	}

	// 获取最新游戏状态并广播
	s.BroadcastGameState(msg.RoomID)
}

// StartGame 开始游戏
func (s *WSService) StartGame(roomID uint) {
	// 初始化游戏
	err := s.pokerService.InitGame(roomID)
	if err != nil {
		// 发送错误消息给房间所有玩家
		s.BroadcastError(roomID, err.Error())
		return
	}

	// 获取游戏状态
	gameState, err := s.pokerService.GetGameState(roomID)
	if err != nil {
		return
	}

	// 广播游戏状态
	s.BroadcastGameState(roomID)

	// 给每个玩家发送其手牌
	for _, player := range gameState.Players {
		playerState, err := s.pokerService.GetPlayerState(roomID, player.ID)
		if err == nil && len(playerState.HandCards) > 0 {
			s.SendCardsToPlayer(roomID, player.ID, playerState.HandCards)
		}
	}
}

// BroadcastRoomUpdate 广播房间更新
func (s *WSService) BroadcastRoomUpdate(roomID uint) {
	// 获取房间信息
	room, err := s.roomService.GetRoomByID(roomID)
	if err != nil {
		return
	}

	// 获取玩家准备状态
	readyStatus, err := s.roomDAO.GetPlayerReadyStatus(roomID)
	if err != nil {
		return
	}

	// 构建包含准备状态的玩家列表
	type PlayerWithReady struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Ready    bool   `json:"ready"`
	}

	roomUpdate := struct {
		ID         uint              `json:"id"`
		Name       string            `json:"name"`
		GameID     uint              `json:"game_id"`
		Players    []PlayerWithReady `json:"players"`
		Status     string            `json:"status"`
		MaxPlayers int               `json:"max_players"`
	}{
		ID:         room.ID,
		Name:       room.Name,
		GameID:     room.GameID,
		Players:    make([]PlayerWithReady, 0),
		Status:     room.Status,
		MaxPlayers: room.MaxPlayers,
	}

	for _, player := range room.Players {
		playerWithReady := PlayerWithReady{
			ID:       player.ID,
			Username: player.Username,
			Ready:    readyStatus[player.ID],
		}
		roomUpdate.Players = append(roomUpdate.Players, playerWithReady)
	}

	// 广播给房间中的所有玩家
	message, _ := json.Marshal(map[string]interface{}{
		"type": "roomUpdate",
		"data": roomUpdate,
	})

	s.BroadcastToRoom(roomID, message)
}

// BroadcastPlayerReady 广播玩家准备状态
func (s *WSService) BroadcastPlayerReady(roomID, playerID uint, ready bool) {
	// 获取玩家信息
	player, err := s.userDAO.GetByID(playerID)
	if err != nil {
		return
	}

	// 构建准备状态消息
	readyMessage, _ := json.Marshal(map[string]interface{}{
		"type": "playerReady",
		"data": map[string]interface{}{
			"playerId": playerID,
			"username": player.Username,
			"ready":    ready,
		},
	})

	// 广播给房间中的所有玩家
	s.BroadcastToRoom(roomID, readyMessage)
}

// BroadcastGameState 广播游戏状态
func (s *WSService) BroadcastGameState(roomID uint) {
	// 获取游戏状态
	gameState, err := s.pokerService.GetGameState(roomID)
	if err != nil {
		return
	}

	// 对于每个玩家，发送适合他们的游戏状态
	s.mu.RLock()
	defer s.mu.RUnlock()

	for playerID, conn := range s.connections[roomID] {
		// 创建一个游戏状态的副本，移除其他玩家的手牌
		stateCopy := *gameState
		for i := range stateCopy.Players {
			if stateCopy.Players[i].ID != playerID {
				// 对其他玩家，只保留手牌数量信息，不显示具体牌面
				stateCopy.Players[i].HandCards = nil
			}
		}

		// 发送游戏状态
		message, _ := json.Marshal(map[string]interface{}{
			"type": "gameState",
			"data": stateCopy,
		})

		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}
}

// SendCardsToPlayer 给玩家发送手牌
func (s *WSService) SendCardsToPlayer(roomID, playerID uint, cards []models.Card) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查玩家连接是否存在
	connections, ok := s.connections[roomID]
	if !ok {
		return
	}

	conn, ok := connections[playerID]
	if !ok {
		return
	}

	// 发送手牌
	message, _ := json.Marshal(map[string]interface{}{
		"type": "dealCards",
		"data": cards,
	})

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("WebSocket write error:", err)
	}
}

// BroadcastToRoom 广播消息给房间中的所有玩家
func (s *WSService) BroadcastToRoom(roomID uint, message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查房间是否存在
	connections, ok := s.connections[roomID]
	if !ok {
		return
	}

	// 向房间中所有玩家发送消息
	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}
}

// SendErrorToUser 发送错误消息给指定用户
func (s *WSService) SendErrorToUser(roomID, userID uint, errorMsg string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查玩家连接是否存在
	connections, ok := s.connections[roomID]
	if !ok {
		return
	}

	conn, ok := connections[userID]
	if !ok {
		return
	}

	// 发送错误消息
	message, _ := json.Marshal(map[string]interface{}{
		"type": "error",
		"data": map[string]string{
			"message": errorMsg,
		},
	})

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("WebSocket write error:", err)
	}
}

// BroadcastError 广播错误消息给房间中的所有玩家
func (s *WSService) BroadcastError(roomID uint, errorMsg string) {
	// 构建错误消息
	message, _ := json.Marshal(map[string]interface{}{
		"type": "error",
		"data": map[string]string{
			"message": errorMsg,
		},
	})

	// 广播给房间中的所有玩家
	s.BroadcastToRoom(roomID, message)
}
