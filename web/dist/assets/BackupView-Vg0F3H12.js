import { x as createLucideIcon, d as defineComponent, D as useToastStore, u as useAppStore, a as onMounted, c as createElementBlock, e as createBaseVNode, t as toDisplayString, F as Fragment, w as withDirectives, M as vModelSelect, r as renderList, i as createBlock, g as unref, W as Archive, k as createTextVNode, f as createVNode, E as vModelText, j as createCommentVNode, n as normalizeClass, R as RefreshCw, H as withCtx, b as createStaticVNode, m as ref, p as computed, B as api, L as reactive, s as openBlock, N as resolveDynamicComponent, K as withModifiers, O as Layers, U as CircleCheck, z as normalizeStyle, X as vModelRadio } from "./index-i5rTgH_z.js";
import { a as formatBytes, b as relativeTime } from "./format-NZu8lNOP.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-DEzw5Eax.js";
import { _ as _sfc_main$3 } from "./Modal.vue_vue_type_script_setup_true_lang-JvIieDIo.js";
import { _ as _sfc_main$2, a as _sfc_main$4 } from "./ToggleSwitch.vue_vue_type_script_setup_true_lang-shjKTjwG.js";
import { L as LoaderCircle } from "./loader-circle-DYS8mkHw.js";
import { T as Trash2 } from "./trash-2-dzgiYvHT.js";
import { S as Save } from "./save-DwQjWTaP.js";
import { H as HardDrive } from "./hard-drive-C-fs18wQ.js";
import { T as TriangleAlert } from "./triangle-alert-BGzBla1S.js";
import { C as Copy } from "./copy-DcCnIKQT.js";
import { D as Download } from "./download-Ekma_IIF.js";
import { E as Eye } from "./eye-CaiNCax7.js";
const ChevronDown = createLucideIcon("ChevronDownIcon", [
  ["path", { d: "m6 9 6 6 6-6", key: "qrunsl" }]
]);
const ChevronRight = createLucideIcon("ChevronRightIcon", [
  ["path", { d: "m9 18 6-6-6-6", key: "mthhwq" }]
]);
const FileCode2 = createLucideIcon("FileCode2Icon", [
  ["path", { d: "M4 22h14a2 2 0 0 0 2-2V7l-5-5H6a2 2 0 0 0-2 2v4", key: "1pf5j1" }],
  ["path", { d: "M14 2v4a2 2 0 0 0 2 2h4", key: "tnqrlb" }],
  ["path", { d: "m5 12-3 3 3 3", key: "oke12k" }],
  ["path", { d: "m9 18 3-3-3-3", key: "112psh" }]
]);
const FolderTree = createLucideIcon("FolderTreeIcon", [
  [
    "path",
    {
      d: "M20 10a1 1 0 0 0 1-1V6a1 1 0 0 0-1-1h-2.5a1 1 0 0 1-.8-.4l-.9-1.2A1 1 0 0 0 15 3h-2a1 1 0 0 0-1 1v5a1 1 0 0 0 1 1Z",
      key: "hod4my"
    }
  ],
  [
    "path",
    {
      d: "M20 21a1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1h-2.9a1 1 0 0 1-.88-.55l-.42-.85a1 1 0 0 0-.92-.6H13a1 1 0 0 0-1 1v5a1 1 0 0 0 1 1Z",
      key: "w4yl2u"
    }
  ],
  ["path", { d: "M3 5a2 2 0 0 0 2 2h3", key: "f2jnh7" }],
  ["path", { d: "M3 3v13a2 2 0 0 0 2 2h3", key: "k8epm1" }]
]);
const RotateCcw = createLucideIcon("RotateCcwIcon", [
  ["path", { d: "M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8", key: "1357e3" }],
  ["path", { d: "M3 3v5h5", key: "1xhq8a" }]
]);
const Upload = createLucideIcon("UploadIcon", [
  ["path", { d: "M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4", key: "ih7n3h" }],
  ["polyline", { points: "17 8 12 3 7 8", key: "t8dd8p" }],
  ["line", { x1: "12", x2: "12", y1: "3", y2: "15", key: "widbto" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "min-w-0 flex-1" };
const _hoisted_4 = { class: "dh-sub" };
const _hoisted_5 = { class: "flex flex-wrap items-center gap-2" };
const _hoisted_6 = ["disabled"];
const _hoisted_7 = ["value"];
const _hoisted_8 = ["disabled"];
const _hoisted_9 = ["disabled"];
const _hoisted_10 = { class: "dh-seg w-fit" };
const _hoisted_11 = ["data-on"];
const _hoisted_12 = ["data-on"];
const _hoisted_13 = { class: "grid grid-cols-2 gap-3 lg:grid-cols-3" };
const _hoisted_14 = { class: "dh-card p-3.5" };
const _hoisted_15 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_16 = { class: "mt-1 text-[11px] text-text-5" };
const _hoisted_17 = { class: "dh-card p-3.5" };
const _hoisted_18 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_19 = ["title"];
const _hoisted_20 = { class: "break-all" };
const _hoisted_21 = { class: "dh-card p-3.5" };
const _hoisted_22 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_23 = { class: "mt-1 text-[11px] text-text-5" };
const _hoisted_24 = { class: "dh-card" };
const _hoisted_25 = { class: "dh-card-head" };
const _hoisted_26 = { class: "text-[11.5px] font-normal text-text-5" };
const _hoisted_27 = { class: "ml-auto flex flex-wrap items-center gap-2 font-normal" };
const _hoisted_28 = { class: "dh-seg" };
const _hoisted_29 = ["data-on"];
const _hoisted_30 = ["data-on"];
const _hoisted_31 = { class: "flex flex-wrap items-center gap-2 border-b border-line-1 px-3.5 py-2" };
const _hoisted_32 = ["disabled"];
const _hoisted_33 = ["disabled"];
const _hoisted_34 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_35 = { class: "dh-table" };
const _hoisted_36 = ["onClick"];
const _hoisted_37 = { class: "flex items-center gap-1.5 text-[12.5px] font-medium text-text-1" };
const _hoisted_38 = { class: "ml-5 text-[10.5px] text-text-6" };
const _hoisted_39 = { class: "ml-1.5 text-[11.5px] text-text-4" };
const _hoisted_40 = { class: "font-mono text-[11.5px] text-text-4" };
const _hoisted_41 = { class: "flex justify-end gap-1.5" };
const _hoisted_42 = ["onClick"];
const _hoisted_43 = ["onClick"];
const _hoisted_44 = ["onClick"];
const _hoisted_45 = ["onClick"];
const _hoisted_46 = { key: 0 };
const _hoisted_47 = {
  colspan: "5",
  class: "!py-2.5"
};
const _hoisted_48 = { class: "flex flex-col gap-1.5" };
const _hoisted_49 = { class: "text-[11px] text-text-5" };
const _hoisted_50 = { class: "grid gap-1.5 sm:grid-cols-2 lg:grid-cols-3" };
const _hoisted_51 = ["title"];
const _hoisted_52 = { class: "dh-badge dh-badge-plain" };
const _hoisted_53 = ["onClick"];
const _hoisted_54 = ["onClick"];
const _hoisted_55 = ["onClick"];
const _hoisted_56 = {
  key: 2,
  class: "overflow-x-auto"
};
const _hoisted_57 = { class: "dh-table" };
const _hoisted_58 = { class: "text-[12.5px] font-medium text-text-1" };
const _hoisted_59 = ["title"];
const _hoisted_60 = { class: "text-[12px] text-text-3" };
const _hoisted_61 = { class: "text-[10.5px] text-text-6" };
const _hoisted_62 = { class: "font-mono text-[11.5px] text-text-4" };
const _hoisted_63 = { class: "flex justify-end gap-1.5" };
const _hoisted_64 = ["onClick"];
const _hoisted_65 = ["onClick"];
const _hoisted_66 = ["onClick"];
const _hoisted_67 = { key: 0 };
const _hoisted_68 = { class: "grid gap-3 lg:grid-cols-2" };
const _hoisted_69 = {
  key: 0,
  class: "dh-card"
};
const _hoisted_70 = { class: "dh-card-head" };
const _hoisted_71 = ["disabled"];
const _hoisted_72 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_73 = { class: "dh-card" };
const _hoisted_74 = { class: "dh-card-head" };
const _hoisted_75 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_76 = { class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5" };
const _hoisted_77 = { class: "font-mono text-[12px] text-text-2" };
const _hoisted_78 = { class: "border-t border-line-1 pt-3" };
const _hoisted_79 = { class: "mb-2 flex items-center gap-2" };
const _hoisted_80 = {
  key: 0,
  class: "text-[11.5px] leading-relaxed text-text-5"
};
const _hoisted_81 = {
  key: 1,
  class: "flex flex-col gap-1.5"
};
const _hoisted_82 = { class: "font-mono text-text-3" };
const _hoisted_83 = { class: "font-mono text-text-3" };
const _hoisted_84 = {
  key: 0,
  class: "text-[11px] text-text-6"
};
const _hoisted_85 = { class: "dh-card" };
const _hoisted_86 = { class: "dh-card-head" };
const _hoisted_87 = { class: "ml-auto flex items-center gap-2 font-normal" };
const _hoisted_88 = ["disabled"];
const _hoisted_89 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_90 = { class: "dh-table" };
const _hoisted_91 = { class: "text-[12.5px] font-medium text-text-1" };
const _hoisted_92 = { class: "text-[10.5px] text-text-6" };
const _hoisted_93 = { class: "flex flex-wrap gap-1.5" };
const _hoisted_94 = ["title"];
const _hoisted_95 = ["title"];
const _hoisted_96 = {
  key: 0,
  class: "text-[11px] text-text-5"
};
const _hoisted_97 = { class: "text-[12px] text-text-3" };
const _hoisted_98 = { class: "text-[10.5px] text-text-6" };
const _hoisted_99 = {
  key: 1,
  class: "text-[11.5px] text-text-5"
};
const _hoisted_100 = { class: "flex justify-end gap-1.5" };
const _hoisted_101 = ["disabled", "onClick"];
const _hoisted_102 = ["onClick"];
const _hoisted_103 = ["onClick"];
const _hoisted_104 = ["onClick"];
const _hoisted_105 = { class: "grid gap-3 lg:grid-cols-2" };
const _hoisted_106 = { class: "dh-card" };
const _hoisted_107 = { class: "dh-card-head" };
const _hoisted_108 = { class: "flex flex-col gap-2.5 p-3.5" };
const _hoisted_109 = {
  key: 0,
  class: "border-t border-line-1 pt-2.5"
};
const _hoisted_110 = { class: "flex flex-wrap gap-1.5" };
const _hoisted_111 = ["onClick"];
const _hoisted_112 = { class: "dh-card" };
const _hoisted_113 = { class: "dh-card-head" };
const _hoisted_114 = { class: "flex flex-col gap-2.5 p-3.5" };
const _hoisted_115 = { class: "flex items-center gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-2.5 py-2" };
const _hoisted_116 = {
  key: 0,
  class: "flex flex-col gap-3"
};
const _hoisted_117 = { class: "dh-scroll max-h-[260px] overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_118 = {
  key: 1,
  class: "flex flex-col gap-3"
};
const _hoisted_119 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text" };
const _hoisted_120 = { class: "mb-2 flex items-center gap-2 text-[12px] text-text-3" };
const _hoisted_121 = {
  key: 0,
  class: "text-[12px] text-text-5"
};
const _hoisted_122 = {
  key: 1,
  class: "rounded-[9px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] text-text-4"
};
const _hoisted_123 = {
  key: 2,
  class: "dh-scroll max-h-[240px] overflow-auto rounded-[10px] border border-line-1"
};
const _hoisted_124 = { class: "dh-table" };
const _hoisted_125 = { class: "text-[12px] text-text-3" };
const _hoisted_126 = { class: "max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-accent-text" };
const _hoisted_127 = { class: "max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-warn-text" };
const _hoisted_128 = ["disabled"];
const _hoisted_129 = { class: "flex flex-col gap-3" };
const _hoisted_130 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text" };
const _hoisted_131 = {
  key: 0,
  class: "flex flex-col gap-1.5"
};
const _hoisted_132 = { class: "flex items-center gap-2 text-[11.5px] text-text-4" };
const _hoisted_133 = { class: "h-[5px] flex-1 overflow-hidden rounded-full bg-ink-800" };
const _hoisted_134 = { class: "dh-scroll max-h-[300px] overflow-auto rounded-[10px] border border-line-1" };
const _hoisted_135 = { class: "min-w-0 flex-1 truncate text-[12px] text-text-2" };
const _hoisted_136 = ["title"];
const _hoisted_137 = ["disabled"];
const _hoisted_138 = ["disabled"];
const _hoisted_139 = {
  key: 0,
  class: "flex flex-col gap-3"
};
const _hoisted_140 = { class: "overflow-hidden rounded-[10px] border border-line-1" };
const _hoisted_141 = { class: "dh-table" };
const _hoisted_142 = { class: "font-mono text-[11.5px] text-text-2" };
const _hoisted_143 = ["title"];
const _hoisted_144 = {
  key: 1,
  class: "flex flex-col gap-3"
};
const _hoisted_145 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text" };
const _hoisted_146 = { class: "flex flex-col gap-1.5" };
const _hoisted_147 = ["value"];
const _hoisted_148 = { class: "min-w-0 flex-1" };
const _hoisted_149 = { class: "text-[12.5px] text-text-2" };
const _hoisted_150 = { class: "ml-2 text-[11px] text-text-5" };
const _hoisted_151 = { class: "block truncate font-mono text-[10.5px] text-text-6" };
const _hoisted_152 = { class: "dh-badge dh-badge-plain" };
const _hoisted_153 = ["disabled", "onClick"];
const _hoisted_154 = {
  key: 0,
  class: "text-[12px] text-text-5"
};
const _hoisted_155 = ["disabled"];
const _hoisted_156 = {
  key: 0,
  class: "flex flex-col gap-3"
};
const _hoisted_157 = { class: "flex items-center gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5" };
const _hoisted_158 = { class: "min-w-0 flex-1 break-all font-mono text-[11.5px] text-accent-text" };
const _hoisted_159 = { class: "flex flex-col gap-3" };
const _hoisted_160 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-4" };
const _hoisted_161 = { class: "grid grid-cols-2 gap-3" };
const _hoisted_162 = ["disabled"];
const _hoisted_163 = ["disabled"];
const _hoisted_164 = ["disabled"];
const _hoisted_165 = { class: "flex flex-col gap-2 text-[12.5px] text-text-3" };
const _hoisted_166 = { class: "text-text-2" };
const _hoisted_167 = ["disabled"];
const _hoisted_168 = ["disabled"];
const _hoisted_169 = { class: "flex flex-col gap-3" };
const _hoisted_170 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12.5px] leading-relaxed text-warn-text" };
const _hoisted_171 = { class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] leading-relaxed text-text-3" };
const _hoisted_172 = { class: "text-[12px] text-text-4" };
const _hoisted_173 = ["disabled"];
const _hoisted_174 = ["disabled"];
const _hoisted_175 = { class: "dh-scroll max-h-[520px] overflow-auto whitespace-pre-wrap break-all rounded-[10px] border border-line-1 bg-ink-800 p-3 font-mono text-[11.5px] leading-[1.7] text-text-2" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "BackupView",
  setup(__props) {
    const toast = useToastStore();
    const app = useAppStore();
    const tab = ref("snapshots");
    const backups = ref([]);
    const stats = ref(null);
    const pathMappings = computed(() => stats.value?.pathMappings ?? []);
    const containers = ref([]);
    const projects = ref([]);
    const p0Readable = computed(() => projects.value.find((p) => p.readable?.length)?.readable ?? []);
    const loading = ref(true);
    const busy = ref("");
    const snapshotTarget = ref("");
    const view = ref("batch");
    const filter = ref("");
    const expanded = reactive({});
    const batchFull = reactive({});
    const removeTarget = ref(null);
    const removeBatch = ref(null);
    const removing = ref(false);
    const restoreTarget = ref(null);
    const diff = ref([]);
    const diffLoading = ref(false);
    const restoreResult = ref(null);
    const restoring = ref(false);
    const showRestore = ref(false);
    const batchTarget = ref(null);
    const batchRunning = ref(false);
    const batchProgress = ref(0);
    const batchResults = ref([]);
    const projectRestoreTarget = ref(null);
    const projectChosen = ref("");
    const projectRestoring = ref(false);
    const projectResult = ref(null);
    const mountHelpTarget = ref(null);
    const showFile = ref(null);
    const policy = ref(null);
    const savingPolicy = ref(false);
    const confirmPrune = ref(false);
    const batches = computed(() => {
      const map = /* @__PURE__ */ new Map();
      for (const it of backups.value) {
        const arr = map.get(it.ts);
        if (arr) arr.push(it);
        else map.set(it.ts, [it]);
      }
      const out = [];
      for (const [ts, items] of map) {
        const sorted = [...items].sort((a, b) => a.container.localeCompare(b.container));
        const head = sorted[0];
        out.push({
          ts,
          items: sorted,
          size: sorted.reduce((a, i) => a + i.size, 0),
          reason: head?.reason ?? "manual",
          created: sorted.reduce((a, i) => i.created > a ? i.created : a, head?.created ?? "")
        });
      }
      return out.sort((a, b) => b.ts.localeCompare(a.ts));
    });
    function firstOf(b) {
      const first = b.items[0];
      if (first) return first;
      return {
        container: "",
        ts: b.ts,
        path: "",
        size: 0,
        image: "",
        running: false,
        created: b.created,
        reason: b.reason
      };
    }
    function latestOf(p) {
      return p.backups?.[0] ?? null;
    }
    const flatRows = computed(() => {
      const kw = filter.value.trim().toLowerCase();
      const all = [...backups.value].sort((a, b) => b.ts.localeCompare(a.ts));
      return kw ? all.filter((i) => i.container.toLowerCase().includes(kw)) : all;
    });
    const pageSub = computed(() => {
      if (tab.value === "projects") {
        const visible = projects.value.filter((p) => (p.readable?.length ?? 0) > 0).length;
        const latest2 = projects.value.map((p) => p.backups?.[0]?.ts ?? "").sort().pop();
        return `${projects.value.length} 个 compose 项目 · ${visible} 个可备份${latest2 ? " · 最近备份 " + formatStamp(latest2) : ""}`;
      }
      const latest = batches.value[0]?.ts;
      return `快照 ${stats.value?.snapshots ?? 0} 份 · 覆盖 ${stats.value?.containers ?? 0} 个容器 · 占用 ${formatBytes(
        stats.value?.sizeBytes
      )}${latest ? " · 最近一次 " + formatStamp(latest) : ""}`;
    });
    function formatStamp(ts) {
      const m = /^(\d{4})(\d{2})(\d{2})-(\d{2})(\d{2})(\d{2})/.exec(ts);
      if (!m) return ts;
      const t = new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]), Number(m[4]), Number(m[5]), Number(m[6]));
      const now = /* @__PURE__ */ new Date();
      const hm = `${m[4]}:${m[5]}`;
      const sameDay = (a, b) => a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
      if (sameDay(t, now)) return `今天 ${hm}`;
      const y = new Date(now.getTime() - 864e5);
      if (sameDay(t, y)) return `昨天 ${hm}`;
      return `${m[2]}-${m[3]} ${hm}`;
    }
    function isFullBatch(b) {
      return b.items.length > 1;
    }
    function toggleBatch(ts) {
      expanded[ts] = !expanded[ts];
    }
    function batchItems(b) {
      if (batchFull[b.ts] || b.items.length <= 6) return b.items;
      return b.items.slice(0, 6);
    }
    async function load() {
      loading.value = true;
      const [b, s, c] = await Promise.allSettled([
        api.get("/api/backups"),
        api.get("/api/backups/stats"),
        api.get("/api/containers")
      ]);
      const errText = (e) => e instanceof Error ? e.message : String(e);
      if (b.status === "fulfilled") backups.value = b.value.backups ?? [];
      else toast.error("读取快照失败", errText(b.reason));
      if (s.status === "fulfilled") stats.value = s.value;
      else toast.error("读取备份统计失败", errText(s.reason));
      if (c.status === "fulfilled") containers.value = c.value.containers ?? [];
      else toast.error("读取容器列表失败", errText(c.reason));
      loading.value = false;
    }
    async function loadProjects() {
      try {
        const res = await api.get("/api/backups/projects");
        projects.value = res.projects ?? [];
      } catch (e) {
        toast.error("读取项目失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function loadPolicy() {
      try {
        policy.value = await api.get("/api/settings");
      } catch {
        policy.value = null;
      }
    }
    async function snapshotOne() {
      const name = snapshotTarget.value;
      if (!name) return;
      snapshotTarget.value = "";
      busy.value = "snapshot";
      try {
        const item = await api.post("/api/backups/snapshot", { container: name, reason: "manual" });
        toast.success(`已备份 ${name}`, `快照 ${item.ts} · ${formatBytes(item.size)}`);
        await load();
      } catch (e) {
        toast.error("备份失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    async function snapshotAll() {
      busy.value = "snapshot-all";
      try {
        const res = await api.post("/api/backups/snapshot-all", { reason: "manual" });
        const detail = `${res.items.length} 个容器成功 · ${formatBytes(
          res.items.reduce((a, i) => a + i.size, 0)
        )}`;
        if (res.failed?.length) {
          toast.info(`已备份 ${res.items.length}/${res.total} 个容器`, `${detail}；${res.failed.length} 个失败`);
        } else {
          toast.success(`已备份全部 ${res.items.length} 个容器`, detail);
        }
        expanded[res.ts] = true;
        await load();
      } catch (e) {
        toast.error("全量备份失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    function restoreFromError(e) {
      const p = e.payload;
      if (p && typeof p === "object" && "container" in p && "steps" in p) return p;
      return null;
    }
    let diffReq = 0;
    async function openRestore(item) {
      const req = ++diffReq;
      restoreTarget.value = item;
      restoreResult.value = null;
      showRestore.value = true;
      diffLoading.value = true;
      diff.value = [];
      try {
        const res = await api.get("/api/backups/diff", {
          container: item.container,
          ts: item.ts
        });
        if (req !== diffReq) return;
        diff.value = res.diff ?? [];
      } catch (e) {
        if (req !== diffReq) return;
        toast.error("无法比对差异", e instanceof Error ? e.message : String(e));
      } finally {
        if (req === diffReq) diffLoading.value = false;
      }
    }
    async function doRestore() {
      const item = restoreTarget.value;
      if (!item) return;
      restoring.value = true;
      try {
        const res = await api.post("/api/backups/restore", {
          container: item.container,
          ts: item.ts,
          preSnapshot: true,
          keepBackupContainer: true
        });
        restoreResult.value = res;
        if (res.ok) toast.success("还原完成");
        else toast.error("还原未成功", res.message);
        await load();
      } catch (e) {
        const res = restoreFromError(e);
        if (res) {
          restoreResult.value = res;
          toast.error("还原未成功", res.message);
        } else {
          toast.error("还原失败", e instanceof Error ? e.message : String(e));
        }
        await load();
      } finally {
        restoring.value = false;
      }
    }
    function openBatchRestore(b) {
      batchTarget.value = b;
      batchResults.value = [];
      batchProgress.value = 0;
    }
    async function runBatchRestore() {
      const b = batchTarget.value;
      if (!b || batchRunning.value) return;
      batchRunning.value = true;
      batchResults.value = b.items.map((i) => ({ container: i.container, ok: false, message: "待还原" }));
      batchProgress.value = 0;
      for (let i = 0; i < b.items.length; i++) {
        const it = b.items[i];
        if (!it) continue;
        batchResults.value[i] = { container: it.container, ok: false, message: "还原中…" };
        try {
          const res = await api.post("/api/backups/restore", {
            container: it.container,
            ts: it.ts,
            preSnapshot: true,
            keepBackupContainer: true
          });
          batchResults.value[i] = { container: it.container, ok: res.ok, message: res.message };
        } catch (e) {
          const res = restoreFromError(e);
          batchResults.value[i] = {
            container: it.container,
            ok: false,
            message: res?.message ?? (e instanceof Error ? e.message : String(e))
          };
        }
        batchProgress.value = i + 1;
      }
      batchRunning.value = false;
      await load();
      const bad = batchResults.value.filter((r) => !r.ok).length;
      if (bad) toast.error(`本批有 ${bad} 个容器没还原成功`, "原因见列表");
      else toast.success(`本批 ${batchProgress.value} 个容器已全部还原`);
    }
    async function downloadFile(url, filename) {
      const res = await fetch(url, { credentials: "same-origin" });
      if (!res.ok) {
        let msg = `下载失败（HTTP ${res.status}）`;
        try {
          const j = await res.json();
          if (j?.error) msg = j.error;
        } catch {
        }
        throw new Error(msg);
      }
      const blob = await res.blob();
      const a = document.createElement("a");
      a.href = URL.createObjectURL(blob);
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      URL.revokeObjectURL(a.href);
    }
    async function exportSnapshot(item) {
      try {
        await downloadFile(
          `/api/backups/export?container=${encodeURIComponent(item.container)}&ts=${encodeURIComponent(item.ts)}`,
          `${item.container}-${item.ts}.json`
        );
      } catch (e) {
        toast.error("导出失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function batchExport(b) {
      if (!isFullBatch(b)) {
        await exportSnapshot(firstOf(b));
        return;
      }
      try {
        await downloadFile(`/api/backups/export-batch?ts=${encodeURIComponent(b.ts)}`, `snapshots-${b.ts}.zip`);
        toast.success("已导出整批快照", `共 ${b.items.length} 个容器`);
      } catch (e) {
        toast.error("导出失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function confirmRemoveBatch() {
      const b = removeBatch.value;
      if (!b || removing.value) return;
      removing.value = true;
      let done = 0;
      try {
        for (const it of b.items) {
          try {
            await api.del(`/api/backups/${encodeURIComponent(it.container)}/${encodeURIComponent(it.ts)}`);
            done++;
          } catch {
          }
        }
        toast.success(`已删除本批 ${done}/${b.items.length} 份快照`);
        removeBatch.value = null;
        await load();
      } finally {
        removing.value = false;
      }
    }
    function batchTone(r) {
      if (r.ok) return "dh-badge-run";
      if (r.message === "待还原" || r.message === "还原中…") return "dh-badge-plain";
      return "dh-badge-err";
    }
    function batchLabel(r) {
      if (r.ok) return "成功";
      if (r.message === "待还原") return "待还原";
      if (r.message === "还原中…") return "进行中";
      return "失败";
    }
    async function confirmRemove() {
      const item = removeTarget.value;
      if (!item || removing.value) return;
      removing.value = true;
      try {
        await api.del(`/api/backups/${encodeURIComponent(item.container)}/${encodeURIComponent(item.ts)}`);
        toast.success("快照已删除");
        removeTarget.value = null;
        await load();
      } catch (e) {
        toast.error("删除失败", e instanceof Error ? e.message : String(e));
      } finally {
        removing.value = false;
      }
    }
    async function doPrune() {
      confirmPrune.value = false;
      busy.value = "prune";
      try {
        const res = await api.post("/api/backups/prune", {});
        toast.success(`清理了 ${res.removed} 份备份`, `释放 ${formatBytes(res.freedBytes)}`);
        await load();
      } catch (e) {
        toast.error("清理失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    const showImport = ref(false);
    const importContainer = ref("");
    const importTS = ref("");
    const importText = ref("");
    const importing = ref(false);
    function openImport() {
      importContainer.value = "";
      importTS.value = (/* @__PURE__ */ new Date()).toLocaleString("sv-SE").replace(/[-: ]/g, "").slice(0, 15);
      importText.value = "";
      showImport.value = true;
    }
    function guessContainer() {
      try {
        const doc = JSON.parse(importText.value);
        const name = String(doc?.inspect?.Name ?? "").replace(/^\//, "");
        if (name && !importContainer.value) importContainer.value = name;
      } catch {
      }
    }
    async function doImport() {
      if (!importContainer.value.trim() || !importTS.value.trim() || !importText.value.trim()) return;
      importing.value = true;
      try {
        const item = await api.post("/api/backups/import", {
          container: importContainer.value.trim(),
          ts: importTS.value.trim(),
          content: importText.value
        });
        toast.success("快照已导入", `${item.container} · ${item.ts}`);
        showImport.value = false;
        await load();
      } catch (e) {
        toast.error("导入失败", e instanceof Error ? e.message : String(e));
      } finally {
        importing.value = false;
      }
    }
    async function savePolicy() {
      if (!policy.value) return;
      savingPolicy.value = true;
      try {
        policy.value = await api.patch("/api/settings", {
          backupKeepPerContainer: policy.value.backupKeepPerContainer,
          backupMaxAgeDays: policy.value.backupMaxAgeDays,
          backupMaxTotalMB: policy.value.backupMaxTotalMB,
          backupKeepPreUpdate: policy.value.backupKeepPreUpdate
        });
        void app.loadSettings();
        toast.success("保留策略已保存", "每天自动清理与「清理过期」都会按它执行");
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      } finally {
        savingPolicy.value = false;
      }
    }
    async function backupProject(p) {
      busy.value = "proj:" + p.project;
      try {
        const item = await api.post("/api/backups/projects/snapshot", { project: p.project });
        toast.success(`已备份项目 ${p.project}`, `${(item.files ?? []).join("、")} · ${formatBytes(item.size)}`);
        await loadProjects();
        await load();
      } catch (e) {
        toast.error("项目备份失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    async function backupAllProjects() {
      busy.value = "proj-all";
      try {
        const res = await api.post(
          "/api/backups/projects/snapshot-all",
          {}
        );
        if (res.failed?.length) {
          toast.info(`已备份 ${res.items.length} 个项目`, `${res.failed.length} 个跳过（文件看不见）`);
        } else {
          toast.success(`已备份全部 ${res.items.length} 个项目`);
        }
        await loadProjects();
        await load();
      } catch (e) {
        toast.error("批量备份项目失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    async function downloadProjectBackup(p, b) {
      const target = p.backups?.[0];
      if (!target) return;
      const files = target.files ?? [];
      const single = files.length === 1;
      const only = files[0];
      const q = new URLSearchParams({ project: p.project, ts: target.ts });
      if (single && only) q.set("file", only);
      const name = single && only ? only : `${p.project}-${target.ts}.zip`;
      try {
        await downloadFile(`/api/backups/projects/download?${q.toString()}`, name);
        toast.success(single ? "已下载 yaml" : "已下载项目备份包", name);
      } catch (e) {
        toast.error("下载失败", e instanceof Error ? e.message : String(e));
      }
    }
    function openProjectRestore(p) {
      projectRestoreTarget.value = p;
      projectChosen.value = p.backups?.[0]?.ts ?? "";
      projectResult.value = null;
    }
    async function doProjectRestore() {
      const p = projectRestoreTarget.value;
      if (!p || !projectChosen.value) return;
      projectRestoring.value = true;
      try {
        const res = await api.post("/api/backups/projects/restore", {
          project: p.project,
          ts: projectChosen.value,
          preSnapshot: true
        });
        projectResult.value = res;
        if (res.ok) toast.success("项目文件已还原", res.message);
        else toast.error("还原未完成", res.message);
        await loadProjects();
        await load();
      } catch (e) {
        const p2 = e.payload;
        if (p2 && typeof p2 === "object" && "files" in p2) {
          projectResult.value = p2;
          toast.error("还原未完成", p2.message);
        } else {
          toast.error("还原失败", e instanceof Error ? e.message : String(e));
        }
      } finally {
        projectRestoring.value = false;
      }
    }
    async function deleteProjectBackup(p, b) {
      try {
        await api.del(
          `/api/backups/projects/${encodeURIComponent(p.project)}/${encodeURIComponent(b.ts)}`
        );
        toast.success("已删除这份备份");
        await loadProjects();
        await load();
      } catch (e) {
        toast.error("删除失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function openFile(f) {
      try {
        const res = await api.get("/api/backups/file", { path: f.hostPath });
        showFile.value = { path: res.path, content: res.content };
      } catch (e) {
        toast.error("无法读取文件", e instanceof Error ? e.message : String(e));
      }
    }
    function suggestedMount(p) {
      const host = p.unreadable?.[0] ?? p.configFiles?.[0] ?? p.workDir ?? "";
      const dir = host.includes("/") ? host.slice(0, host.lastIndexOf("/")) : "";
      if (!dir) return `- /volume1/docker/${p.project}:/host/docker/${p.project}`;
      return `- ${dir}:/host/docker/${p.project}`;
    }
    async function copyText(s) {
      try {
        await navigator.clipboard.writeText(s);
        toast.success("已复制");
      } catch {
        toast.info("复制失败", "请手动选中复制");
      }
    }
    function reasonLabel(reason) {
      switch (reason) {
        case "pre_update":
          return "更新前";
        case "scheduled":
          return "定时";
        case "manual":
          return "手动";
        case "restore-pre":
          return "还原前";
        default:
          return "未知来源";
      }
    }
    function reasonTone(reason) {
      switch (reason) {
        case "pre_update":
          return "dh-badge-warn";
        case "restore-pre":
          return "dh-badge-warn";
        case "scheduled":
          return "dh-badge-plain";
        default:
          return "dh-badge-plain";
      }
    }
    onMounted(async () => {
      await load();
      await loadProjects();
      await loadPolicy();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          createBaseVNode("div", _hoisted_3, [
            _cache[35] || (_cache[35] = createBaseVNode("div", { class: "dh-h1" }, "备份与恢复", -1)),
            createBaseVNode("div", _hoisted_4, toDisplayString(pageSub.value), 1)
          ]),
          createBaseVNode("div", _hoisted_5, [
            tab.value === "snapshots" ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
              withDirectives(createBaseVNode("select", {
                "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => snapshotTarget.value = $event),
                class: "dh-select !mb-0 !w-auto !min-w-[170px] !py-[5px] !text-[12px]",
                disabled: busy.value === "snapshot",
                onChange: snapshotOne
              }, [
                _cache[36] || (_cache[36] = createBaseVNode("option", { value: "" }, "备份单个容器…", -1)),
                (openBlock(true), createElementBlock(Fragment, null, renderList(containers.value, (c) => {
                  return openBlock(), createElementBlock("option", {
                    key: c.id,
                    value: c.name
                  }, toDisplayString(c.name), 9, _hoisted_7);
                }), 128))
              ], 40, _hoisted_6), [
                [vModelSelect, snapshotTarget.value]
              ]),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-primary",
                disabled: busy.value === "snapshot-all",
                onClick: snapshotAll
              }, [
                busy.value === "snapshot-all" ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "h-3.5 w-3.5 dh-spin"
                })) : (openBlock(), createBlock(unref(Archive), {
                  key: 1,
                  class: "h-3.5 w-3.5"
                })),
                _cache[37] || (_cache[37] = createTextVNode("立即备份全部容器 ", -1))
              ], 8, _hoisted_8)
            ], 64)) : (openBlock(), createElementBlock("button", {
              key: 1,
              class: "dh-btn dh-btn-primary",
              disabled: busy.value === "proj-all",
              onClick: backupAllProjects
            }, [
              busy.value === "proj-all" ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(Archive), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              _cache[38] || (_cache[38] = createTextVNode("备份全部项目 ", -1))
            ], 8, _hoisted_9))
          ])
        ]),
        createBaseVNode("div", _hoisted_10, [
          createBaseVNode("button", {
            "data-on": tab.value === "snapshots",
            onClick: _cache[1] || (_cache[1] = ($event) => tab.value = "snapshots")
          }, "容器配置快照", 8, _hoisted_11),
          createBaseVNode("button", {
            "data-on": tab.value === "projects",
            onClick: _cache[2] || (_cache[2] = ($event) => tab.value = "projects")
          }, "compose 项目", 8, _hoisted_12)
        ]),
        tab.value === "snapshots" ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
          createBaseVNode("div", _hoisted_13, [
            createBaseVNode("div", _hoisted_14, [
              _cache[39] || (_cache[39] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "快照份数", -1)),
              createBaseVNode("div", _hoisted_15, toDisplayString(stats.value?.snapshots ?? 0), 1),
              createBaseVNode("div", _hoisted_16, "覆盖 " + toDisplayString(stats.value?.containers ?? 0) + " 个容器", 1)
            ]),
            createBaseVNode("div", _hoisted_17, [
              _cache[40] || (_cache[40] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "占用空间", -1)),
              createBaseVNode("div", _hoisted_18, toDisplayString(unref(formatBytes)(stats.value?.sizeBytes)), 1),
              createBaseVNode("div", {
                class: "mt-1 text-[11px] text-text-5",
                title: String(stats.value?.dir ?? "")
              }, [
                createBaseVNode("span", _hoisted_20, "位于 " + toDisplayString(stats.value?.dir), 1)
              ], 8, _hoisted_19)
            ]),
            createBaseVNode("div", _hoisted_21, [
              _cache[41] || (_cache[41] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "项目配置备份", -1)),
              createBaseVNode("div", _hoisted_22, toDisplayString(stats.value?.projectSnapshots ?? 0), 1),
              createBaseVNode("div", _hoisted_23, " compose 项目 yaml · " + toDisplayString(unref(formatBytes)(stats.value?.projectSizeBytes)), 1)
            ])
          ]),
          createBaseVNode("div", _hoisted_24, [
            createBaseVNode("div", _hoisted_25, [
              createVNode(unref(Archive), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[42] || (_cache[42] = createBaseVNode("span", null, "快照", -1)),
              createBaseVNode("span", _hoisted_26, toDisplayString(view.value === "batch" ? "点击「全量」行展开该批的容器明细" : "一份快照一行"), 1),
              createBaseVNode("span", _hoisted_27, [
                createBaseVNode("div", _hoisted_28, [
                  createBaseVNode("button", {
                    "data-on": view.value === "batch",
                    onClick: _cache[3] || (_cache[3] = ($event) => view.value = "batch")
                  }, "按批次", 8, _hoisted_29),
                  createBaseVNode("button", {
                    "data-on": view.value === "container",
                    onClick: _cache[4] || (_cache[4] = ($event) => view.value = "container")
                  }, "按容器", 8, _hoisted_30)
                ]),
                view.value === "container" ? withDirectives((openBlock(), createElementBlock("input", {
                  key: 0,
                  "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => filter.value = $event),
                  class: "dh-input !mb-0 !w-[150px] !py-[5px] !text-[11.5px]",
                  placeholder: "筛选容器…"
                }, null, 512)), [
                  [vModelText, filter.value]
                ]) : createCommentVNode("", true)
              ])
            ]),
            createBaseVNode("div", _hoisted_31, [
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-ghost",
                onClick: openImport
              }, [
                createVNode(unref(Upload), { class: "h-3 w-3" }),
                _cache[43] || (_cache[43] = createTextVNode("导入备份包 ", -1))
              ]),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-ghost",
                disabled: busy.value === "prune",
                onClick: _cache[6] || (_cache[6] = ($event) => confirmPrune.value = true)
              }, [
                createVNode(unref(Trash2), { class: "h-3 w-3" }),
                _cache[44] || (_cache[44] = createTextVNode("清理过期 ", -1))
              ], 8, _hoisted_32),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-ghost ml-auto",
                disabled: loading.value,
                onClick: load
              }, [
                createVNode(unref(RefreshCw), {
                  class: normalizeClass(["h-3 w-3", loading.value ? "dh-spin" : ""])
                }, null, 8, ["class"]),
                _cache[45] || (_cache[45] = createTextVNode("刷新 ", -1))
              ], 8, _hoisted_33)
            ]),
            !backups.value.length ? (openBlock(), createBlock(_sfc_main$1, {
              key: 0,
              icon: unref(Archive),
              title: loading.value ? "正在载入…" : "还没有任何快照",
              description: "更新容器时 Dockhelm 会自动写一份快照（用于失败回滚），也可以点右上角「立即备份全部容器」把当前所有容器的配置各存一份。"
            }, null, 8, ["icon", "title"])) : view.value === "batch" ? (openBlock(), createElementBlock("div", _hoisted_34, [
              createBaseVNode("table", _hoisted_35, [
                _cache[50] || (_cache[50] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", null, "时间"),
                    createBaseVNode("th", null, "包含"),
                    createBaseVNode("th", { class: "w-[90px]" }, "来源"),
                    createBaseVNode("th", { class: "w-[90px]" }, "大小"),
                    createBaseVNode("th", { class: "w-[210px]" })
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(batches.value, (b) => {
                    return openBlock(), createElementBlock(Fragment, {
                      key: b.ts
                    }, [
                      createBaseVNode("tr", {
                        class: "cursor-pointer",
                        onClick: ($event) => toggleBatch(b.ts)
                      }, [
                        createBaseVNode("td", null, [
                          createBaseVNode("div", _hoisted_37, [
                            (openBlock(), createBlock(resolveDynamicComponent(expanded[b.ts] ? unref(ChevronDown) : unref(ChevronRight)), { class: "h-3.5 w-3.5 flex-none text-text-5" })),
                            createTextVNode(" " + toDisplayString(formatStamp(b.ts)), 1)
                          ]),
                          createBaseVNode("div", _hoisted_38, toDisplayString(unref(relativeTime)(b.created)), 1)
                        ]),
                        createBaseVNode("td", null, [
                          createBaseVNode("span", {
                            class: normalizeClass(["dh-badge", isFullBatch(b) ? "dh-badge-accent" : "dh-badge-plain"])
                          }, toDisplayString(isFullBatch(b) ? "全量" : "单个"), 3),
                          createBaseVNode("span", _hoisted_39, toDisplayString(isFullBatch(b) ? `${b.items.length} 个容器` : firstOf(b).container), 1)
                        ]),
                        createBaseVNode("td", null, [
                          createBaseVNode("span", {
                            class: normalizeClass(["dh-badge", reasonTone(b.reason)])
                          }, toDisplayString(reasonLabel(b.reason)), 3)
                        ]),
                        createBaseVNode("td", _hoisted_40, toDisplayString(unref(formatBytes)(b.size)), 1),
                        createBaseVNode("td", {
                          onClick: _cache[7] || (_cache[7] = withModifiers(() => {
                          }, ["stop"]))
                        }, [
                          createBaseVNode("div", _hoisted_41, [
                            isFullBatch(b) ? (openBlock(), createElementBlock("button", {
                              key: 0,
                              class: "dh-btn dh-btn-sm",
                              onClick: ($event) => openBatchRestore(b)
                            }, [
                              createVNode(unref(Layers), { class: "h-3 w-3" }),
                              _cache[46] || (_cache[46] = createTextVNode("还原本批 ", -1))
                            ], 8, _hoisted_42)) : (openBlock(), createElementBlock("button", {
                              key: 1,
                              class: "dh-btn dh-btn-sm",
                              onClick: ($event) => openRestore(firstOf(b))
                            }, [
                              createVNode(unref(RotateCcw), { class: "h-3 w-3" }),
                              _cache[47] || (_cache[47] = createTextVNode("还原… ", -1))
                            ], 8, _hoisted_43)),
                            createBaseVNode("button", {
                              class: "dh-btn dh-btn-sm",
                              onClick: ($event) => batchExport(b)
                            }, [
                              createVNode(unref(Download), { class: "h-3 w-3" }),
                              _cache[48] || (_cache[48] = createTextVNode("导出 ", -1))
                            ], 8, _hoisted_44),
                            createBaseVNode("button", {
                              class: "dh-btn dh-btn-sm dh-btn-danger",
                              onClick: ($event) => isFullBatch(b) ? removeBatch.value = b : removeTarget.value = firstOf(b)
                            }, [
                              createVNode(unref(Trash2), { class: "h-3 w-3" })
                            ], 8, _hoisted_45)
                          ])
                        ])
                      ], 8, _hoisted_36),
                      expanded[b.ts] ? (openBlock(), createElementBlock("tr", _hoisted_46, [
                        createBaseVNode("td", _hoisted_47, [
                          createBaseVNode("div", _hoisted_48, [
                            createBaseVNode("div", _hoisted_49, " 这一批里 " + toDisplayString(b.items.length) + " 个容器 · 每份都是完整配置，可单独还原 ", 1),
                            createBaseVNode("div", _hoisted_50, [
                              (openBlock(true), createElementBlock(Fragment, null, renderList(batchItems(b), (it) => {
                                return openBlock(), createElementBlock("div", {
                                  key: it.container,
                                  class: "flex items-center gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-2.5 py-1.5"
                                }, [
                                  createBaseVNode("span", {
                                    class: "min-w-0 flex-1 truncate text-[11.5px] text-text-3",
                                    title: it.container
                                  }, toDisplayString(it.container), 9, _hoisted_51),
                                  createBaseVNode("span", _hoisted_52, toDisplayString(unref(formatBytes)(it.size)), 1),
                                  createBaseVNode("button", {
                                    class: "dh-btn dh-btn-sm",
                                    onClick: ($event) => openRestore(it)
                                  }, [
                                    createVNode(unref(RotateCcw), { class: "h-3 w-3" }),
                                    _cache[49] || (_cache[49] = createTextVNode("还原 ", -1))
                                  ], 8, _hoisted_53),
                                  createBaseVNode("button", {
                                    class: "dh-btn dh-btn-sm dh-btn-danger",
                                    title: "只删这一份",
                                    onClick: ($event) => removeTarget.value = it
                                  }, [
                                    createVNode(unref(Trash2), { class: "h-3 w-3" })
                                  ], 8, _hoisted_54)
                                ]);
                              }), 128))
                            ]),
                            b.items.length > 6 && !batchFull[b.ts] ? (openBlock(), createElementBlock("button", {
                              key: 0,
                              class: "dh-btn dh-btn-sm dh-btn-ghost w-fit",
                              onClick: ($event) => batchFull[b.ts] = true
                            }, " 展开全部（还有 " + toDisplayString(b.items.length - 6) + " 个） ", 9, _hoisted_55)) : createCommentVNode("", true)
                          ])
                        ])
                      ])) : createCommentVNode("", true)
                    ], 64);
                  }), 128))
                ])
              ])
            ])) : (openBlock(), createElementBlock("div", _hoisted_56, [
              createBaseVNode("table", _hoisted_57, [
                _cache[54] || (_cache[54] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", null, "容器"),
                    createBaseVNode("th", { class: "w-[140px]" }, "快照时间"),
                    createBaseVNode("th", { class: "w-[80px]" }, "来源"),
                    createBaseVNode("th", { class: "w-[80px]" }, "大小"),
                    createBaseVNode("th", { class: "w-[100px]" }, "快照时状态"),
                    createBaseVNode("th", { class: "w-[220px]" })
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(flatRows.value, (item) => {
                    return openBlock(), createElementBlock("tr", {
                      key: item.container + "/" + item.ts
                    }, [
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_58, toDisplayString(item.container), 1),
                        createBaseVNode("div", {
                          class: "max-w-[220px] truncate font-mono text-[10.5px] text-text-6",
                          title: item.image
                        }, toDisplayString(item.image || "—"), 9, _hoisted_59)
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_60, toDisplayString(formatStamp(item.ts)), 1),
                        createBaseVNode("div", _hoisted_61, toDisplayString(unref(relativeTime)(item.created)), 1)
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge", reasonTone(item.reason)])
                        }, toDisplayString(reasonLabel(item.reason)), 3)
                      ]),
                      createBaseVNode("td", _hoisted_62, toDisplayString(unref(formatBytes)(item.size)), 1),
                      createBaseVNode("td", null, [
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge whitespace-nowrap", item.running ? "dh-badge-run" : "dh-badge-stop"]),
                          title: "拍下这份快照时容器是运行中还是已停止。容器现在是否在运行请看容器页"
                        }, toDisplayString(item.running ? "运行中" : "已停止"), 3)
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_63, [
                          createBaseVNode("button", {
                            class: "dh-btn dh-btn-sm",
                            onClick: ($event) => openRestore(item)
                          }, [
                            createVNode(unref(RotateCcw), { class: "h-3 w-3" }),
                            _cache[51] || (_cache[51] = createTextVNode("还原… ", -1))
                          ], 8, _hoisted_64),
                          createBaseVNode("button", {
                            class: "dh-btn dh-btn-sm",
                            onClick: ($event) => exportSnapshot(item)
                          }, [
                            createVNode(unref(Download), { class: "h-3 w-3" }),
                            _cache[52] || (_cache[52] = createTextVNode("导出 ", -1))
                          ], 8, _hoisted_65),
                          createBaseVNode("button", {
                            class: "dh-btn dh-btn-sm dh-btn-danger",
                            onClick: ($event) => removeTarget.value = item
                          }, [
                            createVNode(unref(Trash2), { class: "h-3 w-3" })
                          ], 8, _hoisted_66)
                        ])
                      ])
                    ]);
                  }), 128)),
                  !flatRows.value.length ? (openBlock(), createElementBlock("tr", _hoisted_67, [..._cache[53] || (_cache[53] = [
                    createBaseVNode("td", {
                      colspan: "6",
                      class: "text-[12px] text-text-4"
                    }, "没有匹配的快照。", -1)
                  ])])) : createCommentVNode("", true)
                ])
              ])
            ]))
          ]),
          createBaseVNode("div", _hoisted_68, [
            policy.value ? (openBlock(), createElementBlock("div", _hoisted_69, [
              createBaseVNode("div", _hoisted_70, [
                createVNode(unref(Trash2), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[56] || (_cache[56] = createBaseVNode("span", null, "保留策略", -1)),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm dh-btn-primary ml-auto",
                  disabled: savingPolicy.value,
                  onClick: savePolicy
                }, [
                  savingPolicy.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                    key: 0,
                    class: "h-3 w-3 dh-spin"
                  })) : (openBlock(), createBlock(unref(Save), {
                    key: 1,
                    class: "h-3 w-3"
                  })),
                  _cache[55] || (_cache[55] = createTextVNode("保存 ", -1))
                ], 8, _hoisted_71)
              ]),
              createBaseVNode("div", _hoisted_72, [
                createVNode(_sfc_main$2, {
                  title: "每容器保留最近份数",
                  sub: "超出后按时间滚动覆盖（更新前快照不计入这个额度）"
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("select", {
                      "onUpdate:modelValue": _cache[8] || (_cache[8] = ($event) => policy.value.backupKeepPerContainer = $event),
                      class: "dh-select !mb-0 !w-[110px] !py-[5px] !text-[11.5px]"
                    }, [..._cache[57] || (_cache[57] = [
                      createBaseVNode("option", { value: 0 }, "不限制", -1),
                      createBaseVNode("option", { value: 5 }, "5 份", -1),
                      createBaseVNode("option", { value: 10 }, "10 份", -1),
                      createBaseVNode("option", { value: 20 }, "20 份", -1),
                      createBaseVNode("option", { value: 50 }, "50 份", -1)
                    ])], 512), [
                      [
                        vModelSelect,
                        policy.value.backupKeepPerContainer,
                        void 0,
                        { number: true }
                      ]
                    ])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$2, {
                  title: "快照保留期",
                  sub: "到期自动清理"
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("select", {
                      "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => policy.value.backupMaxAgeDays = $event),
                      class: "dh-select !mb-0 !w-[110px] !py-[5px] !text-[11.5px]"
                    }, [..._cache[58] || (_cache[58] = [
                      createBaseVNode("option", { value: 0 }, "不限制", -1),
                      createBaseVNode("option", { value: 7 }, "7 天", -1),
                      createBaseVNode("option", { value: 30 }, "30 天", -1),
                      createBaseVNode("option", { value: 90 }, "90 天", -1),
                      createBaseVNode("option", { value: 365 }, "365 天", -1)
                    ])], 512), [
                      [
                        vModelSelect,
                        policy.value.backupMaxAgeDays,
                        void 0,
                        { number: true }
                      ]
                    ])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$2, {
                  title: "快照总容量上限",
                  sub: "达到上限后先清最旧的快照"
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("select", {
                      "onUpdate:modelValue": _cache[10] || (_cache[10] = ($event) => policy.value.backupMaxTotalMB = $event),
                      class: "dh-select !mb-0 !w-[110px] !py-[5px] !text-[11.5px]"
                    }, [..._cache[59] || (_cache[59] = [
                      createBaseVNode("option", { value: 0 }, "不限制", -1),
                      createBaseVNode("option", { value: 512 }, "512 MB", -1),
                      createBaseVNode("option", { value: 2048 }, "2 GB", -1),
                      createBaseVNode("option", { value: 8192 }, "8 GB", -1),
                      createBaseVNode("option", { value: 20480 }, "20 GB", -1)
                    ])], 512), [
                      [
                        vModelSelect,
                        policy.value.backupMaxTotalMB,
                        void 0,
                        { number: true }
                      ]
                    ])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$2, {
                  title: "更新前快照永不自动清理",
                  sub: "这是回滚的底牌 —— 自动更新出事后唯一能救回来的东西"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$4, {
                      modelValue: policy.value.backupKeepPreUpdate,
                      "onUpdate:modelValue": _cache[11] || (_cache[11] = ($event) => policy.value.backupKeepPreUpdate = $event),
                      label: "更新前快照永不自动清理"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                })
              ])
            ])) : createCommentVNode("", true),
            createBaseVNode("div", _hoisted_73, [
              createBaseVNode("div", _hoisted_74, [
                createVNode(unref(HardDrive), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[60] || (_cache[60] = createBaseVNode("span", null, "存储位置", -1))
              ]),
              createBaseVNode("div", _hoisted_75, [
                createBaseVNode("div", _hoisted_76, [
                  createBaseVNode("div", _hoisted_77, toDisplayString(stats.value?.dir || "/data/backups"), 1),
                  _cache[61] || (_cache[61] = createBaseVNode("div", { class: "mt-1.5 text-[11.5px] leading-relaxed text-text-5" }, [
                    createTextVNode(" 这是 Dockhelm "),
                    createBaseVNode("b", { class: "text-text-4" }, "容器内"),
                    createTextVNode("的路径。建议把宿主目录映射进来， 这样容器重建、换镜像都还在自己的快照。 ")
                  ], -1))
                ]),
                createBaseVNode("div", _hoisted_78, [
                  createBaseVNode("div", _hoisted_79, [
                    createVNode(unref(FolderTree), { class: "h-3.5 w-3.5 text-text-4" }),
                    _cache[62] || (_cache[62] = createBaseVNode("span", { class: "text-[12px] text-text-4" }, "宿主机路径映射（启动时从自身容器自动识别）", -1))
                  ]),
                  !pathMappings.value.length ? (openBlock(), createElementBlock("div", _hoisted_80, [..._cache[63] || (_cache[63] = [
                    createTextVNode(" 没有识别到任何挂载映射 —— Dockhelm 读不到宿主机的 docker 目录， compose 项目那份 yaml 也就无从备份。在 compose 里挂一行即可，例如 ", -1),
                    createBaseVNode("code", { class: "font-mono text-accent" }, "/volume1/docker:/host/docker", -1),
                    createTextVNode(" （右边叫什么名字都行，启动时会自动识别）。 ", -1)
                  ])])) : (openBlock(), createElementBlock("div", _hoisted_81, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(pathMappings.value, (m, i) => {
                      return openBlock(), createElementBlock("div", {
                        key: i,
                        class: "flex flex-wrap items-center gap-2 text-[11.5px]"
                      }, [
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge", m.visible ? "dh-badge-run" : "dh-badge-warn"])
                        }, toDisplayString(m.visible ? "可见" : "不可见"), 3),
                        createBaseVNode("code", _hoisted_82, toDisplayString(m.host), 1),
                        _cache[64] || (_cache[64] = createBaseVNode("span", { class: "text-text-6" }, "→", -1)),
                        createBaseVNode("code", _hoisted_83, toDisplayString(m.container), 1),
                        m.source === "env" ? (openBlock(), createElementBlock("span", _hoisted_84, "来自 DOCKHELM_HOST_ROOTS")) : createCommentVNode("", true)
                      ]);
                    }), 128))
                  ]))
                ])
              ])
            ])
          ])
        ], 64)) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
          createBaseVNode("div", _hoisted_85, [
            createBaseVNode("div", _hoisted_86, [
              createVNode(unref(FolderTree), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[67] || (_cache[67] = createBaseVNode("span", null, "项目", -1)),
              createBaseVNode("span", _hoisted_87, [
                _cache[66] || (_cache[66] = createBaseVNode("span", { class: "text-[11.5px] text-text-5" }, "按项目名排序", -1)),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm dh-btn-ghost",
                  disabled: loading.value,
                  onClick: loadProjects
                }, [
                  createVNode(unref(RefreshCw), {
                    class: normalizeClass(["h-3 w-3", loading.value ? "dh-spin" : ""])
                  }, null, 8, ["class"]),
                  _cache[65] || (_cache[65] = createTextVNode("刷新 ", -1))
                ], 8, _hoisted_88)
              ])
            ]),
            !projects.value.length ? (openBlock(), createBlock(_sfc_main$1, {
              key: 0,
              icon: unref(FolderTree),
              title: loading.value ? "正在载入…" : "没有发现 compose 项目",
              description: "只有由 docker compose 创建的容器才会带上项目标签。"
            }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_89, [
              createBaseVNode("table", _hoisted_90, [
                _cache[70] || (_cache[70] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", null, "项目"),
                    createBaseVNode("th", null, "文件"),
                    createBaseVNode("th", { class: "w-[150px]" }, "最近备份"),
                    createBaseVNode("th", { class: "w-[280px]" })
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(projects.value, (p) => {
                    return openBlock(), createElementBlock("tr", {
                      key: p.project
                    }, [
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_91, toDisplayString(p.project), 1),
                        createBaseVNode("div", _hoisted_92, toDisplayString(p.containers?.length ?? 0) + " 个容器", 1)
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_93, [
                          (openBlock(true), createElementBlock(Fragment, null, renderList(p.readable, (f) => {
                            return openBlock(), createElementBlock("span", {
                              key: f.path,
                              class: "inline-flex items-center gap-1 rounded-md border border-line-1 bg-ink-800 px-1.5 py-[2px] font-mono text-[10.5px] text-text-4",
                              title: f.hostPath
                            }, toDisplayString(f.name) + " ✓ ", 9, _hoisted_94);
                          }), 128)),
                          (openBlock(true), createElementBlock(Fragment, null, renderList(p.unreadable, (u) => {
                            return openBlock(), createElementBlock("span", {
                              key: u,
                              class: "inline-flex items-center gap-1 rounded-md border border-line-warn bg-soft-warn px-1.5 py-[2px] font-mono text-[10.5px] text-warn-text",
                              title: u
                            }, toDisplayString(u.split("/").pop()) + " 不可见 ", 9, _hoisted_95);
                          }), 128)),
                          !p.readable.length && !p.unreadable.length ? (openBlock(), createElementBlock("span", _hoisted_96, " 没读到任何 yaml ")) : createCommentVNode("", true)
                        ])
                      ]),
                      createBaseVNode("td", null, [
                        p.backups?.length ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                          createBaseVNode("div", _hoisted_97, toDisplayString(formatStamp(latestOf(p)?.ts ?? "")), 1),
                          createBaseVNode("div", _hoisted_98, toDisplayString(p.backups.length) + " 份历史", 1)
                        ], 64)) : (openBlock(), createElementBlock("span", _hoisted_99, "从未备份"))
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_100, [
                          createBaseVNode("button", {
                            class: "dh-btn dh-btn-sm",
                            disabled: !p.readable.length || busy.value === "proj:" + p.project,
                            onClick: ($event) => backupProject(p)
                          }, [
                            busy.value === "proj:" + p.project ? (openBlock(), createBlock(unref(LoaderCircle), {
                              key: 0,
                              class: "h-3 w-3 dh-spin"
                            })) : (openBlock(), createBlock(unref(Archive), {
                              key: 1,
                              class: "h-3 w-3"
                            })),
                            _cache[68] || (_cache[68] = createTextVNode("立即备份 ", -1))
                          ], 8, _hoisted_101),
                          p.backups?.length ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                            createBaseVNode("button", {
                              class: "dh-btn dh-btn-sm",
                              onClick: ($event) => downloadProjectBackup(p)
                            }, [
                              createVNode(unref(Download), { class: "h-3 w-3" }),
                              createTextVNode(toDisplayString((latestOf(p)?.files?.length ?? 0) === 1 ? "下载 yaml" : "下载 zip"), 1)
                            ], 8, _hoisted_102),
                            createBaseVNode("button", {
                              class: "dh-btn dh-btn-sm",
                              onClick: ($event) => openProjectRestore(p)
                            }, [
                              createVNode(unref(RotateCcw), { class: "h-3 w-3" }),
                              _cache[69] || (_cache[69] = createTextVNode("还原 ", -1))
                            ], 8, _hoisted_103)
                          ], 64)) : !p.readable.length ? (openBlock(), createElementBlock("button", {
                            key: 1,
                            class: "dh-btn dh-btn-sm",
                            onClick: ($event) => mountHelpTarget.value = p
                          }, " 看挂载方法 ", 8, _hoisted_104)) : createCommentVNode("", true)
                        ])
                      ])
                    ]);
                  }), 128))
                ])
              ])
            ]))
          ]),
          createBaseVNode("div", _hoisted_105, [
            createBaseVNode("div", _hoisted_106, [
              createBaseVNode("div", _hoisted_107, [
                createVNode(unref(FileCode2), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[71] || (_cache[71] = createBaseVNode("span", null, "一个项目的备份里有什么", -1))
              ]),
              createBaseVNode("div", _hoisted_108, [
                _cache[73] || (_cache[73] = createStaticVNode('<div class="flex items-start gap-2 text-[11.5px] leading-relaxed text-text-4"><span class="dh-badge dh-badge-plain">1</span><span>项目目录下的 yaml 文本（compose / docker-compose / override）</span></div><div class="flex items-start gap-2 text-[11.5px] leading-relaxed text-text-4"><span class="dh-badge dh-badge-plain">2</span><span><code class="font-mono">.env</code>（容器密码、端口这些变量都在这里）</span></div><div class="flex items-start gap-2 text-[11.5px] leading-relaxed text-text-4"><span class="dh-badge dh-badge-plain">3</span><span>每个文件的原路径 + 校验值，还原时逐个核对</span></div><div class="border-t border-line-1 pt-2.5 text-[11.5px] leading-relaxed text-text-5"> 不打包项目目录里的其它文件（数据库、缓存、上传的文件）—— 那些属于<b class="text-text-4">数据</b>， 仍然要你自己备份。 </div>', 4)),
                p0Readable.value.length ? (openBlock(), createElementBlock("div", _hoisted_109, [
                  _cache[72] || (_cache[72] = createBaseVNode("div", { class: "mb-1.5 text-[11.5px] text-text-5" }, "随便看一个项目的文件内容：", -1)),
                  createBaseVNode("div", _hoisted_110, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(p0Readable.value.slice(0, 4), (f) => {
                      return openBlock(), createElementBlock("button", {
                        key: f.path,
                        class: "dh-btn dh-btn-sm",
                        onClick: ($event) => openFile(f)
                      }, [
                        createVNode(unref(Eye), { class: "h-3 w-3" }),
                        createTextVNode(toDisplayString(f.name), 1)
                      ], 8, _hoisted_111);
                    }), 128))
                  ])
                ])) : createCommentVNode("", true)
              ])
            ]),
            createBaseVNode("div", _hoisted_112, [
              createBaseVNode("div", _hoisted_113, [
                createVNode(unref(TriangleAlert), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[74] || (_cache[74] = createBaseVNode("span", null, "「不可见」是什么意思", -1))
              ]),
              createBaseVNode("div", _hoisted_114, [
                _cache[77] || (_cache[77] = createBaseVNode("div", { class: "text-[11.5px] leading-relaxed text-text-4" }, [
                  createTextVNode(" 意思是 Dockhelm 在"),
                  createBaseVNode("b", { class: "text-text-3" }, "自己的容器里"),
                  createTextVNode("读不到这个文件。 把项目的宿主目录挂进 Dockhelm 容器（右边随便叫什么），重建容器后即可备份： ")
                ], -1)),
                createBaseVNode("div", _hoisted_115, [
                  _cache[76] || (_cache[76] = createBaseVNode("code", { class: "min-w-0 flex-1 truncate font-mono text-[11px] text-text-3" }, " - /volume1/docker/<项目>:/host/docker/<项目> ", -1)),
                  createBaseVNode("button", {
                    class: "dh-btn dh-btn-sm flex-none",
                    onClick: _cache[12] || (_cache[12] = ($event) => copyText("- /volume1/docker/<项目>:/host/docker/<项目>"))
                  }, [
                    createVNode(unref(Copy), { class: "h-3 w-3" }),
                    _cache[75] || (_cache[75] = createTextVNode("复制 ", -1))
                  ])
                ]),
                _cache[78] || (_cache[78] = createBaseVNode("div", { class: "text-[11.5px] leading-relaxed text-text-5" }, [
                  createTextVNode(" 宿主机的 docker 目录整体挂进来最省事，映射会在启动时自动识别。 它"),
                  createBaseVNode("b", { class: "text-text-4" }, "只影响 compose 项目页"),
                  createTextVNode("能不能读到 yaml —— 容器快照走的是 Docker 接口，一个目录都不用挂。 ")
                ], -1))
              ])
            ])
          ])
        ], 64)),
        createVNode(_sfc_main$3, {
          open: showRestore.value,
          title: "还原容器配置",
          subtitle: restoreTarget.value ? `${restoreTarget.value.container} · ${restoreTarget.value.ts}` : "",
          width: "660px",
          busy: restoring.value,
          onClose: _cache[14] || (_cache[14] = ($event) => {
            showRestore.value = false;
            restoreResult.value = null;
          })
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[13] || (_cache[13] = ($event) => {
                showRestore.value = false;
                restoreResult.value = null;
              })
            }, toDisplayString(restoreResult.value ? "关闭" : "取消"), 1),
            !restoreResult.value ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn dh-btn-danger",
              disabled: restoring.value,
              onClick: doRestore
            }, [
              createVNode(unref(RotateCcw), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(restoring.value ? "还原中…" : "确认还原"), 1)
            ], 8, _hoisted_128)) : createCommentVNode("", true)
          ]),
          default: withCtx(() => [
            restoreResult.value ? (openBlock(), createElementBlock("div", _hoisted_116, [
              createBaseVNode("div", {
                class: normalizeClass([
                  "flex items-center gap-2 rounded-[10px] border px-3 py-2.5 text-[12.5px]",
                  restoreResult.value.ok ? "border-line-ok bg-soft-ok text-run-text" : "border-line-err bg-soft-err text-err-text"
                ])
              }, [
                (openBlock(), createBlock(resolveDynamicComponent(restoreResult.value.ok ? unref(CircleCheck) : unref(TriangleAlert)), { class: "h-4 w-4 flex-none" })),
                createTextVNode(" " + toDisplayString(restoreResult.value.message), 1)
              ], 2),
              createBaseVNode("div", _hoisted_117, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(restoreResult.value.steps, (s, i) => {
                  return openBlock(), createElementBlock("div", {
                    key: i,
                    class: "py-[3px] font-mono text-[11.5px] text-text-3"
                  }, toDisplayString(s), 1);
                }), 128))
              ])
            ])) : (openBlock(), createElementBlock("div", _hoisted_118, [
              createBaseVNode("div", _hoisted_119, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none" }),
                _cache[79] || (_cache[79] = createBaseVNode("div", null, [
                  createTextVNode(" 还原会用快照里的配置创建一个"),
                  createBaseVNode("b", null, "新的同名容器"),
                  createTextVNode("，当前容器会被停止并改名成 "),
                  createBaseVNode("code", null, "<名字>__restorebak_<时间>"),
                  createTextVNode(" 保留下来（不会删除）。 还原前会自动为"),
                  createBaseVNode("b", null, "当前状态"),
                  createTextVNode("再存一份快照，改错了还能退回来。 ")
                ], -1))
              ]),
              createBaseVNode("div", null, [
                createBaseVNode("div", _hoisted_120, [
                  _cache[80] || (_cache[80] = createBaseVNode("span", null, "与当前状态的差异", -1)),
                  diffLoading.value ? (openBlock(), createBlock(unref(RefreshCw), {
                    key: 0,
                    class: "h-3 w-3 dh-spin text-text-5"
                  })) : createCommentVNode("", true)
                ]),
                diffLoading.value ? (openBlock(), createElementBlock("div", _hoisted_121, "正在比对…")) : !diff.value.length ? (openBlock(), createElementBlock("div", _hoisted_122, " 快照与当前配置没有差异。 ")) : (openBlock(), createElementBlock("div", _hoisted_123, [
                  createBaseVNode("table", _hoisted_124, [
                    _cache[81] || (_cache[81] = createBaseVNode("thead", null, [
                      createBaseVNode("tr", null, [
                        createBaseVNode("th", { class: "w-[100px]" }, "字段"),
                        createBaseVNode("th", null, "快照里的值"),
                        createBaseVNode("th", null, "当前值")
                      ])
                    ], -1)),
                    createBaseVNode("tbody", null, [
                      (openBlock(true), createElementBlock(Fragment, null, renderList(diff.value, (d) => {
                        return openBlock(), createElementBlock("tr", {
                          key: d.field
                        }, [
                          createBaseVNode("td", _hoisted_125, toDisplayString(d.field), 1),
                          createBaseVNode("td", _hoisted_126, toDisplayString(d.snapshot || "（空）"), 1),
                          createBaseVNode("td", _hoisted_127, toDisplayString(d.current || "（空）"), 1)
                        ]);
                      }), 128))
                    ])
                  ])
                ]))
              ])
            ]))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!batchTarget.value,
          title: "还原本批",
          subtitle: batchTarget.value ? `${batchTarget.value.items.length} 个容器 · 快照 ${batchTarget.value.ts}` : "",
          width: "620px",
          busy: batchRunning.value,
          onClose: _cache[16] || (_cache[16] = ($event) => batchTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: batchRunning.value,
              onClick: _cache[15] || (_cache[15] = ($event) => batchTarget.value = null)
            }, toDisplayString(batchResults.value.length ? "关闭" : "取消"), 9, _hoisted_137),
            !batchResults.value.length || batchProgress.value < (batchTarget.value?.items.length ?? 0) ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn dh-btn-danger",
              disabled: batchRunning.value,
              onClick: runBatchRestore
            }, [
              createVNode(unref(Layers), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(batchRunning.value ? "还原中…" : "开始还原"), 1)
            ], 8, _hoisted_138)) : createCommentVNode("", true)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_129, [
              createBaseVNode("div", _hoisted_130, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none" }),
                _cache[82] || (_cache[82] = createBaseVNode("div", null, [
                  createTextVNode(" 会把这个批次里的容器"),
                  createBaseVNode("b", null, "逐个"),
                  createTextVNode("按快照重建（串行执行，不并发）。每个容器还原前都会先为它 "),
                  createBaseVNode("b", null, "当前的状态"),
                  createTextVNode("存一份快照，所以单个失败也能退回去。 ")
                ], -1))
              ]),
              batchResults.value.length ? (openBlock(), createElementBlock("div", _hoisted_131, [
                createBaseVNode("div", _hoisted_132, [
                  createBaseVNode("span", null, "进度 " + toDisplayString(batchProgress.value) + " / " + toDisplayString(batchTarget.value?.items.length ?? 0), 1),
                  createBaseVNode("div", _hoisted_133, [
                    createBaseVNode("div", {
                      class: "h-full rounded-full bg-accent transition-all",
                      style: normalizeStyle({
                        width: batchProgress.value / Math.max(1, batchTarget.value?.items.length ?? 1) * 100 + "%"
                      })
                    }, null, 4)
                  ])
                ]),
                createBaseVNode("div", _hoisted_134, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(batchResults.value, (r) => {
                    return openBlock(), createElementBlock("div", {
                      key: r.container,
                      class: "flex items-center gap-2 border-b border-line-row px-3 py-2 last:border-b-0"
                    }, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", r.ok ? "dh-badge-run" : batchTone(r)])
                      }, toDisplayString(batchLabel(r)), 3),
                      createBaseVNode("span", _hoisted_135, toDisplayString(r.container), 1),
                      createBaseVNode("span", {
                        class: "max-w-[220px] truncate text-[11px] text-text-5",
                        title: r.message
                      }, toDisplayString(r.message), 9, _hoisted_136)
                    ]);
                  }), 128))
                ])
              ])) : createCommentVNode("", true)
            ])
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!projectRestoreTarget.value,
          title: "还原项目文件",
          subtitle: projectRestoreTarget.value?.project,
          width: "660px",
          busy: projectRestoring.value,
          onClose: _cache[19] || (_cache[19] = ($event) => projectRestoreTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[18] || (_cache[18] = ($event) => projectRestoreTarget.value = null)
            }, toDisplayString(projectResult.value ? "关闭" : "取消"), 1),
            !projectResult.value ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn dh-btn-danger",
              disabled: projectRestoring.value || !projectChosen.value,
              onClick: doProjectRestore
            }, [
              createVNode(unref(RotateCcw), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(projectRestoring.value ? "还原中…" : "确认还原"), 1)
            ], 8, _hoisted_155)) : createCommentVNode("", true)
          ]),
          default: withCtx(() => [
            projectResult.value ? (openBlock(), createElementBlock("div", _hoisted_139, [
              createBaseVNode("div", {
                class: normalizeClass([
                  "flex items-center gap-2 rounded-[10px] border px-3 py-2.5 text-[12.5px]",
                  projectResult.value.ok ? "border-line-ok bg-soft-ok text-run-text" : "border-line-err bg-soft-err text-err-text"
                ])
              }, [
                (openBlock(), createBlock(resolveDynamicComponent(projectResult.value.ok ? unref(CircleCheck) : unref(TriangleAlert)), { class: "h-4 w-4 flex-none" })),
                createTextVNode(" " + toDisplayString(projectResult.value.message), 1)
              ], 2),
              createBaseVNode("div", _hoisted_140, [
                createBaseVNode("table", _hoisted_141, [
                  _cache[83] || (_cache[83] = createBaseVNode("thead", null, [
                    createBaseVNode("tr", null, [
                      createBaseVNode("th", { class: "w-[160px]" }, "文件"),
                      createBaseVNode("th", null, "写回路径"),
                      createBaseVNode("th", { class: "w-[180px]" }, "结果")
                    ])
                  ], -1)),
                  createBaseVNode("tbody", null, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(projectResult.value.files, (f) => {
                      return openBlock(), createElementBlock("tr", {
                        key: f.name
                      }, [
                        createBaseVNode("td", _hoisted_142, toDisplayString(f.name), 1),
                        createBaseVNode("td", {
                          class: "max-w-[240px] truncate font-mono text-[11px] text-text-5",
                          title: f.hostPath
                        }, toDisplayString(f.hostPath), 9, _hoisted_143),
                        createBaseVNode("td", {
                          class: normalizeClass(["text-[11.5px]", f.written ? "text-run-text" : "text-err-text"])
                        }, toDisplayString(f.note), 3)
                      ]);
                    }), 128))
                  ])
                ])
              ])
            ])) : (openBlock(), createElementBlock("div", _hoisted_144, [
              createBaseVNode("div", _hoisted_145, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none" }),
                _cache[84] || (_cache[84] = createBaseVNode("div", null, [
                  createTextVNode(" 还原会把这些文件"),
                  createBaseVNode("b", null, "写回它们原来的宿主路径"),
                  createTextVNode("，覆盖现在的同名文件。动手前 Dockhelm 会先给 "),
                  createBaseVNode("b", null, "当前状态"),
                  createTextVNode("存一份项目备份，写错了还能退回来。 ")
                ], -1))
              ]),
              createBaseVNode("div", null, [
                _cache[85] || (_cache[85] = createBaseVNode("div", { class: "mb-2 text-[12px] text-text-3" }, "选择要还原的那一份", -1)),
                createBaseVNode("div", _hoisted_146, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(projectRestoreTarget.value?.backups ?? [], (b) => {
                    return openBlock(), createElementBlock("label", {
                      key: b.ts,
                      class: normalizeClass([
                        "flex cursor-pointer items-center gap-2.5 rounded-[9px] border px-3 py-2",
                        projectChosen.value === b.ts ? "border-line-accent-soft bg-soft-accent" : "border-line-1 bg-ink-800"
                      ])
                    }, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[17] || (_cache[17] = ($event) => projectChosen.value = $event),
                        type: "radio",
                        value: b.ts,
                        class: "h-[13px] w-[13px] accent-accent"
                      }, null, 8, _hoisted_147), [
                        [vModelRadio, projectChosen.value]
                      ]),
                      createBaseVNode("span", _hoisted_148, [
                        createBaseVNode("span", _hoisted_149, toDisplayString(formatStamp(b.ts)), 1),
                        createBaseVNode("span", _hoisted_150, toDisplayString(unref(relativeTime)(b.created)), 1),
                        createBaseVNode("span", _hoisted_151, toDisplayString((b.files ?? []).join("、")), 1)
                      ]),
                      createBaseVNode("span", _hoisted_152, toDisplayString(unref(formatBytes)(b.size)), 1),
                      createBaseVNode("button", {
                        class: "dh-btn dh-btn-sm dh-btn-danger flex-none",
                        disabled: projectRestoring.value,
                        title: "删除这一份备份",
                        onClick: withModifiers(($event) => projectRestoreTarget.value && deleteProjectBackup(projectRestoreTarget.value, b), ["prevent"])
                      }, [
                        createVNode(unref(Trash2), { class: "h-3 w-3" })
                      ], 8, _hoisted_153)
                    ], 2);
                  }), 128))
                ])
              ]),
              projectRestoreTarget.value?.backups?.[0] === void 0 ? (openBlock(), createElementBlock("div", _hoisted_154, " 这个项目还没有任何备份，先去上面点「立即备份」。 ")) : createCommentVNode("", true)
            ]))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!mountHelpTarget.value,
          title: "让 Dockhelm 读到这个项目的文件",
          subtitle: mountHelpTarget.value?.project,
          width: "640px",
          onClose: _cache[22] || (_cache[22] = ($event) => mountHelpTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[21] || (_cache[21] = ($event) => mountHelpTarget.value = null)
            }, "关闭")
          ]),
          default: withCtx(() => [
            mountHelpTarget.value ? (openBlock(), createElementBlock("div", _hoisted_156, [
              _cache[87] || (_cache[87] = createBaseVNode("div", { class: "text-[12.5px] leading-relaxed text-text-3" }, [
                createTextVNode(" 这个项目的文件在 Dockhelm 容器里看不见，所以备份不了。在 Dockhelm 的 compose 里加一行挂载即可 —— "),
                createBaseVNode("b", { class: "text-text-2" }, "冒号右边叫什么名字都行"),
                createTextVNode("，启动时会自动识别。 ")
              ], -1)),
              createBaseVNode("div", _hoisted_157, [
                createBaseVNode("code", _hoisted_158, toDisplayString(suggestedMount(mountHelpTarget.value)), 1),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm flex-none",
                  onClick: _cache[20] || (_cache[20] = ($event) => copyText(suggestedMount(mountHelpTarget.value)))
                }, [
                  createVNode(unref(Copy), { class: "h-3 w-3" }),
                  _cache[86] || (_cache[86] = createTextVNode("复制 ", -1))
                ])
              ]),
              _cache[88] || (_cache[88] = createBaseVNode("div", { class: "text-[11.5px] leading-relaxed text-text-5" }, [
                createTextVNode(" 改完 "),
                createBaseVNode("code", { class: "font-mono" }, "docker compose up -d"),
                createTextVNode(" 重建 Dockhelm 容器， 回到这一页就会变成可备份。宿主机的整个 docker 目录挂进来最省事。 ")
              ], -1))
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        }, 8, ["open", "subtitle"]),
        createVNode(_sfc_main$3, {
          open: showImport.value,
          title: "导入备份包",
          subtitle: "粘贴一份快照 JSON",
          width: "640px",
          busy: importing.value,
          onClose: _cache[27] || (_cache[27] = ($event) => showImport.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[26] || (_cache[26] = ($event) => showImport.value = false)
            }, "取消"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: importing.value || !importContainer.value.trim() || !importTS.value.trim() || !importText.value.trim(),
              onClick: doImport
            }, [
              importing.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(Download), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              _cache[95] || (_cache[95] = createTextVNode("导入 ", -1))
            ], 8, _hoisted_162)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_159, [
              createBaseVNode("div", _hoisted_160, [
                createVNode(unref(Download), { class: "mt-[2px] h-3.5 w-3.5 flex-none" }),
                _cache[89] || (_cache[89] = createBaseVNode("span", null, [
                  createTextVNode(" 接受 Dockhelm 写出的快照文件（内容形如 "),
                  createBaseVNode("code", null, '{"_dockhelm":{…},"inspect":{…}}'),
                  createTextVNode("）。 导入只写文件、不动任何容器；导入后就能像本地快照一样在还原弹窗里使用。 同名时间戳不会覆盖已有快照，会自动加后缀。 ")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_161, [
                createBaseVNode("div", null, [
                  _cache[90] || (_cache[90] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "容器名", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[23] || (_cache[23] = ($event) => importContainer.value = $event),
                    class: "dh-input",
                    placeholder: "redis",
                    onChange: guessContainer
                  }, null, 544), [
                    [vModelText, importContainer.value]
                  ]),
                  _cache[91] || (_cache[91] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-5" }, "决定这份快照归到哪个容器的列表下。", -1))
                ]),
                createBaseVNode("div", null, [
                  _cache[92] || (_cache[92] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "快照时间戳", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[24] || (_cache[24] = ($event) => importTS.value = $event),
                    class: "dh-input font-mono",
                    placeholder: "20261009-094200"
                  }, null, 512), [
                    [vModelText, importTS.value]
                  ]),
                  _cache[93] || (_cache[93] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-5" }, "格式 20060102-150405。", -1))
                ])
              ]),
              createBaseVNode("div", null, [
                _cache[94] || (_cache[94] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "快照 JSON", -1)),
                withDirectives(createBaseVNode("textarea", {
                  "onUpdate:modelValue": _cache[25] || (_cache[25] = ($event) => importText.value = $event),
                  rows: "10",
                  class: "dh-input dh-scroll h-auto resize-y font-mono text-[11.5px] leading-relaxed",
                  placeholder: '{"_dockhelm":{...},"inspect":{...}}',
                  onChange: guessContainer
                }, null, 544), [
                  [vModelText, importText.value]
                ])
              ])
            ])
          ]),
          _: 1
        }, 8, ["open", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!removeTarget.value,
          title: "删除这一份快照",
          subtitle: removeTarget.value ? `${removeTarget.value.container} · ${removeTarget.value.ts}` : "",
          width: "400px",
          busy: removing.value,
          onClose: _cache[29] || (_cache[29] = ($event) => removeTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[28] || (_cache[28] = ($event) => removeTarget.value = null)
            }, "取消", 8, _hoisted_163),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemove
            }, [
              removing.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(Trash2), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              createTextVNode(toDisplayString(removing.value ? "删除中…" : "确认删除"), 1)
            ], 8, _hoisted_164)
          ]),
          default: withCtx(() => [
            _cache[96] || (_cache[96] = createBaseVNode("div", { class: "text-[12.5px] text-text-3" }, " 删除后这份配置快照就无法再用于还原。已经运行中的容器不受影响。 ", -1))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!removeBatch.value,
          title: "删除这一整批快照",
          subtitle: removeBatch.value ? `${removeBatch.value.items.length} 个容器 · ${removeBatch.value.ts}` : "",
          width: "440px",
          busy: removing.value,
          onClose: _cache[31] || (_cache[31] = ($event) => removeBatch.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[30] || (_cache[30] = ($event) => removeBatch.value = null)
            }, "取消", 8, _hoisted_167),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemoveBatch
            }, [
              removing.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(Trash2), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              createTextVNode(toDisplayString(removing.value ? "删除中…" : "删除整批"), 1)
            ], 8, _hoisted_168)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_165, [
              createBaseVNode("div", null, [
                _cache[97] || (_cache[97] = createTextVNode(" 会删掉这一批里的 ", -1)),
                createBaseVNode("b", _hoisted_166, toDisplayString(removeBatch.value?.items.length ?? 0) + " 份快照", 1),
                _cache[98] || (_cache[98] = createTextVNode("（每个容器各一份）， 删掉后它们都无法再用于还原。已经运行中的容器不受影响。 ", -1))
              ]),
              _cache[99] || (_cache[99] = createBaseVNode("div", { class: "text-[11.5px] text-text-5" }, " 只想删其中某一个容器的话，先把这一行展开，在它那一行点「还原」右边的删除按钮。 ", -1))
            ])
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: confirmPrune.value,
          title: "清理过期备份",
          width: "460px",
          busy: busy.value === "prune",
          onClose: _cache[33] || (_cache[33] = ($event) => confirmPrune.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: busy.value === "prune",
              onClick: _cache[32] || (_cache[32] = ($event) => confirmPrune.value = false)
            }, "取消", 8, _hoisted_173),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: busy.value === "prune",
              onClick: doPrune
            }, [
              busy.value === "prune" ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(Trash2), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              _cache[103] || (_cache[103] = createTextVNode("确认清理 ", -1))
            ], 8, _hoisted_174)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_169, [
              createBaseVNode("div", _hoisted_170, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none" }),
                _cache[100] || (_cache[100] = createBaseVNode("div", null, [
                  createTextVNode(" 会按下面的保留策略"),
                  createBaseVNode("b", null, "直接删除备份文件"),
                  createTextVNode("，"),
                  createBaseVNode("b", null, "不可恢复"),
                  createTextVNode("。容器快照与项目备份都会被清到， 「更新前快照永不自动清理」打开时，那部分不会被碰到。 ")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_171, [
                _cache[101] || (_cache[101] = createTextVNode(" 当前策略：每容器保留最近 ", -1)),
                createBaseVNode("b", null, toDisplayString(policy.value?.backupKeepPerContainer ? policy.value.backupKeepPerContainer + " 份" : "不限"), 1),
                createTextVNode(" · 保留 " + toDisplayString(policy.value?.backupMaxAgeDays ? policy.value.backupMaxAgeDays + " 天" : "不限") + " · 总容量上限 " + toDisplayString(policy.value?.backupMaxTotalMB ? policy.value.backupMaxTotalMB + " MB" : "不限") + " ", 1),
                _cache[102] || (_cache[102] = createBaseVNode("div", { class: "mt-1 text-[11.5px] text-text-5" }, "要改策略请在上面的「保留策略」卡里改并保存。", -1))
              ]),
              createBaseVNode("div", _hoisted_172, " 当前共 " + toDisplayString(backups.value.length) + " 份容器快照（" + toDisplayString(unref(formatBytes)(stats.value?.sizeBytes)) + "）+ " + toDisplayString(stats.value?.projectSnapshots ?? 0) + " 份项目备份（" + toDisplayString(unref(formatBytes)(stats.value?.projectSizeBytes)) + "）。 ", 1)
            ])
          ]),
          _: 1
        }, 8, ["open", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!showFile.value,
          title: "文件内容",
          subtitle: showFile.value?.path,
          width: "720px",
          onClose: _cache[34] || (_cache[34] = ($event) => showFile.value = null)
        }, {
          default: withCtx(() => [
            createBaseVNode("pre", _hoisted_175, toDisplayString(showFile.value?.content), 1)
          ]),
          _: 1
        }, 8, ["open", "subtitle"])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
