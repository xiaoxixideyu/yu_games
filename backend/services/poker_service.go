package services

import (
	"errors"
	"games/models"
	"math/rand"
	"sync"
	"time"
)

// PokerService 德州扑克游戏服务
type PokerService struct {
	// 游戏状态缓存，使用房间ID作为键
	gameStates     map[uint]*models.PokerGameState
	handCards      map[uint]map[uint][]models.Card // 玩家手牌，外层map是房间ID，内层map是玩家ID
	mu             sync.RWMutex                    // 保护上述map的互斥锁
	roomService    *RoomService
	gameEventsChan map[uint]chan GameEvent          // 每个房间对应一个事件通道
	gameResults    map[uint]*models.PokerGameResult // 游戏结果
}

// GameEvent 游戏事件接口
type GameEvent interface {
	GetType() string
}

// PlayerActionEvent 玩家操作事件
type PlayerActionEvent struct {
	PlayerID uint
	Action   string
	Amount   int
}

func (e PlayerActionEvent) GetType() string {
	return "PlayerAction"
}

// NewPokerService 创建德州扑克游戏服务实例
func NewPokerService(roomService *RoomService) *PokerService {
	return &PokerService{
		gameStates:     make(map[uint]*models.PokerGameState),
		handCards:      make(map[uint]map[uint][]models.Card),
		roomService:    roomService,
		gameEventsChan: make(map[uint]chan GameEvent),
		gameResults:    make(map[uint]*models.PokerGameResult),
	}
}

// InitGame 初始化游戏
func (s *PokerService) InitGame(roomID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查房间是否存在
	room, err := s.roomService.GetRoomByID(roomID)
	if err != nil {
		return errors.New("房间不存在")
	}

	// 检查房间状态
	if room.Status != "waiting" {
		return errors.New("房间状态不正确")
	}

	// 检查玩家数量
	if len(room.Players) < 2 {
		return errors.New("玩家数量不足")
	}

	// 创建新的游戏状态
	gameState := &models.PokerGameState{
		Status:         "playing",
		Players:        make([]models.PokerPlayerState, 0),
		CommunityCards: make([]models.Card, 0),
		Pot:            0,
		CurrentBet:     0,
		DealerPosition: rand.Intn(len(room.Players)),
		Round:          models.RoundPreflop,
	}

	// 设置玩家状态
	s.handCards[roomID] = make(map[uint][]models.Card)
	initialChips := 1000 // 初始筹码
	for _, player := range room.Players {
		playerState := models.PokerPlayerState{
			ID:        player.ID,
			Username:  player.Username,
			Chips:     initialChips,
			Bet:       0,
			Folded:    false,
			IsAllIn:   false,
			HandSize:  2,
			HandCards: nil,
		}
		gameState.Players = append(gameState.Players, playerState)
	}

	// 保存游戏状态
	s.gameStates[roomID] = gameState

	// 更新房间状态
	err = s.roomService.roomDAO.UpdateRoomStatus(roomID, "playing")
	if err != nil {
		return err
	}

	// 创建事件通道
	s.gameEventsChan[roomID] = make(chan GameEvent, 100)

	// 启动游戏循环
	go s.gameLoop(roomID)

	return nil
}

// DealCards 发牌
func (s *PokerService) DealCards(roomID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查游戏状态
	gameState, ok := s.gameStates[roomID]
	if !ok {
		return errors.New("游戏未初始化")
	}

	// 创建一副牌
	var deck []models.Card
	suits := []models.CardSuit{models.SuitHearts, models.SuitDiamonds, models.SuitClubs, models.SuitSpades}
	ranks := []models.CardRank{
		models.RankAce, models.RankTwo, models.RankThree, models.RankFour, models.RankFive,
		models.RankSix, models.RankSeven, models.RankEight, models.RankNine, models.RankTen,
		models.RankJack, models.RankQueen, models.RankKing,
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, models.Card{Rank: rank, Suit: suit})
		}
	}

	// 洗牌
	rand.Seed(time.Now().UnixNano())
	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	// 给每个玩家发两张牌
	for i := range gameState.Players {
		if gameState.Players[i].Folded {
			continue
		}

		playerID := gameState.Players[i].ID
		s.handCards[roomID][playerID] = []models.Card{deck[0], deck[1]}
		deck = deck[2:] // 从牌堆中移除这两张牌
	}

	// 保存社区牌（5张）
	gameState.CommunityCards = []models.Card{} // 清空现有的社区牌
	communityCards := deck[:5]
	deck = deck[5:] // 从牌堆中移除这五张牌

	// 根据游戏轮次显示部分或全部社区牌
	if gameState.Round == models.RoundFlop {
		gameState.CommunityCards = communityCards[:3]
	} else if gameState.Round == models.RoundTurn {
		gameState.CommunityCards = communityCards[:4]
	} else if gameState.Round == models.RoundRiver || gameState.Round == models.RoundShowdown {
		gameState.CommunityCards = communityCards
	}

	return nil
}

