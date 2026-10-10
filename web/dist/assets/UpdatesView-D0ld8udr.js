import { x as createLucideIcon, d as defineComponent, u as useAppStore, D as useToastStore, y as watch, a as onMounted, o as onUnmounted, c as createElementBlock, e as createBaseVNode, f as createVNode, g as unref, O as Info, k as createTextVNode, t as toDisplayString, L as Download, n as normalizeClass, P as CircleCheck, i as createBlock, F as Fragment, r as renderList, j as createCommentVNode, H as withCtx, J as RouterLink, m as ref, p as computed, B as api, s as openBlock, C as openStream } from "./index-txkAKhMR.js";
import { b as relativeTime, d as formatDateTime, s as shortImage, e as checkLabel } from "./format-CNNCbTDO.js";
import { _ as _sfc_main$1 } from "./Modal.vue_vue_type_script_setup_true_lang--iVMFZPS.js";
import { L as LoaderCircle } from "./loader-circle-CpKqd4Ou.js";
import { Z as Zap } from "./zap-DYqJdGUs.js";
import { P as Play } from "./play-DNeEkkSs.js";
import { T as TriangleAlert } from "./triangle-alert-sVtxwz6t.js";
const CircleHelp = createLucideIcon("CircleHelpIcon", [
  ["circle", { cx: "12", cy: "12", r: "10", key: "1mglay" }],
  ["path", { d: "M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3", key: "1u773s" }],
  ["path", { d: "M12 17h.01", key: "p32p05" }]
]);
const Clock = createLucideIcon("ClockIcon", [
  ["circle", { cx: "12", cy: "12", r: "10", key: "1mglay" }],
  ["polyline", { points: "12 6 12 12 16 14", key: "68esgv" }]
]);
const ShieldQuestion = createLucideIcon("ShieldQuestionIcon", [
  [
    "path",
    {
      d: "M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z",
      key: "oel41y"
    }
  ],
  ["path", { d: "M9.1 9a3 3 0 0 1 5.82 1c0 2-3 3-3 3", key: "mhlwft" }],
  ["path", { d: "M12 17h.01", key: "p32p05" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "dh-sub" };
const _hoisted_4 = { class: "grid grid-cols-2 gap-3 lg:grid-cols-4" };
const _hoisted_5 = { class: "dh-card p-3.5" };
const _hoisted_6 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_7 = { class: "dh-card p-3.5" };
const _hoisted_8 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_9 = { class: "mt-1.5 text-[22px] font-semibold leading-none" };
const _hoisted_10 = { class: "dh-card p-3.5" };
const _hoisted_11 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_12 = { class: "mt-1.5 text-[22px] font-semibold leading-none" };
const _hoisted_13 = { class: "dh-card p-3.5" };
const _hoisted_14 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_15 = { class: "mt-1.5 text-[22px] font-semibold leading-none" };
const _hoisted_16 = {
  key: 0,
  class: "dh-card"
};
const _hoisted_17 = { class: "dh-card-head" };
const _hoisted_18 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_19 = { class: "flex flex-col" };
const _hoisted_20 = { class: "w-[170px] flex-none truncate font-mono text-[11.5px] text-text-3" };
const _hoisted_21 = { class: "min-w-0 flex-1 truncate text-text-2" };
const _hoisted_22 = { class: "dh-card" };
const _hoisted_23 = { class: "dh-card-head" };
const _hoisted_24 = { class: "ml-auto flex flex-wrap items-center gap-x-3 gap-y-1" };
const _hoisted_25 = ["title"];
const _hoisted_26 = ["disabled"];
const _hoisted_27 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_28 = { class: "rounded-[10px] border border-line-1 bg-ink-800" };
const _hoisted_29 = { class: "flex flex-wrap items-center gap-2 border-b border-line-1 px-3 py-2" };
const _hoisted_30 = { class: "dh-badge dh-badge-warn" };
const _hoisted_31 = {
  key: 0,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_32 = { class: "dh-badge dh-badge-plain" };
const _hoisted_33 = {
  key: 0,
  class: "px-3 py-3 text-[12px] text-text-5"
};
const _hoisted_34 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_35 = { class: "dh-table" };
const _hoisted_36 = { class: "max-w-[220px] truncate font-mono text-[11px] text-text-5" };
const _hoisted_37 = { class: "text-[11.5px]" };
const _hoisted_38 = {
  key: 0,
  class: "dh-badge dh-badge-warn"
};
const _hoisted_39 = {
  key: 1,
  class: "dh-badge dh-badge-accent"
};
const _hoisted_40 = {
  key: 2,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_41 = {
  key: 4,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_42 = { class: "ml-2 text-text-5" };
const _hoisted_43 = { class: "flex flex-col gap-3 text-[12.5px] leading-relaxed text-text-3" };
const _hoisted_44 = { class: "flex items-start gap-2.5" };
const _hoisted_45 = { class: "flex items-start gap-2.5" };
const _hoisted_46 = { class: "flex items-start gap-2.5" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "UpdatesView",
  setup(__props) {
    const app = useAppStore();
    const toast = useToastStore();
    const results = ref([]);
    const checkedAt = ref("");
    const progress = ref([]);
    const showInfo = ref(false);
    let closeStream = null;
    let disposed = false;
    const auto = ref(null);
    const autoBusy = ref(false);
    const available = computed(() => results.value.filter((r) => r.status === "update_available"));
    const resultByName = computed(() => {
      const m = /* @__PURE__ */ new Map();
      for (const r of results.value) m.set(r.container, r);
      return m;
    });
    const willUpdate = computed(() => (auto.value?.candidates ?? []).filter((c) => c.willUpdate));
    const skipped = computed(() => (auto.value?.candidates ?? []).filter((c) => !c.willUpdate && (c.protected || c.excluded)));
    async function load() {
      try {
        const res = await api.get("/api/updates");
        results.value = res.results ?? [];
        checkedAt.value = res.checkedAt;
      } catch (e) {
        toast.error("读取巡检结果失败", e instanceof Error ? e.message : String(e));
      }
    }
    const toneOf = (status) => checkLabel(status).tone;
    async function loadAuto() {
      try {
        auto.value = await api.get("/api/updates/auto");
      } catch {
        auto.value = null;
      }
    }
    async function runAutoUpdate() {
      autoBusy.value = true;
      try {
        await api.post("/api/updates/auto-run", {});
        toast.info("已开始自动更新", "完成后本页会自动刷新");
        for (let i = 0; i < 40 && !disposed; i++) {
          await new Promise((r) => window.setTimeout(r, 1500));
          if (disposed) return;
          const info = await api.get("/api/updates/auto").catch(() => null);
          if (disposed) return;
          if (info) auto.value = info;
          if (info && !info.running) break;
        }
        if (disposed) return;
        await Promise.all([load(), loadAuto()]);
        const last = auto.value?.lastRun;
        if (last) {
          if (last.updated + last.failed === 0) {
            toast.info("本轮没有需要更新的容器", `巡检 ${last.checked} 个，均已是新版`);
          } else {
            toast.success("自动更新完成", `更新 ${last.updated} 个，失败 ${last.failed} 个`);
          }
        }
      } catch (e) {
        if (disposed) return;
        toast.error("启动失败", e instanceof Error ? e.message : String(e));
      } finally {
        if (!disposed) autoBusy.value = false;
      }
    }
    watch(
      () => app.checkTick,
      () => {
        void load();
        void loadAuto();
      }
    );
    onMounted(() => {
      void load();
      void loadAuto();
      closeStream = openStream("/api/events/stream", (topic, ev) => {
        if (topic !== "update") return;
        const d = ev.data ?? {};
        const name = String(d.container ?? "");
        if (ev.kind === "step" || ev.kind === "container_status" || ev.kind === "pull_progress") {
          const msg = String(d.message ?? d.status ?? "");
          if (!msg) return;
          const idx = progress.value.findIndex((p) => p.name === name);
          const entry = { name, message: msg, status: String(ev.status ?? "running") };
          if (idx >= 0) progress.value[idx] = entry;
          else progress.value.unshift(entry);
          if (progress.value.length > 12) progress.value.pop();
        }
        if (ev.kind === "container_status") {
          const st = String(d.status ?? "");
          if (["updated", "failed", "broken", "up_to_date"].includes(st)) {
            const entry = progress.value.find((p) => p.name === name);
            if (entry) entry.status = st === "updated" || st === "up_to_date" ? "success" : "failed";
          }
        }
        if (ev.kind === "batch_done") {
          void load();
        }
        if (ev.kind === "auto_check_done" || ev.kind === "auto_done") {
          void loadAuto();
          if (ev.kind === "auto_done") void load();
        }
      });
    });
    onUnmounted(() => {
      disposed = true;
      closeStream?.();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[4] || (_cache[4] = createBaseVNode("div", { class: "dh-h1" }, "更新中心", -1)),
          createBaseVNode("button", {
            class: "dh-tap inline-flex items-center gap-1 rounded-md px-1.5 py-[3px] text-[11.5px] text-text-5 transition-colors hover:bg-ink-750 hover:text-accent",
            title: "为什么 Dockhelm 不会误停容器",
            onClick: _cache[0] || (_cache[0] = ($event) => showInfo.value = true)
          }, [
            createVNode(unref(Info), { class: "h-3.5 w-3.5" }),
            _cache[3] || (_cache[3] = createTextVNode("了解更多 ", -1))
          ]),
          createBaseVNode("div", _hoisted_3, toDisplayString(checkedAt.value ? `上次检测 ${unref(relativeTime)(checkedAt.value)} · 共比对 ${results.value.length} 个容器` : "还没有检测过"), 1)
        ]),
        createBaseVNode("div", _hoisted_4, [
          createBaseVNode("div", _hoisted_5, [
            createBaseVNode("div", _hoisted_6, [
              createVNode(unref(Download), { class: "h-3.5 w-3.5" }),
              _cache[5] || (_cache[5] = createTextVNode("有可用更新", -1))
            ]),
            createBaseVNode("div", {
              class: normalizeClass(["mt-1.5 text-[22px] font-semibold leading-none", available.value.length ? "text-warn-text" : ""])
            }, toDisplayString(available.value.length), 3)
          ]),
          createBaseVNode("div", _hoisted_7, [
            createBaseVNode("div", _hoisted_8, [
              createVNode(unref(CircleCheck), { class: "h-3.5 w-3.5" }),
              _cache[6] || (_cache[6] = createTextVNode("已是最新", -1))
            ]),
            createBaseVNode("div", _hoisted_9, toDisplayString(results.value.filter((r) => r.status === "up_to_date").length), 1)
          ]),
          createBaseVNode("div", _hoisted_10, [
            createBaseVNode("div", _hoisted_11, [
              createVNode(unref(ShieldQuestion), { class: "h-3.5 w-3.5" }),
              _cache[7] || (_cache[7] = createTextVNode("无法判定", -1))
            ]),
            createBaseVNode("div", _hoisted_12, toDisplayString(results.value.filter((r) => r.status === "unknown").length), 1),
            _cache[8] || (_cache[8] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-6" }, "网络/认证问题导致", -1))
          ]),
          createBaseVNode("div", _hoisted_13, [
            createBaseVNode("div", _hoisted_14, [
              createVNode(unref(CircleHelp), { class: "h-3.5 w-3.5" }),
              _cache[9] || (_cache[9] = createTextVNode("本地镜像", -1))
            ]),
            createBaseVNode("div", _hoisted_15, toDisplayString(results.value.filter((r) => r.status === "no_upstream").length), 1),
            _cache[10] || (_cache[10] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-6" }, "本地构建，无远端可比对", -1))
          ])
        ]),
        progress.value.length ? (openBlock(), createElementBlock("div", _hoisted_16, [
          createBaseVNode("div", _hoisted_17, [
            autoBusy.value ? (openBlock(), createBlock(unref(LoaderCircle), {
              key: 0,
              class: "h-3.5 w-3.5 dh-spin text-accent"
            })) : (openBlock(), createBlock(unref(Zap), {
              key: 1,
              class: "h-3.5 w-3.5 text-text-4"
            })),
            createBaseVNode("span", null, toDisplayString(autoBusy.value ? "更新进行中" : "最近一次更新进度"), 1),
            createBaseVNode("span", _hoisted_18, toDisplayString(progress.value.length) + " 条", 1)
          ]),
          createBaseVNode("div", _hoisted_19, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(progress.value, (p, i) => {
              return openBlock(), createElementBlock("div", {
                key: i,
                class: "flex items-center gap-2.5 border-b border-line-row px-3.5 py-2 text-[12px] last:border-b-0"
              }, [
                createBaseVNode("span", {
                  class: normalizeClass(["h-[6px] w-[6px] flex-none rounded-full", p.status === "failed" ? "bg-err" : p.status === "success" ? "bg-run" : "bg-accent"])
                }, null, 2),
                createBaseVNode("span", _hoisted_20, toDisplayString(p.name), 1),
                createBaseVNode("span", _hoisted_21, toDisplayString(p.message), 1)
              ]);
            }), 128))
          ])
        ])) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_22, [
          createBaseVNode("div", _hoisted_23, [
            createVNode(unref(Zap), {
              class: normalizeClass(["h-3.5 w-3.5", auto.value?.enabled ? "text-warn-text" : "text-text-4"])
            }, null, 8, ["class"]),
            _cache[12] || (_cache[12] = createBaseVNode("span", null, "自动更新", -1)),
            createVNode(unref(RouterLink), {
              to: "/settings",
              class: normalizeClass(["dh-tap dh-badge", auto.value?.enabled ? "dh-badge-warn" : "dh-badge-plain"]),
              title: auto.value?.enabled ? "已开启 · 去「设置 → 更新与检测」调整" : "已关闭 · 去「设置 → 更新与检测」打开"
            }, {
              default: withCtx(() => [
                createTextVNode(toDisplayString(auto.value?.enabled ? "已开启" : "已关闭"), 1)
              ]),
              _: 1
            }, 8, ["class", "title"]),
            createBaseVNode("div", _hoisted_24, [
              createBaseVNode("span", {
                class: "flex items-center gap-1.5 text-[11.5px] text-text-5",
                title: auto.value?.lastCheckAt ? `上次巡检 ${unref(formatDateTime)(auto.value.lastCheckAt)}` : "还没有巡检记录"
              }, [
                createVNode(unref(Clock), { class: "h-3 w-3" }),
                createTextVNode(" " + toDisplayString(auto.value?.nextCheckAt ? `下次巡检 ${unref(formatDateTime)(auto.value.nextCheckAt)}` : "下次巡检尚未排期") + " ", 1),
                auto.value?.checkIntervalHours ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                  createTextVNode("· 每 " + toDisplayString(auto.value.checkIntervalHours) + " 小时一次", 1)
                ], 64)) : createCommentVNode("", true),
                auto.value?.lastCheckAt ? (openBlock(), createElementBlock(Fragment, { key: 1 }, [
                  createTextVNode("· 上次 " + toDisplayString(unref(relativeTime)(auto.value.lastCheckAt)), 1)
                ], 64)) : createCommentVNode("", true)
              ], 8, _hoisted_25),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-primary",
                disabled: autoBusy.value || !auto.value?.enabled,
                onClick: runAutoUpdate
              }, [
                createVNode(unref(Play), { class: "h-3 w-3" }),
                _cache[11] || (_cache[11] = createTextVNode("立即执行自动更新 ", -1))
              ], 8, _hoisted_26)
            ])
          ]),
          createBaseVNode("div", _hoisted_27, [
            createBaseVNode("div", _hoisted_28, [
              createBaseVNode("div", _hoisted_29, [
                _cache[14] || (_cache[14] = createBaseVNode("span", { class: "text-[12px] font-medium text-text-3" }, "本轮会发生什么", -1)),
                createBaseVNode("span", _hoisted_30, toDisplayString(auto.value?.enabled ? `${willUpdate.value.length} 个容器将被更新` : `${willUpdate.value.length} 个容器有可用更新`), 1),
                auto.value?.enabled === false ? (openBlock(), createElementBlock("span", _hoisted_31, "自动更新已关闭 · 手动更新请进容器页")) : createCommentVNode("", true),
                createBaseVNode("span", _hoisted_32, toDisplayString(skipped.value.length) + " 个容器被保护/排除", 1),
                createVNode(unref(RouterLink), {
                  to: "/settings",
                  class: "dh-tap-txt ml-auto text-[11.5px] text-text-5 hover:text-accent"
                }, {
                  default: withCtx(() => [..._cache[13] || (_cache[13] = [
                    createTextVNode(" 去设置里调整检测频率与排除列表 → ", -1)
                  ])]),
                  _: 1
                })
              ]),
              !auto.value?.candidates?.length ? (openBlock(), createElementBlock("div", _hoisted_33, " 还没有巡检结果，点顶栏右上角的「检查更新」先跑一次。 ")) : (openBlock(), createElementBlock("div", _hoisted_34, [
                createBaseVNode("table", _hoisted_35, [
                  _cache[15] || (_cache[15] = createBaseVNode("thead", null, [
                    createBaseVNode("tr", null, [
                      createBaseVNode("th", { class: "w-[200px]" }, "容器"),
                      createBaseVNode("th", { class: "w-[220px]" }, "镜像"),
                      createBaseVNode("th", { class: "w-[110px]" }, "状态"),
                      createBaseVNode("th", null, "去向")
                    ])
                  ], -1)),
                  createBaseVNode("tbody", null, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(auto.value.candidates, (c) => {
                      return openBlock(), createElementBlock("tr", {
                        key: c.name
                      }, [
                        createBaseVNode("td", null, [
                          createVNode(unref(RouterLink), {
                            to: `/containers/${encodeURIComponent(c.name)}`,
                            class: "dh-tap-txt font-mono text-[11.5px] text-text-3 hover:text-accent"
                          }, {
                            default: withCtx(() => [
                              createTextVNode(toDisplayString(c.name), 1)
                            ]),
                            _: 2
                          }, 1032, ["to"])
                        ]),
                        createBaseVNode("td", _hoisted_36, toDisplayString(unref(shortImage)(c.image)), 1),
                        createBaseVNode("td", null, [
                          createBaseVNode("span", {
                            class: normalizeClass(["dh-badge", c.running ? "dh-badge-run" : "dh-badge-stop"])
                          }, toDisplayString(c.running ? "运行中" : "已停止"), 3)
                        ]),
                        createBaseVNode("td", _hoisted_37, [
                          c.willUpdate ? (openBlock(), createElementBlock("span", _hoisted_38, toDisplayString(auto.value?.enabled ? "将更新" : "有可用更新"), 1)) : c.protected ? (openBlock(), createElementBlock("span", _hoisted_39, "受保护")) : c.excluded ? (openBlock(), createElementBlock("span", _hoisted_40, "已排除")) : resultByName.value.get(c.name) ? (openBlock(), createElementBlock("span", {
                            key: 3,
                            class: normalizeClass(["dh-badge", `dh-badge-${toneOf(resultByName.value.get(c.name).status)}`])
                          }, toDisplayString(unref(checkLabel)(resultByName.value.get(c.name).status).text), 3)) : (openBlock(), createElementBlock("span", _hoisted_41, "未检测到更新")),
                          createBaseVNode("span", _hoisted_42, toDisplayString(resultByName.value.get(c.name)?.reason || c.reason), 1)
                        ])
                      ]);
                    }), 128))
                  ])
                ])
              ]))
            ])
          ])
        ]),
        createVNode(_sfc_main$1, {
          open: showInfo.value,
          title: "为什么 Dockhelm 不会误停容器",
          width: "620px",
          onClose: _cache[2] || (_cache[2] = ($event) => showInfo.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: _cache[1] || (_cache[1] = ($event) => showInfo.value = false)
            }, "明白了")
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_43, [
              createBaseVNode("div", _hoisted_44, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none text-warn-text" }),
                _cache[16] || (_cache[16] = createBaseVNode("div", null, [
                  createBaseVNode("b", { class: "text-text-1" }, "dockerCopilot 的两个缺陷。"),
                  createTextVNode(" 一是它的检测路径自拼 registry 请求、用一份硬编码的加速站列表，从不读取守护进程 "),
                  createBaseVNode("code", { class: "text-text-2" }, "daemon.json"),
                  createTextVNode(" 里真正生效的 "),
                  createBaseVNode("code", { class: "text-text-2" }, "registry-mirrors"),
                  createTextVNode("， 于是「检测说有新版本、拉取说已是最新」会永久互相矛盾。 二是它的更新流程无条件执行 "),
                  createBaseVNode("code", { class: "text-text-2" }, "pull → stop → rename → create → start"),
                  createTextVNode("， "),
                  createBaseVNode("b", { class: "text-text-1" }, "算出了新镜像 ID 却只用来决定要不要删旧镜像，从不用于决定要不要停容器"),
                  createTextVNode("。 两者叠加，就出现了「一次更新十几个，连没更新的容器也全被停掉」。 ")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_45, [
                createVNode(unref(CircleCheck), { class: "mt-[2px] h-4 w-4 flex-none text-run-text" }),
                _cache[17] || (_cache[17] = createBaseVNode("div", null, [
                  createBaseVNode("b", { class: "text-text-1" }, "Dockhelm 的做法。"),
                  createTextVNode(" 检测走守护进程的 "),
                  createBaseVNode("code", { class: "text-text-2" }, "/distribution"),
                  createTextVNode(" 接口（与 pull 同一套仓库端点解析）， 检测失败一律标记为「无法判定」而不是「有新版本」。 更新时先拉取，再比对容器使用的镜像 ID："),
                  createBaseVNode("b", { class: "text-text-1" }, "一致就直接返回，容器一个字节都不碰"),
                  createTextVNode("； 不一致才停止旧容器、改名保留（"),
                  createBaseVNode("code", { class: "text-text-2" }, "<名字>__bak_<时间>"),
                  createTextVNode("）、 用原配置创建同名新容器、启动并做健康检查，失败自动把旧容器改回原名并启动。 ")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_46, [
                createVNode(unref(Info), { class: "mt-[2px] h-4 w-4 flex-none text-accent" }),
                _cache[18] || (_cache[18] = createBaseVNode("div", null, [
                  createBaseVNode("b", { class: "text-text-1" }, "检测是只读的。"),
                  createTextVNode(" 检测只下 manifest、比对摘要，"),
                  createBaseVNode("b", { class: "text-text-1" }, "不拉层、不停容器"),
                  createTextVNode("； 真正会动容器的只有「立即执行自动更新」与容器详情页的更新。 ")
                ], -1))
              ])
            ])
          ]),
          _: 1
        }, 8, ["open"])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
