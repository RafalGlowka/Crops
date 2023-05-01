package crops

import (
	"context"
	"fmt"
	"testing"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

func verifyUserByEmail(ctx context.Context, email string) bool {
	user, _ := getUserByEmail(ctx, email)
	if user != nil {
		user.isVerified = true
		updateUser(ctx, *user)
		return true
	}
	return false
}

func TestUser(t *testing.T) {
	ctx := context.Background()
	createTests := []struct {
		params   CreateUserParams
		response *UserResponse
		error    *errs.Error
	}{
		{CreateUserParams{Email: "", Password: ""}, nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "'Email' is required",
		}},
		{CreateUserParams{Email: "12345", Password: ""}, nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "Incorrect email format",
		}},
		{CreateUserParams{Email: "test@test.pl", Password: ""}, nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "'Password' is required",
		}},
		{CreateUserParams{Email: "test@test.pl", Password: "123"}, nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "password should have at least 5 characters",
		}},
		{CreateUserParams{Email: "test@test.pl", Password: "12345"}, &UserResponse{Email: "test@test.pl"}, nil},
		{CreateUserParams{Email: "test@test.pl", Password: "12345"}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: " Cannot create user",
		}},
	}

	for _, test := range createTests {
		resp, err := createUser(ctx, &test.params)
		if notEqual(test.error, err) {
			t.Errorf("wrong error %v, expected error: %v", err, test.error)
		} else if (resp != nil && test.response == nil) || (resp == nil && test.response != nil) || (resp != nil && test.response != nil && *resp != *test.response) {
			t.Errorf("wrong response for %v: got %v, want %v", test.params, resp, test.response)
		}
	}

	loginTests1 := []struct {
		params   LoginParams
		response *LoginResponse
		error    *errs.Error
	}{
		{LoginParams{Email: "", Password: ""}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
		{LoginParams{Email: "12345", Password: ""}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
		{LoginParams{Email: "test@test.pl", Password: ""}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
		//  Account require email verification first
		{LoginParams{Email: "test@test.pl", Password: "12345"}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
	}

	for _, test := range loginTests1 {
		resp, err := loginUser(ctx, &test.params)
		if notEqual(test.error, err) {
			t.Errorf("wrong error %v, expected error: %v", err, test.error)
		} else if (resp != nil && test.response == nil) || (resp == nil && test.response != nil) || (resp != nil && test.response != nil && *resp != *test.response) {
			t.Errorf("wrong response for %v: got %v, want %v", test.params, resp, test.response)
		}
	}

	verifyUserByEmail(ctx, "test@test.pl")

	loginTests2 := []struct {
		params   LoginParams
		response *LoginResponse
		error    *errs.Error
	}{
		{LoginParams{Email: "", Password: ""}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
		{LoginParams{Email: "12345", Password: ""}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
		{LoginParams{Email: "test@test.pl", Password: ""}, nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}},
	}

	for _, test := range loginTests2 {
		resp, err := loginUser(ctx, &test.params)
		if notEqual(test.error, err) {
			t.Errorf("wrong error %v, expected error: %v", err, test.error)
		} else if (resp != nil && test.response == nil) || (resp == nil && test.response != nil) || (resp != nil && test.response != nil && *resp != *test.response) {
			t.Errorf("wrong response for %v: got %v, want %v", test.params, resp, test.response)
		}
	}

	_, err := loginUser(ctx, &LoginParams{Email: "test@test.pl", Password: "12345"})
	if err != nil {
		t.Errorf("wrong error %v, expected no error", err)
	}

	// token := resp.Token

	ctx = auth.WithContext(ctx, auth.UID(fmt.Sprint(1)), &SessionData{UserId: 1})

	_, err = TopUpUser(ctx, &TopUpParams{Amount: 10000})
	if err != nil {
		t.Errorf("wrong error %v, expected no error", err)
	}
}
