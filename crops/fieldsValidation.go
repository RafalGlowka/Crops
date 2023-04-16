package crops

func createRegistrationNumberValidator() ComposeStringValidator {
	v := ComposeStringValidator{validators: []StringValidator{
		NotEmptyStringValidator{errorMessage: "'RegistrationNumber' is required"},
		LengthStringValidator{minLength: 5, minLengthError: "RegistrationNumber should have at least 5 characters"},
	}}
	return v
}
