# 游戏平台 Docker 部署指南

## 项目结构

- `frontend/`: 前端 Vue.js 项目
- `backend/`: 后端 Go 项目
- `nginx/`: Nginx 配置和 Dockerfile
- `docker-compose.yml`: Docker Compose 配置文件
- `start.sh`: 启动脚本

## 系统架构

- **前端**: Vue.js + Element Plus
- **后端**: Go + Gin + GORM
- **数据库**: MySQL
- **代理**: Nginx (反向代理前后端服务)

## 容器配置

本项目使用 Docker Compose 部署以下服务:

1. **MySQL 数据库**
   - 容器名称: game-db
   - 端口: 3306

2. **后端 API 服务**
   - 容器名称: game-backend
   - 端口: 8080

3. **前端构建服务**
   - 容器名称: game-frontend-build

4. **Nginx 服务**
   - 容器名称: game-nginx
   - 端口: 80 (对外提供服务)

## 部署步骤

### 前提条件

确保已安装:
- Docker (20.10.0+)
- Docker Compose (2.0.0+)

### 部署命令

1. 克隆项目代码:
   ```bash
   git clone <项目仓库URL>
   cd <项目目录>
   ```

2. 运行启动脚本:
   ```bash
   ./start.sh
   ```

部署完成后，访问 http://localhost 即可使用应用。

### 查看日志

```bash
# 查看所有容器日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f <服务名>  # 例如: backend, nginx, db
```

### 停止服务

```bash
docker-compose down
```

## 配置说明

- 数据库配置位于 `backend/configs/docker_config.yaml`
- Nginx 配置位于 `nginx/nginx.conf`
- Docker 相关配置位于各个目录的 Dockerfile 和根目录的 docker-compose.yml

## 开发环境

如需在本地开发测试，参考:
- 前端: `frontend/README.md`
- 后端: `backend/README.md` 