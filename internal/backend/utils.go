package backend

import (
	"fmt"

	"github.com/jahvon/tbox/internal/backend/consts"
)

func StrToAuthMode(str string) consts.AuthMode {
	authMode := consts.AuthMode(str)
	if authMode != consts.ModePassword && authMode != consts.ModeMasterKey {
		log.Fatal().Msg(fmt.Sprintf("invalid auth mode %s", str))
	}
	return authMode
}
