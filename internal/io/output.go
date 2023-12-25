package io

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/rs/zerolog"

	"github.com/jahvon/flow/internal/utils"
)

var log = Log()

type OutputFormat string

const (
	OutputFormatJSON       OutputFormat = "json"
	OutputFormatPrettyJSON OutputFormat = "jsonp"
	OutputFormatYAML       OutputFormat = "yaml"
	OutputFormatDefault    OutputFormat = "default"
)

func Log() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
	context := zerolog.New(writer).With().Timestamp()
	return context.Logger()
}

func PrintInfo(message string) {
	log.Info().Msg(message)
}

func PrintNotice(message string) {
	color.HiBlue(message)
}

func PrintQuestion(question string) {
	color.HiCyan(question)
}

func PrintSuccess(message string) {
	color.HiGreen(message)
}

func PrintWarning(message string) {
	color.HiYellow(message)
}

func PrintErrorAndExit(err error) {
	PrintError(err)
	os.Exit(1)
}

func PrintError(err error) {
	color.HiRed(err.Error())
}

func PrintMap(m map[string]string) {
	for k, v := range m {
		if utils.IsMultiLine(v) {
			fmt.Printf("%s:\n%s\n", color.HiBlueString(k), v)
		} else {
			fmt.Printf("%s: %s\n", color.HiBlueString(k), v)
		}
	}
}

func PrintTableData(data [][]string) {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, row := range data {
		_, _ = fmt.Fprintln(tabWriter, strings.Join(row, "\t"))
	}
	_ = tabWriter.Flush()
}

type StdOutWriter struct {
	LogAsDebug bool
	LogFields  map[string]interface{}

	structuredLogBreak bool
}

func (w StdOutWriter) Write(p []byte) (n int, err error) {
	stdOutLog := log.With().Fields(w.LogFields).Logger()
	trimmedP := strings.TrimSpace(string(p))
	if trimmedP == "" {
		return len(p), nil
	}
	splitP := strings.Split(trimmedP, "\n")
	for _, line := range splitP {
		if line == "---break" {
			w.structuredLogBreak = true
			continue
		} else if line == "---endbreak" {
			w.structuredLogBreak = false
			continue
		}

		switch {
		case w.structuredLogBreak:
			_, _ = fmt.Fprintln(os.Stdout, line)
		case w.LogAsDebug:
			stdOutLog.Debug().Msg(line)
		default:
			stdOutLog.Info().Msg(line)
		}
	}

	return len(p), nil
}

type StdErrWriter struct {
	LogAsDebug bool
	LogFields  map[string]interface{}
}

func (w StdErrWriter) Write(p []byte) (n int, err error) {
	stdOutLog := log.With().Fields(w.LogFields).Logger()
	trimmedP := strings.TrimSpace(string(p))
	if trimmedP == "" {
		return len(p), nil
	}

	if w.LogAsDebug {
		stdOutLog.Debug().Msg(string(p))
		return len(p), nil
	}

	stdOutLog.Error().Msg(string(p))
	return len(p), nil
}

func DocsURL(docID string) string {
	return fmt.Sprintf("https://github.com/jahvon/flow/blob/main/docs/%s.md", docID)
}
