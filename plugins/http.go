package plugins

import (
	"log"
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

	if req.Port == "" {
		req.Port = "8080"
	}

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	if req.Verbose {
		log.Println("GET " + "http://" + req.Address + ":" + req.Port + "/")
	}

	resp, err := client.Get("http://" + req.Address + ":" + req.Port + "/")

	if err != nil {
		ch <- false
		return
	} else {
		if req.Verbose {
			log.Printf("HTTP Response: %s\n", resp.Status)
		}
		if resp.Status == "200 OK" {
			ch <- true
		} else {
			ch <- false
		}
	}
}
