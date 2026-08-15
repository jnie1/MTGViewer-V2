package users

import (
	"context"

	"github.com/jnie1/MTGViewer-V2/database"
)

func GetUser(ctx context.Context, email string) (UserInfo, error) {
	db := database.Instance()
	row := db.QueryRowContext(ctx, `
		SELECT name, password_hash, role
		FROM users
		WHERE email = $1`, email)

	user := UserInfo{}
	if err := row.Scan(&user.Name, &user.PasswordHash, &user.Role); err != nil {
		return user, err
	}

	user.Email = email
	return user, nil
}

func CreateUser(ctx context.Context, user UserInfo) error {
	db := database.Instance()
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)`,
		user.Name, user.Email, user.PasswordHash, user.Role)

	return err
}
