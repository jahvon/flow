package utils

import (
	execUtils "github.com/jahvon/flow/internal/utils/executables"
	"github.com/jahvon/flow/types/executable"
)

// FindSubExecs returns a list of sub-executables for the given root executable.
func FindSubExecs(rootExec *executable.Executable, flowFiles executable.FlowFileList) executable.ExecutableList {
	switch {
	case rootExec.Serial != nil:
		return findSerialSubExecs(rootExec, flowFiles)
	case rootExec.Parallel != nil:
		return findParallelSubExecs(rootExec, flowFiles)
	default:
		return nil
	}
}

func findSerialSubExecs(root *executable.Executable, flowFiles executable.FlowFileList) []*executable.Executable {
	serial := root.Serial
	var subExecs []*executable.Executable
	for _, ref := range serial.Refs {
		for _, flowFile := range flowFiles {
			if e, _ := flowFile.Executables.FindByVerbAndID(ref.GetVerb(), ref.GetID()); e != nil {
				subExecs = append(subExecs, e)
				break
			}
		}
	}
	for i, refCfg := range serial.Execs {
		if refCfg.Cmd != "" {
			subExecs = append(subExecs, execUtils.ExecutableForCmd(root, refCfg.Cmd, i))
		}

		for _, flowFile := range flowFiles {
			if e, _ := flowFile.Executables.FindByVerbAndID(refCfg.Ref.GetVerb(), refCfg.Ref.GetID()); e != nil {
				subExecs = append(subExecs, e)
				break
			}
		}
	}
	return subExecs
}

func findParallelSubExecs(root *executable.Executable, flowFiles executable.FlowFileList) []*executable.Executable {
	parallel := root.Parallel
	var subExecs []*executable.Executable
	for _, ref := range parallel.Refs {
		for _, flowFile := range flowFiles {
			if e, _ := flowFile.Executables.FindByVerbAndID(ref.GetVerb(), ref.GetID()); e != nil {
				subExecs = append(subExecs, e)
				break
			}
		}
	}
	for i, refCfg := range parallel.Execs {
		if refCfg.Cmd != "" {
			subExecs = append(subExecs, execUtils.ExecutableForCmd(root, refCfg.Cmd, i))
		}

		for _, flowFile := range flowFiles {
			if e, _ := flowFile.Executables.FindByVerbAndID(refCfg.Ref.GetVerb(), refCfg.Ref.GetID()); e != nil {
				subExecs = append(subExecs, e)
				break
			}
		}
	}
	return subExecs
}
