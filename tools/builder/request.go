package main

import (
	"github.com/jahvon/flow/types/executable"
)

const (
	requestBaseDesc = "Request executables send HTTP requests with the specified request and response settings.\n"
)

func RequestExec(opts ...Option) *executable.Executable {
	name := "request"
	docstring := "The `url` field is required and must be a valid URL. " +
		"The `method` field is optional and defaults to `GET`.\n" +
		"The `headers` field is optional and can be used to set request headers."
	e := &executable.Executable{
		Verb:        "send",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Request: &executable.RequestExecutableType{
			URL:    "https://httpbin.org/get",
			Method: "GET",
			Headers: map[string]string{
				"Authorization": "Bearer token",
				"User-Agent":    "flow",
			},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func RequestExecWithBody(opts ...Option) *executable.Executable {
	name := "request-with-body"
	docstring := "The `body` field is optional and can be used to send a request body."
	e := &executable.Executable{
		Verb:        "send",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Request: &executable.RequestExecutableType{
			URL:    "https://httpbin.org/post",
			Method: "POST",
			Body:   `{"key": "value"}`,
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func RequestExecWithTransform(opts ...Option) *executable.Executable {
	name := "request-with-transform"
	docstring := "The `transformResponse` field is optional and can be used to transform the response using an Expr expression."
	e := &executable.Executable{
		Verb:        "send",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Request: &executable.RequestExecutableType{
			URL:               "https://httpbin.org/get",
			TransformResponse: "status",
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func RequestExecWithTimeout(opts ...Option) *executable.Executable {
	name := "request-with-timeout"
	docstring := "The `timeout` field is optional and can be used to set the request timeout."
	e := &executable.Executable{
		Verb:        "send",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Request: &executable.RequestExecutableType{
			URL:     "https://httpbin.org/delay/3",
			Timeout: 1,
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}

func RequestExecWithValidatedStatus(opts ...Option) *executable.Executable {
	name := "request-with-validated-status"
	docstring := "The `validStatusCodes` field is optional and can be used to specify the valid status codes. " +
		"If the response status code is not in the list, the executable will fail."
	e := &executable.Executable{
		Verb:        "send",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Request: &executable.RequestExecutableType{
			URL:              "https://httpbin.org/status/400",
			ValidStatusCodes: []int{200},
		},
	}
	if len(opts) > 0 {
		vals := NewOptionValues(opts...)
		e.SetContext(vals.WorkspaceName, vals.WorkspacePath, vals.NamespaceName, vals.FlowFilePath)
	}
	return e
}
