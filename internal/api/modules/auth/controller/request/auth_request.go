package request

type SignUpRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3,max=100"`
	LastName  string `json:"last_name" binding:"required,min=3,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}
