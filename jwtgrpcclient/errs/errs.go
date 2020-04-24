package errs

import (
	"errors"
)

var ConfigDecodeError error = errors.New("unable to decode config into struct")
var ConfigVerifyError error = errors.New("config not satisfied the schema")

var SignerTypeError error = errors.New("unknown algo type key")
var LoadKeyError error = errors.New("couldn't read key")
var ParseClaimsToJsonError error = errors.New("Couldn't parse claims JSON")
var SignTokenError error = errors.New("Error signing token")
var TokenInvalidError error = errors.New("Token is invalid")
var VerifyTokenError error = errors.New("Verify Token error")
