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
	"github.com/mitchellh/mapstructure"

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
		m, err := readMap(c)
		if err != nil {
			return nil, err
		}
		conf := newConf()
		err = mapstructure.WeakDecode(m, conf)
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

//read single key/value pairs from block into a map
//expects:
// - 1 and only one directive with this name
// - No top level arguments
// - No repeat config lines
// really for mapping simple key/value into a struct
func readMap(c *setup.Controller) (map[string]interface{}, error) {

	r := map[string]interface{}{}
	if !c.Next() {
		return nil, c.Err("Expect directive")
	}
	if args := c.RemainingArgs(); len(args) != 0 {
		return nil, c.ArgErr()
	}
	for c.NextBlock() {
		k := c.Val()
		if r[k] != nil {
			return nil, c.Errf("Repeat key %s", k)
		}
		args := c.RemainingArgs()
		if len(args) != 1 {
			return nil, c.ArgErr()
		}
		r[k] = args[0]
	}
	if c.Next() {
		return nil, c.Err("Expect only one directive")
	}
	return r, nil
}

type myMux struct {
	*http.ServeMux
}

func (m myMux) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	m.ServeMux.ServeHTTP(w, r)
	return 0, nil
}
