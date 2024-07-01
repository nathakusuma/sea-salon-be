package model

type RegisterRequest struct {
	FullName    string `json:"fullName" validate:"required,max=50,min=1"`
	Email       string `json:"email" validate:"required,email,max=320"`
	PhoneNumber string `json:"phoneNumber" validate:"required,max=15,min=1"`
	Password    string `json:"password" validate:"required,max=128,min=8"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=320"`
	Password string `json:"password" validate:"required,max=128"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	FullName string `json:"fullName"`
	IsAdmin  bool   `json:"isAdmin"`
}
