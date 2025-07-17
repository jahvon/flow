package fileparser_test

import (
	"os"
	"path/filepath"

	"github.com/flowexec/flow/internal/fileparser"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecConfigMapFromFile", func() {
	DescribeTable("should correctly parse configurations",
		func(file string, expected map[string]string) {
			filePath := filepath.Join("testdata", file)
			fileBytes, err := os.ReadFile(filepath.Clean(filePath))
			Expect(err).ToNot(HaveOccurred())

			result, err := fileparser.ExtractExecConfigMap(string(fileBytes), "# ")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		},
		Entry(
			"simple key-value pairs",
			"simple.sh",
			map[string]string{
				fileparser.NameConfigurationKey: "hello",
				fileparser.VerbConfigurationKey: "show",
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
