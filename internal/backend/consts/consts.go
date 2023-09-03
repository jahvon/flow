package consts

type BackendName string
type AuthMode string

const (
	NoAuthBackendName  BackendName = "none"
	KeyringBackendName BackendName = "keyring"
	EnvFileBackendName BackendName = "envFile"

	ModeMasterKey AuthMode = "masterKey"
	ModePassword  AuthMode = "password"
)

var (
	AuthModes      = []AuthMode{ModeMasterKey, ModePassword}
	AuthBackends   = []BackendName{NoAuthBackendName, KeyringBackendName}
	SecretBackends = []BackendName{EnvFileBackendName, KeyringBackendName}
)
