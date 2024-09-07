package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tuikitIO "github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/types/executable"
)

func copyAllArtifacts(
	logger tuikitIO.Logger,
	artifacts []executable.Artifact,
	wsDir, flowfileDir, dstDir string,
	data, envMap map[string]string,
) error {
	var errs []error
	for i, a := range artifacts {
		if err := copyArtifact(
			logger, fmt.Sprintf("artifact-%d", i), wsDir, flowfileDir, dstDir, a, data, envMap,
		); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Errorf("errors copying artifacts: %v", errs)
	}
	return nil
}

//nolint:gocognit
func copyArtifact(
	logger tuikitIO.Logger,
	name, wsPath, flowfileDir, dstDir string,
	artifact executable.Artifact,
	data, envMap map[string]string,
) error {
	srcPath, err := parseSourcePath(logger, name, flowfileDir, wsPath, artifact, data, envMap)
	if err != nil {
		return errors.Wrap(err, "unable to parse source path")
	}

	if artifact.If != nil {
		eval, err := goTemplateEvaluatedTrue(name, *artifact.If, data)
		if err != nil {
			return errors.Wrap(err, "unable to evaluate if condition")
		}
		if !eval {
			logger.Debugf("skipping artifact %s", name)
			return nil
		}
	}

	srcName := filepath.Base(srcPath)
	if strings.Contains(srcName, "*") {
		matches, err := filepath.Glob(srcPath)
		if err != nil {
			return errors.Wrap(err, "unable to glob source path")
		}
		var errs []error
		for i, match := range matches {
			m := artifact
			m.SrcName = filepath.Base(match)
			m.SrcDir = filepath.Dir(match)
			mErr := copyArtifact(logger, fmt.Sprintf("%s-%d", name, i), wsPath, flowfileDir, dstDir, m, data, envMap)
			if mErr != nil {
				errs = append(errs, mErr)
			}
		}
		if len(errs) > 0 {
			return errors.Errorf("errors copying artifact from pattern: %v", errs)
		}
	}

	info, err := os.Stat(srcPath)
	switch {
	case os.IsNotExist(err):
		return errors.Errorf("file does not exist: %s", srcPath)
	case err != nil:
		return errors.Wrap(err, "unable to stat src file")
	case info.IsDir():
		err := filepath.WalkDir(srcPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			a := artifact
			a.SrcName = filepath.Base(path)
			a.SrcDir = filepath.Dir(path)
			aName := fmt.Sprintf("%s-%s", name, a.SrcName)
			return copyArtifact(logger, aName, wsPath, flowfileDir, dstDir, a, data, envMap)
		})
		if err != nil {
			return errors.Wrap(err, "unable to walk directory")
		}
	}
	relPath, err := filepath.Rel(flowfileDir, srcPath)
	if err != nil {
		return errors.Wrap(err, "unable to get relative path")
	}
	dstPath, err := parseDestinationPath(
		logger,
		name,
		dstDir, relPath, wsPath,
		artifact,
		data, envMap,
	)
	if err != nil {
		return errors.Wrap(err, "unable to parse destination path")
	}

	if err := os.MkdirAll(dstDir, 0750); err != nil {
		if !os.IsExist(err) {
			return errors.Wrap(err, "unable to create destination directory")
		}
		return errors.Wrap(err, "unable to create destination directory")
	}

	logger.Debugx("copying artifact", "name", name, "src", srcPath, "dst", dstPath)
	if _, e := os.Stat(dstPath); e == nil {
		// TODO: Add a flag to overwrite existing files
		logger.Warnx("Overwriting existing file", "dst", dstPath)
	}
	if err := filesystem.CopyFile(srcPath, dstPath); err != nil {
		return errors.Wrap(err, "unable to copy artifact")
	}
	return nil
}
