package logs

import (
	"encoding/json"
	"fmt"

	"github.com/flowexec/tuikit/io"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io/common"
)

type entry struct {
	Args string `json:"args" yaml:"args"`
	Time string `json:"time" yaml:"time"`
	File string `json:"file" yaml:"file"`
}

type entryResponse struct {
	Logs []entry `json:"logs" yaml:"logs"`
}

func tuikitToEntry(e io.ArchiveEntry) entry {
	return entry{
		Args: e.Args,
		Time: e.Time.String(),
		File: e.Path,
	}
}

func marshalEntriesJSON(entries []io.ArchiveEntry) ([]byte, error) {
	entriesJSON := make([]entry, len(entries))
	for i, e := range entries {
		entriesJSON[i] = tuikitToEntry(e)
	}
	entriesResponse := entryResponse{Logs: entriesJSON}
	return json.MarshalIndent(entriesResponse, "", "  ")
}

func marshalEntriesYAML(entries []io.ArchiveEntry) ([]byte, error) {
	entriesYAML := make([]entry, len(entries))
	for i, e := range entries {
		entriesYAML[i] = tuikitToEntry(e)
	}
	entriesResponse := entryResponse{Logs: entriesYAML}
	return yaml.Marshal(entriesResponse)
}

func PrintEntries(ctx *context.Context, format string, entries []io.ArchiveEntry) {
	logger := ctx.Logger
	logger.Debugf("listing %d log entries", len(entries))
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := marshalEntriesYAML(entries)
		if err != nil {
			logger.Fatalf("Failed to marshal log entries - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), string(str))
	case common.JSONFormat:
		str, err := marshalEntriesJSON(entries)
		if err != nil {
			logger.Fatalf("Failed to marshal log entries - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), string(str))
	}
}
