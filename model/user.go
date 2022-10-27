package model

type User struct {
	Email    string `json:"email,omitempty" validate:"required" binding:"required"`
	Password string `json:"password,omitempty" validate:"required" binding:"required"`
	Role     string `json:"role,omitempty" validate:"required" binding:"required"`
}

type SessionData struct {
	Data map[string]any `json:"data,omitempty" validate:"required" binding:"required"`
}
