import { q as createLucideIcon, d as defineComponent, C as useToastStore, o as onMounted, s as onUnmounted, c as createElementBlock, b as createBaseVNode, f as unref, e as createVNode, t as toDisplayString, i as createTextVNode, n as normalizeClass, E as Download, h as createCommentVNode, w as withDirectives, M as vModelSelect, R as RefreshCw, y as createBlock, N as resolveDynamicComponent, F as Fragment, r as renderList, H as Box, m as useRouter, l as computed, k as ref, z as api, p as openBlock, J as useRoute, B as openStream } from "./index-HPgVGrf3.js";
import { r as relativeTime } from "./format-mcaWYSwR.js";
import { S as Square, R as RotateCw, _ as _sfc_main$1 } from "./PortChips.vue_vue_type_script_setup_true_lang-BW6YhVU7.js";
import { P as Play } from "./play-BDui4TlS.js";
import { H as HardDrive } from "./hard-drive-C4eHO2oq.js";
import { E as Eye } from "./eye-DjjlCgHi.js";
import { S as ScrollText } from "./scroll-text-tdGP435f.js";
const ArrowLeft = createLucideIcon("ArrowLeftIcon", [
  ["path", { d: "m12 19-7-7 7-7", key: "1l729n" }],
  ["path", { d: "M19 12H5", key: "x3x0zl" }]
]);
const EyeOff = createLucideIcon("EyeOffIcon", [
  [
    "path",
    {
      d: "M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49",
      key: "ct8e1f"
    }
  ],
  ["path", { d: "M14.084 14.158a3 3 0 0 1-4.242-4.242", key: "151rxh" }],
  [
    "path",
    {
      d: "M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143",
      key: "13bj9a"
    }
  ],
  ["path", { d: "m2 2 20 20", key: "1ooewy" }]
]);
const Network = createLucideIcon("NetworkIcon", [
  ["rect", { x: "16", y: "16", width: "6", height: "6", rx: "1", key: "4q2zg0" }],
  ["rect", { x: "2", y: "16", width: "6", height: "6", rx: "1", key: "8cvhb9" }],
  ["rect", { x: "9", y: "2", width: "6", height: "6", rx: "1", key: "1egb70" }],
  ["path", { d: "M5 16v-3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1v3", key: "1jsf9p" }],
  ["path", { d: "M12 12V8", key: "2874zd" }]
]);
const Terminal = createLucideIcon("TerminalIcon", [
  ["polyline", { points: "4 17 10 11 4 5", key: "akl6gq" }],
  ["line", { x1: "12", x2: "20", y1: "19", y2: "19", key: "q2wloq" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "grid h-[30px] w-[30px] flex-none place-items-center rounded-[9px] bg-line-2 text-[12px] font-semibold text-[#5eead4]" };
const _hoisted_4 = { class: "min-w-0" };
const _hoisted_5 = { class: "dh-h1 truncate" };
const _hoisted_6 = { class: "dh-sub truncate font-mono" };
const _hoisted_7 = { class: "ml-auto flex flex-wrap items-center gap-2" };
const _hoisted_8 = ["disabled"];
const _hoisted_9 = ["disabled"];
const _hoisted_10 = ["disabled"];
const _hoisted_11 = {
  key: 0,
  class: "dh-card flex flex-wrap items-center gap-x-3 gap-y-2 px-3.5 py-2.5"
};
const _hoisted_12 = { class: "ml-auto text-[11.5px] text-text-6" };
const _hoisted_13 = { class: "grid grid-cols-2 gap-3 lg:grid-cols-4" };
const _hoisted_14 = { class: "dh-card p-3.5" };
const _hoisted_15 = { class: "mt-1.5 flex flex-wrap items-center gap-1.5" };
const _hoisted_16 = { class: "mt-1.5 text-[11px] text-text-5" };
const _hoisted_17 = { class: "dh-card p-3.5" };
const _hoisted_18 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_19 = { class: "mt-1.5 text-[11px] text-text-5" };
const _hoisted_20 = { class: "dh-card p-3.5" };
const _hoisted_21 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_22 = { class: "mt-1.5 text-[11px] text-text-5" };
const _hoisted_23 = { class: "dh-card p-3.5" };
const _hoisted_24 = { class: "mt-1.5 text-[16px] font-semibold leading-tight" };
const _hoisted_25 = { class: "mt-[3px] text-[16px] font-semibold leading-tight text-text-3" };
const _hoisted_26 = { class: "dh-card" };
const _hoisted_27 = { class: "dh-card-head" };
const _hoisted_28 = { class: "dh-seg" };
const _hoisted_29 = ["data-on"];
const _hoisted_30 = ["data-on"];
const _hoisted_31 = ["data-on"];
const _hoisted_32 = ["data-on"];
const _hoisted_33 = ["data-on"];
const _hoisted_34 = {
  key: 0,
  class: "ml-auto flex items-center gap-2"
};
const _hoisted_35 = { key: 0 };
const _hoisted_36 = {
  key: 0,
  class: "grid h-[200px] place-items-center"
};
const _hoisted_37 = {
  key: 1,
  class: "dh-scroll max-h-[520px] overflow-auto whitespace-pre-wrap break-all px-3.5 py-3 font-mono text-[11.5px] leading-[1.65] text-text-2"
};
const _hoisted_38 = {
  key: 1,
  class: "grid grid-cols-1 gap-0 lg:grid-cols-2"
};
const _hoisted_39 = { class: "border-b border-line-1 p-3.5 lg:border-b-0 lg:border-r" };
const _hoisted_40 = { class: "grid grid-cols-[92px_1fr] gap-x-3 gap-y-2 text-[12px]" };
const _hoisted_41 = { class: "break-all font-mono text-[11.5px] text-text-2" };
const _hoisted_42 = { class: "font-mono text-[11.5px] text-text-2" };
const _hoisted_43 = { class: "break-all font-mono text-[11.5px] text-text-2" };
const _hoisted_44 = { class: "font-mono text-[11.5px] text-text-2" };
const _hoisted_45 = { class: "font-mono text-[11.5px] text-text-2" };
const _hoisted_46 = { class: "text-text-2" };
const _hoisted_47 = { class: "text-text-2" };
const _hoisted_48 = { class: "text-text-2" };
const _hoisted_49 = { class: "text-text-2" };
const _hoisted_50 = { class: "p-3.5" };
const _hoisted_51 = { class: "mb-2 flex items-center gap-2" };
const _hoisted_52 = { class: "text-[11px] text-text-6" };
const _hoisted_53 = { class: "dh-scroll max-h-[360px] overflow-auto" };
const _hoisted_54 = { class: "w-full" };
const _hoisted_55 = { class: "py-1.5 pr-3 align-top font-mono text-[11px] text-text-4" };
const _hoisted_56 = { class: "break-all py-1.5 font-mono text-[11px] text-text-2" };
const _hoisted_57 = {
  key: 0,
  class: "text-text-6"
};
const _hoisted_58 = {
  key: 2,
  class: "p-3.5"
};
const _hoisted_59 = {
  key: 0,
  class: "text-[12.5px] text-text-4"
};
const _hoisted_60 = {
  key: 1,
  class: "flex flex-col gap-2.5"
};
const _hoisted_61 = { class: "flex items-start gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2 text-[11.5px] leading-relaxed text-text-4" };
const _hoisted_62 = { class: "text-text-2" };
const _hoisted_63 = { class: "text-text-2" };
const _hoisted_64 = { class: "overflow-hidden rounded-[10px] border border-line-1" };
const _hoisted_65 = { class: "w-full text-[12px]" };
const _hoisted_66 = { class: "px-3 py-2 font-mono text-[11.5px] text-text-3" };
const _hoisted_67 = { class: "px-3 py-2 font-mono text-[11.5px] text-text-2" };
const _hoisted_68 = { class: "px-3 py-2 font-mono text-[11.5px] text-text-4" };
const _hoisted_69 = { class: "px-3 py-2" };
const _hoisted_70 = {
  key: 3,
  class: "p-3.5"
};
const _hoisted_71 = {
  key: 0,
  class: "text-[12.5px] text-text-4"
};
const _hoisted_72 = {
  key: 1,
  class: "flex flex-col gap-2.5"
};
const _hoisted_73 = { class: "flex items-start gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2 text-[11.5px] leading-relaxed text-text-4" };
const _hoisted_74 = { class: "flex flex-wrap items-center gap-2" };
const _hoisted_75 = { class: "ml-auto font-mono text-[11px] text-text-3" };
const _hoisted_76 = { class: "mt-1.5 break-all font-mono text-[11px] text-text-2" };
const _hoisted_77 = { class: "mt-1 text-[11px] text-text-6" };
const _hoisted_78 = {
  key: 4,
  class: "p-3.5"
};
const _hoisted_79 = {
  key: 0,
  class: "text-[12.5px] text-text-4"
};
const _hoisted_80 = {
  key: 1,
  class: "grid grid-cols-1 gap-2.5 md:grid-cols-2"
};
const _hoisted_81 = { class: "flex items-center gap-2" };
const _hoisted_82 = { class: "text-[12.5px] font-medium" };
const _hoisted_83 = { class: "mt-1.5 grid grid-cols-[62px_1fr] gap-x-3 gap-y-1 text-[11.5px]" };
const _hoisted_84 = { class: "font-mono text-text-2" };
const _hoisted_85 = { class: "font-mono text-text-2" };
const _hoisted_86 = { class: "break-all font-mono text-text-2" };
const _hoisted_87 = { class: "flex items-center gap-2 text-[11.5px] text-text-6" };
const _hoisted_88 = { class: "ml-auto inline-flex items-center gap-1" };
const _hoisted_89 = { class: "text-text-4" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "ContainerDetailView",
  setup(__props) {
    const route = useRoute();
    const router = useRouter();
    const toast = useToastStore();
    const name = computed(() => decodeURIComponent(String(route.params.name ?? "")));
    const detail = ref(null);
    const summary = ref(null);
    const stats = ref(null);
    const logs = ref("");
    const tab = ref("logs");
    const loading = ref(true);
    const busy = ref(false);
    const showSensitive = ref(false);
    const tail = ref(200);
    let timer;
    let closeStream = null;
    const container = computed(() => summary.value ?? {});
    const envList = computed(() => {
      const list = summary.value?.env ?? [];
      return list;
    });
    const mounts = computed(() => summary.value?.mounts ?? []);
    const networks = computed(() => summary.value?.networks ?? []);
    const portList = computed(() => summary.value?.ports ?? []);
    const publishedCount = computed(() => portList.value.filter((p) => p.published).length);
    async function load() {
      loading.value = true;
      try {
        const res = await api.get(
          `/api/containers/${encodeURIComponent(name.value)}`
        );
        detail.value = res.inspect;
        summary.value = res.summary;
      } catch (e) {
        toast.error("读取容器失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function loadLogs() {
      try {
        logs.value = await api.text(`/api/containers/${encodeURIComponent(name.value)}/logs`, { tail: tail.value });
      } catch (e) {
        logs.value = "读取日志失败：" + (e instanceof Error ? e.message : String(e));
      }
    }
    async function loadStats() {
      try {
        const res = await api.get(
          `/api/containers/${encodeURIComponent(name.value)}/stats`
        );
        stats.value = res.summary;
      } catch {
        stats.value = null;
      }
    }
    async function act(action) {
      busy.value = true;
      try {
        await api.post(`/api/containers/${encodeURIComponent(name.value)}/action`, { action });
        toast.success("操作已执行");
        await load();
        await loadStats();
      } catch (e) {
        toast.error("操作失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = false;
      }
    }
    async function snapshot() {
      try {
        await api.post("/api/backups/snapshot", { container: name.value, reason: "manual" });
        toast.success("已保存配置快照");
      } catch (e) {
        toast.error("备份失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function updateNow() {
      try {
        await api.post("/api/updates/apply", { names: [name.value] });
        toast.info("已提交更新任务", "镜像没变化的话容器不会被停止");
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
      }
    }
    const healthTone = computed(() => {
      const h = container.value.health;
      if (h === "healthy") return "dh-badge-run";
      if (h === "unhealthy") return "dh-badge-err";
      return "dh-badge-plain";
    });
    onMounted(async () => {
      await load();
      await loadLogs();
      await loadStats();
      timer = window.setInterval(() => {
        if (tab.value === "logs") void loadLogs();
        void loadStats();
      }, 8e3);
      closeStream = openStream("/api/events/stream", (topic, ev) => {
        if (topic === "container" && String(ev.data?.name ?? "") === name.value) void load();
      });
    });
    onUnmounted(() => {
      if (timer) window.clearInterval(timer);
      closeStream?.();
    });
    const running = computed(() => container.value.health !== void 0 && detail.value?.["State"]?.Running);
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          createBaseVNode("button", {
            class: "dh-iconbtn",
            title: "返回列表",
            onClick: _cache[0] || (_cache[0] = ($event) => unref(router).back())
          }, [
            createVNode(unref(ArrowLeft), { class: "h-4 w-4" })
          ]),
          createBaseVNode("div", _hoisted_3, toDisplayString(name.value.slice(0, 2).toUpperCase()), 1),
          createBaseVNode("div", _hoisted_4, [
            createBaseVNode("div", _hoisted_5, toDisplayString(name.value), 1),
            createBaseVNode("div", _hoisted_6, toDisplayString(container.value.id), 1)
          ]),
          createBaseVNode("div", _hoisted_7, [
            running.value ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn dh-btn-sm",
              disabled: busy.value,
              onClick: _cache[1] || (_cache[1] = ($event) => act("stop"))
            }, [
              createVNode(unref(Square), { class: "h-3 w-3" }),
              _cache[11] || (_cache[11] = createTextVNode("停止 ", -1))
            ], 8, _hoisted_8)) : (openBlock(), createElementBlock("button", {
              key: 1,
              class: "dh-btn dh-btn-sm",
              disabled: busy.value,
              onClick: _cache[2] || (_cache[2] = ($event) => act("start"))
            }, [
              createVNode(unref(Play), { class: "h-3 w-3" }),
              _cache[12] || (_cache[12] = createTextVNode("启动 ", -1))
            ], 8, _hoisted_9)),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              disabled: busy.value,
              onClick: _cache[3] || (_cache[3] = ($event) => act("restart"))
            }, [
              createVNode(unref(RotateCw), {
                class: normalizeClass(["h-3 w-3", busy.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[13] || (_cache[13] = createTextVNode("重启 ", -1))
            ], 8, _hoisted_10),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              onClick: updateNow
            }, [
              createVNode(unref(Download), { class: "h-3 w-3" }),
              _cache[14] || (_cache[14] = createTextVNode("检查并更新 ", -1))
            ]),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              onClick: snapshot
            }, [
              createVNode(unref(HardDrive), { class: "h-3 w-3" }),
              _cache[15] || (_cache[15] = createTextVNode("备份配置 ", -1))
            ])
          ])
        ]),
        portList.value.length ? (openBlock(), createElementBlock("div", _hoisted_11, [
          _cache[16] || (_cache[16] = createBaseVNode("span", { class: "text-[12px] text-text-4" }, "端口", -1)),
          createVNode(_sfc_main$1, {
            ports: portList.value,
            max: 0
          }, null, 8, ["ports"]),
          createBaseVNode("span", _hoisted_12, toDisplayString(publishedCount.value) + " 个已发布到宿主机 · " + toDisplayString(portList.value.length - publishedCount.value) + " 个仅容器内 ", 1)
        ])) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_13, [
          createBaseVNode("div", _hoisted_14, [
            _cache[17] || (_cache[17] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "状态", -1)),
            createBaseVNode("div", _hoisted_15, [
              createBaseVNode("span", {
                class: normalizeClass(["dh-badge", detail.value?.["State"] && detail.value["State"].Running ? "dh-badge-run" : "dh-badge-stop"])
              }, toDisplayString(detail.value?.["State"] && detail.value["State"].Running ? "运行中" : "已停止"), 3),
              container.value.health ? (openBlock(), createElementBlock("span", {
                key: 0,
                class: normalizeClass(["dh-badge", healthTone.value])
              }, toDisplayString(container.value.health), 3)) : createCommentVNode("", true)
            ]),
            createBaseVNode("div", _hoisted_16, " 退出码 " + toDisplayString(container.value.exitCode ?? "—") + " · 重启 " + toDisplayString(container.value.restarts ?? 0) + " 次 ", 1)
          ]),
          createBaseVNode("div", _hoisted_17, [
            _cache[19] || (_cache[19] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "CPU", -1)),
            createBaseVNode("div", _hoisted_18, [
              createTextVNode(toDisplayString(stats.value ? (stats.value.cpuPercent ?? 0).toFixed(1) : "—"), 1),
              _cache[18] || (_cache[18] = createBaseVNode("span", { class: "ml-0.5 text-[12px] text-text-5" }, "%", -1))
            ]),
            createBaseVNode("div", _hoisted_19, toDisplayString(stats.value?.onlineCpus ?? "—") + " 核可用", 1)
          ]),
          createBaseVNode("div", _hoisted_20, [
            _cache[21] || (_cache[21] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "内存", -1)),
            createBaseVNode("div", _hoisted_21, [
              createTextVNode(toDisplayString(stats.value ? (stats.value.memPercent ?? 0).toFixed(1) : "—"), 1),
              _cache[20] || (_cache[20] = createBaseVNode("span", { class: "ml-0.5 text-[12px] text-text-5" }, "%", -1))
            ]),
            createBaseVNode("div", _hoisted_22, toDisplayString(stats.value ? Math.round((stats.value.memUsage ?? 0) / 1048576) + " MB" : "—") + " / " + toDisplayString(stats.value ? Math.round((stats.value.memLimit ?? 0) / 1048576) + " MB" : "—"), 1)
          ]),
          createBaseVNode("div", _hoisted_23, [
            _cache[22] || (_cache[22] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "网络", -1)),
            createBaseVNode("div", _hoisted_24, " ↓ " + toDisplayString(stats.value ? Math.round((stats.value.netRx ?? 0) / 1048576) + " MB" : "—"), 1),
            createBaseVNode("div", _hoisted_25, " ↑ " + toDisplayString(stats.value ? Math.round((stats.value.netTx ?? 0) / 1048576) + " MB" : "—"), 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_26, [
          createBaseVNode("div", _hoisted_27, [
            createBaseVNode("div", _hoisted_28, [
              createBaseVNode("button", {
                "data-on": tab.value === "logs",
                onClick: _cache[4] || (_cache[4] = ($event) => tab.value = "logs")
              }, "日志", 8, _hoisted_29),
              createBaseVNode("button", {
                "data-on": tab.value === "config",
                onClick: _cache[5] || (_cache[5] = ($event) => tab.value = "config")
              }, "配置", 8, _hoisted_30),
              createBaseVNode("button", {
                "data-on": tab.value === "ports",
                onClick: _cache[6] || (_cache[6] = ($event) => tab.value = "ports")
              }, "端口 " + toDisplayString(portList.value.length), 9, _hoisted_31),
              createBaseVNode("button", {
                "data-on": tab.value === "mounts",
                onClick: _cache[7] || (_cache[7] = ($event) => tab.value = "mounts")
              }, "挂载 " + toDisplayString(mounts.value.length), 9, _hoisted_32),
              createBaseVNode("button", {
                "data-on": tab.value === "networks",
                onClick: _cache[8] || (_cache[8] = ($event) => tab.value = "networks")
              }, "网络 " + toDisplayString(networks.value.length), 9, _hoisted_33)
            ]),
            tab.value === "logs" ? (openBlock(), createElementBlock("div", _hoisted_34, [
              withDirectives(createBaseVNode("select", {
                "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => tail.value = $event),
                class: "dh-select !w-[110px] !py-[5px] !text-[11.5px]",
                onChange: loadLogs
              }, [..._cache[23] || (_cache[23] = [
                createBaseVNode("option", { value: 100 }, "最近 100 行", -1),
                createBaseVNode("option", { value: 200 }, "最近 200 行", -1),
                createBaseVNode("option", { value: 500 }, "最近 500 行", -1),
                createBaseVNode("option", { value: 2e3 }, "最近 2000 行", -1)
              ])], 544), [
                [
                  vModelSelect,
                  tail.value,
                  void 0,
                  { number: true }
                ]
              ]),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm",
                onClick: loadLogs
              }, [
                createVNode(unref(RefreshCw), { class: "h-3 w-3" }),
                _cache[24] || (_cache[24] = createTextVNode("刷新", -1))
              ])
            ])) : createCommentVNode("", true)
          ]),
          tab.value === "logs" ? (openBlock(), createElementBlock("div", _hoisted_35, [
            loading.value ? (openBlock(), createElementBlock("div", _hoisted_36, [
              createVNode(unref(RefreshCw), { class: "h-5 w-5 dh-spin text-text-5" })
            ])) : (openBlock(), createElementBlock("pre", _hoisted_37, toDisplayString(logs.value || "（没有日志输出）"), 1))
          ])) : tab.value === "config" ? (openBlock(), createElementBlock("div", _hoisted_38, [
            createBaseVNode("div", _hoisted_39, [
              _cache[34] || (_cache[34] = createBaseVNode("div", { class: "mb-2 text-[12px] font-medium text-text-3" }, "基本信息", -1)),
              createBaseVNode("div", _hoisted_40, [
                _cache[25] || (_cache[25] = createBaseVNode("div", { class: "text-text-5" }, "镜像", -1)),
                createBaseVNode("div", _hoisted_41, toDisplayString(container.value.image), 1),
                _cache[26] || (_cache[26] = createBaseVNode("div", { class: "text-text-5" }, "镜像 ID", -1)),
                createBaseVNode("div", _hoisted_42, toDisplayString(container.value.imageId), 1),
                _cache[27] || (_cache[27] = createBaseVNode("div", { class: "text-text-5" }, "启动命令", -1)),
                createBaseVNode("div", _hoisted_43, toDisplayString([...container.value.entrypoint ?? [], ...container.value.cmd ?? []].join(" ") || "（镜像默认）"), 1),
                _cache[28] || (_cache[28] = createBaseVNode("div", { class: "text-text-5" }, "工作目录", -1)),
                createBaseVNode("div", _hoisted_44, toDisplayString(container.value.workingDir || "—"), 1),
                _cache[29] || (_cache[29] = createBaseVNode("div", { class: "text-text-5" }, "运行用户", -1)),
                createBaseVNode("div", _hoisted_45, toDisplayString(container.value.user || "root"), 1),
                _cache[30] || (_cache[30] = createBaseVNode("div", { class: "text-text-5" }, "重启策略", -1)),
                createBaseVNode("div", _hoisted_46, toDisplayString(container.value.restart), 1),
                _cache[31] || (_cache[31] = createBaseVNode("div", { class: "text-text-5" }, "网络模式", -1)),
                createBaseVNode("div", _hoisted_47, toDisplayString(container.value.networkMode), 1),
                _cache[32] || (_cache[32] = createBaseVNode("div", { class: "text-text-5" }, "特权模式", -1)),
                createBaseVNode("div", _hoisted_48, toDisplayString(container.value.privileged ? "是" : "否"), 1),
                _cache[33] || (_cache[33] = createBaseVNode("div", { class: "text-text-5" }, "启动时间", -1)),
                createBaseVNode("div", _hoisted_49, toDisplayString(unref(relativeTime)(container.value.startedAt)), 1)
              ])
            ]),
            createBaseVNode("div", _hoisted_50, [
              createBaseVNode("div", _hoisted_51, [
                _cache[35] || (_cache[35] = createBaseVNode("span", { class: "text-[12px] font-medium text-text-3" }, "环境变量", -1)),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-ghost !p-1",
                  onClick: _cache[10] || (_cache[10] = ($event) => showSensitive.value = !showSensitive.value)
                }, [
                  (openBlock(), createBlock(resolveDynamicComponent(showSensitive.value ? unref(EyeOff) : unref(Eye)), { class: "h-3.5 w-3.5" }))
                ]),
                createBaseVNode("span", _hoisted_52, toDisplayString(showSensitive.value ? "敏感值已显示" : "敏感值已打码"), 1)
              ]),
              createBaseVNode("div", _hoisted_53, [
                createBaseVNode("table", _hoisted_54, [
                  createBaseVNode("tbody", null, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(envList.value, (e) => {
                      return openBlock(), createElementBlock("tr", {
                        key: e.key,
                        class: "border-b border-[#171f2a] last:border-b-0"
                      }, [
                        createBaseVNode("td", _hoisted_55, toDisplayString(e.key), 1),
                        createBaseVNode("td", _hoisted_56, [
                          e.sensitive === "true" && !showSensitive.value ? (openBlock(), createElementBlock("span", _hoisted_57, "••••••••")) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
                            createTextVNode(toDisplayString(e.value), 1)
                          ], 64))
                        ])
                      ]);
                    }), 128))
                  ])
                ])
              ])
            ])
          ])) : tab.value === "ports" ? (openBlock(), createElementBlock("div", _hoisted_58, [
            !portList.value.length ? (openBlock(), createElementBlock("div", _hoisted_59, " 这个容器没有声明任何端口。如果它只需要被同一网络里的其他容器访问，这是正常的。 ")) : (openBlock(), createElementBlock("div", _hoisted_60, [
              createBaseVNode("div", _hoisted_61, [
                createVNode(unref(Network), { class: "mt-[1px] h-3.5 w-3.5 flex-none" }),
                createBaseVNode("span", null, [
                  _cache[36] || (_cache[36] = createTextVNode(" 共 ", -1)),
                  createBaseVNode("b", _hoisted_62, toDisplayString(portList.value.length), 1),
                  _cache[37] || (_cache[37] = createTextVNode(" 个端口，其中 ", -1)),
                  createBaseVNode("b", _hoisted_63, toDisplayString(publishedCount.value), 1),
                  _cache[38] || (_cache[38] = createTextVNode(" 个已发布到宿主机。 实底徽标表示", -1)),
                  _cache[39] || (_cache[39] = createBaseVNode("b", null, "可从宿主机访问", -1)),
                  _cache[40] || (_cache[40] = createTextVNode("（左边是宿主端口）；虚线徽标表示", -1)),
                  _cache[41] || (_cache[41] = createBaseVNode("b", null, "只在容器网络内可见", -1)),
                  _cache[42] || (_cache[42] = createTextVNode("， 外部连不上。 ", -1))
                ])
              ]),
              createBaseVNode("div", _hoisted_64, [
                createBaseVNode("table", _hoisted_65, [
                  _cache[43] || (_cache[43] = createBaseVNode("thead", null, [
                    createBaseVNode("tr", { class: "border-b border-line-1 bg-ink-800 text-[11px] text-text-5" }, [
                      createBaseVNode("th", { class: "px-3 py-2 text-left font-medium" }, "宿主地址"),
                      createBaseVNode("th", { class: "px-3 py-2 text-left font-medium" }, "宿主端口"),
                      createBaseVNode("th", { class: "px-3 py-2 text-left font-medium" }, "容器端口"),
                      createBaseVNode("th", { class: "px-3 py-2 text-left font-medium" }, "协议"),
                      createBaseVNode("th", { class: "px-3 py-2 text-left font-medium" }, "可见性")
                    ])
                  ], -1)),
                  createBaseVNode("tbody", null, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(portList.value, (p) => {
                      return openBlock(), createElementBlock("tr", {
                        key: `${p.hostIp}:${p.hostPort}>${p.innerPort}/${p.proto}`,
                        class: "border-b border-[#171f2a] last:border-b-0"
                      }, [
                        createBaseVNode("td", _hoisted_66, toDisplayString(p.published ? p.hostIp || "0.0.0.0" : "—"), 1),
                        createBaseVNode("td", {
                          class: normalizeClass(["px-3 py-2 font-mono text-[11.5px]", p.published ? "font-semibold text-text-1" : "text-text-5"])
                        }, toDisplayString(p.published ? p.hostPort : "—"), 3),
                        createBaseVNode("td", _hoisted_67, toDisplayString(p.innerPort), 1),
                        createBaseVNode("td", _hoisted_68, toDisplayString(p.proto), 1),
                        createBaseVNode("td", _hoisted_69, [
                          createBaseVNode("span", {
                            class: normalizeClass(["dh-badge", p.published ? "dh-badge-run" : "dh-badge-plain"])
                          }, toDisplayString(p.published ? "对宿主机发布" : "仅容器内"), 3)
                        ])
                      ]);
                    }), 128))
                  ])
                ])
              ])
            ]))
          ])) : tab.value === "mounts" ? (openBlock(), createElementBlock("div", _hoisted_70, [
            !mounts.value.length ? (openBlock(), createElementBlock("div", _hoisted_71, " 这个容器没有任何挂载。 ")) : (openBlock(), createElementBlock("div", _hoisted_72, [
              createBaseVNode("div", _hoisted_73, [
                createVNode(unref(Box), { class: "mt-[1px] h-3.5 w-3.5 flex-none" }),
                _cache[44] || (_cache[44] = createBaseVNode("span", null, [
                  createTextVNode(" 「绑定挂载」的数据在"),
                  createBaseVNode("b", null, "宿主机目录"),
                  createTextVNode("上，Dockhelm 容器默认看不见 —— 要备份这份数据， 需要把对应宿主目录也挂进 Dockhelm（冒号右边叫什么名字都可以，启动时会自动识别）。「命名卷」的数据在 "),
                  createBaseVNode("code", { class: "text-text-3" }, "/var/lib/docker/volumes"),
                  createTextVNode(" 下，只读挂载该目录即可备份。 ")
                ], -1))
              ]),
              (openBlock(true), createElementBlock(Fragment, null, renderList(mounts.value, (m, i) => {
                return openBlock(), createElementBlock("div", {
                  key: i,
                  class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5"
                }, [
                  createBaseVNode("div", _hoisted_74, [
                    createBaseVNode("span", {
                      class: normalizeClass(["dh-badge", m.type === "volume" ? "dh-badge-accent" : "dh-badge-plain"])
                    }, toDisplayString(m.type === "volume" ? "命名卷/匿名卷" : m.type), 3),
                    createBaseVNode("span", {
                      class: normalizeClass(["dh-badge", m.rw ? "dh-badge-run" : "dh-badge-warn"])
                    }, toDisplayString(m.rw ? "读写" : "只读"), 3),
                    createBaseVNode("code", _hoisted_75, toDisplayString(m.destination), 1)
                  ]),
                  createBaseVNode("div", _hoisted_76, toDisplayString(m.name || m.source), 1),
                  createBaseVNode("div", _hoisted_77, toDisplayString(m.note), 1)
                ]);
              }), 128))
            ]))
          ])) : (openBlock(), createElementBlock("div", _hoisted_78, [
            !networks.value.length ? (openBlock(), createElementBlock("div", _hoisted_79, "没有网络信息")) : (openBlock(), createElementBlock("div", _hoisted_80, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(networks.value, (n) => {
                return openBlock(), createElementBlock("div", {
                  key: n.name,
                  class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5"
                }, [
                  createBaseVNode("div", _hoisted_81, [
                    createVNode(unref(Network), { class: "h-3.5 w-3.5 text-text-3" }),
                    createBaseVNode("span", _hoisted_82, toDisplayString(n.name), 1)
                  ]),
                  createBaseVNode("div", _hoisted_83, [
                    _cache[45] || (_cache[45] = createBaseVNode("div", { class: "text-text-5" }, "容器 IP", -1)),
                    createBaseVNode("div", _hoisted_84, toDisplayString(n.ip || "—"), 1),
                    _cache[46] || (_cache[46] = createBaseVNode("div", { class: "text-text-5" }, "网关", -1)),
                    createBaseVNode("div", _hoisted_85, toDisplayString(n.gateway || "—"), 1),
                    _cache[47] || (_cache[47] = createBaseVNode("div", { class: "text-text-5" }, "别名", -1)),
                    createBaseVNode("div", _hoisted_86, toDisplayString((n.aliases ?? []).join(", ") || "—"), 1)
                  ])
                ]);
              }), 128))
            ]))
          ]))
        ]),
        createBaseVNode("div", _hoisted_87, [
          createVNode(unref(ScrollText), { class: "h-3.5 w-3.5" }),
          _cache[49] || (_cache[49] = createTextVNode(" 日志每 8 秒自动刷新一次 ", -1)),
          createBaseVNode("span", _hoisted_88, [
            createVNode(unref(Terminal), { class: "h-3.5 w-3.5" }),
            _cache[48] || (_cache[48] = createTextVNode("容器名 ", -1)),
            createBaseVNode("code", _hoisted_89, toDisplayString(name.value), 1)
          ])
        ])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
