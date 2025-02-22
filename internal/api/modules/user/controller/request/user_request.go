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
}

type UserUpdatePasswordRequest struct {
	Password    string `json:"password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
