#!/usr/bin/env bash
# 本机构建 Dockhelm 镜像（需要本机装了 docker 且 buildx 可用）。
#
#   ./scripts/build-image.sh                 # 只建当前架构；打 <用户名>/dockhelm:<版本> 和 :latest
#   ./scripts/build-image.sh -p linux/amd64,linux/arm64
#   ./scripts/build-image.sh --push          # 构建后推送到 Docker Hub（多架构必须加这个）
#   ./scripts/build-image.sh -t myrepo/dockhelm:test   # 自定义 tag（给了 -t 就不再附 :latest）
#   ./scripts/build-image.sh --no-latest     # 版本 tag 之外不要 :latest
#   ./scripts/build-image.sh --load          # 单架构时把镜像装进本地 docker
#
# 版本号自动从 internal/version/version.go 读取，与「关于」页显示的一致。

set -euo pipefail

cd "$(dirname "$0")/.."

IMAGE_OWNER="${DOCKERHUB_USERNAME:-aaron2024s}"
IMAGE_NAME="dockhelm"
PLATFORMS=""
TAG=""
TAG_EXPLICIT=""
NO_LATEST=""
LOAD=""
PUSH=""

while [ $# -gt 0 ]; do
  case "$1" in
    -p|--platforms) PLATFORMS="$2";   shift 2 ;;
    -t|--tag)       TAG="$2"; TAG_EXPLICIT=1; shift 2 ;;
    --no-latest)    NO_LATEST=1;      shift ;;
    --load)         LOAD="--load";    shift ;;
    --push)         PUSH="--push";    shift ;;
    -h|--help)      sed -n '2,11p' "$0"; exit 0 ;;
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

TAG="${TAG:-${IMAGE_OWNER}/${IMAGE_NAME}:${VERSION}}"

# 默认再附一个 :latest（与 CI 的行为一致）；显式给了 -t 或加了 --no-latest 就不附。
LATEST_TAG=""
if [ -z "$TAG_EXPLICIT" ] && [ -z "$NO_LATEST" ]; then
  LATEST_TAG="${IMAGE_OWNER}/${IMAGE_NAME}:latest"
fi

echo "镜像   : $TAG"
[ -n "$LATEST_TAG" ] && echo "附带   : $LATEST_TAG"
echo "版本   : $VERSION"
echo "提交   : $(echo "$COMMIT" | cut -c1-7)"
echo "时间   : $BUILD_TIME"
[ -n "$PLATFORMS" ] && echo "平台   : $PLATFORMS"
[ -n "$PUSH" ] && echo "推送   : 是"
echo

BUILD_ARGS=(
  --build-arg "VERSION=$VERSION"
  --build-arg "COMMIT=$COMMIT"
  --build-arg "BUILD_TIME=$BUILD_TIME"
  --tag "$TAG"
)
[ -n "$LATEST_TAG" ] && BUILD_ARGS+=(--tag "$LATEST_TAG")

if [ -n "$PLATFORMS" ]; then
  # 多架构必须用 buildx（且不能 --load）
  docker buildx build ${LOAD} ${PUSH} --platform "$PLATFORMS" "${BUILD_ARGS[@]}" .
else
  docker buildx build ${LOAD} ${PUSH} "${BUILD_ARGS[@]}" .
fi

echo
echo "完成：$TAG${LATEST_TAG:+  +  $LATEST_TAG}"
echo "跑起来："
echo "  docker run -d --name dockhelm --restart unless-stopped \\"
echo "    -p 8080:8080 \\"
echo "    -v ./data:/data \\"
echo "    -v /var/run/docker.sock:/var/run/docker.sock \\"
echo "    -v /volume1/docker:/host/docker \\"
echo "    $TAG"
