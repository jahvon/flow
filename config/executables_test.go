package config_test

import (
	"fmt"
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/config"
)

var _ = Describe("Executable", func() {
	const (
		testWsName        = "workspace"
		testNsName        = "namespace"
		testWorkspacePath = "workspace-path"
		testExecDefPath   = "exec-definition-path"
	)
	var exec *config.Executable

	BeforeEach(func() {
		execType := &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				Command: "echo hello",
			},
		}
		exec = &config.Executable{
			Verb:        "run",
			Name:        "test",
			Aliases:     []string{"t"},
			Description: "test description",
			Type:        execType,
		}
		exec.SetDefaults()
		exec.SetContext(testWsName, testWorkspacePath, testNsName, testExecDefPath)
		Expect(exec.Validate()).To(Succeed())
	})

	Describe("Format Methods", func() {
		It("JSON should return the JSON representation of the executable", func() {
			str, err := exec.JSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(str).ToNot(BeEmpty())
		})
		It("YAML should return the YAML representation of the executable", func() {
			str, err := exec.YAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(str).ToNot(BeEmpty())
		})
		It("Markdown should return the Markdown representation of the executable", func() {
			str := exec.Markdown()
			Expect(str).ToNot(BeEmpty())
		})
	})

	Describe("Metadata Methods", func() {
		It("Ref should return the reference of the executable", func() {
			Expect(exec.Ref().String()).
				To(Equal(fmt.Sprintf("%s %s/%s:%s", exec.Verb, testWsName, testNsName, exec.Name)))
		})
		It("ID should return the ID of the executable", func() {
			Expect(exec.ID()).
				To(Equal(fmt.Sprintf("%s/%s:%s", testWsName, testNsName, exec.Name)))
		})
		It("WorkspacePath should return the workspace path of the executable", func() {
			Expect(exec.WorkspacePath()).To(Equal(testWorkspacePath))
		})
		It("DefinitionPath should return the exec definition path of the executable", func() {
			Expect(exec.DefinitionPath()).To(Equal(testExecDefPath))
		})
	})

	Describe("AliasesIDs", func() {
		It("should return the correct aliases IDs", func() {
			exec.Aliases = []string{"alias1", "alias2"}
			aliasesIDs := exec.AliasesIDs()
			Expect(aliasesIDs).To(ConsistOf(
				fmt.Sprintf("%s/%s:alias1", testWsName, testNsName),
				fmt.Sprintf("%s/%s:alias2", testWsName, testNsName),
			))
		})

		It("should return nil if there are no aliases", func() {
			exec.Aliases = nil
			aliasesIDs := exec.AliasesIDs()
			Expect(aliasesIDs).To(BeNil())
		})
	})

	Describe("NameEquals", func() {
		It("should return the expected value", func() {
			By("having a matching name")
			Expect(exec.NameEquals(exec.Name)).To(BeTrue())

			By("having a matching alias")
			Expect(exec.NameEquals(exec.Aliases[0])).To(BeTrue())

			By("not having a matching name or alias")
			Expect(exec.NameEquals("nonexistent")).To(BeFalse())
		})
	})

	Describe("MergeTags", func() {
		It("should merge the given tags with the executable's existing tags", func() {
			exec.MergeTags(config.Tags{"tag1", "tag2"})
			Expect(exec.Tags).To(ConsistOf("tag1", "tag2"))
		})

		It("should remove duplicate tags", func() {
			exec.MergeTags(config.Tags{"tag1", "tag1"})
			compact := slices.Compact(exec.Tags)
			Expect(compact).To(HaveLen(len(exec.Tags)))
		})
	})

	DescribeTable("IsVisibleFromWorkspace", func(visibility *config.Visibility, wsMatch, expected bool) {
		exec.Visibility = visibility
		if wsMatch {
			Expect(exec.IsVisibleFromWorkspace(testWsName)).To(Equal(expected))
		} else {
			Expect(exec.IsVisibleFromWorkspace("another-ws")).To(Equal(expected))
		}
	},
		Entry("public from ws", vPtr(config.VisibilityPublic), true, true),
		Entry("public from another ws", vPtr(config.VisibilityPublic), false, true),
		Entry("private from ws", vPtr(config.VisibilityPrivate), true, true),
		Entry("private from another ws", vPtr(config.VisibilityPrivate), false, false),
		Entry("internal from ws", vPtr(config.VisibilityInternal), true, false),
		Entry("internal from another ws", vPtr(config.VisibilityInternal), false, false),
		Entry("hidden from ws", vPtr(config.VisibilityHidden), true, false),
		Entry("hidden from another ws", vPtr(config.VisibilityHidden), false, false),
	)

	DescribeTable("IsExecutableFromWorkspace", func(visibility *config.Visibility, wsMatch, expected bool) {
		exec.Visibility = visibility
		if wsMatch {
			Expect(exec.IsExecutableFromWorkspace(testWsName)).To(Equal(expected))
		} else {
			Expect(exec.IsExecutableFromWorkspace("another-ws")).To(Equal(expected))
		}
	},
		Entry("public from ws", vPtr(config.VisibilityPublic), true, true),
		Entry("public from another ws", vPtr(config.VisibilityPublic), false, true),
		Entry("private from ws", vPtr(config.VisibilityPrivate), true, true),
		Entry("private from another ws", vPtr(config.VisibilityPrivate), false, false),
		Entry("internal from ws", vPtr(config.VisibilityInternal), true, true),
		Entry("internal from another ws", vPtr(config.VisibilityInternal), false, false),
		Entry("hidden from ws", vPtr(config.VisibilityHidden), true, false),
		Entry("hidden from another ws", vPtr(config.VisibilityHidden), false, false),
	)
})

