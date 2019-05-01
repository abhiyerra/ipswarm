# IPFS-Compute

Run Serverless loads powered by IPFS.

# Status

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
 - [ ] WASM for isolation and security avoids external dependencies like docker
    - [ ] https://github.com/perlin-network/life
        - [ ] Limit network calls
 - [ ] Maybe start isolation of the VM with Firecracker?
 - [ ] Claim jobs so that only a single node runs the task at a time
 - [ ] Isolation of network resources
 - [ ] Trust. How do we ensure that the person sending requests is trustable?