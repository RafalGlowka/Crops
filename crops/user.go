package crops

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"

	"sync"
)

var sessionManager *SessionManager = &SessionManager{
	passwordSalt: "secret",
	lock:         sync.Mutex{},
	tokens:       map[string]SessionData{},
}

type UserResponse struct {
	Email      string
	IsVerifier bool
	IsSeller   bool
}

// TODO : Session verification
//
// AuthHandler can be named whatever you prefer (but must be exported).
//
//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, *SessionData, error) {
	// Validate the token and look up the user id and user data,
	// for example by calling Firebase Auth.
	sessionData := sessionManager.getSessionData(token)
	if sessionData == nil {
		return "", nil, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "invalid token",
		}
	}
	uid := auth.UID(sessionData.UserId)
	return uid, sessionData, nil
}

type CreateUserParams struct {
	Email    string
	Password string
}

//encore:api public method=POST path=/user/create
func createUser(ctx context.Context, params *CreateUserParams) (*UserResponse, error) {
	emailValidationResult := createEmailValidator().validate(params.Email)
	if emailValidationResult != nil {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: emailValidationResult.errorMessage,
		}
	}

	passwordValidationResult := createPasswordValidator().validate(params.Password)
	if passwordValidationResult != nil {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: passwordValidationResult.errorMessage,
		}
	}

	if err := insertUser(ctx, params.Email, params.Password); err != nil {
		rlog.Error("insertUser error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "Cannot create user",
		}
	}
	userResponse := UserResponse{Email: params.Email, IsVerifier: false, IsSeller: false}
	return &userResponse, nil
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token      string
	Email      string
	IsVerifier bool
	IsSeller   bool
}

//encore:api public method=POST path=/user/login
func loginUser(ctx context.Context, params *LoginParams) (*LoginResponse, error) {
	user, error := verifyUser(ctx, params.Email, params.Password)
	if error != nil {
		rlog.Error("getting user error!", "error", error)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}
	}
	userId, token := sessionManager.startSession(user)
	rlog.Info("authorizatrion", "token", token, "userId", userId)
	loginResponse := LoginResponse{Token: token, Email: user.email, IsVerifier: user.isVerifier, IsSeller: user.isSeller}
	return &loginResponse, nil
}
