package middleware

import (
	"testing"

	"github.com/grafana/grafana/pkg/setting"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddlewareAuth(t *testing.T) {

	Convey("Given the grafana middleware", t, func() {
		reqSignIn := Auth(&AuthOptions{ReqSignedIn: true})

		middlewareScenario(t, "ReqSignIn true and unauthenticated request", func(sc *scenarioContext) {
			sc.m.Get("/secure", reqSignIn, sc.defaultHandler)

			sc.fakeReq("GET", "/secure").exec()

			Convey("Should not redirect", func() {
				So(sc.resp.Code, ShouldEqual, 200)
			})
		})

		middlewareScenario(t, "ReqSignIn true and no such thing as unauthenticated API request", func(sc *scenarioContext) {
			sc.m.Get("/api/secure", reqSignIn, sc.defaultHandler)

			sc.fakeReq("GET", "/api/secure").exec()

			Convey("Should return 401", func() {
				So(sc.resp.Code, ShouldEqual, 200)
			})
		})

		Convey("snapshot public mode or signed in", func() {
			middlewareScenario(t, "Snapshot public mode disabled and no such thing as unauthenticated request", func(sc *scenarioContext) {
				sc.m.Get("/api/snapshot", SnapshotPublicModeOrSignedIn(), sc.defaultHandler)
				sc.fakeReq("GET", "/api/snapshot").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})

			middlewareScenario(t, "Snapshot public mode enabled and unauthenticated request should return 200", func(sc *scenarioContext) {
				setting.SnapshotPublicMode = true
				sc.m.Get("/api/snapshot", SnapshotPublicModeOrSignedIn(), sc.defaultHandler)
				sc.fakeReq("GET", "/api/snapshot").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})
		})
	})
}
