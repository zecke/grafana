package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	macaron "gopkg.in/macaron.v1"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

var getTime = time.Now

const (
	errStringInvalidUsernamePassword = "Invalid username or password"
	errStringInvalidAPIKey           = "Invalid API key"
)

var (
	ReqGrafanaAdmin = Auth(&AuthOptions{
		ReqSignedIn:     true,
		ReqGrafanaAdmin: true,
	})
	ReqSignedIn   = Auth(&AuthOptions{ReqSignedIn: true})
	ReqEditorRole = RoleAuth(models.ROLE_EDITOR, models.ROLE_ADMIN)
	ReqOrgAdmin   = RoleAuth(models.ROLE_ADMIN)
)

func GetContextHandler(
	ats models.UserTokenService,
	remoteCache *remotecache.RemoteCache,
) macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &models.ReqContext{
			Context:        c,
			SignedInUser:   &models.SignedInUser{},
			IsSignedIn:     true,
			AllowAnonymous: true,
			SkipCache:      false,
			Logger:         log.New("context"),
		}

		// this is the fake admin/admin user that currently exist; everyone is admin going forward.
		// {UserId:1 OrgId:1 OrgName:Main Org. OrgRole:Admin Login:admin Name: Email:admin@localhost ApiKeyId:0 OrgCount:1 IsGrafanaAdmin:true IsAnonymous:false HelpFlags1:0 LastSeenAt:2019-10-28 11:12:44 +0000 UTC Teams:[]}
		initContextWithAnonymousUser(ctx)
		ctx.SignedInUser.IsGrafanaAdmin = true
		ctx.SignedInUser.OrgId = 1
		ctx.SignedInUser.UserId = 1
		ctx.SignedInUser.OrgName = "Main Org."
		ctx.SignedInUser.OrgRole = "Admin"
		ctx.SignedInUser.Login = "admin"
		ctx.SignedInUser.Email = "admin@locallhost"
		ctx.SignedInUser.OrgCount = 1

		ctx.Logger = log.New("context", "userId", ctx.UserId, "orgId", ctx.OrgId, "uname", ctx.Login)
		ctx.Data["ctx"] = ctx

		c.Map(ctx)

		// update last seen every 5min
		if ctx.ShouldUpdateLastSeenAt() {
			ctx.Logger.Debug("Updating last user_seen_at", "user_id", ctx.UserId)
			if err := bus.Dispatch(&models.UpdateUserLastSeenAtCommand{UserId: ctx.UserId}); err != nil {
				ctx.Logger.Error("Failed to update last_seen_at", "error", err)
			}
		}
	}
}

func initContextWithAnonymousUser(ctx *models.ReqContext) bool {
	if !setting.AnonymousEnabled {
		return false
	}

	orgQuery := models.GetOrgByNameQuery{Name: setting.AnonymousOrgName}
	if err := bus.Dispatch(&orgQuery); err != nil {
		log.Error(3, "Anonymous access organization error: '%s': %s", setting.AnonymousOrgName, err)
		return false
	}

	ctx.IsSignedIn = false
	ctx.AllowAnonymous = true
	ctx.SignedInUser = &models.SignedInUser{IsAnonymous: true}
	ctx.OrgRole = models.RoleType(setting.AnonymousOrgRole)
	ctx.OrgId = orgQuery.Result.Id
	ctx.OrgName = orgQuery.Result.Name
	return true
}

func WriteSessionCookie(ctx *models.ReqContext, value string, maxLifetimeDays int) {
	if setting.Env == setting.DEV {
		ctx.Logger.Info("New token", "unhashed token", value)
	}

	var maxAge int
	if maxLifetimeDays <= 0 {
		maxAge = -1
	} else {
		maxAgeHours := (time.Duration(setting.LoginMaxLifetimeDays) * 24 * time.Hour) + time.Hour
		maxAge = int(maxAgeHours.Seconds())
	}

	ctx.Resp.Header().Del("Set-Cookie")
	cookie := http.Cookie{
		Name:     setting.LoginCookieName,
		Value:    url.QueryEscape(value),
		HttpOnly: true,
		Path:     setting.AppSubUrl + "/",
		Secure:   setting.CookieSecure,
		MaxAge:   maxAge,
	}
	if setting.CookieSameSite != http.SameSiteDefaultMode {
		cookie.SameSite = setting.CookieSameSite
	}

	http.SetCookie(ctx.Resp, &cookie)
}

func AddDefaultResponseHeaders() macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Resp.Before(func(w macaron.ResponseWriter) {
			if !strings.HasPrefix(ctx.Req.URL.Path, "/api/datasources/proxy/") {
				AddNoCacheHeaders(ctx.Resp)
			}

			if !setting.AllowEmbedding {
				AddXFrameOptionsDenyHeader(w)
			}

			AddSecurityHeaders(w)
		})
	}
}

// AddSecurityHeaders adds various HTTP(S) response headers that enable various security protections behaviors in the client's browser.
func AddSecurityHeaders(w macaron.ResponseWriter) {
	if (setting.Protocol == setting.HTTPS || setting.Protocol == setting.HTTP2) && setting.StrictTransportSecurity {
		strictHeaderValues := []string{fmt.Sprintf("max-age=%v", setting.StrictTransportSecurityMaxAge)}
		if setting.StrictTransportSecurityPreload {
			strictHeaderValues = append(strictHeaderValues, "preload")
		}
		if setting.StrictTransportSecuritySubDomains {
			strictHeaderValues = append(strictHeaderValues, "includeSubDomains")
		}
		w.Header().Add("Strict-Transport-Security", strings.Join(strictHeaderValues, "; "))
	}

	if setting.ContentTypeProtectionHeader {
		w.Header().Add("X-Content-Type-Options", "nosniff")
	}

	if setting.XSSProtectionHeader {
		w.Header().Add("X-XSS-Protection", "1; mode=block")
	}
}

func AddNoCacheHeaders(w macaron.ResponseWriter) {
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "-1")
}

func AddXFrameOptionsDenyHeader(w macaron.ResponseWriter) {
	w.Header().Add("X-Frame-Options", "deny")
}
