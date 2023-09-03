package parameter

import (
	"fmt"
	"regexp"
	"strings"

	tboxIO "github.com/jahvon/tbox/internal/io"
)

type Destination string

type Value struct {
	Text string `yaml:"text"`
	Ref  string `yaml:"secretRef"`
}

type Parameter struct {
	Key      string                 `yaml:"key"`
	Value    *Value                 `yaml:"value"`
	ExposeAs map[Destination]string `yaml:"exposeAs"`
}

const (
	DestinationEnv  Destination = "env"
	DestinationFile Destination = "file"

	ReservedPrefix = "TBOX_"
)

var log = tboxIO.Log()

func (p *Parameter) Validate() error {
	if p.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !re.MatchString(p.Key) {
		return fmt.Errorf("key must be alphanumeric and can only contain underscores characters")
	}

	if strings.HasPrefix(NormalizeKey(p.Key), NormalizeKey(ReservedPrefix)) {
		return fmt.Errorf("key cannot start with reserved prefix '%s'", ReservedPrefix)
	}

	if p.Value == nil {
		return fmt.Errorf("must set parameter value for key %s", p.Key)
	}

	if p.Value.Text == "" && p.Value.Ref == "" {
		return fmt.Errorf("must set either text or secretRef for key %s", p.Key)
	} else if p.Value.Text != "" && p.Value.Ref != "" {
		return fmt.Errorf("cannot set both text and secretRef for key %s", p.Key)
	}

	return nil
}

func LookupParameter(parameters []*Parameter, key string) (*Parameter, bool) {
	for _, p := range parameters {
		if KeysEqual(p.Key, NormalizeKey(key)) {
			return p, true
		}
	}
	return nil, false
}

func StrToDestination(destination string) (Destination, bool) {
	switch destination {
	case string(DestinationEnv):
		return DestinationEnv, true
	case string(DestinationFile):
		return DestinationFile, true
	default:
		return "", false
	}
}

func NormalizeKey(key string) string {
	return strings.ToUpper(strings.Trim(key, "_ "))
}

func KeysEqual(key1, key2 string) bool {
	return NormalizeKey(key1) == NormalizeKey(key2)
}
