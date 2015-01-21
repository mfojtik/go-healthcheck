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

func (HttpPlugin) Perform(req *api.StatusRequest, ch chan bool) {

	uri := "/"

	if len(req.Args) > 1 {
		uri = req.Args[1]
	}

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
		log.Println("GET " + "http://" + req.Address + ":" + req.Port + uri)
	}

	resp, err := client.Get("http://" + req.Address + ":" + req.Port + uri)

	if err != nil {
		if req.Verbose {
			log.Printf("HTTP Error: %s\n", err)
		}
		ch <- false
		return
	}

	if req.Verbose {
		log.Printf("HTTP Response: %s\n", resp.Status)
	}

	// FIXME: Support more HTTP status codes
	//
	if resp.StatusCode != http.StatusOK {
		ch <- true
	} else {
		ch <- false
	}
}
