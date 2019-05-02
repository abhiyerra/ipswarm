package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

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

const (
	DockerType    string = "docker"
	AwsLambdaType        = "aws-lambda"
)

type Job struct {
	Version string
	Type    string
	Docker  struct {
		Image string
		Cmd   []string
	}

	AwsLambda struct {
		ZipFileCid string
		Runtime    string
		Handler    string
		Event      string
	}

	// Wasm
	// Event
	// Context
	// AllowFlood
}

func (j *Job) Start() error {
	switch j.Type {
	case DockerType:
		return j.startDocker()
	case AwsLambdaType:
		return j.startAwsLambda()
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

func (j *Job) startAwsLambda() error {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	// defer os.Remove(tmpfile.Name()) // clean up

	sh.Get(j.AwsLambda.ZipFileCid, tmpfile.Name())

	dir, err := ioutil.TempDir("", j.AwsLambda.ZipFileCid)
	if err != nil {
		log.Fatal(err)
	}
	// defer os.RemoveAll(dir) // clean up

	filenames, err := unzip(tmpfile.Name(), dir)

	log.Println(err)
	log.Println(filenames)
	log.Println(dir)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, "registry.hub.docker.com/lambci/lambda:"+j.AwsLambda.Runtime, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: j.Docker.Image,
		Tty:   true,
		Cmd: []string{
			j.AwsLambda.Handler,
			j.AwsLambda.Event,
		},
	}, &container.HostConfig{
		Binds: []string{
			dir + ":/tasks",
		},
	}, nil, "")
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

func publish(c *cli.Context) error {
	j := Job{Version: "1.0"}

	switch c.String("type") {
	case DockerType:
		j.Type = DockerType
		j.Docker.Image = c.String("image")
		j.Docker.Cmd = c.StringSlice("cmd")
	case AwsLambdaType:
		j.Type = AwsLambdaType
		j.AwsLambda.Runtime = c.String("runtime")
		j.AwsLambda.Event = c.String("event")

		content, err := ioutil.ReadFile(c.String("zip-file"))
		if err != nil {
			log.Fatal(err)
		}

		cid, err := sh.Add(bytes.NewReader(content))
		if err != nil {
			log.Fatal(err)
		}
		j.AwsLambda.ZipFileCid = cid
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

		log.Println(j.Start())
	}
	return nil
}

func apiGateway() {
	// TODO:
	// Http Request comes in and is sent into pub sub
	// Block on object get
	// Can be an api for any request
	// /ipfs-request/path
}

func main() {
	sh = shell.NewShell("localhost:5001")

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
				cli.StringFlag{
					Name:  "runtime",
					Value: "nodejs8.10",
					Usage: "AWS Lambda Runtime",
				},
				cli.StringFlag{
					Name:  "event",
					Usage: "An event that is specified as a json file.",
				},
				cli.StringFlag{
					Name:  "zip-file",
					Usage: "File containing the AWS Lambda code",
				},
				cli.StringFlag{
					Name:  "handler",
					Usage: "handler for AWS Lambda",
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
