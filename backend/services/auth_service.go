package services

import (
	"errors"
	"games/config"
	"games/dao"
	"games/models"
	"games/utils"
)

// AuthService 用户认证服务
type AuthService struct {
	userDAO *dao.UserDAO
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{
		userDAO: dao.NewUserDAO(),
	}
}

// Register 用户注册
func (s *AuthService) Register(username, password string) error {
	// 检查用户名是否已存在
	if s.userDAO.UsernameExists(username) {
		return errors.New("用户名已存在")
	}

	// 创建新用户
	user := &models.User{
		Username: username,
		Password: password,
	}

	return s.userDAO.Create(user)
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	// 根据用户名获取用户
	user, err := s.userDAO.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	appConfig := config.LoadConfig()
	jwt := utils.NewJWT(appConfig.JWT.Secret)
	token, err := jwt.CreateToken(user.ID, appConfig.JWT.ExpireTime)
	if err != nil {
		return "", nil, errors.New("生成令牌失败")
	}

	return token, user, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userDAO.GetByID(id)
}
