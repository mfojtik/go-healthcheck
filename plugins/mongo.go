package plugins

import (
	"log"

	"github.com/openshift/go-healthcheck/api"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoPlugin struct{}

func (MongoPlugin) Name() string {
	return "mongo"
}

func (MongoPlugin) Usage() string {
	return "Usage: TBD"
}

func (MongoPlugin) Check(req *api.StatusRequest, ch chan bool) {
	url := req.Address + ":" + req.Port

	if req.Verbose {
		log.Printf("Connecting to MongoDB (%s) \n", url)
	}
	session, err := mgo.DialWithTimeout(url, api.Timeout)
	if err != nil {
		ch <- false
		return
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	// Use 'test' database by default. If you want to specify another
	// database, then pass it as an argument to the 'status' command
	//
	database := "test"
	if len(req.Args) > 1 {
		database = req.Args[1]
	}

	// Execute 'db.stats' command on selected database
	//
	retval := &bson.M{}
	if err := session.DB(database).Run("dbStats", &retval); err != nil {
		if req.Verbose {
			log.Printf("Failed to query MongoDB: %s\n", err)
		}
		ch <- false
		return
	}

	if (*retval)["ok"] != 1 {
		ch <- true
	} else {
		if req.Verbose {
			log.Printf("The database is not ready yet: %q", *retval)
		}
		ch <- false
	}
}
