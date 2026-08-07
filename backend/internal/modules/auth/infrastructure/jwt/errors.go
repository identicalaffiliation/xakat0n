package jwt

import "errors"

var ErrPrivateKeyRequired = errors.New("jwt private key is required")
