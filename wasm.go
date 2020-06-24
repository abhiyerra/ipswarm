package main

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

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

	return "'"
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
