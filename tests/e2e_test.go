package tests_test

import (
	"fmt"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "End-to-end Test Suite")
}

const PrintEnvVar = "PRINT_TEST_STDOUT"

func readFileContent(f *os.File) (string, error) {
	out, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	outStr := string(out)
	if truthy, _ := strconv.ParseBool(os.Getenv(PrintEnvVar)); truthy {
		fmt.Println(outStr)
	}
	return outStr, nil
}
