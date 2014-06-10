package plugins

import (
	"github.com/openshift/go-healthcheck/api"
	"labix.org/v2/mgo"
)

type MongoPlugin struct{}

func (MongoPlugin) Name() string {
	return "mongo"
}

func (MongoPlugin) Usage() string {
	return "Usage: TBD"
}

func (MongoPlugin) Check(req *api.StatusRequest, ch chan bool) {
	session, err := mgo.DialWithTimeout(req.Address, api.Timeout)
	if err != nil {
		ch <- false
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ch <- true
}
