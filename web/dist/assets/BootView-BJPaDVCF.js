import { d as defineComponent, E as useToastStore, a as onMounted, o as onUnmounted, c as createElementBlock, e as createBaseVNode, j as createCommentVNode, i as createBlock, g as unref, k as createTextVNode, t as toDisplayString, f as createVNode, n as normalizeClass, z as normalizeStyle, F as Fragment, r as renderList, I as withCtx, m as ref, p as computed, q as useRouter, B as api, s as openBlock } from "./index-CgBqIvFv.js";
import { c as containerStateLabel, f as formatDayTime } from "./format-AUhC3z7I.js";
import { _ as _sfc_main$1 } from "./Modal.vue_vue_type_script_setup_true_lang-C_6g4ysF.js";
import { _ as _sfc_main$2 } from "./UiSwitch.vue_vue_type_script_setup_true_lang-bxaO3AFa.js";
import { L as LoaderCircle } from "./loader-circle-BrV2CePY.js";
import { P as Play } from "./play-3GoAoMNi.js";
import { H as History } from "./history-x1gpnYeW.js";
import { G as GripVertical } from "./grip-vertical-DxIEv0Yq.js";
const _hoisted_1 = { class: "flex flex-col" };
const _hoisted_2 = { class: "sticky top-0 z-20 border-b border-line-1 bg-ink-950 px-[18px] pb-2.5 pt-[18px]" };
const _hoisted_3 = { class: "flex flex-wrap items-center gap-x-2.5 gap-y-2" };
const _hoisted_4 = {
  key: 0,
  class: "dh-sub"
};
const _hoisted_5 = ["disabled"];
const _hoisted_6 = { class: "order-2 flex min-w-0 flex-wrap items-center gap-x-3 gap-y-1 max-md:w-full md:order-3" };
const _hoisted_7 = {
  key: 0,
  class: "flex items-center gap-1.5 text-[11.5px] text-text-5"
};
const _hoisted_8 = ["title"];
const _hoisted_9 = {
  key: 0,
  class: "mt-2.5 flex items-center gap-2.5"
};
const _hoisted_10 = { class: "h-[3px] flex-1 overflow-hidden rounded-full bg-line-1" };
const _hoisted_11 = { class: "whitespace-nowrap text-[11.5px] text-text-5" };
const _hoisted_12 = { class: "flex flex-col gap-3.5 px-[18px] pb-[18px] pt-3" };
const _hoisted_13 = {
  key: 0,
  class: "dh-banner-warn"
};
const _hoisted_14 = { class: "dh-card" };
const _hoisted_15 = { class: "dh-card-body" };
const _hoisted_16 = { class: "flex flex-wrap items-center gap-x-7 gap-y-3" };
const _hoisted_17 = { class: "flex items-center gap-2 text-[12.5px] text-text-4" };
const _hoisted_18 = ["value"];
const _hoisted_19 = {
  class: "flex items-center gap-2 text-[12.5px] text-text-4",
  title: "发出 start 后最多等多久它进入 running；到点还没起来就记这个容器失败，接着起下一个"
};
const _hoisted_20 = ["value"];
const _hoisted_21 = { class: "flex items-center gap-2.5 rounded-[10px] border border-line-4 bg-ink-650 py-1.5 pl-3 pr-2" };
const _hoisted_22 = { class: "flex items-center gap-2 text-[12.5px] text-text-4" };
const _hoisted_23 = ["value"];
const _hoisted_24 = ["disabled", "title"];
const _hoisted_25 = { class: "dh-card" };
const _hoisted_26 = { class: "dh-card-head items-baseline" };
const _hoisted_27 = { class: "shrink-0" };
const _hoisted_28 = {
  key: 0,
  class: "ml-1 text-[11px] font-normal leading-[1.55] text-text-6"
};
const _hoisted_29 = ["title"];
const _hoisted_30 = {
  key: 0,
  class: "dh-card-body text-[12.5px] text-text-4"
};
const _hoisted_31 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_32 = ["data-boot-name", "onPointerdown"];
const _hoisted_33 = {
  class: "dh-c-grip",
  "data-grip": "",
  title: "按住拖动手柄排序（触屏也可以）"
};
const _hoisted_34 = { class: "dh-c-idx text-text-5" };
const _hoisted_35 = ["title"];
const _hoisted_36 = { class: "dh-c-meta" };
const _hoisted_37 = ["title"];
const _hoisted_38 = { class: "dh-c-delay" };
const _hoisted_39 = ["value", "title", "onChange"];
const _hoisted_40 = { class: "dh-c-switch flex justify-end" };
const _hoisted_41 = { class: "dh-card" };
const _hoisted_42 = { class: "dh-card-head items-baseline" };
const _hoisted_43 = { class: "shrink-0" };
const _hoisted_44 = {
  key: 0,
  class: "ml-1 text-[11px] font-normal leading-[1.55] text-text-6"
};
const _hoisted_45 = {
  key: 0,
  class: "dh-card-body text-[12.5px] text-text-4"
};
const _hoisted_46 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_47 = { class: "min-w-0 px-3 pb-2 pt-2" };
const _hoisted_48 = ["title"];
const _hoisted_49 = { class: "dh-c-meta" };
const _hoisted_50 = { class: "font-mono text-[11px]" };
const DRAG_THRESHOLD = 4;
const EDGE_ZONE = 64;
const EDGE_MAX_SPEED = 18;
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
    const hasProgress = computed(() => !!status.value?.total);
    const fmtSecs = (s) => s >= 60 ? `${Math.floor(s / 60)} 分 ${s % 60} 秒` : `${s} 秒`;
    const estimate = computed(() => {
      const rows2 = managedRows.value.filter((r) => r.auto);
      if (!rows2.length) return null;
      const init = plan.value.initialDelaySec;
      const startup = rows2.length * 8;
      const delays = rows2.reduce((n, r) => n + r.delaySec, 0);
      const secs = init + startup + delays;
      return {
        text: fmtSecs(secs),
        title: `按当前设置估算：初始等待 ${init} 秒 + ${rows2.length} 个容器 × 约 8 秒（启动并等到 running） + 各容器延时合计 ${delays} 秒 = ${fmtSecs(secs)}。已在运行的容器会被跳过，实际用时更短。`
      };
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
        applyPlanOrder();
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
        await load();
      } finally {
        saving.value = false;
      }
    }
    function applyPlanOrder() {
      const rank = new Map(plan.value.items.map((it, i) => [it.name, i]));
      const pos = (r) => rank.get(r.name) ?? Number.MAX_SAFE_INTEGER;
      rows.value = [...rows.value].sort((a, b) => pos(a) - pos(b));
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
      const n = plan.value.items.length;
      if (!n) return;
      for (const it of plan.value.items) it.delaySec = v;
      for (const r of rows.value) r.delaySec = v;
      void save();
      toast.success("已应用到全部行", `${n} 个容器的延时都设为 ${v} 秒`);
    }
    function setParam(key, ev) {
      const el = ev.target;
      const raw = Number(el.value);
      if (!Number.isFinite(raw)) return;
      plan.value[key] = Math.max(0, Math.round(raw));
      void save();
    }
    const listRef = ref(null);
    const dragFrom = ref(-1);
    const dropSlot = ref(-1);
    const dragActive = ref(false);
    let dragPointerId = -1;
    let dragStartY = 0;
    let lastPointerY = 0;
    let dragRects = [];
    let scrollEl = null;
    let scrollSpeed = 0;
    let rafId = 0;
    let captureEl = null;
    function onRowPointerDown(i, e) {
      const target = e.target;
      if (e.pointerType === "mouse") {
        if (e.button !== 0) return;
        if (target.closest("input, button, a, label, select, textarea")) return;
      } else if (!target.closest("[data-grip]")) {
        return;
      }
      if (managedRows.value.length < 2) return;
      dragFrom.value = i;
      dropSlot.value = i;
      dragPointerId = e.pointerId;
      dragStartY = e.clientY;
      lastPointerY = e.clientY;
      dragActive.value = false;
      dragRects = measureRows();
      scrollEl = findScrollParent(e.currentTarget);
      scrollSpeed = 0;
      window.addEventListener("pointermove", onDragMove, { passive: false });
      window.addEventListener("pointerup", onDragEnd);
      window.addEventListener("pointercancel", onDragCancel);
      rafId = requestAnimationFrame(autoScrollTick);
      captureEl = e.currentTarget ?? null;
      try {
        captureEl?.setPointerCapture(e.pointerId);
      } catch {
      }
      document.body.classList.add("dh-drag-lock");
      e.preventDefault();
    }
    function measureRows() {
      const host = listRef.value;
      if (!host) return [];
      return Array.from(host.querySelectorAll("[data-boot-name]")).map((el) => {
        const r = el.getBoundingClientRect();
        return { top: r.top, bottom: r.bottom };
      });
    }
    function findScrollParent(from) {
      let p = from?.parentElement ?? null;
      while (p) {
        const st = getComputedStyle(p);
        if (/(auto|scroll)/.test(st.overflowY) && p.scrollHeight > p.clientHeight + 1) return p;
        p = p.parentElement;
      }
      return null;
    }
    function slotAt(y) {
      const n = dragRects.length;
      for (let i = 0; i < n; i++) {
        const r = dragRects[i];
        if (y < r.bottom) return y < (r.top + r.bottom) / 2 ? i : i + 1;
      }
      return n;
    }
    function onDragMove(e) {
      if (e.pointerId !== dragPointerId) return;
      lastPointerY = e.clientY;
      if (!dragActive.value) {
        if (Math.abs(e.clientY - dragStartY) < DRAG_THRESHOLD) return;
        dragActive.value = true;
      }
      e.preventDefault();
      dropSlot.value = slotAt(e.clientY);
      updateScrollSpeed(e.clientY);
    }
    function updateScrollSpeed(y) {
      if (!scrollEl) {
        scrollSpeed = 0;
        return;
      }
      const rect = scrollEl.getBoundingClientRect();
      const top = rect.top + 72;
      const bottom = rect.bottom - 20;
      if (y < top + EDGE_ZONE) {
        scrollSpeed = -Math.ceil((top + EDGE_ZONE - y) / EDGE_ZONE * EDGE_MAX_SPEED);
      } else if (y > bottom - EDGE_ZONE) {
        scrollSpeed = Math.ceil((y - (bottom - EDGE_ZONE)) / EDGE_ZONE * EDGE_MAX_SPEED);
      } else {
        scrollSpeed = 0;
      }
      scrollSpeed = Math.max(-EDGE_MAX_SPEED, Math.min(EDGE_MAX_SPEED, scrollSpeed));
    }
    function autoScrollTick() {
      if (dragFrom.value < 0) return;
      if (scrollEl && scrollSpeed !== 0) {
        const before = scrollEl.scrollTop;
        scrollEl.scrollTop = before + scrollSpeed;
        const moved = scrollEl.scrollTop - before;
        if (moved !== 0) {
          for (const r of dragRects) {
            r.top -= moved;
            r.bottom -= moved;
          }
          dropSlot.value = slotAt(lastPointerY);
        }
      }
      rafId = requestAnimationFrame(autoScrollTick);
    }
    function onDragEnd(e) {
      if (e.pointerId !== dragPointerId) return;
      const from = dragFrom.value;
      const slot = dropSlot.value;
      const active = dragActive.value;
      cleanupDrag();
      if (active && from >= 0 && slot >= 0) commitReorder(from, slot);
    }
    function onDragCancel(e) {
      if (e.pointerId !== dragPointerId) return;
      cleanupDrag();
    }
    function cleanupDrag() {
      window.removeEventListener("pointermove", onDragMove);
      window.removeEventListener("pointerup", onDragEnd);
      window.removeEventListener("pointercancel", onDragCancel);
      if (rafId) cancelAnimationFrame(rafId);
      rafId = 0;
      if (captureEl && dragPointerId >= 0) {
        try {
          if (captureEl.hasPointerCapture(dragPointerId)) captureEl.releasePointerCapture(dragPointerId);
        } catch {
        }
      }
      captureEl = null;
      document.body.classList.remove("dh-drag-lock");
      dragFrom.value = -1;
      dropSlot.value = -1;
      dragActive.value = false;
      dragRects = [];
      scrollEl = null;
      scrollSpeed = 0;
      dragPointerId = -1;
    }
    function commitReorder(from, slot) {
      const names = managedRows.value.map((r) => r.name);
      if (slot === from || slot === from + 1) return;
      const moved = names.splice(from, 1)[0];
      if (!moved) return;
      names.splice(slot > from ? slot - 1 : slot, 0, moved);
      const byName = new Map(plan.value.items.map((it) => [it.name, it]));
      const managedSet = new Set(names);
      const queue = [];
      for (const n of names) {
        const it = byName.get(n);
        if (it) queue.push(it);
      }
      let k = 0;
      plan.value.items = plan.value.items.map((it) => {
        if (!managedSet.has(it.name)) return it;
        return queue[k++] ?? it;
      });
      applyPlanOrder();
      void save();
    }
    function dropClass(i) {
      if (!dragActive.value || dropSlot.value < 0) return "";
      if (dropSlot.value === i) return "dh-boot-drop-top";
      if (dropSlot.value === dragRects.length && i === dragRects.length - 1) return "dh-boot-drop-bottom";
      return "";
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
    onUnmounted(() => {
      stopPolling();
      cleanupDrag();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          createBaseVNode("div", _hoisted_3, [
            _cache[7] || (_cache[7] = createBaseVNode("div", { class: "dh-h1" }, "顺序启动", -1)),
            saving.value ? (openBlock(), createElementBlock("span", _hoisted_4, "保存中…")) : createCommentVNode("", true),
            createBaseVNode("button", {
              class: "dh-btn order-1 ml-auto",
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
            ], 8, _hoisted_5),
            createBaseVNode("div", _hoisted_6, [
              running.value && !hasProgress.value ? (openBlock(), createElementBlock("span", _hoisted_7, [
                createVNode(unref(LoaderCircle), { class: "h-3 w-3 dh-spin" }),
                createTextVNode(toDisplayString(progressText.value), 1)
              ])) : (openBlock(), createElementBlock("span", {
                key: 1,
                class: normalizeClass(["text-[11.5px]", lastFailed.value ? "cursor-pointer text-err-text hover:underline" : "text-text-5"]),
                title: lastRun.value ? lastRun.value.message : "还没有执行过",
                onClick: openRunLogs
              }, toDisplayString(lastRunText.value), 11, _hoisted_8)),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-ghost dh-btn-sm !gap-1.5 text-[11.5px] max-md:ml-auto max-md:!px-2",
                title: "打开设置页的运行记录",
                onClick: gotoRunLogs
              }, [
                createVNode(unref(History), { class: "h-3.5 w-3.5" }),
                _cache[6] || (_cache[6] = createBaseVNode("span", { class: "max-md:hidden" }, "运行记录", -1))
              ])
            ])
          ]),
          running.value && hasProgress.value ? (openBlock(), createElementBlock("div", _hoisted_9, [
            createBaseVNode("div", _hoisted_10, [
              createBaseVNode("div", {
                class: "h-full rounded-full bg-accent transition-all",
                style: normalizeStyle({ width: progressPercent.value + "%" })
              }, null, 4)
            ]),
            createBaseVNode("span", _hoisted_11, toDisplayString(progressText.value), 1)
          ])) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", _hoisted_12, [
          _cache[23] || (_cache[23] = createBaseVNode("div", { class: "dh-sub" }, " 本页负责将启动策略非 unless-stopped 且非 always 的容器按顺序启动。 ", -1)),
          selfWarn.value ? (openBlock(), createElementBlock("div", _hoisted_13, toDisplayString(selfWarn.value), 1)) : createCommentVNode("", true),
          createBaseVNode("div", _hoisted_14, [
            _cache[14] || (_cache[14] = createBaseVNode("div", { class: "dh-card-head items-baseline" }, [
              createBaseVNode("span", { class: "shrink-0" }, "执行参数"),
              createBaseVNode("span", { class: "text-[11px] font-normal leading-[1.55] text-text-6" }, " 延时 = 容器启动成功后等 N 秒再启动下一个。新加入的容器默认取这个值；「应用到全部行」只改每行的延时，不动初始等待与容器启动超时。 ")
            ], -1)),
            createBaseVNode("div", _hoisted_15, [
              createBaseVNode("div", _hoisted_16, [
                createBaseVNode("label", _hoisted_17, [
                  _cache[8] || (_cache[8] = createTextVNode(" 初始等待 ", -1)),
                  createBaseVNode("input", {
                    class: "dh-input !w-[76px] text-center font-mono",
                    type: "number",
                    min: "0",
                    max: "1800",
                    value: plan.value.initialDelaySec,
                    onChange: _cache[1] || (_cache[1] = ($event) => setParam("initialDelaySec", $event))
                  }, null, 40, _hoisted_18),
                  _cache[9] || (_cache[9] = createTextVNode(" 秒 ", -1))
                ]),
                createBaseVNode("label", _hoisted_19, [
                  _cache[10] || (_cache[10] = createTextVNode(" 容器启动超时 ", -1)),
                  createBaseVNode("input", {
                    class: "dh-input !w-[76px] text-center font-mono",
                    type: "number",
                    min: "5",
                    max: "600",
                    value: plan.value.waitTimeoutSec,
                    onChange: _cache[2] || (_cache[2] = ($event) => setParam("waitTimeoutSec", $event))
                  }, null, 40, _hoisted_20),
                  _cache[11] || (_cache[11] = createTextVNode(" 秒 ", -1))
                ]),
                createBaseVNode("div", _hoisted_21, [
                  createBaseVNode("label", _hoisted_22, [
                    _cache[12] || (_cache[12] = createTextVNode(" 默认延时 ", -1)),
                    createBaseVNode("input", {
                      class: "dh-input !w-[76px] text-center font-mono",
                      type: "number",
                      min: "0",
                      max: "300",
                      value: plan.value.defaultDelaySec,
                      onChange: _cache[3] || (_cache[3] = ($event) => setParam("defaultDelaySec", $event))
                    }, null, 40, _hoisted_23),
                    _cache[13] || (_cache[13] = createTextVNode(" 秒 ", -1))
                  ]),
                  createBaseVNode("button", {
                    class: "dh-btn dh-btn-sm text-[11.5px]",
                    disabled: !managedRows.value.length || saving.value,
                    title: `把全部 ${managedRows.value.length} 个容器的延时都设成 ${plan.value.defaultDelaySec} 秒`,
                    onClick: applyAllDelay
                  }, " 应用到全部行 ", 8, _hoisted_24)
                ])
              ])
            ])
          ]),
          createBaseVNode("div", _hoisted_25, [
            createBaseVNode("div", _hoisted_26, [
              createBaseVNode("span", _hoisted_27, "由 dockhelm 按顺序启动 · " + toDisplayString(managedRows.value.length) + " 个", 1),
              managedRows.value.length ? (openBlock(), createElementBlock("span", _hoisted_28, " 每次 Dockhelm 启动都会执行一次，已在运行的跳过，只有开机启动已打开的未运行容器会被拉起。 ")) : createCommentVNode("", true),
              estimate.value ? (openBlock(), createElementBlock("span", {
                key: 1,
                class: "ml-auto shrink-0 pl-3 text-[11.5px] font-normal text-text-5",
                title: estimate.value.title
              }, " 预计全部启动完成需 " + toDisplayString(estimate.value.text), 9, _hoisted_29)) : createCommentVNode("", true)
            ]),
            !managedRows.value.length ? (openBlock(), createElementBlock("div", _hoisted_30, toDisplayString(loading.value ? "正在载入…" : "还没有归本页管的容器。把容器的重启策略改成 on-failure（或 no），刷新后就会出现在这里。"), 1)) : (openBlock(), createElementBlock("div", _hoisted_31, [
              createBaseVNode("div", {
                ref_key: "listRef",
                ref: listRef,
                class: "min-w-0 px-3 pb-2 pt-2"
              }, [
                _cache[17] || (_cache[17] = createBaseVNode("div", { class: "dh-boot-grid dh-boot-head px-2 pb-1.5 text-[11px] text-text-6" }, [
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
                    "data-boot-name": row.name,
                    class: normalizeClass(["dh-boot-row dh-boot-grid rounded-lg px-2 py-1.5 text-[12.5px]", [
                      dropClass(i),
                      dragActive.value && dragFrom.value === i ? "dh-boot-dragging" : "",
                      row.state === "running" ? "dh-boot-run" : ""
                    ]]),
                    onPointerdown: ($event) => onRowPointerDown(i, $event)
                  }, [
                    createBaseVNode("span", _hoisted_33, [
                      createVNode(unref(GripVertical), { class: "h-3.5 w-3.5" })
                    ]),
                    createBaseVNode("span", _hoisted_34, toDisplayString(i + 1), 1),
                    createBaseVNode("span", {
                      class: "dh-c-name truncate font-medium",
                      title: row.name
                    }, toDisplayString(row.name), 9, _hoisted_35),
                    createBaseVNode("span", _hoisted_36, [
                      createBaseVNode("span", {
                        class: normalizeClass(row.state === "running" ? "text-run-text" : "text-text-5")
                      }, toDisplayString(unref(containerStateLabel)(row.state)), 3),
                      createBaseVNode("span", {
                        class: "font-mono text-[11px] text-text-4",
                        title: row.restartPolicy
                      }, toDisplayString(row.restartPolicy || "—"), 9, _hoisted_37)
                    ]),
                    createBaseVNode("span", _hoisted_38, [
                      _cache[15] || (_cache[15] = createBaseVNode("span", { class: "dh-delay-label" }, "延时", -1)),
                      createBaseVNode("input", {
                        class: "dh-input !w-[52px] !px-1 text-center font-mono text-[11.5px] max-md:!w-[44px] max-md:!py-0",
                        type: "number",
                        min: "0",
                        max: "300",
                        value: row.delaySec,
                        title: `${row.name} 启动成功后等多少秒再起下一个`,
                        onChange: ($event) => setDelay(row, $event)
                      }, null, 40, _hoisted_39),
                      _cache[16] || (_cache[16] = createBaseVNode("span", { class: "dh-delay-label" }, "秒", -1))
                    ]),
                    createBaseVNode("div", _hoisted_40, [
                      createVNode(_sfc_main$2, {
                        "model-value": row.auto,
                        label: `${row.name} 开机启动`,
                        "onUpdate:modelValue": (v) => setAuto(row, v)
                      }, null, 8, ["model-value", "label", "onUpdate:modelValue"])
                    ])
                  ], 42, _hoisted_32);
                }), 128))
              ], 512)
            ]))
          ]),
          createBaseVNode("div", _hoisted_41, [
            createBaseVNode("div", _hoisted_42, [
              createBaseVNode("span", _hoisted_43, "由 Docker 自行启动 · " + toDisplayString(dockerRows.value.length) + " 个", 1),
              dockerRows.value.length ? (openBlock(), createElementBlock("span", _hoisted_44, " 这些容器由 Docker 在开机时自动启动，不受本程序控制。想让它按顺序启动，去 NAS 面板把重启策略改成 on-failure，刷新后会自动加入到本程序的顺序启动管理列表。 ")) : createCommentVNode("", true)
            ]),
            !dockerRows.value.length ? (openBlock(), createElementBlock("div", _hoisted_45, toDisplayString(loading.value ? "正在载入…" : "没有 always / unless-stopped 的容器 —— 全部受本页控制。"), 1)) : (openBlock(), createElementBlock("div", _hoisted_46, [
              createBaseVNode("div", _hoisted_47, [
                _cache[22] || (_cache[22] = createBaseVNode("div", { class: "dh-boot-grid dh-boot-head px-2 pb-1.5 text-[11px] text-text-6" }, [
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
                    class: "dh-boot-row dh-boot-grid px-2 py-1.5 text-[12.5px] text-text-5"
                  }, [
                    _cache[18] || (_cache[18] = createBaseVNode("span", { class: "dh-c-grip" }, null, -1)),
                    _cache[19] || (_cache[19] = createBaseVNode("span", { class: "dh-c-idx" }, "—", -1)),
                    createBaseVNode("span", {
                      class: "dh-c-name truncate",
                      title: row.name
                    }, toDisplayString(row.name), 9, _hoisted_48),
                    createBaseVNode("span", _hoisted_49, [
                      createBaseVNode("span", {
                        class: normalizeClass(row.state === "running" ? "" : "text-text-6")
                      }, toDisplayString(unref(containerStateLabel)(row.state)), 3),
                      createBaseVNode("span", _hoisted_50, toDisplayString(row.restartPolicy || "—"), 1)
                    ]),
                    _cache[20] || (_cache[20] = createBaseVNode("span", { class: "dh-c-delay" }, null, -1)),
                    _cache[21] || (_cache[21] = createBaseVNode("span", { class: "dh-c-filler" }, null, -1))
                  ]);
                }), 128))
              ])
            ]))
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
              _cache[24] || (_cache[24] = createTextVNode("执行 ", -1))
            ])
          ]),
          default: withCtx(() => [
            _cache[25] || (_cache[25] = createBaseVNode("div", { class: "text-[12.5px] leading-relaxed text-text-4" }, " 已在运行的、开关关掉的、以及归 Docker 管的容器都会被跳过。本次会跳过开头的初始等待，立刻开始。 ", -1))
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
const BootView = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-4cecf1c6"]]);
export {
  BootView as default
};