// GetPlayerState 获取玩家状态
func (s *PokerService) GetPlayerState(roomID, playerID uint) (*models.PokerPlayerState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查游戏状态
	gameState, ok := s.gameStates[roomID]
	if !ok {
		return nil, errors.New("游戏未初始化")
	}

	// 查找玩家
	for i, player := range gameState.Players {
		if player.ID == playerID {
			// 创建玩家状态的副本
			state := gameState.Players[i]

			// 添加玩家手牌
			if cards, ok := s.handCards[roomID][playerID]; ok {
				state.HandCards = cards
			}

			return &state, nil
		}
	}

	return nil, errors.New("玩家不在游戏中")
}

// GetGameState 获取游戏状态
func (s *PokerService) GetGameState(roomID uint) (*models.PokerGameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查游戏状态
	gameState, ok := s.gameStates[roomID]
	if !ok {
		return nil, errors.New("游戏未初始化")
	}

	// 返回游戏状态的副本，确保不会被外部修改
	stateCopy := *gameState
	return &stateCopy, nil
}

// GetGameResult 获取游戏结果
func (s *PokerService) GetGameResult(roomID uint) (*models.PokerGameResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.gameResults[roomID]
	if !ok {
		return nil, errors.New("游戏结果不存在")
	}

	return result, nil
}

// PlayerAction 处理玩家操作
func (s *PokerService) PlayerAction(roomID, playerID uint, action string, amount int) error {
	// 检查游戏状态
	gameState, ok := s.gameStates[roomID]
	if !ok {
		return errors.New("游戏未初始化")
	}

	// 检查是否为当前玩家的回合
	if gameState.CurrentPlayer == nil || gameState.CurrentPlayer.ID != playerID {
		return errors.New("不是你的回合")
	}

	// 创建玩家操作事件
	event := PlayerActionEvent{
		PlayerID: playerID,
		Action:   action,
		Amount:   amount,
	}

	// 将事件发送到游戏循环
	s.gameEventsChan[roomID] <- event

	return nil
}

// gameLoop 游戏循环
func (s *PokerService) gameLoop(roomID uint) {
	// 获取事件通道
	eventChan, ok := s.gameEventsChan[roomID]
	if !ok {
		return
	}

	// 初始化游戏
	err := s.DealCards(roomID)
	if err != nil {
		return
	}

	// 开始第一轮操作
	s.startRound(roomID)

	// 游戏主循环
	for {
		select {
		case event := <-eventChan:
			switch e := event.(type) {
			case PlayerActionEvent:
				s.handlePlayerAction(roomID, e)
			}
		case <-time.After(30 * time.Minute): // 超时保护，防止协程泄漏
			s.endGame(roomID, "游戏超时")
			return
		}
	}
}

// startRound 开始新一轮
func (s *PokerService) startRound(roomID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameState, ok := s.gameStates[roomID]
	if !ok {
		return
	}

	// 重置当前回合下注
	for i := range gameState.Players {
		gameState.Players[i].Bet = 0
	}

	// 根据轮次更新社区牌
	switch gameState.Round {
	case models.RoundPreflop:
		// 已经在DealCards中处理
	case models.RoundFlop:
		// 翻牌，显示前三张社区牌
		gameState.CommunityCards = gameState.CommunityCards[:3]
	case models.RoundTurn:
		// 转牌，显示第四张社区牌
		gameState.CommunityCards = gameState.CommunityCards[:4]
	case models.RoundRiver:
		// 河牌，显示第五张社区牌
		gameState.CommunityCards = gameState.CommunityCards
	case models.RoundShowdown:
		// 摊牌，结束游戏
		s.endGame(roomID, "")
		return
	}

	// 设置当前玩家为庄家后的第一个未弃牌玩家
	activePlayers := 0
	currentPlayerIdx := (gameState.DealerPosition + 1) % len(gameState.Players)
	for i := 0; i < len(gameState.Players); i++ {
		idx := (gameState.DealerPosition + 1 + i) % len(gameState.Players)
		if !gameState.Players[idx].Folded && gameState.Players[idx].Chips > 0 {
			currentPlayerIdx = idx
			activePlayers++
			break
		}
	}

	// 如果活跃玩家不足，结束游戏
	if activePlayers <= 1 {
		s.endGame(roomID, "只剩一名玩家")
		return
	}

	gameState.CurrentPlayer = &models.Player{
		ID:       gameState.Players[currentPlayerIdx].ID,
		Username: gameState.Players[currentPlayerIdx].Username,
	}
}

