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

func readFileContent(f *os.File) (string, error) {
	out, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	outStr := string(out)
	if truthy, _ := strconv.ParseBool(os.Getenv("PRINT_TEST_STDOUT")); truthy {
		fmt.Println(outStr)
	}
	return outStr, nil
}

func writeUserInput(f *os.File, input string) error {
	if _, err := f.WriteString(input); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
