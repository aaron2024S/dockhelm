import { d as defineComponent, y as createBlock, e as createVNode, Y as Transition, I as withCtx, Z as Teleport, p as openBlock, c as createElementBlock, b as createBaseVNode, t as toDisplayString, h as createCommentVNode, f as unref, _ as X, $ as renderSlot, x as normalizeStyle, L as withModifiers } from "./index-HPgVGrf3.js";
const _hoisted_1 = { class: "flex items-start gap-3 border-b border-line-1 px-4 py-3" };
const _hoisted_2 = { class: "min-w-0 flex-1" };
const _hoisted_3 = { class: "truncate text-[13.5px] font-semibold" };
const _hoisted_4 = {
  key: 0,
  class: "mt-0.5 text-[11.5px] leading-relaxed text-text-4"
};
const _hoisted_5 = { class: "px-4 py-4" };
const _hoisted_6 = {
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
    closeOnBackdrop: { type: Boolean, default: true }
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
              onClick: _cache[1] || (_cache[1] = withModifiers(($event) => __props.closeOnBackdrop && emit("close"), ["self"]))
            }, [
              createBaseVNode("div", {
                class: "dh-card w-full shadow-[0_20px_60px_rgba(0,0,0,.5)]",
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
                    class: "grid h-6 w-6 flex-none place-items-center rounded-md text-text-5 transition-colors hover:bg-ink-650 hover:text-text-1",
                    onClick: _cache[0] || (_cache[0] = ($event) => emit("close"))
                  }, [
                    createVNode(unref(X), { class: "h-3.5 w-3.5" })
                  ])
                ]),
                createBaseVNode("div", _hoisted_5, [
                  renderSlot(_ctx.$slots, "default")
                ]),
                _ctx.$slots.footer ? (openBlock(), createElementBlock("div", _hoisted_6, [
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
  _sfc_main as _
};
