# 游戏平台后端服务

基于Gin框架实现的游戏平台后端服务，提供用户注册、登录、游戏列表、房间管理等功能，同时包含德州扑克游戏的实现。

## 功能特性

- 用户认证：注册、登录
- 游戏管理：获取游戏列表
- 房间管理：创建房间、获取房间列表、加入房间
- 德州扑克游戏：基于WebSocket的实时游戏体验

## 技术栈

- Go 语言
- Gin Web框架
- GORM ORM框架
- MySQL数据库
- JWT认证
- WebSocket实时通信

## 环境要求

- Go 1.16+
- MySQL 5.7+

## 安装与运行

1. 克隆代码库

```bash
git clone https://github.com/yourusername/game-server.git
cd game-server
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置

项目使用YAML配置文件，位于`configs`目录：
- `local_config.yaml`：本地开发环境配置
- `docker_config.yaml`：Docker环境配置

通过环境变量`DOCKER_ENV`控制使用哪个配置文件：
```bash
# 使用docker_config.yaml配置
export DOCKER_ENV=true
```

不设置此环境变量或值不为"true"时，默认使用`local_config.yaml`配置。

4. 运行服务

```bash
go run main.go
```

## API接口

详细API接口文档请参考 [API接口文档.md](API接口文档.md)。

### 主要接口列表

- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `GET /api/games` - 获取游戏列表
- `GET /api/games/{gameId}/rooms` - 获取游戏房间列表
- `POST /api/games/{gameId}/rooms` - 创建游戏房间
- `POST /api/rooms/{roomId}/join` - 加入游戏房间
- `GET /ws?roomId={roomId}` - WebSocket连接

## 项目结构

```
.
├── config/         # 配置文件和数据库初始化
├── controllers/    # 控制器
├── dao/            # 数据访问对象
├── middleware/     # 中间件
├── models/         # 数据模型
├── routes/         # 路由配置
├── services/       # 业务逻辑
├── utils/          # 工具函数
├── API接口文档.md    # API接口文档
├── go.mod          # Go模块文件
├── go.sum          # Go模块依赖校验文件
├── main.go         # 主程序入口
└── README.md       # 项目说明
```

## 游戏规则

### 德州扑克

德州扑克是一种风靡全球的扑克游戏。每个玩家获得两张底牌，桌面上有五张公共牌。玩家可以使用自己的两张底牌和五张公共牌中的任意五张组成最好的牌型。

游戏流程：
1. 每个玩家获得两张底牌
2. 第一轮下注
3. 发三张公共牌（翻牌）
4. 第二轮下注
5. 发第四张公共牌（转牌）
6. 第三轮下注
7. 发第五张公共牌（河牌）
8. 最后一轮下注
9. 摊牌比较牌型，确定赢家

牌型从高到低：
- 皇家同花顺：同花色的A-K-Q-J-10
- 同花顺：同花色的顺子
- 四条：四张同点数的牌
- 葫芦：三条加一对
- 同花：五张同花色的牌
- 顺子：五张连续点数的牌
- 三条：三张同点数的牌
- 两对：两个对子
- 一对：一个对子
- 高牌：不构成任何牌型，以最高点数牌比较

## 贡献

欢迎通过Issue和Pull Request的方式参与项目贡献。

## 许可证

MIT许可证 