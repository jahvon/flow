package examples_test

import "github.com/jahvon/flow/config"

var (
	requestBaseDesc = "Request flows sends HTTP requests with the specified request and response settings."
)

var SimpleRequestExec = &config.Executable{
	Verb:       "run",
	Name:       "simple-request",
	Visibility: config.VisibilityPrivate,
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

var RequestWithHeadersExec = &config.Executable{
	Verb:       "run",
	Name:       "request-with-headers",
	Visibility: config.VisibilityPrivate,
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

var RequestWithBodyExec = &config.Executable{
	Verb:       "run",
	Name:       "request-with-body",
	Visibility: config.VisibilityPrivate,
	Description: requestBaseDesc +
		"\n\nThe `body` field is optional and can be used to send a request body.",
	Type: &config.ExecutableTypeSpec{
		Request: &config.RequestExecutableType{
			URL:  "https://httpbin.org/post",
			Body: `{"key": "value"}`,
		},
	},
}

var RequestWithTransformExec = &config.Executable{
	Verb:       "run",
	Name:       "request-with-transform",
	Visibility: config.VisibilityPrivate,
	Description: requestBaseDesc +
		"\n\nThe `transformResponse` field is optional and can be used to transform the response using a jq query.",
	Type: &config.ExecutableTypeSpec{
		Request: &config.RequestExecutableType{
			URL:               "https://httpbin.org/get",
			TransformResponse: ".headers",
		},
	},
}

var RequestWithTimeoutExec = &config.Executable{
	Verb:       "run",
	Name:       "request-with-timeout",
	Visibility: config.VisibilityPrivate,
	Description: requestBaseDesc +
		"\n\nThe `timeout` field is optional and can be used to set the request timeout.",
	Type: &config.ExecutableTypeSpec{
		Request: &config.RequestExecutableType{
			URL:     "https://httpbin.org/delay/3",
			Timeout: 1,
		},
	},
}

var RequestWithValidStatusCodesExec = &config.Executable{
	Verb:       "run",
	Name:       "request-with-valid-status-codes",
	Visibility: config.VisibilityPrivate,
	Description: requestBaseDesc +
		"\n\nThe `validStatusCodes` field is optional and can be used to specify the valid status codes.",
	Type: &config.ExecutableTypeSpec{
		Request: &config.RequestExecutableType{
			URL:              "https://httpbin.org/status/200",
			ValidStatusCodes: []int{200},
		},
	},
}

var RequestWithInvalidStatusCodeExec = &config.Executable{
	Verb:       "run",
	Name:       "request-with-invalid-status-code",
	Visibility: config.VisibilityPrivate,
	Description: requestBaseDesc +
		"\n\nThe `validStatusCodes` field is optional and can be used to specify the valid status codes.",
	Type: &config.ExecutableTypeSpec{
		Request: &config.RequestExecutableType{
			URL:              "https://httpbin.org/status/400",
			ValidStatusCodes: []int{200},
		},
	},
}
