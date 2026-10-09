import { q as createLucideIcon, d as defineComponent, C as useToastStore, o as onMounted, c as createElementBlock, b as createBaseVNode, e as createVNode, n as normalizeClass, f as unref, T as Gauge, i as createTextVNode, O as Info, t as toDisplayString, F as Fragment, r as renderList, U as Rocket, w as withDirectives, G as vModelText, y as createBlock, M as vModelSelect, N as resolveDynamicComponent, P as CircleCheck, h as createCommentVNode, I as withCtx, R as RefreshCw, k as ref, l as computed, z as api, p as openBlock, V as CircleX } from "./index-HPgVGrf3.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-BG9jdsHP.js";
import { _ as _sfc_main$2, a as _sfc_main$3 } from "./ToggleSwitch.vue_vue_type_script_setup_true_lang-DegosQKp.js";
import { T as TriangleAlert } from "./triangle-alert-CI3k_cUc.js";
import { P as Plus } from "./plus-BTTUntJL.js";
import { Z as Zap } from "./zap-CV2RIe0A.js";
import { L as LoaderCircle } from "./loader-circle-C-XNYId3.js";
import { S as Save } from "./save-Dvb5gkQY.js";
import { T as Trash2 } from "./trash-2-C3mj8LF4.js";
const Copy = createLucideIcon("CopyIcon", [
  ["rect", { width: "14", height: "14", x: "8", y: "8", rx: "2", ry: "2", key: "17jyea" }],
  ["path", { d: "M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2", key: "zix9uf" }]
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
  class: "text-[12px] text-[#fca5a5]"
};
const _hoisted_10 = {
  key: 1,
  class: "flex flex-wrap gap-1.5"
};
const _hoisted_11 = {
  key: 2,
  class: "flex items-center gap-2 text-[12px] text-text-5"
};
const _hoisted_12 = { class: "grid grid-cols-1 gap-3.5 xl:grid-cols-[1.3fr_1fr]" };
const _hoisted_13 = { class: "dh-card" };
const _hoisted_14 = { class: "dh-card-head" };
const _hoisted_15 = { class: "ml-auto flex gap-2" };
const _hoisted_16 = ["disabled"];
const _hoisted_17 = ["disabled"];
const _hoisted_18 = {
  id: "mirror-add",
  class: "flex flex-col gap-2 border-b border-line-1 p-3"
};
const _hoisted_19 = { class: "flex flex-wrap gap-2" };
const _hoisted_20 = ["disabled"];
const _hoisted_21 = {
  key: 1,
  class: "flex flex-col"
};
const _hoisted_22 = ["checked", "onChange"];
const _hoisted_23 = { class: "min-w-0 flex-1" };
const _hoisted_24 = { class: "truncate font-mono text-[11.5px] text-text-2" };
const _hoisted_25 = {
  key: 0,
  class: "text-[11px] text-text-5"
};
const _hoisted_26 = ["disabled", "onClick"];
const _hoisted_27 = ["onClick"];
const _hoisted_28 = { class: "flex flex-col gap-3.5" };
const _hoisted_29 = { class: "dh-card" };
const _hoisted_30 = { class: "dh-card-head" };
const _hoisted_31 = { class: "dh-card-body flex flex-col gap-3" };
const _hoisted_32 = ["value"];
const _hoisted_33 = ["disabled"];
const _hoisted_34 = { class: "dh-card" };
const _hoisted_35 = { class: "dh-card-head" };
const _hoisted_36 = { class: "dh-card-body" };
const _hoisted_37 = { class: "mt-2.5 overflow-x-auto rounded-[9px] border border-line-3 bg-ink-800 px-3 py-2.5 font-mono text-[11.5px] leading-[1.8] text-text-2" };
const _hoisted_38 = {
  key: 0,
  class: "dh-card"
};
const _hoisted_39 = { class: "dh-card-head" };
const _hoisted_40 = { class: "grid grid-cols-1 gap-2 p-3 md:grid-cols-2 xl:grid-cols-3" };
const _hoisted_41 = { class: "min-w-0 flex-1" };
const _hoisted_42 = { class: "truncate font-mono text-[11.5px] text-text-2" };
const _hoisted_43 = { class: "text-[11px] text-text-5" };
const _hoisted_44 = ["onClick"];
const _hoisted_45 = ["onClick"];
const _hoisted_46 = {
  key: 1,
  class: "grid grid-cols-1 gap-3.5 xl:grid-cols-2"
};
const _hoisted_47 = { class: "dh-card" };
const _hoisted_48 = { class: "dh-card-head" };
const _hoisted_49 = ["disabled"];
const _hoisted_50 = { class: "flex flex-col gap-3 p-3.5" };
const _hoisted_51 = { class: "dh-card" };
const _hoisted_52 = { class: "dh-card-head" };
const _hoisted_53 = { class: "flex flex-col gap-3 p-3.5" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "RegistriesView",
  setup(__props) {
    const toast = useToastStore();
    const data = ref(null);
    const loading = ref(true);
    const testing = ref(false);
    const testingUrl = ref("");
    const newUrl = ref("");
    const newNote = ref("");
    const copied = ref(false);
    const saving = ref(false);
    const mirrors = computed(() => data.value?.settings.mirrors ?? []);
    const suggestions = computed(() => data.value?.suggestions ?? []);
    const daemonMirrors = computed(() => data.value?.daemonMirrors ?? []);
    const pullMirror = computed({
      get: () => data.value?.settings.pullMirror ?? "",
      set: (v) => {
        if (data.value) data.value.settings.pullMirror = v;
      }
    });
    async function load() {
      loading.value = true;
      try {
        data.value = await api.get("/api/registries");
      } catch (e) {
        toast.error("读取加速源配置失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function save() {
      if (!data.value) return;
      saving.value = true;
      try {
        await api.put("/api/registries", { settings: data.value.settings });
        toast.success("已保存");
        await load();
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      } finally {
        saving.value = false;
      }
    }
    function addMirror() {
      const url = newUrl.value.trim();
      if (!url || !data.value) return;
      data.value.settings.mirrors.push({ url, note: newNote.value.trim(), enabled: true });
      newUrl.value = "";
      newNote.value = "";
    }
    function addSuggestion(s) {
      if (!data.value) return;
      data.value.settings.mirrors.push({ ...s, enabled: true });
      data.value.suggestions = data.value.suggestions.filter((x) => x.url !== s.url);
    }
    function removeMirror(idx) {
      data.value?.settings.mirrors.splice(idx, 1);
    }
    async function testOne(url) {
      testingUrl.value = url;
      try {
        const res = await api.post("/api/registries/test", { url });
        if (res.ok) {
          toast.success(`${res.url} 可用`, `延迟 ${res.latencyMs} ms`);
        } else {
          toast.error(`${res.url} 不可用`, res.err || "连接失败");
        }
        await load();
      } catch (e) {
        toast.error("测速失败", e instanceof Error ? e.message : String(e));
      } finally {
        testingUrl.value = "";
      }
    }
    async function testAll() {
      testing.value = true;
      try {
        const urls = [
          ...mirrors.value.map((m) => m.url),
          ...suggestions.value.map((s) => s.url)
        ];
        await api.post("/api/registries/test-all", { urls });
        toast.success("测速完成", "结果已按延迟排序并保存");
        await load();
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
      if (m.ok === void 0) return "dh-badge-plain";
      if (!m.ok) return "dh-badge-err";
      if ((m.latencyMs ?? 9999) < 300) return "dh-badge-run";
      if ((m.latencyMs ?? 9999) < 1200) return "dh-badge-warn";
      return "dh-badge-plain";
    };
    const policy = ref(null);
    const savingPolicy = ref(false);
    async function loadPolicy() {
      try {
        policy.value = await api.get("/api/settings");
      } catch {
        policy.value = null;
      }
    }
    async function savePolicy() {
      if (!policy.value) return;
      savingPolicy.value = true;
      try {
        policy.value = await api.put("/api/settings", policy.value);
        toast.success("已保存", "拉取与检测设置立即生效");
      } catch (e) {
        toast.error("保存失败", e instanceof Error ? e.message : String(e));
      } finally {
        savingPolicy.value = false;
      }
    }
    onMounted(() => {
      void load();
      void loadPolicy();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[8] || (_cache[8] = createBaseVNode("div", { class: "dh-h1" }, "镜像加速源", -1)),
          _cache[9] || (_cache[9] = createBaseVNode("div", { class: "dh-sub" }, "按顺序优先使用，失败自动顺延下一个", -1)),
          createBaseVNode("div", _hoisted_3, [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: testing.value,
              onClick: testAll
            }, [
              createVNode(unref(Gauge), {
                class: normalizeClass(["h-3.5 w-3.5", testing.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[7] || (_cache[7] = createTextVNode("测试全部 ", -1))
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
            _cache[10] || (_cache[10] = createBaseVNode("span", null, "Docker 守护进程当前生效的加速源", -1)),
            _cache[11] || (_cache[11] = createBaseVNode("span", { class: "ml-auto text-[11.5px] font-normal text-text-5" }, "来源：daemon.json", -1))
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
              _cache[12] || (_cache[12] = createTextVNode(" 守护进程没有配置任何加速源，拉取镜像会直连 Docker Hub。 ", -1))
            ]))
          ])
        ]),
        createBaseVNode("div", _hoisted_12, [
          createBaseVNode("div", _hoisted_13, [
            createBaseVNode("div", _hoisted_14, [
              createVNode(unref(Rocket), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[14] || (_cache[14] = createBaseVNode("span", null, "我的加速源", -1)),
              createBaseVNode("div", _hoisted_15, [
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm",
                  disabled: testing.value,
                  onClick: testAll
                }, [
                  createVNode(unref(Gauge), {
                    class: normalizeClass(["h-3 w-3", testing.value ? "dh-spin" : ""])
                  }, null, 8, ["class"]),
                  _cache[13] || (_cache[13] = createTextVNode("批量测速 ", -1))
                ], 8, _hoisted_16),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm dh-btn-primary",
                  disabled: saving.value,
                  onClick: save
                }, "保存", 8, _hoisted_17)
              ])
            ]),
            createBaseVNode("div", _hoisted_18, [
              createBaseVNode("div", _hoisted_19, [
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
                  _cache[15] || (_cache[15] = createTextVNode("添加 ", -1))
                ], 8, _hoisted_20)
              ])
            ]),
            !mirrors.value.length ? (openBlock(), createBlock(_sfc_main$1, {
              key: 0,
              icon: unref(Rocket),
              title: loading.value ? "正在载入…" : "还没有添加加速源",
              description: "可以从下方的推荐列表一键添加，也可以手动填写自建加速站地址。"
            }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_21, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(mirrors.value, (m, i) => {
                return openBlock(), createElementBlock("div", {
                  key: m.url,
                  class: "flex flex-wrap items-center gap-2 border-b border-[#171f2a] px-3 py-2.5 last:border-b-0"
                }, [
                  createBaseVNode("input", {
                    type: "checkbox",
                    class: "h-[14px] w-[14px] accent-[#2dd4bf]",
                    checked: m.enabled,
                    onChange: ($event) => m.enabled = !m.enabled,
                    title: "启用这个加速源"
                  }, null, 40, _hoisted_22),
                  createBaseVNode("div", _hoisted_23, [
                    createBaseVNode("div", _hoisted_24, toDisplayString(m.url), 1),
                    m.note ? (openBlock(), createElementBlock("div", _hoisted_25, toDisplayString(m.note), 1)) : createCommentVNode("", true)
                  ]),
                  m.ok !== void 0 ? (openBlock(), createElementBlock("span", {
                    key: 0,
                    class: normalizeClass(["dh-badge", latencyTone(m)])
                  }, [
                    m.ok ? (openBlock(), createBlock(unref(CircleCheck), {
                      key: 0,
                      class: "h-3 w-3"
                    })) : (openBlock(), createBlock(unref(CircleX), {
                      key: 1,
                      class: "h-3 w-3"
                    })),
                    createTextVNode(" " + toDisplayString(m.ok ? m.latencyMs + " ms" : "不可用"), 1)
                  ], 2)) : createCommentVNode("", true),
                  createBaseVNode("button", {
                    class: "dh-btn dh-btn-sm",
                    disabled: testingUrl.value === m.url,
                    onClick: ($event) => testOne(m.url)
                  }, [
                    createVNode(unref(Zap), {
                      class: normalizeClass(["h-3 w-3", testingUrl.value === m.url ? "dh-spin" : ""])
                    }, null, 8, ["class"]),
                    _cache[16] || (_cache[16] = createTextVNode("测速 ", -1))
                  ], 8, _hoisted_26),
                  createBaseVNode("button", {
                    class: "dh-btn dh-btn-sm dh-btn-danger",
                    onClick: ($event) => removeMirror(i)
                  }, [
                    createVNode(unref(Trash2), { class: "h-3 w-3" })
                  ], 8, _hoisted_27)
                ]);
              }), 128))
            ]))
          ]),
          createBaseVNode("div", _hoisted_28, [
            createBaseVNode("div", _hoisted_29, [
              createBaseVNode("div", _hoisted_30, [
                createVNode(unref(Zap), { class: "h-3.5 w-3.5 text-text-4" }),
                _cache[17] || (_cache[17] = createBaseVNode("span", null, "Dockhelm 自己的拉取策略", -1))
              ]),
              createBaseVNode("div", _hoisted_31, [
                createBaseVNode("div", null, [
                  _cache[19] || (_cache[19] = createBaseVNode("label", { class: "dh-label" }, "拉取加速源", -1)),
                  withDirectives(createBaseVNode("select", {
                    "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => pullMirror.value = $event),
                    class: "dh-select"
                  }, [
                    _cache[18] || (_cache[18] = createBaseVNode("option", { value: "" }, "不指定（完全交给 Docker 守护进程）", -1)),
                    (openBlock(true), createElementBlock(Fragment, null, renderList(mirrors.value, (m) => {
                      return openBlock(), createElementBlock("option", {
                        key: m.url,
                        value: m.url
                      }, toDisplayString(m.url), 9, _hoisted_32);
                    }), 128))
                  ], 512), [
                    [vModelSelect, pullMirror.value]
                  ]),
                  _cache[20] || (_cache[20] = createBaseVNode("div", { class: "mt-1.5 text-[11.5px] leading-relaxed text-text-5" }, [
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
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-primary",
                  disabled: saving.value,
                  onClick: save
                }, "保存设置", 8, _hoisted_33)
              ])
            ]),
            createBaseVNode("div", _hoisted_34, [
              createBaseVNode("div", _hoisted_35, [
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
              createBaseVNode("div", _hoisted_36, [
                _cache[22] || (_cache[22] = createBaseVNode("div", { class: "text-[11.5px] leading-relaxed text-text-5" }, [
                  createTextVNode(" 把这段贴进 NAS 上的 "),
                  createBaseVNode("code", { class: "text-text-3" }, "/etc/docker/daemon.json"),
                  createTextVNode("（群晖在 Docker 套件的设置里，或 Container Manager 的「注册表镜像」），重启 Docker 服务后生效。 内容是上表里"),
                  createBaseVNode("b", { class: "text-text-3" }, "勾选启用"),
                  createTextVNode("的加速源。 ")
                ], -1)),
                createBaseVNode("pre", _hoisted_37, toDisplayString(data.value?.snippet), 1)
              ])
            ])
          ])
        ]),
        suggestions.value.length ? (openBlock(), createElementBlock("div", _hoisted_38, [
          createBaseVNode("div", _hoisted_39, [
            createVNode(unref(Rocket), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[23] || (_cache[23] = createBaseVNode("span", null, "推荐加速源", -1)),
            _cache[24] || (_cache[24] = createBaseVNode("span", { class: "ml-auto text-[11.5px] font-normal text-text-5" }, " 点「添加」加入列表；这些只是起点，可用性请用测速确认 ", -1))
          ]),
          createBaseVNode("div", _hoisted_40, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(suggestions.value, (s) => {
              return openBlock(), createElementBlock("div", {
                key: s.url,
                class: "flex items-center gap-2 rounded-[10px] border border-line-1 bg-ink-800 px-3 py-2"
              }, [
                createBaseVNode("div", _hoisted_41, [
                  createBaseVNode("div", _hoisted_42, toDisplayString(s.url), 1),
                  createBaseVNode("div", _hoisted_43, toDisplayString(s.note), 1)
                ]),
                s.ok !== void 0 ? (openBlock(), createElementBlock("span", {
                  key: 0,
                  class: normalizeClass(["dh-badge", latencyTone(s)])
                }, toDisplayString(s.ok ? s.latencyMs + " ms" : "不可用"), 3)) : createCommentVNode("", true),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm",
                  onClick: ($event) => testOne(s.url)
                }, [
                  createVNode(unref(Zap), { class: "h-3 w-3" })
                ], 8, _hoisted_44),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm",
                  onClick: ($event) => addSuggestion(s)
                }, [
                  createVNode(unref(Plus), { class: "h-3 w-3" }),
                  _cache[25] || (_cache[25] = createTextVNode("添加 ", -1))
                ], 8, _hoisted_45)
              ]);
            }), 128))
          ])
        ])) : createCommentVNode("", true),
        policy.value ? (openBlock(), createElementBlock("div", _hoisted_46, [
          createBaseVNode("div", _hoisted_47, [
            createBaseVNode("div", _hoisted_48, [
              createVNode(unref(Gauge), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[27] || (_cache[27] = createBaseVNode("span", null, "拉取与更新", -1)),
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
                _cache[26] || (_cache[26] = createTextVNode("保存 ", -1))
              ], 8, _hoisted_49)
            ]),
            createBaseVNode("div", _hoisted_50, [
              createVNode(_sfc_main$2, {
                title: "显式域名优先，不套用加速",
                sub: "ghcr.io、私有仓库等直连 —— 加速站通常只镜像 Docker Hub"
              }, {
                default: withCtx(() => [
                  createVNode(_sfc_main$3, {
                    modelValue: policy.value.directFirst,
                    "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => policy.value.directFirst = $event),
                    label: "显式域名优先"
                  }, null, 8, ["modelValue"])
                ]),
                _: 1
              }),
              createVNode(_sfc_main$2, {
                title: "并发更新",
                sub: "串行最稳，可调 1–8（越高越容易触发仓库限流）"
              }, {
                default: withCtx(() => [
                  withDirectives(createBaseVNode("select", {
                    "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => policy.value.concurrency = $event),
                    class: "dh-select !w-[130px] !py-[5px] !text-[11.5px]"
                  }, [..._cache[28] || (_cache[28] = [
                    createBaseVNode("option", { value: 1 }, "1（串行）", -1),
                    createBaseVNode("option", { value: 2 }, "2（推荐）", -1),
                    createBaseVNode("option", { value: 3 }, "3", -1),
                    createBaseVNode("option", { value: 4 }, "4", -1),
                    createBaseVNode("option", { value: 6 }, "6", -1),
                    createBaseVNode("option", { value: 8 }, "8", -1)
                  ])], 512), [
                    [
                      vModelSelect,
                      policy.value.concurrency,
                      void 0,
                      { number: true }
                    ]
                  ])
                ]),
                _: 1
              }),
              createVNode(_sfc_main$2, {
                title: "更新前保留原容器",
                sub: "重建期间旧容器改名保留，失败时以原名与状态复活"
              }, {
                default: withCtx(() => [..._cache[29] || (_cache[29] = [
                  createBaseVNode("span", { class: "dh-badge dh-badge-plain" }, "始终开启", -1)
                ])]),
                _: 1
              })
            ])
          ]),
          createBaseVNode("div", _hoisted_51, [
            createBaseVNode("div", _hoisted_52, [
              createVNode(unref(RefreshCw), { class: "h-3.5 w-3.5 text-text-4" }),
              _cache[30] || (_cache[30] = createBaseVNode("span", null, "检测", -1))
            ]),
            createBaseVNode("div", _hoisted_53, [
              createVNode(_sfc_main$2, {
                title: "检测频率",
                sub: policy.value.checkIntervalHours > 0 ? `每 ${policy.value.checkIntervalHours} 小时自动扫描一次` : "关闭后只在手动点「重新检测」时才检查"
              }, {
                default: withCtx(() => [
                  withDirectives(createBaseVNode("select", {
                    "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => policy.value.checkIntervalHours = $event),
                    class: "dh-select !w-[130px] !py-[5px] !text-[11.5px]"
                  }, [..._cache[31] || (_cache[31] = [
                    createBaseVNode("option", { value: 0 }, "关闭", -1),
                    createBaseVNode("option", { value: 1 }, "每 1 小时", -1),
                    createBaseVNode("option", { value: 3 }, "每 3 小时", -1),
                    createBaseVNode("option", { value: 6 }, "每 6 小时", -1),
                    createBaseVNode("option", { value: 12 }, "每 12 小时", -1),
                    createBaseVNode("option", { value: 24 }, "每 24 小时", -1)
                  ])], 512), [
                    [
                      vModelSelect,
                      policy.value.checkIntervalHours,
                      void 0,
                      { number: true }
                    ]
                  ])
                ]),
                _: 1
              }, 8, ["sub"]),
              createVNode(_sfc_main$2, {
                title: "摘要来源",
                sub: "由本机 Docker 守护进程解析，与 docker pull 同源"
              }, {
                default: withCtx(() => [..._cache[32] || (_cache[32] = [
                  createBaseVNode("span", { class: "dh-badge dh-badge-plain" }, "推荐", -1)
                ])]),
                _: 1
              }),
              createVNode(_sfc_main$2, {
                title: "检测完成后通知",
                sub: "支持企微 / Telegram / Bark / Webhook 等已配置渠道"
              }, {
                default: withCtx(() => [
                  createVNode(_sfc_main$3, {
                    modelValue: policy.value.notifyOnCheck,
                    "onUpdate:modelValue": _cache[6] || (_cache[6] = ($event) => policy.value.notifyOnCheck = $event),
                    label: "检测完成后通知"
                  }, null, 8, ["modelValue"])
                ]),
                _: 1
              })
            ])
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
export {
  _sfc_main as default
};
