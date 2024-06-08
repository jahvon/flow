package args_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/config/args"
)

func TestParseArgs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Args Suite")
}

var _ = Describe("ParseArgs", func() {
	It("should correctly parse flag arguments", func() {
		flagArgs, posArgs := args.ParseArgs([]string{"flag1=value1", "flag2=value2"})

		Expect(flagArgs).To(Equal(map[string]string{"flag1": "value1", "flag2": "value2"}))
		Expect(posArgs).To(Equal([]string{}))
	})

	It("should correctly parse positional arguments", func() {
		flagArgs, posArgs := args.ParseArgs([]string{"pos1", "pos2"})

		Expect(flagArgs).To(Equal(map[string]string{}))
		Expect(posArgs).To(Equal([]string{"pos1", "pos2"}))
	})

	It("should correctly parse mixed arguments", func() {
		flagArgs, posArgs := args.ParseArgs([]string{"flag1=value1", "pos1"})

		Expect(flagArgs).To(Equal(map[string]string{"flag1": "value1"}))
		Expect(posArgs).To(Equal([]string{"pos1"}))
	})

	It("should correctly parse flag arguments with equal sign in value", func() {
		flagArgs, posArgs := args.ParseArgs([]string{"flag1=value1=value2"})

		Expect(flagArgs).To(Equal(map[string]string{"flag1": "value1=value2"}))
		Expect(posArgs).To(Equal([]string{}))
	})
})
