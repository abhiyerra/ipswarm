# IPFS-Compute

Run Serverless workloads loads powered by IPFS.

# Status: For the love of god don't use this in production. :)

# Use Cases

This could enable the following.

 - Internet Scale Data Processing
 - IoT Processing
 - Idle Processes
 - Distributed Web Processing

# Dependencies

 - Docker Daemon
 - IPFS With PubSub Enabled `ipfs daemon --enable-pubsub-experiment`

# Usage

## Worker

```
ipfs-compute worker
```

## TODO HTTP API Gateway

```
ipfs-compute api-gateway
```

## TODO Submit WASM Job

```
ipfs-compute submit --type wasm --event event.json --wasm-file <file>
```

 - Three functions will be made available within WASM:
    - [ ] ipfsComputeLog
    - [ ] ipfsComputeGet
    - [ ] ipfsComputeLs
    - [ ] ipfsComputeAdd
    - [ ] ipfsComputeCurl
 - [ ] WASM for isolation and security avoids external dependencies like docker
    - [ ] https://github.com/perlin-network/life
        - [ ] Limit network calls

## TODO Submit AWS Lamba Job

```
ipfs-compute submit --type aws-lambda  --zip-file <file> --runtime ruby2.5 --event event.json --handler name
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
 - [ ] Claim jobs so that only a single node runs the task at a time
 - [ ] Isolation of network resources
 - [ ] Trust. How do we ensure that the person sending requests is trustable?
 - [ ] Docker for some level of security and isolation
    - [ ] Maybe start isolation of the VM with Firecracker?


# License

The MIT License (MIT)

Copyright (c) 2019 Abhi Yerra

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE..