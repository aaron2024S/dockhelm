package api

import (
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/aaron2024s/dockhelm/internal/auth"
)

// hSession 返回登录状态（前端启动时第一个调用的接口）。
func (s *Server) hSession(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	if c, err := r.Cookie(auth.CookieName); err == nil && s.auth.Validate(c.Value) {
		loggedIn = true
	}
	ip := s.clientIP(r)
	left := s.auth.LockedFor(ip)
	writeOK(w, map[string]any{
		"initialized":  s.auth.Initialized(),
		"loggedIn":     loggedIn,
		"failures":     s.auth.Failures(ip),
		"maxFailures":  auth.MaxFailures,
		"lockedFor":    retryAfterSeconds(left),
		"lockedHint":   auth.LockedHint(left),
		"minPassword":  auth.MinPasswordLen,
		"sessionCount": s.st.SessionCount(),
	})
}

// retryAfterSeconds 把剩余锁定时间换算成给前端的整数秒（向上取整）。
//
// 向上取整是为了「宁可多显示 1 秒」：向下取整会在还剩 0.4 秒时给出 0，
// 前端据此把界面解锁、按钮点亮，用户一点又是 429 —— 看着像坏了。
func retryAfterSeconds(left time.Duration) int {
	if left <= 0 {
		return 0
	}
	return int(math.Ceil(left.Seconds()))
}

type loginReq struct {
	Password string `json:"password"`
	Keep     bool   `json:"keep"`
}

// hLogin 登录。
func (s *Server) hLogin(w http.ResponseWriter, r *http.Request) {
	var in loginReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	ip := s.clientIP(r)
	ua := r.UserAgent()

	token, err := s.auth.Verify(in.Password, ip, ua, in.Keep)
	if err != nil {
		// 登录失败也要通知（含来源 IP），这是「面板被爆破」的第一手线索
		s.nt.Emit("login_failed", map[string]string{
			"result":  "登录失败",
			"message": "来源 IP " + ip + "，原因：" + err.Error(),
		})
		if err == auth.ErrLocked {
			// 锁定提示要带「还要等多久」，前端拿 retryAfter 起倒计时
			left := s.auth.LockedFor(ip)
			writeJSON(w, http.StatusTooManyRequests, map[string]any{
				"error":      auth.LockedHint(left),
				"code":       "locked",
				"retryAfter": retryAfterSeconds(left),
			})
			return
		}
		remaining := auth.MaxFailures - s.auth.Failures(ip)
		if remaining < 0 {
			remaining = 0
		}
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"error":     err.Error(),
			"code":      "bad_password",
			"remaining": remaining,
		})
		return
	}

	// Secure 由反向代理决定；局域网 HTTP 部署下强制 Secure 会导致登录不了
	secure := strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	http.SetCookie(w, sessionCookie(token, in.Keep, secure))

	s.nt.Emit("login_success", map[string]string{
		"result":  "登录成功",
		"message": "来源 IP " + ip,
	})
	writeOK(w, map[string]any{"ok": true, "expiresIn": int(sessionTTL(in.Keep).Seconds())})
}

// sessionTTL 会话在服务端的有效期。
//
// 勾「保持登录」⇒ 7 天；不勾 ⇒ 24 小时上限兜底 —— 浏览器端那时下发的是会话
// Cookie（关浏览器即失效），这个上限只是防止「恢复标签页」让会话无限续命。
func sessionTTL(keep bool) time.Duration {
	if keep {
		return auth.SessionTTLKeep
	}
	return auth.SessionTTL
}

// sessionCookie 构造登录会话 Cookie。
//
// 勾「保持登录」⇒ 下发 MaxAge，浏览器持久保存（7 天）；
// 不勾 ⇒ MaxAge 留 0，浏览器按会话 Cookie 处理，关掉浏览器即失效。
func sessionCookie(token string, keep, secure bool) *http.Cookie {
	c := &http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	}
	if keep {
		c.MaxAge = int(auth.SessionTTLKeep.Seconds())
	}
	return c
}

// hLogout 退出登录。
func (s *Server) hLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(auth.CookieName); err == nil {
		_ = s.auth.Logout(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name: auth.CookieName, Value: "", Path: "/", HttpOnly: true, MaxAge: -1,
	})
	writeOK(w, map[string]any{"ok": true})
}

