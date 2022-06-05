package main

import (
	"path/filepath"
	"time"
)

type defaulter interface {
	SetDefault(key string, value interface{})
}

func setDefaults(d defaulter) {
	d.SetDefault("rpc.host", "127.0.0.1")
	d.SetDefault("rpc.user", "")
	d.SetDefault("rpc.password", "")
	d.SetDefault("rpc.https", false)
	d.SetDefault("rpc.port", uint32(0))
	d.SetDefault("rpc.uri", "")
	d.SetDefault("rpc.httptimeout", time.Duration(0))
	d.SetDefault("rpc.useragent", "transmission-monitor/"+version)
	d.SetDefault("rpc.debug", false)
	defDBPath := filepath.Join(home, ".transmission-monitor", "db")
	d.SetDefault("db.path", defDBPath)
	d.SetDefault("notify.command", []string{})
	d.SetDefault("notify.timeout", 30*time.Second)
}
