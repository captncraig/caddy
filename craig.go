package main

import (
	"fmt"
	"net/http"

	"github.com/abiosoft/caddy-git"
	"github.com/captncraig/caddy-realip"
	"github.com/captncraig/caddy-stats"
	"github.com/mholt/caddy/caddy"
	"github.com/mholt/caddy/caddy/setup"
	"github.com/mholt/caddy/middleware"

	"github.com/captncraig/hugomail/web"
)

func init() {
	caddy.RegisterDirective("git", git.Setup, "internal")     //can go toward bottom
	caddy.RegisterDirective("realip", realip.Setup, "tls")    //goes almost at very top
	caddy.RegisterDirective("stats", stats.Setup, "shutdown") //as high as I dare

	hugomailSetup := appToDirective(
		func() interface{} { return &web.Config{} },
		func(i interface{}) *http.ServeMux { return i.(*web.Config).Serve() },
	)
	caddy.RegisterDirective("hugomail", hugomailSetup, "") //last is fine
}

// convert a new struct factory and servemux creator into an official middleware
func appToDirective(newConf func() interface{}, getMux func(interface{}) *http.ServeMux) caddy.SetupFunc {
	return func(c *setup.Controller) (middleware.Middleware, error) {
		conf := newConf()
		err := c.Unmarshal(conf)
		if err != nil {
			return nil, err
		}
		fmt.Println(conf)
		mux := getMux(conf)
		return func(next middleware.Handler) middleware.Handler {
			return myMux{mux}
		}, nil
	}
}

type myMux struct {
	*http.ServeMux
}

func (m myMux) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	m.ServeMux.ServeHTTP(w, r)
	return 0, nil
}
