package crops

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"

	"sync"
	"time"
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
	uid := auth.UID(fmt.Sprint(sessionData.UserId))
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

	user, err := insertUser(ctx, params.Email, params.Password)
	if err != nil {
		rlog.Error("insertUser error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "Cannot create user",
		}
	}
	// TODO : Add send verification email to some queue.
	rlog.Info("Sending verification email", "verificationCode", user.verificationCode)

	userResponse := UserResponse{Email: params.Email, IsVerifier: false, IsSeller: false}
	return &userResponse, nil
}

type VerifyUserResponse struct {
	string
}

//encore:api public raw method=GET path=/user/verify/:code
func verifyUser(w http.ResponseWriter, req *http.Request) {

	idSplit := strings.Split(req.URL.Path, "/user/verify/")

	if len(idSplit) < 2 {
		rlog.Error("not enough path parts", "idSplit", idSplit)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := idSplit[1]

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(req.FormValue("timeout"))
	if err == nil {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires.
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	user, err := getUserByVerificationCode(ctx, code)
	if err != nil {
		rlog.Error("userByVerificationCode error!", "err", err)

		response := `<!DOCTYPE html>
			<html>
			    <head>
			        <title>Account verification</title>
    			</head>
			    <body>
			        <p>Account verification failed </p>
			    </body>
			</html>`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}

	user.isVerified = true

	user, err = updateUser(ctx, *user)
	if err != nil {
		rlog.Error("update VerificationCode error!", "err", err)

		response := `<!DOCTYPE html>
			<html>
			    <head>
			        <title>Account verification</title>
    			</head>
			    <body>
			        <p>Account verification failed </p>
			    </body>
			</html>`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}

	response := `<!DOCTYPE html>
			<html>
			    <head>
			        <title>Account verification</title>
			    </head>
    			<body>
			        <p>Your account was verified, You can login.</p>
    			</body>
			</html>`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	return
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token      string
	Email      string
	IsVerifier bool
	IsVerified bool
}

//encore:api public method=POST path=/user/login
func loginUser(ctx context.Context, params *LoginParams) (*LoginResponse, error) {
	user, error := verifyUserLogin(ctx, params.Email, params.Password)
	if error != nil {
		rlog.Error("getting user error!", "error", error)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "cannot login",
		}
	}
	userId, token := sessionManager.startSession(user)
	rlog.Info("authorizatrion", "token", token, "userId", userId)
	loginResponse := LoginResponse{Token: token, Email: user.email, IsVerifier: user.isVerifier, IsVerified: user.isVerified}
	return &loginResponse, nil
}

type TopUpParams struct {
	Amount int64
}

type TopUpResponse struct {
	BalanceAfter int64
}

//encore:api auth method=POST path=/user/topUp
func TopUpUser(ctx context.Context, params *TopUpParams) (*TopUpResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	user, err := getUserById(ctx, sessionData.UserId)
	if err != nil {
		rlog.Error("cannot get user from DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}
	user.balance += params.Amount

	updatedUser, err := updateUser(ctx, *user)
	if err != nil {
		rlog.Error("cannot update in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	response := TopUpResponse{BalanceAfter: updatedUser.balance}
	return &response, nil

}
