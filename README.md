# Dockhelm · 容器舵手

> 看得清、控得准的 Docker 容器管理面板

Dockhelm 是一个自托管的 Docker 容器管理面板：容器与镜像管理、**同源比对的镜像更新检测与自动更新**、
定时启停计划任务、可自由配置的镜像加速源、容器配置快照与 compose 项目文件备份还原，以及多渠道事件通知。

它的核心设计目标是修掉同类工具（尤其是 dockerCopilot）最让人头疼的那个毛病：
**「一键更新十几个容器，把根本没更新的容器也全停掉了」。**

- 单二进制 + `//go:embed` 内嵌前端，一个容器就能跑完
- 后端只有两个第三方依赖（`robfig/cron` 与 `golang.org/x/crypto`），HTTP 用标准库
- 持久化是原子写入的 JSON，不需要外部数据库，也没有 CGO
- 因此可以 `CGO_ENABLED=0` 交叉编译 `amd64` / `arm64`

---

## 它到底修了什么

### 1. 检测与拉取必须同源

dockerCopilot 的检测路径自己拼 registry 请求、用一份写死在它代码里、用户既看不见也改不了的加速站列表，
**从不读取守护进程 `daemon.json` 里真正生效的 `registry-mirrors`**。
于是「检测说有新版本、拉取说已是最新」会永久互相矛盾 —— 这就是你反复看到
`redis:alpine` 提示有更新、但怎么拉都没有更新的根因。

Dockhelm 的检测走 Docker 守护进程自己的 `/distribution/{name}/json`
（与 `docker pull` 完全同一套仓库端点解析），**检测失败一律标记为「无法判定」，
绝不退化成「有新版本」**。所以永远不会出现「检测说有、拉取说没有」的永久误报。

### 2. 镜像没变就不碰容器

dockerCopilot 的更新流程是无条件执行
`pull → stop → rename → create → start`，算出了新镜像 ID 只用来决定要不要删旧镜像，
**从不用于决定要不要停容器**。两者叠加，就是「一次更新十几个，连没更新的也全被停掉」。

Dockhelm 的流程是：

```
拉取镜像
  └─ 比对容器当前使用的镜像 ID
       ├─ 一致  → 直接返回。容器一个字节都不碰（结果标记为「已是最新·已跳过」）
       └─ 不一致 → 停旧容器 → 改名保留为 <名字>__bak_<时间>
                   → 用原配置创建同名新容器 → 启动 → 健康检查
                        ├─ 通过 → 清理备份容器（可选保留）
                        └─ 失败 → 把旧容器改回原名并启动（自动回滚）
```

检测是**只读**的：只向守护进程的 `/distribution` 要一份 manifest 摘要，与本地摘要比对，
不拉层、不停容器。真正会动容器的只有容器页的「更新 N 个容器」、容器详情页的单容器更新，
以及按周期自动跑的自动更新（开关与频率在设置页）。

### 3. 重建容器时配置要「一模一样」

不只是照抄 `Config` 与 `HostConfig` 那么简单，下面这些坑都踩过了：

- 从**顶层 `Mounts`** 重建挂载（它才是守护进程解析后的真相），
  删掉 `Config` 里已废弃的 `Binds` / `Mounts` / `ContainerIDFile` / `MacAddress` / `Shell`
- **匿名卷按 Name 复用** —— 否则每次重建容器都会产生一批新的匿名卷，老数据变孤儿
- 只在 `Hostname` 等于容器短 ID 时才删掉它（否则会保留你自定义的 hostname）
- 从 `NetworkSettings.Networks` 恢复网络别名与静态 IP，但**跳过已经不存在的网络**
- 处理 `AutoRemove: true` 的容器（它在 stop 时就会把自己删掉）

### 4. 备份要说实话

绑定挂载的数据在宿主机目录上，Dockhelm 默认看不见它。
界面上会明确标注「知道路径但看不见」，**绝不假装备份成功** ——
只记路径，并告诉你去挂载哪个目录才能看到内容。

---

## 功能

