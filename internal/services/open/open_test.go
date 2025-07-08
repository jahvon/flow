package open_test

import (
	"testing"

	og "github.com/jahvon/open-golang/open"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/flowexec/flow/internal/services/open"
)

func TestOpen(t *testing.T) {
	RegisterFailHandler(Fail)
	t.Setenv(og.OpenDisabledEnvKey, "true")
	RunSpecs(t, "Open Suite")
}

var _ = Describe("Open", func() {
	Context("when wait is true", func() {
		It("should not return an error", func() {
			err := open.Open("http://example.com", true)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when wait is false", func() {
		It("should not return an error", func() {
			err := open.Open("http://example.com", false)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("OpenWith", func() {
	Context("when wait is true", func() {
		It("should not return an error", func() {
			err := open.OpenWith("firefox", "http://example.com", true)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when wait is false", func() {
		It("should not return an error", func() {
			err := open.OpenWith("firefox", "http://example.com", false)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
