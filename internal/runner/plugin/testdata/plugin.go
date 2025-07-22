package main

import (
	"runtime"
	"unsafe"
)

//go:wasmexport exec
func exec(ptr, size uint32) {
	str := ptrToString(ptr, size)
	_log(ptr, size)
	runtime.KeepAlive(str)
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}

// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

//go:wasmimport env log
func _log(ptr, size uint32)
