fsinformer
=======================================================

**fsinformer** is extension to [github.com/fsnotify](https://github.com/fsnotify/fsnotify) that provide convenient interface to register
event handlers for on-disk file changes. In addition to fsnotify it also implements a "relist" behavior (similar to Kubernetes informers) that
periodically resync the monitored files independent of fsnotify.

To install:

```console
go get -u github.com/mfojtik/fsinformer
```

For complete example look at `example.go` file.

License
-------

Licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/).
