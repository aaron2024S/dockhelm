import { x as createLucideIcon, d as defineComponent, E as useToastStore, a as onMounted, c as createElementBlock, e as createBaseVNode, t as toDisplayString, f as createVNode, n as normalizeClass, g as unref, R as RefreshCw, k as createTextVNode, j as createCommentVNode, T as Info, F as Fragment, r as renderList, m as ref, p as computed, B as api, s as openBlock, i as createBlock, O as resolveDynamicComponent, P as Layers, a5 as CalendarClock, U as Rocket } from "./index-D794l11H.js";
import { g as formatDuration } from "./format-NZu8lNOP.js";
import { Z as Zap } from "./zap-DrIHHST9.js";
import { H as HardDrive } from "./hard-drive-C1nxKFID.js";
const Boxes = createLucideIcon("BoxesIcon", [
  [
    "path",
    {
      d: "M2.97 12.92A2 2 0 0 0 2 14.63v3.24a2 2 0 0 0 .97 1.71l3 1.8a2 2 0 0 0 2.06 0L12 19v-5.5l-5-3-4.03 2.42Z",
      key: "lc1i9w"
    }
  ],
  ["path", { d: "m7 16.5-4.74-2.85", key: "1o9zyk" }],
  ["path", { d: "m7 16.5 5-3", key: "va8pkn" }],
  ["path", { d: "M7 16.5v5.17", key: "jnp8gn" }],
  [
    "path",
    {
      d: "M12 13.5V19l3.97 2.38a2 2 0 0 0 2.06 0l3-1.8a2 2 0 0 0 .97-1.71v-3.24a2 2 0 0 0-.97-1.71L17 10.5l-5 3Z",
      key: "8zsnat"
    }
  ],
  ["path", { d: "m17 16.5-5-3", key: "8arw3v" }],
  ["path", { d: "m17 16.5 4.74-2.85", key: "8rfmw" }],
  ["path", { d: "M17 16.5v5.17", key: "k6z78m" }],
  [
    "path",
    {
      d: "M7.97 4.42A2 2 0 0 0 7 6.13v4.37l5 3 5-3V6.13a2 2 0 0 0-.97-1.71l-3-1.8a2 2 0 0 0-2.06 0l-3 1.8Z",
      key: "1xygjf"
    }
  ],
  ["path", { d: "M12 8 7.26 5.15", key: "1vbdud" }],
  ["path", { d: "m12 8 4.74-2.85", key: "3rx089" }],
  ["path", { d: "M12 13.5V8", key: "1io7kd" }]
]);
const Cpu = createLucideIcon("CpuIcon", [
  ["rect", { width: "16", height: "16", x: "4", y: "4", rx: "2", key: "14l7u7" }],
  ["rect", { width: "6", height: "6", x: "9", y: "9", rx: "1", key: "5aljv4" }],
  ["path", { d: "M15 2v2", key: "13l42r" }],
  ["path", { d: "M15 20v2", key: "15mkzm" }],
  ["path", { d: "M2 15h2", key: "1gxd5l" }],
  ["path", { d: "M2 9h2", key: "1bbxkp" }],
  ["path", { d: "M20 15h2", key: "19e6y8" }],
  ["path", { d: "M20 9h2", key: "19tzq7" }],
  ["path", { d: "M9 2v2", key: "165o2o" }],
  ["path", { d: "M9 20v2", key: "i2bqo8" }]
]);
const ExternalLink = createLucideIcon("ExternalLinkIcon", [
  ["path", { d: "M15 3h6v6", key: "1q9fwt" }],
  ["path", { d: "M10 14 21 3", key: "gplh6r" }],
  ["path", { d: "M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6", key: "a6xqqp" }]
]);
const Github = createLucideIcon("GithubIcon", [
  [
    "path",
    {
      d: "M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4",
      key: "tonef"
    }
  ],
  ["path", { d: "M9 18c-4.51 2-5-2-7-2", key: "9comsn" }]
]);
const Heart = createLucideIcon("HeartIcon", [
  [
    "path",
    {
      d: "M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z",
      key: "c3ymky"
    }
  ]
]);
const Server = createLucideIcon("ServerIcon", [
  ["rect", { width: "20", height: "8", x: "2", y: "2", rx: "2", ry: "2", key: "ngkwjq" }],
  ["rect", { width: "20", height: "8", x: "2", y: "14", rx: "2", ry: "2", key: "iecqi9" }],
  ["line", { x1: "6", x2: "6.01", y1: "6", y2: "6", key: "16zg32" }],
  ["line", { x1: "6", x2: "6.01", y1: "18", y2: "18", key: "nzw8ys" }]
]);
const ShieldCheck = createLucideIcon("ShieldCheckIcon", [
  [
    "path",
    {
      d: "M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z",
      key: "oel41y"
    }
  ],
  ["path", { d: "m9 12 2 2 4-4", key: "dzmm74" }]
]);
const version = "0.4.4";
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "dh-sub" };
const _hoisted_4 = { class: "ml-auto flex gap-2" };
const _hoisted_5 = ["disabled"];
const _hoisted_6 = { class: "dh-card overflow-hidden" };
const _hoisted_7 = {
  class: "flex flex-col items-start gap-4 p-5 lg:flex-row lg:items-center",
  style: { "background": "radial-gradient(620px 200px at 12% 0%, var(--color-glow), transparent 70%)" }
};
const _hoisted_8 = { class: "min-w-0 flex-1" };
const _hoisted_9 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_10 = { class: "text-[20px] font-semibold tracking-[-0.2px]" };
const _hoisted_11 = { class: "dh-badge dh-badge-accent" };
const _hoisted_12 = { class: "dh-badge dh-badge-plain" };
const _hoisted_13 = { class: "mt-1.5 text-[13px] text-text-2" };
const _hoisted_14 = { class: "mt-1 text-[12px] leading-relaxed text-text-4" };
const _hoisted_15 = { class: "flex flex-none flex-col gap-2" };
const _hoisted_16 = ["href"];
const _hoisted_17 = ["href"];
const _hoisted_18 = { class: "grid grid-cols-2 gap-x-6 gap-y-3 border-t border-line-1 p-4 lg:grid-cols-4" };
const _hoisted_19 = { class: "mt-0.5 text-[13px] font-medium" };
const _hoisted_20 = { class: "mt-0.5 text-[13px] font-medium" };
const _hoisted_21 = { class: "mt-0.5 text-[13px] font-medium" };
const _hoisted_22 = { class: "mt-0.5 text-[13px] font-medium" };
const _hoisted_23 = { class: "mt-0.5 break-all font-mono text-[12px]" };
const _hoisted_24 = { class: "mt-0.5 font-mono text-[12px]" };
const _hoisted_25 = { class: "mt-0.5 font-mono text-[12px]" };
const _hoisted_26 = { class: "mt-0.5 text-[13px] font-medium" };
const _hoisted_27 = { class: "dh-card" };
const _hoisted_28 = { class: "dh-card-head" };
const _hoisted_29 = { class: "grid grid-cols-1 gap-2.5 p-3.5 md:grid-cols-2 xl:grid-cols-4" };
const _hoisted_30 = { class: "flex items-center gap-2" };
const _hoisted_31 = { class: "text-[12.5px] font-medium" };
const _hoisted_32 = { class: "mt-1.5 text-[11.5px] leading-relaxed text-text-4" };
const _hoisted_33 = { class: "grid grid-cols-1 gap-3.5 xl:grid-cols-2" };
const _hoisted_34 = { class: "dh-card" };
const _hoisted_35 = { class: "dh-card-head" };
const _hoisted_36 = { class: "grid grid-cols-2 gap-x-5 gap-y-2.5 p-3.5 text-[12px]" };
const _hoisted_37 = { class: "break-all font-mono text-[11.5px]" };
const _hoisted_38 = { class: "break-all font-mono text-[11.5px]" };
const _hoisted_39 = {
  key: 0,
  class: "text-text-5"
};
const _hoisted_40 = { class: "break-all font-mono text-[11.5px]" };
const _hoisted_41 = { class: "text-text-2" };
const _hoisted_42 = {
  key: 0,
  class: "text-text-5"
};
const _hoisted_43 = { class: "truncate text-text-2" };
const _hoisted_44 = ["title"];
const _hoisted_45 = { class: "text-text-2" };
const _hoisted_46 = { class: "dh-card" };
const _hoisted_47 = { class: "dh-card-head" };
const _hoisted_48 = { class: "flex flex-col p-3.5" };
const _hoisted_49 = { class: "w-[76px] flex-none text-text-5" };
const _hoisted_50 = { class: "min-w-0 flex-1 text-text-2" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "AboutView",
  setup(__props) {
    const toast = useToastStore();
    const data = ref(null);
    const loading = ref(true);
    const about = computed(() => data.value?.about);
    const runtime = computed(() => data.value?.runtime ?? {});
    const ver = computed(() => about.value?.version || version);
    const uptimeSeconds = computed(() => Number(runtime.value.uptimeSeconds ?? 0));
    async function load() {
      loading.value = true;
      try {
        data.value = await api.get("/api/about");
      } catch (e) {
        toast.error("读取版本信息失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function copyCommit() {
      const c = about.value?.commit;
      if (!c) return;
      try {
        await navigator.clipboard.writeText(c);
        toast.success("已复制完整提交号");
      } catch {
        toast.error("复制失败");
      }
    }
    const features = [
      { icon: Boxes, title: "容器管理", desc: "启停、重启、重命名、日志、实时资源占用、环境变量与挂载一览" },
      { icon: Layers, title: "镜像管理", desc: "列表、清理未使用镜像、一键删除；显示体积与被引用情况" },
      { icon: Zap, title: "更新检测与自动更新", desc: "先拉取再比对镜像 ID，镜像没变就完全不动容器；按周期自动更新、失败自动回滚；容器页可一键批量更新" },
      { icon: CalendarClock, title: "计划任务", desc: "标准 cron 表达式，定时启停 / 重启 / 更新 / 备份容器" },
      { icon: Rocket, title: "加速源", desc: "显示守护进程真实生效的镜像站、批量测速、生成 daemon.json 片段" },
      { icon: HardDrive, title: "备份与恢复", desc: "容器配置快照（含 docker run 创建的容器）、compose 项目文件整包备份与下载、差异预览、一键还原" },
      { icon: ShieldCheck, title: "登录鉴权", desc: "单用户密码 + bcrypt 哈希 + HttpOnly 会话 + 连续失败锁定" },
      { icon: Heart, title: "事件通知", desc: "10 类渠道预设 + 自定义 Webhook，事件订阅、静默时段、防轰炸" }
    ];
    const stack = [
      { label: "后端", value: "Go + 标准库 net/http（无 Web 框架）" },
      { label: "容器对接", value: "直接使用 Docker Engine HTTP API" },
      { label: "调度", value: "robfig/cron v3" },
      { label: "密码哈希", value: "golang.org/x/crypto/bcrypt" },
      { label: "持久化", value: "原子写入的 JSON 文档（无 CGO、无外部数据库）" },
      { label: "前端", value: "Vue 3 + TypeScript + Vite + Tailwind CSS + Pinia + Vue Router" },
      { label: "图标", value: "lucide" },
      { label: "分发", value: "单二进制 + //go:embed 内嵌前端，单容器部署" }
    ];
    onMounted(() => void load());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[1] || (_cache[1] = createBaseVNode("div", { class: "dh-h1" }, "关于", -1)),
          createBaseVNode("div", _hoisted_3, " Dockhelm v" + toDisplayString(ver.value) + " · 自托管 Docker 容器管理面板 ", 1),
          createBaseVNode("div", _hoisted_4, [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: loading.value,
              onClick: load
            }, [
              createVNode(unref(RefreshCw), {
                class: normalizeClass(["h-3.5 w-3.5", loading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[0] || (_cache[0] = createTextVNode("刷新 ", -1))
            ], 8, _hoisted_5)
          ])
        ]),
        createBaseVNode("div", _hoisted_6, [
          createBaseVNode("div", _hoisted_7, [
            _cache[4] || (_cache[4] = createBaseVNode("div", { class: "grid h-[58px] w-[58px] flex-none place-items-center rounded-[17px] bg-accent text-accent-ink" }, [
              createBaseVNode("svg", {
                viewBox: "0 0 32 32",
                class: "h-7 w-7"
              }, [
                createBaseVNode("circle", {
                  cx: "16",
                  cy: "16",
                  r: "12.3",
                  fill: "none",
                  stroke: "currentColor",
                  "stroke-width": "2.5"
                }),
                createBaseVNode("path", {
                  "fill-rule": "evenodd",
                  fill: "currentColor",
                  d: "M16 6.6 18.5 13.5 25.4 16 18.5 18.5 16 25.4 13.5 18.5 6.6 16 13.5 13.5ZM18 16A2 2 0 1 0 14 16A2 2 0 1 0 18 16Z"
                })
              ])
            ], -1)),
            createBaseVNode("div", _hoisted_8, [
              createBaseVNode("div", _hoisted_9, [
                createBaseVNode("span", _hoisted_10, toDisplayString(about.value?.name ?? "Dockhelm"), 1),
                createBaseVNode("span", _hoisted_11, "v" + toDisplayString(ver.value), 1),
                createBaseVNode("span", _hoisted_12, toDisplayString(about.value?.license), 1),
                about.value?.commitShort ? (openBlock(), createElementBlock("button", {
                  key: 0,
                  type: "button",
                  class: "dh-tap dh-badge dh-badge-plain cursor-pointer",
                  title: "点击复制完整提交号",
                  onClick: copyCommit
                }, [
                  createVNode(unref(Github), { class: "h-3 w-3" }),
                  createTextVNode(toDisplayString(about.value.commitShort), 1)
                ])) : createCommentVNode("", true)
              ]),
              createBaseVNode("div", _hoisted_13, toDisplayString(about.value?.tagline), 1),
              createBaseVNode("div", _hoisted_14, toDisplayString(about.value?.description), 1)
            ]),
            createBaseVNode("div", _hoisted_15, [
              createBaseVNode("a", {
                href: about.value?.repoURL,
                target: "_blank",
                rel: "noreferrer noopener",
                class: "dh-btn dh-btn-primary"
              }, [
                createVNode(unref(Github), { class: "h-3.5 w-3.5" }),
                _cache[2] || (_cache[2] = createTextVNode("项目主页 ", -1)),
                createVNode(unref(ExternalLink), { class: "h-3 w-3" })
              ], 8, _hoisted_16),
              createBaseVNode("a", {
                href: about.value?.authorURL,
                target: "_blank",
                rel: "noreferrer noopener",
                class: "dh-btn"
              }, [
                createVNode(unref(Heart), { class: "h-3.5 w-3.5" }),
                _cache[3] || (_cache[3] = createTextVNode("作者主页 ", -1)),
                createVNode(unref(ExternalLink), { class: "h-3 w-3" })
              ], 8, _hoisted_17)
            ])
          ]),
          createBaseVNode("div", _hoisted_18, [
            createBaseVNode("div", null, [
              _cache[5] || (_cache[5] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "版本", -1)),
              createBaseVNode("div", _hoisted_19, toDisplayString(ver.value), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[6] || (_cache[6] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "作者", -1)),
              createBaseVNode("div", _hoisted_20, toDisplayString(about.value?.author ?? "—"), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[7] || (_cache[7] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "构建时间", -1)),
              createBaseVNode("div", _hoisted_21, toDisplayString(about.value?.buildTime || "本地构建未注入"), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[8] || (_cache[8] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "许可协议", -1)),
              createBaseVNode("div", _hoisted_22, toDisplayString(about.value?.license), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[9] || (_cache[9] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "提交", -1)),
              createBaseVNode("div", _hoisted_23, toDisplayString(about.value?.commit || "未知"), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[10] || (_cache[10] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "Go 版本", -1)),
              createBaseVNode("div", _hoisted_24, toDisplayString(about.value?.goVersion || "—"), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[11] || (_cache[11] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "运行平台", -1)),
              createBaseVNode("div", _hoisted_25, toDisplayString(about.value?.platform || "—"), 1)
            ]),
            createBaseVNode("div", null, [
              _cache[12] || (_cache[12] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, "已运行", -1)),
              createBaseVNode("div", _hoisted_26, toDisplayString(loading.value ? "载入中…" : unref(formatDuration)(uptimeSeconds.value)), 1)
            ])
          ])
        ]),
        createBaseVNode("div", _hoisted_27, [
          createBaseVNode("div", _hoisted_28, [
            createVNode(unref(Info), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[13] || (_cache[13] = createBaseVNode("span", null, "它能做什么", -1))
          ]),
          createBaseVNode("div", _hoisted_29, [
            (openBlock(), createElementBlock(Fragment, null, renderList(features, (f) => {
              return createBaseVNode("div", {
                key: f.title,
                class: "rounded-[10px] border border-line-1 bg-ink-800 p-3"
              }, [
                createBaseVNode("div", _hoisted_30, [
                  (openBlock(), createBlock(resolveDynamicComponent(f.icon), { class: "h-3.5 w-3.5 text-accent" })),
                  createBaseVNode("span", _hoisted_31, toDisplayString(f.title), 1)
                ]),
                createBaseVNode("div", _hoisted_32, toDisplayString(f.desc), 1)
              ]);
            }), 64))
          ])
        ]),
        createBaseVNode("div", _hoisted_33, [
          createBaseVNode("div", _hoisted_34, [
            createBaseVNode("div", _hoisted_35, [
              createVNode(unref(Server), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[14] || (_cache[14] = createBaseVNode("span", null, "运行环境", -1))
            ]),
            createBaseVNode("div", _hoisted_36, [
              _cache[15] || (_cache[15] = createBaseVNode("div", { class: "text-text-5" }, "数据目录", -1)),
              createBaseVNode("div", _hoisted_37, toDisplayString(runtime.value.dataDir ?? "—"), 1),
              _cache[16] || (_cache[16] = createBaseVNode("div", { class: "text-text-5" }, "监听地址", -1)),
              createBaseVNode("div", _hoisted_38, [
                createTextVNode(toDisplayString(runtime.value.listen ?? "—") + " ", 1),
                runtime.value.listenSource ? (openBlock(), createElementBlock("span", _hoisted_39, "（来自 " + toDisplayString(runtime.value.listenSource) + "）", 1)) : createCommentVNode("", true)
              ]),
              _cache[17] || (_cache[17] = createBaseVNode("div", { class: "text-text-5" }, "Docker 地址", -1)),
              createBaseVNode("div", _hoisted_40, toDisplayString(runtime.value.dockerHost ?? "—"), 1),
              _cache[18] || (_cache[18] = createBaseVNode("div", { class: "text-text-5" }, "Docker 版本", -1)),
              createBaseVNode("div", _hoisted_41, [
                createTextVNode(toDisplayString(runtime.value.dockerVersion ?? "未连接") + " ", 1),
                runtime.value.dockerApiVersion ? (openBlock(), createElementBlock("span", _hoisted_42, "（API v" + toDisplayString(runtime.value.dockerApiVersion) + "）", 1)) : createCommentVNode("", true)
              ]),
              _cache[19] || (_cache[19] = createBaseVNode("div", { class: "text-text-5" }, "Docker 系统", -1)),
              createBaseVNode("div", _hoisted_43, toDisplayString([runtime.value.dockerOs, runtime.value.dockerArch, runtime.value.dockerKernel].filter(Boolean).join(" · ") || "—"), 1),
              _cache[20] || (_cache[20] = createBaseVNode("div", { class: "text-text-5" }, "Docker 数据根", -1)),
              createBaseVNode("div", {
                class: "truncate font-mono text-[11.5px]",
                title: String(runtime.value.dockerRoot ?? "")
              }, toDisplayString(runtime.value.dockerRoot ?? "—"), 9, _hoisted_44),
              _cache[21] || (_cache[21] = createBaseVNode("div", { class: "text-text-5" }, "容器 / 镜像总数", -1)),
              createBaseVNode("div", _hoisted_45, toDisplayString(runtime.value.containersTotal ?? "—") + " 个容器 · " + toDisplayString(runtime.value.imagesTotal ?? "—") + " 个镜像 ", 1)
            ])
          ]),
          createBaseVNode("div", _hoisted_46, [
            createBaseVNode("div", _hoisted_47, [
              createVNode(unref(Cpu), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[22] || (_cache[22] = createBaseVNode("span", null, "技术栈", -1))
            ]),
            createBaseVNode("div", _hoisted_48, [
              (openBlock(), createElementBlock(Fragment, null, renderList(stack, (s) => {
                return createBaseVNode("div", {
                  key: s.label,
                  class: "flex gap-4 border-b border-line-row py-2 text-[12px] last:border-b-0"
                }, [
                  createBaseVNode("div", _hoisted_49, toDisplayString(s.label), 1),
                  createBaseVNode("div", _hoisted_50, toDisplayString(s.value), 1)
                ]);
              }), 64))
            ])
          ])
        ])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
