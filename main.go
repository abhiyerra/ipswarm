package main

import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/bytecodealliance/wasmtime-go"
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

func WasmRun(wasmBytes []byte) string {
	store := wasmtime.NewStore(wasmtime.NewEngine())

	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModule(store, wasmBytes)
	check(err)

	// Next up we instantiate a module which is where we link in all our
	// imports. We've got one import so we pass that in here.
	instance, err := wasmtime.NewInstance(store, module, []*wasmtime.Extern{})
	check(err)

	// After we've instantiated we can lookup our `run` function and call
	// it.
	run := instance.GetExport("handler").Func()
	results, err := run.Call()
	check(err)

	return results.(string)
}

func WasmWasiRun(wasmBytes []byte) string {
	store := wasmtime.NewStore(wasmtime.NewEngine())

	config := wasmtime.NewWasiConfig()
	config.InheritStdout()
	// config.SetStdoutFile("./stdout")
	instance, err := wasmtime.NewWasiInstance(store, config, "wasi_snapshot_preview1")
	check(err)

	module, err := wasmtime.NewModule(store, wasmBytes)
	check(err)

	linker := wasmtime.NewLinker(store)
	linker.DefineWasi(instance)
	instance1, err := linker.Instantiate(module)

	run := instance1.GetExport("main").Func()
	_, err = run.Call(0, 0)

	// fmt.Println(results)
	// fmt.Println(err)

	return ""
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
