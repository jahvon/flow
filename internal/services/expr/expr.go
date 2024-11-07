package expr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/types/executable"
)

func IsTruthy(ex string, env *ExpressionData) (bool, error) {
	program, err := expr.Compile(ex, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	switch v := output.(type) {
	case bool:
		return v, nil
	case int, int64, float64, uint, uint64:
		return v != 0, nil
	case string:
		truthy, err := strconv.ParseBool(strings.Trim(v, `"' `))
		if err != nil {
			return false, err
		}
		return truthy || v != "", nil
	default:
		return false, nil
	}
}

func Evaluate(ex string, env *ExpressionData) (interface{}, error) {
	program, err := expr.Compile(ex, expr.Env(env))
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func EvaluateString(ex string, env *ExpressionData) (string, error) {
	output, err := Evaluate(ex, env)
	if err != nil {
		return "", err
	}
	str, ok := output.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", output)
	}
	return str, nil
}

type CtxData struct {
	Workspace     string
	Namespace     string
	WorkspacePath string
	FlowFilePath  string
	FlowFileDir   string
}

type ExpressionData struct {
	OS   string
	Arch string
	Ctx  *CtxData
	Data map[string]string
	Env  map[string]string
}

func ExpressionEnv(
	ctx *context.Context,
	executable *executable.Executable,
	dataMap, envMap map[string]string,
) ExpressionData {
	return ExpressionData{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		Ctx: &CtxData{
			Workspace:     ctx.CurrentWorkspace.AssignedName(),
			Namespace:     ctx.CurrentWorkspace.AssignedName(),
			WorkspacePath: executable.WorkspacePath(),
			FlowFilePath:  executable.FlowFilePath(),
			FlowFileDir:   filepath.Dir(executable.FlowFilePath()),
		},
		Data: dataMap,
		Env:  envMap,
	}
}
