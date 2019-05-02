package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

const TOPIC = "ipfs-compute"

var (
	sh *shell.Shell
)

type Job struct {
	Version string
	Type    string
	Docker  struct {
		Image string
		Cmd   []string
	}
}

func (j *Job) Start() error {
	switch j.Type {
	case DockerType:
		return j.startDocker()
	}

	return nil
}

func (j *Job) startDocker() error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, j.Docker.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: j.Docker.Image,
		Cmd:   j.Docker.Cmd,
		Tty:   true,
	}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)

	return nil
}

const (
	DockerType string = "docker"
)

func publish(c *cli.Context) error {
	j := Job{Version: "1.0"}

	switch c.String("type") {
	case DockerType:
		j.Type = DockerType
		j.Docker.Image = c.String("image")
		j.Docker.Cmd = c.StringSlice("cmd")

	}

	b, err := json.Marshal(j)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sending ", string(b))
	sh.PubSubPublish(TOPIC, string(b))

	return nil
}

func worker(c *cli.Context) error {
	sub, err := sh.PubSubSubscribe(TOPIC)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := sub.Next()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(string(msg.Data))

		var j Job
		if err = json.Unmarshal(msg.Data, &j); err != nil {
			log.Println(err)
			continue
		}

		j.Start()
	}
	return nil
}

func main() {
	sh = shell.NewShell("localhost:5001")
	cid, err := sh.Add(strings.NewReader("hello world!"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("added %s", cid)

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "submit",
			Usage: "add jobs",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type",
					Usage: "docker | aws-lambda | wasm",
				},
				cli.StringFlag{
					Name:  "image",
					Usage: "docker image",
				},
				cli.StringSliceFlag{
					Name:  "cmd",
					Usage: "command to send to docker",
				},
			},
			Action: publish,
		},
		{
			Name:   "worker",
			Usage:  "run jobs",
			Action: worker,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
