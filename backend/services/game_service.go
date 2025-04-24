package services

import (
	"errors"
	"games/dao"
	"games/models"
)

// GameService 游戏服务
type GameService struct {
	gameDAO *dao.GameDAO
}

// NewGameService 创建游戏服务实例
func NewGameService() *GameService {
	return &GameService{
		gameDAO: dao.NewGameDAO(),
	}
}

// GetGameList 获取所有游戏
func (s *GameService) GetGameList() ([]models.Game, error) {
	return s.gameDAO.GetAll()
}

// GetGameByID 根据ID获取游戏
func (s *GameService) GetGameByID(id uint) (*models.Game, error) {
	return s.gameDAO.GetByID(id)
}

// RoomService 游戏房间服务
type RoomService struct {
	roomDAO *dao.RoomDAO
	gameDAO *dao.GameDAO
	userDAO *dao.UserDAO
}

// NewRoomService 创建房间服务实例
func NewRoomService() *RoomService {
	return &RoomService{
		roomDAO: dao.NewRoomDAO(),
		gameDAO: dao.NewGameDAO(),
		userDAO: dao.NewUserDAO(),
	}
}

// GetRoomList 根据游戏ID获取房间列表
func (s *RoomService) GetRoomList(gameID uint) ([]models.Room, error) {
	// 检查游戏是否存在
	_, err := s.gameDAO.GetByID(gameID)
	if err != nil {
		return nil, errors.New("游戏不存在")
	}

	return s.roomDAO.GetByGameID(gameID)
}

// CreateRoom 创建房间
func (s *RoomService) CreateRoom(name string, gameID, creatorID uint) (*models.Room, error) {
	// 检查游戏是否存在
	game, err := s.gameDAO.GetByID(gameID)
	if err != nil {
		return nil, errors.New("游戏不存在")
	}

	// 创建房间
	room := &models.Room{
		Name:       name,
		GameID:     gameID,
		CreatorID:  creatorID,
		Status:     "waiting",
		MaxPlayers: game.MaxPlayers,
	}

	err = s.roomDAO.Create(room)
	if err != nil {
		return nil, err
	}

	// 将创建者添加到房间
	err = s.roomDAO.AddPlayerToRoom(room.ID, creatorID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// JoinRoom 加入房间
func (s *RoomService) JoinRoom(roomID, userID uint) (*models.Room, error) {
	// 检查房间是否存在
	room, err := s.roomDAO.GetByID(roomID)
	if err != nil {
		return nil, errors.New("房间不存在")
	}

	// 检查房间状态
	if room.Status != "waiting" {
		return nil, errors.New("房间已开始游戏或已结束")
	}

	// 检查玩家是否已在房间中
	if s.roomDAO.IsPlayerInRoom(roomID, userID) {
		return nil, errors.New("已在房间中")
	}

	// 检查房间是否已满
	playersCount := s.roomDAO.CountPlayersInRoom(roomID)
	if int(playersCount) >= room.MaxPlayers {
		return nil, errors.New("房间已满")
	}

	// 将玩家添加到房间
	err = s.roomDAO.AddPlayerToRoom(roomID, userID)
	if err != nil {
		return nil, err
	}

	// 重新获取最新的房间信息
	return s.roomDAO.GetByID(roomID)
}

// LeaveRoom 离开房间
func (s *RoomService) LeaveRoom(roomID, userID uint) error {
	// 检查房间是否存在
	room, err := s.roomDAO.GetByID(roomID)
	if err != nil {
		return errors.New("房间不存在")
	}

	// 检查玩家是否在房间中
	if !s.roomDAO.IsPlayerInRoom(roomID, userID) {
		return errors.New("不在房间中")
	}

	// 将玩家从房间移除
	err = s.roomDAO.RemovePlayerFromRoom(roomID, userID)
	if err != nil {
		return err
	}

	// 如果是房主离开且房间中还有其他玩家，则转移房主权限给最早加入的玩家
	if room.CreatorID == userID {
		playersCount := s.roomDAO.CountPlayersInRoom(roomID)
		if playersCount > 0 {
			players, err := s.roomDAO.GetPlayersInRoom(roomID)
			if err != nil {
				return err
			}
			if len(players) > 0 {
				// 更新房主
				room.CreatorID = players[0].ID
				err = s.roomDAO.Create(room) // 使用Create更新（Gorm会根据主键更新）
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GetRoomByID 根据ID获取房间
func (s *RoomService) GetRoomByID(id uint) (*models.Room, error) {
	return s.roomDAO.GetByID(id)
}

// ToggleReady 切换玩家准备状态
func (s *RoomService) ToggleReady(roomID, userID uint) (bool, error) {
	// 检查房间是否存在
	_, err := s.roomDAO.GetByID(roomID)
	if err != nil {
		return false, errors.New("房间不存在")
	}

	// 检查玩家是否在房间中
	if !s.roomDAO.IsPlayerInRoom(roomID, userID) {
		return false, errors.New("不在房间中")
	}

	// 获取当前准备状态
	readyStatus, err := s.roomDAO.GetPlayerReadyStatus(roomID)
	if err != nil {
		return false, err
	}

	// 切换准备状态
	newStatus := !readyStatus[userID]
	err = s.roomDAO.UpdatePlayerReadyStatus(roomID, userID, newStatus)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}
