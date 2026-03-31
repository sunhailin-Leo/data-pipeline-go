
# 构建阶段
FROM golang:1.26-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates gcc musl-dev

# 设置中国镜像加速
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 设置 CGO 并构建
ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o data-pipeline-go ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的工具和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/data-pipeline-go /app/data-pipeline-go

# 暴露 Prometheus metrics 端口
EXPOSE 8080

# 设置工作目录
WORKDIR /app

# 入口点
ENTRYPOINT ["/app/data-pipeline-go"]
