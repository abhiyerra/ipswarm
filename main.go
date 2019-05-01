package main

import (
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

const TOPIC = "cloudless"

func main() {
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(strings.NewReader("hello world!"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("added %s", cid)

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "publish",
			Usage: "add jobs",
			Action: func(c *cli.Context) error {
				log.Println("Sending ", c.Args().First())
				sh.PubSubPublish(TOPIC, c.Args().First())

				return nil
			},
		},
		{
			Name:  "server",
			Usage: "run jobs",
			Action: func(c *cli.Context) error {
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

					ctx := context.Background()
					cli, err := client.NewEnvClient()
					if err != nil {
						log.Println(err)
						continue
					}

					reader, err := cli.ImagePull(ctx, string(msg.Data), types.ImagePullOptions{})
					if err != nil {
						log.Println(err)
						continue
					}
					io.Copy(os.Stdout, reader)

					resp, err := cli.ContainerCreate(ctx, &container.Config{
						Image: string(msg.Data),
						Cmd:   []string{},
						Tty:   true,
					}, nil, nil, "")
					if err != nil {
						log.Println(err)
						continue
					}

					if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
						log.Println(err)
						continue
					}

					out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
					if err != nil {
						panic(err)
					}

					io.Copy(os.Stdout, out)
				}
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// ipfs daemon --enable-pubsub-experiment
// go build
// ./cloudless server
// ./cloudless publish <dockerimage>
// # registry.hub.docker.com/library/nginx
