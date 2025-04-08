package dto

type CreateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type GetJWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessTokenOutput struct {
	AccessToken string `json:"access_token"`
}
