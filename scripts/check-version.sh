#!/usr/bin/env bash
# 版本号一致性闸门 —— 在「登录 / 推送 / 构建镜像」之前跑，参数错了别浪费一次发布。
#
#   ./scripts/check-version.sh                # 校验三处声明一致
#   ./scripts/check-version.sh --expect-tag v0.3.0   # 顺带校验 git tag 与声明一致
#
# 为什么需要它：Dockhelm 的版本号有三个声明处 ——
#   1. internal/version/version.go   ← 唯一事实源，「关于」页与 /api/about 用它
#   2. web/package.json              ← 前端包元数据
#   3. web/src/config.ts             ← 设置页页脚显示的 vX.Y.Z
# 第 3 个是容易漏的那一类「看起来是常量、其实是版本号」的第二出处：
# 只改 version.go 不改它，设置页会一直显示旧版本号，而且**没有任何地方会报错**。
#
# --expect-tag 防的是另一种事故：tag 打成了 v0.3.1 而 version.go 还是 0.3.0，
# 于是 CI 推出 Docker 镜像 :0.3.0，git 里却挂着 v0.3.1 —— 版本无法按 tag 回溯。

set -euo pipefail

cd "$(dirname "$0")/.."

EXPECT_TAG=""
while [ $# -gt 0 ]; do
  case "$1" in
    --expect-tag) EXPECT_TAG="$2"; shift 2 ;;
    -h|--help) sed -n '2,18p' "$0"; exit 0 ;;
    *) echo "未知参数：$1" >&2; exit 2 ;;
  esac
done

# ---- 取值：三处都用「带上下文」的精确匹配，不做全库替换 ----

V_GO=$(grep -oE '^var Version = "[^"]+"' internal/version/version.go | sed 's/.*"\(.*\)"/\1/' || true)
# package.json 顶层 version 是 2 空格缩进；嵌套依赖的 version 缩进更深，不会被这条命中
V_PKG=$(grep -m1 -oE '^  "version": "[^"]+"' web/package.json | sed 's/.*"\(.*\)"/\1/' || true)
V_TS=$(grep -oE "^export const version = '[^']+'" web/src/config.ts | sed "s/.*'\(.*\)'/\1/" || true)

fail=0
check_nonempty() {
  if [ -z "$2" ]; then
    echo "  ✗ 读不到版本号：$1"
    fail=1
  fi
}
check_nonempty "internal/version/version.go" "$V_GO"
check_nonempty "web/package.json"             "$V_PKG"
check_nonempty "web/src/config.ts"            "$V_TS"
if [ "$fail" -ne 0 ]; then
  echo
  echo "版本号读取失败，可能是文件结构变了 —— 请同步更新 scripts/check-version.sh" >&2
  exit 1
fi

echo "声明处："
echo "  internal/version/version.go : $V_GO"
echo "  web/package.json            : $V_PKG"
echo "  web/src/config.ts           : $V_TS"

# ---- 格式校验：挡住 latest / v0.3.0 / 0.2 这类手填值 ----
if ! printf '%s' "$V_GO" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+([-+][0-9A-Za-z.-]+)?$'; then
  echo "  ✗ 版本号 '$V_GO' 不是 x.y.z 形式" >&2
  exit 1
fi

# ---- 三处必须一致 ----
if [ "$V_GO" != "$V_PKG" ] || [ "$V_GO" != "$V_TS" ]; then
  echo
  echo "✗ 版本号不一致 —— 三处必须相同，否则页面显示的版本与实际构建的版本会对不上：" >&2
  echo "    internal/version/version.go = $V_GO" >&2
  echo "    web/package.json            = $V_PKG" >&2
  echo "    web/src/config.ts           = $V_TS" >&2
  exit 1
fi

# ---- CI 上额外校验 tag（tag 是 vX.Y.Z，版本号是 X.Y.Z）----
if [ -n "$EXPECT_TAG" ]; then
  TAG_VER="${EXPECT_TAG#v}"
  if [ "$TAG_VER" != "$V_GO" ]; then
    echo
    echo "✗ git tag '$EXPECT_TAG' 与版本号 '$V_GO' 不一致。" >&2
    echo "  继续跑会推出 Docker 镜像 :$V_GO，而 git 里挂着的是 $EXPECT_TAG，版本将无法按 tag 回溯。" >&2
    echo "  要么改 internal/version/version.go，要么换一个 tag。" >&2
    exit 1
  fi
  echo "tag 校验             : $EXPECT_TAG ✓"
fi

echo
echo "✓ 版本号 $V_GO 已核对一致"
