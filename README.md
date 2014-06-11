go-healthcheck
------------------

This tool can be used to perform a health check operations against the Docker
containers. It has support for plugins, thus it is easely extensible.

There are three plugins atm.:

* **http** - Perform a HTTP health check against container.
* **mongo** - Perform a MongoDB connection test against container
* **tcp** - Perform a generic TCP connection test

The connection based plugins will guess the port from the container by default
(from `EXPOSE`).

### Usage

```

# Check if the container respond to HTTP requests
$ healthck status 3e3a1ebbb4dd -P http

# Check if the MongoDB container is ready for connections
$ healthck status 3e3a1ebbb4dd -P mongo

# A generic TCP check
$ healthck status 3e3a1ebbb4dd -P tcp
```

### Extending

To add a new plugin, all you have to do is to create a `struct` that implements
two methods: `Name()` and `Perform()`. See the existing plugins for an example.


### License

Apache Software License (ASL) 2.0.
