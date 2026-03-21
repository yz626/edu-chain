# ===========================================
# edu-chain Docker构建文件
# ===========================================

# 构建阶段
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev git

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -o edu-chain ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/edu-chain .

# 复制配置文件
COPY config/config.yaml ./config/

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080 9090

# 启动命令
CMD ["./edu-chain"]