| 模块 | 能力 |
| --- | --- |
| **总览** | 容器/镜像数量、可用更新、数据盘占用、Docker 环境、实时活动流、最近操作记录 |
| **容器** | 启停/重启/暂停、重命名、删除、日志（跟随模式）、实时资源占用、环境变量与挂载一览、环境变量自动脱敏；**更新徽标与一键批量更新**（先拉取再比对镜像 ID，镜像没变就完全不动容器） |
| **更新检测与自动更新** | 顶栏一键只读巡检（检测在容器页出徽标）、按周期自动更新（频率与排除列表在设置页）、镜像 ID 比对、跳过与回滚结果明示、更新进度实时推送 |
| **镜像** | 列表（体积、引用的容器数、未使用标记）、删除、清理未使用镜像 |
| **计划任务** | 标准 5 段 cron、六种动作（启动/停止/重启/更新/备份/清理镜像）、目标可多选、下次执行时间预览 |
| **加速源** | 显示守护进程**真实生效**的镜像站、预置常用加速站（首次启动自动写入「我的加速源」，不想要就逐条删、删空可一键找回）、批量测速、一键生成 `daemon.json` 片段 |
| **备份与恢复** | 容器配置快照（单容器 / 一键全量，按批次聚合）、差异预览（镜像/命令/端口/挂载/环境变量…）、一键还原、自动保留策略；compose 项目的 yaml/.env 备份、下载与还原 |
| **通知** | 10 类渠道预设 + 自定义 Webhook、按事件订阅、静默时段、去重窗口、每日上限、失败重试、推送历史 |
| **登录鉴权** | 单用户密码（bcrypt）、HttpOnly 会话、连续失败锁定、改密后强制下线 |
| **关于** | 应用介绍、版本号、提交、构建时间、作者与项目主页、功能与运行环境一览 |

---

## 快速开始

### 方式一：用现成镜像（推荐）

```bash
mkdir -p /volume1/docker/dockhelm && cd /volume1/docker/dockhelm
curl -O https://raw.githubusercontent.com/aaron2024S/dockhelm/main/docker-compose.yml
# 改一下镜像 tag 和 docker 目录路径（端口默认 5923，要换见下面的「端口」一节），然后：
docker compose up -d
```

打开 `http://<你的NAS_IP>:5923`，首次进入会让你设置访问密码。

> **首次登录没有默认密码**，也不会有 `admin/admin` 这种出厂账号 —— 容器起来后打开面板，会直接进入「设置访问密码」引导页，密码由你自己设定（至少 6 位，只保存 bcrypt 哈希，明文不落盘）。想跳过这一步，可以在 compose 里加一行 `DOCKHELM_PASSWORD=你的密码`，容器启动时会直接设好。忘了密码：删掉 `data/auth.json` 后重启容器即可重设。

镜像 tag 有两个，按需选一个：

| tag | 含义 |
|---|---|
| `aaron2024s/dockhelm:0.4.6` | 固定版本。**推荐**，升级由你决定，行为可复现 |
| `aaron2024s/dockhelm:latest` | 总是指向最近一次正式发版。图省事用它，代价是每次 `docker compose pull` 都可能变版本 |

### 方式二：自己构建

```bash
git clone https://github.com/aaron2024S/dockhelm.git
cd dockhelm
./scripts/build-image.sh --load          # 默认打 :<版本> 和 :latest
                                         # 或 -p linux/amd64,linux/arm64 出多架构（--push 推到 Docker Hub）

# 或者直接 compose 构建
docker compose up -d --build
```

### 方式三：不用 Docker 跑后端（开发用）

```bash
cd web && npm install && npm run build && cd ..
go build -o dockhelm .
DOCKHELM_DATA=./data DOCKER_HOST=unix:///var/run/docker.sock ./dockhelm
```

---

## 端口

默认监听 **`5923`**。

没有用 8080，是因为它太容易撞车 —— 群晖 DSM 的登录反代、各种路由器管理页、
以及一大半自托管应用都默认占着它。5923 在 IANA 注册表里没有常见占用，日常不会打架。

