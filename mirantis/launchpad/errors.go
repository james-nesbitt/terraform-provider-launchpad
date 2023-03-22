package launchpad

import "errors"

var (
	ErrEmptyUsernamePass      = errors.New("no username or password provided in MKE client")
	ErrEmptyEndpoint          = errors.New("no endpoint provided in MKE client")
	ErrRequestCreation        = errors.New("error creating request in MKE client")
	ErrMarshaling             = errors.New("error occured while marshalling struct in MKE client")
	ErrUnmarshaling           = errors.New("error occured while unmarshalling struct in MKE client")
	ErrEmptyResError          = errors.New("request returned empty ResponseError struct in MKE client")
	ErrResponseError          = errors.New("request returned ResponseError in MKE client")
	ErrUnauthorizedReq        = errors.New("unauthorized request in MKE client")
	ErrUnknownTarget          = errors.New("unknown API target")
	ErrServerError            = errors.New("server error occured")
	ErrEmptyStruct            = errors.New("empty struct passed in MKE client")
	ErrInvalidFilter          = errors.New("passing invalid account retrieval filter in MKE client")
	ErrMissingConnectionBlock = errors.New("missing connection block for host")
)
