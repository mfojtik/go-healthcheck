package plugins

import (
	"fmt"
	"net"
	"net/http"

	"github.com/openshift/go-healthcheck/api"
)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, api.Timeout)
}

type HttpPlugin struct{}

func (HttpPlugin) Name() string {
	return "http"
}

func (HttpPlugin) Usage() string {
	return "Usage: TBD"
}

func (HttpPlugin) Check(req *api.StatusRequest, ch chan bool) {

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get("http://" + req.Address + ":" + req.Port + "/")

	if req.Verbose {
		fmt.Println(resp.Status)
	}

	if err != nil {
		ch <- false
		return
	} else {
		ch <- true
	}
}
