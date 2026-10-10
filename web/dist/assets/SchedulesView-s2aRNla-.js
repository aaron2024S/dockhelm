import { d as defineComponent, D as useToastStore, a as onMounted, o as onUnmounted, c as createElementBlock, e as createBaseVNode, t as toDisplayString, f as createVNode, n as normalizeClass, g as unref, R as RefreshCw, k as createTextVNode, i as createBlock, F as Fragment, r as renderList, j as createCommentVNode, w as withDirectives, E as vModelText, M as vModelSelect, H as withCtx, m as ref, p as computed, B as api, s as openBlock, C as openStream } from "./index-CFwukrz3.js";
import { e as explainCron, f as formatDayTime } from "./format-NZu8lNOP.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-D-X0dDSu.js";
import { _ as _sfc_main$3 } from "./Modal.vue_vue_type_script_setup_true_lang-gJUwROSk.js";
import { _ as _sfc_main$2 } from "./UiSwitch.vue_vue_type_script_setup_true_lang-Dh4LtQcL.js";
import { P as Plus } from "./plus-ww1wQ7a-.js";
import { L as LoaderCircle } from "./loader-circle-BEPIJzEK.js";
import { P as Play } from "./play-C6_Iwr67.js";
import { T as Trash2 } from "./trash-2-DjE-zQXz.js";
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "dh-sub" };
const _hoisted_4 = { class: "ml-auto flex gap-2" };
const _hoisted_5 = ["disabled"];
const _hoisted_6 = { class: "dh-card" };
const _hoisted_7 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_8 = { class: "dh-table" };
const _hoisted_9 = { class: "text-[13px] font-semibold" };
const _hoisted_10 = ["title"];
const _hoisted_11 = { class: "text-[11.5px] text-text-5" };
const _hoisted_12 = { class: "text-[11.5px] text-text-5" };
const _hoisted_13 = { class: "font-mono text-[11px]" };
const _hoisted_14 = { class: "text-[11.5px] text-text-5" };
const _hoisted_15 = { class: "text-[11.5px] text-text-5" };
const _hoisted_16 = { class: "text-right" };
const _hoisted_17 = { class: "flex justify-end gap-1.5" };
const _hoisted_18 = ["disabled", "onClick"];
const _hoisted_19 = ["onClick"];
const _hoisted_20 = ["onClick"];
const _hoisted_21 = { class: "grid grid-cols-1 gap-3.5 xl:grid-cols-2" };
const _hoisted_22 = {
  id: "schedule-form",
  class: "dh-card"
};
const _hoisted_23 = { class: "dh-card-head" };
const _hoisted_24 = { class: "dh-card-body flex flex-col gap-3.5" };
const _hoisted_25 = { key: 0 };
const _hoisted_26 = { class: "dh-scroll flex max-h-[132px] flex-wrap gap-1.5 overflow-auto rounded-[10px] border border-line-1 bg-ink-800 p-2" };
const _hoisted_27 = ["data-on", "disabled", "title", "onClick"];
const _hoisted_28 = {
  key: 0,
  class: "text-[10px] text-text-5"
};
const _hoisted_29 = {
  key: 1,
  class: "text-[10px] text-text-5"
};
const _hoisted_30 = {
  key: 0,
  class: "p-1 text-[11.5px] text-text-5"
};
const _hoisted_31 = { class: "flex flex-wrap gap-1.5" };
const _hoisted_32 = ["data-on", "title", "onClick"];
const _hoisted_33 = { class: "grid grid-cols-[1.4fr_1fr] gap-2.5" };
const _hoisted_34 = ["value"];
const _hoisted_35 = { class: "dh-label" };
const _hoisted_36 = { class: "text-[11.5px] text-text-5" };
const _hoisted_37 = { class: "font-mono text-text-3" };
const _hoisted_38 = { key: 0 };
const _hoisted_39 = { class: "flex cursor-pointer items-center gap-2.5 text-[12.5px] text-text-2" };
const _hoisted_40 = {
  key: 1,
  class: "dh-banner dh-banner-err py-2 text-[12px]"
};
const _hoisted_41 = ["disabled"];
const _hoisted_42 = { class: "dh-card" };
const _hoisted_43 = {
  key: 0,
  class: "dh-card-body text-[12px] text-text-5"
};
const _hoisted_44 = {
  key: 1,
  class: "dh-card-body flex flex-col gap-2.5"
};
const _hoisted_45 = { class: "w-[74px] flex-none font-mono text-[11.5px] text-accent-text" };
const _hoisted_46 = { class: "min-w-0 flex-1 truncate text-[12.5px]" };
const _hoisted_47 = ["disabled"];
const _hoisted_48 = ["disabled"];
const _hoisted_49 = {
  key: 0,
  class: "grid h-24 place-items-center text-text-5"
};
const _hoisted_50 = {
  key: 1,
  class: "py-8 text-center text-[12.5px] text-text-5"
};
const _hoisted_51 = {
  key: 2,
  class: "overflow-x-auto"
};
const _hoisted_52 = { class: "dh-table" };
const _hoisted_53 = { class: "whitespace-nowrap text-[11.5px] text-text-5" };
const _hoisted_54 = { class: "max-w-[180px] truncate text-[12px] text-text-2" };
const _hoisted_55 = { class: "text-[12px] text-text-3" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "SchedulesView",
  setup(__props) {
    const toast = useToastStore();
    const schedules = ref([]);
    const actions = ref([]);
    const containers = ref([]);
    const excluded = ref([]);
    const loading = ref(true);
    const running = ref(null);
    const removeTarget = ref(null);
    const saveError = ref("");
    const saving = ref(false);
    const removing = ref(false);
    let closeStream = null;
    const editingId = ref(0);
    const form = ref({
      name: "",
      action: "restart",
      targets: [],
      enabled: true,
      /** 重复方式 + 时间的组合，最终换算成 cron；repeat='custom' 时直接用 customCron。 */
      repeat: "daily",
      time: "03:00",
      customCron: "0 3 * * *"
    });
    function blankForm() {
      return {
        name: "",
        action: "restart",
        targets: [],
        enabled: true,
        repeat: "daily",
        time: "03:00",
        customCron: "0 3 * * *"
      };
    }
    const repeatOptions = [
      { key: "daily", label: "每天" },
      { key: "hourly", label: "每小时" },
      { key: "weekdays", label: "工作日（周一至周五）" },
      { key: "weekly-1", label: "每周一" },
      { key: "weekly-2", label: "每周二" },
      { key: "weekly-3", label: "每周三" },
      { key: "weekly-4", label: "每周四" },
      { key: "weekly-5", label: "每周五" },
      { key: "weekly-6", label: "每周六" },
      { key: "weekly-0", label: "每周日" },
      { key: "monthly", label: "每月 1 日" },
      { key: "custom", label: "自定义 cron" }
    ];
    const SHORT_ACTION = {
      start: "启动",
      stop: "停止",
      restart: "重启",
      update: "更新镜像",
      backup: "备份快照",
      prune_images: "清理旧镜像"
    };
    const ACTION_TONE = {
      start: "dh-badge-run",
      stop: "dh-badge-stop",
      restart: "dh-badge-accent",
      update: "dh-badge-warn",
      backup: "dh-badge-plain",
      prune_images: "dh-badge-plain"
    };
    function shortAction(key) {
      return SHORT_ACTION[key] ?? actions.value.find((a) => a.key === key)?.label ?? key;
    }
    function actionTone(key) {
      return ACTION_TONE[key] ?? "dh-badge-plain";
    }
    const composedCron = computed(() => {
      if (form.value.repeat === "custom") return form.value.customCron.trim();
      const [h = "0", m = "0"] = form.value.time.split(":");
      const hh = String(Number(h));
      const mm = String(Number(m));
      switch (form.value.repeat) {
        case "daily":
          return `${mm} ${hh} * * *`;
        case "hourly":
          return `${mm} * * * *`;
        case "weekdays":
          return `${mm} ${hh} * * 1-5`;
        case "monthly":
          return `${mm} ${hh} 1 * *`;
        default: {
          const dow = form.value.repeat.split("-")[1] ?? "0";
          return `${mm} ${hh} * * ${dow}`;
        }
      }
    });
    function applyCron(cron) {
      const t = cron.trim();
      const hhmm = (h, m2) => `${String(h).padStart(2, "0")}:${String(m2).padStart(2, "0")}`;
      let m = /^(\d+)\s+(\d+)\s+\*\s+\*\s+\*$/.exec(t);
      if (m) return { repeat: "daily", time: hhmm(m[2], m[1]) };
      m = /^(\d+)\s+(\d+)\s+\*\s+\*\s+1-5$/.exec(t);
      if (m) return { repeat: "weekdays", time: hhmm(m[2], m[1]) };
      m = /^(\d+)\s+(\d+)\s+\*\s+\*\s+([0-7])$/.exec(t);
      if (m) return { repeat: `weekly-${m[3] === "7" ? "0" : m[3]}`, time: hhmm(m[2], m[1]) };
      m = /^(\d+)\s+(\d+)\s+1\s+\*\s+\*$/.exec(t);
      if (m) return { repeat: "monthly", time: hhmm(m[2], m[1]) };
      m = /^(\d+)\s+\*\s+\*\s+\*\s+\*$/.exec(t);
      if (m) return { repeat: "hourly", time: hhmm("0", m[1]) };
      return { repeat: "custom", time: "03:00" };
    }
    const actionMeta = computed(() => actions.value.find((a) => a.key === form.value.action));
    const needsTargets = computed(() => actionMeta.value?.needsTargets !== false);
    function toPayload(enabled) {
      return {
        name: form.value.name.trim(),
        cron: composedCron.value,
        action: form.value.action,
        targets: needsTargets.value ? [...form.value.targets] : [],
        enabled
      };
    }
    function blockedTarget(c) {
      return !!c.self || excluded.value.includes(c.name);
    }
    function blockedReason(c) {
      if (c.self) return "这是 Dockhelm 自己，任何计划任务都不会作用到它";
      if (excluded.value.includes(c.name)) return "在设置页的排除列表里，任何计划任务都会跳过它";
      return "选中这个容器";
    }
    async function load() {
      loading.value = true;
      try {
        const [s, a, c, st] = await Promise.all([
          api.get("/api/schedules"),
          api.get("/api/schedules/actions"),
          api.get("/api/containers"),
          api.get("/api/settings").catch(() => null)
        ]);
        schedules.value = s.schedules ?? [];
        actions.value = a.actions ?? [];
        containers.value = c.containers ?? [];
        excluded.value = st?.exclude ?? [];
      } catch (e) {
        toast.error("读取计划任务失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    function resetForm() {
      form.value = blankForm();
      editingId.value = 0;
      saveError.value = "";
    }
    function focusForm() {
      resetForm();
      document.getElementById("schedule-form")?.scrollIntoView({ behavior: "smooth", block: "nearest" });
    }
    function startEdit(s) {
      const mapped = applyCron(s.cron);
      form.value = {
        name: s.name,
        action: s.action,
        targets: [...s.targets ?? []],
        enabled: s.enabled,
        repeat: mapped.repeat,
        time: mapped.time,
        customCron: s.cron
      };
      editingId.value = s.id;
      saveError.value = "";
      document.getElementById("schedule-form")?.scrollIntoView({ behavior: "smooth", block: "nearest" });
    }
    async function save() {
      if (!form.value.name.trim()) {
        saveError.value = "请填写任务名称";
        return;
      }
      if (needsTargets.value && !form.value.targets.length) {
        saveError.value = "请至少选择一个容器";
        return;
      }
      saveError.value = "";
      saving.value = true;
      try {
        if (editingId.value) {
          await api.put(`/api/schedules/${editingId.value}`, toPayload(form.value.enabled));
          toast.success("任务已更新");
        } else {
          await api.post("/api/schedules", toPayload(form.value.enabled));
          toast.success("任务已创建");
        }
        resetForm();
        await load();
      } catch (e) {
        saveError.value = e instanceof Error ? e.message : String(e);
      } finally {
        saving.value = false;
      }
    }
    async function toggleEnabled(s) {
      try {
        await api.post(`/api/schedules/${s.id}/enabled`, { enabled: !s.enabled });
        await load();
      } catch (e) {
        toast.error("更新失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function runNow(s) {
      running.value = s.id;
      try {
        const res = await api.post(`/api/schedules/${s.id}/run`);
        if (res.ok) toast.success("执行完成", res.message);
        else toast.error("执行失败", res.message);
        await load();
      } catch (e) {
        toast.error("执行失败", e instanceof Error ? e.message : String(e));
      } finally {
        running.value = null;
      }
    }
    async function confirmRemove() {
      const s = removeTarget.value;
      if (!s || removing.value) return;
      removing.value = true;
      try {
        await api.del(`/api/schedules/${s.id}`);
        toast.success("任务已删除");
        if (editingId.value === s.id) resetForm();
        removeTarget.value = null;
        await load();
      } catch (e) {
        toast.error("删除失败", e instanceof Error ? e.message : String(e));
      } finally {
        removing.value = false;
      }
    }
    function toggleTarget(name) {
      const list = form.value.targets;
      if (list.includes(name)) form.value.targets = list.filter((t) => t !== name);
      else form.value.targets = [...list, name];
    }
    const statusBadge = (s) => {
      if (!s.lastStatus) return { text: "未执行", cls: "dh-badge-plain" };
      if (s.lastStatus === "success") return { text: "成功", cls: "dh-badge-run" };
      if (s.lastStatus === "failed") return { text: "失败", cls: "dh-badge-err" };
      return { text: s.lastStatus, cls: "dh-badge-plain" };
    };
    const summary = computed(() => {
      const todayEnd = /* @__PURE__ */ new Date();
      todayEnd.setHours(23, 59, 59, 999);
      const today = schedules.value.filter((s) => {
        if (!s.enabled || !s.nextRun) return false;
        const t = Date.parse(s.nextRun);
        return !Number.isNaN(t) && t <= todayEnd.getTime();
      }).length;
      return `${schedules.value.length} 个任务 · ${today} 个今天待执行`;
    });
    const upcoming = computed(() => {
      const now = Date.now();
      return schedules.value.filter((s) => s.enabled && s.nextRun).map((s) => ({ s, t: Date.parse(s.nextRun) })).filter((x) => !Number.isNaN(x.t) && x.t <= now + 24 * 3600 * 1e3).sort((a, b) => a.t - b.t).slice(0, 6);
    });
    function targetsText(s) {
      if (s.targets?.length) return s.targets;
      if (s.action === "prune_images") return ["不适用（全局动作）"];
      return ["未指定 —— 任务不会执行"];
    }
    const historyOpen = ref(false);
    const historyLogs = ref([]);
    const historyLoading = ref(false);
    async function openHistory() {
      historyOpen.value = true;
      historyLoading.value = true;
      try {
        const res = await api.get("/api/logs?limit=200");
        historyLogs.value = (res.logs ?? []).filter((l) => l.kind === "schedule").slice(0, 60);
      } catch (e) {
        toast.error("读取执行历史失败", e instanceof Error ? e.message : String(e));
      } finally {
        historyLoading.value = false;
      }
    }
    function historyStatus(l) {
      if (l.status === "success") return { text: "成功", cls: "dh-badge-run" };
      if (l.status === "failed") return { text: "失败", cls: "dh-badge-err" };
      return { text: l.status || "—", cls: "dh-badge-plain" };
    }
    onMounted(() => {
      void load();
      closeStream = openStream("/api/events/stream", (topic) => {
        if (topic === "schedule") void load();
      });
    });
    onUnmounted(() => closeStream?.());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[10] || (_cache[10] = createBaseVNode("div", { class: "dh-h1" }, "计划任务", -1)),
          createBaseVNode("div", _hoisted_3, toDisplayString(summary.value), 1),
          createBaseVNode("div", _hoisted_4, [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: openHistory
            }, "执行历史"),
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: loading.value,
              onClick: load
            }, [
              createVNode(unref(RefreshCw), {
                class: normalizeClass(["h-3.5 w-3.5", loading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[8] || (_cache[8] = createTextVNode("刷新 ", -1))
            ], 8, _hoisted_5),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: focusForm
            }, [
              createVNode(unref(Plus), { class: "h-3.5 w-3.5" }),
              _cache[9] || (_cache[9] = createTextVNode("新建任务 ", -1))
            ])
          ])
        ]),
        createBaseVNode("div", _hoisted_6, [
          !schedules.value.length ? (openBlock(), createBlock(_sfc_main$1, {
            key: 0,
            icon: unref(Plus),
            title: loading.value ? "正在载入…" : "还没有计划任务",
            description: "在下面的「新建任务」里挑几个容器、选个时间和动作，就能定时启停或更新它们。"
          }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_7, [
            createBaseVNode("table", _hoisted_8, [
              _cache[12] || (_cache[12] = createBaseVNode("thead", null, [
                createBaseVNode("tr", null, [
                  createBaseVNode("th", null, "任务名称"),
                  createBaseVNode("th", { class: "w-[150px]" }, "目标"),
                  createBaseVNode("th", { class: "w-[110px]" }, "动作"),
                  createBaseVNode("th", { class: "w-[170px]" }, "计划"),
                  createBaseVNode("th", { class: "w-[150px]" }, "上次执行"),
                  createBaseVNode("th", { class: "w-[150px]" }, "下次执行"),
                  createBaseVNode("th", { class: "w-[90px] text-right" }, "启用"),
                  createBaseVNode("th", { class: "w-[140px]" })
                ])
              ], -1)),
              createBaseVNode("tbody", null, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(schedules.value, (s) => {
                  return openBlock(), createElementBlock("tr", {
                    key: s.id,
                    class: normalizeClass(s.enabled ? "" : "opacity-55")
                  }, [
                    createBaseVNode("td", null, [
                      createBaseVNode("div", _hoisted_9, toDisplayString(s.name), 1),
                      s.lastMessage ? (openBlock(), createElementBlock("div", {
                        key: 0,
                        class: "mt-0.5 max-w-[280px] truncate text-[11.5px] text-text-5",
                        title: s.lastMessage
                      }, toDisplayString(s.lastMessage), 9, _hoisted_10)) : createCommentVNode("", true)
                    ]),
                    createBaseVNode("td", _hoisted_11, [
                      createBaseVNode("div", {
                        class: normalizeClass(["line-clamp-2", !s.targets?.length && s.action !== "prune_images" ? "text-err-text" : ""])
                      }, toDisplayString(targetsText(s).join("、")), 3)
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", actionTone(s.action)])
                      }, toDisplayString(shortAction(s.action)), 3)
                    ]),
                    createBaseVNode("td", _hoisted_12, [
                      createBaseVNode("div", null, toDisplayString(unref(explainCron)(s.cron)), 1),
                      createBaseVNode("div", _hoisted_13, toDisplayString(s.cron), 1)
                    ]),
                    createBaseVNode("td", _hoisted_14, [
                      createBaseVNode("div", null, toDisplayString(s.lastRun ? unref(formatDayTime)(s.lastRun) : "—"), 1),
                      s.lastStatus ? (openBlock(), createElementBlock("div", {
                        key: 0,
                        class: normalizeClass(s.lastStatus === "failed" ? "text-err-text" : "text-run-text")
                      }, [
                        createTextVNode(toDisplayString(statusBadge(s).text), 1),
                        s.lastStatus === "success" && s.targets?.length ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                          createTextVNode(toDisplayString(" " + s.targets.length), 1)
                        ], 64)) : createCommentVNode("", true)
                      ], 2)) : createCommentVNode("", true)
                    ]),
                    createBaseVNode("td", _hoisted_15, toDisplayString(s.enabled ? s.nextRun ? unref(formatDayTime)(s.nextRun) : "—" : "已停用"), 1),
                    createBaseVNode("td", _hoisted_16, [
                      createVNode(_sfc_main$2, {
                        "model-value": s.enabled,
                        label: `启用 ${s.name}`,
                        "onUpdate:modelValue": ($event) => toggleEnabled(s)
                      }, null, 8, ["model-value", "label", "onUpdate:modelValue"])
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("div", _hoisted_17, [
                        createBaseVNode("button", {
                          class: "dh-btn dh-btn-sm",
                          disabled: running.value === s.id,
                          onClick: ($event) => runNow(s)
                        }, [
                          running.value === s.id ? (openBlock(), createBlock(unref(LoaderCircle), {
                            key: 0,
                            class: "h-3 w-3 dh-spin"
                          })) : (openBlock(), createBlock(unref(Play), {
                            key: 1,
                            class: "h-3 w-3"
                          })),
                          _cache[11] || (_cache[11] = createTextVNode("运行 ", -1))
                        ], 8, _hoisted_18),
                        createBaseVNode("button", {
                          class: "dh-btn dh-btn-sm",
                          onClick: ($event) => startEdit(s)
                        }, "编辑", 8, _hoisted_19),
                        createBaseVNode("button", {
                          class: "dh-btn dh-btn-sm dh-btn-danger",
                          onClick: ($event) => removeTarget.value = s
                        }, [
                          createVNode(unref(Trash2), { class: "h-3 w-3" })
                        ], 8, _hoisted_20)
                      ])
                    ])
                  ], 2);
                }), 128))
              ])
            ])
          ]))
        ]),
        createBaseVNode("div", _hoisted_21, [
          createBaseVNode("div", _hoisted_22, [
            createBaseVNode("div", _hoisted_23, [
              createBaseVNode("span", null, toDisplayString(editingId.value ? "编辑任务" : "新建任务"), 1),
              editingId.value ? (openBlock(), createElementBlock("button", {
                key: 0,
                class: "ml-auto text-[11.5px] font-normal text-text-5 hover:text-accent",
                onClick: resetForm
              }, " 取消编辑 ")) : createCommentVNode("", true)
            ]),
            createBaseVNode("div", _hoisted_24, [
              createBaseVNode("div", null, [
                _cache[13] || (_cache[13] = createBaseVNode("label", { class: "dh-label" }, "任务名称", -1)),
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => form.value.name = $event),
                  class: "dh-input",
                  placeholder: "例如：夜间重启下载器"
                }, null, 512), [
                  [vModelText, form.value.name]
                ])
              ]),
              needsTargets.value ? (openBlock(), createElementBlock("div", _hoisted_25, [
                _cache[14] || (_cache[14] = createBaseVNode("label", { class: "dh-label" }, [
                  createTextVNode(" 选择容器 "),
                  createBaseVNode("span", { class: "text-text-6" }, "（必选；Dockhelm 自身与排除列表里的容器选不了）")
                ], -1)),
                createBaseVNode("div", _hoisted_26, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(containers.value, (c) => {
                    return openBlock(), createElementBlock("button", {
                      key: c.id,
                      type: "button",
                      class: "dh-chip",
                      "data-on": form.value.targets.includes(c.name),
                      disabled: blockedTarget(c),
                      title: blockedReason(c),
                      onClick: ($event) => toggleTarget(c.name)
                    }, [
                      createTextVNode(toDisplayString(c.name) + " ", 1),
                      c.self ? (openBlock(), createElementBlock("span", _hoisted_28, "自身")) : excluded.value.includes(c.name) ? (openBlock(), createElementBlock("span", _hoisted_29, "已排除")) : createCommentVNode("", true)
                    ], 8, _hoisted_27);
                  }), 128)),
                  !containers.value.length ? (openBlock(), createElementBlock("div", _hoisted_30, "读不到容器列表")) : createCommentVNode("", true)
                ]),
                createBaseVNode("div", {
                  class: normalizeClass(["mt-1 text-[11px]", form.value.targets.length ? "text-text-5" : "text-err-text"])
                }, [
                  form.value.targets.length ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                    createTextVNode("已选 " + toDisplayString(form.value.targets.length) + " 个", 1)
                  ], 64)) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
                    createTextVNode("还没选容器 —— 至少要选一个，任务才会执行")
                  ], 64))
                ], 2)
              ])) : createCommentVNode("", true),
              createBaseVNode("div", null, [
                _cache[15] || (_cache[15] = createBaseVNode("label", { class: "dh-label" }, "执行动作", -1)),
                createBaseVNode("div", _hoisted_31, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(actions.value, (a) => {
                    return openBlock(), createElementBlock("button", {
                      key: a.key,
                      type: "button",
                      class: "dh-chip",
                      "data-on": form.value.action === a.key,
                      title: a.description,
                      onClick: ($event) => form.value.action = a.key
                    }, toDisplayString(shortAction(a.key)), 9, _hoisted_32);
                  }), 128))
                ])
              ]),
              createBaseVNode("div", _hoisted_33, [
                createBaseVNode("div", null, [
                  _cache[16] || (_cache[16] = createBaseVNode("label", { class: "dh-label" }, "重复", -1)),
                  withDirectives(createBaseVNode("select", {
                    "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => form.value.repeat = $event),
                    class: "dh-select"
                  }, [
                    (openBlock(), createElementBlock(Fragment, null, renderList(repeatOptions, (r) => {
                      return createBaseVNode("option", {
                        key: r.key,
                        value: r.key
                      }, toDisplayString(r.label), 9, _hoisted_34);
                    }), 64))
                  ], 512), [
                    [vModelSelect, form.value.repeat]
                  ])
                ]),
                createBaseVNode("div", null, [
                  createBaseVNode("label", _hoisted_35, toDisplayString(form.value.repeat === "hourly" ? "每小时的第几分" : "时间"), 1),
                  form.value.repeat !== "custom" ? withDirectives((openBlock(), createElementBlock("input", {
                    key: 0,
                    "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => form.value.time = $event),
                    type: "time",
                    class: "dh-input font-mono"
                  }, null, 512)), [
                    [vModelText, form.value.time]
                  ]) : withDirectives((openBlock(), createElementBlock("input", {
                    key: 1,
                    "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => form.value.customCron = $event),
                    class: "dh-input font-mono",
                    placeholder: "0 3 * * *"
                  }, null, 512)), [
                    [vModelText, form.value.customCron]
                  ])
                ])
              ]),
              createBaseVNode("div", _hoisted_36, [
                _cache[17] || (_cache[17] = createTextVNode(" 解析结果：", -1)),
                createBaseVNode("b", _hoisted_37, toDisplayString(composedCron.value), 1),
                unref(explainCron)(composedCron.value) && unref(explainCron)(composedCron.value) !== composedCron.value ? (openBlock(), createElementBlock("span", _hoisted_38, " · " + toDisplayString(unref(explainCron)(composedCron.value)), 1)) : createCommentVNode("", true)
              ]),
              createBaseVNode("label", _hoisted_39, [
                createVNode(_sfc_main$2, {
                  modelValue: form.value.enabled,
                  "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => form.value.enabled = $event)
                }, null, 8, ["modelValue"]),
                _cache[18] || (_cache[18] = createTextVNode(" 创建后立即启用 ", -1))
              ]),
              saveError.value ? (openBlock(), createElementBlock("div", _hoisted_40, toDisplayString(saveError.value), 1)) : createCommentVNode("", true),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-primary",
                disabled: saving.value,
                onClick: save
              }, [
                saving.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "h-3.5 w-3.5 dh-spin"
                })) : createCommentVNode("", true),
                createTextVNode(" " + toDisplayString(editingId.value ? "保存修改" : "创建任务"), 1)
              ], 8, _hoisted_41)
            ])
          ]),
          createBaseVNode("div", _hoisted_42, [
            _cache[19] || (_cache[19] = createBaseVNode("div", { class: "dh-card-head" }, "时间轴 · 未来 24 小时", -1)),
            !upcoming.value.length ? (openBlock(), createElementBlock("div", _hoisted_43, " 未来 24 小时内没有安排。 ")) : (openBlock(), createElementBlock("div", _hoisted_44, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(upcoming.value, (u) => {
                return openBlock(), createElementBlock("div", {
                  key: u.s.id,
                  class: "flex items-center gap-2.5"
                }, [
                  createBaseVNode("span", _hoisted_45, toDisplayString(unref(formatDayTime)(u.t)), 1),
                  createBaseVNode("span", _hoisted_46, toDisplayString(u.s.name), 1),
                  createBaseVNode("span", {
                    class: normalizeClass(["dh-badge flex-none", actionTone(u.s.action)])
                  }, toDisplayString(shortAction(u.s.action)), 3)
                ]);
              }), 128))
            ]))
          ])
        ]),
        createVNode(_sfc_main$3, {
          open: !!removeTarget.value,
          title: "删除计划任务",
          subtitle: removeTarget.value?.name,
          width: "400px",
          busy: removing.value,
          onClose: _cache[6] || (_cache[6] = ($event) => removeTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[5] || (_cache[5] = ($event) => removeTarget.value = null)
            }, "取消", 8, _hoisted_47),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemove
            }, toDisplayString(removing.value ? "删除中…" : "确认删除"), 9, _hoisted_48)
          ]),
          default: withCtx(() => [
            _cache[20] || (_cache[20] = createBaseVNode("div", { class: "text-[12.5px] text-text-3" }, "删除后该任务不再自动执行，已有的执行记录仍然保留。", -1))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: historyOpen.value,
          title: "执行历史",
          subtitle: "仅计划任务",
          width: "720px",
          onClose: _cache[7] || (_cache[7] = ($event) => historyOpen.value = false)
        }, {
          default: withCtx(() => [
            historyLoading.value ? (openBlock(), createElementBlock("div", _hoisted_49, [
              createVNode(unref(LoaderCircle), { class: "h-4 w-4 dh-spin" })
            ])) : !historyLogs.value.length ? (openBlock(), createElementBlock("div", _hoisted_50, " 还没有计划任务的执行记录。 ")) : (openBlock(), createElementBlock("div", _hoisted_51, [
              createBaseVNode("table", _hoisted_52, [
                _cache[21] || (_cache[21] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", { class: "w-[150px]" }, "时间"),
                    createBaseVNode("th", { class: "w-[180px]" }, "任务"),
                    createBaseVNode("th", null, "说明"),
                    createBaseVNode("th", { class: "w-[90px]" }, "结果")
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(historyLogs.value, (l) => {
                    return openBlock(), createElementBlock("tr", {
                      key: l.id
                    }, [
                      createBaseVNode("td", _hoisted_53, toDisplayString(unref(formatDayTime)(l.ts)), 1),
                      createBaseVNode("td", _hoisted_54, toDisplayString(l.ref || "—"), 1),
                      createBaseVNode("td", _hoisted_55, toDisplayString(l.message), 1),
                      createBaseVNode("td", null, [
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge", historyStatus(l).cls])
                        }, toDisplayString(historyStatus(l).text), 3)
                      ])
                    ]);
                  }), 128))
                ])
              ])
            ]))
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
