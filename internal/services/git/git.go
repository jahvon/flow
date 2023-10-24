package git

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/jahvon/flow/internal/services/run"
)

func Pull(repoDir string) error {
	if info, err := os.Stat(repoDir); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("git repo %s does not exist", repoDir)
	} else if err != nil {
		return fmt.Errorf("unable to check for git repo %s - %w", repoDir, err)
	} else if !info.IsDir() {
		return fmt.Errorf("git repo %s is not a directory", repoDir)
	}

	if err := run.RunCmd("git pull", repoDir, nil); err != nil {
		return fmt.Errorf("unable to pull git repo %s - %w", repoDir, err)
	}

	log.Info().Msgf("successfully pulled git repo %s", repoDir)
	return nil
}
