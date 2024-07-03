package builder

import (
	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

var (
	requestBaseDesc = "Request flows sends HTTP requests with the specified request and response settings."
)

func RequestExec(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `url` field is required and must be a valid URL." +
			"\n\nThe `method` field is optional and defaults to `GET`.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL:    "https://httpbin.org/get",
				Method: "GET",
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func RequestExecWithHeaders(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `headers` field is optional and can be used to set custom headers.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL: "https://httpbin.org/get",
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"User-Agent":    "flow",
				},
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func RequestExecWithBody(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `body` field is optional and can be used to send a request body.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL:  "https://httpbin.org/post",
				Body: `{"key": "value"}`,
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func RequestExecWithTransform(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `transformResponse` field is optional and can be used to transform the response using a jq query.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL:               "https://httpbin.org/get",
				TransformResponse: ".headers",
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func RequestExecWithTimeout(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `timeout` field is optional and can be used to set the request timeout.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL:     "https://httpbin.org/delay/3",
				Timeout: 1,
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func RequestExecWithValidStatusCodes(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `validStatusCodes` field is optional and can be used to specify the valid status codes.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL:              "https://httpbin.org/status/200",
				ValidStatusCodes: []int{200},
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}

func RequestExecWithInvalidStatusCode(ctx *context.Context, name, definitionPath string) *config.Executable {
	e := &config.Executable{
		Verb:       "run",
		Name:       name,
		Visibility: config.VisibilityInternal.NewPointer(),
		Description: requestBaseDesc +
			"\n\nThe `validStatusCodes` field is optional and can be used to specify the valid status codes.",
		Type: &config.ExecutableTypeSpec{
			Request: &config.RequestExecutableType{
				URL:              "https://httpbin.org/status/400",
				ValidStatusCodes: []int{200},
			},
		},
	}
	e.SetContext(
		ctx.CurrentWorkspace.AssignedName(), ctx.CurrentWorkspace.Location(),
		ctx.UserConfig.CurrentNamespace, definitionPath,
	)
	return e
}
