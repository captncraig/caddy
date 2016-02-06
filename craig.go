package main

import (
	"github.com/abiosoft/caddy-git"
	"github.com/captncraig/caddy-realip"
	"github.com/captncraig/caddy-stats"
	"github.com/mholt/caddy/caddy"
)

func init() {
	caddy.RegisterDirective("git", git.Setup, "internal")     //can go toward bottom
	caddy.RegisterDirective("realip", realip.Setup, "tls")    //goes almost at very top
	caddy.RegisterDirective("stats", stats.Setup, "shutdown") //as high as I dare
}
