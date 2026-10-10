import { d as defineComponent, c as createElementBlock, n as normalizeClass, e as createBaseVNode, t as toDisplayString, j as createCommentVNode, $ as renderSlot, s as openBlock, a0 as useModel, a1 as mergeModels } from "./index-M6BbY2BT.js";
const _hoisted_1$1 = { class: "min-w-0 flex-1" };
const _hoisted_2 = { class: "text-[12.5px] leading-snug text-text-1" };
const _hoisted_3 = {
  key: 0,
  class: "mt-[2px] text-[11px] leading-relaxed text-text-5"
};
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  __name: "SettingRow",
  props: {
    title: {},
    sub: {},
    stack: { type: Boolean }
  },
  setup(__props) {
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(["flex items-center gap-3", __props.stack ? "max-md:flex-col max-md:items-stretch" : ""])
      }, [
        createBaseVNode("div", _hoisted_1$1, [
          createBaseVNode("div", _hoisted_2, toDisplayString(__props.title), 1),
          __props.sub ? (openBlock(), createElementBlock("div", _hoisted_3, toDisplayString(__props.sub), 1)) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", {
          class: normalizeClass(["flex flex-none items-center gap-2", __props.stack ? "max-md:w-full" : ""])
        }, [
          renderSlot(_ctx.$slots, "default")
        ], 2)
      ], 2);
    };
  }
});
const _hoisted_1 = ["aria-checked", "aria-label", "data-on", "disabled"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "ToggleSwitch",
  props: /* @__PURE__ */ mergeModels({
    disabled: { type: Boolean },
    label: {}
  }, {
    "modelValue": { type: Boolean, ...{ required: true } },
    "modelModifiers": {}
  }),
  emits: ["update:modelValue"],
  setup(__props) {
    const model = useModel(__props, "modelValue");
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("button", {
        type: "button",
        class: "dh-switch",
        role: "switch",
        "aria-checked": model.value,
        "aria-label": __props.label,
        "data-on": model.value,
        disabled: __props.disabled,
        onClick: _cache[0] || (_cache[0] = ($event) => model.value = !model.value)
      }, null, 8, _hoisted_1);
    };
  }
});
export {
  _sfc_main$1 as _,
  _sfc_main as a
};
