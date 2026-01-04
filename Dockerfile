# 使用官方Go镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bus-backend .

# 使用轻量级Alpine镜像作为运行环境
FROM alpine:latest

# 安装ca证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/bus-backend .

# 复制配置文件
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./bus-backend"]
