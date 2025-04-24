package models

import (
	"time"

	"gorm.io/gorm"
)

// Game 游戏模型
type Game struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	ImageURL    string         `json:"image_url" gorm:"size:255"`
	MinPlayers  int            `json:"min_players" gorm:"default:2"`
	MaxPlayers  int            `json:"max_players" gorm:"default:10"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Rooms       []Room         `json:"rooms,omitempty" gorm:"foreignKey:GameID"`
}

// Room 游戏房间模型
type Room struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"size:100;not null"`
	GameID     uint           `json:"game_id"`
	Game       Game           `json:"game,omitempty" gorm:"foreignKey:GameID"`
	Status     string         `json:"status" gorm:"size:20;default:'waiting'"` // waiting, playing, ended
	MaxPlayers int            `json:"max_players"`
	CreatorID  uint           `json:"creator_id"`
	Creator    User           `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Players    []*User        `json:"players,omitempty" gorm:"many2many:room_players;"`
}

// RoomPlayer 房间玩家关联表
type RoomPlayer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoomID    uint      `json:"room_id"`
	UserID    uint      `json:"user_id"`
	IsReady   bool      `json:"is_ready" gorm:"default:false"`
	JoinedAt  time.Time `json:"joined_at" gorm:"autoCreateTime"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
