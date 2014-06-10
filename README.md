go-healthcheck
------------------

This tool can be used to health-check Docker container using the network
connection or via inspecting the Docker container logs.

### Usage

```
# See if the HTTP server is running on port 8080

$ healthck status $CID -p 8080 -P http
$ echo $? # -> 0

# Check if the MongoDB accepts connection on 27017

$ healthck status $CID -p 27017 -P mongo
$ echo $? # -> 0

# Generic TCP check if the port 9292 is open

$ healthck status $CID -p 9292 -P tcp
$ echo $? # -> 1

```

The `$CID` variable represents the Docker container ID.
