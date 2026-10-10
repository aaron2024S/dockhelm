import { d as defineComponent, c as createElementBlock, s as openBlock } from "./index-Cpl09g5v.js";
const _hoisted_1 = ["data-on", "disabled", "aria-checked", "aria-label"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "UiSwitch",
  props: {
    modelValue: { type: Boolean },
    disabled: { type: Boolean, default: false },
    label: { default: "" }
  },
  emits: ["update:modelValue"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    function toggle() {
      if (props.disabled) return;
      emit("update:modelValue", !props.modelValue);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("button", {
        type: "button",
        role: "switch",
        class: "dh-switch",
        "data-on": __props.modelValue ? "true" : "false",
        disabled: __props.disabled,
        "aria-checked": __props.modelValue,
        "aria-label": __props.label || void 0,
        onClick: toggle
      }, null, 8, _hoisted_1);
    };
  }
});
export {
  _sfc_main as _
};
