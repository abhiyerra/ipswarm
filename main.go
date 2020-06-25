package main

import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	host := "localhost:5001"
	if os.Getenv("INSIDE_DOCKER") == "true" {
		host = "host.docker.internal:5001"
	}

	r := mux.NewRouter()
	r.HandleFunc("/{ipfsRef}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		sh := shell.NewShell(host)

		o, err := sh.Cat(vars["ipfsRef"])
		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(o)
		if err != nil {
			fmt.Println(err)
		}

		tmpfile, err := ioutil.TempFile("", "ipswarm.*.wasm")
		if err != nil {
			log.Println(err)
		}
		defer tmpfile.Close()

		if _, err := tmpfile.Write(body); err != nil {
			log.Println(err)
		}

		cmd := exec.Command("./wasmtime", tmpfile.Name())
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}

		fmt.Fprintf(w, string(stdoutStderr))
	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8585", nil))
}
