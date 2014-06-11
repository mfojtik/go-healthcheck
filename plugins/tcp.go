package plugins

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/openshift/go-healthcheck/api"
)

type TcpPlugin struct{}

func (TcpPlugin) Name() string {
	return "tcp"
}

func (TcpPlugin) Perform(req *api.StatusRequest, ch chan bool) {
	if req.Verbose {
		log.Println("Probing TCP port " + req.Address + ":" + req.Port)
	}
	// Connect to the remote TCP port
	//
	c, err := net.Dial("tcp", req.Address+":"+req.Port)
	if err != nil {
		if req.Verbose {
			log.Println("net.Dial: %s", err)
		}
		ch <- false
		return
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(time.Second * 10))

	// Test writing to TCP connection
	//
	fmt.Fprintf(c, "GET / HTTP/1.0\r\n\r\n")

	// Read 1 byte from the TCP connection.
	//
	b := make([]byte, 1)
	if _, err := c.Read(b); err != nil {
		if req.Verbose {
			log.Println(err)
		}
		ch <- false
		return
	}

	ch <- true
}
