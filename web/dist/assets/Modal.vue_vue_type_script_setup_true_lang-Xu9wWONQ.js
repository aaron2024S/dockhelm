import { x as createLucideIcon, d as defineComponent, i as createBlock, f as createVNode, a2 as Transition, I as withCtx, a3 as Teleport, s as openBlock, c as createElementBlock, e as createBaseVNode, t as toDisplayString, j as createCommentVNode, g as unref, a4 as X, $ as renderSlot, z as normalizeStyle, L as withModifiers } from "./index-DC-2UtNU.js";
const Trash2 = createLucideIcon("Trash2Icon", [
  ["path", { d: "M3 6h18", key: "d0wm0j" }],
  ["path", { d: "M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6", key: "4alrt4" }],
  ["path", { d: "M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2", key: "v07s0e" }],
  ["line", { x1: "10", x2: "10", y1: "11", y2: "17", key: "1uufr5" }],
  ["line", { x1: "14", x2: "14", y1: "11", y2: "17", key: "xtxkd" }]
]);
const _hoisted_1 = { class: "flex items-start gap-3 border-b border-line-1 px-4 py-3" };
const _hoisted_2 = { class: "min-w-0 flex-1" };
const _hoisted_3 = { class: "truncate text-[13.5px] font-semibold" };
const _hoisted_4 = {
  key: 0,
  class: "mt-0.5 text-[11.5px] leading-relaxed text-text-4"
};
const _hoisted_5 = ["disabled"];
const _hoisted_6 = { class: "px-4 py-4" };
const _hoisted_7 = {
  key: 0,
  class: "flex items-center justify-end gap-2 border-t border-line-1 px-4 py-3"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "Modal",
  props: {
    open: { type: Boolean },
    title: { default: "" },
    subtitle: { default: "" },
    width: { default: "520px" },
    closeOnBackdrop: { type: Boolean, default: true },
    busy: { type: Boolean, default: false }
  },
  emits: ["close"],
  setup(__props, { emit: __emit }) {
    const emit = __emit;
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, {
          "enter-active-class": "transition duration-150 ease-out",
          "enter-from-class": "opacity-0",
          "leave-active-class": "transition duration-120 ease-in",
          "leave-to-class": "opacity-0"
        }, {
          default: withCtx(() => [
            __props.open ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "fixed inset-0 z-[70] grid place-items-center overflow-y-auto bg-black/60 p-4 backdrop-blur-[2px]",
              onClick: _cache[1] || (_cache[1] = withModifiers(($event) => __props.closeOnBackdrop && !__props.busy && emit("close"), ["self"]))
            }, [
              createBaseVNode("div", {
                class: "dh-card w-full shadow-[var(--shadow-modal)]",
                style: normalizeStyle({ maxWidth: __props.width }),
                role: "dialog",
                "aria-modal": "true"
              }, [
                createBaseVNode("div", _hoisted_1, [
                  createBaseVNode("div", _hoisted_2, [
                    createBaseVNode("div", _hoisted_3, toDisplayString(__props.title), 1),
                    __props.subtitle ? (openBlock(), createElementBlock("div", _hoisted_4, toDisplayString(__props.subtitle), 1)) : createCommentVNode("", true)
                  ]),
                  createBaseVNode("button", {
                    type: "button",
                    disabled: __props.busy,
                    class: "dh-tap grid h-6 w-6 flex-none place-items-center rounded-md text-text-5 transition-colors hover:bg-ink-650 hover:text-text-1 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent",
                    onClick: _cache[0] || (_cache[0] = ($event) => emit("close"))
                  }, [
                    createVNode(unref(X), { class: "h-3.5 w-3.5" })
                  ], 8, _hoisted_5)
                ]),
                createBaseVNode("div", _hoisted_6, [
                  renderSlot(_ctx.$slots, "default")
                ]),
                _ctx.$slots.footer ? (openBlock(), createElementBlock("div", _hoisted_7, [
                  renderSlot(_ctx.$slots, "footer")
                ])) : createCommentVNode("", true)
              ], 4)
            ])) : createCommentVNode("", true)
          ]),
          _: 3
        })
      ]);
    };
  }
});
export {
  Trash2 as T,
  _sfc_main as _
};
