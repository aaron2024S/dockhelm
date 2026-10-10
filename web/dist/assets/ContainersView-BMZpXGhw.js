import { x as createLucideIcon, d as defineComponent, u as useAppStore, D as useToastStore, m as ref, y as watch, a as onMounted, o as onUnmounted, c as createElementBlock, e as createBaseVNode, t as toDisplayString, f as createVNode, n as normalizeClass, g as unref, R as RefreshCw, k as createTextVNode, S as Search, w as withDirectives, E as vModelText, j as createCommentVNode, i as createBlock, G as Box, F as Fragment, r as renderList, H as withCtx, I as useRoute, p as computed, B as api, s as openBlock, z as normalizeStyle, J as RouterLink, K as withModifiers, L as reactive, l as vModelCheckbox, C as openStream } from "./index-NwrwvQBT.js";
import { s as shortImage, c as containerStateLabel, b as relativeTime } from "./format-NZu8lNOP.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-DzPLq4FU.js";
import { _ as _sfc_main$2 } from "./Modal.vue_vue_type_script_setup_true_lang-CCEyRSd2.js";
import { _ as _sfc_main$3, S as Square, R as RotateCw } from "./PortChips.vue_vue_type_script_setup_true_lang-BFjUn_sb.js";
import { P as Play } from "./play-BZZncmvg.js";
import { T as Trash2 } from "./trash-2-Dz4Spa1V.js";
const EllipsisVertical = createLucideIcon("EllipsisVerticalIcon", [
  ["circle", { cx: "12", cy: "12", r: "1", key: "41hilf" }],
  ["circle", { cx: "12", cy: "5", r: "1", key: "gxeob9" }],
  ["circle", { cx: "12", cy: "19", r: "1", key: "lyex9k" }]
]);
const _hoisted_1 = { class: "dh-phead" };
const _hoisted_2 = { class: "flex flex-wrap gap-1.5" };
const _hoisted_3 = ["data-on"];
const _hoisted_4 = ["data-on"];
const _hoisted_5 = ["data-on"];
const _hoisted_6 = ["data-on"];
const _hoisted_7 = { class: "ml-auto flex flex-wrap gap-2" };
const _hoisted_8 = ["disabled"];
const _hoisted_9 = ["disabled"];
const _hoisted_10 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_11 = { class: "relative min-w-[180px] flex-1 sm:max-w-[280px]" };
const _hoisted_12 = {
  key: 0,
  class: "ml-auto flex cursor-pointer items-center gap-2 text-[12px] text-text-5"
};
const _hoisted_13 = ["checked"];
const _hoisted_14 = {
  key: 1,
  class: "ml-auto flex flex-wrap items-center gap-2"
};
const _hoisted_15 = { class: "text-[12px] text-text-4" };
const _hoisted_16 = ["disabled"];
const _hoisted_17 = ["disabled"];
const _hoisted_18 = ["disabled"];
const _hoisted_19 = {
  key: 0,
  class: "dh-card grid h-[240px] place-items-center"
};
const _hoisted_20 = {
  key: 2,
  class: "grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"
};
const _hoisted_21 = { class: "flex items-start gap-2.5" };
const _hoisted_22 = { class: "mt-[3px] flex-none cursor-pointer" };
const _hoisted_23 = ["checked", "onChange"];
const _hoisted_24 = { class: "grid h-[34px] w-[34px] flex-none place-items-center rounded-[10px] bg-line-2 text-[13px] font-semibold text-accent-text" };
const _hoisted_25 = { class: "min-w-0 flex-1" };
const _hoisted_26 = { class: "truncate" };
const _hoisted_27 = ["title"];
const _hoisted_28 = { class: "relative flex-none" };
const _hoisted_29 = ["onClick"];
const _hoisted_30 = ["onClick"];
const _hoisted_31 = ["onClick"];
const _hoisted_32 = ["onClick"];
const _hoisted_33 = { class: "flex flex-wrap items-center gap-1.5" };
const _hoisted_34 = {
  key: 0,
  class: "dh-badge dh-badge-accent"
};
const _hoisted_35 = {
  key: 1,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_36 = {
  key: 2,
  class: "dh-badge dh-badge-warn"
};
const _hoisted_37 = ["title"];
const _hoisted_38 = ["title"];
const _hoisted_39 = { class: "flex-none font-medium tabular-nums" };
const _hoisted_40 = { class: "mt-auto flex gap-1.5 border-t border-line-2 pt-2.5" };
const _hoisted_41 = ["disabled", "onClick"];
const _hoisted_42 = ["disabled", "onClick"];
const _hoisted_43 = ["disabled", "onClick"];
const _hoisted_44 = { class: "text-[10.5px] text-text-6" };
const _hoisted_45 = { class: "flex flex-col gap-3" };
const _hoisted_46 = { class: "max-h-[200px] overflow-auto rounded-[10px] border border-line-2 bg-ink-800 p-2.5" };
const _hoisted_47 = ["disabled"];
const _hoisted_48 = ["disabled"];
const _hoisted_49 = { class: "flex flex-col gap-3" };
const _hoisted_50 = { class: "flex cursor-pointer items-start gap-2.5 text-[12.5px] text-text-2" };
const _hoisted_51 = ["disabled"];
const _hoisted_52 = ["disabled"];
const PROGRESS_LINGER_MS = 2500;
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "ContainersView",
  setup(__props) {
    const app = useAppStore();
    const toast = useToastStore();
    const route = useRoute();
    const containers = ref([]);
    const loading = ref(true);
    const keyword = ref(route.query.q ?? "");
    const filter = ref("all");
    const busyName = ref("");
    const menuFor = ref("");
    const removeTarget = ref(null);
    const removeVolumes = ref(false);
    const removing = ref(false);
    const selected = ref(/* @__PURE__ */ new Set());
    const bulkBusy = ref(false);
    const applying = ref(false);
    const updateProgress = reactive({});
    const STEP_MILESTONES = [
      ["开始处理", 5],
      ["拉取", 10],
      ["复用", 55],
      ["拉取完成", 55],
      ["重建", 62],
      ["快照", 64],
      ["已停止", 72],
      ["无需停止", 72],
      ["改名", 76],
      ["新容器已创建", 86],
      ["已启动", 90],
      ["健康检查通过", 96],
      ["清理", 98]
    ];
    function parsePullBytes(s) {
      const m = /([\d.]+)\s*([KMGTP]?i?B)\s*\/\s*([\d.]+)\s*([KMGTP]?i?B)/i.exec(s);
      if (!m) return [0, 0];
      const scale = (u) => {
        switch (u.toUpperCase().charAt(0)) {
          case "T":
            return 1e12;
          case "G":
            return 1e9;
          case "M":
            return 1e6;
          case "K":
            return 1e3;
          default:
            return 1;
        }
      };
      const cur = Number(m[1]) * scale(m[2] ?? "");
      const total = Number(m[3]) * scale(m[4] ?? "");
      return [cur, total];
    }
    function applyUpdateEvent(ev) {
      const name = typeof ev.data?.container === "string" ? ev.data.container : "";
      if (!name) return;
      const p = updateProgress[name] ??= { message: "", percent: 0, done: "" };
      switch (ev.kind) {
        case "step": {
          const msg = String(ev.data?.message ?? "");
          if (msg) p.message = msg;
          let hit = 0;
          for (const [kw, pct] of STEP_MILESTONES) {
            if (msg.includes(kw)) hit = Math.max(hit, pct);
          }
          p.percent = hit > 0 ? Math.max(p.percent, hit) : Math.min(96, p.percent + 2);
          break;
        }
        case "pull_progress": {
          const [cur, total] = parsePullBytes(String(ev.data?.progress ?? ""));
          if (total > 0) {
            p.percent = Math.max(p.percent, Math.min(55, 10 + Math.round(cur / total * 45)));
          }
          const st = String(ev.data?.status ?? "");
          if (st) p.message = `拉取镜像：${st}`;
          break;
        }
        case "container_status": {
          const st = String(ev.data?.status ?? "");
          const msg = String(ev.data?.message ?? "");
          if (msg) p.message = msg;
          if (st === "pulling") {
            p.percent = Math.max(p.percent, 10);
          } else if (st === "updated") {
            p.done = "updated";
            p.percent = 100;
            scheduleClear(name);
          } else if (st === "failed") {
            p.done = "failed";
            p.percent = Math.max(p.percent, 1);
            scheduleClear(name);
          } else if (st === "up_to_date" || st === "skipped") {
            p.done = "skipped";
            p.percent = 100;
            scheduleClear(name);
          }
          break;
        }
        case "batch_done": {
          for (const k of Object.keys(updateProgress)) delete updateProgress[k];
          break;
        }
      }
    }
    function scheduleClear(name) {
      const snapshot = updateProgress[name];
      window.setTimeout(() => {
        if (updateProgress[name] === snapshot) delete updateProgress[name];
      }, PROGRESS_LINGER_MS);
    }
    function isUpdating(name) {
      return updateProgress[name] !== void 0;
    }
    function pctOf(name) {
      return updateProgress[name]?.percent ?? 0;
    }
    function msgOf(name) {
      return updateProgress[name]?.message ?? "";
    }
    function doneOf(name) {
      return updateProgress[name]?.done ?? "";
    }
    function stateText(name) {
      switch (doneOf(name)) {
        case "updated":
          return "更新完成";
        case "failed":
          return "更新失败";
        case "skipped":
          return "已跳过";
        default:
          return "更新中";
      }
    }
    function barFillClass(name) {
      switch (doneOf(name)) {
        case "updated":
          return "bg-run/15";
        case "failed":
          return "bg-err/15";
        default:
          return "bg-accent/12";
      }
    }
    function barEdgeClass(name) {
      switch (doneOf(name)) {
        case "updated":
          return "bg-gradient-to-l from-run/40 to-transparent";
        case "failed":
          return "bg-gradient-to-l from-err/40 to-transparent";
        default:
          return "bg-gradient-to-l from-accent/40 to-transparent";
      }
    }
    function progressTextClass(name) {
      switch (doneOf(name)) {
        case "updated":
          return "text-run-text";
        case "failed":
          return "text-err-text";
        case "skipped":
          return "text-text-3";
        default:
          return "text-accent-text";
      }
    }
    function progressBorderClass(name) {
      switch (doneOf(name)) {
        case "updated":
          return "border-line-ok bg-run/8";
        case "failed":
          return "border-line-err bg-err/8";
        default:
          return "border-line-accent-soft bg-accent/6";
      }
    }
    const filtered = computed(() => {
      let list = containers.value;
      const k = keyword.value.trim().toLowerCase();
      if (k) {
        list = list.filter(
          (c) => c.name.toLowerCase().includes(k) || c.image.toLowerCase().includes(k) || c.project.toLowerCase().includes(k)
        );
      }
      if (filter.value === "running") list = list.filter((c) => c.state === "running");
      if (filter.value === "stopped") list = list.filter((c) => c.state !== "running");
      if (filter.value === "update") list = list.filter((c) => c.hasUpdate);
      return list;
    });
    const counts = computed(() => ({
      all: containers.value.length,
      running: containers.value.filter((c) => c.state === "running").length,
      stopped: containers.value.filter((c) => c.state !== "running").length,
      update: containers.value.filter((c) => c.hasUpdate).length
    }));
    const updatableSelected = computed(
      () => containers.value.filter((c) => selected.value.has(c.name) && !c.self && !c.excluded).map((c) => c.name)
    );
    const updatableAll = computed(() => containers.value.filter((c) => c.hasUpdate && !c.self && !c.excluded));
    const allSelected = computed(
      () => filtered.value.length > 0 && filtered.value.every((c) => selected.value.has(c.name))
    );
    async function load(silent = false) {
      if (!silent) loading.value = true;
      try {
        const res = await api.get("/api/containers");
        containers.value = res.containers ?? [];
      } catch (e) {
        if (!silent) toast.error("读取容器失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function act(c, action) {
      busyName.value = c.name;
      menuFor.value = "";
      try {
        await api.post(`/api/containers/${encodeURIComponent(c.name)}/action`, { action });
        toast.success(`${c.name} 已${labelOf(action)}`);
        await load(true);
      } catch (e) {
        toast.error(`操作失败`, e instanceof Error ? e.message : String(e));
      } finally {
        busyName.value = "";
      }
    }
    function labelOf(action) {
      return { start: "启动", stop: "停止", restart: "重启", pause: "暂停", unpause: "恢复" }[action] ?? action;
    }
    async function updateOne(c) {
      busyName.value = c.name;
      menuFor.value = "";
      try {
        await api.post("/api/updates/apply", { names: [c.name] });
        toast.info(`已提交更新：${c.name}`, "如果镜像没有变化，容器不会被停止");
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
      } finally {
        busyName.value = "";
      }
    }
    async function updateSelected() {
      const names = updatableSelected.value;
      if (!names.length || applying.value) return;
      applying.value = true;
      try {
        await api.post("/api/updates/apply", { names });
        toast.info(`已提交 ${names.length} 个容器的更新`, "镜像未变化的容器会被自动跳过");
        selected.value = /* @__PURE__ */ new Set();
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
      } finally {
        applying.value = false;
      }
    }
    async function updateAll() {
      const names = updatableAll.value.map((c) => c.name);
      if (!names.length || applying.value) return;
      applying.value = true;
      try {
        await api.post("/api/updates/apply", { names });
        toast.info(`已提交 ${names.length} 个容器的更新`, "镜像未变化的容器会被自动跳过");
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
      } finally {
        applying.value = false;
      }
    }
    async function bulkAct(action) {
      const names = containers.value.filter((c) => selected.value.has(c.name) && !c.self).map((c) => c.name);
      if (!names.length) return;
      bulkPending.value = { action, names };
    }
    const bulkPending = ref(null);
    async function confirmBulkAct() {
      const p = bulkPending.value;
      if (!p) return;
      bulkPending.value = null;
      bulkBusy.value = true;
      let ok = 0;
      for (const name of p.names) {
        try {
          await api.post(`/api/containers/${encodeURIComponent(name)}/action`, { action: p.action });
          ok += 1;
        } catch {
        }
      }
      bulkBusy.value = false;
      selected.value = /* @__PURE__ */ new Set();
      toast.success(
        `已${labelOf(p.action)} ${ok}/${p.names.length} 个`,
        ok < p.names.length ? "部分容器操作失败，详见日志" : void 0
      );
      await load(true);
    }
    async function confirmRemove() {
      const c = removeTarget.value;
      if (!c) return;
      removing.value = true;
      try {
        const res = await api.del(
          `/api/containers/${encodeURIComponent(c.name)}`,
          { force: true, volumes: removeVolumes.value }
        );
        toast.success(res.message || "已删除");
        removeTarget.value = null;
        removeVolumes.value = false;
        await load(true);
      } catch (e) {
        toast.error("删除失败", e instanceof Error ? e.message : String(e));
      } finally {
        removing.value = false;
      }
    }
    function toggleSelect(name) {
      const next = new Set(selected.value);
      if (next.has(name)) next.delete(name);
      else next.add(name);
      selected.value = next;
    }
    function toggleAll() {
      if (allSelected.value) {
        const next = new Set(selected.value);
        for (const c of filtered.value) next.delete(c.name);
        selected.value = next;
      } else {
        const next = new Set(selected.value);
        for (const c of filtered.value) next.add(c.name);
        selected.value = next;
      }
    }
    let inner = null;
    watch(
      () => route.query.q,
      (q) => {
        keyword.value = q ?? "";
      }
    );
    watch(
      () => app.checkTick,
      () => void load(true)
    );
    onMounted(() => {
      void load();
      inner = openStream("/api/events/stream", (topic, ev) => {
        if (topic === "update") {
          applyUpdateEvent(ev);
          if (ev.kind === "batch_done" || ev.kind === "container_status") {
            void load(true);
          }
        }
        if (topic === "container" && ev.kind === "action") {
          void load(true);
        }
      });
    });
    onUnmounted(() => inner?.());
    function closeMenu() {
      menuFor.value = "";
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: "flex flex-col gap-3.5 p-[18px]",
        onClick: closeMenu
      }, [
        createBaseVNode("div", _hoisted_1, [
          _cache[15] || (_cache[15] = createBaseVNode("div", { class: "dh-h1" }, "容器", -1)),
          createBaseVNode("div", _hoisted_2, [
            createBaseVNode("button", {
              class: "dh-chip",
              "data-on": filter.value === "all",
              onClick: _cache[0] || (_cache[0] = ($event) => filter.value = "all")
            }, " 全部 " + toDisplayString(counts.value.all), 9, _hoisted_3),
            createBaseVNode("button", {
              class: "dh-chip",
              "data-on": filter.value === "running",
              onClick: _cache[1] || (_cache[1] = ($event) => filter.value = "running")
            }, " 运行中 " + toDisplayString(counts.value.running), 9, _hoisted_4),
            createBaseVNode("button", {
              class: "dh-chip",
              "data-on": filter.value === "stopped",
              onClick: _cache[2] || (_cache[2] = ($event) => filter.value = "stopped")
            }, " 已停止 " + toDisplayString(counts.value.stopped), 9, _hoisted_5),
            createBaseVNode("button", {
              class: "dh-chip",
              "data-on": filter.value === "update",
              onClick: _cache[3] || (_cache[3] = ($event) => filter.value = "update")
            }, " 有更新 " + toDisplayString(counts.value.update), 9, _hoisted_6)
          ]),
          createBaseVNode("div", _hoisted_7, [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: loading.value,
              onClick: _cache[4] || (_cache[4] = ($event) => load())
            }, [
              createVNode(unref(RefreshCw), {
                class: normalizeClass(["h-3.5 w-3.5", loading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[14] || (_cache[14] = createTextVNode("刷新 ", -1))
            ], 8, _hoisted_8),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: !updatableAll.value.length || applying.value,
              onClick: updateAll
            }, toDisplayString(applying.value ? "提交中…" : `更新 ${updatableAll.value.length} 个容器`), 9, _hoisted_9)
          ])
        ]),
        createBaseVNode("div", _hoisted_10, [
          createBaseVNode("div", _hoisted_11, [
            createVNode(unref(Search), { class: "pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" }),
            withDirectives(createBaseVNode("input", {
              "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => keyword.value = $event),
              class: "dh-input !pl-8",
              placeholder: "按名称 / 镜像 / 项目筛选"
            }, null, 512), [
              [vModelText, keyword.value]
            ])
          ]),
          filtered.value.length && !selected.value.size ? (openBlock(), createElementBlock("label", _hoisted_12, [
            createBaseVNode("input", {
              type: "checkbox",
              class: "h-[14px] w-[14px] accent-accent",
              checked: allSelected.value,
              onChange: toggleAll
            }, null, 40, _hoisted_13),
            createTextVNode(" 全选当前列表（" + toDisplayString(filtered.value.length) + " 个） ", 1)
          ])) : createCommentVNode("", true),
          selected.value.size ? (openBlock(), createElementBlock("div", _hoisted_14, [
            createBaseVNode("span", _hoisted_15, "已选 " + toDisplayString(selected.value.size) + " 个", 1),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              onClick: toggleAll
            }, toDisplayString(allSelected.value ? "取消全选" : "全选本页"), 1),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              disabled: bulkBusy.value,
              onClick: _cache[6] || (_cache[6] = ($event) => bulkAct("restart"))
            }, "重启", 8, _hoisted_16),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              disabled: bulkBusy.value,
              onClick: _cache[7] || (_cache[7] = ($event) => bulkAct("stop"))
            }, "停止", 8, _hoisted_17),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm dh-btn-primary",
              disabled: !updatableSelected.value.length || applying.value,
              onClick: updateSelected
            }, " 更新选中 ", 8, _hoisted_18)
          ])) : createCommentVNode("", true)
        ]),
        loading.value && !containers.value.length ? (openBlock(), createElementBlock("div", _hoisted_19, [
          createVNode(unref(RefreshCw), { class: "h-5 w-5 dh-spin text-text-5" })
        ])) : !filtered.value.length ? (openBlock(), createBlock(_sfc_main$1, {
          key: 1,
          icon: unref(Box),
          title: "没有匹配的容器",
          description: "换个筛选条件，或者确认 Docker 守护进程是否正常。"
        }, null, 8, ["icon"])) : (openBlock(), createElementBlock("div", _hoisted_20, [
          (openBlock(true), createElementBlock(Fragment, null, renderList(filtered.value, (c) => {
            return openBlock(), createElementBlock("div", {
              key: c.id,
              class: normalizeClass(["relative isolate flex flex-col gap-2.5 overflow-hidden rounded-[14px] border bg-ink-700 p-3 transition-colors", [
                c.hasUpdate ? "border-line-warn" : "border-line-1 hover:border-line-4",
                selected.value.has(c.name) ? "!border-accent-line bg-accent-soft" : ""
              ]])
            }, [
              isUpdating(c.name) ? (openBlock(), createElementBlock("div", {
                key: 0,
                class: normalizeClass(["pointer-events-none absolute inset-y-0 left-0 -z-10 transition-[width] duration-500 ease-out", barFillClass(c.name)]),
                style: normalizeStyle({ width: pctOf(c.name) + "%" })
              }, [
                createBaseVNode("div", {
                  class: normalizeClass(["absolute inset-y-0 right-0 w-3", barEdgeClass(c.name)])
                }, null, 2)
              ], 6)) : createCommentVNode("", true),
              createBaseVNode("div", _hoisted_21, [
                createBaseVNode("label", _hoisted_22, [
                  createBaseVNode("input", {
                    type: "checkbox",
                    class: "h-[14px] w-[14px] accent-accent",
                    checked: selected.value.has(c.name),
                    onChange: ($event) => toggleSelect(c.name)
                  }, null, 40, _hoisted_23)
                ]),
                createBaseVNode("div", _hoisted_24, toDisplayString(c.name.slice(0, 2).toUpperCase()), 1),
                createBaseVNode("div", _hoisted_25, [
                  createVNode(unref(RouterLink), {
                    to: `/containers/${encodeURIComponent(c.name)}`,
                    class: "dh-tap-txt block min-w-0 text-[13px] font-semibold hover:text-accent",
                    title: c.name
                  }, {
                    default: withCtx(() => [
                      createBaseVNode("span", _hoisted_26, toDisplayString(c.name), 1)
                    ]),
                    _: 2
                  }, 1032, ["to", "title"]),
                  createBaseVNode("div", {
                    class: "truncate font-mono text-[11.5px] text-text-5",
                    title: c.image
                  }, toDisplayString(unref(shortImage)(c.image)), 9, _hoisted_27)
                ]),
                createBaseVNode("div", _hoisted_28, [
                  createBaseVNode("button", {
                    class: "dh-tap grid h-6 w-6 place-items-center rounded-md text-text-5 hover:bg-ink-650 hover:text-text-1",
                    onClick: withModifiers(($event) => menuFor.value = menuFor.value === c.name ? "" : c.name, ["stop"])
                  }, [
                    createVNode(unref(EllipsisVertical), { class: "h-3.5 w-3.5" })
                  ], 8, _hoisted_29),
                  menuFor.value === c.name ? (openBlock(), createElementBlock("div", {
                    key: 0,
                    class: "absolute right-0 top-7 z-20 w-[150px] overflow-hidden rounded-[10px] border border-line-3 bg-ink-750 py-1 shadow-[var(--shadow-pop)]",
                    onClick: _cache[8] || (_cache[8] = withModifiers(() => {
                    }, ["stop"]))
                  }, [
                    createBaseVNode("button", {
                      class: "block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650",
                      onClick: ($event) => act(c, "restart")
                    }, " 重启 ", 8, _hoisted_30),
                    createBaseVNode("button", {
                      class: "block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650",
                      onClick: ($event) => updateOne(c)
                    }, " 检查并更新 ", 8, _hoisted_31),
                    createBaseVNode("button", {
                      class: "block w-full px-3 py-1.5 text-left text-[12px] text-err-text hover:bg-ink-650",
                      onClick: ($event) => {
                        removeTarget.value = c;
                        menuFor.value = "";
                      }
                    }, " 删除容器… ", 8, _hoisted_32)
                  ])) : createCommentVNode("", true)
                ])
              ]),
              createBaseVNode("div", _hoisted_33, [
                c.self ? (openBlock(), createElementBlock("span", _hoisted_34, "Dockhelm 自身")) : c.excluded ? (openBlock(), createElementBlock("span", _hoisted_35, "已排除")) : createCommentVNode("", true),
                createBaseVNode("span", {
                  class: normalizeClass(["dh-badge", c.state === "running" ? "dh-badge-run" : "dh-badge-stop"])
                }, [
                  createBaseVNode("span", {
                    class: normalizeClass(["h-[7px] w-[7px] rounded-full", c.state === "running" ? "bg-run" : "bg-stop"])
                  }, null, 2),
                  createTextVNode(" " + toDisplayString(unref(containerStateLabel)(c.state)), 1)
                ], 2),
                c.hasUpdate ? (openBlock(), createElementBlock("span", _hoisted_36, "有新版本")) : createCommentVNode("", true),
                c.project ? (openBlock(), createElementBlock("span", {
                  key: 3,
                  class: "dh-badge dh-badge-plain",
                  title: `compose 项目 ${c.project}`
                }, toDisplayString(c.project), 9, _hoisted_37)) : createCommentVNode("", true)
              ]),
              createVNode(_sfc_main$3, {
                ports: c.portList ?? [],
                max: 3
              }, null, 8, ["ports"]),
              isUpdating(c.name) ? (openBlock(), createElementBlock("div", {
                key: 1,
                class: normalizeClass(["flex items-center gap-1.5 rounded-md border px-2 py-1.5 text-[11.5px]", [progressBorderClass(c.name), progressTextClass(c.name)]])
              }, [
                createVNode(unref(RefreshCw), {
                  class: normalizeClass(["h-3 w-3 flex-none", doneOf(c.name) === "" ? "dh-spin" : ""])
                }, null, 8, ["class"]),
                createBaseVNode("span", {
                  class: "min-w-0 flex-1 truncate",
                  title: msgOf(c.name)
                }, [
                  createTextVNode(toDisplayString(stateText(c.name)), 1),
                  msgOf(c.name) ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                    createTextVNode(" · " + toDisplayString(msgOf(c.name)), 1)
                  ], 64)) : createCommentVNode("", true)
                ], 8, _hoisted_38),
                createBaseVNode("span", _hoisted_39, toDisplayString(pctOf(c.name)) + "%", 1)
              ], 2)) : createCommentVNode("", true),
              createBaseVNode("div", _hoisted_40, [
                c.state === "running" ? (openBlock(), createElementBlock("button", {
                  key: 0,
                  class: "dh-btn dh-btn-sm flex-1",
                  disabled: busyName.value === c.name || isUpdating(c.name),
                  onClick: ($event) => act(c, "stop")
                }, [
                  createVNode(unref(Square), { class: "h-3 w-3" }),
                  _cache[16] || (_cache[16] = createTextVNode("停止 ", -1))
                ], 8, _hoisted_41)) : (openBlock(), createElementBlock("button", {
                  key: 1,
                  class: "dh-btn dh-btn-sm flex-1",
                  disabled: busyName.value === c.name || isUpdating(c.name),
                  onClick: ($event) => act(c, "start")
                }, [
                  createVNode(unref(Play), { class: "h-3 w-3" }),
                  _cache[17] || (_cache[17] = createTextVNode("启动 ", -1))
                ], 8, _hoisted_42)),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm flex-1",
                  disabled: busyName.value === c.name || c.self || isUpdating(c.name),
                  onClick: ($event) => act(c, "restart")
                }, [
                  createVNode(unref(RotateCw), {
                    class: normalizeClass(["h-3 w-3", busyName.value === c.name ? "dh-spin" : ""])
                  }, null, 8, ["class"]),
                  _cache[18] || (_cache[18] = createTextVNode("重启 ", -1))
                ], 8, _hoisted_43)
              ]),
              createBaseVNode("div", _hoisted_44, "创建于 " + toDisplayString(unref(relativeTime)(c.created)), 1)
            ], 2);
          }), 128))
        ])),
        createVNode(_sfc_main$2, {
          open: !!bulkPending.value,
          title: `批量${bulkPending.value ? labelOf(bulkPending.value.action) : ""}容器`,
          subtitle: bulkPending.value ? `即将对 ${bulkPending.value.names.length} 个容器执行${labelOf(bulkPending.value.action)}` : "",
          width: "470px",
          busy: bulkBusy.value,
          onClose: _cache[10] || (_cache[10] = ($event) => bulkPending.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: bulkBusy.value,
              onClick: _cache[9] || (_cache[9] = ($event) => bulkPending.value = null)
            }, "取消", 8, _hoisted_47),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: bulkBusy.value,
              onClick: confirmBulkAct
            }, toDisplayString(bulkBusy.value ? "执行中…" : `确认${bulkPending.value ? labelOf(bulkPending.value.action) : ""}`), 9, _hoisted_48)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_45, [
              _cache[19] || (_cache[19] = createBaseVNode("div", { class: "dh-banner dh-banner-warn" }, [
                createBaseVNode("span", { class: "min-w-0 flex-1" }, [
                  createTextVNode(" 批量操作会逐个执行，"),
                  createBaseVNode("b", null, "单个容器失败不会中断其余容器"),
                  createTextVNode("，结束后可在运行记录里查看结果。 ")
                ])
              ], -1)),
              createBaseVNode("div", _hoisted_46, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(bulkPending.value?.names ?? [], (n) => {
                  return openBlock(), createElementBlock("div", {
                    key: n,
                    class: "font-mono text-[11.5px] leading-relaxed text-text-3"
                  }, toDisplayString(n), 1);
                }), 128))
              ])
            ])
          ]),
          _: 1
        }, 8, ["open", "title", "subtitle", "busy"]),
        createVNode(_sfc_main$2, {
          open: !!removeTarget.value,
          title: "删除容器",
          subtitle: removeTarget.value ? `即将删除 ${removeTarget.value.name}` : "",
          width: "470px",
          busy: removing.value,
          onClose: _cache[13] || (_cache[13] = ($event) => {
            removeTarget.value = null;
            removeVolumes.value = false;
          })
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[12] || (_cache[12] = ($event) => {
                removeTarget.value = null;
                removeVolumes.value = false;
              })
            }, " 取消 ", 8, _hoisted_51),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemove
            }, [
              createVNode(unref(Trash2), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(removing.value ? "删除中…" : "确认删除"), 1)
            ], 8, _hoisted_52)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_49, [
              _cache[21] || (_cache[21] = createBaseVNode("div", { class: "rounded-[10px] border border-line-err bg-soft-err px-3 py-2.5 text-[12px] leading-relaxed text-err-text" }, " 此操作不可撤销。删除容器不会删除它的镜像，但容器自身的可写层数据会一并消失。 ", -1)),
              createBaseVNode("label", _hoisted_50, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[11] || (_cache[11] = ($event) => removeVolumes.value = $event),
                  type: "checkbox",
                  class: "mt-[3px] h-[14px] w-[14px] accent-err"
                }, null, 512), [
                  [vModelCheckbox, removeVolumes.value]
                ]),
                _cache[20] || (_cache[20] = createBaseVNode("span", null, [
                  createTextVNode(" 同时删除该容器的"),
                  createBaseVNode("b", { class: "text-err-text" }, "匿名卷"),
                  createBaseVNode("br"),
                  createBaseVNode("span", { class: "text-[11.5px] text-text-5" }, " 勾选后 docker 会连匿名卷一起删掉，卷里的数据将无法找回。命名卷不会被删除。 ")
                ], -1))
              ])
            ])
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
