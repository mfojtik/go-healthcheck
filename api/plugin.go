package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var Timeout = time.Duration(10 * time.Second)

type PluginRepository interface {
	Register(p Plugin)
	FindByName(name string) Plugin
	List() []Plugin
}

type Plugin interface {
	Name() string
	Perform(statusRequest *StatusRequest, ch chan bool)
}

type Repository struct {
	plugins []Plugin
}

func ParsePlugins(s string, repo Repository) (plugins []*Plugin) {
	if s == "" {
		fmt.Printf("Please specify a list of plugins you want to use (-P)\n")
		return
	}
	pluginNames := strings.Split(s, ",")
	for i := 0; i < len(pluginNames); i++ {
		if plugin := repo.FindByName(pluginNames[i]); plugin == nil {
			log.Printf("The plugin '%s' is invalid\n", pluginNames[i])
			continue
		} else {
			plugins = append(plugins, &plugin)
		}
	}
	return plugins
}

func (r *Repository) Add(p Plugin) {
	r.plugins = append(r.plugins, p)
}

func (r *Repository) List() []Plugin {
	return r.plugins
}

func (r *Repository) FindByName(name string) Plugin {
	for i := 0; i < len(r.plugins); i++ {
		if r.plugins[i].Name() == name {
			return r.plugins[i]
		}
	}
	return nil
}
