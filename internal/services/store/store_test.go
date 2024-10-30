package store_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/services/store"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = Describe("Store", func() {
	var s *store.Store
	var err error

	BeforeEach(func() {
		s, err = store.NewStore()
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())
	})

	AfterEach(func() {
		err = s.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("CreateBucket", func() {
		It("should create a new bucket", func() {
			err := s.CreateBucket()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("DeleteBucket", func() {
		It("should delete an existing bucket", func() {
			err := s.CreateBucket()
			Expect(err).NotTo(HaveOccurred())

			err = s.DeleteBucket()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Set and Get", func() {
		It("should set and get a key-value pair", func() {
			err := s.CreateBucket()
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, err := s.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("value"))
		})

		It("should set and get key-value pairs across multiple buckets", func() {
			By("setting a key-value pair in the root bucket")
			err := s.CreateBucket()
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key-inherited", "value2")
			Expect(err).NotTo(HaveOccurred())

			By("setting a key-value pair in a process bucket")
			err = store.SetProcessBucketID("process2", true)
			Expect(err).NotTo(HaveOccurred())

			err = s.CreateBucket()
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key", "value3")
			Expect(err).NotTo(HaveOccurred())

			By("getting key-value pairs from a process bucket")
			err = store.SetProcessBucketID("process2", true)
			Expect(err).NotTo(HaveOccurred())

			value, err := s.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("value3"))

			_, err = s.Get("key-unset")
			Expect(err).To(HaveOccurred())

			value, err = s.Get("key-inherited")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("value2"))

			By("getting key-value pairs from the root bucket")
			err = store.SetProcessBucketID(store.RootBucket, true)
			Expect(err).NotTo(HaveOccurred())

			value, err = s.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("value"))
		})
	})

	Describe("Delete", func() {
		It("should delete a key from the bucket", func() {
			err := s.CreateBucket()
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = s.Delete("key")
			Expect(err).NotTo(HaveOccurred())

			_, err = s.Get("key")
			Expect(err).To(HaveOccurred())
		})
	})
})
