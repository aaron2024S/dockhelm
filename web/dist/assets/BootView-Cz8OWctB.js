import { d as defineComponent, E as useToastStore, a as onMounted, o as onUnmounted, c as createElementBlock, e as createBaseVNode, j as createCommentVNode, i as createBlock, g as unref, k as createTextVNode, t as toDisplayString, n as normalizeClass, z as normalizeStyle, F as Fragment, r as renderList, f as createVNode, I as withCtx, m as ref, p as computed, q as useRouter, B as api, s as openBlock, L as withModifiers } from "./index-OaLdLt7C.js";
import { c as containerStateLabel, f as formatDayTime } from "./format-DCfs2dCQ.js";
import { _ as _sfc_main$1 } from "./Modal.vue_vue_type_script_setup_true_lang-eEVxbfW_.js";
import { _ as _sfc_main$2 } from "./UiSwitch.vue_vue_type_script_setup_true_lang-DIvBPz9y.js";
import { L as LoaderCircle } from "./loader-circle-DIVuXAiz.js";
import { P as Play } from "./play-CcpKWWGe.js";
import { G as GripVertical } from "./grip-vertical-8jUXoaJW.js";
const _hoisted_1 = { class: "flex flex-col" };
const _hoisted_2 = { class: "sticky top-0 z-20 border-b border-line-1 bg-ink-950 px-[18px] pb-2.5 pt-[18px]" };
const _hoisted_3 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_4 = {
  key: 0,
  class: "dh-sub"
};
const _hoisted_5 = { class: "ml-auto flex flex-wrap items-center gap-3" };
const _hoisted_6 = ["disabled"];
const _hoisted_7 = ["title"];
const _hoisted_8 = {
  key: 0,
  class: "mt-2.5 flex items-center gap-2.5"
};
const _hoisted_9 = { class: "h-[3px] flex-1 overflow-hidden rounded-full bg-line-1" };
const _hoisted_10 = { class: "whitespace-nowrap text-[11.5px] text-text-5" };
const _hoisted_11 = { class: "flex flex-col gap-3.5 px-[18px] pb-[18px] pt-3" };
const _hoisted_12 = {
  key: 0,
  class: "dh-banner-warn"
};
const _hoisted_13 = { class: "dh-card" };
const _hoisted_14 = { class: "dh-card-body flex flex-wrap items-center gap-x-7 gap-y-3" };
const _hoisted_15 = { class: "flex items-center gap-2 text-[12.5px] text-text-4" };
const _hoisted_16 = ["value"];
const _hoisted_17 = { class: "flex items-center gap-2 text-[12.5px] text-text-4" };
const _hoisted_18 = ["value"];
const _hoisted_19 = { class: "flex items-center gap-2 text-[12.5px] text-text-4" };
const _hoisted_20 = ["value"];
const _hoisted_21 = { class: "dh-card" };
const _hoisted_22 = { class: "dh-card-head items-start" };
const _hoisted_23 = { class: "shrink-0 pt-[2px]" };
const _hoisted_24 = {
  key: 0,
  class: "ml-1 text-[11px] font-normal leading-[1.55] text-text-6"
};
const _hoisted_25 = { class: "ml-auto shrink-0 pl-3 pt-[1px] text-[11.5px] font-normal text-text-5" };
const _hoisted_26 = {
  key: 0,
  class: "dh-card-body text-[12.5px] text-text-4"
};
const _hoisted_27 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_28 = { class: "min-w-[660px] px-3 pb-2 pt-2" };
const _hoisted_29 = ["onDragstart", "onDragover", "onDrop"];
const _hoisted_30 = {
  class: "cursor-grab text-text-6 active:cursor-grabbing",
  title: "按住拖动排序"
};
const _hoisted_31 = { class: "text-text-5" };
const _hoisted_32 = ["title"];
const _hoisted_33 = ["title"];
const _hoisted_34 = ["value", "title", "onChange"];
const _hoisted_35 = { class: "flex justify-end" };
const _hoisted_36 = { class: "dh-card" };
const _hoisted_37 = { class: "dh-card-head items-start" };
const _hoisted_38 = { class: "shrink-0 pt-[2px]" };
const _hoisted_39 = {
  key: 0,
  class: "ml-1 text-[11px] font-normal leading-[1.55] text-text-6"
};
const _hoisted_40 = {
  key: 0,
  class: "dh-card-body text-[12.5px] text-text-4"
};
const _hoisted_41 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_42 = { class: "min-w-[660px] px-3 pb-2 pt-2" };
const _hoisted_43 = ["title"];
const _hoisted_44 = { class: "font-mono text-[11px]" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "BootView",
  setup(__props) {
    const toast = useToastStore();
    const router = useRouter();
    const loading = ref(true);
    const saving = ref(false);
    const confirming = ref(false);
    const plan = ref({
      enabled: true,
      initialDelaySec: 30,
      waitTimeoutSec: 60,
      defaultDelaySec: 3,
      items: []
    });
    const rows = ref([]);
    const self = ref({ name: "", restartPolicy: "" });
    const lastRun = ref(null);
    const status = ref(null);
    const managedRows = computed(() => rows.value.filter((r) => r.managed));
    const dockerRows = computed(() => rows.value.filter((r) => !r.managed));
    const pendingStart = computed(() => managedRows.value.filter((r) => r.auto && r.state !== "running"));
    const running = computed(() => !!status.value?.running);
    const estimateText = computed(() => {
      const starts = pendingStart.value.length;
      if (!starts) return managedRows.value.length ? "全部已在运行" : "—";
      const secs = plan.value.initialDelaySec + pendingStart.value.reduce((n, r) => n + 8 + r.delaySec, 0);
      return secs >= 60 ? `${Math.floor(secs / 60)} 分 ${secs % 60} 秒` : `${secs} 秒`;
    });
    const lastRunText = computed(() => {
      const l = lastRun.value;
      if (!l) return "上次执行：从未执行";
      const when = formatDayTime(l.ts);
      if (!l.total) return `上次执行：${when} · ${l.message}`;
      const parts = [`成功 ${l.success ?? 0}`, `跳过 ${l.skipped ?? 0}`];
      if (l.failed) parts.push(`失败 ${l.failed}`);
      return `上次执行：${when} · ${parts.join("、")}`;
    });
    const lastFailed = computed(() => lastRun.value?.status === "failed" || (lastRun.value?.failed ?? 0) > 0);
    const progressText = computed(() => {
      const s = status.value;
      if (!s?.running) return "";
      if (s.phase === "waiting-docker") return "等待 Docker 就绪…";
      if (s.phase === "initial-delay") return `初始等待 ${plan.value.initialDelaySec} 秒后开始…`;
      if (!s.total) return "正在读取容器…";
      const tail = s.phase === "starting" && s.current ? `正在启动 ${s.current}` : "准备中";
      const counts = [`成功 ${s.success}`, `跳过 ${s.skipped}`];
      if (s.failed) counts.push(`失败 ${s.failed}`);
      return `第 ${s.index}/${s.total} · ${tail} · ${counts.join("、")}`;
    });
    const progressPercent = computed(() => {
      const s = status.value;
      if (!s?.total) return 0;
      return Math.min(100, Math.round(s.index / s.total * 100));
    });
    const selfWarn = computed(() => {
      const p = self.value.restartPolicy;
      if (!self.value.name || !p) return "";
      if (p !== "no" && p !== "on-failure") return "";
      return `Dockhelm 自身容器（${self.value.name}）的重启策略是 ${p}：开机时它不会自己起来，面板起不来，这一页也就不会执行。建议在 NAS 面板把它改回 unless-stopped。`;
    });
    async function load() {
      loading.value = true;
      try {
        const d = await api.get("/api/boot");
        plan.value = d.plan;
        rows.value = d.containers ?? [];
        self.value = d.self ?? { name: "", restartPolicy: "" };
        lastRun.value = d.lastRun ?? null;
        status.value = d.status;
        if (d.status?.running) startPolling();
      } catch (e) {
        toast.error("读取顺序启动设置失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function save() {
      saving.value = true;
      try {
        const r = await api.put("/api/boot", plan.value);
        plan.value = r.plan;
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
        await load();
      } finally {
        saving.value = false;
      }
    }
    function setAuto(row, v) {
      const it = plan.value.items.find((x) => x.name === row.name);
      if (!it) return;
      it.auto = v;
      row.auto = v;
      void save();
    }
    function setDelay(row, ev) {
      const el = ev.target;
      const raw = Number(el.value);
      const v = Number.isFinite(raw) ? Math.min(300, Math.max(0, Math.round(raw))) : plan.value.defaultDelaySec;
      el.value = String(v);
      const it = plan.value.items.find((x) => x.name === row.name);
      if (!it || it.delaySec === v) return;
      it.delaySec = v;
      row.delaySec = v;
      void save();
    }
    function applyAllDelay() {
      const v = plan.value.defaultDelaySec;
      for (const it of plan.value.items) it.delaySec = v;
      for (const r of rows.value) r.delaySec = v;
      void save();
    }
    function setParam(key, ev) {
      const el = ev.target;
      const raw = Number(el.value);
      if (!Number.isFinite(raw)) return;
      plan.value[key] = Math.max(0, Math.round(raw));
      void save();
    }
    const dragFrom = ref(-1);
    const dragOver = ref(-1);
    function onDragStart(i, e) {
      dragFrom.value = i;
      if (e.dataTransfer) {
        e.dataTransfer.effectAllowed = "move";
        e.dataTransfer.setData("text/plain", String(i));
      }
    }
    function onDragOver(i, e) {
      if (dragFrom.value < 0) return;
      e.preventDefault();
      if (e.dataTransfer) e.dataTransfer.dropEffect = "move";
      dragOver.value = i;
    }
    function resetDrag() {
      dragFrom.value = -1;
      dragOver.value = -1;
    }
    function onDrop(i) {
      const from = dragFrom.value;
      resetDrag();
      if (from < 0 || from === i) return;
      const fromName = managedRows.value[from]?.name;
      const toName = managedRows.value[i]?.name;
      if (!fromName || !toName) return;
      const arr = [...plan.value.items];
      const fi = arr.findIndex((x) => x.name === fromName);
      const ti = arr.findIndex((x) => x.name === toName);
      if (fi < 0 || ti < 0) return;
      const [moved] = arr.splice(fi, 1);
      if (!moved) return;
      arr.splice(ti, 0, moved);
      plan.value.items = arr;
      void save();
    }
    let timer;
    function startPolling() {
      if (timer) return;
      timer = window.setInterval(() => void refreshStatus(), 2e3);
    }
    function stopPolling() {
      if (!timer) return;
      window.clearInterval(timer);
      timer = void 0;
    }
    async function refreshStatus() {
      try {
        const d = await api.get("/api/boot/status");
        const wasRunning = running.value;
        status.value = d.status;
        if (wasRunning && !d.status.running) {
          stopPolling();
          await load();
          toast.success("顺序启动已结束", "结果已写入运行记录");
        }
      } catch {
        stopPolling();
      }
    }
    async function doRun() {
      confirming.value = false;
      try {
        await api.post("/api/boot/run", { skipDelay: true });
        startPolling();
        await refreshStatus();
      } catch (e) {
        toast.error("无法开始", e instanceof Error ? e.message : String(e));
      }
    }
    function gotoRunLogs() {
      void router.push("/settings").then(() => {
        window.setTimeout(() => {
          document.getElementById("runlogs")?.scrollIntoView({ behavior: "smooth", block: "start" });
        }, 120);
      });
    }
    function openRunLogs() {
      if (!lastFailed.value) return;
      gotoRunLogs();
    }
    onMounted(() => void load());
    onUnmounted(() => stopPolling());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          createBaseVNode("div", _hoisted_3, [
            _cache[6] || (_cache[6] = createBaseVNode("div", { class: "dh-h1" }, "顺序启动", -1)),
            saving.value ? (openBlock(), createElementBlock("span", _hoisted_4, "保存中…")) : createCommentVNode("", true),
            createBaseVNode("div", _hoisted_5, [
              createBaseVNode("button", {
                class: "dh-btn",
                disabled: running.value || loading.value,
                onClick: _cache[0] || (_cache[0] = ($event) => confirming.value = true)
              }, [
                running.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "h-3.5 w-3.5 dh-spin"
                })) : (openBlock(), createBlock(unref(Play), {
                  key: 1,
                  class: "h-3.5 w-3.5"
                })),
                createTextVNode(" " + toDisplayString(running.value ? "执行中…" : "立即执行一次"), 1)
              ], 8, _hoisted_6),
              createBaseVNode("span", {
                class: normalizeClass(["text-[11.5px]", lastFailed.value ? "cursor-pointer text-err-text hover:underline" : "text-text-5"]),
                title: lastRun.value ? lastRun.value.message : "还没有执行过",
                onClick: openRunLogs
              }, toDisplayString(lastRunText.value), 11, _hoisted_7)
            ])
          ]),
          running.value ? (openBlock(), createElementBlock("div", _hoisted_8, [
            createBaseVNode("div", _hoisted_9, [
              createBaseVNode("div", {
                class: "h-full rounded-full bg-accent transition-all",
                style: normalizeStyle({ width: progressPercent.value + "%" })
              }, null, 4)
            ]),
            createBaseVNode("span", _hoisted_10, toDisplayString(progressText.value), 1)
          ])) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", _hoisted_11, [
          _cache[22] || (_cache[22] = createBaseVNode("div", { class: "dh-sub" }, " 本页负责将启动策略非 unless-stopped 且非 always 的容器按顺序启动。 ", -1)),
          selfWarn.value ? (openBlock(), createElementBlock("div", _hoisted_12, toDisplayString(selfWarn.value), 1)) : createCommentVNode("", true),
          createBaseVNode("div", _hoisted_13, [
            createBaseVNode("div", { class: "dh-card-head" }, [
              _cache[7] || (_cache[7] = createBaseVNode("span", null, "执行参数", -1)),
              createBaseVNode("button", {
                class: "ml-auto text-[11.5px] font-normal text-text-5 hover:text-accent",
                onClick: applyAllDelay
              }, " 批量应用到全部 ")
            ]),
            createBaseVNode("div", _hoisted_14, [
              createBaseVNode("label", _hoisted_15, [
                _cache[8] || (_cache[8] = createTextVNode(" 初始等待 ", -1)),
                createBaseVNode("input", {
                  class: "dh-input !w-[76px] text-center font-mono",
                  type: "number",
                  min: "0",
                  max: "1800",
                  value: plan.value.initialDelaySec,
                  onChange: _cache[1] || (_cache[1] = ($event) => setParam("initialDelaySec", $event))
                }, null, 40, _hoisted_16),
                _cache[9] || (_cache[9] = createTextVNode(" 秒 ", -1))
              ]),
              createBaseVNode("label", _hoisted_17, [
                _cache[10] || (_cache[10] = createTextVNode(" 单容器上限 ", -1)),
                createBaseVNode("input", {
                  class: "dh-input !w-[76px] text-center font-mono",
                  type: "number",
                  min: "5",
                  max: "600",
                  value: plan.value.waitTimeoutSec,
                  onChange: _cache[2] || (_cache[2] = ($event) => setParam("waitTimeoutSec", $event))
                }, null, 40, _hoisted_18),
                _cache[11] || (_cache[11] = createTextVNode(" 秒 ", -1))
              ]),
              createBaseVNode("label", _hoisted_19, [
                _cache[12] || (_cache[12] = createTextVNode(" 新行默认延时 ", -1)),
                createBaseVNode("input", {
                  class: "dh-input !w-[76px] text-center font-mono",
                  type: "number",
                  min: "0",
                  max: "300",
                  value: plan.value.defaultDelaySec,
                  onChange: _cache[3] || (_cache[3] = ($event) => setParam("defaultDelaySec", $event))
                }, null, 40, _hoisted_20),
                _cache[13] || (_cache[13] = createTextVNode(" 秒 ", -1))
              ]),
              _cache[14] || (_cache[14] = createBaseVNode("span", { class: "text-[11.5px] text-text-6" }, "延时 = 这一行启动成功后等 N 秒再起下一个", -1))
            ])
          ]),
          createBaseVNode("div", _hoisted_21, [
            createBaseVNode("div", _hoisted_22, [
              createBaseVNode("span", _hoisted_23, "由 dockhelm 按顺序启动 · " + toDisplayString(managedRows.value.length) + " 个", 1),
              managedRows.value.length ? (openBlock(), createElementBlock("span", _hoisted_24, " 每次 Dockhelm 启动都会执行一次，已在运行的跳过，只有开机启动已打开的未运行容器会被拉起。 ")) : createCommentVNode("", true),
              createBaseVNode("span", _hoisted_25, " 预计 " + toDisplayString(estimateText.value), 1)
            ]),
            !managedRows.value.length ? (openBlock(), createElementBlock("div", _hoisted_26, toDisplayString(loading.value ? "正在载入…" : "还没有归本页管的容器。把容器的重启策略改成 on-failure（或 no），刷新后就会出现在这里。"), 1)) : (openBlock(), createElementBlock("div", _hoisted_27, [
              createBaseVNode("div", _hoisted_28, [
                _cache[15] || (_cache[15] = createBaseVNode("div", { class: "grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 px-2 pb-1.5 text-[11px] text-text-6" }, [
                  createBaseVNode("span"),
                  createBaseVNode("span", null, "#"),
                  createBaseVNode("span", null, "容器"),
                  createBaseVNode("span", null, "状态"),
                  createBaseVNode("span", null, "策略"),
                  createBaseVNode("span", null, "延时"),
                  createBaseVNode("span", { class: "text-right whitespace-nowrap" }, "开机启动")
                ], -1)),
                (openBlock(true), createElementBlock(Fragment, null, renderList(managedRows.value, (row, i) => {
                  return openBlock(), createElementBlock("div", {
                    key: row.name,
                    class: normalizeClass(["dh-boot-row grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 rounded-lg px-2 py-1.5 text-[12.5px]", [
                      dragOver.value === i && dragFrom.value >= 0 && dragFrom.value !== i ? "border-t-2 border-accent" : "",
                      dragFrom.value === i ? "opacity-45" : "",
                      row.state === "running" ? "dh-boot-run" : ""
                    ]]),
                    draggable: "true",
                    onDragstart: ($event) => onDragStart(i, $event),
                    onDragover: ($event) => onDragOver(i, $event),
                    onDrop: withModifiers(($event) => onDrop(i), ["prevent"]),
                    onDragend: resetDrag
                  }, [
                    createBaseVNode("span", _hoisted_30, [
                      createVNode(unref(GripVertical), { class: "h-3.5 w-3.5" })
                    ]),
                    createBaseVNode("span", _hoisted_31, toDisplayString(i + 1), 1),
                    createBaseVNode("span", {
                      class: "truncate font-medium",
                      title: row.name
                    }, toDisplayString(row.name), 9, _hoisted_32),
                    createBaseVNode("span", {
                      class: normalizeClass(row.state === "running" ? "text-run-text" : "text-text-5")
                    }, toDisplayString(unref(containerStateLabel)(row.state)), 3),
                    createBaseVNode("span", {
                      class: "font-mono text-[11px] text-text-4",
                      title: row.restartPolicy
                    }, toDisplayString(row.restartPolicy || "—"), 9, _hoisted_33),
                    createBaseVNode("input", {
                      class: "dh-input !w-[52px] !px-1 text-center font-mono text-[11.5px]",
                      type: "number",
                      min: "0",
                      max: "300",
                      value: row.delaySec,
                      title: `${row.name} 启动成功后等多少秒再起下一个`,
                      onChange: ($event) => setDelay(row, $event)
                    }, null, 40, _hoisted_34),
                    createBaseVNode("div", _hoisted_35, [
                      createVNode(_sfc_main$2, {
                        "model-value": row.auto,
                        label: `${row.name} 开机启动`,
                        "onUpdate:modelValue": (v) => setAuto(row, v)
                      }, null, 8, ["model-value", "label", "onUpdate:modelValue"])
                    ])
                  ], 42, _hoisted_29);
                }), 128))
              ]),
              _cache[16] || (_cache[16] = createBaseVNode("div", { class: "px-5 pb-3 text-[11px] text-text-6" }, " 拖动手柄整行拖动，松手即保存，不用点「保存」。打开开关 = 开机时由本页启动它；关掉 = 永远不碰它。 ", -1))
            ]))
          ]),
          createBaseVNode("div", _hoisted_36, [
            createBaseVNode("div", _hoisted_37, [
              createBaseVNode("span", _hoisted_38, "由 Docker 自行启动 · " + toDisplayString(dockerRows.value.length) + " 个", 1),
              dockerRows.value.length ? (openBlock(), createElementBlock("span", _hoisted_39, " 这些容器由 Docker 在开机时自动启动，不受本程序控制。想让它按顺序启动，去 NAS 面板把重启策略改成 on-failure，刷新后会自动加入到本程序的顺序启动管理列表。 ")) : createCommentVNode("", true)
            ]),
            !dockerRows.value.length ? (openBlock(), createElementBlock("div", _hoisted_40, toDisplayString(loading.value ? "正在载入…" : "没有 always / unless-stopped 的容器 —— 全部受本页控制。"), 1)) : (openBlock(), createElementBlock("div", _hoisted_41, [
              createBaseVNode("div", _hoisted_42, [
                _cache[21] || (_cache[21] = createBaseVNode("div", { class: "grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 px-2 pb-1.5 text-[11px] text-text-6" }, [
                  createBaseVNode("span"),
                  createBaseVNode("span"),
                  createBaseVNode("span", null, "容器"),
                  createBaseVNode("span", null, "状态"),
                  createBaseVNode("span", null, "策略"),
                  createBaseVNode("span"),
                  createBaseVNode("span")
                ], -1)),
                (openBlock(true), createElementBlock(Fragment, null, renderList(dockerRows.value, (row) => {
                  return openBlock(), createElementBlock("div", {
                    key: row.name,
                    class: "grid grid-cols-[18px_22px_minmax(0,1fr)_72px_92px_54px_56px] items-center gap-2 px-2 py-1.5 text-[12.5px] text-text-5"
                  }, [
                    _cache[17] || (_cache[17] = createBaseVNode("span", null, null, -1)),
                    _cache[18] || (_cache[18] = createBaseVNode("span", null, "—", -1)),
                    createBaseVNode("span", {
                      class: "truncate",
                      title: row.name
                    }, toDisplayString(row.name), 9, _hoisted_43),
                    createBaseVNode("span", {
                      class: normalizeClass(row.state === "running" ? "" : "text-text-6")
                    }, toDisplayString(unref(containerStateLabel)(row.state)), 3),
                    createBaseVNode("span", _hoisted_44, toDisplayString(row.restartPolicy || "—"), 1),
                    _cache[19] || (_cache[19] = createBaseVNode("span", null, null, -1)),
                    _cache[20] || (_cache[20] = createBaseVNode("span", null, null, -1))
                  ]);
                }), 128))
              ])
            ]))
          ]),
          createBaseVNode("div", { class: "flex items-center text-[11.5px] text-text-6" }, [
            createBaseVNode("button", {
              class: "ml-auto text-text-5 hover:text-accent",
              onClick: gotoRunLogs
            }, "查看运行记录")
          ])
        ]),
        createVNode(_sfc_main$1, {
          open: confirming.value,
          title: "立即执行一次顺序启动？",
          subtitle: `将按当前顺序启动 ${pendingStart.value.length} 个已停止的容器，全程串行，随时可在页头看到进度。`,
          width: "460px",
          onClose: _cache[5] || (_cache[5] = ($event) => confirming.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[4] || (_cache[4] = ($event) => confirming.value = false)
            }, "取消"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: doRun
            }, [
              createVNode(unref(Play), { class: "h-3.5 w-3.5" }),
              _cache[23] || (_cache[23] = createTextVNode("执行 ", -1))
            ])
          ]),
          default: withCtx(() => [
            _cache[24] || (_cache[24] = createBaseVNode("div", { class: "text-[12.5px] leading-relaxed text-text-4" }, " 已在运行的、开关关掉的、以及归 Docker 管的容器都会被跳过。本次会跳过开头的初始等待，立刻开始。 ", -1))
          ]),
          _: 1
        }, 8, ["open", "subtitle"])
      ]);
    };
  }
});
const _export_sfc = (sfc, props) => {
  const target = sfc.__vccOpts || sfc;
  for (const [key, val] of props) {
    target[key] = val;
  }
  return target;
};
const BootView = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-32e424b6"]]);
export {
  BootView as default
};
