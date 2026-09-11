package auth

type RegisterRequest struct {
    Email           string    `json:"email" validate:"required,email,unique=users=email"`
    Title           string    `json:"title" validate:"required,title"`
    FirstName       string    `json:"first_name" validate:"required"`
    MiddleName      string    `json:"middle_name"`
    LastName        string    `json:"last_name"`
    Gender          string    `json:"gender" validate:"required,gender"`
    Bod             string 	  `json:"bod" validate:"required,date_format"`
    Pob             string    `json:"pob" validate:"required"`
    PhoneNumber     string    `json:"phone_number" validate:"required"`
    Password        string    `json:"password" validate:"required,password"`
    ConfirmPassword string    `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token	string	`json:"token"`
}

// @Description Response structure for successful user registration
type RegisterResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Register data berhasil"`
	Data    struct {
		Bod         string `json:"bod" example:"2021-12-01T00:00:00Z"`
		Email       string `json:"email" example:"a@b1qdsd.cd"`
		FirstName   string `json:"first_name" example:"Ujang"`
		Gender      string `json:"gender" example:"L"`
		ID          int    `json:"id" example:"32"`
		PhoneNumber string `json:"phone_number" example:"09"`
		Pob         string `json:"pob" example:"Jalan"`
		Title       string `json:"title" example:"Mrs"`
	} `json:"data"`
}
