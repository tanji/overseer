package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mgutz/logxi/v1"
	"github.com/tanji/overseer/alert"
	"github.com/tanji/overseer/config"
	"github.com/tanji/overseer/monitor"
	"github.com/tanji/overseer/server"
)

var (
	l         log.Logger
	rootAlert alert.Alert
	srvlist   server.Servers
	conf      config.Config
	mon       monitor.Monitors

	rootInterval time.Duration
)

func main() {
	l = log.New("overseer")
	l.SetLevel(log.LevelDebug)
	l.Info("Overseer has started")

	var err error
	conf, err = config.Load("config.toml")
	if err != nil {
		l.Fatal("Could not load config:", err)
	}

	l.Debug(fmt.Sprintf("%v", conf))

	conf.SetDefaults()

	rootAlert.Notify = conf.Alert.Notify
	rootAlert.Emails = conf.Alert.Emails
	rootAlert.Sender = conf.Alert.Sender
	rootAlert.Destination = conf.Alert.Destination
	rootAlert.Attempts = conf.Alert.Attempts

	rootInterval = time.Second * time.Duration(conf.Settings.Interval)

	srvlist = conf.InitServers()
	mon = conf.InitMonitors()

	for _, s := range srvlist {
		err := s.Open()
		if err != nil {
			l.Fatal("Could not initiate connection:", err)
		}
		l.Info("Opened connection", "server", s.Name)
	}

	c := time.Tick(rootInterval)
	for _ = range c {
		l.Debug("Entering monitoring loop")
		go alertMonitor()
		monitorUpDown()
		monitorStatus()
	}
}

// monitor server ping
func monitorUpDown() {
	wg := new(sync.WaitGroup)
	for _, sv := range srvlist {
		wg.Add(1)
		go func(s *server.Server) {
			defer wg.Done()
			if err := s.Check(); err != nil {
				l.Warn("Error checking server", "name", s.Name, "err", err)
				if s.FailureCount < rootAlert.Attempts {
					s.FailureCount++
					l.Warn("Incrementing failure counter for server", "name", s.Name, "FailureCount", s.FailureCount)
				}
			} else if s.StatusDown {
				l.Info("Recovery event for server", "name", s.Name)
				s.StatusDown = false
				s.FailureCount = 0
			}
		}(sv)
	}
	wg.Wait()
}

// monitor status vars
func monitorStatus() {
	wg := new(sync.WaitGroup)
	for _, s := range srvlist {
		if s.StatusDown || s.FailureCount > 0 {
			continue
		}
		wg.Add(1)
		go func(sv *server.Server) {
			defer wg.Done()
			st, err := sv.GetStatus()
			if err == nil {
				l.Info("Collected statuses for server", "name", sv.Name, "connections", st["threads_connected"])
			} else {
				l.Warn("Could not get status for server", "name", sv.Name)
			}
		}(s)
	}
	wg.Wait()
}

func alertMonitor() {
	for name, s := range srvlist {
		if s.FailureCount >= 5 && !s.StatusDown {
			l.Info("Sending Alert", "server", name)
			a := rootAlert
			l.Debug(fmt.Sprintf("%v", a))
			a.Origin = s.Name
			a.Type = "Ping"
			a.Value = "Down"
			l.Debug(fmt.Sprintf("%v", conf.Smtp))
			err := a.Email(conf.Smtp)
			if err != nil {
				l.Warn("Could not send alert email", "err", err)
			}
			s.StatusDown = true
		}
	}
	return
}
