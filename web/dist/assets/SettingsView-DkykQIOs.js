import { x as createLucideIcon, d as defineComponent, E as useToastStore, u as useAppStore, a as onMounted, c as createElementBlock, e as createBaseVNode, f as createVNode, g as unref, R as RefreshCw, t as toDisplayString, i as createBlock, k as createTextVNode, I as withCtx, n as normalizeClass, w as withDirectives, N as vModelSelect, F as Fragment, r as renderList, j as createCommentVNode, _ as LogOut, v as vModelDynamic, m as ref, p as computed, B as api, s as openBlock, G as vModelText } from "./index-DBOlVmiV.js";
import { d as formatDateTime, b as relativeTime, r as runKindLabel } from "./format-NZu8lNOP.js";
import { T as Trash2, _ as _sfc_main$2 } from "./Modal.vue_vue_type_script_setup_true_lang-CSO2rwvT.js";
import { _ as _sfc_main$1, a as _sfc_main$3 } from "./ToggleSwitch.vue_vue_type_script_setup_true_lang-CCzcTs2h.js";
import { L as LoaderCircle } from "./loader-circle-CFBmSG65.js";
import { S as Save } from "./save-DsALYqad.js";
import { Z as Zap } from "./zap-C9j9lw7d.js";
import { S as ScrollText } from "./scroll-text-3D8KZ7I7.js";
import { E as EyeOff } from "./eye-off-snJIpsfm.js";
import { E as Eye } from "./eye-33ATqYZj.js";
const Ban = createLucideIcon("BanIcon", [
  ["circle", { cx: "12", cy: "12", r: "10", key: "1mglay" }],
  ["path", { d: "m4.9 4.9 14.2 14.2", key: "1m5liu" }]
]);
const KeyRound = createLucideIcon("KeyRoundIcon", [
  [
    "path",
    {
      d: "M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z",
      key: "1s6t7t"
    }
  ],
  ["circle", { cx: "16.5", cy: "7.5", r: ".5", fill: "currentColor", key: "w0ekpg" }]
]);
const SearchCheck = createLucideIcon("SearchCheckIcon", [
  ["path", { d: "m8 11 2 2 4-4", key: "1sed1v" }],
  ["circle", { cx: "11", cy: "11", r: "8", key: "4ej97u" }],
  ["path", { d: "m21 21-4.3-4.3", key: "1qie3q" }]
]);
const SlidersHorizontal = createLucideIcon("SlidersHorizontalIcon", [
  ["line", { x1: "21", x2: "14", y1: "4", y2: "4", key: "obuewd" }],
  ["line", { x1: "10", x2: "3", y1: "4", y2: "4", key: "1q6298" }],
  ["line", { x1: "21", x2: "12", y1: "12", y2: "12", key: "1iu8h1" }],
  ["line", { x1: "8", x2: "3", y1: "12", y2: "12", key: "ntss68" }],
  ["line", { x1: "21", x2: "16", y1: "20", y2: "20", key: "14d8ph" }],
  ["line", { x1: "12", x2: "3", y1: "20", y2: "20", key: "m0wm8r" }],
  ["line", { x1: "14", x2: "14", y1: "2", y2: "6", key: "14e1ph" }],
  ["line", { x1: "8", x2: "8", y1: "10", y2: "14", key: "1i6ji0" }],
  ["line", { x1: "16", x2: "16", y1: "18", y2: "22", key: "1lctlv" }]
]);
const User = createLucideIcon("UserIcon", [
  ["path", { d: "M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2", key: "975kel" }],
  ["circle", { cx: "12", cy: "7", r: "4", key: "17ys0d" }]
]);
const Users = createLucideIcon("UsersIcon", [
  ["path", { d: "M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2", key: "1yyitq" }],
  ["circle", { cx: "9", cy: "7", r: "4", key: "nufk8" }],
  ["path", { d: "M22 21v-2a4 4 0 0 0-3-3.87", key: "kshegd" }],
  ["path", { d: "M16 3.13a4 4 0 0 1 0 7.75", key: "1da9ce" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-card" };
const _hoisted_3 = { class: "dh-card-head" };
const _hoisted_4 = { class: "ml-2 text-[11.5px] font-normal text-text-5" };
const _hoisted_5 = ["disabled"];
const _hoisted_6 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_7 = { class: "rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_8 = { class: "mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3" };
const _hoisted_9 = { class: "flex flex-col gap-2.5" };
const _hoisted_10 = ["value"];
const _hoisted_11 = { class: "mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3" };
const _hoisted_12 = { class: "flex flex-col gap-2.5" };
const _hoisted_13 = { class: "rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_14 = { class: "mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3" };
const _hoisted_15 = { class: "ml-auto text-[11px] font-normal text-text-5" };
const _hoisted_16 = { class: "flex flex-col gap-2.5" };
const _hoisted_17 = { class: "flex flex-wrap gap-2" };
const _hoisted_18 = ["value"];
const _hoisted_19 = ["disabled"];
const _hoisted_20 = {
  key: 0,
  class: "text-[12px] text-text-5"
};
const _hoisted_21 = {
  key: 1,
  class: "text-[12px] text-text-5"
};
const _hoisted_22 = {
  key: 2,
  class: "flex flex-wrap gap-1.5"
};
const _hoisted_23 = ["onClick"];
const _hoisted_24 = { class: "rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_25 = { class: "mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3" };
const _hoisted_26 = { class: "flex flex-col gap-2.5" };
const _hoisted_27 = { class: "rounded-[10px] border border-line-1 bg-ink-800 p-3" };
const _hoisted_28 = { class: "mb-2.5 flex items-center gap-2 text-[12px] font-medium text-text-3" };
const _hoisted_29 = { class: "flex flex-col gap-2.5" };
const _hoisted_30 = { class: "grid grid-cols-1 gap-3.5 xl:grid-cols-2" };
const _hoisted_31 = { class: "dh-card" };
const _hoisted_32 = { class: "dh-card-head" };
const _hoisted_33 = ["disabled"];
const _hoisted_34 = { class: "flex flex-col gap-3.5 p-3.5" };
const _hoisted_35 = { class: "grid grid-cols-2 gap-3 text-[12px]" };
const _hoisted_36 = { class: "text-text-2" };
const _hoisted_37 = { class: "text-text-2" };
const _hoisted_38 = { class: "text-text-2" };
const _hoisted_39 = { class: "border-t border-line-1 pt-3.5" };
const _hoisted_40 = { class: "mb-2.5 flex items-center gap-2 text-[12.5px] font-medium" };
const _hoisted_41 = { class: "flex flex-col gap-[11px]" };
const _hoisted_42 = { class: "dh-field" };
const _hoisted_43 = ["type"];
const _hoisted_44 = ["title", "aria-label"];
const _hoisted_45 = { class: "dh-field" };
const _hoisted_46 = ["type", "placeholder"];
const _hoisted_47 = ["title", "aria-label"];
const _hoisted_48 = { class: "dh-field" };
const _hoisted_49 = ["type"];
const _hoisted_50 = {
  key: 0,
  class: "text-[11.5px] text-err-text"
};
const _hoisted_51 = { class: "flex justify-end" };
const _hoisted_52 = ["disabled"];
const _hoisted_53 = { class: "border-t border-line-1 pt-3.5" };
const _hoisted_54 = {
  key: 0,
  class: "text-[11.5px] text-text-5"
};
const _hoisted_55 = {
  key: 1,
  class: "flex flex-col gap-1"
};
const _hoisted_56 = { class: "w-[86px] flex-none text-text-5" };
const _hoisted_57 = { class: "font-mono text-text-3" };
const _hoisted_58 = { class: "ml-auto text-text-5" };
const _hoisted_59 = { class: "dh-card" };
const _hoisted_60 = { class: "dh-card-head" };
const _hoisted_61 = { class: "ml-auto flex gap-2" };
const _hoisted_62 = ["disabled"];
const _hoisted_63 = ["disabled"];
const _hoisted_64 = {
  key: 0,
  class: "dh-card-body text-[12.5px] text-text-4"
};
const _hoisted_65 = {
  key: 1,
  class: "dh-scroll max-h-[560px] overflow-auto"
};
const _hoisted_66 = { class: "dh-table" };
const _hoisted_67 = { class: "whitespace-nowrap text-[11.5px] text-text-5" };
const _hoisted_68 = { class: "text-[12px] text-text-2" };
const _hoisted_69 = {
  key: 0,
  class: "font-mono text-[10.5px] text-text-6"
};
const _hoisted_70 = ["disabled"];
const _hoisted_71 = ["disabled"];
const _hoisted_72 = { class: "text-[12.5px] leading-relaxed text-text-3" };
const _hoisted_73 = { class: "text-text-1" };
const _hoisted_74 = ["disabled"];
const _hoisted_75 = ["disabled"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "SettingsView",
  setup(__props) {
    const toast = useToastStore();
    const app = useAppStore();
    const settings = ref({
      exclude: [],
      panelURL: "",
      concurrency: 2,
      logRetention: 100,
      checkOnStart: true,
      checkIntervalHours: 1,
      notifyOnCheck: false,
      autoApply: false,
      pullOnce: true,
      backupBefore: true,
      cleanupAfter: true,
      directFirst: true,
      backupKeepPerContainer: 10,
      backupMaxAgeDays: 30,
      backupMaxTotalMB: 2048,
      backupKeepPreUpdate: true
    });
    const containers = ref([]);
    const logs = ref([]);
    const account = ref(null);
    const auto = ref(null);
    const loading = ref(true);
    const saving = ref(false);
    const newExclude = ref("");
    const oldPw = ref("");
    const newPw = ref("");
    const confirmPw = ref("");
    const showOld = ref(false);
    const showNew = ref(false);
    const changing = ref(false);
    const confirmClearLogs = ref(false);
    const confirmRevoke = ref(false);
    const revoking = ref(false);
    async function revokeOthers() {
      if (revoking.value) return;
      revoking.value = true;
      confirmRevoke.value = false;
      try {
        const r = await api.post("/api/account/sessions/logout-others", {});
        await load();
        toast.success("已登出其他会话", `共作废 ${r.revoked} 个会话，当前登录保持不变`);
      } catch (e) {
        toast.error("操作失败", e instanceof Error ? e.message : String(e));
      } finally {
        revoking.value = false;
      }
    }
    const clearingLogs = ref(false);
    const recentLogins = computed(
      () => account.value?.recentLogins ?? []
    );
    const pwError = computed(() => {
      if (!newPw.value) return "";
      if (newPw.value.length < app.minPasswordLength) return `新密码至少 ${app.minPasswordLength} 位`;
      if (newPw.value !== confirmPw.value) return "两次输入的新密码不一致";
      return "";
    });
    async function load() {
      loading.value = true;
      try {
        const [s, c, l, a, au] = await Promise.all([
          api.get("/api/settings"),
          api.get("/api/containers"),
          api.get("/api/logs", { limit: 120 }),
          api.get("/api/account"),
          // 自动更新的运行状态（下次巡检时间）—— 原来只在已删除的「更新中心」页展示，
          // 那个信息对「多久检测一次」有用，挪到本页检测频率那行。
          api.get("/api/updates/auto").catch(() => null)
        ]);
        settings.value = s;
        savedSnapshot.value = JSON.stringify(patchBody(s));
        containers.value = c.containers ?? [];
        logs.value = l.logs ?? [];
        account.value = a;
        auto.value = au;
      } catch (e) {
        toast.error("读取设置失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function loadAuto() {
      auto.value = await api.get("/api/updates/auto").catch(() => null);
    }
    const checkIntervalSub = computed(() => {
      const base = `每 ${settings.value.checkIntervalHours} 小时自动扫描一次镜像仓库`;
      return auto.value?.nextCheckAt ? `${base} · 下次巡检 ${formatDateTime(auto.value.nextCheckAt)}` : base;
    });
    const intervalOptions = computed(() => {
      const list = auto.value?.intervalChoices;
      return list && list.length ? list : [settings.value.checkIntervalHours];
    });
    const FIELD_LABEL = {
      exclude: "排除列表",
      panelURL: "容器面板地址",
      concurrency: "并发度",
      logRetention: "运行记录保留条数",
      checkOnStart: "启动自动巡检",
      checkIntervalHours: "检测频率",
      notifyOnCheck: "检测完成后通知",
      autoApply: "自动更新",
      pullOnce: "同一镜像只拉取一次",
      backupBefore: "更新前备份",
      cleanupAfter: "更新后清理"
    };
    function patchBody(s) {
      return {
        exclude: s.exclude,
        panelURL: s.panelURL,
        concurrency: s.concurrency,
        logRetention: s.logRetention,
        checkOnStart: s.checkOnStart,
        checkIntervalHours: s.checkIntervalHours,
        notifyOnCheck: s.notifyOnCheck,
        autoApply: s.autoApply,
        pullOnce: s.pullOnce,
        backupBefore: s.backupBefore,
        cleanupAfter: s.cleanupAfter
      };
    }
    const savedSnapshot = ref("");
    async function save() {
      const body = patchBody(settings.value);
      const snapshot = JSON.stringify(body);
      if (snapshot === savedSnapshot.value) {
        toast.info("没有改动", "设置与上次保存时一致");
        return;
      }
      saving.value = true;
      try {
        settings.value = await api.patch("/api/settings", body);
        await app.loadSettings();
        void loadAuto();
        const prev = JSON.parse(savedSnapshot.value || "{}");
        const changed = Object.keys(body).filter((k) => JSON.stringify(body[k]) !== JSON.stringify(prev[k]));
        const names = changed.map((k) => FIELD_LABEL[k] ?? k);
        savedSnapshot.value = snapshot;
        toast.success("设置已保存", names.length ? `已更新：${names.join("、")}` : "设置已更新");
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      } finally {
        saving.value = false;
      }
    }
    function addExclude() {
      const n = newExclude.value.trim();
      if (!n || settings.value.exclude.includes(n)) return;
      settings.value.exclude.push(n);
      newExclude.value = "";
    }
    function removeExclude(name) {
      settings.value.exclude = settings.value.exclude.filter((x) => x !== name);
    }
    async function changePassword() {
      if (pwError.value || !oldPw.value || !newPw.value) return;
      changing.value = true;
      try {
        await api.post("/api/account/password", {
          oldPassword: oldPw.value,
          newPassword: newPw.value,
          confirm: confirmPw.value
        });
        toast.success("密码已修改", "所有设备上的登录会话都已失效，请重新登录");
        oldPw.value = "";
        newPw.value = "";
        confirmPw.value = "";
        window.setTimeout(() => void doLogout(), 1200);
      } catch (e) {
        toast.error("修改失败", e instanceof Error ? e.message : String(e));
      } finally {
        changing.value = false;
      }
    }
    async function doLogout() {
      try {
        await app.logout();
      } finally {
        window.location.replace("/login");
      }
    }
    async function clearLogs() {
      if (clearingLogs.value) return;
      clearingLogs.value = true;
      confirmClearLogs.value = false;
      try {
        await api.del("/api/logs");
        logs.value = [];
        toast.success("运行记录已清空");
      } catch (e) {
        toast.error("清空失败", e instanceof Error ? e.message : String(e));
      } finally {
        clearingLogs.value = false;
      }
    }
    onMounted(() => void load());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        _cache[54] || (_cache[54] = createBaseVNode("div", { class: "dh-phead" }, [
          createBaseVNode("div", { class: "dh-h1" }, "设置"),
          createBaseVNode("div", { class: "dh-sub" }, "更新与检测（含排除列表）与账户安全")
        ], -1)),
        createBaseVNode("div", _hoisted_2, [
          createBaseVNode("div", _hoisted_3, [
            createVNode(unref(RefreshCw), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[23] || (_cache[23] = createBaseVNode("span", null, "更新与检测", -1)),
            createBaseVNode("span", _hoisted_4, toDisplayString(`每 ${settings.value.checkIntervalHours} 小时巡检一次`) + " · 自动更新" + toDisplayString(settings.value.autoApply ? "已开启" : "已关闭"), 1),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm dh-btn-primary ml-auto",
              disabled: saving.value,
              onClick: save
            }, [
              saving.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3 w-3 dh-spin"
              })) : (openBlock(), createBlock(unref(Save), {
                key: 1,
                class: "h-3 w-3"
              })),
              _cache[22] || (_cache[22] = createTextVNode("保存 ", -1))
            ], 8, _hoisted_5)
          ]),
          createBaseVNode("div", _hoisted_6, [
            createBaseVNode("div", _hoisted_7, [
              createBaseVNode("div", _hoisted_8, [
                createVNode(unref(SearchCheck), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[24] || (_cache[24] = createTextVNode("检测 ", -1))
              ]),
              createBaseVNode("div", _hoisted_9, [
                createVNode(_sfc_main$1, {
                  title: "检测频率",
                  sub: checkIntervalSub.value
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("select", {
                      "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => settings.value.checkIntervalHours = $event),
                      class: "dh-select !w-[120px] !py-[5px] !text-[11.5px]"
                    }, [
                      (openBlock(true), createElementBlock(Fragment, null, renderList(intervalOptions.value, (h) => {
                        return openBlock(), createElementBlock("option", {
                          key: h,
                          value: h
                        }, "每 " + toDisplayString(h) + " 小时", 9, _hoisted_10);
                      }), 128))
                    ], 512), [
                      [
                        vModelSelect,
                        settings.value.checkIntervalHours,
                        void 0,
                        { number: true }
                      ]
                    ])
                  ]),
                  _: 1
                }, 8, ["sub"]),
                createVNode(_sfc_main$1, {
                  title: "检测完成后通知",
                  sub: "巡检发现问题时推一条通知到已配置的渠道"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$3, {
                      modelValue: settings.value.notifyOnCheck,
                      "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => settings.value.notifyOnCheck = $event),
                      label: "检测完成后通知"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$1, {
                  title: "启动后自动巡检一次",
                  sub: "延迟 8 秒执行，只读检查，不会停止任何容器"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$3, {
                      modelValue: settings.value.checkOnStart,
                      "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => settings.value.checkOnStart = $event),
                      label: "启动后自动巡检一次"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                })
              ])
            ]),
            createBaseVNode("div", {
              class: normalizeClass(["rounded-[10px] border p-3", settings.value.autoApply ? "border-line-warn bg-soft-warn" : "border-line-1 bg-ink-800"])
            }, [
              createBaseVNode("div", _hoisted_11, [
                createVNode(unref(Zap), {
                  class: normalizeClass(["h-3.5 w-3.5", settings.value.autoApply ? "text-warn-text" : "text-text-4"])
                }, null, 8, ["class"]),
                _cache[25] || (_cache[25] = createTextVNode("自动更新 ", -1))
              ]),
              createBaseVNode("div", _hoisted_12, [
                createVNode(_sfc_main$1, {
                  title: "检测到新版本后自动更新",
                  sub: settings.value.autoApply ? "每轮巡检结束后，把有更新的容器（排除列表与自己除外）自动重建到新镜像；重启后的第一轮排在启动 15 分钟后，不必等满一个检测周期" : "当前只检测、不动手 —— 发现更新后去容器页手动更新"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$3, {
                      modelValue: settings.value.autoApply,
                      "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => settings.value.autoApply = $event),
                      label: "自动更新"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                }, 8, ["sub"])
              ])
            ], 2),
            createBaseVNode("div", _hoisted_13, [
              createBaseVNode("div", _hoisted_14, [
                createVNode(unref(Ban), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[26] || (_cache[26] = createTextVNode("排除列表 ", -1)),
                createBaseVNode("span", _hoisted_15, toDisplayString(settings.value.exclude.length) + " 个容器永不自动更新 ", 1)
              ]),
              createBaseVNode("div", _hoisted_16, [
                createBaseVNode("div", _hoisted_17, [
                  withDirectives(createBaseVNode("select", {
                    "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => newExclude.value = $event),
                    class: "dh-select !w-auto !min-w-[190px]"
                  }, [
                    _cache[27] || (_cache[27] = createBaseVNode("option", { value: "" }, "选择容器…", -1)),
                    (openBlock(true), createElementBlock(Fragment, null, renderList(containers.value.filter((x) => !settings.value.exclude.includes(x.name)), (c) => {
                      return openBlock(), createElementBlock("option", {
                        key: c.id,
                        value: c.name
                      }, toDisplayString(c.name), 9, _hoisted_18);
                    }), 128))
                  ], 512), [
                    [vModelSelect, newExclude.value]
                  ]),
                  createBaseVNode("button", {
                    class: "dh-btn",
                    disabled: !newExclude.value,
                    onClick: addExclude
                  }, "加入排除", 8, _hoisted_19)
                ]),
                unref(app).settings === null && loading.value ? (openBlock(), createElementBlock("div", _hoisted_20, "正在载入…")) : createCommentVNode("", true),
                !settings.value.exclude.length ? (openBlock(), createElementBlock("div", _hoisted_21, "排除列表为空。")) : (openBlock(), createElementBlock("div", _hoisted_22, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(settings.value.exclude, (n) => {
                    return openBlock(), createElementBlock("span", {
                      key: n,
                      class: "inline-flex items-center gap-1.5 rounded-full border border-line-3 px-2.5 py-[3px] font-mono text-[11.5px] text-text-3"
                    }, [
                      createTextVNode(toDisplayString(n) + " ", 1),
                      createBaseVNode("button", {
                        class: "dh-tap text-text-5 hover:text-err-text",
                        onClick: ($event) => removeExclude(n)
                      }, "×", 8, _hoisted_23)
                    ]);
                  }), 128))
                ])),
                _cache[28] || (_cache[28] = createBaseVNode("div", { class: "text-[11px] text-text-5" }, [
                  createTextVNode(" 加进这里的容器不会被自动更新，也不会被任何计划任务作用到（计划任务必须逐个勾选容器）。 改动和上面的开关一样，"),
                  createBaseVNode("b", { class: "text-text-4" }, "点右上角「保存」后才生效"),
                  createTextVNode("。 ")
                ], -1))
              ])
            ]),
            createBaseVNode("div", _hoisted_24, [
              createBaseVNode("div", _hoisted_25, [
                createVNode(unref(SlidersHorizontal), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[29] || (_cache[29] = createTextVNode("执行策略 ", -1))
              ]),
              createBaseVNode("div", _hoisted_26, [
                createVNode(_sfc_main$1, {
                  title: "同一镜像的多个容器只拉取一次",
                  sub: "4 个容器共用 nginx:alpine 时，下载 1 次、重建 4 个"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$3, {
                      modelValue: settings.value.pullOnce,
                      "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => settings.value.pullOnce = $event),
                      label: "同一镜像只拉取一次"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$1, {
                  title: "更新前自动备份容器配置",
                  sub: "失败可一键回滚到更新前"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$3, {
                      modelValue: settings.value.backupBefore,
                      "onUpdate:modelValue": _cache[6] || (_cache[6] = ($event) => settings.value.backupBefore = $event),
                      label: "更新前自动备份容器配置"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$1, {
                  title: "更新后清理旧镜像",
                  sub: "确认没有任何容器再引用后才删除"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$3, {
                      modelValue: settings.value.cleanupAfter,
                      "onUpdate:modelValue": _cache[7] || (_cache[7] = ($event) => settings.value.cleanupAfter = $event),
                      label: "更新后清理旧镜像"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$1, {
                  title: "批量更新并发度",
                  sub: "并发越高越快，但更容易触发镜像仓库限流"
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("select", {
                      "onUpdate:modelValue": _cache[8] || (_cache[8] = ($event) => settings.value.concurrency = $event),
                      class: "dh-select !w-[150px] !py-[5px] !text-[11.5px]"
                    }, [..._cache[30] || (_cache[30] = [
                      createBaseVNode("option", { value: 1 }, "1（串行，最稳）", -1),
                      createBaseVNode("option", { value: 2 }, "2（推荐）", -1),
                      createBaseVNode("option", { value: 3 }, "3", -1),
                      createBaseVNode("option", { value: 4 }, "4", -1),
                      createBaseVNode("option", { value: 6 }, "6（激进）", -1)
                    ])], 512), [
                      [
                        vModelSelect,
                        settings.value.concurrency,
                        void 0,
                        { number: true }
                      ]
                    ])
                  ]),
                  _: 1
                })
              ])
            ]),
            createBaseVNode("div", _hoisted_27, [
              createBaseVNode("div", _hoisted_28, [
                createVNode(unref(ScrollText), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[31] || (_cache[31] = createTextVNode("面板与记录 ", -1))
              ]),
              createBaseVNode("div", _hoisted_29, [
                createVNode(_sfc_main$1, {
                  stack: "",
                  title: "面板地址",
                  sub: "通知模板里的 {{url}} 用它拼可点击的链接"
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("input", {
                      "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => settings.value.panelURL = $event),
                      class: "dh-input !w-[260px] max-md:!w-full",
                      placeholder: "http://192.168.1.10:5923"
                    }, null, 512), [
                      [vModelText, settings.value.panelURL]
                    ])
                  ]),
                  _: 1
                }),
                createVNode(_sfc_main$1, {
                  title: "运行记录保留条数",
                  sub: "超出后按时间滚动覆盖，只影响面板里的历史列表"
                }, {
                  default: withCtx(() => [
                    withDirectives(createBaseVNode("input", {
                      "onUpdate:modelValue": _cache[10] || (_cache[10] = ($event) => settings.value.logRetention = $event),
                      type: "number",
                      min: "50",
                      class: "dh-input !w-[110px]"
                    }, null, 512), [
                      [
                        vModelText,
                        settings.value.logRetention,
                        void 0,
                        { number: true }
                      ]
                    ])
                  ]),
                  _: 1
                })
              ])
            ])
          ])
        ]),
        createBaseVNode("div", _hoisted_30, [
          createBaseVNode("div", _hoisted_31, [
            createBaseVNode("div", _hoisted_32, [
              createVNode(unref(User), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[32] || (_cache[32] = createBaseVNode("span", null, "账户", -1)),
              (account.value?.sessionCount ?? 0) > 1 ? (openBlock(), createElementBlock("button", {
                key: 0,
                class: "dh-btn dh-btn-sm ml-auto font-normal",
                disabled: revoking.value,
                onClick: _cache[11] || (_cache[11] = ($event) => confirmRevoke.value = true)
              }, [
                createVNode(unref(LogOut), { class: "h-3 w-3" }),
                createTextVNode(" 登出其他会话（" + toDisplayString((account.value?.sessionCount ?? 0) - 1) + "） ", 1)
              ], 8, _hoisted_33)) : createCommentVNode("", true)
            ]),
            createBaseVNode("div", _hoisted_34, [
              createBaseVNode("div", _hoisted_35, [
                _cache[33] || (_cache[33] = createBaseVNode("div", { class: "text-text-5" }, "当前会话数", -1)),
                createBaseVNode("div", _hoisted_36, [
                  createVNode(unref(Users), { class: "mr-1 inline h-3 w-3 text-text-5" }),
                  createTextVNode(toDisplayString(account.value?.sessionCount ?? "—"), 1)
                ]),
                _cache[34] || (_cache[34] = createBaseVNode("div", { class: "text-text-5" }, "当前会话创建于", -1)),
                createBaseVNode("div", _hoisted_37, toDisplayString(account.value?.currentSession ? unref(formatDateTime)(account.value.currentSession.createdAt) : "—"), 1),
                _cache[35] || (_cache[35] = createBaseVNode("div", { class: "text-text-5" }, "当前会话过期于", -1)),
                createBaseVNode("div", _hoisted_38, toDisplayString(account.value?.currentSession ? unref(formatDateTime)(account.value.currentSession.expiresAt) : "—"), 1)
              ]),
              createBaseVNode("div", _hoisted_39, [
                createBaseVNode("div", _hoisted_40, [
                  createVNode(unref(KeyRound), { class: "h-3.5 w-3.5 text-text-4" }),
                  _cache[36] || (_cache[36] = createTextVNode("修改登录密码 ", -1))
                ]),
                createBaseVNode("div", _hoisted_41, [
                  createBaseVNode("div", null, [
                    _cache[37] || (_cache[37] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "当前密码", -1)),
                    createBaseVNode("div", _hoisted_42, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[12] || (_cache[12] = ($event) => oldPw.value = $event),
                        type: showOld.value ? "text" : "password",
                        placeholder: "当前使用的密码",
                        autocomplete: "current-password"
                      }, null, 8, _hoisted_43), [
                        [vModelDynamic, oldPw.value]
                      ]),
                      createBaseVNode("button", {
                        type: "button",
                        class: "dh-eye",
                        title: showOld.value ? "隐藏密码" : "显示密码",
                        "aria-label": showOld.value ? "隐藏密码" : "显示密码",
                        onClick: _cache[13] || (_cache[13] = ($event) => showOld.value = !showOld.value)
                      }, [
                        showOld.value ? (openBlock(), createBlock(unref(EyeOff), {
                          key: 0,
                          class: "h-[15px] w-[15px]"
                        })) : (openBlock(), createBlock(unref(Eye), {
                          key: 1,
                          class: "h-[15px] w-[15px]"
                        }))
                      ], 8, _hoisted_44)
                    ])
                  ]),
                  createBaseVNode("div", null, [
                    _cache[38] || (_cache[38] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "新密码", -1)),
                    createBaseVNode("div", _hoisted_45, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[14] || (_cache[14] = ($event) => newPw.value = $event),
                        type: showNew.value ? "text" : "password",
                        placeholder: `至少 ${unref(app).minPasswordLength} 位`,
                        autocomplete: "new-password"
                      }, null, 8, _hoisted_46), [
                        [vModelDynamic, newPw.value]
                      ]),
                      createBaseVNode("button", {
                        type: "button",
                        class: "dh-eye",
                        title: showNew.value ? "隐藏密码" : "显示密码",
                        "aria-label": showNew.value ? "隐藏密码" : "显示密码",
                        onClick: _cache[15] || (_cache[15] = ($event) => showNew.value = !showNew.value)
                      }, [
                        showNew.value ? (openBlock(), createBlock(unref(EyeOff), {
                          key: 0,
                          class: "h-[15px] w-[15px]"
                        })) : (openBlock(), createBlock(unref(Eye), {
                          key: 1,
                          class: "h-[15px] w-[15px]"
                        }))
                      ], 8, _hoisted_47)
                    ])
                  ]),
                  createBaseVNode("div", null, [
                    _cache[39] || (_cache[39] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "确认新密码", -1)),
                    createBaseVNode("div", _hoisted_48, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[16] || (_cache[16] = ($event) => confirmPw.value = $event),
                        type: showNew.value ? "text" : "password",
                        placeholder: "再输一次新密码",
                        autocomplete: "new-password"
                      }, null, 8, _hoisted_49), [
                        [vModelDynamic, confirmPw.value]
                      ])
                    ])
                  ]),
                  pwError.value ? (openBlock(), createElementBlock("div", _hoisted_50, toDisplayString(pwError.value), 1)) : createCommentVNode("", true),
                  _cache[41] || (_cache[41] = createBaseVNode("div", { class: "dh-banner dh-banner-info !gap-2.5 !py-[9px] !pl-3 !pr-3 !text-[12px]" }, [
                    createBaseVNode("span", { class: "h-[13px] w-[13px] flex-none rounded-[4px] bg-current opacity-50" }),
                    createBaseVNode("span", { class: "min-w-0 flex-1" }, "改密成功后其他设备上的登录会立即失效，需重新登录。")
                  ], -1)),
                  createBaseVNode("div", _hoisted_51, [
                    createBaseVNode("button", {
                      class: "dh-btn dh-btn-primary",
                      disabled: changing.value || !oldPw.value || !!pwError.value || !newPw.value,
                      onClick: changePassword
                    }, [
                      changing.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                        key: 0,
                        class: "h-3.5 w-3.5 dh-spin"
                      })) : (openBlock(), createBlock(unref(KeyRound), {
                        key: 1,
                        class: "h-3.5 w-3.5"
                      })),
                      _cache[40] || (_cache[40] = createTextVNode(" 保存新密码 ", -1))
                    ], 8, _hoisted_52)
                  ]),
                  _cache[42] || (_cache[42] = createBaseVNode("div", { class: "text-[11px] leading-relaxed text-text-5" }, [
                    createTextVNode(" 忘记密码时：在 NAS 面板里删除 "),
                    createBaseVNode("code", { class: "text-text-3" }, "data/auth.json"),
                    createTextVNode(" 后重启容器， 或在 compose 里临时加 "),
                    createBaseVNode("code", { class: "text-text-3" }, "DOCKHELM_PASSWORD=新密码"),
                    createTextVNode("。 ")
                  ], -1))
                ])
              ]),
              createBaseVNode("div", _hoisted_53, [
                _cache[43] || (_cache[43] = createBaseVNode("div", { class: "mb-2 text-[12.5px] font-medium" }, "最近登录记录", -1)),
                !recentLogins.value.length ? (openBlock(), createElementBlock("div", _hoisted_54, "暂无记录")) : (openBlock(), createElementBlock("div", _hoisted_55, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(recentLogins.value, (r, i) => {
                    return openBlock(), createElementBlock("div", {
                      key: i,
                      class: "flex items-center gap-2 text-[11.5px]"
                    }, [
                      createBaseVNode("span", {
                        class: normalizeClass(["h-[6px] w-[6px] flex-none rounded-full", r.ok ? "bg-run" : "bg-err"])
                      }, null, 2),
                      createBaseVNode("span", _hoisted_56, toDisplayString(unref(relativeTime)(r.ts)), 1),
                      createBaseVNode("span", _hoisted_57, toDisplayString(r.ip), 1),
                      createBaseVNode("span", _hoisted_58, toDisplayString(r.ok ? "成功" : "失败"), 1)
                    ]);
                  }), 128))
                ]))
              ]),
              createBaseVNode("button", {
                class: "dh-btn",
                onClick: doLogout
              }, [
                createVNode(unref(LogOut), { class: "h-3.5 w-3.5" }),
                _cache[44] || (_cache[44] = createTextVNode("退出登录 ", -1))
              ])
            ])
          ]),
          createBaseVNode("div", _hoisted_59, [
            createBaseVNode("div", _hoisted_60, [
              createVNode(unref(ScrollText), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[46] || (_cache[46] = createBaseVNode("span", null, "运行记录", -1)),
              createBaseVNode("div", _hoisted_61, [
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm",
                  disabled: loading.value,
                  onClick: load
                }, [
                  createVNode(unref(RefreshCw), {
                    class: normalizeClass(["h-3 w-3", loading.value ? "dh-spin" : ""])
                  }, null, 8, ["class"])
                ], 8, _hoisted_62),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm dh-btn-danger",
                  disabled: !logs.value.length || clearingLogs.value,
                  onClick: _cache[17] || (_cache[17] = ($event) => confirmClearLogs.value = true)
                }, [
                  createVNode(unref(Trash2), { class: "h-3 w-3" }),
                  _cache[45] || (_cache[45] = createTextVNode("清空 ", -1))
                ], 8, _hoisted_63)
              ])
            ]),
            !logs.value.length ? (openBlock(), createElementBlock("div", _hoisted_64, "暂无记录")) : (openBlock(), createElementBlock("div", _hoisted_65, [
              createBaseVNode("table", _hoisted_66, [
                _cache[47] || (_cache[47] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", { class: "w-[120px]" }, "时间"),
                    createBaseVNode("th", { class: "w-[86px]" }, "类型"),
                    createBaseVNode("th", null, "说明")
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(logs.value, (l) => {
                    return openBlock(), createElementBlock("tr", {
                      key: l.id
                    }, [
                      createBaseVNode("td", _hoisted_67, toDisplayString(unref(relativeTime)(l.ts)), 1),
                      createBaseVNode("td", null, [
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge", l.status === "failed" ? "dh-badge-err" : l.status === "success" || l.status === "up_to_date" ? "dh-badge-accent" : "dh-badge-plain"])
                        }, toDisplayString(unref(runKindLabel)(l.kind)), 3)
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("div", _hoisted_68, toDisplayString(l.message), 1),
                        l.ref ? (openBlock(), createElementBlock("div", _hoisted_69, toDisplayString(l.ref), 1)) : createCommentVNode("", true)
                      ])
                    ]);
                  }), 128))
                ])
              ])
            ]))
          ])
        ]),
        createVNode(_sfc_main$2, {
          open: confirmClearLogs.value,
          title: "清空运行记录",
          width: "430px",
          busy: clearingLogs.value,
          onClose: _cache[19] || (_cache[19] = ($event) => confirmClearLogs.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: clearingLogs.value,
              onClick: _cache[18] || (_cache[18] = ($event) => confirmClearLogs.value = false)
            }, "取消", 8, _hoisted_70),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: clearingLogs.value,
              onClick: clearLogs
            }, [
              createVNode(unref(Trash2), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(clearingLogs.value ? "清空中…" : "确认清空"), 1)
            ], 8, _hoisted_71)
          ]),
          default: withCtx(() => [
            _cache[48] || (_cache[48] = createBaseVNode("div", { class: "text-[12.5px] leading-relaxed text-text-3" }, [
              createTextVNode("清空后全部更新 / 巡检 / 备份的"),
              createBaseVNode("b", { class: "text-err-text" }, "执行记录"),
              createTextVNode("会立刻消失，不可恢复。容器、镜像与快照本身都不受影响。")
            ], -1))
          ]),
          _: 1
        }, 8, ["open", "busy"]),
        createVNode(_sfc_main$2, {
          open: confirmRevoke.value,
          title: "强制登出其他会话",
          width: "430px",
          busy: revoking.value,
          onClose: _cache[21] || (_cache[21] = ($event) => confirmRevoke.value = false)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: revoking.value,
              onClick: _cache[20] || (_cache[20] = ($event) => confirmRevoke.value = false)
            }, "取消", 8, _hoisted_74),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: revoking.value,
              onClick: revokeOthers
            }, [
              revoking.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                key: 0,
                class: "h-3.5 w-3.5 dh-spin"
              })) : (openBlock(), createBlock(unref(LogOut), {
                key: 1,
                class: "h-3.5 w-3.5"
              })),
              _cache[53] || (_cache[53] = createTextVNode(" 确认登出 ", -1))
            ], 8, _hoisted_75)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_72, [
              _cache[49] || (_cache[49] = createTextVNode(" 其他浏览器 / 设备上的 ", -1)),
              createBaseVNode("b", _hoisted_73, toDisplayString((account.value?.sessionCount ?? 0) - 1) + " 个会话", 1),
              _cache[50] || (_cache[50] = createTextVNode("会立刻失效， 那边再打开面板时需要重新输密码。", -1)),
              _cache[51] || (_cache[51] = createBaseVNode("b", { class: "text-text-1" }, "你当前这个登录不受影响", -1)),
              _cache[52] || (_cache[52] = createTextVNode("，也不会改动密码本身。 ", -1))
            ])
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
