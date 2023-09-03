package parameter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jahvon/tbox/internal/backend"
)

func (p *Parameter) Expose(context, masterKey string, secretBackend backend.SecretBackend) error {
	value, err := p.Data(context, masterKey, secretBackend)
	if err != nil {
		return err
	}

	for destination, setting := range p.ExposeAs {
		if destination == DestinationEnv {
			envVar := NormalizeKey(setting)
			return os.Setenv(envVar, value)
		} else if destination == DestinationFile {
			filename := strings.ToLower(p.Key)
			filename = strings.Trim(filename, "_")
			destinationPath := filepath.Join(setting, filename)
			return os.WriteFile(destinationPath, []byte(value), 0644)
		}
	}

	return nil
}
