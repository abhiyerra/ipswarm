package main

// #include <stdlib.h>
//
// extern void ipfsComputeLog(void *context, _GoString_ msg);
// extern int64_t ipfsComputeGet(void *context);
// extern int64_t ipfsComputeLs(void *context);
// extern int64_t ipfsComputeAdd(void *context);
// extern int64_t ipfsComputeCurl(void *context);
import "C"

import (
	"fmt"
	"log"
	"unsafe"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

//export ipfsComputeLog
func ipfsComputeLog(context unsafe.Pointer, msg string) {
	log.Println(msg)
}

//export ipfsComputeGet
func ipfsComputeGet(context unsafe.Pointer) int64 {
	// ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	// msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	// msg := vm.Memory[ptr : ptr+msgLen]
	// fmt.Printf("[app] %s\n", string(msg))
	return 0
}

//export ipfsComputeLs
func ipfsComputeLs(context unsafe.Pointer) int64 {
	// ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	// msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	// msg := vm.Memory[ptr : ptr+msgLen]
	// fmt.Printf("[app] %s\n", string(msg))
	return 0
}

//export ipfsComputeAdd
func ipfsComputeAdd(context unsafe.Pointer) int64 {
	// ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	// msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	// msg := vm.Memory[ptr : ptr+msgLen]
	// fmt.Printf("[app] %s\n", string(msg))
	return 0
}

//export ipfsComputeCurl
func ipfsComputeCurl(context unsafe.Pointer) int64 {
	// ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	// msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	// msg := vm.Memory[ptr : ptr+msgLen]
	// fmt.Printf("[app] %s\n", string(msg))
	return 0
}

type Job struct {
	WasmCid string
	Handler string
}

func (j *Job) Start() {
	imports := wasm.NewImports()
	imports.Append("ipfsComputeLog", ipfsComputeLog, C.ipfsComputeLog)
	imports.Append("ipfsComputeGet", ipfsComputeGet, C.ipfsComputeGet)
	imports.Append("ipfsComputeLs", ipfsComputeLs, C.ipfsComputeLs)
	imports.Append("ipfsComputeAdd", ipfsComputeAdd, C.ipfsComputeAdd)
	imports.Append("ipfsComputeCurl", ipfsComputeCurl, C.ipfsComputeCurl)

	bytes, _ := wasm.ReadBytes("imported_function.wasm")
	instance, _ := wasm.NewInstanceWithImports(bytes, imports)
	defer instance.Close()

	// Gets and calls the `add1` exported function from the WebAssembly instance.
	results, _ := instance.Exports["add1"](1, 2)

	// TODO: Read the Event and Context
	// Store the Envent and Context back in IPFS

	fmt.Println(results)
}
