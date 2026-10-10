import { d as defineComponent, c as createElementBlock, i as createBlock, O as resolveDynamicComponent, j as createCommentVNode, e as createBaseVNode, t as toDisplayString, s as openBlock } from "./index-M6BbY2BT.js";
const _hoisted_1 = { class: "flex flex-col items-center justify-center gap-2 px-6 py-12 text-center" };
const _hoisted_2 = {
  key: 0,
  class: "grid h-10 w-10 place-items-center rounded-xl border border-line-1 bg-ink-750 text-text-5"
};
const _hoisted_3 = { class: "text-[13px] font-medium text-text-2" };
const _hoisted_4 = {
  key: 1,
  class: "max-w-[420px] text-[12px] leading-relaxed text-text-5"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "EmptyState",
  props: {
    icon: { default: void 0 },
    title: {},
    description: { default: "" },
    actionLabel: { default: "" }
  },
  emits: ["action"],
  setup(__props, { emit: __emit }) {
    const emit = __emit;
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        __props.icon ? (openBlock(), createElementBlock("div", _hoisted_2, [
          (openBlock(), createBlock(resolveDynamicComponent(__props.icon), { class: "h-[18px] w-[18px]" }))
        ])) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_3, toDisplayString(__props.title), 1),
        __props.description ? (openBlock(), createElementBlock("div", _hoisted_4, toDisplayString(__props.description), 1)) : createCommentVNode("", true),
        __props.actionLabel ? (openBlock(), createElementBlock("button", {
          key: 2,
          type: "button",
          class: "dh-btn mt-1.5",
          onClick: _cache[0] || (_cache[0] = ($event) => emit("action"))
        }, toDisplayString(__props.actionLabel), 1)) : createCommentVNode("", true)
      ]);
    };
  }
});
export {
  _sfc_main as _
};
