import { c as createLucideIcon, d as defineComponent, u as useAppStore, o as onUnmounted, a as onMounted, b as openBlock, e as createElementBlock, f as createBaseVNode, g as createStaticVNode, h as createTextVNode, F as Fragment, r as renderList, w as withModifiers, i as createVNode, j as unref, t as toDisplayString, k as withDirectives, v as vModelDynamic, l as createBlock, m as createCommentVNode, n as normalizeClass, p as vModelCheckbox, q as ref, s as computed, x as useRouter, y as useRoute, z as resolveDynamicComponent, L as Layers, R as RefreshCw, C as CalendarClock, A as ListOrdered, B as Bell, D as nextTick, E as ApiError, G as safeRedirect } from "./index-BhImyB54.js";
import { v as version, B as Boxes } from "./config-D4g4wMyX.js";
import { L as LoaderCircle } from "./loader-circle-gCbFc6UW.js";
import { E as EyeOff } from "./eye-off-4_9PmBRg.js";
import { E as Eye } from "./eye-y0z_kbHn.js";
import { T as TriangleAlert } from "./triangle-alert-OoinH6Fb.js";
import { A as ArrowRight } from "./arrow-right-BKIeTJ3r.js";
const ArchiveRestore = createLucideIcon("ArchiveRestoreIcon", [
  ["rect", { width: "20", height: "5", x: "2", y: "3", rx: "1", key: "1wp1u1" }],
  ["path", { d: "M4 8v11a2 2 0 0 0 2 2h2", key: "tvwodi" }],
  ["path", { d: "M20 8v11a2 2 0 0 1-2 2h-2", key: "1gkqxj" }],
  ["path", { d: "m9 15 3-3 3 3", key: "1pd0qc" }],
  ["path", { d: "M12 12v9", key: "192myk" }]
]);
const Lock = createLucideIcon("LockIcon", [
  ["rect", { width: "18", height: "11", x: "3", y: "11", rx: "2", ry: "2", key: "1w4ew1" }],
  ["path", { d: "M7 11V7a5 5 0 0 1 10 0v4", key: "fwvmzm" }]
]);
const _hoisted_1 = { class: "dh-auth" };
const _hoisted_2 = { class: "dh-auth-brand" };
const _hoisted_3 = { class: "dh-auth-mid" };
const _hoisted_4 = { class: "dh-auth-feats" };
const _hoisted_5 = { class: "dh-auth-ico" };
const _hoisted_6 = { class: "d" };
const _hoisted_7 = { class: "dh-auth-form" };
const _hoisted_8 = {
  key: 2,
  class: "grid h-[120px] place-items-center text-text-5"
};
const _hoisted_9 = { class: "mb-[15px]" };
const _hoisted_10 = {
  class: "mb-[7px] block text-[12.5px] font-semibold text-text-3",
  for: "dh-pw"
};
const _hoisted_11 = { class: "dh-field !h-[42px]" };
const _hoisted_12 = ["type", "placeholder", "autocomplete"];
const _hoisted_13 = ["title", "aria-label"];
const _hoisted_14 = { class: "mt-[8px] flex gap-1" };
const _hoisted_15 = { class: "mt-[6px] text-[11.5px] text-text-6" };
const _hoisted_16 = {
  key: 0,
  class: "mb-[15px]"
};
const _hoisted_17 = { class: "dh-field !h-[42px]" };
const _hoisted_18 = ["type"];
const _hoisted_19 = {
  key: 0,
  class: "mt-[6px] text-[11.5px] text-err-text"
};
const _hoisted_20 = {
  key: 2,
  class: "flex cursor-pointer items-center gap-2 text-[12px] text-text-3"
};
const _hoisted_21 = ["disabled", "aria-busy"];
const _hoisted_22 = {
  key: 3,
  class: "dh-auth-note"
};
const _hoisted_23 = { class: "dh-auth-note" };
const _hoisted_24 = { class: "dh-auth-forgot" };
const _hoisted_25 = {
  key: 0,
  class: "dh-auth-tip"
};
const _hoisted_26 = { class: "dh-auth-formfoot" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "LoginView",
  setup(__props) {
    const app = useAppStore();
    const router = useRouter();
    const route = useRoute();
    const mode = ref("loading");
    const password = ref("");
    const confirm = ref("");
    const showPw = ref(false);
    const pwInput = ref(null);
    const showForgot = ref(false);
    const keep = ref(true);
    const busy = ref(false);
    const errorMsg = ref("");
    const remaining = ref(null);
    const locked = ref(false);
    const lockDeadline = ref(0);
    const nowTick = ref(Date.now());
    let lockTimer;
    const features = [
      { icon: Boxes, name: "容器", desc: "启停重启、日志跟随、实时资源占用，环境变量自动脱敏" },
      { icon: Layers, name: "镜像与加速源", desc: "体积与引用数、清理未使用；实测生效的加速站与 daemon.json" },
      { icon: RefreshCw, name: "更新与自动更新", desc: "巡检后标出可更新，按周期自动更新；镜像没变就不动容器" },
      { icon: CalendarClock, name: "计划任务", desc: "5 段 cron：启动 / 停止 / 重启 / 更新 / 备份 / 清理" },
      { icon: ListOrdered, name: "顺序启动", desc: "开机后按你排的顺序串行拉起：拖拽排序、逐行延时、逐行开关" },
      { icon: ArchiveRestore, name: "备份与恢复", desc: "配置快照、差异预览、一键还原；compose / .env 一并备份" },
      { icon: Bell, name: "通知", desc: "11 类渠道、按事件订阅；静默时段、去重、失败重试与推送历史" }
    ];
    const lockLeft = computed(() => {
      const ms = lockDeadline.value - nowTick.value;
      return ms > 0 ? Math.ceil(ms / 1e3) : 0;
    });
    function mmss(total) {
      const m = Math.floor(total / 60);
      const s = total % 60;
      return `${m}:${String(s).padStart(2, "0")}`;
    }
    function ensureLockTicker() {
      if (lockTimer !== void 0) return;
      lockTimer = window.setInterval(() => {
        nowTick.value = Date.now();
        if (lockDeadline.value <= Date.now()) clearLock();
      }, 1e3);
    }
    function applyLock(msg, seconds) {
      errorMsg.value = msg;
      locked.value = true;
      if (seconds > 0) {
        lockDeadline.value = Date.now() + seconds * 1e3;
        nowTick.value = Date.now();
        ensureLockTicker();
      } else {
        lockDeadline.value = 0;
      }
    }
    function clearLock() {
      locked.value = false;
      lockDeadline.value = 0;
      errorMsg.value = "";
      remaining.value = app.maxFailures;
      if (lockTimer !== void 0) {
        window.clearInterval(lockTimer);
        lockTimer = void 0;
      }
    }
    onUnmounted(() => {
      if (lockTimer !== void 0) window.clearInterval(lockTimer);
    });
    const isSetup = computed(() => mode.value === "setup");
    const strength = computed(() => {
      const p = password.value;
      if (!p) return 0;
      let s = 0;
      if (p.length >= app.minPasswordLength) s++;
      if (p.length >= 10) s++;
      if (/[A-Za-z]/.test(p) && /\d/.test(p)) s++;
      if (/[^A-Za-z0-9]/.test(p)) s++;
      return Math.min(s, 4);
    });
    const strengthLabel = computed(() => ["", "偏弱", "一般", "较强", "很强"][strength.value]);
    const strengthHint = computed(
      () => password.value ? `强度：${strengthLabel.value} · 至少 ${app.minPasswordLength} 位，建议含数字与符号` : `至少 ${app.minPasswordLength} 位，建议含数字与符号`
    );
    const canSubmit = computed(() => {
      if (busy.value || locked.value) return false;
      if (isSetup.value)
        return password.value.length >= app.minPasswordLength && password.value === confirm.value;
      return password.value.length > 0;
    });
    const bannerTone = computed(
      () => locked.value || isSetup.value ? "dh-banner-err" : "dh-banner-warn"
    );
    const bannerText = computed(() => {
      if (locked.value) {
        return lockLeft.value > 0 ? `${errorMsg.value}（剩余 ${mmss(lockLeft.value)}）` : errorMsg.value;
      }
      if (!errorMsg.value) return "";
      if (remaining.value !== null && remaining.value > 0) {
        return `${errorMsg.value}，还可尝试 ${remaining.value} 次`;
      }
      return errorMsg.value;
    });
    const submitLabel = computed(() => {
      if (busy.value) return "请稍候…";
      if (locked.value) return lockLeft.value > 0 ? `已锁定 ${mmss(lockLeft.value)}` : "已锁定";
      return isSetup.value ? "设置并进入" : "登 录";
    });
    function goAfterLogin() {
      const back = safeRedirect(route.query.redirect);
      router.replace(back ?? "/overview");
    }
    onMounted(async () => {
      if (!app.ready) await app.bootstrap();
      if (!app.initialized) {
        mode.value = "setup";
        await nextTick();
        pwInput.value?.focus();
        return;
      }
      if (app.loggedIn) {
        goAfterLogin();
        return;
      }
      mode.value = "login";
      remaining.value = app.maxFailures - app.failures;
      if (app.lockedFor > 0) applyLock(app.lockedHint, app.lockedFor);
      await nextTick();
      pwInput.value?.focus();
    });
    async function submit() {
      if (!canSubmit.value) return;
      busy.value = true;
      errorMsg.value = "";
      try {
        if (isSetup.value) {
          const res = await app.setup(password.value, confirm.value);
          if (res.autoLogin) {
            goAfterLogin();
          } else {
            mode.value = "login";
            password.value = "";
            confirm.value = "";
            await nextTick();
            pwInput.value?.focus();
          }
        } else {
          await app.login(password.value, keep.value);
          goAfterLogin();
        }
      } catch (e) {
        if (e instanceof ApiError) {
          const payload = e.payload;
          if (e.code === "locked") {
            applyLock(e.message, typeof payload?.retryAfter === "number" ? payload.retryAfter : 0);
          } else {
            errorMsg.value = e.message;
            if (typeof payload?.remaining === "number") remaining.value = payload.remaining;
            if (e.code === "bad_password") password.value = "";
            await nextTick();
            pwInput.value?.focus();
          }
        } else {
          errorMsg.value = e instanceof Error ? e.message : String(e);
        }
      } finally {
        busy.value = false;
      }
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("aside", _hoisted_2, [
          _cache[8] || (_cache[8] = createStaticVNode('<div class="dh-auth-top"><div class="dh-auth-mark"><svg viewBox="0 0 32 32" class="h-[21px] w-[21px]"><circle cx="16" cy="16" r="12.3" fill="none" stroke="currentColor" stroke-width="2.6"></circle><path fill-rule="evenodd" fill="currentColor" d="M16 5.9 18.76 13.24 26.1 16 18.76 18.76 16 26.1 13.24 18.76 5.9 16 13.24 13.24ZM18.2 16A2.2 2.2 0 1 0 13.8 16A2.2 2.2 0 1 0 18.2 16Z"></path></svg></div><div class="text-[16px] font-semibold tracking-[-.2px]">Dockhelm</div></div>', 1)),
          createBaseVNode("div", _hoisted_3, [
            _cache[5] || (_cache[5] = createBaseVNode("h1", null, [
              createTextVNode("容器、镜像、任务，"),
              createBaseVNode("br"),
              createTextVNode("一个面板全收住")
            ], -1)),
            _cache[6] || (_cache[6] = createBaseVNode("p", { class: "dh-auth-lead" }, [
              createTextVNode(" 打开就是总览：容器、镜像、更新、任务都在这儿。"),
              createBaseVNode("br"),
              createTextVNode(" 不用再记一堆 docker 命令，也不用 SSH 翻 compose。 ")
            ], -1)),
            createBaseVNode("ul", _hoisted_4, [
              (openBlock(), createElementBlock(Fragment, null, renderList(features, (f) => {
                return createBaseVNode("li", {
                  key: f.name
                }, [
                  createBaseVNode("span", _hoisted_5, [
                    (openBlock(), createBlock(resolveDynamicComponent(f.icon), { class: "h-[15px] w-[15px]" }))
                  ]),
                  createBaseVNode("div", null, [
                    createBaseVNode("b", null, toDisplayString(f.name), 1),
                    createBaseVNode("span", _hoisted_6, toDisplayString(f.desc), 1)
                  ])
                ]);
              }), 64))
            ]),
            _cache[7] || (_cache[7] = createBaseVNode("p", { class: "dh-auth-extra" }, [
              createTextVNode(" 登录只要一个密码："),
              createBaseVNode("b", null, "bcrypt"),
              createTextVNode(" 哈希落盘、HttpOnly 会话、连续失败锁定。 ")
            ], -1))
          ]),
          _cache[9] || (_cache[9] = createBaseVNode("div", { class: "dh-auth-brandfoot" }, "© 2026 Dockhelm · MIT 许可", -1))
        ]),
        createBaseVNode("main", _hoisted_7, [
          createBaseVNode("form", {
            class: "dh-auth-box",
            onSubmit: withModifiers(submit, ["prevent"])
          }, [
            isSetup.value ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
              _cache[10] || (_cache[10] = createBaseVNode("span", { class: "dh-auth-eyebrow" }, [
                createBaseVNode("i"),
                createTextVNode("首次启动")
              ], -1)),
              _cache[11] || (_cache[11] = createBaseVNode("h2", null, "设置访问密码", -1)),
              _cache[12] || (_cache[12] = createBaseVNode("p", { class: "dh-auth-sub" }, "这台 Dockhelm 还没有密码。先设一个，之后用它进入面板。", -1))
            ], 64)) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
              _cache[13] || (_cache[13] = createBaseVNode("h2", null, "登录 Dockhelm", -1)),
              _cache[14] || (_cache[14] = createBaseVNode("p", { class: "dh-auth-sub" }, "输入访问密码继续。", -1))
            ], 64)),
            mode.value === "loading" ? (openBlock(), createElementBlock("div", _hoisted_8, [
              createVNode(unref(LoaderCircle), { class: "h-5 w-5 dh-spin" })
            ])) : (openBlock(), createElementBlock(Fragment, { key: 3 }, [
              createBaseVNode("div", _hoisted_9, [
                createBaseVNode("label", _hoisted_10, toDisplayString(isSetup.value ? "密码" : "访问密码"), 1),
                createBaseVNode("div", _hoisted_11, [
                  createVNode(unref(Lock), { class: "h-[15px] w-[15px] flex-none text-text-6" }),
                  withDirectives(createBaseVNode("input", {
                    id: "dh-pw",
                    ref_key: "pwInput",
                    ref: pwInput,
                    "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => password.value = $event),
                    type: showPw.value ? "text" : "password",
                    placeholder: isSetup.value ? `至少 ${unref(app).minPasswordLength} 位` : "请输入密码",
                    autocomplete: isSetup.value ? "new-password" : "current-password",
                    autofocus: ""
                  }, null, 8, _hoisted_12), [
                    [vModelDynamic, password.value]
                  ]),
                  createBaseVNode("button", {
                    type: "button",
                    class: "dh-eye",
                    title: showPw.value ? "隐藏密码" : "显示密码",
                    "aria-label": showPw.value ? "隐藏密码" : "显示密码",
                    onClick: _cache[1] || (_cache[1] = ($event) => showPw.value = !showPw.value)
                  }, [
                    showPw.value ? (openBlock(), createBlock(unref(EyeOff), {
                      key: 0,
                      class: "h-[15px] w-[15px]"
                    })) : (openBlock(), createBlock(unref(Eye), {
                      key: 1,
                      class: "h-[15px] w-[15px]"
                    }))
                  ], 8, _hoisted_13)
                ]),
                isSetup.value ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                  createBaseVNode("div", _hoisted_14, [
                    (openBlock(), createElementBlock(Fragment, null, renderList(4, (i) => {
                      return createBaseVNode("i", {
                        key: i,
                        class: normalizeClass(["h-[4px] flex-1 rounded-[2px]", i <= strength.value ? "bg-run" : "bg-line-1"])
                      }, null, 2);
                    }), 64))
                  ]),
                  createBaseVNode("p", _hoisted_15, toDisplayString(strengthHint.value), 1)
                ], 64)) : createCommentVNode("", true)
              ]),
              isSetup.value ? (openBlock(), createElementBlock("div", _hoisted_16, [
                _cache[15] || (_cache[15] = createBaseVNode("label", {
                  class: "mb-[7px] block text-[12.5px] font-semibold text-text-3",
                  for: "dh-pw2"
                }, " 确认密码 ", -1)),
                createBaseVNode("div", _hoisted_17, [
                  createVNode(unref(Lock), { class: "h-[15px] w-[15px] flex-none text-text-6" }),
                  withDirectives(createBaseVNode("input", {
                    id: "dh-pw2",
                    "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => confirm.value = $event),
                    type: showPw.value ? "text" : "password",
                    placeholder: "再输一次",
                    autocomplete: "new-password"
                  }, null, 8, _hoisted_18), [
                    [vModelDynamic, confirm.value]
                  ])
                ]),
                confirm.value && confirm.value !== password.value ? (openBlock(), createElementBlock("p", _hoisted_19, " 两次输入不一致 ")) : createCommentVNode("", true)
              ])) : createCommentVNode("", true),
              bannerText.value ? (openBlock(), createElementBlock("div", {
                key: 1,
                class: normalizeClass(["dh-banner mb-[14px] !items-start !gap-2 !py-[9px] !pl-[11px] !pr-[11px] !text-[12.5px]", bannerTone.value]),
                role: "alert"
              }, [
                createVNode(unref(TriangleAlert), { class: "h-[14px] w-[14px] flex-none translate-y-[1px]" }),
                createBaseVNode("span", null, toDisplayString(bannerText.value), 1)
              ], 2)) : createCommentVNode("", true),
              !isSetup.value ? (openBlock(), createElementBlock("label", _hoisted_20, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => keep.value = $event),
                  type: "checkbox",
                  class: "h-[14px] w-[14px] accent-accent"
                }, null, 512), [
                  [vModelCheckbox, keep.value]
                ]),
                _cache[16] || (_cache[16] = createTextVNode(" 保持登录（7 天） ", -1))
              ])) : createCommentVNode("", true),
              createBaseVNode("button", {
                type: "submit",
                class: "dh-auth-submit",
                disabled: !canSubmit.value,
                "aria-busy": busy.value
              }, [
                busy.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "h-[15px] w-[15px] dh-spin"
                })) : createCommentVNode("", true),
                createBaseVNode("span", null, toDisplayString(submitLabel.value), 1),
                !busy.value ? (openBlock(), createBlock(unref(ArrowRight), {
                  key: 1,
                  class: "arrow"
                })) : createCommentVNode("", true)
              ], 8, _hoisted_21),
              isSetup.value ? (openBlock(), createElementBlock("p", _hoisted_22, [..._cache[17] || (_cache[17] = [
                createTextVNode(" 密码只存 ", -1),
                createBaseVNode("code", { class: "font-mono" }, "bcrypt", -1),
                createTextVNode(" 哈希，明文不落盘、也不写日志。", -1),
                createBaseVNode("br", null, null, -1),
                createTextVNode(" Dockhelm 持有 Docker 套接字，等于宿主机 root 权限，请务必设置密码。 ", -1)
              ])])) : (openBlock(), createElementBlock(Fragment, { key: 4 }, [
                createBaseVNode("p", _hoisted_23, [
                  createTextVNode(" 连续错 " + toDisplayString(unref(app).maxFailures) + " 次会临时锁定，最长 5 分钟", 1),
                  _cache[18] || (_cache[18] = createBaseVNode("br", null, null, -1)),
                  _cache[19] || (_cache[19] = createTextVNode(" 会话连续 48 小时没使用也会自动失效 ", -1))
                ]),
                createBaseVNode("div", _hoisted_24, [
                  createBaseVNode("button", {
                    type: "button",
                    class: "dh-auth-link",
                    onClick: _cache[4] || (_cache[4] = ($event) => showForgot.value = !showForgot.value)
                  }, toDisplayString(showForgot.value ? "收起" : "忘记密码？"), 1),
                  showForgot.value ? (openBlock(), createElementBlock("div", _hoisted_25, [..._cache[20] || (_cache[20] = [
                    createBaseVNode("p", null, "面板不联网、也不存明文密码，所以没有「找回」入口，只能重置：", -1),
                    createBaseVNode("ol", null, [
                      createBaseVNode("li", null, "在 NAS 的文件管理里找到 Dockhelm 的数据目录（compose 里挂载的那个）。"),
                      createBaseVNode("li", null, [
                        createTextVNode("删掉目录里的 "),
                        createBaseVNode("code", null, "auth.json"),
                        createTextVNode("。")
                      ]),
                      createBaseVNode("li", null, "重启容器，再打开面板就会回到「设置密码」。")
                    ], -1),
                    createBaseVNode("p", { class: "dh-auth-tip-warn" }, "重置后，所有设备上已有的登录会一并失效。", -1)
                  ])])) : createCommentVNode("", true)
                ])
              ], 64))
            ], 64))
          ], 32),
          createBaseVNode("div", _hoisted_26, [
            createTextVNode(" Dockhelm " + toDisplayString(unref(version)), 1),
            _cache[21] || (_cache[21] = createBaseVNode("span", { class: "sep" }, "·", -1)),
            _cache[22] || (_cache[22] = createTextVNode("单用户密码 ", -1))
          ])
        ])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
