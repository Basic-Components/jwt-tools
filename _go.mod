module github.com/Basic-Components/jwt-signer

require (
	github.com/Basic-Components/jwt-signer/config v0.0.0
	github.com/Basic-Components/jwt-signer/errs v0.0.0
	github.com/Basic-Components/jwt-signer/jwtrpc v0.0.0
	github.com/Basic-Components/jwt-signer/jwtsigner v0.0.0
	github.com/Basic-Components/jwt-signer/jwtverifier v0.0.0
	github.com/Basic-Components/jwt-signer/keygen v0.0.0
	github.com/Basic-Components/jwt-signer/logger v0.0.0
	github.com/Basic-Components/jwt-signer/server v0.0.0
	github.com/Basic-Components/jwt-signer/signals v0.0.0
	github.com/Basic-Components/jwt-signer/utils v0.0.0
)

replace (
	github.com/Basic-Components/jwt-signer/config v0.0.0 => ./config
	github.com/Basic-Components/jwt-signer/errs v0.0.0 => ./errs
	github.com/Basic-Components/jwt-signer/jwtrpc v0.0.0 => ./jwtrpc
	github.com/Basic-Components/jwt-signer/jwtsigner v0.0.0 => ./jwtsigner
	github.com/Basic-Components/jwt-signer/jwtverifier v0.0.0 => ./jwtverifier
	github.com/Basic-Components/jwt-signer/keygen v0.0.0 => ./keygen
	github.com/Basic-Components/jwt-signer/logger v0.0.0 => ./logger
	github.com/Basic-Components/jwt-signer/server v0.0.0 => ./server
	github.com/Basic-Components/jwt-signer/signals v0.0.0 => ./signals
	github.com/Basic-Components/jwt-signer/utils v0.0.0 => ./utils
)

go 1.12
