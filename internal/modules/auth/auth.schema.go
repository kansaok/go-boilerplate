package auth

// RegisterRequestSchema defines the structure for a registration request
// @Description Register a new user with email, password, and profile details
type RegisterRequestSchema struct {
	Email                string `json:"email" example:"test@mail.com"`
	Password             string `json:"password" example:"2345678232!Aa"`
	PasswordConfirmation string `json:"password_confirmation" example:"2345678232!Aa"`
	Title                string `json:"title" example:"Mr"`
	FirstName            string `json:"first_name" example:"John"`
	MiddleName           string `json:"middle_name" example:"Lemon"`
	LastName             string `json:"last_name" example:"Djarot"`
	Gender               string `json:"gender" example:"L"`
	Bod                  string `json:"bod" example:"2021-12-01"`
	Pob                  string `json:"pob" example:"2027 North Shannon Drive"`
	PhoneNumber          string `json:"phone_number" example:"689-421-5103"`
}

// RegisterResponseSchema defines the structure for a successful registration response
// @Description Response structure for successful user registration
type RegisterResponseSchema struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Register data berhasil"`
	Data    struct {
		Bod         string `json:"bod" example:"2021-12-01"`
		Email       string `json:"email" example:"a@b1qdsd.cd"`
		FirstName   string `json:"first_name" example:"John"`
		Gender      string `json:"gender" example:"L"`
		ID          int    `json:"id" example:"1"`
		PhoneNumber string `json:"phone_number" example:"689-421-5103"`
		Pob         string `json:"pob" example:"2027 North Shannon Drive"`
		Title       string `json:"title" example:"Mr"`
	} `json:"data"`
}


// LoginRequestSchema defines the structure for a login request
// @Description Login user with email and password
type LoginRequestSchema struct {
	Email                string `json:"email" example:"test@mail.com"`
	Password             string `json:"password" example:"2345678232!Aa"`
}

// LoginResponseSchema defines the structure for a successful login response
// @Description Response structure for successful user login
type LoginResponseSchema struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Register data berhasil"`
	Data    struct {
		Token         string `json:"token" example:"JKJHKJKBJJKjkhjJHBJHbhjBHJGi&Y^*^*&6*&GyuFt&65"`
	} `json:"data"`
}
