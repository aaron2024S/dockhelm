import { d as defineComponent, u as useAppStore, o as onUnmounted, a as onMounted, c as createElementBlock, b as createStaticVNode, e as createBaseVNode, t as toDisplayString, f as createVNode, g as unref, F as Fragment, w as withDirectives, v as vModelDynamic, h as withKeys, i as createBlock, r as renderList, j as createCommentVNode, n as normalizeClass, k as createTextVNode, l as vModelCheckbox, m as ref, p as computed, A as ApiError, q as useRouter, s as openBlock } from "./index-txkAKhMR.js";
import { L as LoaderCircle } from "./loader-circle-CpKqd4Ou.js";
import { E as EyeOff } from "./eye-off-DR72r1Lc.js";
import { E as Eye } from "./eye-BGv27FL5.js";
const _hoisted_1 = { class: "flex min-h-screen flex-col bg-ink-900" };
const _hoisted_2 = { class: "relative flex flex-1 flex-col items-center justify-center gap-[14px] px-5 py-10" };
const _hoisted_3 = { class: "relative flex w-[312px] flex-col gap-[13px] rounded-[16px] border border-line-1 bg-ink-700 p-[22px]" };
const _hoisted_4 = { class: "text-center" };
const _hoisted_5 = { class: "text-[16px] font-semibold" };
const _hoisted_6 = { class: "mt-[3px] text-[12px] text-text-4" };
const _hoisted_7 = {
  key: 0,
  class: "grid h-[110px] place-items-center text-text-5"
};
const _hoisted_8 = { class: "dh-field" };
const _hoisted_9 = ["type", "placeholder"];
const _hoisted_10 = ["title", "aria-label"];
const _hoisted_11 = { class: "mt-[7px] flex gap-1" };
const _hoisted_12 = { class: "mt-[5px] text-[11.5px] text-text-5" };
const _hoisted_13 = { class: "dh-field" };
const _hoisted_14 = ["type"];
const _hoisted_15 = {
  key: 0,
  class: "mt-[5px] text-[11.5px] text-err-text"
};
const _hoisted_16 = {
  key: 1,
  class: "dh-field"
};
const _hoisted_17 = ["type"];
const _hoisted_18 = ["title", "aria-label"];
const _hoisted_19 = {
  key: 3,
  class: "flex cursor-pointer items-center gap-2 text-[12px] text-text-3"
};
const _hoisted_20 = ["disabled"];
const _hoisted_21 = {
  key: 4,
  class: "text-center text-[11.5px] text-text-6"
};
const _hoisted_22 = {
  key: 0,
  class: "relative w-[312px] text-center text-[11.5px] leading-[1.7] text-text-6"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "LoginView",
  setup(__props) {
    const app = useAppStore();
    const router = useRouter();
    const mode = ref("loading");
    const password = ref("");
    const confirm = ref("");
    const showPw = ref(false);
    const keep = ref(true);
    const busy = ref(false);
    const errorMsg = ref("");
    const remaining = ref(null);
    const locked = ref(false);
    const lockDeadline = ref(0);
    const nowTick = ref(Date.now());
    let lockTimer;
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
    const title = computed(() => isSetup.value ? "首次启动" : "Dockhelm");
    const subtitle = computed(() => isSetup.value ? "设置访问密码后才能使用" : "请输入密码以继续");
    const strength = computed(() => {
      const p = password.value;
      if (!p) return 0;
      let s = 0;
      if (p.length >= 6) s++;
      if (p.length >= 10) s++;
      if (/[A-Za-z]/.test(p) && /\d/.test(p)) s++;
      if (/[^A-Za-z0-9]/.test(p)) s++;
      return Math.min(s, 4);
    });
    const strengthLabel = computed(() => ["", "偏弱", "一般", "较强", "很强"][strength.value]);
    const canSubmit = computed(() => {
      if (busy.value || locked.value) return false;
      if (isSetup.value) return password.value.length >= app.minPasswordLength && password.value === confirm.value;
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
    onMounted(async () => {
      if (!app.ready) await app.bootstrap();
      if (!app.initialized) {
        mode.value = "setup";
        return;
      }
      if (app.loggedIn) {
        router.replace("/overview");
        return;
      }
      mode.value = "login";
      remaining.value = app.maxFailures - app.failures;
      if (app.lockedFor > 0) applyLock(app.lockedHint, app.lockedFor);
    });
    async function submit() {
      if (!canSubmit.value) return;
      busy.value = true;
      errorMsg.value = "";
      try {
        if (isSetup.value) {
          const res = await app.setup(password.value, confirm.value);
          if (res.autoLogin) {
            router.replace("/overview");
          } else {
            mode.value = "login";
            password.value = "";
            confirm.value = "";
          }
        } else {
          await app.login(password.value, keep.value);
          router.replace("/overview");
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
          }
        } else {
          errorMsg.value = e instanceof Error ? e.message : String(e);
        }
      } finally {
        busy.value = false;
      }
    }
    function onEnter() {
      void submit();
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        _cache[13] || (_cache[13] = createStaticVNode('<header class="flex flex-none items-center gap-2.5 border-b border-line-2 bg-ink-850 px-[18px] py-3"><div class="grid h-[27px] w-[27px] flex-none place-items-center rounded-[9px] bg-accent text-accent-ink"><svg viewBox="0 0 32 32" class="h-[15px] w-[15px]"><circle cx="16" cy="16" r="12.3" fill="none" stroke="currentColor" stroke-width="2.6"></circle><path fill-rule="evenodd" fill="currentColor" d="M16 6.6 18.5 13.5 25.4 16 18.5 18.5 16 25.4 13.5 18.5 6.6 16 13.5 13.5ZM18 16A2 2 0 1 0 14 16A2 2 0 1 0 18 16Z"></path></svg></div><div class="text-[14px] font-semibold tracking-[.2px]">Dockhelm</div><div class="text-[12px] text-text-4">登录与账户</div></header>', 1)),
        createBaseVNode("div", _hoisted_2, [
          _cache[12] || (_cache[12] = createBaseVNode("div", {
            class: "pointer-events-none absolute inset-x-0 top-0 h-[420px]",
            style: { "background": "radial-gradient(560px 260px at 50% 0%, var(--color-glow), transparent 72%)" }
          }, null, -1)),
          createBaseVNode("div", _hoisted_3, [
            _cache[10] || (_cache[10] = createBaseVNode("div", { class: "dh-lglogo" }, [
              createBaseVNode("svg", {
                viewBox: "0 0 32 32",
                class: "h-[23px] w-[23px]"
              }, [
                createBaseVNode("circle", {
                  cx: "16",
                  cy: "16",
                  r: "12.3",
                  fill: "none",
                  stroke: "currentColor",
                  "stroke-width": "2.5"
                }),
                createBaseVNode("path", {
                  "fill-rule": "evenodd",
                  fill: "currentColor",
                  d: "M16 6.6 18.5 13.5 25.4 16 18.5 18.5 16 25.4 13.5 18.5 6.6 16 13.5 13.5ZM18 16A2 2 0 1 0 14 16A2 2 0 1 0 18 16Z"
                })
              ])
            ], -1)),
            createBaseVNode("div", _hoisted_4, [
              createBaseVNode("div", _hoisted_5, toDisplayString(title.value), 1),
              createBaseVNode("div", _hoisted_6, toDisplayString(subtitle.value), 1)
            ]),
            mode.value === "loading" ? (openBlock(), createElementBlock("div", _hoisted_7, [
              createVNode(unref(LoaderCircle), { class: "h-5 w-5 dh-spin" })
            ])) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
              isSetup.value ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                createBaseVNode("div", null, [
                  _cache[6] || (_cache[6] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "新密码", -1)),
                  createBaseVNode("div", _hoisted_8, [
                    withDirectives(createBaseVNode("input", {
                      "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => password.value = $event),
                      type: showPw.value ? "text" : "password",
                      placeholder: `至少 ${unref(app).minPasswordLength} 位`,
                      autocomplete: "new-password",
                      autofocus: "",
                      onKeyup: withKeys(onEnter, ["enter"])
                    }, null, 40, _hoisted_9), [
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
                    ], 8, _hoisted_10)
                  ]),
                  createBaseVNode("div", _hoisted_11, [
                    (openBlock(), createElementBlock(Fragment, null, renderList(4, (i) => {
                      return createBaseVNode("i", {
                        key: i,
                        class: normalizeClass(["h-[4px] flex-1 rounded-[2px]", i <= strength.value ? "bg-run" : "bg-line-1"])
                      }, null, 2);
                    }), 64))
                  ]),
                  createBaseVNode("div", _hoisted_12, " 强度：" + toDisplayString(strengthLabel.value) + " · 至少 " + toDisplayString(unref(app).minPasswordLength) + " 位，建议含数字与符号 ", 1)
                ]),
                createBaseVNode("div", null, [
                  _cache[7] || (_cache[7] = createBaseVNode("label", { class: "mb-[5px] block text-[12px] text-text-4" }, "确认密码", -1)),
                  createBaseVNode("div", _hoisted_13, [
                    withDirectives(createBaseVNode("input", {
                      "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => confirm.value = $event),
                      type: showPw.value ? "text" : "password",
                      placeholder: "再输一次",
                      autocomplete: "new-password",
                      onKeyup: withKeys(onEnter, ["enter"])
                    }, null, 40, _hoisted_14), [
                      [vModelDynamic, confirm.value]
                    ])
                  ]),
                  confirm.value && confirm.value !== password.value ? (openBlock(), createElementBlock("div", _hoisted_15, " 两次输入不一致 ")) : createCommentVNode("", true)
                ])
              ], 64)) : (openBlock(), createElementBlock("div", _hoisted_16, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => password.value = $event),
                  type: showPw.value ? "text" : "password",
                  placeholder: "请输入密码",
                  autocomplete: "current-password",
                  autofocus: "",
                  onKeyup: withKeys(onEnter, ["enter"])
                }, null, 40, _hoisted_17), [
                  [vModelDynamic, password.value]
                ]),
                createBaseVNode("button", {
                  type: "button",
                  class: "dh-eye",
                  title: showPw.value ? "隐藏密码" : "显示密码",
                  "aria-label": showPw.value ? "隐藏密码" : "显示密码",
                  onClick: _cache[4] || (_cache[4] = ($event) => showPw.value = !showPw.value)
                }, [
                  showPw.value ? (openBlock(), createBlock(unref(EyeOff), {
                    key: 0,
                    class: "h-[15px] w-[15px]"
                  })) : (openBlock(), createBlock(unref(Eye), {
                    key: 1,
                    class: "h-[15px] w-[15px]"
                  }))
                ], 8, _hoisted_18)
              ])),
              bannerText.value ? (openBlock(), createElementBlock("div", {
                key: 2,
                class: normalizeClass(["dh-banner !gap-2.5 !py-2 !pl-[11px] !pr-[11px] !text-[12px]", bannerTone.value])
              }, [
                _cache[8] || (_cache[8] = createBaseVNode("span", { class: "h-[13px] w-[13px] flex-none rounded-[4px] bg-current opacity-50" }, null, -1)),
                createTextVNode(" " + toDisplayString(bannerText.value), 1)
              ], 2)) : createCommentVNode("", true),
              !isSetup.value ? (openBlock(), createElementBlock("label", _hoisted_19, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => keep.value = $event),
                  type: "checkbox",
                  class: "h-[14px] w-[14px] accent-accent"
                }, null, 512), [
                  [vModelCheckbox, keep.value]
                ]),
                _cache[9] || (_cache[9] = createTextVNode(" 保持登录（7 天） ", -1))
              ])) : createCommentVNode("", true),
              createBaseVNode("button", {
                type: "button",
                class: "w-full rounded-[10px] bg-accent px-3 py-2.5 text-center text-[13px] font-semibold text-accent-ink transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-45",
                disabled: !canSubmit.value,
                onClick: submit
              }, toDisplayString(busy.value ? "请稍候…" : locked.value ? `已锁定 ${mmss(lockLeft.value)}` : isSetup.value ? "完成并进入" : "登 录"), 9, _hoisted_20),
              !isSetup.value ? (openBlock(), createElementBlock("div", _hoisted_21, " 连续错 " + toDisplayString(unref(app).maxFailures) + " 次锁定 5 分钟 ", 1)) : createCommentVNode("", true)
            ], 64))
          ]),
          isSetup.value ? (openBlock(), createElementBlock("div", _hoisted_22, [..._cache[11] || (_cache[11] = [
            createTextVNode(" 密码只存 ", -1),
            createBaseVNode("span", { class: "font-mono" }, "bcrypt", -1),
            createTextVNode(" 哈希，明文不落盘、也不写日志。 ", -1),
            createBaseVNode("br", null, null, -1),
            createTextVNode(" Dockhelm 持有 Docker 套接字，等于拥有宿主机 root 权限，请务必设置密码。 ", -1)
          ])])) : createCommentVNode("", true)
        ])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
