package users

type UserInfo struct {
	Username     string
	PasswordHash string
	Role         string
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
