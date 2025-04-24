package dao

import (
	"games/config"
	"games/models"
)

// GameDAO 游戏数据访问对象
type GameDAO struct{}

// NewGameDAO 创建游戏DAO实例
func NewGameDAO() *GameDAO {
	return &GameDAO{}
}

// GetAll 获取所有游戏
func (dao *GameDAO) GetAll() ([]models.Game, error) {
	var games []models.Game
	err := config.DB.Find(&games).Error
	return games, err
}

// GetByID 根据ID获取游戏
func (dao *GameDAO) GetByID(id uint) (*models.Game, error) {
	var game models.Game
	err := config.DB.Where("id = ?", id).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

// RoomDAO 房间数据访问对象
type RoomDAO struct{}

// NewRoomDAO 创建房间DAO实例
func NewRoomDAO() *RoomDAO {
	return &RoomDAO{}
}

// Create 创建房间
func (dao *RoomDAO) Create(room *models.Room) error {
	return config.DB.Create(room).Error
}

// GetByID 根据ID获取房间
func (dao *RoomDAO) GetByID(id uint) (*models.Room, error) {
	var room models.Room
	err := config.DB.Preload("Game").Preload("Creator").Preload("Players").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// GetByGameID 根据游戏ID获取房间列表
func (dao *RoomDAO) GetByGameID(gameID uint) ([]models.Room, error) {
	var rooms []models.Room
	err := config.DB.Preload("Creator").Where("game_id = ?", gameID).Find(&rooms).Error
	return rooms, err
}

// AddPlayerToRoom 将玩家添加到房间
func (dao *RoomDAO) AddPlayerToRoom(roomID, userID uint) error {
	roomPlayer := models.RoomPlayer{
		RoomID: roomID,
		UserID: userID,
	}
	return config.DB.Create(&roomPlayer).Error
}

// RemovePlayerFromRoom 将玩家从房间移除
func (dao *RoomDAO) RemovePlayerFromRoom(roomID, userID uint) error {
	return config.DB.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&models.RoomPlayer{}).Error
}

// IsPlayerInRoom 检查玩家是否在房间中
func (dao *RoomDAO) IsPlayerInRoom(roomID, userID uint) bool {
	var count int64
	config.DB.Model(&models.RoomPlayer{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count)
	return count > 0
}

// CountPlayersInRoom 计算房间中的玩家数量
func (dao *RoomDAO) CountPlayersInRoom(roomID uint) int64 {
	var count int64
	config.DB.Model(&models.RoomPlayer{}).Where("room_id = ?", roomID).Count(&count)
	return count
}

// UpdateRoomStatus 更新房间状态
func (dao *RoomDAO) UpdateRoomStatus(roomID uint, status string) error {
	return config.DB.Model(&models.Room{}).Where("id = ?", roomID).Update("status", status).Error
}

// UpdatePlayerReadyStatus 更新玩家准备状态
func (dao *RoomDAO) UpdatePlayerReadyStatus(roomID, userID uint, isReady bool) error {
	return config.DB.Model(&models.RoomPlayer{}).Where("room_id = ? AND user_id = ?", roomID, userID).Update("is_ready", isReady).Error
}

// GetPlayersInRoom 获取房间中的所有玩家
func (dao *RoomDAO) GetPlayersInRoom(roomID uint) ([]models.User, error) {
	var users []models.User
	err := config.DB.Table("users").
		Joins("JOIN room_players ON users.id = room_players.user_id").
		Where("room_players.room_id = ?", roomID).
		Find(&users).Error
	return users, err
}

// GetPlayerReadyStatus 获取房间中玩家的准备状态
func (dao *RoomDAO) GetPlayerReadyStatus(roomID uint) (map[uint]bool, error) {
	var roomPlayers []models.RoomPlayer
	err := config.DB.Where("room_id = ?", roomID).Find(&roomPlayers).Error
	if err != nil {
		return nil, err
	}

	readyStatus := make(map[uint]bool)
	for _, rp := range roomPlayers {
		readyStatus[rp.UserID] = rp.IsReady
	}
	return readyStatus, nil
}
