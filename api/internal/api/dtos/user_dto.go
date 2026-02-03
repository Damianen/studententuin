package dtos

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name string `json:"name" binding:"required"`
}

type LoginUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
	Name  *string `json:"name,omitempty" binding:"omitempty"`
}
