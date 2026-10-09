#!/usr/bin/env bash
# 本机构建 Dockhelm 镜像（需要本机装了 docker 且 buildx 可用）。
#
#   ./scripts/build-image.sh                 # 只建当前架构，tag 用版本号
#   ./scripts/build-image.sh -p linux/amd64,linux/arm64
#   ./scripts/build-image.sh -t myrepo/dockhelm:test
#   ./scripts/build-image.sh --load           # 单架构时把镜像装进本地 docker
#
# 版本号自动从 internal/version/version.go 读取，与「关于」页显示的一致。

set -euo pipefail

cd "$(dirname "$0")/.."

IMAGE_OWNER="${DOCKERHUB_USERNAME:-aaron2024s}"
IMAGE_NAME="dockhelm"
PLATFORMS=""
TAG=""
LOAD=""

while [ $# -gt 0 ]; do
  case "$1" in
    -p|--platforms) PLATFORMS="$2"; shift 2 ;;
    -t|--tag)       TAG="$2";       shift 2 ;;
    --load)         LOAD="--load";  shift ;;
    -h|--help)      sed -n '2,12p' "$0"; exit 0 ;;
    *) echo "未知参数：$1" >&2; exit 2 ;;
  esac
done

VERSION=$(grep -oE '^var Version = "[^"]+"' internal/version/version.go | sed 's/.*"\(.*\)"/\1/')
if [ -z "$VERSION" ]; then
  echo "读不到版本号，请检查 internal/version/version.go" >&2
  exit 1
fi

COMMIT=$(git rev-parse HEAD 2>/dev/null || echo unknown)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# 只打一个 tag，不附 :latest
TAG="${TAG:-${IMAGE_OWNER}/${IMAGE_NAME}:${VERSION}}"

echo "镜像   : $TAG"
echo "版本   : $VERSION"
echo "提交   : $(echo "$COMMIT" | cut -c1-7)"
echo "时间   : $BUILD_TIME"
[ -n "$PLATFORMS" ] && echo "平台   : $PLATFORMS"
echo

BUILD_ARGS=(
  --build-arg "VERSION=$VERSION"
  --build-arg "COMMIT=$COMMIT"
  --build-arg "BUILD_TIME=$BUILD_TIME"
  --tag "$TAG"
)

if [ -n "$PLATFORMS" ]; then
  # 多架构必须用 buildx（且不能 --load）
  docker buildx build ${LOAD} --platform "$PLATFORMS" "${BUILD_ARGS[@]}" .
else
  docker buildx build ${LOAD} "${BUILD_ARGS[@]}" .
fi

echo
echo "完成：$TAG"
echo "跑起来："
echo "  docker run -d --name dockhelm --restart unless-stopped \\"
echo "    -p 8080:8080 \\"
echo "    -v ./data:/data \\"
echo "    -v /var/run/docker.sock:/var/run/docker.sock \\"
echo "    -v /volume1/docker:/volume1/docker \\"
echo "    $TAG"
