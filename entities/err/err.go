package err

import "errors"

var (
	NoRecord             = errors.New("no records found")
	InternalServiceError = errors.New("internal service error")
	DuplicateEntry       = errors.New("duplicate entry")
	LoadPkg              = errors.New("load package error")
	RCMustBeSuccess      = errors.New("RC must be 0000")
	FieldMustBeExist     = errors.New("Field cannot be empty")
)
