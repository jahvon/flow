package config

type OutputFormat string

const (
	JSON          OutputFormat = "json"
	FormattedJSON OutputFormat = "jsonp"
	YAML          OutputFormat = "yaml"
	UNSET         OutputFormat = ""
)
