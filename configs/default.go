package configs

import "math"

var defaultSecurityConfig = securityConfig{
	PasswordEncKey:   []byte("CHANGE_ME_ABCDEFGHIJKLMNOPQRSTUV"),
	PasswordCost:     0,
	TokenSecret:      []byte("CHANGE_ME"),
	TokenDuration:    math.MaxInt64,
	DisableTokenAuth: false,
	MjApiKeyPublic:   "4ab89bb2a87bc4d190e57c73c53c6705",
	MjApiKeyPrivate:  "2919a718ae663255682c48afbad358a5",
}
