# IPFS-Compute

Run Serverless workloads loads powered by IPFS.

# Status

Experimental way to run Docker images by subscribing to IPFS pubsub.

For the love of god don't use this in production. :)


Run a docker daemon

```
ipfs daemon --enable-pubsub-experiment
```

# Usage

## Worker

```
ipfs-compute worker
```

## Submit AWS Lambda Job

```
ipfs-compute submit --type aws-lambda --runtime ruby2.5 --event event.json --zip <file>
```

## Submit Docker Job

Full canonical registry name. eg registry.hub.docker.com/library/nginx

```
ipfs-compute submit --type docker --image <image> --cmd <cmd>
```

```
ipfs-compute submit --type docker --image registry.hub.docker.com/library/nginx
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