package crops

import (
	"encore.dev/beta/errs"
)

func notEqual(expectedError *errs.Error, actualError error) bool {
	if expectedError == nil && actualError == nil {
		return false
	}

	return (actualError != nil && expectedError == nil) ||
		(actualError == nil && expectedError != nil) ||
		(actualError != nil && expectedError != nil && (errs.Code(actualError) != expectedError.Code))
}
