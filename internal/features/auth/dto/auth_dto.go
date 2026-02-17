package authdto

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user1@yopmail.com"`
	Password string `json:"password" validate:"required" example:"muslimlypass"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
