package dto

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type UserIDRequest struct {
	ID string `json:"id" validate:"required"`
}

type GetDataRequest struct {
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	Sort    string                 `json:"sort"`
	Filters map[string]interface{} `json:"filters"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
	Limit       int   `json:"limit"`
}

type ListUserResponse struct {
	List []UserResponse `json:"list"`
	Meta PaginationMeta `json:"meta"`
}
