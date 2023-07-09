package validate

import (
	"strings"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/sri"
)

func ValidateStat(dir string, stat api.Stat, validationMode validationMode) []StatValidationError {
	var statErrors []StatValidationError
	// only root node should have empty name
	if len(stat.Name) == 0 && (dir != "" || validationMode == ValidationModeFlat) {
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: NodeNameEmpty})
	}
	startsWithSlash := strings.HasPrefix(stat.Name, "/")
	switch {
	case validationMode == ValidationModeFlat && !startsWithSlash:
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: FlatNameNotRooted})
	case validationMode == ValidationModeTree && startsWithSlash:
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: TreeNameRooted})
	}
	containsSlashes := strings.Contains(stat.Name, "/")
	if containsSlashes && validationMode == ValidationModeTree {
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: NodeNameContainsSlash})
	}
	if validationMode == ValidationModeFlat {
		unrooted := stat.Name
		if strings.HasPrefix(unrooted, "/") {
			unrooted = unrooted[1:]
		}
		path := strings.Split(unrooted, "/")
		if len(path) == 1 && path[0] == "" {
			path = nil
		}
		for _, name := range path {
			if len(name) == 0 {
				statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: PathComponentEmpty})
			}
			if name == "." {
				statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: PathComponentDot})
			}
			if name == ".." {
				statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: PathComponentDotDot})
			}
		}
	}

	kindErrors := validateKind(stat)
	for _, err := range kindErrors {
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: err})
	}

	payloadError := validatePayload(stat)
	for _, err := range payloadError {
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: err})
	}

	sizeError := validateSize(stat)
	for _, err := range sizeError {
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: err})
	}

	attributeErrors := validateAttributes(stat)
	for _, err := range attributeErrors {
		statErrors = append(statErrors, StatValidationError{Dir: dir, Stat: stat, Reason: err})
	}

	if len(statErrors) == 0 {
		return nil
	}
	return statErrors
}

func validateKind(stat api.Stat) []ValidationErrorReason {
	// TODO: extend kinds
	switch stat.Kind {
	case api.KindDirectory, api.KindRegular, api.KindSymlink:
		return nil
	}
	return []ValidationErrorReason{KindInvalid}
}

func validateSize(stat api.Stat) []ValidationErrorReason {
	if stat.Kind == api.KindRegular || stat.Size == 0 {
		return nil
	}
	return []ValidationErrorReason{SizeInvalid}
}

func validatePayload(stat api.Stat) []ValidationErrorReason {
	var errors []ValidationErrorReason
	if len(stat.Payload) == 0 &&
		(stat.Kind == api.KindRegular || stat.Kind == api.KindSymlink) {
		errors = append(errors, PayloadEmpty)
	}
	if len(stat.Payload) > 0 &&
		stat.Kind == api.KindDirectory {
		errors = append(errors, PayloadNotEmpty)
	}
	if _, err := sri.FromString(stat.Payload); stat.Kind == api.KindRegular && err != nil {
		errors = append(errors, PayloadInvalidSRI)
	}
	return errors
}

func validateAttributes(stat api.Stat) []ValidationErrorReason {
	// TODO: validate attributes
	return nil
}
