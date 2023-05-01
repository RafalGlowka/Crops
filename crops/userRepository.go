package crops

import (
	"context"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type User struct {
	id               int64  // id
	email            string // email address
	passwordHash     string // passwordHash
	isVerifier       bool
	isVerified       bool
	verificationCode string
	balance          int64
}

// insert inserts user into the database.
func insertUser(ctx context.Context, email, password string) (*User, error) {
	passwordHash := sessionManager.getPasswordHash(password)
	rows, err := sqldb.Query(ctx, `
		INSERT INTO users (email, passwordHash, verifier, verified, balance)
		VALUES ($1, $2, false, false, 0) RETURNING id, verificationcode
	`, email, passwordHash)
	if err != nil {
		rlog.Error("inserting user failed", "err", err)
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	u := User{email: email, isVerified: false, isVerifier: false, balance: 0}
	err = rows.Scan(&u.id, &u.verificationCode)
	if err != nil {
		rlog.Error("reading user failed", "err", err)
		return nil, err
	}

	return &u, nil

}

func verifyUserLogin(ctx context.Context, email, password string) (*User, error) {
	u := &User{}
	err := sqldb.QueryRow(ctx, `
		SELECT id, email, passwordHash, verifier, verified, verificationcode, balance FROM users
		WHERE email = $1
	`, email).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isVerified, &u.verificationCode, &u.balance)
	if err == nil {
		if u.isVerified == false {
			return nil, &errs.Error{
				Code:    errs.PermissionDenied,
				Message: "Verifie account first !",
			}
		}

		passwordHash := sessionManager.getPasswordHash(password)
		if u.passwordHash == passwordHash {
			return u, nil
		}
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Incorrect password",
		}
	}

	return u, err
}

func getUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := sqldb.QueryRow(ctx, `
		SELECT id, email, passwordHash, verifier, verified, verificationcode, balance FROM users
		WHERE email = $1
	`, email).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isVerified, &u.verificationCode, &u.balance)
	return u, err
}

func getUserById(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	err := sqldb.QueryRow(ctx, `
		SELECT id, email, passwordHash, verifier, verified, verificationcode, balance FROM users
		WHERE id = $1
	`, id).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isVerified, &u.verificationCode, &u.balance)
	return u, err
}

func getUserByVerificationCode(ctx context.Context, verificationCode string) (*User, error) {
	u := &User{}
	err := sqldb.QueryRow(ctx, `
		SELECT id, email, passwordHash, verifier, verified, verificationcode, balance FROM users
		WHERE verificationcode = $1
	`, verificationCode).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isVerified, &u.verificationCode, &u.balance)
	return u, err
}

func updateUser(ctx context.Context, user User) (*User, error) {
	_, err := sqldb.Exec(ctx, `
		UPDATE users 
		SET email = $1, passwordHash = $2, verifier = $3, verified = $4, verificationcode = $5, balance = $6
		WHERE id = $7
	`, user.email, user.passwordHash, user.isVerifier, user.isVerified, user.verificationCode, user.balance, user.id)
	if err != nil {
		rlog.Error("update user failed", "err", err)
		return nil, err
	}

	return &user, nil
}
