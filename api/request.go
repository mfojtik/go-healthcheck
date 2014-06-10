package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

type Request struct {
	Address string
	Verbose bool
}

type StatusRequest struct {
	Container  *docker.Container
	Plugins    []*Plugin
	Verbose    bool
	Address    string
	Port       string
	Args       []string
	Socket     string
	PluginList string
}

func (s *StatusRequest) SetArgs(args []string) {
	s.Args = args
}

func (s *StatusRequest) FindContainer(containerId string) (err error) {
	client, err := docker.NewClient(s.Socket)
	if err != nil {
		return
	}
	s.Container, err = client.InspectContainer(containerId)
	s.Address = s.Container.NetworkSettings.IPAddress
	return
}

func (s *StatusRequest) InitializePlugins(repo *Repository) {
	pluginNames := strings.Split(s.PluginList, ",")
	for i := 0; i < len(pluginNames); i++ {
		if plugin := repo.FindByName(pluginNames[i]); plugin == nil {
			fmt.Printf("The plugin '%s' is invalid\n", pluginNames[i])
			continue
		} else {
			if s.Verbose {
				log.Printf("Enabling %s plugin", pluginNames[i])
			}
			s.Plugins = append(s.Plugins, &plugin)
		}
	}
}

func (s *StatusRequest) Execute() bool {
	healthChan := make(chan bool)
	for i := 0; i < len(s.Plugins); i++ {
		go (*s.Plugins[i]).Check(s, healthChan)
	}
	for i := 0; i < len(s.Plugins); i++ {
		ok := <-healthChan
		if s.Verbose {
			log.Printf("<- %t\n", ok)
		}
		if !ok {
			return false
		}
	}
	return true
}
