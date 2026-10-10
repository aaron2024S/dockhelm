import { x as createLucideIcon, d as defineComponent, c as createElementBlock, F as Fragment, r as renderList, t as toDisplayString, j as createCommentVNode, p as computed, s as openBlock, n as normalizeClass, e as createBaseVNode } from "./index-BoPP5Zol.js";
const RotateCw = createLucideIcon("RotateCwIcon", [
  ["path", { d: "M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8", key: "1p45f6" }],
  ["path", { d: "M21 3v5h-5", key: "1q7to0" }]
]);
const Square = createLucideIcon("SquareIcon", [
  ["rect", { width: "18", height: "18", x: "3", y: "3", rx: "2", key: "afitv7" }]
]);
const _hoisted_1 = {
  key: 0,
  class: "flex flex-wrap items-center gap-1"
};
const _hoisted_2 = ["title"];
const _hoisted_3 = { class: "font-semibold text-text-1" };
const _hoisted_4 = { key: 1 };
const _hoisted_5 = ["title"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "PortChips",
  props: {
    ports: {},
    max: { default: 3 }
  },
  setup(__props) {
    const props = __props;
    const visible = computed(() => props.max > 0 ? props.ports.slice(0, props.max) : props.ports);
    const hidden = computed(() => props.max > 0 ? props.ports.length - props.max : 0);
    function keyOf(p) {
      return `${p.hostIp}:${p.hostPort}>${p.innerPort}/${p.proto}`;
    }
    function tipOf(p) {
      return p.published ? `已发布到宿主机 ${p.hostIp || "0.0.0.0"}:${p.hostPort} → 容器内 ${p.innerPort}/${p.proto}` : `仅容器内可见（未发布到宿主机）：${p.innerPort}/${p.proto}`;
    }
    const allTip = computed(
      () => props.ports.map((p) => p.published ? `${p.hostIp || "0.0.0.0"}:${p.hostPort} → ${p.innerPort}/${p.proto}` : `${p.innerPort}/${p.proto}（仅容器内）`).join("\n")
    );
    return (_ctx, _cache) => {
      return __props.ports.length ? (openBlock(), createElementBlock("div", _hoisted_1, [
        (openBlock(true), createElementBlock(Fragment, null, renderList(visible.value, (p) => {
          return openBlock(), createElementBlock("span", {
            key: keyOf(p),
            class: normalizeClass([
              "inline-flex items-center gap-1 rounded-md px-1.5 py-[2px] font-mono text-[10.5px]",
              p.published ? "bg-ink-850 text-text-2 ring-1 ring-inset ring-line-3" : "border border-dashed border-line-3 text-text-5"
            ]),
            title: tipOf(p)
          }, [
            p.published ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
              createBaseVNode("span", _hoisted_3, toDisplayString(p.hostPort), 1),
              _cache[0] || (_cache[0] = createBaseVNode("span", { class: "text-text-6" }, "→", -1)),
              createBaseVNode("span", null, toDisplayString(p.innerPort) + "/" + toDisplayString(p.proto), 1)
            ], 64)) : (openBlock(), createElementBlock("span", _hoisted_4, toDisplayString(p.innerPort) + "/" + toDisplayString(p.proto), 1))
          ], 10, _hoisted_2);
        }), 128)),
        hidden.value > 0 ? (openBlock(), createElementBlock("span", {
          key: 0,
          class: "px-1 text-[10.5px] text-text-5",
          title: allTip.value
        }, " +" + toDisplayString(hidden.value), 9, _hoisted_5)) : createCommentVNode("", true)
      ])) : createCommentVNode("", true);
    };
  }
});
export {
  RotateCw as R,
  Square as S,
  _sfc_main as _
};