var _ = Describe("ExecutableList", func() {
	const (
		exec1Ws = "ws1"
		exec2Ws = "ws2"
		exec1Ns = "ns1"
		exec2Ns = "ns2"
	)
	var (
		exec1 *config.Executable
		exec2 *config.Executable
		execs config.ExecutableList
	)

	BeforeEach(func() {
		exec1 = &config.Executable{
			Verb: "run",
			Name: "test1",
			Type: &config.ExecutableTypeSpec{
				Exec: &config.ExecExecutableType{
					Command: "echo hello",
				},
			},
		}
		exec1.SetDefaults()
		exec1.SetContext(exec1Ws, "workspace-path", exec1Ns, "exec-definition-path")
		exec2 = &config.Executable{
			Verb: "start",
			Name: "test2",
			Type: &config.ExecutableTypeSpec{
				Exec: &config.ExecExecutableType{
					Command: "echo hello",
				},
			},
		}
		exec2.SetDefaults()
		exec2.SetContext(exec2Ws, "workspace-path", exec2Ns, "exec-definition-path")
		execs = config.ExecutableList{exec1, exec2}
	})

	Describe("Format Methods", func() {
		It("JSON should return the JSON representation of the executables", func() {
			str, err := execs.JSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(str).ToNot(BeEmpty())
		})
		It("YAML should return the YAML representation of the executables", func() {
			str, err := execs.YAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(str).ToNot(BeEmpty())
		})
		It("Items should return the Markdown representation of the executables", func() {
			items := execs.Items()
			// TODO: test the markdown content
			Expect(items).To(HaveLen(2))
		})
	})

	Describe("FilterByNamespace", func() {
		It("should return only executables with the given namespace", func() {
			filtered := execs.FilterByNamespace(exec1Ns)
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
		})
	})

	Describe("FilterByTag", func() {
		BeforeEach(func() {
			exec1.Tags = config.Tags{"tag1", "tag2"}
			exec2.Tags = config.Tags{"tag3", "tag4"}
			execs = config.ExecutableList{exec1, exec2}
		})

		It("should return executables with the given tag", func() {
			filtered := execs.FilterByTags(config.Tags{"tag2"})
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
		})

		It("should return no executables if the tag is not found", func() {
			filtered := execs.FilterByTags(config.Tags{"tag5"})
			Expect(filtered).To(BeEmpty())
		})

		It("should return executables with multiple tags", func() {
			filtered := execs.FilterByTags(config.Tags{"tag2", "tag3"})
			Expect(filtered).To(HaveLen(2))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
			Expect(filtered[1].Name).To(Equal(exec2.Name))
		})
	})

	Describe("FilterByVerb", func() {
		BeforeEach(func() {
			exec1.Verb = "run"
			exec2.Verb = "launch"
			execs = config.ExecutableList{exec1, exec2}
		})

		It("should return executables with the given verb", func() {
			filtered := execs.FilterByVerb("run")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
		})

		It("should return no executables if the verb is not found", func() {
			filtered := execs.FilterByVerb("uninstall")
			Expect(filtered).To(BeEmpty())
		})
	})

	Describe("FilterByWorkspace", func() {
		It("should return executables with the given workspace", func() {
			filtered := execs.FilterByWorkspace(exec1Ws)
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
		})
	})

	Describe("FilterBySubstring", func() {
		It("should return executables when the reference matches", func() {
			exec1.Name = "abcdefgh"
			filtered := execs.FilterBySubstring("def")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
		})

		It("should return executables when the description matches", func() {
			exec1.Description = "abcdefgh"
			filtered := execs.FilterBySubstring("def")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal(exec1.Name))
		})
	})

	Describe("FindByVerbAndID", func() {
		It("should return the executable with the given verb and id", func() {
			exec, err := execs.FindByVerbAndID(exec1.Verb, exec1.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(exec).To(Equal(exec1))
		})

		It("should return nil if there is no match", func() {
			exec, err := execs.FindByVerbAndID("nonexistent", "nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(exec).To(BeNil())
		})
	})
})

func vPtr(v config.Visibility) *config.Visibility { return &v }
