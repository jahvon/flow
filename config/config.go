package config

import "github.com/jahvon/flow/internal/io"

var log = io.Log().With().Str("scope", "config").Logger()
