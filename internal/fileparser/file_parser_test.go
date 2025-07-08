package fileparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flowexec/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/flowexec/flow/internal/fileparser"
)

func TestFileParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileParser Suite")
}

var _ = Describe("ExecConfigMapFromFile", func() {
	var (
		ctrl       *gomock.Controller
		mockLogger *mocks.MockLogger
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockLogger = mocks.NewMockLogger(ctrl)
	})

	DescribeTable("should error when the file is invalid", func(file string) {
		mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
		wd, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		filePath := filepath.Join(wd, "testdata", file)
		result, err := fileparser.ExecConfigMapFromFile(mockLogger, filePath)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeNil())
	},
		Entry("non-shell file", "invalidfile"),
		Entry("dir instead of file", "invaliddir.sh"),
		Entry("file without configs", "empty.sh"),
		Entry("non-existent file", "nonexistent.sh"),
	)

	It("should log a warning when configuration key is not recognized", func() {
		mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)
		wd, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		filePath := filepath.Join(wd, "testdata", "unknownkey.sh")
		result, err := fileparser.ExecConfigMapFromFile(mockLogger, filePath)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeNil())
	})

	DescribeTable("should correctly parse configurations",
		func(file string, expected map[string]string) {
			mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
			wd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())
			filePath := filepath.Join(wd, "testdata", file)
			result, err := fileparser.ExecConfigMapFromFile(mockLogger, filePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		},
		Entry(
			"simple key-value pairs",
			"simple.sh",
			map[string]string{
				fileparser.NameConfigurationKey: "value1",
				fileparser.VerbConfigurationKey: "value2",
			}),
		Entry(
			"values with spaces in quotes",
			"quoted.sh",
			map[string]string{
				fileparser.NameConfigurationKey:        "value 1",
				fileparser.VerbConfigurationKey:        "value2",
				fileparser.DescriptionConfigurationKey: "value 3",
			}),
		Entry(
			"mixed separators",
			"mixed.sh",
			map[string]string{
				fileparser.NameConfigurationKey:        "value1",
				fileparser.VerbConfigurationKey:        "value2",
				fileparser.DescriptionConfigurationKey: "value 3",
			}),
		Entry(
			"values with escaped characters",
			"escaped.sh",
			map[string]string{
				fileparser.NameConfigurationKey:        "value 1, one",
				fileparser.DescriptionConfigurationKey: "'value two'",
				fileparser.TagConfigurationKey:         "tag1,tag2",
			}),
		Entry(
			"repeated key configurations",
			"repeated.sh",
			map[string]string{
				fileparser.TagConfigurationKey:         "tag1,tag2,tag3,tag4,tag5",
				fileparser.AliasConfigurationKey:       "alias",
				fileparser.DescriptionConfigurationKey: "first line\nsecond line\nthird line",
			}),
		Entry(
			"complex configuration parsing",
			"complex.sh",
			map[string]string{
				fileparser.NameConfigurationKey:        "name",
				fileparser.VerbConfigurationKey:        "verb",
				fileparser.DescriptionConfigurationKey: "first line\nsecond line\nthird line with 'quotes', and commas\nclosin'",
			}),
		Entry(
			"multi-line description",
			"multiline.sh",
			map[string]string{
				fileparser.DescriptionConfigurationKey: expectedMultiLineDescription,
			}),
	)
})

const expectedMultiLineDescription = `first line
second line
    third line, with commas and 'quotes'
fourth line
fifth line
sixth line
seventh line
eighth line`
