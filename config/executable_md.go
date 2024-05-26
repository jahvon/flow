package config

import (
	"fmt"
	"strings"
)

func execMarkdown(e *Executable) string {
	var mkdwn string
	mkdwn += fmt.Sprintf("# [Executable] %s\n", e.Ref())
	if e.Description != "" {
		mkdwn += "| \n"
		lines := strings.Split(e.Description, "\n")
		for _, line := range lines {
			mkdwn += fmt.Sprintf("| %s\n", line)
		}
		mkdwn += "| \n\n"
	}
	if e.Visibility != nil {
		mkdwn += fmt.Sprintf("**Visibility:** %s\n", e.Visibility)
	}
	if e.Timeout != 0 {
		mkdwn += fmt.Sprintf("**Timeout:** %s\n", e.Timeout.String())
	}
	if len(e.Aliases) > 0 {
		mkdwn += "**Aliases**\n"
		for _, alias := range e.Aliases {
			mkdwn += fmt.Sprintf("- `%s`\n", alias)
		}
		mkdwn += "\n"
	}

	if len(e.Tags) > 0 {
		mkdwn += "**Tags**\n"
		for _, tag := range e.Tags {
			mkdwn += fmt.Sprintf("- `%s`\n", tag)
		}
		mkdwn += "\n"
	}

	mkdwn += execTypeMarkdown(e.Type)
	mkdwn += fmt.Sprintf("\n\n_Executable can be found in_ [%s](%s)\n", e.definitionPath, e.definitionPath)
	return mkdwn
}

func execTypeMarkdown(spec *ExecutableTypeSpec) string {
	var mkdwn string
	switch {
	case spec == nil:
		mkdwn += "No executable type found\n"
	case spec.Exec != nil:
		mkdwn += shellExecMarkdown(spec.Exec)
	case spec.Launch != nil:
		mkdwn += launchExecMarkdown(spec.Launch)
	case spec.Request != nil:
		mkdwn += requestExecMarkdown(spec.Request)
	case spec.Render != nil:
		mkdwn += renderExecMarkdown(spec.Render)
	case spec.Serial != nil:
		mkdwn += serialExecMarkdown(spec.Serial)
	case spec.Parallel != nil:
		mkdwn += parallelExecMarkdown(spec.Parallel)
	default:
		mkdwn += "**generated markdown not supported for type**\n"
	}
	return mkdwn
}

func shellExecMarkdown(s *ExecExecutableType) string {
	if s == nil {
		return ""
	}
	mkdwn := "## Shell Configuration\n"
	if s.Directory != "" {
		mkdwn += fmt.Sprintf("**Executed from:** `%s`\n", s.Directory)
	}
	if s.LogMode != "" {
		mkdwn += fmt.Sprintf("**Log Mode:** %s\n", s.LogMode)
	}
	if s.Command != "" {
		mkdwn += fmt.Sprintf("**Command**\n```sh\n%s\n```\n", s.Command)
	} else if s.File != "" {
		mkdwn += fmt.Sprintf("**File:** `%s`\n", s.File)
	}
	mkdwn += execEnvTable(&s.ExecutableEnvironment)

	return mkdwn
}

func launchExecMarkdown(l *LaunchExecutableType) string {
	if l == nil {
		return ""
	}
	mkdwn := "## Launch Configuration\n"
	if l.App != "" {
		mkdwn += fmt.Sprintf("**App:** `%s`\n", l.App)
	}
	if l.URI != "" {
		mkdwn += fmt.Sprintf("**URI:** [%s](%s)\n", l.URI, l.URI)
	}
	if l.Wait {
		mkdwn += "**Wait:** enabled\n"
	}
	mkdwn += execEnvTable(&l.ExecutableEnvironment)
	return mkdwn
}