三种改法，按省事程度排：

| 想干什么 | 怎么做 |
| --- | --- |
| **只改对外端口**（推荐） | 改 `docker-compose.yml` 里 `ports` 的**左边**，例如 `"8090:5923"`。容器内还是 5923，别的都不用配 |
| **用 `.env` 一次改完** | 把仓库里的 `.env.example` 复制成 `.env`，写一行 `DOCKHELM_PORT=8090`。模板里 `ports` 与容器内监听引用的就是同一个变量，会一起跟随 |
| **纯 `docker run`** | 加 `-e PORT=8090 -p 8090:8090` 即可（三个变量任选一个） |

优先级从高到低：

```
DOCKHELM_LISTEN   >   DOCKHELM_PORT   >   PORT   >   默认 5923
```

`DOCKHELM_LISTEN` 收**完整地址**（`:5923`、`0.0.0.0:5923`、`127.0.0.1:5923`），
另两个只收**端口号**（`8090` 或 `:8090` 都行）。**值非法（不是 1–65535 的数字）会被忽略、
回落到默认端口** —— 不会因为环境里躺着一个乱写的 `PORT` 就让面板起不来。

启动日志会打出实际生效的地址和它的来源，不用猜：

```
HTTP 服务监听 :5923（来自 默认值）
HTTP 服务监听 :8090（来自 DOCKHELM_PORT）
```

「关于」页的「运行环境」里也会显示当前监听地址与它来自哪个变量。

⚠ 镜像里**故意没有**给 `DOCKHELM_LISTEN` 设 ENV 默认值。它优先级最高，一旦在镜像层写死，
用户再加 `-e PORT=8080` 就永远不会生效 —— 这种坑排查起来很费时间。

---

## 目录挂载 —— 这一节一定要看

Dockhelm 能做到「看见」你的 compose 文件（项目页的 yaml 备份全靠它），靠的是这两条：

```yaml
volumes:
  - /var/run/docker.sock:/var/run/docker.sock     # 必填：管理容器与镜像
  - /volume1/docker:/host/docker                 # 强烈建议：宿主机的 docker 目录
```

**冒号右边叫什么名字都可以。** 左边改成你自己的宿主机目录（群晖一般是 `/volume1/docker`，
威联通可能是 `/share/Container`，普通 Linux 通常是 `/opt/docker`），右边叫 `/host/docker`、
`/compose`、`/whatever` 都行。

Dockhelm 启动时会读**自己这个容器**的挂载表（`Mounts`），自动得到「宿主路径 → 容器内路径」的对应关系。
于是容器 label 里记着的 `/volume1/docker/moontv/docker-compose.yml`，会被换算成
`/host/docker/moontv/docker-compose.yml` 再去读 —— 你不需要为了迁就它把容器内路径也写成 `/volume1/docker`。

- **「备份与恢复」页面顶部**会列出当前生效的全部映射，以及每一条**是否真的可见**，启动日志里也会打一行，
  所以不用猜它到底看见了什么。
- 「两边写一样」（`/volume1/docker:/volume1/docker`）作为一种特例依然有效，只是容器里会出现
  `/volume1/...` 这种只属于宿主机的目录结构，看日志时容易误以为那是容器自己的路径。

万一自动识别不出来（例如你在 compose 里自定义了 `hostname`，Dockhelm 认不出自己那个容器），
把映射显式写一遍即可：

```yaml
environment:
  DOCKHELM_HOST_ROOTS: /volume1/docker=/host/docker
```

语法是 `宿主机路径=容器内路径`（逗号分隔可写多组）；只写一个路径就表示「两边一致」。
**换算失败时它会明说「看不见」，不会假装读到了。**

