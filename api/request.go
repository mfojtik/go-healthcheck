package api

import (
	"log"

	"github.com/fsouza/go-dockerclient"
)

type StatusRequest struct {
	Container *docker.Container
	Verbose   bool
	Address   string
	Port      string
	Args      []string
	Socket    string
	Plugins   []*Plugin
}

func (s *StatusRequest) SetArgs(args []string) {
	s.Args = args
}

func (s *StatusRequest) SetPlugins(plugins []*Plugin) {
	s.Plugins = plugins
}

func (s *StatusRequest) FindContainer(containerId string) (err error) {
	client, err := docker.NewClient(s.Socket)
	if err != nil {
		return
	}
	s.Container, err = client.InspectContainer(containerId)
	s.Address = s.Container.NetworkSettings.IPAddress
	if s.Port == "" {
		for key, _ := range s.Container.Config.ExposedPorts {
			s.Port = key.Port()
		}
		if s.Verbose {
			log.Printf("Using port exposed by container: %s\n", s.Port)
		}
	}
	return
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
