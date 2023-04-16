package crops

func createEmailValidator() ComposeStringValidator {
	v := ComposeStringValidator{validators: []StringValidator{
		NotEmptyStringValidator{errorMessage: " 'Email' is required"},
		RegexpStringValidator{
			regexp:       "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$",
			errorMessage: "Incorrect email format",
		},
	}}
	return v
}

func createPasswordValidator() ComposeStringValidator {
	v := ComposeStringValidator{validators: []StringValidator{
		NotEmptyStringValidator{errorMessage: "'Password' is required"},
		LengthStringValidator{minLength: 5, minLengthError: "password should have at least 5 characters"},
	}}
	return v
}