func requestExecMarkdown(r *RequestExecutableType) string {
	if r == nil {
		return ""
	}
	mkdwn := "## Request Configuration\n"
	mkdwn += fmt.Sprintf("**Method:** %s\n", r.Method)
	mkdwn += fmt.Sprintf("**URL:** [%s](%s)\n", r.URL, r.URL)

	if r.Timeout != 0 {
		mkdwn += fmt.Sprintf("**Request Timeout:** %s\n", r.Timeout)
	}
	if r.LogResponse {
		mkdwn += "**Log Response:** enabled\n"
	}
	if r.Body != "" {
		mkdwn += fmt.Sprintf("**Body:**\n```\n%s\n```\n", r.Body)
	}

	if len(r.Headers) > 0 {
		mkdwn += "\n**Headers**\n"
		for k, v := range r.Headers {
			mkdwn += fmt.Sprintf("- %s: %s\n", k, v)
		}
	}
	if len(r.ValidStatusCodes) > 0 {
		mkdwn += "**Accepted Status Codes**\n"
		for _, code := range r.ValidStatusCodes {
			mkdwn += fmt.Sprintf("- %d\n", code)
		}
	}

	if r.ResponseFile != nil {
		mkdwn += fmt.Sprintf("**Resonse Saved To:** %s\n", r.ResponseFile.Filename)
		if r.ResponseFile.SaveAs != "" {
			mkdwn += fmt.Sprintf("**Response Saved As:** %s\n", r.ResponseFile.SaveAs)
		}
	}
	if r.TransformResponse != "" {
		mkdwn += fmt.Sprintf("**Transformation Expression:** ```\n%s\n```\n", r.TransformResponse)
	}

	mkdwn += execEnvTable(&r.ExecutableEnvironment)
	return mkdwn
}

func renderExecMarkdown(r *RenderExecutableType) string {
	if r == nil {
		return ""
	}

	mkdwn := "## Render Configuration\n"
	if r.Directory != "" {
		mkdwn += fmt.Sprintf("**Executed from:** `%s`\n", r.Directory)
	}
	if r.TemplateFile != "" {
		mkdwn += fmt.Sprintf("**Template File:** `%s`\n", r.TemplateFile)
	}
	if r.TemplateDataFile != "" {
		mkdwn += fmt.Sprintf("**Template Data File:** `%s`\n", r.TemplateDataFile)
	}
	mkdwn += execEnvTable(&r.ExecutableEnvironment)
	return mkdwn
}

func serialExecMarkdown(s *SerialExecutableType) string {
	if s == nil {
		return ""
	}
	mkdwn := "## Serial Configuration\n"
	if s.FailFast {
		mkdwn += "**Fail Fast:** enabled\n"
	}
	mkdwn += "**Executables**\n"
	for i, ref := range s.ExecutableRefs {
		mkdwn += fmt.Sprintf("%d. %s\n", i+1, ref)
	}
	mkdwn += execEnvTable(&s.ExecutableEnvironment)
	return mkdwn
}

func parallelExecMarkdown(p *ParallelExecutableType) string {
	if p == nil {
		return ""
	}
	mkdwn := "## Parallel Configuration\n"
	if p.MaxThreads > 0 {
		mkdwn += fmt.Sprintf("**Max Threads:** %d\n", p.MaxThreads)
	}
	if p.FailFast {
		mkdwn += "**Fail Fast:** enabled\n"
	}
	mkdwn += "**Executables**\n"
	for _, ref := range p.ExecutableRefs {
		mkdwn += fmt.Sprintf("- %s\n", ref)
	}
	mkdwn += execEnvTable(&p.ExecutableEnvironment)
	return mkdwn
}

func execEnvTable(env *ExecutableEnvironment) string {
	var table string
	if len(env.Parameters) > 0 {
		table += "### Parameters\n"
		table += "| Env Key | Type | Value |\n| --- | --- | --- |\n"
		for _, p := range env.Parameters {
			var valueType, valueInput string
			switch {
			case p.Text != "":
				valueType = "text"
				valueInput = p.Text
			case p.SecretRef != "":
				valueType = "secret"
				valueInput = p.SecretRef
			case p.Prompt != "":
				valueType = "prompt"
				valueInput = p.Prompt
			}
			table += fmt.Sprintf("| `%s` | %s | %s |\n", p.EnvKey, valueType, valueInput)
		}
	}

	if len(env.Args) > 0 {
		table += "### Arguments\n"
		table += "| Env Key | Arg Type | Input Type | Default | Required |\n| --- | --- | --- | --- | --- |\n"
		for _, a := range env.Args {
			var argType string
			switch {
			case a.Pos != 0:
				argType = "positional"
			case a.Flag != "":
				argType = "flag"
			}
			table += fmt.Sprintf(
				"| `%s` | %s | %s | %s | %t |\n",
				a.EnvKey, argType, a.Type, a.Default, a.Required,
			)
		}
	}
	return table
}
