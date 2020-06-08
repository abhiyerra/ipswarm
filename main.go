package main

import "C"

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli"
)

const (
	TOPIC = "ipswarm"
)

var (
	sh *shell.Shell
)

func Publish(c *cli.Context) error {
	j := Job{
		WasmCid: c.String("wasm-cid"),
		Handler: c.String("handler"),
	}

	b, err := json.Marshal(j)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sending ", string(b))
	sh.PubSubPublish(TOPIC, string(b))

	return nil
}

func Worker(c *cli.Context) error {
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

func APIGateway() {
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8585", nil))
}

func main() {
	sh = shell.NewShell("localhost:5001")

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "submit",
			Usage: "submit jobs to the cluster",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wasm-cid",
					Usage: "wasm file cid on ipfs",
				},
				cli.StringFlag{
					Name:  "handler",
					Usage: "handler to call in the WASM library. It should not have empty parametes",
				},
			},
			Action: Publish,
		},
		{
			Name:   "worker",
			Usage:  "run jobs",
			Action: Worker,
		},
		{
			Name:   "api-gateway",
			Usage:  "HTTP API Gateway",
			Action: APIGateway,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
