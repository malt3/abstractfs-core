package validate

import (
	"fmt"
	"path"

	"github.com/malt3/abstractfs-core/api"
)

type TreeValidationError struct {
	RootErrors []ValidationErrorReason
	NodeErrors []StatValidationError
}

func (e TreeValidationError) Error() string {
	var errStr string
	for i, err := range e.RootErrors {
		if i > 0 {
			errStr += ", "
		}
		errStr += err.Error()
	}
	for _, err := range e.NodeErrors {
		if errStr != "" {
			errStr += ", "
		}
		errStr += err.Error()
	}
	return "validating tree: " + errStr
}

type StatValidationError struct {
	Dir    string
	Stat   api.Stat
	Reason ValidationErrorReason
}

func (e StatValidationError) Error() string {
	return fmt.Sprintf("validating %q: %s", path.Join("/", e.Dir, e.Stat.Name), e.Reason)
}

type ValidationErrorReason string

func (e ValidationErrorReason) Error() string {
	return string(e)
}

const (
	RootNodeMissing       ValidationErrorReason = "root node is missing"
	RootNodeKindIsNotDir  ValidationErrorReason = "root node kind is not dir"
	RootNodeNameNotEmpty  ValidationErrorReason = "root node name is not empty"
	NodeNameEmpty         ValidationErrorReason = "node name is empty"
	NodeNameContainsSlash ValidationErrorReason = "node name contains slash"
	FlatNameNotRooted     ValidationErrorReason = "flat name must be rooted (start with slash)"
	TreeNameRooted        ValidationErrorReason = "tree name must not start with slash"
	PathComponentEmpty    ValidationErrorReason = "path component is empty"
	PathComponentDot      ValidationErrorReason = "path component is '.'"
	PathComponentDotDot   ValidationErrorReason = "path component is '..'"
	KindInvalid           ValidationErrorReason = "kind is invalid"
	SizeInvalid           ValidationErrorReason = "size is invalid"
	PayloadEmpty          ValidationErrorReason = "payload is empty"
	PayloadNotEmpty       ValidationErrorReason = "payload must be empty"
	PayloadInvalidSRI     ValidationErrorReason = "payload is invalid SRI"
)

type validationMode string

const (
	ValidationModeFlat validationMode = "flat"
	ValidationModeTree validationMode = "tree"
)
