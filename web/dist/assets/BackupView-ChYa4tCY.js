import { p as createLucideIcon, d as defineComponent, B as useToastStore, u as useAppStore, o as onMounted, c as createElementBlock, a as createBaseVNode, t as toDisplayString, e as unref, b as createVNode, h as createTextVNode, F as Fragment, w as withDirectives, L as vModelSelect, r as renderList, i as vModelCheckbox, V as Archive, g as createCommentVNode, n as normalizeClass, R as RefreshCw, x as createBlock, H as withCtx, W as createStaticVNode, j as ref, k as computed, y as api, m as openBlock, M as resolveDynamicComponent, O as CircleCheck, E as vModelText, X as resolveComponent } from "./index-BcifnPnm.js";
import { a as formatBytes, b as relativeTime } from "./format-CNNCbTDO.js";
import { _ as _sfc_main$2 } from "./EmptyState.vue_vue_type_script_setup_true_lang-BXLSYawR.js";
import { _ as _sfc_main$3 } from "./Modal.vue_vue_type_script_setup_true_lang-CRE-FNvg.js";
import { _ as _sfc_main$1, a as _sfc_main$4 } from "./ToggleSwitch.vue_vue_type_script_setup_true_lang-QgvFxivt.js";
import { H as HardDrive } from "./hard-drive-XRsVBbXV.js";
import { T as Trash2 } from "./trash-2-BIPQPW5w.js";
import { L as LoaderCircle } from "./loader-circle-DVmrmPX9.js";
import { S as Save } from "./save-9dA5AYTE.js";
import { E as Eye } from "./eye-Cv4Psnu1.js";
import { T as TriangleAlert } from "./triangle-alert-Cr7FcHAI.js";
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
const _hoisted_3 = { class: "dh-sub" };
const _hoisted_4 = { class: "dh-banner dh-banner-info items-start" };
const _hoisted_5 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_6 = { class: "dh-seg" };
const _hoisted_7 = ["data-on"];
const _hoisted_8 = ["data-on"];
const _hoisted_9 = ["value"];
const _hoisted_10 = ["title"];
const _hoisted_11 = ["disabled"];
const _hoisted_12 = { class: "ml-auto flex flex-wrap items-center gap-2" };
const _hoisted_13 = ["disabled"];
const _hoisted_14 = ["disabled"];
const _hoisted_15 = { class: "grid grid-cols-2 gap-3 lg:grid-cols-4" };
const _hoisted_16 = { class: "dh-card p-3.5" };
const _hoisted_17 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_18 = { class: "mt-1 text-[11px] text-text-5" };
const _hoisted_19 = { class: "dh-card p-3.5" };
const _hoisted_20 = { class: "mt-1.5 text-[20px] font-semibold leading-none" };
const _hoisted_21 = ["title"];
const _hoisted_22 = { class: "break-all" };
const _hoisted_23 = { class: "dh-card p-3.5" };
const _hoisted_24 = { class: "mt-1.5" };
const _hoisted_25 = { class: "mt-1 truncate font-mono text-[11px] text-text-5" };
const _hoisted_26 = { class: "dh-card p-3.5" };
const _hoisted_27 = { class: "mt-1.5" };
const _hoisted_28 = { class: "mt-1 text-[11px] leading-relaxed text-text-5" };
const _hoisted_29 = { class: "dh-card" };
const _hoisted_30 = { class: "dh-card-head" };
const _hoisted_31 = { class: "dh-card-body" };
const _hoisted_32 = {
  key: 0,
  class: "text-[12px] leading-relaxed text-text-4"
};
const _hoisted_33 = {
  key: 1,
  class: "space-y-1.5"
};
const _hoisted_34 = { class: "font-mono text-[11.5px] text-text-3" };
const _hoisted_35 = { class: "font-mono text-[11.5px] text-accent" };
const _hoisted_36 = {
  key: 0,
  class: "text-[11px] text-text-5"
};
const _hoisted_37 = {
  key: 1,
  class: "text-[11px] text-text-5"
};
const _hoisted_38 = {
  key: 0,
  class: "dh-card"
};
const _hoisted_39 = { class: "dh-card-head" };
const _hoisted_40 = ["disabled"];
const _hoisted_41 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_42 = {
  key: 1,
  class: "dh-card"
};
const _hoisted_43 = { class: "dh-card-head" };
const _hoisted_44 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_45 = { class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5" };
const _hoisted_46 = { class: "font-mono text-[12px] text-text-2" };
const _hoisted_47 = {
  key: 0,
  class: "border-t border-line-1 pt-3"
};
const _hoisted_48 = { class: "flex flex-col gap-1.5" };
const _hoisted_49 = { class: "font-mono text-text-3" };
const _hoisted_50 = { class: "font-mono text-text-3" };
const _hoisted_51 = {
  key: 2,
  class: "dh-card"
};
const _hoisted_52 = { class: "dh-card-head" };
const _hoisted_53 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_54 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_55 = { class: "dh-table" };
const _hoisted_56 = { class: "text-[12.5px] font-medium text-text-1" };
const _hoisted_57 = ["title"];
const _hoisted_58 = { class: "text-[12px] text-text-3" };
const _hoisted_59 = { class: "font-mono text-[10.5px] text-text-6" };
const _hoisted_60 = { class: "font-mono text-[11.5px] text-text-4" };
const _hoisted_61 = { class: "flex justify-end gap-1.5" };
const _hoisted_62 = ["onClick"];
const _hoisted_63 = ["onClick"];
const _hoisted_64 = {
  key: 3,
  class: "dh-card"
};
const _hoisted_65 = { class: "dh-card-head" };
const _hoisted_66 = {
  key: 1,
  class: "flex flex-col"
};
const _hoisted_67 = { class: "flex flex-wrap items-center gap-2" };
const _hoisted_68 = { class: "text-[13px] font-semibold" };
const _hoisted_69 = { class: "dh-badge dh-badge-plain" };
const _hoisted_70 = ["title"];
const _hoisted_71 = { class: "mt-1.5 flex flex-wrap gap-1" };
const _hoisted_72 = {
  key: 0,
  class: "mt-2.5 flex flex-col gap-1.5"
};
const _hoisted_73 = { class: "min-w-0 flex-1" };
const _hoisted_74 = { class: "text-[12px] text-text-2" };
const _hoisted_75 = { class: "truncate font-mono text-[10.5px] text-text-6" };
const _hoisted_76 = { class: "text-[11px] text-text-5" };
const _hoisted_77 = ["onClick"];
const _hoisted_78 = {
  key: 1,
  class: "mt-2.5 flex flex-col gap-1.5 rounded-[9px] border border-line-warn bg-soft-warn px-2.5 py-2"
};
const _hoisted_79 = { class: "flex items-center gap-2 text-[11.5px] text-warn-text" };
const _hoisted_80 = ["title"];
const _hoisted_81 = {
  key: 0,
  class: "flex flex-col gap-3"
};
const _hoisted_82 = { class: "dh-scroll max-h-[260px] overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_83 = {
  key: 1,
  class: "flex flex-col gap-3"
};
const _hoisted_84 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12px] leading-relaxed text-warn-text" };
const _hoisted_85 = { class: "mb-2 flex items-center gap-2 text-[12px] text-text-3" };
const _hoisted_86 = {
  key: 0,
  class: "text-[12px] text-text-5"
};
const _hoisted_87 = {
  key: 1,
  class: "rounded-[9px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] text-text-4"
};
const _hoisted_88 = {
  key: 2,
  class: "dh-scroll max-h-[240px] overflow-auto rounded-[10px] border border-line-1"
};
const _hoisted_89 = { class: "dh-table" };
const _hoisted_90 = { class: "text-[12px] text-text-3" };
const _hoisted_91 = { class: "max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-accent-text" };
const _hoisted_92 = { class: "max-w-[220px] whitespace-pre-wrap break-all font-mono text-[11px] text-warn-text" };
const _hoisted_93 = {
  key: 0,
  class: "border-t border-line-1 pt-3"
};
const _hoisted_94 = ["disabled"];
const _hoisted_95 = { class: "text-[12.5px]" };
const _hoisted_96 = { class: "mt-[2px] text-[11px] leading-relaxed text-text-5" };
const _hoisted_97 = { class: "font-mono text-text-4" };
const _hoisted_98 = {
  key: 0,
  class: "mt-2 flex flex-col gap-1"
};
const _hoisted_99 = { class: "flex-none font-mono text-text-4" };
const _hoisted_100 = { class: "text-text-6" };
const _hoisted_101 = {
  key: 1,
  class: "mt-2.5 flex items-start gap-2 rounded-[9px] border border-line-err bg-soft-err px-3 py-2 text-[11.5px] leading-relaxed text-err-text"
};
const _hoisted_102 = ["disabled"];
const _hoisted_103 = { class: "flex flex-col gap-3" };
const _hoisted_104 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-4" };
const _hoisted_105 = { class: "grid grid-cols-2 gap-3" };
const _hoisted_106 = ["disabled"];
const _hoisted_107 = ["disabled"];
const _hoisted_108 = ["disabled"];
const _hoisted_109 = { class: "flex flex-col gap-3" };
const _hoisted_110 = { class: "flex items-start gap-2.5 rounded-[10px] border border-line-warn bg-soft-warn px-3 py-2.5 text-[12.5px] leading-relaxed text-warn-text" };
const _hoisted_111 = { class: "rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[12px] leading-relaxed text-text-3" };
const _hoisted_112 = { class: "text-[12px] text-text-4" };
const _hoisted_113 = ["disabled"];
const _hoisted_114 = ["disabled"];
const _hoisted_115 = { class: "dh-scroll max-h-[520px] overflow-auto whitespace-pre-wrap break-all rounded-[10px] border border-line-1 bg-ink-800 p-3 font-mono text-[11.5px] leading-[1.7] text-text-2" };
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
    const loading = ref(true);
    const busy = ref("");
    const snapshotTarget = ref("");
    const snapshotWithVolumes = ref(false);
    const removeTarget = ref(null);
    const restoreTarget = ref(null);
    const diff = ref([]);
    const diffLoading = ref(false);
    const restoreResult = ref(null);
    const restoring = ref(false);
    const removing = ref(false);
    const confirmPrune = ref(false);
    const showRestore = ref(false);
    const volMeta = ref(null);
    const withVolumes = ref(false);
    const showFile = ref(null);
    const pruneKeep = ref(10);
    const pruneDays = ref(30);
    const policy = ref(null);
    const savingPolicy = ref(false);
    const rows = computed(() => [...backups.value].sort((a, b) => b.ts.localeCompare(a.ts)));
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
    async function doSnapshot() {
      const name = snapshotTarget.value;
      if (!name) return;
      busy.value = "snapshot";
      try {
        const item = await api.post("/api/backups/snapshot", {
          container: name,
          reason: "manual",
          withVolumes: snapshotWithVolumes.value
        });
        toast.success(
          `已备份 ${name}`,
          snapshotWithVolumes.value ? `快照 ${item.ts} · ${formatBytes(item.size)}${item.withData ? "（含卷数据）" : "（卷数据未打包，见快照说明）"}` : `快照 ${item.ts} · ${formatBytes(item.size)}`
        );
        await load();
      } catch (e) {
        toast.error("备份失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    let diffReq = 0;
    async function openRestore(item) {
      const req = ++diffReq;
      restoreTarget.value = item;
      restoreResult.value = null;
      showRestore.value = true;
      diffLoading.value = true;
      diff.value = [];
      volMeta.value = null;
      withVolumes.value = false;
      try {
        const res = await api.get("/api/backups/diff", {
          container: item.container,
          ts: item.ts
        });
        if (req !== diffReq) return;
        diff.value = res.diff ?? [];
        volMeta.value = res.volumes ?? null;
      } catch (e) {
        if (req !== diffReq) return;
        toast.error("无法比对差异", e instanceof Error ? e.message : String(e));
      } finally {
        if (req === diffReq) diffLoading.value = false;
      }
    }
    const packedVolumes = computed(() => (volMeta.value?.items ?? []).filter((v) => v.packed));
    const unpackedVolumes = computed(() => (volMeta.value?.items ?? []).filter((v) => !v.packed));
    async function doRestore() {
      const item = restoreTarget.value;
      if (!item) return;
      restoring.value = true;
      try {
        const res = await api.post("/api/backups/restore", {
          container: item.container,
          ts: item.ts,
          preSnapshot: true,
          keepBackupContainer: true,
          withVolumes: withVolumes.value
        });
        restoreResult.value = res;
        if (res.ok) toast.success("还原完成");
        else toast.error("还原未成功", res.message);
        await load();
      } catch (e) {
        toast.error("还原失败", e instanceof Error ? e.message : String(e));
      } finally {
        restoring.value = false;
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
        toast.success(`清理了 ${res.removed} 份快照`, `释放 ${formatBytes(res.freedBytes)}`);
        await load();
      } catch (e) {
        toast.error("清理失败", e instanceof Error ? e.message : String(e));
      } finally {
        busy.value = "";
      }
    }
    async function loadPolicy() {
      try {
        policy.value = await api.get("/api/settings");
        if (policy.value) {
          pruneKeep.value = policy.value.backupKeepPerContainer;
          pruneDays.value = policy.value.backupMaxAgeDays;
        }
      } catch {
        policy.value = null;
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
        pruneKeep.value = policy.value.backupKeepPerContainer;
        pruneDays.value = policy.value.backupMaxAgeDays;
        void app.loadSettings();
        toast.success("保留策略已保存", "每天自动清理与「清理过期」都会按它执行");
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      } finally {
        savingPolicy.value = false;
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
    function reasonLabel(reason) {
      switch (reason) {
        case "pre_update":
          return "更新前";
        case "scheduled":
          return "定时";
        case "manual":
          return "手动";
        default:
          return "未知来源";
      }
    }
    function reasonTone(reason) {
      switch (reason) {
        case "pre_update":
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
      const _component_Download = resolveComponent("Download");
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[22] || (_cache[22] = createBaseVNode("div", { class: "dh-h1" }, "备份与恢复", -1)),
          createBaseVNode("div", _hoisted_3, " 快照 " + toDisplayString(stats.value?.snapshots ?? 0) + " 份 · 占用 " + toDisplayString(unref(formatBytes)(stats.value?.sizeBytes)) + " · 覆盖 " + toDisplayString(stats.value?.containers ?? 0) + " 个容器 ", 1)
        ]),
        createBaseVNode("div", _hoisted_4, [
          createVNode(unref(HardDrive), { class: "mt-[2px] h-4 w-4 flex-none" }),
          _cache[23] || (_cache[23] = createBaseVNode("div", { class: "min-w-0 flex-1 leading-relaxed" }, [
            createTextVNode(" Dockhelm 备份的是"),
            createBaseVNode("b", null, "容器配置快照"),
            createTextVNode("（"),
            createBaseVNode("code", null, "docker inspect"),
            createTextVNode(" 的结果），可以一键还原成同名容器。 "),
            createBaseVNode("b", null, "绑定挂载的数据在宿主机目录上"),
            createTextVNode("，Dockhelm 容器默认看不见 —— 要备份那份数据，得把宿主目录也挂进来。 compose 项目真正的来源是那份 "),
            createBaseVNode("code", null, "yaml"),
            createTextVNode("，用「compose 项目」标签页查看。 ")
          ], -1))
        ]),
        createBaseVNode("div", _hoisted_5, [
          createBaseVNode("div", _hoisted_6, [
            createBaseVNode("button", {
              "data-on": tab.value === "snapshots",
              onClick: _cache[0] || (_cache[0] = ($event) => tab.value = "snapshots")
            }, "容器配置快照", 8, _hoisted_7),
            createBaseVNode("button", {
              "data-on": tab.value === "projects",
              onClick: _cache[1] || (_cache[1] = ($event) => tab.value = "projects")
            }, "compose 项目", 8, _hoisted_8)
          ]),
          tab.value === "snapshots" ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
            withDirectives(createBaseVNode("select", {
              "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => snapshotTarget.value = $event),
              class: "dh-select !w-auto !min-w-[190px]"
            }, [
              _cache[24] || (_cache[24] = createBaseVNode("option", { value: "" }, "选择要备份的容器…", -1)),
              (openBlock(true), createElementBlock(Fragment, null, renderList(containers.value, (c) => {
                return openBlock(), createElementBlock("option", {
                  key: c.id,
                  value: c.name
                }, toDisplayString(c.name), 9, _hoisted_9);
              }), 128))
            ], 512), [
              [vModelSelect, snapshotTarget.value]
            ]),
            createBaseVNode("label", {
              class: "flex cursor-pointer items-center gap-1.5 text-[11.5px] text-text-4",
              title: stats.value?.volumeRootMounted ? "连 named volume 的数据一起打包（体积可能很大）" : "宿主机的 /var/lib/docker/volumes 没有映射进来，卷数据打不了包"
            }, [
              withDirectives(createBaseVNode("input", {
                "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => snapshotWithVolumes.value = $event),
                type: "checkbox",
                class: "h-[13px] w-[13px] accent-warn"
              }, null, 512), [
                [vModelCheckbox, snapshotWithVolumes.value]
              ]),
              _cache[25] || (_cache[25] = createTextVNode(" 含卷数据 ", -1))
            ], 8, _hoisted_10),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: !snapshotTarget.value || busy.value === "snapshot",
              onClick: doSnapshot
            }, [
              createVNode(unref(Archive), { class: "h-3.5 w-3.5" }),
              _cache[26] || (_cache[26] = createTextVNode("立即备份 ", -1))
            ], 8, _hoisted_11)
          ], 64)) : createCommentVNode("", true),
          createBaseVNode("div", _hoisted_12, [
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              onClick: openImport
            }, [
              createVNode(unref(Upload), { class: "h-3 w-3" }),
              _cache[27] || (_cache[27] = createTextVNode("导入备份包 ", -1))
            ]),
            tab.value === "snapshots" ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn dh-btn-sm",
              disabled: busy.value === "prune",
              onClick: _cache[4] || (_cache[4] = ($event) => confirmPrune.value = true)
            }, [
              createVNode(unref(Trash2), { class: "h-3 w-3" }),
              _cache[28] || (_cache[28] = createTextVNode("清理过期 ", -1))
            ], 8, _hoisted_13)) : createCommentVNode("", true),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              disabled: loading.value,
              onClick: load
            }, [
              createVNode(unref(RefreshCw), {
                class: normalizeClass(["h-3 w-3", loading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[29] || (_cache[29] = createTextVNode("刷新 ", -1))
            ], 8, _hoisted_14)
          ])
        ]),
        createBaseVNode("div", _hoisted_15, [
          createBaseVNode("div", _hoisted_16, [
            _cache[30] || (_cache[30] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "快照数量", -1)),
            createBaseVNode("div", _hoisted_17, toDisplayString(stats.value?.snapshots ?? 0), 1),
            createBaseVNode("div", _hoisted_18, "覆盖 " + toDisplayString(stats.value?.containers ?? 0) + " 个容器", 1)
          ]),
          createBaseVNode("div", _hoisted_19, [
            _cache[31] || (_cache[31] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "占用空间", -1)),
            createBaseVNode("div", _hoisted_20, toDisplayString(unref(formatBytes)(stats.value?.sizeBytes)), 1),
            createBaseVNode("div", {
              class: "mt-1 text-[11px] text-text-5",
              title: String(stats.value?.dir ?? "")
            }, [
              createBaseVNode("span", _hoisted_22, "位于 " + toDisplayString(stats.value?.dir), 1)
            ], 8, _hoisted_21)
          ]),
          createBaseVNode("div", _hoisted_23, [
            _cache[32] || (_cache[32] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "Docker 数据根目录", -1)),
            createBaseVNode("div", _hoisted_24, [
              createBaseVNode("span", {
                class: normalizeClass(["dh-badge", stats.value?.dockerRootVisible ? "dh-badge-run" : "dh-badge-warn"])
              }, toDisplayString(stats.value?.dockerRootVisible ? "可见" : "不可见"), 3)
            ]),
            createBaseVNode("div", _hoisted_25, toDisplayString(stats.value?.dockerRoot || "—"), 1)
          ]),
          createBaseVNode("div", _hoisted_26, [
            _cache[33] || (_cache[33] = createBaseVNode("div", { class: "text-[12px] text-text-4" }, "命名卷目录", -1)),
            createBaseVNode("div", _hoisted_27, [
              createBaseVNode("span", {
                class: normalizeClass(["dh-badge", stats.value?.volumeRootMounted ? "dh-badge-run" : "dh-badge-warn"])
              }, toDisplayString(stats.value?.volumeRootMounted ? "已挂载" : "未挂载"), 3)
            ]),
            createBaseVNode("div", _hoisted_28, toDisplayString(stats.value?.volumeRootMounted ? "可以读取卷数据" : "想备份卷数据需只读挂载 /var/lib/docker/volumes"), 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_29, [
          createBaseVNode("div", _hoisted_30, [
            createVNode(unref(FolderTree), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[34] || (_cache[34] = createBaseVNode("span", null, "宿主路径映射", -1)),
            _cache[35] || (_cache[35] = createBaseVNode("span", { class: "ml-auto text-[11.5px] font-normal text-text-5" }, " 冒号两边不必写一样，右边叫什么都可以 ", -1))
          ]),
          createBaseVNode("div", _hoisted_31, [
            !pathMappings.value.length ? (openBlock(), createElementBlock("div", _hoisted_32, [..._cache[36] || (_cache[36] = [
              createTextVNode(" 没有识别到任何挂载映射 —— 说明 Dockhelm 看不到宿主机的 docker 目录，读不到你的 compose 文件。 在 compose 里挂一行即可，例如 ", -1),
              createBaseVNode("code", { class: "font-mono text-accent" }, "/volume1/docker:/host/docker", -1),
              createTextVNode(" （右边叫什么名字都行，启动时会自动识别）。 ", -1)
            ])])) : (openBlock(), createElementBlock("div", _hoisted_33, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(pathMappings.value, (m) => {
                return openBlock(), createElementBlock("div", {
                  key: m.host + ">" + m.container,
                  class: "flex flex-wrap items-center gap-x-2 gap-y-1"
                }, [
                  createBaseVNode("span", {
                    class: normalizeClass(["dh-badge", m.visible ? "dh-badge-run" : "dh-badge-warn"])
                  }, toDisplayString(m.visible ? "可见" : "不可见"), 3),
                  createBaseVNode("span", _hoisted_34, toDisplayString(m.host), 1),
                  _cache[37] || (_cache[37] = createBaseVNode("span", { class: "text-text-5" }, "→", -1)),
                  createBaseVNode("span", _hoisted_35, toDisplayString(m.container), 1),
                  m.host === m.container ? (openBlock(), createElementBlock("span", _hoisted_36, "两边一致")) : createCommentVNode("", true),
                  m.source === "env" ? (openBlock(), createElementBlock("span", _hoisted_37, "来自 DOCKHELM_HOST_ROOTS")) : createCommentVNode("", true)
                ]);
              }), 128))
            ]))
          ])
        ]),
        tab.value === "snapshots" && policy.value ? (openBlock(), createElementBlock("div", _hoisted_38, [
          createBaseVNode("div", _hoisted_39, [
            createVNode(unref(Trash2), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[39] || (_cache[39] = createBaseVNode("span", null, "保留策略", -1)),
            _cache[40] || (_cache[40] = createBaseVNode("span", { class: "ml-2 text-[11.5px] font-normal text-text-5" }, " 定时清理每天跑一次 · 也可随时点「清理过期」立即执行 ", -1)),
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
              _cache[38] || (_cache[38] = createTextVNode("保存策略 ", -1))
            ], 8, _hoisted_40)
          ]),
          createBaseVNode("div", _hoisted_41, [
            createVNode(_sfc_main$1, {
              title: "每容器保留最近份数",
              sub: "超出后按时间滚动覆盖（更新前快照不计入这个额度）"
            }, {
              default: withCtx(() => [
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => policy.value.backupKeepPerContainer = $event),
                  class: "dh-select !w-[110px] !py-[5px] !text-[11.5px]"
                }, [..._cache[41] || (_cache[41] = [
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
            createVNode(_sfc_main$1, {
              title: "快照保留期",
              sub: "到期自动清理"
            }, {
              default: withCtx(() => [
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[6] || (_cache[6] = ($event) => policy.value.backupMaxAgeDays = $event),
                  class: "dh-select !w-[110px] !py-[5px] !text-[11.5px]"
                }, [..._cache[42] || (_cache[42] = [
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
            createVNode(_sfc_main$1, {
              title: "快照总容量上限",
              sub: "达到上限后先清最旧的快照"
            }, {
              default: withCtx(() => [
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[7] || (_cache[7] = ($event) => policy.value.backupMaxTotalMB = $event),
                  class: "dh-select !w-[110px] !py-[5px] !text-[11.5px]"
                }, [..._cache[43] || (_cache[43] = [
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
            createVNode(_sfc_main$1, {
              title: "更新前快照永不自动清理",
              sub: "这是回滚的底牌 —— 自动更新出事后唯一能救回来的东西"
            }, {
              default: withCtx(() => [
                createVNode(_sfc_main$4, {
                  modelValue: policy.value.backupKeepPreUpdate,
                  "onUpdate:modelValue": _cache[8] || (_cache[8] = ($event) => policy.value.backupKeepPreUpdate = $event),
                  label: "更新前快照永不自动清理"
                }, null, 8, ["modelValue"])
              ]),
              _: 1
            })
          ])
        ])) : createCommentVNode("", true),
        tab.value === "snapshots" ? (openBlock(), createElementBlock("div", _hoisted_42, [
          createBaseVNode("div", _hoisted_43, [
            createVNode(unref(HardDrive), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[44] || (_cache[44] = createBaseVNode("span", null, "存储位置", -1))
          ]),
          createBaseVNode("div", _hoisted_44, [
            createBaseVNode("div", _hoisted_45, [
              createBaseVNode("div", _hoisted_46, toDisplayString(stats.value?.dir || "/data/backups"), 1),
              _cache[45] || (_cache[45] = createBaseVNode("div", { class: "mt-1.5 text-[11.5px] leading-relaxed text-text-5" }, [
                createTextVNode(" 这是 Dockhelm "),
                createBaseVNode("b", { class: "text-text-4" }, "容器内"),
                createTextVNode("的路径。建议把宿主目录映射进来， 这样容器重建、换镜像都还在自己的快照。 ")
              ], -1))
            ]),
            _cache[48] || (_cache[48] = createStaticVNode('<div class="border-t border-line-1"></div><div class="flex items-start gap-2.5"><span class="mt-[3px] h-[15px] w-[15px] flex-none rounded-[5px] bg-warn opacity-60"></span><div class="text-[11.5px] leading-relaxed text-text-4"><b class="text-text-3">bind mount</b> 的数据在宿主机目录上，容器内默认看不见 —— Dockhelm 只能把路径记进快照，没法替你打包那份数据； <b class="text-text-3">named volume</b> 的内容才在 <code class="text-text-3">/var/lib/docker/volumes</code> 下，可以真正读到。 </div></div>', 2)),
            pathMappings.value.length ? (openBlock(), createElementBlock("div", _hoisted_47, [
              _cache[47] || (_cache[47] = createBaseVNode("div", { class: "mb-2 text-[12px] text-text-4" }, "宿主机路径映射（启动时从自身容器自动识别）", -1)),
              createBaseVNode("div", _hoisted_48, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(pathMappings.value, (m, i) => {
                  return openBlock(), createElementBlock("div", {
                    key: i,
                    class: "flex flex-wrap items-center gap-2 text-[11.5px]"
                  }, [
                    createBaseVNode("span", {
                      class: normalizeClass(["dh-badge", m.source === "auto" ? "dh-badge-accent" : "dh-badge-plain"])
                    }, toDisplayString(m.source === "auto" ? "自动识别" : "环境变量"), 3),
                    createBaseVNode("code", _hoisted_49, toDisplayString(m.host), 1),
                    _cache[46] || (_cache[46] = createBaseVNode("span", { class: "text-text-6" }, "→", -1)),
                    createBaseVNode("code", _hoisted_50, toDisplayString(m.container), 1),
                    createBaseVNode("span", {
                      class: normalizeClass(["dh-badge", m.visible ? "dh-badge-run" : "dh-badge-err"])
                    }, toDisplayString(m.visible ? "可见" : "看不见"), 3)
                  ]);
                }), 128))
              ])
            ])) : createCommentVNode("", true)
          ])
        ])) : createCommentVNode("", true),
        tab.value === "snapshots" ? (openBlock(), createElementBlock("div", _hoisted_51, [
          createBaseVNode("div", _hoisted_52, [
            createVNode(unref(Archive), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[49] || (_cache[49] = createBaseVNode("span", null, "配置快照", -1)),
            createBaseVNode("span", _hoisted_53, toDisplayString(backups.value.length) + " 份 · 覆盖 " + toDisplayString(stats.value?.containers ?? 0) + " 个容器 ", 1)
          ]),
          !backups.value.length ? (openBlock(), createBlock(_sfc_main$2, {
            key: 0,
            icon: unref(Archive),
            title: loading.value ? "正在载入…" : "还没有任何快照",
            description: "更新容器时 Dockhelm 会自动写一份快照（用于失败回滚），你也可以在这里手动备份。"
          }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_54, [
            createBaseVNode("table", _hoisted_55, [
              _cache[51] || (_cache[51] = createBaseVNode("thead", null, [
                createBaseVNode("tr", null, [
                  createBaseVNode("th", null, "容器"),
                  createBaseVNode("th", { class: "w-[130px]" }, "快照时间"),
                  createBaseVNode("th", { class: "w-[80px]" }, "来源"),
                  createBaseVNode("th", { class: "w-[110px]" }, "内容"),
                  createBaseVNode("th", { class: "w-[80px]" }, "大小"),
                  createBaseVNode("th", { class: "w-[90px]" }, "快照时状态"),
                  createBaseVNode("th", { class: "w-[130px]" })
                ])
              ], -1)),
              createBaseVNode("tbody", null, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(rows.value, (item) => {
                  return openBlock(), createElementBlock("tr", {
                    key: item.container + "/" + item.ts
                  }, [
                    createBaseVNode("td", null, [
                      createBaseVNode("div", _hoisted_56, toDisplayString(item.container), 1),
                      createBaseVNode("div", {
                        class: "max-w-[220px] truncate font-mono text-[10.5px] text-text-6",
                        title: item.image
                      }, toDisplayString(item.image || "—"), 9, _hoisted_57)
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("div", _hoisted_58, toDisplayString(unref(relativeTime)(item.created)), 1),
                      createBaseVNode("div", _hoisted_59, toDisplayString(item.ts), 1)
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", reasonTone(item.reason)])
                      }, toDisplayString(reasonLabel(item.reason)), 3)
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", item.withData ? "dh-badge-accent" : "dh-badge-plain"])
                      }, toDisplayString(item.withData ? "配置 + 卷数据" : "仅配置"), 3)
                    ]),
                    createBaseVNode("td", _hoisted_60, toDisplayString(unref(formatBytes)(item.size)), 1),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge whitespace-nowrap", item.running ? "dh-badge-run" : "dh-badge-stop"]),
                        title: "拍下这份快照时容器是运行中还是已停止。容器现在是否在运行请看容器页"
                      }, toDisplayString(item.running ? "运行中" : "已停止"), 3)
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("div", _hoisted_61, [
                        createBaseVNode("button", {
                          class: "dh-btn dh-btn-sm",
                          onClick: ($event) => openRestore(item)
                        }, [
                          createVNode(unref(RotateCcw), { class: "h-3 w-3" }),
                          _cache[50] || (_cache[50] = createTextVNode("还原… ", -1))
                        ], 8, _hoisted_62),
                        createBaseVNode("button", {
                          class: "dh-btn dh-btn-sm dh-btn-danger",
                          onClick: ($event) => removeTarget.value = item
                        }, [
                          createVNode(unref(Trash2), { class: "h-3 w-3" })
                        ], 8, _hoisted_63)
                      ])
                    ])
                  ]);
                }), 128))
              ])
            ])
          ]))
        ])) : (openBlock(), createElementBlock("div", _hoisted_64, [
          createBaseVNode("div", _hoisted_65, [
            createVNode(unref(FolderTree), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[52] || (_cache[52] = createBaseVNode("span", null, "compose 项目", -1)),
            _cache[53] || (_cache[53] = createBaseVNode("span", { class: "ml-auto text-[11.5px] font-normal text-text-5" }, " 来自容器标签 com.docker.compose.project.config_files ", -1))
          ]),
          !projects.value.length ? (openBlock(), createBlock(_sfc_main$2, {
            key: 0,
            icon: unref(FolderTree),
            title: loading.value ? "正在载入…" : "没有发现 compose 项目",
            description: "只有由 docker compose 创建的容器才会带上项目标签。"
          }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_66, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(projects.value, (p) => {
              return openBlock(), createElementBlock("div", {
                key: p.project,
                class: "border-b border-line-1 p-3.5 last:border-b-0"
              }, [
                createBaseVNode("div", _hoisted_67, [
                  createBaseVNode("span", _hoisted_68, toDisplayString(p.project), 1),
                  createBaseVNode("span", _hoisted_69, toDisplayString(p.containers?.length ?? 0) + " 个容器", 1),
                  p.workDir ? (openBlock(), createElementBlock("span", {
                    key: 0,
                    class: "truncate font-mono text-[11px] text-text-6",
                    title: p.workDir
                  }, toDisplayString(p.workDir), 9, _hoisted_70)) : createCommentVNode("", true)
                ]),
                createBaseVNode("div", _hoisted_71, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(p.containers, (c) => {
                    return openBlock(), createElementBlock("span", {
                      key: c,
                      class: "rounded-md bg-ink-800 px-1.5 py-[2px] font-mono text-[10.5px] text-text-4"
                    }, toDisplayString(c), 1);
                  }), 128))
                ]),
                p.readable?.length ? (openBlock(), createElementBlock("div", _hoisted_72, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(p.readable, (f) => {
                    return openBlock(), createElementBlock("div", {
                      key: f.path,
                      class: "flex items-center gap-2 rounded-[9px] border border-line-1 bg-ink-800 px-2.5 py-1.5"
                    }, [
                      createVNode(unref(FileCode2), { class: "h-3.5 w-3.5 flex-none text-accent" }),
                      createBaseVNode("div", _hoisted_73, [
                        createBaseVNode("div", _hoisted_74, toDisplayString(f.name), 1),
                        createBaseVNode("div", _hoisted_75, toDisplayString(f.hostPath), 1)
                      ]),
                      createBaseVNode("span", _hoisted_76, toDisplayString(unref(formatBytes)(f.size)), 1),
                      createBaseVNode("button", {
                        class: "dh-btn dh-btn-sm",
                        onClick: ($event) => openFile(f)
                      }, [
                        createVNode(unref(Eye), { class: "h-3 w-3" }),
                        _cache[54] || (_cache[54] = createTextVNode("查看 ", -1))
                      ], 8, _hoisted_77)
                    ]);
                  }), 128))
                ])) : createCommentVNode("", true),
                p.unreadable?.length ? (openBlock(), createElementBlock("div", _hoisted_78, [
                  createBaseVNode("div", _hoisted_79, [
                    createVNode(unref(TriangleAlert), { class: "h-3.5 w-3.5" }),
                    _cache[55] || (_cache[55] = createTextVNode(" 以下文件", -1)),
                    _cache[56] || (_cache[56] = createBaseVNode("b", null, "知道路径但在 Dockhelm 容器里看不见", -1)),
                    _cache[57] || (_cache[57] = createTextVNode("，请自行备份： ", -1))
                  ]),
                  (openBlock(true), createElementBlock(Fragment, null, renderList(p.unreadable, (u) => {
                    return openBlock(), createElementBlock("div", {
                      key: u,
                      class: "truncate font-mono text-[10.5px] text-text-5",
                      title: u
                    }, toDisplayString(u), 9, _hoisted_80);
                  }), 128)),
                  _cache[58] || (_cache[58] = createBaseVNode("div", { class: "text-[11px] leading-relaxed text-text-5" }, [
                    createTextVNode(" 解决办法：在 Dockhelm 的 compose 里把这个目录也挂进来即可 —— "),
                    createBaseVNode("b", { class: "text-text-3" }, "右边叫什么名字都行"),
                    createTextVNode("（例如 "),
                    createBaseVNode("code", null, "- /volume1/docker:/host/docker"),
                    createTextVNode("），启动时会自动识别， 容器标签里的宿主路径就能换算过去。 ")
                  ], -1))
                ])) : createCommentVNode("", true)
              ]);
            }), 128))
          ]))
        ])),
        createVNode(_sfc_main$3, {
          open: showRestore.value,
          title: "还原容器配置",
          subtitle: restoreTarget.value ? `${restoreTarget.value.container} · ${restoreTarget.value.ts}` : "",
          width: "660px",
          busy: restoring.value,
          onClose: _cache[11] || (_cache[11] = ($event) => {
            showRestore.value = false;
            restoreResult.value = null;
          })
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[10] || (_cache[10] = ($event) => {
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
            ], 8, _hoisted_102)) : createCommentVNode("", true)
          ]),
          default: withCtx(() => [
            restoreResult.value ? (openBlock(), createElementBlock("div", _hoisted_81, [
              createBaseVNode("div", {
                class: normalizeClass([
                  "flex items-center gap-2 rounded-[10px] border px-3 py-2.5 text-[12.5px]",
                  restoreResult.value.ok ? "border-line-ok bg-soft-ok text-run-text" : "border-line-err bg-soft-err text-err-text"
                ])
              }, [
                (openBlock(), createBlock(resolveDynamicComponent(restoreResult.value.ok ? unref(CircleCheck) : unref(TriangleAlert)), { class: "h-4 w-4 flex-none" })),
                createTextVNode(" " + toDisplayString(restoreResult.value.message), 1)
              ], 2),
              createBaseVNode("div", _hoisted_82, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(restoreResult.value.steps, (s, i) => {
                  return openBlock(), createElementBlock("div", {
                    key: i,
                    class: "py-[3px] font-mono text-[11.5px] text-text-3"
                  }, toDisplayString(s), 1);
                }), 128))
              ])
            ])) : (openBlock(), createElementBlock("div", _hoisted_83, [
              createBaseVNode("div", _hoisted_84, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none" }),
                _cache[59] || (_cache[59] = createBaseVNode("div", null, [
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
                createBaseVNode("div", _hoisted_85, [
                  _cache[60] || (_cache[60] = createBaseVNode("span", null, "与当前状态的差异", -1)),
                  diffLoading.value ? (openBlock(), createBlock(unref(RefreshCw), {
                    key: 0,
                    class: "h-3 w-3 dh-spin text-text-5"
                  })) : createCommentVNode("", true)
                ]),
                diffLoading.value ? (openBlock(), createElementBlock("div", _hoisted_86, "正在比对…")) : !diff.value.length ? (openBlock(), createElementBlock("div", _hoisted_87, " 快照与当前配置没有差异。 ")) : (openBlock(), createElementBlock("div", _hoisted_88, [
                  createBaseVNode("table", _hoisted_89, [
                    _cache[61] || (_cache[61] = createBaseVNode("thead", null, [
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
                          createBaseVNode("td", _hoisted_90, toDisplayString(d.field), 1),
                          createBaseVNode("td", _hoisted_91, toDisplayString(d.snapshot || "（空）"), 1),
                          createBaseVNode("td", _hoisted_92, toDisplayString(d.current || "（空）"), 1)
                        ]);
                      }), 128))
                    ])
                  ])
                ]))
              ]),
              volMeta.value?.items?.length ? (openBlock(), createElementBlock("div", _hoisted_93, [
                _cache[65] || (_cache[65] = createBaseVNode("div", { class: "mb-2 text-[12px] text-text-3" }, "卷数据", -1)),
                createBaseVNode("label", {
                  class: normalizeClass([
                    "flex cursor-pointer items-start gap-2.5 rounded-[9px] border px-3 py-2.5",
                    packedVolumes.value.length ? "border-line-1 bg-ink-800" : "cursor-not-allowed border-line-1 bg-ink-800 opacity-60"
                  ])
                }, [
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => withVolumes.value = $event),
                    type: "checkbox",
                    class: "mt-[3px] h-[14px] w-[14px] accent-warn",
                    disabled: !packedVolumes.value.length
                  }, null, 8, _hoisted_94), [
                    [vModelCheckbox, withVolumes.value]
                  ]),
                  createBaseVNode("span", _hoisted_95, [
                    createBaseVNode("b", {
                      class: normalizeClass(packedVolumes.value.length ? "text-warn-text" : "text-text-4")
                    }, " 含卷数据（覆盖现有文件） ", 2),
                    createBaseVNode("div", _hoisted_96, [
                      packedVolumes.value.length ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                        createTextVNode(" 这份快照打包了 " + toDisplayString(packedVolumes.value.length) + " 个卷： ", 1),
                        createBaseVNode("span", _hoisted_97, toDisplayString(packedVolumes.value.map((v) => v.name).join("、")), 1),
                        createTextVNode(" （共 " + toDisplayString(unref(formatBytes)(packedVolumes.value.reduce((a, v) => a + v.bytes, 0))) + "）。 勾选后会把卷", 1),
                        _cache[62] || (_cache[62] = createBaseVNode("b", null, "当前的内容整个覆盖掉", -1)),
                        _cache[63] || (_cache[63] = createTextVNode(" —— 这是不可逆的。 ", -1))
                      ], 64)) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
                        createTextVNode(" 这份快照里没有任何被打包的卷数据，只能还原容器配置。 ")
                      ], 64))
                    ])
                  ])
                ], 2),
                unpackedVolumes.value.length ? (openBlock(), createElementBlock("div", _hoisted_98, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(unpackedVolumes.value, (v) => {
                    return openBlock(), createElementBlock("div", {
                      key: v.destination,
                      class: "flex items-start gap-2 text-[11.5px]"
                    }, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge flex-none", v.type === "bind" ? "dh-badge-warn" : "dh-badge-plain"])
                      }, toDisplayString(v.type === "bind" ? "绑定挂载" : v.type), 3),
                      createBaseVNode("code", _hoisted_99, toDisplayString(v.destination), 1),
                      createBaseVNode("span", _hoisted_100, toDisplayString(v.note), 1)
                    ]);
                  }), 128))
                ])) : createCommentVNode("", true),
                withVolumes.value && packedVolumes.value.length ? (openBlock(), createElementBlock("div", _hoisted_101, [
                  createVNode(unref(TriangleAlert), { class: "mt-[1px] h-3.5 w-3.5 flex-none" }),
                  _cache[64] || (_cache[64] = createBaseVNode("span", null, "已勾选「含卷数据」：还原过程中会把这些卷里的现有文件覆盖成快照里的版本，无法撤销。", -1))
                ])) : createCommentVNode("", true)
              ])) : createCommentVNode("", true)
            ]))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: showImport.value,
          title: "导入备份包",
          subtitle: "粘贴一份快照 JSON",
          width: "640px",
          busy: importing.value,
          onClose: _cache[16] || (_cache[16] = ($event) => showImport.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[15] || (_cache[15] = ($event) => showImport.value = false)
            }, "取消"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: importing.value || !importContainer.value.trim() || !importTS.value.trim() || !importText.value.trim(),
              onClick: doImport
            }, [
              importing.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(_component_Download, {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              _cache[72] || (_cache[72] = createTextVNode("导入 ", -1))
            ], 8, _hoisted_106)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_103, [
              createBaseVNode("div", _hoisted_104, [
                createVNode(_component_Download, { class: "mt-[2px] h-3.5 w-3.5 flex-none" }),
                _cache[66] || (_cache[66] = createBaseVNode("span", null, [
                  createTextVNode(" 接受 Dockhelm 写出的快照文件（内容形如 "),
                  createBaseVNode("code", null, '{"_dockhelm":{…},"inspect":{…}}'),
                  createTextVNode("）。 导入只写文件、不动任何容器；导入后就能像本地快照一样在还原弹窗里使用。 同名时间戳不会覆盖已有快照，会自动加后缀。 ")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_105, [
                createBaseVNode("div", null, [
                  _cache[67] || (_cache[67] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "容器名", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[12] || (_cache[12] = ($event) => importContainer.value = $event),
                    class: "dh-input",
                    placeholder: "redis",
                    onChange: guessContainer
                  }, null, 544), [
                    [vModelText, importContainer.value]
                  ]),
                  _cache[68] || (_cache[68] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-5" }, "决定这份快照归到哪个容器的列表下。", -1))
                ]),
                createBaseVNode("div", null, [
                  _cache[69] || (_cache[69] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "快照时间戳", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[13] || (_cache[13] = ($event) => importTS.value = $event),
                    class: "dh-input font-mono",
                    placeholder: "20261009-094200"
                  }, null, 512), [
                    [vModelText, importTS.value]
                  ]),
                  _cache[70] || (_cache[70] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-5" }, "格式 20060102-150405。", -1))
                ])
              ]),
              createBaseVNode("div", null, [
                _cache[71] || (_cache[71] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "快照 JSON", -1)),
                withDirectives(createBaseVNode("textarea", {
                  "onUpdate:modelValue": _cache[14] || (_cache[14] = ($event) => importText.value = $event),
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
          title: "删除快照",
          subtitle: removeTarget.value ? `${removeTarget.value.container} · ${removeTarget.value.ts}` : "",
          width: "400px",
          busy: removing.value,
          onClose: _cache[18] || (_cache[18] = ($event) => removeTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[17] || (_cache[17] = ($event) => removeTarget.value = null)
            }, "取消", 8, _hoisted_107),
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
            ], 8, _hoisted_108)
          ]),
          default: withCtx(() => [
            _cache[73] || (_cache[73] = createBaseVNode("div", { class: "text-[12.5px] text-text-3" }, " 删除后这份配置快照就无法再用于还原。已经运行中的容器不受影响。 ", -1))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: confirmPrune.value,
          title: "清理过期快照",
          width: "460px",
          busy: busy.value === "prune",
          onClose: _cache[20] || (_cache[20] = ($event) => confirmPrune.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: busy.value === "prune",
              onClick: _cache[19] || (_cache[19] = ($event) => confirmPrune.value = false)
            }, "取消", 8, _hoisted_113),
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
              _cache[77] || (_cache[77] = createTextVNode("确认清理 ", -1))
            ], 8, _hoisted_114)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_109, [
              createBaseVNode("div", _hoisted_110, [
                createVNode(unref(TriangleAlert), { class: "mt-[2px] h-4 w-4 flex-none" }),
                _cache[74] || (_cache[74] = createBaseVNode("div", null, [
                  createTextVNode(" 会按下面的保留策略"),
                  createBaseVNode("b", null, "直接删除快照文件"),
                  createTextVNode("，连同打包进去的卷数据一起，"),
                  createBaseVNode("b", null, "不可恢复"),
                  createTextVNode("。 「更新前快照永不自动清理」打开时，那部分不会被碰到。 ")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_111, [
                _cache[75] || (_cache[75] = createTextVNode(" 当前策略：每个容器保留最近 ", -1)),
                createBaseVNode("b", null, toDisplayString(policy.value?.backupKeepPerContainer ? policy.value.backupKeepPerContainer + " 份" : "不限"), 1),
                createTextVNode(" · 快照保留 " + toDisplayString(policy.value?.backupMaxAgeDays ? policy.value.backupMaxAgeDays + " 天" : "不限") + " · 总容量上限 " + toDisplayString(policy.value?.backupMaxTotalMB ? policy.value.backupMaxTotalMB + " MB" : "不限") + " ", 1),
                _cache[76] || (_cache[76] = createBaseVNode("div", { class: "mt-1 text-[11.5px] text-text-5" }, "要改策略请在上面的「保留策略」卡里改并保存。", -1))
              ]),
              createBaseVNode("div", _hoisted_112, "当前共 " + toDisplayString(backups.value.length) + " 份快照，占用 " + toDisplayString(unref(formatBytes)(stats.value?.sizeBytes)) + "。", 1)
            ])
          ]),
          _: 1
        }, 8, ["open", "busy"]),
        createVNode(_sfc_main$3, {
          open: !!showFile.value,
          title: "文件内容",
          subtitle: showFile.value?.path,
          width: "720px",
          onClose: _cache[21] || (_cache[21] = ($event) => showFile.value = null)
        }, {
          default: withCtx(() => [
            createBaseVNode("pre", _hoisted_115, toDisplayString(showFile.value?.content), 1)
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
