package crops

import (
	"context"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type User struct {
	id           int64  // id
	email        string // email address
	passwordHash string // passwordHash
	isVerifier   bool
	isSeller     bool
}

// insert inserts user into the database.
func insertUser(ctx context.Context, email, password string) error {
	passwordHash := sessionManager.getPasswordHash(password)
	_, err := sqldb.Exec(ctx, `
		INSERT INTO users (email, passwordHash, verifier, seller)
		VALUES ($1, $2, false, false)
	`, email, passwordHash)
	return err
}

func verifyUser(ctx context.Context, email, password string) (*User, error) {
	u := &User{}
	err := sqldb.QueryRow(ctx, `
		SELECT id, email, passwordHash, verifier, seller FROM users
		WHERE email = $1
	`, email).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isSeller)
	if err == nil {
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
		SELECT id, email, passwordHash, verifier, seller FROM users
		WHERE email = $1
	`, email).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isSeller)
	return u, err
}

func getUserById(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := sqldb.QueryRow(ctx, `
		SELECT id, email, passwordHash, verifier, seller FROM users
		WHERE id = $1
	`, id).Scan(&u.id, &u.email, &u.passwordHash, &u.isVerifier, &u.isSeller)
	return u, err
}
