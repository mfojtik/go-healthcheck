package plugins

import (
	"log"

	"github.com/openshift/geard/cmd/switchns/namespace"
	docker "github.com/openshift/geard/docker"
	"github.com/openshift/go-healthcheck/api"
)

type FilePlugin struct{}

func (FilePlugin) Name() string {
	return "file"
}

func (FilePlugin) Perform(req *api.StatusRequest, ch chan bool) {

	client, err := docker.GetConnection(req.Socket)

	if err != nil {
		ch <- false
		return
	}

	filePath := "/tmp/.ready"

	if len(req.Args) > 1 {
		filePath = req.Args[1]
	}

	cmd := []string{"/bin/ls", filePath}

	containerNsPID, err := client.ChildProcessForContainer(req.Container)
	exitCode, err := namespace.RunIn(req.Container.Name, containerNsPID, cmd, req.Container.Config.Env)

	if req.Verbose {
		log.Printf("[file] The file '%s' exists == %t\n", filePath, exitCode)
	}

	if exitCode == 0 {
		ch <- true
	} else {
		ch <- false
	}

}
