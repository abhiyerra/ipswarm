package main

import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{ipfsRef}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		sh := shell.NewShell("localhost:5001")

		o, err := sh.Cat(vars["ipfsRef"])
		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(o)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, WasmWasiRun(body))
	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8585", nil))
}
