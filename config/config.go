package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/tanji/overseer/monitor"
	"github.com/tanji/overseer/server"
)

// Config defines a configuration file structure
type Config struct {
	Settings Settings
	Servers  map[string]server.Server
	Alert    Alert
	Smtp     SMTP
	Monitor  []monitor.Monitor
}

type Alert struct {
	Notify      bool
	Sender      string
	Emails      []string
	Destination string
	Attempts    uint
}

type SMTP struct {
	Username string
	Password string
	Server   string
}

type Settings struct {
	Interval uint64
}

var (
	defaultInterval uint64 = 10
)

func Load(cf string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(cf, &config)
	if err != nil {
		return config, fmt.Errorf("Could not read configuration file: %s", err)
	}
	return config, nil
}

func (cf *Config) InitServers() server.Servers {
	sl := make(server.Servers, len(cf.Servers))
	k := 0
	for n, s := range cf.Servers {
		ns := server.Server{}
		ns = s
		ns.Name = n
		sl[k] = &ns
		k++
	}
	return sl
}

func (cf *Config) InitMonitors() monitor.Monitors {
	m := make(monitor.Monitors, len(cf.Monitor))
	for k, v := range cf.Monitor {
		m[k] = v
	}
	return m
}

func (cf *Config) SetDefaults() {
	if cf.Settings.Interval == 0 {
		cf.Settings.Interval = defaultInterval
	}
	return
}
