import { q as createLucideIcon, d as defineComponent, C as useToastStore, k as ref, D as watch, o as onMounted, s as onUnmounted, c as createElementBlock, b as createBaseVNode, t as toDisplayString, e as createVNode, n as normalizeClass, f as unref, R as RefreshCw, i as createTextVNode, E as Download, S as Search, w as withDirectives, G as vModelText, h as createCommentVNode, y as createBlock, H as Box, F as Fragment, r as renderList, I as withCtx, J as useRoute, l as computed, z as api, p as openBlock, K as RouterLink, L as withModifiers, j as vModelCheckbox, B as openStream } from "./index-HPgVGrf3.js";
import { s as shortImage, c as containerStateLabel, r as relativeTime } from "./format-mcaWYSwR.js";
import { _ as _sfc_main$1 } from "./EmptyState.vue_vue_type_script_setup_true_lang-BG9jdsHP.js";
import { _ as _sfc_main$2 } from "./Modal.vue_vue_type_script_setup_true_lang-CGk8pClZ.js";
import { _ as _sfc_main$3, S as Square, R as RotateCw } from "./PortChips.vue_vue_type_script_setup_true_lang-BW6YhVU7.js";
import { P as Play } from "./play-BDui4TlS.js";
import { T as Trash2 } from "./trash-2-C3mj8LF4.js";
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
const _hoisted_10 = ["disabled"];
const _hoisted_11 = { class: "flex flex-wrap items-center gap-2.5" };
const _hoisted_12 = { class: "relative min-w-[180px] flex-1 sm:max-w-[280px]" };
const _hoisted_13 = {
  key: 0,
  class: "ml-auto flex flex-wrap items-center gap-2"
};
const _hoisted_14 = { class: "text-[12px] text-text-4" };
const _hoisted_15 = ["disabled"];
const _hoisted_16 = ["disabled"];
const _hoisted_17 = ["disabled"];
const _hoisted_18 = {
  key: 0,
  class: "dh-card grid h-[240px] place-items-center"
};
const _hoisted_19 = {
  key: 2,
  class: "grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"
};
const _hoisted_20 = { class: "flex items-start gap-2.5" };
const _hoisted_21 = { class: "mt-[3px] flex-none cursor-pointer" };
const _hoisted_22 = ["checked", "onChange"];
const _hoisted_23 = { class: "grid h-[34px] w-[34px] flex-none place-items-center rounded-[10px] bg-line-2 text-[13px] font-semibold text-[#5eead4]" };
const _hoisted_24 = { class: "min-w-0 flex-1" };
const _hoisted_25 = ["title"];
const _hoisted_26 = { class: "relative flex-none" };
const _hoisted_27 = ["onClick"];
const _hoisted_28 = ["onClick"];
const _hoisted_29 = ["onClick"];
const _hoisted_30 = ["onClick"];
const _hoisted_31 = { class: "flex flex-wrap items-center gap-1.5" };
const _hoisted_32 = {
  key: 0,
  class: "dh-badge dh-badge-accent"
};
const _hoisted_33 = {
  key: 1,
  class: "dh-badge dh-badge-plain"
};
const _hoisted_34 = {
  key: 2,
  class: "dh-badge dh-badge-warn"
};
const _hoisted_35 = ["title"];
const _hoisted_36 = { class: "mt-auto flex gap-1.5 border-t border-line-2 pt-2.5" };
const _hoisted_37 = ["disabled", "onClick"];
const _hoisted_38 = ["disabled", "onClick"];
const _hoisted_39 = ["disabled", "onClick"];
const _hoisted_40 = { class: "text-[10.5px] text-text-6" };
const _hoisted_41 = {
  key: 3,
  class: "flex items-center gap-2 text-[12px] text-text-5"
};
const _hoisted_42 = { class: "flex cursor-pointer items-center gap-2" };
const _hoisted_43 = { class: "flex flex-col gap-3" };
const _hoisted_44 = { class: "flex cursor-pointer items-start gap-2.5 text-[12.5px] text-text-2" };
const _hoisted_45 = ["disabled"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "ContainersView",
  setup(__props) {
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
    const checking = ref(false);
    const bulkBusy = ref(false);
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
    async function load(silent = false) {
      if (!silent) loading.value = true;
      try {
        const res = await api.get("/api/containers");
        containers.value = res.containers ?? [];
      } catch (e) {
        toast.error("读取容器失败", e instanceof Error ? e.message : String(e));
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
      if (!names.length) return;
      try {
        await api.post("/api/updates/apply", { names });
        toast.info(`已提交 ${names.length} 个容器的更新`, "镜像未变化的容器会被自动跳过");
        selected.value = /* @__PURE__ */ new Set();
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function checkAll() {
      checking.value = true;
      try {
        const res = await api.post("/api/updates/check", { deep: false });
        const n = (res.results ?? []).filter((r) => r.status === "update_available").length;
        toast.success("巡检完成", n ? `发现 ${n} 个有可用更新` : "所有容器都是最新的");
        await load(true);
      } catch (e) {
        toast.error("巡检失败", e instanceof Error ? e.message : String(e));
      } finally {
        checking.value = false;
      }
    }
    async function updateAll() {
      const names = containers.value.filter((c) => c.hasUpdate && !c.self && !c.excluded).map((c) => c.name);
      if (!names.length) return;
      try {
        await api.post("/api/updates/apply", { names });
        toast.info(`已提交 ${names.length} 个容器的更新`, "镜像未变化的容器会被自动跳过");
      } catch (e) {
        toast.error("提交失败", e instanceof Error ? e.message : String(e));
      }
    }
    async function bulkAct(action) {
      const names = containers.value.filter((c) => selected.value.has(c.name) && !c.self).map((c) => c.name);
      if (!names.length) return;
      bulkBusy.value = true;
      let ok = 0;
      for (const name of names) {
        try {
          await api.post(`/api/containers/${encodeURIComponent(name)}/action`, { action });
          ok += 1;
        } catch {
        }
      }
      bulkBusy.value = false;
      selected.value = /* @__PURE__ */ new Set();
      toast.success(`已${labelOf(action)} ${ok}/${names.length} 台`, ok < names.length ? "部分容器操作失败，详见日志" : void 0);
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
      if (selected.value.size === filtered.value.length) {
        selected.value = /* @__PURE__ */ new Set();
      } else {
        selected.value = new Set(filtered.value.map((c) => c.name));
      }
    }
    let inner = null;
    watch(
      () => route.query.q,
      (q) => {
        keyword.value = q ?? "";
      }
    );
    onMounted(() => {
      void load();
      inner = openStream("/api/events/stream", (topic, ev) => {
        if (topic === "update" && (ev.kind === "batch_done" || ev.kind === "container_status")) {
          void load(true);
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
          _cache[14] || (_cache[14] = createBaseVNode("div", { class: "dh-h1" }, "容器", -1)),
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
              _cache[12] || (_cache[12] = createTextVNode("刷新 ", -1))
            ], 8, _hoisted_8),
            createBaseVNode("button", {
              class: "dh-btn",
              disabled: checking.value,
              onClick: checkAll
            }, [
              createVNode(unref(Download), {
                class: normalizeClass(["h-3.5 w-3.5", checking.value ? "dh-spin" : ""])
              }, null, 8, ["class"]),
              _cache[13] || (_cache[13] = createTextVNode("检测更新 ", -1))
            ], 8, _hoisted_9),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-primary",
              disabled: !counts.value.update,
              onClick: updateAll
            }, " 更新 " + toDisplayString(counts.value.update) + " 台有更新的容器 ", 9, _hoisted_10)
          ])
        ]),
        createBaseVNode("div", _hoisted_11, [
          createBaseVNode("div", _hoisted_12, [
            createVNode(unref(Search), { class: "pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-5" }),
            withDirectives(createBaseVNode("input", {
              "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => keyword.value = $event),
              class: "dh-input !pl-8",
              placeholder: "按名称 / 镜像 / 项目筛选"
            }, null, 512), [
              [vModelText, keyword.value]
            ])
          ]),
          selected.value.size ? (openBlock(), createElementBlock("div", _hoisted_13, [
            createBaseVNode("span", _hoisted_14, "已选 " + toDisplayString(selected.value.size) + " 台", 1),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              onClick: toggleAll
            }, "全选本页"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              disabled: bulkBusy.value,
              onClick: _cache[6] || (_cache[6] = ($event) => bulkAct("restart"))
            }, "重启", 8, _hoisted_15),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm",
              disabled: bulkBusy.value,
              onClick: _cache[7] || (_cache[7] = ($event) => bulkAct("stop"))
            }, "停止", 8, _hoisted_16),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-sm dh-btn-primary",
              disabled: !updatableSelected.value.length,
              onClick: updateSelected
            }, " 更新选中 ", 8, _hoisted_17)
          ])) : createCommentVNode("", true)
        ]),
        _cache[20] || (_cache[20] = createBaseVNode("div", { class: "dh-banner dh-banner-warn" }, [
          createTextVNode(" 批量更新前会先核对镜像摘要："),
          createBaseVNode("b", null, "只有镜像真的变了才会重启容器"),
          createTextVNode("，未变化的容器会原样跳过。 ")
        ], -1)),
        loading.value && !containers.value.length ? (openBlock(), createElementBlock("div", _hoisted_18, [
          createVNode(unref(RefreshCw), { class: "h-5 w-5 dh-spin text-text-5" })
        ])) : !filtered.value.length ? (openBlock(), createBlock(_sfc_main$1, {
          key: 1,
          icon: unref(Box),
          title: "没有匹配的容器",
          description: "换个筛选条件，或者确认 Docker 守护进程是否正常。"
        }, null, 8, ["icon"])) : (openBlock(), createElementBlock("div", _hoisted_19, [
          (openBlock(true), createElementBlock(Fragment, null, renderList(filtered.value, (c) => {
            return openBlock(), createElementBlock("div", {
              key: c.id,
              class: normalizeClass(["flex flex-col gap-2.5 rounded-[14px] border bg-ink-700 p-3 transition-colors", [
                c.hasUpdate ? "border-[rgba(245,165,36,.42)]" : "border-line-1 hover:border-line-4",
                selected.value.has(c.name) ? "!border-[rgba(45,212,191,.5)] bg-[#12201f]" : ""
              ]])
            }, [
              createBaseVNode("div", _hoisted_20, [
                createBaseVNode("label", _hoisted_21, [
                  createBaseVNode("input", {
                    type: "checkbox",
                    class: "h-[14px] w-[14px] accent-[#2dd4bf]",
                    checked: selected.value.has(c.name),
                    onChange: ($event) => toggleSelect(c.name)
                  }, null, 40, _hoisted_22)
                ]),
                createBaseVNode("div", _hoisted_23, toDisplayString(c.name.slice(0, 2).toUpperCase()), 1),
                createBaseVNode("div", _hoisted_24, [
                  createVNode(unref(RouterLink), {
                    to: `/containers/${encodeURIComponent(c.name)}`,
                    class: "block truncate text-[13px] font-semibold hover:text-accent"
                  }, {
                    default: withCtx(() => [
                      createTextVNode(toDisplayString(c.name), 1)
                    ]),
                    _: 2
                  }, 1032, ["to"]),
                  createBaseVNode("div", {
                    class: "truncate font-mono text-[11.5px] text-text-5",
                    title: c.image
                  }, toDisplayString(unref(shortImage)(c.image)), 9, _hoisted_25)
                ]),
                createBaseVNode("div", _hoisted_26, [
                  createBaseVNode("button", {
                    class: "grid h-6 w-6 place-items-center rounded-md text-text-5 hover:bg-ink-650 hover:text-text-1",
                    onClick: withModifiers(($event) => menuFor.value = menuFor.value === c.name ? "" : c.name, ["stop"])
                  }, [
                    createVNode(unref(EllipsisVertical), { class: "h-3.5 w-3.5" })
                  ], 8, _hoisted_27),
                  menuFor.value === c.name ? (openBlock(), createElementBlock("div", {
                    key: 0,
                    class: "absolute right-0 top-7 z-20 w-[150px] overflow-hidden rounded-[10px] border border-line-3 bg-ink-750 py-1 shadow-[0_12px_30px_rgba(0,0,0,.5)]",
                    onClick: _cache[8] || (_cache[8] = withModifiers(() => {
                    }, ["stop"]))
                  }, [
                    createBaseVNode("button", {
                      class: "block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650",
                      onClick: ($event) => act(c, "restart")
                    }, " 重启 ", 8, _hoisted_28),
                    createBaseVNode("button", {
                      class: "block w-full px-3 py-1.5 text-left text-[12px] text-text-2 hover:bg-ink-650",
                      onClick: ($event) => updateOne(c)
                    }, " 检查并更新 ", 8, _hoisted_29),
                    createBaseVNode("button", {
                      class: "block w-full px-3 py-1.5 text-left text-[12px] text-[#fca5a5] hover:bg-ink-650",
                      onClick: ($event) => {
                        removeTarget.value = c;
                        menuFor.value = "";
                      }
                    }, " 删除容器… ", 8, _hoisted_30)
                  ])) : createCommentVNode("", true)
                ])
              ]),
              createBaseVNode("div", _hoisted_31, [
                c.self ? (openBlock(), createElementBlock("span", _hoisted_32, "Dockhelm 自身")) : c.excluded ? (openBlock(), createElementBlock("span", _hoisted_33, "已排除")) : createCommentVNode("", true),
                createBaseVNode("span", {
                  class: normalizeClass(["dh-badge", c.state === "running" ? "dh-badge-run" : "dh-badge-stop"])
                }, [
                  createBaseVNode("span", {
                    class: normalizeClass(["h-[7px] w-[7px] rounded-full", c.state === "running" ? "bg-run" : "bg-stop"])
                  }, null, 2),
                  createTextVNode(" " + toDisplayString(unref(containerStateLabel)(c.state)), 1)
                ], 2),
                c.hasUpdate ? (openBlock(), createElementBlock("span", _hoisted_34, "有新版本")) : createCommentVNode("", true),
                c.project ? (openBlock(), createElementBlock("span", {
                  key: 3,
                  class: "dh-badge dh-badge-plain",
                  title: `compose 项目 ${c.project}`
                }, toDisplayString(c.project), 9, _hoisted_35)) : createCommentVNode("", true)
              ]),
              createVNode(_sfc_main$3, {
                ports: c.portList ?? [],
                max: 3
              }, null, 8, ["ports"]),
              createBaseVNode("div", _hoisted_36, [
                c.state === "running" ? (openBlock(), createElementBlock("button", {
                  key: 0,
                  class: "dh-btn dh-btn-sm flex-1",
                  disabled: busyName.value === c.name,
                  onClick: ($event) => act(c, "stop")
                }, [
                  createVNode(unref(Square), { class: "h-3 w-3" }),
                  _cache[15] || (_cache[15] = createTextVNode("停止 ", -1))
                ], 8, _hoisted_37)) : (openBlock(), createElementBlock("button", {
                  key: 1,
                  class: "dh-btn dh-btn-sm flex-1",
                  disabled: busyName.value === c.name,
                  onClick: ($event) => act(c, "start")
                }, [
                  createVNode(unref(Play), { class: "h-3 w-3" }),
                  _cache[16] || (_cache[16] = createTextVNode("启动 ", -1))
                ], 8, _hoisted_38)),
                createBaseVNode("button", {
                  class: "dh-btn dh-btn-sm flex-1",
                  disabled: busyName.value === c.name || c.self,
                  onClick: ($event) => act(c, "restart")
                }, [
                  createVNode(unref(RotateCw), {
                    class: normalizeClass(["h-3 w-3", busyName.value === c.name ? "dh-spin" : ""])
                  }, null, 8, ["class"]),
                  _cache[17] || (_cache[17] = createTextVNode("重启 ", -1))
                ], 8, _hoisted_39)
              ]),
              createBaseVNode("div", _hoisted_40, "创建于 " + toDisplayString(unref(relativeTime)(c.created)), 1)
            ], 2);
          }), 128))
        ])),
        filtered.value.length && !selected.value.size ? (openBlock(), createElementBlock("div", _hoisted_41, [
          createBaseVNode("label", _hoisted_42, [
            createBaseVNode("input", {
              type: "checkbox",
              class: "h-[14px] w-[14px] accent-[#2dd4bf]",
              checked: false,
              onChange: toggleAll
            }, null, 32),
            createTextVNode(" 全选当前列表（" + toDisplayString(filtered.value.length) + " 台） ", 1)
          ])
        ])) : createCommentVNode("", true),
        createVNode(_sfc_main$2, {
          open: !!removeTarget.value,
          title: "删除容器",
          subtitle: removeTarget.value ? `即将删除 ${removeTarget.value.name}` : "",
          width: "470px",
          onClose: _cache[11] || (_cache[11] = ($event) => {
            removeTarget.value = null;
            removeVolumes.value = false;
          })
        }, {
          footer: withCtx(() => [
            createBaseVNode("button", {
              class: "dh-btn",
              onClick: _cache[10] || (_cache[10] = ($event) => {
                removeTarget.value = null;
                removeVolumes.value = false;
              })
            }, "取消"),
            createBaseVNode("button", {
              class: "dh-btn dh-btn-danger",
              disabled: removing.value,
              onClick: confirmRemove
            }, [
              createVNode(unref(Trash2), { class: "h-3.5 w-3.5" }),
              createTextVNode(toDisplayString(removing.value ? "删除中…" : "确认删除"), 1)
            ], 8, _hoisted_45)
          ]),
          default: withCtx(() => [
            createBaseVNode("div", _hoisted_43, [
              _cache[19] || (_cache[19] = createBaseVNode("div", { class: "rounded-[10px] border border-[rgba(248,113,113,.35)] bg-[rgba(248,113,113,.07)] px-3 py-2.5 text-[12px] leading-relaxed text-[#fca5a5]" }, " 此操作不可撤销。删除容器不会删除它的镜像，但容器自身的可写层数据会一并消失。 ", -1)),
              createBaseVNode("label", _hoisted_44, [
                withDirectives(createBaseVNode("input", {
                  "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => removeVolumes.value = $event),
                  type: "checkbox",
                  class: "mt-[3px] h-[14px] w-[14px] accent-[#f87171]"
                }, null, 512), [
                  [vModelCheckbox, removeVolumes.value]
                ]),
                _cache[18] || (_cache[18] = createBaseVNode("span", null, [
                  createTextVNode(" 同时删除该容器的"),
                  createBaseVNode("b", { class: "text-[#fca5a5]" }, "匿名卷"),
                  createBaseVNode("br"),
                  createBaseVNode("span", { class: "text-[11.5px] text-text-5" }, " 勾选后 docker 会连匿名卷一起删掉，卷里的数据将无法找回。命名卷不会被删除。 ")
                ], -1))
              ])
            ])
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
