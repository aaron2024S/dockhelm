import { x as createLucideIcon, d as defineComponent, D as useToastStore, a as onMounted, c as createElementBlock, e as createBaseVNode, t as toDisplayString, f as createVNode, g as unref, k as createTextVNode, j as createCommentVNode, F as Fragment, P as Info, i as createBlock, Y as Bell, r as renderList, w as withDirectives, E as vModelText, L as vModelSelect, H as withCtx, m as ref, p as computed, B as api, s as openBlock, n as normalizeClass, T as CircleCheck, U as CircleX } from "./index-BWcOCF5V.js";
import { b as relativeTime } from "./format-NZu8lNOP.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-B2RxRpEh.js";
import { _ as _sfc_main$3 } from "./Modal.vue_vue_type_script_setup_true_lang-CY0B9bsm.js";
import { _ as _sfc_main$2 } from "./UiSwitch.vue_vue_type_script_setup_true_lang-BYQIoEEw.js";
import { P as Plus } from "./plus-CWvN5_LN.js";
import { T as Trash2 } from "./trash-2-BvtUCprF.js";
import { S as Save } from "./save-CGjumFrh.js";
import { L as LoaderCircle } from "./loader-circle-BEgatAk7.js";
const BellOff = createLucideIcon("BellOffIcon", [
  ["path", { d: "M10.268 21a2 2 0 0 0 3.464 0", key: "vwvbt9" }],
  [
    "path",
    {
      d: "M17 17H4a1 1 0 0 1-.74-1.673C4.59 13.956 6 12.499 6 8a6 6 0 0 1 .258-1.742",
      key: "178tsu"
    }
  ],
  ["path", { d: "m2 2 20 20", key: "1ooewy" }],
  ["path", { d: "M8.668 3.01A6 6 0 0 1 18 8c0 2.687.77 4.653 1.707 6.05", key: "1hqiys" }]
]);
const History = createLucideIcon("HistoryIcon", [
  ["path", { d: "M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8", key: "1357e3" }],
  ["path", { d: "M3 3v5h5", key: "1xhq8a" }],
  ["path", { d: "M12 7v5l4 2", key: "1fdv2h" }]
]);
const Send = createLucideIcon("SendIcon", [
  [
    "path",
    {
      d: "M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z",
      key: "1ffxy3"
    }
  ],
  ["path", { d: "m21.854 2.147-10.94 10.939", key: "12cjpa" }]
]);
const Settings2 = createLucideIcon("Settings2Icon", [
  ["path", { d: "M20 7h-9", key: "3s1dr2" }],
  ["path", { d: "M14 17H5", key: "gfn3mx" }],
  ["circle", { cx: "17", cy: "17", r: "3", key: "18b49y" }],
  ["circle", { cx: "7", cy: "7", r: "3", key: "dfmy0x" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "dh-sub" };
const _hoisted_4 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_5 = { class: "dh-seg" };
const _hoisted_6 = ["data-on"];
const _hoisted_7 = ["data-on"];
const _hoisted_8 = ["data-on"];
const _hoisted_9 = { class: "ml-auto flex items-center gap-2" };
const _hoisted_10 = ["disabled"];
const _hoisted_11 = {
  key: 1,
  class: "dh-card"
};
const _hoisted_12 = { class: "dh-card-head" };
const _hoisted_13 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_14 = { class: "flex flex-col" };
const _hoisted_15 = { class: "min-w-[150px] flex-1" };
const _hoisted_16 = { class: "text-[12.5px] font-medium" };
const _hoisted_17 = { class: "text-[11px] text-text-5" };
const _hoisted_18 = { class: "flex gap-1.5" };
const _hoisted_19 = ["disabled", "onClick"];
const _hoisted_20 = ["onClick"];
const _hoisted_21 = ["onClick"];
const _hoisted_22 = {
  key: 2,
  class: "dh-card"
};
const _hoisted_23 = { class: "dh-card-head" };
const _hoisted_24 = ["disabled"];
const _hoisted_25 = { class: "grid grid-cols-1 gap-x-8 gap-y-4 p-3.5 lg:grid-cols-2" };
const _hoisted_26 = { class: "flex items-start gap-2.5" };
const _hoisted_27 = { class: "dh-label" };
const _hoisted_28 = { class: "flex items-start gap-2.5" };
const _hoisted_29 = { class: "flex items-end gap-2" };
const _hoisted_30 = { class: "flex-1" };
const _hoisted_31 = { class: "flex-1" };
const _hoisted_32 = { class: "flex items-start gap-2.5" };
const _hoisted_33 = { class: "mt-1 text-[11px] text-text-5" };
const _hoisted_34 = { class: "lg:col-span-2" };
const _hoisted_35 = { class: "rounded-[9px] border border-line-1 bg-ink-800 px-3 py-2.5 text-[11.5px] leading-relaxed text-text-5" };
const _hoisted_36 = { class: "flex items-start gap-3 rounded-[14px] border border-line-1 bg-ink-750 px-4 py-3 text-[12px] leading-relaxed text-text-4" };
const _hoisted_37 = { class: "dh-card-head" };
const _hoisted_38 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_39 = { class: "flex flex-col" };
const _hoisted_40 = { class: "min-w-[200px] flex-1" };
const _hoisted_41 = { class: "flex items-center gap-2" };
const _hoisted_42 = { class: "text-[12.5px] font-medium" };
const _hoisted_43 = { class: "text-[10.5px] text-text-6" };
const _hoisted_44 = { class: "mt-0.5 text-[11px] leading-relaxed text-text-5" };
const _hoisted_45 = {
  key: 2,
  class: "dh-card"
};
const _hoisted_46 = { class: "dh-card-head" };
const _hoisted_47 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_48 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_49 = { class: "dh-table" };
const _hoisted_50 = { class: "whitespace-nowrap text-[11.5px] text-text-5" };
const _hoisted_51 = { class: "text-[12px] text-text-3" };
const _hoisted_52 = { class: "max-w-[420px]" };
const _hoisted_53 = ["title"];
const _hoisted_54 = ["title"];
const _hoisted_55 = {
  key: 0,
  class: "flex flex-col gap-3.5"
};
const _hoisted_56 = { class: "grid grid-cols-2 gap-3" };
const _hoisted_57 = ["disabled"];
const _hoisted_58 = ["value"];
const _hoisted_59 = {
  key: 0,
  class: "text-[11.5px] leading-relaxed text-text-5"
};
const _hoisted_60 = { class: "flex flex-col gap-3" };
const _hoisted_61 = { class: "dh-label" };
const _hoisted_62 = {
  key: 0,
  class: "text-err-text"
};
const _hoisted_63 = ["value", "placeholder", "onInput"];
const _hoisted_64 = ["value", "type", "placeholder", "onInput"];
const _hoisted_65 = {
  key: 2,
  class: "mt-1 text-[11px] text-text-5"
};
const _hoisted_66 = { class: "flex cursor-pointer items-center gap-2.5 text-[12.5px] text-text-2" };
const _hoisted_67 = { class: "text-[12px] leading-relaxed text-text-4" };
const _hoisted_68 = { class: "text-text-2" };
const _hoisted_69 = { class: "mt-3 overflow-hidden rounded-[10px] border border-line-1" };
const _hoisted_70 = { class: "dh-table" };
const _hoisted_71 = { class: "font-mono text-[11.5px] text-accent" };
const _hoisted_72 = { class: "text-[12px] text-text-3" };
const _hoisted_73 = ["disabled"];
const _hoisted_74 = ["disabled"];
const _hoisted_75 = ["disabled"];
const _hoisted_76 = ["disabled"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "NotifyView",
  setup(__props) {
    const toast = useToastStore();
    const tab = ref("channels");
    const presets = ref([]);
    const templateVars = ref([]);
    const channels = ref([]);
    const catalog = ref([]);
    const eventRows = ref({});
    const settings = ref(null);
    const history = ref([]);
    const loading = ref(true);
    const saving = ref(false);
    const removing = ref(false);
    const confirmClearHistory = ref(false);
    const clearingHistory = ref(false);
    const testingId = ref(null);
    const showEditor = ref(false);
    const editingChannel = ref(null);
    const removeTarget = ref(null);
    const showVars = ref(false);
    const groups = computed(() => {
      const out = [];
      for (const e of catalog.value) {
        let g = out.find((x) => x.name === e.group);
        if (!g) {
          g = { name: e.group, items: [] };
          out.push(g);
        }
        g.items.push(e);
      }
      return out;
    });
    const enabledCount = computed(
      () => catalog.value.filter((e) => eventRows.value[e.event]?.enabled).length
    );
    const currentPreset = computed(
      () => presets.value.find((p) => p.type === editingChannel.value?.type)
    );
    async function load() {
      loading.value = true;
      try {
        const [p, e, s, ch, h] = await Promise.all([
          api.get("/api/notify/presets"),
          api.get("/api/notify/events"),
          api.get("/api/notify/settings"),
          api.get("/api/notify/channels"),
          api.get("/api/notify/history", { limit: 80 })
        ]);
        presets.value = p.presets ?? [];
        templateVars.value = p.vars ?? [];
        catalog.value = e.catalog ?? [];
        const map = {};
        for (const def of e.catalog ?? []) {
          const cur = (e.current ?? []).find((c) => c.event === def.event);
          map[def.event] = cur ?? { event: def.event, enabled: def.default, level: def.level };
        }
        eventRows.value = map;
        settings.value = s.settings;
        channels.value = ch.channels ?? [];
        history.value = h.history ?? [];
      } catch (err) {
        toast.error("读取通知配置失败", err instanceof Error ? err.message : String(err));
      } finally {
        loading.value = false;
      }
    }
    function openCreate() {
      const def = presets.value[0];
      editingChannel.value = {
        id: 0,
        name: def?.label ?? "新渠道",
        type: def?.type ?? "webhook",
        enabled: true,
        config: {},
        createdAt: ""
      };
      showEditor.value = true;
    }
    function openEdit(c) {
      editingChannel.value = { ...c, config: { ...c.config ?? {} } };
      showEditor.value = true;
    }
    function needField(key) {
      const p = currentPreset.value;
      const v = editingChannel.value?.config?.[key];
      return p?.fields.find((f) => f.key === key)?.required && !String(v ?? "").trim();
    }
    async function saveChannel() {
      const ch = editingChannel.value;
      if (!ch) return;
      if (!ch.name.trim()) {
        toast.error("请填写渠道名称");
        return;
      }
      try {
        if (ch.id) {
          await api.put(`/api/notify/channels/${ch.id}`, ch);
          toast.success("渠道已更新");
        } else {
          await api.post("/api/notify/channels", ch);
          toast.success("渠道已创建");
        }
        showEditor.value = false;
        await load();
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function toggleChannel(c, v) {
      const prev = c.enabled;
      c.enabled = v;
      try {
        await api.put(`/api/notify/channels/${c.id}`, c);
      } catch (e) {
        c.enabled = prev;
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function testChannel(c) {
      testingId.value = c.id;
      try {
        const res = await api.post(`/api/notify/channels/${c.id}/test`);
        if (res.ok) toast.success("测试消息已发送", `${c.name} 配置可用`);
      } catch (e) {
        toast.error("测试失败", e instanceof Error ? e.message : String(e));
      } finally {
        testingId.value = null;
        void loadHistory();
      }
    }
    async function confirmRemove() {
      const c = removeTarget.value;
      if (!c || removing.value) return;
      removing.value = true;
      try {
        await api.del(`/api/notify/channels/${c.id}`);
        toast.success("渠道已删除");
        removeTarget.value = null;
        await load();
      } catch (e) {
        toast.error("删除失败", e instanceof Error ? e.message : String(e));
      } finally {
        removing.value = false;
      }
    }
    async function toggleEvent(def) {
      const cur = eventRows.value[def.event];
      if (!cur) return;
      cur.enabled = !cur.enabled;
      try {
        await api.put("/api/notify/events", { events: [cur] });
      } catch (e) {
        cur.enabled = !cur.enabled;
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function bulkEvents(preset) {
      const rows = catalog.value.map((def) => ({
        event: def.event,
        enabled: preset === "all" ? true : preset === "none" ? false : def.default,
        level: eventRows.value[def.event]?.level ?? def.level
      }));
      try {
        await api.put("/api/notify/events", { events: rows });
        for (const r of rows) eventRows.value[r.event] = r;
        toast.success("已应用");
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function saveSettings() {
      if (!settings.value) return;
      saving.value = true;
      try {
        const res = await api.put("/api/notify/settings", settings.value);
        settings.value = res.settings;
        toast.success("通知设置已保存");
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      } finally {
        saving.value = false;
      }
    }
    async function loadHistory() {
      try {
        const res = await api.get("/api/notify/history", { limit: 80 });
        history.value = res.history ?? [];
      } catch {
      }
    }
    async function clearHistory() {
      if (clearingHistory.value) return;
      clearingHistory.value = true;
      confirmClearHistory.value = false;
      try {
        await api.del("/api/notify/history");
        history.value = [];
        toast.success("推送历史已清空");
      } catch (e) {
        toast.error("清空失败", e instanceof Error ? e.message : String(e));
      } finally {
        clearingHistory.value = false;
      }
    }
    const levelBadge = (lv) => lv === "urgent" ? "dh-badge-err" : "dh-badge-plain";
    const tplVar = (name) => `{{${name}}}`;
    onMounted(() => void load());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[28] || (_cache[28] = createBaseVNode("div", { class: "dh-h1" }, "通知", -1)),
          createBaseVNode("div", _hoisted_3, toDisplayString(channels.value.filter((c) => c.enabled).length) + " 个渠道已启用 · 已订阅 " + toDisplayString(enabledCount.value) + " 个事件 ", 1)
        ]),
        createBaseVNode("div", _hoisted_4, [
          createBaseVNode("div", _hoisted_5, [
            createBaseVNode("button", {
              "data-on": tab.value === "channels",
              onClick: _cache[0] || (_cache[0] = ($event) => tab.value = "channels")
            }, " 渠道 " + toDisplayString(channels.value.length), 9, _hoisted_6),
            createBaseVNode("button", {
              "data-on": tab.value === "events",
              onClick: _cache[1] || (_cache[1] = ($event) => tab.value = "events")
            }, " 事件订阅 " + toDisplayString(enabledCount.value), 9, _hoisted_7),
            createBaseVNode("button", {
              "data-on": tab.value === "history",
              onClick: _cache[2] || (_cache[2] = ($event) => tab.value = "history")
            }, "推送历史", 8, _hoisted_8)
          ]),
          createBaseVNode("div", _hoisted_9, [
            tab.value === "channels" ? (openBlock(), createElementBlock("button", {
              key: 0,
              class: "dh-btn dh-btn-primary",
              onClick: openCreate
            }, [
              createVNode(unref(Plus), { class: "h-3.5 w-3.5" }),
              _cache[29] || (_cache[29] = createTextVNode("添加渠道 ", -1))
            ])) : createCommentVNode("", true),
            tab.value === "events" ? (openBlock(), createElementBlock(Fragment, { key: 1 }, [
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm",
                onClick: _cache[3] || (_cache[3] = ($event) => bulkEvents("recommended"))
              }, "恢复推荐"),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm",
                onClick: _cache[4] || (_cache[4] = ($event) => bulkEvents("none"))
              }, "全部关闭"),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm",
                onClick: _cache[5] || (_cache[5] = ($event) => showVars.value = true)
              }, [
                createVNode(unref(Info), { class: "h-3 w-3" }),
                _cache[30] || (_cache[30] = createTextVNode("模板变量 ", -1))
              ])
            ], 64)) : createCommentVNode("", true),
            tab.value === "history" ? (openBlock(), createElementBlock(Fragment, { key: 2 }, [
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm",
                onClick: loadHistory
              }, "刷新"),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-danger",
                disabled: !history.value.length || clearingHistory.value,
                onClick: _cache[6] || (_cache[6] = ($event) => confirmClearHistory.value = true)
              }, [
                createVNode(unref(Trash2), { class: "h-3 w-3" }),
                _cache[31] || (_cache[31] = createTextVNode("清空 ", -1))
              ], 8, _hoisted_10)
            ], 64)) : createCommentVNode("", true)
          ])
        ]),
        tab.value === "channels" ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
          !channels.value.length ? (openBlock(), createBlock(_sfc_main$1, {
            key: 0,
            icon: unref(BellOff),
            title: loading.value ? "正在载入…" : "还没有配置通知渠道",
            description: "支持 Telegram、Bark、ntfy、企业微信、钉钉、飞书、Server 酱、PushPlus、邮件 SMTP，以及任意自定义 Webhook。",
            "action-label": "添加第一个渠道",
            onAction: openCreate
          }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_11, [
            createBaseVNode("div", _hoisted_12, [
              createVNode(unref(Bell), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[32] || (_cache[32] = createBaseVNode("span", null, "通知渠道", -1)),
              createBaseVNode("span", _hoisted_13, toDisplayString(channels.value.filter((c) => c.enabled).length) + " / " + toDisplayString(channels.value.length) + " 已启用 ", 1)
            ]),
            createBaseVNode("div", _hoisted_14, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(channels.value, (c) => {
                return openBlock(), createElementBlock("div", {
                  key: c.id,
                  class: "flex flex-wrap items-center gap-2.5 border-b border-line-row px-3.5 py-3 last:border-b-0"
                }, [
                  createVNode(_sfc_main$2, {
                    "model-value": c.enabled,
                    label: `启用 ${c.name}`,
                    "onUpdate:modelValue": (v) => toggleChannel(c, v)
                  }, null, 8, ["model-value", "label", "onUpdate:modelValue"]),
                  createBaseVNode("div", _hoisted_15, [
                    createBaseVNode("div", _hoisted_16, toDisplayString(c.name), 1),
                    createBaseVNode("div", _hoisted_17, toDisplayString(presets.value.find((p) => p.type === c.type)?.label ?? c.type), 1)
                  ]),
                  createBaseVNode("div", _hoisted_18, [
                    createBaseVNode("button", {
                      class: "dh-btn dh-btn-sm",
                      disabled: testingId.value === c.id,
                      onClick: ($event) => testChannel(c)
                    }, [
                      testingId.value === c.id ? (openBlock(), createBlock(unref(LoaderCircle), {
                        key: 0,
                        class: "h-3 w-3 dh-spin"
                      })) : (openBlock(), createBlock(unref(Send), {
                        key: 1,
                        class: "h-3 w-3"
                      })),
                      _cache[33] || (_cache[33] = createTextVNode("测试 ", -1))
                    ], 8, _hoisted_19),
                    createBaseVNode("button", {
                      class: "dh-btn dh-btn-sm",
                      onClick: ($event) => openEdit(c)
                    }, "编辑", 8, _hoisted_20),
                    createBaseVNode("button", {
                      class: "dh-btn dh-btn-sm dh-btn-danger",
                      onClick: ($event) => removeTarget.value = c
                    }, [
                      createVNode(unref(Trash2), { class: "h-3 w-3" })
                    ], 8, _hoisted_21)
                  ])
                ]);
              }), 128))
            ])
          ])),
          settings.value ? (openBlock(), createElementBlock("div", _hoisted_22, [
            createBaseVNode("div", _hoisted_23, [
              createVNode(unref(Settings2), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[35] || (_cache[35] = createBaseVNode("span", null, "投递策略", -1)),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm dh-btn-primary ml-auto",
                disabled: saving.value,
                onClick: saveSettings
              }, [
                createVNode(unref(Save), { class: "h-3 w-3" }),
                _cache[34] || (_cache[34] = createTextVNode("保存 ", -1))
              ], 8, _hoisted_24)
            ]),
            createBaseVNode("div", _hoisted_25, [
              createBaseVNode("label", _hoisted_26, [
                createVNode(_sfc_main$2, {
                  modelValue: settings.value.enabled,
                  "onUpdate:modelValue": _cache[7] || (_cache[7] = ($event) => settings.value.enabled = $event)
                }, null, 8, ["modelValue"]),
                _cache[36] || (_cache[36] = createBaseVNode("span", { class: "text-[12.5px]" }, [
                  createTextVNode(" 启用通知 "),
                  createBaseVNode("div", { class: "text-[11px] text-text-5" }, "总开关。关闭后不发送任何通知，也不写推送历史。")
                ], -1))
              ]),
              createBaseVNode("div", null, [
                createBaseVNode("label", _hoisted_27, [
                  _cache[37] || (_cache[37] = createTextVNode("面板地址（用于模板变量 ", -1)),
                  createBaseVNode("code", null, toDisplayString(tplVar("url")), 1),
                  _cache[38] || (_cache[38] = createTextVNode("）", -1))
                ]),
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[8] || (_cache[8] = ($event) => settings.value.panelURL = $event),
                  class: "dh-input",
                  placeholder: "http://192.168.1.10:5923"
                }, null, 512), [
                  [vModelText, settings.value.panelURL]
                ])
              ]),
              createBaseVNode("label", _hoisted_28, [
                createVNode(_sfc_main$2, {
                  modelValue: settings.value.quietEnabled,
                  "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => settings.value.quietEnabled = $event)
                }, null, 8, ["modelValue"]),
                _cache[39] || (_cache[39] = createBaseVNode("span", { class: "text-[12.5px]" }, [
                  createTextVNode(" 启用静默时段 "),
                  createBaseVNode("div", { class: "text-[11px] text-text-5" }, "支持跨天，例如 23:00–07:00。")
                ], -1))
              ]),
              createBaseVNode("div", _hoisted_29, [
                createBaseVNode("div", _hoisted_30, [
                  _cache[40] || (_cache[40] = createBaseVNode("label", { class: "dh-label" }, "开始", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[10] || (_cache[10] = ($event) => settings.value.quietStart = $event),
                    type: "time",
                    class: "dh-input"
                  }, null, 512), [
                    [vModelText, settings.value.quietStart]
                  ])
                ]),
                createBaseVNode("div", _hoisted_31, [
                  _cache[41] || (_cache[41] = createBaseVNode("label", { class: "dh-label" }, "结束", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[11] || (_cache[11] = ($event) => settings.value.quietEnd = $event),
                    type: "time",
                    class: "dh-input"
                  }, null, 512), [
                    [vModelText, settings.value.quietEnd]
                  ])
                ])
              ]),
              createBaseVNode("div", null, [
                _cache[43] || (_cache[43] = createBaseVNode("label", { class: "dh-label" }, "静默期的普通事件", -1)),
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[12] || (_cache[12] = ($event) => settings.value.quietNormalMode = $event),
                  class: "dh-select"
                }, [..._cache[42] || (_cache[42] = [
                  createBaseVNode("option", { value: "digest" }, "攒着，时段结束后汇总发一条", -1),
                  createBaseVNode("option", { value: "drop" }, "直接丢弃", -1)
                ])], 512), [
                  [vModelSelect, settings.value.quietNormalMode]
                ])
              ]),
              createBaseVNode("label", _hoisted_32, [
                createVNode(_sfc_main$2, {
                  modelValue: settings.value.quietUrgentSend,
                  "onUpdate:modelValue": _cache[13] || (_cache[13] = ($event) => settings.value.quietUrgentSend = $event)
                }, null, 8, ["modelValue"]),
                _cache[44] || (_cache[44] = createBaseVNode("span", { class: "text-[12.5px]" }, [
                  createTextVNode(" 紧急事件在静默期照常发送 "),
                  createBaseVNode("div", { class: "text-[11px] text-text-5" }, "更新失败、容器意外退出、登录告警属于紧急事件。")
                ], -1))
              ]),
              createBaseVNode("div", null, [
                _cache[45] || (_cache[45] = createBaseVNode("label", { class: "dh-label" }, "同容器同事件去重窗口（分钟）", -1)),
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[14] || (_cache[14] = ($event) => settings.value.dedupeWindow = $event),
                  type: "number",
                  min: "0",
                  class: "dh-input"
                }, null, 512), [
                  [
                    vModelText,
                    settings.value.dedupeWindow,
                    void 0,
                    { number: true }
                  ]
                ]),
                _cache[46] || (_cache[46] = createBaseVNode("div", { class: "mt-1 text-[11px] text-text-5" }, " 崩溃循环的容器一分钟能产生几十条事件，靠这个窗口压住。 ", -1))
              ]),
              createBaseVNode("div", null, [
                _cache[47] || (_cache[47] = createBaseVNode("label", { class: "dh-label" }, "每日推送上限（条）", -1)),
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[15] || (_cache[15] = ($event) => settings.value.dailyLimit = $event),
                  type: "number",
                  min: "0",
                  class: "dh-input"
                }, null, 512), [
                  [
                    vModelText,
                    settings.value.dailyLimit,
                    void 0,
                    { number: true }
                  ]
                ]),
                createBaseVNode("div", _hoisted_33, "今日已发送 " + toDisplayString(settings.value.sentToday) + " 条，超出上限后丢弃。", 1)
              ]),
              createBaseVNode("div", _hoisted_34, [
                createBaseVNode("div", _hoisted_35, [
                  _cache[48] || (_cache[48] = createTextVNode(" 投递失败时会自动重试 3 次并按指数退避，失败只记入推送历史， ", -1)),
                  _cache[49] || (_cache[49] = createBaseVNode("b", { class: "text-text-3" }, "绝不影响更新主流程", -1)),
                  createTextVNode("。 当前状态：" + toDisplayString(settings.value.inQuietHours ? "处于静默时段" : "非静默时段") + "。 ", 1)
                ])
              ])
            ])
          ])) : createCommentVNode("", true)
        ], 64)) : tab.value === "events" ? (openBlock(), createElementBlock(Fragment, { key: 1 }, [
          createBaseVNode("div", _hoisted_36, [
            createVNode(unref(Info), { class: "mt-[2px] h-4 w-4 flex-none text-accent" }),
            _cache[50] || (_cache[50] = createBaseVNode("div", null, [
              createTextVNode(" 订阅是按"),
              createBaseVNode("b", null, "事件"),
              createTextVNode("而不是按渠道设置的：某个事件一旦开启，会同时发往所有已启用的渠道。 标注「紧急」的事件可以在静默时段照常发送。登录成功与登录失败是两个独立事件，可以分别开关。 ")
            ], -1))
          ]),
          (openBlock(true), createElementBlock(Fragment, null, renderList(groups.value, (g) => {
            return openBlock(), createElementBlock("div", {
              key: g.name,
              class: "dh-card"
            }, [
              createBaseVNode("div", _hoisted_37, [
                createVNode(unref(Bell), { class: "h-3.5 w-3.5 text-text-4" }),
                createBaseVNode("span", null, toDisplayString(g.name), 1),
                createBaseVNode("span", _hoisted_38, toDisplayString(g.items.filter((i) => eventRows.value[i.event]?.enabled).length) + " / " + toDisplayString(g.items.length) + " 已开启 ", 1)
              ]),
              createBaseVNode("div", _hoisted_39, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(g.items, (e) => {
                  return openBlock(), createElementBlock("div", {
                    key: e.event,
                    class: "flex flex-wrap items-center gap-3 border-b border-line-row px-3.5 py-2.5 last:border-b-0"
                  }, [
                    createVNode(_sfc_main$2, {
                      "model-value": eventRows.value[e.event]?.enabled ?? false,
                      label: e.label,
                      "onUpdate:modelValue": ($event) => toggleEvent(e)
                    }, null, 8, ["model-value", "label", "onUpdate:modelValue"]),
                    createBaseVNode("div", _hoisted_40, [
                      createBaseVNode("div", _hoisted_41, [
                        createBaseVNode("span", _hoisted_42, toDisplayString(e.label), 1),
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge", levelBadge(e.level)])
                        }, toDisplayString(e.level === "urgent" ? "紧急" : "普通"), 3),
                        createBaseVNode("code", _hoisted_43, toDisplayString(e.event), 1)
                      ]),
                      createBaseVNode("div", _hoisted_44, toDisplayString(e.description), 1)
                    ])
                  ]);
                }), 128))
              ])
            ]);
          }), 128))
        ], 64)) : (openBlock(), createElementBlock("div", _hoisted_45, [
          createBaseVNode("div", _hoisted_46, [
            createVNode(unref(History), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[51] || (_cache[51] = createBaseVNode("span", null, "推送历史", -1)),
            createBaseVNode("span", _hoisted_47, toDisplayString(history.value.length) + " 条", 1)
          ]),
          !history.value.length ? (openBlock(), createBlock(_sfc_main$1, {
            key: 0,
            icon: unref(History),
            title: "还没有推送记录",
            description: "通知发出后（无论成功或失败）都会在这里留下一条记录。"
          }, null, 8, ["icon"])) : (openBlock(), createElementBlock("div", _hoisted_48, [
            createBaseVNode("table", _hoisted_49, [
              _cache[52] || (_cache[52] = createBaseVNode("thead", null, [
                createBaseVNode("tr", null, [
                  createBaseVNode("th", { class: "w-[140px]" }, "时间"),
                  createBaseVNode("th", { class: "w-[150px]" }, "渠道"),
                  createBaseVNode("th", { class: "w-[150px]" }, "事件"),
                  createBaseVNode("th", { class: "w-[70px]" }, "结果"),
                  createBaseVNode("th", null, "内容")
                ])
              ], -1)),
              createBaseVNode("tbody", null, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(history.value, (h) => {
                  return openBlock(), createElementBlock("tr", {
                    key: h.id
                  }, [
                    createBaseVNode("td", _hoisted_50, toDisplayString(unref(relativeTime)(h.ts)), 1),
                    createBaseVNode("td", _hoisted_51, toDisplayString(h.channel || "—"), 1),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", levelBadge(h.level)])
                      }, toDisplayString(h.event), 3)
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("span", {
                        class: normalizeClass(["dh-badge", h.ok ? "dh-badge-run" : "dh-badge-err"])
                      }, [
                        h.ok ? (openBlock(), createBlock(unref(CircleCheck), {
                          key: 0,
                          class: "h-3 w-3"
                        })) : (openBlock(), createBlock(unref(CircleX), {
                          key: 1,
                          class: "h-3 w-3"
                        })),
                        createTextVNode(" " + toDisplayString(h.ok ? "成功" : "失败"), 1)
                      ], 2)
                    ]),
                    createBaseVNode("td", _hoisted_52, [
                      createBaseVNode("div", {
                        class: "truncate text-[12px] text-text-2",
                        title: h.title
                      }, toDisplayString(h.title), 9, _hoisted_53),
                      !h.ok && h.errmsg ? (openBlock(), createElementBlock("div", {
                        key: 0,
                        class: "mt-0.5 truncate text-[11px] text-err-text",
                        title: h.errmsg
                      }, toDisplayString(h.errmsg), 9, _hoisted_54)) : createCommentVNode("", true)
                    ])
                  ]);
                }), 128))
              ])
            ])
          ]))
        ])),
        createVNode(_sfc_main$3, {
          open: showEditor.value,
          title: editingChannel.value?.id ? "编辑渠道" : "添加通知渠道",
          width: "620px",
          busy: saving.value,
          onClose: _cache[21] || (_cache[21] = ($event) => showEditor.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[20] || (_cache[20] = ($event) => showEditor.value = false)
            }, "取消"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: saveChannel
            }, toDisplayString(editingChannel.value?.id ? "保存修改" : "创建渠道"), 1)
          ]),
          default: withCtx(() => [
            editingChannel.value ? (openBlock(), createElementBlock("div", _hoisted_55, [
              createBaseVNode("div", _hoisted_56, [
                createBaseVNode("div", null, [
                  _cache[53] || (_cache[53] = createBaseVNode("label", { class: "dh-label" }, "渠道名称", -1)),
                  withDirectives(createBaseVNode("input", {
                    "onUpdate:modelValue": _cache[16] || (_cache[16] = ($event) => editingChannel.value.name = $event),
                    class: "dh-input",
                    placeholder: "例如：家庭群通知"
                  }, null, 512), [
                    [vModelText, editingChannel.value.name]
                  ])
                ]),
                createBaseVNode("div", null, [
                  _cache[54] || (_cache[54] = createBaseVNode("label", { class: "dh-label" }, "类型", -1)),
                  withDirectives(createBaseVNode("select", {
                    "onUpdate:modelValue": _cache[17] || (_cache[17] = ($event) => editingChannel.value.type = $event),
                    class: "dh-select",
                    disabled: !!editingChannel.value.id,
                    onChange: _cache[18] || (_cache[18] = ($event) => editingChannel.value.config = {})
                  }, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(presets.value, (p) => {
                      return openBlock(), createElementBlock("option", {
                        key: p.type,
                        value: p.type
                      }, toDisplayString(p.label), 9, _hoisted_58);
                    }), 128))
                  ], 40, _hoisted_57), [
                    [vModelSelect, editingChannel.value.type]
                  ])
                ])
              ]),
              currentPreset.value ? (openBlock(), createElementBlock("div", _hoisted_59, toDisplayString(currentPreset.value.description), 1)) : createCommentVNode("", true),
              createBaseVNode("div", _hoisted_60, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(currentPreset.value?.fields ?? [], (f) => {
                  return openBlock(), createElementBlock("div", {
                    key: f.key
                  }, [
                    createBaseVNode("label", _hoisted_61, [
                      createTextVNode(toDisplayString(f.label) + " ", 1),
                      f.required ? (openBlock(), createElementBlock("span", _hoisted_62, "*")) : createCommentVNode("", true)
                    ]),
                    f.type === "textarea" ? (openBlock(), createElementBlock("textarea", {
                      key: 0,
                      value: String(editingChannel.value.config[f.key] ?? ""),
                      class: "dh-textarea font-mono text-[11.5px]",
                      placeholder: f.placeholder,
                      onInput: ($event) => editingChannel.value.config[f.key] = $event.target.value
                    }, null, 40, _hoisted_63)) : (openBlock(), createElementBlock("input", {
                      key: 1,
                      value: String(editingChannel.value.config[f.key] ?? ""),
                      type: f.type === "password" ? "password" : "text",
                      class: normalizeClass(["dh-input", needField(f.key) ? "!border-line-err" : ""]),
                      placeholder: f.placeholder,
                      onInput: ($event) => editingChannel.value.config[f.key] = $event.target.value
                    }, null, 42, _hoisted_64)),
                    f.help ? (openBlock(), createElementBlock("div", _hoisted_65, toDisplayString(f.help), 1)) : createCommentVNode("", true)
                  ]);
                }), 128))
              ]),
              createBaseVNode("label", _hoisted_66, [
                createVNode(_sfc_main$2, {
                  modelValue: editingChannel.value.enabled,
                  "onUpdate:modelValue": _cache[19] || (_cache[19] = ($event) => editingChannel.value.enabled = $event)
                }, null, 8, ["modelValue"]),
                _cache[55] || (_cache[55] = createTextVNode(" 启用这个渠道 ", -1))
              ])
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        }, 8, ["open", "title", "busy"]),
        createVNode(_sfc_main$3, {
          open: showVars.value,
          title: "模板变量",
          width: "560px",
          onClose: _cache[23] || (_cache[23] = ($event) => showVars.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: _cache[22] || (_cache[22] = ($event) => showVars.value = false)
            }, "关闭")
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_67, [
              _cache[56] || (_cache[56] = createTextVNode(" 自定义 Webhook 的 URL、请求头、请求体里都可以使用这些变量，写法是 ", -1)),
              createBaseVNode("code", _hoisted_68, toDisplayString(tplVar("变量名")), 1),
              _cache[57] || (_cache[57] = createTextVNode("。URL 中的变量会自动做 URL 转义。 ", -1))
            ]),
            createBaseVNode("div", _hoisted_69, [
              createBaseVNode("table", _hoisted_70, [
                _cache[58] || (_cache[58] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", { class: "w-[140px]" }, "变量"),
                    createBaseVNode("th", null, "说明")
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(templateVars.value, (v) => {
                    return openBlock(), createElementBlock("tr", {
                      key: v.key
                    }, [
                      createBaseVNode("td", _hoisted_71, toDisplayString(tplVar(v.key)), 1),
                      createBaseVNode("td", _hoisted_72, toDisplayString(v.label), 1)
                    ]);
                  }), 128))
                ])
              ])
            ])
          ]),
          _: 1
        }, 8, ["open"]),
        createVNode(_sfc_main$3, {
          open: !!removeTarget.value,
          title: "删除渠道",
          subtitle: removeTarget.value?.name,
          width: "400px",
          busy: removing.value,
          onClose: _cache[25] || (_cache[25] = ($event) => removeTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[24] || (_cache[24] = ($event) => removeTarget.value = null)
            }, "取消", 8, _hoisted_73),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemove
            }, toDisplayString(removing.value ? "删除中…" : "确认删除"), 9, _hoisted_74)
          ]),
          default: withCtx(() => [
            _cache[59] || (_cache[59] = createBaseVNode("div", { class: "text-[12.5px] text-text-3" }, "删除后该渠道不再接收任何通知。", -1))
          ]),
          _: 1
        }, 8, ["open", "subtitle", "busy"]),
        createVNode(_sfc_main$3, {
          open: confirmClearHistory.value,
          title: "清空推送历史",
          width: "430px",
          busy: clearingHistory.value,
          onClose: _cache[27] || (_cache[27] = ($event) => confirmClearHistory.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: clearingHistory.value,
              onClick: _cache[26] || (_cache[26] = ($event) => confirmClearHistory.value = false)
            }, "取消", 8, _hoisted_75),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: clearingHistory.value,
              onClick: clearHistory
            }, [
              createVNode(unref(Trash2), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(clearingHistory.value ? "清空中…" : "确认清空"), 1)
            ], 8, _hoisted_76)
          ]),
          default: withCtx(() => [
            _cache[60] || (_cache[60] = createBaseVNode("div", { class: "text-[12.5px] leading-relaxed text-text-3" }, [
              createTextVNode("清空后所有渠道的"),
              createBaseVNode("b", { class: "text-err-text" }, "推送记录"),
              createTextVNode("会立刻消失，不可恢复。已配置的通知渠道与事件订阅不受影响。")
            ], -1))
          ]),
          _: 1
        }, 8, ["open", "busy"])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
