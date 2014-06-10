package api

import "time"

var Timeout = time.Duration(10 * time.Second)

type PluginRepository interface {
	Register(p Plugin)
	FindByName(name string) Plugin
	List() []Plugin
}

type Plugin interface {
	Name() string
	Usage() string
	Check(statusRequest *StatusRequest, ch chan bool)
}

type Repository struct {
	plugins []Plugin
}

func (r *Repository) Register(p Plugin) {
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
