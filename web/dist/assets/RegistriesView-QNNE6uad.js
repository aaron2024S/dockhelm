import { x as createLucideIcon, d as defineComponent, E as useToastStore, u as useAppStore, a as onMounted, c as createElementBlock, e as createBaseVNode, f as createVNode, n as normalizeClass, g as unref, Q as Gauge, k as createTextVNode, T as Info, t as toDisplayString, F as Fragment, r as renderList, U as Rocket, j as createCommentVNode, w as withDirectives, G as vModelText, i as createBlock, N as vModelSelect, I as withCtx, O as resolveDynamicComponent, V as CircleCheck, m as ref, p as computed, B as api, s as openBlock } from "./index-DBOlVmiV.js";
import { f as formatDayTime } from "./format-NZu8lNOP.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-CbWMlniA.js";
import { _ as _sfc_main$2, a as _sfc_main$4 } from "./ToggleSwitch.vue_vue_type_script_setup_true_lang-CCzcTs2h.js";
import { P as Plus, _ as _sfc_main$3 } from "./UiSwitch.vue_vue_type_script_setup_true_lang-oDTqX6ny.js";
import { T as TriangleAlert } from "./triangle-alert-CC8PYdGF.js";
import { Z as Zap } from "./zap-C9j9lw7d.js";
import { L as LoaderCircle } from "./loader-circle-CFBmSG65.js";
import { S as Save } from "./save-DsALYqad.js";
import { C as Copy } from "./copy-Cx9zOmPR.js";
const GripVertical = createLucideIcon("GripVerticalIcon", [
  ["circle", { cx: "9", cy: "12", r: "1", key: "1vctgf" }],
  ["circle", { cx: "9", cy: "5", r: "1", key: "hp0tcf" }],
  ["circle", { cx: "9", cy: "19", r: "1", key: "fkjjf6" }],
  ["circle", { cx: "15", cy: "12", r: "1", key: "1tmaij" }],
  ["circle", { cx: "15", cy: "5", r: "1", key: "19l28e" }],
  ["circle", { cx: "15", cy: "19", r: "1", key: "f4zoj3" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "ml-auto flex gap-2" };
const _hoisted_4 = ["disabled"];
const _hoisted_5 = { class: "dh-card" };
const _hoisted_6 = { class: "dh-card-head" };
const _hoisted_7 = { class: "dh-card-body flex flex-col gap-2.5" };
const _hoisted_8 = { class: "text-[12px] leading-relaxed text-text-4" };
const _hoisted_9 = {
  key: 0,
  class: "text-[12px] text-err-text"
};
const _hoisted_10 = {
  key: 1,
  class: "flex flex-wrap gap-1.5"
};
const _hoisted_11 = {
  key: 2,
  class: "flex items-center gap-2 text-[12px] text-text-5"
};
const _hoisted_12 = { class: "grid grid-cols-1 gap-3.5 xl:grid-cols-2" };
const _hoisted_13 = { class: "dh-card xl:col-span-2" };
const _hoisted_14 = { class: "dh-card-head" };
const _hoisted_15 = { class: "ml-auto flex items-center gap-2" };
const _hoisted_16 = {
  key: 0,
  class: "text-[11.5px] text-warn-text"
};
const _hoisted_17 = ["disabled"];
const _hoisted_18 = ["disabled"];
const _hoisted_19 = {
  id: "mirror-add",
  class: "flex flex-col gap-2 border-b border-line-1 p-3"
};
const _hoisted_20 = { class: "flex flex-wrap gap-2" };
const _hoisted_21 = ["disabled"];
const _hoisted_22 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_23 = { class: "dh-table dh-table-fixed" };
const _hoisted_24 = { class: "!px-2" };
const _hoisted_25 = ["aria-label", "onPointerdown", "onKeydown"];
const _hoisted_26 = ["title"];
const _hoisted_27 = { class: "mt-0.5 flex items-center gap-1.5 text-[11px] text-text-5" };
const _hoisted_28 = {
  key: 0,
  class: "dh-badge dh-badge-plain flex-none"
};
const _hoisted_29 = {
  key: 1,
  class: "truncate"
};
const _hoisted_30 = { class: "text-[11.5px] text-text-5" };
const _hoisted_31 = { class: "text-right" };
const _hoisted_32 = { class: "text-right whitespace-nowrap" };
const _hoisted_33 = ["disabled", "onClick"];
const _hoisted_34 = ["onClick"];
const _hoisted_35 = { class: "dh-card" };
const _hoisted_36 = { class: "dh-card-head" };
const _hoisted_37 = { class: "dh-card-body flex flex-col gap-3" };
const _hoisted_38 = ["value"];
const _hoisted_39 = {
  key: 0,
  class: "border-t border-line-1 pt-3"
};
const _hoisted_40 = ["disabled"];
const _hoisted_41 = { class: "dh-card" };
const _hoisted_42 = { class: "dh-card-head" };
const _hoisted_43 = { class: "dh-card-body" };
const _hoisted_44 = { class: "mt-2.5 overflow-x-auto rounded-[9px] border border-line-3 bg-ink-800 px-3 py-2.5 font-mono text-[11.5px] leading-[1.8] text-text-2" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "RegistriesView",
  setup(__props) {
    const toast = useToastStore();
    const app = useAppStore();
    const data = ref(null);
    const loading = ref(true);
    const testing = ref(false);
    const testingUrl = ref("");
    const newUrl = ref("");
    const newNote = ref("");
    const copied = ref(false);
    const saving = ref(false);
    const mirrors = computed(() => data.value?.settings?.mirrors ?? []);
    const presets = computed(() => data.value?.presets ?? []);
    const daemonMirrors = computed(() => data.value?.daemonMirrors ?? []);
    function urlKey(u) {
      return u.trim().toLowerCase().replace(/\/+$/, "");
    }
    const pullMirror = computed({
      get: () => data.value?.settings?.pullMirror ?? "",
      set: (v) => {
        if (data.value?.settings) data.value.settings.pullMirror = v;
      }
    });
    async function load() {
      loading.value = true;
      try {
        data.value = await api.get("/api/registries");
        savedSig.value = sigOf(data.value.settings);
      } catch (e) {
        toast.error("读取加速源配置失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function save() {
      if (!data.value) return;
      saving.value = true;
      let registrySaved = false;
      try {
        await api.put("/api/registries", { settings: data.value.settings });
        registrySaved = true;
        if (policy.value) {
          policy.value = await api.patch("/api/settings", { directFirst: policy.value.directFirst });
        }
        void app.loadSettings();
        toast.success("已保存");
        await load();
      } catch (e) {
        toast.error(
          registrySaved ? "加速源已保存，但 directFirst 策略没保存成功" : "保存失败",
          e instanceof Error ? e.message : String(e)
        );
      } finally {
        saving.value = false;
      }
    }
    function addMirror() {
      const url = newUrl.value.trim();
      if (!url || !data.value) return;
      if (mirrors.value.some((m) => urlKey(m.url) === urlKey(url))) {
        toast.error("这个地址已经在列表里了");
        return;
      }
      data.value.settings.mirrors.push({ url, note: newNote.value.trim(), enabled: true });
      newUrl.value = "";
      newNote.value = "";
    }
    async function addPresets() {
      if (!data.value) return;
      let added = 0;
      for (const p of presets.value) {
        if (mirrors.value.some((m) => urlKey(m.url) === urlKey(p.url))) continue;
        data.value.settings.mirrors.push({ ...p });
        added++;
      }
      if (added) await save();
    }
    function removeMirror(idx) {
      data.value?.settings.mirrors.splice(idx, 1);
    }
    const dragFrom = ref(-1);
    const dragOver = ref(-1);
    let dragMoved = false;
    function moveMirror(from, to) {
      const list = data.value?.settings.mirrors;
      if (!list) return;
      if (from === to || from < 0 || to < 0 || from >= list.length || to >= list.length) return;
      const [row] = list.splice(from, 1);
      if (!row) return;
      list.splice(to, 0, row);
      dragFrom.value = to;
      dragOver.value = to;
      dragMoved = true;
    }
    function startDrag(i, ev) {
      if (mirrors.value.length < 2) return;
      ev.preventDefault();
      dragFrom.value = i;
      dragOver.value = i;
      dragMoved = false;
      window.addEventListener("pointermove", onDragMove);
      window.addEventListener("pointerup", endDrag);
      window.addEventListener("pointercancel", endDrag);
    }
    function onDragMove(ev) {
      const from = dragFrom.value;
      if (from < 0) return;
      const rows = Array.from(document.querySelectorAll("[data-mirror-row]"));
      for (let k = 0; k < rows.length; k++) {
        if (k === from) continue;
        const el = rows[k];
        if (!el) continue;
        const rect = el.getBoundingClientRect();
        const mid = rect.top + rect.height / 2;
        if (k < from && ev.clientY < mid) return moveMirror(from, k);
        if (k > from && ev.clientY > mid) return moveMirror(from, k);
      }
    }
    function endDrag() {
      dragFrom.value = -1;
      dragOver.value = -1;
      window.removeEventListener("pointermove", onDragMove);
      window.removeEventListener("pointerup", endDrag);
      window.removeEventListener("pointercancel", endDrag);
      if (dragMoved) toast.success("顺序已调整", "点「保存」后生效");
      dragMoved = false;
    }
    function onGripKey(i, ev) {
      if (ev.key !== "ArrowUp" && ev.key !== "ArrowDown") return;
      ev.preventDefault();
      const to = ev.key === "ArrowUp" ? i - 1 : i + 1;
      if (to < 0 || to >= mirrors.value.length) return;
      moveMirror(i, to);
      toast.success("顺序已调整", "点「保存」后生效");
      dragMoved = false;
    }
    const savedSig = ref("");
    function sigOf(s) {
      if (!s) return "";
      return JSON.stringify({
        mirrors: s.mirrors.map((m) => [m.url, m.note, m.enabled]),
        pullMirror: s.pullMirror,
        insecure: s.insecure
      });
    }
    const dirty = computed(() => !!data.value && sigOf(data.value.settings) !== savedSig.value);
    function applyResults(list, results) {
      for (const r of results) {
        const hit = list.find((x) => urlKey(x.url) === urlKey(r.url));
        if (!hit) continue;
        hit.latencyMs = r.latencyMs;
        hit.ok = r.ok;
        hit.err = r.err;
        hit.lastTested = r.lastTested;
      }
    }
    async function testOne(url) {
      testingUrl.value = url;
      try {
        const res = await api.post("/api/registries/test", { url });
        applyResults(mirrors.value, [res]);
        if (res.ok) {
          toast.success(`${res.url} 可用`, `延迟 ${res.latencyMs} ms`);
        } else {
          toast.error(`${res.url} 不可用`, res.err || "连接失败");
        }
      } catch (e) {
        toast.error("测速失败", e instanceof Error ? e.message : String(e));
      } finally {
        testingUrl.value = "";
      }
    }
    async function testAll() {
      testing.value = true;
      try {
        const urls = [...new Set(mirrors.value.map((m) => m.url))];
        const res = await api.post("/api/registries/test-all", { urls });
        applyResults(mirrors.value, res.results ?? []);
        toast.success("测速完成", "结果已写回各自那一行");
      } catch (e) {
        toast.error("批量测速失败", e instanceof Error ? e.message : String(e));
      } finally {
        testing.value = false;
      }
    }
    function focusAdd() {
      document.getElementById("mirror-add")?.scrollIntoView({ behavior: "smooth", block: "center" });
      document.getElementById("mirror-add")?.querySelector("input")?.focus();
    }
    async function copySnippet() {
      const text = data.value?.snippet ?? "";
      try {
        await navigator.clipboard.writeText(text);
        copied.value = true;
        window.setTimeout(() => copied.value = false, 1800);
        toast.success("已复制到剪贴板");
      } catch {
        toast.error("复制失败", "请手动选中代码块复制");
      }
    }
    const latencyTone = (m) => {
      if (!m.lastTested) return "dh-badge-plain";
      if (!m.ok) return "dh-badge-err";
      if ((m.latencyMs ?? 9999) < 300) return "dh-badge-run";
      if ((m.latencyMs ?? 9999) < 1200) return "dh-badge-warn";
      return "dh-badge-plain";
    };
    function connText(m) {
      if (!m.lastTested) return "未测速";
      if (m.ok) return `正常 · ${m.latencyMs ?? 0} ms`;
      return shortErr(m.err);
    }
    function shortErr(err) {
      const s = (err ?? "").trim();
      if (!s) return "连接失败";
      if (/no such host|server misbehaving|lookup .* on /i.test(s)) return "域名解析失败";
      if (/deadline exceeded|timed out|timeout/i.test(s)) return "连接超时";
      if (/connection refused/i.test(s)) return "拒绝连接";
      if (/certificate|x509|tls:/i.test(s)) return "证书错误";
      if (/connection reset|unexpected EOF/i.test(s)) return "连接被重置";
      const http = s.match(/HTTP\s*(\d{3})/);
      if (http) {
        const code = Number(http[1]);
        if (code === 401 || code === 403) return `${code} 拒绝访问`;
        if (code === 404) return "404 未找到";
        if (code >= 500) return `${code} 服务异常`;
        return `HTTP ${code}`;
      }
      return s.length > 18 ? `${s.slice(0, 18)}…` : s;
    }
    const policy = ref(null);
    async function loadPolicy() {
      try {
        policy.value = await api.get("/api/settings");
      } catch {
        policy.value = null;
      }
    }
    onMounted(() => {
      void load();
      void loadPolicy();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[5] || (_cache[5] = createBaseVNode("div", { class: "dh-h1" }, "镜像加速源", -1)),
          _cache[6] || (_cache[6] = createBaseVNode("div", { class: "dh-sub" }, "按顺序优先使用，失败自动顺延下一个", -1)),
          createBaseVNode("div", _hoisted_3, [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: testing.value,
              onClick: testAll
            }, [
              createVNode(unref(Gauge), {
                class: normalizeClass(["h-3.5 w-3.5", testing.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[4] || (_cache[4] = createTextVNode("测试全部 ", -1))
            ], 8, _hoisted_4),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              onClick: focusAdd
            }, "添加加速源")
          ])
        ]),
        createBaseVNode("div", _hoisted_5, [
          createBaseVNode("div", _hoisted_6, [
            createVNode(unref(Info), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[7] || (_cache[7] = createBaseVNode("span", null, "Docker 守护进程当前生效的加速源", -1)),
            _cache[8] || (_cache[8] = createBaseVNode("span", { class: "ml-auto text-[11.5px] font-normal text-text-5" }, "来源：daemon.json", -1))
          ]),
          createBaseVNode("div", _hoisted_7, [
            createBaseVNode("div", _hoisted_8, toDisplayString(data.value?.explain), 1),
            data.value?.daemonError ? (openBlock(), createElementBlock("div", _hoisted_9, " 无法读取守护进程信息：" + toDisplayString(data.value.daemonError), 1)) : daemonMirrors.value.length ? (openBlock(), createElementBlock("div", _hoisted_10, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(daemonMirrors.value, (m) => {
                return openBlock(), createElementBlock("span", {
                  key: m,
                  class: "rounded-md bg-ink-800 px-2 py-[3px] font-mono text-[11.5px] text-accent"
                }, toDisplayString(m), 1);
              }), 128))
            ])) : (openBlock(), createElementBlock("div", _hoisted_11, [
              createVNode(unref(TriangleAlert), { class: "h-3.5 w-3.5" }),
              _cache[9] || (_cache[9] = createTextVNode(" 守护进程没有配置任何加速源，拉取镜像会直连 Docker Hub。 ", -1))
            ]))
          ])
        ]),
        createBaseVNode("div", _hoisted_12, [
          createBaseVNode("div", _hoisted_13, [
            createBaseVNode("div", _hoisted_14, [
              createVNode(unref(Rocket), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[11] || (_cache[11] = createBaseVNode("span", null, "我的加速源", -1)),
              createBaseVNode("div", _hoisted_15, [
                dirty.value ? (openBlock(), createElementBlock("span", _hoisted_16, "有未保存的改动")) : createCommentVNode("", true),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm",
                  disabled: testing.value,
                  onClick: testAll
                }, [
                  createVNode(unref(Gauge), {
                    class: normalizeClass(["h-3 w-3", testing.value ? "dh-spin" : ""])
                  }, null, 8, ["class"]),
                  _cache[10] || (_cache[10] = createTextVNode("批量测速 ", -1))
                ], 8, _hoisted_17),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm dh-btn-primary",
                  disabled: saving.value,
                  onClick: save
                }, "保存", 8, _hoisted_18)
              ])
            ]),
            createBaseVNode("div", _hoisted_19, [
              createBaseVNode("div", _hoisted_20, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => newUrl.value = $event),
                  class: "dh-input flex-1 !min-w-[180px]",
                  placeholder: "https://你的加速站地址"
                }, null, 512), [
                  [vModelText, newUrl.value]
                ]),
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => newNote.value = $event),
                  class: "dh-input !w-[110px]",
                  placeholder: "备注"
                }, null, 512), [
                  [vModelText, newNote.value]
                ]),
                createBaseVNode("button", {
                  class: "dh-btn",
                  disabled: !newUrl.value.trim(),
                  onClick: addMirror
                }, [
                  createVNode(unref(Plus), { class: "h-3.5 w-3.5" }),
                  _cache[12] || (_cache[12] = createTextVNode("添加 ", -1))
                ], 8, _hoisted_21)
              ])
            ]),
            !mirrors.value.length ? (openBlock(), createBlock(_sfc_main$1, {
              key: 0,
              icon: unref(Rocket),
              title: loading.value ? "正在载入…" : "还没有添加加速源",
              description: loading.value ? "" : "预置的常用加速站随首次启动就已经写进来了，不想要哪条直接删；这里被删空之后可以一键找回，也可以手动填写自建加速站地址。",
              "action-label": !loading.value && presets.value.length ? "加入预置的常用加速站" : "",
              onAction: addPresets
            }, null, 8, ["icon", "title", "description", "action-label"])) : (openBlock(), createElementBlock("div", _hoisted_22, [
              createBaseVNode("table", _hoisted_23, [
                _cache[14] || (_cache[14] = createBaseVNode("thead", null, [
                  createBaseVNode("tr", null, [
                    createBaseVNode("th", { class: "w-[38px]" }),
                    createBaseVNode("th", null, "加速源地址"),
                    createBaseVNode("th", { class: "w-[175px]" }, "连通性"),
                    createBaseVNode("th", { class: "w-[130px]" }, "上次测速"),
                    createBaseVNode("th", { class: "w-[70px] text-right" }, "启用"),
                    createBaseVNode("th", { class: "w-[130px] text-right" }, "操作")
                  ])
                ], -1)),
                createBaseVNode("tbody", null, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(mirrors.value, (m, i) => {
                    return openBlock(), createElementBlock("tr", {
                      key: m.url,
                      "data-mirror-row": "",
                      class: normalizeClass({
                        "dh-drag-row": dragFrom.value === i,
                        "dh-drag-over": dragOver.value === i && dragFrom.value >= 0 && dragFrom.value !== i
                      })
                    }, [
                      createBaseVNode("td", _hoisted_24, [
                        createBaseVNode("button", {
                          type: "button",
                          class: "dh-grip",
                          "aria-label": `拖动调整 ${m.url} 的优先级，也可用上下方向键`,
                          title: "拖动调整优先级（也可用 ↑ ↓）",
                          onPointerdown: ($event) => startDrag(i, $event),
                          onKeydown: ($event) => onGripKey(i, $event)
                        }, [
                          createVNode(unref(GripVertical), { class: "h-3.5 w-3.5" })
                        ], 40, _hoisted_25)
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("div", {
                          class: "truncate font-mono text-[12px] text-text-2",
                          title: m.url
                        }, toDisplayString(m.url), 9, _hoisted_26),
                        createBaseVNode("div", _hoisted_27, [
                          m.builtin ? (openBlock(), createElementBlock("span", _hoisted_28, "预置")) : createCommentVNode("", true),
                          m.note ? (openBlock(), createElementBlock("span", _hoisted_29, toDisplayString(m.note), 1)) : createCommentVNode("", true)
                        ])
                      ]),
                      createBaseVNode("td", null, [
                        createBaseVNode("span", {
                          class: normalizeClass(["dh-badge", latencyTone(m)])
                        }, toDisplayString(connText(m)), 3)
                      ]),
                      createBaseVNode("td", _hoisted_30, toDisplayString(m.lastTested ? unref(formatDayTime)(m.lastTested) : "—"), 1),
                      createBaseVNode("td", _hoisted_31, [
                        createVNode(_sfc_main$3, {
                          "model-value": m.enabled,
                          label: `启用 ${m.url}`,
                          "onUpdate:modelValue": (v) => m.enabled = v
                        }, null, 8, ["model-value", "label", "onUpdate:modelValue"])
                      ]),
                      createBaseVNode("td", _hoisted_32, [
                        createBaseVNode("button", {
                          class: "dh-link",
                          disabled: testingUrl.value === m.url,
                          onClick: ($event) => testOne(m.url)
                        }, toDisplayString(testingUrl.value === m.url ? "测速中…" : "测速"), 9, _hoisted_33),
                        _cache[13] || (_cache[13] = createBaseVNode("span", { class: "mx-1.5 text-text-6" }, "·", -1)),
                        createBaseVNode("button", {
                          class: "dh-link dh-link-danger",
                          onClick: ($event) => removeMirror(i)
                        }, "删除", 8, _hoisted_34)
                      ])
                    ], 2);
                  }), 128))
                ])
              ]),
              _cache[15] || (_cache[15] = createBaseVNode("div", { class: "border-t border-line-row px-4 py-2 text-[11px] text-text-6" }, " 拖动左侧把手调整优先级：从上到下依次尝试，前面失败会自动顺延下一个。 ", -1))
            ]))
          ]),
          createBaseVNode("div", _hoisted_35, [
            createBaseVNode("div", _hoisted_36, [
              createVNode(unref(Zap), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[16] || (_cache[16] = createBaseVNode("span", null, "Dockhelm 自己的拉取策略", -1))
            ]),
            createBaseVNode("div", _hoisted_37, [
              createBaseVNode("div", null, [
                _cache[18] || (_cache[18] = createBaseVNode("label", { class: "dh-label" }, "拉取加速源", -1)),
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => pullMirror.value = $event),
                  class: "dh-select"
                }, [
                  _cache[17] || (_cache[17] = createBaseVNode("option", { value: "" }, "不指定（完全交给 Docker 守护进程）", -1)),
                  (openBlock(true), createElementBlock(Fragment, null, renderList(mirrors.value, (m) => {
                    return openBlock(), createElementBlock("option", {
                      key: m.url,
                      value: m.url
                    }, toDisplayString(m.url), 9, _hoisted_38);
                  }), 128))
                ], 512), [
                  [vModelSelect, pullMirror.value]
                ]),
                _cache[19] || (_cache[19] = createBaseVNode("div", { class: "mt-1.5 text-[11.5px] leading-relaxed text-text-5" }, [
                  createTextVNode(" 保持「不指定」时，拉取行为与手动执行 "),
                  createBaseVNode("code", { class: "text-text-3" }, "docker pull"),
                  createTextVNode(" 完全一致。 如果你没法改 daemon.json，可以在这里指定一个加速源：Dockhelm 会用 "),
                  createBaseVNode("code", { class: "text-text-3" }, "<加速站>/<仓库>:<标签>"),
                  createTextVNode(" 拉取， 拉完再打回原始标签，compose 与其它工具仍然按原来的名字找得到镜像。 "),
                  createBaseVNode("br"),
                  createBaseVNode("b", { class: "text-text-3" }, "注意"),
                  createTextVNode("：只对来自 Docker Hub 的镜像生效。 ")
                ], -1))
              ]),
              policy.value ? (openBlock(), createElementBlock("div", _hoisted_39, [
                createVNode(_sfc_main$2, {
                  title: "显式域名优先，不套用加速",
                  sub: "ghcr.io、私有仓库等已经写明域名的镜像直连 —— 加速站通常只镜像 Docker Hub"
                }, {
                  default: withCtx(() => [
                    createVNode(_sfc_main$4, {
                      modelValue: policy.value.directFirst,
                      "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => policy.value.directFirst = $event),
                      label: "显式域名优先"
                    }, null, 8, ["modelValue"])
                  ]),
                  _: 1
                })
              ])) : createCommentVNode("", true),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-primary",
                disabled: saving.value,
                onClick: save
              }, [
                saving.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "h-3.5 w-3.5 dh-spin"
                })) : (openBlock(), createBlock(unref(Save), {
                  key: 1,
                  class: "h-3.5 w-3.5"
                })),
                _cache[20] || (_cache[20] = createTextVNode("保存设置 ", -1))
              ], 8, _hoisted_40)
            ])
          ]),
          createBaseVNode("div", _hoisted_41, [
            createBaseVNode("div", _hoisted_42, [
              createVNode(unref(Copy), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[21] || (_cache[21] = createBaseVNode("span", null, "daemon.json 片段", -1)),
              createBaseVNode("button", {
                class: "dh-btn dh-btn-sm ml-auto",
                onClick: copySnippet
              }, [
                (openBlock(), createBlock(resolveDynamicComponent(copied.value ? unref(CircleCheck) : unref(Copy)), { class: "h-3 w-3" })),
                createTextVNode(" " + toDisplayString(copied.value ? "已复制" : "复制"), 1)
              ])
            ]),
            createBaseVNode("div", _hoisted_43, [
              _cache[22] || (_cache[22] = createBaseVNode("div", { class: "text-[11.5px] leading-relaxed text-text-5" }, [
                createTextVNode(" 把这段贴进 NAS 上的 "),
                createBaseVNode("code", { class: "text-text-3" }, "/etc/docker/daemon.json"),
                createTextVNode("（群晖在 Docker 套件的设置里，或 Container Manager 的「注册表镜像」），重启 Docker 服务后生效。 内容是上表里"),
                createBaseVNode("b", { class: "text-text-3" }, "已启用"),
                createTextVNode("的加速源，顺序与上表一致。 ")
              ], -1)),
              createBaseVNode("pre", _hoisted_44, toDisplayString(data.value?.snippet), 1)
            ])
          ])
        ])
      ]);
    };
  }
});
export {
  _sfc_main as default
};
