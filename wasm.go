package main

// import (
// 	"fmt"

// 	"github.com/perlin-network/life/exec"
// )

// type Resolver struct{}

// func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
// 	switch module {
// 	case "env":
// 		switch field {
// 		case "ipfsComputeLog":
// 			return func(vm *exec.VirtualMachine) int64 {
// 				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
// 				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
// 				msg := vm.Memory[ptr : ptr+msgLen]
// 				fmt.Printf("[app] %s\n", string(msg))
// 				return 0
// 			}
// 		case "ipfsComputeGet":
// 			return func(vm *exec.VirtualMachine) int64 {
// 				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
// 				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
// 				msg := vm.Memory[ptr : ptr+msgLen]
// 				fmt.Printf("[app] %s\n", string(msg))
// 				return 0
// 			}
// 		case "ipfsComputeLs":
// 			return func(vm *exec.VirtualMachine) int64 {
// 				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
// 				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
// 				msg := vm.Memory[ptr : ptr+msgLen]
// 				fmt.Printf("[app] %s\n", string(msg))
// 				return 0
// 			}

// 		case "ipfsComputeAdd":
// 			return func(vm *exec.VirtualMachine) int64 {
// 				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
// 				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
// 				msg := vm.Memory[ptr : ptr+msgLen]
// 				fmt.Printf("[app] %s\n", string(msg))
// 				return 0
// 			}

// 		case "ipfsComputeCurl":
// 			return func(vm *exec.VirtualMachine) int64 {
// 				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
// 				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
// 				msg := vm.Memory[ptr : ptr+msgLen]
// 				fmt.Printf("[app] %s\n", string(msg))
// 				return 0
// 			}
// 		default:
// 			panic(fmt.Errorf("unknown import resolved: %s", field))
// 		}
// 	default:
// 		panic(fmt.Errorf("unknown module: %s", module))
// 	}
// }

// func (r *Resolver) ResolveGlobal(module, field string) int64 {
// 	panic("we're not resolving global variables for now")
// }
