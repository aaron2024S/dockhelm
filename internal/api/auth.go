package api

import (
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
	writeOK(w, map[string]any{
		"initialized":  s.auth.Initialized(),
		"loggedIn":     loggedIn,
		"failures":     s.auth.Failures(ip),
		"maxFailures":  auth.MaxFailures,
		"minPassword":  auth.MinPasswordLen,
		"sessionCount": s.st.SessionCount(),
	})
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
		remaining := auth.MaxFailures - s.auth.Failures(ip)
		if remaining < 0 {
			remaining = 0
		}
		if err == auth.ErrLocked {
			writeErrCode(w, http.StatusTooManyRequests, err.Error(), "locked")
			return
		}
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"error":     err.Error(),
			"code":      "bad_password",
			"remaining": remaining,
		})
		return
	}

	ttl := auth.SessionTTL
	if in.Keep {
		ttl = auth.SessionTTLKeep
	}
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
		// Secure 由反向代理决定；局域网 HTTP 部署下强制 Secure 会导致登录不了
		Secure: strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https"),
	})

	s.nt.Emit("login_success", map[string]string{
		"result":  "登录成功",
		"message": "来源 IP " + ip,
	})
	writeOK(w, map[string]any{"ok": true, "expiresIn": int(ttl.Seconds())})
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
	if err := s.auth.ForceSetPassword(in.Password); err != nil {
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
	http.SetCookie(w, &http.Cookie{
		Name: auth.CookieName, Value: token, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode,
		MaxAge: int(auth.SessionTTLKeep.Seconds()),
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
