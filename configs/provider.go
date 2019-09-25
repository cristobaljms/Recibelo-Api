package configs

import "time"

type securityConfig struct {
	PasswordEncKey   []byte
	PasswordCost     int64
	TokenSecret      []byte
	TokenDuration    time.Duration
	DisableTokenAuth bool
	MjApiKeyPublic   string
	MjApiKeyPrivate  string
}

var SecurityCfg = defaultSecurityConfig
