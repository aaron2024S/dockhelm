const RUN_KIND_LABEL = {
  update: "更新",
  container: "容器",
  image: "镜像",
  volume: "卷",
  schedule: "计划任务",
  backup: "备份",
  restore: "还原",
  system: "系统",
  check: "检测",
  auto_update: "自动更新",
  auto_check: "自动巡检"
};
function runKindLabel(kind) {
  return RUN_KIND_LABEL[kind] ?? kind;
}
function formatBytes(n, digits = 1) {
  if (n === void 0 || n === null || Number.isNaN(n) || n <= 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  let i = 0;
  let v = n;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i++;
  }
  const fixed = i === 0 ? 0 : digits;
  return `${v.toFixed(fixed)} ${units[i]}`;
}
function relativeTime(input) {
  if (!input) return "—";
  const t = typeof input === "number" ? input * 1e3 : Date.parse(input);
  if (Number.isNaN(t)) return "—";
  const diff = Date.now() - t;
  const abs = Math.abs(diff);
  const future = diff < 0;
  const sec = Math.round(abs / 1e3);
  const fmt = (n, unit) => future ? `${n} ${unit}后` : `${n} ${unit}前`;
  if (sec < 45) return future ? "即将" : "刚刚";
  const min = Math.round(sec / 60);
  if (min < 60) return fmt(min, "分钟");
  const hr = Math.round(min / 60);
  if (hr < 24) return fmt(hr, "小时");
  const day = Math.round(hr / 24);
  if (day < 30) return fmt(day, "天");
  return formatDateTime(t);
}
function formatDateTime(input) {
  if (!input) return "—";
  const t = typeof input === "number" ? input > 1e12 ? input : input * 1e3 : Date.parse(input);
  if (Number.isNaN(t)) return "—";
  const d = new Date(t);
  const p = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}
function formatDuration(seconds) {
  if (!seconds || seconds < 0) return "—";
  const d = Math.floor(seconds / 86400);
  const h = Math.floor(seconds % 86400 / 3600);
  const m = Math.floor(seconds % 3600 / 60);
  const parts = [];
  if (d) parts.push(`${d} 天`);
  if (h) parts.push(`${h} 小时`);
  if (!d && m) parts.push(`${m} 分钟`);
  if (!parts.length) parts.push(`${Math.floor(seconds)} 秒`);
  return parts.join(" ");
}
const STATUS_LABEL = {
  running: "运行中",
  exited: "已停止",
  created: "已创建",
  paused: "已暂停",
  restarting: "重启中",
  dead: "异常",
  removing: "删除中"
};
function containerStateLabel(state) {
  return STATUS_LABEL[state] ?? state;
}
function shortImage(image) {
  if (!image) return "—";
  return image.replace(/^docker\.io\/library\//, "").replace(/^library\//, "");
}
function formatDayTime(input, mode = "datetime") {
  if (!input) return "—";
  const t = typeof input === "number" ? input > 1e12 ? input : input * 1e3 : Date.parse(input);
  if (Number.isNaN(t)) return "—";
  const d = new Date(t);
  const p = (n) => String(n).padStart(2, "0");
  const hm = `${p(d.getHours())}:${p(d.getMinutes())}`;
  const day0 = new Date(t);
  day0.setHours(0, 0, 0, 0);
  const today0 = /* @__PURE__ */ new Date();
  today0.setHours(0, 0, 0, 0);
  const diffDays = Math.round((day0.getTime() - today0.getTime()) / 864e5);
  let day;
  if (diffDays === 0) day = "今天";
  else if (diffDays === 1) day = "明天";
  else if (diffDays === 2) day = "后天";
  else if (diffDays > 2 && diffDays < 7) day = `周${"日一二三四五六"[d.getDay()]}`;
  else day = `${d.getMonth() + 1} 月 ${d.getDate()} 日`;
  return mode === "date" ? day : `${day} ${hm}`;
}
function explainCron(expr) {
  const e = expr.trim();
  if (!e) return "";
  if (e === "@daily" || e === "@midnight") return "每天 00:00";
  if (e === "@hourly") return "每小时";
  if (e === "@weekly") return "每周";
  if (e === "@monthly") return "每月";
  if (e === "@yearly" || e === "@annually") return "每年";
  const every = /^@every\s+(.+)$/.exec(e);
  if (every) return `每 ${every[1]}`;
  const parts = e.split(/\s+/);
  if (parts.length !== 5) return e;
  const min = parts[0];
  const hour = parts[1];
  const dom = parts[2];
  const mon = parts[3];
  const dow = parts[4];
  const isNum = (s) => /^\d+$/.test(s);
  const pad = (s) => s.padStart(2, "0");
  if (isNum(min) && isNum(hour)) {
    const hm = `${pad(hour)}:${pad(min)}`;
    if (dom === "*" && mon === "*" && dow === "*") return `每天 ${hm}`;
    if (dom === "*" && mon === "*" && dow !== "*") {
      const names = {
        "0": "周日",
        "1": "周一",
        "2": "周二",
        "3": "周三",
        "4": "周四",
        "5": "周五",
        "6": "周六",
        "7": "周日"
      };
      const list = dow.split(",").map((d) => names[d] ?? `周${d}`).join("、");
      return `每${list} ${hm}`;
    }
    if (mon === "*" && dow === "*" && isNum(dom)) return `每月 ${dom} 号 ${hm}`;
  }
  if (min === "*" && isNum(hour)) return `每小时的第 ${hour} 分（${pad(hour)}:${min.replace("*", "00")} 起每分钟）`;
  if (min.startsWith("*/")) return `每 ${min.slice(2)} 分钟`;
  if (hour.startsWith("*/")) return `每 ${hour.slice(2)} 小时`;
  return e;
}
export {
  formatBytes as a,
  relativeTime as b,
  containerStateLabel as c,
  formatDateTime as d,
  explainCron as e,
  formatDayTime as f,
  formatDuration as g,
  runKindLabel as r,
  shortImage as s
};
