package dao

import (
	"games/config"
	"games/models"
)

// UserDAO 用户数据访问对象
type UserDAO struct{}

// NewUserDAO 创建用户DAO实例
func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

// Create 创建用户
func (dao *UserDAO) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// GetByID 根据ID获取用户
func (dao *UserDAO) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (dao *UserDAO) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UsernameExists 检查用户名是否已存在
func (dao *UserDAO) UsernameExists(username string) bool {
	var count int64
	config.DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}
