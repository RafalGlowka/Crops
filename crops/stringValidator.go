package crops

import (
	"regexp"
	"strings"
)

type ValidationIssue struct {
	errorMessage string
}

type StringValidator interface {
	validate(string) *ValidationIssue
}

type NotEmptyStringValidator struct {
	errorMessage string
}

func (v NotEmptyStringValidator) validate(data string) *ValidationIssue {
	if len(strings.TrimSpace(data)) < 1 {
		return &ValidationIssue{errorMessage: v.errorMessage}
	}
	return nil
}

type RegexpStringValidator struct {
	regexp       string
	errorMessage string
}

func (v RegexpStringValidator) validate(data string) *ValidationIssue {
	matched, _ := regexp.MatchString(v.regexp, data)
	if matched {
		return nil
	}
	return &ValidationIssue{errorMessage: v.errorMessage}
}

type LengthStringValidator struct {
	minLength      int
	minLengthError string
	maxLength      int
	maxLengthError string
}

func (v LengthStringValidator) validate(data string) *ValidationIssue {
	if v.minLength > 0 && len(data) < v.minLength {
		return &ValidationIssue{errorMessage: v.minLengthError}
	}
	if v.maxLength > 0 && len(data) > v.maxLength {
		return &ValidationIssue{errorMessage: v.maxLengthError}
	}
	return nil
}

type ComposeStringValidator struct {
	validators []StringValidator
}

func (v ComposeStringValidator) validate(data string) *ValidationIssue {
	for i := 0; i < len(v.validators); i++ {
		result := v.validators[i].validate(data)
		if result != nil {
			return result
		}
	}
	return nil
}
