import { d as defineComponent, c as createElementBlock, a as createBaseVNode, t as toDisplayString, g as createCommentVNode, _ as renderSlot, m as openBlock, $ as useModel, a0 as mergeModels } from "./index-D8YHJPqb.js";
const _hoisted_1$1 = { class: "flex items-center gap-3" };
const _hoisted_2 = { class: "min-w-0 flex-1" };
const _hoisted_3 = { class: "text-[12.5px] leading-snug text-text-1" };
const _hoisted_4 = {
  key: 0,
  class: "mt-[2px] text-[11px] leading-relaxed text-text-5"
};
const _hoisted_5 = { class: "flex flex-none items-center gap-2" };
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  __name: "SettingRow",
  props: {
    title: {},
    sub: {}
  },
  setup(__props) {
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$1, [
        createBaseVNode("div", _hoisted_2, [
          createBaseVNode("div", _hoisted_3, toDisplayString(__props.title), 1),
          __props.sub ? (openBlock(), createElementBlock("div", _hoisted_4, toDisplayString(__props.sub), 1)) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", _hoisted_5, [
          renderSlot(_ctx.$slots, "default")
        ])
      ]);
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
