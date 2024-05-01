package domain

type SignUpInput struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type ConfirmEmailInput struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,max=30"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=30"`
}

type UpdateProfileInput struct {
	Name       *string `json:"name"`
	ProfileImg *string `json:"profile_image"`
}

type UpdatePasswordInput struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}
