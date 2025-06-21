package main

import (
	"github.com/jahvon/flow/types/common"
	"github.com/jahvon/flow/types/executable"
)

func sharedExecTags() executable.ExecutableTags {
	return executable.ExecutableTags{"generated"}
}

func privateFlowFileVisibility() *executable.FlowFileVisibility {
	v := executable.FlowFileVisibility(common.VisibilityPrivate)
	return &v
}

func privateExecVisibility() *executable.ExecutableVisibility {
	v := executable.ExecutableVisibility(common.VisibilityPrivate)
	return &v
}
