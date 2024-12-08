package erro

import (
	"errors"
)

var (
	ErrCertRevoked = errors.New("unauthorized cert revoked")
	ErrParseCert = errors.New("unable to parse x509 cert")
	ErrDecodeCert = errors.New("failed to decode pem-encoded cert")
	ErrDecodeKey = errors.New("error decode rsa key")
	ErrTokenExpired	= errors.New("token expired")
	ErrStatusUnauthorized = errors.New("invalid Token")
	ErrArnMalFormad = errors.New("unauthorized arn scoped malformed")
	ErrBearTokenFormad = errors.New("unauthorized token not informed")
)