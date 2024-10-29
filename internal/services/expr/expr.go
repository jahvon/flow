package expr

import (
	"fmt"
	"strconv"

	"github.com/expr-lang/expr"
)

func IsTruthy(ex string, env map[string]interface{}) (bool, error) {
	program, err := expr.Compile(ex, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	switch v := output.(type) {
	case bool:
		return v, nil
	case int, int64, float64, uint, uint64:
		return v != 0, nil
	case string:
		truthy, err := strconv.ParseBool(v)
		if err != nil {
			return false, nil
		}
		return truthy || v != "", nil
	default:
		return false, nil
	}
}

func Evaluate(ex string, env map[string]interface{}) (interface{}, error) {
	program, err := expr.Compile(ex, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	return output, nil
}

func EvaluateString(ex string, env map[string]interface{}) (string, error) {
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
