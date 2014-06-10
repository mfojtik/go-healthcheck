package plugins

import (
	"net"

	"github.com/openshift/go-healthcheck/api"
)

type TcpPlugin struct{}

func (TcpPlugin) Name() string {
	return "tcp"
}

func (TcpPlugin) Usage() string {
	return "Usage: TBD"
}

func (TcpPlugin) Check(req *api.StatusRequest, ch chan bool) {
	c, err := net.Dial("tcp", req.Address+":"+req.Port)
	if err != nil {
		ch <- false
		return
	}
	defer c.Close()
	ch <- true
}
