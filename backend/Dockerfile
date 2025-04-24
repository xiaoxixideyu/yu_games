FROM golang:1.22 as build

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 设置环境变量
ENV DOCKER_ENV=true

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# 生产环境
FROM alpine:latest

# 安装CA证书
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=build /app/server /app/
COPY --from=build /app/configs /app/configs

# 暴露端口
EXPOSE 8080

CMD ["./server"] 