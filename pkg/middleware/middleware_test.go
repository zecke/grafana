package middleware

import (
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/macaron.v1"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/auth"
	"github.com/grafana/grafana/pkg/setting"
)

const errorTemplate = "error-template"

func mockGetTime() {
	var timeSeed int64
	getTime = func() time.Time {
		fakeNow := time.Unix(timeSeed, 0)
		timeSeed++
		return fakeNow
	}
}

func resetGetTime() {
	getTime = time.Now
}

func TestMiddlewareContext(t *testing.T) {
	setting.ERR_TEMPLATE_NAME = errorTemplate

	Convey("Given the grafana middleware", t, func() {
		middlewareScenario(t, "middleware should add context to injector", func(sc *scenarioContext) {
			sc.fakeReq("GET", "/").exec()
			So(sc.context, ShouldNotBeNil)
		})

		middlewareScenario(t, "Default middleware should allow get request", func(sc *scenarioContext) {
			sc.fakeReq("GET", "/").exec()
			So(sc.resp.Code, ShouldEqual, 200)
		})

		middlewareScenario(t, "middleware should add Cache-Control header for requests to API", func(sc *scenarioContext) {
			sc.fakeReq("GET", "/api/search").exec()
			So(sc.resp.Header().Get("Cache-Control"), ShouldEqual, "no-cache")
			So(sc.resp.Header().Get("Pragma"), ShouldEqual, "no-cache")
			So(sc.resp.Header().Get("Expires"), ShouldEqual, "-1")
		})

		middlewareScenario(t, "middleware should not add Cache-Control header for requests to datasource proxy API", func(sc *scenarioContext) {
			sc.fakeReq("GET", "/api/datasources/proxy/1/test").exec()
			So(sc.resp.Header().Get("Cache-Control"), ShouldBeEmpty)
			So(sc.resp.Header().Get("Pragma"), ShouldBeEmpty)
			So(sc.resp.Header().Get("Expires"), ShouldBeEmpty)
		})

		middlewareScenario(t, "middleware should add Cache-Control header for requests with html response", func(sc *scenarioContext) {
			sc.handler(func(c *models.ReqContext) {
				data := &dtos.IndexViewData{
					User:     &dtos.CurrentUser{},
					Settings: map[string]interface{}{},
					NavTree:  []*dtos.NavLink{},
				}
				c.HTML(200, "index-template", data)
			})
			sc.fakeReq("GET", "/").exec()
			So(sc.resp.Code, ShouldEqual, 200)
			So(sc.resp.Header().Get("Cache-Control"), ShouldEqual, "no-cache")
			So(sc.resp.Header().Get("Pragma"), ShouldEqual, "no-cache")
			So(sc.resp.Header().Get("Expires"), ShouldEqual, "-1")
		})

		middlewareScenario(t, "middleware should add X-Frame-Options header with deny for request when not allowing embedding", func(sc *scenarioContext) {
			sc.fakeReq("GET", "/api/search").exec()
			So(sc.resp.Header().Get("X-Frame-Options"), ShouldEqual, "deny")
		})

		middlewareScenario(t, "middleware should not add X-Frame-Options header for request when allowing embedding", func(sc *scenarioContext) {
			setting.AllowEmbedding = true
			sc.fakeReq("GET", "/api/search").exec()
			So(sc.resp.Header().Get("X-Frame-Options"), ShouldBeEmpty)
		})
	})
}

func middlewareScenario(t *testing.T, desc string, fn scenarioFunc) {
	Convey(desc, func() {
		defer bus.ClearBusHandlers()

		setting.LoginCookieName = "grafana_session"
		setting.LoginMaxLifetimeDays = 30

		sc := &scenarioContext{}

		viewsPath, _ := filepath.Abs("../../public/views")

		sc.m = macaron.New()
		sc.m.Use(AddDefaultResponseHeaders())
		sc.m.Use(macaron.Renderer(macaron.RenderOptions{
			Directory: viewsPath,
			Delims:    macaron.Delims{Left: "[[", Right: "]]"},
		}))

		sc.userAuthTokenService = auth.NewFakeUserAuthTokenService()
		sc.remoteCacheService = remotecache.NewFakeStore(t)

		sc.m.Use(GetContextHandler(sc.userAuthTokenService, sc.remoteCacheService))

		sc.m.Use(OrgRedirect())

		sc.defaultHandler = func(c *models.ReqContext) {
			sc.context = c
			if sc.handlerFunc != nil {
				sc.handlerFunc(sc.context)
			}
		}

		sc.m.Get("/", sc.defaultHandler)

		fn(sc)
	})
}
