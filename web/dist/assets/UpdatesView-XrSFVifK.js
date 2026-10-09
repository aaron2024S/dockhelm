import { p as createLucideIcon, d as defineComponent, B as useToastStore, o as onMounted, q as onUnmounted, c as createElementBlock, a as createBaseVNode, b as createVNode, e as unref, N as Info, h as createTextVNode, t as toDisplayString, n as normalizeClass, R as RefreshCw, D as Download, g as createCommentVNode, O as CircleCheck, x as createBlock, F as Fragment, r as renderList, H as withCtx, J as RouterLink, j as ref, k as computed, y as api, m as openBlock, w as withDirectives, i as vModelCheckbox, z as openStream } from "./index-D8YHJPqb.js";
import { r as relativeTime, b as formatDateTime, s as shortImage, d as checkLabel } from "./format-mcaWYSwR.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-C598PQOA.js";
import { _ as _sfc_main$2 } from "./Modal.vue_vue_type_script_setup_true_lang-BpU2S3gf.js";
import { Z as Zap } from "./zap-X75_V2ew.js";
import { L as LoaderCircle } from "./loader-circle-DquAvDQP.js";
import { P as Play } from "./play-DTB9Nwud.js";
import { T as TriangleAlert } from "./triangle-alert-BiZ3M3e1.js";
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
const _hoisted_4 = { class: "ml-auto flex gap-2" };
const _hoisted_5 = ["disabled"];
const _hoisted_6 = ["disabled"];
const _hoisted_7 = ["disabled"];
const _hoisted_8 = { class: "grid grid-cols-2 gap-3 lg:grid-cols-4" };
const _hoisted_9 = { class: "dh-card p-3.5" };
const _hoisted_10 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_11 = { class: "dh-card p-3.5" };
const _hoisted_12 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_13 = { class: "mt-1.5 text-[22px] font-semibold leading-none" };
const _hoisted_14 = { class: "dh-card p-3.5" };
const _hoisted_15 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_16 = { class: "mt-1.5 text-[22px] font-semibold leading-none" };
const _hoisted_17 = { class: "dh-card p-3.5" };
const _hoisted_18 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_19 = { class: "mt-1.5 text-[22px] font-semibold leading-none" };
const _hoisted_20 = {
  key: 0,
  class: "dh-card"
};
const _hoisted_21 = { class: "dh-card-head" };
const _hoisted_22 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_23 = { class: "flex flex-col" };
const _hoisted_24 = { class: "w-[170px] flex-none truncate font-mono text-[11.5px] text-text-3" };
const _hoisted_25 = { class: "min-w-0 flex-1 truncate text-text-2" };
const _hoisted_26 = { class: "dh-card" };
const _hoisted_27 = { class: "dh-card-head" };
const _hoisted_28 = { class: "text-[11.5px] font-normal text-text-5" };
const _hoisted_29 = {
  key: 0,
  class: "ml-auto flex items-center gap-2"
};
const _hoisted_30 = { class: "flex cursor-pointer items-center gap-2 text-[11.5px] font-normal text-text-4" };
const _hoisted_31 = ["checked"];
const _hoisted_32 = ["disabled"];
const _hoisted_33 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_34 = { class: "dh-table" };
const _hoisted_35 = ["checked", "onChange"];
const _hoisted_36 = { class: "max-w-[220px] truncate font-mono text-[11.5px] text-text-3" };
const _hoisted_37 = { class: "font-mono text-[11px] text-text-5" };
const _hoisted_38 = { class: "font-mono text-[11px] text-text-5" };
const _hoisted_39 = { class: "dh-card" };
const _hoisted_40 = { class: "dh-card-head" };
const _hoisted_41 = { class: "ml-auto flex items-center gap-2" };
const _hoisted_42 = ["disabled"];
const _hoisted_43 = ["disabled"];
const _hoisted_44 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_45 = { class: "grid grid-cols-1 gap-3 lg:grid-cols-3" };
const _hoisted_46 = { class: "flex items-center gap-2 text-[12px] text-text-4" };
const _hoisted_47 = { class: "mt-1.5 text-[12.5px] leading-relaxed text-text-3" };
const _hoisted_48 = { class: "rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_49 = { class: "flex items-center gap-1.5 text-[12px] text-text-4" };
const _hoisted_50 = { class: "mt-1.5 text-[14px] font-semibold text-text-1" };
const _hoisted_51 = { class: "mt-1 text-[11px] text-text-5" };
const _hoisted_52 = { class: "rounded-[10px] border border-line-1 bg-ink-800" };
const _hoisted_53 = { class: "flex flex-wrap items-center gap-2 border-b border-line-1 px-3 py-2" };
const _hoisted_54 = { class: "dh-badge dh-badge-warn" };
const _hoisted_55 = { class: "dh-badge dh-badge-plain" };
const _hoisted_56 = {
  key: 0,
  class: "px-3 py-3 text-[12px] text-text-5"
};
const _hoisted_57 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_58 = { class: "dh-table" };
const _hoisted_59 = { class: "font-mono text-[11.5px] text-text-3" };
const _hoisted_60 = { class: "max-w-[220px] truncate font-mono text-[11px] text-text-5" };
const _hoisted_61 = { class: "text-[11.5px]" };
const _hoisted_62 = {
  key: 0,
  class: "dh-badge dh-badge-warn"
};
const _hoisted_63 = {
  key: 1,
  class: "dh-badge dh-badge-accent"
};
const _hoisted_64 = {
  key: 2,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_65 = {
  key: 3,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_66 = {
  key: 4,
  class: "ml-2 text-text-5"
};
const _hoisted_67 = {
  key: 0,
  class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5"
};
const _hoisted_68 = { class: "flex flex-wrap items-center gap-2 text-[12px]" };
const _hoisted_69 = { class: "font-medium text-text-3" };
const _hoisted_70 = { class: "dh-badge dh-badge-plain" };
const _hoisted_71 = {
  key: 0,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_72 = {
  key: 1,
  class: "dh-badge dh-badge-run"
};
const _hoisted_73 = {
  key: 2,
  class: "dh-badge dh-badge-err"
};
const _hoisted_74 = { class: "ml-auto text-[11.5px] text-text-5" };
const _hoisted_75 = {
  key: 1,
  class: "flex items-start gap-2 rounded-[8px] border border-line-1 px-2.5 py-2 text-[11.5px] leading-relaxed text-text-4"
};
const _hoisted_76 = {
  key: 1,
  class: "dh-card"
};
const _hoisted_77 = { class: "dh-card-head" };
const _hoisted_78 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_79 = { class: "overflow-x-auto" };
const _hoisted_80 = { class: "dh-table" };
const _hoisted_81 = { class: "flex items-center gap-1.5" };
const _hoisted_82 = {
  key: 0,
  class: "dh-badge dh-badge-accent"
};
const _hoisted_83 = { class: "max-w-[240px] truncate font-mono text-[11.5px] text-text-3" };
const _hoisted_84 = { class: "text-[11.5px] text-text-4" };
const _hoisted_85 = { class: "flex flex-col gap-3" };
const _hoisted_86 = { class: "dh-scroll max-h-[200px] overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-2.5" };
const _hoisted_87 = { class: "flex cursor-pointer items-start gap-2.5 text-[12px] text-text-2" };
const _hoisted_88 = ["disabled"];
const _hoisted_89 = { class: "flex flex-col gap-3 text-[12.5px] leading-relaxed text-text-3" };
const _hoisted_90 = { class: "flex items-start gap-2.5" };
const _hoisted_91 = { class: "flex items-start gap-2.5" };
const _hoisted_92 = { class: "flex items-start gap-2.5" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "UpdatesView",
  setup(__props) {
    const toast = useToastStore();
    const results = ref([]);
    const checkedAt = ref("");
    const selfName = ref("");
    const loading = ref(false);
    const deepLoading = ref(false);
    const applying = ref(false);
    const selected = ref(/* @__PURE__ */ new Set());
    const progress = ref([]);
    const showConfirm = ref(false);
    const forceUpdate = ref(false);
    const showInfo = ref(false);
    let closeStream = null;
    let disposed = false;
    const auto = ref(null);
    const autoBusy = ref(false);
    const available = computed(() => results.value.filter((r) => r.status === "update_available"));
    const other = computed(() => results.value.filter((r) => r.status !== "update_available"));
    const selectedNames = computed(() => [...selected.value]);
    const willUpdate = computed(() => (auto.value?.candidates ?? []).filter((c) => c.willUpdate));
    const skipped = computed(() => (auto.value?.candidates ?? []).filter((c) => !c.willUpdate && (c.protected || c.excluded)));
    function toggle(name) {
      const s = new Set(selected.value);
      if (s.has(name)) s.delete(name);
      else s.add(name);
      selected.value = s;
    }
    function selectAllAvailable() {
      if (selected.value.size === available.value.length) selected.value = /* @__PURE__ */ new Set();
      else selected.value = new Set(available.value.map((r) => r.container));
    }
    async function load() {
      try {
        const res = await api.get("/api/updates");
        results.value = res.results ?? [];
        checkedAt.value = res.checkedAt;
        selfName.value = res.selfName;
      } catch (e) {
        toast.error("读取巡检结果失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function check(deep) {
      if (deep) deepLoading.value = true;
      else loading.value = true;
      try {
        const res = await api.post("/api/updates/check", { deep });
        results.value = res.results ?? [];
        checkedAt.value = (/* @__PURE__ */ new Date()).toISOString();
        const n = (res.results ?? []).filter((r) => r.status === "update_available").length;
        toast.success(deep ? "深度检测完成" : "巡检完成", `发现 ${n} 个有可用更新`);
      } catch (e) {
        toast.error("巡检失败", e instanceof Error ? e.message : String(e));
      } finally {
        deepLoading.value = false;
        loading.value = false;
      }
    }
    let applyWatchdog;
    function armApplyWatchdog() {
      clearApplyWatchdog();
      applyWatchdog = window.setTimeout(() => {
        if (!applying.value) return;
        applying.value = false;
        toast.info("更新状态已超时解锁", "长时间没有收到批次结束事件，已重新拉取状态");
        void load();
      }, 20 * 60 * 1e3);
    }
    function clearApplyWatchdog() {
      if (applyWatchdog) {
        window.clearTimeout(applyWatchdog);
        applyWatchdog = void 0;
      }
    }
    async function apply() {
      const names = selectedNames.value;
      if (!names.length) return;
      applying.value = true;
      progress.value = [];
      try {
        await api.post("/api/updates/apply", { names, force: forceUpdate.value });
        showConfirm.value = false;
        armApplyWatchdog();
        toast.info(`已提交 ${names.length} 个容器的更新任务`, "正在后台执行，进度见下方");
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
        applying.value = false;
        clearApplyWatchdog();
      }
    }
    function releaseApply() {
      applying.value = false;
      clearApplyWatchdog();
      void load();
    }
    const toneOf = (status) => checkLabel(status).tone;
    async function loadAuto() {
      try {
        auto.value = await api.get("/api/updates/auto");
      } catch {
        auto.value = null;
      }
    }
    async function runAutoCycle(dryRun) {
      autoBusy.value = true;
      try {
        await api.post("/api/updates/auto-run", { dryRun });
        toast.info(dryRun ? "已开始巡检（不动容器）" : "已开始自动更新", "完成后本页会自动刷新");
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
        if (last && !last.dryRun) {
          toast.success("自动更新完成", `更新 ${last.updated} 个，失败 ${last.failed} 个`);
        }
      } catch (e) {
        if (disposed) return;
        toast.error("启动失败", e instanceof Error ? e.message : String(e));
      } finally {
        if (!disposed) autoBusy.value = false;
      }
    }
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
          applying.value = false;
          clearApplyWatchdog();
          selected.value = /* @__PURE__ */ new Set();
          void load();
          const msg = `更新 ${d.updated ?? 0} 个，已是最新/跳过 ${d.skipped ?? 0} 个，失败 ${d.failed ?? 0} 个`;
          toast.success("批量更新完成", msg);
        }
        if (ev.kind === "auto_check_done" || ev.kind === "auto_done") {
          void loadAuto();
          if (ev.kind === "auto_done") void load();
        }
      });
    });
    onUnmounted(() => {
      disposed = true;
      clearApplyWatchdog();
      closeStream?.();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[16] || (_cache[16] = createBaseVNode("div", { class: "dh-h1" }, "更新中心", -1)),
          createBaseVNode("button", {
            class: "dh-tap inline-flex items-center gap-1 rounded-md px-1.5 py-[3px] text-[11.5px] text-text-5 transition-colors hover:bg-ink-750 hover:text-accent",
            title: "为什么 Dockhelm 不会误停容器",
            onClick: _cache[0] || (_cache[0] = ($event) => showInfo.value = true)
          }, [
            createVNode(unref(Info), { class: "h-3.5 w-3.5" }),
            _cache[12] || (_cache[12] = createTextVNode("了解更多 ", -1))
          ]),
          createBaseVNode("div", _hoisted_3, toDisplayString(checkedAt.value ? `上次检测 ${unref(relativeTime)(checkedAt.value)} · 共比对 ${results.value.length} 个容器` : "还没有检测过"), 1),
          createBaseVNode("div", _hoisted_4, [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: loading.value || deepLoading.value,
              onClick: _cache[1] || (_cache[1] = ($event) => check(true))
            }, [
              createVNode(unref(Zap), {
                class: normalizeClass(["h-3.5 w-3.5", deepLoading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[13] || (_cache[13] = createTextVNode("深度检测 ", -1))
            ], 8, _hoisted_5),
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: loading.value || deepLoading.value,
              onClick: _cache[2] || (_cache[2] = ($event) => check(false))
            }, [
              createVNode(unref(RefreshCw), {
                class: normalizeClass(["h-3.5 w-3.5", loading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[14] || (_cache[14] = createTextVNode("重新检测 ", -1))
            ], 8, _hoisted_6),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: !available.value.length || applying.value,
              onClick: _cache[3] || (_cache[3] = ($event) => showConfirm.value = true)
            }, [
              createVNode(unref(Download), { class: "h-3.5 w-3.5" }),
              createTextVNode("更新 " + toDisplayString(available.value.length) + " 个镜像 ", 1)
            ], 8, _hoisted_7),
            applying.value ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn",
              title: "长时间没有收到批次结束事件时，点这里重新拉取状态",
              onClick: releaseApply
            }, [
              createVNode(unref(RefreshCw), { class: "h-3.5 w-3.5" }),
              _cache[15] || (_cache[15] = createTextVNode("刷新状态 ", -1))
            ])) : createCommentVNode("", true)
          ])
        ]),
        createBaseVNode("div", _hoisted_8, [
          createBaseVNode("div", _hoisted_9, [
            createBaseVNode("div", _hoisted_10, [
              createVNode(unref(Download), { class: "h-3.5 w-3.5" }),
              _cache[17] || (_cache[17] = createTextVNode("有可用更新", -1))
            ]),
            createBaseVNode("div", {
              class: normalizeClass(["mt-1.5 text-[22px] font-semibold leading-none", available.value.length ? "text-warn-text" : ""])
            }, toDisplayString(available.value.length), 3)
          ]),
          createBaseVNode("div", _hoisted_11, [
            createBaseVNode("div", _hoisted_12, [
              createVNode(unref(CircleCheck), { class: "h-3.5 w-3.5" }),
              _cache[18] || (_cache[18] = createTextVNode("已是最新", -1))
            ]),
            createBaseVNode("div", _hoisted_13, toDisplayString(results.value.filter((r) => r.status === "up_to_date").length), 1)
          ]),
          createBaseVNode("div", _hoisted_14, [
            createBaseVNode("div", _hoisted_15, [
              createVNode(unref(ShieldQuestion), { class: "h-3.5 w-3.5" }),
              _cache[19] || (_cache[19] = createTextVNode("无法判定", -1))
            ]),
            createBaseVNode("div", _hoisted_16, toDisplayString(results.value.filter((r) => r.status === "unknown").length), 1),
            _cache[20] || (_cache[20] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-6" }, "网络/认证问题导致，绝不当作有更新", -1))
          ]),
          createBaseVNode("div", _hoisted_17, [
            createBaseVNode("div", _hoisted_18, [
              createVNode(unref(CircleHelp), { class: "h-3.5 w-3.5" }),
              _cache[21] || (_cache[21] = createTextVNode("本地镜像", -1))
            ]),
            createBaseVNode("div", _hoisted_19, toDisplayString(results.value.filter((r) => r.status === "no_upstream").length), 1),
            _cache[22] || (_cache[22] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-6" }, "本地构建，无远端可比对", -1))
          ])
        ]),
        progress.value.length ? (openBlock(), createElementBlock("div", _hoisted_20, [
          createBaseVNode("div", _hoisted_21, [
            applying.value ? (openBlock(), createBlock(unref(LoaderCircle), {
              key: 0,
              class: "h-3.5 w-3.5 dh-spin text-accent"
            })) : (openBlock(), createBlock(unref(Zap), {
              key: 1,
              class: "h-3.5 w-3.5 text-text-4"
            })),
            createBaseVNode("span", null, toDisplayString(applying.value ? "更新进行中" : "最近一次更新进度"), 1),
            createBaseVNode("span", _hoisted_22, toDisplayString(progress.value.length) + " 条", 1)
          ]),
          createBaseVNode("div", _hoisted_23, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(progress.value, (p, i) => {
              return openBlock(), createElementBlock("div", {
                key: i,
                class: "flex items-center gap-2.5 border-b border-line-row px-3.5 py-2 text-[12px] last:border-b-0"
              }, [
                createBaseVNode("span", {
                  class: normalizeClass(["h-[6px] w-[6px] flex-none rounded-full", p.status === "failed" ? "bg-err" : p.status === "success" ? "bg-run" : "bg-accent"])
                }, null, 2),
                createBaseVNode("span", _hoisted_24, toDisplayString(p.name), 1),
                createBaseVNode("span", _hoisted_25, toDisplayString(p.message), 1)
              ]);
            }), 128))
          ])
        ])) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_26, [
          createBaseVNode("div", _hoisted_27, [
            createVNode(unref(Download), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[24] || (_cache[24] = createBaseVNode("span", null, "有可用更新的容器", -1)),
            createBaseVNode("span", _hoisted_28, toDisplayString(checkedAt.value ? unref(relativeTime)(checkedAt.value) + "检查" : "尚未巡检"), 1),
            available.value.length ? (openBlock(), createElementBlock("div", _hoisted_29, [
              createBaseVNode("label", _hoisted_30, [
                createBaseVNode("input", {
                  type: "checkbox",
                  class: "h-[13px] w-[13px] accent-accent",
                  checked: selected.value.size === available.value.length && available.value.length > 0,
                  onChange: selectAllAvailable
                }, null, 40, _hoisted_31),
                _cache[23] || (_cache[23] = createTextVNode(" 全选 ", -1))
              ]),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-primary",
                disabled: !selected.value.size,
                onClick: _cache[4] || (_cache[4] = ($event) => showConfirm.value = true)
              }, " 更新选中的 " + toDisplayString(selected.value.size) + " 个 ", 9, _hoisted_32)
            ])) : createCommentVNode("", true)
          ]),
          !available.value.length ? (openBlock(), createBlock(_sfc_main$1, {
            key: 0,
            icon: unref(CircleCheck),
            title: checkedAt.value ? "没有可用更新" : "还没有巡检过",
            description: checkedAt.value ? "所有可比对的容器都与仓库摘要一致。" : "点击右上角「巡检」只读检查一遍（不会动任何容器）。"
          }, null, 8, ["icon", "title", "description"])) : (openBlock(), createElementBlock("div", _hoisted_33, [
            createBaseVNode("table", _hoisted_34, [
              _cache[26] || (_cache[26] = createBaseVNode("thead", null, [
                createBaseVNode("tr", null, [
                  createBaseVNode("th", { class: "w-[34px]" }),
                  createBaseVNode("th", null, "容器"),
                  createBaseVNode("th", null, "镜像"),
                  createBaseVNode("th", { class: "w-[150px]" }, "本地摘要"),
                  createBaseVNode("th", { class: "w-[150px]" }, "仓库摘要"),
                  createBaseVNode("th", { class: "w-[110px]" }, "判定")
                ])
              ], -1)),
              createBaseVNode("tbody", null, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(available.value, (r) => {
                  return openBlock(), createElementBlock("tr", {
                    key: r.container
                  }, [
                    createBaseVNode("td", null, [
                      createBaseVNode("input", {
                        type: "checkbox",
                        class: "h-[14px] w-[14px] accent-accent",
                        checked: selected.value.has(r.container),
                        onChange: ($event) => toggle(r.container)
                      }, null, 40, _hoisted_35)
                    ]),
                    createBaseVNode("td", null, [
                      createVNode(unref(RouterLink), {
                        to: `/containers/${encodeURIComponent(r.container)}`,
                        class: "dh-tap-txt text-[12.5px] font-medium hover:text-accent"
                      }, {
                        default: withCtx(() => [
                          createTextVNode(toDisplayString(r.container), 1)
                        ]),
                        _: 2
                      }, 1032, ["to"])
                    ]),
                    createBaseVNode("td", _hoisted_36, toDisplayString(unref(shortImage)(r.image)), 1),
                    createBaseVNode("td", _hoisted_37, toDisplayString((r.localDigest || "—").slice(0, 19)), 1),
                    createBaseVNode("td", _hoisted_38, toDisplayString((r.remoteDigest || "—").slice(0, 19)), 1),
                    _cache[25] || (_cache[25] = createBaseVNode("td", null, [
                      createBaseVNode("span", { class: "dh-badge dh-badge-warn" }, "有新版本")
                    ], -1))
                  ]);
                }), 128))
              ])
            ])
          ]))
        ]),
        createBaseVNode("div", _hoisted_39, [
          createBaseVNode("div", _hoisted_40, [
            createVNode(unref(Zap), {
              class: normalizeClass(["h-3.5 w-3.5", auto.value?.enabled ? "text-warn-text" : "text-text-4"])
            }, null, 8, ["class"]),
            _cache[29] || (_cache[29] = createBaseVNode("span", null, "自动更新", -1)),
            createBaseVNode("span", {
              class: normalizeClass(["ml-2 text-[11.5px] font-normal", auto.value?.enabled ? "text-warn-text" : "text-text-5"])
            }, toDisplayString(auto.value?.enabled ? "检测到新版本会自动更新" : "只检测，不会自动动容器"), 3),
            createBaseVNode("div", _hoisted_41, [
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm",
                disabled: autoBusy.value,
                onClick: _cache[5] || (_cache[5] = ($event) => runAutoCycle(true))
              }, [
                createVNode(unref(RefreshCw), {
                  class: normalizeClass(["h-3 w-3", autoBusy.value ? "dh-spin" : ""])
                }, null, 8, ["class"]),
                _cache[27] || (_cache[27] = createTextVNode("立即巡检一轮 ", -1))
              ], 8, _hoisted_42),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-primary",
                disabled: autoBusy.value || !auto.value?.enabled,
                onClick: _cache[6] || (_cache[6] = ($event) => runAutoCycle(false))
              }, [
                createVNode(unref(Play), { class: "h-3 w-3" }),
                _cache[28] || (_cache[28] = createTextVNode("立即执行自动更新 ", -1))
              ], 8, _hoisted_43)
            ])
          ]),
          createBaseVNode("div", _hoisted_44, [
            createBaseVNode("div", _hoisted_45, [
              createBaseVNode("div", {
                class: normalizeClass(["rounded-[10px] border p-3 lg:col-span-2", auto.value?.enabled ? "border-line-warn bg-soft-warn" : "border-line-1 bg-ink-800"])
              }, [
                createBaseVNode("div", _hoisted_46, [
                  createVNode(unref(Zap), {
                    class: normalizeClass(["h-3.5 w-3.5", auto.value?.enabled ? "text-warn-text" : "text-text-4"])
                  }, null, 8, ["class"]),
                  _cache[30] || (_cache[30] = createTextVNode(" 自动更新总开关 ", -1)),
                  createBaseVNode("span", {
                    class: normalizeClass(["dh-badge", auto.value?.enabled ? "dh-badge-warn" : "dh-badge-plain"])
                  }, toDisplayString(auto.value?.enabled ? "已开启" : "已关闭"), 3)
                ]),
                createBaseVNode("div", _hoisted_47, toDisplayString(auto.value?.enabled ? "每轮巡检结束后，自动把有更新的容器（排除列表与自己除外）重建到新镜像。" : "只检测、不动手 —— 发现更新后要你在下面手动勾选并更新。"), 1),
                createVNode(unref(RouterLink), {
                  to: "/settings",
                  class: "dh-tap-txt mt-2 inline-flex items-center gap-1 text-[11.5px] text-accent hover:underline"
                }, {
                  default: withCtx(() => [..._cache[31] || (_cache[31] = [
                    createTextVNode(" 去「设置 → 更新与检测」开关它 → ", -1)
                  ])]),
                  _: 1
                })
              ], 2),
              createBaseVNode("div", _hoisted_48, [
                createBaseVNode("div", _hoisted_49, [
                  createVNode(unref(Clock), { class: "h-3.5 w-3.5" }),
                  _cache[32] || (_cache[32] = createTextVNode("下次巡检", -1))
                ]),
                createBaseVNode("div", _hoisted_50, toDisplayString(auto.value?.nextCheckAt ? unref(formatDateTime)(auto.value.nextCheckAt) : "未开启周期巡检"), 1),
                createBaseVNode("div", _hoisted_51, toDisplayString(auto.value?.lastCheckAt ? `上次 ${unref(relativeTime)(auto.value.lastCheckAt)}` : "还没有巡检记录") + " · 每 " + toDisplayString(auto.value?.checkIntervalHours ?? 0) + " 小时一次 ", 1)
              ])
            ]),
            createBaseVNode("div", _hoisted_52, [
              createBaseVNode("div", _hoisted_53, [
                _cache[34] || (_cache[34] = createBaseVNode("span", { class: "text-[12px] font-medium text-text-3" }, "本轮会发生什么", -1)),
                createBaseVNode("span", _hoisted_54, toDisplayString(willUpdate.value.length) + " 个容器将被更新", 1),
                createBaseVNode("span", _hoisted_55, toDisplayString(skipped.value.length) + " 个容器被保护/排除", 1),
                createVNode(unref(RouterLink), {
                  to: "/settings",
                  class: "dh-tap-txt ml-auto text-[11.5px] text-text-5 hover:text-accent"
                }, {
                  default: withCtx(() => [..._cache[33] || (_cache[33] = [
                    createTextVNode(" 去设置里调整检测频率与排除列表 → ", -1)
                  ])]),
                  _: 1
                })
              ]),
              !auto.value?.candidates?.length ? (openBlock(), createElementBlock("div", _hoisted_56, " 还没有巡检结果，点右上角「立即巡检一轮」先跑一次。 ")) : (openBlock(), createElementBlock("div", _hoisted_57, [
                createBaseVNode("table", _hoisted_58, [
                  _cache[35] || (_cache[35] = createBaseVNode("thead", null, [
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
                        createBaseVNode("td", _hoisted_59, toDisplayString(c.name), 1),
                        createBaseVNode("td", _hoisted_60, toDisplayString(unref(shortImage)(c.image)), 1),
                        createBaseVNode("td", null, [
                          createBaseVNode("span", {
                            class: normalizeClass(["dh-badge", c.running ? "dh-badge-run" : "dh-badge-stop"])
                          }, toDisplayString(c.running ? "运行中" : "已停止"), 3)
                        ]),
                        createBaseVNode("td", _hoisted_61, [
                          c.willUpdate ? (openBlock(), createElementBlock("span", _hoisted_62, "将更新")) : c.protected ? (openBlock(), createElementBlock("span", _hoisted_63, "受保护")) : c.excluded ? (openBlock(), createElementBlock("span", _hoisted_64, "已排除")) : (openBlock(), createElementBlock("span", _hoisted_65, "无需更新")),
                          c.reason ? (openBlock(), createElementBlock("span", _hoisted_66, toDisplayString(c.reason), 1)) : createCommentVNode("", true)
                        ])
                      ]);
                    }), 128))
                  ])
                ])
              ]))
            ]),
            auto.value?.lastRun ? (openBlock(), createElementBlock("div", _hoisted_67, [
              createBaseVNode("div", _hoisted_68, [
                createBaseVNode("span", _hoisted_69, "上一轮（" + toDisplayString(auto.value.lastRun.trigger === "manual" ? "手动触发" : "定时触发") + "）", 1),
                createBaseVNode("span", _hoisted_70, toDisplayString(unref(relativeTime)(auto.value.lastRun.startedAt)), 1),
                auto.value.lastRun.dryRun ? (openBlock(), createElementBlock("span", _hoisted_71, "仅巡检")) : (openBlock(), createElementBlock("span", _hoisted_72, "更新 " + toDisplayString(auto.value.lastRun.updated) + " 个容器", 1)),
                auto.value.lastRun.failed ? (openBlock(), createElementBlock("span", _hoisted_73, "失败 " + toDisplayString(auto.value.lastRun.failed) + " 个容器", 1)) : createCommentVNode("", true),
                createBaseVNode("span", _hoisted_74, " 巡检 " + toDisplayString(auto.value.lastRun.checked) + " 个容器 · 发现 " + toDisplayString(auto.value.lastRun.available) + " 个有更新 · 耗时 " + toDisplayString((auto.value.lastRun.durationMs / 1e3).toFixed(1)) + "s ", 1)
              ])
            ])) : createCommentVNode("", true),
            auto.value && !auto.value.enabled ? (openBlock(), createElementBlock("div", _hoisted_75, [
              createVNode(unref(Info), { class: "mt-[1px] h-3.5 w-3.5 flex-none" }),
              createBaseVNode("span", null, [
                _cache[37] || (_cache[37] = createTextVNode(" 自动更新默认关闭：检测到新版本只会打上「有新版本」标记，要不要更新由你决定。 打开开关后，每轮巡检结束就会按上面的清单自动重建容器 —— 开启前建议先看清清单。 开关与执行策略（并发度、同一镜像只拉一次、更新前备份、更新后清理）都在 ", -1)),
                createVNode(unref(RouterLink), {
                  to: "/settings",
                  class: "text-accent hover:underline"
                }, {
                  default: withCtx(() => [..._cache[36] || (_cache[36] = [
                    createTextVNode("设置页", -1)
                  ])]),
                  _: 1
                }),
                _cache[38] || (_cache[38] = createTextVNode("统一维护，本页只做预览与手动触发。 ", -1))
              ])
            ])) : createCommentVNode("", true)
          ])
        ]),
        other.value.length ? (openBlock(), createElementBlock("div", _hoisted_76, [
          createBaseVNode("div", _hoisted_77, [
            createVNode(unref(CircleCheck), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[39] || (_cache[39] = createBaseVNode("span", null, "其余容器", -1)),
            createBaseVNode("span", _hoisted_78, toDisplayString(other.value.length) + " 个", 1)
          ]),
          createBaseVNode("div", _hoisted_79, [
            createBaseVNode("table", _hoisted_80, [
              _cache[40] || (_cache[40] = createBaseVNode("thead", null, [
                createBaseVNode("tr", null, [
                  createBaseVNode("th", null, "容器"),
                  createBaseVNode("th", null, "镜像"),
                  createBaseVNode("th", { class: "w-[110px]" }, "判定"),
                  createBaseVNode("th", null, "说明")
                ])
              ], -1)),
              createBaseVNode("tbody", null, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(other.value, (r) => {
                  return openBlock(), createElementBlock("tr", {
                    key: r.container
                  }, [
                    createBaseVNode("td", null, [
                      createBaseVNode("div", _hoisted_81, [
                        createVNode(unref(RouterLink), {
                          to: `/containers/${encodeURIComponent(r.container)}`,
                          class: "dh-tap-txt text-[12.5px] font-medium hover:text-accent"
                        }, {
                          default: withCtx(() => [
                            createTextVNode(toDisplayString(r.container), 1)
                          ]),
                          _: 2
                        }, 1032, ["to"]),
                        r.container === selfName.value ? (openBlock(), createElementBlock("span", _hoisted_82, "自身")) : createCommentVNode("", true)
                      ])
                    ]),
                    createBaseVNode("td", _hoisted_83, toDisplayString(unref(shortImage)(r.image)), 1),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", `dh-badge-${toneOf(r.status)}`])
                      }, toDisplayString(unref(checkLabel)(r.status).text), 3)
                    ]),
                    createBaseVNode("td", _hoisted_84, toDisplayString(r.reason), 1)
                  ]);
                }), 128))
              ])
            ])
          ])
        ])) : createCommentVNode("", true),
        createVNode(_sfc_main$2, {
          open: showConfirm.value,
          title: "确认执行更新",
          subtitle: `共 ${selected.value.size} 个容器`,
          busy: applying.value,
          onClose: _cache[9] || (_cache[9] = ($event) => showConfirm.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[8] || (_cache[8] = ($event) => showConfirm.value = false)
            }, "取消"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: applying.value,
              onClick: apply
            }, [
              applying.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(Download), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              _cache[43] || (_cache[43] = createTextVNode(" 开始更新 ", -1))
            ], 8, _hoisted_88)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_85, [
              _cache[42] || (_cache[42] = createBaseVNode("div", { class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] leading-relaxed text-text-3" }, [
                createTextVNode(" Docker 会按顺序对每个容器执行：拉取镜像 → 比对镜像 ID。 "),
                createBaseVNode("b", { class: "text-text-1" }, "ID 没变就完全跳过"),
                createTextVNode("，容器不会被停止或重建。 只有镜像真的变化时才会走「停旧 → 改名保留 → 建新 → 健康检查（失败自动回滚）」。 ")
              ], -1)),
              createBaseVNode("div", _hoisted_86, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(selectedNames.value, (n) => {
                  return openBlock(), createElementBlock("div", {
                    key: n,
                    class: "px-1 py-[3px] font-mono text-[11.5px] text-text-3"
                  }, toDisplayString(n), 1);
                }), 128))
              ]),
              createBaseVNode("label", _hoisted_87, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[7] || (_cache[7] = ($event) => forceUpdate.value = $event),
                  type: "checkbox",
                  class: "mt-[3px] h-[14px] w-[14px] accent-warn"
                }, null, 512), [
                  [vModelCheckbox, forceUpdate.value]
                ]),
                _cache[41] || (_cache[41] = createBaseVNode("span", null, [
                  createTextVNode(" 强制重建（即使镜像 ID 未变化也重建容器） "),
                  createBaseVNode("br"),
                  createBaseVNode("span", { class: "text-[11.5px] text-text-5" }, " 一般不需要。勾选后连「已是最新」的容器也会被停掉重建 —— 这正是 dockerCopilot 曾经的行为。 ")
                ], -1))
              ])
            ])
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$2, {
          open: showInfo.value,
          title: "为什么 Dockhelm 不会误停容器",
          width: "620px",
          onClose: _cache[11] || (_cache[11] = ($event) => showInfo.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: _cache[10] || (_cache[10] = ($event) => showInfo.value = false)
            }, "明白了")
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_89, [
              createBaseVNode("div", _hoisted_90, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none text-warn-text" }),
                _cache[44] || (_cache[44] = createBaseVNode("div", null, [
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
              createBaseVNode("div", _hoisted_91, [
                createVNode(unref(CircleCheck), { class: "mt-[2px] h-4 w-4 flex-none text-run-text" }),
                _cache[45] || (_cache[45] = createBaseVNode("div", null, [
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
              createBaseVNode("div", _hoisted_92, [
                createVNode(unref(Info), { class: "mt-[2px] h-4 w-4 flex-none text-accent" }),
                _cache[46] || (_cache[46] = createBaseVNode("div", null, [
                  createBaseVNode("b", { class: "text-text-1" }, "「深度检测」是什么。"),
                  createTextVNode(" 它真的拉一次镜像再比对镜像 ID —— 这是唯一 100% 同源的判定。 镜像已最新时守护进程只会下载 manifest 与 config（几 KB），不会下载层文件。 如果你怀疑某个容器被误报，用它确认即可。 ")
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
