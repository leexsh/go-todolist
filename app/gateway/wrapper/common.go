package wrapper

import (
	"sync"
	"time"
)

type Config struct {
	Namespace              string
	Timeout                time.Duration
	MaxConcurrentRequests  int
	RequestVolumeThreshold uint64
	SleepWindow            time.Duration
	ErrorPercentThreshold  int
}

type Group struct {
	sync.RWMutex
	namespace string
	settings   map[string]bool
	conf      *Config
}

var (
	_mu   sync.RWMutex
	_conf = &Config{
		Namespace:              "default",
		Timeout:                time.Duration(hystrix.DefaultTimeout),
		MaxConcurrentRequests:  hystrix.DefaultMaxConcurrent,
		RequestVolumeThreshold: uint64(hystrix.DefaultVolumeThreshold),
		SleepWindow:            time.Duration(hystrix.DefaultSleepWindow),
		ErrorPercentThreshold:  hystrix.DefaultErrorPercentThreshold,
	}
}

func (conf *Config) fix() {
	if conf.Namespace == "" {
		conf.Namespace = "default"
	}
	if conf.Timeout <= 0 {
		conf.Timeout = time.Duration(hystrix.DefaultTimeout)
	}
	if conf.MaxConcurrentRequests <= 0 {
		conf.MaxConcurrentRequests = hystrix.DefaultMaxConcurrent
	}
	if conf.RequestVolumeThreshold <= 0 {
		conf.RequestVolumeThreshold = uint64(hystrix.DefaultVolumeThreshold)
	}
	if conf.SleepWindow == 0 {
		conf.SleepWindow = time.Duration(hystrix.DefaultSleepWindow)
	}
	if conf.ErrorPercentThreshold <= 0 {
		conf.ErrorPercentThreshold = hystrix.DefaultErrorPercentThreshold
	}
}

func NewGroup(conf *Config) *Group {
	if conf == nil {
		_mu.RLock()
		conf = _conf
		_mu.RUnlock()
	} else {
		conf.fix()
	}
	return &Group{
		namespace: conf.Namespace,
		settings:  make(map[string]bool),
		conf:      conf,
	}
}

func (g *Group) Reload(conf *Config) {
	if conf == nil {
		return
	}
	conf.fix()
	g.Lock()
	g.conf = conf
	g.Unlock()
}

func (g *Group) Do(name string, run func() error) (err error) {
	name = g.namespace + "-" + name
	g.setBreakerConfig(name)
	return hystrix.Do(name, func() error {
		return run()
	}, nil)
}


func (g *Group) setBreakerConfig(name string) {
	if _, ok := g.settings[name]; !ok {
		g.Lock()
		defer g.Unlock()

		if _, ok := g.settings[name]; !ok {
			hystrix.ConfigureCommand(name, hystrix.CommandConfig{
				Timeout:                int(time.Duration(g.conf.Timeout) / time.Millisecond),
				MaxConcurrentRequests:  g.conf.MaxConcurrentRequests,
				RequestVolumeThreshold: int(g.conf.RequestVolumeThreshold),
				SleepWindow:            int(time.Duration(g.conf.SleepWindow) / time.Millisecond),
				ErrorPercentThreshold:  g.conf.ErrorPercentThreshold,
			})

			copy := make(map[string]bool)
			for key, val := range g.settings {
				copy[key] = val
			}
			copy[name] = true
			g.settings = copy
		}
	}
}