type setupReq struct {
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

// hSetup 首启设置密码。只在尚未初始化时可用。
func (s *Server) hSetup(w http.ResponseWriter, r *http.Request) {
	if s.auth.Initialized() {
		writeErrCode(w, http.StatusConflict, "密码已设置，如需重置请在 NAS 面板里删除 data/auth.json 后重启容器", "already_initialized")
		return
	}
	var in setupReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if in.Password != in.Confirm {
		writeErr(w, http.StatusBadRequest, "两次输入的密码不一致")
		return
	}
	// 走 SetupPassword 而不是 ForceSetPassword：只靠上面那次 Initialized() 判断
	// 是「检查后使用」，两个并发请求会双双通过，后来者覆盖先来者设的密码。
	if err := s.auth.SetupPassword(in.Password); err != nil {
		if errors.Is(err, auth.ErrAlreadyInitialized) {
			writeErrCode(w, http.StatusConflict, "密码已设置，如需重置请在 NAS 面板里删除 data/auth.json 后重启容器", "already_initialized")
			return
		}
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	// 设置完直接登录，省一次输入
	ip := s.clientIP(r)
	token, err := s.auth.Verify(in.Password, ip, r.UserAgent(), true)
	if err != nil {
		writeOK(w, map[string]any{"ok": true, "autoLogin": false})
		return
	}
	// 首启自动登录下发**会话 Cookie**（不设 MaxAge）：与登录页「不勾保持登录」
	// 同一语义 —— 关掉浏览器就失效，不悄悄留下一个持久登录。
	http.SetCookie(w, &http.Cookie{
		Name: auth.CookieName, Value: token, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	writeOK(w, map[string]any{"ok": true, "autoLogin": true})
}

// hAccount 账户信息。
func (s *Server) hAccount(w http.ResponseWriter, r *http.Request) {
	info := map[string]any{
		"initialized":  s.auth.Initialized(),
		"sessionCount": s.st.SessionCount(),
	}
	if c, err := r.Cookie(auth.CookieName); err == nil {
		if sess := s.st.GetSession(c.Value); sess != nil {
			info["currentSession"] = map[string]any{
				"createdAt": sess.CreatedAt,
				"expiresAt": sess.ExpiresAt,
				"ip":        sess.IP,
			}
		}
	}
	recent := s.st.RecentLogins(10)
	list := make([]map[string]any, 0, len(recent))
	for _, it := range recent {
		list = append(list, map[string]any{"ts": it.TS, "ip": it.IP, "ok": it.OK})
	}
	info["recentLogins"] = list
	info["now"] = time.Now().UTC().Format(time.RFC3339)
	writeOK(w, info)
}

// hLogoutOthers 强制登出当前会话之外的所有登录会话。
//
// 会话没有「登出按钮」，被忘在旧浏览器/旧设备上时只能等过期（最长 7 天）。
// 这个接口给用户一个主动清场出口：以请求自带的会话 Cookie 为「当前」，
// 其余全部作废 —— 改密码也能达到同样效果，但对只是想清理残留的用户来说太重。
func (s *Server) hLogoutOthers(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(auth.CookieName)
	if err != nil || c.Value == "" {
		writeErr(w, http.StatusUnauthorized, "未登录")
		return
	}
	if s.st.GetSession(c.Value) == nil {
		writeErr(w, http.StatusUnauthorized, "当前会话已失效，请重新登录")
		return
	}
	n, err := s.st.DeleteOtherSessions(c.Value)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"ok": true, "revoked": n})
}

type changePwReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Confirm     string `json:"confirm"`
}

// hChangePassword 修改密码。改完所有既有会话失效。
func (s *Server) hChangePassword(w http.ResponseWriter, r *http.Request) {
	var in changePwReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if in.NewPassword != in.Confirm {
		writeErr(w, http.StatusBadRequest, "两次输入的新密码不一致")
		return
	}
	if err := s.auth.SetPassword(in.OldPassword, in.NewPassword); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	// 旧 Cookie 已随会话一起失效，清掉它
	http.SetCookie(w, &http.Cookie{
		Name: auth.CookieName, Value: "", Path: "/", HttpOnly: true, MaxAge: -1,
	})
	s.nt.Emit("password_changed", map[string]string{
		"result":  "密码已修改",
		"message": "来源 IP " + s.clientIP(r) + "，所有已有登录会话已失效",
	})
	writeOK(w, map[string]any{"ok": true, "relogin": true})
}
