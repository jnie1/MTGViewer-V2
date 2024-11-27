package users

import "github.com/jnie1/MTGViewer-V2/database"

func GetUser(email string) (UserInfo, error) {
	user := UserInfo{}

	db := database.Instance()
	row := db.QueryRow(`
		SELECT name, password_hash, role
		FROM users
		WHERE email = $1`, email)

	err := row.Scan(&user.Name, &user.PasswordHash, &user.Role)

	if err != nil {
		return user, err
	}

	user.Email = email
	return user, nil
}

func CreateUser(user UserInfo) error {
	db := database.Instance()

	_, err := db.Exec(`
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)`,
		user.Name, user.Email, user.PasswordHash, user.Role)

	return err
}
