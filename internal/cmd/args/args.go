package args

import (
	"strings"
)

func ParseArgs(args []string) (flagArgs map[string]string, posArgs []string) {
	flagArgs = make(map[string]string)
	posArgs = make([]string, 0)
	for i := 0; i < len(args); i++ {
		split := strings.Split(args[i], "=")
		if len(split) >= 2 {
			flagArgs[split[0]] = strings.Join(split[1:], "=")
			continue
		}
		posArgs = append(posArgs, args[i])
	}
	return
}
