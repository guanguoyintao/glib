package uhttp

import (
	"time"
)

type HttpConf struct {
	TimeoutString       string        `yaml:"timeout" json:"timeout"`
	Timeout             time.Duration `yaml:"-" json:"-"`
	MaxIdleConns        int32         `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxConnsPerHost     int32         `yaml:"max_conns_per_host" json:"max_conns_per_host"`
	MaxIdleConnsPerHost int32         `yaml:"max_idle_conns_per_host" json:"max_idle_conns_per_host"`
	Debug               bool          `yaml:"debug" json:"debug"`
}

type HttpConfig struct {
	Http HttpConf `yaml:"http" json:"http"`
}
