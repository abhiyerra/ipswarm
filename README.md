# Cloudless

Experimental way to run Docker images by subscribing to IPFS pubsub.

For the love of god don't use this in production. :)


Run a docker daemon

```
ipfs daemon --enable-pubsub-experiment
go build
```

# Usage

```
./cloudless server
./cloudless publish <dockerimage>
# Full canoncical registry name. eg registry.hub.docker.com/library/nginx
```


# Future

 - [ ] Handle that the task has been finished
 - [ ] WASM for isolation and security
 - [ ] Maybe start isolation of the VM with Firecracker?