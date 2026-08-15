package users

import "github.com/jnie1/MTGViewer-V2/database"

func GetUser(email string) (UserInfo, error) {
	db := database.Instance()
	row := db.QueryRow(`
		SELECT email, password_hash, user_role
		FROM users
		WHERE email = $1`, email)

	user := UserInfo{}
	if err := row.Scan(&user.Email, &user.PasswordHash, &user.Role); err != nil {
		return user, err
	}

	user.Email = email
	return user, nil
}

func CreateUser(user UserInfo) error {
	db := database.Instance()
	_, err := db.Exec(`
		INSERT INTO users ( email, password_hash, user_role)
		VALUES ($1, $2, $3)`,
		user.Email, user.PasswordHash, user.Role)

	return err
}
