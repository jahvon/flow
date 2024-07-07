package builder

import (
	"github.com/jahvon/flow/types/executable"
)

const (
	requestBaseDesc = "Request executables send HTTP requests with the specified request and response settings.\n"
)

func RequestExec(opts ...Option) *executable.Executable {
	name := "request"
	docstring := requestBaseDesc +
		"The `url` field is required and must be a valid URL. " +
		"The `method` field is optional and defaults to `GET`.\n" +
		"The `headers` field is optional and can be used to set request headers."
	e := &executable.Executable{
		Verb:        "run",
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
	docstring := requestBaseDesc +
		"The `body` field is optional and can be used to send a request body."
	e := &executable.Executable{
		Verb:        "run",
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
	docstring := requestBaseDesc +
		"The `transformResponse` field is optional and can be used to transform the response using a jq query."
	e := &executable.Executable{
		Verb:        "run",
		Name:        name,
		Visibility:  privateExecVisibility(),
		Description: docstring,
		Request: &executable.RequestExecutableType{
			URL:               "https://httpbin.org/get",
			TransformResponse: ".headers",
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
	docstring := requestBaseDesc +
		"The `timeout` field is optional and can be used to set the request timeout."
	e := &executable.Executable{
		Verb:        "run",
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
	docstring := requestBaseDesc +
		"The `validStatusCodes` field is optional and can be used to specify the valid status codes. " +
		"If the response status code is not in the list, the executable will fail."
	e := &executable.Executable{
		Verb:        "run",
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
