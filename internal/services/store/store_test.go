package store_test

import (
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/services/store"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BoltStore Suite")
}

var _ = Describe("BoltStore", func() {
	var s store.Store
	var err error

	BeforeEach(func() {
		path := filepath.Join(GinkgoT().TempDir(), fmt.Sprintf("test_%s.db", GinkgoT().Name()))
		s, err = store.NewStore(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())
	})

	AfterEach(func() {
		err = s.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("CreateBucket", func() {
		It("should create a new bucket", func() {
			err := s.CreateBucket(store.RootBucket)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("DeleteBucket", func() {
		It("should delete an existing bucket", func() {
			err := s.CreateBucket("test")
			Expect(err).NotTo(HaveOccurred())

			err = s.DeleteBucket("test")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Set and Get", func() {
		It("should set and get a key-value pair", func() {
			err = s.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, err := s.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("value"))
		})

		It("should set and get key-value pairs across multiple buckets", func() {
			By("setting a key-value pair in the root bucket")
			err = s.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key-inherited", "value2")
			Expect(err).NotTo(HaveOccurred())

			By("setting a key-value pair in a process bucket")
			id, err := s.CreateAndSetBucket("process2")
			Expect(id).To(Equal("process2"))
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key", "value3")
			Expect(err).NotTo(HaveOccurred())

			By("getting key-value pairs from a process bucket")
			_, err = s.CreateAndSetBucket(store.EnvironmentBucket())
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
			id, err = s.CreateAndSetBucket(store.RootBucket)
			Expect(id).To(Equal(store.RootBucket))
			Expect(err).NotTo(HaveOccurred())

			value, err = s.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("value"))
		})
	})

	Describe("GetAll", func() {
		It("should get all key-value pairs from the bucket", func() {
			err = s.CreateBucket("test")
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key1", "value1")
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key2", "value2")
			Expect(err).NotTo(HaveOccurred())

			pairs, err := s.GetAll()
			Expect(err).NotTo(HaveOccurred())
			Expect(pairs).To(HaveLen(2))
			Expect(pairs["key1"]).To(Equal("value1"))
			Expect(pairs["key2"]).To(Equal("value2"))
		})
	})

	Describe("GetKeys", func() {
		It("should get all keys from the bucket", func() {
			err = s.CreateBucket("test")
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key1", "value1")
			Expect(err).NotTo(HaveOccurred())

			err = s.Set("key2", "value2")
			Expect(err).NotTo(HaveOccurred())

			keys, err := s.GetKeys()
			Expect(err).NotTo(HaveOccurred())
			Expect(keys).To(HaveLen(2))
			Expect(keys).To(ContainElement("key1"))
			Expect(keys).To(ContainElement("key2"))
		})
	})

	Describe("Delete", func() {
		It("should delete a key from the bucket", func() {
			err := s.CreateBucket("test")
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
