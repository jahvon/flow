package config

type OutputFormat string

const (
	JSON        OutputFormat = "json"
	JSONP       OutputFormat = "jsonp"
	YAML        OutputFormat = "yaml"
	INTERACTIVE OutputFormat = "interactive"
)
