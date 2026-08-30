package users

import (
	"context"

	"github.com/jnie1/MTGViewer-V2/database"
)

func GetUser(ctx context.Context, username string) (UserInfo, error) {
	db := database.Instance()
	row := db.QueryRowContext(ctx, `
		SELECT username, password_hash, user_role
		FROM users
		WHERE username = $1`, username)

	var user UserInfo
	if err := row.Scan(&user.Username, &user.PasswordHash, &user.Role); err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(ctx context.Context, user UserInfo) error {
	db := database.Instance()
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, user_role)
		VALUES ($1, $2, $3)`,
		user.Username, user.PasswordHash, user.Role)

	return err
}
