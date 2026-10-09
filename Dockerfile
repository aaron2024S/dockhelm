# syntax=docker/dockerfile:1

# ─────────────────────────────────────────────────────────────
# 1) 前端构建
# ─────────────────────────────────────────────────────────────
FROM node:22-alpine AS web
WORKDIR /src/web

# 先只拷贝清单，让依赖层可以被缓存
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund

COPY web/ ./
RUN npm run build


# ─────────────────────────────────────────────────────────────
# 2) 后端构建（CGO 关掉，才能静态链接、才能交叉编译 arm64）
# ─────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS build
WORKDIR /src

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 用真实的前端产物覆盖仓库里的那一份
COPY --from=web /src/web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags "-s -w \
        -X github.com/aaron2024s/dockhelm/internal/version.Version=${VERSION} \
        -X github.com/aaron2024s/dockhelm/internal/version.Commit=${COMMIT} \
        -X github.com/aaron2024s/dockhelm/internal/version.BuildTime=${BUILD_TIME}" \
      -o /out/dockhelm .


# ─────────────────────────────────────────────────────────────
# 3) 运行镜像
# ─────────────────────────────────────────────────────────────
FROM alpine:3.21

# ca-certificates：访问 registry / 通知渠道的 HTTPS 要用
# tzdata：cron 计划任务与通知里的时间要用本地时区
RUN apk add --no-cache ca-certificates tzdata \
 && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo "Asia/Shanghai" > /etc/timezone

COPY --from=build /out/dockhelm /usr/local/bin/dockhelm

ENV DOCKHELM_DATA=/data \
    DOCKHELM_LISTEN=:8080 \
    TZ=Asia/Shanghai

VOLUME ["/data"]
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/health || exit 1

# 需要读写 /var/run/docker.sock 与宿主机 docker 目录，所以以 root 运行。
# 这不是懒 —— 是 Docker 套接字本身就等价于宿主机 root，换非 root 用户并不会更安全，
# 只会让挂进来的目录读不到。真正的边界是「别把这个端口暴露到公网」。
ENTRYPOINT ["/usr/local/bin/dockhelm"]
