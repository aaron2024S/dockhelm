# syntax=docker/dockerfile:1

# ─────────────────────────────────────────────────────────────
# 关于 --platform=$BUILDPLATFORM（两段构建都钉在构建机上）：
# 不加这一句时，多架构构建会把「前端安装+打包」「Go 编译」在 QEMU 模拟下
# 对 arm64 再完完整整跑一遍 —— 一次发版 5 分半钟基本都耗在这上面，而且
# 日志长时间不动，看起来像卡死。钉到构建机后改成交叉编译：产物一样，
# 但少了整轮模拟执行，只有最后那段 apk 装证书还在模拟里跑（几秒钟）。
# ─────────────────────────────────────────────────────────────

# ─────────────────────────────────────────────────────────────
# 1) 前端构建（产物与架构无关，只构建一次、两个平台共用）
# ─────────────────────────────────────────────────────────────
FROM --platform=$BUILDPLATFORM node:22-alpine AS web
WORKDIR /src/web

# 先只拷贝清单，让依赖层可以被缓存
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund

COPY web/ ./
RUN npm run build


# ─────────────────────────────────────────────────────────────
# 2) 后端构建（CGO 关掉，才能静态链接、才能交叉编译 arm64）
# ─────────────────────────────────────────────────────────────
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build
WORKDIR /src

# BuildKit 会按目标平台自动注入这两个参数（amd64/arm64 各一次）。
# 必须在这里重新 ARG 声明一次，否则 FROM 之后拿不到。
ARG TARGETOS
ARG TARGETARCH

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 用真实的前端产物覆盖仓库里的那一份
COPY --from=web /src/web/dist ./web/dist

# CGO 关掉 + 显式指定 GOARCH ⇒ 在构建机上直接交叉编译出目标架构的二进制，
# 不需要为 arm64 起模拟器。这就是上面钉 $BUILDPLATFORM 换来的收益。
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build \
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

# 注意：这里**故意不设** DOCKHELM_LISTEN。
# 它的优先级最高，一旦在镜像里写死，用户再加 -e PORT=8081 / -e DOCKHELM_PORT=8081
# 就永远不会生效（试半天找不出原因的那种坑）。默认端口由程序内部提供。
ENV DOCKHELM_DATA=/data \
    TZ=Asia/Shanghai

VOLUME ["/data"]
# 只是声明，方便 --publish-all 与面板自动提示；实际端口仍以 DOCKHELM_LISTEN /
# DOCKHELM_PORT / PORT 为准。
EXPOSE 5923

# 健康检查的端口要和程序里的解析顺序保持一致：
# DOCKHELM_LISTEN 的冒号后半段 → DOCKHELM_PORT → PORT → 5923
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD sh -c 'p="${DOCKHELM_LISTEN##*:}"; [ -n "$p" ] || p="${DOCKHELM_PORT:-${PORT:-5923}}"; wget -qO- "http://127.0.0.1:$p/api/health" || exit 1'

# 需要读写 /var/run/docker.sock 与宿主机 docker 目录，所以以 root 运行。
# 这不是懒 —— 是 Docker 套接字本身就等价于宿主机 root，换非 root 用户并不会更安全，
# 只会让挂进来的目录读不到。真正的边界是「别把这个端口暴露到公网」。
ENTRYPOINT ["/usr/local/bin/dockhelm"]
