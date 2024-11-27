package users

type UserInfo struct {
	Name         string
	Email        string
	PasswordHash string
	Role         string
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
