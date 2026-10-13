import { d as runKindLabel, g as runRefText, f as formatDayTime, h as runRefShown } from "./format-BfNgxWdW.js";
import { d as defineComponent, s as openBlock, c as createElementBlock, e as createBaseVNode, t as toDisplayString, n as normalizeClass, g as unref, j as createCommentVNode, k as createTextVNode } from "./index-BpYw31T2.js";
const _hoisted_1 = { class: "dh-tl" };
const _hoisted_2 = { class: "w-[86px] flex-none text-[11.5px] text-text-5 max-md:order-1" };
const _hoisted_3 = ["title"];
const _hoisted_4 = {
  key: 0,
  class: "mr-1.5 font-mono text-[11.5px] text-text-3"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "RunLogRow",
  props: {
    log: {}
  },
  setup(__props) {
    const props = __props;
    function logTime(ts) {
      const t = Date.parse(ts);
      if (Number.isNaN(t)) return "—";
      const d = new Date(t);
      const p = (n) => String(n).padStart(2, "0");
      const hm = `${p(d.getHours())}:${p(d.getMinutes())}`;
      const today0 = /* @__PURE__ */ new Date();
      today0.setHours(0, 0, 0, 0);
      const day0 = new Date(t);
      day0.setHours(0, 0, 0, 0);
      const diff = Math.round((day0.getTime() - today0.getTime()) / 864e5);
      if (diff === 0) return hm;
      if (diff === -1) return `昨天 ${hm}`;
      return formatDayTime(ts, "date");
    }
    function kindClass(kind, status) {
      if (status === "failed") return "dh-badge-err";
      if (kind === "schedule") return "dh-badge-accent";
      return "dh-badge-plain";
    }
    function resultBadge(l) {
      if (l.status === "success") return { text: "成功", cls: "dh-badge-run" };
      if (l.status === "up_to_date") return { text: "跳过", cls: "dh-badge-accent" };
      if (l.status === "failed") return { text: "失败", cls: "dh-badge-err" };
      if (l.status === "done") return { text: "完成", cls: "dh-badge-run" };
      return { text: l.status || "—", cls: "dh-badge-plain" };
    }
    function refShown() {
      return runRefShown(props.log.kind, props.log.ref);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, toDisplayString(logTime(__props.log.ts)), 1),
        createBaseVNode("span", {
          class: normalizeClass(["dh-badge min-w-[64px] flex-none justify-center max-md:order-2", kindClass(__props.log.kind, __props.log.status)])
        }, toDisplayString(unref(runKindLabel)(__props.log.kind)), 3),
        createBaseVNode("span", {
          class: normalizeClass(["dh-badge min-w-[42px] flex-none justify-center max-md:order-3", resultBadge(__props.log).cls])
        }, toDisplayString(resultBadge(__props.log).text), 3),
        createBaseVNode("span", {
          class: "min-w-0 flex-1 truncate text-[12.5px] text-text-2 max-md:order-4",
          title: refShown() ? `${unref(runRefText)(__props.log.kind, __props.log.ref)} ${__props.log.message}` : __props.log.message
        }, [
          refShown() ? (openBlock(), createElementBlock("span", _hoisted_4, toDisplayString(unref(runRefText)(__props.log.kind, __props.log.ref)), 1)) : createCommentVNode("", true),
          createTextVNode(toDisplayString(__props.log.message), 1)
        ], 8, _hoisted_3)
      ]);
    };
  }
});
export {
  _sfc_main as _
};