---

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `DOCKHELM_DATA` | `/data` | 配置与备份的持久化目录 |
| `DOCKHELM_LISTEN` | `:5923` | 完整监听地址（优先级最高）。详见上面的「[端口](#端口)」一节 |
| `DOCKHELM_PORT` | 空 | 只写端口号；不填就用默认 `5923`。优先级高于通用的 `PORT` |
| `PORT` | 空 | 只写端口号（通用约定，Railway / Zeabur / Render 这类平台会自动注入） |
| `DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker 守护进程地址（支持 `unix://` 与 `tcp://`） |
| `DOCKHELM_PASSWORD` | 空 | 设置后每次启动都强制把密码重置为该值。**忘记密码的兜底手段，用完请删掉这行** |
| `DOCKHELM_HOST_ROOTS` | 空 | 手动声明路径映射，写法 `宿主机路径=容器内路径`（逗号分隔可多组；只写一个路径表示两边一致）。默认自动识别，一般不用填 |
| `DOCKHELM_SELF` | 空 | 显式指定 Dockhelm 自身容器名。默认靠 hostname 反查 |
| `DOCKHELM_TRUSTED_PROXIES` | 空 | 反代 IP 白名单。**只有填了它，后端才会信任 `X-Forwarded-For`** |
| `DOCKHELM_LOG_LEVEL` | `info` | `debug` / `info` / `warn` / `error` |

---

## 通知

### 支持的渠道

`Telegram`、`Bark`、`ntfy`、`企业微信机器人`、`钉钉机器人`（含加签）、
`飞书机器人`（含签名）、`Server 酱`、`PushPlus`、`邮件 SMTP`（SSL / STARTTLS）、
`自定义 Webhook`（任意 URL + 请求头 + 请求体模板，都支持变量）。

渠道内部是「一个 URL + 一个请求体模板」的数据：加渠道就是加一行数据，不需要改代码。

### 事件列表

**更新**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `update_available` | 检测到有更新 | 普通 | 关 |
| `update_success` | 容器更新成功 | 普通 | 关 |
| `update_failed` | 容器更新失败 | 紧急 | 开 |
| `batch_update_done` | 批量更新完成 | 普通 | 开 |

**容器**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `container_died` | 容器意外退出 | 紧急 | 开 |
| `container_crashloop` | 容器崩溃循环 | 紧急 | 开 |
| `health_failed` | 健康检查失败 | 紧急 | 开 |

**计划任务**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `schedule_success` | 计划任务执行成功 | 普通 | 关 |
| `schedule_failed` | 计划任务执行失败 | 紧急 | 开 |

**备份**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `backup_success` | 备份完成 | 普通 | 关 |
| `backup_failed` | 备份失败 | 紧急 | 开 |
| `restore_done` | 还原完成 | 紧急 | 开 |

**系统**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `docker_disconnected` | Docker 连接断开 | 紧急 | 开 |
| `docker_reconnected` | Docker 连接恢复 | 紧急 | 开 |

**登录与账户**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `login_success` | 登录成功 | 普通 | 关 |
| `login_failed` | 登录失败 | 紧急 | 开 |
| `login_locked` | 登录连续失败告警 | 紧急 | 开 |
| `password_changed` | 登录密码已修改 | 紧急 | 开 |

**其它**

| 事件代码 | 名称 | 级别 | 默认 |
| --- | --- | --- | --- |
| `test` | 测试消息 | 普通 | 开 |

### 投递策略

- **按事件订阅**，不按渠道 —— 某个事件打开后同时发往所有已启用的渠道
- **静默时段**支持跨天（例如 `23:00–07:00`），普通事件可以「攒着结束后汇总发一条」或「直接丢弃」，
  紧急事件可以设置照常发送
- **去重窗口**：同一容器同一事件在 N 分钟内只发一条（崩溃循环的容器一分钟能产生几十条事件）
- **每日上限**：超过后丢弃，防止打爆你的通知渠道
- **失败重试 3 次并按指数退避**，失败只记入推送历史，**绝不影响更新主流程**
- 能识别「HTTP 200 但响应体里报错」的情况（企业微信、钉钉、飞书、Telegram 都存在这个坑）

### 模板变量

自定义 Webhook 的 URL、请求头、请求体里都能用 `{{变量名}}`，URL 中的变量会自动 URL 转义：

`event`、`eventLabel`、`level`、`title`、`text`、`message`、`container`、`image`、
`result`、`time`、`host`、`url`、`sound`、`priority`、`tags`、`titleEncoded`、`textEncoded`

---

## 数据与备份

配置全部放在 `DOCKHELM_DATA`（默认 `/data`）下：

```
data/
├── dockhelm.db          # JSON 文档：设置、计划任务、渠道、事件订阅、更新记录、登录记录
├── auth.json            # bcrypt 密码哈希 + 会话（权限 0600）
├── backups/
│   ├── containers/<容器名>/<时间戳>.json   # 容器配置快照
│   └── projects/                          # 项目文件备份
└── logs/
```

- 写入一律是「临时文件 → `fsync` → `rename`」的原子替换，断电不会写出半截文件
- 数据量有明确上限：运行记录默认保留 100 条（可在设置页改）、推送历史 500 条、登录记录 500 条
- 快照自动保留：每容器 10 份、最多 30 天、总额 2 GB
- 备份/还原整个配置：直接拷贝 `data` 目录即可

**关于数据文件**：Dockhelm 备份的是**容器配置**（镜像、命令、挂载、环境变量、端口、网络…）
与 **compose 项目的源文件**（yaml / `.env` / override 的原样副本），**不备份任何业务数据** ——
数据库、缓存、上传的文件请照旧用你自己的方式备份。

`compose.yaml` 的备份在「备份与恢复 → compose 项目」标签页：可逐项目备份、下载、
把文件原样写回（还原前会自动再存一份当前状态），前提是上面「目录挂载」那一节做对了。

---

## 安全提醒

Dockhelm 挂载了 `/var/run/docker.sock`，**这等价于宿主机的 root 权限**。

- 务必设置登录密码，不要把它暴露到公网
- 需要在公网访问时，请套一层带 HTTPS 与额外认证的反向代理，
  在设置里填好「面板地址」，并把 `DOCKHELM_TRUSTED_PROXIES` 设成你的反代 IP
- 忘记密码：在 NAS 面板里删掉 `data/auth.json` 后重启容器即可重设；
  或者在 compose 里临时加一行 `DOCKHELM_PASSWORD=新密码`

---

## 技术栈

| 层 | 选型 |
| --- | --- |
| 后端 | Go + 标准库 `net/http`（Go 1.22+ 的 method+pattern 路由），无 Web 框架 |
| 容器对接 | 直接使用 Docker Engine HTTP API（不依赖 docker SDK） |
| 调度 | `github.com/robfig/cron/v3` |
| 密码哈希 | `golang.org/x/crypto/bcrypt` |
| 持久化 | 原子写入的 JSON 文档（无 CGO、无外部数据库） |
| 实时推送 | SSE + 进程内发布订阅（重连时回放每个主题的最后一条事件） |
| 前端 | Vue 3 + TypeScript + Vite + Tailwind CSS v4 + Pinia + Vue Router |
| 图标 | `lucide-vue-next` |

## 开发

```bash
# 后端
go build ./... && go vet ./... && go test ./...

# 前端
cd web
npm install
npm run dev          # 开发服务器，/api 代理到 127.0.0.1:5923
npm run type-check   # vue-tsc
npm run build
```

发新版本：改 `internal/version/version.go` 里的 `Version`，提交并打 tag，CI 就会构建 amd64 + arm64 并推送**两个** tag（`:<版本>` 与 `:latest`，latest 总是指向最近一次发版）：

```bash
git tag v0.4.6 && git push origin v0.4.6
```

推 `main` 不会触发任何 CI —— 只有打 tag 或在 Actions 页面手动 Run workflow 才会构建并推送。手动 Run workflow 同样会移动 `latest`，所以别拿它做试探性构建。

需要在仓库 Secrets 里配置 `DOCKERHUB_USERNAME` 与 `DOCKERHUB_TOKEN`。

---

## 作者

**aaron2024S** · [GitHub](https://github.com/aaron2024S)

## 许可

[MIT](LICENSE)
