import { x as createLucideIcon, d as defineComponent, D as useToastStore, a as onMounted, c as createElementBlock, e as createBaseVNode, t as toDisplayString, g as unref, f as createVNode, n as normalizeClass, R as RefreshCw, k as createTextVNode, S as Search, w as withDirectives, E as vModelText, l as vModelCheckbox, O as Layers, i as createBlock, F as Fragment, r as renderList, H as withCtx, m as ref, p as computed, B as api, s as openBlock } from "./index-D0P9SfJT.js";
import { a as formatBytes, b as relativeTime } from "./format-NZu8lNOP.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-lqZI6iCq.js";
import { _ as _sfc_main$2 } from "./Modal.vue_vue_type_script_setup_true_lang-BeECE25y.js";
import { T as Trash2 } from "./trash-2-BuikYSkM.js";
const Sparkles = createLucideIcon("SparklesIcon", [
  [
    "path",
    {
      d: "M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z",
      key: "4pj2yx"
    }
  ],
  ["path", { d: "M20 3v4", key: "1olli1" }],
  ["path", { d: "M22 5h-4", key: "1gvqau" }],
  ["path", { d: "M4 17v2", key: "vumght" }],
  ["path", { d: "M5 18H3", key: "zchphs" }]
]);
const _hoisted_1 = { class: "flex flex-col gap-3.5 p-[18px]" };
const _hoisted_2 = { class: "dh-phead" };
const _hoisted_3 = { class: "dh-sub" };
const _hoisted_4 = { class: "ml-auto flex gap-2" };
const _hoisted_5 = ["disabled"];
const _hoisted_6 = ["disabled"];
const _hoisted_7 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_8 = { class: "relative min-w-[170px] flex-1 sm:max-w-[260px]" };
const _hoisted_9 = { class: "flex cursor-pointer items-center gap-2 text-[12px] text-text-3" };
const _hoisted_10 = { class: "dh-card" };
const _hoisted_11 = { class: "dh-card-head" };
const _hoisted_12 = { class: "ml-auto text-[11.5px] font-normal text-text-5" };
const _hoisted_13 = {
  key: 1,
  class: "overflow-x-auto"
};
const _hoisted_14 = { class: "dh-table" };
const _hoisted_15 = {
  key: 0,
  class: "flex flex-wrap gap-1"
};
const _hoisted_16 = {
  key: 1,
  class: "flex flex-wrap items-center gap-1.5"
};
const _hoisted_17 = { class: "font-mono text-[11.5px] text-text-2" };
const _hoisted_18 = { class: "font-mono text-[11.5px] text-text-4" };
const _hoisted_19 = { class: "text-[12px] text-text-2" };
const _hoisted_20 = { class: "text-[11.5px] text-text-4" };
const _hoisted_21 = {
  key: 0,
  class: "dh-badge dh-badge-accent"
};
const _hoisted_22 = {
  key: 1,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_23 = ["disabled", "title", "onClick"];
const _hoisted_24 = ["disabled"];
const _hoisted_25 = ["disabled"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "ImagesView",
  setup(__props) {
    const toast = useToastStore();
    const images = ref([]);
    const totalSize = ref(0);
    const loading = ref(false);
    const pruning = ref(false);
    const removing = ref(false);
    const keyword = ref("");
    const onlyDangling = ref(false);
    const removeTarget = ref(null);
    const filtered = computed(() => {
      let list = images.value;
      const k = keyword.value.trim().toLowerCase();
      if (k) {
        list = list.filter(
          (i) => i.tags.some((t) => t.toLowerCase().includes(k)) || (i.repo ?? "").toLowerCase().includes(k) || i.shortId.includes(k)
        );
      }
      if (onlyDangling.value) list = list.filter((i) => i.dangling);
      return list;
    });
    const danglingCount = computed(() => images.value.filter((i) => i.dangling).length);
    const danglingSize = computed(
      () => images.value.filter((i) => i.dangling).reduce((sum, i) => sum + (i.size || 0), 0)
    );
    async function load() {
      loading.value = true;
      try {
        const res = await api.get("/api/images");
        images.value = res.images ?? [];
        totalSize.value = res.totalSize ?? 0;
      } catch (e) {
        toast.error("读取镜像失败", e instanceof Error ? e.message : String(e));
      } finally {
        loading.value = false;
      }
    }
    async function prune() {
      pruning.value = true;
      try {
        const res = await api.post("/api/images/prune");
        toast.success("已清理未使用镜像", `释放 ${formatBytes(res.freedBytes)}`);
        await load();
      } catch (e) {
        toast.error("清理失败", e instanceof Error ? e.message : String(e));
      } finally {
        pruning.value = false;
      }
    }
    async function confirmRemove() {
      const img = removeTarget.value;
      if (!img || removing.value) return;
      removing.value = true;
      try {
        await api.del(`/api/images/${encodeURIComponent(img.id)}`, { force: true });
        toast.success("镜像已删除");
        removeTarget.value = null;
        await load();
      } catch (e) {
        toast.error("删除失败", e instanceof Error ? e.message : String(e));
      } finally {
        removing.value = false;
      }
    }
    onMounted(() => void load());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[6] || (_cache[6] = createBaseVNode("div", { class: "dh-h1" }, "镜像", -1)),
          createBaseVNode("div", _hoisted_3, toDisplayString(images.value.length) + " 个 · 占用 " + toDisplayString(unref(formatBytes)(totalSize.value)) + " · 未使用 " + toDisplayString(danglingCount.value) + " 个（可回收 " + toDisplayString(unref(formatBytes)(danglingSize.value)) + "） ", 1),
          createBaseVNode("div", _hoisted_4, [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: loading.value,
              onClick: load
            }, [
              createVNode(unref(RefreshCw), {
                class: normalizeClass(["h-3.5 w-3.5", loading.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[4] || (_cache[4] = createTextVNode("刷新 ", -1))
            ], 8, _hoisted_5),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: pruning.value || !danglingCount.value,
              onClick: prune
            }, [
              createVNode(unref(Sparkles), {
                class: normalizeClass(["h-3.5 w-3.5", pruning.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[5] || (_cache[5] = createTextVNode(" 清理未使用镜像 ", -1))
            ], 8, _hoisted_6)
          ])
        ]),
        createBaseVNode("div", _hoisted_7, [
          createBaseVNode("div", _hoisted_8, [
            createVNode(unref(Search), { class: "pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" }),
            withDirectives(createBaseVNode("input", {
              "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => keyword.value = $event),
              class: "dh-input !pl-8",
              placeholder: "按标签、仓库名或 ID 搜索"
            }, null, 512), [
              [vModelText, keyword.value]
            ])
          ]),
          createBaseVNode("label", _hoisted_9, [
            withDirectives(createBaseVNode("input", {
              "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => onlyDangling.value = $event),
              type: "checkbox",
              class: "h-[14px] w-[14px] accent-accent"
            }, null, 512), [
              [vModelCheckbox, onlyDangling.value]
            ]),
            createTextVNode(" 只看未使用镜像 (" + toDisplayString(danglingCount.value) + ") ", 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_10, [
          createBaseVNode("div", _hoisted_11, [
            createVNode(unref(Layers), { class: "h-3.5 w-3.5 text-text-4" }),
            _cache[7] || (_cache[7] = createBaseVNode("span", null, "镜像列表", -1)),
            createBaseVNode("span", _hoisted_12, toDisplayString(filtered.value.length) + " 个", 1)
          ]),
          !filtered.value.length ? (openBlock(), createBlock(_sfc_main$1, {
            key: 0,
            icon: unref(Layers),
            title: loading.value ? "正在载入…" : "没有匹配的镜像",
            description: "镜像会随着容器更新不断积累，定期清理未使用镜像可以回收空间。"
          }, null, 8, ["icon", "title"])) : (openBlock(), createElementBlock("div", _hoisted_13, [
            createBaseVNode("table", _hoisted_14, [
              _cache[9] || (_cache[9] = createBaseVNode("thead", null, [
                createBaseVNode("tr", null, [
                  createBaseVNode("th", null, "标签"),
                  createBaseVNode("th", { class: "w-[130px]" }, "镜像 ID"),
                  createBaseVNode("th", { class: "w-[100px]" }, "大小"),
                  createBaseVNode("th", { class: "w-[130px]" }, "创建时间"),
                  createBaseVNode("th", { class: "w-[90px]" }, "被引用"),
                  createBaseVNode("th", { class: "w-[70px]" })
                ])
              ], -1)),
              createBaseVNode("tbody", null, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(filtered.value, (img) => {
                  return openBlock(), createElementBlock("tr", {
                    key: img.id
                  }, [
                    createBaseVNode("td", null, [
                      img.tags.length ? (openBlock(), createElementBlock("div", _hoisted_15, [
                        (openBlock(true), createElementBlock(Fragment, null, renderList(img.tags, (t) => {
                          return openBlock(), createElementBlock("span", {
                            key: t,
                            class: "rounded-md bg-ink-800 px-1.5 py-[2px] font-mono text-[11px] text-text-2"
                          }, toDisplayString(t), 1);
                        }), 128))
                      ])) : (openBlock(), createElementBlock("div", _hoisted_16, [
                        createBaseVNode("span", _hoisted_17, toDisplayString(img.repo || img.shortId), 1),
                        _cache[8] || (_cache[8] = createBaseVNode("span", { class: "dh-badge dh-badge-plain" }, "未使用镜像", -1))
                      ]))
                    ]),
                    createBaseVNode("td", _hoisted_18, toDisplayString(img.shortId), 1),
                    createBaseVNode("td", _hoisted_19, toDisplayString(unref(formatBytes)(img.size)), 1),
                    createBaseVNode("td", _hoisted_20, toDisplayString(unref(relativeTime)(img.created)), 1),
                    createBaseVNode("td", null, [
                      img.containers ? (openBlock(), createElementBlock("span", _hoisted_21, toDisplayString(img.containers) + " 个容器", 1)) : (openBlock(), createElementBlock("span", _hoisted_22, "未使用"))
                    ]),
                    createBaseVNode("td", null, [
                      createBaseVNode("button", {
                        class: "dh-btn dh-btn-sm dh-btn-danger",
                        disabled: img.containers > 0,
                        title: img.containers > 0 ? "还有容器在使用这个镜像，需要先删除容器" : "删除镜像",
                        onClick: ($event) => removeTarget.value = img
                      }, [
                        createVNode(unref(Trash2), { class: "h-3 w-3" })
                      ], 8, _hoisted_23)
                    ])
                  ]);
                }), 128))
              ])
            ])
          ]))
        ]),
        createVNode(_sfc_main$2, {
          open: !!removeTarget.value,
          title: "删除镜像",
          subtitle: removeTarget.value?.tags.join(", ") || removeTarget.value?.repo || removeTarget.value?.shortId,
          width: "440px",
          busy: removing.value,
          onClose: _cache[3] || (_cache[3] = ($event) => removeTarget.value = null)
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: removing.value,
              onClick: _cache[2] || (_cache[2] = ($event) => removeTarget.value = null)
            }, "取消", 8, _hoisted_24),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemove
            }, toDisplayString(removing.value ? "删除中…" : "确认删除"), 9, _hoisted_25)
          ]),
          default: withCtx(() => [
            _cache[10] || (_cache[10] = createBaseVNode("div", { class: "rounded-[10px] border border-line-err bg-soft-err px-3 py-2.5 text-[12px] leading-relaxed text-err-text" }, " 删除镜像本身不会删除容器，但如果这个镜像还在被容器使用，容器下次启动时会失败。 如果只是想回收空间，用「清理未使用镜像」更安全。 ", -1))
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
