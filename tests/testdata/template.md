# {{ .header }}

{{ env "GREETING" }} {{ if eq (env "NAME") "" }}friend{{ else }}{{ env "NAME" }}{{ end }},

{{ .body }}