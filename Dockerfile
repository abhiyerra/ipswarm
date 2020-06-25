FROM golang:1-stretch

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build .

RUN apt-get update -y && apt-get install   -y xz-utils
RUN curl https://wasmtime.dev/install.sh -sSf | bash

RUN mv /root/.wasmtime/bin/wasmtime .

EXPOSE 8585

ENV INSIDE_DOCKER true

CMD ["./ipswarm"]