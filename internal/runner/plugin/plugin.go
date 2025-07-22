package plugin

import (
	stdCtx "context"
	_ "embed"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/runner/engine"
	"github.com/flowexec/flow/types/executable"
)

//go:embed testdata/plugin.wasm
var plugin []byte

type PluginRunner struct {
}

func NewRunner() *PluginRunner {
	return &PluginRunner{}
}
func (p *PluginRunner) Name() string {
	return "plugin"
}

func (p *PluginRunner) IsCompatible(executable *executable.Executable) bool {
	return executable != nil && executable.Plugin != nil
}

func (p *PluginRunner) Exec(ctx *context.Context, e *executable.Executable, eng engine.Engine, inputEnv map[string]string) error {
	r := wazero.NewRuntime(ctx.Ctx)
	defer r.Close(ctx.Ctx)

	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logString).Export("log").
		Instantiate(ctx.Ctx)
	if err != nil {
		log.Panicln(err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctx.Ctx, r)
	mod, err := r.InstantiateWithConfig(ctx.Ctx, plugin, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	exec := mod.ExportedFunction("exec")
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")
	str := *e.Plugin.Data
	strSize := uint64(len(str))
	results, err := malloc.Call(ctx.Ctx, uint64(strSize))
	if err != nil {
		log.Panicln(err)
	}
	strPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx.Ctx, strPtr)

	if !mod.Memory().Write(uint32(strPtr), []byte(str)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			strPtr, strSize, mod.Memory().Size())
	}

	results, err = exec.Call(ctx.Ctx, strPtr, strSize)
	if err != nil {
		log.Panicf("failed to call exec: %v", err)
	}
	return nil
}

func logString(_ stdCtx.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}