// handlePlayerAction 处理玩家操作
func (s *PokerService) handlePlayerAction(roomID uint, event PlayerActionEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameState, ok := s.gameStates[roomID]
	if !ok {
		return
	}

	// 查找玩家索引
	playerIdx := -1
	for i, player := range gameState.Players {
		if player.ID == event.PlayerID {
			playerIdx = i
			break
		}
	}

	if playerIdx == -1 {
		return // 玩家不存在
	}

	// 处理不同的操作
	switch event.Action {
	case "check":
		// 看牌，不需要额外处理
	case "call":
		// 跟注
		callAmount := gameState.CurrentBet - gameState.Players[playerIdx].Bet
		if callAmount > gameState.Players[playerIdx].Chips {
			callAmount = gameState.Players[playerIdx].Chips
			gameState.Players[playerIdx].IsAllIn = true
		}
		gameState.Players[playerIdx].Chips -= callAmount
		gameState.Players[playerIdx].Bet += callAmount
		gameState.Pot += callAmount
	case "raise":
		// 加注
		raiseAmount := event.Amount
		if raiseAmount > gameState.Players[playerIdx].Chips {
			raiseAmount = gameState.Players[playerIdx].Chips
			gameState.Players[playerIdx].IsAllIn = true
		}
		gameState.Players[playerIdx].Chips -= raiseAmount
		gameState.Players[playerIdx].Bet += raiseAmount
		gameState.Pot += raiseAmount
		gameState.CurrentBet = gameState.Players[playerIdx].Bet
	case "fold":
		// 弃牌
		gameState.Players[playerIdx].Folded = true
	}

	// 找到下一个未弃牌玩家
	nextPlayerIdx := -1
	for i := 1; i <= len(gameState.Players); i++ {
		idx := (playerIdx + i) % len(gameState.Players)
		if !gameState.Players[idx].Folded && gameState.Players[idx].Chips > 0 && !gameState.Players[idx].IsAllIn {
			nextPlayerIdx = idx
			break
		}
	}

	// 检查是否所有玩家都已行动
	allPlayersActed := true
	for _, player := range gameState.Players {
		if !player.Folded && player.Chips > 0 && !player.IsAllIn && player.Bet < gameState.CurrentBet {
			allPlayersActed = false
			break
		}
	}

	// 检查是否有足够的活跃玩家
	activePlayers := 0
	for _, player := range gameState.Players {
		if !player.Folded {
			activePlayers++
		}
	}

	// 如果只剩一个玩家未弃牌，结束游戏
	if activePlayers <= 1 {
		s.endGame(roomID, "只剩一名玩家")
		return
	}

	// 如果所有玩家都已行动或没有下一个玩家，进入下一轮
	if allPlayersActed || nextPlayerIdx == -1 {
		// 进入下一轮
		switch gameState.Round {
		case models.RoundPreflop:
			gameState.Round = models.RoundFlop
		case models.RoundFlop:
			gameState.Round = models.RoundTurn
		case models.RoundTurn:
			gameState.Round = models.RoundRiver
		case models.RoundRiver:
			gameState.Round = models.RoundShowdown
		}
		s.startRound(roomID)
	} else {
		// 设置下一个玩家
		gameState.CurrentPlayer = &models.Player{
			ID:       gameState.Players[nextPlayerIdx].ID,
			Username: gameState.Players[nextPlayerIdx].Username,
		}
	}
}

// endGame 结束游戏
func (s *PokerService) endGame(roomID uint, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameState, ok := s.gameStates[roomID]
	if !ok {
		return
	}

	// 更新游戏状态
	gameState.Status = "ended"
	gameState.CurrentPlayer = nil

	// 计算获胜者
	var winner *models.Player
	var resultMessage string

	// 如果只有一名玩家未弃牌，则直接宣布获胜
	if message == "只剩一名玩家" {
		for _, player := range gameState.Players {
			if !player.Folded {
				winner = &models.Player{
					ID:       player.ID,
					Username: player.Username,
				}
				resultMessage = "只剩一名玩家，" + winner.Username + "获胜"
				break
			}
		}
	} else {
		// 摊牌比较牌型，找出获胜者
		// 注：这里简化处理，实际德州扑克比牌规则较复杂
		// 应该实现完整的牌型判断和比较逻辑
		winner = &models.Player{
			ID:       gameState.Players[0].ID,
			Username: gameState.Players[0].Username,
		}
		resultMessage = winner.Username + "获胜"
	}

	// 创建游戏结果
	result := &models.PokerGameResult{
		Message: resultMessage,
		Winner:  winner,
		Players: make([]models.PlayerResult, 0),
	}

	// 保存游戏结果
	s.gameResults[roomID] = result

	// 更新房间状态
	s.roomService.roomDAO.UpdateRoomStatus(roomID, "ended")

	// 清理资源
	delete(s.gameStates, roomID)
	delete(s.handCards, roomID)
	close(s.gameEventsChan[roomID])
	delete(s.gameEventsChan, roomID)
}
