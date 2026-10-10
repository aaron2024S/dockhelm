import { x as createLucideIcon, d as defineComponent, u as useAppStore, y as watch, a as onMounted, o as onUnmounted, c as createElementBlock, f as createVNode, g as unref, e as createBaseVNode, t as toDisplayString, j as createCommentVNode, k as createTextVNode, F as Fragment, n as normalizeClass, z as normalizeStyle, r as renderList, q as useRouter, m as ref, p as computed, s as openBlock, B as api, C as openStream } from "./index-BoPP5Zol.js";
import { f as formatDayTime, a as formatBytes, r as runKindLabel } from "./format-NZu8lNOP.js";
import { T as TriangleAlert } from "./triangle-alert-BwT196BN.js";
const ArrowRight = createLucideIcon("ArrowRightIcon", [
  ["path", { d: "M5 12h14", key: "1ays0h" }],
  ["path", { d: "m12 5 7 7-7 7", key: "xquz4c" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = {
  key: 0,
  class: "dh-banner dh-banner-err"
};
const _hoisted_3 = { class: "min-w-0 flex-1" };
const _hoisted_4 = {
  key: 1,
  class: "dh-banner dh-banner-err"
};
const _hoisted_5 = { class: "min-w-0 flex-1" };
const _hoisted_6 = { class: "grid grid-cols-2 gap-3 lg:grid-cols-4" };
const _hoisted_7 = { class: "dh-metric" };
const _hoisted_8 = { class: "v" };
const _hoisted_9 = { class: "s" };
const _hoisted_10 = { class: "dh-metric" };
const _hoisted_11 = { class: "s" };
const _hoisted_12 = { class: "dh-metric" };
const _hoisted_13 = { class: "v" };
const _hoisted_14 = { class: "s" };
const _hoisted_15 = { class: "dh-metric" };
const _hoisted_16 = { class: "v" };
const _hoisted_17 = { class: "s" };
const _hoisted_18 = { class: "grid grid-cols-1 gap-3.5 xl:grid-cols-2" };
const _hoisted_19 = { class: "dh-card" };
const _hoisted_20 = { class: "dh-card-body flex items-center gap-[22px]" };
const _hoisted_21 = {
  viewBox: "0 0 120 120",
  class: "h-[118px] w-[118px] flex-none"
};
const _hoisted_22 = ["stroke-dasharray"];
const _hoisted_23 = ["stroke-dasharray", "stroke-dashoffset"];
const _hoisted_24 = {
  x: "60",
  y: "56",
  "text-anchor": "middle",
  style: { "fill": "var(--color-text-1)" },
  "font-size": "21",
  "font-weight": "600"
};
const _hoisted_25 = { class: "flex min-w-0 flex-1 flex-col gap-[11px]" };
const _hoisted_26 = { class: "flex gap-4 text-[12px] text-text-4" };
const _hoisted_27 = { class: "flex items-center gap-1.5" };
const _hoisted_28 = { class: "flex items-center gap-1.5" };
const _hoisted_29 = { class: "flex flex-col gap-1.5" };
const _hoisted_30 = { class: "flex justify-between text-[11.5px] text-text-4" };
const _hoisted_31 = { class: "dh-bar" };
const _hoisted_32 = { class: "flex flex-col gap-1.5" };
const _hoisted_33 = { class: "flex justify-between text-[11.5px] text-text-4" };
const _hoisted_34 = { class: "dh-bar" };
const _hoisted_35 = {
  key: 0,
  class: "flex flex-col gap-1.5"
};
const _hoisted_36 = { class: "flex justify-between text-[11.5px] text-text-4" };
const _hoisted_37 = { class: "dh-bar" };
const _hoisted_38 = { class: "dh-card flex flex-col" };
const _hoisted_39 = { class: "dh-card-head" };
const _hoisted_40 = {
  key: 0,
  class: "dh-card-body flex flex-1 flex-col items-center justify-center gap-2 py-6 text-center"
};
const _hoisted_41 = { class: "text-[12.5px] text-text-3" };
const _hoisted_42 = { class: "text-[11.5px] leading-relaxed text-text-5" };
const _hoisted_43 = {
  key: 1,
  class: "dh-card-body flex flex-col gap-[11px]"
};
const _hoisted_44 = { class: "grid h-[28px] w-[28px] flex-none place-items-center rounded-[9px] bg-line-2 text-[11px] font-semibold text-accent-text" };
const _hoisted_45 = { class: "min-w-0 flex-1" };
const _hoisted_46 = ["title"];
const _hoisted_47 = ["title"];
const _hoisted_48 = { class: "dh-card" };
const _hoisted_49 = {
  key: 0,
  class: "dh-card-body text-[12.5px] text-text-4"
};
const _hoisted_50 = {
  key: 1,
  class: "dh-card-body flex flex-col py-1"
};
const _hoisted_51 = { class: "w-[86px] flex-none text-[11.5px] text-text-5" };
const _hoisted_52 = ["title"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "OverviewView",
  setup(__props) {
    const app = useAppStore();
    const router = useRouter();
    const data = ref(null);
    const schedules = ref([]);
    const loading = ref(true);
    const errorMsg = ref("");
    let closeStream = null;
    let timers = [];
    function later(fn, ms) {
      timers.push(window.setTimeout(fn, ms));
    }
    async function load() {
      loading.value = true;
      errorMsg.value = "";
      try {
        const [ov, sch] = await Promise.all([
          api.get("/api/overview"),
          api.get("/api/schedules")
        ]);
        data.value = ov;
        schedules.value = sch.schedules ?? [];
      } catch (e) {
        errorMsg.value = e instanceof Error ? e.message : String(e);
      } finally {
        loading.value = false;
      }
    }
    let usageTimer;
    async function pollUsage() {
      if (document.hidden) return;
      try {
        const u = await api.get("/api/usage");
        if (!data.value) return;
        data.value.usage = { cpuPercent: u.cpuPercent, memUsed: u.memUsed, memTotal: u.memTotal };
        if (u.diskTotal) data.value.disk = { free: u.diskFree, total: u.diskTotal };
      } catch {
      }
    }
    function sizeParts(n) {
      const [v = "0", u = "B"] = formatBytes(n).split(" ");
      return { v, u };
    }
    const donut = computed(() => {
      const total = data.value?.containers?.total ?? 0;
      const running = data.value?.containers?.running ?? 0;
      const C = 2 * Math.PI * 46;
      const runArc = total ? running / total * C : 0;
      const stopArc = total ? (total - running) / total * C : 0;
      return { C, runArc, stopArc };
    });
    const memTotal = computed(
      () => data.value?.usage?.memTotal || data.value?.docker?.memTotal || 0
    );
    const memUsed = computed(() => data.value?.usage?.memUsed ?? 0);
    const diskUsed = computed(() => {
      const d = data.value?.disk;
      if (!d?.total) return 0;
      return d.total - d.free;
    });
    const pct = (used, total) => total > 0 ? Math.min(used / total * 100, 100) : 0;
    const cpuPercent = computed(() => Math.min(data.value?.usage?.cpuPercent ?? 0, 100));
    const pendingCount = computed(() => data.value?.updates?.items?.length ?? 0);
    const updateGroups = computed(() => {
      const map = /* @__PURE__ */ new Map();
      for (const it of data.value?.updates?.items ?? []) {
        const entry = map.get(it.image);
        if (entry) {
          entry.containers.push(it.container);
        } else {
          map.set(it.image, {
            image: it.image,
            containers: [it.container],
            localDigest: it.localDigest,
            remoteDigest: it.remoteDigest
          });
        }
      }
      return [...map.values()];
    });
    const nextRun = computed(() => {
      const times = schedules.value.filter((s) => s.enabled && s.nextRun).map((s) => Date.parse(s.nextRun)).filter((t) => !Number.isNaN(t)).sort((a, b) => a - b);
      return times.length ? times[0] : 0;
    });
    const recentRows = computed(() => (data.value?.recent ?? []).slice(0, 8));
    function initial(name) {
      return (name[0] ?? "?").toUpperCase();
    }
    function imageShort(image) {
      const noTag = (image.split("@")[0] ?? image).split("/").pop() ?? image;
      const colon = noTag.indexOf(":");
      return colon > 0 ? noTag.slice(0, colon) : noTag;
    }
    function kindClass(kind, status) {
      if (status === "failed") return "dh-badge-err";
      if (kind === "schedule") return "dh-badge-accent";
      return "dh-badge-plain";
    }
    function resultBadge(l) {
      if (l.status === "success") return { text: "成功", cls: "dh-badge-run" };
      if (l.status === "up_to_date") return { text: "跳过", cls: "dh-badge-accent" };
      if (l.status === "failed") return { text: "失败", cls: "dh-badge-err" };
      if (l.status === "done") return { text: "完成", cls: "dh-badge-run" };
      return { text: l.status || "—", cls: "dh-badge-plain" };
    }
    function logTime(ts) {
      const t = Date.parse(ts);
      if (Number.isNaN(t)) return "—";
      const d = new Date(t);
      const p = (n) => String(n).padStart(2, "0");
      const hm = `${p(d.getHours())}:${p(d.getMinutes())}`;
      const today0 = /* @__PURE__ */ new Date();
      today0.setHours(0, 0, 0, 0);
      const day0 = new Date(t);
      day0.setHours(0, 0, 0, 0);
      const diff = Math.round((day0.getTime() - today0.getTime()) / 864e5);
      if (diff === 0) return hm;
      if (diff === -1) return `昨天 ${hm}`;
      return formatDayTime(ts, "date");
    }
    watch(
      () => app.checkTick,
      () => void load()
    );
    onMounted(() => {
      void load();
      usageTimer = window.setInterval(() => void pollUsage(), 1e4);
      closeStream = openStream("/api/events/stream", (topic, ev) => {
        if (topic !== "update" && topic !== "schedule") return;
        if (ev.kind === "batch_done" || ev.kind === "check_done") {
          later(() => void load(), 800);
        }
      });
    });
    onUnmounted(() => {
      closeStream?.();
      if (usageTimer !== void 0) window.clearInterval(usageTimer);
      timers.forEach((t) => window.clearTimeout(t));
      timers = [];
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        errorMsg.value ? (openBlock(), createElementBlock("div", _hoisted_2, [
          createVNode(unref(TriangleAlert), { class: "h-4 w-4 flex-none" }),
          createBaseVNode("span", _hoisted_3, toDisplayString(errorMsg.value), 1),
          createBaseVNode("button", {
            class: "dh-btn dh-btn-sm",
            onClick: load
          }, "重试")
        ])) : createCommentVNode("", true),
        data.value?.dockerError ? (openBlock(), createElementBlock("div", _hoisted_4, [
          createVNode(unref(TriangleAlert), { class: "h-4 w-4 flex-none" }),
          createBaseVNode("span", _hoisted_5, "无法连接 Docker 守护进程：" + toDisplayString(data.value.dockerError), 1)
        ])) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_6, [
          createBaseVNode("div", _hoisted_7, [
            _cache[2] || (_cache[2] = createBaseVNode("div", { class: "k" }, "容器", -1)),
            createBaseVNode("div", _hoisted_8, [
              createTextVNode(toDisplayString(data.value?.containers?.total ?? "—") + " ", 1),
              _cache[1] || (_cache[1] = createBaseVNode("small", null, "个", -1))
            ]),
            createBaseVNode("div", _hoisted_9, [
              createTextVNode(" 运行 " + toDisplayString(data.value?.containers?.running ?? 0) + " · 停止 " + toDisplayString(data.value?.containers?.stopped ?? 0) + " ", 1),
              data.value?.containers?.unhealthy ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                createTextVNode(" · 异常 " + toDisplayString(data.value.containers.unhealthy), 1)
              ], 64)) : createCommentVNode("", true)
            ])
          ]),
          createBaseVNode("div", _hoisted_10, [
            _cache[4] || (_cache[4] = createBaseVNode("div", { class: "k" }, "待更新镜像", -1)),
            createBaseVNode("div", {
              class: normalizeClass(["v", updateGroups.value.length ? "text-warn-text" : ""])
            }, [
              createTextVNode(toDisplayString(updateGroups.value.length) + " ", 1),
              _cache[3] || (_cache[3] = createBaseVNode("small", null, "个", -1))
            ], 2),
            createBaseVNode("div", _hoisted_11, toDisplayString(data.value?.updates?.checkedAt ? `影响 ${data.value.updates.items?.length ?? 0} 个容器` : "还没做过巡检"), 1)
          ]),
          createBaseVNode("div", _hoisted_12, [
            _cache[6] || (_cache[6] = createBaseVNode("div", { class: "k" }, "计划任务", -1)),
            createBaseVNode("div", _hoisted_13, [
              createTextVNode(toDisplayString(schedules.value.length) + " ", 1),
              _cache[5] || (_cache[5] = createBaseVNode("small", null, "个", -1))
            ]),
            createBaseVNode("div", _hoisted_14, "下次 " + toDisplayString(nextRun.value ? unref(formatDayTime)(nextRun.value) : "暂无启用中的任务"), 1)
          ]),
          createBaseVNode("div", _hoisted_15, [
            _cache[7] || (_cache[7] = createBaseVNode("div", { class: "k" }, "镜像占用", -1)),
            createBaseVNode("div", _hoisted_16, [
              createTextVNode(toDisplayString(sizeParts(data.value?.images?.sizeBytes).v) + " ", 1),
              createBaseVNode("small", null, toDisplayString(sizeParts(data.value?.images?.sizeBytes).u), 1)
            ]),
            createBaseVNode("div", _hoisted_17, "可回收 " + toDisplayString(unref(formatBytes)(data.value?.images?.reclaimable)), 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_18, [
          createBaseVNode("div", _hoisted_19, [
            _cache[13] || (_cache[13] = createBaseVNode("div", { class: "dh-card-head" }, [
              createBaseVNode("span", null, "容器状态"),
              createBaseVNode("span", { class: "ml-auto dh-badge dh-badge-plain" }, "实时")
            ], -1)),
            createBaseVNode("div", _hoisted_20, [
              (openBlock(), createElementBlock("svg", _hoisted_21, [
                _cache[8] || (_cache[8] = createBaseVNode("circle", {
                  cx: "60",
                  cy: "60",
                  r: "46",
                  fill: "none",
                  style: { "stroke": "var(--color-line-1)" },
                  "stroke-width": "13"
                }, null, -1)),
                createBaseVNode("circle", {
                  cx: "60",
                  cy: "60",
                  r: "46",
                  fill: "none",
                  style: { "stroke": "var(--color-run)" },
                  "stroke-width": "13",
                  "stroke-linecap": "round",
                  "stroke-dasharray": `${donut.value.runArc} ${donut.value.C}`,
                  transform: "rotate(-90 60 60)"
                }, null, 8, _hoisted_22),
                donut.value.stopArc > 0 ? (openBlock(), createElementBlock("circle", {
                  key: 0,
                  cx: "60",
                  cy: "60",
                  r: "46",
                  fill: "none",
                  style: { "stroke": "var(--color-stop)" },
                  "stroke-width": "13",
                  "stroke-dasharray": `${donut.value.stopArc} ${donut.value.C}`,
                  "stroke-dashoffset": -donut.value.runArc,
                  transform: "rotate(-90 60 60)"
                }, null, 8, _hoisted_23)) : createCommentVNode("", true),
                createBaseVNode("text", _hoisted_24, toDisplayString(data.value?.containers?.running ?? 0), 1),
                _cache[9] || (_cache[9] = createBaseVNode("text", {
                  x: "60",
                  y: "74",
                  "text-anchor": "middle",
                  style: { "fill": "var(--color-text-4)" },
                  "font-size": "11"
                }, " 运行中 ", -1))
              ])),
              createBaseVNode("div", _hoisted_25, [
                createBaseVNode("div", _hoisted_26, [
                  createBaseVNode("span", _hoisted_27, [
                    _cache[10] || (_cache[10] = createBaseVNode("i", { class: "h-[7px] w-[7px] rounded-full bg-run" }, null, -1)),
                    createTextVNode("运行中 " + toDisplayString(data.value?.containers?.running ?? 0), 1)
                  ]),
                  createBaseVNode("span", _hoisted_28, [
                    _cache[11] || (_cache[11] = createBaseVNode("i", { class: "h-[7px] w-[7px] rounded-full bg-stop" }, null, -1)),
                    createTextVNode("已停止 " + toDisplayString(data.value?.containers?.stopped ?? 0), 1)
                  ])
                ]),
                createBaseVNode("div", _hoisted_29, [
                  createBaseVNode("div", _hoisted_30, [
                    _cache[12] || (_cache[12] = createBaseVNode("span", null, "CPU 总占用", -1)),
                    createBaseVNode("span", null, toDisplayString(cpuPercent.value.toFixed(0)) + "%", 1)
                  ]),
                  createBaseVNode("div", _hoisted_31, [
                    createBaseVNode("i", {
                      style: normalizeStyle({ width: `${cpuPercent.value}%` })
                    }, null, 4)
                  ])
                ]),
                createBaseVNode("div", _hoisted_32, [
                  createBaseVNode("div", _hoisted_33, [
                    createBaseVNode("span", null, "内存 " + toDisplayString(unref(formatBytes)(memUsed.value)) + " / " + toDisplayString(unref(formatBytes)(memTotal.value)), 1),
                    createBaseVNode("span", null, toDisplayString(pct(memUsed.value, memTotal.value).toFixed(0)) + "%", 1)
                  ]),
                  createBaseVNode("div", _hoisted_34, [
                    createBaseVNode("i", {
                      style: normalizeStyle({ width: `${pct(memUsed.value, memTotal.value)}%`, background: "var(--color-accent-text)" })
                    }, null, 4)
                  ])
                ]),
                data.value?.disk?.total ? (openBlock(), createElementBlock("div", _hoisted_35, [
                  createBaseVNode("div", _hoisted_36, [
                    createBaseVNode("span", null, "磁盘 " + toDisplayString(unref(formatBytes)(diskUsed.value)) + " / " + toDisplayString(unref(formatBytes)(data.value?.disk?.total)), 1),
                    createBaseVNode("span", null, toDisplayString(pct(diskUsed.value, data.value?.disk?.total ?? 0).toFixed(0)) + "%", 1)
                  ]),
                  createBaseVNode("div", _hoisted_37, [
                    createBaseVNode("i", {
                      style: normalizeStyle({ width: `${pct(diskUsed.value, data.value?.disk?.total ?? 0)}%`, background: "var(--color-chart-indigo)" })
                    }, null, 4)
                  ])
                ])) : createCommentVNode("", true)
              ])
            ])
          ]),
          createBaseVNode("div", _hoisted_38, [
            createBaseVNode("div", _hoisted_39, [
              _cache[14] || (_cache[14] = createBaseVNode("span", null, "待更新镜像", -1)),
              createBaseVNode("span", {
                class: normalizeClass(["ml-auto dh-badge", pendingCount.value ? "dh-badge-warn" : "dh-badge-plain"])
              }, toDisplayString(pendingCount.value) + " 个 ", 3)
            ]),
            !pendingCount.value ? (openBlock(), createElementBlock("div", _hoisted_40, [
              createBaseVNode("div", _hoisted_41, toDisplayString(data.value?.updates?.checkedAt ? "所有镜像都是最新的" : "还没有做过巡检"), 1),
              createBaseVNode("div", _hoisted_42, toDisplayString(data.value?.updates?.checkedAt ? "没有任何容器需要更新。" : "点顶栏右上角的「检查更新」可以只读地检查一遍所有容器。"), 1)
            ])) : (openBlock(), createElementBlock("div", _hoisted_43, [
              (openBlock(true), createElementBlock(Fragment, null, renderList((data.value?.updates?.items ?? []).slice(0, 4), (it) => {
                return openBlock(), createElementBlock("div", {
                  key: it.container,
                  class: "flex items-center gap-2.5"
                }, [
                  createBaseVNode("div", _hoisted_44, toDisplayString(initial(it.container)), 1),
                  createBaseVNode("div", _hoisted_45, [
                    createBaseVNode("div", {
                      class: "truncate text-[12.5px] font-semibold",
                      title: it.container
                    }, toDisplayString(it.container), 9, _hoisted_46),
                    createBaseVNode("div", {
                      class: "truncate text-[11.5px] text-text-5",
                      title: it.image
                    }, toDisplayString(imageShort(it.image)), 9, _hoisted_47)
                  ])
                ]);
              }), 128)),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-primary mt-0.5",
                onClick: _cache[0] || (_cache[0] = ($event) => unref(router).push("/containers"))
              }, [
                _cache[15] || (_cache[15] = createTextVNode(" 前往容器页面 ", -1)),
                createVNode(unref(ArrowRight), { class: "h-3.5 w-3.5" })
              ])
            ]))
          ])
        ]),
        createBaseVNode("div", _hoisted_48, [
          _cache[16] || (_cache[16] = createBaseVNode("div", { class: "dh-card-head" }, [
            createBaseVNode("span", null, "最近执行记录"),
            createBaseVNode("span", { class: "ml-auto text-[12px] font-normal text-text-4" }, "最近 15 条")
          ], -1)),
          !recentRows.value.length ? (openBlock(), createElementBlock("div", _hoisted_49, "暂无记录")) : (openBlock(), createElementBlock("div", _hoisted_50, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(recentRows.value, (l) => {
              return openBlock(), createElementBlock("div", {
                key: l.id,
                class: "dh-tl"
              }, [
                createBaseVNode("div", _hoisted_51, toDisplayString(logTime(l.ts)), 1),
                createBaseVNode("span", {
                  class: normalizeClass(["dh-badge flex-none", kindClass(l.kind, l.status)])
                }, toDisplayString(unref(runKindLabel)(l.kind)), 3),
                createBaseVNode("span", {
                  class: "min-w-0 flex-1 truncate text-[12.5px] text-text-2",
                  title: l.message
                }, toDisplayString(l.message), 9, _hoisted_52),
                createBaseVNode("span", {
                  class: normalizeClass(["dh-badge flex-none", resultBadge(l).cls])
                }, toDisplayString(resultBadge(l).text), 3)
              ]);
            }), 128))
          ]))
        ])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
