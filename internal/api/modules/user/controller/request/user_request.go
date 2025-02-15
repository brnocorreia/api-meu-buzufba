package request

type UserRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3,max=100"`
	LastName  string `json:"last_name" binding:"required,min=3,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=3,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,min=3,max=100"`
	Password  string `json:"password" binding:"omitempty,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
